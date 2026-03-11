# PicoClaw Data Flow Documentation

## 1. Message Processing Flow (Chi tiết)

### Step-by-Step Flow

```
┌─────────────────────────────────────────────────────────────┐
│ 1. USER INPUT                                               │
│    User sends: "Build a REST API for authentication"       │
│    Platform: Telegram                                       │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 2. CHANNEL LAYER                                            │
│    - Telegram channel receives update                       │
│    - Parse message content                                  │
│    - Extract sender info                                    │
│    - Normalize to InboundMessage                           │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 3. MESSAGE BUS - INBOUND                                    │
│    msg := InboundMessage{                                   │
│        Channel: "telegram",                                 │
│        SenderID: "123456",                                  │
│        ChatID: "789",                                       │
│        Content: "Build a REST API...",                      │
│        SessionKey: "",                                      │
│    }                                                        │
│    bus.PublishInbound(ctx, msg)                            │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 4. AGENT LOOP - CONSUME                                     │
│    msg, ok := bus.ConsumeInbound(ctx)                      │
│    if !ok { continue }                                      │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 5. ROUTING                                                  │
│    route := registry.ResolveRoute(RouteInput{              │
│        Channel: "telegram",                                 │
│        AccountID: "123456",                                 │
│    })                                                       │
│    agent := registry.GetAgent(route.AgentID)               │
│    sessionKey := route.SessionKey                          │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 6. LOAD SESSION HISTORY                                     │
│    history := agent.Sessions.GetHistory(sessionKey)        │
│    summary := agent.Sessions.GetSummary(sessionKey)        │
│    // history = [                                           │
│    //   {role: "user", content: "previous message"},       │
│    //   {role: "assistant", content: "previous response"}  │
│    // ]                                                     │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 7. SEARCH VECTOR MEMORY (if enabled)                       │
│    if memoryEnabled {                                       │
│        embedding := embeddingService.Generate(query)        │
│        results := vectorStore.Search(embedding, topK=3)     │
│        memoryContext := formatResults(results)              │
│    }                                                        │
│    // memoryContext = "[Relevant past context:]            │
│    //   1. User asked about auth before (0.92)             │
│    //   2. Discussed JWT tokens (0.87)"                    │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 8. BUILD CONTEXT                                            │
│    messages := contextBuilder.BuildMessages(                │
│        history,                                             │
│        summary,                                             │
│        userMessage,                                         │
│        media,                                               │
│        channel,                                             │
│        chatID                                               │
│    )                                                        │
│    if memoryContext != "" {                                 │
│        injectMemoryContext(messages, memoryContext)         │
│    }                                                        │
│    // messages = [                                          │
│    //   {role: "system", content: "You are..."},           │
│    //   {role: "user", content: "previous"},               │
│    //   {role: "assistant", content: "previous response"}, │
│    //   {role: "user", content: "Build a REST API..."}     │
│    // ]                                                     │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 9. LLM CALL WITH TOOLS                                      │
│    response := provider.Chat(messages, tools)               │
│    // LLM thinks: "This is complex, I need dev-team"       │
│    // response = {                                          │
│    //   content: "",                                        │
│    //   tool_calls: [{                                      │
│    //     name: "delegate_to_team",                         │
│    //     args: {                                           │
│    //       team_id: "dev-team",                            │
│    //       task: "Build REST API for auth"                 │
│    //     }                                                 │
│    //   }]                                                  │
│    // }                                                     │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 10. TOOL EXECUTION                                          │
│     for each tool_call:                                     │
│         tool := toolRegistry.Get(tool_call.name)            │
│         result := tool.Execute(ctx, tool_call.args)         │
│         toolResults.append(result)                          │
│     // result = {                                           │
│     //   success: true,                                     │
│     //   content: "API implemented with 5 endpoints..."     │
│     // }                                                    │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 11. FEED RESULTS BACK TO LLM                                │
│     messages.append({                                       │
│         role: "assistant",                                  │
│         tool_calls: [...]                                   │
│     })                                                      │
│     messages.append({                                       │
│         role: "tool",                                       │
│         content: result.content                             │
│     })                                                      │
│     response := provider.Chat(messages, tools)              │
│     // LLM: "Great! The team completed the API..."         │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 12. FINAL RESPONSE                                          │
│     finalContent := response.content                        │
│     // "✓ Dev-team completed the API implementation        │
│     //  with JWT authentication and 5 endpoints..."        │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 13. STORE IN VECTOR MEMORY (async)                         │
│     go func() {                                             │
│         combined := "User: " + userMsg + "\n" +             │
│                     "Assistant: " + finalContent            │
│         embedding := embeddingService.Generate(combined)    │
│         vector := Vector{                                   │
│             ID: sessionKey + ":" + timestamp,               │
│             Embedding: embedding,                           │
│             Metadata: {                                     │
│                 session: sessionKey,                        │
│                 content: finalContent,                      │
│                 timestamp: now                              │
│             }                                               │
│         }                                                   │
│         vectorStore.Upsert([vector])                        │
│     }()                                                     │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 14. SAVE SESSION                                            │
│     agent.Sessions.AddMessage(sessionKey, "user", userMsg)  │
│     agent.Sessions.AddMessage(sessionKey, "assistant",      │
│                               finalContent)                 │
│     agent.Sessions.Save(sessionKey)                         │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 15. PUBLISH OUTBOUND                                        │
│     bus.PublishOutbound(ctx, OutboundMessage{              │
│         Channel: "telegram",                                │
│         ChatID: "789",                                      │
│         Content: finalContent                               │
│     })                                                      │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 16. CHANNEL SENDS                                           │
│     telegram.Send(chatID, content)                          │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 17. USER RECEIVES                                           │
│     User sees: "✓ Dev-team completed the API..."           │
└─────────────────────────────────────────────────────────────┘
```

---

## 2. Team Delegation Flow (Chi tiết)

### Scenario: Agent delegates task to dev-team

```
┌─────────────────────────────────────────────────────────────┐
│ 1. AGENT DECIDES TO DELEGATE                                │
│    LLM: "This task needs multiple skills, I'll use team"   │
│    Tool call: delegate_to_team({                            │
│        team_id: "dev-team",                                 │
│        task: "Build REST API for authentication"            │
│    })                                                       │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 2. TOOL EXECUTION                                           │
│    tool := TeamDelegationTool                               │
│    result := tool.Execute(ctx, args)                        │
│    // Calls teamManager.ExecuteTask()                       │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 3. TEAM MANAGER - GET TEAM                                  │
│    team := teamManager.GetTeam("dev-team")                  │
│    // team = {                                              │
│    //   id: "dev-team",                                     │
│    //   pattern: "hierarchical",                            │
│    //   roles: [architect, developer, tester, manager]      │
│    // }                                                     │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 4. INTELLIGENT ROUTING                                      │
│    router := NewIntelligentRouter(agentExecutor)            │
│    role := router.DetermineRole(task, team.Roles)          │
│    // Analyzes task: "Build REST API"                      │
│    // Keywords: "build", "API", "authentication"            │
│    // Best match: "developer" role                          │
│    // role = "developer"                                    │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 5. CREATE COORDINATOR                                       │
│    coordinator := NewCoordinatorAgent(                      │
│        id: "coordinator-dev-team",                          │
│        teamID: "dev-team",                                  │
│        team: team,                                          │
│        pattern: "hierarchical"                              │
│    )                                                        │
│    coordinator.Executor = agentExecutor                     │
│    defer coordinator.Shutdown()                             │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 6. CREATE TASK                                              │
│    task := NewTask(                                         │
│        description: "Build REST API for authentication",    │
│        requiredRole: "developer",                           │
│        dependencies: nil                                    │
│    )                                                        │
│    task.ID = "task-12345"                                   │
│    task.Status = TaskStatusPending                          │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 7. EXECUTE HIERARCHICAL PATTERN                             │
│    result := coordinator.ExecuteHierarchical(ctx, task)     │
│    // Hierarchical: Coordinator delegates to role           │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 8. FIND AGENT FOR ROLE                                      │
│    agentID := "dev-team-developer"                          │
│    agent := team.Agents[agentID]                            │
│    // agent = {                                             │
│    //   id: "dev-team-developer",                           │
│    //   role: "developer",                                  │
│    //   model: "claude-sonnet-4",                           │
│    //   capabilities: ["coding", "implementation"]          │
│    // }                                                     │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 9. EXECUTE WITH AGENT                                       │
│    result := agentExecutor.Execute(                         │
│        ctx,                                                 │
│        agentID: "dev-team-developer",                       │
│        task: "Build REST API for authentication"            │
│    )                                                        │
│    // This calls AgentLoop.ProcessWithAgent()               │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 10. AGENT PROCESSES TASK                                    │
│     - Load developer's session                              │
│     - Build context with task description                   │
│     - Call LLM (claude-sonnet-4)                            │
│     - Execute tools (file_write, shell_exec, etc.)          │
│     - Generate code                                         │
│     - Return result                                         │
│     // result = "API implemented with:                      │
│     //   - 5 endpoints (login, register, refresh, etc.)     │
│     //   - JWT authentication                               │
│     //   - Password hashing with bcrypt                     │
│     //   - Rate limiting                                    │
│     //   - Input validation"                                │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 11. COORDINATOR COLLECTS RESULT                             │
│     task.Status = TaskStatusCompleted                       │
│     task.Result = result                                    │
│     task.CompletedAt = time.Now()                           │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 12. RETURN TO TEAM MANAGER                                  │
│     return result, nil                                      │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 13. TOOL RETURNS RESULT                                     │
│     toolResult := ToolResult{                               │
│         Success: true,                                      │
│         Content: result,                                    │
│         ForUser: "✓ Dev-team completed the task"            │
│     }                                                       │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 14. AGENT INCORPORATES RESULT                               │
│     - Feed tool result back to LLM                          │
│     - LLM formats response for user                         │
│     - "The development team has successfully..."            │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 15. USER RECEIVES FINAL RESPONSE                            │
│     "✓ Dev-team completed the API implementation with       │
│      JWT authentication and 5 endpoints. The code is        │
│      ready for testing."                                    │
└─────────────────────────────────────────────────────────────┘
```

---

## 3. Mention Cascading Flow (Chi tiết)

### Scenario: User mentions architect, who mentions developer, who mentions tester

```
┌─────────────────────────────────────────────────────────────┐
│ 1. USER SENDS MESSAGE WITH MENTION                          │
│    User (in Telegram): "@architect design the database"    │
│    Platform: Telegram collaborative chat                    │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 2. TELEGRAM CHANNEL DETECTS MENTION                         │
│    mentions := ExtractMentions(content)                     │
│    // mentions = ["architect"]                              │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 3. COLLABORATIVE MANAGER - HANDLE MENTIONS                  │
│    ManagerV2.HandleMentions(                                │
│        ctx,                                                 │
│        platform: telegram,                                  │
│        chatID: 123456,                                      │
│        teamID: "dev-team",                                  │
│        content: "@architect design the database",           │
│        mentions: ["architect"],                             │
│        sender: {username: "user123"},                       │
│        maxContext: 50                                       │
│    )                                                        │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 4. GET OR CREATE SESSION                                    │
│    session := GetOrCreateSession(chatID, teamID, 50)       │
│    // session = {                                           │
│    //   sessionID: "chat123456789",                         │
│    //   chatID: 123456,                                     │
│    //   teamID: "dev-team",                                 │
│    //   context: [],                                        │
│    //   mentionDepth: 0,                                    │
│    //   maxMentionDepth: 20                                 │
│    // }                                                     │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 5. CHECK DEPTH LIMIT                                        │
│    if session.MentionDepth >= maxMentionDepth {            │
│        return error("depth limit reached")                  │
│    }                                                        │
│    // OK: depth = 0, limit = 20                            │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 6. ADD USER MESSAGE TO CONTEXT                              │
│    session.AddMessage("user", content, mentions)            │
│    // session.context = [                                   │
│    //   {role: "user", content: "@architect design...",     │
│    //    mentions: ["architect"]}                           │
│    // ]                                                     │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 7. ENQUEUE MENTION REQUEST                                  │
│    req := MentionRequest{                                   │
│        Role: "architect",                                   │
│        Prompt: "@architect design the database",            │
│        SessionID: session.SessionID,                        │
│        ChatID: 123456,                                      │
│        TeamID: "dev-team",                                  │
│        Depth: 0,                                            │
│        MentionedBy: ""  // No one mentioned architect       │
│    }                                                        │
│    queueManager.Enqueue(req)                                │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 8. QUEUE PROCESSES REQUEST                                  │
│    queue := GetOrCreateQueue("architect")                   │
│    // Rate limit: wait 2s since last execution              │
│    // Execute with retry (max 3 attempts)                   │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 9. EXECUTE ARCHITECT AGENT                                  │
│    executeAgentAndCascadeWithError(                         │
│        platform, chatID, teamID, session,                   │
│        role: "architect",                                   │
│        triggerContent: "@architect design the database",    │
│        depth: 0,                                            │
│        mentionedBy: ""                                      │
│    )                                                        │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 10. CHECK IDEMPOTENCY                                       │
│     messageID := GenerateMessageID(chatID, sessionID,       │
│                                    role, content)           │
│     if dispatchTracker.IsDispatched(messageID) {            │
│         return // Skip duplicate                            │
│     }                                                       │
│     dispatchTracker.MarkDispatched(messageID)               │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 11. MARK AGENT IN CASCADE                                   │
│     session.IncrementMentionDepth()  // depth = 1           │
│     session.MarkAgentInCascade("architect")                 │
│     defer session.DecrementMentionDepth()                   │
│     defer session.UnmarkAgentInCascade("architect")         │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 12. BUILD PROMPT FOR ARCHITECT                              │
│     contextStr := session.GetContextAsString()              │
│     prompt := fmt.Sprintf(`                                 │
│         %s                                                  │
│         === Team Information ===                            │
│         %s                                                  │
│         User message: %s                                    │
│         You are @architect. Respond considering history.    │
│     `, contextStr, teamRoster, triggerContent)              │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 13. EXECUTE ARCHITECT AGENT                                 │
│     result := teamManager.ExecuteTaskWithRole(              │
│         ctx, "dev-team", prompt, "architect"                │
│     )                                                       │
│     // Architect thinks and responds:                       │
│     // "I'll design a schema with users, roles, sessions.   │
│     //  @developer can you implement the migrations?"       │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 14. ADD ARCHITECT RESPONSE TO CONTEXT                       │
│     session.AddMessage("architect", responseStr, nil)       │
│     session.UpdateAgentStatus("architect", "idle")          │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 15. SEND ARCHITECT RESPONSE                                 │
│     formattedMsg := FormatMessage(sessionID, "architect",   │
│                                   responseStr)              │
│     // "[chat123] 🏗️ ARCHITECT: I'll design..."            │
│     platform.SendMessage(ctx, chatID, formattedMsg)         │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 16. DETECT NEW MENTIONS IN RESPONSE                         │
│     mentionsInResponse := ExtractMentions(responseStr)      │
│     // mentionsInResponse = ["developer"]                   │
│     // Filter out self-mentions and ack-mentions            │
│     newMentions := []                                       │
│     for _, mentioned := range mentionsInResponse {          │
│         if mentioned == "architect" { continue }            │
│         if mentionedBy != "" && mentioned == mentionedBy {  │
│             continue  // Ack-loop prevention                │
│         }                                                   │
│         newMentions.append(mentioned)                       │
│     }                                                       │
│     // newMentions = ["developer"]                          │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 17. ENQUEUE CASCADED MENTION                                │
│     req := MentionRequest{                                  │
│         Role: "developer",                                  │
│         Prompt: responseStr,                                │
│         SessionID: session.SessionID,                       │
│         ChatID: 123456,                                     │
│         TeamID: "dev-team",                                 │
│         Depth: 1,  // Incremented                           │
│         MentionedBy: "architect"  // Track who mentioned    │
│     }                                                       │
│     queueManager.Enqueue(req)                               │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 18. EXECUTE DEVELOPER AGENT (depth=1)                       │
│     // Similar process as architect                         │
│     // Developer responds:                                  │
│     // "Sure! I'll create migrations. @tester verify"       │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 19. ENQUEUE TESTER MENTION (depth=2)                        │
│     req := MentionRequest{                                  │
│         Role: "tester",                                     │
│         Depth: 2,                                           │
│         MentionedBy: "developer"                            │
│     }                                                       │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 20. EXECUTE TESTER AGENT (depth=2)                          │
│     // Tester responds:                                     │
│     // "Will do! I'll write integration tests"              │
│     // No new mentions → cascade ends                       │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│ 21. ALL RESPONSES SENT TO USER                              │
│     User sees:                                              │
│     [chat123] 🏗️ ARCHITECT: I'll design...                 │
│     [chat123] 💻 DEVELOPER: Sure! I'll create...            │
│     [chat123] 🧪 TESTER: Will do! I'll write...             │
└─────────────────────────────────────────────────────────────┘
```

---

## Next Steps

Đọc thêm:
- [Configuration Guide](./CONFIGURATION_GUIDE.md)
- [Development Guide](./DEVELOPMENT_GUIDE.md)
