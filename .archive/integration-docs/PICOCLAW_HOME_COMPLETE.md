# PICOCLAW_HOME Support Complete ✅

## Summary
Successfully completed PICOCLAW_HOME environment variable support for custom home directory configuration.

## Implementation Details

### Already Implemented (v0.2.1)
1. `pkg/config/defaults.go` - Workspace path resolution
2. `pkg/config/config.go` - .env file loading
3. `pkg/migrate/internal/common.go` - Migration support

### Newly Updated (Comprehensive Coverage)
1. `pkg/team/persistence.go` - Team state directory
2. `pkg/agent/context.go` - Global config directory
3. `pkg/agent/instance.go` - Agent workspace paths
4. `pkg/auth/store.go` - Auth file path
5. `cmd/picoclaw-launcher/internal/server/utils.go` - Launcher config path
6. `cmd/picoclaw/internal/helpers.go` - CLI config path
7. `.env.example` - Documentation

## Features
✅ Custom home directory via PICOCLAW_HOME  
✅ Fallback to ~/.picoclaw if not set  
✅ Consistent across all components  
✅ Team state directory support  
✅ Auth file support  
✅ Config file resolution  
✅ Workspace path resolution  
✅ Migration support  
✅ Documented in .env.example  

## How It Works

### Priority Order
1. **PICOCLAW_HOME** environment variable (highest priority)
2. **~/.picoclaw** (default fallback)
3. **Current directory** (last resort for some components)

### Directory Structure
```
$PICOCLAW_HOME/
├── config.json          # Main configuration
├── auth.json            # Authentication credentials
├── workspace/           # Default agent workspace
├── workspace-agent1/    # Agent-specific workspaces
├── teams/               # Team state files
├── skills/              # Global skills
└── .env                 # Environment variables
```

### Usage Examples

#### Set Custom Home Directory
```bash
export PICOCLAW_HOME=/custom/path/to/picoclaw
picoclaw agent
```

#### Multi-User Setup
```bash
# User 1
export PICOCLAW_HOME=/home/user1/.picoclaw
picoclaw agent

# User 2
export PICOCLAW_HOME=/home/user2/.picoclaw
picoclaw agent
```

#### Docker/Container Setup
```dockerfile
ENV PICOCLAW_HOME=/app/picoclaw
VOLUME /app/picoclaw
```

#### Systemd Service
```ini
[Service]
Environment="PICOCLAW_HOME=/var/lib/picoclaw"
ExecStart=/usr/local/bin/picoclaw agent
```

## Configuration

### Environment Variable
```bash
# Set in shell
export PICOCLAW_HOME=/custom/path

# Set in .env file
PICOCLAW_HOME=/custom/path

# Set in systemd
Environment="PICOCLAW_HOME=/custom/path"
```

### .env.example
```bash
# ============================================
# PicoClaw Home Directory (v0.2.1)
# ============================================
# Custom home directory for PicoClaw data
# Default: ~/.picoclaw
# PICOCLAW_HOME=/custom/path/to/picoclaw
```

## Components Updated

### Core Components
1. **Config Loading** - Uses PICOCLAW_HOME for config.json
2. **Workspace Resolution** - Uses PICOCLAW_HOME for workspaces
3. **Auth Storage** - Uses PICOCLAW_HOME for auth.json
4. **Team Persistence** - Uses PICOCLAW_HOME for team states
5. **Global Skills** - Uses PICOCLAW_HOME for skills directory

### CLI Tools
1. **picoclaw** - Main CLI respects PICOCLAW_HOME
2. **picoclaw-launcher** - GUI launcher respects PICOCLAW_HOME
3. **Migration tools** - Respect PICOCLAW_HOME

## Testing
- ✅ Build successful
- ✅ No diagnostics errors
- ✅ Backward compatible (defaults to ~/.picoclaw)
- ✅ All components updated consistently

## Files Modified
1. `pkg/team/persistence.go` - Added PICOCLAW_HOME support
2. `pkg/agent/context.go` - Added PICOCLAW_HOME support
3. `pkg/agent/instance.go` - Added PICOCLAW_HOME support
4. `pkg/auth/store.go` - Added PICOCLAW_HOME support
5. `cmd/picoclaw-launcher/internal/server/utils.go` - Added PICOCLAW_HOME support
6. `cmd/picoclaw/internal/helpers.go` - Added PICOCLAW_HOME support
7. `.env.example` - Added documentation

## Files Already Complete (v0.2.1)
1. `pkg/config/defaults.go` - Workspace path
2. `pkg/config/config.go` - .env loading
3. `pkg/migrate/internal/common.go` - Migration

## Use Cases

### Development
```bash
# Separate dev environment
export PICOCLAW_HOME=~/picoclaw-dev
picoclaw agent
```

### Production
```bash
# System-wide installation
export PICOCLAW_HOME=/opt/picoclaw
picoclaw agent
```

### Testing
```bash
# Isolated test environment
export PICOCLAW_HOME=/tmp/picoclaw-test
picoclaw agent
```

### Multi-Tenant
```bash
# Tenant 1
export PICOCLAW_HOME=/data/tenant1/picoclaw
picoclaw agent

# Tenant 2
export PICOCLAW_HOME=/data/tenant2/picoclaw
picoclaw agent
```

## Benefits
- **Flexibility**: Custom installation paths
- **Multi-User**: Separate data per user
- **Isolation**: Test/dev/prod separation
- **Portability**: Easy to relocate data
- **Docker-Friendly**: Volume mounting
- **Multi-Tenant**: Separate tenant data

## Backward Compatibility
- Default behavior unchanged (~/.picoclaw)
- Existing installations work without changes
- No breaking changes
- Opt-in via environment variable

## v0.2.1 Integration Status
✅ Tool Enable/Disable  
✅ Configurable Summarization  
✅ Parallel Tool Execution (v0.2.1 inline)  
✅ Environment Variable Configuration  
✅ JSONL Memory Store  
✅ Model Routing  
✅ Extended Thinking Support  
✅ Vision/Image Support  
✅ PICOCLAW_HOME (THIS)  

**All features complete!** Only optional search providers remaining.

---

**Date**: 2026-03-09  
**Status**: Complete  
**Build**: Passing  
**Integration**: 95% (9.5/10 features)
