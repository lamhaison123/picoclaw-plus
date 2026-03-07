// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package collaborative

import (
	"strings"
	"testing"
	"time"
)

func TestGetContextAsString_WithoutSummary(t *testing.T) {
	session := NewSession(12345, "team1", 50)

	// Add some messages
	session.AddMessage("user", "Hello team", nil)
	session.AddMessage("architect", "Hi, I'm the architect", nil)
	session.AddMessage("developer", "Ready to code", nil)

	contextStr := session.GetContextAsString()

	// Should contain header
	if !strings.Contains(contextStr, "=== Collaborative Chat Context ===") {
		t.Error("Expected context header")
	}

	// Should contain session info
	if !strings.Contains(contextStr, session.SessionID) {
		t.Error("Expected session ID in context")
	}

	// Should contain "Conversation History" (not "Recent Messages")
	if !strings.Contains(contextStr, "=== Conversation History ===") {
		t.Error("Expected 'Conversation History' header")
	}

	// Should NOT contain summary section
	if strings.Contains(contextStr, "=== Context Summary ===") {
		t.Error("Should not contain summary section when no summary present")
	}

	// Should contain messages
	if !strings.Contains(contextStr, "Hello team") {
		t.Error("Expected user message in context")
	}
	if !strings.Contains(contextStr, "architect") {
		t.Error("Expected architect message in context")
	}
}

func TestGetContextAsString_WithSummary(t *testing.T) {
	session := NewSession(12345, "team1", 50)

	// Add compacted context
	session.CompactedContext = &CompactedContext{
		Summary:        "Previous discussion about REST API design with JWT authentication",
		SummaryVersion: 1,
		CompactedAt:    time.Now(),
		CompactedCount: 25,
		OriginalSize:   5000,
		CompressedSize: 100,
	}

	// Add recent messages
	session.AddMessage("user", "Let's add rate limiting", nil)
	session.AddMessage("architect", "I'll design the rate limiter", nil)

	contextStr := session.GetContextAsString()

	// Should contain summary section
	if !strings.Contains(contextStr, "=== Context Summary ===") {
		t.Error("Expected summary section")
	}

	// Should contain summary text
	if !strings.Contains(contextStr, "REST API design") {
		t.Error("Expected summary content")
	}

	// Should show compacted count
	if !strings.Contains(contextStr, "25 earlier messages") {
		t.Error("Expected compacted count in summary")
	}

	// Should contain "Recent Messages" header
	if !strings.Contains(contextStr, "=== Recent Messages ===") {
		t.Error("Expected 'Recent Messages' header when summary present")
	}

	// Should contain recent messages
	if !strings.Contains(contextStr, "rate limiting") {
		t.Error("Expected recent message in context")
	}
}

func TestGetContextAsString_WithEmptySummary(t *testing.T) {
	session := NewSession(12345, "team1", 50)

	// Add compacted context with empty summary
	session.CompactedContext = &CompactedContext{
		Summary:        "", // Empty summary
		SummaryVersion: 0,
		CompactedCount: 0,
	}

	// Add messages
	session.AddMessage("user", "Hello", nil)

	contextStr := session.GetContextAsString()

	// Should NOT contain summary section when summary is empty
	if strings.Contains(contextStr, "=== Context Summary ===") {
		t.Error("Should not contain summary section when summary is empty")
	}

	// Should contain "Conversation History" instead
	if !strings.Contains(contextStr, "=== Conversation History ===") {
		t.Error("Expected 'Conversation History' header when summary is empty")
	}
}

func TestGetContextAsString_ThreadSafety(t *testing.T) {
	session := NewSession(12345, "team1", 50)

	// Add initial messages
	for i := 0; i < 10; i++ {
		session.AddMessage("user", "test message", nil)
	}

	// Add compacted context
	session.CompactedContext = &CompactedContext{
		Summary:        "Test summary",
		CompactedCount: 5,
	}

	// Concurrent reads
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			_ = session.GetContextAsString()
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Should not panic or race
}

func TestGetContextAsString_Format(t *testing.T) {
	session := NewSession(12345, "team1", 50)

	// Add messages with different roles
	session.AddMessage("user", "Hello", nil)
	session.AddMessage("architect", "Hi", nil)
	session.AddMessage("developer", "Ready", nil)

	contextStr := session.GetContextAsString()

	// Check format includes timestamps
	if !strings.Contains(contextStr, "[") {
		t.Error("Expected timestamp brackets in format")
	}

	// Check format includes role names in uppercase
	if !strings.Contains(contextStr, "ARCHITECT") {
		t.Error("Expected uppercase role name")
	}
	if !strings.Contains(contextStr, "DEVELOPER") {
		t.Error("Expected uppercase role name")
	}

	// User messages should not have emoji prefix (checked by absence of double space)
	lines := strings.Split(contextStr, "\n")
	for _, line := range lines {
		if strings.Contains(line, "USER:") {
			// Should be "[timestamp] USER:" not "[timestamp]  USER:"
			if strings.Contains(line, "]  USER:") {
				t.Error("User messages should not have emoji prefix")
			}
		}
	}
}

func TestGetContextAsString_SummaryVersionTracking(t *testing.T) {
	session := NewSession(12345, "team1", 50)

	// Add compacted context with version
	session.CompactedContext = &CompactedContext{
		Summary:        "Version 1 summary",
		SummaryVersion: 1,
		CompactedCount: 10,
	}

	contextStr1 := session.GetContextAsString()

	// Update summary
	session.CompactedContext.mu.Lock()
	session.CompactedContext.Summary = "Version 2 summary"
	session.CompactedContext.SummaryVersion = 2
	session.CompactedContext.mu.Unlock()

	contextStr2 := session.GetContextAsString()

	// Should contain different summaries
	if !strings.Contains(contextStr1, "Version 1") {
		t.Error("Expected version 1 summary in first context")
	}
	if !strings.Contains(contextStr2, "Version 2") {
		t.Error("Expected version 2 summary in second context")
	}
}

func TestGetContextAsString_LargeCompactedCount(t *testing.T) {
	session := NewSession(12345, "team1", 50)

	// Add compacted context with large count
	session.CompactedContext = &CompactedContext{
		Summary:        "Summary of many messages",
		CompactedCount: 1000,
	}

	contextStr := session.GetContextAsString()

	// Should show large count
	if !strings.Contains(contextStr, "1000 earlier messages") {
		t.Error("Expected large compacted count in summary")
	}
}
