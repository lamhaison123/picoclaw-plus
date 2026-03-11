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
	"time"

	"github.com/sipeed/picoclaw/pkg/bus"
	"github.com/sipeed/picoclaw/pkg/logger"
)

// ManagerV2 manages all active collaborative chat sessions with improved mention handling
type ManagerV2 struct {
	sessions          map[int64]*Session
	mu                sync.RWMutex
	dispatchTracker   *DispatchTracker
	queueManager      *QueueManager      // Queue manager for rate limiting and retry
	compactionManager *CompactionManager // Compaction manager for context compression
	maxMentionDepth   int                // Maximum depth for cascading mentions (default: 5)
	config            *Config            // Configuration
}

// NewManagerV2 creates a new improved collaborative chat manager with default config
func NewManagerV2() *ManagerV2 {
	return NewManagerV2WithConfig(&Config{
		MentionQueueSize:    20,
		MentionRateLimit:    2 * time.Second,
		MentionMaxRetries:   3,
		MentionRetryBackoff: 1 * time.Second,
	})
}

// NewManagerV2WithConfig creates a manager with custom config
func NewManagerV2WithConfig(config *Config) *ManagerV2 {
	// Set defaults if not provided
	if config.MentionQueueSize == 0 {
		config.MentionQueueSize = 20
	}
	if config.MentionRateLimit == 0 {
		config.MentionRateLimit = 2 * time.Second
	}
	if config.MentionMaxRetries == 0 {
		config.MentionMaxRetries = 3
	}
	if config.MentionRetryBackoff == 0 {
		config.MentionRetryBackoff = 1 * time.Second
	}
	if config.MaxMentionDepth == 0 {
		config.MaxMentionDepth = 20 // Increased from 3 to 20 for more flexible workflows
	}

	// Initialize compaction manager if enabled
	var compactionMgr *CompactionManager
	if config.CompactionEnabled && config.LLMProvider != nil {
		// Set compaction config defaults
		if config.CompactionConfig.TriggerThreshold == 0 {
			config.CompactionConfig.TriggerThreshold = 40
		}
		if config.CompactionConfig.KeepRecentCount == 0 {
			config.CompactionConfig.KeepRecentCount = 15
		}
		if config.CompactionConfig.CompactBatchSize == 0 {
			config.CompactionConfig.CompactBatchSize = 25
		}
		if config.CompactionConfig.MinInterval == 0 {
			config.CompactionConfig.MinInterval = 5 * time.Minute
		}
		if config.CompactionConfig.SummaryMaxLength == 0 {
			config.CompactionConfig.SummaryMaxLength = 2000
		}
		if config.CompactionConfig.LLMModel == "" {
			config.CompactionConfig.LLMModel = "gpt-4o-mini"
		}
		if config.CompactionConfig.LLMTimeout == 0 {
			config.CompactionConfig.LLMTimeout = 90 * time.Second // Increased from 30s to 90s for summarization
		}
		if config.CompactionConfig.LLMMaxRetries == 0 {
			config.CompactionConfig.LLMMaxRetries = 3
		}

		config.CompactionConfig.Enabled = true

		// Create summarizer
		summarizer := NewLLMSummarizer(config.CompactionConfig, config.LLMProvider)
		compactionMgr = NewCompactionManager(config.CompactionConfig, summarizer)

		logger.InfoCF("collaborative", "Compaction manager initialized", map[string]any{
			"trigger_threshold": config.CompactionConfig.TriggerThreshold,
			"keep_recent":       config.CompactionConfig.KeepRecentCount,
			"llm_model":         config.CompactionConfig.LLMModel,
		})
	}

	return &ManagerV2{
		sessions:          make(map[int64]*Session),
		dispatchTracker:   NewDispatchTracker(),
		queueManager:      NewQueueManager(config.MentionQueueSize, config.MentionRateLimit, config.MentionMaxRetries, config.MentionRetryBackoff),
		compactionManager: compactionMgr,
		maxMentionDepth:   config.MaxMentionDepth,
		config:            config,
	}
}

// HandleMentions processes mentions and triggers agents with idempotency
func (m *ManagerV2) HandleMentions(
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

	logger.InfoCF("collaborative", "Processing mentions", map[string]any{
		"chat_id":      chatID,
		"team_id":      teamID,
		"mentions":     mentions,
		"content_len":  len(content),
		"sender":       sender.Username,
		"content_utf8": len([]rune(content)), // UTF-8 character count
	})

	// Get or create session
	session := m.GetOrCreateSession(chatID, teamID, maxContext)

	// BUG FIX: Check if session creation failed
	if session == nil {
		logger.ErrorCF("collaborative", "Failed to create session", map[string]any{
			"chat_id": chatID,
			"team_id": teamID,
		})
		return fmt.Errorf("failed to create session for chat %d", chatID)
	}

	// Check mention depth to prevent infinite loops
	if session.MentionDepth >= m.maxMentionDepth {
		logger.WarnCF("collaborative", "Mention depth limit reached", map[string]any{
			"chat_id":       chatID,
			"session_id":    session.SessionID,
			"current_depth": session.MentionDepth,
			"max_depth":     m.maxMentionDepth,
		})
		return fmt.Errorf("mention depth limit reached (%d)", m.maxMentionDepth)
	}

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

	// Execute each mentioned role via queue
	for _, role := range mentions {
		req := &MentionRequest{
			Role:        role,
			Prompt:      content,
			SessionID:   session.SessionID,
			ChatID:      chatID,
			TeamID:      teamID,
			Timestamp:   time.Now(),
			Context:     ctx,
			Platform:    platform,
			Session:     session,
			TeamRoster:  teamRoster,
			Depth:       0,
			ExecuteFunc: m.executeMentionRequest,
		}

		err := m.queueManager.Enqueue(req)
		if err != nil {
			logger.ErrorCF("collaborative", "Failed to enqueue mention", map[string]any{
				"role":       role,
				"chat_id":    chatID,
				"session_id": session.SessionID,
				"error":      err.Error(),
			})

			// Send error message to user
			platform.SendMessage(ctx, fmt.Sprintf("%d", chatID),
				fmt.Sprintf("⚠️ Queue full for @%s, please try again later", role))
		}
	}

	return nil
}

// executeAgentAndCascadeWithError handles agent execution and returns error for retry
func (m *ManagerV2) executeAgentAndCascadeWithError(
	platform Platform,
	chatID int64,
	teamID string,
	session *Session,
	role string,
	triggerContent string,
	teamRoster string,
	currentDepth int,
	mentionedBy string, // NEW: Who mentioned this role (for ack-loop prevention)
) error {
	// Check if we've reached max depth before cascading
	if currentDepth >= m.maxMentionDepth {
		logger.WarnCF("collaborative", "Max mention depth reached, skipping cascade", map[string]any{
			"role":          role,
			"session_id":    session.SessionID,
			"current_depth": currentDepth,
			"max_depth":     m.maxMentionDepth,
		})
		return nil
	}

	// Generate unique message ID for idempotency
	messageID := GenerateMessageID(chatID, session.SessionID, role, triggerContent)

	// Check and mark as dispatched atomically to prevent race conditions
	if !m.dispatchTracker.TryMarkDispatched(messageID) {
		logger.WarnCF("collaborative", "Skipping duplicate mention dispatch", map[string]any{
			"message_id": messageID,
			"role":       role,
			"chat_id":    chatID,
			"session_id": session.SessionID,
		})
		return nil
	}

	logger.InfoCF("collaborative", "Dispatching mention to agent", map[string]any{
		"message_id":          messageID,
		"role":                role,
		"chat_id":             chatID,
		"session_id":          session.SessionID,
		"dispatch_queue_size": m.dispatchTracker.Size(),
		"depth":               currentDepth,
	})

	// BUG FIX #8: Check context before starting cascade
	if platform.GetContext().Err() != nil {
		logger.WarnCF("collaborative", "Context cancelled before cascade", map[string]any{
			"role":       role,
			"session_id": session.SessionID,
		})
		return platform.GetContext().Err()
	}

	// Mark agent as in cascade
	session.IncrementMentionDepth()
	defer session.DecrementMentionDepth()
	session.MarkAgentInCascade(role)

	// Update agent status
	session.UpdateAgentStatus(role, "thinking")

	// Build context
	contextStr := session.GetContextAsString()

	var prompt string
	if currentDepth == 0 {
		prompt = fmt.Sprintf(`%s

=== Team Information ===
%s

User message: %s

You are @%s. Respond to the user's message considering the conversation history above.
You can mention other team members using @role format (e.g., @architect, @developer).`,
			contextStr, teamRoster, triggerContent, role)
	} else {
		// Build ack-loop prevention warning
		ackWarning := ""
		if mentionedBy != "" {
			ackWarning = fmt.Sprintf("\n\nIMPORTANT: Do NOT mention @%s back immediately (they just mentioned you). Only mention them if you genuinely need their input for a NEW task.", mentionedBy)
		}

		prompt = fmt.Sprintf(`%s

=== Team Information ===
%s

You were mentioned in the conversation. Please respond if needed.
Trigger context: %s

You are @%s. You can mention other team members using @role format.%s`,
			contextStr, teamRoster, triggerContent, role, ackWarning)
	}

	logger.DebugCF("collaborative", "Starting agent execution", map[string]any{
		"message_id": messageID,
		"role":       role,
		"chat_id":    chatID,
		"depth":      currentDepth,
	})

	teamManager := platform.GetTeamManager()
	if teamManager == nil {
		// Unmark agent from cascade if team manager not available
		session.UnmarkAgentInCascade(role)
		logger.ErrorCF("collaborative", "Team manager not available", map[string]any{
			"role":       role,
			"session_id": session.SessionID,
		})
		return fmt.Errorf("team manager not available")
	}

	// Use platform's context
	result, err := teamManager.ExecuteTaskWithRole(platform.GetContext(), teamID, prompt, role)

	if err != nil {
		logger.ErrorCF("collaborative", "Agent execution failed", map[string]any{
			"message_id": messageID,
			"role":       role,
			"chat_id":    chatID,
			"session_id": session.SessionID,
			"error":      err.Error(),
		})
		session.UpdateAgentStatus(role, "error")

		// Unmark agent from cascade on error
		session.UnmarkAgentInCascade(role)

		// Send user-friendly error message
		var errorMsg string
		if strings.Contains(err.Error(), "not found in team configuration") {
			errorMsg = fmt.Sprintf("❌ Role @%s không tồn tại trong team configuration.\n\n💡 Các role có sẵn: @architect, @developer, @tester, @manager", role)
		} else {
			errorMsg = fmt.Sprintf("❌ @%s encountered an error: %v", role, err)
		}
		platform.SendMessage(platform.GetContext(), fmt.Sprintf("%d", chatID), errorMsg)
		return err
	}

	// Extract response text from result
	// Note: We avoid importing collaborative/manager.go directly; extractResponseText should be available in package scope
	if result == nil {
		logger.ErrorCF("collaborative", "Nil result from agent execution", map[string]any{
			"message_id": messageID,
			"role":       role,
			"session_id": session.SessionID,
		})
		session.UpdateAgentStatus(role, "error")
		session.UnmarkAgentInCascade(role)
		return fmt.Errorf("nil result from agent")
	}

	responseStr := extractResponseText(result)

	logger.InfoCF("collaborative", "Agent execution completed", map[string]any{
		"message_id":    messageID,
		"role":          role,
		"chat_id":       chatID,
		"session_id":    session.SessionID,
		"response_len":  len(responseStr),
		"response_utf8": len([]rune(responseStr)),
	})

	session.AddMessage(role, responseStr, nil)
	session.UpdateAgentStatus(role, "idle")

	// Format with IRC-style prefix
	formattedMsg := FormatMessage(session.SessionID, role, responseStr)

	// Send to platform
	err = platform.SendMessage(platform.GetContext(), fmt.Sprintf("%d", chatID), formattedMsg)
	if err != nil {
		logger.ErrorCF("collaborative", "Failed to send agent response", map[string]any{
			"message_id": messageID,
			"role":       role,
			"chat_id":    chatID,
			"error":      err.Error(),
		})
		// Unmark agent from cascade even if send fails
		session.UnmarkAgentInCascade(role)
		return err
	}

	logger.InfoCF("collaborative", "Agent response sent successfully", map[string]any{
		"message_id": messageID,
		"role":       role,
		"session_id": session.SessionID,
		"chat_id":    chatID,
	})

	// Trigger compaction if enabled and manager is available
	if m.compactionManager != nil {
		m.compactionManager.CompactAsync(platform.GetContext(), session)
	}

	// Unmark agent from cascade BEFORE checking for new mentions
	// This allows other agents to mention this agent again if needed
	session.UnmarkAgentInCascade(role)

	// Check if agent mentioned other agents in their response
	mentionsInResponse := ExtractMentionsImproved(responseStr)
	if len(mentionsInResponse) > 0 {
		// Filter out self-mentions AND ack-mentions (ack-loop prevention)
		newMentions := []string{}
		for _, mentioned := range mentionsInResponse {
			// Don't mention self
			if mentioned == role {
				continue
			}

			// Don't mention back the person who just mentioned you (ack-loop prevention)
			// This prevents: A mentions B → B mentions A → A mentions B → ...
			if mentionedBy != "" && mentioned == mentionedBy {
				logger.InfoCF("collaborative", "Filtered ack-mention to prevent loop", map[string]any{
					"from_role":  role,
					"to_role":    mentioned,
					"session_id": session.SessionID,
					"depth":      currentDepth,
				})
				continue
			}

			newMentions = append(newMentions, mentioned)
		}

		if len(newMentions) > 0 {
			// Check if we've reached max depth before cascading
			if currentDepth+1 >= m.maxMentionDepth {
				logger.WarnCF("collaborative", "Max mention depth reached, skipping cascade", map[string]any{
					"from_role":     role,
					"mentions":      newMentions,
					"session_id":    session.SessionID,
					"current_depth": currentDepth,
					"max_depth":     m.maxMentionDepth,
				})
				return nil
			}

			logger.InfoCF("collaborative", "Agent mentioned other agents", map[string]any{
				"from_role":     role,
				"mentions":      newMentions,
				"session_id":    session.SessionID,
				"current_depth": currentDepth,
			})

			for _, mentionedRole := range newMentions {
				// Enqueue cascaded mention with mentionedBy for ack-loop prevention
				req := &MentionRequest{
					Role:        mentionedRole,
					Prompt:      responseStr,
					SessionID:   session.SessionID,
					ChatID:      chatID,
					TeamID:      teamID,
					Timestamp:   time.Now(),
					Context:     platform.GetContext(),
					Platform:    platform,
					Session:     session,
					TeamRoster:  teamRoster,
					Depth:       currentDepth + 1,
					MentionedBy: role, // NEW: Track who mentioned this role
					ExecuteFunc: m.executeMentionRequest,
				}

				err := m.queueManager.Enqueue(req)
				if err != nil {
					logger.ErrorCF("collaborative", "Failed to enqueue cascaded mention", map[string]any{
						"from_role":   role,
						"target_role": mentionedRole,
						"session_id":  session.SessionID,
						"error":       err.Error(),
					})
				}
			}
		}
	}

	return nil
}

// GetOrCreateSession gets an existing session or creates a new one
func (m *ManagerV2) GetOrCreateSession(chatID int64, teamID string, maxContext int) *Session {
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
func (m *ManagerV2) GetSession(chatID int64) (*Session, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	session, exists := m.sessions[chatID]
	return session, exists
}

// RemoveSession removes a session
func (m *ManagerV2) RemoveSession(chatID int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.sessions, chatID)

	logger.InfoCF("collaborative", "Removed session", map[string]any{
		"chat_id": chatID,
	})
}

// CleanupDispatchTracker clears old dispatch records (call periodically)
func (m *ManagerV2) CleanupDispatchTracker() {
	oldSize := m.dispatchTracker.Size()
	m.dispatchTracker.Clear()

	logger.InfoCF("collaborative", "Cleaned up dispatch tracker", map[string]any{
		"old_size": oldSize,
		"new_size": m.dispatchTracker.Size(),
	})
}

// executeMentionRequest is the callback function for queue execution
func (m *ManagerV2) executeMentionRequest(req *MentionRequest) error {
	// Check depth limit before execution
	if req.Depth >= m.maxMentionDepth {
		logger.WarnCF("collaborative", "Max mention depth reached, skipping execution", map[string]any{
			"role":          req.Role,
			"session_id":    req.SessionID,
			"current_depth": req.Depth,
			"max_depth":     m.maxMentionDepth,
		})
		return nil // Not an error, just skip
	}

	// Execute and capture any errors for retry
	return m.executeAgentAndCascadeWithError(
		req.Platform,
		req.ChatID,
		req.TeamID,
		req.Session,
		req.Role,
		req.Prompt,
		req.TeamRoster,
		req.Depth,
		req.MentionedBy, // Pass mentionedBy for ack-loop prevention
	)
}

// GetQueueMetrics returns metrics for a specific role
func (m *ManagerV2) GetQueueMetrics(role string) *QueueMetrics {
	return m.queueManager.GetMetrics(role)
}

// GetAllQueueMetrics returns metrics for all roles
func (m *ManagerV2) GetAllQueueMetrics() map[string]QueueMetrics {
	return m.queueManager.GetAllMetrics()
}

// Stop gracefully stops the manager and all queues
func (m *ManagerV2) Stop() {
	logger.InfoC("collaborative", "Stopping ManagerV2")
	m.queueManager.Stop()

	if m.compactionManager != nil {
		logger.InfoC("collaborative", "Stopping compaction manager")
		m.compactionManager.Stop()
	}
}

// GetCompactionManager returns the compaction manager (for testing)
func (m *ManagerV2) GetCompactionManager() *CompactionManager {
	return m.compactionManager
}

// GetEnhancedMetrics returns the enhanced metrics collector
func (m *ManagerV2) GetEnhancedMetrics() *EnhancedMetrics {
	if m.dispatchTracker != nil {
		return m.dispatchTracker.GetMetrics()
	}
	return nil
}
