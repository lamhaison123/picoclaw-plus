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

// MentionRequest represents a queued mention request
type MentionRequest struct {
	Role       string
	Prompt     string
	SessionID  string
	ChatID     int64
	TeamID     string
	Timestamp  time.Time
	RetryCount int
	Context    context.Context

	// Additional fields for execution
	Platform   Platform // Platform interface for sending messages
	Session    *Session // Session for context management
	TeamRoster string   // Team roster information
	Depth      int      // Current cascade depth

	// Callback for execution
	ExecuteFunc func(*MentionRequest) error
}

// MentionQueue manages queued mentions for a specific role
type MentionQueue struct {
	role         string
	queue        chan *MentionRequest
	rateLimit    time.Duration
	maxRetries   int
	retryBackoff time.Duration
	lastExecTime time.Time
	metrics      *QueueMetrics
	mu           sync.RWMutex
	stopCh       chan struct{}
	wg           sync.WaitGroup
}

// QueueMetrics tracks queue performance metrics
type QueueMetrics struct {
	QueueLength    int
	DroppedCount   int64
	ProcessedCount int64
	RetryCount     int64
	FailureCount   int64
	TotalWaitTime  time.Duration
	mu             sync.RWMutex
}

// NewMentionQueue creates a new mention queue for a role
func NewMentionQueue(role string, queueSize int, rateLimit time.Duration, maxRetries int, retryBackoff time.Duration) *MentionQueue {
	mq := &MentionQueue{
		role:         role,
		queue:        make(chan *MentionRequest, queueSize),
		rateLimit:    rateLimit,
		maxRetries:   maxRetries,
		retryBackoff: retryBackoff,
		metrics:      &QueueMetrics{},
		stopCh:       make(chan struct{}),
	}

	// Start worker
	mq.wg.Add(1)
	go mq.worker()

	return mq
}

// Enqueue adds a mention request to the queue
func (mq *MentionQueue) Enqueue(req *MentionRequest) error {
	select {
	case mq.queue <- req:
		mq.updateQueueLength(len(mq.queue))
		logger.InfoCF("collaborative", "Mention queued", map[string]any{
			"role":         mq.role,
			"session_id":   req.SessionID,
			"queue_length": len(mq.queue),
		})
		return nil
	default:
		// Queue is full
		mq.metrics.mu.Lock()
		mq.metrics.DroppedCount++
		mq.metrics.mu.Unlock()

		logger.WarnCF("collaborative", "Mention queue full, request dropped", map[string]any{
			"role":          mq.role,
			"session_id":    req.SessionID,
			"queue_size":    cap(mq.queue),
			"dropped_count": mq.metrics.DroppedCount,
		})

		return fmt.Errorf("mention queue for @%s is full (size: %d)", mq.role, cap(mq.queue))
	}
}

// worker processes queued mentions with rate limiting
func (mq *MentionQueue) worker() {
	defer mq.wg.Done()

	for {
		select {
		case <-mq.stopCh:
			logger.InfoCF("collaborative", "Mention queue worker stopping", map[string]any{
				"role": mq.role,
			})
			return

		case req := <-mq.queue:
			mq.updateQueueLength(len(mq.queue))
			mq.processRequest(req)
		}
	}
}

// processRequest processes a single mention request with rate limiting and retry
func (mq *MentionQueue) processRequest(req *MentionRequest) {
	// Calculate wait time
	waitStart := time.Now()

	// Rate limiting: ensure minimum time between executions
	mq.mu.Lock()
	timeSinceLastExec := time.Since(mq.lastExecTime)
	if timeSinceLastExec < mq.rateLimit {
		sleepDuration := mq.rateLimit - timeSinceLastExec
		mq.mu.Unlock()

		logger.InfoCF("collaborative", "Rate limiting mention", map[string]any{
			"role":           mq.role,
			"sleep_duration": sleepDuration.String(),
		})

		select {
		case <-req.Context.Done():
			return
		case <-time.After(sleepDuration):
		}
	} else {
		mq.mu.Unlock()
	}

	// Update last execution time
	mq.mu.Lock()
	mq.lastExecTime = time.Now()
	mq.mu.Unlock()

	// Track wait time
	waitTime := time.Since(waitStart)
	mq.metrics.mu.Lock()
	mq.metrics.TotalWaitTime += waitTime
	mq.metrics.mu.Unlock()

	// Execute with retry logic
	err := mq.executeWithRetry(req)

	if err != nil {
		mq.metrics.mu.Lock()
		mq.metrics.FailureCount++
		mq.metrics.mu.Unlock()

		logger.ErrorCF("collaborative", "Mention execution failed after retries", map[string]any{
			"role":        mq.role,
			"session_id":  req.SessionID,
			"retry_count": req.RetryCount,
			"error":       err.Error(),
		})
	} else {
		mq.metrics.mu.Lock()
		mq.metrics.ProcessedCount++
		mq.metrics.mu.Unlock()
	}
}

// executeWithRetry executes a mention request with exponential backoff retry
func (mq *MentionQueue) executeWithRetry(req *MentionRequest) error {
	var lastErr error

	for attempt := 0; attempt <= mq.maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff: 1s, 2s, 4s
			backoffDuration := mq.retryBackoff * time.Duration(1<<uint(attempt-1))

			logger.InfoCF("collaborative", "Retrying mention after backoff", map[string]any{
				"role":    mq.role,
				"attempt": attempt,
				"backoff": backoffDuration.String(),
			})

			select {
			case <-req.Context.Done():
				return fmt.Errorf("context cancelled during retry backoff: %w", req.Context.Err())
			case <-time.After(backoffDuration):
			}

			mq.metrics.mu.Lock()
			mq.metrics.RetryCount++
			mq.metrics.mu.Unlock()
		}

		// Execute the mention
		err := mq.execute(req)
		if err == nil {
			return nil
		}

		lastErr = err
		req.RetryCount = attempt + 1

		logger.WarnCF("collaborative", "Mention execution attempt failed", map[string]any{
			"role":    mq.role,
			"attempt": attempt + 1,
			"error":   err.Error(),
		})
	}

	return fmt.Errorf("mention failed after %d attempts: %w", mq.maxRetries+1, lastErr)
}

// execute performs the actual mention execution
func (mq *MentionQueue) execute(req *MentionRequest) error {
	if req.ExecuteFunc == nil {
		return fmt.Errorf("no execute function provided")
	}

	logger.InfoCF("collaborative", "Executing mention from queue", map[string]any{
		"role":       mq.role,
		"session_id": req.SessionID,
		"depth":      req.Depth,
		"retry":      req.RetryCount,
	})

	return req.ExecuteFunc(req)
}

// GetMetrics returns current queue metrics
func (mq *MentionQueue) GetMetrics() QueueMetrics {
	mq.metrics.mu.RLock()
	defer mq.metrics.mu.RUnlock()

	return QueueMetrics{
		QueueLength:    mq.metrics.QueueLength,
		DroppedCount:   mq.metrics.DroppedCount,
		ProcessedCount: mq.metrics.ProcessedCount,
		RetryCount:     mq.metrics.RetryCount,
		FailureCount:   mq.metrics.FailureCount,
		TotalWaitTime:  mq.metrics.TotalWaitTime,
	}
}

// GetAverageWaitTime returns average wait time per request
func (mq *MentionQueue) GetAverageWaitTime() time.Duration {
	mq.metrics.mu.RLock()
	defer mq.metrics.mu.RUnlock()

	if mq.metrics.ProcessedCount == 0 {
		return 0
	}

	return mq.metrics.TotalWaitTime / time.Duration(mq.metrics.ProcessedCount)
}

// updateQueueLength updates the queue length metric
func (mq *MentionQueue) updateQueueLength(length int) {
	mq.metrics.mu.Lock()
	mq.metrics.QueueLength = length
	mq.metrics.mu.Unlock()
}

// Stop gracefully stops the queue worker
func (mq *MentionQueue) Stop() {
	close(mq.stopCh)
	mq.wg.Wait()

	logger.InfoCF("collaborative", "Mention queue stopped", map[string]any{
		"role":            mq.role,
		"processed_count": mq.metrics.ProcessedCount,
		"dropped_count":   mq.metrics.DroppedCount,
	})
}
