// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package providers

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// QueuedProvider wraps an LLMProvider with request queuing and worker pool.
type QueuedProvider struct {
	delegate   LLMProvider
	queue      RequestQueue
	workerPool *WorkerPool
	provider   string
}

// NewQueuedProvider creates a new provider with request queue and worker pool.
func NewQueuedProvider(delegate LLMProvider, providerName string, queueSize, workerCount, maxRetries int) *QueuedProvider {
	queue := NewInMemoryQueue(queueSize)
	backoffFunc := ExponentialBackoff(1*time.Second, 16*time.Second)
	workerPool := NewWorkerPool(queue, delegate, workerCount, maxRetries, backoffFunc)
	
	// Start workers
	workerPool.Start()
	
	return &QueuedProvider{
		delegate:   delegate,
		queue:      queue,
		workerPool: workerPool,
		provider:   providerName,
	}
}

// Chat executes the chat request through the queue system.
func (p *QueuedProvider) Chat(
	ctx context.Context,
	messages []Message,
	tools []ToolDefinition,
	model string,
	options map[string]any,
) (*LLMResponse, error) {
	// Create queued request
	req := &QueuedRequest{
		ID:        uuid.New().String(),
		Ctx:       ctx,
		Messages:  messages,
		Tools:     tools,
		Model:     model,
		Options:   options,
		ResultCh:  make(chan *QueuedResult, 1),
		EnqueueAt: time.Now(),
		Retries:   0,
	}
	
	// Enqueue request
	if err := p.queue.Enqueue(req); err != nil {
		return nil, fmt.Errorf("failed to enqueue request: %w", err)
	}
	
	p.workerPool.GetMetrics().RecordEnqueue()
	
	// Wait for result
	select {
	case result := <-req.ResultCh:
		if result.Error != nil {
			return nil, result.Error
		}
		return result.Response, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("request cancelled: %w", ctx.Err())
	}
}

// GetDefaultModel returns the delegate's default model.
func (p *QueuedProvider) GetDefaultModel() string {
	return p.delegate.GetDefaultModel()
}

// Close stops the worker pool and closes the queue.
func (p *QueuedProvider) Close() {
	p.workerPool.Stop()
	p.queue.Close()
	
	if sp, ok := p.delegate.(StatefulProvider); ok {
		sp.Close()
	}
}

// GetMetrics returns current queue metrics.
func (p *QueuedProvider) GetMetrics() *QueueMetrics {
	return p.workerPool.GetMetrics()
}

// GetQueueDepth returns current queue depth.
func (p *QueuedProvider) GetQueueDepth() int {
	return p.queue.Depth()
}
