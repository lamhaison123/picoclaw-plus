# ⚠️ DEPRECATED - IRC Gateway Python Skill

This Python skill has been **replaced by native Go implementation** in PicoClaw v1.1.0+.

## Why Deprecated?

The native implementation offers significant advantages:
- ✅ Single Telegram bot token (no need for 3-4 separate bots)
- ✅ Better integration with PicoClaw's team system
- ✅ Lower memory footprint (~5MB vs ~50MB)
- ✅ Faster response times (native Go vs Python subprocess)
- ✅ More reliable (no external process management)
- ✅ Easier configuration (integrated into main config)

## Migration Guide

### Step 1: Remove Python Skill

```bash
# Stop any running gateway
systemctl stop picoclaw.service  # if using systemd

# Remove the skill directory
rm -rf ~/.picoclaw/workspace/skills/irc-gateway
```

### Step 2: Configure Native Collaborative Chat

Edit `~/.picoclaw/config.json`:

```json
{
  "channels": {
    "telegram": {
      "enabled": true,
      "token": "YOUR_SINGLE_BOT_TOKEN",
      "allow_from": ["YOUR_USER_ID"],
      "collaborative_chat": {
        "enabled": true,
        "default_team_id": "dev-team",
        "max_context_length": 50
      }
    }
  }
}
```

### Step 3: Create Team Configuration

Copy the example:

```bash
cp templates/teams/collaborative-dev-team.json ~/.picoclaw/teams/dev-team.json
```

Or create your own team config with the roles you need.

### Step 4: Restart Gateway

```bash
picoclaw gateway
```

## Documentation

- **Quick Start**: [docs/COLLABORATIVE_CHAT_QUICKSTART.md](../../../docs/COLLABORATIVE_CHAT_QUICKSTART.md)
- **Full Guide**: [docs/COLLABORATIVE_CHAT.md](../../../docs/COLLABORATIVE_CHAT.md)
- **Multi-Agent Guide**: [docs/MULTI_AGENT_GUIDE.md](../../../docs/MULTI_AGENT_GUIDE.md)

## Usage Comparison

### Old (Python Skill)
```
# Required 3-4 separate Telegram bots
# External Python process
# Manual process management
# Separate configuration file
```

### New (Native Go)
```
User: @architect @developer How to implement auth?

[abc123] 🏗️ ARCHITECT: I recommend JWT tokens...
[abc123] 💻 DEVELOPER: I can implement...

# Single bot token
# Native integration
# Automatic management
# Integrated configuration
```

## Support

If you encounter issues during migration:
1. Check the [troubleshooting guide](../../../docs/COLLABORATIVE_CHAT.md#troubleshooting)
2. Review your team configuration
3. Check logs: `tail -f ~/.picoclaw/logs/picoclaw.log`
4. Open an issue on GitHub

## Timeline

- **v1.0.x**: Python skill supported
- **v1.1.0+**: Native implementation available, Python skill deprecated
- **v1.2.0+**: Python skill will be removed

Please migrate to the native implementation as soon as possible.
