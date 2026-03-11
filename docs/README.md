# PicoClaw Documentation

Welcome to the PicoClaw documentation! This directory contains comprehensive guides, references, and technical documentation.

## 📖 Documentation Index

**→ [INDEX.md](./INDEX.md)** - Complete documentation index with all available resources

## 🗂️ Directory Structure

```
docs/
├── INDEX.md              # Complete documentation index
├── README.md            # This file
│
├── architecture/        # System architecture and design
│   ├── README.md
│   ├── ARCHITECTURE_OVERVIEW.md
│   ├── COMPONENT_DETAILS.md
│   ├── DATA_FLOW.md
│   └── ...
│
├── guides/              # User guides and tutorials
│   ├── README.md
│   ├── MULTI_AGENT_GUIDE.md
│   ├── COLLABORATIVE_CHAT.md
│   ├── TEAM_AGENT_USAGE.md
│   └── ...
│
├── reference/           # API and configuration reference
│   ├── README.md
│   ├── API_REFERENCE.md
│   ├── tools_configuration.md
│   ├── SAFETY_LEVELS.md
│   └── ...
│
├── development/         # Development and troubleshooting
│   ├── README.md
│   ├── DEVELOPER_GUIDE.md
│   ├── COLLABORATIVE_CHAT_ARCHITECTURE.md
│   └── troubleshooting.md
│
├── channels/            # Platform-specific integrations
│   ├── telegram/
│   ├── discord/
│   ├── slack/
│   └── ...
│
├── memory/              # Memory system documentation
│   ├── CONFIG_GUIDE.md
│   └── vector/
│
├── skills/              # Skills and plugins
│   └── IRC_GATEWAY.md
│
├── design/              # Design documents and RFCs
│   └── ...
│
└── migration/           # Migration guides
    └── model-list-migration.md
```

## 🚀 Quick Start

### New Users
1. Read [Architecture Overview](./architecture/ARCHITECTURE_OVERVIEW.md)
2. Follow [Multi-Agent Guide](./guides/MULTI_AGENT_GUIDE.md)
3. Configure your [Channel](./channels/)

### Developers
1. Check [Developer Guide](./development/DEVELOPER_GUIDE.md)
2. Study [Component Details](./architecture/COMPONENT_DETAILS.md)
3. Review [Data Flow](./architecture/DATA_FLOW.md)

### System Administrators
1. Review [Safety Levels](./reference/SAFETY_LEVELS.md)
2. Configure [Tools](./reference/tools_configuration.md)
3. Setup [Systemd Service](../SYSTEMD_SETUP.md)

## 📚 Documentation Categories

### 🏗️ Architecture
Understanding PicoClaw's internal design and structure.
- System architecture
- Component specifications
- Data flow diagrams
- Design patterns

**→ [architecture/](./architecture/)**

### 📖 Guides
Step-by-step guides for using PicoClaw features.
- Getting started
- **v0.2.1 Features** (NEW)
- **Vision Support** (NEW)
- **Model Routing** (NEW)
- Multi-agent setup
- Team collaboration
- Security configuration

**→ [guides/](./guides/)**

### 📋 Reference
API documentation and configuration references.
- API endpoints
- Configuration options
- **JSONL Storage** (NEW)
- **Search Providers** (NEW)
- **PICOCLAW_HOME** (NEW)
- Security levels
- Tool system

**→ [reference/](./reference/)**

### 🔧 Development
Resources for developers and contributors.
- Development setup
- Code style
- Testing guidelines
- Troubleshooting

**→ [development/](./development/)**

### 🔌 Channels
Platform-specific integration guides.
- Telegram, Discord, Slack
- DingTalk, Feishu, WeCom
- QQ, LINE, OneBot
- MaixCAM

**→ [channels/](./channels/)**

### 💾 Memory
Vector memory and storage documentation.
- Configuration guide
- Vector store setup
- Embedding services

**→ [memory/](./memory/)**

### 🎯 Skills
Extending PicoClaw with skills and plugins.
- Skill development
- Plugin system
- IRC Gateway

**→ [skills/](./skills/)**

### 📐 Design
Technical design documents and proposals.
- RFCs
- Architecture decisions
- Refactoring plans

**→ [design/](./design/)**

### 🔄 Migration
Guides for upgrading between versions.
- Version migration
- Breaking changes
- Configuration updates

**→ [migration/](./migration/)**

## 🔍 Finding Documentation

### By Topic

| Topic | Document |
|-------|----------|
| **v0.2.1 Features** | [guides/V0.2.1_FEATURES.md](./guides/V0.2.1_FEATURES.md) |
| **Vision Support** | [guides/VISION_SUPPORT.md](./guides/VISION_SUPPORT.md) |
| **Model Routing** | [guides/MODEL_ROUTING.md](./guides/MODEL_ROUTING.md) |
| **JSONL Storage** | [reference/JSONL_STORAGE.md](./reference/JSONL_STORAGE.md) |
| **Search Providers** | [reference/SEARCH_PROVIDERS.md](./reference/SEARCH_PROVIDERS.md) |
| **MindGraph Memory** | [reference/MINDGRAPH_INTEGRATION.md](./reference/MINDGRAPH_INTEGRATION.md) |
| **Mem0 Memory** | [reference/MEM0_INTEGRATION.md](./reference/MEM0_INTEGRATION.md) |
| System Overview | [architecture/ARCHITECTURE_OVERVIEW.md](./architecture/ARCHITECTURE_OVERVIEW.md) |
| Getting Started | [guides/MULTI_AGENT_GUIDE.md](./guides/MULTI_AGENT_GUIDE.md) |
| API Reference | [reference/API_REFERENCE.md](./reference/API_REFERENCE.md) |
| Development | [development/DEVELOPER_GUIDE.md](./development/DEVELOPER_GUIDE.md) |
| Security | [reference/SAFETY_LEVELS.md](./reference/SAFETY_LEVELS.md) |
| Troubleshooting | [development/troubleshooting.md](./development/troubleshooting.md) |

### By Role

| Role | Start Here |
|------|------------|
| New User | [guides/MULTI_AGENT_GUIDE.md](./guides/MULTI_AGENT_GUIDE.md) |
| Developer | [development/DEVELOPER_GUIDE.md](./development/DEVELOPER_GUIDE.md) |
| System Admin | [reference/SAFETY_LEVELS.md](./reference/SAFETY_LEVELS.md) |
| Architect | [architecture/ARCHITECTURE_OVERVIEW.md](./architecture/ARCHITECTURE_OVERVIEW.md) |

## 📝 Documentation Standards

### Language
- Primary: English
- Secondary: Vietnamese (selected documents)
- Additional: Chinese, Japanese, French, Portuguese (README files)

### Format
- Markdown (.md) for all documentation
- ASCII diagrams for architecture
- Code examples in fenced blocks
- Clear section headers

### Structure
- Each directory has a README.md
- Related documents grouped together
- Clear navigation links
- Consistent formatting

## 🔗 External Resources

### Project Links
- **Repository**: https://github.com/sipeed/picoclaw
- **Issues**: https://github.com/sipeed/picoclaw/issues
- **Discussions**: https://github.com/sipeed/picoclaw/discussions
- **Releases**: https://github.com/sipeed/picoclaw/releases

### Community
- **Contributing**: [../CONTRIBUTING.md](../CONTRIBUTING.md)
- **Code of Conduct**: [../CODE_OF_CONDUCT.md](../CODE_OF_CONDUCT.md)
- **License**: [../LICENSE](../LICENSE)

## 📞 Support

### Documentation Issues
Found a problem with the documentation?
- Report: [GitHub Issues](https://github.com/sipeed/picoclaw/issues)
- Label: `documentation`

### Questions
Have questions about PicoClaw?
- Ask: [GitHub Discussions](https://github.com/sipeed/picoclaw/discussions)
- Category: Q&A

### Contributions
Want to improve the documentation?
- Read: [Contributing Guide](../CONTRIBUTING.md)
- Submit: Pull Request with `documentation` label

## 📊 Documentation Status

| Category | Status | Last Updated |
|----------|--------|--------------|
| Architecture | ✅ Complete | 2026-03-09 |
| Guides | ✅ Complete | 2026-03-09 |
| Reference | ✅ Complete | 2026-03-09 |
| Development | ✅ Complete | 2026-03-09 |
| Channels | ⚠️ Partial | Various |
| Memory | ✅ Complete | 2026-03-09 |
| Skills | ⚠️ Partial | Various |
| Design | ⚠️ Partial | Various |
| Migration | ⚠️ Partial | Various |

Legend:
- ✅ Complete and up-to-date
- ⚠️ Partial or needs updates
- ❌ Missing or outdated

## 🎯 Documentation Roadmap

### Planned Improvements
- [ ] Add more code examples
- [ ] Create video tutorials
- [ ] Translate more documents to Vietnamese
- [ ] Add troubleshooting flowcharts
- [ ] Expand channel-specific guides
- [ ] Add performance tuning guide
- [ ] Create deployment guide

### Recent Updates
- ✅ Added Mem0 integration (2026-03-09)
- ✅ Added MindGraph integration (2026-03-09)
- ✅ Added v0.2.1 features documentation (2026-03-09)
- ✅ Added Vision Support guide (2026-03-09)
- ✅ Added Model Routing guide (2026-03-09)
- ✅ Added JSONL Storage reference (2026-03-09)
- ✅ Added Search Providers reference (2026-03-09)
- ✅ Reorganized documentation structure (2026-03-09)
- ✅ Added comprehensive architecture docs (2026-03-09)
- ✅ Created category READMEs (2026-03-09)
- ✅ Added navigation index (2026-03-09)

## 📄 License

Documentation is licensed under [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/).  
Code examples in documentation follow the project's MIT License.

---

**Documentation Version**: 2.0.7  
**Last Updated**: 2026-03-09  
**Maintained by**: PicoClaw Contributors
