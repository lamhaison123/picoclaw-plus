// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package collaborative

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sipeed/picoclaw/pkg/logger"
)

// CompactionManager manages context compaction
type CompactionManager struct {
	config     CompactionConfig
	metrics    *CompactionMetrics
	summarizer Summarizer
	mu         sync.RWMutex
	stopCh     chan struct{}
	wg         sync.WaitGroup
}

// Summarizer interface for generating summaries
type Summarizer interface {
	Summarize(ctx context.Context, req *CompactionRequest) (*CompactionResult, error)
}

// NewCompactionManager creates a new compaction manager
func NewCompactionManager(config CompactionConfig, summarizer Summarizer) *CompactionManager {
	return &CompactionManager{
		config:     config,
		metrics:    &CompactionMetrics{},
		summarizer: summarizer,
		stopCh:     make(chan struct{}),
	}
}

// ShouldCompact checks if compaction should be triggered
func (cm *CompactionManager) ShouldCompact(session *Session) bool {
	if !cm.config.Enabled {
		return false
	}

	// Check message count threshold
	session.mu.RLock()
	contextLen := len(session.Context)
	session.mu.RUnlock()

	if contextLen < cm.config.TriggerThreshold {
		return false
	}

	// Check minimum interval
	session.mu.RLock()
	lastCompaction := session.LastCompaction
	session.mu.RUnlock()

	if time.Since(lastCompaction) < cm.config.MinInterval {
		return false
	}

	// Check if already compacting (non-blocking check)
	if !session.CompactionMutex.TryLock() {
		return false
	}
	session.CompactionMutex.Unlock()

	return true
}

// CompactAsync triggers async compaction
func (cm *CompactionManager) CompactAsync(ctx context.Context, session *Session) {
	cm.wg.Add(1)
	go func() {
		defer cm.wg.Done()

		if err := cm.Compact(ctx, session); err != nil {
			logger.ErrorCF("compaction", "Compaction failed", map[string]any{
				"session_id": session.SessionID,
				"error":      err.Error(),
			})
		}
	}()
}

// Compact performs the actual compaction
func (cm *CompactionManager) Compact(ctx context.Context, session *Session) error {
	// Lock to prevent concurrent compaction
	session.CompactionMutex.Lock()
	defer session.CompactionMutex.Unlock()

	startTime := time.Now()

	logger.InfoCF("compaction", "Starting compaction", map[string]any{
		"session_id":        session.SessionID,
		"message_count":     len(session.Context),
		"trigger_threshold": cm.config.TriggerThreshold,
	})

	// Extract messages to compact
	messagesToCompact := cm.extractMessagesToCompact(session)
	if len(messagesToCompact) == 0 {
		return fmt.Errorf("no messages to compact")
	}

	// Build compaction request
	req := &CompactionRequest{
		SessionID:       session.SessionID,
		Messages:        messagesToCompact,
		ExistingSummary: cm.getExistingSummary(session),
		Config:          cm.config,
		Timestamp:       time.Now(),
	}

	// Generate summary
	result, err := cm.summarizer.Summarize(ctx, req)
	if err != nil {
		cm.updateMetrics(false, 0, 0, time.Since(startTime))
		return fmt.Errorf("summarization failed: %w", err)
	}

	// Update session context
	cm.updateSessionContext(session, result)

	// Update metrics
	cm.updateMetrics(true, result.OriginalSize, result.CompressedSize, time.Since(startTime))

	logger.InfoCF("compaction", "Compaction completed", map[string]any{
		"session_id":         session.SessionID,
		"messages_compacted": len(messagesToCompact),
		"original_size":      result.OriginalSize,
		"compressed_size":    result.CompressedSize,
		"compression_ratio":  float64(result.OriginalSize) / float64(result.CompressedSize),
		"duration":           time.Since(startTime).String(),
	})

	return nil
}

// extractMessagesToCompact extracts oldest messages to compact
func (cm *CompactionManager) extractMessagesToCompact(session *Session) []Message {
	session.mu.RLock()
	defer session.mu.RUnlock()

	totalMessages := len(session.Context)
	if totalMessages <= cm.config.KeepRecentCount {
		return nil
	}

	// Calculate how many to compact
	compactCount := min(cm.config.CompactBatchSize, totalMessages-cm.config.KeepRecentCount)

	// Extract oldest messages (make a copy)
	messages := make([]Message, compactCount)
	copy(messages, session.Context[:compactCount])
	return messages
}

// getExistingSummary gets existing summary if present
func (cm *CompactionManager) getExistingSummary(session *Session) string {
	if session.CompactedContext != nil {
		session.CompactedContext.mu.RLock()
		defer session.CompactedContext.mu.RUnlock()
		return session.CompactedContext.Summary
	}
	return ""
}

// updateSessionContext updates session with compacted context
func (cm *CompactionManager) updateSessionContext(session *Session, result *CompactionResult) {
	session.mu.Lock()
	defer session.mu.Unlock()

	// Calculate how many messages were compacted
	compactCount := result.MessagesCount

	// Create or update compacted context
	if session.CompactedContext == nil {
		session.CompactedContext = &CompactedContext{}
	}

	session.CompactedContext.mu.Lock()
	session.CompactedContext.Summary = result.Summary
	session.CompactedContext.SummaryVersion++
	session.CompactedContext.CompactedAt = time.Now()
	session.CompactedContext.CompactedCount += compactCount
	session.CompactedContext.OriginalSize = result.OriginalSize
	session.CompactedContext.CompressedSize = result.CompressedSize
	session.CompactedContext.mu.Unlock()

	// Remove compacted messages, keep recent ones
	session.Context = session.Context[compactCount:]
	session.LastCompaction = time.Now()
}

// updateMetrics updates compaction metrics
func (cm *CompactionManager) updateMetrics(success bool, originalSize, compressedSize int, duration time.Duration) {
	cm.metrics.mu.Lock()
	defer cm.metrics.mu.Unlock()

	cm.metrics.TotalCompactions++
	if success {
		cm.metrics.SuccessCount++
		cm.metrics.TotalBytesSaved += int64(originalSize - compressedSize)
		cm.metrics.TotalTimeSaved += duration

		if compressedSize > 0 {
			ratio := float64(originalSize) / float64(compressedSize)
			cm.metrics.AverageCompression = (cm.metrics.AverageCompression*float64(cm.metrics.SuccessCount-1) + ratio) / float64(cm.metrics.SuccessCount)
		}
	} else {
		cm.metrics.FailureCount++
	}
	cm.metrics.LastCompaction = time.Now()
}

// GetMetrics returns current metrics (copy without mutex)
func (cm *CompactionManager) GetMetrics() CompactionMetrics {
	cm.metrics.mu.RLock()
	defer cm.metrics.mu.RUnlock()

	// Return copy of fields without mutex
	return CompactionMetrics{
		TotalCompactions:   cm.metrics.TotalCompactions,
		SuccessCount:       cm.metrics.SuccessCount,
		FailureCount:       cm.metrics.FailureCount,
		TotalTimeSaved:     cm.metrics.TotalTimeSaved,
		TotalBytesSaved:    cm.metrics.TotalBytesSaved,
		AverageCompression: cm.metrics.AverageCompression,
		LastCompaction:     cm.metrics.LastCompaction,
		// Note: mu is NOT copied to avoid mutex copy warning
	}
}

// Stop gracefully stops the compaction manager
func (cm *CompactionManager) Stop() {
	close(cm.stopCh)
	cm.wg.Wait()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
