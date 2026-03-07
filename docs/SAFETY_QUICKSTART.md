# Safety Levels - Quick Start Guide

## TL;DR

PicoClaw now has 4 safety levels for LLM command execution:

```json
{
  "tools": {
    "exec": {
      "safety_level": "moderate"
    }
  }
}
```

**Options:** `"strict"` | `"moderate"` (default) | `"permissive"` | `"off"`

---

## Quick Comparison

| Level | Best For | LLM Can Do | LLM Cannot Do |
|-------|----------|------------|---------------|
| **strict** | Production | Read files, run tests, build code | Install packages, modify system, use sudo |
| **moderate** | Development | Everything in strict + install packages, docker, git force push | Delete root filesystem, wipe disks |
| **permissive** | DevOps/Admin | Almost everything | Only catastrophic operations |
| **off** | Testing only | Everything (no restrictions) | Nothing blocked ⚠️ |

---

## 5-Minute Setup

### 1. Choose Your Safety Level

**For most developers (recommended):**
```json
{
  "tools": {
    "exec": {
      "safety_level": "moderate"
    }
  }
}
```

**For production environments:**
```json
{
  "tools": {
    "exec": {
      "safety_level": "strict"
    }
  }
}
```

**For system administration:**
```json
{
  "tools": {
    "exec": {
      "safety_level": "permissive"
    }
  }
}
```

### 2. Add Custom Rules (Optional)

**Allow specific commands in strict mode:**
```json
{
  "tools": {
    "exec": {
      "safety_level": "strict",
      "custom_allow_patterns": [
        "\\bgit\\s+push\\s+--force\\b",
        "\\bdocker\\s+run\\b"
      ]
    }
  }
}
```

**Block additional commands in permissive mode:**
```json
{
  "tools": {
    "exec": {
      "safety_level": "permissive",
      "custom_deny_patterns": [
        "\\brm\\s+-rf\\s+/home",
        "\\bchmod\\s+-R\\s+000"
      ]
    }
  }
}
```

### 3. Test Your Configuration

```bash
# Set environment variable
export PICOCLAW_TOOLS_EXEC_SAFETY_LEVEL=moderate

# Or use config file
cat > config.json << EOF
{
  "tools": {
    "exec": {
      "safety_level": "moderate"
    }
  }
}
EOF

# Run tests
go test ./pkg/tools -run TestSafetyLevel -v
```

---

## Common Use Cases

### Use Case 1: Web Development

```json
{
  "tools": {
    "exec": {
      "safety_level": "moderate"
    }
  }
}
```

✅ LLM can:
- `npm install` / `yarn add`
- `npm run build` / `npm test`
- `git commit` / `git push`
- `docker-compose up`

❌ LLM cannot:
- `rm -rf /`
- `dd if=/dev/zero of=/dev/sda`

### Use Case 2: CI/CD Pipeline

```json
{
  "tools": {
    "exec": {
      "safety_level": "strict",
      "custom_allow_patterns": [
        "\\bgit\\s+push\\s+origin\\b",
        "\\bdocker\\s+push\\b",
        "\\bnpm\\s+publish\\b"
      ]
    }
  }
}
```

✅ LLM can:
- `git push origin main` (explicitly allowed)
- `docker push myimage` (explicitly allowed)
- `npm publish` (explicitly allowed)
- `npm test` / `go build`

❌ LLM cannot:
- `sudo apt install`
- `chmod 777`
- Any command not explicitly allowed

### Use Case 3: Server Administration

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

✅ LLM can:
- `sudo systemctl restart nginx`
- `apt install package`
- `chmod 755 script.sh`
- `kill -9 12345`

❌ LLM cannot:
- `rm -rf /bin` (custom deny)
- `rm -rf /` (catastrophic)
- `dd if=/dev/zero of=/dev/sda` (catastrophic)

---

## Migration Guide

### From Old Config

**Before (deprecated):**
```json
{
  "tools": {
    "exec": {
      "enable_deny_patterns": false
    }
  }
}
```

**After (recommended):**
```json
{
  "tools": {
    "exec": {
      "safety_level": "off"
    }
  }
}
```

The old config still works but will show a deprecation warning.

---

## Troubleshooting

### Problem: Command blocked unexpectedly

**Solution 1:** Check safety level
```bash
# Look for [Safety: XXX] in tool description
```

**Solution 2:** Add custom allow pattern
```json
{
  "tools": {
    "exec": {
      "safety_level": "strict",
      "custom_allow_patterns": [
        "\\byour-command-here\\b"
      ]
    }
  }
}
```

**Solution 3:** Use more permissive level
```json
{
  "tools": {
    "exec": {
      "safety_level": "moderate"
    }
  }
}
```

### Problem: Dangerous command allowed

**Solution 1:** Use stricter level
```json
{
  "tools": {
    "exec": {
      "safety_level": "strict"
    }
  }
}
```

**Solution 2:** Add custom deny pattern
```json
{
  "tools": {
    "exec": {
      "safety_level": "moderate",
      "custom_deny_patterns": [
        "\\bdangerous-command\\b"
      ]
    }
  }
}
```

---

## Testing Your Configuration

```bash
# Test strict mode
cat > test_strict.json << EOF
{
  "tools": {
    "exec": {
      "safety_level": "strict"
    }
  }
}
EOF

# Test moderate mode
cat > test_moderate.json << EOF
{
  "tools": {
    "exec": {
      "safety_level": "moderate"
    }
  }
}
EOF

# Run tests
go test ./pkg/tools -run TestSafetyLevel_Strict -v
go test ./pkg/tools -run TestSafetyLevel_Moderate -v
```

---

## Best Practices

1. ✅ **Start with `moderate`** - good default for most use cases
2. ✅ **Use `strict` in production** - maximum protection
3. ✅ **Test custom patterns** - ensure they work as expected
4. ✅ **Review LLM commands** - monitor what the LLM executes
5. ❌ **Never use `off` in production** - extremely dangerous
6. ❌ **Don't over-allow** - be specific with custom allow patterns

---

## Need More Details?

See [SAFETY_LEVELS.md](./SAFETY_LEVELS.md) for complete documentation.

---

## Examples

See [config/safety_examples.json](../config/safety_examples.json) for more configuration examples.

---

## Support

- GitHub Issues: Report bugs or request features
- Documentation: See `docs/SAFETY_LEVELS.md`
- Tests: See `pkg/tools/shell_safety_test.go`
