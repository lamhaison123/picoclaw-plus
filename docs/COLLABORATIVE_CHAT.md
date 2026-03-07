# Collaborative Multi-Agent Chat

Native IRC-style collaborative chat integration for PicoClaw, allowing multiple agents to participate in conversations within a single chat across different messaging platforms.

## Overview

Collaborative chat enables users to mention multiple agents in a single message (e.g., `@architect @developer How should we implement this?`), and all mentioned agents will respond in parallel with full conversation context. Responses are formatted in IRC-style with session IDs and role emojis.

The collaborative chat system is built on a platform-agnostic architecture (`pkg/collaborative/`) that can be integrated with any messaging platform (Telegram, IRC, Discord, etc.).

## Features

- **@mention-based routing**: Mention agents by role (e.g., `@architect`, `@developer`)
- **Parallel execution**: Multiple agents respond simultaneously
- **Shared context**: All agents see the full conversation history
- **IRC-style formatting**: Messages prefixed with `[session-id] emoji ROLE: message`
- **Session management**: Automatic session creation per chat
- **Context window**: Configurable conversation history length
- **Agent-to-agent communication**: Agents can mention and trigger other agents
- **Platform-agnostic**: Works with Telegram, IRC, and other platforms

## Architecture

```
User Message with @mentions
    ↓
Platform Channel (Telegram/IRC/etc.)
    ↓
Extract mentions (@architect, @developer, etc.)
    ↓
collaborative.Manager.HandleMentions()
    ↓
Get/Create Session
    ↓
Add message to context
    ↓
Get TeamManager from Platform
    ↓
For each mentioned role (parallel):
    - Build prompt with full context + team roster
    - Execute TeamManager.ExecuteTaskWithRole()
    - Extract clean response text
    - Format with IRC-style prefix
    - Send via Platform.SendMessage()
    - Check for agent-to-agent mentions
    - Trigger mentioned agents recursively
```

### Package Structure

```
pkg/collaborative/          # Platform-agnostic collaborative chat
├── types.go               # Core types and interfaces
├── session.go             # Session management
├── manager.go             # Collaborative manager
├── mention.go             # @mention extraction
├── formatting.go          # IRC-style formatting
└── roster.go              # Team roster building

pkg/channels/telegram/     # Telegram integration
pkg/channels/irc/          # IRC integration (future)
```

## Configuration

### 1. Enable Collaborative Chat in Config

Edit `~/.picoclaw/config.json`:

```json
{
  "channels": {
    "telegram": {
      "enabled": true,
      "token": "YOUR_BOT_TOKEN",
      "allow_from": ["YOUR_USER_ID"],
      "collaborative_chat": {
        "enabled": true,
        "default_team_id": "dev-team",
        "max_context_length": 50,
        "role_map": {
          "architect": "architect",
          "developer": "developer",
          "tester": "tester",
          "manager": "manager"
        }
      }
    }
  }
}
```

**Configuration Fields:**
- `enabled`: Enable/disable collaborative chat
- `default_team_id`: Team ID to use for agent execution
- `max_context_length`: Maximum number of messages to keep in context (default: 50)
- `role_map`: Optional mapping of mention names to role names

### 2. Create Team Configuration

Create a team configuration file (e.g., `~/.picoclaw/teams/dev-team.json`):

```json
{
  "team_id": "dev-team",
  "name": "Development Team",
  "description": "Collaborative development team",
  "pattern": "parallel",
  "roles": [
    {
      "name": "architect",
      "description": "System architect and designer",
      "capabilities": ["design", "architecture", "planning"],
      "model": "claude-3-5-sonnet-20241022",
      "tools": ["*"]
    },
    {
      "name": "developer",
      "description": "Software developer",
      "capabilities": ["coding", "implementation", "debugging"],
      "model": "claude-3-5-sonnet-20241022",
      "tools": ["*"]
    },
    {
      "name": "tester",
      "description": "QA and testing specialist",
      "capabilities": ["testing", "qa", "validation"],
      "model": "claude-3-5-sonnet-20241022",
      "tools": ["*"]
    },
    {
      "name": "manager",
      "description": "Project manager",
      "capabilities": ["planning", "coordination", "documentation"],
      "model": "claude-3-5-sonnet-20241022",
      "tools": ["*"]
    }
  ],
  "coordinator": {
    "role": "manager"
  },
  "settings": {
    "max_delegation_depth": 3,
    "agent_timeout_seconds": 300
  }
}
```

### 3. Load Team Configuration

Use the `picoclaw team` command to load the team:

```bash
picoclaw team create ~/.picoclaw/teams/dev-team.json
```

Or place the file in `~/.picoclaw/teams/` and it will be auto-loaded on startup.

## Usage

### Basic Usage

Start the gateway:

```bash
picoclaw gateway
```

In Telegram, mention agents in your messages:

```
User: @architect How should we structure the authentication system?

[abc123] 🏗️ ARCHITECT: I recommend a layered approach with JWT tokens...

User: @architect @developer Can you work together on this?

[abc123] 🏗️ ARCHITECT: I'll design the architecture...
[abc123] 💻 DEVELOPER: I can implement the JWT handling...

User: @tester What test cases do we need?

[abc123] 🧪 TESTER: We should test login, logout, token refresh...
```

### Role Emojis

Default role emojis (defined in `collaborative_chat.go`):
- 🏗️ `architect`
- 💻 `developer`
- 🧪 `tester`
- 📋 `manager`
- 🎨 `designer`
- ⚙️ `devops`
- 🤖 (default for unknown roles)

### Session Management

- **Automatic session creation**: Sessions are created automatically per chat
- **Session ID format**: Short hash based on chat ID (e.g., `abc123`)
- **Context window**: Configurable via `max_context_length`
- **Session persistence**: Currently in-memory (optional: add persistence)

## Implementation Details

### Key Components

1. **CollaborativeChatManager** (`pkg/channels/telegram/collaborative_chat.go`)
   - Manages sessions per chat
   - Tracks agent states
   - Maintains conversation context

2. **CollaborativeChatSession**
   - Stores conversation history
   - Tracks active agents
   - Provides context formatting for LLMs

3. **TelegramChannel Integration** (`pkg/channels/telegram/telegram.go`)
   - Detects @mentions in messages
   - Routes to collaborative handler
   - Formats and sends responses

4. **TeamManager Integration** (`pkg/channels/manager.go`, `cmd/picoclaw/internal/gateway/helpers.go`)
   - TeamManager created in gateway
   - Injected into ChannelManager
   - Accessible to all channels

### Message Flow

1. User sends message with @mentions
2. `handleMessage` detects mentions via `extractMentions()`
3. `handleCollaborativeMessage` is called
4. Session is retrieved/created
5. Message added to context
6. For each mentioned role:
   - Agent status updated to "thinking"
   - Full context + prompt built
   - `TeamManager.ExecuteTaskWithRole()` called (async)
   - Response formatted with IRC-style prefix
   - Sent to Telegram
   - Agent status updated to "idle"

## Advanced Features (Future)

### Commands (Optional)

Future enhancement: Add Telegram commands for session management:

- `/start_team <team-id>` - Start collaborative session with specific team
- `/who` - List active agents in current session
- `/status` - Show session status and agent states
- `/context` - Display conversation history
- `/reset` - Clear session and start fresh

### Auto-join Rules

Future enhancement: Agents automatically join based on keywords:

```json
{
  "collaborative_chat": {
    "auto_join_rules": [
      {
        "role": "architect",
        "keywords": ["design", "architecture", "structure"]
      },
      {
        "role": "tester",
        "keywords": ["test", "bug", "qa"]
      }
    ]
  }
}
```

### Agent-to-Agent Communication

Future enhancement: Allow agents to @mention each other:

```
[abc123] 🏗️ ARCHITECT: @developer Can you implement this design?
[abc123] 💻 DEVELOPER: @architect Yes, I'll start with the auth module...
```

## Troubleshooting

### Agents not responding

1. Check team configuration is loaded:
   ```bash
   picoclaw team list
   ```

2. Verify collaborative chat is enabled in config:
   ```json
   "collaborative_chat": {
     "enabled": true,
     "default_team_id": "dev-team"
   }
   ```

3. Check logs for errors:
   ```bash
   tail -f ~/.picoclaw/logs/picoclaw.log
   ```

### Wrong team being used

Ensure `default_team_id` matches your team configuration's `team_id` field.

### Context not maintained

Check `max_context_length` setting. Increase if needed:
```json
"max_context_length": 100
```

## Performance Considerations

- **Parallel execution**: Agents execute in parallel, so response time is limited by the slowest agent
- **Context size**: Larger context windows increase token usage and latency
- **Session cleanup**: Sessions are kept in memory; consider implementing TTL-based cleanup for production

## Security

- **Access control**: Use `allow_from` in Telegram config to restrict access
- **Tool permissions**: Configure role-specific tool access in team configuration
- **Rate limiting**: Consider implementing rate limits for collaborative chat

## Migration from Python Skill

If you were using the Python IRC gateway skill:

1. Remove the Python skill from `~/.picoclaw/workspace/skills/irc-gateway/`
2. Configure collaborative chat in `config.json` as shown above
3. Create team configuration file
4. Restart gateway

Benefits of native implementation:
- Single Telegram bot token (no need for multiple bots)
- Better integration with PicoClaw's team system
- Lower memory footprint
- Faster response times
- Native Go performance

## See Also

- [Multi-Agent Guide](MULTI_AGENT_GUIDE.md)
- [Team Tool Access](TEAM_TOOL_ACCESS.md)
- [Multi-Agent Model Selection](MULTI_AGENT_MODEL_SELECTION.md)


## Developer Guide

### Package Architecture

The collaborative chat system is built on the `pkg/collaborative/` package, which provides platform-agnostic functionality:

```
pkg/collaborative/
├── types.go          # Core types and Platform interface
├── session.go        # Session management and context
├── manager.go        # Collaborative manager and mention handling
├── mention.go        # @mention extraction utilities
├── formatting.go     # IRC-style message formatting
└── roster.go         # Team roster building
```

### Implementing for New Platforms

To add collaborative chat to a new messaging platform:

#### 1. Implement Platform Interface

```go
import "github.com/sipeed/picoclaw/pkg/collaborative"

type YourChannel struct {
    chatManager *collaborative.Manager
    // ... other fields
}

// Implement Platform interface
func (c *YourChannel) SendMessage(ctx context.Context, chatID string, content string) error {
    // Send message to your platform
    return nil
}

func (c *YourChannel) GetTeamManager() collaborative.TeamManager {
    // Return your team manager
    return c.teamManager
}

func (c *YourChannel) GetContext() context.Context {
    return c.ctx
}
```

#### 2. Initialize Manager

```go
func NewYourChannel(cfg *config.Config) *YourChannel {
    return &YourChannel{
        chatManager: collaborative.NewManager(),
        // ... initialize other fields
    }
}
```

#### 3. Handle Messages with Mentions

```go
func (c *YourChannel) handleMessage(ctx context.Context, chatID int64, content string) {
    // Extract mentions
    mentions := collaborative.ExtractMentions(content)
    
    if len(mentions) > 0 {
        // Delegate to collaborative manager
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

### API Reference

See [pkg/collaborative/README.md](../pkg/collaborative/README.md) for complete API documentation.

### Key Functions

- `collaborative.NewManager()` - Create manager
- `collaborative.ExtractMentions(text)` - Extract @mentions
- `collaborative.FormatMessage(sessionID, role, content)` - Format IRC-style
- `collaborative.GetRoleEmoji(role)` - Get role emoji
- `collaborative.BuildTeamRoster(teamInfo)` - Build team roster

### Session Management

Sessions are automatically created and managed:

```go
// Get or create session
session := manager.GetOrCreateSession(chatID, teamID, maxContext)

// Add message to context
session.AddMessage("user", content, mentions)

// Get formatted context
contextStr := session.GetContextAsString()

// Update agent status
session.UpdateAgentStatus("developer", "thinking")
```

### Response Extraction

The package handles complex response structures from team execution:

```go
// Automatically extracts clean text from:
// - Direct strings
// - Nested maps with "result" field
// - Arrays from parallel execution
// - Task metadata structures
responseText := extractResponseText(result)
```

### Error Handling

User-friendly error messages in Vietnamese:

```go
// Role not found
"❌ Role @designer không tồn tại trong team configuration.
💡 Các role có sẵn: @architect, @developer, @tester, @manager"

// Team manager unavailable
"⚠️ Collaborative chat is not properly configured."
```

## Related Documentation

- [Collaborative Package README](../pkg/collaborative/README.md) - API documentation
- [Quick Start Guide](COLLABORATIVE_CHAT_QUICKSTART.md) - 5-minute setup
- [Architecture](COLLABORATIVE_CHAT_ARCHITECTURE.md) - System design
- [Multi-Agent Guide](MULTI_AGENT_GUIDE.md) - Team configuration
- [PR #1138 Analysis](../ANALYSIS_PR1138_IRC_INTEGRATION.md) - IRC integration plan

## Changelog

### v0.2.0 (2026-03-05)

- **Refactored to platform-agnostic package** (`pkg/collaborative/`)
- Added Platform interface for multi-platform support
- Improved response text extraction
- Added agent-to-agent mention support
- Better error messages
- Prepared for IRC integration

### v0.1.0 (Initial Release)

- Basic @mention routing
- Parallel agent execution
- IRC-style formatting
- Session management
- Telegram integration


## How Agents Know Each Other

### Team Roster in Prompts

Every agent receives a **team roster** in their prompt, so they know who else is available in the team. This enables natural agent-to-agent communication.

#### Example Prompt Structure

When an agent is triggered, they receive:

```
=== Collaborative Chat Context ===
Session: chat51263350 | Team: dev-team
Started: 12:46:10
=== Conversation History ===

[12:46:10] USER: @developer can you help?
[12:46:15] 💻 DEVELOPER: Sure! Let me check...

=== Team Information ===
Team: Development Team
Members:
  • @architect - System architect and designer
  • @developer - Software developer
  • @tester - QA and testing specialist
  • @manager - Project manager
  • @designer - UI/UX designer
  • @devops - DevOps engineer

User message: Can you review this code?

You are @developer. Respond to the user's message considering the conversation history above.
You can mention other team members using @role format (e.g., @architect, @developer).
```

### How It Works

1. **Team Config Loading**: System reads team configuration from `~/.picoclaw/teams/*.json`
2. **Roster Building**: `BuildTeamRoster()` extracts role names and descriptions
3. **Prompt Injection**: Roster is included in every agent's prompt
4. **Agent Awareness**: Agents can see all available roles and their purposes

### Configuration

Team roster is automatically built from your team configuration:

```json
{
  "name": "dev-team",
  "description": "Development team",
  "roles": [
    {
      "name": "architect",
      "description": "System architect and designer",
      "model": "claude-3-5-sonnet-20241022"
    },
    {
      "name": "developer", 
      "description": "Software developer",
      "model": "claude-3-5-sonnet-20241022"
    },
    {
      "name": "tester",
      "description": "QA and testing specialist",
      "model": "claude-3-5-sonnet-20241022"
    }
  ]
}
```

### Benefits

1. **Natural Communication**: Agents know who to ask for help
2. **Role Awareness**: Each agent understands others' specializations
3. **Automatic Discovery**: No manual configuration needed
4. **Dynamic Updates**: Roster updates when team config changes

### Example Interactions

#### Agent Asking for Help

```
User: @developer implement login feature
Developer: Sure! I'll implement the backend. @designer can you create the UI mockups?
Designer: [automatically triggered] Of course! I'll design the login screen...
```

#### Agent Delegating Work

```
User: @manager we need to deploy the new feature
Manager: Got it! @developer please prepare the release. @devops can you handle deployment?
Developer: [triggered] Release is ready for deployment
DevOps: [triggered] I'll deploy to staging first for testing
```

#### Agent Collaboration

```
User: @architect design the new API
Architect: Here's the design... @developer can you implement it? @tester please prepare test cases
Developer: [triggered] I'll start implementation
Tester: [triggered] I'll create comprehensive test scenarios
```

### Fallback Roster

If team configuration cannot be loaded, a default roster is provided:

```
Team: Development Team
Members:
  • @architect - System architect and designer
  • @developer - Software developer
  • @tester - QA and testing specialist
  • @manager - Project manager
```

### Customization

You can customize role descriptions in your team config to help agents understand each other better:

```json
{
  "name": "specialist",
  "description": "Security specialist - handles authentication, authorization, and security audits",
  "model": "claude-3-5-sonnet-20241022"
}
```

The more descriptive your role descriptions, the better agents can collaborate!

### Technical Implementation

The roster building happens in `pkg/collaborative/roster.go`:

```go
// BuildTeamRoster creates a formatted string with team member information
func BuildTeamRoster(teamInfo any) string {
    // Extracts roles from team config
    // Formats as readable list
    // Returns formatted roster string
}
```

This is automatically called by `collaborative.Manager.HandleMentions()` and included in every agent's prompt.
