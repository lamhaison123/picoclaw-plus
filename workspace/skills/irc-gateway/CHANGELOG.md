# IRC Gateway Changelog

All notable changes to the IRC Gateway skill will be documented in this file.

## [1.0.1] - 2026-03-05

### Changed
- **Automatic Token Loading**: Gateway now automatically reads Telegram token from `~/.picoclaw/config.json`
- Configuration priority: Environment variable → .env file → PicoClaw config
- Improved startup logging to show token source

### Improved
- No need to duplicate Telegram token if already configured in PicoClaw
- Better integration with existing PicoClaw Telegram setup
- Clearer error messages when token is missing

## [1.0.0] - 2026-03-05

### Added
- Initial release of IRC Gateway skill
- Mention-based role routing (@architect, @developer, @tester, @manager)
- Integration with PicoClaw team system via `picoclaw team execute`
- IRC-style message formatting `[session-id] emoji ROLE: message`
- Parallel role execution support
- Session management per Telegram chat
- Telegram bot commands: /start, /who, /status, /team, /clear
- Team configuration template (irc-dev-team.json)
- Comprehensive test suite (test_simple.py, test_gateway.py)
- Setup automation script (setup.sh)
- Windows launcher (start.bat)
- Full documentation (README.md, INTEGRATION.md)

### Features
- **Team Integration**: Routes to PicoClaw's pkg/team system
- **Parallel Processing**: Multiple roles execute simultaneously
- **Session Context**: Maintains conversation history per chat
- **Status Monitoring**: Real-time role busy/idle tracking
- **Flexible Configuration**: Customizable roles and team patterns
- **Per-Role Models**: Different LLM models for each role

### Architecture
- Lightweight Python gateway (~300 lines)
- CLI-based integration with PicoClaw
- No direct Go package dependencies
- Minimal memory footprint (~20MB)

### Documentation
- User guide: workspace/skills/irc-gateway/README.md
- Integration details: workspace/skills/irc-gateway/INTEGRATION.md
- Official docs: docs/skills/IRC_GATEWAY.md
- Team config template: workspace/skills/irc-gateway/irc-dev-team.json

### Testing
- Unit tests for mention parsing
- Session ID generation tests
- Role configuration validation
- Integration tests with mocked PicoClaw
- 12 test cases, 100% pass rate

### Requirements
- PicoClaw installed and configured
- Python 3.8+
- python-telegram-bot library
- Telegram bot token

### Known Limitations
- Requires PicoClaw team to be pre-configured
- CLI-based integration (not native Go)
- Limited to Telegram platform
- No built-in rate limiting (relies on Telegram API)

### Future Enhancements
- Direct Go integration using pkg/channels
- WebSocket support for real-time updates
- Rich media support (images, files)
- Consensus voting integration
- Task queuing for busy roles
- Metrics dashboard

## Version History

### [1.0.1] - 2026-03-05
- Automatic token loading from PicoClaw config
- Improved configuration flexibility
- Better startup logging

### [1.0.0] - 2026-03-05
- Initial stable release
- Full integration with PicoClaw v0.1.1+
- Production-ready for team collaboration workflows

---

## Upgrade Guide

### From Scratch to 1.0.0

1. **Install dependencies:**
   ```bash
   pip install python-telegram-bot
   ```

2. **Setup team:**
   ```bash
   cp workspace/skills/irc-gateway/irc-dev-team.json ~/.picoclaw/workspace/teams/
   picoclaw team create ~/.picoclaw/workspace/teams/irc-dev-team.json
   ```

3. **Configure environment:**
   ```bash
   echo "TELEGRAM_BOT_TOKEN=your_token" >> .env
   ```

4. **Run gateway:**
   ```bash
   cd workspace/skills/irc-gateway
   python gateway.py
   ```

## Support

For issues, questions, or contributions:
- GitHub Issues: https://github.com/sipeed/picoclaw/issues
- Documentation: docs/skills/IRC_GATEWAY.md
- Discord: https://discord.gg/V4sAZ9XWpN

## License

MIT License - Part of PicoClaw ecosystem
