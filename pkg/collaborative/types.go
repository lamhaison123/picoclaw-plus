// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package collaborative

import (
	"context"
	"sync"
	"time"

	"github.com/sipeed/picoclaw/pkg/bus"
)

// Session represents a collaborative chat session
type Session struct {
	SessionID    string
	TeamID       string
	ChatID       int64
	StartTime    time.Time
	LastActivity time.Time
	Context      []Message
	ActiveAgents map[string]*AgentState
	MaxContext   int
	mu           sync.RWMutex
}

// Message represents a message in the conversation
type Message struct {
	Role      string
	Content   string
	Timestamp time.Time
	Mentions  []string
}

// AgentState tracks agent status and activity
type AgentState struct {
	Role         string
	Status       string // "idle", "thinking", "busy", "error"
	MessageCount int
	LastSeen     time.Time
}

// SenderInfo contains information about the message sender
type SenderInfo struct {
	Platform    string
	PlatformID  string
	CanonicalID string
	Username    string
	DisplayName string
}

// Platform interface for platform-specific operations
type Platform interface {
	// SendMessage sends a message to the specified chat
	SendMessage(ctx context.Context, chatID string, content string) error

	// GetTeamManager returns the team manager for executing agent tasks
	GetTeamManager() TeamManager

	// GetContext returns the platform's context
	GetContext() context.Context
}

// TeamManager interface for team operations
type TeamManager interface {
	// ExecuteTaskWithRole executes a task with a specific role
	ExecuteTaskWithRole(ctx context.Context, teamID, prompt, role string) (any, error)

	// GetTeam returns team information
	GetTeam(teamID string) (any, error)
}

// Config holds configuration for collaborative chat
type Config struct {
	Enabled          bool
	DefaultTeamID    string
	MaxContextLength int
}

// MessageHandler handles collaborative messages
type MessageHandler interface {
	// HandleMentions processes mentions and triggers agents
	HandleMentions(
		ctx context.Context,
		platform Platform,
		chatID int64,
		teamID string,
		content string,
		mentions []string,
		sender bus.SenderInfo,
	) error

	// HandleCommand processes collaborative commands
	HandleCommand(
		ctx context.Context,
		platform Platform,
		chatID int64,
		command string,
	) bool
}
