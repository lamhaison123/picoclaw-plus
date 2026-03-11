# PicoClaw Architecture Overview

## Giới thiệu

PicoClaw là một **multi-agent AI system** được thiết kế để:
- Xử lý tin nhắn từ nhiều platform (Telegram, Discord, Slack, etc.)
- Quản lý nhiều AI agents với workspace riêng biệt
- Hỗ trợ team collaboration với nhiều agents làm việc cùng nhau
- Lưu trữ và tìm kiếm semantic memory qua vector database
- Thực thi tools và commands một cách an toàn

## Kiến trúc tổng quan

```
┌─────────────────────────────────────────────────────────────┐
│                    USER INTERFACES                          │
│  Telegram │ Discord │ Slack │ WhatsApp │ CLI │ HTTP API    │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│                   CHANNEL LAYER                             │
│  - Platform adapters                                        │
│  - Message normalization                                    │
│  - Media handling                                           │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│                   MESSAGE BUS                               │
│  - Inbound queue  (User → Agent)                           │
│  - Outbound queue (Agent → User)                           │
│  - Media queue    (Agent → User)                           │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│                   AGENT LOOP                                │
│  - Message routing                                          │
│  - Session management                                       │
│  - Context building                                         │
│  - LLM interaction                                          │
│  - Tool execution                                           │
└─────┬────────────────┬────────────────┬─────────────────────┘
      │                │                │
      ▼                ▼                ▼
┌──────────┐    ┌──────────┐    ┌──────────────┐
│  MEMORY  │    │  TOOLS   │    │  TEAM        │
│  SYSTEM  │    │  SYSTEM  │    │  MANAGER     │
└──────────┘    └──────────┘    └──────────────┘
```

## Core Components

### 1. Message Bus (`pkg/bus/`)
**Vai trò**: Hub trung tâm cho tất cả communication

**Channels**:
- `inbound`: User messages → Agent
- `outbound`: Agent responses → User  
- `outboundMedia`: Media attachments → User

**Đặc điểm**:
- Thread-safe với atomic operations
- Context cancellation support
- Buffered channels (default: 64)
- Graceful shutdown

### 2. Agent Loop (`pkg/agent/loop.go`)
**Vai trò**: Core execution engine

**Chức năng chính**:
- Consume messages từ bus
- Route messages đến đúng agent
- Load session history
- Search vector memory cho context
- Build context cho LLM
- Execute LLM với tools
- Store conversation vào memory
- Publish response về bus

**Flow**:
```
Message → Route → Load History → Search Memory → 
Build Context → LLM Call → Tool Execution → 
Store Memory → Save Session → Response
```

### 3. Agent Registry (`pkg/agent/registry.go`)
**Vai trò**: Quản lý nhiều agent instances

**Chức năng**:
- Load agent configs
- Create agent instances
- Resolve routing rules
- Provide default agent

**Agent Instance** bao gồm:
- ID và Name
- Model và Fallbacks
- Workspace (isolated)
- Session Manager
- Tool Registry
- Context Builder

### 4. Team Manager (`pkg/team/manager.go`)
**Vai trò**: Orchestrate multi-agent teams

**Chức năng**:
- Create/dissolve teams
- Delegate tasks to roles
- Coordinate execution
- Manage shared context
- Track metrics

**Team Structure**:
```json
{
  "team_id": "dev-team",
  "pattern": "hierarchical",
  "roles": [
    {"name": "architect", "capabilities": ["design"]},
    {"name": "developer", "capabilities": ["coding"]},
    {"name": "tester", "capabilities": ["testing"]}
  ]
}
```

### 5. Collaborative Chat (`pkg/collaborative/`)
**Vai trò**: Multi-agent conversations với mentions

**Chức năng**:
- Detect @mentions
- Cascade mentions (A→B→C)
- Prevent infinite loops
- Rate limiting
- Context compaction

**Flow**:
```
User: "@architect design API"
  ↓
Architect: "I'll design... @developer implement"
  ↓
Developer: "Implementing... @tester verify"
  ↓
Tester: "Testing complete"
```

### 6. Memory System

#### Session Memory (`pkg/session/`)
- File-based JSON storage
- Per-session conversation history
- Summary tracking
- Thread-safe operations

#### Vector Memory (`pkg/memory/vector/`)
- Semantic search qua embeddings
- Qdrant/LanceDB support (với LanceDB chạy native qua CGO)
- Circuit breaker cho fault tolerance
- Async storage

**Flow**:
```
User Query → Generate Embedding → Search Similar → 
Inject Context → LLM Response → Store Embedding
```

### 7. Tools System (`pkg/tools/`)
**Vai trò**: Provide capabilities cho agents

**Built-in Tools**:
- File operations (read, write, edit, append, list)
- Shell execution (với restrictions)
- Web search (Brave, Tavily, DuckDuckGo)
- Web fetch
- Team delegation
- Subagent spawning
- MCP tools

**Tool Execution**:
```
LLM requests tool → Validate permissions → 
Execute → Return result → Feed back to LLM
```

### 8. Channel System (`pkg/channels/`)
**Vai trò**: Platform integrations

**Supported Platforms**:
- Telegram
- Discord
- Slack
- WhatsApp
- Matrix
- IRC
- CLI
- HTTP API

**Channel Interface**:
```go
type Channel interface {
    Start(ctx) error
    Stop() error
    Send(message) error
}
```

## Data Flow

### Message Processing Flow
```
1. User sends message via Telegram
2. Telegram channel receives message
3. Channel publishes to MessageBus.inbound
4. AgentLoop consumes from inbound
5. Route to appropriate agent
6. Load session history
7. Search vector memory
8. Build context with history + memory
9. Call LLM with tools
10. Execute tool calls (if any)
11. Get final response
12. Store in vector memory (async)
13. Save to session
14. Publish to MessageBus.outbound
15. Channel sends to Telegram
16. User receives response
```

### Team Delegation Flow
```
1. Agent uses delegate_to_team tool
2. TeamManager.ExecuteTask(teamID, task)
3. Create Coordinator for team
4. Coordinator analyzes task
5. Route to appropriate role
6. Execute with role's agent
7. Collect result
8. Return to original agent
9. Agent incorporates result
```

### Mention Cascading Flow
```
1. User: "@architect design API"
2. ManagerV2.HandleMentions(["architect"])
3. Architect processes and responds
4. Architect: "@developer implement"
5. ManagerV2 detects new mention
6. Check depth limit (max: 20)
7. Developer processes and responds
8. Developer: "@tester verify"
9. Tester processes and responds
10. All responses sent to user
```

## Design Patterns

### 1. Registry Pattern
- AgentRegistry: Manages agents
- ToolRegistry: Manages tools
- ChannelRegistry: Manages channels

### 2. Factory Pattern
- NewAgentInstance()
- NewTeamManager()
- initializeVectorStore()

### 3. Strategy Pattern
- LLMProvider interface
- VectorStore interface
- EmbeddingService interface

### 4. Observer/Pub-Sub
- MessageBus as event hub
- Channels subscribe/publish

### 5. Chain of Responsibility
- FallbackChain for LLM providers
- Tool execution loop

### 6. Adapter Pattern
- TeamManagerAdapter
- Channel adapters

### 7. Circuit Breaker
- Wraps vector store calls
- Prevents cascading failures

### 8. State Machine
- Team status transitions
- Agent status transitions
- Circuit breaker states

## Configuration

### Hierarchy
```
config.json
├─ agents.defaults (global defaults)
├─ agents.list[] (per-agent overrides)
├─ channels[] (platform configs)
├─ providers (LLM credentials)
├─ memory (vector DB config)
├─ tools (restrictions)
└─ heartbeat (periodic tasks)
```

### Feature Flags
- `memory.enabled`: Vector memory on/off
- `channels.*.enabled`: Per-channel on/off
- `tools.restrict_to_workspace`: File restrictions
- `heartbeat.enabled`: Periodic tasks

## Performance

### Latency
- Message processing: ~100-500ms
- Tool execution: ~10-100ms
- Vector search: ~10-50ms
- Embedding: ~100-500ms
- Total with memory: ~200-1000ms

### Throughput
- Single agent: ~10-20 msg/s
- Multiple agents: ~50-100 msg/s
- Team execution: ~5-10 tasks/s

### Memory Usage
- Base agent: ~50-100MB
- Per session: ~1-5MB
- Vector store: ~100MB-1GB

## Security

### Tool Restrictions
- File ops restricted to workspace
- Shell execution whitelist
- Path traversal prevention
- Configurable allow/deny lists

### Team Isolation
- Shared context with access control
- Role-based capabilities
- Task validation

### Channel Security
- Platform authentication
- Message validation
- Rate limiting
- Mention depth limiting

## Extensibility

### Easy to Extend
1. New channels: Implement Channel interface
2. New tools: Register in ToolRegistry
3. New LLM providers: Implement LLMProvider
4. New vector stores: Implement VectorStore
5. New team roles: Define in config

### Extension Points
- Tool system: Add custom tools
- Channel system: Add new platforms
- Memory system: Add new backends
- Team patterns: Add new collaboration modes

## Next Steps

Đọc thêm:
- [Component Details](./COMPONENT_DETAILS.md)
- [Data Flow](./DATA_FLOW.md)
- [Configuration Guide](./CONFIGURATION_GUIDE.md)
- [Development Guide](./DEVELOPMENT_GUIDE.md)
