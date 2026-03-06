// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package team

import (
	"sync"
	"testing"
	"time"
)

// TestSharedContextCreation tests SharedContext initialization
func TestSharedContextCreation(t *testing.T) {
	sc := NewSharedContext("team-001")

	if sc == nil {
		t.Fatal("Expected non-nil SharedContext")
	}

	if sc.teamID != "team-001" {
		t.Errorf("Expected teamID 'team-001', got '%s'", sc.teamID)
	}

	if len(sc.entries) != 0 {
		t.Errorf("Expected empty entries map, got %d entries", len(sc.entries))
	}

	if len(sc.history) != 0 {
		t.Errorf("Expected empty history, got %d entries", len(sc.history))
	}
}

// TestSharedContextSetGet tests basic set and get operations
func TestSharedContextSetGet(t *testing.T) {
	sc := NewSharedContext("team-001")

	// Set a value
	err := sc.Set("key1", "value1", "agent-001")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Get the value
	value, exists := sc.Get("key1")
	if !exists {
		t.Fatal("Expected key1 to exist")
	}

	if value != "value1" {
		t.Errorf("Expected value 'value1', got '%v'", value)
	}

	// Get non-existent key
	_, exists = sc.Get("nonexistent")
	if exists {
		t.Error("Expected nonexistent key to not exist")
	}
}

// TestSharedContextGetAll tests retrieving all entries
func TestSharedContextGetAll(t *testing.T) {
	sc := NewSharedContext("team-001")

	// Add multiple entries
	sc.Set("key1", "value1", "agent-001")
	sc.Set("key2", "value2", "agent-002")
	sc.Set("key3", "value3", "agent-003")

	all := sc.GetAll()

	if len(all) != 3 {
		t.Errorf("Expected 3 entries, got %d", len(all))
	}

	if all["key1"].Value != "value1" {
		t.Errorf("Expected key1 value 'value1', got '%v'", all["key1"].Value)
	}

	if all["key1"].AgentID != "agent-001" {
		t.Errorf("Expected key1 agentID 'agent-001', got '%s'", all["key1"].AgentID)
	}
}

// TestSharedContextHistory tests history tracking
func TestSharedContextHistory(t *testing.T) {
	sc := NewSharedContext("team-001")

	// Perform operations
	sc.Set("key1", "value1", "agent-001")
	sc.Set("key2", "value2", "agent-002")
	sc.AddHistoryEntry("agent-003", "custom_action", map[string]string{"detail": "test"})

	history := sc.GetHistory()

	if len(history) != 3 {
		t.Errorf("Expected 3 history entries, got %d", len(history))
	}

	// Check first entry
	if history[0].AgentID != "agent-001" {
		t.Errorf("Expected first entry agentID 'agent-001', got '%s'", history[0].AgentID)
	}

	if history[0].Action != "set" {
		t.Errorf("Expected first entry action 'set', got '%s'", history[0].Action)
	}

	// Check custom entry
	if history[2].Action != "custom_action" {
		t.Errorf("Expected third entry action 'custom_action', got '%s'", history[2].Action)
	}
}

// TestSharedContextSnapshot tests snapshot creation
func TestSharedContextSnapshot(t *testing.T) {
	sc := NewSharedContext("team-001")

	// Add entries
	sc.Set("key1", "value1", "agent-001")
	sc.Set("key2", 42, "agent-002")
	sc.Set("key3", true, "agent-003")

	snapshot := sc.Snapshot()

	if len(snapshot) != 3 {
		t.Errorf("Expected 3 entries in snapshot, got %d", len(snapshot))
	}

	if snapshot["key1"] != "value1" {
		t.Errorf("Expected snapshot key1 'value1', got '%v'", snapshot["key1"])
	}

	if snapshot["key2"] != 42 {
		t.Errorf("Expected snapshot key2 42, got '%v'", snapshot["key2"])
	}

	if snapshot["key3"] != true {
		t.Errorf("Expected snapshot key3 true, got '%v'", snapshot["key3"])
	}
}

// TestSharedContextConcurrentReads tests concurrent read operations
func TestSharedContextConcurrentReads(t *testing.T) {
	sc := NewSharedContext("team-001")

	// Populate context
	for i := 0; i < 10; i++ {
		sc.Set(string(rune('a'+i)), i, "agent-001")
	}

	// Concurrent reads
	var wg sync.WaitGroup
	numReaders := 100

	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				_, _ = sc.Get(string(rune('a' + j)))
			}
		}()
	}

	wg.Wait()
	// If we get here without deadlock, test passes
}

// TestSharedContextConcurrentWrites tests concurrent write operations
func TestSharedContextConcurrentWrites(t *testing.T) {
	sc := NewSharedContext("team-001")

	var wg sync.WaitGroup
	numWriters := 50

	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			sc.Set(string(rune('a'+id%10)), id, "agent-001")
		}(i)
	}

	wg.Wait()

	// Verify all writes completed
	all := sc.GetAll()
	if len(all) != 10 {
		t.Errorf("Expected 10 unique keys, got %d", len(all))
	}
}

// TestSharedContextTimestamps tests that timestamps are set correctly
func TestSharedContextTimestamps(t *testing.T) {
	sc := NewSharedContext("team-001")

	before := time.Now()
	time.Sleep(10 * time.Millisecond)

	sc.Set("key1", "value1", "agent-001")

	time.Sleep(10 * time.Millisecond)
	after := time.Now()

	all := sc.GetAll()
	timestamp := all["key1"].Timestamp

	if timestamp.Before(before) || timestamp.After(after) {
		t.Errorf("Timestamp %v not between %v and %v", timestamp, before, after)
	}
}

// TestSharedContextHistoryOrdering tests chronological order of history
func TestSharedContextHistoryOrdering(t *testing.T) {
	sc := NewSharedContext("team-001")

	// Add entries with small delays
	sc.Set("key1", "value1", "agent-001")
	time.Sleep(5 * time.Millisecond)

	sc.Set("key2", "value2", "agent-002")
	time.Sleep(5 * time.Millisecond)

	sc.Set("key3", "value3", "agent-003")

	history := sc.GetHistory()

	// Verify chronological order
	for i := 1; i < len(history); i++ {
		if history[i].Timestamp.Before(history[i-1].Timestamp) {
			t.Errorf("History entry %d timestamp before entry %d", i, i-1)
		}
	}
}

// Mock session manager for testing
type mockSessionManager struct {
	sessions map[string][]mockMessage
	mu       sync.RWMutex
}

type mockMessage struct {
	role    string
	content string
}

func (m *mockSessionManager) GetOrCreate(key string) any {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.sessions[key]; !exists {
		m.sessions[key] = []mockMessage{}
	}
	return nil
}

func (m *mockSessionManager) AddMessage(sessionKey, role, content string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.sessions[sessionKey]; !exists {
		m.sessions[sessionKey] = []mockMessage{}
	}
	m.sessions[sessionKey] = append(m.sessions[sessionKey], mockMessage{role: role, content: content})
}

func (m *mockSessionManager) Save(key string) error {
	return nil
}

func (m *mockSessionManager) GetHistory(key string) []any {
	m.mu.RLock()
	defer m.mu.RUnlock()
	msgs, exists := m.sessions[key]
	if !exists {
		return []any{}
	}
	result := make([]any, len(msgs))
	for i, msg := range msgs {
		result[i] = msg
	}
	return result
}

func (m mockMessage) GetContent() string {
	return m.content
}

func TestSharedContext_PersistToSession(t *testing.T) {
	ctx := NewSharedContext("test-team")
	ctx.Set("key1", "value1", "agent1")
	ctx.Set("key2", 123, "agent1")

	sm := &mockSessionManager{
		sessions: make(map[string][]mockMessage),
	}

	err := ctx.PersistToSession(sm, "team1")
	if err != nil {
		t.Fatalf("PersistToSession failed: %v", err)
	}

	// Check session was created
	sessionKey := "team:team1:context"
	msgs, exists := sm.sessions[sessionKey]
	if !exists {
		t.Error("Expected session to be created")
	}

	if len(msgs) != 1 {
		t.Errorf("Expected 1 message, got %d", len(msgs))
	}

	if msgs[0].role != "system" {
		t.Errorf("Expected role 'system', got '%s'", msgs[0].role)
	}
}

func TestSharedContext_LoadFromSession(t *testing.T) {
	// Create context and persist
	ctx1 := NewSharedContext("test-team")
	ctx1.Set("key1", "value1", "agent1")
	ctx1.Set("key2", float64(123), "agent1") // JSON unmarshals numbers as float64

	sm := &mockSessionManager{
		sessions: make(map[string][]mockMessage),
	}

	err := ctx1.PersistToSession(sm, "team1")
	if err != nil {
		t.Fatalf("PersistToSession failed: %v", err)
	}

	// Create new context and load
	ctx2 := NewSharedContext("test-team")
	err = ctx2.LoadFromSession(sm, "team1")
	if err != nil {
		t.Fatalf("LoadFromSession failed: %v", err)
	}

	// Check data was loaded
	val1, exists := ctx2.Get("key1")
	if !exists {
		t.Error("Expected key1 to exist")
	}
	if val1 != "value1" {
		t.Errorf("Expected value1, got %v", val1)
	}

	val2, exists := ctx2.Get("key2")
	if !exists {
		t.Error("Expected key2 to exist")
	}
	if val2 != float64(123) {
		t.Errorf("Expected 123, got %v", val2)
	}
}

func TestSharedContext_LoadFromSession_NoData(t *testing.T) {
	ctx := NewSharedContext("test-team")
	sm := &mockSessionManager{
		sessions: make(map[string][]mockMessage),
	}

	err := ctx.LoadFromSession(sm, "team1")
	if err != nil {
		t.Errorf("Expected no error for empty session, got %v", err)
	}

	// Context should remain empty
	allData := ctx.GetAll()
	if len(allData) != 0 {
		t.Errorf("Expected empty context, got %d items", len(allData))
	}
}
