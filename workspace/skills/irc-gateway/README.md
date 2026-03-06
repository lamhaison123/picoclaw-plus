# IRC-Style Communication Gateway

IRC-style Telegram bot gateway that integrates with PicoClaw's Team system for mention-based role routing.

## Features

- 🎯 **Team Integration**: Routes mentions to PicoClaw's multi-agent team system
- 🏷️ **Mention-Based Routing**: `@architect`, `@developer`, `@tester`, `@manager`
- ⚡ **Parallel Processing**: Tag multiple roles for concurrent execution
- 💬 **IRC Formatting**: Responses prefixed with `[session-id] emoji ROLE: message`
- 📊 **Session Management**: Maintains conversation context per chat
- 🔍 **Status Monitoring**: Real-time role status tracking
- 🤝 **Native Integration**: Uses `picoclaw team execute` under the hood

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

## Quick Start

> **Upgrading from v1.0.0?** See [MIGRATION.md](MIGRATION.md) for upgrade guide.

### 1. Setup Team Configuration

First, create the IRC team in PicoClaw:

```bash
# Copy team template to PicoClaw workspace
cp workspace/skills/irc-gateway/irc-dev-team.json ~/.picoclaw/workspace/teams/

# Create the team
picoclaw team create ~/.picoclaw/workspace/teams/irc-dev-team.json

# Verify team is created
picoclaw team status irc-dev-team
```

### 2. Configure Telegram Bot

The gateway automatically reads the Telegram token from your PicoClaw configuration.

**Option 1: Use existing PicoClaw Telegram config (Recommended)**

If you already have Telegram configured in `~/.picoclaw/config.json`:

```json
{
  "channels": {
    "telegram": {
      "enabled": true,
      "token": "123456789:ABCdefGHIjklMNOpqrsTUVwxyz"
    }
  }
}
```

The gateway will automatically use this token. No additional configuration needed!

**Option 2: Use separate bot token**

Create `.env` file in repo root:

```env
TELEGRAM_BOT_TOKEN=your_separate_bot_token
PICOCLAW_BIN=picoclaw
IRC_TEAM_ID=irc-dev-team
```

**Get a Telegram Bot Token:**

1. Open Telegram and search for [@BotFather](https://t.me/BotFather)
2. Send `/newbot` and follow instructions
3. Copy the token (format: `123456789:ABCdefGHIjklMNOpqrsTUVwxyz`)
4. Add to `~/.picoclaw/config.json` or `.env` file

**Configuration Priority:**
1. Environment variable `TELEGRAM_BOT_TOKEN`
2. `.env` file in repo root
3. `~/.picoclaw/config.json` (channels.telegram.token)

### 4. Run Gateway

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

## Usage

### Basic Commands

- `/start` - Initialize session and show help
- `/who` - List all available roles
- `/status` - Check which roles are busy/idle
- `/team` - Show team configuration and status
- `/clear` - Clear session history

### Routing Examples

**Single Role:**
```
@architect design a REST API for user management
```
Response:
```
[cmmd1234] 🏗️ ARCHITECT: Here's the API design...
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

**Default (Manager):**
```
what's the project status?
```
(Routes to @manager when no mention)

### Response Format

```
[session-id] emoji ROLE: response content
```

Example:
```
[cmmd8145] 🏗️ ARCHITECT: Based on your requirements...
[cmmd8145] 💻 DEVELOPER: I've implemented the following...
```

## Role Definitions

| Role | Emoji | Description | Capabilities |
|------|-------|-------------|--------------|
| Architect | 🏗️ | System design & architecture | design, architecture, planning |
| Developer | 💻 | Code implementation | code, implementation, debugging |
| Tester | 🧪 | Testing & QA | testing, qa, validation |
| Manager | 📋 | Project coordination | coordination, planning, reporting |

## Integration with PicoClaw

This gateway is a **thin routing layer** that:

1. Parses @mentions from Telegram messages
2. Routes to appropriate team roles via `picoclaw team execute`
3. Formats responses in IRC style
4. Maintains session context

**Key Integration Points:**
- ✅ Uses `pkg/team` for multi-agent coordination
- ✅ Leverages `picoclaw team execute` CLI
- ✅ Respects team configuration in `~/.picoclaw/workspace/teams/`
- ✅ Works with existing PicoClaw agent system
- ✅ Supports all configured LLM providers

## Configuration

### Automatic Token Loading

The gateway automatically loads the Telegram bot token from multiple sources (in priority order):

1. **Environment variable**: `TELEGRAM_BOT_TOKEN`
2. **Repo .env file**: `.env` in repository root
3. **PicoClaw config**: `~/.picoclaw/config.json` → `channels.telegram.token`

**Example PicoClaw config:**
```json
{
  "channels": {
    "telegram": {
      "enabled": true,
      "token": "123456789:ABCdefGHIjklMNOpqrsTUVwxyz",
      "allow_from": ["your_user_id"]
    }
  }
}
```

**Example .env file:**
```env
TELEGRAM_BOT_TOKEN=123456789:ABCdefGHIjklMNOpqrsTUVwxyz
PICOCLAW_BIN=picoclaw
IRC_TEAM_ID=irc-dev-team
```

### Team Configuration

Edit `~/.picoclaw/workspace/teams/irc-dev-team.json`:

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
      "model": "gpt-4"  // Optional: specify model per role
    },
    {
      "name": "developer",
      "capabilities": ["code", "implementation"],
      "tools": ["*"],
      "model": "claude-sonnet-4.6"
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

### Custom Roles

To add new roles:

1. Edit `irc-dev-team.json` and add role definition
2. Update `gateway.py` ROLE_MAP:

```python
ROLE_MAP = {
    "architect": {"emoji": "🏗️", "description": "System design"},
    "devops": {"emoji": "🚀", "description": "DevOps & deployment"},  # New role
}
```

3. Recreate team:
```bash
picoclaw team dissolve irc-dev-team
picoclaw team create ~/.picoclaw/workspace/teams/irc-dev-team.json
```

### Environment Variables

```env
# Optional - only needed if not using PicoClaw config
TELEGRAM_BOT_TOKEN=your_token

# Optional overrides
PICOCLAW_BIN=picoclaw              # Path to picoclaw binary
IRC_TEAM_ID=irc-dev-team           # Team ID to use
```

**Note**: If you already have Telegram configured in `~/.picoclaw/config.json`, you don't need to set `TELEGRAM_BOT_TOKEN` separately.

## Troubleshooting

**Bot not responding:**
- Check if token is configured in `~/.picoclaw/config.json` or `.env`
- Verify gateway loaded token: check startup logs for "✅ Loaded Telegram token"
- Verify bot is running: `ps aux | grep gateway.py`
- Check logs for errors

**Team not found:**
```bash
# Check if team exists
picoclaw team list

# Create team if missing
picoclaw team create ~/.picoclaw/workspace/teams/irc-dev-team.json
```

**Role timeout:**
- Increase timeout in team config: `"task_timeout": 180`
- Check picoclaw binary is accessible: `which picoclaw`

**Dependencies error:**
```bash
pip install python-telegram-bot
```

**PicoClaw binary not found:**
```bash
# Check if picoclaw is in PATH
which picoclaw

# Or specify full path in .env
PICOCLAW_BIN=/path/to/picoclaw
```

## Development

**Test locally:**
```bash
python gateway.py
```

**View logs:**
```bash
# Logs show routing decisions and role responses
[14:23:45] INFO: [cmmd0123] Routing to @architect: design a REST API...
[14:23:50] INFO: Team execution completed (5.2s)
```

**Debug mode:**
```python
# In gateway.py, change log level
logging.basicConfig(level=logging.DEBUG)
```

## Comparison: Gateway vs Native Telegram Channel

| Feature | IRC Gateway | Native `pkg/channels/telegram` |
|---------|-------------|-------------------------------|
| **Purpose** | Mention-based role routing | Full-featured chat integration |
| **Integration** | Thin layer over team system | Deep integration with agent loop |
| **Use Case** | Multi-role collaboration | General chat assistant |
| **Session** | Per-chat sessions | Full session management |
| **Formatting** | IRC-style `[id] ROLE: msg` | Standard markdown |
| **Routing** | Explicit @mentions | Automatic |

**When to use IRC Gateway:**
- You want explicit role-based routing
- Working with development teams
- Need IRC-style formatting
- Want parallel role execution

**When to use Native Channel:**
- General purpose AI assistant
- Single agent conversations
- Full PicoClaw feature set
- Production deployments

## Advanced Usage

### Combining with Cron

Schedule periodic team tasks:

```bash
# Add cron job to check project status daily
picoclaw cron add "0 9 * * *" "python gateway.py --task '@manager daily standup report'"
```

### Integration with Skills

The gateway can trigger PicoClaw skills through team roles:

```
@developer use the code-review skill on main.go
```

### Custom Team Patterns

Change collaboration pattern in team config:

- `"pattern": "sequential"` - Roles execute in order
- `"pattern": "parallel"` - Roles execute simultaneously (default)
- `"pattern": "hierarchical"` - Dynamic task decomposition

## License

Part of PicoClaw ecosystem - MIT License

## Contributing

This is a skill/extension for PicoClaw. To contribute:

1. Test changes locally
2. Ensure compatibility with PicoClaw team system
3. Update documentation
4. Submit PR to main PicoClaw repo

## See Also

- [PicoClaw Multi-Agent Guide](../../../docs/MULTI_AGENT_GUIDE.md)
- [Team Configuration](../../../templates/teams/)
- [Telegram Channel Implementation](../../../pkg/channels/telegram/)
