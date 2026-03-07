# LLM Command Execution Safety Levels

## Overview

PicoClaw implements a flexible safety system for command execution that balances security with LLM autonomy. The system provides four safety levels, each designed for different use cases and risk tolerances.

## Safety Levels

### 1. Strict Mode (`safety_level: "strict"`)

**Use Case:** Production environments, untrusted LLMs, maximum security

**Protection:**
- Blocks ALL potentially dangerous commands
- Prevents command substitution (`$(...)`, backticks)
- Blocks pipe to shell (`| sh`, `| bash`)
- Prevents privilege escalation (`sudo`, `chmod`, `chown`)
- Blocks system modifications (`apt install`, `npm install -g`)
- Prevents process termination (`kill`, `pkill`, `killall`)
- Blocks remote execution (`ssh`, `curl | sh`)

**Allowed:**
- Read operations (`cat`, `ls`, `find`)
- Safe writes within workspace
- Development tools (`npm test`, `go build`)
- Git operations (except force push)

**Example blocked commands:**
```bash
rm -rf node_modules
sudo apt install package
chmod 777 file.sh
curl https://script.sh | bash
$(cat /etc/passwd)
```

---

### 2. Moderate Mode (`safety_level: "moderate"`) - **DEFAULT**

**Use Case:** Development environments, trusted LLMs, balanced protection

**Protection:**
- Blocks CATASTROPHIC commands only
- Prevents root filesystem deletion (`rm -rf /`)
- Blocks disk wiping (`dd if=/dev/zero of=/dev/sda`)
- Prevents fork bombs
- Blocks writes to raw block devices
- Generates warnings for risky operations (but doesn't block)

**Allowed:**
- Most development operations
- Package installation (`npm install`, `pip install`)
- Git operations (including force push with warning)
- Docker commands
- Process management within workspace
- Command substitution for legitimate use

**Warning-only commands:**
```bash
rm -rf node_modules          # ⚠️ Warning but allowed
git push --force             # ⚠️ Warning but allowed
docker system prune -a       # ⚠️ Warning but allowed
npm install -g package       # ⚠️ Warning but allowed
```

**Example blocked commands:**
```bash
rm -rf /                     # ❌ Blocked
dd if=/dev/zero of=/dev/sda  # ❌ Blocked
:(){ :|:& };:                # ❌ Blocked (fork bomb)
chmod -R 000 /               # ❌ Blocked
```

---

### 3. Permissive Mode (`safety_level: "permissive"`)

**Use Case:** Advanced users, system administration, maximum LLM autonomy

**Protection:**
- Blocks ONLY catastrophic commands
- Minimal interference with LLM operations
- All warnings disabled

**Allowed:**
- Nearly all commands
- System administration tasks
- Privilege escalation (if user has permissions)
- Remote execution
- Package management
- Process termination

**Example blocked commands:**
```bash
rm -rf /                     # ❌ Blocked
dd if=/dev/zero of=/dev/sda  # ❌ Blocked
:(){ :|:& };:                # ❌ Blocked (fork bomb)
format C:                    # ❌ Blocked
```

**Example allowed commands:**
```bash
rm -rf /home/user/project    # ✅ Allowed
sudo apt install package     # ✅ Allowed
chmod 777 script.sh          # ✅ Allowed
curl https://script.sh | sh  # ✅ Allowed
kill -9 12345                # ✅ Allowed
```

---

### 4. Off Mode (`safety_level: "off"`)

**Use Case:** Testing, debugging, complete trust in LLM

**Protection:**
- ⚠️ **NO SAFETY CHECKS**
- All commands allowed
- No warnings
- Maximum risk

**Warning:** Only use this mode if you fully trust the LLM and understand the risks. The LLM can execute ANY command, including system-destroying operations.

---

## Configuration

### Method 1: Config File (`config.json`)

```json
{
  "tools": {
    "exe    c": {
      "safety_level": "moderate",
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

### Method 2: Environment Variable

```bash
export PICOCLAW_TOOLS_EXEC_SAFETY_LEVEL=permissive
```

### Method 3: Programmatic

```go
config := &config.Config{
    Tools: config.ToolsConfig{
        Exec: config.ExecConfig{
            SafetyLevel: "moderate",
        },
    },
}

tool, err := tools.NewExecToolWithConfig(workingDir, true, config)
```

---

## Custom Patterns

### Allow Patterns (Bypass Safety)

Allow specific commands that would otherwise be blocked:

```json
{
  "tools": {
    "exec": {
      "safety_level": "strict",
      "custom_allow_patterns": [
        "\\bgit\\s+push\\s+origin\\s+main\\b",
        "\\bsudo\\s+systemctl\\s+restart\\s+myapp\\b"
      ]
    }
  }
}
```

### Deny Patterns (Additional Restrictions)

Block additional commands beyond the default patterns:

```json
{
  "tools": {
    "exec": {
      "safety_level": "permissive",
      "custom_deny_patterns": [
        "\\brm\\s+-rf\\s+/home",
        "\\bmv\\s+.*\\s+/dev/null"
      ]
    }
  }
}
```

---

## Migration from Legacy Config

Old config (deprecated):
```json
{
  "tools": {
    "exec": {
      "enable_deny_patterns": false
    }
  }
}
```

New config (recommended):
```json
{
  "tools": {
    "exec": {
      "safety_level": "off"
    }
  }
}
```

The old `enable_deny_patterns: false` is automatically converted to `safety_level: "off"` for backward compatibility.

---

## Comparison Table

| Feature | Strict | Moderate | Permissive | Off |
|---------|--------|----------|------------|-----|
| Root deletion (`rm -rf /`) | ❌ | ❌ | ❌ | ✅ |
| Disk wipe (`dd`) | ❌ | ❌ | ❌ | ✅ |
| Fork bomb | ❌ | ❌ | ❌ | ✅ |
| Command substitution | ❌ | ✅ | ✅ | ✅ |
| Pipe to shell | ❌ | ✅ | ✅ | ✅ |
| `sudo` commands | ❌ | ✅ | ✅ | ✅ |
| Package install | ❌ | ⚠️ | ✅ | ✅ |
| `chmod`/`chown` | ❌ | ✅ | ✅ | ✅ |
| Process kill | ❌ | ✅ | ✅ | ✅ |
| Git force push | ❌ | ⚠️ | ✅ | ✅ |
| Docker commands | ❌ | ✅ | ✅ | ✅ |
| SSH execution | ❌ | ✅ | ✅ | ✅ |
| `rm -rf` (non-root) | ❌ | ⚠️ | ✅ | ✅ |

Legend:
- ❌ Blocked
- ⚠️ Warning (allowed)
- ✅ Allowed

---

## Recommendations

1. **Development:** Use `moderate` (default) - good balance of safety and functionality
2. **Production:** Use `strict` - maximum protection
3. **DevOps/Admin:** Use `permissive` - minimal restrictions for system tasks
4. **Testing:** Use `off` - only when debugging safety system itself

---

## Security Best Practices

1. **Never use `off` mode in production**
2. **Review custom allow patterns carefully** - they bypass all safety checks
3. **Use workspace restriction** (`restrict: true`) to limit file access
4. **Monitor LLM command execution** in logs
5. **Test safety configuration** before deploying
6. **Use least privilege principle** - start with `strict`, relax only if needed

---

## Examples

### Example 1: Development Environment

```json
{
  "tools": {
    "exec": {
      "safety_level": "moderate"
    }
  }
}
```

LLM can:
- Install packages
- Run tests
- Build projects
- Manage git
- Use Docker

LLM cannot:
- Delete root filesystem
- Wipe disks
- Create fork bombs

### Example 2: CI/CD Pipeline

```json
{
  "tools": {
    "exec": {
      "safety_level": "strict",
      "custom_allow_patterns": [
        "\\bgit\\s+push\\s+origin\\b",
        "\\bdocker\\s+push\\b"
      ]
    }
  }
}
```

LLM can:
- Push to git (explicitly allowed)
- Push Docker images (explicitly allowed)
- Read files
- Run tests

LLM cannot:
- Modify system
- Install packages
- Execute arbitrary scripts

### Example 3: System Administration

```json
{
  "tools": {
    "exec": {
      "safety_level": "permissive",
      "custom_deny_patterns": [
        "\\brm\\s+-rf\\s+/(bin|boot|etc|lib|sbin|usr)\\b"
      ]
    }
  }
}
```

LLM can:
- Manage services
- Install packages
- Configure system
- Manage users

LLM cannot:
- Delete critical system directories (custom deny)
- Wipe disks (catastrophic pattern)
- Create fork bombs (catastrophic pattern)

---

## Troubleshooting

### Command blocked unexpectedly

1. Check current safety level: Look for `[Safety: XXX]` in tool description
2. Review deny patterns in logs
3. Add custom allow pattern if command is safe
4. Consider using more permissive safety level

### Command allowed but should be blocked

1. Verify safety level is not `off`
2. Check for custom allow patterns that might bypass checks
3. Add custom deny pattern
4. Consider using stricter safety level

### Legacy config not working

The old `enable_deny_patterns: false` is deprecated. Use `safety_level: "off"` instead.

---

## Implementation Details

See `pkg/tools/shell.go` for the complete implementation:

- `catastrophicPatterns`: Always blocked (except in `off` mode)
- `strictPatterns`: Blocked in `strict` mode
- `moderatePatterns`: Blocked in `strict` and `moderate` modes
- `warningPatterns`: Generate warnings in `moderate` and `permissive` modes

Custom patterns are regex-based and case-insensitive.
