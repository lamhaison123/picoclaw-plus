<div align="center">
  <img src="assets/logo.jpg" alt="PicoClaw" width="512">

  <h1>PicoClaw: Trợ lý AI Siêu Hiệu Năng bằng Go</h1>

  <h3>Phần cứng $10 · 10MB RAM · Khởi động 1s · 皮皮虾， chúng ta đi!</h3>

  <p>
    <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go&logoColor=white" alt="Go">
    <img src="https://img.shields.io/badge/Arch-x86__64%2C%20ARM64%2C%20RISC--V-blue" alt="Hardware">
    <img src="https://img.shields.io/badge/Phiên_bản-v2.1.0-orange" alt="Version">
    <img src="https://img.shields.io/badge/license-MIT-green" alt="License">
    <br>
    <a href="https://picoclaw.io"><img src="https://img.shields.io/badge/Website-picoclaw.io-blue?style=flat&logo=google-chrome&logoColor=white" alt="Website"></a>
    <a href="https://x.com/SipeedIO"><img src="https://img.shields.io/badge/X_(Twitter)-SipeedIO-black?style=flat&logo=x&logoColor=white" alt="Twitter"></a>
    <br>
    <a href="./assets/wechat.png"><img src="https://img.shields.io/badge/WeChat-Group-41d56b?style=flat&logo=wechat&logoColor=white"></a>
    <a href="https://discord.gg/V4sAZ9XWpN"><img src="https://img.shields.io/badge/Discord-Community-4c60eb?style=flat&logo=discord&logoColor=white" alt="Discord"></a>
  </p>

[Tiếng Việt] | [English](README.md) | [中文](README.zh.md) | [日本語](README.ja.md) | [Português](README.pt-br.md) | [Français](README.fr.md)

</div>

---

## 🚀 PicoClaw là gì?

PicoClaw là một trợ lý AI cá nhân siêu nhẹ được xây dựng bằng Go, được thiết kế để chạy trên phần cứng tối thiểu với hiệu suất tối đa.

⚡️ **Chạy trên phần cứng $10 với <50MB RAM** — Hỗ trợ phân cấp đa tác nhân lên đến 20 cấp độ!

🦐 Lấy cảm hứng từ [nanobot](https://github.com/HKUDS/nanobot), được tái cấu trúc hoàn toàn thông qua quá trình tự khởi tạo do AI điều khiển.

---

## ✨ Các tính năng chính

| Tính năng | Mô tả |
|-----------|-------|
| 🪶 **Siêu nhẹ** | <10MB RAM — nhỏ hơn 99% so với các lựa chọn thay thế |
| 🧠 **Bộ nhớ Vector** | Bộ nhớ dài hạn với **Qdrant** & **LanceDB** (Local/Cloud) |
| 👁️ **Hỗ trợ Thị giác** | Hiểu hình ảnh đa phương thức (GPT-4o, Claude 3.5) |
| 🔍 **Tìm kiếm Nâng cao** | 7+ nhà cung cấp: Perplexity, Exa, SearXNG, GLM, Brave, Tavily |
| ⚡️ **Cực nhanh** | Thời gian khởi động 1 giây ngay cả trên lõi đơn 0.6GHz |
| 🌍 **Tính di động cao** | Một file binary duy nhất cho RISC-V, ARM, x86 |
| 🤖 **AI-Bootstrapped** | 95% do AI tạo ra với sự tinh chỉnh của con người |
| 🤝 **Nhóm Đa tác nhân** | Phối hợp các tác nhân AI chuyên biệt với các khả năng dựa trên vai trò |
| 🔒 **An toàn Linh hoạt** | Hệ thống an toàn 4 cấp độ để kiểm soát LLM |

---

## 🚀 Những cải tiến gần đây

### v2.1.0 (Mới nhất) - Bản cập nhật Bộ nhớ & Thị giác

✅ **Hệ sinh thái Bộ nhớ Nâng cao**
- **Tìm kiếm Vector**: Tích hợp **Qdrant** (production) và **LanceDB** (nhúng cục bộ).
- **Khả năng phục hồi**: Mô hình Circuit Breaker & Exponential Backoff cho tất cả các nhà cung cấp bộ nhớ.
- **Mem0 & MindGraph**: Hỗ trợ bộ nhớ cá nhân hóa và đồ thị tri thức.

✅ **Thị giác Đa phương thức (Vision)**
- **Hiểu hình ảnh**: Hỗ trợ đầy đủ các mô hình Vision của OpenAI và Anthropic.
- **Streaming hiệu quả**: Mã hóa Base64 với tự động phát hiện loại MIME.

✅ **Tìm kiếm & Điều hướng Nâng cao**
- **7 Nhà cung cấp Tìm kiếm**: Thêm SearXNG, GLM Search và Exa AI.
- **Điều hướng thông minh**: Lựa chọn mô hình dựa trên độ phức tạp để tiết kiệm chi phí.
- **Trích xuất suy nghĩ**: Hỗ trợ các khối "Extended Thinking" của Claude 3.5.

✅ **Khả năng phục hồi doanh nghiệp**
- **PICOCLAW_HOME**: Cấu hình thư mục chủ cho Docker/Đa người dùng.
- **Lưu trữ JSONL**: Theo dõi phiên làm việc an toàn, chống crash.

📖 **Chi tiết**: [Release Notes v2.1](Release_note_2.1.md) | [Changelog](CHANGELOG.md)

---

## 📦 Bắt đầu nhanh

### 1. Cài đặt

**Tải bản binary đã biên dịch sẵn:**
```bash
# Tải từ GitHub Releases
wget https://github.com/lamhaison123/picoclaw-plus/releases/download/v2.1/picoclaw-linux-amd64
chmod +x picoclaw-linux-amd64
```

### 2. Khởi tạo

```bash
./picoclaw-linux-amd64 onboard
```

### 3. Cấu hình

Chỉnh sửa `~/.picoclaw/config.json` để thêm API key cho OpenAI, Anthropic hoặc Qdrant.

---

## 💬 Tích hợp ứng dụng Chat

Kết nối với Telegram, Discord, WhatsApp, WeCom và nhiều ứng dụng khác.

```bash
picoclaw gateway
```

---

## 📄 Giấy phép

Giấy phép MIT - xem [LICENSE](LICENSE) để biết chi tiết

**Được tạo ra với ❤️ bởi cộng đồng PicoClaw**
