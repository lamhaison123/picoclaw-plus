# Search Providers Guide

## Available Providers (7 total)

### 1. Perplexity (LLM-based)
**Priority**: Highest  
**API Key**: Required  
**Cost**: $$$

```json
{
  "tools": {
    "web": {
      "perplexity": {
        "enabled": true,
        "api_key": "pplx-...",
        "max_results": 5
      }
    }
  }
}
```

### 2. Exa AI (NEW - AI-powered)
**Priority**: High  
**API Key**: Required  
**Cost**: $$

```json
{
  "tools": {
    "web": {
      "exa": {
        "enabled": true,
        "api_key": "your-exa-key",
        "max_results": 5
      }
    }
  }
}
```

### 3. GLM Search (NEW - Chinese)
**Priority**: High  
**API Key**: Required  
**Cost**: $$

```json
{
  "tools": {
    "web": {
      "glm": {
        "enabled": true,
        "api_key": "your-glm-key",
        "max_results": 5
      }
    }
  }
}
```

### 4. Brave
**Priority**: Medium  
**API Key**: Required  
**Cost**: $

```json
{
  "tools": {
    "web": {
      "brave": {
        "enabled": true,
        "api_key": "your-brave-key",
        "max_results": 5
      }
    }
  }
}
```

### 5. Tavily
**Priority**: Medium  
**API Key**: Required  
**Cost**: $$

```json
{
  "tools": {
    "web": {
      "tavily": {
        "enabled": true,
        "api_key": "tvly-...",
        "max_results": 5
      }
    }
  }
}
```

### 6. SearXNG (NEW - Privacy)
**Priority**: Low  
**API Key**: Not required  
**Cost**: Free (self-hosted)

```json
{
  "tools": {
    "web": {
      "searxng": {
        "enabled": true,
        "base_url": "https://searx.example.com",
        "max_results": 5
      }
    }
  }
}
```

### 7. DuckDuckGo (Fallback)
**Priority**: Lowest  
**API Key**: Not required  
**Cost**: Free

```json
{
  "tools": {
    "web": {
      "duckduckgo": {
        "enabled": true,
        "max_results": 5
      }
    }
  }
}
```

## Priority Chain

```
Perplexity > Exa > GLM > Brave > Tavily > SearXNG > DuckDuckGo
```

First enabled provider in chain is used.

## Use Cases

### Privacy-Focused
Enable: SearXNG + DuckDuckGo

### Chinese Market
Enable: GLM + Brave

### AI-Powered
Enable: Exa + Perplexity

### Cost-Optimized
Enable: SearXNG + DuckDuckGo

### Best Quality
Enable: Perplexity + Exa

---

**Updated**: 2026-03-09
