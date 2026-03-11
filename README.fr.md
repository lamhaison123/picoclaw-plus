<div align="center">
  <img src="assets/logo.jpg" alt="PicoClaw" width="512">

  <h1>PicoClaw : Assistant IA Ultra-Efficace en Go</h1>

  <h3>Matériel à 10$ · 10 Mo de RAM · Démarrage en 1s · 皮皮虾， chúng ta đi!</h3>

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

[Français] | [English](README.md) | [中文](README.zh.md) | [日本語](README.ja.md) | [Português](README.pt-br.md) | [Tiếng Việt](README.vi.md)

</div>

---

## 🚀 Qu'est-ce que PicoClaw ?

PicoClaw est un assistant personnel IA ultra-léger construit en Go, conçu pour fonctionner sur du matériel minimal avec une efficacité maximale.

⚡️ **Fonctionne sur du matériel à 10$ avec <50 Mo de RAM** — Supporte des cascades multi-agents jusqu'à 20 niveaux !

🦐 Inspiré de [nanobot](https://github.com/HKUDS/nanobot), refactorisé de zéro via un processus d'auto-amorçage piloté par l'IA.

---

## ✨ Fonctionnalités Clés

| Fonctionnalité | Description |\
|---------|-------------|\
| 🪶 **Ultra-Léger** | <10 Mo RAM — 99% plus petit que les alternatives |\
| 🧠 **Mémoire Vectorielle** | Mémoire à long terme avec **Qdrant** & **LanceDB** (Local/Cloud) |\
| 👁️ **Support Vision** | Compréhension d'image multi-modale (GPT-4o, Claude 3.5) |\
| 🔍 **Recherche Avancée** | 7+ fournisseurs : Perplexity, Exa, SearXNG, GLM, Brave, Tavily |\
| ⚡️ **Ultra Rapide** | Démarrage en 1s même sur CPU single-core 0.6GHz |\
| 🌍 **Vraiment Portable** | Binaire unique pour RISC-V, ARM, x86 |\
| 🤖 **Auto-Construit par IA** | 95% du code généré par agent avec raffinement humain |\
| 🤝 **Équipes Multi-Agents** | Coordonnez des agents IA spécialisés avec des capacités basées sur les rôles |\
| 🔒 **Sécurité Flexible** | Système de sécurité à 4 niveaux pour contrôler le LLM |

---

## 🚀 Améliorations Récentes

### v2.1.0 (Dernière) - Mise à jour Mémoire & Vision

✅ **Écosystème de Mémoire Avancé**
- **Recherche Vectorielle** : Intégration de **Qdrant** (production) et **LanceDB** (local embarqué).
- **Résilience** : Modèle Circuit Breaker & Exponential Backoff pour tous les fournisseurs de mémoire.
- **Mem0 & MindGraph** : Support pour la mémoire personnalisée et les graphes de connaissances.

✅ **Vision Multi-Modale**
- **Compréhension d'Image** : Support complet pour les modèles Vision d'OpenAI et Anthropic.
- **Streaming Efficace** : Encodage Base64 avec auto-détection MIME.

✅ **Recherche & Routage Améliorés**
- **7 Fournisseurs de Recherche** : Ajout de SearXNG, GLM Search et Exa AI.
- **Routage Intelligent** : Sélection du modèle basée sur la complexité pour économiser les coûts.
- **Extraction de Raisonnement** : Support pour les blocs "Extended Thinking" de Claude 3.5.

✅ **Résilience Entreprise**
- **PICOCLAW_HOME** : Répertoire home configurable pour Docker/Multi-tenant.
- **Stockage JSONL** : Suivi de session sécurisé contre les crashs, en mode ajout uniquement.

📖 **Détails** : [Notes de version v2.1](Release_note_2.1.md) | [Journal des modifications](CHANGELOG.md)

---

## 📦 Démarrage Rapide

### 1. Installer

**Télécharger le binaire précompilé :**
```bash
# Télécharger depuis GitHub Releases
wget https://github.com/lamhaison123/picoclaw-plus/releases/download/v2.1/picoclaw-linux-amd64
chmod +x picoclaw-linux-amd64
```

### 2. Initialiser

```bash
./picoclaw-linux-amd64 onboard
```

### 3. Configurer

Éditer `~/.picoclaw/config.json` pour ajouter vos clés API pour OpenAI, Anthropic ou Qdrant.

---

## 💬 Intégration d'Applications de Chat

Connectez-vous à Telegram, Discord, WhatsApp, WeCom et plus.

```bash
picoclaw gateway
```

---

## 📄 Licence

Licence MIT - voir [LICENSE](LICENSE) pour les détails

**Fait avec ❤️ par la communauté PicoClaw**
