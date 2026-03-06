package telegram

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/sipeed/picoclaw/pkg/logger"
	"github.com/sipeed/picoclaw/pkg/team"
)

// CollaborativeChatSession manages a multi-agent conversation in a Telegram chat
type CollaborativeChatSession struct {
	SessionID    string
	ChatID       int64
	TeamID       string
	ActiveAgents map[string]*AgentState
	Context      []ChatMessage
	StartTime    time.Time
	LastActivity time.Time
	MaxContext   int
	mu           sync.RWMutex
}

// AgentState tracks the state of an agent in the session
type AgentState struct {
	Role         string
	Status       string // idle, thinking, busy
	LastSeen     time.Time
	MessageCount int
}

// ChatMessage represents a message in the collaborative chat
type ChatMessage struct {
	ID        string
	Author    string // "user" or role name
	Content   string
	Timestamp time.Time
	Mentions  []string
}

var (
	mentionRegex = regexp.MustCompile(`@(\w+)`)
	roleEmojis   = map[string]string{
		"architect": "🏗️",
		"developer": "💻",
		"tester":    "🧪",
		"manager":   "📋",
		"designer":  "🎨",
		"devops":    "⚙️",
	}
)

// NewCollaborativeChatSession creates a new collaborative chat session
func NewCollaborativeChatSession(chatID int64, teamID string, maxContext int) *CollaborativeChatSession {
	if maxContext <= 0 {
		maxContext = 50 // default
	}

	return &CollaborativeChatSession{
		SessionID:    generateSessionID(chatID),
		ChatID:       chatID,
		TeamID:       teamID,
		ActiveAgents: make(map[string]*AgentState),
		Context:      make([]ChatMessage, 0, maxContext),
		StartTime:    time.Now(),
		LastActivity: time.Now(),
		MaxContext:   maxContext,
	}
}

// AddMessage adds a message to the conversation context
func (s *CollaborativeChatSession) AddMessage(author, content string, mentions []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	msg := ChatMessage{
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		Author:    author,
		Content:   content,
		Timestamp: time.Now(),
		Mentions:  mentions,
	}

	s.Context = append(s.Context, msg)
	s.LastActivity = time.Now()

	// Trim context if too long
	if len(s.Context) > s.MaxContext {
		s.Context = s.Context[len(s.Context)-s.MaxContext:]
	}

	logger.DebugCF("collaborative_chat", "Added message to context", map[string]any{
		"session_id":  s.SessionID,
		"author":      author,
		"mentions":    mentions,
		"context_len": len(s.Context),
	})
}

// GetFullContext returns the complete conversation history
func (s *CollaborativeChatSession) GetFullContext() []ChatMessage {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Return a copy to avoid race conditions
	contextCopy := make([]ChatMessage, len(s.Context))
	copy(contextCopy, s.Context)
	return contextCopy
}

// GetContextAsString formats the context as a string for LLM
func (s *CollaborativeChatSession) GetContextAsString() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var sb strings.Builder
	sb.WriteString("=== Collaborative Chat Context ===\n")
	sb.WriteString(fmt.Sprintf("Session: %s | Team: %s\n", s.SessionID, s.TeamID))
	sb.WriteString(fmt.Sprintf("Started: %s\n", s.StartTime.Format("15:04:05")))
	sb.WriteString("=== Conversation History ===\n\n")

	for _, msg := range s.Context {
		emoji := ""
		if msg.Author != "user" {
			emoji = getRoleEmoji(msg.Author) + " "
		}
		sb.WriteString(fmt.Sprintf("[%s] %s%s: %s\n",
			msg.Timestamp.Format("15:04:05"),
			emoji,
			strings.ToUpper(msg.Author),
			msg.Content,
		))
	}

	return sb.String()
}

// UpdateAgentStatus updates the status of an agent
func (s *CollaborativeChatSession) UpdateAgentStatus(role, status string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if agent, exists := s.ActiveAgents[role]; exists {
		agent.Status = status
		agent.LastSeen = time.Now()
	} else {
		s.ActiveAgents[role] = &AgentState{
			Role:     role,
			Status:   status,
			LastSeen: time.Now(),
		}
	}
}

// GetAgentStatus returns the status of an agent
func (s *CollaborativeChatSession) GetAgentStatus(role string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if agent, exists := s.ActiveAgents[role]; exists {
		return agent.Status
	}
	return "unknown"
}

// GetActiveAgents returns a list of active agents
func (s *CollaborativeChatSession) GetActiveAgents() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	agents := make([]string, 0, len(s.ActiveAgents))
	for role := range s.ActiveAgents {
		agents = append(agents, role)
	}
	return agents
}

// IncrementMessageCount increments the message count for an agent
func (s *CollaborativeChatSession) IncrementMessageCount(role string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if agent, exists := s.ActiveAgents[role]; exists {
		agent.MessageCount++
	}
}

// extractMentions extracts @mentions from a message
func extractMentions(text string) []string {
	matches := mentionRegex.FindAllStringSubmatch(text, -1)
	mentions := make([]string, 0, len(matches))
	seen := make(map[string]bool)

	for _, match := range matches {
		if len(match) > 1 {
			mention := strings.ToLower(match[1])
			if !seen[mention] {
				mentions = append(mentions, mention)
				seen[mention] = true
			}
		}
	}

	return mentions
}

// generateSessionID generates a unique session ID for a chat
func generateSessionID(chatID int64) string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("chat%d%d", chatID%10000, timestamp%10000)
}

// formatAgentMessage formats an agent's message in IRC style
func formatAgentMessage(role, content, sessionID string) string {
	emoji := getRoleEmoji(role)
	return fmt.Sprintf("[%s] %s %s: %s",
		sessionID,
		emoji,
		strings.ToUpper(role),
		content,
	)
}

// getRoleEmoji returns the emoji for a role
func getRoleEmoji(role string) string {
	if emoji, exists := roleEmojis[strings.ToLower(role)]; exists {
		return emoji
	}
	return "🤖" // default
}

// shouldAutoJoin determines if an agent should auto-join based on keywords
func shouldAutoJoin(role, message string, autoJoinRules []AutoJoinRule) bool {
	message = strings.ToLower(message)
	for _, rule := range autoJoinRules {
		if rule.Role == role {
			for _, keyword := range rule.Keywords {
				if strings.Contains(message, strings.ToLower(keyword)) {
					return true
				}
			}
		}
	}
	return false
}

// AutoJoinRule defines when an agent should automatically join a conversation
type AutoJoinRule struct {
	Role     string   `json:"role"`
	Keywords []string `json:"keywords"`
}

// CollaborativeChatManager manages all active collaborative chat sessions
type CollaborativeChatManager struct {
	sessions map[int64]*CollaborativeChatSession
	mu       sync.RWMutex
}

// NewCollaborativeChatManager creates a new manager
func NewCollaborativeChatManager() *CollaborativeChatManager {
	return &CollaborativeChatManager{
		sessions: make(map[int64]*CollaborativeChatSession),
	}
}

// GetOrCreateSession gets an existing session or creates a new one
func (m *CollaborativeChatManager) GetOrCreateSession(chatID int64, teamID string, maxContext int) *CollaborativeChatSession {
	m.mu.Lock()
	defer m.mu.Unlock()

	if session, exists := m.sessions[chatID]; exists {
		return session
	}

	session := NewCollaborativeChatSession(chatID, teamID, maxContext)
	m.sessions[chatID] = session

	logger.InfoCF("collaborative_chat", "Created new session", map[string]any{
		"session_id": session.SessionID,
		"chat_id":    chatID,
		"team_id":    teamID,
	})

	return session
}

// GetSession gets an existing session
func (m *CollaborativeChatManager) GetSession(chatID int64) (*CollaborativeChatSession, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	session, exists := m.sessions[chatID]
	return session, exists
}

// RemoveSession removes a session
func (m *CollaborativeChatManager) RemoveSession(chatID int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.sessions, chatID)

	logger.InfoCF("collaborative_chat", "Removed session", map[string]any{
		"chat_id": chatID,
	})
}

// ExecuteAgentWithContext executes an agent with full conversation context
func (s *CollaborativeChatSession) ExecuteAgentWithContext(
	ctx context.Context,
	teamManager *team.TeamManager,
	role string,
	userMessage string,
) (string, error) {
	// Update agent status
	s.UpdateAgentStatus(role, "thinking")
	defer s.UpdateAgentStatus(role, "idle")

	// Build prompt with full context
	contextStr := s.GetContextAsString()
	fullPrompt := fmt.Sprintf("%s\n\nNew message: %s\n\nYou are the %s. Respond to the conversation above.",
		contextStr,
		userMessage,
		role,
	)

	// Execute with team manager
	result, err := teamManager.ExecuteTaskWithRole(ctx, s.TeamID, fullPrompt, role)
	if err != nil {
		logger.ErrorCF("collaborative_chat", "Agent execution failed", map[string]any{
			"session_id": s.SessionID,
			"role":       role,
			"error":      err.Error(),
		})
		return "", err
	}

	// Extract response
	response := fmt.Sprintf("%v", result)

	// Increment message count
	s.IncrementMessageCount(role)

	return response, nil
}
