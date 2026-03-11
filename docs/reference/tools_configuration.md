# Tools Configuration

PicoClaw's tools configuration is located in the `tools` field of `config.json`.

## Directory Structure

```json
{
  "tools": {
    "web": { ... },
    "mcp": { ... },
    "exec": { ... },
    "cron": { ... },
    "skills": { ... }
  }
}
```

## Web Tools

Web tools are used for web search and fetching.

### Brave

| Config        | Type   | Default | Description               |
| ------------- | ------ | ------- | ------------------------- |
| `enabled`     | bool   | false   | Enable Brave search       |
| `api_key`     | string | -       | Brave Search API key      |
| `max_results` | int    | 5       | Maximum number of results |

### DuckDuckGo

| Config        | Type | Default | Description               |
| ------------- | ---- | ------- | ------------------------- |
| `enabled`     | bool | true    | Enable DuckDuckGo search  |
| `max_results` | int  | 5       | Maximum number of results |

### Perplexity

| Config        | Type   | Default | Description               |
| ------------- | ------ | ------- | ------------------------- |
| `enabled`     | bool   | false   | Enable Perplexity search  |
| `api_key`     | string | -       | Perplexity API key        |
| `max_results` | int    | 5       | Maximum number of results |

## Exec Tool

The exec tool is used to execute shell commands with a 4-level safety system.

| Config                    | Type   | Default      | Description                                      |
| ------------------------- | ------ | ------------ | ------------------------------------------------ |
| `safety_level`            | string | `"moderate"` | Safety level: `strict`, `moderate`, `permissive`, `off` |
| `custom_allow_patterns`   | array  | []           | Custom allow patterns (bypass safety checks)     |
| `custom_deny_patterns`    | array  | []           | Custom deny patterns (additional restrictions)   |
| `enable_deny_patterns`    | bool   | true         | (Deprecated) Use `safety_level` instead          |

### Safety Levels

PicoClaw provides 4 safety levels for command execution:

#### 1. Strict (`"strict"`)
**Use Case:** Production environments, maximum security

- Blocks ALL potentially dangerous commands
- Prevents command substitution, pipe to shell, privilege escalation
- Blocks system modifications, package installation, process termination
- Allows: Read operations, safe writes, development tools, git (except force push)

#### 2. Moderate (`"moderate"`) - **DEFAULT**
**Use Case:** Development environments, balanced protection

- Blocks CATASTROPHIC commands only (root deletion, disk wipe, fork bombs)
- Generates warnings for risky operations (but allows them)
- Allows: Most development operations, package installation, docker, git force push (with warning)

#### 3. Permissive (`"permissive"`)
**Use Case:** System administration, maximum autonomy

- Blocks ONLY catastrophic commands
- Minimal interference, all warnings disabled
- Allows: Nearly all commands including sudo, system admin tasks

#### 4. Off (`"off"`)
**Use Case:** Testing only (⚠️ DANGEROUS)

- NO SAFETY CHECKS - all commands allowed
- Only use if you fully trust the LLM

### Configuration Examples

#### Basic Configuration
```json
{
  "tools": {
    "exec": {
      "safety_level": "moderate"
    }
  }
}
```

#### With Custom Patterns
```json
{
  "tools": {
    "exec": {
      "safety_level": "strict",
      "custom_allow_patterns": [
        "\\bgit\\s+push\\s+--force\\b"
      ],
      "custom_deny_patterns": [
        "\\brm\\s+-rf\\s+/home"
      ]
    }
  }
}
```

#### Legacy Configuration (Deprecated)
```json
{
  "tools": {
    "exec": {
      "enable_deny_patterns": false  // Automatically converted to safety_level: "off"
    }
  }
}
```

### Comparison Table

| Feature | Strict | Moderate | Permissive | Off |
|---------|--------|----------|------------|-----|
| Root deletion | ❌ | ❌ | ❌ | ✅ |
| Disk wipe | ❌ | ❌ | ❌ | ✅ |
| Fork bomb | ❌ | ❌ | ❌ | ✅ |
| Command substitution | ❌ | ✅ | ✅ | ✅ |
| `sudo` commands | ❌ | ✅ | ✅ | ✅ |
| Package install | ❌ | ⚠️ | ✅ | ✅ |
| Git force push | ❌ | ⚠️ | ✅ | ✅ |
| Docker commands | ❌ | ✅ | ✅ | ✅ |

Legend: ❌ Blocked | ⚠️ Warning (allowed) | ✅ Allowed

📖 **Full Documentation:** See [SAFETY_LEVELS.md](SAFETY_LEVELS.md) for complete details

## Cron Tool

The cron tool is used for scheduling periodic tasks.

| Config                 | Type | Default | Description                                    |
| ---------------------- | ---- | ------- | ---------------------------------------------- |
| `exec_timeout_minutes` | int  | 5       | Execution timeout in minutes, 0 means no limit |

## MCP Tool

The MCP tool enables integration with external Model Context Protocol servers.

### Global Config

| Config    | Type   | Default | Description                         |
| --------- | ------ | ------- | ----------------------------------- |
| `enabled` | bool   | false   | Enable MCP integration globally     |
| `servers` | object | `{}`    | Map of server name to server config |

### Per-Server Config

| Config     | Type   | Required | Description                                |
| ---------- | ------ | -------- | ------------------------------------------ |
| `enabled`  | bool   | yes      | Enable this MCP server                     |
| `type`     | string | no       | Transport type: `stdio`, `sse`, `http`     |
| `command`  | string | stdio    | Executable command for stdio transport     |
| `args`     | array  | no       | Command arguments for stdio transport      |
| `env`      | object | no       | Environment variables for stdio process    |
| `env_file` | string | no       | Path to environment file for stdio process |
| `url`      | string | sse/http | Endpoint URL for `sse`/`http` transport    |
| `headers`  | object | no       | HTTP headers for `sse`/`http` transport    |

### Transport Behavior

- If `type` is omitted, transport is auto-detected:
  - `url` is set → `sse`
  - `command` is set → `stdio`
- `http` and `sse` both use `url` + optional `headers`.
- `env` and `env_file` are only applied to `stdio` servers.

### Configuration Examples

#### 1) Stdio MCP server

```json
{
  "tools": {
    "mcp": {
      "enabled": true,
      "servers": {
        "filesystem": {
          "enabled": true,
          "command": "npx",
          "args": ["-y", "@modelcontextprotocol/server-filesystem", "/tmp"]
        }
      }
    }
  }
}
```

#### 2) Remote SSE/HTTP MCP server

```json
{
  "tools": {
    "mcp": {
      "enabled": true,
      "servers": {
        "remote-mcp": {
          "enabled": true,
          "type": "sse",
          "url": "https://example.com/mcp",
          "headers": {
            "Authorization": "Bearer YOUR_TOKEN"
          }
        }
      }
    }
  }
}
```

## Skills Tool

The skills tool configures skill discovery and installation via registries like ClawHub.

### Registries

| Config                             | Type   | Default              | Description             |
| ---------------------------------- | ------ | -------------------- | ----------------------- |
| `registries.clawhub.enabled`       | bool   | true                 | Enable ClawHub registry |
| `registries.clawhub.base_url`      | string | `https://clawhub.ai` | ClawHub base URL        |
| `registries.clawhub.search_path`   | string | `/api/v1/search`     | Search API path         |
| `registries.clawhub.skills_path`   | string | `/api/v1/skills`     | Skills API path         |
| `registries.clawhub.download_path` | string | `/api/v1/download`   | Download API path       |

### Configuration Example

```json
{
  "tools": {
    "skills": {
      "registries": {
        "clawhub": {
          "enabled": true,
          "base_url": "https://clawhub.ai",
          "search_path": "/api/v1/search",
          "skills_path": "/api/v1/skills",
          "download_path": "/api/v1/download"
        }
      }
    }
  }
}
```

## Environment Variables

All configuration options can be overridden via environment variables with the format `PICOCLAW_TOOLS_<SECTION>_<KEY>`:

For example:

- `PICOCLAW_TOOLS_WEB_BRAVE_ENABLED=true`
- `PICOCLAW_TOOLS_EXEC_ENABLE_DENY_PATTERNS=false`
- `PICOCLAW_TOOLS_CRON_EXEC_TIMEOUT_MINUTES=10`
- `PICOCLAW_TOOLS_MCP_ENABLED=true`

Note: Nested map-style config (for example `tools.mcp.servers.<name>.*`) is configured in `config.json` rather than environment variables.
