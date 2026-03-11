# Search Providers Reference

Complete guide to all 7 search providers available in PicoClaw.

## Overview

PicoClaw supports 7 search providers with different strengths and use cases.

## Provider Comparison

| Provider | API Key | Cost | Privacy | Language | Quality |
|----------|---------|------|---------|----------|---------|
| Perplexity | ✅ | $$$ | ⭐⭐⭐ | All | ⭐⭐⭐⭐⭐ |
| Exa AI | ✅ | $$ | ⭐⭐⭐ | All | ⭐⭐⭐⭐⭐ |
| GLM | ✅ | $$ | ⭐⭐⭐ | Chinese | ⭐⭐⭐⭐ |
| Brave | ✅ | $ | ⭐⭐⭐⭐ | All | ⭐⭐⭐⭐ |
| Tavily | ✅ | $$ | ⭐⭐⭐ | All | ⭐⭐⭐⭐ |
| SearXNG | ❌ | Free | ⭐⭐⭐⭐⭐ | All | ⭐⭐⭐ |
| DuckDuckGo | ❌ | Free | ⭐⭐⭐⭐ | All | ⭐⭐⭐ |

## Priority Chain

```
Perplexity > Exa > GLM > Brave > Tavily > SearXNG > DuckDuckGo
```

First enabled provider in chain is used.

## Providers

### 1. Perplexity (LLM-based)

**Best For**: Highest quality results  
**Type**: LLM-powered search  
**Cost**: $$$

**Configuration**:
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

**Environment**:
```bash
PICOCLAW_TOOLS_WEB_PERPLEXITY_ENABLED=true
PICOCLAW_TOOLS_WEB_PERPLEXITY_API_KEY=pplx-...
```

---

### 2. Exa AI (AI-powered)

**Best For**: Semantic search, AI-optimized  
**Type**: AI-native search engine  
**Cost**: $$

**Features**:
- Autoprompt optimization
- Live crawling
- Content extraction
- Highlight extraction

**Configuration**:
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

**Environment**:
```bash
PICOCLAW_TOOLS_WEB_EXA_ENABLED=true
PICOCLAW_TOOLS_WEB_EXA_API_KEY=your-key
```

**Get API Key**: https://exa.ai/

---

### 3. GLM Search (Chinese)

**Best For**: Chinese language queries  
**Type**: AI-powered Chinese search  
**Cost**: $$  
**Provider**: Zhipu AI (智谱AI)

**Configuration**:
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

**Environment**:
```bash
PICOCLAW_TOOLS_WEB_GLM_ENABLED=true
PICOCLAW_TOOLS_WEB_GLM_API_KEY=your-key
```

**Get API Key**: https://open.bigmodel.cn/

---

### 4. Brave

**Best For**: Privacy + quality balance  
**Type**: Privacy-focused search  
**Cost**: $

**Configuration**:
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

**Environment**:
```bash
PICOCLAW_TOOLS_WEB_BRAVE_ENABLED=true
PICOCLAW_TOOLS_WEB_BRAVE_API_KEY=your-key
```

---

### 5. Tavily

**Best For**: AI-focused search  
**Type**: AI-optimized search  
**Cost**: $$

**Configuration**:
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

**Environment**:
```bash
PICOCLAW_TOOLS_WEB_TAVILY_ENABLED=true
PICOCLAW_TOOLS_WEB_TAVILY_API_KEY=tvly-...
```

---

### 6. SearXNG (Privacy)

**Best For**: Maximum privacy, self-hosted  
**Type**: Metasearch engine  
**Cost**: Free (self-hosted)

**Features**:
- No tracking
- No profiling
- Aggregates multiple engines
- Open source

**Configuration**:
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

**Environment**:
```bash
PICOCLAW_TOOLS_WEB_SEARXNG_ENABLED=true
PICOCLAW_TOOLS_WEB_SEARXNG_BASE_URL=https://searx.example.com
```

**Setup**: Requires your own SearXNG instance  
**Public Instances**: https://searx.space/

---

### 7. DuckDuckGo (Fallback)

**Best For**: No API key needed  
**Type**: Privacy-focused search  
**Cost**: Free

**Configuration**:
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

**Environment**:
```bash
PICOCLAW_TOOLS_WEB_DUCKDUCKGO_ENABLED=true
```

---

## Use Cases

### Privacy-Focused
```json
{
  "tools": {
    "web": {
      "searxng": {"enabled": true, "base_url": "..."},
      "duckduckgo": {"enabled": true}
    }
  }
}
```

### Chinese Market
```json
{
  "tools": {
    "web": {
      "glm": {"enabled": true, "api_key": "..."}
    }
  }
}
```

### AI-Powered
```json
{
  "tools": {
    "web": {
      "exa": {"enabled": true, "api_key": "..."},
      "perplexity": {"enabled": true, "api_key": "..."}
    }
  }
}
```

### Cost-Optimized
```json
{
  "tools": {
    "web": {
      "searxng": {"enabled": true, "base_url": "..."},
      "duckduckgo": {"enabled": true}
    }
  }
}
```

### Best Quality
```json
{
  "tools": {
    "web": {
      "perplexity": {"enabled": true, "api_key": "..."},
      "exa": {"enabled": true, "api_key": "..."}
    }
  }
}
```

## Proxy Support

All providers support proxy configuration:

```json
{
  "tools": {
    "web": {
      "proxy": "http://proxy.example.com:8080"
    }
  }
}
```

**Environment**:
```bash
PICOCLAW_TOOLS_WEB_PROXY=http://proxy.example.com:8080
```

## Troubleshooting

### No Search Results
Check:
- Provider is enabled
- API key is valid (if required)
- Network connectivity
- Proxy settings

### Poor Quality Results
Try:
- Different provider
- More specific query
- Adjust max_results

### Rate Limiting
- Check API quota
- Use different provider
- Add delay between requests

## Best Practices

### Production
- Use paid providers for reliability
- Configure multiple providers as fallback
- Monitor API usage and costs

### Development
- Use free providers (DuckDuckGo)
- Test with different providers
- Verify API keys work

### Privacy
- Use SearXNG or DuckDuckGo
- Avoid logging search queries
- Use proxy if needed

## See Also

- [Configuration Guide](CONFIGURATION.md)
- [Web Tools](tools_configuration.md)
- [v0.2.1 Features](../guides/V0.2.1_FEATURES.md)

---

**Version**: v0.2.1  
**Last Updated**: 2026-03-09
