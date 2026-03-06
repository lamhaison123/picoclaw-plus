#!/bin/bash
# IRC Gateway Setup Script for PicoClaw

set -e

echo "========================================="
echo "  PicoClaw IRC Gateway Setup"
echo "========================================="
echo ""

# Check if picoclaw is installed
if ! command -v picoclaw &> /dev/null; then
    echo "❌ Error: picoclaw not found in PATH"
    echo "   Please install picoclaw first"
    exit 1
fi

echo "✅ PicoClaw found: $(which picoclaw)"
echo ""

# Check Python
if ! command -v python3 &> /dev/null; then
    echo "❌ Error: Python 3 not found"
    echo "   Please install Python 3.8+"
    exit 1
fi

echo "✅ Python found: $(python3 --version)"
echo ""

# Install dependencies
echo "[1/4] Installing Python dependencies..."
pip3 install python-telegram-bot --quiet
echo "✅ Dependencies installed"
echo ""

# Setup team configuration
echo "[2/4] Setting up team configuration..."
TEAM_DIR="$HOME/.picoclaw/workspace/teams"
mkdir -p "$TEAM_DIR"

if [ ! -f "$TEAM_DIR/irc-dev-team.json" ]; then
    cp irc-dev-team.json "$TEAM_DIR/"
    echo "✅ Team config copied to $TEAM_DIR"
else
    echo "⚠️  Team config already exists, skipping"
fi
echo ""

# Create team
echo "[3/4] Creating team in PicoClaw..."
if picoclaw team status irc-dev-team &> /dev/null; then
    echo "⚠️  Team 'irc-dev-team' already exists"
else
    picoclaw team create "$TEAM_DIR/irc-dev-team.json"
    echo "✅ Team created successfully"
fi
echo ""

# Check .env file
echo "[4/4] Checking environment configuration..."
ENV_FILE="../../../../.env"

if [ ! -f "$ENV_FILE" ]; then
    echo "⚠️  .env file not found in repo root"
    echo ""
    echo "Please create $ENV_FILE with:"
    echo ""
    echo "TELEGRAM_BOT_TOKEN=your_token_here"
    echo "PICOCLAW_BIN=picoclaw"
    echo "IRC_TEAM_ID=irc-dev-team"
    echo ""
    echo "Get your bot token from @BotFather on Telegram"
else
    if grep -q "TELEGRAM_BOT_TOKEN" "$ENV_FILE"; then
        echo "✅ .env file configured"
    else
        echo "⚠️  TELEGRAM_BOT_TOKEN not found in .env"
        echo "   Please add your bot token"
    fi
fi
echo ""

echo "========================================="
echo "  Setup Complete!"
echo "========================================="
echo ""
echo "Next steps:"
echo "1. Add TELEGRAM_BOT_TOKEN to .env file (if not done)"
echo "2. Run: python3 gateway.py"
echo "3. Start chatting with your bot on Telegram!"
echo ""
echo "Commands:"
echo "  /start  - Initialize bot"
echo "  /who    - List roles"
echo "  /team   - Show team status"
echo ""
