// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package collaborative

import "fmt"

// BuildTeamRoster creates a formatted string with team member information
func BuildTeamRoster(teamInfo any) string {
	// Try to extract team information from the interface
	if m, ok := teamInfo.(map[string]any); ok {
		if config, exists := m["Config"]; exists {
			if configMap, ok := config.(map[string]any); ok {
				if roles, exists := configMap["roles"]; exists {
					if rolesList, ok := roles.([]any); ok {
						roster := "Team: " + fmt.Sprintf("%v", m["Name"]) + "\n"
						roster += "Members:\n"
						for _, role := range rolesList {
							if roleMap, ok := role.(map[string]any); ok {
								roleName := roleMap["name"]
								roleDesc := roleMap["description"]
								roster += fmt.Sprintf("  • @%s - %s\n", roleName, roleDesc)
							}
						}
						return roster
					}
				}
			}
		}
	}

	// Ultimate fallback - just list common roles
	return `Team: Development Team
Members:
  • @architect - System architect and designer
  • @developer - Software developer
  • @tester - QA and testing specialist
  • @manager - Project manager`
}

// BuildTeamRosterFromTeam creates a formatted string from a Team struct
// This is a helper for when you have a concrete Team type
func BuildTeamRosterFromTeam(team interface{}) string {
	// This function accepts any type and delegates to BuildTeamRoster
	// which handles the type conversion
	return BuildTeamRoster(team)
}
