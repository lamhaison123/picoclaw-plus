# Reference Documentation

API documentation, configuration references, and technical specifications.

## 📚 Contents

### v0.2.1 Features (NEW)
- **[JSONL_STORAGE.md](./JSONL_STORAGE.md)** - Crash-safe storage
  - File format
  - Migration guide
  - Performance
  - Troubleshooting

- **[SEARCH_PROVIDERS.md](./SEARCH_PROVIDERS.md)** - All 7 search providers
  - Provider comparison
  - Configuration
  - Use cases
  - API keys

- **[PICOCLAW_HOME.md](./PICOCLAW_HOME.md)** - Custom home directory
  - Multi-user setup
  - Docker deployment
  - Migration
  - Best practices

- **[MINDGRAPH_INTEGRATION.md](./MINDGRAPH_INTEGRATION.md)** - MindGraph memory provider
  - Knowledge graph memory
  - Configuration (self-hosted & cloud)
  - API usage
  - Best practices

- **[MEM0_INTEGRATION.md](./MEM0_INTEGRATION.md)** - Mem0 memory provider
  - Personalized memory
  - Configuration (self-hosted & cloud)
  - API usage
  - Best practices

### API Documentation
- **[API_REFERENCE.md](./API_REFERENCE.md)** - Complete API documentation
  - REST endpoints
  - Request/response formats
  - Authentication
  - Error codes
  - Examples

### Configuration
- **[tools_configuration.md](./tools_configuration.md)** - Tool system configuration
  - Tool registration
  - Permission settings
  - Workspace restrictions
  - Custom tools

### Security
- **[SAFETY_LEVELS.md](./SAFETY_LEVELS.md)** - Security levels explained
  - Level definitions (0-4)
  - Tool restrictions per level
  - Workspace isolation
  - Best practices

- **[ANTIGRAVITY_AUTH.md](./ANTIGRAVITY_AUTH.md)** - Authentication system
  - OAuth flow
  - Token management
  - API keys
  - Security considerations

### System Components
- **[CIRCUIT_BREAKER.md](./CIRCUIT_BREAKER.md)** - Fault tolerance system
  - Circuit states
  - Failure thresholds
  - Recovery strategies
  - Configuration

- **[TEAM_TOOL_ACCESS.md](./TEAM_TOOL_ACCESS.md)** - Team tool permissions
  - Role-based access
  - Tool allowlists
  - Capability mapping
  - Security model

- **[MULTI_AGENT_MODEL_SELECTION.md](./MULTI_AGENT_MODEL_SELECTION.md)** - Model routing
  - Model selection logic
  - Fallback chains
  - Cost optimization
  - Performance tuning

## 🎯 Quick Reference

### Configuration Files

```
config/
├── config.json              # Main configuration
├── config.memory.*.json     # Memory presets
└── agents/                  # Agent definitions
```

### Safety Levels

| Level | Description | Use Case |
|-------|-------------|----------|
| 0 | Unrestricted | Development only |
| 1 | Basic restrictions | Trusted environments |
| 2 | Moderate restrictions | Default production |
| 3 | Strict restrictions | Sensitive data |
| 4 | Maximum restrictions | Public/untrusted |

### API Endpoints

```
POST   /api/v1/chat          # Send message
GET    /api/v1/agents        # List agents
POST   /api/v1/teams         # Create team
GET    /api/v1/status        # System status
```

### Tool Categories

- **File Operations**: read, write, edit, append, list
- **Shell Execution**: execute, spawn
- **Web Access**: search, fetch
- **Team Coordination**: delegate, status
- **Memory**: search, store
- **MCP**: External tool protocols

## 📖 Usage Examples

### API Call
```bash
curl -X POST http://localhost:8080/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Hello, PicoClaw!",
    "agent_id": "default",
    "session_key": "user123"
  }'
```

### Tool Configuration
```json
{
  "tools": {
    "restrict_to_workspace": true,
    "allowed_commands": ["ls", "cat", "grep"],
    "denied_paths": ["/etc", "/sys"]
  }
}
```

### Safety Level
```json
{
  "agents": {
    "defaults": {
      "safety_level": 2
    }
  }
}
```

## 🔍 Finding Information

### I need to...

#### Configure tools
→ [tools_configuration.md](./tools_configuration.md)

#### Set security levels
→ [SAFETY_LEVELS.md](./SAFETY_LEVELS.md)

#### Use the API
→ [API_REFERENCE.md](./API_REFERENCE.md)

#### Configure memory providers
→ [MINDGRAPH_INTEGRATION.md](./MINDGRAPH_INTEGRATION.md)

#### Configure fault tolerance
→ [CIRCUIT_BREAKER.md](./CIRCUIT_BREAKER.md)

#### Set up team permissions
→ [TEAM_TOOL_ACCESS.md](./TEAM_TOOL_ACCESS.md)

#### Choose LLM models
→ [MULTI_AGENT_MODEL_SELECTION.md](./MULTI_AGENT_MODEL_SELECTION.md)

#### Authenticate with cloud
→ [ANTIGRAVITY_AUTH.md](./ANTIGRAVITY_AUTH.md)

## 📊 Configuration Schema

### Main Config Structure
```json
{
  "agents": {
    "defaults": {},
    "list": []
  },
  "channels": [],
  "providers": {},
  "memory": {},
  "tools": {},
  "heartbeat": {}
}
```

### Agent Definition
```json
{
  "id": "agent-id",
  "name": "Agent Name",
  "model": "gpt-4",
  "workspace": "./workspace",
  "safety_level": 2,
  "tools": [],
  "routing": []
}
```

### Memory Configuration
```json
{
  "enabled": true,
  "embedding": {
    "provider": "openai",
    "model": "text-embedding-3-small",
    "dimension": 1536
  },
  "vector_store": {
    "provider": "qdrant",
    "qdrant": {
      "enabled": true,
      "url": "http://localhost:6333"
    }
  },
  "memory_provider": {
    "provider": "mindgraph",
    "mindgraph": {
      "enabled": true,
      "url": "https://api.mindgraph.cloud",
      "api_key": "${MINDGRAPH_API_KEY}"
    }
  }
}
```

## 🔗 Related Documentation

- [Architecture Overview](../architecture/ARCHITECTURE_OVERVIEW.md) - System design
- [User Guides](../guides/) - How-to guides
- [Developer Guide](../development/DEVELOPER_GUIDE.md) - Development
- [Troubleshooting](../development/troubleshooting.md) - Common issues

## 📝 API Versioning

PicoClaw uses semantic versioning for the API:
- **v1**: Current stable API
- Breaking changes increment major version
- New features increment minor version
- Bug fixes increment patch version

## 🛠️ Configuration Validation

Validate your configuration:
```bash
picoclaw validate config.json
```

Check agent definitions:
```bash
picoclaw agent list
picoclaw agent validate <agent-id>
```

## 📞 Support

- **API Issues**: [GitHub Issues](https://github.com/sipeed/picoclaw/issues)
- **Configuration Help**: [Discussions](https://github.com/sipeed/picoclaw/discussions)
- **Security Concerns**: security@sipeed.com

---

**Last Updated**: 2026-03-09  
**API Version**: v1
