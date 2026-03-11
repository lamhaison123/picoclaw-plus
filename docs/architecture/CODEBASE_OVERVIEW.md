# PicoClaw Codebase Overview

## Introduction

This document provides a comprehensive overview of the PicoClaw codebase, helping developers understand the project structure, key components, and how everything fits together.

## Project Statistics

- **Language**: Go 1.25.5
- **Total Packages**: 30+
- **Lines of Code**: ~50,000+
- **Test Coverage**: 180+ unit tests, 9 integration tests, 9 benchmarks
- **Dependencies**: 40+ external packages
- **Supported Platforms**: Linux (x86_64, ARM64, RISC-V, LoongArch), macOS, Windows

## Directory Structure

```
picoclaw/
├── cmd/                          # Command-line applications
│   ├── picoclaw/                 # Main CLI application
│   │   ├── main.go              # Entry point
│   │   └── internal/            # CLI command implementations
│   │       ├── agent/           # Agent command (chat with AI)
│   │       ├── gateway/         # Gateway command (messaging platforms)
│   │       ├── teamcmd/         # Team command (multi-agent teams)
│   │       ├── cron/            # Cron command (scheduled tasks)
│   │       ├── status/          # Status command (system info)
│   │       ├── auth/            # Auth command (OAuth management)
│   │       ├── onboard/         # Onboard command (initialization)
│   │       ├── migrate/         # Migrate command (config migration)
│   │       ├── skills/          # Skills command (skill management)
│   │       └── version/         # Version command
│   ├── picoclaw-launcher/       # GUI launcher (Fyne-based)
│   └── picoclaw-launcher-tui/   # TUI launcher (tview-based)
│
├── pkg/                          # Reusable packages (public API)
│   ├── agent/                   # Core agent system
│   │   ├── instance.go          # Agent instance
│   │   ├── loop.go              # Agent execution loop
│   │   ├── context.go           # Context management
│   │   ├── memory.go            # Memory management
│   │   └── registry.go          # Agent registry
│   │
│   ├── team/                    # Multi-agent framework
│   │   ├── manager.go           # Team manager
│   │   ├── team.go              # Team structure
│   │   ├── coordinator.go       # Coordination patterns
│   │   ├── consensus.go         # Consensus mechanisms
│   │   ├── memory.go            # Team memory
│   │   └── metrics.go           # Performance metrics
│   │
│   ├── collaborative/           # Collaborative chat system
│   │   ├── types.go             # Core types and interfaces
│   │   ├── manager.go           # Session and mention management
│   │   ├── session.go           # Session context
│   │   ├── mention.go           # @mention extraction
│   │   ├── formatting.go        # IRC-style formatting
│   │   └── roster.go            # Team roster building
│   │
│   ├── providers/               # LLM provider implementations
│   │   ├── factory.go           # Provider factory
│   │   ├── openai.go            # OpenAI provider
│   │   ├── anthropic.go         # Anthropic provider
│   │   ├── zhipu.go             # Zhipu provider
│   │   ├── ollama.go            # Ollama provider
│   │   ├── gemini.go            # Gemini provider
│   │   └── ...                  # Other providers
│   │
│   ├── channels/                # Messaging platform integrations
│   │   ├── manager.go           # Channel manager
│   │   ├── interfaces.go        # Channel interfaces
│   │   ├── telegram/            # Telegram integration
│   │   ├── discord/             # Discord integration
│   │   ├── whatsapp/            # WhatsApp integration
│   │   ├── qq/                  # QQ integration
│   │   ├── dingtalk/            # DingTalk integration
│   │   ├── slack/               # Slack integration
│   │   └── ...                  # Other channels
│   │
│   ├── tools/                   # Tool system
│   │   ├── registry.go          # Tool registry
│   │   ├── interfaces.go        # Tool interfaces
│   │   ├── file_read.go         # File reading tool
│   │   ├── file_write.go        # File writing tool
│   │   ├── exec.go              # Shell execution tool
│   │   ├── web_search.go        # Web search tool
│   │   ├── mcp.go               # MCP tool integration
│   │   └── ...                  # Other tools
│   │
│   ├── bus/                     # Message bus (event system)
│   │   ├── bus.go               # Message bus implementation
│   │   └── types.go             # Message types
│   │
│   ├── config/                  # Configuration system
│   │   ├── config.go            # Main configuration
│   │   ├── agent.go             # Agent configuration
│   │   ├── provider.go          # Provider configuration
│   │   ├── channel.go           # Channel configuration
│   │   └── tools.go             # Tool configuration
│   │
│   ├── logger/                  # Logging system
│   │   └── logger.go            # Structured logging
│   │
│   ├── session/                 # Session management
│   │   └── manager.go           # Session manager
│   │
│   ├── auth/                    # Authentication
│   │   ├── oauth.go             # OAuth 2.0 implementation
│   │   ├── pkce.go              # PKCE flow
│   │   ├── store.go             # Token storage
│   │   └── token.go             # Token management
│   │
│   ├── fileutil/                # File utilities
│   ├── httputil/                # HTTP utilities
│   └── ...                      # Other utilities
│
├── internal/                     # Internal packages (private)
│   ├── memory/                  # Memory pool
│   ├── retry/                   # Retry policies
│   └── ratelimit/               # Rate limiting
│
├── docs/                         # Documentation
│   ├── ARCHITECTURE.md          # Architecture overview
│   ├── DEVELOPER_GUIDE.md       # Developer guide
│   ├── API_REFERENCE.md         # API reference
│   ├── MULTI_AGENT_GUIDE.md     # Multi-agent guide
│   ├── COLLABORATIVE_CHAT.md    # Collaborative chat guide
│   ├── SAFETY_LEVELS.md         # Safety system guide
│   └── ...                      # Other documentation
│
├── templates/                    # Configuration templates
│   └── teams/                   # Team configuration templates
│       ├── development-team.json
│       ├── research-team.json
│       └── ...
│
├── skills/                       # Built-in skills
│   └── ...                      # Skill implementations
│
├── config/                       # Example configurations
│   ├── config.example.json      # Example config
│   ├── collaborative-chat-*.json # Collaborative chat examples
│   └── ...
│
├── docker/                       # Docker files
│   ├── Dockerfile               # Minimal Alpine image
│   ├── Dockerfile.full          # Full-featured image
│   ├── docker-compose.yml       # Docker Compose config
│   └── entrypoint.sh            # Container entrypoint
│
├── examples/                     # Example code
│   └── teams/                   # Team usage examples
│
├── build/                        # Build artifacts (generated)
├── Makefile                      # Build automation
├── go.mod                        # Go module definition
├── go.sum                        # Go module checksums
├── README.md                     # Project README
├── CHANGELOG.md                  # Version history
├── CONTRIBUTING.md               # Contribution guide
└── LICENSE                       # MIT License
```

## Key Components

### 1. Agent System (`pkg/agent/`)

**Purpose**: Core agent execution system

**Key Files**:
- `instance.go`: Agent instance with configuration and tools
- `loop.go`: Main execution loop for agent-LLM interaction
- `context.go`: Context building and caching
- `memory.go`: Conversation memory management
- `registry.go`: Agent registration and lookup

**Responsibilities**:
- Execute tasks with LLM
- Manage conversation context
- Execute tools
- Handle memory and summarization
- Cache context for performance

### 2. Multi-Agent Framework (`pkg/team/`)

**Purpose**: Coordinate multiple agents as teams

**Key Files**:
- `manager.go`: Team orchestration and lifecycle
- `team.go`: Team structure and configuration
- `coordinator.go`: Collaboration patterns (sequential, parallel, hierarchical)
- `consensus.go`: Voting and consensus mechanisms
- `memory.go`: Persistent team memory

**Responsibilities**:
- Create and manage teams
- Delegate tasks to roles
- Coordinate agent execution
- Implement collaboration patterns
- Track team performance

### 3. Collaborative Chat (`pkg/collaborative/`)

**Purpose**: Platform-agnostic IRC-style multi-agent chat

**Key Files**:
- `manager.go`: Session and mention management
- `session.go`: Conversation context per chat
- `mention.go`: @mention extraction
- `formatting.go`: IRC-style message formatting
- `roster.go`: Team roster building

**Responsibilities**:
- Extract @mentions from messages
- Route mentions to appropriate agents
- Manage conversation sessions
- Format responses with IRC-style prefixes
- Enable agent-to-agent communication

### 4. Provider System (`pkg/providers/`)

**Purpose**: LLM provider abstraction and management

**Key Files**:
- `factory.go`: Provider selection and creation
- `openai.go`, `anthropic.go`, etc.: Provider implementations
- `interfaces.go`: Provider interface definitions

**Responsibilities**:
- Abstract LLM API differences
- Handle authentication
- Implement retry and circuit breaker
- Support fallback chains
- Manage rate limiting

### 5. Channel Manager (`pkg/channels/`)

**Purpose**: Messaging platform integrations

**Key Files**:
- `manager.go`: Channel lifecycle management
- `interfaces.go`: Channel interface definitions
- `telegram/`, `discord/`, etc.: Platform implementations

**Responsibilities**:
- Connect to messaging platforms
- Handle inbound/outbound messages
- Support typing indicators, reactions, etc.
- Manage platform-specific features

### 6. Tool System (`pkg/tools/`)

**Purpose**: Extensible tool registry and execution

**Key Files**:
- `registry.go`: Tool registration and execution
- `interfaces.go`: Tool interface definitions
- `file_*.go`, `exec.go`, etc.: Tool implementations

**Responsibilities**:
- Register and manage tools
- Execute tools with safety checks
- Support async tools
- Handle tool permissions

### 7. Configuration System (`pkg/config/`)

**Purpose**: Comprehensive configuration management

**Key Files**:
- `config.go`: Main configuration structure
- `agent.go`, `provider.go`, etc.: Component configs

**Responsibilities**:
- Load/save configuration
- Environment variable support
- Validation
- Default values

### 8. Message Bus (`pkg/bus/`)

**Purpose**: Event-driven inter-component communication

**Key Files**:
- `bus.go`: Message bus implementation
- `types.go`: Message type definitions

**Responsibilities**:
- Route inbound/outbound messages
- Support media messages
- Handle concurrent access
- Graceful shutdown

## Data Flow

### CLI Chat Flow

```
User Input
    ↓
CLI Command Parser
    ↓
Agent Instance
    ↓
Agent Loop
    ↓
Context Builder → Build Prompt
    ↓
Provider → LLM API Call
    ↓
Response Parser
    ↓
Tool Execution (if needed)
    ↓
Loop Continue (if tool calls)
    ↓
Final Response
    ↓
Display to User
```

### Gateway Message Flow

```
Messaging Platform (Telegram/Discord/etc.)
    ↓
Channel Handler
    ↓
Message Bus (Inbound)
    ↓
Gateway Loop
    ↓
Mention Detection
    ↓
Collaborative Manager
    ↓
Team Manager
    ↓
Agent Execution (parallel)
    ↓
Response Formatting
    ↓
Message Bus (Outbound)
    ↓
Channel Handler
    ↓
Messaging Platform
```

### Team Execution Flow

```
User Task
    ↓
Team Manager
    ↓
Coordinator
    ↓
Pattern Selection (Sequential/Parallel/Hierarchical)
    ↓
Role Assignment
    ↓
Agent Execution
    ↓
Result Aggregation
    ↓
Consensus (if needed)
    ↓
Final Result
```

## Key Design Patterns

### 1. Factory Pattern

Used in provider system for creating LLM providers:

```go
func NewProvider(cfg *config.Config) (LLMProvider, error) {
    // Select provider based on configuration
    // Return appropriate implementation
}
```

### 2. Registry Pattern

Used in agent and tool systems:

```go
type Registry struct {
    items map[string]Item
    mu    sync.RWMutex
}

func (r *Registry) Register(item Item)
func (r *Registry) Get(name string) (Item, bool)
```

### 3. Strategy Pattern

Used in collaboration patterns:

```go
type Coordinator interface {
    Execute(ctx context.Context, task string) (any, error)
}

type SequentialCoordinator struct { /* ... */ }
type ParallelCoordinator struct { /* ... */ }
type HierarchicalCoordinator struct { /* ... */ }
```

### 4. Observer Pattern

Used in message bus:

```go
type MessageBus struct {
    inbound  chan InboundMessage
    outbound chan OutboundMessage
}

func (mb *MessageBus) PublishInbound(msg InboundMessage)
func (mb *MessageBus) ConsumeInbound() (InboundMessage, bool)
```

### 5. Adapter Pattern

Used in channel integrations:

```go
type Channel interface {
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
}

type TelegramChannel struct { /* ... */ }
type DiscordChannel struct { /* ... */ }
```

## Concurrency Model

### Goroutines

- Each channel runs in its own goroutine
- Agent execution can be parallel (in teams)
- Message bus uses goroutines for routing
- Tool execution can be async

### Synchronization

- `sync.RWMutex` for read-heavy registries
- `sync.Mutex` for exclusive access
- `atomic` for lock-free counters
- Channels for communication

### Context Management

- `context.Context` for cancellation
- Timeout support
- Graceful shutdown

## Testing Strategy

### Unit Tests

- Package-level tests (`*_test.go`)
- Table-driven tests
- Mock providers and tools
- 180+ tests with 100% pass rate

### Integration Tests

- End-to-end workflows
- Multi-component interaction
- 9 integration tests

### Benchmarks

- Performance measurements
- Memory profiling
- 9 benchmarks

### Property-Based Testing

- Gopter for randomized testing
- Edge case discovery

## Build System

### Makefile Targets

```bash
make build          # Build for current platform
make build-all      # Build for all platforms
make test           # Run tests
make test-coverage  # Run tests with coverage
make lint           # Run linter
make install        # Install to ~/.local/bin
make clean          # Clean build artifacts
```

### Cross-Compilation

Supports multiple platforms:
- Linux: x86_64, ARM64, RISC-V, LoongArch
- macOS: x86_64, ARM64
- Windows: x86_64

### Docker

- Minimal Alpine-based image (<50MB)
- Full-featured image with Node.js
- Multi-stage builds for optimization

## Dependencies

### Core Dependencies

- **Cobra**: CLI framework
- **Anthropic SDK**: Claude API
- **OpenAI SDK**: GPT API
- **Telego**: Telegram bot API
- **DiscordGo**: Discord bot API
- **SQLite**: Session storage
- **UUID**: Unique identifiers

### Optional Dependencies

- **Fyne**: GUI launcher
- **tview**: TUI launcher
- **QR Terminal**: QR code display
- **MCP SDK**: Model Context Protocol

## Performance Characteristics

### Memory Usage

- Baseline: <10MB
- Per agent: ~1-2MB
- Context cache: 1-5MB
- Total: 10-20MB typical

### Response Time

- Local models: 1-5s
- Cloud APIs: 2-10s
- With cache: 20-50% faster

### Concurrency

- Thousands of goroutines
- Buffered channels (64 default)
- Parallel tool execution

## Security Considerations

### Command Execution

- 4-level safety system
- Pattern-based filtering
- Workspace restrictions

### Authentication

- OAuth 2.0 with PKCE
- Token storage
- API key management

### Data Protection

- Workspace isolation
- Path validation
- Sensitive data filtering

## Future Enhancements

### Planned Features

1. Streaming responses
2. Advanced caching
3. More collaboration patterns
4. Additional integrations
5. Enterprise features

### Technical Debt

1. Improve test coverage
2. Add more benchmarks
3. Optimize memory usage
4. Enhance documentation

## Resources

- [Architecture](ARCHITECTURE.md)
- [Developer Guide](DEVELOPER_GUIDE.md)
- [API Reference](API_REFERENCE.md)
- [Contributing](../CONTRIBUTING.md)

---

**Last Updated**: 2026-03-07  
**Version**: 1.1.1
