// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package providers

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// WorkerPool manages a pool of workers that process queued requests.
type WorkerPool struct {
	queue       RequestQueue
	provider    LLMProvider
	workerCount int
	maxRetries  int
	backoffFunc BackoffFunc
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
	metrics     *QueueMetrics
}

// BackoffFunc calculates the backoff duration for a given retry attempt.
type BackoffFunc func(attempt int) time.Duration

// QueueMetrics tracks queue performance metrics.
type QueueMetrics struct {
	mu              sync.RWMutex
	Enqueued        int64
	Processed       int64
	Failed          int64
	Retried         int64
	CurrentDepth    int
	AvgProcessTime  time.Duration
	processTimes    []time.Duration
	maxProcessTimes int
}

// NewQueueMetrics creates a new metrics tracker.
func NewQueueMetrics() *QueueMetrics {
	return &QueueMetrics{
		maxProcessTimes: 100, // Keep last 100 samples
		processTimes:    make([]time.Duration, 0, 100),
	}
}

// RecordEnqueue increments enqueued counter.
func (m *QueueMetrics) RecordEnqueue() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Enqueued++
}

// RecordProcessed increments processed counter and updates avg time.
func (m *QueueMetrics) RecordProcessed(duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Processed++
	
	m.processTimes = append(m.processTimes, duration)
	if len(m.processTimes) > m.maxProcessTimes {
		m.processTimes = m.processTimes[1:]
	}
	
	var total time.Duration
	for _, d := range m.processTimes {
		total += d
	}
	m.AvgProcessTime = total / time.Duration(len(m.processTimes))
}

// RecordFailed increments failed counter.
func (m *QueueMetrics) RecordFailed() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Failed++
}

// RecordRetry increments retry counter.
func (m *QueueMetrics) RecordRetry() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Retried++
}

// UpdateDepth updates current queue depth.
func (m *QueueMetrics) UpdateDepth(depth int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.CurrentDepth = depth
}

// GetStats returns current metrics snapshot.
func (m *QueueMetrics) GetStats() (enqueued, processed, failed, retried int64, depth int, avgTime time.Duration) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.Enqueued, m.Processed, m.Failed, m.Retried, m.CurrentDepth, m.AvgProcessTime
}

// ExponentialBackoff returns a backoff function with exponential delays.
func ExponentialBackoff(initial, max time.Duration) BackoffFunc {
	return func(attempt int) time.Duration {
		backoff := initial * time.Duration(1<<uint(attempt))
		if backoff > max {
			backoff = max
		}
		return backoff
	}
}

// NewWorkerPool creates a new worker pool.
func NewWorkerPool(queue RequestQueue, provider LLMProvider, workerCount, maxRetries int, backoffFunc BackoffFunc) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		queue:       queue,
		provider:    provider,
		workerCount: workerCount,
		maxRetries:  maxRetries,
		backoffFunc: backoffFunc,
		ctx:         ctx,
		cancel:      cancel,
		metrics:     NewQueueMetrics(),
	}
}

// Start starts all workers in the pool.
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workerCount; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

// Stop stops all workers and waits for them to finish.
func (wp *WorkerPool) Stop() {
	wp.cancel()
	wp.wg.Wait()
}

// worker processes requests from the queue.
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()
	
	for {
		select {
		case <-wp.ctx.Done():
			return
		default:
		}
		
		// Dequeue with timeout
		dequeueCtx, cancel := context.WithTimeout(wp.ctx, 5*time.Second)
		req, err := wp.queue.Dequeue(dequeueCtx)
		cancel()
		
		if err != nil {
			if err == context.Canceled || err == context.DeadlineExceeded {
				continue
			}
			if err == ErrQueueClosed {
				return
			}
			log.Printf("[Worker %d] Dequeue error: %v", id, err)
			continue
		}
		
		if req == nil {
			continue
		}
		
		// Update queue depth metric
		wp.metrics.UpdateDepth(wp.queue.Depth())
		
		// Process request with retries
		wp.processRequest(id, req)
	}
}

// processRequest processes a single request with retry logic.
func (wp *WorkerPool) processRequest(workerID int, req *QueuedRequest) {
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime)
		wp.metrics.RecordProcessed(duration)
	}()
	
	var lastErr error
	
	for attempt := 0; attempt <= wp.maxRetries; attempt++ {
		// Check if request context is still valid
		select {
		case <-req.Ctx.Done():
			result := &QueuedResult{
				Error: fmt.Errorf("request context cancelled: %w", req.Ctx.Err()),
			}
			select {
			case req.ResultCh <- result:
			default:
			}
			wp.metrics.RecordFailed()
			return
		default:
		}
		
		// Apply backoff for retries
		if attempt > 0 {
			backoff := wp.backoffFunc(attempt - 1)
			log.Printf("[Worker %d] Retry %d/%d for request %s after %v", 
				workerID, attempt, wp.maxRetries, req.ID, backoff)
			wp.metrics.RecordRetry()
			
			select {
			case <-time.After(backoff):
			case <-wp.ctx.Done():
				return
			case <-req.Ctx.Done():
				result := &QueuedResult{
					Error: fmt.Errorf("request context cancelled during backoff: %w", req.Ctx.Err()),
				}
				select {
				case req.ResultCh <- result:
				default:
				}
				wp.metrics.RecordFailed()
				return
			}
		}
		
		// Execute LLM call with progressive timeout (increases on retry)
		// Base: 120s, increases by 15s per retry, cap at 150s
		// Attempt 0: 120s, Attempt 1: 135s, Attempt 2: 150s
		baseTimeout := 120 * time.Second
		progressiveIncrease := time.Duration(attempt) * 15 * time.Second
		callTimeout := baseTimeout + progressiveIncrease
		if callTimeout > 150*time.Second {
			callTimeout = 150 * time.Second // Cap at 2.5 minutes per architecture review
		}
		
		callCtx, cancel := context.WithTimeout(req.Ctx, callTimeout)
		resp, err := wp.provider.Chat(callCtx, req.Messages, req.Tools, req.Model, req.Options)
		cancel()
		
		if err == nil {
			// Success - send result
			result := &QueuedResult{
				Response: resp,
			}
			select {
			case req.ResultCh <- result:
			default:
				log.Printf("[Worker %d] Failed to send result for request %s (channel closed)", workerID, req.ID)
			}
			return
		}
		
		lastErr = err
		
		// Check if error is retriable
		if failErr, ok := err.(*FailoverError); ok {
			if !failErr.IsRetriable() {
				// Non-retriable error (e.g., format error)
				result := &QueuedResult{
					Error: err,
				}
				select {
				case req.ResultCh <- result:
				default:
				}
				wp.metrics.RecordFailed()
				return
			}
		}
		
		log.Printf("[Worker %d] Attempt %d/%d failed for request %s: %v", 
			workerID, attempt+1, wp.maxRetries+1, req.ID, err)
	}
	
	// All retries exhausted
	result := &QueuedResult{
		Error: fmt.Errorf("max retries (%d) exhausted: %w", wp.maxRetries, lastErr),
	}
	select {
	case req.ResultCh <- result:
	default:
		log.Printf("[Worker %d] Failed to send error result for request %s (channel closed)", workerID, req.ID)
	}
	wp.metrics.RecordFailed()
}

// GetMetrics returns the current metrics.
func (wp *WorkerPool) GetMetrics() *QueueMetrics {
	return wp.metrics
}
