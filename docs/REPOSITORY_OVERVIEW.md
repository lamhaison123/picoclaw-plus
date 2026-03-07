# PicoClaw Repository Overview

## 📋 Project Overview

**PicoClaw** is an ultra-lightweight AI Assistant written in Go, designed to run on minimal hardware with maximum efficiency.

### Technical Specifications
- **Language**: Go 1.25+
- **RAM**: <10MB (99% smaller than alternatives)
- **Boot time**: <1s on 0.6GHz single core CPU
- **Cost**: Runs on $10 hardware
- **Architecture**: x86_64, ARM64, RISC-V

### Key Features
- 🤖 Multi-agent collaboration with role-based specialization
- 💬 IRC-style collaborative chat in Telegram
- 🔒 4-level safety system (strict, moderate, permissive, off)
- 🌐 Multi-platform chat integration (Telegram, Discord, WhatsApp, QQ, etc.)
- 🔌 Model Context Protocol (MCP) support
- 🤝 Flexible model provider support (OpenAI, Anthropic, Zhipu, etc.)

---

## 📁 Directory Structure

```
picoclaw/
├── cmd/                          # Command-line applications
│   ├── picoclaw/                # Main CLI binary
│   │   ├── main.go             # Entry point
│   │   └── internal/           # CLI-specific logic
│   │       ├── agent/          # Agent command handlers
│   │       ├── auth/           # Authentication commands
│   │       ├── cron/           # Cron job management
│   │       ├── gateway/        # Gateway server
│   │       ├── migrate/        # Config migration
│   │       ├── onboard/        # Initial setup
│   │       ├── skills/         # Skills management
│   │       ├── status/         # System status
│   │       ├── teamcmd/        # Multi-agent team commands
│   │       └── version/        # Version info
│   ├── picoclaw-launcher/      # GUI launcher (Windows)
│   └── picoclaw-launcher-tui/  # TUI launcher
│
├── pkg/                         # Reusable packages
│   ├── agent/                  # Agent core logic
│   │   ├── context.go         # Context management
│   │   ├── instance.go        # Agent instances
│   │   ├── loop.go            # Main agent loop
│   │   ├── memory.go          # Memory management
│   │   └── registry.go        # Agent registry
│   │
│   ├── auth/                   # Authentication & OAuth
│   │   ├── oauth.go           # OAuth flow
│   │   ├── pkce.go            # PKCE implementation
│   │   ├── store.go           # Token storage
│   │   └── token.go           # Token management
│   │
│   ├── bus/                    # Message bus
│   │   ├── bus.go             # Message bus implementation
│   │   └── types.go           # Message types
│   │
│   ├── channels/               # Chat platform integrations
│   │   ├── base.go            # Base channel implementation
│   │   ├── manager.go         # Channel manager
│   │   ├── telegram/          # Telegram integration
│   │   │   ├── telegram.go   # Main Telegram handler
│   │   │   └── collaborative_chat.go  # Collaborative chat
│   │   ├── discord/           # Discord integration
│   │   ├── whatsapp/          # WhatsApp integration
│   │   ├── qq/                # QQ integration
│   │   ├── dingtalk/          # DingTalk integration
│   │   ├── slack/             # Slack integration
│   │   ├── line/              # LINE integration
│   │   ├── feishu/            # Feishu integration
│   │   ├── wecom/             # WeCom integration
│   │   └── ...                # Other platforms
│   │
│   ├── collaborative/          # Platform-agnostic collaborative chat
│   │   ├── types.go           # Core types and Platform interface
│   │   ├── session.go         # Session management
│   │   ├── manager.go         # Collaborative manager
│   │   ├── mention.go         # @mention extraction
│   │   ├── formatting.go      # IRC-style formatting
│   │   └── roster.go          # Team roster building
│   │
│   ├── config/                 # Configuration management
│   │   ├── config.go          # Config loading
│   │   ├── defaults.go        # Default values
│   │   └── migration.go       # Config migration
│   │
│   ├── cron/                   # Job scheduling
│   │   └── service.go         # Cron service
│   │
│   ├── devices/                # Device integrations
│   │   ├── service.go         # Device service
│   │   ├── events/            # Device events
│   │   └── sources/           # Device sources
│   │
│   ├── health/                 # Health checks
│   │   └── server.go          # Health server
│   │
│   ├── heartbeat/              # Heartbeat service
│   │   └── service.go         # Heartbeat implementation
│   │
│   ├── identity/               # Identity management
│   │   └── identity.go        # Canonical ID system
│   │
│   ├── logger/                 # Structured logging
│   │   └── logger.go          # Logger implementation
│   │
│   ├── mcp/                    # Model Context Protocol
│   │   ├── manager.go         # MCP manager
│   │   ├── client.go          # MCP client
│   │   └── types.go           # MCP types
│   │
│   ├── media/                  # Media handling
│   │   └── media.go           # Media utilities
│   │
│   ├── providers/              # LLM provider implementations
│   │   ├── openai/            # OpenAI provider
│   │   ├── anthropic/         # Anthropic provider
│   │   ├── zhipu/             # Zhipu provider
│   │   ├── openrouter/        # OpenRouter provider
│   │   ├── gemini/            # Gemini provider
│   │   ├── groq/              # Groq provider
│   │   ├── ollama/            # Ollama provider
│   │   └── ...                # Other providers
│   │
│   ├── routing/                # Message routing
│   │   ├── router.go          # Router implementation
│   │   └── session_key.go     # Session key management
│   │
│   ├── session/                # Session management
│   │   └── manager.go         # Session manager
│   │
│   ├── skills/                 # Skills system
│   │   ├── manager.go         # Skills manager
│   │   ├── installer.go       # Skills installer
│   │   └── registry.go        # Skills registry
│   │
│   ├── state/                  # State persistence
│   │   └── store.go           # State store
│   │
│   ├── team/                   # Multi-agent teams
│   │   ├── types.go           # Team types
│   │   ├── manager.go         # Team manager
│   │   ├── coordinator.go     # Coordination patterns
│   │   ├── consensus.go       # Consensus mechanisms
│   │   ├── config.go          # Team configuration
│   │   ├── memory.go          # Team memory
│   │   ├── health.go          # Health checks
│   │   └── metrics.go         # Metrics collection
│   │
│   ├── tools/                  # Agent tools
│   │   ├── exec.go            # Shell execution
│   │   ├── filesystem.go      # File operations
│   │   ├── web.go             # Web search
│   │   ├── subagent.go        # Subagent spawning
│   │   └── ...                # Other tools
│   │
│   ├── utils/                  # Utilities
│   │   ├── download.go        # Download utilities
│   │   ├── zip.go             # ZIP utilities
│   │   ├── http_retry.go      # HTTP retry logic
│   │   └── ...                # Other utilities
│   │
│   └── voice/                  # Voice processing
│       └── transcriber.go     # Voice transcription
│
├── templates/                   # Configuration templates
│   └── teams/                  # Team configurations
│       ├── development-team.json
│       ├── research-team.json
│       ├── analysis-team.json
│       └── collaborative-dev-team.json
│
├── docs/                        # Documentation
│   ├── MULTI_AGENT_GUIDE.md
│   ├── COLLABORATIVE_CHAT.md
│   ├── SAFETY_LEVELS.md
│   ├── channels/               # Channel-specific guides
│   └── ...
│
├── config/                      # Example configurations
│   ├── config.example.json
│   └── safety_examples.json
│
├── docker/                      # Docker configurations
│   ├── Dockerfile
│   ├── Dockerfile.full
│   ├── docker-compose.yml
│   └── docker-compose.full.yml
│
└── assets/                      # Static assets
    ├── logo.jpg
    └── ...
```

---

## 🏗️ Architecture Overview

### Core Components

1. **Agent System** (`pkg/agent/`)
   - Agent lifecycle management
   - Context and memory management
   - Main execution loop
   - Agent registry

2. **Multi-Agent Teams** (`pkg/team/`)
   - Team coordination (sequential, parallel, hierarchical)
   - Consensus mechanisms
   - Team memory and state
   - Health monitoring

3. **Collaborative Chat** (`pkg/collaborative/`)
   - Platform-agnostic implementation
   - @mention routing
   - Session management
   - IRC-style formatting

4. **LLM Providers** (`pkg/providers/`)
   - Multiple provider support
   - Unified interface
   - Streaming support
   - Error handling

5. **Tools System** (`pkg/tools/`)
   - Shell execution with safety levels
   - File operations
   - Web search
   - Subagent spawning

6. **Channel Integrations** (`pkg/channels/`)
   - Telegram, Discord, WhatsApp, QQ, etc.
   - Unified channel interface
   - Message routing
   - Media handling

---

## 🔄 Data Flow

### Single Agent Flow
```
User Input → Channel → Agent → LLM Provider → Tools → Response → Channel → User
```

### Multi-Agent Flow
```
User Input → Channel → Team Manager → Coordinator → Multiple Agents → Consensus → Response
```

### Collaborative Chat Flow
```
User @mention → Telegram → Collaborative Manager → Team → Parallel Agents → Formatted Responses
```

---

## 🛠️ Development Workflow

### Building
```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Build with WhatsApp support
make build-whatsapp-native
```

### Testing
```bash
# Run all tests
make test

# Run specific package tests
go test ./pkg/agent/...

# Run with coverage
go test -cover ./...
```

### Code Quality
```bash
# Format code
make fmt

# Run linters
make lint

# Fix linting issues
make fix

# Run all checks
make check
```

---

## 📦 Key Packages

### Core Packages

- **agent**: Agent lifecycle and execution
- **team**: Multi-agent coordination
- **collaborative**: Platform-agnostic collaborative chat
- **providers**: LLM provider implementations
- **tools**: Agent tools and capabilities

### Integration Packages

- **channels**: Chat platform integrations
- **mcp**: Model Context Protocol
- **skills**: Extensibility system
- **cron**: Job scheduling

### Infrastructure Packages

- **config**: Configuration management
- **auth**: Authentication and OAuth
- **logger**: Structured logging
- **bus**: Message bus for inter-agent communication
- **routing**: Message routing and session management

---

## 🔐 Security

### Safety Levels

PicoClaw implements a 4-level safety system for shell command execution:

1. **Strict**: Maximum protection, blocks most dangerous commands
2. **Moderate**: Balanced protection (default)
3. **Permissive**: Minimal protection
4. **Off**: No safety checks (use with caution)

See [SAFETY_LEVELS.md](SAFETY_LEVELS.md) for details.

### Authentication

- OAuth 2.0 with PKCE support
- Token storage and management
- Per-channel authentication

---

## 📚 Documentation

### User Guides
- [Multi-Agent Guide](MULTI_AGENT_GUIDE.md)
- [Collaborative Chat Guide](COLLABORATIVE_CHAT.md)
- [Safety Levels](SAFETY_LEVELS.md)
- [Channel Setup Guides](channels/)

### Developer Guides
- [Contributing Guide](../CONTRIBUTING.md)
- [API Documentation](../pkg/collaborative/README.md)
- [Architecture](COLLABORATIVE_CHAT_ARCHITECTURE.md)

### Configuration
- [Configuration Examples](../config/)
- [Team Templates](../templates/teams/)

---

## 🚀 Getting Started

1. **Install**: Download binary or build from source
2. **Initialize**: Run `picoclaw onboard`
3. **Configure**: Edit `~/.picoclaw/config.json`
4. **Run**: Start with `picoclaw agent` or `picoclaw gateway`

See [README.md](../README.md) for detailed instructions.

---

## 🤝 Contributing

We welcome contributions! See [CONTRIBUTING.md](../CONTRIBUTING.md) for guidelines.

---

## 📄 License

MIT License - see [LICENSE](../LICENSE) for details.

---

## 🔗 Links

- **Website**: https://picoclaw.io
- **GitHub**: https://github.com/sipeed/picoclaw
- **Discord**: https://discord.gg/V4sAZ9XWpN
- **Twitter**: @SipeedIO
