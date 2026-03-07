# PicoClaw Systemd Service

Quản lý PicoClaw như một system service trên Linux với systemd.

---

## 🚀 Quick Start

### 1. Cài đặt service

```bash
# Chạy script cài đặt tự động
sudo bash install-systemd.sh
```

### 2. Khởi động service

```bash
sudo systemctl start picoclaw
```

### 3. Kiểm tra trạng thái

```bash
sudo systemctl status picoclaw
```

### 4. Xem logs

```bash
sudo journalctl -u picoclaw -f
```

---

## 📁 Files

- **`picoclaw.service`** - Systemd service definition file
- **`install-systemd.sh`** - Automatic installation script
- **`picoclaw-ctl.sh`** - Helper control script
- **`SYSTEMD_SETUP.md`** - Detailed documentation

---

## 🎮 Control Script

Sử dụng `picoclaw-ctl.sh` để quản lý service dễ dàng hơn:

```bash
# Make executable
chmod +x picoclaw-ctl.sh

# Show status
sudo ./picoclaw-ctl.sh status

# Start service
sudo ./picoclaw-ctl.sh start

# Stop service
sudo ./picoclaw-ctl.sh stop

# Restart service
sudo ./picoclaw-ctl.sh restart

# Enable auto-start
sudo ./picoclaw-ctl.sh enable

# Disable auto-start
sudo ./picoclaw-ctl.sh disable

# Show logs
sudo ./picoclaw-ctl.sh logs 100

# Follow logs
sudo ./picoclaw-ctl.sh follow
```

---

## 📋 Common Commands

### Service Management

```bash
# Start
sudo systemctl start picoclaw

# Stop
sudo systemctl stop picoclaw

# Restart
sudo systemctl restart picoclaw

# Status
sudo systemctl status picoclaw

# Enable auto-start
sudo systemctl enable picoclaw

# Disable auto-start
sudo systemctl disable picoclaw
```

### Logs

```bash
# Follow logs
sudo journalctl -u picoclaw -f

# Last 100 lines
sudo journalctl -u picoclaw -n 100

# Since 1 hour ago
sudo journalctl -u picoclaw --since "1 hour ago"

# Only errors
sudo journalctl -u picoclaw -p err
```

---

## ⚙️ Configuration

### Service File Location
```
/etc/systemd/system/picoclaw.service
```

### Edit Service
```bash
sudo systemctl edit --full picoclaw.service
```

### After Editing
```bash
sudo systemctl daemon-reload
sudo systemctl restart picoclaw
```

---

## 🔧 Customization

### Change User

Edit service file:
```ini
[Service]
User=your-username
Group=your-group
WorkingDirectory=/home/your-username
Environment="HOME=/home/your-username"
ReadWritePaths=/home/your-username/.picoclaw
```

### Add Environment Variables

```ini
[Service]
Environment="PICOCLAW_LOG_LEVEL=debug"
Environment="PICOCLAW_PORT=8080"
```

### Resource Limits

```ini
[Service]
LimitNOFILE=100000
LimitNPROC=8192
CPUQuota=200%
MemoryLimit=2G
```

---

## 🛡️ Security Features

Service được cấu hình với các security options:

- ✅ `NoNewPrivileges=true` - Không escalate privileges
- ✅ `PrivateTmp=true` - Isolated /tmp
- ✅ `ProtectSystem=strict` - Read-only system directories
- ✅ `ProtectHome=read-only` - Read-only home directories
- ✅ `ReadWritePaths=/root/.picoclaw` - Chỉ ghi vào workspace

---

## 🔍 Troubleshooting

### Service không start

```bash
# Check status
sudo systemctl status picoclaw

# View logs
sudo journalctl -u picoclaw -n 50

# Verify service file
sudo systemd-analyze verify picoclaw.service
```

### Port conflict

```bash
# Check port 18790
sudo netstat -tulpn | grep 18790
sudo lsof -i :18790

# Kill conflicting process
sudo kill -9 <PID>
```

### Permission issues

```bash
# Fix workspace permissions
sudo chown -R root:root /root/.picoclaw
sudo chmod -R 755 /root/.picoclaw
```

---

## 📊 Monitoring

### Check Uptime
```bash
systemctl show picoclaw --property=ActiveEnterTimestamp
```

### Check Memory Usage
```bash
systemctl status picoclaw | grep Memory
```

### Check CPU Usage
```bash
systemctl status picoclaw | grep CPU
```

### All Metrics
```bash
systemctl show picoclaw
```

---

## 🔄 Updates

### Update Binary

```bash
# 1. Stop service
sudo systemctl stop picoclaw

# 2. Update binary
sudo cp /path/to/new/picoclaw /usr/local/bin/picoclaw
sudo chmod +x /usr/local/bin/picoclaw

# 3. Start service
sudo systemctl start picoclaw

# 4. Verify
sudo systemctl status picoclaw
```

---

## 📦 Uninstall

```bash
# Stop and disable
sudo systemctl stop picoclaw
sudo systemctl disable picoclaw

# Remove service file
sudo rm /etc/systemd/system/picoclaw.service

# Reload systemd
sudo systemctl daemon-reload

# (Optional) Remove binary
sudo rm /usr/local/bin/picoclaw

# (Optional) Remove workspace
sudo rm -rf /root/.picoclaw
```

---

## 📚 Documentation

Xem `SYSTEMD_SETUP.md` để biết thêm chi tiết về:
- Advanced configuration
- Security hardening
- Resource management
- Monitoring and alerting
- Best practices

---

## 🆘 Support

- GitHub: https://github.com/sipeed/picoclaw
- Issues: https://github.com/sipeed/picoclaw/issues
- Documentation: https://github.com/sipeed/picoclaw/docs

---

**Version**: 1.0  
**Last Updated**: 2026-03-07  
**Tested On**: Ubuntu 20.04+, Debian 11+, CentOS 8+
