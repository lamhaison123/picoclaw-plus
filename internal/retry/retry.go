package retry

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"
)

// CalculateBackoff computes the backoff duration for a given attempt with jitter
func (p *RetryPolicy) CalculateBackoff(attempt int) time.Duration {
	// Exponential backoff: initial * (multiplier ^ attempt)
	backoff := float64(p.InitialBackoff) * math.Pow(p.BackoffMultiplier, float64(attempt))

	// Cap at MaxBackoff
	if backoff > float64(p.MaxBackoff) {
		backoff = float64(p.MaxBackoff)
	}

	// Add jitter to avoid thundering herd
	jitter := backoff * p.JitterFactor * (rand.Float64()*2 - 1) // [-jitter, +jitter]

	result := time.Duration(backoff + jitter)
	if result < 0 {
		result = p.InitialBackoff
	}

	return result
}

// WithRetry executes a function with retry logic using exponential backoff
func WithRetry(ctx context.Context, policy RetryPolicy, fn func() error) error {
	var lastErr error

	for attempt := 0; attempt < policy.MaxAttempts; attempt++ {
		// Try to execute the function
		err := fn()

		if err == nil {
			return nil // Success
		}

		lastErr = err

		// Check if error is retryable
		if !IsRetryable(err) {
			return fmt.Errorf("non-retryable error: %w", err)
		}

		// If this is the last attempt, don't sleep
		if attempt == policy.MaxAttempts-1 {
			break
		}

		// Calculate backoff and sleep
		backoff := policy.CalculateBackoff(attempt)

		select {
		case <-ctx.Done():
			return fmt.Errorf("retry cancelled: %w", ctx.Err())
		case <-time.After(backoff):
			// Continue to next attempt
		}
	}

	return fmt.Errorf("max retry attempts (%d) reached: %w", policy.MaxAttempts, lastErr)
}
