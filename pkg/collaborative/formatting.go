// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package collaborative

import (
	"fmt"
	"strings"
)

var roleEmojis = map[string]string{
	"architect": "🏗️",
	"developer": "💻",
	"tester":    "🧪",
	"manager":   "📋",
	"designer":  "🎨",
	"devops":    "⚙️",
}

// FormatMessage formats a message without session ID prefix
func FormatMessage(sessionID, role, content string) string {
	emoji := GetRoleEmoji(role)
	return fmt.Sprintf("%s %s: %s",
		emoji, strings.ToUpper(role), content)
}

// GetRoleEmoji returns emoji for a role
func GetRoleEmoji(role string) string {
	if emoji, ok := roleEmojis[strings.ToLower(role)]; ok {
		return emoji
	}
	return "🤖" // default
}
