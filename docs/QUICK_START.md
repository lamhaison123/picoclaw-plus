# Documentation Quick Start

## 🚀 Where to Start?

### I'm a New User
**Start here**: [Multi-Agent Guide](./guides/MULTI_AGENT_GUIDE.md)

Then read:
1. [Architecture Overview](./architecture/ARCHITECTURE_OVERVIEW.md) - Understand the system
2. [Safety Quickstart](./guides/SAFETY_QUICKSTART.md) - Secure your setup
3. [Collaborative Chat Quickstart](./guides/COLLABORATIVE_CHAT_QUICKSTART.md) - Enable multi-agent chat

### I'm a Developer
**Start here**: [Developer Guide](./development/DEVELOPER_GUIDE.md)

Then read:
1. [Component Details](./architecture/COMPONENT_DETAILS.md) - Implementation details
2. [Data Flow](./architecture/DATA_FLOW.md) - Message processing
3. [Troubleshooting](./development/troubleshooting.md) - Debug techniques

### I'm a System Administrator
**Start here**: [Safety Levels](./reference/SAFETY_LEVELS.md)

Then read:
1. [Tools Configuration](./reference/tools_configuration.md) - Configure tools
2. [Systemd Quick Reference](./guides/SYSTEMD_QUICK_REFERENCE.md) - Service setup
3. [Circuit Breaker](./reference/CIRCUIT_BREAKER.md) - Fault tolerance

### I'm an Architect
**Start here**: [Architecture Overview](./architecture/ARCHITECTURE_OVERVIEW.md)

Then read:
1. [Component Details](./architecture/COMPONENT_DETAILS.md) - System components
2. [Data Flow](./architecture/DATA_FLOW.md) - Processing flows
3. [Collaborative Chat Architecture](./development/COLLABORATIVE_CHAT_ARCHITECTURE.md) - Chat system

## 📚 Documentation Structure

```
docs/
├── README.md            ← Overview and navigation
├── INDEX.md             ← Complete documentation index
├── QUICK_START.md       ← This file
│
├── architecture/        ← System design
│   ├── README.md
│   ├── ARCHITECTURE_OVERVIEW.md    ⭐ Start here for architecture
│   ├── COMPONENT_DETAILS.md
│   └── DATA_FLOW.md
│
├── guides/              ← User guides
│   ├── README.md
│   ├── MULTI_AGENT_GUIDE.md        ⭐ Start here for users
│   ├── COLLABORATIVE_CHAT.md
│   └── TEAM_AGENT_USAGE.md
│
├── reference/           ← API & config reference
│   ├── README.md
│   ├── API_REFERENCE.md
│   ├── SAFETY_LEVELS.md            ⭐ Start here for security
│   └── tools_configuration.md
│
└── development/         ← Development resources
    ├── README.md
    ├── DEVELOPER_GUIDE.md          ⭐ Start here for developers
    ├── COLLABORATIVE_CHAT_ARCHITECTURE.md
    └── troubleshooting.md
```

## 🎯 Common Tasks

### Setup PicoClaw
1. Read [Multi-Agent Guide](./guides/MULTI_AGENT_GUIDE.md)
2. Configure agents in `config.json`
3. Follow [Safety Quickstart](./guides/SAFETY_QUICKSTART.md)
4. Setup your [Channel](./channels/)

### Enable Multi-Agent Chat
1. Read [Collaborative Chat Quickstart](./guides/COLLABORATIVE_CHAT_QUICKSTART.md)
2. Configure mention system
3. Test with @mentions
4. Review [Collaborative Chat](./guides/COLLABORATIVE_CHAT.md) for advanced features

### Create Agent Teams
1. Read [Team Agent Usage](./guides/TEAM_AGENT_USAGE.md)
2. Define team roles in config
3. Use `delegate_to_team` tool
4. Monitor team execution

### Secure Your Installation
1. Read [Safety Levels](./reference/SAFETY_LEVELS.md)
2. Set appropriate safety level (2-3 for production)
3. Configure [Tools](./reference/tools_configuration.md)
4. Enable workspace isolation

### Deploy as Service
1. Read [Systemd Quick Reference](./guides/SYSTEMD_QUICK_REFERENCE.md)
2. Create systemd service file
3. Enable auto-start
4. Monitor logs

### Develop Features
1. Read [Developer Guide](./development/DEVELOPER_GUIDE.md)
2. Study [Component Details](./architecture/COMPONENT_DETAILS.md)
3. Write tests
4. Submit PR

### Troubleshoot Issues
1. Check [Troubleshooting](./development/troubleshooting.md)
2. Enable debug logging
3. Review [Data Flow](./architecture/DATA_FLOW.md)
4. Ask in [Discussions](https://github.com/sipeed/picoclaw/discussions)

## 🔍 Finding Specific Information

### Configuration
- Main config: [Tools Configuration](./reference/tools_configuration.md)
- Memory config: [Memory Config Guide](./memory/CONFIG_GUIDE.md)
- Safety config: [Safety Levels](./reference/SAFETY_LEVELS.md)

### API
- REST API: [API Reference](./reference/API_REFERENCE.md)
- Tool API: [Tools Configuration](./reference/tools_configuration.md)
- Team API: [Team Tool Access](./reference/TEAM_TOOL_ACCESS.md)

### Architecture
- System overview: [Architecture Overview](./architecture/ARCHITECTURE_OVERVIEW.md)
- Components: [Component Details](./architecture/COMPONENT_DETAILS.md)
- Message flow: [Data Flow](./architecture/DATA_FLOW.md)

### Features
- Multi-agent: [Multi-Agent Guide](./guides/MULTI_AGENT_GUIDE.md)
- Collaboration: [Collaborative Chat](./guides/COLLABORATIVE_CHAT.md)
- Teams: [Team Agent Usage](./guides/TEAM_AGENT_USAGE.md)
- Memory: [Memory Config Guide](./memory/CONFIG_GUIDE.md)

## 📖 Reading Paths

### 30-Minute Overview
1. [README](../README.md) - 5 min
2. [Architecture Overview](./architecture/ARCHITECTURE_OVERVIEW.md) - 15 min
3. [Multi-Agent Guide](./guides/MULTI_AGENT_GUIDE.md) - 10 min

### 2-Hour Deep Dive
1. [Architecture Overview](./architecture/ARCHITECTURE_OVERVIEW.md) - 20 min
2. [Component Details](./architecture/COMPONENT_DETAILS.md) - 30 min
3. [Data Flow](./architecture/DATA_FLOW.md) - 20 min
4. [Multi-Agent Guide](./guides/MULTI_AGENT_GUIDE.md) - 15 min
5. [Collaborative Chat](./guides/COLLABORATIVE_CHAT.md) - 15 min
6. [Developer Guide](./development/DEVELOPER_GUIDE.md) - 20 min

### Full Documentation (1 Day)
1. Architecture section - 2 hours
2. Guides section - 2 hours
3. Reference section - 2 hours
4. Development section - 1 hour
5. Channel-specific docs - 1 hour

## 🔗 External Resources

### Project
- **Repository**: https://github.com/sipeed/picoclaw
- **Issues**: https://github.com/sipeed/picoclaw/issues
- **Discussions**: https://github.com/sipeed/picoclaw/discussions

### Community
- **Contributing**: [../CONTRIBUTING.md](../CONTRIBUTING.md)
- **Changelog**: [../CHANGELOG.md](../CHANGELOG.md)
- **Roadmap**: [../ROADMAP.md](../ROADMAP.md)

## 💡 Tips

### Navigation
- Use `docs/INDEX.md` to find all documents
- Each directory has a `README.md` with overview
- Follow links in documents for related topics

### Learning
- Start with overview documents
- Dive into details as needed
- Try examples in guides
- Reference API docs when coding

### Contributing
- Read [Developer Guide](./development/DEVELOPER_GUIDE.md) first
- Check [Contributing Guide](../CONTRIBUTING.md)
- Follow code style guidelines
- Write tests for new features

## 📞 Getting Help

### Documentation Issues
- Report: [GitHub Issues](https://github.com/sipeed/picoclaw/issues)
- Label: `documentation`

### Questions
- Ask: [GitHub Discussions](https://github.com/sipeed/picoclaw/discussions)
- Category: Q&A

### Bugs
- Report: [GitHub Issues](https://github.com/sipeed/picoclaw/issues)
- Include: Version, config, logs

---

**Last Updated**: 2026-03-09  
**Documentation Version**: 2.0.7
