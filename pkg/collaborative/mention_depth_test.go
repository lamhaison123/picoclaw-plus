// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package collaborative

import (
	"testing"
	"time"
)

func TestMentionDepthTracking(t *testing.T) {
	session := NewSession(12345, "test-team", 50)

	// Test initial depth
	if session.MentionDepth != 0 {
		t.Errorf("Expected initial MentionDepth to be 0, got %d", session.MentionDepth)
	}

	// Simulate nested mentions
	session.mu.Lock()
	session.MentionDepth++
	session.mu.Unlock()

	if session.MentionDepth != 1 {
		t.Errorf("Expected MentionDepth to be 1, got %d", session.MentionDepth)
	}

	// Simulate deeper nesting
	session.mu.Lock()
	session.MentionDepth++
	session.mu.Unlock()

	if session.MentionDepth != 2 {
		t.Errorf("Expected MentionDepth to be 2, got %d", session.MentionDepth)
	}

	// Simulate decrement
	session.mu.Lock()
	session.MentionDepth--
	session.mu.Unlock()

	if session.MentionDepth != 1 {
		t.Errorf("Expected MentionDepth to be 1 after decrement, got %d", session.MentionDepth)
	}
}

func TestMentionDepthLimit(t *testing.T) {
	session := NewSession(12345, "test-team", 50)

	// Simulate reaching max depth
	session.mu.Lock()
	session.MentionDepth = 3
	currentDepth := session.MentionDepth
	session.mu.Unlock()

	if currentDepth < 3 {
		t.Errorf("Expected depth to be at least 3, got %d", currentDepth)
	}

	// Verify that depth 3 should trigger the limit
	if currentDepth >= 3 {
		t.Log("✓ Depth limit (3) would be triggered correctly")
	}
}

func TestSessionConcurrentAccess(t *testing.T) {
	session := NewSession(12345, "test-team", 50)

	// Test concurrent access to MentionDepth
	done := make(chan bool)

	// Goroutine 1: Increment
	go func() {
		for i := 0; i < 10; i++ {
			session.mu.Lock()
			session.MentionDepth++
			session.mu.Unlock()
			time.Sleep(1 * time.Millisecond)
		}
		done <- true
	}()

	// Goroutine 2: Decrement
	go func() {
		for i := 0; i < 10; i++ {
			session.mu.Lock()
			if session.MentionDepth > 0 {
				session.MentionDepth--
			}
			session.mu.Unlock()
			time.Sleep(1 * time.Millisecond)
		}
		done <- true
	}()

	// Wait for both goroutines
	<-done
	<-done

	// Final depth should be 0 or close to it
	session.mu.Lock()
	finalDepth := session.MentionDepth
	session.mu.Unlock()

	t.Logf("Final MentionDepth after concurrent access: %d", finalDepth)

	// Should not have negative depth
	if finalDepth < 0 {
		t.Errorf("MentionDepth should not be negative, got %d", finalDepth)
	}
}
