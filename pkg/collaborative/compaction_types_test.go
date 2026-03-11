// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package collaborative

import (
	"testing"
	"time"
)

func TestCompactedContext_Creation(t *testing.T) {
	ctx := &CompactedContext{
		Summary:        "Test summary",
		SummaryVersion: 1,
		CompactedAt:    time.Now(),
		CompactedCount: 10,
		OriginalSize:   1000,
		CompressedSize: 100,
	}

	if ctx.CompactedAt.IsZero() {
		t.Error("Expected CompactedAt to be set")
	}
	if ctx.OriginalSize != 1000 {
		t.Errorf("Expected OriginalSize 1000, got %d", ctx.OriginalSize)
	}
	if ctx.CompressedSize != 100 {
		t.Errorf("Expected CompressedSize 100, got %d", ctx.CompressedSize)
	}

	if ctx.Summary != "Test summary" {
		t.Errorf("Expected summary 'Test summary', got '%s'", ctx.Summary)
	}
	if ctx.SummaryVersion != 1 {
		t.Errorf("Expected version 1, got %d", ctx.SummaryVersion)
	}
	if ctx.CompactedCount != 10 {
		t.Errorf("Expected count 10, got %d", ctx.CompactedCount)
	}
}

func TestCompactionConfig_Defaults(t *testing.T) {
	config := CompactionConfig{
		Enabled:          true,
		TriggerThreshold: 40,
		KeepRecentCount:  15,
		CompactBatchSize: 25,
		MinInterval:      5 * time.Minute,
		SummaryMaxLength: 2000,
		LLMProvider:      "openai",
		LLMModel:         "gpt-4o-mini",
		LLMMaxRetries:    3,
	}

	if config.LLMTimeout != 30*time.Second {
		t.Errorf("Expected LLMTimeout 30s, got %v", config.LLMTimeout)
	}
	if config.LLMMaxRetries != 3 {
		t.Errorf("Expected LLMMaxRetries 3, got %d", config.LLMMaxRetries)
	}

	if !config.Enabled {
		t.Error("Expected Enabled to be true")
	}
	if config.TriggerThreshold != 40 {
		t.Errorf("Expected TriggerThreshold 40, got %d", config.TriggerThreshold)
	}
	if config.KeepRecentCount != 15 {
		t.Errorf("Expected KeepRecentCount 15, got %d", config.KeepRecentCount)
	}
	if config.CompactBatchSize != 25 {
		t.Errorf("Expected CompactBatchSize 25, got %d", config.CompactBatchSize)
	}
	if config.MinInterval != 5*time.Minute {
		t.Errorf("Expected MinInterval 5m, got %v", config.MinInterval)
	}
	if config.SummaryMaxLength != 2000 {
		t.Errorf("Expected SummaryMaxLength 2000, got %d", config.SummaryMaxLength)
	}
	if config.LLMProvider != "openai" {
		t.Errorf("Expected LLMProvider 'openai', got '%s'", config.LLMProvider)
	}
	if config.LLMModel != "gpt-4o-mini" {
		t.Errorf("Expected LLMModel 'gpt-4o-mini', got '%s'", config.LLMModel)
	}
}

func TestCompactionMetrics_Initialization(t *testing.T) {
	metrics := &CompactionMetrics{}

	if metrics.TotalCompactions != 0 {
		t.Errorf("Expected TotalCompactions 0, got %d", metrics.TotalCompactions)
	}
	if metrics.SuccessCount != 0 {
		t.Errorf("Expected SuccessCount 0, got %d", metrics.SuccessCount)
	}
	if metrics.FailureCount != 0 {
		t.Errorf("Expected FailureCount 0, got %d", metrics.FailureCount)
	}
	if metrics.AverageCompression != 0 {
		t.Errorf("Expected AverageCompression 0, got %f", metrics.AverageCompression)
	}
}

func TestCompactionRequest_Creation(t *testing.T) {
	messages := []Message{
		{Role: "user", Content: "Hello", Timestamp: time.Now()},
		{Role: "architect", Content: "Hi", Timestamp: time.Now()},
	}

	config := CompactionConfig{
		Enabled:          true,
		TriggerThreshold: 40,
	}

	req := &CompactionRequest{
		SessionID:       "test123",
		Messages:        messages,
		ExistingSummary: "Previous summary",
		Config:          config,
		Timestamp:       time.Now(),
	}

	if req.Config.Enabled != true {
		t.Error("Expected Config.Enabled to be true")
	}
	if req.Timestamp.IsZero() {
		t.Error("Expected Timestamp to be set")
	}

	if req.SessionID != "test123" {
		t.Errorf("Expected SessionID 'test123', got '%s'", req.SessionID)
	}
	if len(req.Messages) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(req.Messages))
	}
	if req.ExistingSummary != "Previous summary" {
		t.Errorf("Expected ExistingSummary 'Previous summary', got '%s'", req.ExistingSummary)
	}
}

func TestCompactionResult_Success(t *testing.T) {
	result := &CompactionResult{
		MessagesCount: 25,
	}

	if result.Summary != "Compacted summary" {
		t.Errorf("Expected Summary 'Compacted summary', got '%s'", result.Summary)
	}
	if result.Duration != 2*time.Second {
		t.Errorf("Expected Duration 2s, got %v", result.Duration)
	}

	if !result.Success {
		t.Error("Expected Success to be true")
	}
	if result.Error != nil {
		t.Errorf("Expected no error, got %v", result.Error)
	}
	if result.OriginalSize != 1000 {
		t.Errorf("Expected OriginalSize 1000, got %d", result.OriginalSize)
	}
	if result.CompressedSize != 100 {
		t.Errorf("Expected CompressedSize 100, got %d", result.CompressedSize)
	}
	if result.MessagesCount != 25 {
		t.Errorf("Expected MessagesCount 25, got %d", result.MessagesCount)
	}

	// Check compression ratio
	ratio := float64(result.OriginalSize) / float64(result.CompressedSize)
	if ratio != 10.0 {
		t.Errorf("Expected compression ratio 10.0, got %f", ratio)
	}
}

func TestCompactedContext_ThreadSafety(t *testing.T) {
	ctx := &CompactedContext{
		Summary:        "Initial",
		SummaryVersion: 1,
	}

	// Test concurrent reads
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			ctx.mu.RLock()
			_ = ctx.Summary
			ctx.mu.RUnlock()
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Test concurrent writes
	for i := 0; i < 10; i++ {
		go func(version int) {
			ctx.mu.Lock()
			ctx.SummaryVersion = version
			ctx.mu.Unlock()
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify final state
	if ctx.SummaryVersion < 0 || ctx.SummaryVersion >= 10 {
		t.Errorf("Expected SummaryVersion between 0-9, got %d", ctx.SummaryVersion)
	}
}

func TestCompactionMetrics_ThreadSafety(t *testing.T) {
	metrics := &CompactionMetrics{}

	// Test concurrent updates
	done := make(chan bool)
	for i := 0; i < 100; i++ {
		go func() {
			metrics.mu.Lock()
			metrics.TotalCompactions++
			metrics.SuccessCount++
			metrics.mu.Unlock()
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 100; i++ {
		<-done
	}

	// Verify final counts
	if metrics.TotalCompactions != 100 {
		t.Errorf("Expected TotalCompactions 100, got %d", metrics.TotalCompactions)
	}
	if metrics.SuccessCount != 100 {
		t.Errorf("Expected SuccessCount 100, got %d", metrics.SuccessCount)
	}
}
