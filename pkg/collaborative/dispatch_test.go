// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package collaborative

import (
	"strings"
	"sync"
	"testing"
)

func TestDispatchTracker(t *testing.T) {
	tracker := NewDispatchTracker()

	// Test initial state
	if tracker.Size() != 0 {
		t.Errorf("Expected initial size 0, got %d", tracker.Size())
	}

	// Test marking as dispatched
	msgID := "test-message-1"
	if tracker.IsDispatched(msgID) {
		t.Error("Message should not be dispatched initially")
	}

	tracker.MarkDispatched(msgID)
	if !tracker.IsDispatched(msgID) {
		t.Error("Message should be marked as dispatched")
	}

	if tracker.Size() != 1 {
		t.Errorf("Expected size 1, got %d", tracker.Size())
	}

	// Test duplicate marking (idempotent)
	tracker.MarkDispatched(msgID)
	if tracker.Size() != 1 {
		t.Errorf("Expected size 1 after duplicate mark, got %d", tracker.Size())
	}

	// Test clearing
	tracker.Clear()
	if tracker.Size() != 0 {
		t.Errorf("Expected size 0 after clear, got %d", tracker.Size())
	}

	if tracker.IsDispatched(msgID) {
		t.Error("Message should not be dispatched after clear")
	}
}

func TestGenerateMessageID(t *testing.T) {
	tests := []struct {
		name      string
		chatID    int64
		sessionID string
		role      string
		content   string
	}{
		{
			name:      "Basic message",
			chatID:    12345,
			sessionID: "session1",
			role:      "developer",
			content:   "@developer please fix this",
		},
		{
			name:      "Long content",
			chatID:    12345,
			sessionID: "session1",
			role:      "developer",
			content:   strings.Repeat("This is a long message. ", 50), // Long but meaningful content
		},
		{
			name:      "UTF-8 content",
			chatID:    12345,
			sessionID: "session1",
			role:      "developer",
			content:   "Cảm ơn @developer đã hoàn thành",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msgID := GenerateMessageID(tt.chatID, tt.sessionID, tt.role, tt.content)

			// Check format
			if msgID == "" {
				t.Error("Message ID should not be empty")
			}

			// Check consistency (same input = same output)
			msgID2 := GenerateMessageID(tt.chatID, tt.sessionID, tt.role, tt.content)
			if msgID != msgID2 {
				t.Errorf("Message IDs should be consistent: %s != %s", msgID, msgID2)
			}

			// Check uniqueness (different content = different ID)
			// For long content, change at the beginning to ensure difference is detected
			differentContent := tt.content
			if len(tt.content) > 100 {
				differentContent = "DIFFERENT" + tt.content[9:] // Change first 9 chars
			} else {
				differentContent = tt.content + "different"
			}
			msgID3 := GenerateMessageID(tt.chatID, tt.sessionID, tt.role, differentContent)
			if msgID == msgID3 {
				t.Errorf("Different content should produce different message IDs: %s == %s", msgID, msgID3)
			}
		})
	}
}

func TestDispatchTrackerConcurrency(t *testing.T) {
	tracker := NewDispatchTracker()
	const numGoroutines = 100
	const numMessages = 10

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Concurrent marking
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numMessages; j++ {
				msgID := GenerateMessageID(int64(id), "session", "developer", string(rune(j)))
				tracker.MarkDispatched(msgID)
			}
		}(i)
	}

	wg.Wait()

	// Check that all messages were tracked
	expectedSize := numGoroutines * numMessages
	if tracker.Size() != expectedSize {
		t.Errorf("Expected size %d, got %d", expectedSize, tracker.Size())
	}

	// Concurrent reading
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numMessages; j++ {
				msgID := GenerateMessageID(int64(id), "session", "developer", string(rune(j)))
				if !tracker.IsDispatched(msgID) {
					t.Errorf("Message %s should be dispatched", msgID)
				}
			}
		}(i)
	}

	wg.Wait()
}

func TestDispatchTrackerDuplicatePrevention(t *testing.T) {
	tracker := NewDispatchTracker()

	chatID := int64(12345)
	sessionID := "session1"
	role := "developer"
	content := "@developer please help"

	msgID := GenerateMessageID(chatID, sessionID, role, content)

	// First dispatch should succeed
	if tracker.IsDispatched(msgID) {
		t.Error("Message should not be dispatched initially")
	}

	tracker.MarkDispatched(msgID)

	// Second dispatch should be detected as duplicate
	if !tracker.IsDispatched(msgID) {
		t.Error("Duplicate dispatch should be detected")
	}

	// Same content, different role = different message
	msgID2 := GenerateMessageID(chatID, sessionID, "tester", content)
	if tracker.IsDispatched(msgID2) {
		t.Error("Different role should create different message ID")
	}

	// Same content, different session = different message
	msgID3 := GenerateMessageID(chatID, "session2", role, content)
	if tracker.IsDispatched(msgID3) {
		t.Error("Different session should create different message ID")
	}
}

func BenchmarkGenerateMessageID(b *testing.B) {
	chatID := int64(12345)
	sessionID := "session1"
	role := "developer"
	content := "@developer please fix this bug in the code"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GenerateMessageID(chatID, sessionID, role, content)
	}
}

func BenchmarkDispatchTrackerMark(b *testing.B) {
	tracker := NewDispatchTracker()
	chatID := int64(12345)
	sessionID := "session1"
	role := "developer"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		msgID := GenerateMessageID(chatID, sessionID, role, string(rune(i)))
		tracker.MarkDispatched(msgID)
	}
}

func BenchmarkDispatchTrackerCheck(b *testing.B) {
	tracker := NewDispatchTracker()
	chatID := int64(12345)
	sessionID := "session1"
	role := "developer"

	// Pre-populate with some messages
	for i := 0; i < 1000; i++ {
		msgID := GenerateMessageID(chatID, sessionID, role, string(rune(i)))
		tracker.MarkDispatched(msgID)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		msgID := GenerateMessageID(chatID, sessionID, role, string(rune(i%1000)))
		tracker.IsDispatched(msgID)
	}
}
