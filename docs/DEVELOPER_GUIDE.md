# PicoClaw Developer Guide

## Table of Contents

1. [Getting Started](#getting-started)
2. [Development Environment](#development-environment)
3. [Project Structure](#project-structure)
4. [Building and Testing](#building-and-testing)
5. [Core Concepts](#core-concepts)
6. [Adding New Features](#adding-new-features)
7. [Code Style and Standards](#code-style-and-standards)
8. [Debugging](#debugging)
9. [Contributing](#contributing)

## Getting Started

### Prerequisites

- **Go**: 1.25.5 or later
- **Git**: For version control
- **Make**: For build automation
- **golangci-lint**: For code linting (optional)

### Clone and Build

```bash
# Clone repository
git clone https://github.com/sipeed/picoclaw.git
cd picoclaw

# Build
make build

# Run tests
make test

# Run linter
make lint
```

### First Run

```bash
# Initialize configuration
./build/picoclaw-linux-amd64 onboard

# Test agent
./build/picoclaw-linux-amd64 agent -m "Hello, world!"
```

## Development Environment

### Recommended Tools

- **IDE**: VS Code, GoLand, or Vim with Go plugins
- **Debugger**: Delve (dlv)
- **Profiler**: pprof
- **Testing**: Go test with coverage

### VS Code Setup

Install extensions:
- Go (golang.go)
- Go Test Explorer
- Error Lens

`.vscode/settings.json`:
```json
{
  "go.useLanguageServer": true,
  "go.lintTool": "golangci-lint",
  "go.lintOnSave": "package",
  "go.testFlags": ["-v"],
  "go.coverOnSave": true
}
```

### Environment Variables

```bash
# Development
export PICOCLAW_HOME=~/.picoclaw
export PICOCLAW_CONFIG=~/.picoclaw/config.json
export PICOCLAW_LOG_LEVEL=debug

# Testing
export PICOCLAW_TEST_MODE=true
export PICOCLAW_TEST_WORKSPACE=/tmp/picoclaw-test
```

## Project Structure

```
picoclaw/
├── cmd/                          # Command-line applications
│   ├── picoclaw/                 # Main CLI
│   │   ├── main.go              # Entry point
│   │   └── internal/            # CLI commands
│   │       ├── agent/           # Agent command
│   │       ├── gateway/         # Gateway command
│   │       ├── team/            # Team command
│   │       └── ...
│   ├── picoclaw-launcher/       # GUI launcher
│   └── picoclaw-launcher-tui/   # TUI launcher
├── pkg/                          # Reusable packages
│   ├── agent/                   # Agent system
│   ├── team/                    # Multi-agent framework
│   ├── collaborative/           # Collaborative chat
│   ├── providers/               # LLM providers
│   ├── channels/                # Messaging platforms
│   ├── tools/                   # Tool system
│   ├── bus/                     # Message bus
│   ├── config/                  # Configuration
│   ├── logger/                  # Logging
│   └── ...
├── internal/                     # Internal packages
│   ├── memory/                  # Memory management
│   ├── retry/                   # Retry policies
│   └── ratelimit/               # Rate limiting
├── docs/                         # Documentation
├── templates/                    # Configuration templates
├── skills/                       # Built-in skills
├── docker/                       # Docker files
├── build/                        # Build artifacts
├── Makefile                      # Build automation
├── go.mod                        # Go dependencies
└── README.md                     # Project README
```

### Key Directories

- **cmd/**: Command-line applications (main packages)
- **pkg/**: Reusable packages (can be imported by other projects)
- **internal/**: Internal packages (cannot be imported externally)
- **docs/**: Documentation files
- **templates/**: Configuration templates
- **skills/**: Built-in skills (extensibility)

## Building and Testing

### Build Commands

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Build with specific flags
make build GOFLAGS="-v -tags stdjson"

# Install to ~/.local/bin
make install

# Uninstall
make uninstall

# Clean build artifacts
make clean
```

### Testing Commands

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific package tests
go test -v ./pkg/agent/...

# Run specific test
go test -v -run TestAgentLoop ./pkg/agent/

# Run benchmarks
go test -bench=. ./pkg/...

# Run with race detector
go test -race ./...
```

### Linting

```bash
# Run golangci-lint
make lint

# Auto-fix issues
golangci-lint run --fix

# Run specific linters
golangci-lint run --enable-only=gofmt,govet
```

## Core Concepts

### 1. Agent System

**AgentInstance** represents a single agent:

```go
// Create agent instance
agent := agent.NewAgentInstance(
    agentCfg,
    defaults,
    cfg,
    provider,
)

// Execute task
result, err := agent.Execute(ctx, "Write a hello world function")
```

**AgentLoop** handles the execution loop:

```go
// Main loop
for {
    // Get LLM response
    response := provider.Complete(ctx, messages)
    
    // Parse tool calls
    toolCalls := parseToolCalls(response)
    
    // Execute tools
    for _, call := range toolCalls {
        result := toolRegistry.Execute(ctx, call.Name, call.Args)
        messages = append(messages, result)
    }
    
    // Check if done
    if len(toolCalls) == 0 {
        break
    }
}
```

### 2. Multi-Agent Teams

**TeamManager** orchestrates teams:

```go
// Create team manager
tm := team.NewTeamManager(registry, msgBus)

// Load team config
config, _ := team.LoadTeamConfig("dev-team.json")

// Create team
myTeam, _ := tm.CreateTeam(ctx, config)

// Execute task
result, _ := tm.ExecuteTask(ctx, myTeam.ID, "Implement feature X")
```

**Collaboration Patterns**:

```go
// Sequential
type SequentialCoordinator struct {
    roles []string
}

func (c *SequentialCoordinator) Execute(ctx context.Context, task string) {
    for _, role := range c.roles {
        result := executeRole(ctx, role, task)
        task = result // Pass to next role
    }
}

// Parallel
type ParallelCoordinator struct {
    roles []string
}

func (c *ParallelCoordinator) Execute(ctx context.Context, task string) {
    var wg sync.WaitGroup
    results := make(chan Result, len(c.roles))
    
    for _, role := range c.roles {
        wg.Add(1)
        go func(r string) {
            defer wg.Done()
            results <- executeRole(ctx, r, task)
        }(role)
    }
    
    wg.Wait()
    close(results)
}
```

### 3. Collaborative Chat

**Platform Interface**:

```go
type Platform interface {
    SendMessage(ctx context.Context, chatID string, content string) error
    GetTeamManager() TeamManager
    GetContext() context.Context
}
```

**Implementing for New Platform**:

```go
type MyChannel struct {
    chatManager *collaborative.Manager
    teamManager *team.TeamManager
    ctx         context.Context
}

func (c *MyChannel) SendMessage(ctx context.Context, chatID string, content string) error {
    // Send message to platform
    return nil
}

func (c *MyChannel) GetTeamManager() collaborative.TeamManager {
    return c.teamManager
}

func (c *MyChannel) GetContext() context.Context {
    return c.ctx
}

func (c *MyChannel) handleMessage(ctx context.Context, chatID int64, content string) {
    mentions := collaborative.ExtractMentions(content)
    if len(mentions) > 0 {
        c.chatManager.HandleMentions(ctx, c, chatID, teamID, content, mentions, sender, maxContext)
    }
}
```

### 4. Provider System

**Implementing New Provider**:

```go
type MyProvider struct {
    apiKey  string
    apiBase string
}

func (p *MyProvider) Complete(ctx context.Context, messages []Message) (*Response, error) {
    // Call LLM API
    resp, err := http.Post(p.apiBase+"/chat/completions", ...)
    if err != nil {
        return nil, err
    }
    
    // Parse response
    return parseResponse(resp)
}

func (p *MyProvider) CompleteStream(ctx context.Context, messages []Message) (<-chan StreamChunk, error) {
    // Streaming implementation
    chunks := make(chan StreamChunk)
    go func() {
        defer close(chunks)
        // Stream chunks
    }()
    return chunks, nil
}
```

### 5. Tool System

**Creating New Tool**:

```go
type MyTool struct {
    name        string
    description string
}

func (t *MyTool) Name() string {
    return t.name
}

func (t *MyTool) Description() string {
    return t.description
}

func (t *MyTool) Parameters() map[string]any {
    return map[string]any{
        "type": "object",
        "properties": map[string]any{
            "param1": map[string]any{
                "type":        "string",
                "description": "Parameter description",
            },
        },
        "required": []string{"param1"},
    }
}

func (t *MyTool) Execute(ctx context.Context, args map[string]any) *tools.ToolResult {
    param1, ok := args["param1"].(string)
    if !ok {
        return tools.ErrorResult("param1 is required")
    }
    
    // Tool logic
    result := doSomething(param1)
    
    return tools.SuccessResult(result)
}

// Register tool
registry.Register(&MyTool{
    name:        "my_tool",
    description: "Does something useful",
})
```

## Adding New Features

### 1. Adding a New Channel

1. Create package: `pkg/channels/mychannel/`
2. Implement channel interface
3. Add configuration struct
4. Register in channel manager
5. Add documentation
6. Write tests

Example structure:
```
pkg/channels/mychannel/
├── mychannel.go          # Main implementation
├── mychannel_test.go     # Tests
├── config.go             # Configuration
└── README.md             # Documentation
```

### 2. Adding a New Provider

1. Create provider struct in `pkg/providers/`
2. Implement `LLMProvider` interface
3. Add to factory in `factory.go`
4. Add configuration
5. Write tests

### 3. Adding a New Tool

1. Create tool struct in `pkg/tools/`
2. Implement `Tool` interface
3. Register in tool registry
4. Add safety checks if needed
5. Write tests

### 4. Adding a New Collaboration Pattern

1. Create coordinator in `pkg/team/`
2. Implement pattern logic
3. Add to team manager
4. Create example configuration
5. Write tests

## Code Style and Standards

### Go Style Guide

Follow [Effective Go](https://golang.org/doc/effective_go.html) and [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments).

### Naming Conventions

```go
// Packages: lowercase, single word
package agent

// Types: PascalCase
type AgentInstance struct {}

// Functions: camelCase (exported: PascalCase)
func newAgent() *AgentInstance {}
func NewAgent() *AgentInstance {}

// Constants: PascalCase or UPPER_CASE
const MaxRetries = 3
const DEFAULT_TIMEOUT = 30

// Variables: camelCase
var agentCount int
```

### Error Handling

```go
// Always check errors
result, err := doSomething()
if err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}

// Use errors.Is and errors.As
if errors.Is(err, ErrNotFound) {
    // Handle not found
}

// Wrap errors with context
return fmt.Errorf("processing task %s: %w", taskID, err)
```

### Logging

```go
// Use structured logging
logger.InfoCF("agent", "Task started", map[string]any{
    "task_id": taskID,
    "agent_id": agentID,
})

logger.ErrorCF("agent", "Task failed", map[string]any{
    "task_id": taskID,
    "error": err.Error(),
})
```

### Testing

```go
// Table-driven tests
func TestMyFunction(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"valid input", "hello", "HELLO", false},
        {"empty input", "", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := MyFunction(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("MyFunction() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("MyFunction() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Documentation

```go
// Package documentation
// Package agent provides the core agent execution system.
//
// The agent system handles LLM interactions, tool execution,
// and conversation management.
package agent

// Type documentation
// AgentInstance represents a single agent with its configuration,
// tools, and provider.
type AgentInstance struct {
    // ...
}

// Function documentation
// NewAgentInstance creates a new agent instance with the given
// configuration and provider.
//
// Parameters:
//   - agentCfg: Agent-specific configuration
//   - defaults: Default configuration values
//   - cfg: Global configuration
//   - provider: LLM provider for completions
//
// Returns:
//   - *AgentInstance: Configured agent instance
func NewAgentInstance(...) *AgentInstance {
    // ...
}
```

## Debugging

### Using Delve

```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug test
dlv test ./pkg/agent/ -- -test.run TestAgentLoop

# Debug binary
dlv exec ./build/picoclaw-linux-amd64 -- agent -m "test"

# Attach to running process
dlv attach <pid>
```

### Debugging Commands

```
(dlv) break main.main
(dlv) continue
(dlv) next
(dlv) step
(dlv) print variable
(dlv) goroutines
(dlv) goroutine <id>
```

### Profiling

```bash
# CPU profiling
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof

# Memory profiling
go test -memprofile=mem.prof -bench=.
go tool pprof mem.prof

# Trace
go test -trace=trace.out
go tool trace trace.out
```

### Logging

```bash
# Set log level
export PICOCLAW_LOG_LEVEL=debug

# Enable specific component logging
export PICOCLAW_LOG_COMPONENTS=agent,team,provider
```

## Contributing

### Workflow

1. Fork repository
2. Create feature branch: `git checkout -b feature/my-feature`
3. Make changes
4. Run tests: `make test`
5. Run linter: `make lint`
6. Commit: `git commit -m "Add my feature"`
7. Push: `git push origin feature/my-feature`
8. Create Pull Request

### Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add new collaborative chat feature
fix: resolve memory leak in agent loop
docs: update developer guide
test: add tests for team manager
refactor: simplify provider factory
perf: optimize context caching
chore: update dependencies
```

### Pull Request Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Unit tests added/updated
- [ ] Integration tests added/updated
- [ ] Manual testing performed
- [ ] Tool validation tests (if adding/modifying tools)
- [ ] Collaborative chat tests (if modifying mention/cascade logic)

### Test Coverage Requirements

**Critical Components** (require 100% coverage):
- Collaborative chat (mention extraction, depth tracking, cycle detection)
- Tool validation (argument checking, type validation)
- Dispatch tracking (TTL cleanup, concurrency)

**Running Tests**:
```bash
# Run all tests
make test

# Run specific package tests
go test ./pkg/collaborative/... -v
go test ./pkg/tools/... -v

# Run with race detector
go test -race ./pkg/collaborative/...

# Run specific test
go test ./pkg/tools/... -run TestToolArgumentValidation -v
```

**Test Examples**:

1. **Tool Validation Test**:
```go
func TestToolArgumentValidation(t *testing.T) {
    registry := NewToolRegistry()
    tool := &MockToolWithParams{
        name: "test_tool",
        params: map[string]any{
            "required": []string{"param1"},
            "properties": map[string]any{
                "param1": map[string]any{"type": "string"},
            },
        },
    }
    registry.Register(tool)
    
    // Test missing required parameter
    result := registry.Execute(ctx, "test_tool", map[string]any{})
    assert.True(t, result.IsError)
}
```

2. **Collaborative Chat Test**:
```go
func TestMentionDepthLimit(t *testing.T) {
    session := NewSession(123, "team-1", 50)
    
    // Test depth tracking
    session.IncrementMentionDepth()
    assert.Equal(t, 1, session.MentionDepth)
    
    session.DecrementMentionDepth()
    assert.Equal(t, 0, session.MentionDepth)
}
```

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] Tests pass
- [ ] Linter passes
- [ ] No race conditions (verified with -race flag)
- [ ] No memory leaks (verified with profiling if applicable)
```

### Code Review

- Be respectful and constructive
- Focus on code, not person
- Explain reasoning
- Suggest improvements
- Approve when ready

## Resources

### Documentation

- [Architecture](ARCHITECTURE.md)
- [Multi-Agent Guide](MULTI_AGENT_GUIDE.md)
- [Collaborative Chat](COLLABORATIVE_CHAT.md)
- [Safety Levels](SAFETY_LEVELS.md)

### External Resources

- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Go Testing](https://golang.org/pkg/testing/)
- [Delve Debugger](https://github.com/go-delve/delve)

---

**Last Updated**: 2026-03-07  
**Version**: 1.3.0
