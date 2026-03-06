<div align="center">
  <img src="assets/logo.jpg" alt="PicoClaw" width="512">

  <h1>PicoClaw: Ultra-Efficient AI Assistant in Go</h1>

  <h3>$10 Hardware · 10MB RAM · 1s Boot · 皮皮虾，我们走！</h3>

  <p>
    <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go&logoColor=white" alt="Go">
    <img src="https://img.shields.io/badge/Arch-x86__64%2C%20ARM64%2C%20RISC--V-blue" alt="Hardware">
    <img src="https://img.shields.io/badge/license-MIT-green" alt="License">
    <br>
    <a href="https://picoclaw.io"><img src="https://img.shields.io/badge/Website-picoclaw.io-blue?style=flat&logo=google-chrome&logoColor=white" alt="Website"></a>
    <a href="https://x.com/SipeedIO"><img src="https://img.shields.io/badge/X_(Twitter)-SipeedIO-black?style=flat&logo=x&logoColor=white" alt="Twitter"></a>
    <br>
    <a href="./assets/wechat.png"><img src="https://img.shields.io/badge/WeChat-Group-41d56b?style=flat&logo=wechat&logoColor=white"></a>
    <a href="https://discord.gg/V4sAZ9XWpN"><img src="https://img.shields.io/badge/Discord-Community-4c60eb?style=flat&logo=discord&logoColor=white" alt="Discord"></a>
  </p>

[中文](README.zh.md) | [日本語](README.ja.md) | [Português](README.pt-br.md) | [Tiếng Việt](README.vi.md) | [Français](README.fr.md) | **English**

</div>

---

## 🚀 What is PicoClaw?

PicoClaw is an ultra-lightweight personal AI Assistant built in Go, designed to run on minimal hardware with maximum efficiency.

⚡️ **Runs on $10 hardware with <10MB RAM** — 99% less memory than OpenClaw and 98% cheaper than a Mac mini!

🦐 Inspired by [nanobot](https://github.com/HKUDS/nanobot), refactored from the ground up through AI-driven self-bootstrapping.

<table align="center">
  <tr align="center">
    <td align="center" valign="top">
      <p align="center">
        <img src="assets/picoclaw_mem.gif" width="360" height="240">
      </p>
    </td>
    <td align="center" valign="top">
      <p align="center">
        <img src="assets/licheervnano.png" width="400" height="240">
      </p>
    </td>
  </tr>
</table>

> [!CAUTION]
> **🚨 SECURITY & OFFICIAL CHANNELS**
>
> * **NO CRYPTO:** PicoClaw has **NO** official token/coin. All claims are **SCAMS**.
> * **OFFICIAL DOMAIN:** Only **[picoclaw.io](https://picoclaw.io)** and **[sipeed.com](https://sipeed.com)**
> * **Warning:** Early development - not production-ready before v1.0
> * **Note:** Recent PRs may increase memory footprint (10–20MB) temporarily

---

## ✨ Key Features

| Feature | Description |
|---------|-------------|
| 🪶 **Ultra-Lightweight** | <10MB RAM — 99% smaller than alternatives |
| 💰 **Minimal Cost** | Runs on $10 hardware — 98% cheaper |
| ⚡️ **Lightning Fast** | 1s boot time even on 0.6GHz single core |
| 🌍 **True Portability** | Single binary for RISC-V, ARM, x86 |
| 🤖 **AI-Bootstrapped** | 95% agent-generated with human refinement |
| 🤝 **Multi-Agent Teams** | Coordinate specialized AI agents |
| 💬 **Collaborative Chat** | IRC-style multi-agent conversations in Telegram |
| 🔓 **Flexible Safety** | 4-level safety system for LLM control |

---

## 📦 Quick Start

### 1. Install

**Download precompiled binary:**
```bash
# Download from https://github.com/sipeed/picoclaw/releases
wget https://github.com/sipeed/picoclaw/releases/download/v0.1.1/picoclaw-linux-amd64
chmod +x picoclaw-linux-amd64
```

**Or build from source:**
```bash
git clone https://github.com/sipeed/picoclaw.git
cd picoclaw
make build
```

### 2. Initialize

```bash
picoclaw onboard
```

### 3. Configure

Edit `~/.picoclaw/config.json`:

```json
{
  "agents": {
    "defaults": {
      "model": "gpt-5.2"
    }
  },
  "model_list": [
    {
      "model_name": "gpt-5.2",
      "model": "openai/gpt-5.2",
      "api_key": "your-api-key"
    }
  ]
}
```

**Get API Keys:**
- LLM: [OpenRouter](https://openrouter.ai/keys) · [Zhipu](https://open.bigmodel.cn) · [Anthropic](https://console.anthropic.com)
- Search (optional): [Tavily](https://tavily.com) · [Brave](https://brave.com/search/api)

### 4. Chat

```bash
picoclaw agent -m "What is 2+2?"
```

---

## 🤝 Multi-Agent Collaboration

Coordinate teams of specialized AI agents with role-based capabilities:

**Three Collaboration Patterns:**
- 🔄 **Sequential**: Tasks execute in order (design → implement → test → review)
- ⚡ **Parallel**: Tasks run simultaneously for speed
- 🌳 **Hierarchical**: Complex tasks decompose dynamically

**Quick Example:**

```bash
# Create development team
picoclaw team create templates/teams/development-team.json

# Execute task
picoclaw team execute dev-team-001 -t "Create a hello world function"

# Check status
picoclaw team status dev-team-001
```

**Key Features:**
- 👥 Role-based specialization with tool permissions
- 🗳️ Consensus voting (majority/unanimous/weighted)
- 🔄 Dynamic agent composition
- 📊 Comprehensive monitoring
- 💾 Automatic memory persistence

📖 **Learn More**: [Multi-Agent Guide](docs/MULTI_AGENT_GUIDE.md) | [Examples](examples/teams/)

---

## 💬 Collaborative Chat (NEW!)

**IRC-style multi-agent conversations in Telegram** — mention multiple agents in a single message and they'll all respond with full context!

```
User: @architect @developer How should we implement user authentication?

[abc123] 🏗️ ARCHITECT: I recommend using JWT tokens with...
[abc123] 💻 DEVELOPER: I can implement that using...
```

**Quick Setup:**

1. Enable in config:
```json
{
  "channels": {
    "telegram": {
      "collaborative_chat": {
        "enabled": true,
        "default_team_id": "dev-team",
        "max_context_length": 50
      }
    }
  }
}
```

2. Create team config (see [templates/teams/collaborative-dev-team.json](templates/teams/collaborative-dev-team.json))

3. Start gateway: `picoclaw gateway`

**Features:**
- 🎯 @mention-based routing (@architect, @developer, @tester)
- ⚡ Parallel agent execution
- 🧠 Shared conversation context
- 🎨 IRC-style formatting with emojis
- 📝 Session management per chat
- 👥 `/who` command - See all registered agents and active sessions

**Commands:**
- `/who` - Show team status, registered agents, and active agents
- `/help` - Show available commands

📖 **Learn More**: [Quick Start](docs/COLLABORATIVE_CHAT_QUICKSTART.md) | [Full Guide](docs/COLLABORATIVE_CHAT.md)

---

## 🔒 Security & Safety

### 4-Level Safety System

Choose the right balance between security and autonomy:

| Level | Best For | Blocks | Allows |
|-------|----------|--------|--------|
| **strict** | Production | sudo, chmod, docker, package install | Read, build, test, safe git |
| **moderate** | Development (default) | Catastrophic ops only | Most dev operations |
| **permissive** | DevOps/Admin | Only catastrophic ops | Almost everything |
| **off** | Testing ⚠️ | Nothing | Everything (DANGEROUS!) |

**Configuration:**

```json
{
  "tools": {
    "exec": {
      "safety_level": "moderate",
      "custom_allow_patterns": ["\\bgit\\s+push\\s+--force\\b"]
    }
  }
}
```

📖 **Full Documentation**: [Safety Levels Guide](docs/SAFETY_LEVELS.md) | [Quick Start](docs/SAFETY_QUICKSTART.md)

---

## 💬 Chat Apps Integration

Connect to Telegram, Discord, WhatsApp, QQ, DingTalk, LINE, WeCom, and more.

**Quick Setup (Telegram):**

1. Create bot with [@BotFather](https://t.me/BotFather)
2. Configure:
```json
{
  "channels": {
    "telegram": {
      "enabled": true,
      "token": "YOUR_BOT_TOKEN",
      "allow_from": ["YOUR_USER_ID"]
    }
  }
}
```
3. Run: `picoclaw gateway`

📖 **More Channels**: See [README sections](#-chat-apps) for Discord, WhatsApp, QQ, etc.

---

## � Docker Compose

```bash
# Clone repo
git clone https://github.com/sipeed/picoclaw.git
cd picoclaw

# First run (generates config)
docker compose -f docker/docker-compose.yml --profile gateway up

# Edit config
vim docker/data/config.json

# Start
docker compose -f docker/docker-compose.yml --profile gateway up -d
```

---

## ⚙️ Configuration

### Supported Providers

| Provider | Purpose | Get API Key |
|----------|---------|-------------|
| OpenAI | GPT models | [platform.openai.com](https://platform.openai.com) |
| Anthropic | Claude models | [console.anthropic.com](https://console.anthropic.com) |
| Zhipu | GLM models (Chinese) | [bigmodel.cn](https://bigmodel.cn) |
| OpenRouter | All models | [openrouter.ai](https://openrouter.ai) |
| Gemini | Google models | [aistudio.google.com](https://aistudio.google.com) |
| Groq | Fast inference | [console.groq.com](https://console.groq.com) |
| Ollama | Local models | Local (no key) |

### Model Configuration

```json
{
  "model_list": [
    {
      "model_name": "gpt-5.2",
      "model": "openai/gpt-5.2",
      "api_key": "sk-..."
    },
    {
      "model_name": "claude-sonnet-4.6",
      "model": "anthropic/claude-sonnet-4.6",
      "api_key": "sk-ant-..."
    },
    {
      "model_name": "llama3",
      "model": "ollama/llama3"
    }
  ]
}
```

---

## 📱 Deploy Anywhere

### Old Android Phones

```bash
# Install Termux from F-Droid
pkg install proot
wget https://github.com/sipeed/picoclaw/releases/download/v0.1.1/picoclaw-linux-arm64
chmod +x picoclaw-linux-arm64
termux-chroot ./picoclaw-linux-arm64 onboard
```

### Low-Cost Hardware

- $9.9 [LicheeRV-Nano](https://www.aliexpress.com/item/1005006519668532.html) - Minimal home assistant
- $30-50 [NanoKVM](https://www.aliexpress.com/item/1005007369816019.html) - Server maintenance
- $50 [MaixCAM](https://www.aliexpress.com/item/1005008053333693.html) - Smart monitoring

---

## 🛠️ CLI Reference

| Command | Description |
|---------|-------------|
| `picoclaw onboard` | Initialize config & workspace |
| `picoclaw agent -m "..."` | Chat with agent |
| `picoclaw agent` | Interactive mode |
| `picoclaw gateway` | Start gateway for chat apps |
| `picoclaw status` | Show status |
| `picoclaw team create <config>` | Create agent team |
| `picoclaw team list` | List teams |
| `picoclaw team status <id>` | Team status |
| `picoclaw cron list` | List scheduled jobs |

---

## 🤝 Contributing

PRs welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

**Key Points:**
- AI-assisted contributions are welcome (with disclosure)
- Run `make check` before submitting
- Keep PRs focused and small
- Fill out PR template completely

**Roadmap**: [ROADMAP.md](ROADMAP.md)

---

## 📊 Comparison

|  | OpenClaw | NanoBot | **PicoClaw** |
|---|---|---|---|
| **Language** | TypeScript | Python | **Go** |
| **RAM** | >1GB | >100MB | **<10MB** |
| **Startup** (0.8GHz) | >500s | >30s | **<1s** |
| **Cost** | Mac Mini $599 | Linux SBC ~$50 | **Any Linux $10+** |

---

## 📝 Documentation

- [Multi-Agent Guide](docs/MULTI_AGENT_GUIDE.md) - Team collaboration
- [Safety Levels](docs/SAFETY_LEVELS.md) - Security configuration
- [Tool Access Control](docs/TEAM_TOOL_ACCESS.md) - Permission system
- [Model Selection](docs/MULTI_AGENT_MODEL_SELECTION.md) - Per-role models
- [Contributing](CONTRIBUTING.md) - Contribution guidelines
- [Roadmap](ROADMAP.md) - Future plans
- [Changelog](CHANGELOG_MULTI_AGENT.md) - Multi-agent updates

---

## 🐛 Troubleshooting

**Web search not working?**
- Get free API key: [Brave Search](https://brave.com/search/api) (2000/month) or [Tavily](https://tavily.com) (1000/month)
- Or use DuckDuckGo (no key required, auto-fallback)

**Telegram bot conflict?**
- Only one `picoclaw gateway` instance can run at a time

**Content filtering errors?**
- Some providers (Zhipu) have strict filtering - try rephrasing or different model

---

## 📢 Community

- **Discord**: [Join Server](https://discord.gg/V4sAZ9XWpN)
- **WeChat**: <img src="assets/wechat.png" width="200">
- **Twitter**: [@SipeedIO](https://x.com/SipeedIO)
- **Website**: [picoclaw.io](https://picoclaw.io)

---

## 📄 License

MIT License - see [LICENSE](LICENSE) for details

---

**Made with ❤️ by the PicoClaw community**
