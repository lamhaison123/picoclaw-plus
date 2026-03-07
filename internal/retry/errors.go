package retry

import (
	"context"
	"errors"
	"strings"
)

// IsRetryable determines if an error should trigger a retry
func IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	// Context cancelled should not retry
	if errors.Is(err, context.Canceled) {
		return false
	}

	// Timeout can retry with new context
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}

	// Check for circuit breaker errors
	errStr := err.Error()
	if strings.Contains(errStr, "circuit breaker") || 
	   strings.Contains(errStr, "cooldown") ||
	   strings.Contains(errStr, "skipped") {
		return true
	}

	// Network/timeout errors are retryable
	if strings.Contains(errStr, "timeout") ||
	   strings.Contains(errStr, "connection") ||
	   strings.Contains(errStr, "network") {
		return true
	}

	// Rate limit errors are retryable
	if strings.Contains(errStr, "rate limit") ||
	   strings.Contains(errStr, "429") {
		return true
	}

	// Server errors (5xx) are retryable
	if strings.Contains(errStr, "500") ||
	   strings.Contains(errStr, "502") ||
	   strings.Contains(errStr, "503") ||
	   strings.Contains(errStr, "504") {
		return true
	}

	// Client errors (4xx except 429) are NOT retryable
	if strings.Contains(errStr, "400") ||
	   strings.Contains(errStr, "401") ||
	   strings.Contains(errStr, "403") ||
	   strings.Contains(errStr, "404") {
		return false
	}

	// Default: retry for unknown errors
	return true
}
