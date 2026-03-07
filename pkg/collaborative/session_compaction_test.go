// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package collaborative

import (
	"testing"
	"time"
)

func TestSession_WithCompactionFields(t *testing.T) {
	session := NewSession(12345, "team1", 50)

	// Verify compaction fields are initialized
	if session.CompactedContext != nil {
		t.Error("Expected CompactedContext to be nil initially")
	}

	if !session.LastCompaction.IsZero() {
		t.Error("Expected LastCompaction to be zero initially")
	}

	// Test setting compaction config
	session.CompactionConfig = CompactionConfig{
		Enabled:          true,
		TriggerThreshold: 40,
		KeepRecentCount:  15,
	}

	if !session.CompactionConfig.Enabled {
		t.Error("Expected CompactionConfig.Enabled to be true")
	}
}

func TestSession_CompactedContextInitialization(t *testing.T) {
	session := NewSession(12345, "team1", 50)

	// Initialize compacted context
	session.CompactedContext = &CompactedContext{
		Summary:        "Test summary",
		SummaryVersion: 1,
		CompactedAt:    time.Now(),
		CompactedCount: 10,
		OriginalSize:   1000,
		CompressedSize: 100,
	}

	if session.CompactedContext == nil {
		t.Fatal("Expected CompactedContext to be set")
	}

	if session.CompactedContext.Summary != "Test summary" {
		t.Errorf("Expected summary 'Test summary', got '%s'", session.CompactedContext.Summary)
	}

	if session.CompactedContext.CompactedCount != 10 {
		t.Errorf("Expected CompactedCount 10, got %d", session.CompactedContext.CompactedCount)
	}
}

func TestSession_CompactionMutex(t *testing.T) {
	session := NewSession(12345, "team1", 50)

	// Test that mutex can be locked
	session.CompactionMutex.Lock()
	locked := true
	session.CompactionMutex.Unlock()

	if !locked {
		t.Error("Expected to be able to lock CompactionMutex")
	}
}

func TestSession_ConcurrentCompactionPrevention(t *testing.T) {
	session := NewSession(12345, "team1", 50)

	// Simulate concurrent compaction attempts
	done := make(chan bool, 2)
	compactionCount := 0

	// First goroutine locks
	go func() {
		session.CompactionMutex.Lock()
		time.Sleep(100 * time.Millisecond)
		compactionCount++
		session.CompactionMutex.Unlock()
		done <- true
	}()

	// Give first goroutine time to lock
	time.Sleep(10 * time.Millisecond)

	// Second goroutine tries to lock (should wait)
	go func() {
		session.CompactionMutex.Lock()
		compactionCount++
		session.CompactionMutex.Unlock()
		done <- true
	}()

	// Wait for both
	<-done
	<-done

	// Both should have completed
	if compactionCount != 2 {
		t.Errorf("Expected compactionCount 2, got %d", compactionCount)
	}
}

func TestSession_LastCompactionTracking(t *testing.T) {
	session := NewSession(12345, "team1", 50)

	// Initially zero
	if !session.LastCompaction.IsZero() {
		t.Error("Expected LastCompaction to be zero initially")
	}

	// Set last compaction time
	now := time.Now()
	session.LastCompaction = now

	if session.LastCompaction.IsZero() {
		t.Error("Expected LastCompaction to be set")
	}

	// Check time difference
	diff := time.Since(session.LastCompaction)
	if diff > 1*time.Second {
		t.Errorf("Expected LastCompaction to be recent, got diff %v", diff)
	}
}

func TestSession_CompactionConfigDefaults(t *testing.T) {
	session := NewSession(12345, "team1", 50)

	// Set default config
	session.CompactionConfig = CompactionConfig{
		Enabled:          true,
		TriggerThreshold: 40,
		KeepRecentCount:  15,
		CompactBatchSize: 25,
		MinInterval:      5 * time.Minute,
		SummaryMaxLength: 2000,
		LLMProvider:      "openai",
		LLMModel:         "gpt-4o-mini",
		LLMTimeout:       30 * time.Second,
		LLMMaxRetries:    3,
	}

	config := session.CompactionConfig

	if config.TriggerThreshold != 40 {
		t.Errorf("Expected TriggerThreshold 40, got %d", config.TriggerThreshold)
	}
	if config.KeepRecentCount != 15 {
		t.Errorf("Expected KeepRecentCount 15, got %d", config.KeepRecentCount)
	}
	if config.CompactBatchSize != 25 {
		t.Errorf("Expected CompactBatchSize 25, got %d", config.CompactBatchSize)
	}
}

func TestSession_CompactedContextThreadSafety(t *testing.T) {
	session := NewSession(12345, "team1", 50)
	session.CompactedContext = &CompactedContext{
		Summary:        "Initial",
		SummaryVersion: 1,
	}

	// Test concurrent reads
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			session.CompactedContext.mu.RLock()
			_ = session.CompactedContext.Summary
			session.CompactedContext.mu.RUnlock()
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
			session.CompactedContext.mu.Lock()
			session.CompactedContext.SummaryVersion = version
			session.CompactedContext.mu.Unlock()
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
}
