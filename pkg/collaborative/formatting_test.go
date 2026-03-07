package collaborative

import (
	"strings"
	"testing"
	"time"
)

func TestGetRoleEmoji(t *testing.T) {
	tests := []struct {
		name     string
		role     string
		expected string
	}{
		{"architect", "architect", "🏗️"},
		{"developer", "developer", "💻"},
		{"tester", "tester", "🧪"},
		{"manager", "manager", "📋"},
		{"designer", "designer", "🎨"},
		{"devops", "devops", "⚙️"},
		{"unknown", "unknown", "🤖"},
		{"uppercase", "DEVELOPER", "💻"},
		{"mixed case", "DevOps", "⚙️"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetRoleEmoji(tt.role)
			if result != tt.expected {
				t.Errorf("GetRoleEmoji(%s) = %s, want %s", tt.role, result, tt.expected)
			}
		})
	}
}

func TestFormatAgentMessage(t *testing.T) {
	tests := []struct {
		name      string
		role      string
		content   string
		sessionID string
		want      string
	}{
		{
			name:      "developer message",
			role:      "developer",
			content:   "Task completed",
			sessionID: "chat12345678",
			want:      "[chat12345678] 💻 DEVELOPER: Task completed",
		},
		{
			name:      "architect message",
			role:      "architect",
			content:   "Design approved",
			sessionID: "chat87654321",
			want:      "[chat87654321] 🏗️ ARCHITECT: Design approved",
		},
		{
			name:      "unknown role",
			role:      "unknown",
			content:   "Hello",
			sessionID: "chat11111111",
			want:      "[chat11111111] 🤖 UNKNOWN: Hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatAgentMessage(tt.role, tt.content, tt.sessionID)
			if result != tt.want {
				t.Errorf("FormatAgentMessage() = %s, want %s", result, tt.want)
			}
		})
	}
}

func TestGenerateSessionID(t *testing.T) {
	chatID := int64(123456789)
	sessionID := GenerateSessionID(chatID)

	// Check format: chat{chatID%10000}{timestamp%10000}
	if !strings.HasPrefix(sessionID, "chat") {
		t.Errorf("SessionID should start with 'chat', got %s", sessionID)
	}

	if len(sessionID) < 9 { // "chat" + at least 5 digits
		t.Errorf("SessionID too short: %s", sessionID)
	}

	// Generate another one with different chatID to verify format
	sessionID2 := GenerateSessionID(chatID + 1)
	if sessionID == sessionID2 {
		// This is actually OK if timestamps are same, just verify format is correct
		t.Logf("SessionIDs: %s, %s", sessionID, sessionID2)
	}
}

func TestFormatTimeSince(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		time     time.Time
		expected string
	}{
		{"just now", now.Add(-30 * time.Second), "just now"},
		{"5 minutes ago", now.Add(-5 * time.Minute), "5m ago"},
		{"30 minutes ago", now.Add(-30 * time.Minute), "30m ago"},
		{"2 hours ago", now.Add(-2 * time.Hour), "2h ago"},
		{"24 hours ago", now.Add(-24 * time.Hour), "24h ago"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatTimeSince(tt.time)
			if result != tt.expected {
				t.Errorf("FormatTimeSince() = %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestGetStatusEmoji(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected string
	}{
		{"thinking", "thinking", "🤔"},
		{"busy", "busy", "⚡"},
		{"idle", "idle", "✅"},
		{"error", "error", "❌"},
		{"unknown", "unknown", "💤"},
		{"empty", "", "💤"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetStatusEmoji(tt.status)
			if result != tt.expected {
				t.Errorf("GetStatusEmoji(%s) = %s, want %s", tt.status, result, tt.expected)
			}
		})
	}
}

func TestFormatSessionContext(t *testing.T) {
	session := &Session{
		ChatID:    123456,
		TeamID:    "dev-team",
		StartTime: time.Date(2026, 3, 7, 10, 30, 0, 0, time.UTC),
		Context: []Message{
			{
				Role:      "user",
				Content:   "Hello @developer",
				Timestamp: time.Date(2026, 3, 7, 10, 30, 5, 0, time.UTC),
				Mentions:  []string{"developer"},
			},
			{
				Role:      "developer",
				Content:   "Hi! How can I help?",
				Timestamp: time.Date(2026, 3, 7, 10, 30, 10, 0, time.UTC),
				Mentions:  []string{},
			},
		},
	}

	result := FormatSessionContext(session)

	// Check key components
	if !strings.Contains(result, "Collaborative Chat Context") {
		t.Error("Should contain 'Collaborative Chat Context'")
	}
	if !strings.Contains(result, "dev-team") {
		t.Error("Should contain team ID")
	}
	if !strings.Contains(result, "USER: Hello @developer") {
		t.Error("Should contain user message")
	}
	if !strings.Contains(result, "DEVELOPER: Hi! How can I help?") {
		t.Error("Should contain developer message")
	}
	if !strings.Contains(result, "💻") {
		t.Error("Should contain developer emoji")
	}
}
