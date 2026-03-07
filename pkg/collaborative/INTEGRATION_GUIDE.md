# Collaborative Chat Integration Guide

**Version:** 1.0  
**Last Updated:** 2026-03-07

---

## Overview

This guide shows you how to add collaborative chat support to any PicoClaw channel. After refactoring, integration requires only ~50 lines of code per channel.

---

## Quick Start

### 1. Initialize ManagerV2

```go
import "github.com/sipeed/picoclaw/pkg/collaborative"

chatManager := collaborative.NewManagerV2WithConfig(&collaborative.Config{
    Enabled:             cfg.CollaborativeChat.Enabled,
    DefaultTeamID:       cfg.CollaborativeChat.DefaultTeamID,
    MaxContextLength:    cfg.CollaborativeChat.MaxContextLength,
    MentionQueueSize:    20,
    MentionRateLimit:    2 * time.Second,
    MentionMaxRetries:   3,
    MentionRetryBackoff: 1 * time.Second,
})
```

### 2. Implement Platform Interface

Your channel must implement 3 methods:

```go
type Platform interface {
    SendMessage(ctx context.Context, chatID string, content string) error
    GetTeamManager() TeamManager
    GetContext() context.Context
}
```

### 3. Handle Mentions

```go
mentions := collaborative.ExtractMentions(content)
if len(mentions) > 0 {
    return c.chatManager.HandleMentions(
        ctx,
        c,        // Platform interface
        chatID,
        teamID,
        content,
        mentions,
        sender,
        maxContextLength,
    )
}
```

### 4. Handle Commands

```go
if strings.HasPrefix(content, "/") {
    response, handled := collaborative.HandleCommand(
        c.chatManager,
        chatID,
        c,      // Platform interface
        content,
        nil,    // availableRoles (optional)
    )
    if handled {
        c.SendMessage(ctx, chatID, response)
        return nil
    }
}
```

---

## Platform Interface Details

### SendMessage

Send a formatted message to the chat platform.

```go
func (c *YourChannel) SendMessage(ctx context.Context, chatID string, content string) error {
    // Convert chatID string to platform-specific ID
    // Format content for your platform (HTML, Markdown, etc.)
    // Send via platform API
    return nil
}
```

**Example (Telegram):**
```go
func (c *TelegramChannel) SendMessage(ctx context.Context, chatID string, content string) error {
    cid, err := parseChatID(chatID)
    if err != nil {
        return err
    }

    htmlContent := markdownToTelegramHTML(content)
    tgMsg := tu.Message(tu.ID(cid), htmlContent)
    tgMsg.ParseMode = telego.ModeHTML

    _, err = c.bot.SendMessage(ctx, tgMsg)
    return err
}
```

### GetTeamManager

Return the team manager for accessing team configuration.

```go
func (c *YourChannel) GetTeamManager() collaborative.TeamManager {
    recorder := c.GetPlaceholderRecorder()
    if recorder == nil {
        return nil
    }

    // Get team manager from recorder
    // Wrap it to implement collaborative.TeamManager interface
    return &teamManagerAdapter{tm: tm}
}
```

**TeamManager Interface:**
```go
type TeamManager interface {
    ExecuteTaskWithRole(ctx context.Context, teamID, prompt, role string) (any, error)
    GetTeam(teamID string) (any, error)
}
```

### GetContext

Return the channel's context for lifecycle management.

```go
func (c *YourChannel) GetContext() context.Context {
    return c.ctx
}
```

---

## Complete Example: Discord Integration

Here's a complete example showing how to add collaborative chat to Discord:

```go
package discord

import (
    "context"
    "fmt"
    "strings"
    "time"

    "github.com/bwmarrin/discordgo"
    "github.com/sipeed/picoclaw/pkg/collaborative"
    "github.com/sipeed/picoclaw/pkg/config"
)

type DiscordChannel struct {
    session     *discordgo.Session
    config      *config.Config
    ctx         context.Context
    chatManager *collaborative.ManagerV2
}

func NewDiscordChannel(cfg *config.Config) (*DiscordChannel, error) {
    session, err := discordgo.New("Bot " + cfg.Channels.Discord.Token)
    if err != nil {
        return nil, err
    }

    ctx := context.Background()
    
    return &DiscordChannel{
        session: session,
        config:  cfg,
        ctx:     ctx,
        chatManager: collaborative.NewManagerV2WithConfig(&collaborative.Config{
            Enabled:             cfg.Channels.Discord.CollaborativeChat.Enabled,
            DefaultTeamID:       cfg.Channels.Discord.CollaborativeChat.DefaultTeamID,
            MaxContextLength:    cfg.Channels.Discord.CollaborativeChat.MaxContextLength,
            MentionQueueSize:    20,
            MentionRateLimit:    2 * time.Second,
            MentionMaxRetries:   3,
            MentionRetryBackoff: 1 * time.Second,
        }),
    }, nil
}

// Platform interface implementation

func (c *DiscordChannel) SendMessage(ctx context.Context, chatID string, content string) error {
    _, err := c.session.ChannelMessageSend(chatID, content)
    return err
}

func (c *DiscordChannel) GetTeamManager() collaborative.TeamManager {
    // Implementation similar to Telegram
    return nil // or actual team manager
}

func (c *DiscordChannel) GetContext() context.Context {
    return c.ctx
}

// Message handler

func (c *DiscordChannel) handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
    if m.Author.Bot {
        return
    }

    content := m.Content
    channelID := m.ChannelID

    // Check for collaborative commands
    if strings.HasPrefix(content, "/") {
        response, handled := collaborative.HandleCommand(
            c.chatManager,
            channelID,
            c,
            content,
            nil,
        )
        if handled {
            c.SendMessage(c.ctx, channelID, response)
            return
        }
    }

    // Check for mentions
    mentions := collaborative.ExtractMentions(content)
    if len(mentions) > 0 {
        collabCfg := c.config.Channels.Discord.CollaborativeChat
        teamID := collabCfg.DefaultTeamID

        sender := bus.SenderInfo{
            Platform:    "discord",
            PlatformID:  m.Author.ID,
            Username:    m.Author.Username,
            DisplayName: m.Author.Username,
        }

        c.chatManager.HandleMentions(
            c.ctx,
            c,
            channelID,
            teamID,
            content,
            mentions,
            sender,
            collabCfg.MaxContextLength,
        )
        return
    }

    // Handle regular message...
}
```

**Total: ~50 lines for collaborative integration**

---

## Available Helper Functions

### Formatting Helpers

```go
// Get emoji for a role
emoji := collaborative.GetRoleEmoji("developer") // 💻

// Format agent message
msg := collaborative.FormatAgentMessage("developer", "Task done", "chat12345678")
// Output: [chat12345678] 💻 DEVELOPER: Task done

// Generate session ID
sessionID := collaborative.GenerateSessionID(chatID)

// Format time since
timeStr := collaborative.FormatTimeSince(time.Now().Add(-5 * time.Minute))
// Output: "5m ago"

// Get status emoji
emoji := collaborative.GetStatusEmoji("thinking") // 🤔

// Format session context for LLM
context := collaborative.FormatSessionContext(session)
```

### Command Handlers

```go
// Check if message is a command
isCmd := collaborative.IsCollaborativeCommand("/who") // true

// Handle command
response, handled := collaborative.HandleCommand(
    manager,
    chatID,
    platform,
    "/who",
    []string{"developer", "tester"}, // optional available roles
)
```

---

## Configuration

Add to your channel config:

```yaml
channels:
  your_channel:
    collaborative_chat:
      enabled: true
      default_team_id: "dev-team"
      max_context_length: 50
```

Or in Go:

```go
type YourChannelConfig struct {
    CollaborativeChat struct {
        Enabled          bool   `yaml:"enabled"`
        DefaultTeamID    string `yaml:"default_team_id"`
        MaxContextLength int    `yaml:"max_context_length"`
    } `yaml:"collaborative_chat"`
}
```

---

## Testing

### Unit Tests

Test your Platform interface implementation:

```go
func TestSendMessage(t *testing.T) {
    channel := setupTestChannel()
    err := channel.SendMessage(context.Background(), "test-chat", "Hello")
    if err != nil {
        t.Errorf("SendMessage failed: %v", err)
    }
}
```

### Integration Tests

Test collaborative flow:

```go
func TestCollaborativeMentions(t *testing.T) {
    channel := setupTestChannel()
    
    // Simulate message with mentions
    content := "Hey @developer can you help?"
    mentions := collaborative.ExtractMentions(content)
    
    err := channel.chatManager.HandleMentions(
        context.Background(),
        channel,
        "test-chat",
        "test-team",
        content,
        mentions,
        testSender,
        50,
    )
    
    if err != nil {
        t.Errorf("HandleMentions failed: %v", err)
    }
}
```

---

## Checklist for New Channel

- [ ] Initialize `ManagerV2` with config
- [ ] Implement `SendMessage(ctx, chatID, content) error`
- [ ] Implement `GetTeamManager() TeamManager`
- [ ] Implement `GetContext() context.Context`
- [ ] Add mention detection: `collaborative.ExtractMentions(content)`
- [ ] Add mention handling: `chatManager.HandleMentions(...)`
- [ ] Add command detection: `strings.HasPrefix(content, "/")`
- [ ] Add command handling: `collaborative.HandleCommand(...)`
- [ ] Add config section for `collaborative_chat`
- [ ] Test with `/who` and `/help` commands
- [ ] Test with `@developer` mentions
- [ ] Test cascading mentions (agent mentions another agent)
- [ ] Update channel documentation

---

## Troubleshooting

### Mentions Not Working

1. Check if `CollaborativeChat.Enabled = true` in config
2. Verify `DefaultTeamID` is set and team exists
3. Check `ExtractMentions()` is finding mentions correctly
4. Verify Platform interface is implemented correctly

### Commands Not Responding

1. Check if message starts with `/`
2. Verify `HandleCommand()` is called before regular message handling
3. Check `SendMessage()` is working correctly

### Agent Not Responding

1. Check team configuration has the mentioned role
2. Verify `GetTeamManager()` returns valid team manager
3. Check queue is processing (see logs)
4. Verify no errors in `HandleMentions()`

---

## Architecture

```
User Message
    ↓
[Extract Mentions] ← collaborative.ExtractMentions()
    ↓
[Handle Mentions] ← chatManager.HandleMentions()
    ↓
[Queue Manager] ← Queues mention requests
    ↓
[Worker Pool] ← Processes mentions with rate limiting
    ↓
[Team Manager] ← Executes agent with role
    ↓
[Platform.SendMessage()] ← Sends response back
```

---

## Best Practices

1. **Always check config.Enabled** before processing collaborative features
2. **Use shared helpers** from `pkg/collaborative/formatting.go` and `commands.go`
3. **Handle errors gracefully** - don't crash on collaborative failures
4. **Log important events** - mention processing, queue status, errors
5. **Test cascading** - ensure agents can mention each other
6. **Respect rate limits** - use configured `MentionRateLimit`
7. **Clean up sessions** - manager handles TTL automatically

---

## Performance Considerations

- **Queue Size**: Default 20, increase for high-traffic channels
- **Rate Limit**: Default 2s, adjust based on platform limits
- **Max Context**: Default 50 messages, balance memory vs context quality
- **Worker Pool**: Automatically sized, handles concurrent mentions
- **Session TTL**: 1 hour, sessions auto-expire

---

## Support

For questions or issues:
1. Check existing channel implementations (Telegram is reference)
2. Review test files in `pkg/collaborative/*_test.go`
3. Check logs for error messages
4. Open issue on GitHub

---

**Next Steps:**
1. Choose a channel to integrate
2. Follow the checklist above
3. Test with `/who`, `/help`, and `@mentions`
4. Submit PR with your integration

Good luck! 🚀
