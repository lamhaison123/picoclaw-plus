#!/bin/bash
# IRC Gateway Migration Script: v1.0.0 → v1.0.1
# Automatically migrates to new version with token auto-loading

set -e

echo "========================================="
echo "  IRC Gateway Migration: v1.0.0 → v1.0.1"
echo "========================================="
echo ""

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Check if running from correct directory
if [ ! -f "gateway.py" ]; then
    echo -e "${RED}❌ Error: Please run this script from workspace/skills/irc-gateway/${NC}"
    exit 1
fi

echo "Step 1: Checking current setup..."
echo "-----------------------------------"

# Check for existing .env
if [ -f "../../../../.env" ]; then
    echo -e "${GREEN}✅ Found .env file in repo root${NC}"
    HAS_ENV=true
else
    echo -e "${YELLOW}⚠️  No .env file found${NC}"
    HAS_ENV=false
fi

# Check for PicoClaw config
PICOCLAW_CONFIG="$HOME/.picoclaw/config.json"
if [ -f "$PICOCLAW_CONFIG" ]; then
    echo -e "${GREEN}✅ Found PicoClaw config: $PICOCLAW_CONFIG${NC}"
    
    # Check if Telegram is configured
    if grep -q "telegram" "$PICOCLAW_CONFIG"; then
        echo -e "${GREEN}✅ Telegram configured in PicoClaw${NC}"
        HAS_PICOCLAW_TOKEN=true
    else
        echo -e "${YELLOW}⚠️  Telegram not configured in PicoClaw${NC}"
        HAS_PICOCLAW_TOKEN=false
    fi
else
    echo -e "${YELLOW}⚠️  PicoClaw config not found${NC}"
    HAS_PICOCLAW_TOKEN=false
fi

echo ""
echo "Step 2: Determining migration path..."
echo "-----------------------------------"

if [ "$HAS_PICOCLAW_TOKEN" = true ]; then
    echo -e "${GREEN}✅ Recommended: Use PicoClaw Telegram config${NC}"
    echo "   Gateway will automatically use token from PicoClaw"
    echo ""
    read -p "Do you want to remove .env file to use PicoClaw config? (y/N): " -n 1 -r
    echo ""
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        if [ "$HAS_ENV" = true ]; then
            mv ../../../../.env ../../../../.env.backup
            echo -e "${GREEN}✅ Backed up .env to .env.backup${NC}"
        fi
        MIGRATION_PATH="picoclaw"
    else
        echo -e "${YELLOW}⚠️  Keeping .env file (higher priority than PicoClaw config)${NC}"
        MIGRATION_PATH="env"
    fi
elif [ "$HAS_ENV" = true ]; then
    echo -e "${GREEN}✅ Using existing .env file${NC}"
    MIGRATION_PATH="env"
else
    echo -e "${RED}❌ No token source found!${NC}"
    echo ""
    echo "Please configure token in one of these locations:"
    echo "1. ~/.picoclaw/config.json (recommended)"
    echo "2. .env file in repo root"
    echo "3. Environment variable TELEGRAM_BOT_TOKEN"
    echo ""
    exit 1
fi

echo ""
echo "Step 3: Backing up current gateway.py..."
echo "-----------------------------------"

if [ -f "gateway.py" ]; then
    cp gateway.py gateway.py.backup
    echo -e "${GREEN}✅ Backed up to gateway.py.backup${NC}"
fi

echo ""
echo "Step 4: Installing dependencies..."
echo "-----------------------------------"

# Install python-telegram-bot if not already installed
if python3 -c "import telegram" 2>/dev/null; then
    echo -e "${GREEN}✅ python-telegram-bot already installed${NC}"
else
    echo "Installing python-telegram-bot..."
    if pip3 install python-telegram-bot 2>/dev/null; then
        echo -e "${GREEN}✅ python-telegram-bot installed${NC}"
    else
        echo -e "${RED}❌ Failed to install python-telegram-bot${NC}"
        echo "   Please install manually: pip3 install python-telegram-bot"
        mv gateway.py.backup gateway.py
        exit 1
    fi
fi

echo ""
echo "Step 5: Testing new gateway..."
echo "-----------------------------------"

# Test import (without executing main)
if python3 -c "import sys; sys.path.insert(0, '.'); exec(open('gateway.py').read().replace('if __name__ == \"__main__\":', 'if False:'))" 2>/dev/null; then
    echo -e "${GREEN}✅ Gateway syntax check passed${NC}"
else
    echo -e "${RED}❌ Gateway has syntax errors${NC}"
    echo "   Restoring backup..."
    mv gateway.py.backup gateway.py
    exit 1
fi

echo ""
echo "Step 5: Running tests..."
echo "-----------------------------------"

if [ -f "test_simple.py" ]; then
    if python3 test_simple.py; then
        echo -e "${GREEN}✅ All tests passed${NC}"
    else
        echo -e "${RED}❌ Tests failed${NC}"
        echo "   Restoring backup..."
        mv gateway.py.backup gateway.py
        exit 1
    fi
else
    echo -e "${YELLOW}⚠️  test_simple.py not found, skipping tests${NC}"
fi

echo ""
echo "========================================="
echo "  Migration Complete!"
echo "========================================="
echo ""

case $MIGRATION_PATH in
    picoclaw)
        echo -e "${GREEN}✅ Gateway will use token from PicoClaw config${NC}"
        echo "   Location: $PICOCLAW_CONFIG"
        echo ""
        echo "To start gateway:"
        echo "  python3 gateway.py"
        ;;
    env)
        echo -e "${GREEN}✅ Gateway will use token from .env file${NC}"
        echo "   Location: ../../../../.env"
        echo ""
        echo "To start gateway:"
        echo "  python3 gateway.py"
        ;;
esac

echo ""
echo "Verify token source in startup logs:"
echo "  ✅ Loaded Telegram token from <source>"
echo ""
echo "Backup files created:"
echo "  - gateway.py.backup"
if [ -f "../../../../.env.backup" ]; then
    echo "  - ../../../../.env.backup"
fi
echo ""
echo "To rollback:"
echo "  mv gateway.py.backup gateway.py"
if [ -f "../../../../.env.backup" ]; then
    echo "  mv ../../../../.env.backup ../../../../.env"
fi
echo ""
echo "For more details, see MIGRATION.md"
echo ""
