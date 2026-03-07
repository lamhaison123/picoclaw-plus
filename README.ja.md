<div align="center">
  <img src="assets/logo.jpg" alt="PicoClaw" width="512">

  <h1>PicoClaw: Go で書かれた超効率 AI アシスタント</h1>

  <h3>$10 ハードウェア · 10MB RAM · 1秒起動 · 行くぜ、シャコ！</h3>

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

[中文](README.zh.md) | **日本語** | [Português](README.pt-br.md) | [Tiếng Việt](README.vi.md) | [Français](README.fr.md) | [English](README.md)

</div>

---

## 🚀 PicoClawとは？

PicoClawは、Goで構築された超軽量パーソナルAIアシスタントで、最小限のハードウェアで最大の効率で動作するように設計されています。

⚡️ **$10のハードウェアで10MB未満のRAMで動作** — OpenClawより99%少ないメモリ、Mac miniより98%安い！

🦐 [nanobot](https://github.com/HKUDS/nanobot)にインスパイアされ、AI駆動のセルフブートストラッピングプロセスを通じてゼロから再構築されました。

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
> **🚨 セキュリティ & 公式チャンネル**
>
> * **暗号通貨なし:** PicoClawには公式トークン/コインは**ありません**。すべての主張は**詐欺**です。
> * **公式ドメイン:** **[picoclaw.io](https://picoclaw.io)** と **[sipeed.com](https://sipeed.com)** のみ
> * **警告:** 初期開発段階 - v1.0前に本番環境へのデプロイは推奨されません
> * **注意:** 最近のPRにより、メモリ使用量が一時的に10-20MBに増加する可能性があります

---

## ✨ 主な機能

| 機能 | 説明 |
|---------|-------------|
| 🪶 **超軽量** | <10MB RAM — 代替品より99%小さい |
| 💰 **最小コスト** | $10のハードウェアで動作 — 98%安い |
| ⚡️ **超高速** | 0.6GHzシングルコアでも1秒で起動 |
| 🌍 **真のポータビリティ** | RISC-V、ARM、x86用の単一バイナリ |
| 🤖 **AIブートストラップ** | 95%のコードがエージェント生成、人間による洗練 |
| 🤝 **マルチエージェントチーム** | 専門化されたAIエージェントを調整 |
| 💬 **コラボレーティブチャット** | TelegramでのIRCスタイルのマルチエージェント会話 |
| 🔓 **柔軟なセキュリティ** | LLM制御のための4レベルセキュリティシステム |

---

## 📦 クイックスタート

### 1. インストール

**プリコンパイル済みバイナリをダウンロード:**
```bash
# https://github.com/sipeed/picoclaw/releases からダウンロード
wget https://github.com/sipeed/picoclaw/releases/download/v0.1.1/picoclaw-linux-amd64
chmod +x picoclaw-linux-amd64
```

**またはソースからビルド:**
```bash
git clone https://github.com/sipeed/picoclaw.git
cd picoclaw
make build
```

### 2. 初期化

```bash
picoclaw onboard
```

### 3. 設定

`~/.picoclaw/config.json`を編集:

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

**APIキーを取得:**
- LLM: [OpenRouter](https://openrouter.ai/keys) · [Zhipu](https://open.bigmodel.cn) · [Anthropic](https://console.anthropic.com)
- 検索（オプション）: [Tavily](https://tavily.com) · [Brave](https://brave.com/search/api)

### 4. チャット

```bash
picoclaw agent -m "2+2は？"
```

---

## 🤝 マルチエージェントコラボレーション

役割ベースの機能を持つ専門化されたAIエージェントのチームを調整:

**3つのコラボレーションパターン:**
- 🔄 **シーケンシャル**: タスクが順番に実行される（設計 → 実装 → テスト → レビュー）
- ⚡ **パラレル**: タスクが同時に実行されて速度向上
- 🌳 **階層的**: 複雑なタスクが動的に分解される

**クイック例:**

```bash
# 開発チームを作成
picoclaw team create templates/teams/development-team.json

# タスクを実行
picoclaw team execute dev-team-001 -t "hello world関数を作成"

# ステータスを確認
picoclaw team status dev-team-001
```

**主な機能:**
- 👥 ツール権限を持つ役割ベースの専門化
- 🗳️ コンセンサス投票（多数決/全会一致/重み付け）
- 🔄 動的エージェント構成
- 📊 包括的な監視
- 💾 自動メモリ永続化

📖 **詳細**: [マルチエージェントガイド](docs/MULTI_AGENT_GUIDE.md) | [例](examples/teams/)

---

## 💬 コラボレーティブチャット（新機能！）

**TelegramでのIRCスタイルのマルチエージェント会話** — 1つのメッセージで複数のエージェントに言及すると、すべてが完全なコンテキストで応答します！

```
User: @architect @developer ユーザー認証をどのように実装すべきですか？

[abc123] 🏗️ ARCHITECT: JWTトークンを使用することをお勧めします...
[abc123] 💻 DEVELOPER: 次のように実装できます...
```

**クイックセットアップ:**

1. 設定で有効化:
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

2. チーム設定を作成（[templates/teams/collaborative-dev-team.json](templates/teams/collaborative-dev-team.json)を参照）

3. ゲートウェイを起動: `picoclaw gateway`

**機能:**
- 🎯 @メンションベースのルーティング（@architect、@developer、@tester）
- ⚡ 並列エージェント実行
- 🧠 共有会話コンテキスト
- 🎨 絵文字付きIRCスタイルフォーマット
- 📝 チャットごとのセッション管理
- 👥 `/who`コマンド - すべての登録エージェントとアクティブセッションを表示

**コマンド:**
- `/who` - チームステータス、登録エージェント、アクティブエージェントを表示
- `/help` - 利用可能なコマンドを表示

📖 **詳細**: [クイックスタート](docs/COLLABORATIVE_CHAT_QUICKSTART.md) | [完全ガイド](docs/COLLABORATIVE_CHAT.md)

---

## 🔒 セキュリティと安全性

### 4レベルセキュリティシステム

セキュリティと自律性の適切なバランスを選択:

| レベル | 最適な用途 | ブロック | 許可 |
|-------|----------|--------|--------|
| **strict** | 本番環境 | sudo、chmod、docker、パッケージインストール | 読み取り、ビルド、テスト、安全なgit |
| **moderate** | 開発（デフォルト） | 壊滅的な操作のみ | ほとんどの開発操作 |
| **permissive** | DevOps/管理者 | 壊滅的な操作のみ | ほぼすべて |
| **off** | テスト ⚠️ | なし | すべて（危険！） |

**設定:**

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

📖 **完全なドキュメント**: [セキュリティレベルガイド](docs/SAFETY_LEVELS.md) | [クイックスタート](docs/SAFETY_QUICKSTART.md)

---

## 💬 チャットアプリ統合

Telegram、Discord、WhatsApp、QQ、DingTalk、LINE、WeComなどに接続。

**クイックセットアップ（Telegram）:**

1. [@BotFather](https://t.me/BotFather)でボットを作成
2. 設定:
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
3. 実行: `picoclaw gateway`

📖 **その他のチャンネル**: Discord、WhatsApp、QQなどについては[READMEセクション](#-chat-apps)を参照

---

## 🐳 Docker Compose

```bash
# リポジトリをクローン
git clone https://github.com/sipeed/picoclaw.git
cd picoclaw

# 初回実行（設定を生成）
docker compose -f docker/docker-compose.yml --profile gateway up

# 設定を編集
vim docker/data/config.json

# 起動
docker compose -f docker/docker-compose.yml --profile gateway up -d
```

---

## ⚙️ 設定

### サポートされているプロバイダー

| プロバイダー | 目的 | APIキーを取得 |
|----------|---------|-------------|
| OpenAI | GPTモデル | [platform.openai.com](https://platform.openai.com) |
| Anthropic | Claudeモデル | [console.anthropic.com](https://console.anthropic.com) |
| Zhipu | GLMモデル（中国語） | [bigmodel.cn](https://bigmodel.cn) |
| OpenRouter | すべてのモデル | [openrouter.ai](https://openrouter.ai) |
| Gemini | Googleモデル | [aistudio.google.com](https://aistudio.google.com) |
| Groq | 高速推論 | [console.groq.com](https://console.groq.com) |
| Ollama | ローカルモデル | ローカル（キー不要） |

### モデル設定

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

## 📱 どこでもデプロイ

### 古いAndroidフォン

```bash
# F-DroidからTermuxをインストール
pkg install proot
wget https://github.com/sipeed/picoclaw/releases/download/v0.1.1/picoclaw-linux-arm64
chmod +x picoclaw-linux-arm64
termux-chroot ./picoclaw-linux-arm64 onboard
```

### 低コストハードウェア

- $9.9 [LicheeRV-Nano](https://www.aliexpress.com/item/1005006519668532.html) - 最小限のホームアシスタント
- $30-50 [NanoKVM](https://www.aliexpress.com/item/1005007369816019.html) - サーバーメンテナンス
- $50 [MaixCAM](https://www.aliexpress.com/item/1005008053333693.html) - スマート監視

---

## 🛠️ CLIリファレンス

| コマンド | 説明 |
|---------|-------------|
| `picoclaw onboard` | 設定とワークスペースを初期化 |
| `picoclaw agent -m "..."` | エージェントとチャット |
| `picoclaw agent` | インタラクティブモード |
| `picoclaw gateway` | チャットアプリ用ゲートウェイを起動 |
| `picoclaw status` | ステータスを表示 |
| `picoclaw team create <config>` | エージェントチームを作成 |
| `picoclaw team list` | チームをリスト |
| `picoclaw team status <id>` | チームステータス |
| `picoclaw cron list` | スケジュールされたジョブをリスト |

---

## 🤝 貢献

PRを歓迎します！ガイドラインについては[CONTRIBUTING.md](CONTRIBUTING.md)を参照してください。

**重要なポイント:**
- AI支援による貢献を歓迎（開示が必要）
- 提出前に`make check`を実行
- PRを集中的で小さく保つ
- PRテンプレートを完全に記入

**ロードマップ**: [ROADMAP.md](ROADMAP.md)

---

## 📊 比較

|  | OpenClaw | NanoBot | **PicoClaw** |
|---|---|---|---|
| **言語** | TypeScript | Python | **Go** |
| **RAM** | >1GB | >100MB | **<10MB** |
| **起動** (0.8GHz) | >500s | >30s | **<1s** |
| **コスト** | Mac Mini $599 | Linux SBC ~$50 | **任意のLinux $10+** |

---

## 📝 ドキュメント

- [マルチエージェントガイド](docs/MULTI_AGENT_GUIDE.md) - チームコラボレーション
- [セキュリティレベル](docs/SAFETY_LEVELS.md) - セキュリティ設定
- [ツールアクセス制御](docs/TEAM_TOOL_ACCESS.md) - 権限システム
- [モデル選択](docs/MULTI_AGENT_MODEL_SELECTION.md) - 役割ごとのモデル
- [貢献](CONTRIBUTING.md) - 貢献ガイドライン
- [ロードマップ](ROADMAP.md) - 将来の計画
- [変更履歴](CHANGELOG_MULTI_AGENT.md) - マルチエージェント更新

---

## 🐛 トラブルシューティング

**Web検索が機能しない？**
- 無料APIキーを取得: [Brave Search](https://brave.com/search/api)（2000/月）または[Tavily](https://tavily.com)（1000/月）
- またはDuckDuckGoを使用（キー不要、自動フォールバック）

**Telegramボットの競合？**
- 一度に1つの`picoclaw gateway`インスタンスのみ実行可能

**コンテンツフィルタリングエラー？**
- 一部のプロバイダー（Zhipuなど）には厳格なフィルタリングがあります - 言い換えるか別のモデルを使用

---

## 📢 コミュニティ

- **Discord**: [サーバーに参加](https://discord.gg/V4sAZ9XWpN)
- **WeChat**: <img src="assets/wechat.png" width="200">
- **Twitter**: [@SipeedIO](https://x.com/SipeedIO)
- **ウェブサイト**: [picoclaw.io](https://picoclaw.io)

---

## 📄 ライセンス

MITライセンス - 詳細は[LICENSE](LICENSE)を参照

---

**PicoClawコミュニティによって❤️で作成**
