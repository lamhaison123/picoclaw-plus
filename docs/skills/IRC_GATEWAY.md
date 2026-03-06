# IRC Gateway Skill (DEPRECATED)

> **⚠️ DEPRECATED**: This Python skill has been replaced by native Go implementation.
> 
> **Please use the new [Collaborative Chat](../COLLABORATIVE_CHAT.md) feature instead!**
>
> The native implementation offers:
> - Single Telegram bot token (no need for multiple bots)
> - Better integration with PicoClaw's team system
> - Lower memory footprint and faster performance
> - Native Go reliability
>
> **Migration Guide**: See [Collaborative Chat Quick Start](../COLLABORATIVE_CHAT_QUICKSTART.md)

---

# Original Documentation (For Reference Only)

IRC-style Telegram bot gateway for mention-based role routing with PicoClaw's multi-agent team system.

## Overview

The IRC Gateway is a lightweight Python skill that extends PicoClaw's Telegram integration with IRC-style mention-based routing. It allows you to explicitly route messages to specific team roles using @mentions, enabling collaborative multi-agent workflows through a single Telegram bot.

## Features

- 🎯 **Team Integration**: Routes mentions directly to PicoClaw's team system
- 🏷️ **Mention-Based Routing**: Use `@architect`, `@developer`, `@tester`, `@manager`
- ⚡ **Parallel Processing**: Tag multiple roles for concurrent execution
- 💬 **IRC Formatting**: Responses prefixed with `[session-id] emoji ROLE: message`
- 📊 **Session Management**: Maintains conversation context per chat
- 🔍 **Status Monitoring**: Real-time role status tracking

## Architecture

```
Telegram Bot
    ↓
gateway.py (Mention Parser)
    ↓
picoclaw team execute --role <role>
    ↓
PicoClaw Team System (pkg/team)
    ↓
Role-specific Agent Instances
    ↓
IRC-formatted responses → Telegram
```

## Installation

> **Upgrading from v1.0.0?** See [Migration Guide](../../workspace/skills/irc-gateway/MIGRATION.md)

### Prerequisites

- PicoClaw installed and configured
- Python 3.8+
- Telegram bot token

### Quick Setup

```bash
# 1. Navigate to skill directory
cd workspace/skills/irc-gateway

# 2. Run setup script (Linux/Mac)
bash setup.sh

# Or manually:
# Install dependencies
pip install python-telegram-bot

# Copy team config
cp irc-dev-team.json ~/.picoclaw/workspace/teams/

# Create team
picoclaw team create ~/.picoclaw/workspace/teams/irc-dev-team.json
```

### Configuration

Create `.env` file in repo root:

```env
TELEGRAM_BOT_TOKEN=your_bot_token_here
PICOCLAW_BIN=picoclaw
IRC_TEAM_ID=irc-dev-team
```

Get your bot token from [@BotFather](https://t.me/BotFather) on Telegram.

**Note**: The gateway can also automatically read the token from `~/.picoclaw/config.json` if you have Telegram already configured:

```json
{
  "channels": {
    "telegram": {
      "token": "your_token_here"
    }
  }
}
```

## Usage

### Starting the Gateway

**Windows:**
```bash
cd workspace/skills/irc-gateway
start.bat
```

**Linux/Mac:**
```bash
cd workspace/skills/irc-gateway
python3 gateway.py
```

### Telegram Commands

- `/start` - Initialize session and show help
- `/who` - List all available roles
- `/status` - Check which roles are busy/idle
- `/team` - Show team configuration and status
- `/clear` - Clear session history

### Routing Messages

**Single Role:**
```
@architect design a REST API for user management
```

Response:
```
[cmmd1234] 🏗️ ARCHITECT: Here's the API design with endpoints...
```

**Multiple Roles (Parallel):**
```
hey @architect @developer let's build a login system
```

Response:
```
[cmmd1234] 🏗️ ARCHITECT: I'll design the authentication flow...
[cmmd1234] 💻 DEVELOPER: I'll implement the endpoints...
```

**Default Routing:**
```
what's the project status?
```

Routes to @manager when no mention is specified.

## Role Definitions

| Role | Emoji | Description | Capabilities |
|------|-------|-------------|--------------|
| Architect | 🏗️ | System design & architecture | design, architecture, planning |
| Developer | 💻 | Code implementation | code, implementation, debugging |
| Tester | 🧪 | Testing & QA | testing, qa, validation |
| Manager | 📋 | Project coordination | coordination, planning, reporting |

## Team Configuration

The gateway uses a team configuration file at `~/.picoclaw/workspace/teams/irc-dev-team.json`:

```json
{
  "team_id": "irc-dev-team",
  "name": "IRC Development Team",
  "pattern": "parallel",
  "roles": [
    {
      "name": "architect",
      "capabilities": ["design", "architecture", "planning"],
      "tools": ["*"],
      "model": ""
    }
    // ... more roles
  ],
  "coordinator": {
    "role": "manager"
  },
  "settings": {
    "max_delegation_depth": 3,
    "task_timeout": 120
  }
}
```

### Customizing Roles

1. **Edit team configuration** to add/modify roles
2. **Update gateway.py** ROLE_MAP with new role metadata
3. **Recreate team**:
```bash
picoclaw team dissolve irc-dev-team
picoclaw team create ~/.picoclaw/workspace/teams/irc-dev-team.json
```

### Per-Role Model Selection

Specify different models for each role:

```json
{
  "roles": [
    {
      "name": "architect",
      "model": "gpt-4"
    },
    {
      "name": "developer",
      "model": "claude-sonnet-4.6"
    }
  ]
}
```

## Integration with PicoClaw

The IRC Gateway is a **thin routing layer** that integrates with:

- ✅ **Team System** (`pkg/team`) - Multi-agent coordination
- ✅ **Agent Registry** - Role-based agent instances
- ✅ **Message Bus** (`pkg/bus`) - Inter-agent communication
- ✅ **Session Management** - Conversation context
- ✅ **Tool System** - All PicoClaw tools available to roles
- ✅ **Skills System** - Roles can use installed skills
- ✅ **MCP Support** - Model Context Protocol integration

## Comparison with Native Telegram Channel

| Feature | IRC Gateway | Native `pkg/channels/telegram` |
|---------|-------------|-------------------------------|
| **Purpose** | Mention-based role routing | Full-featured chat integration |
| **Integration** | CLI-based (thin layer) | Native Go packages |
| **Use Case** | Multi-role collaboration | General assistant |
| **Routing** | Explicit @mentions | Automatic |
| **Formatting** | IRC-style `[id] ROLE: msg` | Standard markdown |
| **Complexity** | ~300 lines Python | ~1000+ lines Go |

**When to use IRC Gateway:**
- Explicit role-based routing needed
- Working with development teams
- IRC-style formatting preferred
- Parallel role execution required

**When to use Native Channel:**
- General purpose AI assistant
- Single agent conversations
- Full PicoClaw feature set
- Production deployments

## Testing

Run the test suite:

```bash
cd workspace/skills/irc-gateway
python test_simple.py
```

Expected output:
```
============================================================
IRC Gateway - Simple Unit Tests
============================================================
...
Results: 12 passed, 0 failed
============================================================
```

## Troubleshooting

### Bot Not Responding

- Check `TELEGRAM_BOT_TOKEN` in `.env`
- Verify gateway is running: `ps aux | grep gateway.py`
- Check logs for errors

### Team Not Found

```bash
# Check if team exists
picoclaw team list

# Create team if missing
picoclaw team create ~/.picoclaw/workspace/teams/irc-dev-team.json

# Verify team status
picoclaw team status irc-dev-team
```

### Role Execution Timeout

- Increase timeout in team config: `"task_timeout": 180`
- Check PicoClaw binary is accessible: `which picoclaw`
- Test role directly: `picoclaw team execute irc-dev-team -t "test" --role architect`

### Python Dependencies

```bash
pip install python-telegram-bot
```

### PicoClaw Binary Not Found

```bash
# Check if picoclaw is in PATH
which picoclaw

# Or specify full path in .env
PICOCLAW_BIN=/usr/local/bin/picoclaw
```

## Advanced Usage

### Combining with Cron

Schedule periodic team tasks:

```bash
picoclaw cron add "0 9 * * *" "python gateway.py --task '@manager daily standup'"
```

### Custom Team Patterns

Change collaboration pattern in team config:

- `"pattern": "sequential"` - Roles execute in order
- `"pattern": "parallel"` - Roles execute simultaneously (default)
- `"pattern": "hierarchical"` - Dynamic task decomposition

### Integration with Skills

Roles can trigger PicoClaw skills:

```
@developer use the code-review skill on main.go
```

## Performance

- **Gateway overhead**: <50ms (parsing + formatting)
- **Team execution**: 1-10s (depends on LLM provider)
- **Parallel roles**: No additional latency
- **Memory usage**: ~20MB (Python + telegram bot)
- **Concurrent chats**: Limited by Telegram API rate limits

## Security

The gateway inherits PicoClaw's security features:

- Team tool access control
- Safety levels (strict/moderate/permissive)
- Allowlist configuration
- Environment variable protection

**Best Practices:**
1. Use environment variables for sensitive data
2. Restrict bot access via Telegram allowlist
3. Configure tool permissions in team config
4. Enable appropriate safety levels
5. Monitor logs for suspicious activity

## Files

```
workspace/skills/irc-gateway/
├── gateway.py              # Main gateway script
├── start.bat              # Windows launcher
├── setup.sh               # Linux/Mac setup script
├── test_simple.py         # Unit tests
├── test_gateway.py        # Full test suite
├── irc-dev-team.json      # Team configuration template
├── README.md              # Detailed documentation
├── INTEGRATION.md         # Architecture details
└── .env.example           # Environment template
```

## See Also

- [Multi-Agent Guide](../MULTI_AGENT_GUIDE.md)
- [Team Configuration](../../templates/teams/)
- [Telegram Channel Implementation](../../pkg/channels/telegram/)
- [Skills System](../../pkg/skills/)

## Contributing

To contribute improvements:

1. Test changes locally
2. Ensure compatibility with PicoClaw team system
3. Update documentation
4. Run test suite
5. Submit PR to main PicoClaw repo

## License

Part of PicoClaw ecosystem - MIT License
