// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package collaborative

import (
	"context"
	"testing"
	"time"
)

// TestCompactionEndToEnd tests the complete compaction flow
func TestCompactionEndToEnd(t *testing.T) {
	// Setup
	config := CompactionConfig{
		Enabled:          true,
		TriggerThreshold: 40,
		KeepRecentCount:  15,
		CompactBatchSize: 25,
		MinInterval:      1 * time.Millisecond,
		SummaryMaxLength: 2000,
		LLMModel:         "gpt-4o-mini",
		LLMTimeout:       30 * time.Second,
		LLMMaxRetries:    3,
	}

	provider := &MockLLMProvider{
		responseText: "## Project Overview\nBuilding REST API with authentication\n\n## Key Decisions\n- Use JWT tokens\n- PostgreSQL database",
	}

	summarizer := NewLLMSummarizer(config, provider)
	cm := NewCompactionManager(config, summarizer)

	session := NewSession(12345, "team1", 50)

	// Add 45 messages to trigger compaction
	for i := 0; i < 45; i++ {
		session.AddMessage("user", "Message "+string(rune(i)), nil)
	}

	// Verify should compact
	if !cm.ShouldCompact(session) {
		t.Error("Expected ShouldCompact to return true")
	}

	// Perform compaction
	ctx := context.Background()
	err := cm.Compact(ctx, session)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify session was updated
	if session.CompactedContext == nil {
		t.Fatal("Expected CompactedContext to be set")
	}

	if session.CompactedContext.Summary == "" {
		t.Error("Expected non-empty summary")
	}

	if session.CompactedContext.CompactedCount != 25 {
		t.Errorf("Expected 25 compacted messages, got %d", session.CompactedContext.CompactedCount)
	}

	// Verify messages were removed
	if len(session.Context) != 20 { // 45 - 25 = 20
		t.Errorf("Expected 20 remaining messages, got %d", len(session.Context))
	}

	// Verify metrics
	metrics := cm.GetMetrics()
	if metrics.TotalCompactions != 1 {
		t.Errorf("Expected 1 compaction, got %d", metrics.TotalCompactions)
	}
	if metrics.SuccessCount != 1 {
		t.Errorf("Expected 1 success, got %d", metrics.SuccessCount)
	}
	if metrics.AverageCompression <= 0 {
		t.Error("Expected positive compression ratio")
	}
}

// TestCompactionWithMultipleRounds tests multiple compaction rounds
func TestCompactionWithMultipleRounds(t *testing.T) {
	config := CompactionConfig{
		Enabled:          true,
		TriggerThreshold: 40,
		KeepRecentCount:  15,
		CompactBatchSize: 25,
		MinInterval:      1 * time.Millisecond,
		SummaryMaxLength: 2000,
		LLMModel:         "gpt-4o-mini",
		LLMTimeout:       30 * time.Second,
		LLMMaxRetries:    3,
	}

	provider := &MockLLMProvider{
		responseText: "Round 1 summary",
	}

	summarizer := NewLLMSummarizer(config, provider)
	cm := NewCompactionManager(config, summarizer)

	session := NewSession(12345, "team1", 50)

	// First round: Add 45 messages
	for i := 0; i < 45; i++ {
		session.AddMessage("user", "Round 1 message", nil)
	}

	ctx := context.Background()
	err := cm.Compact(ctx, session)
	if err != nil {
		t.Fatalf("Round 1 compaction failed: %v", err)
	}

	firstSummary := session.CompactedContext.Summary
	firstVersion := session.CompactedContext.SummaryVersion

	// Add more messages for second round
	time.Sleep(2 * time.Millisecond) // Wait for MinInterval
	for i := 0; i < 30; i++ {
		session.AddMessage("user", "Round 2 message", nil)
	}

	// Update provider response for round 2
	provider.responseText = "Round 2 summary (updated)"

	err = cm.Compact(ctx, session)
	if err != nil {
		t.Fatalf("Round 2 compaction failed: %v", err)
	}

	// Verify summary was updated
	if session.CompactedContext.Summary == firstSummary {
		t.Error("Expected summary to be updated in round 2")
	}

	if session.CompactedContext.SummaryVersion <= firstVersion {
		t.Error("Expected summary version to increment")
	}

	// Verify metrics
	metrics := cm.GetMetrics()
	if metrics.TotalCompactions != 2 {
		t.Errorf("Expected 2 compactions, got %d", metrics.TotalCompactions)
	}
}

// TestCompactionCompressionRatio tests compression ratio
func TestCompactionCompressionRatio(t *testing.T) {
	config := CompactionConfig{
		Enabled:          true,
		TriggerThreshold: 40,
		KeepRecentCount:  15,
		CompactBatchSize: 25,
		MinInterval:      1 * time.Millisecond,
		SummaryMaxLength: 200, // Small max to ensure good compression
		LLMModel:         "gpt-4o-mini",
		LLMTimeout:       30 * time.Second,
		LLMMaxRetries:    3,
	}

	provider := &MockLLMProvider{
		responseText: "Short summary",
	}

	summarizer := NewLLMSummarizer(config, provider)
	cm := NewCompactionManager(config, summarizer)

	session := NewSession(12345, "team1", 50)

	// Add messages with substantial content
	for i := 0; i < 45; i++ {
		session.AddMessage("user", "This is a longer message with more content to compress", nil)
	}

	ctx := context.Background()
	err := cm.Compact(ctx, session)
	if err != nil {
		t.Fatalf("Compaction failed: %v", err)
	}

	// Check compression ratio
	originalSize := session.CompactedContext.OriginalSize
	compressedSize := session.CompactedContext.CompressedSize

	if originalSize <= compressedSize {
		t.Errorf("Expected compression: original %d should be > compressed %d", originalSize, compressedSize)
	}

	ratio := float64(originalSize) / float64(compressedSize)
	if ratio < 2.0 {
		t.Errorf("Expected compression ratio >= 2.0, got %.2f", ratio)
	}

	t.Logf("Compression ratio: %.2f:1 (original: %d bytes, compressed: %d bytes)",
		ratio, originalSize, compressedSize)
}

// TestCompactionContextFormatting tests context string with summary
func TestCompactionContextFormatting(t *testing.T) {
	config := CompactionConfig{
		Enabled:          true,
		TriggerThreshold: 40,
		KeepRecentCount:  15,
		CompactBatchSize: 25,
		MinInterval:      1 * time.Millisecond,
		SummaryMaxLength: 2000,
		LLMModel:         "gpt-4o-mini",
		LLMTimeout:       30 * time.Second,
		LLMMaxRetries:    3,
	}

	provider := &MockLLMProvider{
		responseText: "Summary of early discussion about REST API design",
	}

	summarizer := NewLLMSummarizer(config, provider)
	cm := NewCompactionManager(config, summarizer)

	session := NewSession(12345, "team1", 50)

	// Add messages
	for i := 0; i < 45; i++ {
		session.AddMessage("user", "Message "+string(rune(i)), nil)
	}

	// Compact
	ctx := context.Background()
	err := cm.Compact(ctx, session)
	if err != nil {
		t.Fatalf("Compaction failed: %v", err)
	}

	// Get context string
	contextStr := session.GetContextAsString()

	// Verify format
	if contextStr == "" {
		t.Fatal("Expected non-empty context string")
	}

	// Should contain summary section
	if !contains(contextStr, "=== Context Summary ===") {
		t.Error("Expected summary section in context")
	}

	// Should contain summary content
	if !contains(contextStr, "REST API design") {
		t.Error("Expected summary content in context")
	}

	// Should show compacted count
	if !contains(contextStr, "25 earlier messages") {
		t.Error("Expected compacted count in context")
	}

	// Should contain recent messages header
	if !contains(contextStr, "=== Recent Messages ===") {
		t.Error("Expected recent messages header")
	}

	t.Logf("Context string length: %d characters", len(contextStr))
}

// TestCompactionErrorRecovery tests error handling and recovery
func TestCompactionErrorRecovery(t *testing.T) {
	config := CompactionConfig{
		Enabled:          true,
		TriggerThreshold: 40,
		KeepRecentCount:  15,
		CompactBatchSize: 25,
		MinInterval:      1 * time.Millisecond,
		SummaryMaxLength: 2000,
		LLMModel:         "gpt-4o-mini",
		LLMTimeout:       30 * time.Second,
		LLMMaxRetries:    2,
	}

	provider := &MockLLMProvider{
		shouldFail: true,
	}

	summarizer := NewLLMSummarizer(config, provider)
	cm := NewCompactionManager(config, summarizer)

	session := NewSession(12345, "team1", 50)

	// Add messages
	for i := 0; i < 45; i++ {
		session.AddMessage("user", "Message", nil)
	}

	// First compaction should fail
	ctx := context.Background()
	err := cm.Compact(ctx, session)
	if err == nil {
		t.Error("Expected error on first compaction")
	}

	// Verify session state is unchanged
	if session.CompactedContext != nil {
		t.Error("Expected CompactedContext to remain nil after failed compaction")
	}

	if len(session.Context) != 45 {
		t.Errorf("Expected 45 messages after failed compaction, got %d", len(session.Context))
	}

	// Verify metrics tracked failure
	metrics := cm.GetMetrics()
	if metrics.FailureCount != 1 {
		t.Errorf("Expected 1 failure, got %d", metrics.FailureCount)
	}

	// Now fix provider and retry
	provider.shouldFail = false
	provider.responseText = "Recovery summary"

	time.Sleep(2 * time.Millisecond) // Wait for MinInterval

	err = cm.Compact(ctx, session)
	if err != nil {
		t.Errorf("Expected success on retry, got %v", err)
	}

	// Verify recovery
	if session.CompactedContext == nil {
		t.Error("Expected CompactedContext to be set after recovery")
	}

	metrics = cm.GetMetrics()
	if metrics.SuccessCount != 1 {
		t.Errorf("Expected 1 success after recovery, got %d", metrics.SuccessCount)
	}
}

// TestCompactionAsync tests async compaction
func TestCompactionAsync(t *testing.T) {
	config := CompactionConfig{
		Enabled:          true,
		TriggerThreshold: 40,
		KeepRecentCount:  15,
		CompactBatchSize: 25,
		MinInterval:      1 * time.Millisecond,
		SummaryMaxLength: 2000,
		LLMModel:         "gpt-4o-mini",
		LLMTimeout:       30 * time.Second,
		LLMMaxRetries:    3,
	}

	provider := &MockLLMProvider{
		responseText: "Async summary",
		delay:        50 * time.Millisecond,
	}

	summarizer := NewLLMSummarizer(config, provider)
	cm := NewCompactionManager(config, summarizer)

	session := NewSession(12345, "team1", 50)

	// Add messages
	for i := 0; i < 45; i++ {
		session.AddMessage("user", "Message", nil)
	}

	// Trigger async compaction
	ctx := context.Background()
	cm.CompactAsync(ctx, session)

	// Should return immediately (not block)
	// Wait for completion
	cm.Stop()

	// Verify compaction happened
	if session.CompactedContext == nil {
		t.Error("Expected CompactedContext to be set after async compaction")
	}

	if session.CompactedContext.Summary != "Async summary" {
		t.Errorf("Expected 'Async summary', got '%s'", session.CompactedContext.Summary)
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
