package providers

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestCircuitBreaker_StateClosed(t *testing.T) {
	cb := NewCircuitBreaker(DefaultBreakerConfig()).(*circuitBreaker)
	
	if cb.State() != StateClosed {
		t.Errorf("expected initial state to be Closed, got %v", cb.State())
	}

	// Successful call should keep circuit closed
	err := cb.Call(context.Background(), func() error {
		return nil
	})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if cb.State() != StateClosed {
		t.Errorf("expected state to remain Closed after success, got %v", cb.State())
	}
}

func TestCircuitBreaker_TransitionToOpen(t *testing.T) {
	config := DefaultBreakerConfig()
	config.FailureThreshold = 3
	cb := NewCircuitBreaker(config).(*circuitBreaker)

	// Trigger failures to open circuit
	testErr := errors.New("test error")
	for i := 0; i < 3; i++ {
		_ = cb.Call(context.Background(), func() error {
			return testErr
		})
	}

	if cb.State() != StateOpen {
		t.Errorf("expected state to be Open after %d failures, got %v", config.FailureThreshold, cb.State())
	}

	// Next call should fail fast
	err := cb.Call(context.Background(), func() error {
		t.Error("function should not be called when circuit is open")
		return nil
	})
	if err == nil {
		t.Error("expected error when circuit is open")
	}
}

func TestCircuitBreaker_TransitionToHalfOpen(t *testing.T) {
	config := DefaultBreakerConfig()
	config.FailureThreshold = 2
	config.OpenTimeout = 100 * time.Millisecond
	cb := NewCircuitBreaker(config).(*circuitBreaker)

	// Open the circuit
	testErr := errors.New("test error")
	for i := 0; i < 2; i++ {
		_ = cb.Call(context.Background(), func() error {
			return testErr
		})
	}

	if cb.State() != StateOpen {
		t.Fatalf("expected state to be Open, got %v", cb.State())
	}

	// Wait for timeout
	time.Sleep(150 * time.Millisecond)

	// Next call should transition to half-open and then closed if successful
	_ = cb.Call(context.Background(), func() error {
		return nil
	})

	if cb.State() != StateClosed {
		t.Errorf("expected state to be Closed after successful half-open call, got %v", cb.State())
	}
}

func TestCircuitBreaker_HalfOpenToClosed(t *testing.T) {
	config := DefaultBreakerConfig()
	config.FailureThreshold = 2
	config.OpenTimeout = 50 * time.Millisecond
	config.HalfOpenMaxCalls = 2
	cb := NewCircuitBreaker(config).(*circuitBreaker)

	// Open the circuit
	testErr := errors.New("test error")
	for i := 0; i < 2; i++ {
		_ = cb.Call(context.Background(), func() error {
			return testErr
		})
	}

	// Wait for timeout to transition to half-open
	time.Sleep(100 * time.Millisecond)

	// Successful calls in half-open should close circuit
	for i := 0; i < 2; i++ {
		err := cb.Call(context.Background(), func() error {
			return nil
		})
		if err != nil {
			t.Errorf("expected no error in half-open, got %v", err)
		}
	}

	if cb.State() != StateClosed {
		t.Errorf("expected state to be Closed after successful half-open calls, got %v", cb.State())
	}
}

func TestCircuitBreaker_HalfOpenToOpen(t *testing.T) {
	config := DefaultBreakerConfig()
	config.FailureThreshold = 2
	config.OpenTimeout = 50 * time.Millisecond
	cb := NewCircuitBreaker(config).(*circuitBreaker)

	// Open the circuit
	testErr := errors.New("test error")
	for i := 0; i < 2; i++ {
		_ = cb.Call(context.Background(), func() error {
			return testErr
		})
	}

	// Wait for timeout to transition to half-open
	time.Sleep(100 * time.Millisecond)

	// Failed call in half-open should reopen circuit
	_ = cb.Call(context.Background(), func() error {
		return testErr
	})

	if cb.State() != StateOpen {
		t.Errorf("expected state to be Open after failed half-open call, got %v", cb.State())
	}
}

func TestCircuitBreaker_FailureRate(t *testing.T) {
	config := DefaultBreakerConfig()
	config.FailureThreshold = 100 // High threshold to test failure rate
	config.FailureRate = 0.5
	config.SamplingWindow = 1 * time.Second
	cb := NewCircuitBreaker(config).(*circuitBreaker)

	testErr := errors.New("test error")

	// 5 successes, 5 failures = 50% failure rate
	// We need 10 calls because of the minimum sample requirement (5)
	for i := 0; i < 10; i++ {
		var err error
		if i%2 == 0 {
			err = nil
		} else {
			err = testErr
		}
		_ = cb.Call(context.Background(), func() error {
			return err
		})
	}

	if cb.State() != StateOpen {
		t.Errorf("expected state to be Open due to failure rate, got %v", cb.State())
	}
}

func TestCircuitBreaker_Metrics(t *testing.T) {
	cb := NewCircuitBreaker(DefaultBreakerConfig())

	// Make some calls
	_ = cb.Call(context.Background(), func() error { return nil })
	_ = cb.Call(context.Background(), func() error { return errors.New("error") })
	_ = cb.Call(context.Background(), func() error { return nil })

	metrics := cb.Metrics()

	if metrics.TotalCalls != 3 {
		t.Errorf("expected 3 total calls, got %d", metrics.TotalCalls)
	}
	if metrics.SuccessfulCalls != 2 {
		t.Errorf("expected 2 successful calls, got %d", metrics.SuccessfulCalls)
	}
	if metrics.FailedCalls != 1 {
		t.Errorf("expected 1 failed call, got %d", metrics.FailedCalls)
	}
}

func TestCircuitBreaker_Reset(t *testing.T) {
	config := DefaultBreakerConfig()
	config.FailureThreshold = 2
	cb := NewCircuitBreaker(config)

	// Open the circuit
	testErr := errors.New("test error")
	for i := 0; i < 2; i++ {
		_ = cb.Call(context.Background(), func() error {
			return testErr
		})
	}

	if cb.State() != StateOpen {
		t.Fatalf("expected state to be Open, got %v", cb.State())
	}

	// Reset should close the circuit
	cb.Reset()

	if cb.State() != StateClosed {
		t.Errorf("expected state to be Closed after reset, got %v", cb.State())
	}

	// Should be able to make calls again
	err := cb.Call(context.Background(), func() error {
		return nil
	})
	if err != nil {
		t.Errorf("expected no error after reset, got %v", err)
	}
}

func TestCircuitBreaker_HalfOpenMaxCalls(t *testing.T) {
	config := DefaultBreakerConfig()
	config.FailureThreshold = 2
	config.OpenTimeout = 50 * time.Millisecond
	config.HalfOpenMaxCalls = 1
	cb := NewCircuitBreaker(config).(*circuitBreaker)

	// Open the circuit
	testErr := errors.New("test error")
	for i := 0; i < 2; i++ {
		_ = cb.Call(context.Background(), func() error {
			return testErr
		})
	}

	// Wait for timeout
	time.Sleep(100 * time.Millisecond)

	// First call should be allowed
	done := make(chan bool)
	go func() {
		_ = cb.Call(context.Background(), func() error {
			time.Sleep(50 * time.Millisecond) // Simulate slow call
			return nil
		})
		done <- true
	}()

	time.Sleep(10 * time.Millisecond) // Let first call start

	// Second concurrent call should be rejected
	err := cb.Call(context.Background(), func() error {
		t.Error("second call should not execute")
		return nil
	})
	if err == nil {
		t.Error("expected error for exceeding max half-open calls")
	}

	<-done // Wait for first call to complete
}

func TestCircuitBreaker_SamplingWindow(t *testing.T) {
	config := DefaultBreakerConfig()
	config.FailureThreshold = 10 // Higher than calls
	config.FailureRate = 1.0     // Disable failure rate for this test
	config.SamplingWindow = 100 * time.Millisecond
	cb := NewCircuitBreaker(config).(*circuitBreaker)

	testErr := errors.New("test error")

	// Add failures
	for i := 0; i < 5; i++ {
		_ = cb.Call(context.Background(), func() error {
			return testErr
		})
	}

	if cb.State() != StateClosed {
		t.Fatalf("expected state to be Closed, got %v", cb.State())
	}

	// Wait for sampling window to expire
	time.Sleep(150 * time.Millisecond)

	// Old failures should be cleared
	_ = cb.Call(context.Background(), func() error {
		return nil
	})

	metrics := cb.Metrics()
	if metrics.FailureRate > 0 {
		t.Errorf("expected failure rate to be 0 after window expire, got %f", metrics.FailureRate)
	}
}

func TestCircuitBreaker_NonSystemError(t *testing.T) {
	config := DefaultBreakerConfig()
	config.FailureThreshold = 2
	cb := NewCircuitBreaker(config)

	// Trigger non-system errors
	for i := 0; i < 5; i++ {
		_ = cb.Call(context.Background(), func() error {
			return &nonSystemError{err: errors.New("auth error")}
		})
	}

	// Circuit should still be Closed
	if cb.State() != StateClosed {
		t.Errorf("expected state to remain Closed after non-system errors, got %v", cb.State())
	}
}
