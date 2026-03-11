package session

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/sipeed/picoclaw/pkg/logger"
	"github.com/sipeed/picoclaw/pkg/memory"
	"github.com/sipeed/picoclaw/pkg/providers"
)

type Session struct {
	Key      string              `json:"key"`
	Messages []providers.Message `json:"messages"`
	Summary  string              `json:"summary,omitempty"`
	Created  time.Time           `json:"created"`
	Updated  time.Time           `json:"updated"`
}

// SessionManager manages conversation sessions with pluggable storage backend.
// v0.2.1: Refactored to use Store interface for JSONL support
type SessionManager struct {
	sessions map[string]*Session
	mu       sync.RWMutex
	storage  string
	store    memory.Store // v0.2.1: Pluggable storage backend
	useStore bool         // Whether to use Store interface or legacy in-memory
}

// NewSessionManager creates a new session manager with legacy in-memory storage.
// Deprecated: Use NewSessionManagerWithStore for production deployments.
func NewSessionManager(storage string) *SessionManager {
	sm := &SessionManager{
		sessions: make(map[string]*Session),
		storage:  storage,
		useStore: false,
	}

	if storage != "" {
		os.MkdirAll(storage, 0o755)
		sm.loadSessions()
	}

	return sm
}

// NewSessionManagerWithStore creates a new session manager with pluggable storage backend.
// v0.2.1: Use Store interface for crash-safe JSONL storage
func NewSessionManagerWithStore(storage string, store memory.Store) *SessionManager {
	sm := &SessionManager{
		sessions: make(map[string]*Session),
		storage:  storage,
		store:    store,
		useStore: true,
	}

	if storage != "" {
		os.MkdirAll(storage, 0o755)
	}

	// Auto-migrate from JSON to JSONL if needed
	if store != nil {
		ctx := context.Background()
		if err := memory.AutoMigrate(ctx, storage, store); err != nil {
			logger.WarnCF("session", "Auto-migration failed",
				map[string]any{"error": err.Error()})
		}
	}

	return sm
}

func (sm *SessionManager) GetOrCreate(key string) *Session {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	session, ok := sm.sessions[key]
	if ok {
		return session
	}

	session = &Session{
		Key:      key,
		Messages: []providers.Message{},
		Created:  time.Now(),
		Updated:  time.Now(),
	}
	sm.sessions[key] = session

	return session
}

func (sm *SessionManager) AddMessage(sessionKey, role, content string) {
	// v0.2.1: Use Store interface if available
	if sm.useStore && sm.store != nil {
		ctx := context.Background()
		if err := sm.store.AddMessage(ctx, sessionKey, role, content); err != nil {
			logger.ErrorCF("session", "Failed to add message to store",
				map[string]any{
					"key":   sessionKey,
					"error": err.Error(),
				})
		}
		return
	}

	// Legacy in-memory storage
	sm.AddFullMessage(sessionKey, providers.Message{
		Role:    role,
		Content: content,
	})
}

// AddFullMessage adds a complete message with tool calls and tool call ID to the session.
// This is used to save the full conversation flow including tool calls and tool results.
func (sm *SessionManager) AddFullMessage(sessionKey string, msg providers.Message) {
	// v0.2.1: Use Store interface if available
	if sm.useStore && sm.store != nil {
		ctx := context.Background()
		if err := sm.store.AddFullMessage(ctx, sessionKey, msg); err != nil {
			logger.ErrorCF("session", "Failed to add full message to store",
				map[string]any{
					"key":   sessionKey,
					"error": err.Error(),
				})
		}
		return
	}

	// Legacy in-memory storage
	sm.mu.Lock()
	defer sm.mu.Unlock()

	session, ok := sm.sessions[sessionKey]
	if !ok {
		session = &Session{
			Key:      sessionKey,
			Messages: []providers.Message{},
			Created:  time.Now(),
		}
		sm.sessions[sessionKey] = session
	}

	session.Messages = append(session.Messages, msg)
	session.Updated = time.Now()
}

func (sm *SessionManager) GetHistory(key string) []providers.Message {
	// v0.2.1: Use Store interface if available
	if sm.useStore && sm.store != nil {
		ctx := context.Background()
		history, err := sm.store.GetHistory(ctx, key)
		if err != nil {
			logger.ErrorCF("session", "Failed to get history from store",
				map[string]any{
					"key":   key,
					"error": err.Error(),
				})
			return []providers.Message{}
		}
		return history
	}

	// Legacy in-memory storage
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	session, ok := sm.sessions[key]
	if !ok {
		return []providers.Message{}
	}

	history := make([]providers.Message, len(session.Messages))
	copy(history, session.Messages)
	return history
}

func (sm *SessionManager) GetSummary(key string) string {
	// v0.2.1: Use Store interface if available
	if sm.useStore && sm.store != nil {
		ctx := context.Background()
		summary, err := sm.store.GetSummary(ctx, key)
		if err != nil {
			logger.ErrorCF("session", "Failed to get summary from store",
				map[string]any{
					"key":   key,
					"error": err.Error(),
				})
			return ""
		}
		return summary
	}

	// Legacy in-memory storage
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	session, ok := sm.sessions[key]
	if !ok {
		return ""
	}
	return session.Summary
}

func (sm *SessionManager) SetSummary(key string, summary string) {
	// v0.2.1: Use Store interface if available
	if sm.useStore && sm.store != nil {
		ctx := context.Background()
		if err := sm.store.SetSummary(ctx, key, summary); err != nil {
			logger.ErrorCF("session", "Failed to set summary in store",
				map[string]any{
					"key":   key,
					"error": err.Error(),
				})
		}
		return
	}

	// Legacy in-memory storage
	sm.mu.Lock()
	defer sm.mu.Unlock()

	session, ok := sm.sessions[key]
	if ok {
		session.Summary = summary
		session.Updated = time.Now()
	}
}

func (sm *SessionManager) TruncateHistory(key string, keepLast int) {
	// v0.2.1: Use Store interface if available
	if sm.useStore && sm.store != nil {
		ctx := context.Background()
		if err := sm.store.TruncateHistory(ctx, key, keepLast); err != nil {
			logger.ErrorCF("session", "Failed to truncate history in store",
				map[string]any{
					"key":   key,
					"error": err.Error(),
				})
		}
		return
	}

	// Legacy in-memory storage
	sm.mu.Lock()
	defer sm.mu.Unlock()

	session, ok := sm.sessions[key]
	if !ok {
		return
	}

	if keepLast <= 0 {
		session.Messages = []providers.Message{}
		session.Updated = time.Now()
		return
	}

	if len(session.Messages) <= keepLast {
		return
	}

	session.Messages = session.Messages[len(session.Messages)-keepLast:]
	session.Updated = time.Now()
}

// sanitizeFilename converts a session key into a cross-platform safe filename.
// Session keys use "channel:chatID" (e.g. "telegram:123456") but ':' is the
// volume separator on Windows, so filepath.Base would misinterpret the key.
// We replace it with '_'. The original key is preserved inside the JSON file,
// so loadSessions still maps back to the right in-memory key.
func sanitizeFilename(key string) string {
	return strings.ReplaceAll(key, ":", "_")
}

func (sm *SessionManager) loadSessions() error {
	files, err := os.ReadDir(sm.storage)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if filepath.Ext(file.Name()) != ".json" {
			continue
		}

		sessionPath := filepath.Join(sm.storage, file.Name())
		data, err := os.ReadFile(sessionPath)
		if err != nil {
			continue
		}

		var session Session
		if err := json.Unmarshal(data, &session); err != nil {
			continue
		}

		sm.sessions[session.Key] = &session
	}

	return nil
}

// SetHistory updates the messages of a session.
func (sm *SessionManager) SetHistory(key string, history []providers.Message) {
	// v0.2.1: Use Store interface if available
	if sm.useStore && sm.store != nil {
		ctx := context.Background()
		if err := sm.store.SetHistory(ctx, key, history); err != nil {
			logger.ErrorCF("session", "Failed to set history in store",
				map[string]any{
					"key":   key,
					"error": err.Error(),
				})
		}
		return
	}

	// Legacy in-memory storage
	sm.mu.Lock()
	defer sm.mu.Unlock()

	session, ok := sm.sessions[key]
	if ok {
		// Create a deep copy to strictly isolate internal state
		// from the caller's slice, including nested slices
		msgs := make([]providers.Message, len(history))
		for i, msg := range history {
			msgs[i] = deepCopyMessage(msg)
		}
		session.Messages = msgs
		session.Updated = time.Now()
	}
}

// Save persists a session to disk (legacy JSON format only).
// v0.2.1: When using Store interface, saves are automatic (no-op).
func (sm *SessionManager) Save(key string) error {
	// v0.2.1: Store interface handles persistence automatically
	if sm.useStore && sm.store != nil {
		return nil
	}

	// Legacy JSON storage
	if sm.storage == "" {
		return nil
	}

	filename := sanitizeFilename(key)

	// filepath.IsLocal rejects empty names, "..", absolute paths, and
	// OS-reserved device names (NUL, COM1 … on Windows).
	// The extra checks reject "." and any directory separators so that
	// the session file is always written directly inside sm.storage.
	if filename == "." || !filepath.IsLocal(filename) || strings.ContainsAny(filename, `/\`) {
		return os.ErrInvalid
	}

	// Snapshot under read lock, then perform slow file I/O after unlock.
	sm.mu.RLock()
	stored, ok := sm.sessions[key]
	if !ok {
		sm.mu.RUnlock()
		return nil
	}

	snapshot := Session{
		Key:     stored.Key,
		Summary: stored.Summary,
		Created: stored.Created,
		Updated: stored.Updated,
	}
	if len(stored.Messages) > 0 {
		snapshot.Messages = make([]providers.Message, len(stored.Messages))
		copy(snapshot.Messages, stored.Messages)
	} else {
		snapshot.Messages = []providers.Message{}
	}
	sm.mu.RUnlock()

	data, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return err
	}

	sessionPath := filepath.Join(sm.storage, filename+".json")
	tmpFile, err := os.CreateTemp(sm.storage, "session-*.tmp")
	if err != nil {
		return err
	}

	tmpPath := tmpFile.Name()
	cleanup := true
	defer func() {
		if cleanup {
			_ = os.Remove(tmpPath)
		}
	}()

	if _, err := tmpFile.Write(data); err != nil {
		_ = tmpFile.Close()
		return err
	}
	if err := tmpFile.Chmod(0o644); err != nil {
		_ = tmpFile.Close()
		return err
	}
	if err := tmpFile.Sync(); err != nil {
		_ = tmpFile.Close()
		return err
	}
	if err := tmpFile.Close(); err != nil {
		return err
	}

	if err := os.Rename(tmpPath, sessionPath); err != nil {
		return err
	}
	cleanup = false
	return nil
}

// deepCopyMessage creates a deep copy of a Message including all nested slices
func deepCopyMessage(msg providers.Message) providers.Message {
	copied := providers.Message{
		Role:             msg.Role,
		Content:          msg.Content,
		ReasoningContent: msg.ReasoningContent,
	}

	// Deep copy ToolCalls
	if len(msg.ToolCalls) > 0 {
		copied.ToolCalls = make([]providers.ToolCall, len(msg.ToolCalls))
		for i, tc := range msg.ToolCalls {
			copied.ToolCalls[i] = providers.ToolCall{
				ID:   tc.ID,
				Name: tc.Name,
				Type: tc.Type,
			}
			// Deep copy Arguments map
			if tc.Arguments != nil {
				copied.ToolCalls[i].Arguments = make(map[string]interface{})
				for k, v := range tc.Arguments {
					copied.ToolCalls[i].Arguments[k] = v
				}
			}
			// Copy ExtraContent (shallow copy is sufficient for struct)
			if tc.ExtraContent != nil {
				extraCopy := *tc.ExtraContent
				copied.ToolCalls[i].ExtraContent = &extraCopy
			}
			// Copy Function if present
			if tc.Function != nil {
				copied.ToolCalls[i].Function = &providers.FunctionCall{
					Name:             tc.Function.Name,
					Arguments:        tc.Function.Arguments,
					ThoughtSignature: tc.Function.ThoughtSignature,
				}
			}
		}
	}

	// Deep copy Media slice
	if len(msg.Media) > 0 {
		copied.Media = make([]string, len(msg.Media))
		copy(copied.Media, msg.Media)
	}

	return copied
}
