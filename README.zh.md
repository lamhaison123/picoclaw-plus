<div align="center">
  <img src="assets/logo.jpg" alt="PicoClaw" width="512">

  <h1>PicoClaw: 基于 Go 语言的超高效 AI 助手</h1>

  <h3>$10 硬件 · 10MB RAM · 1秒启动 · 皮皮虾，我们走！</h3>

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

**中文** | [English](README.md) | [日本語](README.ja.md) | [Português](README.pt-br.md) | [Tiếng Việt](README.vi.md) | [Français](README.fr.md)

</div>

---

## 🚀 什么是 PicoClaw？

PicoClaw 是一个用 Go 构建的超轻量级个人 AI 助手，旨在以最大效率在最小硬件上运行。

⚡️ **在 $10 硬件上运行，RAM <50MB** — 支持高达 20 级的多 Agent 级联！

🦐 受 [nanobot](https://github.com/HKUDS/nanobot) 启发，通过 AI 驱动的自举过程从头重构。

---

## ✨ 核心特性

| 特性 | 描述 |
|---------|-------------|
| 🪶 **超轻量级** | <10MB RAM — 比替代方案小 99% |
| 🧠 **向量内存** | 使用 **Qdrant** & **LanceDB** (本地/云端) 的长期记忆 |
| 👁️ **视觉支持** | 多模态图像理解 (GPT-4o, Claude 3.5) |
| 🔍 **高级搜索** | 7+ 提供商：Perplexity, Exa, SearXNG, GLM, Brave, Tavily |
| ⚡️ **闪电启动** | 即使在 0.6GHz 单核上也能 1 秒启动 |
| 🌍 **真正可移植** | RISC-V、ARM、x86 单一二进制文件 |
| 🤖 **AI 自举** | 95% 代码由 Agent 生成，人工精炼 |
| 🤝 **多 Agent 团队** | 协调具有基于角色能力的专业化 AI Agent |
| 🔒 **灵活安全** | 4 级安全系统控制 LLM |

---

## 🚀 最近改进

### v2.1.0 (最新) - 内存与视觉更新

✅ **高级内存生态系统**
- **向量搜索**：集成了 **Qdrant** (生产级) 和 **LanceDB** (本地嵌入式)。
- **弹性**：所有内存提供商均支持断路器模式和指数退避。
- **Mem0 & MindGraph**：支持个性化记忆和知识图谱。

✅ **多模态视觉**
- **图像理解**：全面支持 OpenAI 和 Anthropic 的视觉模型。
- **高效流式传输**：带 MIME 自动检测的 Base64 编码。

✅ **增强的搜索与路由**
- **7 个搜索提供商**：新增了 SearXNG, GLM Search 和 Exa AI。
- **智能路由**：根据任务复杂度选择模型以节省成本。
- **推理提取**：支持 Claude 3.5 的“深度思考” (Extended Thinking) 模块。

✅ **企业级韧性**
- **PICOCLAW_HOME**：可配置的主目录，适用于 Docker/多租户环境。
- **JSONL 存储**：防崩溃、仅追加的会话追踪。

📖 **详情**: [发布说明 v2.1](Release_note_2.1.md) | [更新日志](CHANGELOG.md)

---

## 📦 快速开始

### 1. 安装

**下载预编译二进制文件:**
```bash
# 从 GitHub Releases 下载
wget https://github.com/lamhaison123/picoclaw-plus/releases/download/v2.1/picoclaw-linux-amd64
chmod +x picoclaw-linux-amd64
```

### 2. 初始化

```bash
./picoclaw-linux-amd64 onboard
```

### 3. 配置

编辑 `~/.picoclaw/config.json` 以添加 OpenAI, Anthropic 或 Qdrant 的 API Key。

---

## 💬 聊天应用集成

连接到 Telegram, Discord, WhatsApp, 微信企业号等。

```bash
picoclaw gateway
```

---

## 📄 许可证

MIT 许可证 - 详见 [LICENSE](LICENSE)

**由 PicoClaw 社区用 ❤️ 制作**
