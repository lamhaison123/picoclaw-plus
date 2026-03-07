// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package collaborative

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/sipeed/picoclaw/pkg/logger"
)

// HandleWhoCommand generates the response for /who command
// Shows team status, registered agents, and active agents in the session
func HandleWhoCommand(manager *ManagerV2, chatID int64, platform Platform) string {
	var sb strings.Builder
	sb.WriteString("🤖 Team Status\n\n")

	// Get team information
	teamManager := platform.GetTeamManager()
	if teamManager != nil {
		// Try to get team info from the session first
		session, exists := manager.GetSession(chatID)
		if exists {
			if teamInfoAny, err := teamManager.GetTeam(session.TeamID); err == nil && teamInfoAny != nil {
				// Use reflection to extract team info
				teamName := session.TeamID // default
				var roles []struct {
					Name        string
					Description string
				}

				v := reflect.ValueOf(teamInfoAny)
				if v.Kind() == reflect.Ptr {
					v = v.Elem()
				}

				if v.Kind() == reflect.Struct {
					// Extract Name field
					if nameField := v.FieldByName("Name"); nameField.IsValid() && nameField.Kind() == reflect.String {
						teamName = nameField.String()
					}

					// Extract Config.Roles
					if configField := v.FieldByName("Config"); configField.IsValid() {
						// Check if field is a pointer and not nil
						if configField.Kind() == reflect.Ptr && !configField.IsNil() {
							configVal := configField.Elem()
							if rolesField := configVal.FieldByName("Roles"); rolesField.IsValid() {
								for i := 0; i < rolesField.Len(); i++ {
									roleVal := rolesField.Index(i)
									var roleName, roleDesc string
									if nameField := roleVal.FieldByName("Name"); nameField.IsValid() {
										roleName = nameField.String()
									}
									if descField := roleVal.FieldByName("Description"); descField.IsValid() {
										roleDesc = descField.String()
									}
									if roleName != "" {
										roles = append(roles, struct {
											Name        string
											Description string
										}{Name: roleName, Description: roleDesc})
									}
								}
							}
						}
					}
				}

				sb.WriteString(fmt.Sprintf("**Team:** %s\n", teamName))
				if len(roles) > 0 {
					sb.WriteString(fmt.Sprintf("**Registered Agents:** %d\n\n", len(roles)))
					sb.WriteString("📋 **Available Agents:**\n")
					for _, role := range roles {
						emoji := GetRoleEmoji(role.Name)
						sb.WriteString(fmt.Sprintf("  %s @%s - %s\n", emoji, role.Name, role.Description))
					}
					sb.WriteString("\n")
				}
			}
		}
	}

	// Show session info if exists
	session, exists := manager.GetSession(chatID)
	if !exists {
		sb.WriteString("💡 **No active session**\n")
		sb.WriteString("Mention an agent to start: @developer @tester\n")
	} else {
		sb.WriteString(fmt.Sprintf("**Session:** %s\n", GenerateSessionID(session.ChatID)))
		sb.WriteString(fmt.Sprintf("**Started:** %s\n\n",
			session.StartTime.Format("15:04:05")))

		activeAgents := session.GetActiveAgents()
		if len(activeAgents) == 0 {
			sb.WriteString("⚡ **Active Agents:** None yet\n")
			sb.WriteString("Mention an agent to activate them!\n")
		} else {
			sb.WriteString(fmt.Sprintf("⚡ **Active Agents:** %d\n", len(activeAgents)))
			for _, role := range activeAgents {
				emoji := GetRoleEmoji(role)
				status := session.GetAgentStatus(role)
				statusEmoji := GetStatusEmoji(status)

				// Get agent state for more details
				agent := session.GetAgentState(role)
				if agent == nil {
					logger.DebugCF("collaborative", "Agent state not found", map[string]any{
						"role":       role,
						"session_id": session.SessionID,
					})
					continue
				}

				timeStr := FormatTimeSince(agent.LastSeen)

				sb.WriteString(fmt.Sprintf("  %s @%s %s [%d msgs, %s]\n",
					emoji, role, statusEmoji, agent.MessageCount, timeStr))
			}

			sb.WriteString(fmt.Sprintf("\n📊 **Context:** %d messages | **Last activity:** %s",
				len(session.Context),
				session.LastActivity.Format("15:04:05")))
		}
	}

	return sb.String()
}

// HandleHelpCommand generates the response for /help command
// Shows usage instructions and available commands
func HandleHelpCommand(availableRoles []string) string {
	var sb strings.Builder
	sb.WriteString("🤖 Collaborative Chat Commands\n\n")
	sb.WriteString("📋 Available Commands:\n")
	sb.WriteString("  /who - Show team status and active agents\n")
	sb.WriteString("  /help - Show this help message\n\n")
	sb.WriteString("💬 How to Use:\n")
	sb.WriteString("  • Mention agents: @developer @tester\n")
	sb.WriteString("  • Agents can mention each other\n")
	sb.WriteString("  • Full conversation context is shared\n\n")
	sb.WriteString("🎯 Example:\n")
	sb.WriteString("  \"Hey @developer can you help @tester with this?\"\n\n")

	if len(availableRoles) > 0 {
		sb.WriteString("Available roles: ")
		for i, role := range availableRoles {
			if i > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(fmt.Sprintf("@%s", role))
		}
	} else {
		sb.WriteString("Available roles: @architect @developer @tester @manager")
	}

	return sb.String()
}

// IsCollaborativeCommand checks if a message is a collaborative command
func IsCollaborativeCommand(content string) bool {
	content = strings.TrimSpace(strings.ToLower(content))
	return content == "/who" || content == "/help"
}

// HandleCommand processes a collaborative command and returns the response
// Returns (response, handled) where handled indicates if it was a valid command
func HandleCommand(manager *ManagerV2, chatID int64, platform Platform, content string, availableRoles []string) (string, bool) {
	content = strings.TrimSpace(strings.ToLower(content))

	switch content {
	case "/who":
		return HandleWhoCommand(manager, chatID, platform), true
	case "/help":
		return HandleHelpCommand(availableRoles), true
	default:
		return "", false
	}
}
