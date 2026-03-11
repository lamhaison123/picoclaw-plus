# PICOCLAW_HOME Reference

Custom home directory configuration for multi-user and container deployments.

## Overview

PICOCLAW_HOME allows you to specify a custom directory for all PicoClaw data, enabling multi-user setups and container-friendly deployments.

## Benefits

### Multi-User Support
- Each user has isolated data
- No conflicts between users
- Separate configurations per user

### Multi-Tenant
- Multiple PicoClaw instances
- Isolated workspaces
- Independent configurations

### Container-Friendly
- Mount custom directory
- Persistent data outside container
- Easy backup and restore

## Configuration

### Environment Variable
```bash
export PICOCLAW_HOME=/custom/path
```

### .env File
```bash
PICOCLAW_HOME=/custom/path
```

### Default Behavior
If not set, uses `~/.picoclaw`

## Directory Structure

### Default (~/.picoclaw)
```
~/.picoclaw/
├── config.json          # Configuration
├── auth.json            # Credentials
├── .env                 # Environment variables
├── workspace/           # Default workspace
├── sessions/            # Session files
├── teams/               # Team states
└── skills/              # Global skills
```

### Custom ($PICOCLAW_HOME)
```
$PICOCLAW_HOME/
├── config.json
├── auth.json
├── .env
├── workspace/
├── sessions/
├── teams/
└── skills/
```

## Use Cases

### Multi-User System
```bash
# User 1
export PICOCLAW_HOME=/home/user1/.picoclaw
picoclaw agent

# User 2
export PICOCLAW_HOME=/home/user2/.picoclaw
picoclaw agent
```

### Docker Container
```dockerfile
FROM golang:1.21
ENV PICOCLAW_HOME=/data/picoclaw
VOLUME /data/picoclaw
```

```bash
docker run -v /host/picoclaw:/data/picoclaw picoclaw
```

### Multiple Instances
```bash
# Development
export PICOCLAW_HOME=/opt/picoclaw-dev
picoclaw agent

# Production
export PICOCLAW_HOME=/opt/picoclaw-prod
picoclaw agent
```

### Testing
```bash
# Temporary test environment
export PICOCLAW_HOME=/tmp/picoclaw-test
picoclaw agent
```

## Components Using PICOCLAW_HOME

### Configuration
- `config.json` location
- `.env` file location
- Default paths

### Authentication
- `auth.json` location
- Credentials storage

### Sessions
- Session files (`.jsonl`)
- Session history

### Teams
- Team state files
- Team configurations

### Workspaces
- Default workspace path
- Agent workspaces

### Skills
- Global skills directory
- Skill installations

## Migration

### Moving Existing Data
```bash
# Copy data to new location
cp -r ~/.picoclaw /custom/path

# Set PICOCLAW_HOME
export PICOCLAW_HOME=/custom/path

# Verify
picoclaw status
```

### Backup
```bash
# Backup current data
tar -czf picoclaw-backup.tar.gz ~/.picoclaw

# Restore to new location
mkdir -p /custom/path
tar -xzf picoclaw-backup.tar.gz -C /custom/path
export PICOCLAW_HOME=/custom/path
```

## Priority

Configuration loading priority:
```
1. PICOCLAW_HOME (if set)
2. ~/.picoclaw (default)
3. Current directory (fallback)
```

## Best Practices

### Production
- Use absolute paths
- Set in systemd service
- Ensure proper permissions
- Regular backups

### Development
- Use separate directory
- Easy to clean/reset
- Isolated from production

### Docker
- Use volume mounts
- Persistent storage
- Easy backup/restore

## Systemd Integration

### Service File
```ini
[Unit]
Description=PicoClaw Agent
After=network.target

[Service]
Type=simple
User=picoclaw
Environment="PICOCLAW_HOME=/opt/picoclaw"
ExecStart=/usr/local/bin/picoclaw agent
Restart=always

[Install]
WantedBy=multi-user.target
```

## Permissions

### Recommended
```bash
# Create directory
mkdir -p /opt/picoclaw

# Set ownership
chown -R picoclaw:picoclaw /opt/picoclaw

# Set permissions
chmod 700 /opt/picoclaw
chmod 600 /opt/picoclaw/auth.json
chmod 600 /opt/picoclaw/.env
```

## Troubleshooting

### Config Not Found
Check:
```bash
echo $PICOCLAW_HOME
ls -la $PICOCLAW_HOME
```

### Permission Denied
Fix:
```bash
chmod 700 $PICOCLAW_HOME
chown -R $USER:$USER $PICOCLAW_HOME
```

### Sessions Not Loading
Verify:
```bash
ls -la $PICOCLAW_HOME/sessions/
```

## Environment Variables

### Related Variables
```bash
# Home directory
PICOCLAW_HOME=/custom/path

# Workspace (relative to PICOCLAW_HOME)
PICOCLAW_AGENTS_DEFAULTS_WORKSPACE=./workspace

# Or absolute path
PICOCLAW_AGENTS_DEFAULTS_WORKSPACE=/data/workspace
```

## API Reference

### Path Resolution
```go
func GetPicoClawHome() string {
    if home := os.Getenv("PICOCLAW_HOME"); home != "" {
        return home
    }
    return filepath.Join(os.Getenv("HOME"), ".picoclaw")
}
```

### Usage
```go
home := config.GetPicoClawHome()
configPath := filepath.Join(home, "config.json")
```

## Examples

### Development Setup
```bash
# Create dev environment
mkdir -p ~/picoclaw-dev
export PICOCLAW_HOME=~/picoclaw-dev

# Initialize
picoclaw agent
```

### Production Setup
```bash
# Create production directory
sudo mkdir -p /opt/picoclaw
sudo chown picoclaw:picoclaw /opt/picoclaw

# Configure systemd
sudo systemctl edit picoclaw
# Add: Environment="PICOCLAW_HOME=/opt/picoclaw"

# Start service
sudo systemctl start picoclaw
```

### Docker Compose
```yaml
version: '3'
services:
  picoclaw:
    image: picoclaw:latest
    environment:
      - PICOCLAW_HOME=/data
    volumes:
      - ./picoclaw-data:/data
```

## See Also

- [Configuration Guide](CONFIGURATION.md)
- [Installation Guide](../guides/INSTALLATION.md)
- [v0.2.1 Features](../guides/V0.2.1_FEATURES.md)

---

**Version**: v0.2.1  
**Last Updated**: 2026-03-09
