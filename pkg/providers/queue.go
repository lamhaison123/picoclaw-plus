// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package providers

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	ErrQueueFull    = errors.New("queue is full")
	ErrQueueClosed  = errors.New("queue is closed")
	ErrQueueTimeout = errors.New("queue operation timeout")
)

// QueuedRequest represents a request waiting in the queue.
type QueuedRequest struct {
	ID        string
	Ctx       context.Context
	Messages  []Message
	Tools     []ToolDefinition
	Model     string
	Options   map[string]any
	ResultCh  chan *QueuedResult
	EnqueueAt time.Time
	Retries   int
}

// QueuedResult represents the result of a queued request.
type QueuedResult struct {
	Response *LLMResponse
	Error    error
}

// RequestQueue defines the interface for request queuing.
type RequestQueue interface {
	Enqueue(req *QueuedRequest) error
	Dequeue(ctx context.Context) (*QueuedRequest, error)
	Depth() int
	Close() error
}

// InMemoryQueue implements RequestQueue using Go channels.
type InMemoryQueue struct {
	queue   chan *QueuedRequest
	closed  bool
	closeMu sync.RWMutex
}

// NewInMemoryQueue creates a new in-memory queue with the specified capacity.
func NewInMemoryQueue(capacity int) *InMemoryQueue {
	return &InMemoryQueue{
		queue: make(chan *QueuedRequest, capacity),
	}
}

// Enqueue adds a request to the queue.
func (q *InMemoryQueue) Enqueue(req *QueuedRequest) error {
	q.closeMu.RLock()
	defer q.closeMu.RUnlock()

	if q.closed {
		return ErrQueueClosed
	}

	select {
	case q.queue <- req:
		return nil
	default:
		return ErrQueueFull
	}
}

// Dequeue removes and returns a request from the queue.
func (q *InMemoryQueue) Dequeue(ctx context.Context) (*QueuedRequest, error) {
	q.closeMu.RLock()
	defer q.closeMu.RUnlock()

	if q.closed {
		return nil, ErrQueueClosed
	}

	select {
	case req := <-q.queue:
		return req, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// Depth returns the current number of items in the queue.
func (q *InMemoryQueue) Depth() int {
	return len(q.queue)
}

// Close closes the queue and prevents new enqueues.
func (q *InMemoryQueue) Close() error {
	q.closeMu.Lock()
	defer q.closeMu.Unlock()

	if q.closed {
		return nil
	}

	q.closed = true
	close(q.queue)
	return nil
}
