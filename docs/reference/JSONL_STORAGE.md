# JSONL Storage Reference

Crash-safe append-only session storage using JSONL format.

## Overview

JSONL (JSON Lines) storage provides crash-safe session persistence by using an append-only file format with automatic fsync.

## Benefits

### Crash Safety
- Automatic fsync after each write
- No data loss on power failure
- Atomic append operations
- Per-session file locking

### Performance
- Append-only writes (no file rewrites)
- Logical truncation (no physical compaction)
- Efficient for large sessions
- Memory-efficient loading

### Compatibility
- Automatic migration from JSON
- Backward compatible
- Old JSON files preserved

## Configuration

### Enable JSONL Storage
```json
{
  "session": {
    "storage_backend": "jsonl",
    "auto_migrate": true
  }
}
```

### Environment Variables
```bash
PICOCLAW_SESSION_STORAGE_BACKEND=jsonl
PICOCLAW_SESSION_AUTO_MIGRATE=true
```

## File Format

### JSONL Structure
Each line is a complete JSON object:
```jsonl
{"type":"message","role":"user","content":"Hello"}
{"type":"message","role":"assistant","content":"Hi!"}
{"type":"truncate","index":0}
```

### Message Types
- `message`: Chat message
- `truncate`: Logical truncation marker

## Migration

### Automatic Migration
Sessions automatically migrate from JSON to JSONL on first access.

**Process**:
1. Load old JSON file
2. Write to new JSONL file
3. Preserve old JSON as backup
4. Use JSONL for future writes

### Manual Migration
Not needed - automatic migration handles everything.

## File Locations

### Default
```
~/.picoclaw/sessions/
├── session-123.jsonl
├── session-456.jsonl
└── old-session.json (backup)
```

### Custom (PICOCLAW_HOME)
```
$PICOCLAW_HOME/sessions/
├── session-123.jsonl
└── session-456.jsonl
```

## Operations

### Append Message
```go
store.Append(sessionID, message)
// Automatically fsynced
```

### Load Session
```go
messages := store.Load(sessionID)
// Handles both JSON and JSONL
```

### Truncate Session
```go
store.Truncate(sessionID, index)
// Logical truncation (no file rewrite)
```

## Locking

### Per-Session Locks
- Uses sync.Map for concurrent access
- One lock per session
- Prevents race conditions
- Automatic cleanup

## Error Handling

### Corrupted Files
- Skips invalid lines
- Continues loading valid data
- Logs errors for debugging

### Disk Full
- Returns error immediately
- No partial writes
- Session remains consistent

## Performance

### Benchmarks
- Append: <10ms per message
- Load: <100ms for 1000 messages
- Fsync: <5ms per write

### Optimization Tips
- Use SSD for better fsync performance
- Increase summarization threshold
- Clean old sessions periodically

## Troubleshooting

### Session Not Loading
Check file permissions:
```bash
ls -la ~/.picoclaw/sessions/
```

### Slow Performance
Check disk I/O:
```bash
iostat -x 1
```

### Migration Issues
Check logs:
```bash
picoclaw agent --log-level debug
```

## Best Practices

### Production
- Enable JSONL storage
- Enable auto-migration
- Use SSD storage
- Monitor disk space

### Development
- Use default settings
- Enable debug logging
- Test migration path

### Backup
```bash
# Backup sessions
cp -r ~/.picoclaw/sessions/ ~/backup/
```

## API Reference

### Store Interface
```go
type Store interface {
    Append(sessionID string, msg Message) error
    Load(sessionID string) ([]Message, error)
    Truncate(sessionID string, index int) error
    Delete(sessionID string) error
}
```

### JSONL Store
```go
store := jsonl.NewStore(baseDir)
```

## See Also

- [Configuration Guide](CONFIGURATION.md)
- [Session Management](../guides/SESSION_MANAGEMENT.md)
- [v0.2.1 Features](../guides/V0.2.1_FEATURES.md)

---

**Version**: v0.2.1  
**Last Updated**: 2026-03-09
