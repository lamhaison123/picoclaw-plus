// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package team

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// SharedContext provides thread-safe storage for team-wide information
type SharedContext struct {
	teamID  string
	entries map[string]*ContextEntry
	history []*HistoryEntry
	mu      sync.RWMutex
}

// ContextEntry represents a single entry in the shared context
type ContextEntry struct {
	Key       string
	Value     any
	Timestamp time.Time
	AgentID   string
}

// HistoryEntry represents a historical action in the shared context
type HistoryEntry struct {
	Timestamp time.Time
	AgentID   string
	Action    string
	Data      any
}

// NewSharedContext creates a new SharedContext for a team
func NewSharedContext(teamID string) *SharedContext {
	return &SharedContext{
		teamID:  teamID,
		entries: make(map[string]*ContextEntry),
		history: make([]*HistoryEntry, 0),
	}
}

// Set stores a value in the shared context
func (sc *SharedContext) Set(key string, value any, agentID string) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	entry := &ContextEntry{
		Key:       key,
		Value:     value,
		Timestamp: time.Now(),
		AgentID:   agentID,
	}

	sc.entries[key] = entry

	// Add to history with size limit (atomic operation)
	const maxHistorySize = 1000
	historyEntry := &HistoryEntry{
		Timestamp: time.Now(),
		AgentID:   agentID,
		Action:    "set",
		Data:      map[string]any{"key": key, "value": value},
	}

	if len(sc.history) >= maxHistorySize {
		// Keep last 80% of entries to avoid frequent trimming
		keepSize := maxHistorySize * 4 / 5
		sc.history = sc.history[len(sc.history)-keepSize:]
	}
	sc.history = append(sc.history, historyEntry)

	return nil
}

// Get retrieves a value from the shared context
func (sc *SharedContext) Get(key string) (any, bool) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	entry, exists := sc.entries[key]
	if !exists {
		return nil, false
	}

	return entry.Value, true
}

// GetAll returns all entries in the shared context
func (sc *SharedContext) GetAll() map[string]*ContextEntry {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	// Create a copy to avoid external modifications
	result := make(map[string]*ContextEntry, len(sc.entries))
	for k, v := range sc.entries {
		result[k] = v
	}

	return result
}

// GetHistory returns the chronological history of actions
func (sc *SharedContext) GetHistory() []*HistoryEntry {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	// Create a copy to avoid external modifications
	result := make([]*HistoryEntry, len(sc.history))
	copy(result, sc.history)

	return result
}

// AddHistoryEntry adds a custom entry to the history
func (sc *SharedContext) AddHistoryEntry(agentID, action string, data any) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	historyEntry := &HistoryEntry{
		Timestamp: time.Now(),
		AgentID:   agentID,
		Action:    action,
		Data:      data,
	}

	// Apply same size limit as Set()
	const maxHistorySize = 1000
	if len(sc.history) >= maxHistorySize {
		keepSize := maxHistorySize * 4 / 5
		sc.history = sc.history[len(sc.history)-keepSize:]
	}

	sc.history = append(sc.history, historyEntry)
}

// Snapshot creates a snapshot of the current context for persistence
func (sc *SharedContext) Snapshot() map[string]any {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	snapshot := make(map[string]any, len(sc.entries))
	for k, v := range sc.entries {
		snapshot[k] = v.Value
	}

	return snapshot
}

// PersistToSession saves shared context to session storage
func (sc *SharedContext) PersistToSession(sessionManager any, teamID string) error {
	// Type assert to session manager interface
	type sessionPersister interface {
		GetOrCreate(key string) any
		AddMessage(sessionKey, role, content string)
		Save(key string) error
	}

	sm, ok := sessionManager.(sessionPersister)
	if !ok {
		return fmt.Errorf("invalid session manager type")
	}

	// Create session key
	sessionKey := fmt.Sprintf("team:%s:context", teamID)

	// Get snapshot of context
	snapshot := sc.Snapshot()

	// Convert to JSON
	data, err := json.Marshal(snapshot)
	if err != nil {
		return fmt.Errorf("failed to marshal context: %w", err)
	}

	// Store in session
	sm.AddMessage(sessionKey, "system", string(data))

	// Save session
	if err := sm.Save(sessionKey); err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}

	return nil
}

// LoadFromSession loads shared context from session storage
func (sc *SharedContext) LoadFromSession(sessionManager any, teamID string) error {
	// Type assert to session manager interface
	type sessionLoader interface {
		GetHistory(key string) []any
	}

	sl, ok := sessionManager.(sessionLoader)
	if !ok {
		return fmt.Errorf("invalid session manager type")
	}

	// Create session key
	sessionKey := fmt.Sprintf("team:%s:context", teamID)

	// Get history from session
	history := sl.GetHistory(sessionKey)

	if len(history) == 0 {
		return nil // No saved context
	}

	// Get last message (most recent context)
	lastMsg := history[len(history)-1]

	// Type assert to get content
	type message interface {
		GetContent() string
	}

	msg, ok := lastMsg.(message)
	if !ok {
		return fmt.Errorf("invalid message type")
	}

	content := msg.GetContent()

	// Parse JSON
	var snapshot map[string]any
	if err := json.Unmarshal([]byte(content), &snapshot); err != nil {
		return fmt.Errorf("failed to unmarshal context: %w", err)
	}

	// Restore context
	sc.mu.Lock()
	defer sc.mu.Unlock()

	for key, value := range snapshot {
		sc.entries[key] = &ContextEntry{
			Key:       key,
			Value:     value,
			Timestamp: time.Now(),
			AgentID:   "session",
		}
	}

	return nil
}
