// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package collaborative

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/sipeed/picoclaw/pkg/bus"
	"github.com/sipeed/picoclaw/pkg/logger"
)

// Manager manages all active collaborative chat sessions
type Manager struct {
	sessions map[int64]*Session
	mu       sync.RWMutex
}

// NewManager creates a new collaborative chat manager
func NewManager() *Manager {
	return &Manager{
		sessions: make(map[int64]*Session),
	}
}

// GetOrCreateSession gets an existing session or creates a new one
func (m *Manager) GetOrCreateSession(chatID int64, teamID string, maxContext int) *Session {
	m.mu.Lock()
	defer m.mu.Unlock()

	if session, exists := m.sessions[chatID]; exists {
		return session
	}

	session := NewSession(chatID, teamID, maxContext)
	m.sessions[chatID] = session

	logger.InfoCF("collaborative", "Created new session", map[string]any{
		"session_id": session.SessionID,
		"chat_id":    chatID,
		"team_id":    teamID,
	})

	return session
}

// GetSession gets an existing session
func (m *Manager) GetSession(chatID int64) (*Session, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	session, exists := m.sessions[chatID]
	return session, exists
}

// RemoveSession removes a session
func (m *Manager) RemoveSession(chatID int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.sessions, chatID)

	logger.InfoCF("collaborative", "Removed session", map[string]any{
		"chat_id": chatID,
	})
}

// HandleMentions processes mentions and triggers agents
func (m *Manager) HandleMentions(
	ctx context.Context,
	platform Platform,
	chatID int64,
	teamID string,
	content string,
	mentions []string,
	sender bus.SenderInfo,
	maxContext int,
) error {
	if len(mentions) == 0 {
		return nil
	}

	// Get or create session
	session := m.GetOrCreateSession(chatID, teamID, maxContext)

	// Add user message to context
	session.AddMessage("user", content, mentions)

	// Get team manager
	teamManager := platform.GetTeamManager()
	if teamManager == nil {
		logger.ErrorC("collaborative", "Team manager not available")
		return fmt.Errorf("team manager not available")
	}

	// Get team roster
	var teamRoster string
	if teamInfo, err := teamManager.GetTeam(teamID); err == nil {
		teamRoster = BuildTeamRoster(teamInfo)
	}

	// Build context for agents
	contextStr := session.GetContextAsString()

	// Execute each mentioned role in parallel
	for _, role := range mentions {
		// Update agent status
		session.UpdateAgentStatus(role, "thinking")

		// Build prompt with full context including team roster
		prompt := fmt.Sprintf(`%s

=== Team Information ===
%s

User message: %s

You are @%s. Respond to the user's message considering the conversation history above.
You can mention other team members using @role format (e.g., @architect, @developer).`,
			contextStr, teamRoster, content, role)

		// Execute agent with role
		go func(r string) {
			// Use platform's context
			result, err := teamManager.ExecuteTaskWithRole(platform.GetContext(), teamID, prompt, r)

			if err != nil {
				logger.ErrorCF("collaborative", "Agent execution failed", map[string]any{
					"role":  r,
					"error": err.Error(),
				})
				session.UpdateAgentStatus(r, "error")

				// Send user-friendly error message
				var errorMsg string
				if strings.Contains(err.Error(), "not found in team configuration") {
					errorMsg = fmt.Sprintf("❌ Role @%s không tồn tại trong team configuration.\n\n💡 Các role có sẵn: @architect, @developer, @tester, @manager", r)
				} else {
					errorMsg = fmt.Sprintf("❌ @%s encountered an error: %v", r, err)
				}
				platform.SendMessage(platform.GetContext(), fmt.Sprintf("%d", chatID), errorMsg)
				return
			}

			// Extract response text from result
			responseStr := extractResponseText(result)
			session.AddMessage(r, responseStr, nil)
			session.UpdateAgentStatus(r, "idle")

			// Format with IRC-style prefix
			formattedMsg := FormatMessage(session.SessionID, r, responseStr)

			// Send to platform
			platform.SendMessage(platform.GetContext(), fmt.Sprintf("%d", chatID), formattedMsg)

			logger.InfoCF("collaborative", "Agent response sent", map[string]any{
				"role":       r,
				"session_id": session.SessionID,
			})

			// Check if agent mentioned other agents in their response
			mentionsInResponse := ExtractMentions(responseStr)
			if len(mentionsInResponse) > 0 {
				// Filter out self-mentions and already active agents
				newMentions := []string{}
				for _, mentioned := range mentionsInResponse {
					if mentioned != r { // Don't mention self
						newMentions = append(newMentions, mentioned)
					}
				}

				if len(newMentions) > 0 {
					logger.InfoCF("collaborative", "Agent mentioned other agents", map[string]any{
						"from_role": r,
						"mentions":  newMentions,
					})

					// Trigger mentioned agents with updated context
					contextStr := session.GetContextAsString()
					for _, mentionedRole := range newMentions {
						session.UpdateAgentStatus(mentionedRole, "thinking")

						prompt := fmt.Sprintf(`%s

=== Team Information ===
%s

@%s mentioned you in their response. Please respond if needed.

You are @%s. You can mention other team members using @role format.`,
							contextStr, teamRoster, r, mentionedRole)

						// Execute in new goroutine
						go func(targetRole string) {
							result, err := teamManager.ExecuteTaskWithRole(platform.GetContext(), teamID, prompt, targetRole)
							if err != nil {
								logger.ErrorCF("collaborative", "Mentioned agent execution failed", map[string]any{
									"role":  targetRole,
									"error": err.Error(),
								})
								session.UpdateAgentStatus(targetRole, "error")
								return
							}

							responseStr := extractResponseText(result)
							session.AddMessage(targetRole, responseStr, nil)
							session.UpdateAgentStatus(targetRole, "idle")

							formattedMsg := FormatMessage(session.SessionID, targetRole, responseStr)
							platform.SendMessage(platform.GetContext(), fmt.Sprintf("%d", chatID), formattedMsg)

							logger.InfoCF("collaborative", "Mentioned agent response sent", map[string]any{
								"role":       targetRole,
								"session_id": session.SessionID,
							})
						}(mentionedRole)
					}
				}
			}
		}(role)
	}

	return nil
}

// extractResponseText extracts the actual response text from team execution result
func extractResponseText(result any) string {
	if result == nil {
		return ""
	}

	// If result is a string, return it directly
	if str, ok := result.(string); ok {
		return str
	}

	// If result is an array (from ExecuteParallel), extract first element
	if arr, ok := result.([]any); ok {
		if len(arr) > 0 {
			return extractResponseText(arr[0]) // Recursive call for first element
		}
		return ""
	}

	// If result is a map, try to extract "result" field
	if m, ok := result.(map[string]any); ok {
		// First try to get "result" field
		if resultField, exists := m["result"]; exists {
			// If result is a string, return it
			if resultStr, ok := resultField.(string); ok {
				return resultStr
			}
			// Otherwise recursively extract
			return extractResponseText(resultField)
		}

		// Try other common field names
		if response, exists := m["response"]; exists {
			if responseStr, ok := response.(string); ok {
				return responseStr
			}
			return extractResponseText(response)
		}

		if output, exists := m["output"]; exists {
			if outputStr, ok := output.(string); ok {
				return outputStr
			}
			return extractResponseText(output)
		}

		// If map has "description" field, it might be the task structure itself
		// Don't return the whole map structure
		if _, hasDesc := m["description"]; hasDesc {
			// This is likely the task metadata, not the actual response
			logger.WarnCF("collaborative", "Received task metadata instead of response", map[string]any{
				"map_keys": fmt.Sprintf("%v", getMapKeys(m)),
			})
			return ""
		}
	}

	// Fallback: only convert simple types to string, not complex structures
	switch v := result.(type) {
	case int, int64, float64, bool:
		return fmt.Sprintf("%v", v)
	default:
		// For complex types, log warning and return empty
		logger.WarnCF("collaborative", "Cannot extract text from complex result type", map[string]any{
			"type": fmt.Sprintf("%T", result),
		})
		return ""
	}
}

// getMapKeys returns the keys of a map for debugging
func getMapKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
