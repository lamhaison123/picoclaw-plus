# PicoClaw Repository Overview

## 📋 Tổng Quan Dự Án

**PicoClaw** là một AI Assistant siêu nhẹ được viết bằng Go, thiết kế để chạy trên phần cứng tối thiểu với hiệu suất tối đa.

### Thông Số Kỹ Thuật
- **Ngôn ngữ**: Go 1.25+
- **RAM**: <10MB (99% nhỏ hơn các giải pháp khác)
- **Boot time**: <1s trên CPU 0.6GHz single core
- **Chi phí**: Chạy trên phần cứng $10
- **Kiến trúc**: x86_64, ARM64, RISC-V

### Tính Năng Chính
- 🤖 Multi-agent collaboration với role-based specialization
- 💬 IRC-style collaborative chat trong Telegram
- 🔒 4-level safety system (strict, moderate, permissive, off)
- 🌐 Multi-platform chat integration (Telegram, Discord, WhatsApp, QQ, etc.)
- 🔌 Model Context Protocol (MCP) support
- 🤝 Flexible model provider support (OpenAI, Anthropic, Zhipu, etc.)

---

## 📁 Cấu Trúc Thư Mục

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
│   │   └── manager.go         # MCP manager
│   │
│   ├── media/                  # Media handling
│   │   └── store.go           # Media storage
│   │
│   ├── providers/              # LLM provider implementations
│   │   ├── factory.go         # Provider factory
│   │   ├── claude_provider.go # Anthropic Claude
│   │   ├── codex_provider.go  # GitHub Copilot
│   │   ├── antigravity_provider.go  # Antigravity
│   │   ├── openai_compat/     # OpenAI-compatible providers
│   │   └── ...                # Other providers
│   │
│   ├── routing/                # Message routing
│   │   ├── route.go           # Route management
│   │   ├── agent_id.go        # Agent ID routing
│   │   └── session_key.go     # Session key routing
│   │
│   ├── session/                # Session management
│   │   └── manager.go         # Session manager
│   │
│   ├── skills/                 # Skills system
│   │   ├── loader.go          # Skill loader
│   │   ├── registry.go        # Skill registry
│   │   └── installer.go       # Skill installer
│   │
│   ├── state/                  # State persistence
│   │   └── state.go           # State manager
│   │
│   ├── team/                   # Multi-agent collaboration
│   │   ├── manager.go         # Team manager
│   │   ├── coordinator.go     # Coordination patterns
│   │   ├── consensus.go       # Consensus mechanisms
│   │   ├── config.go          # Team configuration
│   │   ├── executor.go        # Agent executor
│   │   ├── router.go          # Task routing
│   │   ├── intelligent_router.go  # Smart routing
│   │   ├── memory/            # Team memory
│   │   └── sessions/          # Team sessions
│   │
│   ├── tools/                  # Agent tools
│   │   ├── base.go            # Base tool interface
│   │   ├── filesystem.go      # File operations
│   │   ├── shell.go           # Shell execution
│   │   ├── web.go             # Web search
│   │   ├── edit.go            # Code editing
│   │   ├── message.go         # Messaging
│   │   ├── spawn.go           # Subagent spawning
│   │   ├── cron.go            # Cron management
│   │   ├── mcp_tool.go        # MCP tools
│   │   ├── skills_*.go        # Skills tools
│   │   └── ...                # Other tools
│   │
│   ├── utils/                  # Utilities
│   │   ├── download.go        # File download
│   │   ├── http_retry.go      # HTTP retry logic
│   │   ├── string.go          # String utilities
│   │   └── ...                # Other utilities
│   │
│   └── voice/                  # Voice processing
│       └── transcriber.go     # Voice transcription
│
├── templates/                   # Configuration templates
│   └── teams/                  # Team configurations
│       ├── development-team.json
│       ├── collaborative-dev-team.json
│       ├── research-team.json
│       └── ...
│
├── config/                      # Example configurations
│   ├── config.example.json
│   ├── collaborative-chat-*.json
│   └── safety_examples.json
│
├── docs/                        # Documentation
│   ├── MULTI_AGENT_GUIDE.md
│   ├── COLLABORATIVE_CHAT.md
│   ├── COLLABORATIVE_CHAT_QUICKSTART.md
│   ├── COLLABORATIVE_CHAT_ARCHITECTURE.md
│   ├── SAFETY_LEVELS.md
│   ├── TEAM_TOOL_ACCESS.md
│   ├── channels/               # Platform-specific guides
│   └── ...
│
├── docker/                      # Docker configurations
│   ├── Dockerfile              # Alpine-based minimal
│   ├── Dockerfile.full         # Node.js with MCP
│   └── docker-compose.yml
│
├── examples/                    # Code examples
│   └── teams/                  # Team examples
│
├── .kiro/                       # Kiro IDE workspace
│   ├── specs/                  # Feature specifications
│   └── steering/               # AI assistant rules
│
├── Makefile                     # Build automation
├── go.mod                       # Go dependencies
├── go.sum                       # Dependency checksums
└── README.md                    # Main documentation
```

---

## 🏗️ Kiến Trúc Hệ Thống

### 1. Agent System (`pkg/agent/`)

**Chức năng chính:**
- Quản lý vòng đời agent (lifecycle)
- Context management với caching
- Memory persistence
- Agent registry và instance management

**Luồng hoạt động:**
```
User Input → AgentLoop → Provider (LLM) → Tool Execution → Response
```

### 2. Multi-Agent Team System (`pkg/team/`)

**Coordination Patterns:**
- **Sequential**: Tasks thực thi tuần tự (design → implement → test)
- **Parallel**: Tasks chạy song song để tăng tốc
- **Hierarchical**: Phân rã tasks phức tạp động

**Components:**
- `TeamManager`: Quản lý teams và coordination
- `Coordinator`: Điều phối execution patterns
- `Consensus`: Cơ chế voting (majority/unanimous/weighted)
- `Router`: Task routing và delegation
- `Memory`: Team memory persistence

### 3. Collaborative Chat System (`pkg/channels/telegram/`)

**Architecture:**
```
User Message with @mentions
    ↓
TelegramChannel.handleMessage()
    ↓
Extract mentions (@architect, @developer)
    ↓
handleCollaborativeMessage()
    ↓
CollaborativeChatSession (context management)
    ↓
TeamManager.ExecuteTaskWithRole() (parallel)
    ↓
Format with IRC-style prefix
    ↓
Send to Telegram
```

**Key Features:**
- @mention-based routing
- Parallel agent execution
- Shared conversation context
- IRC-style formatting: `[session-id] emoji ROLE: message`
- Session management per chat

### 4. Channel System (`pkg/channels/`)

**Supported Platforms:**
- Telegram (với collaborative chat)
- Discord
- WhatsApp (native & bridge)
- QQ
- DingTalk
- Slack
- LINE
- Feishu
- WeCom (3 variants)
- OneBot
- Pico

**Architecture:**
- `Manager`: Quản lý tất cả channels
- `BaseChannel`: Base implementation cho tất cả channels
- Channel-specific implementations
- Message bus integration
- Media handling
- Webhook support

### 5. Provider System (`pkg/providers/`)

**Supported Providers:**
- OpenAI (GPT models)
- Anthropic (Claude models)
- Zhipu (GLM models - Chinese)
- OpenRouter (all models)
- Gemini (Google models)
- Groq (fast inference)
- Ollama (local models)
- GitHub Copilot
- Antigravity

**Features:**
- Factory pattern cho provider creation
- Fallback mechanism
- Error classification
- Cooldown management
- Tool call extraction

### 6. Tools System (`pkg/tools/`)

**Available Tools:**
- **Filesystem**: read, write, edit, delete files
- **Shell**: execute commands với safety levels
- **Web**: search, fetch content
- **Edit**: code editing với AST
- **Message**: send messages to channels
- **Spawn**: create subagents
- **Cron**: schedule jobs
- **MCP**: Model Context Protocol tools
- **Skills**: install, search, manage skills
- **I2C/SPI**: hardware interfaces (Linux only)

**Safety System:**
- 4 levels: strict, moderate, permissive, off
- Pattern-based blocking
- Custom allow/block patterns
- Per-tool configuration

### 7. Message Bus (`pkg/bus/`)

**Message Types:**
- `InboundMessage`: Messages từ users
- `OutboundMessage`: Messages tới users
- `OutboundMediaMessage`: Media messages
- `SystemMessage`: Internal system messages

**Features:**
- Pub/sub pattern
- Channel-agnostic routing
- Media handling
- Graceful shutdown

---

## 🔧 Build System

### Common Commands

```bash
# Build for current platform
make build

# Run all checks (deps + fmt + vet + test)
make check

# Generate code (embeds workspace files)
make generate

# Format code
make fmt

# Run linters
make lint

# Fix linting issues
make fix

# Run tests
make test

# Build for all platforms
make build-all

# Build with WhatsApp native support
make build-whatsapp-native

# Install to ~/.local/bin
make install

# Uninstall
make uninstall-all
```

### Cross-Platform Builds

```bash
# Specific platforms
make build-linux-arm      # ARMv7 (32-bit)
make build-linux-arm64    # ARM64
make build-pi-zero        # Both ARM variants

# Docker builds
make docker-build         # Alpine-based minimal
make docker-build-full    # Node.js 24 with MCP support
```

---

## 📝 Configuration

### Config Location
- User config: `~/.picoclaw/config.json`
- Workspace: `~/.picoclaw/workspace/`
- Skills: `~/.picoclaw/workspace/skills/`
- Teams: `~/.picoclaw/teams/`
- Logs: `~/.picoclaw/logs/`

### Config Structure

```json
{
  "provider": {
    "name": "anthropic",
    "api_key": "YOUR_API_KEY"
  },
  "agents": {
    "defaults": {
      "model_name": "claude-3-5-sonnet-20241022",
      "max_tokens": 4096,
      "temperature": 0.7
    }
  },
  "channels": {
    "telegram": {
      "enabled": true,
      "token": "YOUR_BOT_TOKEN",
      "collaborative_chat": {
        "enabled": true,
        "default_team_id": "dev-team",
        "max_context_length": 50
      }
    }
  },
  "tools": {
    "exec": {
      "safety_level": "moderate"
    }
  },
  "teams": [
    {
      "team_id": "dev-team",
      "name": "Development Team",
      "pattern": "parallel",
      "roles": [...]
    }
  ]
}
```

---

## 🧪 Testing

### Test Organization
- Unit tests: `*_test.go` files alongside implementation
- Integration tests: `*_integration_test.go`
- Benchmarks: `*_bench_test.go` or in `*_test.go`
- Property-based tests: Using gopter

### Running Tests

```bash
# All tests
make test

# Specific package
go test ./pkg/agent/...

# With coverage
go test -cover ./...

# Benchmarks
go test -bench=. ./pkg/team/
```

---

## 📊 Performance Targets

- **Memory**: <10MB RAM for core functionality
- **Boot time**: <1s on 0.6GHz single core
- **Binary size**: Minimize through static compilation
- **Startup**: Fast initialization, lazy loading

**Note**: Recent PRs may temporarily increase footprint to 10-20MB

---

## 🔐 Security

### Safety Levels

| Level | Use Case | Blocks | Allows |
|-------|----------|--------|--------|
| **strict** | Production | sudo, chmod, docker, package install | Read, build, test, safe git |
| **moderate** | Development (default) | Catastrophic ops only | Most dev operations |
| **permissive** | DevOps/Admin | Only catastrophic ops | Almost everything |
| **off** | Testing ⚠️ | Nothing | Everything (DANGEROUS!) |

### Access Control
- `allow_from` lists per channel
- Tool permissions per role
- Canonical ID system for identity
- OAuth support for providers

---

## 📚 Key Documentation

### User Guides
- [Multi-Agent Guide](MULTI_AGENT_GUIDE.md) - Team collaboration
- [Collaborative Chat](COLLABORATIVE_CHAT.md) - IRC-style chat
- [Quick Start](COLLABORATIVE_CHAT_QUICKSTART.md) - 5-minute setup
- [Safety Levels](SAFETY_LEVELS.md) - Security configuration

### Developer Guides
- [Contributing](../CONTRIBUTING.md) - Contribution guidelines
- [Architecture](COLLABORATIVE_CHAT_ARCHITECTURE.md) - Technical deep dive
- [Tool Access](TEAM_TOOL_ACCESS.md) - Permission system
- [Model Selection](MULTI_AGENT_MODEL_SELECTION.md) - Per-role models

### Platform Guides
- [Telegram](channels/telegram/README.zh.md)
- [Discord](channels/discord/README.zh.md)
- [WhatsApp](channels/whatsapp/README.zh.md)
- [QQ](channels/qq/README.zh.md)
- And more...

---

## 🚀 Recent Features

### v1.1.0 - Collaborative Chat (Latest)
- ✅ Native IRC-style collaborative chat in Telegram
- ✅ @mention-based routing (@architect, @developer, @tester)
- ✅ Parallel agent execution with shared context
- ✅ Session management per chat
- ✅ IRC-style formatting with emojis
- ✅ TeamManager integration into channel system

### v1.0.0 - Multi-Agent System
- ✅ Role-based agent specialization
- ✅ Three coordination patterns (sequential, parallel, hierarchical)
- ✅ Consensus mechanisms (majority, unanimous, weighted)
- ✅ Dynamic agent composition
- ✅ Team memory persistence
- ✅ Comprehensive monitoring

---

## 🛣️ Roadmap

### Short-term
- [ ] Session persistence for collaborative chat
- [ ] Auto-join rules based on keywords
- [ ] Agent-to-agent communication
- [ ] More team templates
- [ ] Performance optimizations

### Long-term
- [ ] Web UI for team management
- [ ] Advanced routing strategies
- [ ] Multi-team coordination
- [ ] Plugin system
- [ ] Cloud deployment options

---

## 🤝 Contributing

Contributions are welcome! See [CONTRIBUTING.md](../CONTRIBUTING.md) for guidelines.

**Key Points:**
- AI-assisted contributions are welcome (with disclosure)
- Run `make check` before submitting
- Keep PRs focused and small
- Fill out PR template completely

---

## 📄 License

MIT License - see [LICENSE](../LICENSE) for details

---

**Generated**: 2026-03-05
**Version**: v1.1.0
**Maintainer**: PicoClaw Community
