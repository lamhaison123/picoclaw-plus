// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package collaborative

import (
	"context"
	"crypto/sha256"
	"fmt"
	"sync"
	"time"
)

// DispatchEntry tracks when a message was dispatched
type DispatchEntry struct {
	Timestamp time.Time
	MessageID string
}

// DispatchTracker tracks dispatched mentions to prevent duplicates with TTL
type DispatchTracker struct {
	dispatched map[string]*DispatchEntry
	mu         sync.RWMutex
	maxSize    int
	ttl        time.Duration
	metrics    *EnhancedMetrics
}

// NewDispatchTracker creates a new dispatch tracker with default settings
func NewDispatchTracker() *DispatchTracker {
	return NewDispatchTrackerWithConfig(1000, 1*time.Hour)
}

// NewDispatchTrackerWithConfig creates a new dispatch tracker with custom settings
func NewDispatchTrackerWithConfig(maxSize int, ttl time.Duration) *DispatchTracker {
	return &DispatchTracker{
		dispatched: make(map[string]*DispatchEntry),
		maxSize:    maxSize,
		ttl:        ttl,
		metrics:    NewEnhancedMetrics(),
	}
}

// GenerateMessageID generates a unique ID for a mention dispatch
// Format: chatID:sessionID:role:contentHash
func GenerateMessageID(chatID int64, sessionID, role, content string) string {
	// Use first 100 chars of content to avoid huge keys
	contentSnippet := content
	if len(content) > 100 {
		contentSnippet = content[:100]
	}

	// Create hash of content for uniqueness
	hash := sha256.Sum256([]byte(contentSnippet))
	hashStr := fmt.Sprintf("%x", hash[:8]) // Use first 8 bytes

	return fmt.Sprintf("%d:%s:%s:%s", chatID, sessionID, role, hashStr)
}

// IsDispatched checks if a message has been dispatched and is still valid (within TTL)
func (dt *DispatchTracker) IsDispatched(messageID string) bool {
	dt.mu.RLock()
	defer dt.mu.RUnlock()

	entry, exists := dt.dispatched[messageID]
	if !exists {
		if dt.metrics != nil {
			dt.metrics.RecordIdempotencyMiss()
		}
		return false
	}

	// Check if entry has expired
	if time.Since(entry.Timestamp) > dt.ttl {
		if dt.metrics != nil {
			dt.metrics.RecordIdempotencyKeyExpired()
		}
		return false // Expired, treat as not dispatched
	}

	// Hit - duplicate detected
	if dt.metrics != nil {
		dt.metrics.RecordIdempotencyHit()
	}
	return true
}

// TryMarkDispatched atomically checks if a message is dispatched and marks it if not.
// Returns true if successfully marked (not previously dispatched), false if already dispatched.
// This prevents race conditions in check-then-mark patterns.
func (dt *DispatchTracker) TryMarkDispatched(messageID string) bool {
	dt.mu.Lock()
	defer dt.mu.Unlock()

	// Check if already exists and is still valid
	if entry, exists := dt.dispatched[messageID]; exists {
		if time.Since(entry.Timestamp) <= dt.ttl {
			if dt.metrics != nil {
				dt.metrics.RecordIdempotencyHit()
			}
			return false // Already dispatched and still valid
		}
		// Entry expired, will be replaced
		if dt.metrics != nil {
			dt.metrics.RecordIdempotencyKeyExpired()
		}
	} else {
		if dt.metrics != nil {
			dt.metrics.RecordIdempotencyMiss()
		}
	}

	// Mark as dispatched atomically
	dt.dispatched[messageID] = &DispatchEntry{
		Timestamp: time.Now(),
		MessageID: messageID,
	}

	// Trigger cleanup if size exceeds threshold
	if len(dt.dispatched) > dt.maxSize {
		dt.cleanupOldEntriesLocked()
	}

	return true // Successfully marked
}

// MarkDispatched marks a message as dispatched with current timestamp
func (dt *DispatchTracker) MarkDispatched(messageID string) {
	dt.mu.Lock()
	defer dt.mu.Unlock()

	dt.dispatched[messageID] = &DispatchEntry{
		Timestamp: time.Now(),
		MessageID: messageID,
	}

	// Trigger cleanup if size exceeds threshold
	if len(dt.dispatched) > dt.maxSize {
		dt.cleanupOldEntriesLocked()
	}
}

// GetMetrics returns the enhanced metrics collector
func (dt *DispatchTracker) GetMetrics() *EnhancedMetrics {
	return dt.metrics
}

// Clear removes all entries (use for explicit cleanup)
func (dt *DispatchTracker) Clear() {
	dt.mu.Lock()
	defer dt.mu.Unlock()

	// Create a new map to allow GC to reclaim memory
	dt.dispatched = make(map[string]*DispatchEntry)
}

// Size returns the number of tracked dispatches
func (dt *DispatchTracker) Size() int {
	dt.mu.RLock()
	defer dt.mu.RUnlock()
	return len(dt.dispatched)
}

// cleanupOldEntriesLocked removes expired entries (must be called with lock held)
func (dt *DispatchTracker) cleanupOldEntriesLocked() {
	now := time.Now()
	removed := 0

	for id, entry := range dt.dispatched {
		if now.Sub(entry.Timestamp) > dt.ttl {
			delete(dt.dispatched, id)
			removed++
		}
	}

	// If still over max size after cleanup, remove oldest entries
	if len(dt.dispatched) > dt.maxSize {
		// Convert to slice for sorting by timestamp
		entries := make([]*DispatchEntry, 0, len(dt.dispatched))
		for _, entry := range dt.dispatched {
			entries = append(entries, entry)
		}

		// Sort by timestamp (oldest first)
		for i := 0; i < len(entries)-1; i++ {
			for j := i + 1; j < len(entries); j++ {
				if entries[i].Timestamp.After(entries[j].Timestamp) {
					entries[i], entries[j] = entries[j], entries[i]
				}
			}
		}

		// Remove oldest entries until we're under max size
		toRemove := len(dt.dispatched) - dt.maxSize
		for i := 0; i < toRemove && i < len(entries); i++ {
			delete(dt.dispatched, entries[i].MessageID)
			removed++
		}
	}
}

// CleanupOldEntries removes expired entries (public method with locking)
func (dt *DispatchTracker) CleanupOldEntries() int {
	dt.mu.Lock()
	defer dt.mu.Unlock()

	oldSize := len(dt.dispatched)
	dt.cleanupOldEntriesLocked()
	return oldSize - len(dt.dispatched)
}

// StartPeriodicCleanup starts a goroutine that periodically cleans up old entries
func (dt *DispatchTracker) StartPeriodicCleanup(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			removed := dt.CleanupOldEntries()
			if removed > 0 {
				// Log cleanup (would need logger import)
				// logger.DebugCF("collaborative", "Cleaned up dispatch tracker", ...)
			}
		case <-ctx.Done():
			return
		}
	}
}
