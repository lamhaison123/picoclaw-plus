# Recent Updates & New Features

**Last Updated:** 2026-03-06  
**Version:** 1.1.1+

This document summarizes recent feature additions and improvements to PicoClaw.

---

## 🔄 Circuit Breaker for LLM Providers

**Status:** ✅ Implemented and Documented

### Overview

PicoClaw now includes a robust Circuit Breaker pattern to protect against cascading failures when LLM providers experience issues. This ensures system resilience and prevents resource exhaustion from repeated failed API calls.

### Key Features

- 🔒 **Per-Provider Isolation**: Each provider (OpenAI, Anthropic, Gemini, etc.) has its own independent circuit breaker
- 🎯 **Smart Error Classification**: Only system errors (network, timeout, 5xx) trip the breaker; auth/format errors bypass
- 🔄 **Automatic Recovery**: Self-healing with half-open state for gradual recovery testing
- 📊 **Dual Triggers**: Failure threshold (5 consecutive) AND failure rate (>50%) detection
- 🧵 **Thread-Safe**: Concurrent request handling with mutex protection
- 📈 **Observable**: Built-in metrics for monitoring and alerting

### How It Works

**State Machine:**
```
CLOSED (Normal) → OPEN (Failed) → HALF-OPEN (Testing) → CLOSED
```

**Example Scenario:**
1. OpenAI experiences network issues
2. After 5 consecutive failures → Circuit OPENS
3. All OpenAI requests fail immediately (no waiting)
4. Fallback providers (Claude, Gemini) continue working normally
5. After 30 seconds → Circuit enters HALF-OPEN
6. Limited test requests are sent
7. If successful → Circuit CLOSES, OpenAI back to normal
8. If failed → Circuit returns to OPEN

### Configuration

**Default Settings:**
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

### Benefits

- **Fail Fast**: No waiting for timeouts when provider is down
- **Resource Protection**: Prevents wasting resources on failing providers
- **Automatic Recovery**: Self-healing without manual intervention
- **Provider Isolation**: One provider's failure doesn't affect others
- **Better UX**: Faster fallback to working providers

### Documentation

📖 **Full Guide**: [Circuit Breaker Documentation](CIRCUIT_BREAKER.md)

---

## 💬 Collaborative Chat Enhancements

**Status:** ✅ Stable and Production-Ready

### Overview

IRC-style multi-agent conversations in Telegram and other chat platforms. Mention multiple agents in a single message and they'll all respond with full shared context.

### Key Features

- 🎯 **@mention-based Routing**: Tag agents like `@architect @developer @tester`
- ⚡ **Parallel Execution**: All mentioned agents respond simultaneously
- 🧠 **Shared Context**: Full conversation history available to all agents
- 🎨 **IRC-style Formatting**: Clean, readable output with emojis and session IDs
- 📝 **Session Management**: Persistent sessions per chat with context tracking
- 👥 **Team Management**: `/who` command shows all registered agents and active sessions

### Example Usage

```
User: @architect @developer How should we implement user authentication?

[abc123] 🏗️ ARCHITECT: I recommend using JWT tokens with refresh token rotation...
[abc123] 💻 DEVELOPER: I can implement that using the golang-jwt library...
```

### Commands

- `/who` - Show team status, registered agents, and active agents
- `/help` - Show available commands and usage

### Configuration

```json
{
  "channels": {
    "telegram": {
      "collaborative_chat": {
        "enabled": true,
        "default_team_id": "dev-team",
        "max_context_length": 50,
        "session_timeout": "24h"
      }
    }
  }
}
```

### Documentation

📖 **Quick Start**: [Collaborative Chat Quick Start](COLLABORATIVE_CHAT_QUICKSTART.md)  
📖 **Full Guide**: [Collaborative Chat Documentation](COLLABORATIVE_CHAT.md)  
📖 **Architecture**: [Technical Architecture](COLLABORATIVE_CHAT_ARCHITECTURE.md)

---

## 🔒 Enhanced Safety System

### 4-Level Safety System

Choose the right balance between security and autonomy:

| Level | Best For | Blocks | Allows |
|-------|----------|--------|--------|
| **strict** | Production | sudo, chmod, docker, package install | Read, build, test, safe git |
| **moderate** | Development (default) | Catastrophic ops only | Most dev operations |
| **permissive** | DevOps/Admin | Only catastrophic ops | Almost everything |
| **off** | Testing ⚠️ | Nothing | Everything (DANGEROUS!) |

### Documentation

📖 **Full Guide**: [Safety Levels Documentation](SAFETY_LEVELS.md)  
📖 **Quick Start**: [Safety Quick Start](SAFETY_QUICKSTART.md)

---

## 🤝 Multi-Agent System Improvements

### Enhanced Features

- **Role-based Tool Access**: Fine-grained permission control per agent role
- **Consensus Voting**: Majority, unanimous, or weighted voting strategies
- **Dynamic Composition**: Agents can be added/removed during execution
- **Comprehensive Monitoring**: Real-time status tracking and metrics
- **Automatic Memory Persistence**: Conversation history and context saved automatically

### Collaboration Patterns

- 🔄 **Sequential**: Tasks execute in order (design → implement → test → review)
- ⚡ **Parallel**: Tasks run simultaneously for maximum speed
- 🌳 **Hierarchical**: Complex tasks decompose dynamically

### Documentation

📖 **Full Guide**: [Multi-Agent Guide](MULTI_AGENT_GUIDE.md)  
📖 **Tool Access**: [Team Tool Access Control](TEAM_TOOL_ACCESS.md)  
📖 **Model Selection**: [Per-Role Model Configuration](MULTI_AGENT_MODEL_SELECTION.md)

---

## 📊 What's Next

### Planned Features (from Architect's recommendations)

**Phase 1 (v1.1.2 - 2 weeks):**
- Memory profiling and optimization
- Security hardening (input validation, SSRF protection)
- English documentation for top 5 channels

**Phase 2 (v1.2.0 - 1 month):**
- Provider architecture refactor
- Persistent session storage
- Interactive quickstart wizard
- Observability dashboard

**Phase 3 (v2.0.0 - 3 months):**
- MCP (Model Context Protocol) support
- Browser automation capabilities
- Swarm mode (multi-instance collaboration)
- Web UI

See [ROADMAP.md](../ROADMAP.md) for full details.

---

## 🐛 Known Issues & Limitations

### Current Limitations

1. **Memory Footprint**: Recent PRs increased memory usage to 10-20MB (temporary)
   - Target: Optimize back to <10MB
   - See Architect's recommendations for optimization strategies

2. **Documentation Gaps**:
   - Channel integration guides mostly Chinese-only
   - Need English versions for QQ, DingTalk, WeCom, LINE

3. **Testing**:
   - Circuit breaker needs fault-injection testing in production-like environment
   - Load testing for collaborative chat under high concurrency

### Workarounds

- For memory-constrained environments, use single-agent mode
- For channel guides, use translation tools or refer to Telegram/Discord examples
- Monitor circuit breaker metrics during initial deployment

---

## 📚 Documentation Index

### Core Features
- [Circuit Breaker](CIRCUIT_BREAKER.md) - LLM provider resilience
- [Collaborative Chat](COLLABORATIVE_CHAT.md) - Multi-agent conversations
- [Multi-Agent Guide](MULTI_AGENT_GUIDE.md) - Team collaboration
- [Safety Levels](SAFETY_LEVELS.md) - Security configuration

### Quick Starts
- [Collaborative Chat Quick Start](COLLABORATIVE_CHAT_QUICKSTART.md)
- [Safety Quick Start](SAFETY_QUICKSTART.md)

### Reference
- [Tool Access Control](TEAM_TOOL_ACCESS.md)
- [Model Selection](MULTI_AGENT_MODEL_SELECTION.md)
- [Troubleshooting](troubleshooting.md)

### Main Documentation
- [Documentation Index](INDEX.md) - Complete documentation index
- [Repository Overview](REPOSITORY_OVERVIEW.md) - Project structure
- [Contributing Guide](../CONTRIBUTING.md) - How to contribute

---

## 🤝 Contributing

We welcome contributions! Recent features were developed collaboratively by the team:

- **Circuit Breaker**: Implemented by @developer, reviewed by @architect
- **Collaborative Chat**: Designed by @architect, implemented by @developer
- **Documentation**: Maintained by @developer with input from @manager

See [CONTRIBUTING.md](../CONTRIBUTING.md) for guidelines.

---

## 📧 Support

- **GitHub Issues**: https://github.com/sipeed/picoclaw/issues
- **Discord**: https://discord.gg/V4sAZ9XWpN
- **Twitter**: @SipeedIO
- **Website**: https://picoclaw.io

---

**Made with ❤️ by the PicoClaw community**
