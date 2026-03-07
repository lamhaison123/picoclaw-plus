# PicoClaw Systemd Service Setup

Hướng dẫn cài đặt và quản lý PicoClaw như một systemd service trên Linux.

---

## 🚀 Quick Install

### Cách 1: Tự động (Khuyến nghị)

```bash
# Chạy script cài đặt
sudo bash install-systemd.sh
```

### Cách 2: Thủ công

```bash
# 1. Copy service file
sudo cp picoclaw.service /etc/systemd/system/

# 2. Reload systemd
sudo systemctl daemon-reload

# 3. Enable service
sudo systemctl enable picoclaw.service

# 4. Start service
sudo systemctl start picoclaw.service
```

---

## 📋 Quản lý Service

### Khởi động service
```bash
sudo systemctl start picoclaw
```

### Dừng service
```bash
sudo systemctl stop picoclaw
```

### Restart service
```bash
sudo systemctl restart picoclaw
```

### Kiểm tra trạng thái
```bash
sudo systemctl status picoclaw
```

### Enable auto-start khi boot
```bash
sudo systemctl enable picoclaw
```

### Disable auto-start
```bash
sudo systemctl disable picoclaw
```

---

## 📊 Xem Logs

### Xem logs realtime (follow mode)
```bash
sudo journalctl -u picoclaw -f
```

### Xem logs từ lần boot gần nhất
```bash
sudo journalctl -u picoclaw -b
```

### Xem logs với timestamp
```bash
sudo journalctl -u picoclaw --since "1 hour ago"
sudo journalctl -u picoclaw --since "2024-03-07 10:00:00"
```

### Xem logs với priority
```bash
# Chỉ xem errors
sudo journalctl -u picoclaw -p err

# Xem warnings và errors
sudo journalctl -u picoclaw -p warning
```

### Export logs ra file
```bash
sudo journalctl -u picoclaw > picoclaw.log
```

---

## ⚙️ Cấu hình Service

### Vị trí file
```
/etc/systemd/system/picoclaw.service
```

### Chỉnh sửa service
```bash
sudo systemctl edit --full picoclaw.service
```

### Sau khi chỉnh sửa
```bash
sudo systemctl daemon-reload
sudo systemctl restart picoclaw
```

---

## 🔧 Tùy chỉnh Service

### Thay đổi user/group
Mặc định service chạy với user `root`. Để chạy với user khác:

```ini
[Service]
User=your-username
Group=your-group
WorkingDirectory=/home/your-username
Environment="HOME=/home/your-username"
ReadWritePaths=/home/your-username/.picoclaw
```

### Thêm environment variables
```ini
[Service]
Environment="PICOCLAW_LOG_LEVEL=debug"
Environment="PICOCLAW_PORT=8080"
```

### Thay đổi restart policy
```ini
[Service]
Restart=on-failure        # Chỉ restart khi fail
RestartSec=5              # Đợi 5s trước khi restart
StartLimitBurst=5         # Tối đa 5 lần restart
StartLimitIntervalSec=60  # Trong 60s
```

---

## 🛡️ Security Hardening

Service đã được cấu hình với các security options:

- `NoNewPrivileges=true` - Không cho phép escalate privileges
- `PrivateTmp=true` - Sử dụng /tmp riêng biệt
- `ProtectSystem=strict` - Chỉ đọc system directories
- `ProtectHome=read-only` - Chỉ đọc home directories
- `ReadWritePaths=/root/.picoclaw` - Chỉ cho phép ghi vào workspace

### Tăng cường security hơn nữa

```ini
[Service]
# Network isolation
PrivateNetwork=false  # Set true nếu không cần network

# Filesystem isolation
ProtectKernelTunables=true
ProtectKernelModules=true
ProtectControlGroups=true

# System calls filtering
SystemCallFilter=@system-service
SystemCallErrorNumber=EPERM
```

---

## 📈 Resource Limits

### Mặc định
- Max open files: 65536
- Max processes: 4096

### Tùy chỉnh
```ini
[Service]
LimitNOFILE=100000      # Tăng số file descriptors
LimitNPROC=8192         # Tăng số processes
LimitMEMLOCK=infinity   # Unlimited memory lock
CPUQuota=200%           # Giới hạn CPU (200% = 2 cores)
MemoryLimit=2G          # Giới hạn RAM
```

---

## 🔍 Troubleshooting

### Service không start
```bash
# Kiểm tra status
sudo systemctl status picoclaw

# Xem logs chi tiết
sudo journalctl -u picoclaw -n 50 --no-pager

# Kiểm tra syntax của service file
sudo systemd-analyze verify picoclaw.service
```

### Service bị crash liên tục
```bash
# Xem crash logs
sudo journalctl -u picoclaw -p err -n 100

# Kiểm tra resource limits
sudo systemctl show picoclaw | grep Limit

# Tăng restart delay
sudo systemctl edit picoclaw.service
# Thêm: RestartSec=30
```

### Permission issues
```bash
# Kiểm tra quyền của workspace
ls -la /root/.picoclaw

# Fix permissions
sudo chown -R root:root /root/.picoclaw
sudo chmod -R 755 /root/.picoclaw
```

### Port đã được sử dụng
```bash
# Kiểm tra port 18790
sudo netstat -tulpn | grep 18790
sudo lsof -i :18790

# Kill process đang dùng port
sudo kill -9 <PID>
```

---

## 🔄 Update PicoClaw

### Cập nhật binary
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

### Rollback nếu có vấn đề
```bash
# Restore old binary
sudo cp /path/to/old/picoclaw /usr/local/bin/picoclaw
sudo systemctl restart picoclaw
```

---

## 📦 Uninstall

```bash
# 1. Stop và disable service
sudo systemctl stop picoclaw
sudo systemctl disable picoclaw

# 2. Remove service file
sudo rm /etc/systemd/system/picoclaw.service

# 3. Reload systemd
sudo systemctl daemon-reload

# 4. (Optional) Remove binary
sudo rm /usr/local/bin/picoclaw

# 5. (Optional) Remove workspace
sudo rm -rf /root/.picoclaw
```

---

## 📊 Monitoring

### Kiểm tra uptime
```bash
systemctl show picoclaw --property=ActiveEnterTimestamp
```

### Kiểm tra memory usage
```bash
systemctl status picoclaw | grep Memory
```

### Kiểm tra CPU usage
```bash
systemctl status picoclaw | grep CPU
```

### Xem tất cả metrics
```bash
systemctl show picoclaw
```

---

## 🔔 Notifications

### Email khi service fail (với postfix)
```bash
# Install mailutils
sudo apt-get install mailutils

# Tạo script notify
sudo nano /usr/local/bin/picoclaw-notify.sh
```

```bash
#!/bin/bash
echo "PicoClaw service failed at $(date)" | mail -s "PicoClaw Alert" admin@example.com
```

```bash
# Make executable
sudo chmod +x /usr/local/bin/picoclaw-notify.sh

# Add to service
sudo systemctl edit picoclaw.service
```

```ini
[Service]
ExecStopPost=/usr/local/bin/picoclaw-notify.sh
```

---

## 📝 Best Practices

1. **Always check logs** sau khi start/restart
   ```bash
   sudo journalctl -u picoclaw -f
   ```

2. **Test configuration** trước khi restart
   ```bash
   picoclaw gateway --dry-run  # Nếu có option này
   ```

3. **Backup configuration** trước khi update
   ```bash
   cp /root/.picoclaw/config.yaml /root/.picoclaw/config.yaml.backup
   ```

4. **Monitor resource usage** định kỳ
   ```bash
   systemctl status picoclaw
   ```

5. **Rotate logs** để tránh đầy disk
   ```bash
   sudo journalctl --vacuum-time=7d  # Xóa logs > 7 ngày
   sudo journalctl --vacuum-size=500M  # Giữ tối đa 500MB logs
   ```

---

## 🆘 Support

Nếu gặp vấn đề:

1. Check logs: `sudo journalctl -u picoclaw -f`
2. Check status: `sudo systemctl status picoclaw`
3. Check GitHub issues: https://github.com/sipeed/picoclaw/issues
4. Check documentation: https://github.com/sipeed/picoclaw

---

**Created**: 2026-03-07  
**Version**: 1.0  
**Tested on**: Ubuntu 20.04+, Debian 11+, CentOS 8+
