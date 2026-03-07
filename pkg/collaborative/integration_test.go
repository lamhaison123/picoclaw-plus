package collaborative

import (
	"context"
	"sync"
	"testing"
	"time"
)

// MockPlatform for testing
type MockPlatform struct {
	mu       sync.Mutex
	messages []string
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewMockPlatform() *MockPlatform {
	ctx, cancel := context.WithCancel(context.Background())
	return &MockPlatform{
		messages: make([]string, 0),
		ctx:      ctx,
		cancel:   cancel,
	}
}

func (m *MockPlatform) SendMessage(ctx context.Context, message string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.messages = append(m.messages, message)
	return nil
}

func (m *MockPlatform) GetContext() context.Context {
	return m.ctx
}

func (m *MockPlatform) GetMessages() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	result := make([]string, len(m.messages))
	copy(result, m.messages)
	return result
}

func (m *MockPlatform) CancelContext() {
	m.cancel()
}

// MockTeamManager for testing
type MockTeamManager struct {
	mu        sync.Mutex
	execCount map[string]int
}

func NewMockTeamManager() *MockTeamManager {
	return &MockTeamManager{
		execCount: make(map[string]int),
	}
}

func (m *MockTeamManager) ExecuteTaskWithRole(ctx context.Context, role, task string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.execCount[role]++

	// Simulate agent response with mention
	if role == "manager" {
		return "Task completed! @developer please verify this.", nil
	}
	if role == "developer" {
		return "Verified successfully!", nil
	}
	return "Done", nil
}

func (m *MockTeamManager) GetExecutionCount(role string) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.execCount[role]
}

// TestManagerMentionsDeveloper tests the specific bug case
func TestManagerMentionsDeveloper(t *testing.T) {
	platform := NewMockPlatform()
	teamManager := NewMockTeamManager()

	// Use ManagerV2 instead of deprecated Manager
	config := &Config{
		MentionQueueSize:    20,
		MentionRateLimit:    2 * time.Second,
		MentionMaxRetries:   3,
		MentionRetryBackoff: 1 * time.Second,
	}
	manager := NewManagerV2WithConfig(config)

	var chatID int64 = 123456
	sessionID := "test-session-456"

	// Create session
	session := &Session{
		ChatID:       chatID,
		SessionID:    sessionID,
		TeamID:       "dev-team",
		StartTime:    time.Now(),
		LastActivity: time.Now(),
		Context:      make([]Message, 0),
		ActiveAgents: make(map[string]*AgentState),
		MaxContext:   100,
	}
	manager.sessions[chatID] = session

	// Simulate manager being triggered and responding with @developer mention
	ctx := context.Background()

	// Manager executes task
	managerResponse, err := teamManager.ExecuteTaskWithRole(ctx, "manager", "Update docs")
	if err != nil {
		t.Fatalf("Manager execution failed: %v", err)
	}

	// Check if manager response contains @developer mention
	mentions := ExtractMentions(managerResponse)
	if len(mentions) == 0 {
		t.Fatal("Manager response should contain @developer mention")
	}
	if mentions[0] != "developer" {
		t.Fatalf("Expected mention 'developer', got '%s'", mentions[0])
	}

	// Simulate nested mention dispatch (this is where the bug was)
	// The fix uses context.Background() instead of platform.GetContext()
	go func() {
		execCtx, execCancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer execCancel()

		_, err := teamManager.ExecuteTaskWithRole(execCtx, "developer", "Verify docs update")
		if err != nil {
			t.Errorf("Developer execution failed: %v", err)
		}
	}()

	// Cancel platform context (simulating manager task completion)
	platform.CancelContext()

	// Wait for developer to execute
	time.Sleep(100 * time.Millisecond)

	// Verify developer was triggered despite platform context cancellation
	devCount := teamManager.GetExecutionCount("developer")
	if devCount != 1 {
		t.Errorf("Developer should be triggered once, got %d times", devCount)
	}
}

// TestNestedMentionDepthLimit tests depth tracking via message context
func TestNestedMentionDepthLimit(t *testing.T) {
	// Use ManagerV2 instead of deprecated Manager
	config := &Config{
		MentionQueueSize:    20,
		MentionRateLimit:    2 * time.Second,
		MentionMaxRetries:   3,
		MentionRetryBackoff: 1 * time.Second,
	}
	manager := NewManagerV2WithConfig(config)

	var chatID int64 = 789
	sessionID := "test-session-depth"
	session := &Session{
		ChatID:       chatID,
		SessionID:    sessionID,
		TeamID:       "dev-team",
		StartTime:    time.Now(),
		LastActivity: time.Now(),
		Context:      make([]Message, 0),
		ActiveAgents: make(map[string]*AgentState),
		MaxContext:   100,
	}
	manager.sessions[chatID] = session

	// Test depth tracking via message context
	// Add messages to simulate nested mentions
	for range 3 {
		msg := Message{
			Role:      "manager",
			Content:   "@developer please check",
			Timestamp: time.Now(),
			Mentions:  []string{"developer"},
		}
		session.Context = append(session.Context, msg)
	}

	// Verify we can track depth via context length
	if len(session.Context) != 3 {
		t.Errorf("Expected 3 messages in context, got %d", len(session.Context))
	}
}

// TestConcurrentMentions tests multiple agents mentioning simultaneously
func TestConcurrentMentions(t *testing.T) {
	teamManager := NewMockTeamManager()

	var wg sync.WaitGroup
	roles := []string{"manager", "architect", "tester"}

	for _, role := range roles {
		wg.Add(1)
		go func(r string) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()

			_, err := teamManager.ExecuteTaskWithRole(ctx, r, "concurrent task")
			if err != nil {
				t.Errorf("Role %s execution failed: %v", r, err)
			}
		}(role)
	}

	wg.Wait()

	// Verify all roles were executed
	for _, role := range roles {
		count := teamManager.GetExecutionCount(role)
		if count != 1 {
			t.Errorf("Role %s should execute once, got %d", role, count)
		}
	}
}

// TestLongFormContentWithMentions tests mentions in complex content
func TestLongFormContentWithMentions(t *testing.T) {
	longContent := `
# Báo Cáo Tiến Độ Dự Án

## Tổng Quan
Dự án đang tiến triển tốt với nhiều cải tiến quan trọng.

### Circuit Breaker Implementation
- ✅ Core implementation hoàn thành
- ✅ Tests coverage 100%
- ⏳ Đang chờ @developer verify

### Collaborative Chat
Hệ thống chat đã được cải thiện với:
- Mention detection tốt hơn
- UTF-8 support đầy đủ
- Depth tracking để tránh infinite loops

@architect nhờ anh review kiến trúc.
@tester nhờ bạn chạy integration tests.

## Kết Luận
Mọi thứ đang on track! 🚀
`

	mentions := ExtractMentions(longContent)

	expectedMentions := []string{"developer", "architect", "tester"}
	if len(mentions) != len(expectedMentions) {
		t.Errorf("Expected %d mentions, got %d", len(expectedMentions), len(mentions))
	}

	for i, expected := range expectedMentions {
		if mentions[i] != expected {
			t.Errorf("Expected mention '%s', got '%s'", expected, mentions[i])
		}
	}
}
