# 🦐 PicoClaw Roadmap

> **Vision**: Build the ultimate lightweight, secure, and fully autonomous AI Agent infrastructure

---

## 🚀 Core Optimization: Extreme Lightweight

*Fight software bloat to run on smallest embedded devices*

- [**Memory Footprint Reduction**](https://github.com/sipeed/picoclaw/issues/346)
  - **Goal**: Run on 64MB RAM boards with <20MB process
  - **Priority**: Memory > Storage (RAM is expensive on edge)
  - **Action**: Analyze growth, remove redundancies, optimize structures

---

## 🛡️ Security Hardening: Defense in Depth

*Build "Secure-by-Default" agent*

### Input Defense & Permission Control
- Prompt injection defense
- Tool abuse prevention
- SSRF protection (block internal IPs)

### Sandboxing & Isolation
- Filesystem sandbox (restrict R/W)
- Context isolation (prevent data leakage)
- Privacy redaction (auto-redact API keys, PII)

### Authentication & Secrets
- Modern crypto (ChaCha20-Poly1305)
- OAuth 2.0 flow (deprecate hardcoded keys)

---

## 🔌 Connectivity: Protocol-First Architecture

*Connect every model, reach every platform*

### Provider
- [**Architecture Upgrade**](https://github.com/sipeed/picoclaw/issues/283): Protocol-based (OpenAI-compatible, Ollama-compatible)
- **Local Models**: Ollama, vLLM, LM Studio, Mistral
- **Online Models**: Frontier closed-source models

### Channel
- **IM Matrix**: QQ, WeChat, DingTalk, Feishu, Telegram, Discord, WhatsApp, LINE, Slack, Email, Signal
- **Standards**: OneBot protocol
- [**Attachments**](https://github.com/sipeed/picoclaw/issues/348): Images, audio, video

### Skill Marketplace
- [**Discovery**](https://github.com/sipeed/picoclaw/issues/287): `find_skill` to auto-discover and install

---

## 🧠 Advanced Capabilities: From Chatbot to Agentic AI

*Beyond conversation—focus on action and collaboration*

### Operations
- [**MCP Support**](https://github.com/sipeed/picoclaw/issues/290): Model Context Protocol
- [**Browser Automation**](https://github.com/sipeed/picoclaw/issues/293): CDP (Chrome DevTools Protocol)
- [**Mobile Operation**](https://github.com/sipeed/picoclaw/issues/292): Android device control

### Multi-Agent Collaboration
- [**Basic Multi-Agent**](https://github.com/sipeed/picoclaw/issues/294): ✅ Implemented
- [**Model Routing**](https://github.com/sipeed/picoclaw/issues/295): Smart routing (simple→small/local, complex→SOTA)
- [**Swarm Mode**](https://github.com/sipeed/picoclaw/issues/284): Multiple PicoClaw instances collaboration
- [**AIEOS**](https://github.com/sipeed/picoclaw/issues/296): AI-Native OS interaction

---

## 📚 Developer Experience (DevEx)

*Lower barrier to entry*

- [**QuickGuide**](https://github.com/sipeed/picoclaw/issues/350): Interactive CLI wizard (zero-config start)
- **Comprehensive Docs**: Platform guides, step-by-step tutorials
- **AI-Assisted Docs**: Auto-generate API refs (with human verification)

---

## 🤖 Engineering: AI-Powered Open Source

*Use AI to accelerate development*

### AI-Enhanced CI/CD
- AI Code Review, Linting, PR Labeling
- **Bot Noise Reduction**: Clean PR timelines
- **Issue Triage**: AI agents analyze and suggest fixes

---

## 🎨 Brand & Community

- [**Logo Design**](https://github.com/sipeed/picoclaw/issues/297): Looking for **Mantis Shrimp** logo!
  - *Concept*: "Small but Mighty" + "Lightning Fast Strikes"

---

## 🤝 Call for Contributions

We welcome community contributions! Comment on Issues or submit PRs.

**Join Developer Group** after your first merged PR!

---

## 📊 Progress Tracking

| Category | Status | Priority |
|----------|--------|----------|
| Memory Optimization | 🟡 In Progress | 🔴 Critical |
| Security Hardening | 🟡 In Progress | 🔴 Critical |
| Provider Architecture | 🟡 In Progress | 🟠 High |
| Multi-Agent | ✅ Complete | 🟢 Done |
| MCP Support | 🔴 Planned | 🟠 High |
| Browser Automation | 🔴 Planned | 🟡 Medium |
| Swarm Mode | 🔴 Planned | 🟡 Medium |

---

**Let's build the best Edge AI Agent together!** 🦐
