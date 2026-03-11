// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Package: memory/v2.0.7
// Module: Circuit Breaker (from @architect design)
// Sprint 1 Implementation - Phase 1

package memory

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// CircuitState represents the state of the circuit breaker
type CircuitState int

const (
	StateClosed CircuitState = iota
	StateOpen
	StateHalfOpen
)

// CircuitBreaker implements the circuit breaker pattern for fault tolerance
type CircuitBreaker struct {
	mu sync.RWMutex

	// Configuration
	maxFailures  int
	resetTimeout time.Duration
	halfOpenMax  int

	// State
	state         CircuitState
	failures      int
	lastFailTime  time.Time
	halfOpenCalls int

	// Metrics
	totalCalls    int64
	totalFailures int64
	totalSuccess  int64
}

// CircuitBreakerConfig holds circuit breaker configuration
type CircuitBreakerConfig struct {
	MaxFailures   int // default: 5
	ResetTimeoutS int // default: 30
	HalfOpenMax   int // default: 3
}

// NewCircuitBreaker creates a new circuit breaker instance
func NewCircuitBreaker(cfg CircuitBreakerConfig) *CircuitBreaker {
	if cfg.MaxFailures <= 0 {
		cfg.MaxFailures = 5
	}
	if cfg.ResetTimeoutS <= 0 {
		cfg.ResetTimeoutS = 30
	}
	if cfg.HalfOpenMax <= 0 {
		cfg.HalfOpenMax = 3
	}

	// Ensure HalfOpenMax is reasonable relative to MaxFailures
	if cfg.HalfOpenMax > cfg.MaxFailures {
		cfg.HalfOpenMax = cfg.MaxFailures
	}

	return &CircuitBreaker{
		maxFailures:  cfg.MaxFailures,
		resetTimeout: time.Duration(cfg.ResetTimeoutS) * time.Second,
		halfOpenMax:  cfg.HalfOpenMax,
		state:        StateClosed,
	}
}

// Call executes the given function with circuit breaker protection
func (cb *CircuitBreaker) Call(ctx context.Context, fn func() error) error {
	// Check context BEFORE execution to avoid wasting resources
	if err := ctx.Err(); err != nil {
		return err
	}

	// Check if circuit is open
	if !cb.allowRequest() {
		return fmt.Errorf("ERR_CIRCUIT_OPEN: circuit breaker is open")
	}

	// Execute function
	err := fn()

	// Don't record context errors as provider failures
	if err != nil && (errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)) {
		return err // Pass through without recording
	}

	cb.recordResult(err)
	return err
}

// allowRequest checks if the request should be allowed
func (cb *CircuitBreaker) allowRequest() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.totalCalls++

	switch cb.state {
	case StateClosed:
		return true

	case StateOpen:
		// Check if reset timeout has passed
		if time.Since(cb.lastFailTime) > cb.resetTimeout {
			cb.state = StateHalfOpen
			cb.halfOpenCalls = 1 // Count this first request
			return true
		}
		return false

	case StateHalfOpen:
		// Allow limited requests in half-open state
		if cb.halfOpenCalls < cb.halfOpenMax {
			cb.halfOpenCalls++
			return true
		}
		return false

	default:
		return false
	}
}

// recordResult records the result of a call
func (cb *CircuitBreaker) recordResult(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.totalFailures++
		cb.failures++
		cb.lastFailTime = time.Now()

		switch cb.state {
		case StateClosed:
			if cb.failures >= cb.maxFailures {
				cb.state = StateOpen
			}

		case StateHalfOpen:
			// Any failure in half-open state reopens the circuit
			cb.state = StateOpen
			cb.halfOpenCalls = 0
			// Reset failures to start fresh - ensures consistent behavior
			cb.failures = 1 // Count this failure
		}
	} else {
		cb.totalSuccess++

		switch cb.state {
		case StateClosed:
			// Reset failure count on success
			cb.failures = 0

		case StateHalfOpen:
			// If all half-open calls succeed, close the circuit
			if cb.halfOpenCalls >= cb.halfOpenMax {
				cb.state = StateClosed
				cb.failures = 0
				cb.halfOpenCalls = 0
			}
		}
	}
}

// State returns the current state of the circuit breaker
func (cb *CircuitBreaker) State() CircuitState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// IsOpen returns true if the circuit is open
func (cb *CircuitBreaker) IsOpen() bool {
	return cb.State() == StateOpen
}

// Metrics returns current circuit breaker metrics
func (cb *CircuitBreaker) Metrics() CircuitBreakerMetrics {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	return CircuitBreakerMetrics{
		State:         cb.state,
		TotalCalls:    cb.totalCalls,
		TotalFailures: cb.totalFailures,
		TotalSuccess:  cb.totalSuccess,
		Failures:      cb.failures,
	}
}

// CircuitBreakerMetrics holds circuit breaker metrics
type CircuitBreakerMetrics struct {
	State         CircuitState
	TotalCalls    int64
	TotalFailures int64
	TotalSuccess  int64
	Failures      int
}

// Reset resets the circuit breaker to closed state
func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.state = StateClosed
	cb.failures = 0
	cb.halfOpenCalls = 0
}
