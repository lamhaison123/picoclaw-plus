# PicoClaw Documentation Index

Welcome to PicoClaw documentation! This guide will help you navigate through all available documentation.

## 🚀 Quick Start

**New to PicoClaw?** → [QUICK_START.md](./QUICK_START.md) - Fast-track guide based on your role

## 📚 Documentation Structure

```
docs/
├── architecture/     # System architecture and design
├── guides/          # User guides and tutorials
├── reference/       # API and configuration reference
├── development/     # Development and troubleshooting
├── channels/        # Platform-specific channel docs
├── memory/          # Memory system documentation
├── skills/          # Skills and plugins
├── design/          # Design documents and RFCs
└── migration/       # Migration guides
```

## 🏗️ Architecture

Understanding PicoClaw's internal design and structure.

- **[Architecture Overview](./architecture/ARCHITECTURE_OVERVIEW.md)** - Complete system architecture
- **[Component Details](./architecture/COMPONENT_DETAILS.md)** - Detailed component specifications
- **[Data Flow](./architecture/DATA_FLOW.md)** - Message processing flows
- **[Codebase Overview](./architecture/CODEBASE_OVERVIEW.md)** - Code organization
- **[Repository Overview](./architecture/REPOSITORY_OVERVIEW.md)** - Repository structure (EN)
- **[Repository Overview (VI)](./architecture/REPOSITORY_OVERVIEW.vi.md)** - Repository structure (Vietnamese)

## 📖 User Guides

Step-by-step guides for using PicoClaw features.

### v0.2.1 Features (NEW)
- **[v0.2.1 Features](./guides/V0.2.1_FEATURES.md)** - Complete guide to all v0.2.1 features
- **[Vision Support](./guides/VISION_SUPPORT.md)** - Multi-modal AI with images
- **[Model Routing](./guides/MODEL_ROUTING.md)** - Cost optimization

### Getting Started
- **[Multi-Agent Guide](./guides/MULTI_AGENT_GUIDE.md)** - Working with multiple agents
- **[Collaborative Chat](./guides/COLLABORATIVE_CHAT.md)** - Multi-agent conversations
- **[Collaborative Chat Quickstart](./guides/COLLABORATIVE_CHAT_QUICKSTART.md)** - Quick setup

### Team Features
- **[Team Agent Usage](./guides/TEAM_AGENT_USAGE.md)** - Using team delegation
- **[Safety Quickstart](./guides/SAFETY_QUICKSTART.md)** - Security configuration

### Advanced Features
- **[Antigravity Usage](./guides/ANTIGRAVITY_USAGE.md)** - Cloud integration
- **[Systemd Quick Reference](./guides/SYSTEMD_QUICK_REFERENCE.md)** - Linux service setup
- **[Compaction Quick Reference](./guides/COMPACTION_QUICK_REFERENCE.md)** - Memory compaction

## 📋 Reference

API documentation and configuration references.

### v0.2.1 Features (NEW)
- **[JSONL Storage](./reference/JSONL_STORAGE.md)** - Crash-safe storage
- **[Search Providers](./reference/SEARCH_PROVIDERS.md)** - All 7 search providers
- **[PICOCLAW_HOME](./reference/PICOCLAW_HOME.md)** - Custom home directory
- **[MindGraph Integration](./reference/MINDGRAPH_INTEGRATION.md)** - Knowledge graph memory provider

### API & Configuration
- **[API Reference](./reference/API_REFERENCE.md)** - Complete API documentation
- **[Tools Configuration](./reference/tools_configuration.md)** - Tool system configuration
- **[Safety Levels](./reference/SAFETY_LEVELS.md)** - Security levels explained

### System Components
- **[Circuit Breaker](./reference/CIRCUIT_BREAKER.md)** - Fault tolerance system
- **[Team Tool Access](./reference/TEAM_TOOL_ACCESS.md)** - Team tool permissions
- **[Multi-Agent Model Selection](./reference/MULTI_AGENT_MODEL_SELECTION.md)** - Model routing
- **[Antigravity Auth](./reference/ANTIGRAVITY_AUTH.md)** - Authentication system

## 🔧 Development

Resources for developers and contributors.

- **[Developer Guide](./development/DEVELOPER_GUIDE.md)** - Development setup and workflow
- **[Collaborative Chat Architecture](./development/COLLABORATIVE_CHAT_ARCHITECTURE.md)** - Chat system design
- **[Collaborative Chat Flow](./development/COLLABORATIVE_CHAT_FLOW.txt)** - Message flow diagram
- **[Troubleshooting](./development/troubleshooting.md)** - Common issues and solutions

## 🔌 Channels

Platform-specific integration guides.

- **[Telegram](./channels/telegram/)** - Telegram bot setup
- **[Discord](./channels/discord/)** - Discord bot integration
- **[Slack](./channels/slack/)** - Slack app configuration
- **[DingTalk](./channels/dingtalk/)** - DingTalk integration
- **[Feishu](./channels/feishu/)** - Feishu/Lark setup
- **[WeCom](./channels/wecom/)** - WeChat Work integration
- **[QQ](./channels/qq/)** - QQ bot setup
- **[LINE](./channels/line/)** - LINE bot integration
- **[OneBot](./channels/onebot/)** - OneBot protocol
- **[MaixCAM](./channels/maixcam/)** - MaixCAM device integration

## 💾 Memory System

Vector memory and storage documentation.

- **[Memory Configuration Guide](./memory/CONFIG_GUIDE.md)** - Memory system setup
- **[Vector Store](./memory/vector/)** - Vector database integration

## 🎯 Skills & Plugins

Extending PicoClaw with skills and plugins.

- **[IRC Gateway](./skills/IRC_GATEWAY.md)** - IRC bridge skill

## 📐 Design Documents

Technical design documents and proposals.

- **[Provider Refactoring](./design/provider-refactoring.md)** - LLM provider redesign
- **[Provider Refactoring Tests](./design/provider-refactoring-tests.md)** - Test strategy
- **[Issue 783 Investigation](./design/issue-783-investigation-and-fix-plan.zh.md)** - Bug analysis (Chinese)

## 🔄 Migration Guides

Guides for upgrading between versions.

- **[Model List Migration](./migration/model-list-migration.md)** - Model configuration updates

## 📝 Additional Resources

### Root Documentation
- **[README](../README.md)** - Project overview (English)
- **[README (中文)](../README.zh.md)** - Project overview (Chinese)
- **[README (日本語)](../README.ja.md)** - Project overview (Japanese)
- **[README (Français)](../README.fr.md)** - Project overview (French)
- **[README (Português)](../README.pt-br.md)** - Project overview (Portuguese)
- **[README (Tiếng Việt)](../README.vi.md)** - Project overview (Vietnamese)

### Project Management
- **[CHANGELOG](../CHANGELOG.md)** - Version history
- **[ROADMAP](../ROADMAP.md)** - Future plans
- **[CONTRIBUTING](../CONTRIBUTING.md)** - Contribution guidelines
- **[CONTRIBUTING (中文)](../CONTRIBUTING.zh.md)** - Contribution guidelines (Chinese)

### Release Notes
- **[v2.0.7](../RELEASE_NOTES_v2.0.7.md)** - Latest release
- **[v2.0.6](../RELEASE_NOTES_v2.0.6.md)** - Previous release
- **[v2.0.5](../RELEASE_NOTES_v2.0.5.md)** - Previous release
- **[v1.1.1](../RELEASE_NOTES_v1.1.1.md)** - Legacy release

### System Setup
- **[Systemd Setup](../SYSTEMD_SETUP.md)** - Linux service configuration
- **[Systemd README](../README_SYSTEMD.md)** - Systemd integration guide

### Experimental
- **[WASM Plugin POC](../WASM_PLUGIN_POC.md)** - WebAssembly plugin proof of concept
- **[Enhanced Metrics Testing](../ENHANCED_METRICS_TESTING_GUIDE.md)** - Metrics system testing
- **[Project Status](../PROJECT_STATUS.md)** - Current development status

## 🔍 Quick Links

### For New Users
1. Start with [README](../README.md)
2. Read [Architecture Overview](./architecture/ARCHITECTURE_OVERVIEW.md)
3. Follow [Multi-Agent Guide](./guides/MULTI_AGENT_GUIDE.md)
4. Configure your [Channel](./channels/)

### For Developers
1. Read [Developer Guide](./development/DEVELOPER_GUIDE.md)
2. Study [Component Details](./architecture/COMPONENT_DETAILS.md)
3. Review [Data Flow](./architecture/DATA_FLOW.md)
4. Check [Troubleshooting](./development/troubleshooting.md)

### For System Administrators
1. Review [Safety Levels](./reference/SAFETY_LEVELS.md)
2. Configure [Tools](./reference/tools_configuration.md)
3. Setup [Systemd Service](../SYSTEMD_SETUP.md)
4. Monitor with [Circuit Breaker](./reference/CIRCUIT_BREAKER.md)

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/sipeed/picoclaw/issues)
- **Discussions**: [GitHub Discussions](https://github.com/sipeed/picoclaw/discussions)
- **Documentation Issues**: Report in [Issues](https://github.com/sipeed/picoclaw/issues) with `documentation` label

## 📄 License

PicoClaw is licensed under the MIT License. See [LICENSE](../LICENSE) for details.

---

**Last Updated**: 2026-03-09  
**Documentation Version**: 2.0.7
