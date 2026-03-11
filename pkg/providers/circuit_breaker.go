package providers

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// State represents the circuit breaker state.
type State int

const (
	StateClosed   State = iota // Normal operation
	StateOpen                  // Failing fast
	StateHalfOpen              // Testing recovery
)

func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half_open"
	default:
		return "unknown"
	}
}

// BreakerConfig holds circuit breaker configuration.
type BreakerConfig struct {
	FailureThreshold int           // Consecutive failures to open circuit
	FailureRate      float64       // Failure rate (0.0-1.0) to open circuit
	OpenTimeout      time.Duration // Time to wait before half-open
	HalfOpenMaxCalls int           // Max concurrent calls in half-open
	SamplingWindow   time.Duration // Window for failure rate calculation
}

// DefaultBreakerConfig returns sensible defaults.
func DefaultBreakerConfig() BreakerConfig {
	return BreakerConfig{
		FailureThreshold: 5,
		FailureRate:      0.5,
		OpenTimeout:      30 * time.Second,
		HalfOpenMaxCalls: 2,
		SamplingWindow:   10 * time.Second,
	}
}

// BreakerMetrics exposes circuit breaker metrics.
type BreakerMetrics struct {
	State            State
	FailureRate      float64
	ConsecutiveFails int
	LastStateChange  time.Time
	TotalCalls       int64
	SuccessfulCalls  int64
	FailedCalls      int64
}

// CircuitBreaker wraps function calls with circuit breaker pattern.
type CircuitBreaker interface {
	Call(ctx context.Context, fn func() error) error
	State() State
	Metrics() BreakerMetrics
	Reset()
}

// circuitBreaker implements CircuitBreaker interface.
type circuitBreaker struct {
	mu     sync.RWMutex
	config BreakerConfig

	state            State
	consecutiveFails int
	lastStateChange  time.Time
	halfOpenCalls    int

	// Metrics
	totalCalls      int64
	successfulCalls int64
	failedCalls     int64

	// Sampling window for failure rate
	recentCalls []callResult
	windowStart time.Time

	nowFunc func() time.Time // for testing
}

type callResult struct {
	timestamp time.Time
	success   bool
}

// NewCircuitBreaker creates a new circuit breaker with given config.
func NewCircuitBreaker(config BreakerConfig) CircuitBreaker {
	return &circuitBreaker{
		config:          config,
		state:           StateClosed,
		lastStateChange: time.Now(),
		recentCalls:     make([]callResult, 0),
		windowStart:     time.Now(),
		nowFunc:         time.Now,
	}
}

// Call executes the function with circuit breaker protection.
func (cb *circuitBreaker) Call(ctx context.Context, fn func() error) error {
	// Check if we can proceed
	if err := cb.beforeCall(); err != nil {
		return err
	}

	// Execute the function
	err := fn()

	// Record result
	cb.afterCall(err)

	return err
}

// beforeCall checks if the call should proceed based on current state.
func (cb *circuitBreaker) beforeCall() error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	now := cb.nowFunc()

	switch cb.state {
	case StateClosed:
		// Normal operation, allow call
		return nil

	case StateOpen:
		// Check if we should transition to half-open
		if now.Sub(cb.lastStateChange) >= cb.config.OpenTimeout {
			cb.transitionTo(StateHalfOpen, now)
			cb.halfOpenCalls++ // Increment for the probe call
			return nil
		}
		// Still open, fail fast
		return fmt.Errorf("circuit breaker is open")

	case StateHalfOpen:
		// Limit concurrent calls in half-open state
		if cb.halfOpenCalls >= cb.config.HalfOpenMaxCalls {
			return fmt.Errorf("circuit breaker half-open: max concurrent calls reached")
		}
		cb.halfOpenCalls++
		return nil

	default:
		return fmt.Errorf("unknown circuit breaker state")
	}
}

// afterCall records the result and updates state accordingly.
func (cb *circuitBreaker) afterCall(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	now := cb.nowFunc()

	// Skip recording non-system errors (e.g. auth, format) as failures
	// These errors are client-side and shouldn't trip the circuit breaker.
	if err != nil {
		var nse *nonSystemError
		if errors.As(err, &nse) {
			// Treat as success for the circuit breaker's health tracking
			err = nil
		}
	}

	success := err == nil

	// Update metrics
	cb.totalCalls++
	if success {
		cb.successfulCalls++
	} else {
		cb.failedCalls++
	}

	// Add to recent calls for failure rate calculation
	cb.addRecentCall(now, success)

	// Update state based on result
	switch cb.state {
	case StateClosed:
		if success {
			cb.consecutiveFails = 0
		} else {
			cb.consecutiveFails++
			// Check if we should open the circuit
			if cb.shouldOpen() {
				cb.transitionTo(StateOpen, now)
			}
		}

	case StateHalfOpen:
		cb.halfOpenCalls--
		if success {
			// Success in half-open, check if we can close
			cb.consecutiveFails = 0
			if cb.halfOpenCalls == 0 {
				// All probe calls succeeded, close the circuit
				cb.transitionTo(StateClosed, now)
			}
		} else {
			// Failure in half-open, reopen the circuit
			cb.consecutiveFails++
			cb.transitionTo(StateOpen, now)
		}

	case StateOpen:
		// Should not reach here, but handle gracefully
		if !success {
			cb.consecutiveFails++
		}
	}
}

// shouldOpen determines if the circuit should open based on thresholds.
func (cb *circuitBreaker) shouldOpen() bool {
	// Check consecutive failures threshold
	if cb.consecutiveFails >= cb.config.FailureThreshold {
		return true
	}

	// Check failure rate threshold - only if we have enough samples
	if len(cb.recentCalls) < 10 {
		return false
	}

	failureRate := cb.calculateFailureRate()
	if failureRate >= cb.config.FailureRate {
		return true
	}

	return false
}

// calculateFailureRate computes failure rate in the sampling window.
func (cb *circuitBreaker) calculateFailureRate() float64 {
	if len(cb.recentCalls) == 0 {
		return 0.0
	}

	failures := 0
	for _, call := range cb.recentCalls {
		if !call.success {
			failures++
		}
	}

	return float64(failures) / float64(len(cb.recentCalls))
}

// addRecentCall adds a call result to the sampling window.
func (cb *circuitBreaker) addRecentCall(now time.Time, success bool) {
	// Remove old calls outside the sampling window
	cutoff := now.Add(-cb.config.SamplingWindow)
	validCalls := make([]callResult, 0, len(cb.recentCalls)+1)
	for _, call := range cb.recentCalls {
		if call.timestamp.After(cutoff) {
			validCalls = append(validCalls, call)
		}
	}

	// Add new call
	validCalls = append(validCalls, callResult{
		timestamp: now,
		success:   success,
	})

	cb.recentCalls = validCalls
	if len(cb.recentCalls) > 0 && cb.recentCalls[0].timestamp.After(cutoff) {
		cb.windowStart = cb.recentCalls[0].timestamp
	} else {
		cb.windowStart = now
	}
}

// transitionTo changes the circuit breaker state.
func (cb *circuitBreaker) transitionTo(newState State, now time.Time) {
	cb.state = newState
	cb.lastStateChange = now

	switch newState {
	case StateClosed:
		// Reset counters when closing
		cb.consecutiveFails = 0
		cb.recentCalls = make([]callResult, 0)
	case StateHalfOpen:
		cb.halfOpenCalls = 0
	}
}

// State returns the current circuit breaker state.
func (cb *circuitBreaker) State() State {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// Metrics returns current circuit breaker metrics.
func (cb *circuitBreaker) Metrics() BreakerMetrics {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	return BreakerMetrics{
		State:            cb.state,
		FailureRate:      cb.calculateFailureRate(),
		ConsecutiveFails: cb.consecutiveFails,
		LastStateChange:  cb.lastStateChange,
		TotalCalls:       cb.totalCalls,
		SuccessfulCalls:  cb.successfulCalls,
		FailedCalls:      cb.failedCalls,
	}
}

// Reset manually resets the circuit breaker to closed state.
func (cb *circuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	now := cb.nowFunc()
	cb.state = StateClosed
	cb.consecutiveFails = 0
	cb.halfOpenCalls = 0
	cb.lastStateChange = now
	cb.recentCalls = make([]callResult, 0)
	cb.windowStart = now
}
