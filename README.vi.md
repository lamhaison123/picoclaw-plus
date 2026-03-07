<div align="center">
  <img src="assets/logo.jpg" alt="PicoClaw" width="512">

  <h1>PicoClaw: Trợ lý AI Siêu Nhẹ viết bằng Go</h1>

  <h3>Phần cứng $10 · RAM 10MB · Khởi động 1 giây · Nào, xuất phát!</h3>

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

[中文](README.zh.md) | [日本語](README.ja.md) | [Português](README.pt-br.md) | **Tiếng Việt** | [Français](README.fr.md) | [English](README.md)

</div>

---

## 🚀 PicoClaw là gì?

PicoClaw là trợ lý AI cá nhân siêu nhẹ được xây dựng bằng Go, được thiết kế để chạy trên phần cứng tối thiểu với hiệu suất tối đa.

⚡️ **Chạy trên phần cứng $10 với RAM <10MB** — tiết kiệm 99% bộ nhớ so với OpenClaw và rẻ hơn 98% so với Mac mini!

🦐 Lấy cảm hứng từ [nanobot](https://github.com/HKUDS/nanobot), được tái cấu trúc hoàn toàn thông qua quá trình tự khởi tạo (self-bootstrapping) do AI điều khiển.

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
> **🚨 BẢO MẬT & KÊNH CHÍNH THỨC**
>
> * **KHÔNG CÓ CRYPTO:** PicoClaw **KHÔNG** có bất kỳ token/coin chính thức nào. Mọi thông tin đều là **LỪA ĐẢO**.
> * **DOMAIN CHÍNH THỨC:** Chỉ có **[picoclaw.io](https://picoclaw.io)** và **[sipeed.com](https://sipeed.com)**
> * **Cảnh báo:** Đang trong giai đoạn phát triển sớm - không nên triển khai production trước v1.0
> * **Lưu ý:** Các PR gần đây có thể tăng bộ nhớ lên 10-20MB tạm thời

---

## ✨ Tính năng nổi bật

| Tính năng | Mô tả |
|---------|-------------|
| 🪶 **Siêu nhẹ** | <10MB RAM — nhỏ hơn 99% so với các giải pháp khác |
| 💰 **Chi phí tối thiểu** | Chạy trên phần cứng $10 — rẻ hơn 98% |
| ⚡️ **Cực nhanh** | Khởi động 1 giây ngay cả trên CPU 0.6GHz đơn nhân |
| 🌍 **Di động thực sự** | Binary duy nhất cho RISC-V, ARM, x86 |
| 🤖 **Tự xây dựng bằng AI** | 95% code được agent tạo ra với tinh chỉnh của con người |
| 🤝 **Đội nhóm Multi-Agent** | Phối hợp các AI agent chuyên môn hóa |
| 💬 **Collaborative Chat** | Trò chuyện multi-agent kiểu IRC trong Telegram |
| 🔓 **Hệ thống An toàn Linh hoạt** | 4 cấp độ an toàn để kiểm soát LLM |

---

## 📦 Bắt đầu nhanh

### 1. Cài đặt

**Tải binary đã biên dịch:**
```bash
# Tải từ https://github.com/sipeed/picoclaw/releases
wget https://github.com/sipeed/picoclaw/releases/download/v0.1.1/picoclaw-linux-amd64
chmod +x picoclaw-linux-amd64
```

**Hoặc build từ source:**
```bash
git clone https://github.com/sipeed/picoclaw.git
cd picoclaw
make build
```

### 2. Khởi tạo

```bash
picoclaw onboard
```

### 3. Cấu hình

Chỉnh sửa `~/.picoclaw/config.json`:

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

**Lấy API Keys:**
- LLM: [OpenRouter](https://openrouter.ai/keys) · [Zhipu](https://open.bigmodel.cn) · [Anthropic](https://console.anthropic.com)
- Tìm kiếm (tùy chọn): [Tavily](https://tavily.com) · [Brave](https://brave.com/search/api)

### 4. Trò chuyện

```bash
picoclaw agent -m "2+2 bằng mấy?"
```

---

## 🤝 Cộng tác Multi-Agent

Phối hợp các đội nhóm AI agent chuyên môn với khả năng dựa trên vai trò:

**Ba mô hình cộng tác:**
- 🔄 **Sequential**: Các tác vụ thực hiện theo thứ tự (thiết kế → triển khai → test → review)
- ⚡ **Parallel**: Các tác vụ chạy đồng thời để tăng tốc độ
- 🌳 **Hierarchical**: Các tác vụ phức tạp được phân rã động

**Ví dụ nhanh:**

```bash
# Tạo đội phát triển
picoclaw team create templates/teams/development-team.json

# Thực hiện tác vụ
picoclaw team execute dev-team-001 -t "Tạo hàm hello world"

# Kiểm tra trạng thái
picoclaw team status dev-team-001
```

**Tính năng chính:**
- 👥 Chuyên môn hóa theo vai trò với quyền công cụ
- 🗳️ Bỏ phiếu đồng thuận (đa số/nhất trí/có trọng số)
- 🔄 Thành phần agent động
- 📊 Giám sát toàn diện
- 💾 Lưu trữ bộ nhớ tự động

📖 **Tìm hiểu thêm**: [Hướng dẫn Multi-Agent](docs/MULTI_AGENT_GUIDE.md) | [Ví dụ](examples/teams/)

---

## 💬 Collaborative Chat (MỚI!)

**Trò chuyện multi-agent kiểu IRC trong Telegram** — đề cập nhiều agent trong một tin nhắn và tất cả sẽ phản hồi với ngữ cảnh đầy đủ!

```
User: @architect @developer Chúng ta nên triển khai xác thực người dùng như thế nào?

[abc123] 🏗️ ARCHITECT: Tôi khuyên dùng JWT tokens với...
[abc123] 💻 DEVELOPER: Tôi có thể triển khai bằng cách...
```

**Thiết lập nhanh:**

1. Bật trong config:
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

2. Tạo config đội (xem [templates/teams/collaborative-dev-team.json](templates/teams/collaborative-dev-team.json))

3. Khởi động gateway: `picoclaw gateway`

**Tính năng:**
- 🎯 Định tuyến dựa trên @mention (@architect, @developer, @tester)
- ⚡ Thực thi agent song song
- 🧠 Ngữ cảnh hội thoại được chia sẻ
- 🎨 Định dạng kiểu IRC với emoji
- 📝 Quản lý phiên theo chat
- 👥 Lệnh `/who` - Xem tất cả agent đã đăng ký và phiên hoạt động

**Lệnh:**
- `/who` - Hiển thị trạng thái đội, agent đã đăng ký và agent đang hoạt động
- `/help` - Hiển thị các lệnh có sẵn

📖 **Tìm hiểu thêm**: [Bắt đầu nhanh](docs/COLLABORATIVE_CHAT_QUICKSTART.md) | [Hướng dẫn đầy đủ](docs/COLLABORATIVE_CHAT.md)

---

## 🔒 Bảo mật & An toàn

### Hệ thống An toàn 4 Cấp độ

Chọn sự cân bằng phù hợp giữa bảo mật và tự chủ:

| Cấp độ | Phù hợp cho | Chặn | Cho phép |
|-------|----------|--------|--------|
| **strict** | Production | sudo, chmod, docker, cài đặt package | Đọc, build, test, git an toàn |
| **moderate** | Development (mặc định) | Chỉ các thao tác thảm họa | Hầu hết thao tác dev |
| **permissive** | DevOps/Admin | Chỉ các thao tác thảm họa | Gần như mọi thứ |
| **off** | Testing ⚠️ | Không có gì | Mọi thứ (NGUY HIỂM!) |

**Cấu hình:**

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

📖 **Tài liệu đầy đủ**: [Hướng dẫn Cấp độ An toàn](docs/SAFETY_LEVELS.md) | [Bắt đầu nhanh](docs/SAFETY_QUICKSTART.md)

---

## 💬 Tích hợp Ứng dụng Chat

Kết nối với Telegram, Discord, WhatsApp, QQ, DingTalk, LINE, WeCom và nhiều hơn nữa.

**Thiết lập nhanh (Telegram):**

1. Tạo bot với [@BotFather](https://t.me/BotFather)
2. Cấu hình:
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
3. Chạy: `picoclaw gateway`

📖 **Các kênh khác**: Xem [phần README](#-chat-apps) cho Discord, WhatsApp, QQ, v.v.

---

## 🐳 Docker Compose

```bash
# Clone repo
git clone https://github.com/sipeed/picoclaw.git
cd picoclaw

# Chạy lần đầu (tạo config)
docker compose -f docker/docker-compose.yml --profile gateway up

# Chỉnh sửa config
vim docker/data/config.json

# Khởi động
docker compose -f docker/docker-compose.yml --profile gateway up -d
```

---

## ⚙️ Cấu hình

### Các Provider được hỗ trợ

| Provider | Mục đích | Lấy API Key |
|----------|---------|-------------|
| OpenAI | Mô hình GPT | [platform.openai.com](https://platform.openai.com) |
| Anthropic | Mô hình Claude | [console.anthropic.com](https://console.anthropic.com) |
| Zhipu | Mô hình GLM (Trung Quốc) | [bigmodel.cn](https://bigmodel.cn) |
| OpenRouter | Tất cả mô hình | [openrouter.ai](https://openrouter.ai) |
| Gemini | Mô hình Google | [aistudio.google.com](https://aistudio.google.com) |
| Groq | Suy luận nhanh | [console.groq.com](https://console.groq.com) |
| Ollama | Mô hình local | Local (không cần key) |

### Cấu hình Mô hình

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

## 📱 Triển khai mọi nơi

### Điện thoại Android cũ

```bash
# Cài Termux từ F-Droid
pkg install proot
wget https://github.com/sipeed/picoclaw/releases/download/v0.1.1/picoclaw-linux-arm64
chmod +x picoclaw-linux-arm64
termux-chroot ./picoclaw-linux-arm64 onboard
```

### Phần cứng giá rẻ

- $9.9 [LicheeRV-Nano](https://www.aliexpress.com/item/1005006519668532.html) - Trợ lý gia đình tối thiểu
- $30-50 [NanoKVM](https://www.aliexpress.com/item/1005007369816019.html) - Bảo trì server
- $50 [MaixCAM](https://www.aliexpress.com/item/1005008053333693.html) - Giám sát thông minh

---

## 🛠️ Tham chiếu CLI

| Lệnh | Mô tả |
|---------|-------------|
| `picoclaw onboard` | Khởi tạo config & workspace |
| `picoclaw agent -m "..."` | Chat với agent |
| `picoclaw agent` | Chế độ tương tác |
| `picoclaw gateway` | Khởi động gateway cho chat apps |
| `picoclaw status` | Hiển thị trạng thái |
| `picoclaw team create <config>` | Tạo đội agent |
| `picoclaw team list` | Liệt kê các đội |
| `picoclaw team status <id>` | Trạng thái đội |
| `picoclaw cron list` | Liệt kê các công việc đã lên lịch |

---

## 🤝 Đóng góp

Chào đón PR! Xem [CONTRIBUTING.md](CONTRIBUTING.md) để biết hướng dẫn.

**Điểm chính:**
- Chào đón đóng góp có sự hỗ trợ của AI (với công bố)
- Chạy `make check` trước khi submit
- Giữ PR tập trung và nhỏ gọn
- Điền đầy đủ template PR

**Lộ trình**: [ROADMAP.md](ROADMAP.md)

---

## 📊 So sánh

|  | OpenClaw | NanoBot | **PicoClaw** |
|---|---|---|---|
| **Ngôn ngữ** | TypeScript | Python | **Go** |
| **RAM** | >1GB | >100MB | **<10MB** |
| **Khởi động** (0.8GHz) | >500s | >30s | **<1s** |
| **Chi phí** | Mac Mini $599 | Linux SBC ~$50 | **Bất kỳ Linux $10+** |

---

## 📝 Tài liệu

- [Hướng dẫn Multi-Agent](docs/MULTI_AGENT_GUIDE.md) - Cộng tác đội nhóm
- [Cấp độ An toàn](docs/SAFETY_LEVELS.md) - Cấu hình bảo mật
- [Kiểm soát Truy cập Công cụ](docs/TEAM_TOOL_ACCESS.md) - Hệ thống quyền
- [Lựa chọn Mô hình](docs/MULTI_AGENT_MODEL_SELECTION.md) - Mô hình theo vai trò
- [Đóng góp](CONTRIBUTING.md) - Hướng dẫn đóng góp
- [Lộ trình](ROADMAP.md) - Kế hoạch tương lai
- [Changelog](CHANGELOG_MULTI_AGENT.md) - Cập nhật multi-agent

---

## 🐛 Xử lý sự cố

**Tìm kiếm web không hoạt động?**
- Lấy API key miễn phí: [Brave Search](https://brave.com/search/api) (2000/tháng) hoặc [Tavily](https://tavily.com) (1000/tháng)
- Hoặc dùng DuckDuckGo (không cần key, tự động fallback)

**Xung đột Telegram bot?**
- Chỉ một instance `picoclaw gateway` có thể chạy cùng lúc

**Lỗi lọc nội dung?**
- Một số provider (Zhipu) có lọc nghiêm ngặt - thử diễn đạt lại hoặc dùng mô hình khác

---

## 📢 Cộng đồng

- **Discord**: [Tham gia Server](https://discord.gg/V4sAZ9XWpN)
- **WeChat**: <img src="assets/wechat.png" width="200">
- **Twitter**: [@SipeedIO](https://x.com/SipeedIO)
- **Website**: [picoclaw.io](https://picoclaw.io)

---

## 📄 Giấy phép

Giấy phép MIT - xem [LICENSE](LICENSE) để biết chi tiết

---

**Được tạo với ❤️ bởi cộng đồng PicoClaw**
