// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package collaborative

import (
	"fmt"
	"strings"
	"time"
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

// FormatAgentMessage formats an agent's message in IRC style
// Format: [sessionID] emoji ROLE: content
func FormatAgentMessage(role, content, sessionID string) string {
	emoji := GetRoleEmoji(role)
	return fmt.Sprintf("[%s] %s %s: %s",
		sessionID,
		emoji,
		strings.ToUpper(role),
		content,
	)
}

// FormatSessionContext formats the session context as a string for LLM
func FormatSessionContext(session *Session) string {
	if session == nil {
		return "=== No Session ===\n"
	}

	var sb strings.Builder
	sb.WriteString("=== Collaborative Chat Context ===\n")
	sb.WriteString(fmt.Sprintf("Session: chat%d%d | Team: %s\n",
		session.ChatID%10000,
		session.StartTime.Unix()%10000,
		session.TeamID))
	sb.WriteString(fmt.Sprintf("Started: %s\n", session.StartTime.Format("15:04:05")))
	sb.WriteString("=== Conversation History ===\n\n")

	for _, msg := range session.Context {
		emoji := ""
		if msg.Role != "user" {
			emoji = GetRoleEmoji(msg.Role) + " "
		}
		sb.WriteString(fmt.Sprintf("[%s] %s%s: %s\n",
			msg.Timestamp.Format("15:04:05"),
			emoji,
			strings.ToUpper(msg.Role),
			msg.Content,
		))
	}

	return sb.String()
}

// GenerateSessionID generates a unique session ID for a chat
func GenerateSessionID(chatID int64) string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("chat%d%d", chatID%10000, timestamp%10000)
}

// FormatTimeSince formats a duration as a human-readable string
func FormatTimeSince(t time.Time) string {
	timeSince := time.Since(t)
	if timeSince < time.Minute {
		return "just now"
	} else if timeSince < time.Hour {
		return fmt.Sprintf("%dm ago", int(timeSince.Minutes()))
	} else {
		return fmt.Sprintf("%dh ago", int(timeSince.Hours()))
	}
}

// GetStatusEmoji returns the emoji for an agent status
func GetStatusEmoji(status string) string {
	switch status {
	case "thinking":
		return "🤔"
	case "busy":
		return "⚡"
	case "idle":
		return "✅"
	case "error":
		return "❌"
	default:
		return "💤"
	}
}
