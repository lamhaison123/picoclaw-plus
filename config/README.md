# Configuration Examples

This directory contains example configuration files for PicoClaw, including specialized configurations for collaborative chat.

## Quick Start

1. Copy an example config to your home directory:
```bash
cp config/collaborative-chat-minimal.json ~/.picoclaw/config.json
```

2. Edit the config with your API keys:
```bash
nano ~/.picoclaw/config.json
```

3. Start the gateway:
```bash
picoclaw gateway
```

## Configuration Files

### Main Configuration Examples

#### `config.example.json`
Complete configuration example with all available options. Use this as a reference for all PicoClaw features.

**Use when:** You want to see all available configuration options.

### Collaborative Chat Configurations

#### `collaborative-chat-minimal.json`
Minimal configuration for collaborative chat with Anthropic Claude.

**Use when:** 
- First time setup
- You want the simplest working configuration
- Testing collaborative chat

**Includes:**
- Anthropic provider (Claude)
- Telegram channel with collaborative chat
- Basic team configuration

**Setup:**
```bash
# 1. Copy config
cp config/collaborative-chat-minimal.json ~/.picoclaw/config.json

# 2. Edit with your keys
nano ~/.picoclaw/config.json
# - Replace YOUR_ANTHROPIC_API_KEY
# - Replace YOUR_TELEGRAM_BOT_TOKEN
# - Replace YOUR_TELEGRAM_USER_ID

# 3. Copy team config
cp templates/teams/minimal-team.json ~/.picoclaw/teams/dev-team.json

# 4. Start gateway
picoclaw gateway
```

#### `collaborative-chat-full.json`
Full-featured configuration with all collaborative chat options.

**Use when:**
- You need advanced features
- Custom role mappings
- Auto-join rules
- Multiple users

**Includes:**
- All collaborative chat options
- Role mapping (shortcuts like @arch → @architect)
- Auto-join rules (keyword-based agent activation)
- Multiple allowed users
- Group trigger configuration
- Placeholder messages

**Features:**
- **Role Mapping**: Use shortcuts like `@arch` instead of `@architect`
- **Auto-join Rules**: Agents automatically join when keywords detected
- **Session Prefix**: Custom session ID prefix
- **Extended Context**: Up to 50 messages

#### `collaborative-chat-openai.json`
Configuration using OpenAI GPT-4 models.

**Use when:**
- You prefer OpenAI over Anthropic
- You have OpenAI API access
- You want to use GPT-4

**Model:** `gpt-4-turbo-preview`

**Setup:**
```bash
cp config/collaborative-chat-openai.json ~/.picoclaw/config.json
# Edit: Replace YOUR_OPENAI_API_KEY
```

#### `collaborative-chat-ollama.json`
Configuration for local Ollama models (no API key needed).

**Use when:**
- You want to run locally
- Privacy is important
- No internet connection
- Cost-free operation

**Model:** `llama3.1:8b` (or any Ollama model)

**Prerequisites:**
```bash
# Install Ollama
curl -fsSL https://ollama.com/install.sh | sh

# Pull a model
ollama pull llama3.1:8b
```

**Setup:**
```bash
cp config/collaborative-chat-ollama.json ~/.picoclaw/config.json
# No API key needed!
```

**Note:** Ollama models are slower but completely free and private.

## Team Configuration Examples

Located in `templates/teams/`:

### `minimal-team.json`
Simplest team with 2 roles: architect and developer.

**Use when:** Testing or simple use cases.

### `collaborative-dev-team.json`
Full development team with 6 roles.

**Roles:**
- 🏗️ architect - System design
- 💻 developer - Implementation
- 🧪 tester - QA and testing
- 📋 manager - Project management
- ⚙️ devops - Deployment and ops
- 🎨 designer - UI/UX design

**Use when:** Full-featured development collaboration.

### `research-team-collab.json`
Research team for information gathering and analysis.

**Roles:**
- 🔍 researcher - Information gathering
- 📊 analyst - Data analysis
- ✍️ writer - Documentation

**Use when:** Research projects, data analysis, report writing.

### `support-team-collab.json`
Customer support team with hierarchical structure.

**Roles:**
- 💬 support - First-line support
- 🔧 technical - Technical specialist
- 👔 manager - Escalation point

**Use when:** Customer support, help desk, ticket handling.

### `creative-team-collab.json`
Creative team for content and design.

**Roles:**
- ✍️ writer - Content creation
- 🎨 designer - Visual design
- 🎯 strategist - Creative direction
- ✅ reviewer - Quality assurance

**Use when:** Content creation, marketing, branding.

## Configuration Fields Explained

### Provider Configuration

```json
{
  "provider": {
    "name": "anthropic",           // Provider: anthropic, openai, ollama, etc.
    "api_key": "YOUR_KEY",          // API key (not needed for ollama)
    "base_url": "",                 // Custom API endpoint (optional)
    "timeout": 300                  // Request timeout in seconds
  }
}
```

### Telegram Configuration

```json
{
  "channels": {
    "telegram": {
      "enabled": true,
      "token": "YOUR_BOT_TOKEN",
      "allow_from": ["USER_ID"],    // Whitelist of user IDs
      "collaborative_chat": {
        "enabled": true,
        "default_team_id": "dev-team",
        "max_context_length": 50,   // Max messages in context
        "role_map": {                // Shortcuts
          "arch": "architect",
          "dev": "developer"
        },
        "auto_join_rules": [         // Auto-join on keywords
          {
            "role": "architect",
            "keywords": ["design", "architecture"]
          }
        ]
      }
    }
  }
}
```

### Team Configuration

```json
{
  "team_id": "dev-team",
  "name": "Development Team",
  "pattern": "parallel",            // parallel, sequential, hierarchical
  "roles": [
    {
      "name": "architect",
      "description": "System architect",
      "capabilities": ["design"],
      "model": "claude-3-5-sonnet-20241022",
      "tools": ["*"]                // All tools, or specific: ["file_read"]
    }
  ],
  "settings": {
    "max_delegation_depth": 3,
    "agent_timeout_seconds": 300
  }
}
```

## Getting Your API Keys

### Anthropic (Claude)
1. Go to https://console.anthropic.com
2. Sign up / Log in
3. Go to API Keys
4. Create new key
5. Copy the key (starts with `sk-ant-`)

### OpenAI (GPT-4)
1. Go to https://platform.openai.com
2. Sign up / Log in
3. Go to API Keys
4. Create new secret key
5. Copy the key (starts with `sk-`)

### Telegram Bot
1. Open Telegram
2. Search for @BotFather
3. Send `/newbot`
4. Follow instructions
5. Copy the bot token

### Your Telegram User ID
1. Open Telegram
2. Search for @userinfobot
3. Send any message
4. Copy your ID (number)

## Common Patterns

### Development Team
```bash
# Config
cp config/collaborative-chat-minimal.json ~/.picoclaw/config.json

# Team
cp templates/teams/collaborative-dev-team.json ~/.picoclaw/teams/dev-team.json

# Usage
@architect @developer How to implement authentication?
```

### Research Team
```bash
# Config
cp config/collaborative-chat-full.json ~/.picoclaw/config.json

# Team
cp templates/teams/research-team-collab.json ~/.picoclaw/teams/research-team.json

# Update config: "default_team_id": "research-team"

# Usage
@researcher @analyst Research market trends for AI assistants
```

### Local/Offline Setup
```bash
# Install Ollama
curl -fsSL https://ollama.com/install.sh | sh
ollama pull llama3.1:8b

# Config
cp config/collaborative-chat-ollama.json ~/.picoclaw/config.json

# Team
cp templates/teams/minimal-team.json ~/.picoclaw/teams/local-team.json

# Usage (works offline!)
@assistant @coder Help me with this code
```

## Troubleshooting

### "Team not found"
- Check `default_team_id` matches your team file name
- Verify team file is in `~/.picoclaw/teams/`
- Check team file is valid JSON

### "No response from agents"
- Verify API key is correct
- Check `allow_from` includes your user ID
- Check logs: `tail -f ~/.picoclaw/logs/picoclaw.log`

### "Permission denied"
- Add your Telegram user ID to `allow_from` array
- Get your ID from @userinfobot

### Ollama connection failed
- Check Ollama is running: `ollama list`
- Verify base_url: `http://localhost:11434`
- Pull model: `ollama pull llama3.1:8b`

## Advanced Configuration

### Multiple Teams
You can create multiple teams and switch between them:

```json
{
  "collaborative_chat": {
    "default_team_id": "dev-team"
  }
}
```

Teams in `~/.picoclaw/teams/`:
- `dev-team.json` - Development
- `research-team.json` - Research
- `support-team.json` - Support

### Per-Role Models
Different models for different roles:

```json
{
  "roles": [
    {
      "name": "architect",
      "model": "claude-3-5-sonnet-20241022"  // Best model
    },
    {
      "name": "tester",
      "model": "claude-3-haiku-20240307"     // Faster/cheaper
    }
  ]
}
```

### Tool Restrictions
Limit tools per role:

```json
{
  "roles": [
    {
      "name": "architect",
      "tools": ["*"]                    // All tools
    },
    {
      "name": "tester",
      "tools": ["file_read", "exec_safe"]  // Limited
    }
  ]
}
```

## See Also

- [Collaborative Chat Guide](../docs/COLLABORATIVE_CHAT.md)
- [Quick Start](../docs/COLLABORATIVE_CHAT_QUICKSTART.md)
- [Multi-Agent Guide](../docs/MULTI_AGENT_GUIDE.md)
- [Architecture Details](../docs/COLLABORATIVE_CHAT_ARCHITECTURE.md)
