# Environment Variable Configuration Complete - v0.2.1

## ✅ Feature Complete

**Time**: 30 minutes  
**Status**: COMPLETE  
**Build**: PASSING ✅

## What Was Implemented

### .env File Support
Added support for loading configuration from .env files, following 12-factor app principles.

**Benefits**:
- Keep secrets out of config.json
- Environment-specific configuration
- Easy deployment (Docker, Kubernetes, etc.)
- CI/CD friendly
- No code changes needed for different environments

## Implementation Details

### New File: `pkg/config/env.go`
Created dedicated file for .env loading:

**Functions**:
1. `LoadEnvFile(path string)` - Load single .env file
2. `LoadEnvFiles(paths ...string)` - Load multiple .env files

**Features**:
- Parses KEY=VALUE format
- Skips empty lines and comments (#)
- Removes quotes from values
- Validates keys (no empty keys)
- Existing env vars take precedence
- Graceful error handling (logs warnings, doesn't fail)

### Modified: `pkg/config/config.go`
Updated `LoadConfig()` to load .env files before config:

**Load Order**:
1. `$PICOCLAW_HOME/.env` (if PICOCLAW_HOME is set)
2. `.env` (current directory)
3. `.env.local` (local overrides)
4. `config.json` (configuration file)
5. Environment variables (final override)

**Precedence** (highest to lowest):
```
Environment Variables > .env.local > .env > config.json > defaults
```

### New File: `.env.example`
Created comprehensive example with all available variables:

**Sections**:
- LLM Provider API Keys (OpenAI, Anthropic, etc.)
- Agent Configuration
- Tool Configuration
- Session Configuration
- Memory Configuration
- Channel Configuration
- Gateway Configuration
- Heartbeat Configuration
- Skills Configuration
- Advanced Configuration

### New File: `.gitignore`
Created .gitignore to prevent committing secrets:

**Ignored**:
- `.env`
- `.env.local`
- `.env.*.local`
- Build artifacts
- Workspace data
- IDE files

## Usage

### 1. Create .env File
```bash
# Copy example
cp .env.example .env

# Edit with your values
nano .env
```

### 2. Add Your Secrets
```bash
# .env
OPENAI_API_KEY=sk-...
ANTHROPIC_API_KEY=sk-ant-...
DEEPSEEK_API_KEY=sk-...

# Tool configuration
PICOCLAW_TOOLS_SHELL_ENABLED=false
PICOCLAW_TOOLS_WEB_ENABLED=true

# Summarization
PICOCLAW_SESSION_SUMMARIZATION_MESSAGE_THRESHOLD=30
```

### 3. Run PicoClaw
```bash
# .env is loaded automatically
./picoclaw

# Or specify custom config
./picoclaw -config custom.json
```

### 4. Override with Environment Variables
```bash
# Override specific values
PICOCLAW_TOOLS_SHELL_ENABLED=true ./picoclaw

# Multiple overrides
OPENAI_API_KEY=sk-test \
PICOCLAW_AGENTS_DEFAULTS_MODEL=gpt-4 \
./picoclaw
```

## Environment Variable Names

### Naming Convention
All PicoClaw env vars use `PICOCLAW_` prefix:
```
PICOCLAW_<SECTION>_<SUBSECTION>_<FIELD>
```

### Examples
```bash
# Agent defaults
PICOCLAW_AGENTS_DEFAULTS_WORKSPACE=./workspace
PICOCLAW_AGENTS_DEFAULTS_MODEL=gpt-4
PICOCLAW_AGENTS_DEFAULTS_MAX_TOKENS=32768

# Tools
PICOCLAW_TOOLS_FILE_ENABLED=true
PICOCLAW_TOOLS_SHELL_ENABLED=false
PICOCLAW_TOOLS_WEB_ENABLED=true

# Session
PICOCLAW_SESSION_SUMMARIZATION_MESSAGE_THRESHOLD=20
PICOCLAW_SESSION_SUMMARIZATION_TOKEN_PERCENT=0.75

# Memory
PICOCLAW_MEMORY_ENABLED=true
PICOCLAW_MEMORY_QDRANT_URL=http://localhost:6333

# Channels
PICOCLAW_CHANNELS_TELEGRAM_ENABLED=true
PICOCLAW_CHANNELS_TELEGRAM_TOKEN=123456:ABC...
```

## .env File Format

### Basic Format
```bash
# Comments start with #
KEY=VALUE

# Quotes are optional and will be removed
API_KEY="sk-..."
API_KEY='sk-...'
API_KEY=sk-...

# All three are equivalent
```

### Boolean Values
```bash
# Case-insensitive
ENABLED=true
ENABLED=True
ENABLED=TRUE

DISABLED=false
DISABLED=False
DISABLED=FALSE
```

### Numeric Values
```bash
# Integers
PORT=8080
THRESHOLD=20

# Floats
TEMPERATURE=0.75
PERCENT=0.8
```

### String Values
```bash
# No quotes needed
MODEL=gpt-4
WORKSPACE=./workspace

# Quotes optional
API_KEY="sk-..."
URL='http://localhost:6333'
```

## Deployment Scenarios

### 1. Development
```bash
# .env
OPENAI_API_KEY=sk-dev-...
PICOCLAW_AGENTS_DEFAULTS_MODEL=gpt-3.5-turbo
PICOCLAW_TOOLS_SHELL_ENABLED=true
```

### 2. Production
```bash
# .env
OPENAI_API_KEY=sk-prod-...
PICOCLAW_AGENTS_DEFAULTS_MODEL=gpt-4
PICOCLAW_TOOLS_SHELL_ENABLED=false
PICOCLAW_SESSION_SUMMARIZATION_MESSAGE_THRESHOLD=30
```

### 3. Docker
```dockerfile
# Dockerfile
FROM golang:1.21
COPY . /app
WORKDIR /app
RUN go build ./cmd/picoclaw

# Use env vars from docker-compose or -e flags
CMD ["./picoclaw"]
```

```yaml
# docker-compose.yml
services:
  picoclaw:
    build: .
    environment:
      - OPENAI_API_KEY=${OPENAI_API_KEY}
      - PICOCLAW_AGENTS_DEFAULTS_MODEL=gpt-4
    env_file:
      - .env
```

### 4. Kubernetes
```yaml
# secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: picoclaw-secrets
type: Opaque
stringData:
  OPENAI_API_KEY: sk-...
  ANTHROPIC_API_KEY: sk-ant-...

---
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: picoclaw
spec:
  template:
    spec:
      containers:
      - name: picoclaw
        image: picoclaw:latest
        envFrom:
        - secretRef:
            name: picoclaw-secrets
        env:
        - name: PICOCLAW_AGENTS_DEFAULTS_MODEL
          value: "gpt-4"
```

### 5. CI/CD
```yaml
# .github/workflows/test.yml
jobs:
  test:
    runs-on: ubuntu-latest
    env:
      OPENAI_API_KEY: ${{ secrets.OPENAI_API_KEY }}
      PICOCLAW_TOOLS_SHELL_ENABLED: false
    steps:
      - uses: actions/checkout@v2
      - run: go test ./...
```

## Security Best Practices

### 1. Never Commit .env
```bash
# .gitignore already includes:
.env
.env.local
.env.*.local
```

### 2. Use .env.example
```bash
# Commit .env.example with dummy values
# Team members copy to .env and fill in real values
cp .env.example .env
```

### 3. Rotate Keys Regularly
```bash
# Update .env with new keys
# Restart PicoClaw
./picoclaw
```

### 4. Use Secret Management
```bash
# AWS Secrets Manager
aws secretsmanager get-secret-value --secret-id picoclaw/openai | \
  jq -r .SecretString | \
  jq -r .OPENAI_API_KEY

# HashiCorp Vault
vault kv get -field=OPENAI_API_KEY secret/picoclaw
```

## Troubleshooting

### .env Not Loading
```bash
# Check file exists
ls -la .env

# Check file format (no BOM, Unix line endings)
file .env

# Check logs
./picoclaw 2>&1 | grep "Loaded .env"
```

### Environment Variable Not Working
```bash
# Check precedence
echo $OPENAI_API_KEY  # Existing env var takes precedence

# Unset to use .env value
unset OPENAI_API_KEY
./picoclaw
```

### Invalid .env Format
```bash
# Check for errors in logs
./picoclaw 2>&1 | grep "Invalid .env line"

# Common issues:
# - Missing = sign
# - Empty key
# - Special characters not quoted
```

## Files Created/Modified

1. ✅ `pkg/config/env.go` - NEW (100 lines)
2. ✅ `pkg/config/config.go` - MODIFIED (added .env loading)
3. ✅ `.env.example` - NEW (comprehensive example)
4. ✅ `.gitignore` - NEW (prevent committing secrets)

## Build Status

```bash
$ go build -tags=no_qdrant ./cmd/picoclaw
# Exit Code: 0 ✅
```

## Next Steps

1. ✅ Environment variable config - DONE
2. ⬜ Document in user guide
3. ⬜ Add to deployment docs
4. ⬜ Create Docker example
5. ⬜ Create Kubernetes example

---

**Date**: 2026-03-09  
**Status**: COMPLETE ✅  
**Impact**: Better DevOps, 12-factor app compliance  
**Build**: PASSING ✅  
**Next**: Start JSONL Memory Store or Model Routing
