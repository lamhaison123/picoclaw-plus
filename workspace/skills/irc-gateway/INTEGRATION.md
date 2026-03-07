# IRC Gateway Integration with PicoClaw

## Overview

The IRC Gateway is a **lightweight routing layer** that extends PicoClaw's Telegram integration with IRC-style mention-based role routing. It integrates seamlessly with PicoClaw's existing multi-agent team system.

## Architecture Integration

```
┌─────────────────────────────────────────────────────────────┐
│                     Telegram Platform                        │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────┐
│              gateway.py (IRC Gateway)                        │
│  • Parse @mentions                                           │
│  • Session management                                        │
│  • IRC formatting                                            │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼ picoclaw team execute --role <role>
┌─────────────────────────────────────────────────────────────┐
│           PicoClaw Team System (pkg/team)                    │
│  • TeamManager                                               │
│  • Role-based routing                                        │
│  • Parallel/Sequential execution                             │
└──────────────────────┬──────────────────────────────────────┘
                       │
        ┌──────────────┼──────────────┬──────────────┐
        ▼              ▼              ▼              ▼
   ┌─────────┐   ┌─────────┐   ┌─────────┐   ┌─────────┐
   │Architect│   │Developer│   │ Tester  │   │ Manager │
   │ Agent   │   │ Agent   │   │ Agent   │   │ Agent   │
   └─────────┘   └─────────┘   └─────────┘   └─────────┘
        │              │              │              │
        └──────────────┴──────────────┴──────────────┘
                       │
                       ▼
              LLM Provider (OpenAI, Anthropic, etc.)
```

## Key Integration Points

### 1. Team System (`pkg/team`)

The gateway uses PicoClaw's team system for all agent coordination:

- **TeamManager**: Manages team lifecycle and agent registry
- **Role Configuration**: Defined in `~/.picoclaw/workspace/teams/irc-dev-team.json`
- **Parallel Execution**: Multiple roles can process tasks simultaneously
- **Shared Context**: Team members share conversation context

### 2. CLI Integration

Gateway invokes PicoClaw CLI commands:

```bash
picoclaw team execute irc-dev-team \
  -t "design a REST API" \
  --role architect \
  --format json
```

This ensures:
- ✅ Uses existing agent infrastructure
- ✅ Respects configuration in `~/.picoclaw/config.json`
- ✅ Leverages all PicoClaw features (tools, skills, MCP)
- ✅ Maintains session state

### 3. Session Management

Gateway maintains lightweight session tracking:

```python
sessions[chat_id] = {
    "session_id": "cmmd1234",
    "history": [
        {"user": "message", "roles": ["architect"], "timestamp": ...}
    ]
}
```

PicoClaw's session system handles the heavy lifting through the team execute command.

### 4. Message Bus (`pkg/bus`)

While the gateway doesn't directly use the message bus, the underlying team system does:

- **Task Delegation**: Roles communicate via message bus
- **Result Aggregation**: Coordinator collects responses
- **Event Notifications**: Team events are logged

## Comparison with Native Telegram Channel

| Aspect | IRC Gateway | Native `pkg/channels/telegram` |
|--------|-------------|-------------------------------|
| **Architecture** | Thin routing layer | Full channel implementation |
| **Agent Access** | Via team system | Direct agent loop |
| **Message Format** | IRC-style `[id] ROLE: msg` | Standard markdown |
| **Routing** | Explicit @mentions | Automatic |
| **Use Case** | Multi-role collaboration | General assistant |
| **Complexity** | ~300 lines Python | ~1000+ lines Go |
| **Integration** | CLI-based | Native Go packages |

## Data Flow

### Incoming Message

1. **Telegram** → User sends: `@architect design API`
2. **Gateway** → Parses mention: `architect`
3. **Gateway** → Executes: `picoclaw team execute irc-dev-team -t "design API" --role architect`
4. **PicoClaw** → TeamManager routes to architect agent
5. **Agent** → Processes with LLM provider
6. **PicoClaw** → Returns JSON result
7. **Gateway** → Formats: `[cmmd1234] 🏗️ ARCHITECT: Here's the design...`
8. **Telegram** → Sends formatted response

### Parallel Execution

1. **Telegram** → User sends: `@architect @developer build auth`
2. **Gateway** → Parses mentions: `[architect, developer]`
3. **Gateway** → Spawns 2 parallel tasks:
   - Task 1: `picoclaw team execute ... --role architect`
   - Task 2: `picoclaw team execute ... --role developer`
4. **PicoClaw** → Both agents process simultaneously
5. **Gateway** → Collects both responses
6. **Gateway** → Formats both with IRC style
7. **Telegram** → Sends both responses

## Configuration Files

### Team Configuration
**Location**: `~/.picoclaw/workspace/teams/irc-dev-team.json`

Defines:
- Role names and capabilities
- Tool access permissions
- Collaboration pattern (parallel/sequential/hierarchical)
- Timeout settings

### PicoClaw Configuration
**Location**: `~/.picoclaw/config.json`

Used by team system for:
- LLM provider settings
- Model selection per role
- Tool configurations
- Safety levels

### Gateway Configuration
**Location**: `.env` in repo root (optional)

Gateway-specific settings:
- Telegram bot token (optional if using PicoClaw config)
- PicoClaw binary path
- Team ID to use

**Token Loading Priority:**
1. Environment variable `TELEGRAM_BOT_TOKEN`
2. `.env` file in repo root
3. `~/.picoclaw/config.json` → `channels.telegram.token`

## Extending the Gateway

### Adding New Roles

1. **Update team config**:
```json
{
  "roles": [
    {
      "name": "devops",
      "capabilities": ["deployment", "infrastructure"],
      "tools": ["*"],
      "model": "gpt-4"
    }
  ]
}
```

2. **Update gateway ROLE_MAP**:
```python
ROLE_MAP = {
    "devops": {"emoji": "🚀", "description": "DevOps & deployment"}
}
```

3. **Recreate team**:
```bash
picoclaw team dissolve irc-dev-team
picoclaw team create ~/.picoclaw/workspace/teams/irc-dev-team.json
```

### Custom Formatting

Modify `format_response()` in `gateway.py`:

```python
async def format_response(role: str, response: str, session_id: str) -> str:
    emoji = ROLE_MAP[role]["emoji"]
    timestamp = datetime.now().strftime("%H:%M:%S")
    return f"[{timestamp}] [{session_id}] {emoji} {role.upper()}: {response}"
```

### Adding Commands

Add new Telegram command handlers:

```python
async def metrics_command(update: Update, context: ContextTypes.DEFAULT_TYPE):
    """Show team metrics"""
    result = subprocess.run(
        [PICOCLAW_BIN, "team", "metrics", TEAM_ID],
        capture_output=True, text=True
    )
    await update.message.reply_text(result.stdout)

# Register in main()
app.add_handler(CommandHandler("metrics", metrics_command))
```

## Performance Considerations

### Latency

- **Gateway overhead**: <50ms (parsing + formatting)
- **Team execution**: Depends on LLM provider (1-10s typical)
- **Parallel roles**: No additional latency (concurrent execution)

### Resource Usage

- **Gateway**: ~20MB RAM (Python + telegram bot)
- **PicoClaw team**: Managed by PicoClaw (minimal overhead)
- **Total**: Gateway + PicoClaw base (~30-40MB)

### Scalability

- **Concurrent chats**: Limited by Telegram bot API rate limits
- **Parallel roles**: Limited by PicoClaw team configuration
- **Recommended**: 1 gateway instance per bot token

## Security Considerations

### Authentication

Gateway inherits PicoClaw's security:
- Team tool access control
- Safety levels (strict/moderate/permissive)
- Allowlist configuration

### Best Practices

1. **Use environment variables** for sensitive data
2. **Restrict bot access** via Telegram allowlist
3. **Configure tool permissions** in team config
4. **Enable safety levels** in PicoClaw config
5. **Monitor logs** for suspicious activity

## Troubleshooting Integration Issues

### Team Not Found

```bash
# Check team exists
picoclaw team list

# Recreate if missing
picoclaw team create ~/.picoclaw/workspace/teams/irc-dev-team.json
```

### Role Execution Fails

```bash
# Test role directly
picoclaw team execute irc-dev-team -t "test message" --role architect

# Check team status
picoclaw team status irc-dev-team
```

### Gateway Can't Find Token

```bash
# Check if token exists in PicoClaw config
cat ~/.picoclaw/config.json | grep -A 5 telegram

# Or set in .env
echo "TELEGRAM_BOT_TOKEN=your_token" >> .env

# Or set environment variable
export TELEGRAM_BOT_TOKEN=your_token
```

## Future Enhancements

Potential improvements:

1. **Direct Go Integration**: Rewrite gateway in Go using `pkg/channels`
2. **WebSocket Support**: Real-time updates without polling
3. **Rich Media**: Support for images, files in role responses
4. **Consensus Voting**: Integrate team consensus mechanisms
5. **Task Queuing**: Queue management for busy roles
6. **Metrics Dashboard**: Real-time team performance metrics

## References

- [PicoClaw Multi-Agent Guide](../../../docs/MULTI_AGENT_GUIDE.md)
- [Team System Implementation](../../../pkg/team/)
- [Telegram Channel](../../../pkg/channels/telegram/)
- [Message Bus](../../../pkg/bus/)
