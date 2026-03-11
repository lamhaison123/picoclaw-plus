# PicoClaw Migration Guide

Complete guide for migrating between PicoClaw versions.

## Table of Contents

- [v0.2.0 to v0.2.1](#v020-to-v021)
- [JSON to JSONL Storage](#json-to-jsonl-storage)
- [Legacy Providers to Model List](#legacy-providers-to-model-list)
- [Configuration Updates](#configuration-updates)
- [Breaking Changes](#breaking-changes)

---

## v0.2.0 to v0.2.1

### Overview

v0.2.1 introduces 10 new features with full backward compatibility. No breaking changes.

### New Features

1. **JSONL Memory Store** - Crash-safe storage
2. **Vision/Image Support** - Multi-modal AI
3. **Parallel Tool Execution** - 2x faster
4. **Model Routing** - Cost optimization
5. **Environment Configuration** - .env support
6. **Tool Enable/Disable** - Granular control
7. **Extended Thinking** - AI reasoning
8. **Configurable Summarization** - Flexible thresholds
9. **PICOCLAW_HOME** - Custom home directory
10. **New Search Providers** - SearXNG, GLM, Exa

### Migration Steps

#### Step 1: Backup Your Data

```bash
# Backup configuration
cp ~/.picoclaw/config.json ~/.picoclaw/config.json.backup

# Backup sessions
cp -r ~/.picoclaw/sessions ~/.picoclaw/sessions.backup

# Backup auth
cp ~/.picoclaw/auth.json ~/.picoclaw/auth.json.backup
```

#### Step 2: Update Binary

```bash
# Pull latest code
git pull origin main

# Build new version
go build -o build/picoclaw ./cmd/picoclaw

# Verify version
./build/picoclaw version
```

#### Step 3: Update Configuration (Optional)

All new features are opt-in. Your existing config will continue to work.

**Enable JSONL Storage** (Recommended):
```json
{
  "session": {
    "storage_backend": "jsonl",
    "auto_migrate": true
  }
}
```

**Enable Model Routing** (Optional):
```json
{
  "routing": {
    "enabled": true,
    "tiers": [
      {"name": "cheap", "models": ["gpt-4o-mini"]},
      {"name": "expensive", "models": ["gpt-5.2"]}
    ]
  }
}
```

**Configure New Search Providers** (Optional):
```json
{
  "tools": {
    "web": {
      "searxng": {
        "enabled": true,
        "base_url": "https://searx.example.com"
      },
      "glm": {
        "enabled": true,
        "api_key": "your-glm-key"
      },
      "exa": {
        "enabled": true,
        "api_key": "your-exa-key"
      }
    }
  }
}
```

#### Step 4: Restart PicoClaw

```bash
# Stop old version
pkill picoclaw

# Start new version
./build/picoclaw agent
```

#### Step 5: Verify Migration

```bash
# Check sessions migrated to JSONL
ls ~/.picoclaw/sessions/*.jsonl

# Check logs for migration messages
tail -f ~/.picoclaw/logs/picoclaw.log | grep migration

# Test basic functionality
./build/picoclaw agent
> Hello
> Exit
```

### Rollback Procedure

If you encounter issues:

```bash
# Stop new version
pkill picoclaw

# Restore old binary
cp build/picoclaw.old build/picoclaw

# Restore configuration
cp ~/.picoclaw/config.json.backup ~/.picoclaw/config.json

# Restore sessions (if needed)
rm -rf ~/.picoclaw/sessions
cp -r ~/.picoclaw/sessions.backup ~/.picoclaw/sessions

# Start old version
./build/picoclaw agent
```

---

## JSON to JSONL Storage

### Why Migrate?

- **Crash Safety**: No data loss on power failure
- **Performance**: Faster append operations
- **Reliability**: Atomic writes with fsync

### Automatic Migration

**Recommended Method** - Let PicoClaw migrate automatically:

```json
{
  "session": {
    "storage_backend": "jsonl",
    "auto_migrate": true
  }
}
```

**How it works**:
1. On first session access, PicoClaw detects JSON file
2. Reads all messages from JSON
3. Writes to new JSONL file
4. Preserves old JSON as backup
5. Uses JSONL for all future operations

### Manual Migration

If you prefer manual control:

```bash
# Create migration script
cat > migrate-sessions.sh << 'EOF'
#!/bin/bash
for json_file in ~/.picoclaw/sessions/*.json; do
  session_id=$(basename "$json_file" .json)
  echo "Migrating: $session_id"
  
  # PicoClaw will auto-migrate on first access
  # Just ensure auto_migrate is enabled in config
done
EOF

chmod +x migrate-sessions.sh
./migrate-sessions.sh
```

### Verification

```bash
# Check JSONL files created
ls -lh ~/.picoclaw/sessions/*.jsonl

# Check old JSON files preserved
ls -lh ~/.picoclaw/sessions/*.json

# Compare file sizes
du -sh ~/.picoclaw/sessions/
```

### Troubleshooting

**Issue**: Migration not happening
```bash
# Check config
cat ~/.picoclaw/config.json | grep storage_backend

# Check logs
tail -f ~/.picoclaw/logs/picoclaw.log | grep migration
```

**Issue**: Corrupted JSONL file
```bash
# Restore from JSON backup
rm ~/.picoclaw/sessions/session-id.jsonl
# PicoClaw will re-migrate from JSON on next access
```

---

## Legacy Providers to Model List

### Why Migrate?

The old `providers` format is deprecated. Use `model_list` for:
- Better model management
- Load balancing support
- Clearer configuration
- More flexibility

### Old Format (Deprecated)

```json
{
  "providers": {
    "openai": {
      "api_key": "sk-...",
      "model": "gpt-4o"
    },
    "anthropic": {
      "api_key": "sk-ant-...",
      "model": "claude-3-opus"
    }
  }
}
```

### New Format (Recommended)

```json
{
  "model_list": [
    {
      "model_name": "gpt-4o",
      "model": "openai/gpt-4o",
      "api_base": "https://api.openai.com/v1",
      "api_key": "${OPENAI_API_KEY}"
    },
    {
      "model_name": "claude-opus",
      "model": "anthropic/claude-3-opus-20240229",
      "api_base": "https://api.anthropic.com",
      "api_key": "${ANTHROPIC_API_KEY}"
    }
  ]
}
```

### Migration Steps

1. **Identify Current Providers**:
```bash
cat ~/.picoclaw/config.json | grep -A 10 '"providers"'
```

2. **Convert to Model List**:

Use this conversion table:

| Old Provider | New Model Format |
|--------------|------------------|
| `openai` | `openai/gpt-4o` |
| `anthropic` | `anthropic/claude-3-opus-20240229` |
| `deepseek` | `deepseek/deepseek-chat` |
| `qwen` | `qwen/qwen-max` |
| `zhipu` | `zhipu/glm-4` |
| `moonshot` | `moonshot/moonshot-v1-8k` |

3. **Update Configuration**:
```bash
# Backup old config
cp ~/.picoclaw/config.json ~/.picoclaw/config.json.old

# Edit config
nano ~/.picoclaw/config.json
# Replace providers section with model_list
```

4. **Verify**:
```bash
# Test configuration
./build/picoclaw agent --dry-run

# Check model list loaded
./build/picoclaw agent
> /models
```

### Automatic Conversion

PicoClaw automatically converts old format to new format internally. You can continue using the old format, but migration is recommended.

---

## Configuration Updates

### Environment Variables

v0.2.1 adds .env file support. Migrate secrets to .env:

**Before** (config.json):
```json
{
  "model_list": [
    {
      "api_key": "sk-actual-key-here"
    }
  ]
}
```

**After** (.env + config.json):

`.env`:
```bash
OPENAI_API_KEY=sk-actual-key-here
ANTHROPIC_API_KEY=sk-ant-actual-key
```

`config.json`:
```json
{
  "model_list": [
    {
      "api_key": "${OPENAI_API_KEY}"
    }
  ]
}
```

### Tool Configuration

v0.2.1 adds individual tool flags:

**Add to config.json**:
```json
{
  "tools": {
    "file_tools_enabled": true,
    "shell_tools_enabled": true,
    "web_tools_enabled": true,
    "message_tool_enabled": true,
    "spawn_tool_enabled": true,
    "team_tools_enabled": true,
    "skill_tools_enabled": true,
    "hardware_tools_enabled": false
  }
}
```

### Session Configuration

v0.2.1 adds session options:

**Add to config.json**:
```json
{
  "session": {
    "storage_backend": "jsonl",
    "auto_migrate": true,
    "summarization_message_threshold": 20,
    "summarization_token_percent": 0.75
  }
}
```

---

## Breaking Changes

### None in v0.2.1

v0.2.1 has **zero breaking changes**. All existing configurations continue to work.

### Deprecations

The following are deprecated but still supported:

1. **Old Provider Format**: Use `model_list` instead
2. **JSON Storage**: Use JSONL for better reliability
3. **Hardcoded API Keys**: Use environment variables

---

## Platform-Specific Migration

### Docker

Update your Dockerfile:

```dockerfile
FROM golang:1.21

# Set custom home directory
ENV PICOCLAW_HOME=/data/picoclaw

# Create directories
RUN mkdir -p /data/picoclaw

# Copy binary
COPY build/picoclaw /usr/local/bin/

# Volume for persistent data
VOLUME /data/picoclaw

CMD ["picoclaw", "agent"]
```

Update docker-compose.yml:

```yaml
version: '3'
services:
  picoclaw:
    image: picoclaw:v0.2.1
    environment:
      - PICOCLAW_HOME=/data
      - PICOCLAW_SESSION_STORAGE_BACKEND=jsonl
      - OPENAI_API_KEY=${OPENAI_API_KEY}
    volumes:
      - ./picoclaw-data:/data
    ports:
      - "8080:8080"
```

### Systemd

Update service file:

```ini
[Unit]
Description=PicoClaw Agent v0.2.1
After=network.target

[Service]
Type=simple
User=picoclaw
Environment="PICOCLAW_HOME=/opt/picoclaw"
Environment="PICOCLAW_SESSION_STORAGE_BACKEND=jsonl"
EnvironmentFile=/opt/picoclaw/.env
ExecStart=/usr/local/bin/picoclaw agent
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

Reload and restart:

```bash
sudo systemctl daemon-reload
sudo systemctl restart picoclaw
sudo systemctl status picoclaw
```

### Kubernetes

Update deployment:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: picoclaw
spec:
  replicas: 1
  template:
    spec:
      containers:
      - name: picoclaw
        image: picoclaw:v0.2.1
        env:
        - name: PICOCLAW_HOME
          value: /data
        - name: PICOCLAW_SESSION_STORAGE_BACKEND
          value: jsonl
        envFrom:
        - secretRef:
            name: picoclaw-secrets
        volumeMounts:
        - name: data
          mountPath: /data
      volumes:
      - name: data
        persistentVolumeClaim:
          claimName: picoclaw-data
```

---

## Post-Migration Checklist

### Verify Installation

- [ ] Binary version correct: `picoclaw version`
- [ ] Configuration valid: `picoclaw agent --dry-run`
- [ ] Sessions migrated: `ls ~/.picoclaw/sessions/*.jsonl`
- [ ] Logs clean: `tail -f ~/.picoclaw/logs/picoclaw.log`

### Test Functionality

- [ ] Basic chat works
- [ ] Tools execute correctly
- [ ] Sessions persist
- [ ] Models respond
- [ ] Channels connect (if enabled)

### Performance Check

- [ ] Response times normal
- [ ] Memory usage acceptable
- [ ] CPU usage reasonable
- [ ] Disk I/O healthy

### Security Review

- [ ] API keys in .env (not config.json)
- [ ] File permissions correct (600 for sensitive files)
- [ ] .env in .gitignore
- [ ] No secrets in logs

---

## Troubleshooting

### Common Issues

**Issue**: Config not loading
```bash
# Check syntax
cat ~/.picoclaw/config.json | jq .

# Check permissions
ls -la ~/.picoclaw/config.json

# Check PICOCLAW_HOME
echo $PICOCLAW_HOME
```

**Issue**: Sessions not migrating
```bash
# Enable auto-migration
# In config.json:
{
  "session": {
    "auto_migrate": true
  }
}

# Check logs
tail -f ~/.picoclaw/logs/picoclaw.log | grep migration
```

**Issue**: Models not found
```bash
# List available models
picoclaw agent
> /models

# Check model_list in config
cat ~/.picoclaw/config.json | jq .model_list
```

### Getting Help

- **Documentation**: See `docs/` directory
- **Issues**: https://github.com/sipeed/picoclaw/issues
- **Discussions**: https://github.com/sipeed/picoclaw/discussions

---

## See Also

- [Configuration Reference](../reference/CONFIGURATION.md)
- [v0.2.1 Features Guide](../guides/V0.2.1_FEATURES.md)
- [JSONL Storage Reference](../reference/JSONL_STORAGE.md)
- [Quick Start Guide](../QUICK_START.md)

---

**Last Updated**: 2026-03-09  
**Version**: v0.2.1
