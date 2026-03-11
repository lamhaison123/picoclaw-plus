# Team Agent Tool Access Control

## Overview

This document describes the tool access control mechanism for multi-agent teams and identifies current limitations.

## Current Implementation

### Tool Registration

When a team agent is created via `NewAgentInstance()`, the following tools are automatically registered:

- `readFile` - Read file contents
- `writeFile` - Write to files
- `listDir` - List directory contents
- `exec` - Execute shell commands
- `editFile` - Edit file contents
- `appendFile` - Append to files

### Role-Based Tool Configuration

Team configurations support role-based tool restrictions via the `tools` field in role definitions:

```json
{
  "roles": [
    {
      "name": "developer",
      "tools": ["readCode", "editCode", "fsWrite", "executePwsh"],
      "capabilities": ["code_generation", "implementation"]
    },
    {
      "name": "reviewer",
      "tools": ["readCode", "readFile", "grepSearch"],
      "capabilities": ["code_review", "quality_inspection"]
    }
  ]
}
```

### ValidateToolAccess Method

The `TeamManager` provides a `ValidateToolAccess()` method that:

- Checks if an agent is allowed to use a specific tool
- Supports wildcard patterns (e.g., `"file_*"` matches `file_read`, `file_write`)
- Returns `true` if tool is allowed, `false` otherwise

**Location**: `pkg/team/manager.go`

```go
func (tm *TeamManager) ValidateToolAccess(agentID, toolName string) bool
```

## Current Limitation

### Issue: Tool Restrictions Not Enforced

**Status**: ⚠️ Tool access validation exists but is NOT enforced during execution

**Details**:

1. The `ValidateToolAccess()` method is implemented but never called
2. Tool execution happens in `pkg/agent/loop.go` via `agent.Tools.ExecuteWithContext()`
3. `ToolRegistry.ExecuteWithContext()` has no permission checking mechanism
4. Team agents can use ANY tool in their registry regardless of role configuration

**Impact**:

- A `reviewer` role configured with read-only tools can still write files
- A `tester` role can execute arbitrary shell commands even if not permitted
- Tool restrictions in team configuration files are currently ignored

### Example

Given this configuration:

```json
{
  "name": "reviewer",
  "tools": ["readCode", "readFile", "grepSearch"],
  "capabilities": ["code_review"]
}
```

The reviewer agent can STILL use:
- `writeFile` ❌ (should be blocked)
- `editCode` ❌ (should be blocked)
- `executePwsh` ❌ (should be blocked)

## Proposed Solutions

### Option 1: Add Validation to ToolRegistry

Modify `ToolRegistry.ExecuteWithContext()` to accept a validation callback:

```go
func (r *ToolRegistry) ExecuteWithContext(
    ctx context.Context,
    name string,
    args map[string]any,
    channel, chatID string,
    asyncCallback AsyncCallback,
    validateFunc func(toolName string) bool, // NEW
) *ToolResult {
    // Check permission before execution
    if validateFunc != nil && !validateFunc(name) {
        return ErrorResult(fmt.Sprintf("tool %q not permitted for this agent", name))
    }
    
    // ... rest of execution
}
```

**Pros**: Centralized validation, works for all agents
**Cons**: Requires passing validation function through call chain

### Option 2: Filtered Tool Registry per Agent

Create a filtered tool registry when initializing team agents:

```go
func createFilteredRegistry(allowedTools []string, baseRegistry *ToolRegistry) *ToolRegistry {
    filtered := NewToolRegistry()
    for _, toolName := range allowedTools {
        if tool, ok := baseRegistry.Get(toolName); ok {
            filtered.Register(tool)
        }
    }
    return filtered
}
```

**Pros**: Simple, no runtime checks needed
**Cons**: More memory usage, need to manage multiple registries

### Option 3: Wrapper Tool with Permission Check

Wrap each tool with a permission-checking decorator:

```go
type PermissionCheckedTool struct {
    baseTool Tool
    agentID string
    validator func(agentID, toolName string) bool
}

func (t *PermissionCheckedTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
    if !t.validator(t.agentID, t.baseTool.Name()) {
        return ErrorResult("permission denied")
    }
    return t.baseTool.Execute(ctx, args)
}
```

**Pros**: Flexible, can add logging/auditing
**Cons**: More complex, wrapper overhead

## Recommendation

**Option 2 (Filtered Tool Registry)** is recommended because:

1. Zero runtime overhead - validation happens at agent creation time
2. Simple implementation - just filter tools during registration
3. Clear separation - each agent has exactly the tools it needs
4. Easy to debug - inspect agent's tool registry to see what's available

## Implementation Checklist

- [ ] Implement filtered tool registry creation
- [ ] Modify `CreateTeam()` to use filtered registries
- [ ] Modify `RestoreTeamFromState()` to use filtered registries
- [ ] Add tests for tool access control
- [ ] Update documentation with examples
- [ ] Add logging for denied tool access attempts

## Testing

Test cases needed:

1. Agent with read-only tools cannot write files
2. Agent with write tools can write files
3. Wildcard patterns work correctly (`"file_*"` allows `file_read`, `file_write`)
4. Agent with `"*"` can use all tools
5. Attempting to use unpermitted tool returns clear error message

## Related Files

- `pkg/team/manager.go` - Team management and ValidateToolAccess
- `pkg/agent/instance.go` - Agent instance creation and tool registration
- `pkg/agent/loop.go` - Agent execution loop and tool calling
- `pkg/tools/registry.go` - Tool registry implementation
- `pkg/team/executor.go` - Task execution with agents

## References

- Multi-Agent Collaboration Framework: `.kiro/specs/multi-agent-collaboration-framework/`
- Tool Implementation: `pkg/tools/`
- Agent Configuration: `config/config.example.json`
