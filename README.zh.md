<div align="center">
  <img src="assets/logo.jpg" alt="PicoClaw" width="512">

  <h1>PicoClaw: 基于Go语言的超高效 AI 助手</h1>

  <h3>10$硬件 · 10MB内存 · 1秒启动 · 皮皮虾，我们走！</h3>

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

**中文** | [日本語](README.ja.md) | [Português](README.pt-br.md) | [Tiếng Việt](README.vi.md) | [Français](README.fr.md) | [English](README.md)

</div>

---

## 🚀 什么是 PicoClaw？

PicoClaw 是一个用 Go 构建的超轻量级个人 AI 助手，旨在以最大效率在最小硬件上运行。

⚡️ **在 $10 硬件上运行，RAM <10MB** — 比 OpenClaw 节省 99% 内存，比 Mac mini 便宜 98%！

🦐 受 [nanobot](https://github.com/HKUDS/nanobot) 启发，通过 AI 驱动的自举过程从头重构。

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
> **🚨 安全声明 & 官方渠道**
>
> * **无加密货币:** PicoClaw **没有**发行任何官方代币。所有声称均为**诈骗**。
> * **官方域名:** 唯一官方网站是 **[picoclaw.io](https://picoclaw.io)**，公司官网是 **[sipeed.com](https://sipeed.com)**
> * **警告:** 正处于早期开发阶段 - v1.0 之前不建议部署到生产环境
> * **注意:** 最近的 PR 可能暂时增加内存占用（10-20MB）

---

## ✨ 核心特性

| 特性 | 描述 |
|---------|-------------|
| 🪶 **超轻量级** | <10MB RAM — 比替代方案小 99% |
| 💰 **最低成本** | 在 $10 硬件上运行 — 便宜 98% |
| ⚡️ **闪电启动** | 即使在 0.6GHz 单核上也能 1 秒启动 |
| 🌍 **真正可移植** | RISC-V、ARM、x86 单一二进制文件 |
| 🤖 **AI 自举** | 95% 代码由 Agent 生成，人工精炼 |
| 🤝 **多 Agent 团队** | 协调专业化 AI Agent |
| 💬 **协作聊天** | Telegram 中的 IRC 风格多 Agent 对话 |
| 🔓 **灵活安全** | 4 级安全系统控制 LLM |

---

## 📦 快速开始

### 1. 安装

**下载预编译二进制:**
```bash
# 从 https://github.com/sipeed/picoclaw/releases 下载
wget https://github.com/sipeed/picoclaw/releases/download/v0.1.1/picoclaw-linux-amd64
chmod +x picoclaw-linux-amd64
```

**或从源码构建:**
```bash
git clone https://github.com/sipeed/picoclaw.git
cd picoclaw
make build
```

### 2. 初始化

```bash
picoclaw onboard
```

### 3. 配置

编辑 `~/.picoclaw/config.json`:

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

**获取 API Keys:**
- LLM: [OpenRouter](https://openrouter.ai/keys) · [Zhipu](https://open.bigmodel.cn) · [Anthropic](https://console.anthropic.com)
- 搜索（可选）: [Tavily](https://tavily.com) · [Brave](https://brave.com/search/api)

### 4. 聊天

```bash
picoclaw agent -m "2+2 等于几？"
```

---

## 🤝 多 Agent 协作

协调具有基于角色能力的专业 AI Agent 团队：

**三种协作模式:**
- 🔄 **顺序**: 任务按顺序执行（设计 → 实现 → 测试 → 审查）
- ⚡ **并行**: 任务同时运行以提高速度
- 🌳 **层次**: 复杂任务动态分解

**快速示例:**

```bash
# 创建开发团队
picoclaw team create templates/teams/development-team.json

# 执行任务
picoclaw team execute dev-team-001 -t "创建一个 hello world 函数"

# 检查状态
picoclaw team status dev-team-001
```

**关键特性:**
- 👥 基于角色的专业化与工具权限
- 🗳️ 共识投票（多数/一致/加权）
- 🔄 动态 Agent 组合
- 📊 全面监控
- 💾 自动内存持久化

📖 **了解更多**: [多 Agent 指南](docs/MULTI_AGENT_GUIDE.md) | [示例](examples/teams/)

---

## 💬 协作聊天（新功能！）

**Telegram 中的 IRC 风格多 Agent 对话** — 在单条消息中提及多个 Agent，它们都会以完整上下文响应！

```
User: @architect @developer 我们应该如何实现用户认证？

[abc123] 🏗️ ARCHITECT: 我建议使用 JWT tokens...
[abc123] 💻 DEVELOPER: 我可以使用...实现
```

**快速设置:**

1. 在配置中启用:
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

2. 创建团队配置（参见 [templates/teams/collaborative-dev-team.json](templates/teams/collaborative-dev-team.json)）

3. 启动网关: `picoclaw gateway`

**特性:**
- 🎯 基于 @mention 的路由（@architect、@developer、@tester）
- ⚡ 并行 Agent 执行
- 🧠 共享对话上下文
- 🎨 带表情符号的 IRC 风格格式
- 📝 每个聊天的会话管理
- 👥 `/who` 命令 - 查看所有注册的 Agent 和活动会话

**命令:**
- `/who` - 显示团队状态、注册的 Agent 和活动的 Agent
- `/help` - 显示可用命令

📖 **了解更多**: [快速开始](docs/COLLABORATIVE_CHAT_QUICKSTART.md) | [完整指南](docs/COLLABORATIVE_CHAT.md)

---

## 🔒 安全与保护

### 4 级安全系统

在安全性和自主性之间选择合适的平衡：

| 级别 | 适用于 | 阻止 | 允许 |
|-------|----------|--------|--------|
| **strict** | 生产环境 | sudo、chmod、docker、包安装 | 读取、构建、测试、安全 git |
| **moderate** | 开发（默认） | 仅灾难性操作 | 大多数开发操作 |
| **permissive** | DevOps/管理员 | 仅灾难性操作 | 几乎所有操作 |
| **off** | 测试 ⚠️ | 无 | 所有操作（危险！） |

**配置:**

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

📖 **完整文档**: [安全级别指南](docs/SAFETY_LEVELS.md) | [快速开始](docs/SAFETY_QUICKSTART.md)

---

## 💬 聊天应用集成

连接到 Telegram、Discord、WhatsApp、QQ、DingTalk、LINE、WeCom 等。

**快速设置（Telegram）:**

1. 使用 [@BotFather](https://t.me/BotFather) 创建机器人
2. 配置:
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
3. 运行: `picoclaw gateway`

📖 **更多渠道**: 查看 [README 部分](#-chat-apps) 了解 Discord、WhatsApp、QQ 等

---

## 🐳 Docker Compose

```bash
# 克隆仓库
git clone https://github.com/sipeed/picoclaw.git
cd picoclaw

# 首次运行（生成配置）
docker compose -f docker/docker-compose.yml --profile gateway up

# 编辑配置
vim docker/data/config.json

# 启动
docker compose -f docker/docker-compose.yml --profile gateway up -d
```

---

## ⚙️ 配置

### 支持的提供商

| 提供商 | 用途 | 获取 API Key |
|----------|---------|-------------|
| OpenAI | GPT 模型 | [platform.openai.com](https://platform.openai.com) |
| Anthropic | Claude 模型 | [console.anthropic.com](https://console.anthropic.com) |
| Zhipu | GLM 模型（中文） | [bigmodel.cn](https://bigmodel.cn) |
| OpenRouter | 所有模型 | [openrouter.ai](https://openrouter.ai) |
| Gemini | Google 模型 | [aistudio.google.com](https://aistudio.google.com) |
| Groq | 快速推理 | [console.groq.com](https://console.groq.com) |
| Ollama | 本地模型 | 本地（无需密钥） |

### 模型配置

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

## 📱 随处部署

### 旧 Android 手机

```bash
# 从 F-Droid 安装 Termux
pkg install proot
wget https://github.com/sipeed/picoclaw/releases/download/v0.1.1/picoclaw-linux-arm64
chmod +x picoclaw-linux-arm64
termux-chroot ./picoclaw-linux-arm64 onboard
```

### 低成本硬件

- $9.9 [LicheeRV-Nano](https://www.aliexpress.com/item/1005006519668532.html) - 最小家庭助手
- $30-50 [NanoKVM](https://www.aliexpress.com/item/1005007369816019.html) - 服务器维护
- $50 [MaixCAM](https://www.aliexpress.com/item/1005008053333693.html) - 智能监控

---

## 🛠️ CLI 参考

| 命令 | 描述 |
|---------|-------------|
| `picoclaw onboard` | 初始化配置和工作区 |
| `picoclaw agent -m "..."` | 与 Agent 聊天 |
| `picoclaw agent` | 交互模式 |
| `picoclaw gateway` | 启动聊天应用网关 |
| `picoclaw status` | 显示状态 |
| `picoclaw team create <config>` | 创建 Agent 团队 |
| `picoclaw team list` | 列出团队 |
| `picoclaw team status <id>` | 团队状态 |
| `picoclaw cron list` | 列出计划任务 |

---

## 🤝 贡献

欢迎 PR！查看 [CONTRIBUTING.md](CONTRIBUTING.md) 了解指南。

**要点:**
- 欢迎 AI 辅助贡献（需披露）
- 提交前运行 `make check`
- 保持 PR 专注和小型
- 完整填写 PR 模板

**路线图**: [ROADMAP.md](ROADMAP.md)

---

## 📊 对比

|  | OpenClaw | NanoBot | **PicoClaw** |
|---|---|---|---|
| **语言** | TypeScript | Python | **Go** |
| **RAM** | >1GB | >100MB | **<10MB** |
| **启动** (0.8GHz) | >500s | >30s | **<1s** |
| **成本** | Mac Mini $599 | Linux SBC ~$50 | **任何 Linux $10+** |

---

## 📝 文档

- [多 Agent 指南](docs/MULTI_AGENT_GUIDE.md) - 团队协作
- [安全级别](docs/SAFETY_LEVELS.md) - 安全配置
- [工具访问控制](docs/TEAM_TOOL_ACCESS.md) - 权限系统
- [模型选择](docs/MULTI_AGENT_MODEL_SELECTION.md) - 按角色选择模型
- [贡献](CONTRIBUTING.md) - 贡献指南
- [路线图](ROADMAP.md) - 未来计划
- [更新日志](CHANGELOG_MULTI_AGENT.md) - 多 Agent 更新

---

## 🐛 故障排除

**网络搜索不工作？**
- 获取免费 API key: [Brave Search](https://brave.com/search/api)（2000/月）或 [Tavily](https://tavily.com)（1000/月）
- 或使用 DuckDuckGo（无需密钥，自动回退）

**Telegram 机器人冲突？**
- 一次只能运行一个 `picoclaw gateway` 实例

**内容过滤错误？**
- 某些提供商（如 Zhipu）有严格过滤 - 尝试重新表述或使用不同模型

---

## 📢 社区

- **Discord**: [加入服务器](https://discord.gg/V4sAZ9XWpN)
- **微信**: <img src="assets/wechat.png" width="200">
- **Twitter**: [@SipeedIO](https://x.com/SipeedIO)
- **网站**: [picoclaw.io](https://picoclaw.io)

---

## 📄 许可证

MIT 许可证 - 详见 [LICENSE](LICENSE)

---

**由 PicoClaw 社区用 ❤️ 制作**
