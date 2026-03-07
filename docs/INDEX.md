# PicoClaw Documentation Index

Welcome to PicoClaw documentation! This index helps you find the right documentation for your needs.

---

## 🚀 Quick Start

New to PicoClaw? Start here:

1. [README](../README.md) - Project overview and quick start
2. [Installation Guide](../README.md#-quick-start) - Get PicoClaw running
3. [Configuration Guide](../config/README.md) - Configure your setup

---

## 📚 Core Documentation

### Getting Started
- [Repository Overview](REPOSITORY_OVERVIEW.md) - Project structure and architecture
- [Codebase Overview](CODEBASE_OVERVIEW.md) - Complete codebase structure and components
- [Architecture](ARCHITECTURE.md) - System architecture and design principles
- [Contributing Guide](../CONTRIBUTING.md) - How to contribute
- [Changelog](../CHANGELOG.md) - Version history

### Bug Fixes & Quality (v1.1.1+fixes)
- [Completion Report](../COMPLETION_REPORT.md) - Comprehensive summary of all fixes
- [Code Review Report](../CODE_REVIEW_REPORT.md) - Original findings (26 issues)
- [Critical Fixes Summary](../CRITICAL_FIXES_SUMMARY.md) - Critical issues resolved
- [All Fixes Summary](../ALL_FIXES_SUMMARY.md) - Complete fixes overview
- [Final Fix Report](../FINAL_FIX_REPORT.md) - Test results and verification

### Development
- [Developer Guide](DEVELOPER_GUIDE.md) - Development setup and guidelines
- [API Reference](API_REFERENCE.md) - Complete API documentation

### Multi-Agent System
- [Multi-Agent Guide](MULTI_AGENT_GUIDE.md) - Complete guide to multi-agent collaboration
- [Team Agent Usage](TEAM_AGENT_USAGE.md) - Using team agents
- [Team Tool Access](TEAM_TOOL_ACCESS.md) - Tool access control
- [Model Selection](MULTI_AGENT_MODEL_SELECTION.md) - Per-role model configuration

### Collaborative Chat
- [Collaborative Chat Guide](COLLABORATIVE_CHAT.md) - Complete IRC-style chat guide
- [Quick Start](COLLABORATIVE_CHAT_QUICKSTART.md) - 5-minute setup
- [Architecture](COLLABORATIVE_CHAT_ARCHITECTURE.md) - Technical architecture
- [Message Flow](COLLABORATIVE_CHAT_FLOW.txt) - Flow diagram

### Security & Safety
- [Safety Levels](SAFETY_LEVELS.md) - 4-level safety system
- [Safety Quick Start](SAFETY_QUICKSTART.md) - Quick reference
- [Tool Configuration](tools_configuration.md) - Configure tool safety

---

## 🔌 Integration Guides

### Chat Platforms

#### English Guides
- Telegram - See [main README](../README.md#-chat-apps-integration)
- Discord - See [main README](../README.md#-chat-apps-integration)
- WhatsApp - See [main README](../README.md#-chat-apps-integration)

#### Chinese Guides (中文指南)
- [Telegram](channels/telegram/README.zh.md)
- [Discord](channels/discord/README.zh.md)
- [DingTalk](channels/dingtalk/README.zh.md)
- [Feishu](channels/feishu/README.zh.md)
- [LINE](channels/line/README.zh.md)
- [MaixCAM](channels/maixcam/README.zh.md)
- [OneBot](channels/onebot/README.zh.md)
- [QQ](channels/qq/README.zh.md)
- [Slack](channels/slack/README.zh.md)
- [WeCom](channels/wecom/) - Multiple integration types

### LLM Providers
- [Antigravity Auth](ANTIGRAVITY_AUTH.md) - Authentication guide
- [Antigravity Usage](ANTIGRAVITY_USAGE.md) - Usage guide
- OpenAI, Anthropic, Zhipu, etc. - See [main README](../README.md#-configuration)

---

## 🛠️ Developer Documentation

### API Documentation
- [Collaborative Package API](../pkg/collaborative/README.md) - Platform-agnostic collaborative chat API
- [Channels Package](../pkg/channels/README.md) - Channel integration API

### Design Documents
- [Provider Refactoring](design/provider-refactoring.md) - Provider architecture
- [Provider Tests](design/provider-refactoring-tests.md) - Testing strategy

### Migration Guides
- [Model List Migration](migration/model-list-migration.md) - Migrate model configurations

---

## 🔧 Configuration

### Templates
- [Team Templates](../templates/teams/README.md) - Pre-built team configurations
- [Config Examples](../config/README.md) - Configuration examples

### Example Configurations
- `config.example.json` - Full configuration example
- `safety_examples.json` - Safety level examples
- `collaborative-chat-*.json` - Collaborative chat examples

---

## 🐛 Troubleshooting

- [Troubleshooting Guide](troubleshooting.md) - Common issues and solutions
- [GitHub Issues](https://github.com/sipeed/picoclaw/issues) - Report bugs
- [Discord Community](https://discord.gg/V4sAZ9XWpN) - Get help

---

## 📖 Additional Resources

### Multi-Language READMEs
- [English](../README.md)
- [中文](../README.zh.md)
- [日本語](../README.ja.md)
- [Português](../README.pt-br.md)
- [Tiếng Việt](../README.vi.md)
- [Français](../README.fr.md)

### Repository Overviews
- [English](REPOSITORY_OVERVIEW.md)
- [Tiếng Việt](REPOSITORY_OVERVIEW.vi.md)

### Changelogs
- [Main Changelog](../CHANGELOG.md) - All versions
- [Multi-Agent Changelog](../CHANGELOG_MULTI_AGENT.md) - Multi-agent features

### Release Notes
- [v1.1.1 Release Notes](../RELEASE_NOTES_v1.1.1.md) - Latest release

---

## 🎯 Documentation by Use Case

### I want to...

#### Run PicoClaw
→ Start with [README](../README.md) and [Quick Start](../README.md#-quick-start)

#### Set up multi-agent collaboration
→ Read [Multi-Agent Guide](MULTI_AGENT_GUIDE.md) and [Team Agent Usage](TEAM_AGENT_USAGE.md)

#### Enable collaborative chat in Telegram
→ Follow [Collaborative Chat Quick Start](COLLABORATIVE_CHAT_QUICKSTART.md)

#### Integrate with a chat platform
→ Check [Integration Guides](#-integration-guides) for your platform

#### Configure safety levels
→ Read [Safety Levels](SAFETY_LEVELS.md) and [Safety Quick Start](SAFETY_QUICKSTART.md)

#### Contribute to the project
→ Read [Contributing Guide](../CONTRIBUTING.md)

#### Understand the architecture
→ Read [Repository Overview](REPOSITORY_OVERVIEW.md) and [Architecture](COLLABORATIVE_CHAT_ARCHITECTURE.md)

#### Troubleshoot issues
→ Check [Troubleshooting Guide](troubleshooting.md)

---

## 📝 Documentation Status

### ✅ Complete & Up-to-date
- Core documentation (README, guides)
- Multi-agent documentation
- Collaborative chat documentation
- Safety system documentation
- API documentation

### ⚠️ Needs Translation
- Channel integration guides (mostly Chinese-only)
- Some design documents

### 🔄 In Progress
- Video tutorials
- Interactive examples
- Expanded troubleshooting

---

## 🤝 Contributing to Documentation

Found an issue or want to improve the docs? We welcome contributions!

1. Check [Contributing Guide](../CONTRIBUTING.md)
2. Open an issue or PR on [GitHub](https://github.com/sipeed/picoclaw)
3. Join our [Discord](https://discord.gg/V4sAZ9XWpN) to discuss

---

## 📧 Support

- **GitHub Issues**: https://github.com/sipeed/picoclaw/issues
- **Discord**: https://discord.gg/V4sAZ9XWpN
- **Twitter**: @SipeedIO
- **Website**: https://picoclaw.io

---

**Last Updated:** 2026-03-07  
**Version:** 1.1.1+fixes
