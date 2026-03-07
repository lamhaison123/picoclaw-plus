# PicoClaw Systemd - Quick Reference

Cheat sheet nhanh cho các lệnh thường dùng.

---

## 🚀 Installation

```bash
sudo bash install-systemd.sh
```

---

## 🎮 Service Control

| Command | Description |
|---------|-------------|
| `sudo systemctl start picoclaw` | Khởi động service |
| `sudo systemctl stop picoclaw` | Dừng service |
| `sudo systemctl restart picoclaw` | Restart service |
| `sudo systemctl status picoclaw` | Xem trạng thái |
| `sudo systemctl enable picoclaw` | Enable auto-start |
| `sudo systemctl disable picoclaw` | Disable auto-start |

---

## 📊 Logs

| Command | Description |
|---------|-------------|
| `sudo journalctl -u picoclaw -f` | Follow logs realtime |
| `sudo journalctl -u picoclaw -n 100` | Last 100 lines |
| `sudo journalctl -u picoclaw -b` | Since last boot |
| `sudo journalctl -u picoclaw -p err` | Only errors |
| `sudo journalctl -u picoclaw --since "1h ago"` | Last 1 hour |
| `sudo journalctl -u picoclaw > log.txt` | Export to file |

---

## 🔧 Configuration

| File | Location |
|------|----------|
| Service file | `/etc/systemd/system/picoclaw.service` |
| Binary | `/usr/local/bin/picoclaw` |
| Workspace | `/root/.picoclaw` |
| Config | `/root/.picoclaw/config.yaml` |

---

## 🛠️ Helper Script

```bash
# Make executable
chmod +x picoclaw-ctl.sh

# Commands
sudo ./picoclaw-ctl.sh status    # Status
sudo ./picoclaw-ctl.sh start     # Start
sudo ./picoclaw-ctl.sh stop      # Stop
sudo ./picoclaw-ctl.sh restart   # Restart
sudo ./picoclaw-ctl.sh logs 100  # Logs
sudo ./picoclaw-ctl.sh follow    # Follow logs
```

---

## 🔍 Troubleshooting

### Service won't start
```bash
sudo systemctl status picoclaw
sudo journalctl -u picoclaw -n 50
```

### Port conflict
```bash
sudo lsof -i :18790
sudo kill -9 <PID>
```

### Permission issues
```bash
sudo chown -R root:root /root/.picoclaw
sudo chmod -R 755 /root/.picoclaw
```

### Reload after config change
```bash
sudo systemctl daemon-reload
sudo systemctl restart picoclaw
```

---

## 📈 Monitoring

```bash
# Uptime
systemctl show picoclaw --property=ActiveEnterTimestamp

# Memory
systemctl status picoclaw | grep Memory

# CPU
systemctl status picoclaw | grep CPU

# All metrics
systemctl show picoclaw
```

---

## 🔄 Update

```bash
sudo systemctl stop picoclaw
sudo cp new-picoclaw /usr/local/bin/picoclaw
sudo chmod +x /usr/local/bin/picoclaw
sudo systemctl start picoclaw
```

---

## 📦 Uninstall

```bash
sudo systemctl stop picoclaw
sudo systemctl disable picoclaw
sudo rm /etc/systemd/system/picoclaw.service
sudo systemctl daemon-reload
```

---

## 🆘 Emergency

### Force stop
```bash
sudo systemctl kill picoclaw
```

### Prevent auto-restart
```bash
sudo systemctl mask picoclaw
```

### Re-enable
```bash
sudo systemctl unmask picoclaw
```

---

**Quick Help**: `./picoclaw-ctl.sh help`
