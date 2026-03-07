# Collaborative Chat Package

Platform-agnostic collaborative chat system for multi-agent conversations with @mention routing.

## Overview

The `collaborative` package provides a reusable framework for implementing IRC-style collaborative chat across different messaging platforms (Telegram, IRC, Discord, etc.). Agents can communicate with each other using @mentions, creating natural multi-agent conversations.

## Features

- **Platform-Agnostic**: Works with any messaging platform via the `Platform` interface
- **@Mention Routing**: Automatic agent triggering via @mentions
- **Agent-to-Agent Communication**: Agents can mention and trigger other agents
- **Team Awareness**: Agents automatically know about all team members via roster
- **Session Management**: Maintains conversation context and history
- **Parallel Execution**: Multiple agents can respond simultaneously
- **IRC-Style Formatting**: Clean, readable message formatting with role emojis
- **Team Roster**: Automatic team member discovery and listing in prompts

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    Platform Layer                        │
│  (Telegram, IRC, Discord, Slack, etc.)                  │
└────────────────────┬────────────────────────────────────┘
                     │ implements Platform interface
                     ▼
┌─────────────────────────────────────────────────────────┐
│              Collaborative Manager                       │
│  - Session management                                    │
│  - Mention extraction and routing                        │
│  - Agent coordination                                    │
└────────────────────┬────────────────────────────────────┘
                     │ uses
                     ▼
┌─────────────────────────────────────────────────────────┐
│                 Team Manager                             │
│  - Agent execution                                       │
│  - Role-based task delegation                            │
└─────────────────────────────────────────────────────────┘
```

## Quick Start

### 1. Implement Platform Interface

```go
import "github.com/sipeed/picoclaw/pkg/collaborative"

type MyChannel struct {
    // your channel fields
    chatManager *collaborative.Manager
}

// Implement Platform interface
func (c *MyChannel) SendMessage(ctx context.Context, chatID string, content string) error {
    // Send message to your platform
    return nil
}

func (c *MyChannel) GetTeamManager() collaborative.TeamManager {
    // Return your team manager
    return c.teamManager
}

func (c *MyChannel) GetContext() context.Context {
    // Return channel context
    return c.ctx
}
```

### 2. Create Manager

```go
func NewMyChannel() *MyChannel {
    return &MyChannel{
        chatManager: collaborative.NewManager(),
    }
}
```

### 3. Handle Messages

```go
func (c *MyChannel) handleMessage(ctx context.Context, chatID int64, content string) {
    // Extract mentions
    mentions := collaborative.ExtractMentions(content)
    
    if len(mentions) > 0 {
        // Handle collaborative chat
        c.chatManager.HandleMentions(
            ctx,
            c,              // Platform interface
            chatID,
            teamID,
            content,
            mentions,
            sender,
            maxContextLength,
        )
    }
}
```

## Core Types

### Session

Represents a collaborative chat session with conversation history and active agents.

```go
type Session struct {
    SessionID    string
    TeamID       string
    ChatID       int64
    StartTime    time.Time
    LastActivity time.Time
    Context      []Message
    ActiveAgents map[string]*AgentState
    MaxContext   int
}
```

### Platform Interface

Platform-specific operations that must be implemented.

```go
type Platform interface {
    SendMessage(ctx context.Context, chatID string, content string) error
    GetTeamManager() TeamManager
    GetContext() context.Context
}
```

### TeamManager Interface

Team operations for agent execution.

```go
type TeamManager interface {
    ExecuteTaskWithRole(ctx context.Context, teamID, prompt, role string) (any, error)
    GetTeam(teamID string) (any, error)
}
```

## Usage Examples

### Basic Mention Handling

```go
// User: "@developer can you help?"
mentions := collaborative.ExtractMentions(content)
// mentions = ["developer"]

manager.HandleMentions(ctx, platform, chatID, teamID, content, mentions, sender, 50)
// Developer agent is triggered and responds
```

### Agent-to-Agent Communication

```go
// Developer: "Sure! @tester can you verify this?"
// Tester is automatically triggered by the mention
// Both agents can see full conversation context
```

### Session Management

```go
// Get or create session
session := manager.GetOrCreateSession(chatID, teamID, maxContext)

// Add message to context
session.AddMessage("user", "Hello @developer", []string{"developer"})

// Get formatted context for LLM
contextStr := session.GetContextAsString()

// Update agent status
session.UpdateAgentStatus("developer", "thinking")
```

### Message Formatting

```go
// Format agent response
msg := collaborative.FormatMessage(sessionID, "developer", "I can help with that!")
// Output: "[chat51263350] 💻 DEVELOPER: I can help with that!"

// Get role emoji
emoji := collaborative.GetRoleEmoji("tester")
// Output: "🧪"
```

## API Reference

### Manager

#### `NewManager() *Manager`
Creates a new collaborative chat manager.

#### `HandleMentions(ctx, platform, chatID, teamID, content, mentions, sender, maxContext) error`
Processes @mentions and triggers appropriate agents.

#### `GetOrCreateSession(chatID, teamID, maxContext) *Session`
Gets existing session or creates new one.

#### `GetSession(chatID) (*Session, bool)`
Gets existing session if it exists.

#### `RemoveSession(chatID)`
Removes a session.

### Session

#### `AddMessage(author, content string, mentions []string)`
Adds message to conversation context.

#### `GetContextAsString() string`
Formats context as string for LLM prompts.

#### `UpdateAgentStatus(role, status string)`
Updates agent status (idle, thinking, busy, error).

#### `GetAgentStatus(role string) string`
Gets current agent status.

#### `GetActiveAgents() []string`
Returns list of active agent roles.

#### `GetAgentState(role string) *AgentState`
Gets detailed agent state.

### Utilities

#### `ExtractMentions(text string) []string`
Extracts @mentions from text.

#### `FormatMessage(sessionID, role, content string) string`
Formats message in IRC style.

#### `GetRoleEmoji(role string) string`
Returns emoji for role.

#### `BuildTeamRoster(teamInfo any) string`
Builds formatted team roster.

## Configuration

### v2.0.5+ Advanced Configuration

The collaborative package supports advanced configuration for high-performance scenarios:

```json
{
  "channels": {
    "telegram": {
      "collaborative_chat": {
        "enabled": true,
        "default_team_id": "dev-team",
        "max_context_length": 50,
        "max_depth": 20,
        "max_concurrent_tasks": 200,
        "queue_size": 500,
        "circuit_breaker": {
          "enabled": true,
          "threshold": 300,
          "memory_threshold": "40MB",
          "cooldown": "30s"
        }
      }
    }
  }
}
```

**Key Parameters:**

- `max_depth` (default: 20): Maximum cascade depth for agent-to-agent mentions. Prevents infinite loops while supporting deep multi-agent workflows (configurable up to 50).

- `max_concurrent_tasks` (default: 200): Maximum parallel agent executions. Optimized for 50MB RAM budget.

- `circuit_breaker`: Automatic overload protection with memory-based triggers and cooldown periods.

**Performance:** ~8-15MB under normal load, 200 concurrent tasks capacity, automatic circuit breaker prevents overload.

### Basic Usage

For simple use cases, minimal configuration is needed - the package uses sensible defaults. Configuration is handled by the platform layer.

## Error Handling

The package provides user-friendly error messages:

```go
// Role not found
"❌ Role @designer không tồn tại trong team configuration.
💡 Các role có sẵn: @architect, @developer, @tester, @manager"

// Team manager not available
"⚠️ Collaborative chat is not properly configured. Team manager is not available."
```

## Performance

- **Memory**: Minimal overhead, sessions are lightweight
- **Concurrency**: Agents execute in parallel goroutines
- **Context Trimming**: Automatic context length management
- **No Blocking**: All agent execution is asynchronous

## Thread Safety

All operations are thread-safe:
- Session access uses `sync.RWMutex`
- Manager uses `sync.RWMutex` for session map
- Safe for concurrent access from multiple goroutines

## Logging

Uses `pkg/logger` for structured logging:
- `INFO`: Session creation, agent responses
- `DEBUG`: Message additions, context updates
- `WARN`: Missing team info, complex result types
- `ERROR`: Agent execution failures

## Testing

See platform implementations for integration tests:
- `pkg/channels/telegram/` - Telegram integration
- Future: `pkg/channels/irc/` - IRC integration

## Migration Guide

### From Telegram-Specific Code

Old:
```go
session := c.chatManager.GetOrCreateSession(chatID, teamID, maxContext)
mentions := extractMentions(content)
emoji := getRoleEmoji(role)
```

New:
```go
session := c.chatManager.GetOrCreateSession(chatID, teamID, maxContext)
mentions := collaborative.ExtractMentions(content)
emoji := collaborative.GetRoleEmoji(role)
```

## Future Enhancements

- [ ] Persistent session storage
- [ ] Session history and replay
- [ ] Cross-platform message routing
- [ ] Advanced agent coordination patterns
- [ ] Custom role emojis
- [ ] Per-role context limits
- [ ] Session analytics and metrics

## Contributing

When adding new features:
1. Keep the Platform interface minimal
2. Maintain backward compatibility
3. Add comprehensive logging
4. Document all public APIs
5. Write tests for new functionality

## License

MIT License - See LICENSE file for details.

## See Also

- [Collaborative Chat Guide](../../docs/COLLABORATIVE_CHAT.md)
- [Quick Start](../../docs/COLLABORATIVE_CHAT_QUICKSTART.md)
- [Architecture](../../docs/COLLABORATIVE_CHAT_ARCHITECTURE.md)
- [Multi-Agent Guide](../../docs/MULTI_AGENT_GUIDE.md)


## How Agents Know Each Other

### Team Roster in Prompts

A key feature of the collaborative package is that **every agent automatically knows about all other team members**. This is achieved by including a team roster in every agent's prompt.

### Roster Structure

When an agent is triggered, their prompt includes:

```
=== Team Information ===
Team: Development Team
Members:
  • @architect - System architect and designer
  • @developer - Software developer
  • @tester - QA and testing specialist
  • @manager - Project manager
```

### How It Works

1. **Team Config Loading**: `TeamManager.GetTeam(teamID)` retrieves team configuration
2. **Roster Building**: `BuildTeamRoster(teamInfo)` extracts and formats role information
3. **Prompt Injection**: Roster is automatically included in every agent's prompt
4. **Agent Awareness**: Agents can see all available roles and their descriptions

### Example Prompt

```
=== Collaborative Chat Context ===
Session: chat51263350 | Team: dev-team
Started: 12:46:10
=== Conversation History ===

[12:46:10] USER: @developer can you help?

=== Team Information ===
Team: Development Team
Members:
  • @architect - System architect and designer
  • @developer - Software developer
  • @tester - QA and testing specialist
  • @manager - Project manager

User message: Can you review this code?

You are @developer. Respond to the user's message considering the conversation history above.
You can mention other team members using @role format (e.g., @architect, @tester).
```

### Benefits

1. **Natural Collaboration**: Agents know who to ask for help
2. **Role Awareness**: Each agent understands others' specializations
3. **Automatic Discovery**: No manual configuration needed
4. **Dynamic Updates**: Roster updates when team config changes

### Implementation

```go
// In manager.go HandleMentions()

// Get team roster
var teamRoster string
if teamInfo, err := teamManager.GetTeam(teamID); err == nil {
    teamRoster = BuildTeamRoster(teamInfo)
}

// Build prompt with roster
prompt := fmt.Sprintf(`%s

=== Team Information ===
%s

User message: %s

You are @%s. You can mention other team members using @role format.`,
    contextStr, teamRoster, content, role)
```

### Fallback Roster

If team configuration cannot be loaded, a default roster is provided:

```go
// Ultimate fallback in roster.go
return `Team: Development Team
Members:
  • @architect - System architect and designer
  • @developer - Software developer
  • @tester - QA and testing specialist
  • @manager - Project manager`
```

### Customization

Improve agent collaboration by adding detailed role descriptions in your team config:

```json
{
  "roles": [
    {
      "name": "security",
      "description": "Security specialist - handles authentication, authorization, and security audits",
      "model": "claude-3-5-sonnet-20241022"
    }
  ]
}
```

The more descriptive your role descriptions, the better agents can understand each other's capabilities and collaborate effectively!
