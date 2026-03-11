<div align="center">
  <img src="assets/logo.jpg" alt="PicoClaw" width="512">

  <h1>PicoClaw: Assistente de IA Ultra-Eficiente em Go</h1>

  <h3>Hardware de $10 · 10MB de RAM · Boot em 1s · 皮皮虾， chúng ta đi!</h3>

  <p>
    <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go&logoColor=white" alt="Go">
    <img src="https://img.shields.io/badge/Arch-x86__64%2C%20ARM64%2C%20RISC--V-blue" alt="Hardware">
    <img src="https://img.shields.io/badge/Versão-v2.1.0-orange" alt="Version">
    <img src="https://img.shields.io/badge/license-MIT-green" alt="License">
    <br>
    <a href="https://picoclaw.io"><img src="https://img.shields.io/badge/Website-picoclaw.io-blue?style=flat&logo=google-chrome&logoColor=white" alt="Website"></a>
    <a href="https://x.com/SipeedIO"><img src="https://img.shields.io/badge/X_(Twitter)-SipeedIO-black?style=flat&logo=x&logoColor=white" alt="Twitter"></a>
    <br>
    <a href="./assets/wechat.png"><img src="https://img.shields.io/badge/WeChat-Group-41d56b?style=flat&logo=wechat&logoColor=white"></a>
    <a href="https://discord.gg/V4sAZ9XWpN"><img src="https://img.shields.io/badge/Discord-Community-4c60eb?style=flat&logo=discord&logoColor=white" alt="Discord"></a>
  </p>

[Português] | [English](README.md) | [中文](README.zh.md) | [日本語](README.ja.md) | [Tiếng Việt](README.vi.md) | [Français](README.fr.md)

</div>

---

## 🚀 O que é PicoClaw?

PicoClaw é um assistente pessoal de IA ultra-leve construído em Go, projetado para rodar em hardware mínimo com máxima eficiência.

⚡️ **Roda em hardware de $10 com <50MB de RAM** — Suporta cascatas multi-agente de até 20 níveis!

🦐 Inspirado no [nanobot](https://github.com/HKUDS/nanobot), refatorado do zero através de um processo de auto-inicialização (self-bootstrapping) conduzido por IA.

---

## ✨ Recursos Principais

| Recurso | Descrição |
|---------|-------------|
| 🪶 **Ultra-Leve** | <10MB RAM — 99% menor que alternativas |
| 🧠 **Memória Vector** | Memória de longo prazo com **Qdrant** & **LanceDB** (Local/Cloud) |
| 👁️ **Suporte a Visão** | Compreensão de imagem multi-modal (GPT-4o, Claude 3.5) |
| 🔍 **Busca Avançada** | 7+ provedores: Perplexity, Exa, SearXNG, GLM, Brave, Tavily |
| ⚡️ **Super Rápido** | Boot em 1s mesmo em CPU single-core de 0.6GHz |
| 🌍 **Verdadeiramente Portátil** | Binário único para RISC-V, ARM, x86 |
| 🤖 **Auto-Construído por IA** | 95% do código gerado por agente com refinamento humano |
| 🤝 **Times Multi-Agente** | Coordene agentes de IA especializados com capacidades baseadas em funções |
| 🔒 **Segurança Flexível** | Sistema de segurança de 4 níveis para controle de LLM |

---

## 🚀 Melhorias Recentes

### v2.1.0 (Mais Recente) - O Update de Memória & Visão

✅ **Ecossistema de Memória Avançado**
- **Busca Vector**: Integrado **Qdrant** (produção) e **LanceDB** (embutido local).
- **Resiliência**: Padrão Circuit Breaker & Exponential Backoff para todos os provedores de memória.
- **Mem0 & MindGraph**: Suporte para memória personalizada e grafos de conhecimento.

✅ **Visão Multi-Modal**
- **Compreensão de Imagem**: Suporte completo para modelos Vision da OpenAI e Anthropic.
- **Streaming Eficiente**: Codificação Base64 com auto-detecção MIME.

✅ **Busca & Roteamento Aprimorados**
- **7 Provedores de Busca**: Adicionados SearXNG, GLM Search e Exa AI.
- **Roteamento Inteligente**: Seleção de modelo baseada na complexidade para economizar custos.
- **Extração de Raciocínio**: Suporte para blocos "Extended Thinking" do Claude 3.5.

✅ **Resiliência Empresarial**
- **PICOCLAW_HOME**: Diretório home configurável para Docker/Multi-tenant.
- **Armazenamento JSONL**: Rastreamento de sessão anexo-apenas e seguro contra quedas.

📖 **Detalhes**: [Release Notes v2.1](Release_note_2.1.md) | [Changelog](CHANGELOG.md)

---

## 📦 Início Rápido

### 1. Instalar

**Baixar binário pré-compilado:**
```bash
# Baixe do GitHub Releases
wget https://github.com/lamhaison123/picoclaw-plus/releases/download/v2.1/picoclaw-linux-amd64
chmod +x picoclaw-linux-amd64
```

### 2. Inicializar

```bash
./picoclaw-linux-amd64 onboard
```

### 3. Configurar

Edite `~/.picoclaw/config.json` para adicionar suas chaves API para OpenAI, Anthropic ou Qdrant.

---

## 💬 Integração com Apps de Chat

Conecte ao Telegram, Discord, WhatsApp, WeCom e mais.

```bash
picoclaw gateway
```

---

## 📄 Licença

Licença MIT - veja [LICENSE](LICENSE) para detalhes

**Feito com ❤️ pela comunidade PicoClaw**
