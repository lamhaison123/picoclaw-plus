# PicoClaw Component Details

## 1. Agent Loop (`pkg/agent/loop.go`)

### Struct Definition
```go
type AgentLoop struct {
    bus              *bus.MessageBus
    cfg              *config.Config
    registry         *AgentRegistry
    state            *state.Manager
    running          atomic.Bool
    summarizing      sync.Map
    fallback         *providers.FallbackChain
    channelManager   *channels.Manager
    mediaStore       media.MediaStore
    reasoningSem     chan struct{}
    vectorStore      memory.VectorStore
    embeddingService embedding.Service
    memoryEnabled    atomic.Bool
    teamManager      tools.TeamManager
}
```

### Key Methods

#### `Run(ctx context.Context) error`
Main event loop:
1. Initialize MCP servers
2. Loop while running:
   - Consume inbound message
   - Process message
   - Publish outbound response
3. Cleanup on shutdown

#### `processMessage(ctx, msg) (string, error)`
Message processing pipeline:
1. Route to appropriate agent
2. Load session history
3. Search vector memory
4. Build context
5. Run LLM iteration loop
6. Store in memory (async)
7. Save session
8. Return response

#### `runLLMIteration(ctx, agent, messages, opts)`
LLM execution with tool loop:
1. Call LLM with messages and tools
2. Parse tool calls from response
3. Execute each tool
4. Collect results
5. Feed back to LLM
6. Repeat until no more tool calls
7. Return final response

### Vector Memory Integration

#### Search Before LLM
```go
if al.memoryEnabled.Load() && opts.UserMessage != "" {
    memoryContext = al.searchVectorMemory(ctx, query, sessionKey)
    al.injectMemoryContext(messages, memoryContext)
}
```

#### Store After Response
```go
if al.memoryEnabled.Load() && finalContent != "" {
    go al.storeInVectorMemory(ctx, sessionKey, 
        userMsg, assistantMsg, channel)
}
```

### Team Integration

#### Set Team Manager
```go
func (al *AgentLoop) SetTeamManager(tm tools.TeamManager) {
    al.teamManager = tm
    registerSharedTools(al.cfg, al.bus, al.registry, 
        al.registry.GetDefaultAgent().Provider, tm)
}
```

#### Handle Mentions
```go
func (al *AgentLoop) HandleMention(ctx, mentionedID, 
    message, channel, chatID string) (string, error) {
    return al.ProcessWithAgent(ctx, mentionedID, message, 
        fmt.Sprintf("mention:%s", mentionedID), channel, chatID)
}
```

---

## 2. Team Manager (`pkg/team/manager.go`)

### Struct Definition
```go
type TeamManager struct {
    registry         *agent.AgentRegistry
    bus              *bus.MessageBus
    teams            map[string]*Team
    roleCapabilities map[string][]string
    teamMemory       *memory.TeamMemory
    metrics          *MetricsCollector
    workspace        string
    agentExecutor    AgentExecutor
    agentPool        *AgentPool
    roleCache        *RoleCache
    provider         providers.LLMProvider
    cfg              *config.Config
    mu               sync.RWMutex
}
```

### Key Methods

#### `CreateTeam(ctx, teamConfig) (*Team, error)`
Create team from configuration:
1. Validate config
2. Create team instance
3. Register agents for each role
4. Set coordinator
5. Persist team state
6. Return team

#### `ExecuteTask(ctx, teamID, taskDescription) (any, error)`
Execute task with team:
1. Get team
2. Use intelligent router to determine role
3. Create coordinator
4. Execute based on pattern (sequential/parallel/hierarchical)
5. Return result

#### `ExecuteTaskWithRole(ctx, teamID, task, role) (any, error)`
Execute with specific role:
1. Validate role exists
2. Create coordinator
3. Execute task
4. Return result

### Team Patterns

#### Sequential
```go
for _, task := range tasks {
    result, err := executeTask(task)
    if err != nil {
        return nil, err
    }
    results = append(results, result)
}
```

#### Parallel
```go
var wg sync.WaitGroup
for _, task := range tasks {
    wg.Add(1)
    go func(t *Task) {
        defer wg.Done()
        result, err := executeTask(t)
        results <- result
    }(task)
}
wg.Wait()
```

#### Hierarchical
```go
coordinator := GetCoordinator()
result, err := coordinator.DelegateTask(task)
return result, err
```

---

## 3. Collaborative Chat (`pkg/collaborative/manager_improved.go`)

### Struct Definition
```go
type ManagerV2 struct {
    sessions          map[int64]*Session
    mu                sync.RWMutex
    dispatchTracker   *DispatchTracker
    queueManager      *QueueManager
    compactionManager *CompactionManager
    maxMentionDepth   int
    config            *Config
}
```

### Key Methods

#### `HandleMentions(ctx, platform, chatID, teamID, content, mentions, sender, maxContext)`
Process mentions:
1. Get or create session
2. Check mention depth limit
3. Add user message to context
4. For each mention:
   - Create MentionRequest
   - Enqueue to QueueManager
5. Return

#### `executeAgentAndCascadeWithError(...)`
Execute agent and handle cascading:
1. Check depth limit
2. Check context cancellation
3. Mark agent in cascade
4. Build prompt with context
5. Execute agent
6. Send response
7. Trigger compaction if needed
8. Detect new mentions in response
9. Filter ack-mentions (prevent loops)
10. Enqueue cascaded mentions

### Ack-Loop Prevention
```go
// Don't mention back the person who just mentioned you
if mentionedBy != "" && mentioned == mentionedBy {
    logger.InfoCF("collaborative", "Filtered ack-mention", ...)
    continue
}
```

### Queue Management

#### MentionQueue
- Per-role queue
- Rate limiting (default: 2s between executions)
- Retry with exponential backoff
- Metrics tracking

#### QueueManager
- Manages queues for all roles
- Creates queues on-demand
- Provides metrics aggregation

---

## 4. Memory System

### Vector Store (`pkg/memory/vector/qdrant_store.go`)

#### Struct Definition
```go
type QdrantStore struct {
    client    *qdrant.Client
    config    QdrantConfig
    vectorCfg VectorConfig
    breaker   *CircuitBreaker
    idCounter atomic.Uint64
    mu        sync.RWMutex
}
```

#### Key Methods

##### `Upsert(ctx, vectors) error`
Store vectors:
1. Validate collection
2. Batch vectors (max 100)
3. For each batch:
   - Convert to Qdrant format
   - Upsert with retry
   - Handle errors
4. Return

##### `Search(ctx, query, topK) ([]Vector, error)`
Search similar vectors:
1. Validate topK (max 100)
2. Generate embedding
3. Search with circuit breaker
4. Convert results
5. Return

##### Circuit Breaker
```go
type CircuitBreaker struct {
    state        State  // Closed, Open, HalfOpen
    failures     int
    lastFailTime time.Time
    config       CircuitBreakerConfig
    mu           sync.RWMutex
}
```

States:
- **Closed**: Normal operation
- **Open**: Too many failures, reject calls
- **HalfOpen**: Testing if service recovered

### Embedding Service (`pkg/embedding/openai.go`)

#### Struct Definition
```go
type OpenAIService struct {
    config Config
    client *http.Client
    mu     sync.RWMutex
}
```

#### Key Methods

##### `Generate(ctx, text) ([]float32, error)`
Generate embedding:
1. Create request
2. Call OpenAI API
3. Parse response
4. Return embedding

##### `GenerateBatch(ctx, texts) ([][]float32, error)`
Batch generation:
1. Split into batches
2. For each batch:
   - Generate embeddings
   - Collect results
3. Return all embeddings

---

## 5. Tools System (`pkg/tools/`)

### Tool Registry

#### Struct Definition
```go
type ToolRegistry struct {
    tools map[string]Tool
    mu    sync.RWMutex
}
```

#### Tool Interface
```go
type Tool interface {
    Name() string
    Description() string
    InputSchema() map[string]any
    Execute(ctx, args) (ToolResult, error)
}
```

### Built-in Tools

#### File Tools
- `read_file`: Read file content
- `write_file`: Write to file
- `edit_file`: Edit file content
- `append_file`: Append to file
- `list_dir`: List directory

#### Shell Tool
- `shell_exec`: Execute shell command
- Restrictions: workspace, whitelist
- Timeout handling

#### Web Tools
- `web_search`: Search web (Brave/Tavily/DuckDuckGo)
- `web_fetch`: Fetch URL content

#### Team Tools
- `delegate_to_team`: Delegate task to team
- `team_status`: Get team status

#### Subagent Tool
- `spawn_subagent`: Spawn background agent

### Tool Execution Flow
```
1. LLM requests tool call
2. Parse tool name and args
3. Get tool from registry
4. Validate permissions
5. Execute tool
6. Collect result
7. Format for LLM
8. Return to LLM
```

---

## 6. Session Management (`pkg/session/manager.go`)

### Struct Definition
```go
type SessionManager struct {
    sessions map[string]*Session
    mu       sync.RWMutex
    storage  string
}

type Session struct {
    Key      string
    Messages []providers.Message
    Summary  string
    Created  time.Time
    Updated  time.Time
}
```

### Key Methods

#### `GetOrCreate(key) *Session`
Get or create session:
1. Lock
2. Check if exists
3. If not, create new
4. Return session

#### `AddMessage(sessionKey, role, content)`
Add message to session:
1. Lock
2. Get or create session
3. Append message
4. Update timestamp
5. Unlock

#### `Save(sessionKey) error`
Persist session to disk:
1. Get session
2. Marshal to JSON
3. Write to file atomically
4. Return

#### `Load(sessionKey) error`
Load session from disk:
1. Read file
2. Unmarshal JSON
3. Store in memory
4. Return

---

## 7. Configuration (`pkg/config/config.go`)

### Main Config Structure
```go
type Config struct {
    Agents    AgentsConfig
    Bindings  []AgentBinding
    Session   SessionConfig
    Channels  ChannelsConfig
    Providers ProvidersConfig
    ModelList []ModelConfig
    Gateway   GatewayConfig
    Tools     ToolsConfig
    Heartbeat HeartbeatConfig
    Memory    MemoryConfig
    Devices   DevicesConfig
}
```

### Agent Config
```go
type AgentConfig struct {
    ID                  string
    Name                string
    Model               *AgentModelConfig
    Workspace           string
    MaxIterations       int
    MaxTokens           int
    Temperature         float64
    ContextWindow       int
    RestrictToWorkspace *bool
    Subagents           *SubagentsConfig
    SkillsFilter        []string
}
```

### Memory Config
```go
type MemoryConfig struct {
    Enabled     bool
    Embedding   EmbeddingConfig
    VectorStore VectorStoreConfig
    Cache       CacheConfig
}
```

### Loading Config
```go
func LoadConfig(path string) (*Config, error) {
    // 1. Read file
    data, err := os.ReadFile(path)
    
    // 2. Unmarshal JSON
    var cfg Config
    err = json.Unmarshal(data, &cfg)
    
    // 3. Apply defaults
    applyDefaults(&cfg)
    
    // 4. Validate
    err = validateConfig(&cfg)
    
    return &cfg, nil
}
```

---

## 8. Message Bus (`pkg/bus/bus.go`)

### Struct Definition
```go
type MessageBus struct {
    inbound       chan InboundMessage
    outbound      chan OutboundMessage
    outboundMedia chan OutboundMediaMessage
    done          chan struct{}
    closed        atomic.Bool
}
```

### Message Types

#### InboundMessage
```go
type InboundMessage struct {
    Channel    string
    SenderID   string
    ChatID     string
    Content    string
    Media      []string
    SessionKey string
    Metadata   map[string]string
}
```

#### OutboundMessage
```go
type OutboundMessage struct {
    Channel string
    ChatID  string
    Content string
}
```

### Key Methods

#### `PublishInbound(ctx, msg) error`
Publish to inbound:
1. Check if closed
2. Check context
3. Select on channels:
   - Send to inbound
   - Done channel
   - Context done
4. Return

#### `ConsumeInbound(ctx) (InboundMessage, bool)`
Consume from inbound:
1. Select on channels:
   - Receive from inbound
   - Done channel
   - Context done
2. Return message

#### `Close() error`
Graceful shutdown:
1. Set closed flag
2. Close done channel
3. Close all channels
4. Return

---

## Next Steps

Đọc thêm:
- [Data Flow Details](./DATA_FLOW.md)
- [Configuration Guide](./CONFIGURATION_GUIDE.md)
- [Development Guide](./DEVELOPMENT_GUIDE.md)
