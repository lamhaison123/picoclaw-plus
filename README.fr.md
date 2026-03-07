<div align="center">
  <img src="assets/logo.jpg" alt="PicoClaw" width="512">

  <h1>PicoClaw : Assistant IA Ultra-Efficace en Go</h1>

  <h3>Matériel à 10$ · 10 Mo de RAM · Démarrage en 1s · 皮皮虾，我们走！</h3>

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

[中文](README.zh.md) | [日本語](README.ja.md) | [Português](README.pt-br.md) | [Tiếng Việt](README.vi.md) | **Français** | [English](README.md)

</div>

---

## 🚀 Qu'est-ce que PicoClaw ?

PicoClaw est un assistant personnel IA ultra-léger construit en Go, conçu pour fonctionner sur du matériel minimal avec une efficacité maximale.

⚡️ **Fonctionne sur du matériel à 10$ avec <10 Mo de RAM** — 99% moins de mémoire qu'OpenClaw et 98% moins cher qu'un Mac mini !

🦐 Inspiré de [nanobot](https://github.com/HKUDS/nanobot), refactorisé de zéro via un processus d'auto-amorçage piloté par l'IA.

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
> **🚨 SÉCURITÉ & CANAUX OFFICIELS**
>
> * **PAS DE CRYPTO :** PicoClaw n'a **AUCUN** token/jeton officiel. Toutes les affirmations sont des **ARNAQUES**.
> * **DOMAINE OFFICIEL :** Uniquement **[picoclaw.io](https://picoclaw.io)** et **[sipeed.com](https://sipeed.com)**
> * **Avertissement :** En phase de développement précoce - non recommandé pour la production avant v1.0
> * **Note :** Les PR récentes peuvent augmenter temporairement l'utilisation de la mémoire à 10-20 Mo

---

## ✨ Fonctionnalités Clés

| Fonctionnalité | Description |
|---------|-------------|
| 🪶 **Ultra-Léger** | <10 Mo RAM — 99% plus petit que les alternatives |
| 💰 **Coût Minimal** | Fonctionne sur du matériel à 10$ — 98% moins cher |
| ⚡️ **Ultra Rapide** | Démarrage en 1s même sur CPU single-core 0.6GHz |
| 🌍 **Vraiment Portable** | Binaire unique pour RISC-V, ARM, x86 |
| 🤖 **Auto-Construit par IA** | 95% du code généré par agent avec raffinement humain |
| 🤝 **Équipes Multi-Agents** | Coordonnez des agents IA spécialisés |
| 💬 **Chat Collaboratif** | Conversations multi-agents style IRC dans Telegram |
| 🔓 **Sécurité Flexible** | Système de sécurité à 4 niveaux pour contrôler le LLM |

---

## 📦 Démarrage Rapide

### 1. Installer

**Télécharger le binaire précompilé :**
```bash
# Télécharger depuis https://github.com/sipeed/picoclaw/releases
wget https://github.com/sipeed/picoclaw/releases/download/v0.1.1/picoclaw-linux-amd64
chmod +x picoclaw-linux-amd64
```

**Ou compiler depuis les sources :**
```bash
git clone https://github.com/sipeed/picoclaw.git
cd picoclaw
make build
```

### 2. Initialiser

```bash
picoclaw onboard
```

### 3. Configurer

Éditer `~/.picoclaw/config.json` :

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

**Obtenir les Clés API :**
- LLM : [OpenRouter](https://openrouter.ai/keys) · [Zhipu](https://open.bigmodel.cn) · [Anthropic](https://console.anthropic.com)
- Recherche (optionnel) : [Tavily](https://tavily.com) · [Brave](https://brave.com/search/api)

### 4. Discuter

```bash
picoclaw agent -m "Combien font 2+2 ?"
```

---

## 🤝 Collaboration Multi-Agents

Coordonnez des équipes d'agents IA spécialisés avec des capacités basées sur les rôles :

**Trois Modèles de Collaboration :**
- 🔄 **Séquentiel** : Les tâches s'exécutent dans l'ordre (conception → implémentation → test → revue)
- ⚡ **Parallèle** : Les tâches s'exécutent simultanément pour la vitesse
- 🌳 **Hiérarchique** : Les tâches complexes se décomposent dynamiquement

**Exemple Rapide :**

```bash
# Créer une équipe de développement
picoclaw team create templates/teams/development-team.json

# Exécuter une tâche
picoclaw team execute dev-team-001 -t "Créer une fonction hello world"

# Vérifier le statut
picoclaw team status dev-team-001
```

**Fonctionnalités Principales :**
- 👥 Spécialisation basée sur les rôles avec permissions d'outils
- 🗳️ Vote par consensus (majorité/unanime/pondéré)
- 🔄 Composition dynamique d'agents
- 📊 Surveillance complète
- 💾 Persistance automatique de la mémoire

📖 **En savoir plus** : [Guide Multi-Agents](docs/MULTI_AGENT_GUIDE.md) | [Exemples](examples/teams/)

---

## 💬 Chat Collaboratif (NOUVEAU !)

**Conversations multi-agents style IRC dans Telegram** — mentionnez plusieurs agents dans un seul message et tous répondront avec le contexte complet !

```
User: @architect @developer Comment devrions-nous implémenter l'authentification utilisateur ?

[abc123] 🏗️ ARCHITECT: Je recommande d'utiliser des tokens JWT avec...
[abc123] 💻 DEVELOPER: Je peux implémenter cela en utilisant...
```

**Configuration Rapide :**

1. Activer dans la configuration :
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

2. Créer la configuration d'équipe (voir [templates/teams/collaborative-dev-team.json](templates/teams/collaborative-dev-team.json))

3. Démarrer la passerelle : `picoclaw gateway`

**Fonctionnalités :**
- 🎯 Routage basé sur @mention (@architect, @developer, @tester)
- ⚡ Exécution parallèle des agents
- 🧠 Contexte de conversation partagé
- 🎨 Formatage style IRC avec emojis
- 📝 Gestion de session par chat
- 👥 Commande `/who` - Voir tous les agents enregistrés et sessions actives

**Commandes :**
- `/who` - Afficher le statut de l'équipe, les agents enregistrés et les agents actifs
- `/help` - Afficher les commandes disponibles

📖 **En savoir plus** : [Démarrage Rapide](docs/COLLABORATIVE_CHAT_QUICKSTART.md) | [Guide Complet](docs/COLLABORATIVE_CHAT.md)

---

## 🔒 Sécurité & Protection

### Système de Sécurité à 4 Niveaux

Choisissez le bon équilibre entre sécurité et autonomie :

| Niveau | Idéal pour | Bloque | Autorise |
|-------|----------|--------|--------|
| **strict** | Production | sudo, chmod, docker, installation de paquets | Lecture, build, test, git sécurisé |
| **moderate** | Développement (par défaut) | Opérations catastrophiques uniquement | La plupart des opérations de dev |
| **permissive** | DevOps/Admin | Opérations catastrophiques uniquement | Presque tout |
| **off** | Test ⚠️ | Rien | Tout (DANGEREUX !) |

**Configuration :**

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

📖 **Documentation Complète** : [Guide des Niveaux de Sécurité](docs/SAFETY_LEVELS.md) | [Démarrage Rapide](docs/SAFETY_QUICKSTART.md)

---

## 💬 Intégration d'Applications de Chat

Connectez-vous à Telegram, Discord, WhatsApp, QQ, DingTalk, LINE, WeCom et plus.

**Configuration Rapide (Telegram) :**

1. Créer un bot avec [@BotFather](https://t.me/BotFather)
2. Configurer :
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
3. Exécuter : `picoclaw gateway`

📖 **Plus de Canaux** : Voir les [sections README](#-chat-apps) pour Discord, WhatsApp, QQ, etc.

---

## 🐳 Docker Compose

```bash
# Cloner le repo
git clone https://github.com/sipeed/picoclaw.git
cd picoclaw

# Première exécution (génère la config)
docker compose -f docker/docker-compose.yml --profile gateway up

# Éditer la config
vim docker/data/config.json

# Démarrer
docker compose -f docker/docker-compose.yml --profile gateway up -d
```

---

## ⚙️ Configuration

### Fournisseurs Supportés

| Fournisseur | Objectif | Obtenir une Clé API |
|----------|---------|-------------|
| OpenAI | Modèles GPT | [platform.openai.com](https://platform.openai.com) |
| Anthropic | Modèles Claude | [console.anthropic.com](https://console.anthropic.com) |
| Zhipu | Modèles GLM (Chinois) | [bigmodel.cn](https://bigmodel.cn) |
| OpenRouter | Tous les modèles | [openrouter.ai](https://openrouter.ai) |
| Gemini | Modèles Google | [aistudio.google.com](https://aistudio.google.com) |
| Groq | Inférence rapide | [console.groq.com](https://console.groq.com) |
| Ollama | Modèles locaux | Local (pas de clé) |

### Configuration de Modèle

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

## 📱 Déployez Partout

### Anciens Téléphones Android

```bash
# Installer Termux depuis F-Droid
pkg install proot
wget https://github.com/sipeed/picoclaw/releases/download/v0.1.1/picoclaw-linux-arm64
chmod +x picoclaw-linux-arm64
termux-chroot ./picoclaw-linux-arm64 onboard
```

### Matériel à Faible Coût

- 9,9$ [LicheeRV-Nano](https://www.aliexpress.com/item/1005006519668532.html) - Assistant domestique minimal
- 30-50$ [NanoKVM](https://www.aliexpress.com/item/1005007369816019.html) - Maintenance serveur
- 50$ [MaixCAM](https://www.aliexpress.com/item/1005008053333693.html) - Surveillance intelligente

---

## 🛠️ Référence CLI

| Commande | Description |
|---------|-------------|
| `picoclaw onboard` | Initialiser config & workspace |
| `picoclaw agent -m "..."` | Discuter avec l'agent |
| `picoclaw agent` | Mode interactif |
| `picoclaw gateway` | Démarrer la passerelle pour les apps de chat |
| `picoclaw status` | Afficher le statut |
| `picoclaw team create <config>` | Créer une équipe d'agents |
| `picoclaw team list` | Lister les équipes |
| `picoclaw team status <id>` | Statut de l'équipe |
| `picoclaw cron list` | Lister les tâches planifiées |

---

## 🤝 Contribuer

Les PR sont les bienvenues ! Voir [CONTRIBUTING.md](CONTRIBUTING.md) pour les directives.

**Points Clés :**
- Les contributions assistées par IA sont bienvenues (avec divulgation)
- Exécutez `make check` avant de soumettre
- Gardez les PR ciblées et petites
- Remplissez complètement le template de PR

**Feuille de Route** : [ROADMAP.md](ROADMAP.md)

---

## 📊 Comparaison

|  | OpenClaw | NanoBot | **PicoClaw** |
|---|---|---|---|
| **Langage** | TypeScript | Python | **Go** |
| **RAM** | >1GB | >100MB | **<10MB** |
| **Démarrage** (0.8GHz) | >500s | >30s | **<1s** |
| **Coût** | Mac Mini 599$ | Linux SBC ~50$ | **N'importe quel Linux 10$+** |

---

## 📝 Documentation

- [Guide Multi-Agents](docs/MULTI_AGENT_GUIDE.md) - Collaboration d'équipe
- [Niveaux de Sécurité](docs/SAFETY_LEVELS.md) - Configuration de sécurité
- [Contrôle d'Accès aux Outils](docs/TEAM_TOOL_ACCESS.md) - Système de permissions
- [Sélection de Modèle](docs/MULTI_AGENT_MODEL_SELECTION.md) - Modèles par rôle
- [Contribuer](CONTRIBUTING.md) - Directives de contribution
- [Feuille de Route](ROADMAP.md) - Plans futurs
- [Changelog](CHANGELOG_MULTI_AGENT.md) - Mises à jour multi-agents

---

## 🐛 Dépannage

**La recherche web ne fonctionne pas ?**
- Obtenez une clé API gratuite : [Brave Search](https://brave.com/search/api) (2000/mois) ou [Tavily](https://tavily.com) (1000/mois)
- Ou utilisez DuckDuckGo (pas de clé nécessaire, repli automatique)

**Conflit de bot Telegram ?**
- Une seule instance de `picoclaw gateway` peut fonctionner à la fois

**Erreurs de filtrage de contenu ?**
- Certains fournisseurs (Zhipu) ont un filtrage strict - essayez de reformuler ou utilisez un modèle différent

---

## 📢 Communauté

- **Discord** : [Rejoindre le Serveur](https://discord.gg/V4sAZ9XWpN)
- **WeChat** : <img src="assets/wechat.png" width="200">
- **Twitter** : [@SipeedIO](https://x.com/SipeedIO)
- **Site Web** : [picoclaw.io](https://picoclaw.io)

---

## 📄 Licence

Licence MIT - voir [LICENSE](LICENSE) pour les détails

---

**Fait avec ❤️ par la communauté PicoClaw**
