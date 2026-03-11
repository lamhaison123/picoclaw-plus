// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package tools

import (
	"context"
	"fmt"
	"strings"
)

// MentionHandler defines the interface for handling mentions
type MentionHandler interface {
	// HandleMention processes a mention and returns the response
	HandleMention(ctx context.Context, mentionedID, message, channel, chatID string) (string, error)
}

// MentionTool allows agents to mention and communicate with each other
type MentionTool struct {
	handler       MentionHandler
	currentAgent  string // ID of the agent using this tool
	teamManager   TeamManager
	agentRegistry AgentRegistry
}

// AgentRegistry interface for accessing other agents
type AgentRegistry interface {
	GetAgent(agentID string) (interface{}, bool)
	ListAgentIDs() []string
}

// NewMentionTool creates a new mention tool
func NewMentionTool(currentAgent string, handler MentionHandler, teamManager TeamManager, registry AgentRegistry) *MentionTool {
	return &MentionTool{
		handler:       handler,
		currentAgent:  currentAgent,
		teamManager:   teamManager,
		agentRegistry: registry,
	}
}

func (t *MentionTool) Name() string {
	return "mention"
}

func (t *MentionTool) Description() string {
	return "Mention another agent or team to get their input or delegate a task. Use @agent_id or @team_id syntax. Examples: '@picoclaw what do you think?', '@dev-team implement this feature', '@manager review this code'."
}

func (t *MentionTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"target": map[string]any{
				"type":        "string",
				"description": "The agent or team to mention (e.g., 'picoclaw', 'dev-team', 'manager'). Do not include @ symbol.",
			},
			"message": map[string]any{
				"type":        "string",
				"description": "The message to send to the mentioned agent or team",
			},
		},
		"required": []string{"target", "message"},
	}
}

func (t *MentionTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	target, ok := args["target"].(string)
	if !ok || target == "" {
		return ErrorResult("target is required").WithError(fmt.Errorf("target parameter is required"))
	}

	message, ok := args["message"].(string)
	if !ok || message == "" {
		return ErrorResult("message is required").WithError(fmt.Errorf("message parameter is required"))
	}

	// Remove @ prefix if present
	target = strings.TrimPrefix(target, "@")

	// Check if target is a team
	if t.teamManager != nil {
		teams := t.teamManager.GetAllTeams()
		if _, exists := teams[target]; exists {
			// Delegate to team
			result, err := t.teamManager.ExecuteTask(ctx, target, message)
			if err != nil {
				return &ToolResult{
					ForLLM:  fmt.Sprintf("@%s failed to respond: %v", target, err),
					ForUser: fmt.Sprintf("❌ @%s is unavailable: %v", target, err),
					Err:     err,
				}
			}

			resultStr := fmt.Sprintf("%v", result)
			return &ToolResult{
				ForLLM:  fmt.Sprintf("@%s responded: %s", target, resultStr),
				ForUser: fmt.Sprintf("@%s: %s", target, resultStr),
			}
		}
	}

	// Check if target is an agent
	if t.agentRegistry != nil {
		if _, exists := t.agentRegistry.GetAgent(target); exists {
			// Use mention handler to process
			if t.handler != nil {
				response, err := t.handler.HandleMention(ctx, target, message, "system", fmt.Sprintf("mention:%s", t.currentAgent))
				if err != nil {
					return &ToolResult{
						ForLLM:  fmt.Sprintf("@%s failed to respond: %v", target, err),
						ForUser: fmt.Sprintf("❌ @%s is unavailable: %v", target, err),
						Err:     err,
					}
				}

				return &ToolResult{
					ForLLM:  fmt.Sprintf("@%s responded: %s", target, response),
					ForUser: fmt.Sprintf("@%s: %s", target, response),
				}
			}
		}
	}

	// Target not found
	availableTargets := []string{}
	if t.agentRegistry != nil {
		availableTargets = append(availableTargets, t.agentRegistry.ListAgentIDs()...)
	}
	if t.teamManager != nil {
		teams := t.teamManager.GetAllTeams()
		for teamID := range teams {
			availableTargets = append(availableTargets, teamID)
		}
	}

	return ErrorResult(fmt.Sprintf("@%s not found. Available: %s", target, strings.Join(availableTargets, ", "))).
		WithError(fmt.Errorf("target not found: %s", target))
}

// ParseMentions extracts all mentions from a message
// Returns a map of mentioned IDs to their context in the message
func ParseMentions(message string) map[string]string {
	mentions := make(map[string]string)
	words := strings.Fields(message)

	for i, word := range words {
		if strings.HasPrefix(word, "@") {
			mentionedID := strings.TrimPrefix(word, "@")
			mentionedID = strings.Trim(mentionedID, ".,!?;:")

			// Get context (surrounding words)
			start := max(0, i-3)
			end := min(len(words), i+4)
			context := strings.Join(words[start:end], " ")

			mentions[mentionedID] = context
		}
	}

	return mentions
}

// ContainsMention checks if a message contains a mention of a specific ID
func ContainsMention(message, targetID string) bool {
	mentions := ParseMentions(message)
	_, exists := mentions[targetID]
	return exists
}

// StripMentions removes all @mentions from a message
func StripMentions(message string) string {
	words := strings.Fields(message)
	filtered := make([]string, 0, len(words))

	for _, word := range words {
		if !strings.HasPrefix(word, "@") {
			filtered = append(filtered, word)
		}
	}

	return strings.Join(filtered, " ")
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
