# User Guides

Step-by-step guides for using PicoClaw features and capabilities.

## 📚 Available Guides

### v0.2.1 Features (NEW)

- **[V0.2.1_FEATURES.md](./V0.2.1_FEATURES.md)** - Complete guide to all v0.2.1 features
  - JSONL crash-safe storage
  - Vision/image support
  - Model routing
  - Environment configuration
  - Tool enable/disable
  - Extended thinking
  - PICOCLAW_HOME
  - New search providers

- **[VISION_SUPPORT.md](./VISION_SUPPORT.md)** - Multi-modal AI with images
  - Supported models
  - Usage examples
  - Configuration
  - Best practices

- **[MODEL_ROUTING.md](./MODEL_ROUTING.md)** - Cost optimization
  - Complexity scoring
  - Model tiers
  - Configuration
  - Cost analysis

### Getting Started

#### Multi-Agent System
- **[MULTI_AGENT_GUIDE.md](./MULTI_AGENT_GUIDE.md)** - Complete guide to working with multiple agents
  - Agent configuration
  - Routing rules
  - Session management
  - Best practices

#### Collaborative Features
- **[COLLABORATIVE_CHAT.md](./COLLABORATIVE_CHAT.md)** - Multi-agent conversations
  - @mention system
  - Cascade prevention
  - Context sharing
  - Rate limiting

- **[COLLABORATIVE_CHAT_QUICKSTART.md](./COLLABORATIVE_CHAT_QUICKSTART.md)** - Quick setup guide
  - 5-minute setup
  - Basic configuration
  - First conversation
  - Common patterns

### Team Features

- **[TEAM_AGENT_USAGE.md](./TEAM_AGENT_USAGE.md)** - Team delegation and coordination
  - Creating teams
  - Role definitions
  - Task delegation
  - Monitoring execution

### Security & Safety

- **[SAFETY_QUICKSTART.md](./SAFETY_QUICKSTART.md)** - Security configuration
  - Safety levels
  - Tool restrictions
  - Workspace isolation
  - Access control

### Cloud Integration

- **[ANTIGRAVITY_USAGE.md](./ANTIGRAVITY_USAGE.md)** - Cloud service integration
  - Authentication
  - API usage
  - Best practices
  - Troubleshooting

### System Administration

- **[SYSTEMD_QUICK_REFERENCE.md](./SYSTEMD_QUICK_REFERENCE.md)** - Linux service management
  - Service installation
  - Start/stop/restart
  - Log viewing
  - Auto-start configuration

### Performance Optimization

- **[COMPACTION_QUICK_REFERENCE.md](./COMPACTION_QUICK_REFERENCE.md)** - Memory compaction
  - When to compact
  - Configuration options
  - Performance impact
  - Monitoring

## 🎯 Guide Selection

### I want to...

#### Learn about v0.2.1 features
→ Start with [V0.2.1_FEATURES.md](./V0.2.1_FEATURES.md)

#### Use images with AI
→ Read [VISION_SUPPORT.md](./VISION_SUPPORT.md)

#### Optimize costs
→ Follow [MODEL_ROUTING.md](./MODEL_ROUTING.md)

#### Set up PicoClaw
→ Start with [MULTI_AGENT_GUIDE.md](./MULTI_AGENT_GUIDE.md)

#### Enable multi-agent conversations
→ Read [COLLABORATIVE_CHAT_QUICKSTART.md](./COLLABORATIVE_CHAT_QUICKSTART.md)

#### Create agent teams
→ Follow [TEAM_AGENT_USAGE.md](./TEAM_AGENT_USAGE.md)

#### Secure my installation
→ Check [SAFETY_QUICKSTART.md](./SAFETY_QUICKSTART.md)

#### Run as a system service
→ Use [SYSTEMD_QUICK_REFERENCE.md](./SYSTEMD_QUICK_REFERENCE.md)

#### Optimize memory usage
→ See [COMPACTION_QUICK_REFERENCE.md](./COMPACTION_QUICK_REFERENCE.md)

#### Integrate with cloud services
→ Follow [ANTIGRAVITY_USAGE.md](./ANTIGRAVITY_USAGE.md)

## 📖 Reading Order

### For First-Time Users
1. [MULTI_AGENT_GUIDE.md](./MULTI_AGENT_GUIDE.md) - Understand the basics
2. [SAFETY_QUICKSTART.md](./SAFETY_QUICKSTART.md) - Secure your setup
3. [COLLABORATIVE_CHAT_QUICKSTART.md](./COLLABORATIVE_CHAT_QUICKSTART.md) - Enable collaboration

### For Team Administrators
1. [TEAM_AGENT_USAGE.md](./TEAM_AGENT_USAGE.md) - Set up teams
2. [SAFETY_QUICKSTART.md](./SAFETY_QUICKSTART.md) - Configure security
3. [SYSTEMD_QUICK_REFERENCE.md](./SYSTEMD_QUICK_REFERENCE.md) - Deploy as service

### For Power Users
1. [COLLABORATIVE_CHAT.md](./COLLABORATIVE_CHAT.md) - Advanced collaboration
2. [COMPACTION_QUICK_REFERENCE.md](./COMPACTION_QUICK_REFERENCE.md) - Optimize performance
3. [ANTIGRAVITY_USAGE.md](./ANTIGRAVITY_USAGE.md) - Cloud integration

## 💡 Tips & Best Practices

### Configuration
- Always start with default config and modify incrementally
- Test changes in a development environment first
- Keep backups of working configurations

### Security
- Use workspace isolation for untrusted agents
- Configure tool restrictions appropriately
- Review safety levels regularly

### Performance
- Monitor memory usage with compaction
- Adjust session limits based on usage
- Use vector memory for long-term context

### Collaboration
- Define clear agent roles and responsibilities
- Set appropriate cascade depth limits
- Monitor mention patterns for loops

## 🔗 Related Documentation

- [Architecture Overview](../architecture/ARCHITECTURE_OVERVIEW.md) - System design
- [API Reference](../reference/API_REFERENCE.md) - API documentation
- [Configuration Reference](../reference/tools_configuration.md) - Config options
- [Troubleshooting](../development/troubleshooting.md) - Common issues

## 📝 Contributing

Help improve these guides:
1. Report unclear instructions
2. Suggest additional examples
3. Share your use cases
4. Translate to other languages

Submit improvements via [GitHub Issues](https://github.com/sipeed/picoclaw/issues) with the `documentation` label.

---

**Last Updated**: 2026-03-09
