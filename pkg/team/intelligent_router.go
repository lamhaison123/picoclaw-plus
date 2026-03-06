// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package team

import (
	"context"
	"fmt"
	"strings"

	"github.com/sipeed/picoclaw/pkg/logger"
)

// IntelligentRouter uses LLM to intelligently route tasks to appropriate roles
type IntelligentRouter struct {
	executor AgentExecutor
	teamID   string // Team ID for constructing full agent IDs
}

// NewIntelligentRouter creates a new intelligent router
func NewIntelligentRouter(executor AgentExecutor) *IntelligentRouter {
	return &IntelligentRouter{
		executor: executor,
		teamID:   "", // Will be set when used
	}
}

// SetTeamID sets the team ID for this router
func (ir *IntelligentRouter) SetTeamID(teamID string) {
	ir.teamID = teamID
}

// DetermineRole uses LLM to analyze task and determine the best role
func (ir *IntelligentRouter) DetermineRole(ctx context.Context, taskDescription string, availableRoles []RoleConfig) (string, error) {
	if len(availableRoles) == 0 {
		return "", fmt.Errorf("no roles available")
	}

	// Build prompt for role selection
	var rolesDesc strings.Builder
	for _, role := range availableRoles {
		rolesDesc.WriteString(fmt.Sprintf("- %s: %s (capabilities: %s)\n",
			role.Name,
			role.Description,
			strings.Join(role.Capabilities, ", ")))
	}

	prompt := fmt.Sprintf(`You are a team manager. Analyze the following task and determine which role should handle it.

Available roles:
%s

Task: %s

Respond with ONLY the role name (one word, lowercase). Choose the most appropriate role based on the task requirements.

Examples:
- "write code for user authentication" -> developer
- "design system architecture" -> architect
- "write unit tests" -> tester
- "review code quality" -> reviewer
- "plan project timeline" -> manager

Role:`, rolesDesc.String(), taskDescription)

	logger.DebugCF("team", "Determining role for task",
		map[string]any{
			"task":            taskDescription,
			"available_roles": len(availableRoles),
		})

	// Construct full agent ID for manager role
	managerAgentID := "manager"
	if ir.teamID != "" {
		managerAgentID = fmt.Sprintf("%s-manager", ir.teamID)
	}

	// Use manager agent to determine role
	result, err := ir.executor.Execute(ctx, managerAgentID, &Task{
		ID:          "role-determination",
		Description: prompt,
		Status:      TaskStatusPending,
	})

	if err != nil {
		return "", fmt.Errorf("failed to determine role: %w", err)
	}

	// Extract role from result
	// Result is a map with structure: {status, agent_id, task_id, role, description, result, timestamp}
	// We need to extract the "result" field which contains the actual LLM response
	var roleStr string
	if resultMap, ok := result.(map[string]any); ok {
		if resultField, ok := resultMap["result"]; ok {
			roleStr = strings.TrimSpace(fmt.Sprintf("%v", resultField))
		} else {
			// Fallback: use entire result if "result" field not found
			roleStr = strings.TrimSpace(fmt.Sprintf("%v", result))
		}
	} else {
		// Fallback: result is not a map, use as-is
		roleStr = strings.TrimSpace(fmt.Sprintf("%v", result))
	}
	roleStr = strings.ToLower(roleStr)

	// Remove common prefixes/suffixes
	roleStr = strings.TrimPrefix(roleStr, "role:")
	roleStr = strings.TrimSpace(roleStr)

	// Extract first word if multiple words
	words := strings.Fields(roleStr)
	if len(words) > 0 {
		roleStr = words[0]
	}

	// Validate role exists
	for _, role := range availableRoles {
		if strings.EqualFold(role.Name, roleStr) {
			logger.InfoCF("team", "Role determined",
				map[string]any{
					"task":          taskDescription,
					"selected_role": role.Name,
				})
			return role.Name, nil
		}
	}

	// Fallback to first role if no match
	logger.WarnCF("team", "Could not match determined role, using fallback",
		map[string]any{
			"determined_role": roleStr,
			"fallback_role":   availableRoles[0].Name,
		})

	return availableRoles[0].Name, nil
}

// DetermineRoleSimple uses keyword matching for simple role determination (faster, no LLM call)
func (ir *IntelligentRouter) DetermineRoleSimple(taskDescription string, availableRoles []RoleConfig) string {
	if len(availableRoles) == 0 {
		return ""
	}

	taskLower := strings.ToLower(taskDescription)

	// Keyword mapping for common tasks
	keywords := map[string][]string{
		"architect": {"design", "architecture", "system design", "plan", "structure"},
		"developer": {"code", "implement", "write", "develop", "program", "function", "class"},
		"coder":     {"code", "implement", "write", "develop", "program"},
		"tester":    {"test", "testing", "verify", "validate", "check", "bug"},
		"reviewer":  {"review", "audit", "check quality", "inspect", "evaluate"},
		"manager":   {"manage", "coordinate", "organize", "plan", "oversee"},
	}

	// Score each role based on keyword matches
	scores := make(map[string]int)
	for _, role := range availableRoles {
		score := 0
		if kws, exists := keywords[strings.ToLower(role.Name)]; exists {
			for _, kw := range kws {
				if strings.Contains(taskLower, kw) {
					score++
				}
			}
		}
		scores[role.Name] = score
	}

	// Find role with highest score
	maxScore := 0
	selectedRole := availableRoles[0].Name
	for _, role := range availableRoles {
		if scores[role.Name] > maxScore {
			maxScore = scores[role.Name]
			selectedRole = role.Name
		}
	}

	logger.InfoCF("team", "Role determined (simple)",
		map[string]any{
			"task":          taskDescription,
			"selected_role": selectedRole,
			"score":         maxScore,
		})

	return selectedRole
}
