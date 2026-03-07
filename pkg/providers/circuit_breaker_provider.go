package providers

import (
	"context"
)

// CircuitBreakerProvider wraps an LLMProvider with circuit breaker protection.
type CircuitBreakerProvider struct {
	delegate LLMProvider
	breaker  CircuitBreaker
	provider string
}

// NewCircuitBreakerProvider creates a new provider with circuit breaker.
func NewCircuitBreakerProvider(delegate LLMProvider, breaker CircuitBreaker, providerName string) *CircuitBreakerProvider {
	return &CircuitBreakerProvider{
		delegate: delegate,
		breaker:  breaker,
		provider: providerName,
	}
}

// Chat executes the chat request with circuit breaker protection.
func (p *CircuitBreakerProvider) Chat(
	ctx context.Context,
	messages []Message,
	tools []ToolDefinition,
	model string,
	options map[string]any,
) (*LLMResponse, error) {
	var resp *LLMResponse
	err := p.breaker.Call(ctx, func() error {
		var chatErr error
		resp, chatErr = p.delegate.Chat(ctx, messages, tools, model, options)
		
		if chatErr != nil {
			// Classify error to decide if it should count towards circuit breaker failure
			failErr := ClassifyError(chatErr, p.provider, model)
			if failErr != nil {
				// Only system-level failures (timeout, overloaded, rate limit) should trip the breaker
				// Auth or Format errors are client-side and shouldn't trip the breaker for everyone
				switch failErr.Reason {
				case FailoverTimeout, FailoverOverloaded, FailoverRateLimit, FailoverUnknown:
					return chatErr
				default:
					// For other errors, we return nil to the breaker so it doesn't count as a failure,
					// but we still need to return the actual error to the caller.
					// This is a bit tricky with the current CircuitBreaker.Call interface.
					// We'll wrap it in a special error that the breaker ignores but we catch.
					return &nonSystemError{err: chatErr}
				}
			}
			return chatErr
		}
		return nil
	})

	if err != nil {
		if nse, ok := err.(*nonSystemError); ok {
			return resp, nse.err
		}
		return nil, err
	}

	return resp, nil
}

// Unwrap returns the underlying LLMProvider.
func (p *CircuitBreakerProvider) Unwrap() LLMProvider {
	return p.delegate
}

// GetDefaultModel returns the delegate's default model.
func (p *CircuitBreakerProvider) GetDefaultModel() string {
	return p.delegate.GetDefaultModel()
}

// Close closes the delegate if it's a StatefulProvider.
func (p *CircuitBreakerProvider) Close() {
	if sp, ok := p.delegate.(StatefulProvider); ok {
		sp.Close()
	}
}

// nonSystemError is used to pass through errors that shouldn't trip the circuit breaker.
type nonSystemError struct {
	err error
}

func (e *nonSystemError) Error() string {
	return e.err.Error()
}
