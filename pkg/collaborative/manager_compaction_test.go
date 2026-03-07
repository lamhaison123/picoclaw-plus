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

// TestManagerV2WithCompaction tests ManagerV2 with compaction enabled
func TestManagerV2WithCompaction(t *testing.T) {
	provider := &MockLLMProvider{
		responseText: "Summary of conversation",
	}

	config := &Config{
		MentionQueueSize:    20,
		MentionRateLimit:    2 * time.Second,
		MentionMaxRetries:   3,
		MentionRetryBackoff: 1 * time.Second,
		CompactionEnabled:   true,
		LLMProvider:         provider,
		CompactionConfig: CompactionConfig{
			TriggerThreshold: 10, // Low threshold for testing
			KeepRecentCount:  5,
			CompactBatchSize: 5,
			MinInterval:      1 * time.Millisecond,
			SummaryMaxLength: 2000,
		},
	}

	manager := NewManagerV2WithConfig(config)
	defer manager.Stop()

	if manager.GetCompactionManager() == nil {
		t.Fatal("Expected compaction manager to be initialized")
	}

	// Create session
	session := manager.GetOrCreateSession(12345, "team1", 50)

	// Add messages to trigger compaction
	for i := 0; i < 12; i++ {
		session.AddMessage("user", "Test message", nil)
	}

	// Check if compaction should trigger
	if !manager.GetCompactionManager().ShouldCompact(session) {
		t.Error("Expected ShouldCompact to return true")
	}

	// Trigger compaction
	ctx := context.Background()
	err := manager.GetCompactionManager().Compact(ctx, session)
	if err != nil {
		t.Errorf("Compaction failed: %v", err)
	}

	// Verify compaction happened
	if session.CompactedContext == nil {
		t.Error("Expected CompactedContext to be set")
	}

	if len(session.Context) != 7 { // 12 - 5 = 7
		t.Errorf("Expected 7 messages after compaction, got %d", len(session.Context))
	}
}

// TestManagerV2WithoutCompaction tests ManagerV2 with compaction disabled
func TestManagerV2WithoutCompaction(t *testing.T) {
	config := &Config{
		MentionQueueSize:  20,
		MentionRateLimit:  2 * time.Second,
		CompactionEnabled: false, // Disabled
	}

	manager := NewManagerV2WithConfig(config)
	defer manager.Stop()

	if manager.GetCompactionManager() != nil {
		t.Error("Expected compaction manager to be nil when disabled")
	}
}

// TestManagerV2CompactionDefaults tests default compaction config
func TestManagerV2CompactionDefaults(t *testing.T) {
	provider := &MockLLMProvider{}

	config := &Config{
		CompactionEnabled: true,
		LLMProvider:       provider,
		// Don't set CompactionConfig fields - test defaults
	}

	manager := NewManagerV2WithConfig(config)
	defer manager.Stop()

	if manager.GetCompactionManager() == nil {
		t.Fatal("Expected compaction manager to be initialized")
	}

	// Check defaults were set
	cfg := manager.config.CompactionConfig
	if cfg.TriggerThreshold != 40 {
		t.Errorf("Expected default TriggerThreshold 40, got %d", cfg.TriggerThreshold)
	}
	if cfg.KeepRecentCount != 15 {
		t.Errorf("Expected default KeepRecentCount 15, got %d", cfg.KeepRecentCount)
	}
	if cfg.CompactBatchSize != 25 {
		t.Errorf("Expected default CompactBatchSize 25, got %d", cfg.CompactBatchSize)
	}
	if cfg.MinInterval != 5*time.Minute {
		t.Errorf("Expected default MinInterval 5m, got %v", cfg.MinInterval)
	}
	if cfg.SummaryMaxLength != 2000 {
		t.Errorf("Expected default SummaryMaxLength 2000, got %d", cfg.SummaryMaxLength)
	}
	if cfg.LLMModel != "gpt-4o-mini" {
		t.Errorf("Expected default LLMModel 'gpt-4o-mini', got '%s'", cfg.LLMModel)
	}
}

// TestManagerV2CompactionStop tests graceful shutdown with compaction
func TestManagerV2CompactionStop(t *testing.T) {
	provider := &MockLLMProvider{
		delay: 50 * time.Millisecond,
	}

	config := &Config{
		CompactionEnabled: true,
		LLMProvider:       provider,
		CompactionConfig: CompactionConfig{
			TriggerThreshold: 10,
			KeepRecentCount:  5,
			MinInterval:      1 * time.Millisecond,
		},
	}

	manager := NewManagerV2WithConfig(config)

	// Create session and add messages
	session := manager.GetOrCreateSession(12345, "team1", 50)
	for i := 0; i < 12; i++ {
		session.AddMessage("user", "Test", nil)
	}

	// Trigger async compaction
	ctx := context.Background()
	manager.GetCompactionManager().CompactAsync(ctx, session)

	// Stop should wait for compaction to complete
	manager.Stop()

	// Should not panic or hang
}

// TestManagerV2CompactionMetrics tests metrics collection
func TestManagerV2CompactionMetrics(t *testing.T) {
	provider := &MockLLMProvider{
		responseText: "Test summary",
	}

	config := &Config{
		CompactionEnabled: true,
		LLMProvider:       provider,
		CompactionConfig: CompactionConfig{
			TriggerThreshold: 10,
			KeepRecentCount:  5,
			CompactBatchSize: 5,
			MinInterval:      1 * time.Millisecond,
		},
	}

	manager := NewManagerV2WithConfig(config)
	defer manager.Stop()

	// Create session and trigger compaction
	session := manager.GetOrCreateSession(12345, "team1", 50)
	for i := 0; i < 12; i++ {
		session.AddMessage("user", "Test message", nil)
	}

	ctx := context.Background()
	err := manager.GetCompactionManager().Compact(ctx, session)
	if err != nil {
		t.Fatalf("Compaction failed: %v", err)
	}

	// Check metrics
	metrics := manager.GetCompactionManager().GetMetrics()
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

// TestManagerV2CompactionWithNilProvider tests behavior with nil provider
func TestManagerV2CompactionWithNilProvider(t *testing.T) {
	config := &Config{
		CompactionEnabled: true,
		LLMProvider:       nil, // Nil provider
	}

	manager := NewManagerV2WithConfig(config)
	defer manager.Stop()

	// Should not initialize compaction manager without provider
	if manager.GetCompactionManager() != nil {
		t.Error("Expected compaction manager to be nil with nil provider")
	}
}

// TestManagerV2GetOrCreateSessionWithCompaction tests session creation with compaction config
func TestManagerV2GetOrCreateSessionWithCompaction(t *testing.T) {
	provider := &MockLLMProvider{}

	config := &Config{
		CompactionEnabled: true,
		LLMProvider:       provider,
		CompactionConfig: CompactionConfig{
			TriggerThreshold: 40,
		},
	}

	manager := NewManagerV2WithConfig(config)
	defer manager.Stop()

	session := manager.GetOrCreateSession(12345, "team1", 50)

	// Session should have compaction fields initialized
	if session.CompactedContext != nil {
		t.Error("Expected CompactedContext to be nil initially")
	}

	if !session.LastCompaction.IsZero() {
		t.Error("Expected LastCompaction to be zero initially")
	}

	// CompactionConfig should be set from manager config
	// (Note: Currently not set automatically, but could be added)
}

// TestManagerV2CompactionContextFormatting tests context string with compaction
func TestManagerV2CompactionContextFormatting(t *testing.T) {
	provider := &MockLLMProvider{
		responseText: "Summary of early messages",
	}

	config := &Config{
		CompactionEnabled: true,
		LLMProvider:       provider,
		CompactionConfig: CompactionConfig{
			TriggerThreshold: 10,
			KeepRecentCount:  5,
			CompactBatchSize: 5,
			MinInterval:      1 * time.Millisecond,
		},
	}

	manager := NewManagerV2WithConfig(config)
	defer manager.Stop()

	session := manager.GetOrCreateSession(12345, "team1", 50)

	// Add messages
	for i := 0; i < 12; i++ {
		session.AddMessage("user", "Message "+string(rune(i)), nil)
	}

	// Compact
	ctx := context.Background()
	err := manager.GetCompactionManager().Compact(ctx, session)
	if err != nil {
		t.Fatalf("Compaction failed: %v", err)
	}

	// Get context string
	contextStr := session.GetContextAsString()

	// Should contain summary
	if len(contextStr) == 0 {
		t.Error("Expected non-empty context string")
	}

	// Should have summary section
	if !containsStr(contextStr, "=== Context Summary ===") {
		t.Error("Expected summary section in context")
	}

	// Should have recent messages
	if !containsStr(contextStr, "=== Recent Messages ===") {
		t.Error("Expected recent messages section")
	}
}

// Helper function
func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && findSubstr(s, substr)
}

func findSubstr(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
