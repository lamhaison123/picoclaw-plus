#!/bin/bash
# PicoClaw systemd installation script

set -e

echo "🚀 Installing PicoClaw systemd service..."

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    echo "❌ Please run as root (use sudo)"
    exit 1
fi

# Check if picoclaw binary exists
if ! command -v picoclaw &> /dev/null; then
    echo "❌ picoclaw binary not found in PATH"
    echo "Please install picoclaw first or add it to /usr/local/bin/"
    exit 1
fi

# Get picoclaw binary location
PICOCLAW_BIN=$(which picoclaw)
echo "✓ Found picoclaw at: $PICOCLAW_BIN"

# Create systemd service file
SERVICE_FILE="/etc/systemd/system/picoclaw.service"

cat > "$SERVICE_FILE" << 'EOF'
[Unit]
Description=PicoClaw - Ultra-lightweight Personal AI Agent Gateway
Documentation=https://github.com/sipeed/picoclaw
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=root
Group=root
WorkingDirectory=/root
ExecStart=/usr/local/bin/picoclaw gateway
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal
SyslogIdentifier=picoclaw

# Environment variables
Environment="HOME=/root"
Environment="PATH=/usr/local/bin:/usr/bin:/bin"

# Security hardening
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=read-only
ReadWritePaths=/root/.picoclaw

# Resource limits
LimitNOFILE=65536
LimitNPROC=4096

# Graceful shutdown
TimeoutStopSec=30
KillMode=mixed
KillSignal=SIGTERM

[Install]
WantedBy=multi-user.target
EOF

# Update ExecStart with actual binary path
sed -i "s|/usr/local/bin/picoclaw|$PICOCLAW_BIN|g" "$SERVICE_FILE"

echo "✓ Created systemd service file: $SERVICE_FILE"

# Reload systemd
systemctl daemon-reload
echo "✓ Reloaded systemd daemon"

# Enable service
systemctl enable picoclaw.service
echo "✓ Enabled picoclaw service (will start on boot)"

echo ""
echo "✅ Installation complete!"
echo ""
echo "📋 Available commands:"
echo "  sudo systemctl start picoclaw    # Start the service"
echo "  sudo systemctl stop picoclaw     # Stop the service"
echo "  sudo systemctl restart picoclaw  # Restart the service"
echo "  sudo systemctl status picoclaw   # Check service status"
echo "  sudo journalctl -u picoclaw -f   # View logs (follow mode)"
echo ""
echo "🚀 To start PicoClaw now, run:"
echo "  sudo systemctl start picoclaw"
