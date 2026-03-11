// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors
//
// Package: memory/v2.0.7
// Module: Circuit Breaker Unit Tests
// Sprint 1 Implementation - Phase 2

package memory

import (
	"context"
	"errors"
	"testing"
	"time"
)

// TestCircuitBreaker_StateTransitions tests circuit breaker state machine
func TestCircuitBreaker_StateTransitions(t *testing.T) {
	breaker := NewCircuitBreaker(CircuitBreakerConfig{
		MaxFailures:   3,
		ResetTimeoutS: 1,
		HalfOpenMax:   2,
	})

	// Initial state should be Closed
	if breaker.State() != StateClosed {
		t.Errorf("Initial state = %v, want StateClosed", breaker.State())
	}

	ctx := context.Background()

	// Simulate 3 failures to open circuit
	for i := 0; i < 3; i++ {
		err := breaker.Call(ctx, func() error {
			return errors.New("simulated failure")
		})
		if err == nil {
			t.Errorf("Call() expected error, got nil")
		}
	}

	// Circuit should now be Open
	if breaker.State() != StateOpen {
		t.Errorf("After 3 failures, state = %v, want StateOpen", breaker.State())
	}

	// Next call should fail immediately with ERR_CIRCUIT_OPEN
	err := breaker.Call(ctx, func() error {
		t.Error("Function should not be called when circuit is open")
		return nil
	})
	if err == nil {
		t.Error("Call() expected ERR_CIRCUIT_OPEN, got nil")
	}
	if !contains(err.Error(), "ERR_CIRCUIT_OPEN") {
		t.Errorf("Call() error = %v, want ERR_CIRCUIT_OPEN", err)
	}

	// Wait for reset timeout
	time.Sleep(1100 * time.Millisecond)

	// Circuit should transition to HalfOpen
	err = breaker.Call(ctx, func() error {
		return nil // Success
	})
	if err != nil {
		t.Errorf("Call() in half-open state unexpected error = %v", err)
	}

	if breaker.State() != StateHalfOpen {
		t.Errorf("After reset timeout, state = %v, want StateHalfOpen", breaker.State())
	}

	// One more success should close the circuit
	err = breaker.Call(ctx, func() error {
		return nil // Success
	})
	if err != nil {
		t.Errorf("Call() in half-open state unexpected error = %v", err)
	}

	if breaker.State() != StateClosed {
		t.Errorf("After 2 successes in half-open, state = %v, want StateClosed", breaker.State())
	}
}

// TestCircuitBreaker_HalfOpenFailure tests failure in half-open state
func TestCircuitBreaker_HalfOpenFailure(t *testing.T) {
	breaker := NewCircuitBreaker(CircuitBreakerConfig{
		MaxFailures:   2,
		ResetTimeoutS: 1,
		HalfOpenMax:   2,
	})

	ctx := context.Background()

	// Open the circuit with 2 failures
	for i := 0; i < 2; i++ {
		breaker.Call(ctx, func() error {
			return errors.New("failure")
		})
	}

	if breaker.State() != StateOpen {
		t.Errorf("state = %v, want StateOpen", breaker.State())
	}

	// Wait for reset timeout
	time.Sleep(1100 * time.Millisecond)

	// First call in half-open succeeds
	err := breaker.Call(ctx, func() error {
		return nil
	})
	if err != nil {
		t.Errorf("Call() unexpected error = %v", err)
	}

	if breaker.State() != StateHalfOpen {
		t.Errorf("state = %v, want StateHalfOpen", breaker.State())
	}

	// Second call fails - should reopen circuit
	err = breaker.Call(ctx, func() error {
		return errors.New("failure in half-open")
	})
	if err == nil {
		t.Error("Call() expected error, got nil")
	}

	if breaker.State() != StateOpen {
		t.Errorf("After failure in half-open, state = %v, want StateOpen", breaker.State())
	}
}

// TestCircuitBreaker_Metrics tests metrics tracking
func TestCircuitBreaker_Metrics(t *testing.T) {
	breaker := NewCircuitBreaker(CircuitBreakerConfig{
		MaxFailures:   5,
		ResetTimeoutS: 30,
		HalfOpenMax:   3,
	})

	ctx := context.Background()

	// Execute some successful calls
	for i := 0; i < 3; i++ {
		breaker.Call(ctx, func() error {
			return nil
		})
	}

	// Execute some failed calls
	for i := 0; i < 2; i++ {
		breaker.Call(ctx, func() error {
			return errors.New("failure")
		})
	}

	metrics := breaker.Metrics()

	if metrics.TotalCalls != 5 {
		t.Errorf("TotalCalls = %d, want 5", metrics.TotalCalls)
	}
	if metrics.TotalSuccess != 3 {
		t.Errorf("TotalSuccess = %d, want 3", metrics.TotalSuccess)
	}
	if metrics.TotalFailures != 2 {
		t.Errorf("TotalFailures = %d, want 2", metrics.TotalFailures)
	}
	if metrics.Failures != 2 {
		t.Errorf("Failures = %d, want 2", metrics.Failures)
	}
	if metrics.State != StateClosed {
		t.Errorf("State = %v, want StateClosed", metrics.State)
	}
}

// TestCircuitBreaker_Reset tests manual reset
func TestCircuitBreaker_Reset(t *testing.T) {
	breaker := NewCircuitBreaker(CircuitBreakerConfig{
		MaxFailures:   2,
		ResetTimeoutS: 30,
		HalfOpenMax:   3,
	})

	ctx := context.Background()

	// Open the circuit
	for i := 0; i < 2; i++ {
		breaker.Call(ctx, func() error {
			return errors.New("failure")
		})
	}

	if breaker.State() != StateOpen {
		t.Errorf("state = %v, want StateOpen", breaker.State())
	}

	// Manual reset
	breaker.Reset()

	if breaker.State() != StateClosed {
		t.Errorf("After Reset(), state = %v, want StateClosed", breaker.State())
	}

	metrics := breaker.Metrics()
	if metrics.Failures != 0 {
		t.Errorf("After Reset(), Failures = %d, want 0", metrics.Failures)
	}
}

// TestCircuitBreaker_IsOpen tests IsOpen helper
func TestCircuitBreaker_IsOpen(t *testing.T) {
	breaker := NewCircuitBreaker(CircuitBreakerConfig{
		MaxFailures:   2,
		ResetTimeoutS: 30,
		HalfOpenMax:   3,
	})

	ctx := context.Background()

	// Initially closed
	if breaker.IsOpen() {
		t.Error("IsOpen() = true, want false for initial state")
	}

	// Open the circuit
	for i := 0; i < 2; i++ {
		breaker.Call(ctx, func() error {
			return errors.New("failure")
		})
	}

	// Should be open now
	if !breaker.IsOpen() {
		t.Error("IsOpen() = false, want true after failures")
	}
}

// TestCircuitBreaker_ConcurrentCalls tests thread safety
func TestCircuitBreaker_ConcurrentCalls(t *testing.T) {
	breaker := NewCircuitBreaker(CircuitBreakerConfig{
		MaxFailures:   10,
		ResetTimeoutS: 30,
		HalfOpenMax:   3,
	})

	ctx := context.Background()
	done := make(chan bool)

	// Launch 10 concurrent goroutines
	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 10; j++ {
				breaker.Call(ctx, func() error {
					if j%2 == 0 {
						return nil // Success
					}
					return errors.New("failure") // Failure
				})
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	metrics := breaker.Metrics()

	// Should have 100 total calls (10 goroutines * 10 calls)
	if metrics.TotalCalls != 100 {
		t.Errorf("TotalCalls = %d, want 100", metrics.TotalCalls)
	}

	// Should have 50 successes and 50 failures
	if metrics.TotalSuccess != 50 {
		t.Errorf("TotalSuccess = %d, want 50", metrics.TotalSuccess)
	}
	if metrics.TotalFailures != 50 {
		t.Errorf("TotalFailures = %d, want 50", metrics.TotalFailures)
	}
}

// TestCircuitBreaker_DefaultConfig tests default configuration values
func TestCircuitBreaker_DefaultConfig(t *testing.T) {
	tests := []struct {
		name   string
		cfg    CircuitBreakerConfig
		wantMF int
		wantRT time.Duration
		wantHO int
	}{
		{
			name:   "all defaults",
			cfg:    CircuitBreakerConfig{},
			wantMF: 5,
			wantRT: 30 * time.Second,
			wantHO: 3,
		},
		{
			name: "partial defaults",
			cfg: CircuitBreakerConfig{
				MaxFailures: 10,
			},
			wantMF: 10,
			wantRT: 30 * time.Second,
			wantHO: 3,
		},
		{
			name: "no defaults",
			cfg: CircuitBreakerConfig{
				MaxFailures:   7,
				ResetTimeoutS: 60,
				HalfOpenMax:   5,
			},
			wantMF: 7,
			wantRT: 60 * time.Second,
			wantHO: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			breaker := NewCircuitBreaker(tt.cfg)

			if breaker.maxFailures != tt.wantMF {
				t.Errorf("maxFailures = %d, want %d", breaker.maxFailures, tt.wantMF)
			}
			if breaker.resetTimeout != tt.wantRT {
				t.Errorf("resetTimeout = %v, want %v", breaker.resetTimeout, tt.wantRT)
			}
			if breaker.halfOpenMax != tt.wantHO {
				t.Errorf("halfOpenMax = %d, want %d", breaker.halfOpenMax, tt.wantHO)
			}
		})
	}
}
