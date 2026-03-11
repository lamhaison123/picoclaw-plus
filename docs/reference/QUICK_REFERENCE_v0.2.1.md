# PicoClaw v0.2.1 Quick Reference Card

## 🚀 New Features at a Glance

### 🛡️ Crash-Safe Storage
```json
{
  "session": {
    "storage_backend": "jsonl",  // Default: crash-safe
    "auto_migrate": true          // Auto-upgrade from JSON
  }
}
```

### 🖼️ Vision Support
```bash
# Send images to AI
picoclaw agent
> Describe this: /path/to/image.jpg
> What's in this chart: media://chart.png
```

### 💰 Model Routing
```json
{
  "routing": {
    "enabled": true,
    "tiers": [
      {"name": "cheap", "models": ["gpt-4o-mini"]},
      {"name": "expensive", "models": ["gpt-5.2"]}
    ]
  }
}
```

### 🧠 Extended Thinking
```bash
# See AI's reasoning process
picoclaw agent --model claude-sonnet-4.6
> Solve this complex problem...
# Thinking process shown in reasoning channel
```

### ⚡ Parallel Tools
```bash
# Tools execute in parallel automatically
> Search for "AI" and "ML" and "DL"
# All searches run simultaneously
```

## 🔧 Configuration

### Environment Variables (.env)
```bash
# Home directory
PICOCLAW_HOME=/custom/path

# Routing
PICOCLAW_ROUTING_ENABLED=true

# Session
PICOCLAW_SESSION_STORAGE_BACKEND=jsonl
PICOCLAW_SESSION_AUTO_MIGRATE=true
PICOCLAW_SESSION_SUMMARIZATION_MESSAGE_THRESHOLD=20
PICOCLAW_SESSION_SUMMARIZATION_TOKEN_PERCENT=0.75

# Tools
PICOCLAW_TOOLS_FILE_ENABLED=true
PICOCLAW_TOOLS_SHELL_ENABLED=true
PICOCLAW_TOOLS_WEB_ENABLED=true
PICOCLAW_TOOLS_MESSAGE_ENABLED=true
PICOCLAW_TOOLS_SPAWN_ENABLED=true
PICOCLAW_TOOLS_TEAM_ENABLED=true
PICOCLAW_TOOLS_SKILL_ENABLED=true
PICOCLAW_TOOLS_HARDWARE_ENABLED=true

# Media
PICOCLAW_AGENTS_DEFAULTS_MAX_MEDIA_SIZE=20971520
```

### Tool Control
```json
{
  "tools": {
    "file_tools_enabled": true,
    "shell_tools_enabled": true,
    "web_tools_enabled": true,
    "message_tool_enabled": true,
    "spawn_tool_enabled": true,
    "team_tools_enabled": true,
    "skill_tools_enabled": true,
    "hardware_tools_enabled": true
  }
}
```

### Summarization
```json
{
  "session": {
    "summarization_message_threshold": 20,
    "summarization_token_percent": 0.75
  }
}
```

## 📁 Directory Structure

### Default (~/.picoclaw)
```
~/.picoclaw/
├── config.json          # Configuration
├── auth.json            # Credentials
├── .env                 # Environment variables
├── workspace/           # Default workspace
├── sessions/            # Session files (.jsonl)
├── teams/               # Team states
└── skills/              # Global skills
```

### Custom (PICOCLAW_HOME)
```bash
export PICOCLAW_HOME=/custom/path
# All data stored in /custom/path/
```

## 🎯 Common Tasks

### Enable Model Routing
```bash
# 1. Edit config.json
{
  "routing": {
    "enabled": true,
    "tiers": [
      {"name": "cheap", "models": ["gpt-4o-mini"]},
      {"name": "expensive", "models": ["gpt-5.2"]}
    ]
  }
}

# 2. Restart PicoClaw
picoclaw agent
```

### Send Images
```bash
# Method 1: Direct path
picoclaw agent
> Describe /path/to/image.jpg

# Method 2: Media reference
> Analyze media://image-id
```

### Use Custom Home
```bash
# Temporary
export PICOCLAW_HOME=/custom/path
picoclaw agent

# Permanent (add to ~/.bashrc)
echo 'export PICOCLAW_HOME=/custom/path' >> ~/.bashrc
```

### Disable Tools
```bash
# Via config.json
{
  "tools": {
    "web_tools_enabled": false,
    "shell_tools_enabled": false
  }
}

# Via environment
export PICOCLAW_TOOLS_WEB_ENABLED=false
export PICOCLAW_TOOLS_SHELL_ENABLED=false
```

### Migrate from JSON
```bash
# Automatic (recommended)
{
  "session": {
    "storage_backend": "jsonl",
    "auto_migrate": true
  }
}
# Restart - migration happens automatically

# Manual (if needed)
# Old sessions in ~/.picoclaw/sessions/*.json
# Will be migrated on first access
```

## 🐛 Troubleshooting

### Issue: Data loss on crash
**Solution**: Enable JSONL storage (default in v0.2.1)
```json
{"session": {"storage_backend": "jsonl"}}
```

### Issue: High costs
**Solution**: Enable model routing
```json
{"routing": {"enabled": true}}
```

### Issue: Images not working
**Solution**: Check max_media_size
```json
{"agents": {"defaults": {"max_media_size": 20971520}}}
```

### Issue: Tools not available
**Solution**: Check tool flags
```json
{"tools": {"web_tools_enabled": true}}
```

### Issue: Config not loading
**Solution**: Check PICOCLAW_HOME
```bash
echo $PICOCLAW_HOME
# Should point to correct directory
```

## 📊 Performance Tips

### Optimize Costs
- Enable model routing
- Use cheap models for simple queries
- Configure appropriate tiers

### Improve Speed
- Parallel tool execution (automatic)
- Reduce max_media_size if not needed
- Adjust summarization thresholds

### Save Memory
- Lower max_media_size
- Enable summarization
- Clean old sessions

## 🔒 Security Best Practices

### Protect Secrets
```bash
# Use .env for secrets
echo "OPENAI_API_KEY=sk-..." > .env

# Add to .gitignore
echo ".env" >> .gitignore
echo ".env.local" >> .gitignore
```

### File Permissions
```bash
# Secure auth file
chmod 600 ~/.picoclaw/auth.json

# Secure .env
chmod 600 .env
```

### Multi-User Setup
```bash
# User 1
export PICOCLAW_HOME=/home/user1/.picoclaw

# User 2
export PICOCLAW_HOME=/home/user2/.picoclaw
```

## 📚 Resources

### Documentation
- `docs/` - Complete documentation
- `CHANGELOG.md` - Version history
- `.env.example` - Configuration examples
- `INTEGRATION_CHECKLIST.md` - Feature list

### Support
- GitHub Issues
- Documentation
- Examples directory

## 🎓 Examples

### Basic Usage
```bash
picoclaw agent
> Hello
> What can you do?
> Exit
```

### With Vision
```bash
picoclaw agent
> Describe this chart: /path/to/chart.png
> What trends do you see?
```

### With Routing
```bash
# Simple query (cheap model)
> What is 2+2?

# Complex query (expensive model)
> Write a sorting algorithm in Python
```

### Multi-User
```bash
# Terminal 1 (User A)
export PICOCLAW_HOME=/data/user-a
picoclaw agent

# Terminal 2 (User B)
export PICOCLAW_HOME=/data/user-b
picoclaw agent
```

---

**Version**: v0.2.1 Integration  
**Last Updated**: 2026-03-09  
**Status**: Production Ready
