// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package tools

import (
	"context"
	"fmt"
)

// TeamManager interface defines the contract for team management operations
// This allows the tool to work with the team system without tight coupling
type TeamManager interface {
	ExecuteTask(ctx context.Context, teamID string, taskDescription string) (any, error)
	ExecuteTaskWithRole(ctx context.Context, teamID string, taskDescription string, role string) (any, error)
	GetTeamStatus(teamID string) (interface{}, error)
	GetAllTeams() map[string]interface{}
}

// TeamDelegationTool allows PicoClaw to delegate tasks to multi-agent teams
type TeamDelegationTool struct {
	teamManager TeamManager
}

// NewTeamDelegationTool creates a new team delegation tool
func NewTeamDelegationTool(teamManager TeamManager) *TeamDelegationTool {
	return &TeamDelegationTool{
		teamManager: teamManager,
	}
}

func (t *TeamDelegationTool) Name() string {
	return "delegate_to_team"
}

func (t *TeamDelegationTool) Description() string {
	return "Delegate a complex task to a specialized multi-agent team. Use this when a task requires multiple specialized skills or parallel execution. The team will coordinate internally and return the final result."
}

func (t *TeamDelegationTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"team_id": map[string]any{
				"type":        "string",
				"description": "ID of the team to delegate to (e.g., 'dev-team')",
			},
			"task": map[string]any{
				"type":        "string",
				"description": "Description of the task to delegate to the team",
			},
			"role": map[string]any{
				"type":        "string",
				"description": "Optional: specific role within the team to assign the task to (e.g., 'backend', 'frontend', 'reviewer')",
			},
		},
		"required": []string{"team_id", "task"},
	}
}

func (t *TeamDelegationTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	if t.teamManager == nil {
		return ErrorResult("Team system not configured").WithError(fmt.Errorf("team manager is nil"))
	}

	// Extract parameters
	teamID, ok := args["team_id"].(string)
	if !ok || teamID == "" {
		return ErrorResult("team_id is required").WithError(fmt.Errorf("team_id parameter is required"))
	}

	task, ok := args["task"].(string)
	if !ok || task == "" {
		return ErrorResult("task is required").WithError(fmt.Errorf("task parameter is required"))
	}

	role, _ := args["role"].(string)

	// Execute task with team
	var result any
	var err error

	if role != "" {
		// Execute with specific role
		result, err = t.teamManager.ExecuteTaskWithRole(ctx, teamID, task, role)
	} else {
		// Execute with automatic role selection
		result, err = t.teamManager.ExecuteTask(ctx, teamID, task)
	}

	if err != nil {
		return &ToolResult{
			ForLLM:  fmt.Sprintf("Team execution failed: %v", err),
			ForUser: fmt.Sprintf("❌ Team '%s' failed to complete the task: %v", teamID, err),
			Err:     err,
		}
	}

	// Format result for LLM and user
	resultStr := fmt.Sprintf("%v", result)

	return &ToolResult{
		ForLLM:  fmt.Sprintf("Team '%s' completed task successfully. Result: %s", teamID, resultStr),
		ForUser: fmt.Sprintf("✓ Team '%s' completed the task:\n%s", teamID, resultStr),
	}
}

// TeamStatusTool allows checking the status of teams
type TeamStatusTool struct {
	teamManager TeamManager
}

// NewTeamStatusTool creates a new team status tool
func NewTeamStatusTool(teamManager TeamManager) *TeamStatusTool {
	return &TeamStatusTool{
		teamManager: teamManager,
	}
}

func (t *TeamStatusTool) Name() string {
	return "team_status"
}

func (t *TeamStatusTool) Description() string {
	return "Check the status of a team or list all available teams. Use this to see what teams are available and their current state."
}

func (t *TeamStatusTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"team_id": map[string]any{
				"type":        "string",
				"description": "Optional: ID of a specific team to check. If not provided, lists all teams.",
			},
		},
	}
}

func (t *TeamStatusTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	if t.teamManager == nil {
		return ErrorResult("Team system not configured").WithError(fmt.Errorf("team manager is nil"))
	}

	teamID, _ := args["team_id"].(string)

	if teamID != "" {
		// Get specific team status
		status, err := t.teamManager.GetTeamStatus(teamID)
		if err != nil {
			return ErrorResult(fmt.Sprintf("Failed to get team status: %v", err)).WithError(err)
		}

		statusStr := fmt.Sprintf("%v", status)
		return &ToolResult{
			ForLLM:  fmt.Sprintf("Team '%s' status: %s", teamID, statusStr),
			ForUser: fmt.Sprintf("Team '%s' status:\n%s", teamID, statusStr),
		}
	}

	// List all teams
	teams := t.teamManager.GetAllTeams()
	if len(teams) == 0 {
		return &ToolResult{
			ForLLM:  "No teams are currently configured",
			ForUser: "No teams available",
		}
	}

	teamList := ""
	for id := range teams {
		teamList += fmt.Sprintf("- %s\n", id)
	}

	return &ToolResult{
		ForLLM:  fmt.Sprintf("Available teams (%d): %s", len(teams), teamList),
		ForUser: fmt.Sprintf("Available teams (%d):\n%s", len(teams), teamList),
	}
}
