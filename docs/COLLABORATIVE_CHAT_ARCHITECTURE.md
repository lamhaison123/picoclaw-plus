# Collaborative Chat - Cơ Chế Hoạt Động

## Tổng Quan

IRC-style collaborative chat cho phép nhiều AI agents tham gia vào cùng một cuộc trò chuyện Telegram, sử dụng @mentions để định tuyến và duy trì context chung.

## Kiến Trúc Tổng Thể

```
┌─────────────────────────────────────────────────────────────────┐
│                         User (Telegram)                          │
│                  "@architect @developer help me"                 │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                    TelegramChannel.handleMessage()               │
│  1. Nhận message từ Telegram                                    │
│  2. Kiểm tra collaborative_chat.enabled                          │
│  3. Extract @mentions từ message                                 │
│  4. Nếu có mentions → gọi handleCollaborativeMessage()           │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│              TelegramChannel.handleCollaborativeMessage()        │
│  1. Get/Create CollaborativeChatSession cho chat này            │
│  2. Add message vào conversation context                         │
│  3. Get TeamManager từ ChannelManager                            │
│  4. Build prompt với full context                                │
│  5. For each @mention (parallel goroutines):                     │
│     - Update agent status = "thinking"                           │
│     - Execute TeamManager.ExecuteTaskWithRole()                  │
│     - Format response với IRC-style                              │
│     - Send về Telegram                                           │
│     - Update agent status = "idle"                               │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                      TeamManager.ExecuteTaskWithRole()           │
│  1. Tìm team theo team_id                                        │
│  2. Validate role tồn tại trong team                             │
│  3. Create CoordinatorAgent                                      │
│  4. Create Task với role và prompt                               │
│  5. Execute task theo pattern (parallel/sequential/hierarchical) │
│  6. Return kết quả                                               │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Agent Execution (LLM Call)                    │
│  1. Agent nhận prompt với full context                           │
│  2. LLM xử lý và generate response                               │
│  3. Response được return về                                      │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Format & Send Response                        │
│  Format: [session-id] emoji ROLE: response                       │
│  Example: [abc123] 🏗️ ARCHITECT: I recommend...                │
└─────────────────────────────────────────────────────────────────┘
```

## Chi Tiết Từng Bước

### 1. Khởi Tạo (Gateway Startup)

```go
// cmd/picoclaw/internal/gateway/helpers.go

func gatewayCmd() {
    // 1. Tạo AgentLoop với registry và provider
    agentLoop := agent.NewAgentLoop(cfg, msgBus, provider)
    registry := agentLoop.GetRegistry()
    
    // 2. Tạo TeamManager
    teamManager := team.NewTeamManager(registry, msgBus)
    teamManager.SetProvider(provider, cfg)
    
    // 3. Tạo Executor để chạy agents
    executor := team.NewDirectAgentExecutor(agentLoop)
    teamManager.SetAgentExecutor(executor)
    
    // 4. Setup team memory
    teamMemory := memory.NewTeamMemory(cfg.WorkspacePath())
    teamManager.SetTeamMemory(teamMemory)
    
    // 5. Tạo ChannelManager
    channelManager := channels.NewManager(cfg, msgBus, mediaStore)
    
    // 6. INJECT TeamManager vào ChannelManager
    channelManager.SetTeamManager(teamManager)
    
    // Bây giờ TelegramChannel có thể access TeamManager!
}
```

**Tại sao cần inject?**
- TelegramChannel không thể import trực tiếp `pkg/team` (circular dependency)
- Sử dụng interface `TeamManager` để decouple
- ChannelManager giữ reference và cung cấp cho channels

### 2. Nhận Message Từ Telegram

```go
// pkg/channels/telegram/telegram.go

func (c *TelegramChannel) handleMessage(ctx context.Context, message *telego.Message) error {
    // 1. Parse message content
    content := message.Text
    
    // 2. Kiểm tra collaborative chat enabled
    collabCfg := c.config.Channels.Telegram.CollaborativeChat
    if collabCfg.Enabled {
        // 3. Extract @mentions
        mentions := extractMentions(content)
        // mentions = ["architect", "developer"]
        
        if len(mentions) > 0 {
            // 4. Route to collaborative handler
            return c.handleCollaborativeMessage(ctx, chatID, content, mentions, sender)
        }
    }
    
    // 5. Nếu không có mentions, xử lý như message thường
    c.HandleMessage(...)
}
```

**extractMentions() hoạt động như thế nào?**

```go
// pkg/channels/telegram/collaborative_chat.go

func extractMentions(content string) []string {
    // Regex: @(\w+) - tìm tất cả @word
    mentionRegex := regexp.MustCompile(`@(\w+)`)
    matches := mentionRegex.FindAllStringSubmatch(content, -1)
    
    mentions := []string{}
    seen := make(map[string]bool)
    
    for _, match := range matches {
        role := match[1]  // "architect", "developer"
        if !seen[role] {
            mentions = append(mentions, role)
            seen[role] = true
        }
    }
    
    return mentions
}
```

**Ví dụ:**
- Input: `"@architect @developer @architect help me"`
- Output: `["architect", "developer"]` (unique)

### 3. Session Management

```go
func (c *TelegramChannel) handleCollaborativeMessage(...) error {
    // 1. Get hoặc create session cho chat này
    session := c.chatManager.GetOrCreateSession(
        chatID,                              // Telegram chat ID
        collabCfg.DefaultTeamID,            // "dev-team"
        collabCfg.MaxContextLength          // 50 messages
    )
    
    // 2. Add message vào context
    session.AddMessage("user", content, mentions)
    // Context bây giờ có: [user: "@architect help", ...]
}
```

**CollaborativeChatSession Structure:**

```go
type CollaborativeChatSession struct {
    SessionID    string              // "abc123" (hash của chatID)
    ChatID       int64               // Telegram chat ID
    TeamID       string              // "dev-team"
    ActiveAgents map[string]*AgentState
    Context      []ChatMessage       // Conversation history
    MaxContext   int                 // 50 messages max
}

type ChatMessage struct {
    ID        string
    Author    string      // "user" hoặc "architect", "developer"
    Content   string
    Timestamp time.Time
    Mentions  []string    // ["architect", "developer"]
}
```

**Context được lưu như thế nào?**

```go
func (s *CollaborativeChatSession) AddMessage(author, content string, mentions []string) {
    msg := ChatMessage{
        ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
        Author:    author,
        Content:   content,
        Timestamp: time.Now(),
        Mentions:  mentions,
    }
    
    s.Context = append(s.Context, msg)
    
    // Trim nếu quá dài
    if len(s.Context) > s.MaxContext {
        s.Context = s.Context[len(s.Context)-s.MaxContext:]
    }
}
```

### 4. Get TeamManager

```go
func (c *TelegramChannel) handleCollaborativeMessage(...) error {
    // 1. Get PlaceholderRecorder (thực ra là ChannelManager)
    recorder := c.GetPlaceholderRecorder()
    
    // 2. Type assert để get TeamManager
    type teamManagerGetter interface {
        GetTeamManager() channels.TeamManager
    }
    
    tmGetter := recorder.(teamManagerGetter)
    teamManager := tmGetter.GetTeamManager()
    
    // Bây giờ có TeamManager!
}
```

**Tại sao phức tạp thế?**
- TelegramChannel không biết về ChannelManager trực tiếp
- Sử dụng PlaceholderRecorder interface (đã có sẵn)
- Type assertion để access GetTeamManager()
- Tránh circular dependency

### 5. Build Context String

```go
func (s *CollaborativeChatSession) GetContextAsString() string {
    var sb strings.Builder
    sb.WriteString("=== Collaborative Chat Context ===\n")
    sb.WriteString(fmt.Sprintf("Session: %s | Team: %s\n", s.SessionID, s.TeamID))
    sb.WriteString("=== Conversation History ===\n\n")
    
    for _, msg := range s.Context {
        emoji := getRoleEmoji(msg.Author)
        sb.WriteString(fmt.Sprintf("[%s] %s%s: %s\n",
            msg.Timestamp.Format("15:04:05"),
            emoji,
            strings.ToUpper(msg.Author),
            msg.Content,
        ))
    }
    
    return sb.String()
}
```

**Output example:**
```
=== Collaborative Chat Context ===
Session: abc123 | Team: dev-team
=== Conversation History ===

[14:30:15] USER: @architect how to design auth?
[14:30:20] 🏗️ ARCHITECT: I recommend JWT tokens...
[14:30:25] USER: @developer can you implement that?
```

### 6. Parallel Agent Execution

```go
func (c *TelegramChannel) handleCollaborativeMessage(...) error {
    contextStr := session.GetContextAsString()
    
    // Execute mỗi role trong parallel goroutines
    for _, role := range mentions {
        // Update status
        session.UpdateAgentStatus(role, "thinking")
        
        // Build prompt
        prompt := fmt.Sprintf(
            "%s\n\nUser message: %s\n\nYou are @%s. Respond considering the conversation history.",
            contextStr, content, role
        )
        
        // Execute ASYNC (goroutine)
        go func(r string) {
            // Call TeamManager
            result, err := teamManager.ExecuteTaskWithRole(
                ctx,
                teamID,  // "dev-team"
                prompt,  // Full context + instruction
                r        // "architect" hoặc "developer"
            )
            
            if err != nil {
                // Send error message
                errorMsg := fmt.Sprintf("❌ @%s encountered an error: %v", r, err)
                c.bot.SendMessage(ctx, tu.Message(tu.ID(chatID), errorMsg))
                return
            }
            
            // Format response
            responseStr := fmt.Sprintf("%v", result)
            session.AddMessage(r, responseStr, nil)
            
            // IRC-style formatting
            emoji := getRoleEmoji(r)
            formattedMsg := fmt.Sprintf("[%s] %s %s: %s",
                session.SessionID,
                emoji,
                strings.ToUpper(r),
                responseStr
            )
            
            // Send to Telegram
            c.bot.SendMessage(ctx, tu.Message(tu.ID(chatID), formattedMsg))
            
            session.UpdateAgentStatus(r, "idle")
        }(role)
    }
    
    return nil
}
```

**Tại sao dùng goroutines?**
- Agents execute đồng thời (parallel)
- User nhận responses nhanh hơn
- Không block main thread
- Response time = slowest agent (không phải tổng)

### 7. TeamManager Execution

```go
// pkg/team/manager.go

func (tm *TeamManager) ExecuteTaskWithRole(ctx context.Context, teamID, taskDescription, requiredRole string) (any, error) {
    // 1. Get team
    team := tm.teams[teamID]
    
    // 2. Validate role exists
    roleExists := false
    for _, roleConfig := range team.Config.Roles {
        if roleConfig.Name == requiredRole {
            roleExists = true
            break
        }
    }
    
    // 3. Create coordinator
    coordinator := NewCoordinatorAgent(
        fmt.Sprintf("coordinator-%s", teamID),
        teamID,
        team,
        team.Pattern,  // "parallel"
        tm.bus,
        router,
        ctx,
    )
    
    coordinator.Executor = tm.agentExecutor
    
    // 4. Create task
    task := NewTask(taskDescription, requiredRole, nil)
    
    // 5. Execute based on pattern
    var result any
    switch team.Pattern {
    case PatternParallel:
        result, err = coordinator.ExecuteParallel(ctx, []*Task{task})
    case PatternSequential:
        result, err = coordinator.ExecuteSequential(ctx, []*Task{task})
    case PatternHierarchical:
        result, err = coordinator.ExecuteHierarchical(ctx, task)
    }
    
    return result, err
}
```

### 8. Agent Execution (LLM Call)

```go
// pkg/team/executor.go

func (e *DirectAgentExecutor) Execute(ctx context.Context, agentID, prompt string) (map[string]any, error) {
    // 1. Get agent instance từ registry
    instance := e.agentLoop.GetRegistry().GetAgent(agentID)
    
    // 2. Execute với LLM
    response, err := e.agentLoop.ProcessMessage(
        ctx,
        prompt,      // Full context + instruction
        "internal",  // channel
        agentID,     // chat ID
    )
    
    // 3. Return result
    return map[string]any{
        "status":    "completed",
        "agent_id":  agentID,
        "role":      role,
        "result":    response,  // LLM response
        "timestamp": time.Now(),
    }, nil
}
```

**LLM nhận gì?**
```
=== Collaborative Chat Context ===
Session: abc123 | Team: dev-team
=== Conversation History ===

[14:30:15] USER: @architect how to design auth?

User message: @architect how to design auth?

You are @architect. Respond considering the conversation history.
```

**LLM trả về:**
```
I recommend using JWT tokens with the following approach:
1. User login with credentials
2. Server validates and issues JWT
3. Client stores JWT in localStorage
4. Include JWT in Authorization header for API calls
...
```

### 9. Format & Send Response

```go
// Format IRC-style
emoji := getRoleEmoji("architect")  // 🏗️
formattedMsg := fmt.Sprintf("[%s] %s %s: %s",
    "abc123",           // session ID
    emoji,              // 🏗️
    "ARCHITECT",        // role uppercase
    response            // LLM response
)

// Send to Telegram
c.bot.SendMessage(ctx, tu.Message(tu.ID(chatID), formattedMsg))
```

**Output:**
```
[abc123] 🏗️ ARCHITECT: I recommend using JWT tokens with the following approach:
1. User login with credentials
2. Server validates and issues JWT
...
```

## Flow Diagram Chi Tiết

```
User sends: "@architect @developer help me"
    │
    ▼
┌─────────────────────────────────────────┐
│ TelegramChannel.handleMessage()         │
│ - Parse message                         │
│ - Check collaborative_chat.enabled      │
│ - Extract mentions: ["architect", "dev"]│
└────────────┬────────────────────────────┘
             │
             ▼
┌─────────────────────────────────────────┐
│ handleCollaborativeMessage()            │
│ - Get session (abc123)                  │
│ - Add to context                        │
│ - Get TeamManager                       │
└────────────┬────────────────────────────┘
             │
             ├──────────────┬──────────────┐
             ▼              ▼              ▼
    ┌────────────┐  ┌────────────┐  (parallel)
    │ Goroutine  │  │ Goroutine  │
    │ @architect │  │ @developer │
    └──────┬─────┘  └──────┬─────┘
           │                │
           ▼                ▼
    ┌────────────┐  ┌────────────┐
    │ TeamMgr    │  │ TeamMgr    │
    │ Execute    │  │ Execute    │
    │ architect  │  │ developer  │
    └──────┬─────┘  └──────┬─────┘
           │                │
           ▼                ▼
    ┌────────────┐  ┌────────────┐
    │ LLM Call   │  │ LLM Call   │
    │ (Claude)   │  │ (Claude)   │
    └──────┬─────┘  └──────┬─────┘
           │                │
           ▼                ▼
    ┌────────────┐  ┌────────────┐
    │ Format     │  │ Format     │
    │ [abc] 🏗️  │  │ [abc] 💻  │
    └──────┬─────┘  └──────┬─────┘
           │                │
           ▼                ▼
    ┌────────────────────────────┐
    │   Send to Telegram         │
    │   (responses arrive async) │
    └────────────────────────────┘
```

## Các Tính Năng Quan Trọng

### 1. Context Window Management

```go
// Tự động trim khi quá dài
if len(s.Context) > s.MaxContext {
    s.Context = s.Context[len(s.Context)-s.MaxContext:]
}
```

**Tại sao cần?**
- Tránh context quá dài → token limit
- Giữ conversation relevant
- Configurable: `max_context_length: 50`

### 2. Agent Status Tracking

```go
type AgentState struct {
    Role         string
    Status       string  // "idle", "thinking", "busy", "error"
    LastSeen     time.Time
    MessageCount int
}

session.UpdateAgentStatus("architect", "thinking")
// ... execute ...
session.UpdateAgentStatus("architect", "idle")
```

**Dùng để:**
- Track agent activity
- Debug issues
- Future: show status in UI

### 3. Session ID Generation

```go
func generateSessionID(chatID int64) string {
    hash := fmt.Sprintf("%x", chatID)
    if len(hash) > 6 {
        hash = hash[:6]
    }
    return hash
}
```

**Ví dụ:**
- chatID: `123456789`
- sessionID: `75bcd1` (first 6 chars of hex)

### 4. Role Emoji Mapping

```go
var roleEmojis = map[string]string{
    "architect": "🏗️",
    "developer": "💻",
    "tester":    "🧪",
    "manager":   "📋",
    "devops":    "⚙️",
    "designer":  "🎨",
}

func getRoleEmoji(role string) string {
    if emoji, ok := roleEmojis[role]; ok {
        return emoji
    }
    return "🤖"  // default
}
```

## Performance Characteristics

### Memory Usage

```
Per Session:
- SessionID: 8 bytes
- ChatID: 8 bytes
- TeamID: ~20 bytes
- Context (50 messages × ~500 bytes): ~25KB
- ActiveAgents (6 agents × 100 bytes): ~600 bytes

Total per session: ~26KB
```

**10 concurrent chats = ~260KB** (rất nhẹ!)

### Response Time

```
Sequential (old):
Agent1 (3s) + Agent2 (3s) + Agent3 (3s) = 9s total

Parallel (new):
max(Agent1: 3s, Agent2: 3s, Agent3: 3s) = 3s total
```

**3x faster với parallel execution!**

### Scalability

- Mỗi chat có session riêng (isolated)
- Goroutines lightweight (2KB stack)
- No shared state between chats
- Can handle 100+ concurrent chats easily

## Error Handling

```go
// 1. Team not found
if teamID == "" {
    logger.WarnC("telegram", "No default_team_id configured")
    return nil
}

// 2. TeamManager not available
if teamManager == nil {
    logger.WarnC("telegram", "Team manager not available")
    return nil
}

// 3. Agent execution failed
if err != nil {
    errorMsg := fmt.Sprintf("❌ @%s encountered an error: %v", role, err)
    c.bot.SendMessage(ctx, tu.Message(tu.ID(chatID), errorMsg))
    session.UpdateAgentStatus(role, "error")
    return
}
```

## Security Considerations

### 1. Access Control

```go
// Check allowlist
if !c.IsAllowedSender(sender) {
    logger.DebugCF("telegram", "Message rejected by allowlist")
    return nil
}
```

### 2. Tool Permissions

```json
{
  "roles": [
    {
      "name": "architect",
      "tools": ["*"]  // All tools
    },
    {
      "name": "tester",
      "tools": ["file_read", "exec_safe"]  // Limited
    }
  ]
}
```

### 3. Rate Limiting

- Telegram bot API: 30 messages/second
- Context window: 50 messages max
- No additional rate limiting needed

## Debugging Tips

### 1. Enable Debug Logging

```bash
picoclaw gateway --debug
```

### 2. Check Logs

```bash
tail -f ~/.picoclaw/logs/picoclaw.log | grep collaborative
```

### 3. Verify Team Config

```bash
picoclaw team list
picoclaw team status dev-team
```

### 4. Test Mentions

```
# Single mention
@architect help

# Multiple mentions
@architect @developer @tester review this

# Invalid mention (ignored)
@nonexistent help
```

## Kết Luận

Collaborative chat hoạt động thông qua:

1. **Mention Detection**: Regex extract @mentions
2. **Session Management**: Context per chat
3. **TeamManager Integration**: Interface-based injection
4. **Parallel Execution**: Goroutines for speed
5. **IRC Formatting**: Clear, readable output
6. **Context Sharing**: Full history to all agents

Kiến trúc này cho phép:
- ✅ Multiple agents in one chat
- ✅ Shared conversation context
- ✅ Fast parallel responses
- ✅ Clean separation of concerns
- ✅ Easy to extend and maintain
