# Collaborative Chat Quick Start

Get started with IRC-style collaborative multi-agent chat in Telegram in 5 minutes.

## Prerequisites

- PicoClaw installed and configured
- Telegram bot token
- Your Telegram user ID

## Step 1: Configure Telegram

Edit `~/.picoclaw/config.json`:

```json
{
  "channels": {
    "telegram": {
      "enabled": true,
      "token": "YOUR_BOT_TOKEN_HERE",
      "allow_from": ["YOUR_USER_ID_HERE"],
      "collaborative_chat": {
        "enabled": true,
        "default_team_id": "dev-team",
        "max_context_length": 50
      }
    }
  }
}
```

**How to get your Telegram user ID:**
1. Message [@userinfobot](https://t.me/userinfobot) on Telegram
2. Copy the ID number it sends you

## Step 2: Create Team Configuration

Copy the example team configuration:

```bash
# Linux/Mac
cp templates/teams/collaborative-dev-team.json ~/.picoclaw/teams/dev-team.json

# Windows
copy templates\teams\collaborative-dev-team.json %USERPROFILE%\.picoclaw\teams\dev-team.json
```

Or create `~/.picoclaw/teams/dev-team.json` manually:

```json
{
  "team_id": "dev-team",
  "name": "Development Team",
  "pattern": "parallel",
  "roles": [
    {
      "name": "architect",
      "capabilities": ["design", "architecture"],
      "model": "claude-3-5-sonnet-20241022",
      "tools": ["*"]
    },
    {
      "name": "developer",
      "capabilities": ["coding", "implementation"],
      "model": "claude-3-5-sonnet-20241022",
      "tools": ["*"]
    },
    {
      "name": "tester",
      "capabilities": ["testing", "qa"],
      "model": "claude-3-5-sonnet-20241022",
      "tools": ["*"]
    }
  ],
  "coordinator": {
    "role": "architect"
  },
  "settings": {
    "max_delegation_depth": 3,
    "agent_timeout_seconds": 300
  }
}
```

## Step 3: Start Gateway

```bash
picoclaw gateway
```

You should see:
```
✓ Channels enabled: [telegram]
✓ Gateway started on 0.0.0.0:8080
✓ Cron service started
✓ Heartbeat service started
Press Ctrl+C to stop
```

## Step 4: Test in Telegram

Open your Telegram bot and send:

```
@architect How should we structure a REST API?
```

You should get a response like:

```
[abc123] 🏗️ ARCHITECT: I recommend structuring your REST API with the following layers:

1. Controller Layer - Handle HTTP requests
2. Service Layer - Business logic
3. Repository Layer - Data access
4. Model Layer - Data structures

This follows the MVC pattern and ensures separation of concerns...
```

## Step 5: Try Multiple Agents

```
@architect @developer @tester How should we implement user authentication?
```

All three agents will respond in parallel:

```
[abc123] 🏗️ ARCHITECT: For authentication, I recommend JWT tokens with...

[abc123] 💻 DEVELOPER: I can implement this using the following approach...

[abc123] 🧪 TESTER: We should test the following scenarios...
```

## Available Roles

Default roles in the example configuration:

- 🏗️ `@architect` - System design and architecture
- 💻 `@developer` - Code implementation
- 🧪 `@tester` - Testing and QA
- 📋 `@manager` - Project management
- ⚙️ `@devops` - Deployment and infrastructure
- 🎨 `@designer` - UI/UX design

## How Agents Know Each Other

**Agents automatically know about all team members!** Every agent receives a team roster in their prompt:

```
=== Team Information ===
Team: Development Team
Members:
  • @architect - System design and architecture
  • @developer - Code implementation
  • @tester - Testing and QA
  • @manager - Project management
```

This enables natural agent-to-agent communication:

```
User: @developer implement login feature

[abc123] 💻 DEVELOPER: I'll implement the backend. @designer can you create the UI?

[abc123] 🎨 DESIGNER: Sure! I'll design the login screen with...
```

**Key Benefits:**
- Agents know who to ask for help
- Natural collaboration without manual configuration
- Role descriptions help agents understand each other's expertise

**Tip:** Add detailed descriptions in your team config to improve collaboration:

```json
{
  "name": "security",
  "description": "Security specialist - handles authentication, authorization, and security audits",
  "model": "claude-3-5-sonnet-20241022"
}
```

## Tips

### Conversation Context

Agents see the full conversation history:

```
User: @architect Design a user service
[abc123] 🏗️ ARCHITECT: Here's the design...

User: @developer Implement what the architect suggested
[abc123] 💻 DEVELOPER: Based on the architect's design, I'll implement...
```

### Parallel Execution

Multiple agents respond simultaneously:

```
User: @architect @developer @tester Review this code: [paste code]

# All three agents analyze and respond at the same time
```

### Session Management

Each chat has its own session with independent context. Sessions are created automatically and persist for the lifetime of the gateway process.

## Troubleshooting

### "No response from agents"

1. Check gateway is running: `picoclaw gateway`
2. Check logs: `tail -f ~/.picoclaw/logs/picoclaw.log`
3. Verify team is loaded: `picoclaw team list`

### "Team not found"

Ensure `default_team_id` in config matches `team_id` in team configuration file.

### "Permission denied"

Add your Telegram user ID to `allow_from` in config.

## Next Steps

- Read the [full documentation](COLLABORATIVE_CHAT.md)
- Customize team roles and capabilities
- Explore [multi-agent patterns](MULTI_AGENT_GUIDE.md)
- Configure [per-role models](MULTI_AGENT_MODEL_SELECTION.md)

## Example Use Cases

### Code Review

```
User: @architect @developer Review this implementation:
[paste code]

@tester What test cases should we add?
```

### Planning

```
User: @manager @architect Plan the next sprint

@developer What's your capacity?

@tester What testing is needed?
```

### Debugging

```
User: @developer This code is failing: [error]

@tester Can you reproduce this?

@architect Is this a design issue?
```

### Design Discussion

```
User: @architect @designer Design a login page

@developer How complex is the implementation?

@tester What should we test?
```

## Support

For issues or questions:
- Check [troubleshooting guide](COLLABORATIVE_CHAT.md#troubleshooting)
- Review [multi-agent documentation](MULTI_AGENT_GUIDE.md)
- Open an issue on GitHub
