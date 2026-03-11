# Vector Memory Integration Guide

## Quick Start

### 1. Basic Setup

```go
import "github.com/picoclaw/pkg/memory/vector"

// Create Qdrant store with default config
config := vector.DefaultConfig()
config.QdrantURL = "http://localhost:6333"
config.CollectionName = "my_vectors"

store, err := vector.NewQdrantStore(config)
if err != nil {
    log.Fatal(err)
}
defer store.Close()
```

### 2. Upsert Vectors

```go
vectors := []vector.Vector{
    {
        ID:        "doc1",
        Embedding: []float32{0.1, 0.2, 0.3, ...}, // 384 dimensions
        Metadata: map[string]interface{}{
            "text": "Hello world",
            "timestamp": time.Now().Unix(),
        },
    },
}

ctx := context.Background()
err = store.Upsert(ctx, vectors)
if err != nil {
    // Handle error - check if it's timeout or circuit breaker
    log.Printf("Upsert failed: %v", err)
}
```

### 3. Search Similar Vectors

```go
query := []float32{0.1, 0.2, 0.3, ...} // 384 dimensions

results, err := store.Search(ctx, query, 5) // Top 5 results
if err != nil {
    log.Printf("Search failed: %v", err)
}

for _, result := range results {
    fmt.Printf("ID: %s, Score: %.4f\n", result.ID, result.Score)
    fmt.Printf("Metadata: %v\n", result.Metadata)
}
```

## Circuit Breaker Behavior

The system automatically protects against provider failures:

**States:**
- **Closed** (Normal): All requests go through
- **Open** (Failing): Requests fail fast after 5 consecutive failures
- **Half-Open** (Testing): Allows 3 test requests after 30s cooldown

**Error Handling:**
```go
err := store.Upsert(ctx, vectors)
if err != nil {
    if strings.Contains(err.Error(), "ERR_CIRCUIT_OPEN") {
        // Circuit breaker is open - provider is down
        // Wait 30s before retry
    } else if strings.Contains(err.Error(), "ERR_TIMEOUT") {
        // Request timed out (800ms for search, 1200ms for upsert)
        // Safe to retry
    }
}
```

## Retry Logic

System automatically retries on transient failures:
- **Max attempts:** 2
- **Backoff:** 100ms exponential
- **Retries on:** Network errors, temporary Qdrant failures
- **No retry on:** Context cancellation, circuit breaker open

## Timeout Configuration

```go
config := vector.DefaultConfig()
config.SearchTimeout = 800 * time.Millisecond  // Vector search
config.UpsertTimeout = 1200 * time.Millisecond // Memory recall
config.WatchdogTimeout = 5000 * time.Millisecond // Max operation time
```

**Important:** If parent context has deadline, it takes priority:
```go
// This will use 500ms timeout (parent wins)
ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
defer cancel()
err := store.Search(ctx, query, 5)
```

## Monitoring Circuit Breaker

```go
metrics := store.GetMetrics()
fmt.Printf("State: %s\n", metrics.State)
fmt.Printf("Total calls: %d\n", metrics.TotalCalls)
fmt.Printf("Success rate: %.2f%%\n", 
    float64(metrics.TotalSuccess)/float64(metrics.TotalCalls)*100)
```

## Testing Without Qdrant

For development/testing without real Qdrant instance:

```go
// Use in-memory mock (coming in Phase 3)
store := vector.NewMockStore()
```

## Common Issues

**1. Circuit breaker opens immediately**
- Check Qdrant is running: `curl http://localhost:6333/health`
- Verify collection exists
- Check network connectivity

**2. Timeouts on every request**
- Increase timeout in config
- Check Qdrant performance
- Verify vector dimensions match (384)

**3. Context cancellation errors**
- These are NOT counted as provider failures
- Safe to retry with new context
- Check parent context deadline

## Production Checklist

- ✅ Qdrant endpoint configured
- ✅ Collection created with correct dimension (384)
- ✅ Timeouts tuned for your workload
- ✅ Circuit breaker thresholds validated
- ✅ Monitoring/alerting on circuit breaker state
- ✅ Retry logic tested under load

## Support

For issues or questions:
- Check logs for detailed error messages
- Review circuit breaker metrics
- Verify Qdrant health endpoint
- Contact: @developer @architect
