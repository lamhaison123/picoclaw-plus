# PicoClaw API Reference

## Overview

This document provides a comprehensive API reference for PicoClaw's core packages and interfaces.

## Table of Contents

1. [Agent System](#agent-system)
2. [Multi-Agent Framework](#multi-agent-framework)
3. [Collaborative Chat](#collaborative-chat)
4. [Provider System](#provider-system)
5. [Tool System](#tool-system)
6. [Channel Manager](#channel-manager)
7. [Configuration](#configuration)
8. [Message Bus](#message-bus)

## Agent System

### Package: `pkg/agent`

#### AgentInstance

Represents a single agent with its configuration, tools, and provider.

```go
type AgentInstance struct {
    ID              string
    Name            string
    Model           string
    Provider        providers.LLMProvider
    ToolsRegistry   *tools.ToolRegistry
    SessionsManager *session.SessionManager
    ContextBuilder  *ContextBuilder
    // ...
}
```

**Methods:**

```go
// NewAgentInstance creates a new agent instance
func NewAgentInstance(
    agentCfg *config.AgentConfig,
    defaults *config.AgentDefaults,
    cfg *config.Config,
    provider providers.LLMProvider,
) *AgentInstance

// Execute runs the agent with a task
func (a *AgentInstance) Execute(ctx context.Context, task string) (string, error)

// ExecuteWithSession runs the agent with session context
func (a *AgentInstance) ExecuteWithSession(
    ctx context.Context,
    sessionID string,
    task string,
) (string, error)
```

#### AgentRegistry

Manages agent instances and configurations.

```go
type AgentRegistry struct {
    agents map[string]*AgentInstance
    mu     sync.RWMutex
}
```

**Methods:**

```go
// NewAgentRegistry creates a new agent registry
func NewAgentRegistry(cfg *config.Config, provider providers.LLMProvider) *AgentRegistry

// Get retrieves an agent by ID
func (r *AgentRegistry) Get(agentID string) (*AgentInstance, bool)

// Register registers a new agent
func (r *AgentRegistry) Register(agent *AgentInstance)

// List returns all registered agents
func (r *AgentRegistry) List() []*AgentInstance
```

#### ContextBuilder

Builds context for LLM prompts with caching support.

```go
type ContextBuilder struct {
    workspace string
    cache     *ContextCache
}
```

**Methods:**

```go
// NewContextBuilder creates a new context builder
func NewContextBuilder(workspace string) *ContextBuilder

// BuildContext builds context from messages
func (cb *ContextBuilder) BuildContext(
    messages []Message,
    systemPrompt string,
) ([]Message, error)

// CacheContext caches context for reuse
func (cb *ContextBuilder) CacheContext(key string, messages []Message)
```

## Multi-Agent Framework

### Package: `pkg/team`

#### TeamManager

Orchestrates agent teams and task delegation.

```go
type TeamManager struct {
    registry         *agent.AgentRegistry
    bus              *bus.MessageBus
    teams            map[string]*Team
    roleCapabilities map[string][]string
    teamMemory       *memory.TeamMemory
    // ...
}
```

**Methods:**

```go
// NewTeamManager creates a new team manager
func NewTeamManager(
    registry *agent.AgentRegistry,
    msgBus *bus.MessageBus,
) *TeamManager

// CreateTeam creates a new team from configuration
func (tm *TeamManager) CreateTeam(
    ctx context.Context,
    config *TeamConfig,
) (*Team, error)

// ExecuteTask executes a task with the team
func (tm *TeamManager) ExecuteTask(
    ctx context.Context,
    teamID string,
    task string,
) (any, error)

// ExecuteTaskWithRole executes a task with a specific role
func (tm *TeamManager) ExecuteTaskWithRole(
    ctx context.Context,
    teamID string,
    role string,
    task string,
) (any, error)

// DissolveTeam dissolves a team and cleans up resources
func (tm *TeamManager) DissolveTeam(ctx context.Context, teamID string) error

// ListTeams returns all active teams
func (tm *TeamManager) ListTeams() []*Team
```

#### Team

Represents a team with roles and coordination patterns.

```go
type Team struct {
    ID          string
    Name        string
    Description string
    Pattern     string
    Roles       []Role
    Coordinator *Coordinator
    Settings    TeamSettings
}
```

#### TeamConfig

Configuration for creating teams.

```go
type TeamConfig struct {
    TeamID      string       `json:"team_id"`
    Name        string       `json:"name"`
    Description string       `json:"description"`
    Pattern     string       `json:"pattern"`
    Roles       []RoleConfig `json:"roles"`
    Coordinator struct {
        Role    string `json:"role"`
        AgentID string `json:"agent_id,omitempty"`
    } `json:"coordinator"`
    Settings TeamSettings `json:"settings"`
}
```

**Functions:**

```go
// LoadTeamConfig loads team configuration from file
func LoadTeamConfig(path string) (*TeamConfig, error)

// SaveTeamConfig saves team configuration to file
func SaveTeamConfig(config *TeamConfig, path string) error
```

#### ConsensusManager

Handles voting and consensus mechanisms.

```go
type ConsensusManager struct {
    method  string
    weights map[string]float64
}
```

**Methods:**

```go
// NewConsensusManager creates a new consensus manager
func NewConsensusManager(method string, weights map[string]float64) *ConsensusManager

// ReachConsensus determines consensus from votes
func (cm *ConsensusManager) ReachConsensus(votes map[string]any) (any, error)
```

## Collaborative Chat

### Package: `pkg/collaborative`

#### Manager

Manages collaborative sessions and mention routing.

```go
type Manager struct {
    sessions     map[int64]*Session
    queueManager *QueueManager
    config       *Config
    mu           sync.RWMutex
}
```

**Methods:**

```go
// NewManager creates a new collaborative manager
func NewManager() *Manager

// NewManagerWithConfig creates a manager with custom config
func NewManagerWithConfig(config *Config) *Manager

// GetOrCreateSession gets or creates a session
func (m *Manager) GetOrCreateSession(
    chatID int64,
    teamID string,
    maxContext int,
) *Session

// HandleMentions handles @mentions in messages
func (m *Manager) HandleMentions(
    ctx context.Context,
    platform Platform,
    chatID int64,
    teamID string,
    content string,
    mentions []string,
    sender string,
    maxContextLength int,
)
```

#### Session

Maintains conversation context per chat.

```go
type Session struct {
    SessionID    string
    ChatID       int64
    TeamID       string
    Context      []ContextMessage
    AgentStates  map[string]string
    MaxContext   int
    CreatedAt    time.Time
    LastActivity time.Time
    mu           sync.RWMutex
}
```

**Methods:**

```go
// NewSession creates a new session
func NewSession(chatID int64, teamID string, maxContext int) *Session

// AddMessage adds a message to context
func (s *Session) AddMessage(
    role string,
    content string,
    mentions []string,
)

// GetContextAsString returns formatted context
func (s *Session) GetContextAsString() string

// UpdateAgentStatus updates agent status
func (s *Session) UpdateAgentStatus(role string, status string)
```

#### Platform Interface

Abstraction for messaging platforms.

```go
type Platform interface {
    SendMessage(ctx context.Context, chatID string, content string) error
    GetTeamManager() TeamManager
    GetContext() context.Context
}
```

#### Functions

```go
// ExtractMentions extracts @mentions from text
func ExtractMentions(text string) []string

// FormatMessage formats message with IRC-style prefix
func FormatMessage(sessionID string, role string, content string) string

// GetRoleEmoji returns emoji for role
func GetRoleEmoji(role string) string

// BuildTeamRoster builds team roster string
func BuildTeamRoster(teamInfo any) string
```

## Provider System

### Package: `pkg/providers`

#### LLMProvider Interface

Interface for LLM providers.

```go
type LLMProvider interface {
    Complete(ctx context.Context, messages []Message) (*Response, error)
    CompleteStream(ctx context.Context, messages []Message) (<-chan StreamChunk, error)
    Name() string
    Model() string
}
```

#### Message

Represents a message in conversation.

```go
type Message struct {
    Role       string      `json:"role"`
    Content    string      `json:"content,omitempty"`
    ToolCalls  []ToolCall  `json:"tool_calls,omitempty"`
    ToolCallID string      `json:"tool_call_id,omitempty"`
    Name       string      `json:"name,omitempty"`
}
```

#### Response

Represents LLM response.

```go
type Response struct {
    Content    string
    ToolCalls  []ToolCall
    FinishReason string
    Usage      Usage
}
```

#### Functions

```go
// NewProvider creates a provider from configuration
func NewProvider(cfg *config.Config) (LLMProvider, error)

// NewProviderWithFallback creates a provider with fallback chain
func NewProviderWithFallback(
    cfg *config.Config,
    fallbacks []string,
) (LLMProvider, error)
```

## Tool System

### Package: `pkg/tools`

#### Tool Interface

Interface for tools.

```go
type Tool interface {
    Name() string
    Description() string
    Parameters() map[string]any
    Execute(ctx context.Context, args map[string]any) *ToolResult
}
```

#### ToolRegistry

Manages tool registration and execution.

```go
type ToolRegistry struct {
    tools map[string]Tool
    mu    sync.RWMutex
}
```

**Methods:**

```go
// NewToolRegistry creates a new tool registry
func NewToolRegistry() *ToolRegistry

// Register registers a tool
func (r *ToolRegistry) Register(tool Tool)

// Get retrieves a tool by name
func (r *ToolRegistry) Get(name string) (Tool, bool)

// Execute executes a tool
func (r *ToolRegistry) Execute(
    ctx context.Context,
    name string,
    args map[string]any,
) *ToolResult

// ExecuteWithContext executes a tool with context
func (r *ToolRegistry) ExecuteWithContext(
    ctx context.Context,
    name string,
    args map[string]any,
    channel string,
    chatID string,
    asyncCallback AsyncCallback,
) *ToolResult

// List returns all registered tools
func (r *ToolRegistry) List() []Tool
```

#### ToolResult

Represents tool execution result.

```go
type ToolResult struct {
    Success bool
    Content string
    Error   error
    Data    map[string]any
}
```

**Functions:**

```go
// SuccessResult creates a success result
func SuccessResult(content string) *ToolResult

// ErrorResult creates an error result
func ErrorResult(message string) *ToolResult

// DataResult creates a result with data
func DataResult(data map[string]any) *ToolResult
```

#### Built-in Tools

```go
// File operations
func NewReadFileTool(workspace string, restrict bool, allowPaths []string) Tool
func NewWriteFileTool(workspace string, restrict bool, allowPaths []string) Tool
func NewEditFileTool(workspace string, restrict bool, allowPaths []string) Tool
func NewAppendFileTool(workspace string, restrict bool, allowPaths []string) Tool
func NewListDirTool(workspace string, restrict bool, allowPaths []string) Tool

// Shell execution
func NewExecToolWithConfig(workspace string, restrict bool, cfg *config.Config) (Tool, error)

// Web tools
func NewWebSearchTool(apiKey string, provider string) Tool
func NewWebFetchTool() Tool

// Code tools
func NewReadCodeTool(workspace string) Tool
func NewEditCodeTool(workspace string) Tool
func NewGetDiagnosticsTool() Tool
```

## Channel Manager

### Package: `pkg/channels`

#### ChannelManager

Manages messaging platform channels.

```go
type ChannelManager struct {
    channels map[string]Channel
    bus      *bus.MessageBus
    mu       sync.RWMutex
}
```

**Methods:**

```go
// NewChannelManager creates a new channel manager
func NewChannelManager(bus *bus.MessageBus) *ChannelManager

// Register registers a channel
func (cm *ChannelManager) Register(name string, channel Channel)

// Start starts all channels
func (cm *ChannelManager) Start(ctx context.Context) error

// Stop stops all channels
func (cm *ChannelManager) Stop(ctx context.Context) error
```

#### Channel Interface

Interface for messaging channels.

```go
type Channel interface {
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
    Name() string
}
```

#### Optional Interfaces

```go
// TypingCapable - channels that can show typing indicator
type TypingCapable interface {
    StartTyping(ctx context.Context, chatID string) (stop func(), err error)
}

// MessageEditor - channels that can edit messages
type MessageEditor interface {
    EditMessage(ctx context.Context, chatID string, messageID string, content string) error
}

// ReactionCapable - channels that can add reactions
type ReactionCapable interface {
    ReactToMessage(ctx context.Context, chatID, messageID string) (undo func(), err error)
}

// PlaceholderCapable - channels that can send placeholder messages
type PlaceholderCapable interface {
    SendPlaceholder(ctx context.Context, chatID string) (messageID string, err error)
}
```

## Configuration

### Package: `pkg/config`

#### Config

Main configuration structure.

```go
type Config struct {
    Agents    AgentsConfig    `json:"agents"`
    Bindings  []AgentBinding  `json:"bindings,omitempty"`
    Session   SessionConfig   `json:"session,omitempty"`
    Channels  ChannelsConfig  `json:"channels"`
    Providers ProvidersConfig `json:"providers,omitempty"`
    ModelList []ModelConfig   `json:"model_list"`
    Gateway   GatewayConfig   `json:"gateway"`
    Tools     ToolsConfig     `json:"tools"`
    Heartbeat HeartbeatConfig `json:"heartbeat"`
    Devices   DevicesConfig   `json:"devices"`
}
```

**Functions:**

```go
// LoadConfig loads configuration from file
func LoadConfig(path string) (*Config, error)

// SaveConfig saves configuration to file
func SaveConfig(config *Config, path string) error

// DefaultConfig returns default configuration
func DefaultConfig() *Config
```

#### AgentsConfig

Agent configuration.

```go
type AgentsConfig struct {
    Defaults AgentDefaults `json:"defaults"`
    List     []AgentConfig `json:"list,omitempty"`
}

type AgentDefaults struct {
    Model                   string   `json:"model"`
    Fallbacks               []string `json:"fallbacks,omitempty"`
    Provider                string   `json:"provider,omitempty"`
    Workspace               string   `json:"workspace,omitempty"`
    RestrictToWorkspace     bool     `json:"restrict_to_workspace"`
    AllowReadOutsideWorkspace bool   `json:"allow_read_outside_workspace"`
}

type AgentConfig struct {
    ID        string   `json:"id"`
    Name      string   `json:"name,omitempty"`
    Model     string   `json:"model,omitempty"`
    Fallbacks []string `json:"fallbacks,omitempty"`
    Workspace string   `json:"workspace,omitempty"`
}
```

## Message Bus

### Package: `pkg/bus`

#### MessageBus

Event-driven message bus for inter-component communication.

```go
type MessageBus struct {
    inbound       chan InboundMessage
    outbound      chan OutboundMessage
    outboundMedia chan OutboundMediaMessage
    done          chan struct{}
    closed        atomic.Bool
}
```

**Methods:**

```go
// NewMessageBus creates a new message bus
func NewMessageBus() *MessageBus

// PublishInbound publishes an inbound message
func (mb *MessageBus) PublishInbound(ctx context.Context, msg InboundMessage) error

// ConsumeInbound consumes an inbound message
func (mb *MessageBus) ConsumeInbound(ctx context.Context) (InboundMessage, bool)

// PublishOutbound publishes an outbound message
func (mb *MessageBus) PublishOutbound(ctx context.Context, msg OutboundMessage) error

// ConsumeOutbound consumes an outbound message
func (mb *MessageBus) ConsumeOutbound(ctx context.Context) (OutboundMessage, bool)

// Close closes the message bus
func (mb *MessageBus) Close()
```

#### Message Types

```go
type InboundMessage struct {
    Channel   string
    ChatID    string
    UserID    string
    Content   string
    MessageID string
    Timestamp time.Time
}

type OutboundMessage struct {
    Channel   string
    ChatID    string
    Content   string
    ReplyTo   string
    Timestamp time.Time
}

type OutboundMediaMessage struct {
    Channel   string
    ChatID    string
    MediaType string
    MediaData []byte
    Caption   string
    Timestamp time.Time
}
```

## Error Handling

### Common Errors

```go
// Agent errors
var (
    ErrAgentNotFound     = errors.New("agent not found")
    ErrInvalidConfig     = errors.New("invalid configuration")
    ErrProviderFailed    = errors.New("provider failed")
    ErrToolExecutionFailed = errors.New("tool execution failed")
)

// Team errors
var (
    ErrTeamNotFound      = errors.New("team not found")
    ErrRoleNotFound      = errors.New("role not found")
    ErrInvalidPattern    = errors.New("invalid collaboration pattern")
    ErrConsensusNotReached = errors.New("consensus not reached")
)

// Channel errors
var (
    ErrChannelNotFound   = errors.New("channel not found")
    ErrMessageSendFailed = errors.New("message send failed")
    ErrInvalidChatID     = errors.New("invalid chat ID")
)
```

## Examples

### Creating and Using an Agent

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/sipeed/picoclaw/pkg/agent"
    "github.com/sipeed/picoclaw/pkg/config"
    "github.com/sipeed/picoclaw/pkg/providers"
)

func main() {
    // Load configuration
    cfg, err := config.LoadConfig("~/.picoclaw/config.json")
    if err != nil {
        panic(err)
    }
    
    // Create provider
    provider, err := providers.NewProvider(cfg)
    if err != nil {
        panic(err)
    }
    
    // Create agent
    agentInst := agent.NewAgentInstance(
        nil,
        &cfg.Agents.Defaults,
        cfg,
        provider,
    )
    
    // Execute task
    ctx := context.Background()
    result, err := agentInst.Execute(ctx, "Write a hello world function in Go")
    if err != nil {
        panic(err)
    }
    
    fmt.Println(result)
}
```

### Creating and Using a Team

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/sipeed/picoclaw/pkg/agent"
    "github.com/sipeed/picoclaw/pkg/bus"
    "github.com/sipeed/picoclaw/pkg/config"
    "github.com/sipeed/picoclaw/pkg/providers"
    "github.com/sipeed/picoclaw/pkg/team"
)

func main() {
    // Load configuration
    cfg, err := config.LoadConfig("~/.picoclaw/config.json")
    if err != nil {
        panic(err)
    }
    
    // Create provider
    provider, err := providers.NewProvider(cfg)
    if err != nil {
        panic(err)
    }
    
    // Create agent registry
    registry := agent.NewAgentRegistry(cfg, provider)
    
    // Create message bus
    msgBus := bus.NewMessageBus()
    defer msgBus.Close()
    
    // Create team manager
    tm := team.NewTeamManager(registry, msgBus)
    
    // Load team configuration
    teamCfg, err := team.LoadTeamConfig("templates/teams/development-team.json")
    if err != nil {
        panic(err)
    }
    
    // Create team
    ctx := context.Background()
    myTeam, err := tm.CreateTeam(ctx, teamCfg)
    if err != nil {
        panic(err)
    }
    
    // Execute task
    result, err := tm.ExecuteTask(ctx, myTeam.ID, "Implement user authentication")
    if err != nil {
        panic(err)
    }
    
    fmt.Println(result)
    
    // Dissolve team
    tm.DissolveTeam(ctx, myTeam.ID)
}
```

---

**Last Updated**: 2026-03-07  
**Version**: 1.1.1
