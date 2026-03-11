<div align="center">
  <img src="assets/logo.jpg" alt="PicoClaw" width="512">

  <h1>PicoClaw: Ultra-Efficient AI Assistant in Go</h1>

  <h3>$10 Hardware · 10MB RAM · 1s Boot · 皮皮虾， chúng ta đi!</h3>

  <p>
    <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go&logoColor=white" alt="Go">
    <img src="https://img.shields.io/badge/Arch-x86__64%2C%20ARM64%2C%20RISC--V-blue" alt="Hardware">
    <img src="https://img.shields.io/badge/Version-v2.1.0-orange" alt="Version">
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

⚡️ **Runs on $10 hardware with <50MB RAM** — Supports multi-agent cascades up to 20 levels!

🦐 Inspired by [nanobot](https://github.com/HKUDS/nanobot), refactored from the ground up through AI-driven self-bootstrapping.

---

## ✨ Key Features

| Feature | Description |
|---------|-------------|
| 🪶 **Ultra-Lightweight** | <10MB RAM — 99% smaller than alternatives |
| 🧠 **Vector Memory** | Long-term memory with **Qdrant** & **LanceDB** (Local/Cloud) |
| 👁️ **Vision Support** | Multi-modal image understanding (GPT-4o, Claude 3.5) |
| 🔍 **Advanced Search** | 7+ providers: Perplexity, Exa, SearXNG, GLM, Brave, Tavily |
| ⚡️ **Lightning Fast** | 1s boot time even on 0.6GHz single core |
| 🌍 **True Portability** | Single binary for RISC-V, ARM, x86 |
| 🤖 **AI-Bootstrapped** | 95% agent-generated with human refinement |
| 🤝 **Multi-Agent Teams** | Coordinate specialized AI agents with role-based capabilities |
| 🔒 **Flexible Safety** | 4-level safety system for LLM control |

---

## 🚀 Recent Improvements

### v2.1.0 (Latest) - The Memory & Vision Update

✅ **Advanced Memory Ecosystem**
- **Vector Search**: Integrated **Qdrant** (production) and **LanceDB** (local embedded).
- **Resilience**: Circuit breaker pattern & exponential backoff for all memory providers.
- **Mem0 & MindGraph**: Support for personalized memory and knowledge graphs.

✅ **Multi-Modal Vision**
- **Image Understanding**: Full support for OpenAI and Anthropic Vision models.
- **Efficient Streaming**: Base64 encoding with MIME auto-detection.

✅ **Enhanced Search & Routing**
- **7 Search Providers**: Added SearXNG, GLM Search, and Exa AI.
- **Smart Routing**: Complexity-based model selection to save costs.
- **Reasoning Extraction**: Support for Claude 3.5 "Extended Thinking" blocks.

✅ **Enterprise Resilience**
- **PICOCLAW_HOME**: Configurable home directory for Docker/Multi-tenant.
- **JSONL Storage**: Crash-safe, append-only session tracking.

📖 **Details**: [Release Notes v2.1](Release_note_2.1.md) | [Changelog](CHANGELOG.md)

---

## 📦 Quick Start

### 1. Install

**Download precompiled binary:**
```bash
# Download from https://github.com/lamhaison123/picoclaw-plus/releases
wget https://github.com/lamhaison123/picoclaw-plus/releases/download/v2.1/picoclaw-linux-amd64
chmod +x picoclaw-linux-amd64
```

### 2. Initialize

```bash
./picoclaw-linux-amd64 onboard
```

### 3. Configure

Edit `~/.picoclaw/config.json` to add your API keys for OpenAI, Anthropic, or Qdrant.

---

## 💬 Chat Apps Integration

Connect to Telegram, Discord, WhatsApp, WeCom, and more.

```bash
picoclaw gateway
```

---

## 📄 License

MIT License - see [LICENSE](LICENSE) for details

**Made with ❤️ by the PicoClaw community**
