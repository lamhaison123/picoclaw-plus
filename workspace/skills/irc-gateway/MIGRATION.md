# Migration Guide: IRC Gateway v1.0.0 → v1.0.1

Quick guide to upgrade to v1.0.1 with automatic token loading.

## What's New in v1.0.1

✅ **Automatic Token Loading** - Gateway now reads Telegram token from PicoClaw config  
✅ **Flexible Configuration** - Multiple token sources with priority order  
✅ **Better Logging** - Shows where token was loaded from  

## Migration Steps

### Option 1: Use Existing PicoClaw Telegram Config (Recommended)

If you already have Telegram configured in PicoClaw, no changes needed!

**1. Check your PicoClaw config:**
```bash
cat ~/.picoclaw/config.json | grep -A 5 telegram
```

**2. If you see this, you're good to go:**
```json
{
  "channels": {
    "telegram": {
      "enabled": true,
      "token": "123456789:ABCdefGHIjklMNOpqrsTUVwxyz"
    }
  }
}
```

**3. Update gateway.py:**
```bash
cd workspace/skills/irc-gateway
git pull  # Or copy new gateway.py
```

**4. Run gateway:**
```bash
python gateway.py
```

**5. Verify token loaded:**
Look for this in startup logs:
```
✅ Loaded Telegram token from /home/user/.picoclaw/config.json
```

**Done!** No need to maintain separate `.env` file.

---

### Option 2: Keep Using .env File

If you prefer to keep token in `.env` file:

**1. Update gateway.py:**
```bash
cd workspace/skills/irc-gateway
git pull  # Or copy new gateway.py
```

**2. Keep your existing .env:**
```env
TELEGRAM_BOT_TOKEN=your_token_here
PICOCLAW_BIN=picoclaw
IRC_TEAM_ID=irc-dev-team
```

**3. Run gateway:**
```bash
python gateway.py
```

**4. Verify token loaded:**
Look for this in startup logs:
```
✅ Loaded Telegram token from /path/to/repo/.env
```

**Done!** Gateway still works with `.env` file.

---

### Option 3: Use Environment Variable

For containerized or CI/CD environments:

**1. Update gateway.py:**
```bash
cd workspace/skills/irc-gateway
git pull  # Or copy new gateway.py
```

**2. Set environment variable:**
```bash
export TELEGRAM_BOT_TOKEN=your_token_here
```

**3. Run gateway:**
```bash
python gateway.py
```

**4. Verify token loaded:**
Look for this in startup logs:
```
✅ Loaded Telegram token from environment variable
```

**Done!** Gateway uses environment variable.

---

## Configuration Priority

The gateway checks for token in this order:

1. **Environment variable** `TELEGRAM_BOT_TOKEN` (highest priority)
2. **Repo .env file** `.env` in repository root
3. **PicoClaw config** `~/.picoclaw/config.json` (lowest priority)

First found token is used.

---

## Troubleshooting Migration

### Gateway can't find token

**Symptom:**
```
❌ Error: TELEGRAM_BOT_TOKEN not found in .env file
```

**Solution:**
Check all three sources:

```bash
# 1. Check environment variable
echo $TELEGRAM_BOT_TOKEN

# 2. Check .env file
cat .env | grep TELEGRAM_BOT_TOKEN

# 3. Check PicoClaw config
cat ~/.picoclaw/config.json | grep -A 5 telegram
```

Add token to at least one source.

---

### Token loaded from wrong source

**Symptom:**
Gateway loads token from `.env` but you want it from PicoClaw config.

**Solution:**
Remove token from higher priority sources:

```bash
# Remove from environment
unset TELEGRAM_BOT_TOKEN

# Remove from .env
sed -i '/TELEGRAM_BOT_TOKEN/d' .env

# Now gateway will use PicoClaw config
python gateway.py
```

---

### Multiple bots conflict

**Symptom:**
You have different tokens in different places.

**Solution:**
Decide which token to use and remove others:

**For PicoClaw integration (recommended):**
```bash
# Remove from .env
rm .env  # Or edit to remove TELEGRAM_BOT_TOKEN line

# Use token from ~/.picoclaw/config.json
python gateway.py
```

**For separate bot:**
```bash
# Keep .env with separate token
echo "TELEGRAM_BOT_TOKEN=separate_bot_token" > .env

# Gateway will use .env token (higher priority than config.json)
python gateway.py
```

---

## Rollback to v1.0.0

If you need to rollback:

**1. Restore old gateway.py:**
```bash
git checkout v1.0.0 workspace/skills/irc-gateway/gateway.py
```

**2. Ensure .env file exists:**
```bash
cat > .env << EOF
TELEGRAM_BOT_TOKEN=your_token_here
PICOCLAW_BIN=picoclaw
IRC_TEAM_ID=irc-dev-team
EOF
```

**3. Run gateway:**
```bash
python gateway.py
```

---

## Testing Migration

**1. Stop old gateway:**
```bash
# Find process
ps aux | grep gateway.py

# Kill it
kill <pid>
```

**2. Update to v1.0.1:**
```bash
cd workspace/skills/irc-gateway
git pull
```

**3. Test token loading:**
```bash
python gateway.py
```

**4. Check startup logs:**
```
✅ Loaded Telegram token from <source>
🚀 Starting PicoClaw IRC Gateway...
📁 Repo root: /path/to/picoclaw
🤖 Team ID: irc-dev-team
👥 Roles: architect, developer, tester, manager
✅ PicoClaw binary found: picoclaw
✅ Gateway is running. Press Ctrl+C to stop.
```

**5. Test in Telegram:**
Send `/start` to your bot and verify it responds.

---

## Migration Checklist

- [ ] Backup current gateway.py
- [ ] Update to v1.0.1
- [ ] Verify token source (config.json, .env, or env var)
- [ ] Test gateway startup
- [ ] Check startup logs for token source
- [ ] Test `/start` command in Telegram
- [ ] Test mention routing (`@architect test`)
- [ ] Verify team execution works
- [ ] Update documentation if needed

---

## Benefits After Migration

✅ **No duplicate configuration** - Use existing PicoClaw Telegram setup  
✅ **Flexible deployment** - Choose token source based on environment  
✅ **Better debugging** - Logs show where token came from  
✅ **Easier maintenance** - One less config file to manage  

---

## FAQ

**Q: Do I need to change anything if I'm already using .env?**  
A: No, v1.0.1 is backward compatible. Your .env file will still work.

**Q: Can I use the same token for both PicoClaw gateway and IRC gateway?**  
A: Yes, but only one can run at a time. They'll conflict if both try to use the same bot.

**Q: What if I want different tokens for different purposes?**  
A: Use .env for IRC gateway with a separate token, and keep PicoClaw's token in config.json.

**Q: Does this work on Windows?**  
A: Yes, all three token sources work on Windows.

**Q: How do I know which token source is being used?**  
A: Check the startup logs. Gateway prints: `✅ Loaded Telegram token from <source>`

---

## Support

If you encounter issues during migration:

1. Check [Troubleshooting](#troubleshooting-migration) section above
2. Review [CHANGELOG.md](CHANGELOG.md) for detailed changes
3. See [README.md](README.md) for full configuration guide
4. Open issue on GitHub if problem persists

---

## Version Compatibility

| Version | Token Loading | Backward Compatible |
|---------|---------------|---------------------|
| v1.0.0  | .env only | N/A |
| v1.0.1  | .env + config.json + env var | ✅ Yes |

---

*Migration guide last updated: 2026-03-05*
