# Memory Configuration Guide

## Overview

PicoClaw hỗ trợ cấu hình linh hoạt cho memory system với khả năng bật/tắt từng component riêng biệt.

---

## Config Structure

```json
{
  "memory": {
    "enabled": true,              // Master switch - tắt toàn bộ memory
    "embedding": { ... },          // Embedding service config
    "vector_store": { ... },       // Vector database config
    "memory_provider": { ... },    // Memory provider config
    "cache": { ... }               // Cache config
  }
}
```

---

## 1. Master Switch

### Tắt toàn bộ Memory System
```json
{
  "memory": {
    "enabled": false
  }
}
```

**Effect:** Tất cả memory features bị tắt, agent hoạt động như bình thường không có memory.

---

## 2. Embedding Configuration

### OpenAI Embeddings (Recommended)
```json
{
  "memory": {
    "enabled": true,
    "embedding": {
      "provider": "openai",
      "model": "text-embedding-3-small",
      "dimension": 384,
      "api_key": "${PICOCLAW_EMBEDDING_API_KEY}",
      "base_url": "",
      "timeout_ms": 10000
    }
  }
}
```

### Custom / Local Embeddings (vLLM, Ollama, LiteLLM)
```json
{
  "memory": {
    "enabled": true,
    "embedding": {
      "provider": "vllm",
      "model": "mistral/mistral-embed",
      "dimension": 1024,
      "api_key": "sk-dummy",
      "base_url": "https://api.tsdreamer.io.vn/v1",
      "timeout_ms": 10000
    }
  }
}
```
> **Lưu ý quan trọng**: Các cấu hình này nằm ở cấp ngoài cùng của object `embedding`, KHÔNG bọc vào object con (như `"vllm": { ... }`). Key cấu hình đường dẫn là `"base_url"`, không dùng `"api_base"`.
> Nếu dùng endpoint tương thích OpenAI cho các model nội bộ (như Mistral), hãy đảm bảo cung cấp đúng tiền tố nếu router yêu cầu (ví dụ: `"model": "mistral/mistral-embed"`).

### Disable Embeddings
```json
{
  "memory": {
    "enabled": true,
    "embedding": {
      "provider": "none"
    }
  }
}
```

**Models Available:**
- `text-embedding-3-small` (384 dims) - Fast, cheap (OpenAI)
- `text-embedding-3-large` (1536 dims) - Better quality (OpenAI)
- *Local Models* (requires setting correct `dimension` matching the model output)

---

## 3. Vector Store Configuration

### Option A: Qdrant Only (Recommended)
```json
{
  "memory": {
    "enabled": true,
    "vector_store": {
      "provider": "qdrant",
      "qdrant": {
        "enabled": true,
        "url": "http://localhost:6333",
        "api_key": "",
        "collection": "picoclaw_memory",
        "dimension": 384,
        "timeout_ms": 800
      },
      "lancedb": {
        "enabled": false
      }
    }
  }
}
```

**Use when:**
- ✅ Need production-ready vector DB
- ✅ Want gRPC performance
- ✅ Need clustering/sharding
- ✅ Have Qdrant server running

**Start Qdrant:**
```bash
docker run -p 6333:6333 -p 6334:6334 qdrant/qdrant
```
> **Flow tự động lưu trữ (Auto-Memory):**
> Khi Vector Store được bật và cấu hình hợp lệ, `AgentLoop` sẽ **tự động** chạy một luồng bất đồng bộ (goroutine) lưu nội dung trò chuyện (prompt của User và đáp án của Agent) vào Qdrant sau mỗi phản hồi mà không cần gọi API thủ công. 
> 
> *Lưu ý Port kết nối:* Client trong Go mặc định sử dụng gRPC pipeline vì vậy nếu bạn truyền `http://localhost:6333` (REST), PicoClaw sẽ tự động map và kết nối tới port gRPC mặc định (`6334`).

### Option B: LanceDB Only

For embedded local storage with no external dependencies (requires CGO):

```json
{
  "memory": {
    "enabled": true,
    "vector_store": {
      "provider": "lancedb",
      "qdrant": {
        "enabled": false
      },
      "lancedb": {
        "enabled": true,
        "mode": "api",
        "url": "http://localhost:8000",
        "dataset": "picoclaw",
        "timeout_ms": 800
      }
    }
  }
}
```

**Use when:**
- ✅ Need embedded database
- ✅ Want simpler deployment
- ✅ Don't need clustering

This configuration disables Qdrant and purely relies on local LanceDB.

### Option C: Disable Vector Store
```json
{
  "memory": {
    "enabled": true,
    "vector_store": {
      "provider": "none",
      "qdrant": {
        "enabled": false
      },
      "lancedb": {
        "enabled": false
      }
    }
  }
}
```

---

## 4. Memory Provider Configuration

### Option A: Disable (Default)
```json
{
  "memory": {
    "enabled": true,
    "memory_provider": {
      "provider": "none",
      "sidecar": {
        "enabled": false
      },
      "mem0": {
        "enabled": false
      },
      "mindgraph": {
        "enabled": false
      }
    }
  }
}
```

**Use when:**
- ✅ Only need vector search (no advanced memory)
- ✅ Want simplest setup
- ✅ Don't need personalization

### Option B: Sidecar (Python-based)
```json
{
  "memory": {
    "enabled": true,
    "memory_provider": {
      "provider": "sidecar",
      "sidecar": {
        "enabled": true,
        "endpoint": "http://localhost:8765",
        "timeout_ms": 1200,
        "circuit_breaker": {
          "max_failures": 5,
          "reset_timeout_s": 60
        }
      }
    }
  }
}
```

**Use when:**
- ✅ Need Mem0 or MindGraph features
- ✅ Want Python-based memory providers
- ✅ Have sidecar server running

### Option C: Mem0 Direct
```json
{
  "memory": {
    "enabled": true,
    "memory_provider": {
      "provider": "mem0",
      "mem0": {
        "enabled": true,
        "url": "http://localhost:8001",
        "api_key": "",
        "timeout_ms": 1200
      }
    }
  }
}
```

### Option D: MindGraph Direct
```json
{
  "memory": {
    "enabled": true,
    "memory_provider": {
      "provider": "mindgraph",
      "mindgraph": {
        "enabled": true,
        "url": "http://localhost:8002",
        "timeout_ms": 1200
      }
    }
  }
}
```

---

## 5. Cache Configuration

### Enable Cache (Recommended)
```json
{
  "memory": {
    "enabled": true,
    "cache": {
      "enabled": true,
      "max_entries": 1000,
      "ttl_seconds": 3600
    }
  }
}
```

**Benefits:**
- ✅ Faster repeated queries
- ✅ Reduced API calls
- ✅ Lower latency

### Disable Cache
```json
{
  "memory": {
    "enabled": true,
    "cache": {
      "enabled": false
    }
  }
}
```

---

## Common Configurations

### 1. Development (Minimal)
**File:** `config/config.memory.disabled.json`
```json
{
  "memory": {
    "enabled": false
  }
}
```

**Use:** Fastest startup, no dependencies

---

### 2. Production (Qdrant + OpenAI)
**File:** `config/config.memory.qdrant-only.json`
```json
{
  "memory": {
    "enabled": true,
    "embedding": {
      "provider": "openai",
      "model": "text-embedding-3-small",
      "dimension": 384,
      "api_key": "${PICOCLAW_EMBEDDING_API_KEY}"
    },
    "vector_store": {
      "provider": "qdrant",
      "qdrant": {
        "enabled": true,
        "url": "http://localhost:6333",
        "collection": "picoclaw_memory"
      }
    },
    "cache": {
      "enabled": true
    }
  }
}
```

**Use:** Best for production with semantic search

---

### 3. Full Stack (All Features)
**File:** `config/config.memory.flexible.json`
```json
{
  "memory": {
    "enabled": true,
    "embedding": {
      "provider": "openai",
      "model": "text-embedding-3-small",
      "dimension": 384,
      "api_key": "${PICOCLAW_EMBEDDING_API_KEY}"
    },
    "vector_store": {
      "provider": "qdrant",
      "qdrant": {
        "enabled": true,
        "url": "http://localhost:6333"
      }
    },
    "memory_provider": {
      "provider": "sidecar",
      "sidecar": {
        "enabled": true,
        "endpoint": "http://localhost:8765"
      }
    },
    "cache": {
      "enabled": true
    }
  }
}
```

**Use:** Maximum features, requires all services running

---

## Environment Variables

### Override via Environment
```bash
# Embedding
export PICOCLAW_EMBEDDING_API_KEY="sk-..."
export PICOCLAW_EMBEDDING_BASE_URL="https://api.openai.com/v1"

# Qdrant
export PICOCLAW_MEMORY_VECTOR_QDRANT_ENABLED="true"
export PICOCLAW_MEMORY_VECTOR_QDRANT_URL="http://localhost:6333"
export PICOCLAW_MEMORY_VECTOR_QDRANT_COLLECTION="picoclaw_memory"

# LanceDB
export PICOCLAW_MEMORY_VECTOR_LANCEDB_ENABLED="false"

# Sidecar
export PICOCLAW_MEMORY_SIDECAR_ENABLED="false"
```

**Priority:** Environment variables > Config file

---

## Feature Matrix

| Feature | Qdrant | LanceDB | Mem0 | MindGraph | Sidecar |
|---------|--------|---------|------|-----------|---------|
| Vector Search | ✅ | 🚧 | ❌ | ❌ | ✅ |
| Semantic Memory | ✅ | 🚧 | ✅ | ✅ | ✅ |
| Personalization | ❌ | ❌ | ✅ | ❌ | ✅ |
| Knowledge Graph | ❌ | ❌ | ❌ | ✅ | ✅ |
| No External Deps | ❌ | ✅* | ❌ | ❌ | ❌ |
| Production Ready | ✅ | 🚧 | ✅ | 🚧 | 🚧 |

Legend:
- ✅ Available
- 🚧 Coming soon
- ❌ Not supported
- \* CGO mode only

---

## Troubleshooting

### Memory not working?

1. **Check master switch:**
   ```json
   "memory": { "enabled": true }
   ```

2. **Check provider enabled:**
   ```json
   "qdrant": { "enabled": true }
   ```

3. **Check service running:**
   ```bash
   curl http://localhost:6333/health  # Qdrant
   ```

4. **Check API key:**
   ```bash
   echo $PICOCLAW_EMBEDDING_API_KEY
   ```

5. **Check logs:**
   ```
   INFO: Vector memory initialized
   ```

### Build errors?

**With Qdrant:**
```bash
go build -o picoclaw ./cmd/picoclaw
```

**Without Qdrant:**
```bash
go build -tags=no_qdrant -o picoclaw ./cmd/picoclaw
```

---

## Best Practices

### 1. Start Simple
```json
{
  "memory": {
    "enabled": false
  }
}
```
Get agent working first, then enable memory.

### 2. Add Vector Search
```json
{
  "memory": {
    "enabled": true,
    "embedding": { "provider": "openai" },
    "vector_store": { "provider": "qdrant" }
  }
}
```
Most useful feature, minimal complexity.

### 3. Add Advanced Features
```json
{
  "memory": {
    "enabled": true,
    "memory_provider": { "provider": "sidecar" }
  }
}
```
Only when you need personalization/knowledge graph.

---

## Summary

**Flexibility Levels:**

1. **Disabled** - No memory, fastest
2. **Vector Only** - Semantic search, recommended
3. **Full Stack** - All features, most complex

**Control Points:**

- ✅ Master switch: `memory.enabled`
- ✅ Per-provider: `qdrant.enabled`, `lancedb.enabled`
- ✅ Per-feature: `cache.enabled`, `sidecar.enabled`
- ✅ Environment variables for all settings

**Recommendation:** Start with Qdrant-only config, expand as needed.
