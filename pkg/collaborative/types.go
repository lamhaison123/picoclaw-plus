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
	"github.com/sipeed/picoclaw/pkg/providers"
)

// Session represents a collaborative chat session
type Session struct {
	SessionID       string
	TeamID          string
	ChatID          int64
	StartTime       time.Time
	LastActivity    time.Time
	Context         []Message
	ActiveAgents    map[string]*AgentState
	MaxContext      int
	MentionDepth    int             // Track nested mention depth to prevent infinite loops
	MaxMentionDepth int             // Maximum allowed mention depth (default: 5)
	CascadeAgents   map[string]bool // Track agents currently in cascade to detect cycles
	mu              sync.RWMutex

	// Compaction fields
	CompactedContext *CompactedContext
	CompactionConfig CompactionConfig
	LastCompaction   time.Time
	CompactionMutex  sync.Mutex // Prevent concurrent compaction
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
	Enabled             bool
	DefaultTeamID       string
	MaxContextLength    int
	MentionQueueSize    int           // Queue size per role (default: 20)
	MentionRateLimit    time.Duration // Rate limit per role (default: 2s)
	MentionMaxRetries   int           // Max retry attempts (default: 3)
	MentionRetryBackoff time.Duration // Initial backoff duration (default: 1s)
	MaxMentionDepth     int           // Maximum depth for cascading mentions (default: 5)

	// Compaction config
	CompactionEnabled bool                  // Enable/disable compaction (default: false)
	CompactionConfig  CompactionConfig      // Compaction configuration
	LLMProvider       providers.LLMProvider // LLM provider for summarization (optional)
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
