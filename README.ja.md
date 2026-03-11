<div align="center">
  <img src="assets/logo.jpg" alt="PicoClaw" width="512">

  <h1>PicoClaw: Go で書かれた超効率 AI アシスタント</h1>

  <h3>$10 ハードウェア · 10MB RAM · 1秒起動 · 皮皮虾，我们走！</h3>

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

[日本語] | [English](README.md) | [中文](README.zh.md) | [Português](README.pt-br.md) | [Tiếng Việt](README.vi.md) | [Français](README.fr.md)

</div>

---

## 🚀 PicoClaw とは？

PicoClaw は、Go で構築された超軽量パーソナル AI アシスタントで、最小限のハードウェアで最大の効率で動作するように設計されています。

⚡️ **$10 のハードウェアで 50MB 未満の RAM で動作** — 最大 20 レベルのマルチエージェント・カスケードをサポート！

🦐 [nanobot](https://github.com/HKUDS/nanobot) にインスパイアされ、AI 駆動のセルフブートストラッピングプロセスを通じてゼロから再構築されました。

---

## ✨ 主な機能

| 機能 | 説明 |
|---------|-------------|
| 🪶 **超軽量** | <10MB RAM — 代替品より 99% 小さい |
| 🧠 **ベクトルメモリ** | **Qdrant** & **LanceDB** (ローカル/クラウド) による長期記憶 |
| 👁️ **ビジョンサポート** | 多モード画像理解 (GPT-4o, Claude 3.5) |
| 🔍 **高度な検索** | 7 つ以上のプロバイダー：Perplexity, Exa, SearXNG, GLM, Brave, Tavily |
| ⚡️ **超高速** | 0.6GHz シングルコアでも 1 秒で起動 |
| 🌍 **真のポータビリティ** | RISC-V、ARM、x86 用の単一バイナリ |
| 🤖 **AI ブートストラップ** | 95% のコードがエージェント生成、人間による洗練 |
| 🤝 **マルチエージェントチーム** | 役割ベースの機能を持つ専門化された AI エージェントを調整 |
| 🔒 **柔軟なセキュリティ** | LLM 制御のための 4 レベルセキュリティシステム |

---

## 🚀 最近の改善

### v2.1.0 (最新) - メモリとビジョンのアップデート

✅ **高度なメモリ・エコシステム**
- **ベクトル検索**: **Qdrant** (本番環境用) と **LanceDB** (ローカル埋め込み用) を統合。
- **レジリエンス**: すべてのメモリプロバイダーにサーキットブレーカーパターンと指数バックオフを導入。
- **Mem0 & MindGraph**: パーソナライズされた記憶と知識グラフをサポート。

✅ **マルチモーダル・ビジョン**
- **画像理解**: OpenAI と Anthropic のビジョンモデルをフルサポート。
- **効率的なストリーミング**: MIME 自動検出機能付き Base64 エンコーディング。

✅ **検索とルーティングの強化**
- **7 つの検索プロバイダー**: SearXNG、GLM Search、Exa AI を追加。
- **スマートルーティング**: コスト削減のため、複雑さに基づいたモデル選択。
- **推論抽出**: Claude 3.5 の「Extended Thinking」ブロックをサポート。

✅ **エンタープライズ・レジリエンス**
- **PICOCLAW_HOME**: Docker/マルチテナント環境向けに設定可能なホームディレクトリ。
- **JSONL ストレージ**: クラッシュに強い、アペンドオンリーのセッション追跡。

📖 **詳細**: [リリースノート v2.1](Release_note_2.1.md) | [変更履歴](CHANGELOG.md)

---

## 📦 クイックスタート

### 1. インストール

**プリコンパイル済みバイナリをダウンロード:**
```bash
# GitHub Releases からダウンロード
wget https://github.com/lamhaison123/picoclaw-plus/releases/download/v2.1/picoclaw-linux-amd64
chmod +x picoclaw-linux-amd64
```

### 2. 初期化

```bash
./picoclaw-linux-amd64 onboard
```

### 3. 設定

`~/.picoclaw/config.json` を編集して、OpenAI、Anthropic、または Qdrant の API キーを追加します。

---

## 💬 チャットアプリ統合

Telegram、Discord、WhatsApp、WeCom などに接続します。

```bash
picoclaw gateway
```

---

## 📄 ライセンス

MIT ライセンス - 詳細は [LICENSE](LICENSE) を参照

**PicoClaw コミュニティによって ❤️ で作成**
