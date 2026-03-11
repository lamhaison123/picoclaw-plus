# MindGraph Integration Guide

**Version**: v0.2.1+  
**Status**: Implemented  
**Package**: `pkg/memory/vector`

## Overview

MindGraph is a knowledge graph memory provider for AI agents that structures memory into deterministic layers, mirroring how cognitive systems reason from raw observation to final action.

PicoClaw now includes a fully functional MindGraph client that implements the `MemoryProvider` interface, allowing agents to store, recall, update, and delete memories using MindGraph's knowledge graph architecture.

## Features

- ✅ Store memories with metadata
- ✅ Recall memories based on queries
- ✅ Update existing memories
- ✅ Delete memories
- ✅ Health checks
- ✅ Circuit breaker protection
- ✅ Bearer token authentication
- ✅ Configurable timeouts
- ✅ Support for both self-hosted and cloud deployments

## Configuration

### Self-Hosted MindGraph

```json
{
  "memory": {
    "enabled": true,
    "memory_provider": {
      "provider": "mindgraph",
      "mindgraph": {
        "enabled": true,
        "url": "http://localhost:8002",
        "api_key": "${MINDGRAPH_API_KEY}",
        "timeout_ms": 1200
      }
    }
  }
}
```

### MindGraph Cloud

```json
{
  "memory": {
    "enabled": true,
    "memory_provider": {
      "provider": "mindgraph",
      "mindgraph": {
        "enabled": true,
        "url": "https://api.mindgraph.cloud",
        "api_key": "${MINDGRAPH_API_KEY}",
        "timeout_ms": 1200
      }
    }
  }
}
```

### Environment Variables

Add to your `.env` file:

```bash
# MindGraph Configuration
PICOCLAW_MEMORY_MINDGRAPH_ENABLED=true
PICOCLAW_MEMORY_MINDGRAPH_URL=https://api.mindgraph.cloud
MINDGRAPH_API_KEY=your-api-key-here
PICOCLAW_MEMORY_MINDGRAPH_TIMEOUT_MS=1200
```

## API Endpoints

The MindGraph client uses the following REST API endpoints:

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/memory` | Store a new memory |
| POST | `/api/v1/recall` | Recall memories by query |
| PUT | `/api/v1/memory/{id}` | Update existing memory |
| DELETE | `/api/v1/memory/{id}` | Delete a memory |
| GET | `/api/v1/health` | Health check |

## Usage Examples

### Store a Memory

```go
import (
    "context"
    memory "github.com/sipeed/picoclaw/pkg/memory/vector"
)

// Create client
breaker := memory.NewCircuitBreaker(memory.CircuitBreakerConfig{
    MaxFailures:   5,
    ResetTimeoutS: 30,
})

client, err := memory.NewMindGraphClient(memory.MindGraphConfig{
    URL:       "https://api.mindgraph.cloud",
    APIKey:    "your-api-key",
    TimeoutMS: 1200,
}, breaker)
if err != nil {
    log.Fatal(err)
}
defer client.Close()

// Store memory
ctx := context.Background()
metadata := map[string]interface{}{
    "type": "conversation",
    "user": "user-123",
    "timestamp": time.Now().Unix(),
}

id, err := client.Store(ctx, "User prefers dark mode", metadata)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Stored memory with ID: %s\n", id)
```

### Recall Memories

```go
// Recall relevant memories
memories, err := client.Recall(ctx, "user preferences", 10)
if err != nil {
    log.Fatal(err)
}

for _, mem := range memories {
    fmt.Printf("Memory %s: %s\n", mem.ID, mem.Content)
    fmt.Printf("  Metadata: %v\n", mem.Metadata)
    fmt.Printf("  Score: %.2f\n", mem.Score)
}
```

### Update a Memory

```go
// Update existing memory
err = client.Update(ctx, "mem-123", "User prefers dark mode and large fonts", metadata)
if err != nil {
    log.Fatal(err)
}
```

### Delete a Memory

```go
// Delete memory
err = client.Delete(ctx, "mem-123")
if err != nil {
    log.Fatal(err)
}
```

### Health Check

```go
// Check MindGraph health
err = client.Health(ctx)
if err != nil {
    log.Printf("MindGraph unhealthy: %v\n", err)
} else {
    log.Println("MindGraph is healthy")
}
```

## Authentication

The MindGraph client supports Bearer token authentication. The API key is automatically added to requests as:

```
Authorization: Bearer <your-api-key>
```

## Circuit Breaker

The client includes circuit breaker protection to prevent cascading failures:

- **Max Failures**: 5 (default)
- **Reset Timeout**: 30 seconds (default)
- **Half-Open Max**: 3 requests (default)

When the circuit is open, requests will fail immediately with `ERR_CIRCUIT_OPEN` error.

## Error Handling

The client returns standardized error codes:

| Error Code | Description |
|------------|-------------|
| `ERR_CONFIG_INVALID` | Invalid configuration |
| `ERR_PROVIDER_UNAVAILABLE` | MindGraph service unavailable |
| `ERR_TIMEOUT` | Request timeout |
| `ERR_CIRCUIT_OPEN` | Circuit breaker is open |
| `ERR_INTERNAL` | Internal error |

Example error handling:

```go
id, err := client.Store(ctx, content, metadata)
if err != nil {
    if strings.Contains(err.Error(), "ERR_CIRCUIT_OPEN") {
        log.Println("Circuit breaker is open, service temporarily unavailable")
    } else if strings.Contains(err.Error(), "ERR_TIMEOUT") {
        log.Println("Request timed out")
    } else {
        log.Printf("Error: %v\n", err)
    }
    return
}
```

## Memory Structure

MindGraph organizes memories into layers:

1. **Reality Layer**: Raw observations and facts
2. **Epistemic Layer**: Claims, evidence, hypotheses
3. **Intent Layer**: Goals, decisions, questions
4. **Action Layer**: Tasks, flows, controls
5. **Memory Layer**: Sessions, summaries

The PicoClaw client provides a simplified interface that maps to these layers automatically.

## Performance Considerations

### Timeouts

- Default timeout: 1200ms
- Recommended for cloud: 2000-5000ms
- Recommended for local: 500-1000ms

### Caching

Enable memory caching to reduce API calls:

```json
{
  "memory": {
    "cache": {
      "enabled": true,
      "max_entries": 1000,
      "ttl_seconds": 3600
    }
  }
}
```

### Rate Limiting

Be mindful of API rate limits when using MindGraph Cloud:

- Free tier: 10,000 nodes, 50,000 requests/month
- Implement exponential backoff for rate limit errors

## Comparison with Other Providers

| Feature | MindGraph | Mem0 | Qdrant |
|---------|-----------|------|--------|
| Type | Knowledge Graph | Memory Layer | Vector Store |
| Structure | Typed nodes & relationships | Flat memories | Vectors only |
| Provenance | ✅ Built-in | ❌ No | ❌ No |
| Reasoning | ✅ Explicit | ❌ No | ❌ No |
| Semantic Search | ✅ Yes | ✅ Yes | ✅ Yes |
| Self-hosted | 🚧 Coming | ✅ Yes | ✅ Yes |
| Cloud | ✅ Yes | ✅ Yes | ✅ Yes |

## Troubleshooting

### Connection Errors

```bash
# Test connectivity
curl -H "Authorization: Bearer your-api-key" \
     https://api.mindgraph.cloud/api/v1/health
```

### Circuit Breaker Open

If you see `ERR_CIRCUIT_OPEN` errors:

1. Check MindGraph service health
2. Wait for reset timeout (30s default)
3. Verify API key is valid
4. Check network connectivity

### Timeout Issues

If requests timeout frequently:

1. Increase `timeout_ms` in config
2. Check network latency
3. Verify MindGraph service performance
4. Consider using local deployment

## Migration from Other Providers

### From Mem0

```go
// Before (Mem0)
mem0Client.Add(ctx, content, metadata)

// After (MindGraph)
mindgraphClient.Store(ctx, content, metadata)
```

### From Qdrant

MindGraph provides higher-level memory operations compared to Qdrant's vector-only approach. Consider using both:

- **Qdrant**: For fast semantic search
- **MindGraph**: For structured knowledge and reasoning

## Best Practices

1. **Use Metadata**: Add rich metadata to memories for better organization
2. **Enable Caching**: Reduce API calls with local caching
3. **Monitor Health**: Regularly check service health
4. **Handle Errors**: Implement proper error handling and retries
5. **Batch Operations**: Group related memories when possible
6. **Clean Up**: Delete obsolete memories to manage costs

## Resources

- **MindGraph Website**: https://mindgraph.cloud/
- **API Documentation**: (Contact MindGraph for access)
- **PicoClaw Memory Guide**: `docs/memory/CONFIG_GUIDE.md`
- **Source Code**: `pkg/memory/vector/mindgraph_client.go`

## Support

For MindGraph-specific issues:
- Visit: https://mindgraph.cloud/
- Check their documentation and support channels

For PicoClaw integration issues:
- GitHub Issues: https://github.com/sipeed/picoclaw/issues
- Documentation: `docs/`

---

**Last Updated**: 2026-03-09  
**Version**: v0.2.1+  
**Status**: Production Ready
