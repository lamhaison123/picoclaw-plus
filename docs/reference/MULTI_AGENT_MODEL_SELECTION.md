# Multi-Agent Team Model Selection

## Overview

This document describes the implementation of per-role model selection for multi-agent teams in PicoClaw.

## Problem Statement

Previously, all team agents used the default model (`gemini-3-flash`) regardless of their role configuration. This meant:

- A `developer` role configured with `gpt-5.3-codex` would still use `gemini-3-flash`
- A `reviewer` role configured with `kimi-k2.5` would still use `gemini-3-flash`
- Model configurations in team config files were ignored

## Root Causes

Three critical issues were identified:

### 1. Separate Agent Registry Instances

**Problem**: Two separate `AgentRegistry` instances existed:
- One in `AgentLoop` where agents were registered
- One in `TeamManager` where agents were searched

**Impact**: Agents registered in one registry couldn't be found in the other

**Solution**: 
- Added `GetRegistry()` method to `AgentLoop` to expose internal registry
- Modified `TeamManager` to use shared registry from `AgentLoop`

### 2. Missing Agent Registration on Team Restore

**Problem**: When teams were loaded from disk, agents weren't re-registered

**Impact**: Restored teams had no agents in the registry

**Solution**: Modified `RestoreTeamFromState()` to register agents when loading teams

### 3. No Model-Specific Agent Creation

**Problem**: `CreateTeam()` didn't create agent instances with role-specific models

**Impact**: Even if agents were registered, they used default model

**Solution**: 
- Added `provider` and `cfg` fields to `TeamManager`
- Added `SetProvider()` method to inject LLM provider and config
- Modified `CreateTeam()` to create agents with role-specific models

## Implementation

### Changes to TeamManager

**File**: `pkg/team/manager.go`

```go
type TeamManager struct {
    // ... existing fields
    provider providers.LLMProvider  // NEW
    cfg      *config.Config          // NEW
}

// NEW: Set provider for model selection
func (tm *TeamManager) SetProvider(provider providers.LLMProvider, cfg *config.Config) {
    tm.provider = provider
    tm.cfg = cfg
}
```

### Changes to CreateTeam

**File**: `pkg/team/manager.go`

```go
func (tm *TeamManager) CreateTeam(ctx context.Context, teamConfig *TeamConfig) (*Team, error) {
    // ... team creation logic
    
    // Register agents for each role
    for _, roleConfig := range teamConfig.Roles {
        agentID := fmt.Sprintf("%s-%s", teamConfig.TeamID, roleConfig.Name)
        
        // Register agent instance with correct model if provider is available
        if tm.provider != nil && tm.cfg != nil {
            agentCfg := &config.AgentConfig{
                ID:   agentID,
                Name: fmt.Sprintf("%s (%s)", team.Name, roleConfig.Name),
                Model: &config.AgentModelConfig{
                    Primary: roleConfig.Model,  // Use role-specific model
                },
                Workspace: fmt.Sprintf("%s/teams/%s/%s", tm.workspace, teamConfig.TeamID, roleConfig.Name),
            }
            
            instance := agent.NewAgentInstance(agentCfg, &tm.cfg.Agents.Defaults, tm.cfg, tm.provider)
            tm.registry.RegisterTeamAgent(agentID, instance)
        }
    }
}
```

### Changes to RestoreTeamFromState

**File**: `pkg/team/persistence.go`

```go
func (tm *TeamManager) RestoreTeamFromState(state *TeamState) error {
    // ... restore team logic
    
    // Register agents in registry
    for _, roleConfig := range team.Config.Roles {
        agentID := fmt.Sprintf("%s-%s", team.ID, roleConfig.Name)
        
        if tm.provider != nil && tm.cfg != nil {
            agentCfg := &config.AgentConfig{
                ID:    agentID,
                Name:  fmt.Sprintf("%s (%s)", team.Name, roleConfig.Name),
                Model: &config.AgentModelConfig{
                    Primary: roleConfig.Model,
                },
                Workspace: fmt.Sprintf("%s/teams/%s/%s", tm.workspace, team.ID, roleConfig.Name),
            }
            
            instance := agent.NewAgentInstance(agentCfg, &tm.cfg.Agents.Defaults, tm.cfg, tm.provider)
            tm.registry.RegisterTeamAgent(agentID, instance)
        }
    }
}
```

### Changes to AgentLoop

**File**: `pkg/agent/loop.go`

```go
// NEW: Expose internal registry
func (al *AgentLoop) GetRegistry() *AgentRegistry {
    return al.registry
}

// NEW: Process with specific agent (bypass routing)
func (al *AgentLoop) ProcessWithAgent(
    ctx context.Context,
    agentID string,
    userMessage string,
    sessionKey string,
    channel string,
    chatID string,
) (string, error) {
    agent, err := al.registry.GetAgent(agentID)
    if err != nil {
        return "", fmt.Errorf("agent not found: %s", agentID)
    }
    
    return al.runAgentLoop(ctx, agent, processOptions{
        UserMessage: userMessage,
        SessionKey:  sessionKey,
        Channel:     channel,
        ChatID:      chatID,
    })
}
```

### Changes to Task Executor

**File**: `pkg/team/executor.go`

```go
func (e *DirectAgentExecutor) Execute(ctx context.Context, agentID string, task *Task) (any, error) {
    // Use ProcessWithAgent to specify exact agent
    result, err := e.agentLoop.ProcessWithAgent(ctx, agentID, prompt, sessionKey, channel, task.ID)
    // ...
}
```

### Changes to Team Command Initialization

**File**: `cmd/picoclaw/internal/teamcmd/command.go`

```go
func initializeTeamManager(cfg *config.Config, agentLoop *agent.AgentLoop) *team.TeamManager {
    // Use shared registry from AgentLoop
    registry := agentLoop.GetRegistry()
    
    tm := team.NewTeamManager(registry, msgBus)
    
    // Set provider BEFORE setting team memory
    tm.SetProvider(provider, cfg)
    
    // Then set team memory (which loads persisted teams)
    tm.SetTeamMemory(teamMemory)
    
    return tm
}
```

## Configuration Example

Team configuration with role-specific models:

```json
{
  "team_id": "dev-team-001",
  "name": "Development Team",
  "pattern": "sequential",
  "roles": [
    {
      "name": "developer",
      "model": "gpt-5.3-codex",
      "capabilities": ["code_generation", "implementation"]
    },
    {
      "name": "reviewer",
      "model": "kimi-k2.5",
      "capabilities": ["code_review"]
    },
    {
      "name": "tester",
      "model": "gemini-3-flash",
      "capabilities": ["testing"]
    }
  ]
}
```

## Verification

### Build Test
```bash
go build ./cmd/picoclaw
# ✅ Build successful
```

### Unit Tests
```bash
go test ./...
# ✅ 166 tests passed
```

### Production Test
```bash
picoclaw team execute dev-team-001 -t "Create a hello world function"
```

**Expected Log Output**:
```
[INFO] agent: Registered agent {agent_id=dev-team-001-developer, model=gpt-5.3-codex}
[INFO] agent: Registered agent {agent_id=dev-team-001-reviewer, model=kimi-k2.5}
[INFO] agent: Registered agent {agent_id=dev-team-001-tester, model=gemini-3-flash}
```

## Benefits

1. **Role-Specific Optimization**: Each role uses the most appropriate model
   - Developers use code-specialized models
   - Reviewers use reasoning-focused models
   - Testers use cost-effective models

2. **Cost Efficiency**: Use expensive models only where needed

3. **Performance**: Specialized models perform better for their specific tasks

4. **Flexibility**: Easy to experiment with different model combinations

## Related Documentation

- Team Configuration: `config/config.example.json`
- Multi-Agent Framework: `.kiro/specs/multi-agent-collaboration-framework/`
- Tool Access Control: `docs/TEAM_TOOL_ACCESS.md`

## Known Limitations

1. **Tool Access Control**: Tool restrictions in role configs are not enforced (see `docs/TEAM_TOOL_ACCESS.md`)
2. **Model Fallback**: Fallback models not yet implemented for team agents
3. **Dynamic Model Switching**: Cannot change model for running agents

## Future Improvements

- [ ] Implement tool access control enforcement
- [ ] Add model fallback support for team agents
- [ ] Support dynamic model switching
- [ ] Add model usage metrics per role
- [ ] Implement cost tracking per team
