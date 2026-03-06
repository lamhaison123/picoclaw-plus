// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package collaborative

import (
	"fmt"
	"strings"
	"time"

	"github.com/sipeed/picoclaw/pkg/logger"
)

// NewSession creates a new collaborative chat session
func NewSession(chatID int64, teamID string, maxContext int) *Session {
	if maxContext <= 0 {
		maxContext = 50 // default
	}

	return &Session{
		SessionID:    generateSessionID(chatID),
		ChatID:       chatID,
		TeamID:       teamID,
		ActiveAgents: make(map[string]*AgentState),
		Context:      make([]Message, 0, maxContext),
		StartTime:    time.Now(),
		LastActivity: time.Now(),
		MaxContext:   maxContext,
	}
}

// AddMessage adds a message to the conversation context
func (s *Session) AddMessage(author, content string, mentions []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	msg := Message{
		Role:      author,
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

	logger.DebugCF("collaborative", "Added message to context", map[string]any{
		"session_id":  s.SessionID,
		"author":      author,
		"mentions":    mentions,
		"context_len": len(s.Context),
	})
}

// GetFullContext returns the complete conversation history
func (s *Session) GetFullContext() []Message {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Return a copy to avoid race conditions
	contextCopy := make([]Message, len(s.Context))
	copy(contextCopy, s.Context)
	return contextCopy
}

// GetContextAsString formats the context as a string for LLM
func (s *Session) GetContextAsString() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var sb strings.Builder
	sb.WriteString("=== Collaborative Chat Context ===\n")
	sb.WriteString(fmt.Sprintf("Session: %s | Team: %s\n", s.SessionID, s.TeamID))
	sb.WriteString(fmt.Sprintf("Started: %s\n", s.StartTime.Format("15:04:05")))
	sb.WriteString("=== Conversation History ===\n\n")

	for _, msg := range s.Context {
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

// UpdateAgentStatus updates the status of an agent
func (s *Session) UpdateAgentStatus(role, status string) {
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
func (s *Session) GetAgentStatus(role string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if agent, exists := s.ActiveAgents[role]; exists {
		return agent.Status
	}
	return "unknown"
}

// GetActiveAgents returns a list of active agents
func (s *Session) GetActiveAgents() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	agents := make([]string, 0, len(s.ActiveAgents))
	for role := range s.ActiveAgents {
		agents = append(agents, role)
	}
	return agents
}

// IncrementMessageCount increments the message count for an agent
func (s *Session) IncrementMessageCount(role string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if agent, exists := s.ActiveAgents[role]; exists {
		agent.MessageCount++
	}
}

// generateSessionID generates a unique session ID for a chat
func generateSessionID(chatID int64) string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("chat%d%d", chatID%10000, timestamp%10000)
}

// GetAgentState returns the agent state for a role
func (s *Session) GetAgentState(role string) *AgentState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if agent, exists := s.ActiveAgents[role]; exists {
		// Return a copy to avoid race conditions
		return &AgentState{
			Role:         agent.Role,
			Status:       agent.Status,
			MessageCount: agent.MessageCount,
			LastSeen:     agent.LastSeen,
		}
	}
	return nil
}
