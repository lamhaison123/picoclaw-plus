// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package collaborative

import (
	"sync"
	"time"

	"github.com/sipeed/picoclaw/pkg/logger"
)

// QueueManager manages mention queues for all roles
type QueueManager struct {
	queues       map[string]*MentionQueue
	queueSize    int
	rateLimit    time.Duration
	maxRetries   int
	retryBackoff time.Duration
	mu           sync.RWMutex
}

// NewQueueManager creates a new queue manager
func NewQueueManager(queueSize int, rateLimit time.Duration, maxRetries int, retryBackoff time.Duration) *QueueManager {
	return &QueueManager{
		queues:       make(map[string]*MentionQueue),
		queueSize:    queueSize,
		rateLimit:    rateLimit,
		maxRetries:   maxRetries,
		retryBackoff: retryBackoff,
	}
}

// GetOrCreateQueue gets or creates a queue for a role
func (qm *QueueManager) GetOrCreateQueue(role string) *MentionQueue {
	qm.mu.RLock()
	queue, exists := qm.queues[role]
	qm.mu.RUnlock()

	if exists {
		return queue
	}

	qm.mu.Lock()
	defer qm.mu.Unlock()

	// Double-check after acquiring write lock
	if queue, exists := qm.queues[role]; exists {
		return queue
	}

	// Create new queue
	queue = NewMentionQueue(role, qm.queueSize, qm.rateLimit, qm.maxRetries, qm.retryBackoff)
	qm.queues[role] = queue

	logger.InfoCF("collaborative", "Created mention queue", map[string]any{
		"role":       role,
		"queue_size": qm.queueSize,
		"rate_limit": qm.rateLimit.String(),
	})

	return queue
}

// Enqueue adds a mention request to the appropriate role queue
func (qm *QueueManager) Enqueue(req *MentionRequest) error {
	queue := qm.GetOrCreateQueue(req.Role)
	return queue.Enqueue(req)
}

// GetMetrics returns metrics for a specific role
func (qm *QueueManager) GetMetrics(role string) *QueueMetrics {
	qm.mu.RLock()
	queue, exists := qm.queues[role]
	qm.mu.RUnlock()

	if !exists {
		return nil
	}

	metrics := queue.GetMetrics()
	return &metrics
}

// GetAllMetrics returns metrics for all roles
func (qm *QueueManager) GetAllMetrics() map[string]QueueMetrics {
	qm.mu.RLock()
	defer qm.mu.RUnlock()

	result := make(map[string]QueueMetrics)
	for role, queue := range qm.queues {
		result[role] = queue.GetMetrics()
	}

	return result
}

// Stop gracefully stops all queues
func (qm *QueueManager) Stop() {
	qm.mu.Lock()
	defer qm.mu.Unlock()

	for role, queue := range qm.queues {
		logger.InfoCF("collaborative", "Stopping queue", map[string]any{
			"role": role,
		})
		queue.Stop()
	}

	qm.queues = make(map[string]*MentionQueue)
}
