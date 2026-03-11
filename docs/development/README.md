# Development Documentation

Resources for developers contributing to PicoClaw.

## 📚 Contents

### Getting Started
- **[DEVELOPER_GUIDE.md](./DEVELOPER_GUIDE.md)** - Complete development guide
  - Development setup
  - Build instructions
  - Testing guidelines
  - Code style
  - Contribution workflow

### Architecture Deep Dive
- **[COLLABORATIVE_CHAT_ARCHITECTURE.md](./COLLABORATIVE_CHAT_ARCHITECTURE.md)** - Chat system design
  - Mention detection
  - Cascade mechanism
  - Queue management
  - Context compaction
  - Implementation details

- **[COLLABORATIVE_CHAT_FLOW.txt](./COLLABORATIVE_CHAT_FLOW.txt)** - Message flow diagram
  - ASCII flow diagram
  - State transitions
  - Error handling

### Troubleshooting
- **[troubleshooting.md](./troubleshooting.md)** - Common issues and solutions
  - Build errors
  - Runtime issues
  - Configuration problems
  - Performance issues
  - Debug techniques

## 🚀 Quick Start

### Prerequisites
```bash
# Required
go 1.21+
git

# Optional (for full features)
docker          # For Qdrant
python 3.8+     # For MCP servers
```

### Clone and Build
```bash
git clone https://github.com/sipeed/picoclaw.git
cd picoclaw
go mod download
go build ./cmd/picoclaw
```

### Run Tests
```bash
# All tests
go test ./...

# Specific package
go test ./pkg/agent

# With coverage
go test -cover ./...

# Without Qdrant
go test -tags=no_qdrant ./...
```

### Build Tags
```bash
# Without Qdrant (default for development)
go build -tags=no_qdrant ./cmd/picoclaw

# With all features
go build ./cmd/picoclaw
```

## 🏗️ Project Structure

```
picoclaw/
├── cmd/                    # Command-line tools
│   ├── picoclaw/          # Main application
│   └── picoclaw-launcher/ # GUI launcher
├── pkg/                    # Core packages
│   ├── agent/             # Agent loop
│   ├── bus/               # Message bus
│   ├── channels/          # Platform integrations
│   ├── collaborative/     # Multi-agent chat
│   ├── config/            # Configuration
│   ├── embedding/         # Embedding services
│   ├── memory/            # Memory systems
│   ├── providers/         # LLM providers
│   ├── team/              # Team coordination
│   └── tools/             # Tool system
├── internal/              # Internal packages
├── docs/                  # Documentation
├── config/                # Config templates
└── workspace/             # Default workspace
```

## 🔧 Development Workflow

### 1. Create Feature Branch
```bash
git checkout -b feature/my-feature
```

### 2. Make Changes
- Write code following Go conventions
- Add tests for new functionality
- Update documentation
- Run linters and formatters

### 3. Test Locally
```bash
# Format code
go fmt ./...

# Run linters
golangci-lint run

# Run tests
go test ./...

# Build
go build ./cmd/picoclaw
```

### 4. Commit and Push
```bash
git add .
git commit -m "feat: add new feature"
git push origin feature/my-feature
```

### 5. Create Pull Request
- Describe changes clearly
- Reference related issues
- Ensure CI passes
- Request review

## 📝 Code Style

### Go Conventions
- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Run `golangci-lint` before committing
- Write clear comments for exported functions

### Naming
- Packages: lowercase, single word
- Files: lowercase with underscores
- Types: PascalCase
- Functions: camelCase (exported: PascalCase)
- Constants: PascalCase or UPPER_CASE

### Error Handling
```go
// Good: wrap errors with context
if err != nil {
    return fmt.Errorf("failed to process message: %w", err)
}

// Good: check errors immediately
result, err := doSomething()
if err != nil {
    return err
}
```

### Logging
```go
// Use structured logging
logger.InfoCF("component", "message",
    map[string]any{
        "key": value,
    })

// Log levels: Debug, Info, Warn, Error
```

## 🧪 Testing Guidelines

### Unit Tests
- Test file: `*_test.go`
- Test function: `TestFunctionName`
- Use table-driven tests
- Mock external dependencies

```go
func TestProcessMessage(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"valid", "hello", "response", false},
        {"empty", "", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ProcessMessage(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
            }
            if got != tt.want {
                t.Errorf("got %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Integration Tests
- Use build tags: `// +build integration`
- Require external services
- Run separately: `go test -tags=integration`

### Test Coverage
```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View in browser
go tool cover -html=coverage.out
```

## 🐛 Debugging

### Enable Debug Logging
```json
{
  "log_level": "debug"
}
```

### Use Delve Debugger
```bash
# Install
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug
dlv debug ./cmd/picoclaw
```

### Common Issues
See [troubleshooting.md](./troubleshooting.md) for:
- Build errors
- Runtime panics
- Memory leaks
- Goroutine leaks
- Race conditions

## 📊 Performance Profiling

### CPU Profile
```bash
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof
```

### Memory Profile
```bash
go test -memprofile=mem.prof -bench=.
go tool pprof mem.prof
```

### Race Detector
```bash
go test -race ./...
go build -race ./cmd/picoclaw
```

## 🔍 Code Review Checklist

- [ ] Code follows Go conventions
- [ ] Tests added/updated
- [ ] Documentation updated
- [ ] No race conditions
- [ ] Error handling proper
- [ ] Logging appropriate
- [ ] Performance considered
- [ ] Security reviewed
- [ ] Breaking changes noted

## 🔗 Related Documentation

- [Architecture Overview](../architecture/ARCHITECTURE_OVERVIEW.md) - System design
- [Component Details](../architecture/COMPONENT_DETAILS.md) - Implementation details
- [API Reference](../reference/API_REFERENCE.md) - API documentation
- [Contributing Guide](../../CONTRIBUTING.md) - Contribution guidelines

## 📚 Additional Resources

### Go Resources
- [Go Documentation](https://golang.org/doc/)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

### Tools
- [golangci-lint](https://golangci-lint.run/) - Linter aggregator
- [delve](https://github.com/go-delve/delve) - Debugger
- [pprof](https://golang.org/pkg/net/http/pprof/) - Profiler

### Community
- [GitHub Discussions](https://github.com/sipeed/picoclaw/discussions)
- [Issue Tracker](https://github.com/sipeed/picoclaw/issues)
- [Pull Requests](https://github.com/sipeed/picoclaw/pulls)

## 📞 Getting Help

- **Development Questions**: [GitHub Discussions](https://github.com/sipeed/picoclaw/discussions)
- **Bug Reports**: [GitHub Issues](https://github.com/sipeed/picoclaw/issues)
- **Feature Requests**: [GitHub Issues](https://github.com/sipeed/picoclaw/issues)

---

**Last Updated**: 2026-03-09
