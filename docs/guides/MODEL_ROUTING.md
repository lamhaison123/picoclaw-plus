# Model Routing Guide

Automatically route queries to appropriate models based on complexity.

## Overview

Model routing analyzes query complexity and selects the most cost-effective model while maintaining quality.

## Benefits

### Cost Optimization
- Use cheap models for simple queries
- Use expensive models for complex tasks
- Automatic selection based on complexity

### Performance
- Fast models for simple queries
- Powerful models for complex tasks
- Optimal resource utilization

### Flexibility
- Configurable model tiers
- Adjustable complexity thresholds
- Per-agent configuration

## Configuration

### Basic Setup
```json
{
  "routing": {
    "enabled": true,
    "tiers": [
      {
        "name": "cheap",
        "models": ["gpt-4o-mini", "llama-3.3-70b"]
      },
      {
        "name": "expensive",
        "models": ["gpt-5.2", "claude-opus"]
      }
    ]
  }
}
```

### Three-Tier Setup
```json
{
  "routing": {
    "enabled": true,
    "tiers": [
      {
        "name": "cheap",
        "models": ["gpt-4o-mini"]
      },
      {
        "name": "medium",
        "models": ["gpt-4o"]
      },
      {
        "name": "expensive",
        "models": ["gpt-5.2", "claude-opus"]
      }
    ],
    "thresholds": {
      "low": 100,
      "medium": 500
    }
  }
}
```

### Environment Variables
```bash
PICOCLAW_ROUTING_ENABLED=true
```

## Complexity Scoring

### Features Analyzed
1. **Token Estimate**: Message length (CJK-aware)
2. **Code Blocks**: Number of code blocks
3. **Tool Usage**: Recent tool calls
4. **Conversation Depth**: Number of messages
5. **Attachments**: Media files

### Scoring Formula
```
complexity = (
    token_estimate * 1.0 +
    code_blocks * 50 +
    tool_calls * 30 +
    depth * 10 +
    attachments * 100
)
```

### Complexity Levels
- **Low** (<100): Simple questions, greetings
- **Medium** (100-500): Code tasks, analysis
- **High** (>500): Complex reasoning, multiple tools

## Examples

### Simple Query (Cheap Model)
```bash
> What is 2+2?
# Routed to: gpt-4o-mini
```

### Code Task (Medium Model)
```bash
> Write a Python function to sort a list
# Routed to: gpt-4o
```

### Complex Analysis (Expensive Model)
```bash
> Analyze this codebase and suggest improvements
# Routed to: gpt-5.2
```

### With Image (Expensive Model)
```bash
> Describe this image: /path/to/chart.png
# Routed to: gpt-5.2 (has vision)
```

## Model Selection

### Priority
1. Check if routing enabled
2. Calculate complexity score
3. Select tier based on score
4. Choose first available model in tier
5. Fallback to primary model if none available

### Tier Selection
```
if complexity < 100:
    tier = "cheap"
elif complexity < 500:
    tier = "medium"
else:
    tier = "expensive"
```

## Logging

### Routing Decisions
```
[INFO] Model routing: complexity=150, tier=medium, model=gpt-4o
[INFO] Model routing: complexity=50, tier=cheap, model=gpt-4o-mini
```

### Disabled Routing
```
[INFO] Model routing disabled, using primary model: gpt-4o
```

## Best Practices

### Model Tiers
- **Cheap**: Fast, inexpensive models for simple tasks
- **Medium**: Balanced models for most tasks
- **Expensive**: Powerful models for complex reasoning

### Threshold Tuning
- Start with defaults (100, 500)
- Monitor routing decisions
- Adjust based on cost/quality balance

### Model Selection
- Include multiple models per tier
- Order by preference
- Ensure all models are configured

## Cost Analysis

### Example Savings
```
Without Routing:
- 100 queries * $0.01 (gpt-5.2) = $1.00

With Routing:
- 70 simple * $0.001 (gpt-4o-mini) = $0.07
- 20 medium * $0.005 (gpt-4o) = $0.10
- 10 complex * $0.01 (gpt-5.2) = $0.10
Total: $0.27 (73% savings)
```

## Troubleshooting

### Always Using Expensive Model
Check:
- Routing enabled in config
- Thresholds not too low
- Cheap models configured

### Poor Quality Responses
Try:
- Lower thresholds
- Better medium-tier models
- Adjust complexity scoring

### Routing Not Working
Verify:
- `routing.enabled = true`
- Models exist in config
- Logs show routing decisions

## Advanced Configuration

### Per-Agent Routing
```json
{
  "agents": {
    "code-agent": {
      "routing": {
        "enabled": true,
        "tiers": [
          {"name": "cheap", "models": ["codellama-70b"]},
          {"name": "expensive", "models": ["gpt-5.2"]}
        ]
      }
    }
  }
}
```

### Custom Thresholds
```json
{
  "routing": {
    "thresholds": {
      "low": 50,
      "medium": 300
    }
  }
}
```

## API Reference

### Complexity Scorer
```go
type ComplexityScorer struct {
    TokenEstimate    int
    CodeBlockCount   int
    RecentToolCalls  int
    ConversationDepth int
    HasAttachments   bool
}

func (s *ComplexityScorer) Score() int
```

### Router
```go
type Router struct {
    Config RoutingConfig
}

func (r *Router) SelectModel(
    messages []Message,
    primaryModel string,
) string
```

## See Also

- [Configuration Guide](../reference/CONFIGURATION.md)
- [Model Configuration](../reference/MODEL_CONFIGURATION.md)
- [v0.2.1 Features](V0.2.1_FEATURES.md)

---

**Version**: v0.2.1  
**Last Updated**: 2026-03-09
