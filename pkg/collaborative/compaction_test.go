// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package collaborative

import (
	"context"
	"errors"
	"testing"
	"time"
)

// MockSummarizer for testing
type MockSummarizer struct {
	shouldFail bool
	delay      time.Duration
}

func (m *MockSummarizer) Summarize(ctx context.Context, req *CompactionRequest) (*CompactionResult, error) {
	if m.delay > 0 {
		time.Sleep(m.delay)
	}

	if m.shouldFail {
		return nil, errors.New("mock summarization failed")
	}

	// Calculate sizes
	originalSize := 0
	for _, msg := range req.Messages {
		originalSize += len(msg.Content) + len(msg.Role) + 50
	}

	summary := "Mock summary of " + req.SessionID
	if req.ExistingSummary != "" {
		summary = req.ExistingSummary + " + " + summary
	}

	return &CompactionResult{
		Summary:        summary,
		Success:        true,
		Error:          nil,
		OriginalSize:   originalSize,
		CompressedSize: len(summary),
		Duration:       m.delay,
		MessagesCount:  len(req.Messages),
	}, nil
}

func TestNewCompactionManager(t *testing.T) {
	config := CompactionConfig{
		Enabled:          true,
		TriggerThreshold: 40,
	}
	summarizer := &MockSummarizer{}

	cm := NewCompactionManager(config, summarizer)

	if cm == nil {
		t.Fatal("Expected CompactionManager to be created")
	}
	if !cm.config.Enabled {
		t.Error("Expected config.Enabled to be true")
	}
	if cm.metrics == nil {
		t.Error("Expected metrics to be initialized")
	}
}

func TestShouldCompact_Disabled(t *testing.T) {
	config := CompactionConfig{
		Enabled: false,
	}
	cm := NewCompactionManager(config, &MockSummarizer{})

	session := NewSession(12345, "team1", 50)
	// Add many messages
	for i := 0; i < 50; i++ {
		session.AddMessage("user", "test message", nil)
	}

	if cm.ShouldCompact(session) {
		t.Error("Expected ShouldCompact to return false when disabled")
	}
}

func TestShouldCompact_BelowThreshold(t *testing.T) {
	config := CompactionConfig{
		Enabled:          true,
		TriggerThreshold: 40,
	}
	cm := NewCompactionManager(config, &MockSummarizer{})

	session := NewSession(12345, "team1", 50)
	// Add fewer messages than threshold
	for i := 0; i < 30; i++ {
		session.AddMessage("user", "test message", nil)
	}

	if cm.ShouldCompact(session) {
		t.Error("Expected ShouldCompact to return false when below threshold")
	}
}

func TestShouldCompact_AboveThreshold(t *testing.T) {
	config := CompactionConfig{
		Enabled:          true,
		TriggerThreshold: 40,
		MinInterval:      1 * time.Millisecond,
	}
	cm := NewCompactionManager(config, &MockSummarizer{})

	session := NewSession(12345, "team1", 50)
	// Add more messages than threshold
	for i := 0; i < 45; i++ {
		session.AddMessage("user", "test message", nil)
	}

	if !cm.ShouldCompact(session) {
		t.Error("Expected ShouldCompact to return true when above threshold")
	}
}

func TestShouldCompact_MinInterval(t *testing.T) {
	config := CompactionConfig{
		Enabled:          true,
		TriggerThreshold: 40,
		MinInterval:      1 * time.Hour, // Long interval
	}
	cm := NewCompactionManager(config, &MockSummarizer{})

	session := NewSession(12345, "team1", 50)
	// Add messages
	for i := 0; i < 45; i++ {
		session.AddMessage("user", "test message", nil)
	}

	// Set recent compaction
	session.LastCompaction = time.Now()

	if cm.ShouldCompact(session) {
		t.Error("Expected ShouldCompact to return false when within MinInterval")
	}
}

func TestExtractMessagesToCompact(t *testing.T) {
	config := CompactionConfig{
		KeepRecentCount:  15,
		CompactBatchSize: 25,
	}
	cm := NewCompactionManager(config, &MockSummarizer{})

	session := NewSession(12345, "team1", 50)
	// Add 45 messages
	for i := 0; i < 45; i++ {
		session.AddMessage("user", "message "+string(rune(i)), nil)
	}

	messages := cm.extractMessagesToCompact(session)

	// Should extract 25 messages (CompactBatchSize)
	if len(messages) != 25 {
		t.Errorf("Expected 25 messages, got %d", len(messages))
	}

	// Should be oldest messages
	if messages[0].Content != "message \x00" {
		t.Errorf("Expected first message to be oldest, got %s", messages[0].Content)
	}
}

func TestExtractMessagesToCompact_TooFewMessages(t *testing.T) {
	config := CompactionConfig{
		KeepRecentCount:  15,
		CompactBatchSize: 25,
	}
	cm := NewCompactionManager(config, &MockSummarizer{})

	session := NewSession(12345, "team1", 50)
	// Add only 10 messages (less than KeepRecentCount)
	for i := 0; i < 10; i++ {
		session.AddMessage("user", "test message", nil)
	}

	messages := cm.extractMessagesToCompact(session)

	if messages != nil {
		t.Errorf("Expected nil messages, got %d", len(messages))
	}
}

func TestGetExistingSummary(t *testing.T) {
	cm := NewCompactionManager(CompactionConfig{}, &MockSummarizer{})

	session := NewSession(12345, "team1", 50)

	// No existing summary
	summary := cm.getExistingSummary(session)
	if summary != "" {
		t.Errorf("Expected empty summary, got '%s'", summary)
	}

	// With existing summary
	session.CompactedContext = &CompactedContext{
		Summary: "Existing summary",
	}

	summary = cm.getExistingSummary(session)
	if summary != "Existing summary" {
		t.Errorf("Expected 'Existing summary', got '%s'", summary)
	}
}

func TestCompact_Success(t *testing.T) {
	config := CompactionConfig{
		Enabled:          true,
		TriggerThreshold: 40,
		KeepRecentCount:  15,
		CompactBatchSize: 25,
		MinInterval:      1 * time.Millisecond,
	}
	cm := NewCompactionManager(config, &MockSummarizer{})

	session := NewSession(12345, "team1", 50)
	// Add 45 messages
	for i := 0; i < 45; i++ {
		session.AddMessage("user", "test message", nil)
	}

	ctx := context.Background()
	err := cm.Compact(ctx, session)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Check session was updated
	if session.CompactedContext == nil {
		t.Error("Expected CompactedContext to be set")
	}

	// Check messages were removed
	if len(session.Context) != 20 { // 45 - 25 = 20
		t.Errorf("Expected 20 messages remaining, got %d", len(session.Context))
	}

	// Check metrics
	metrics := cm.GetMetrics()
	if metrics.TotalCompactions != 1 {
		t.Errorf("Expected 1 compaction, got %d", metrics.TotalCompactions)
	}
	if metrics.SuccessCount != 1 {
		t.Errorf("Expected 1 success, got %d", metrics.SuccessCount)
	}
}

func TestCompact_Failure(t *testing.T) {
	config := CompactionConfig{
		Enabled:          true,
		TriggerThreshold: 40,
		KeepRecentCount:  15,
		CompactBatchSize: 25,
	}
	cm := NewCompactionManager(config, &MockSummarizer{shouldFail: true})

	session := NewSession(12345, "team1", 50)
	// Add 45 messages
	for i := 0; i < 45; i++ {
		session.AddMessage("user", "test message", nil)
	}

	ctx := context.Background()
	err := cm.Compact(ctx, session)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	// Check metrics
	metrics := cm.GetMetrics()
	if metrics.TotalCompactions != 1 {
		t.Errorf("Expected 1 compaction attempt, got %d", metrics.TotalCompactions)
	}
	if metrics.FailureCount != 1 {
		t.Errorf("Expected 1 failure, got %d", metrics.FailureCount)
	}
}

func TestCompactAsync(t *testing.T) {
	config := CompactionConfig{
		Enabled:          true,
		TriggerThreshold: 40,
		KeepRecentCount:  15,
		CompactBatchSize: 25,
		MinInterval:      1 * time.Millisecond,
	}
	cm := NewCompactionManager(config, &MockSummarizer{delay: 50 * time.Millisecond})

	session := NewSession(12345, "team1", 50)
	// Add 45 messages
	for i := 0; i < 45; i++ {
		session.AddMessage("user", "test message", nil)
	}

	ctx := context.Background()
	cm.CompactAsync(ctx, session)

	// Should return immediately (async)
	// Wait for completion
	cm.Stop()

	// Check compaction happened
	if session.CompactedContext == nil {
		t.Error("Expected CompactedContext to be set after async compaction")
	}
}

func TestUpdateSessionContext(t *testing.T) {
	cm := NewCompactionManager(CompactionConfig{}, &MockSummarizer{})

	session := NewSession(12345, "team1", 50)
	// Add 30 messages
	for i := 0; i < 30; i++ {
		session.AddMessage("user", "test message", nil)
	}

	result := &CompactionResult{
		Summary:        "Test summary",
		Success:        true,
		OriginalSize:   1000,
		CompressedSize: 100,
		MessagesCount:  10,
	}

	cm.updateSessionContext(session, result)

	// Check CompactedContext was created
	if session.CompactedContext == nil {
		t.Fatal("Expected CompactedContext to be created")
	}

	// Check summary
	if session.CompactedContext.Summary != "Test summary" {
		t.Errorf("Expected summary 'Test summary', got '%s'", session.CompactedContext.Summary)
	}

	// Check version
	if session.CompactedContext.SummaryVersion != 1 {
		t.Errorf("Expected version 1, got %d", session.CompactedContext.SummaryVersion)
	}

	// Check messages were removed
	if len(session.Context) != 20 { // 30 - 10 = 20
		t.Errorf("Expected 20 messages, got %d", len(session.Context))
	}

	// Check LastCompaction was updated
	if session.LastCompaction.IsZero() {
		t.Error("Expected LastCompaction to be set")
	}
}

func TestUpdateMetrics_Success(t *testing.T) {
	cm := NewCompactionManager(CompactionConfig{}, &MockSummarizer{})

	cm.updateMetrics(true, 1000, 100, 2*time.Second)

	metrics := cm.GetMetrics()

	if metrics.TotalCompactions != 1 {
		t.Errorf("Expected 1 compaction, got %d", metrics.TotalCompactions)
	}
	if metrics.SuccessCount != 1 {
		t.Errorf("Expected 1 success, got %d", metrics.SuccessCount)
	}
	if metrics.TotalBytesSaved != 900 {
		t.Errorf("Expected 900 bytes saved, got %d", metrics.TotalBytesSaved)
	}
	if metrics.AverageCompression != 10.0 {
		t.Errorf("Expected compression ratio 10.0, got %f", metrics.AverageCompression)
	}
}

func TestUpdateMetrics_Failure(t *testing.T) {
	cm := NewCompactionManager(CompactionConfig{}, &MockSummarizer{})

	cm.updateMetrics(false, 1000, 100, 2*time.Second)

	metrics := cm.GetMetrics()

	if metrics.TotalCompactions != 1 {
		t.Errorf("Expected 1 compaction, got %d", metrics.TotalCompactions)
	}
	if metrics.FailureCount != 1 {
		t.Errorf("Expected 1 failure, got %d", metrics.FailureCount)
	}
	if metrics.SuccessCount != 0 {
		t.Errorf("Expected 0 success, got %d", metrics.SuccessCount)
	}
}

func TestConcurrentCompactionPrevention(t *testing.T) {
	config := CompactionConfig{
		Enabled:          true,
		TriggerThreshold: 40,
		KeepRecentCount:  15,
		CompactBatchSize: 25,
		MinInterval:      1 * time.Millisecond,
	}
	cm := NewCompactionManager(config, &MockSummarizer{delay: 100 * time.Millisecond})

	session := NewSession(12345, "team1", 50)
	// Add 45 messages
	for i := 0; i < 45; i++ {
		session.AddMessage("user", "test message", nil)
	}

	ctx := context.Background()

	// Start first compaction
	go cm.CompactAsync(ctx, session)

	// Give it time to lock
	time.Sleep(10 * time.Millisecond)

	// Try second compaction (should be prevented)
	shouldCompact := cm.ShouldCompact(session)

	if shouldCompact {
		t.Error("Expected ShouldCompact to return false when already compacting")
	}

	// Wait for completion
	cm.Stop()
}

func TestStop(t *testing.T) {
	cm := NewCompactionManager(CompactionConfig{}, &MockSummarizer{})

	// Start some async work
	session := NewSession(12345, "team1", 50)
	for i := 0; i < 45; i++ {
		session.AddMessage("user", "test message", nil)
	}

	ctx := context.Background()
	cm.CompactAsync(ctx, session)

	// Stop should wait for completion
	cm.Stop()

	// Should not panic or hang
}
