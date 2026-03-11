# Circuit Breaker for LLM Providers

## Overview

PicoClaw implements a **Circuit Breaker pattern** to protect against cascading failures when LLM providers experience issues. This ensures system resilience and prevents resource exhaustion from repeated failed API calls.

## Key Features

- 🔒 **Per-Provider Isolation**: Each provider has its own circuit breaker - failures in one don't affect others
- 🔄 **Automatic Recovery**: Self-healing with half-open state for gradual recovery testing
- 🎯 **Smart Error Classification**: Only system errors trip the breaker (auth/format errors bypass)
- 📊 **Dual Triggers**: Failure threshold AND failure rate detection
- 🧵 **Thread-Safe**: Concurrent request handling with mutex protection
- 📈 **Observable**: Built-in metrics for monitoring and alerting

## How It Works

### State Machine

```
┌─────────┐  Failures exceed    ┌──────┐  Timeout    ┌───────────┐
│ CLOSED  │────threshold────────▶│ OPEN │────expires──▶│ HALF-OPEN │
└─────────┘                      └──────┘              └───────────┘
     ▲                              │                        │
     │                              │                        │
     └──────Success probes──────────┴────Failure────────────┘
```

**States:**

1. **CLOSED** (Normal Operation)
   - All requests pass through to provider
   - Tracks failure count and rate
   - Transitions to OPEN when thresholds exceeded

2. **OPEN** (Failure Mode)
   - All requests fail immediately with `ErrCircuitOpen`
   - No calls to provider (prevents resource waste)
   - After timeout, transitions to HALF-OPEN

3. **HALF-OPEN** (Recovery Testing)
   - Limited concurrent probe calls (default: 2)
   - Success → back to CLOSED
   - Failure → back to OPEN

### Trigger Conditions

Circuit opens when **EITHER** condition is met:

1. **Failure Threshold**: 5 consecutive failures
2. **Failure Rate**: >50% failures within 10s sampling window

### Error Classification

Only **system errors** trip the circuit breaker:

- ✅ **Trips Breaker**: Network errors, timeouts, 5xx responses, rate limits
- ❌ **Bypasses Breaker**: Auth errors (401), format errors (400), content filtering

This prevents auth issues from unnecessarily opening the circuit.

## Configuration

### Default Settings

```go
DefaultConfig = CircuitBreakerConfig{
    FailureThreshold:   5,      // Consecutive failures
    FailureRate:        0.5,    // 50% failure rate
    OpenTimeout:        30 * time.Second,
    HalfOpenMaxCalls:   2,      // Concurrent probes
    SamplingWindow:     10 * time.Second,
}
```

### Custom Configuration

```json
{
  "providers": {
    "openai": {
      "circuit_breaker": {
        "enabled": true,
        "failure_threshold": 5,
        "failure_rate": 0.5,
        "open_timeout": "30s",
        "half_open_max_calls": 2,
        "sampling_window": "10s"
      }
    }
  }
}
```

### Per-Provider Override

```go
// In code
config := llm.CircuitBreakerConfig{
    FailureThreshold: 3,  // More sensitive
    OpenTimeout:      60 * time.Second,  // Longer recovery
}

provider := factory.CreateWithCircuitBreaker("openai", config)
```

## Usage

### Automatic (Recommended)

Circuit breaker is automatically enabled for all LLM providers:

```go
// Factory automatically wraps providers
provider, err := factory.CreateProvider("openai")
// Provider is already wrapped with circuit breaker
```

### Manual Wrapping

```go
import "github.com/sipeed/picoclaw/pkg/llm"

// Create base provider
baseProvider := openai.NewProvider(config)

// Wrap with circuit breaker
provider := llm.WrapWithCircuitBreaker(
    baseProvider,
    "openai",
    llm.DefaultCircuitBreakerConfig,
)
```

### Checking State

```go
// Get circuit breaker state
state := provider.GetState()

switch state {
case llm.StateClosed:
    // Normal operation
case llm.StateOpen:
    // Provider unavailable
case llm.StateHalfOpen:
    // Testing recovery
}
```

### Manual Control

```go
// Force reset (use with caution)
provider.Reset()

// Get metrics
metrics := provider.GetMetrics()
fmt.Printf("Failures: %d, Success: %d\n", 
    metrics.Failures, metrics.Successes)
```

## Monitoring

### Metrics

Circuit breaker exposes these metrics:

```go
type Metrics struct {
    State           State
    Failures        int64
    Successes       int64
    ConsecutiveFails int64
    LastStateChange time.Time
}
```

### Logging

Circuit breaker logs state transitions:

```
[INFO] Circuit breaker for provider 'openai' opened after 5 consecutive failures
[INFO] Circuit breaker for provider 'openai' entering half-open state
[INFO] Circuit breaker for provider 'openai' closed after successful recovery
```

### Prometheus Integration (Future)

```go
// Planned metrics
circuit_breaker_state{provider="openai"} 0  // 0=closed, 1=open, 2=half-open
circuit_breaker_failures_total{provider="openai"} 42
circuit_breaker_state_changes_total{provider="openai",from="closed",to="open"} 3
```

## Best Practices

### 1. Provider Isolation

Each provider has independent circuit breaker:

```go
// OpenAI fails → circuit opens for OpenAI only
// Anthropic, Gemini continue working normally
```

### 2. Fallback Strategy

Combine with provider fallback:

```json
{
  "agents": {
    "defaults": {
      "model": "gpt-4",
      "fallback_models": ["claude-3-opus", "gemini-pro"]
    }
  }
}
```

When circuit opens:
1. Primary provider fails immediately
2. System tries fallback providers
3. User gets response from healthy provider

### 3. Tuning for Your Use Case

**High-Volume Production:**
```go
config := CircuitBreakerConfig{
    FailureThreshold: 10,  // More tolerance
    FailureRate:      0.7, // 70% failure rate
    OpenTimeout:      60 * time.Second,
}
```

**Development/Testing:**
```go
config := CircuitBreakerConfig{
    FailureThreshold: 3,   // Fail fast
    OpenTimeout:      10 * time.Second,
}
```

### 4. Monitoring Alerts

Set up alerts for:
- Circuit opens (immediate attention)
- High failure rate (warning)
- Frequent state changes (flapping)

## Troubleshooting

### Circuit Opens Frequently

**Symptoms:** Circuit breaker opens and closes repeatedly

**Causes:**
- Provider having intermittent issues
- Network instability
- Rate limits being hit

**Solutions:**
1. Increase `FailureThreshold` for more tolerance
2. Increase `OpenTimeout` for longer recovery periods
3. Check provider status page
4. Verify rate limits aren't exceeded

### Circuit Stays Open

**Symptoms:** Circuit remains open for extended periods

**Causes:**
- Provider is down
- Network connectivity issues
- Invalid API credentials

**Solutions:**
1. Check provider status
2. Verify network connectivity
3. Validate API keys
4. Check logs for specific error messages
5. Manual reset if needed: `provider.Reset()`

### Auth Errors Not Bypassing

**Symptoms:** Circuit opens due to auth errors

**Issue:** Error classifier not working correctly

**Solution:**
```go
// Verify error classification
if llm.IsSystemError(err) {
    // Should be false for auth errors
}
```

## Testing

### Unit Tests

```bash
cd pkg/llm
go test -v -run TestCircuitBreaker
```

### Integration Tests

```bash
# Test with real provider
go test -v -run TestCircuitBreakerIntegration
```

### Fault Injection

```go
// Simulate provider failures
for i := 0; i < 5; i++ {
    _, err := provider.Complete(ctx, request)
    // Circuit should open after 5 failures
}

state := provider.GetState()
assert.Equal(t, llm.StateOpen, state)
```

## Implementation Details

### Architecture

```
┌─────────────────────────────────────┐
│   LLM Provider Factory              │
│                                     │
│  ┌──────────────────────────────┐  │
│  │ Circuit Breaker Registry     │  │
│  │ (per-provider instances)     │  │
│  └──────────────────────────────┘  │
│           │                         │
│           ▼                         │
│  ┌──────────────────────────────┐  │
│  │ Circuit Breaker Wrapper      │  │
│  │ - State machine              │  │
│  │ - Error classification       │  │
│  │ - Metrics tracking           │  │
│  └──────────────────────────────┘  │
│           │                         │
│           ▼                         │
│  ┌──────────────────────────────┐  │
│  │ Base LLM Provider            │  │
│  │ (OpenAI, Anthropic, etc.)    │  │
│  └──────────────────────────────┘  │
└─────────────────────────────────────┘
```

### Thread Safety

- Mutex protects state transitions
- Atomic operations for counters
- Safe for concurrent requests

### Performance Impact

- Minimal overhead: ~1-2μs per request
- No additional network calls
- Memory: ~1KB per provider

## Related Documentation

- [Multi-Agent Guide](MULTI_AGENT_GUIDE.md) - Team resilience
- [Safety Levels](SAFETY_LEVELS.md) - Error handling
- [Troubleshooting](troubleshooting.md) - Common issues

## Future Enhancements

- [ ] Prometheus metrics export
- [ ] Configurable error classification rules
- [ ] Adaptive timeout based on provider latency
- [ ] Circuit breaker dashboard
- [ ] Webhook notifications on state changes
- [ ] Per-model circuit breakers (within same provider)

## References

- [Circuit Breaker Pattern](https://martinfowler.com/bliki/CircuitBreaker.html) - Martin Fowler
- [Release Stability Patterns](https://www.oreilly.com/library/view/release-it-2nd/9781680504552/) - Michael Nygard
- Implementation: `pkg/llm/circuit_breaker.go`
