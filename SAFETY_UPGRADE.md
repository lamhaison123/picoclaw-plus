# Safety System Upgrade - LLM Omnipotence Mode

## 🎯 Mục tiêu

Nâng cấp hệ thống safety của PicoClaw để cho phép LLM toàn năng hơn, linh hoạt hơn trong việc thực thi lệnh, đồng thời vẫn đảm bảo an toàn cơ bản.

## 🚀 Các thay đổi chính

### 1. Safety Levels (4 cấp độ)

Thay vì chỉ có bật/tắt deny patterns, giờ có 4 cấp độ safety:

| Level | Mô tả | Use Case |
|-------|-------|----------|
| **strict** | Chặn hầu hết lệnh nguy hiểm | Production, môi trường không tin tưởng |
| **moderate** | Cân bằng giữa an toàn và linh hoạt (mặc định) | Development, môi trường tin tưởng |
| **permissive** | Chỉ chặn lệnh cực kỳ nguy hiểm | DevOps, system admin |
| **off** | Không có kiểm tra (nguy hiểm!) | Testing only |

### 2. Phân loại Patterns

**Catastrophic Patterns** (luôn chặn trừ khi off):
- `rm -rf /` - Xóa root filesystem
- `dd if=/dev/zero of=/dev/sda` - Wipe disk
- `:(){ :|:& };:` - Fork bomb
- `chmod -R 000 /` - Remove all permissions on root

**Strict Patterns** (chặn ở strict mode):
- Command substitution: `$(...)`, backticks
- Pipe to shell: `| sh`, `| bash`
- `sudo`, `chmod`, `chown`
- Package installation
- Process termination

**Moderate Patterns** (chặn ở strict và moderate):
- `rm -rf /path` trên root paths
- `dd` với dangerous sources
- Write to block devices
- Pipe curl/wget to shell

**Warning Patterns** (cảnh báo nhưng không chặn):
- `rm -rf` (non-root)
- `git push --force`
- `docker system prune -a`
- Package installation

### 3. Custom Patterns

**Custom Allow Patterns** - Bypass tất cả checks:
```json
{
  "custom_allow_patterns": [
    "\\bgit\\s+push\\s+--force\\b",
    "\\bsudo\\s+systemctl\\s+restart\\s+myapp\\b"
  ]
}
```

**Custom Deny Patterns** - Thêm restrictions:
```json
{
  "custom_deny_patterns": [
    "\\brm\\s+-rf\\s+/home",
    "\\bmv\\s+.*\\s+/dev/null"
  ]
}
```

## 📝 Files Changed

### Core Implementation
- `pkg/tools/shell.go` - Main safety logic
  - Added `SafetyLevel` enum
  - Refactored patterns into categories
  - Added warning system
  - Improved `guardCommand()` function

### Configuration
- `pkg/config/config.go` - Added `SafetyLevel` field
- `pkg/config/defaults.go` - Set default to "moderate"
- `pkg/migrate/sources/openclaw/openclaw_config.go` - Migration support

### Tests
- `pkg/tools/shell_safety_test.go` - New comprehensive tests
  - Test all 4 safety levels
  - Test custom patterns
  - Test backward compatibility

### Documentation
- `docs/SAFETY_LEVELS.md` - Complete documentation
- `docs/SAFETY_QUICKSTART.md` - Quick start guide
- `config/safety_examples.json` - Configuration examples
- `SAFETY_UPGRADE.md` - This file

## 🔧 Configuration Examples

### Example 1: Default (Moderate)
```json
{
  "tools": {
    "exec": {
      "safety_level": "moderate"
    }
  }
}
```

### Example 2: Production (Strict)
```json
{
  "tools": {
    "exec": {
      "safety_level": "strict",
      "custom_allow_patterns": [
        "\\bgit\\s+push\\s+origin\\s+main\\b"
      ]
    }
  }
}
```

### Example 3: DevOps (Permissive)
```json
{
  "tools": {
    "exec": {
      "safety_level": "permissive",
      "custom_deny_patterns": [
        "\\brm\\s+-rf\\s+/(bin|boot|etc|lib)\\b"
      ]
    }
  }
}
```

### Example 4: Testing (Off)
```json
{
  "tools": {
    "exec": {
      "safety_level": "off"
    }
  }
}
```

## 🔄 Backward Compatibility

Old config vẫn hoạt động:
```json
{
  "tools": {
    "exec": {
      "enable_deny_patterns": false
    }
  }
}
```

Sẽ tự động convert thành:
```json
{
  "tools": {
    "exec": {
      "safety_level": "off"
    }
  }
}
```

## 📊 So sánh Old vs New

### Old System
```
✅ Bật deny patterns (chặn tất cả)
❌ Tắt deny patterns (cho phép tất cả)
```

### New System
```
🔴 Strict    - Chặn nhiều nhất
🟡 Moderate  - Cân bằng (default)
🟢 Permissive - Chặn ít nhất
⚪ Off       - Không chặn
```

## 🎯 Benefits

### 1. Linh hoạt hơn
- LLM có thể thực hiện nhiều tác vụ hơn
- Phù hợp với nhiều use case khác nhau
- Dễ dàng tùy chỉnh theo nhu cầu

### 2. An toàn hơn
- Luôn chặn catastrophic commands (trừ khi off)
- Warning system cho các lệnh nguy hiểm
- Custom patterns cho fine-grained control

### 3. Dễ sử dụng hơn
- Clear documentation
- Many examples
- Quick start guide
- Backward compatible

## 🧪 Testing

```bash
# Test all safety levels
go test ./pkg/tools -run TestSafetyLevel -v

# Test specific level
go test ./pkg/tools -run TestSafetyLevel_Strict -v
go test ./pkg/tools -run TestSafetyLevel_Moderate -v
go test ./pkg/tools -run TestSafetyLevel_Permissive -v
go test ./pkg/tools -run TestSafetyLevel_Off -v

# Test custom patterns
go test ./pkg/tools -run TestSafetyLevel_CustomAllowPatterns -v
go test ./pkg/tools -run TestSafetyLevel_CustomDenyPatterns -v
```

## 📚 Documentation

1. **Quick Start**: `docs/SAFETY_QUICKSTART.md`
   - 5-minute setup
   - Common use cases
   - Troubleshooting

2. **Complete Guide**: `docs/SAFETY_LEVELS.md`
   - Detailed explanation of each level
   - All patterns documented
   - Security best practices
   - Migration guide

3. **Examples**: `config/safety_examples.json`
   - 10 real-world configuration examples
   - Different use cases
   - Custom patterns examples

## 🚨 Security Considerations

### Always Blocked (except off mode)
- Root filesystem deletion
- Disk wiping
- Fork bombs
- Force shutdown/reboot

### Configurable
- Command substitution
- Pipe to shell
- Privilege escalation
- Package installation
- Process termination

### Never Blocked
- Read operations
- Safe writes within workspace
- Development tools
- Git operations (except force push in strict)

## 🎓 Usage Recommendations

| Scenario | Recommended Level | Notes |
|----------|------------------|-------|
| Production API | strict | Maximum protection |
| Development | moderate | Good balance |
| CI/CD | strict + custom allow | Explicit whitelist |
| DevOps | permissive | System admin tasks |
| Testing | off | Only for testing safety system |

## 🔮 Future Enhancements

Potential improvements:
1. Interactive confirmation for dangerous commands
2. Audit logging of all executed commands
3. Rate limiting for certain command types
4. Context-aware safety (different levels per workspace)
5. Machine learning-based anomaly detection

## 📞 Support

- Documentation: `docs/SAFETY_LEVELS.md`
- Examples: `config/safety_examples.json`
- Tests: `pkg/tools/shell_safety_test.go`
- Issues: GitHub Issues

## ✅ Checklist

- [x] Implement SafetyLevel enum
- [x] Refactor patterns into categories
- [x] Add warning system
- [x] Update configuration structs
- [x] Add backward compatibility
- [x] Write comprehensive tests
- [x] Create documentation
- [x] Create examples
- [x] Update migration code

## 🎉 Summary

Hệ thống safety mới cho phép LLM toàn năng hơn với 4 cấp độ linh hoạt, từ strict (production) đến off (testing), đồng thời vẫn đảm bảo an toàn cơ bản bằng cách luôn chặn các lệnh catastrophic. Custom patterns cho phép fine-grained control theo nhu cầu cụ thể.

**Default: moderate** - Cân bằng tốt giữa an toàn và linh hoạt cho hầu hết use cases!
