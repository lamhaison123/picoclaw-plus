# Changelog

All notable changes to PicoClaw will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Comprehensive bug review and fixes
- Enhanced error logging for filesystem operations
- Nil safety checks for tool registry

## [1.1.1] - 2026-03-06

### Fixed
- **Grep Exit Code Handling**: grep exit code 1 (no matches) now treated as success instead of error
- **Email Detection**: Email addresses no longer incorrectly detected as @mentions
- **Filesystem Sync**: Added warning logs for directory sync errors
- **Nil Tools Registry**: Added nil check and initialization for subagent tools

### Changed
- **Message Formatting**: Removed session ID prefix from collaborative chat messages for cleaner UX
- **Enhanced /who Command**: Now shows all registered agents with descriptions and emojis

### Added
- Comprehensive test suite for grep exit code handling
- Email detection tests with 12 test cases
- Better error observability in filesystem operations

## [1.1.0] - 2026-03-05

### Added
- **IRC-Style Collaborative Chat**: Native multi-agent chat in Telegram with @mention routing
- **Platform-Agnostic Collaborative Package**: Refactored collaborative chat into reusable `pkg/collaborative/`
- **Team Roster in Prompts**: Agents automatically know about other team members
- **Enhanced /who Command**: Shows team status, all registered agents, and active agents
- **Session Management**: Per-chat sessions with context trimming
- **Parallel Agent Execution**: Multiple agents can respond simultaneously

### Changed
- Telegram channel now uses collaborative package for better maintainability
- Improved agent coordination and message routing

### Documentation
- Added `docs/COLLABORATIVE_CHAT.md` - Complete collaborative chat guide
- Added `docs/COLLABORATIVE_CHAT_QUICKSTART.md` - 5-minute quick start
- Added `pkg/collaborative/README.md` - API documentation
- Updated `CHANGELOG_MULTI_AGENT.md` with detailed changes

## [1.0.0] - 2025-12-XX

### Added
- Initial release of PicoClaw
- Ultra-lightweight design (<10MB RAM)
- Multi-agent collaboration framework
- 4-level safety system
- Multi-platform chat integration (Telegram, Discord, WhatsApp, QQ, etc.)
- Model Context Protocol (MCP) support
- Flexible LLM provider support (OpenAI, Anthropic, Zhipu, etc.)
- Cron job scheduling
- Skills system for extensibility

### Supported Platforms
- Linux (x86_64, ARM64, RISC-V)
- Windows (x86_64)
- macOS (x86_64, ARM64)
- Docker (Alpine-based and Node.js full-featured)

---

## Version History

- **v1.1.1** (2026-03-06) - Bug fixes and improvements
- **v1.1.0** (2026-03-05) - Collaborative chat and platform-agnostic refactoring
- **v1.0.0** (2025-12-XX) - Initial release

---

## Migration Guides

### Upgrading to v1.1.1
No breaking changes. All fixes are backward compatible.

### Upgrading to v1.1.0
No breaking changes. Collaborative chat is opt-in via configuration.

---

## Links

- [Multi-Agent Changelog](CHANGELOG_MULTI_AGENT.md) - Detailed multi-agent feature changes
- [Contributing Guide](CONTRIBUTING.md) - How to contribute
- [Documentation](docs/) - Full documentation
