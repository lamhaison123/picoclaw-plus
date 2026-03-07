package retry

import "time"

// RetryPolicy defines the retry behavior configuration
type RetryPolicy struct {
	MaxAttempts       int           // Maximum number of retry attempts
	InitialBackoff    time.Duration // Initial backoff duration
	MaxBackoff        time.Duration // Maximum backoff duration
	BackoffMultiplier float64       // Multiplier for exponential backoff
	JitterFactor      float64       // Jitter factor (0-1) to add randomness
}

// DefaultLLMRetryPolicy provides sensible defaults for LLM provider retries
var DefaultLLMRetryPolicy = RetryPolicy{
	MaxAttempts:       3,
	InitialBackoff:    200 * time.Millisecond,
	MaxBackoff:        5 * time.Second,
	BackoffMultiplier: 2.0,
	JitterFactor:      0.2,
}
