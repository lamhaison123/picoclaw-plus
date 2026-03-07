// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package team

import (
	"context"
	"fmt"
	"time"

	"github.com/sipeed/picoclaw/pkg/agent"
	"github.com/sipeed/picoclaw/pkg/logger"
)

// AgentExecutor defines interface for executing tasks with agents
type AgentExecutor interface {
	Execute(ctx context.Context, agentID string, task *Task) (any, error)
}

// DirectAgentExecutor executes tasks by calling agent directly
type DirectAgentExecutor struct {
	agentLoop *agent.AgentLoop
}

// NewDirectAgentExecutor creates executor that uses agent loop
func NewDirectAgentExecutor(agentLoop *agent.AgentLoop) *DirectAgentExecutor {
	return &DirectAgentExecutor{
		agentLoop: agentLoop,
	}
}

// Execute runs task using specific agent
func (e *DirectAgentExecutor) Execute(ctx context.Context, agentID string, task *Task) (any, error) {
	if e.agentLoop == nil {
		return nil, fmt.Errorf("agent loop not set")
	}

	// Validate agentID is not empty
	if agentID == "" {
		return nil, fmt.Errorf("agent ID cannot be empty")
	}

	// Build task prompt
	prompt := fmt.Sprintf("Task: %s\n\nRole: %s\n\nPlease complete this task.", task.Description, task.RequiredRole)

	// Add context if available
	if len(task.Context) > 0 {
		prompt += fmt.Sprintf("\n\nContext: %v", task.Context)
	}

	logger.InfoCF("team", "Executing task with agent",
		map[string]any{
			"agent_id": agentID,
			"task_id":  task.ID,
			"role":     task.RequiredRole,
		})

	// Call agent with task using ProcessWithAgent to specify exact agent
	sessionKey := fmt.Sprintf("team:task:%s", task.ID)
	channel := "team"
	result, err := e.agentLoop.ProcessWithAgent(ctx, agentID, prompt, sessionKey, channel, task.ID)
	if err != nil {
		logger.ErrorCF("team", "Task execution failed",
			map[string]any{
				"agent_id": agentID,
				"task_id":  task.ID,
				"error":    err.Error(),
			})
		return nil, err
	}

	logger.InfoCF("team", "Task completed successfully",
		map[string]any{
			"agent_id":   agentID,
			"task_id":    task.ID,
			"result_len": len(result),
		})

	// Return result as structured data
	return map[string]any{
		"status":      "completed",
		"agent_id":    agentID,
		"task_id":     task.ID,
		"role":        task.RequiredRole,
		"description": task.Description,
		"result":      result,
		"timestamp":   time.Now().Format(time.RFC3339),
	}, nil
}

// TeamAwareExecutor executes tasks with team-specific agent configurations
type TeamAwareExecutor struct {
	baseExecutor AgentExecutor
	teamManager  *TeamManager
}

// NewTeamAwareExecutor creates an executor that uses team-specific agent configs
func NewTeamAwareExecutor(baseExecutor AgentExecutor, teamManager *TeamManager) *TeamAwareExecutor {
	return &TeamAwareExecutor{
		baseExecutor: baseExecutor,
		teamManager:  teamManager,
	}
}

// Execute runs task using team-specific agent configuration
func (e *TeamAwareExecutor) Execute(ctx context.Context, agentID string, task *Task) (any, error) {
	// Extract team ID from agent ID (format: teamID-roleName)
	// For now, we'll use the base executor since we need deeper integration
	// to support per-role model selection

	// TODO: Implement model selection based on team role configuration
	// This requires:
	// 1. Parse agentID to extract teamID and role
	// 2. Look up role config from team
	// 3. Create agent instance with role-specific model
	// 4. Execute task with that agent

	logger.InfoCF("team", "Executing with team-aware executor",
		map[string]any{
			"agent_id": agentID,
			"task_id":  task.ID,
		})

	return e.baseExecutor.Execute(ctx, agentID, task)
}
