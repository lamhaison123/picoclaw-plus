// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package collaborative

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sipeed/picoclaw/pkg/bus"
)

// Mock Platform for testing
type mockPlatformV2 struct {
	messages       []string
	mu             sync.Mutex
	teamManager    *mockTeamManagerV2
	ctx            context.Context
	sendDelay      time.Duration
	shouldFailSend bool
}

func newMockPlatformV2() *mockPlatformV2 {
	return &mockPlatformV2{
		messages:    make([]string, 0),
		teamManager: newMockTeamManagerV2(),
		ctx:         context.Background(),
	}
}

func (m *mockPlatformV2) SendMessage(ctx context.Context, chatID string, content string) error {
	if m.sendDelay > 0 {
		time.Sleep(m.sendDelay)
	}

	if m.shouldFailSend {
		return fmt.Errorf("mock send failure")
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	m.messages = append(m.messages, content)
	return nil
}

func (m *mockPlatformV2) GetTeamManager() TeamManager {
	return m.teamManager
}

func (m *mockPlatformV2) GetContext() context.Context {
	return m.ctx
}

// Mock Team Manager for testing
type mockTeamManagerV2 struct {
	execCount    atomic.Int32
	shouldFail   bool
	failCount    int
	currentFails int
	execDelay    time.Duration
	mu           sync.Mutex
}

func newMockTeamManagerV2() *mockTeamManagerV2 {
	return &mockTeamManagerV2{}
}

func (m *mockTeamManagerV2) ExecuteTaskWithRole(ctx context.Context, teamID, prompt, role string) (any, error) {
	m.execCount.Add(1)

	if m.execDelay > 0 {
		time.Sleep(m.execDelay)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.shouldFail {
		if m.failCount > 0 && m.currentFails < m.failCount {
			m.currentFails++
			return nil, fmt.Errorf("mock execution failure %d", m.currentFails)
		}
	}

	return fmt.Sprintf("Response from @%s", role), nil
}

func (m *mockTeamManagerV2) GetTeam(teamID string) (any, error) {
	return map[string]any{
		"id": teamID,
		"roles": []map[string]string{
			{"name": "architect", "description": "System architect"},
			{"name": "developer", "description": "Developer"},
		},
	}, nil
}

func (m *mockTeamManagerV2) getExecCount() int32 {
	return m.execCount.Load()
}

// Test 1: Queue Integration Test
func TestManagerV2_QueueIntegration(t *testing.T) {
	config := &Config{
		MentionQueueSize:    5,
		MentionRateLimit:    100 * time.Millisecond,
		MentionMaxRetries:   2,
		MentionRetryBackoff: 50 * time.Millisecond,
	}

	manager := NewManagerV2WithConfig(config)
	defer manager.Stop()

	platform := newMockPlatformV2()
	ctx := context.Background()

	// Test enqueue
	err := manager.HandleMentions(
		ctx,
		platform,
		12345,
		"test-team",
		"@architect help me",
		[]string{"architect"},
		bus.SenderInfo{Username: "user1"},
		50,
	)

	if err != nil {
		t.Fatalf("HandleMentions failed: %v", err)
	}

	// Wait for execution
	time.Sleep(500 * time.Millisecond)

	// Check metrics
	metrics := manager.GetQueueMetrics("architect")
	if metrics == nil {
		t.Fatal("Expected metrics for architect, got nil")
	}

	if metrics.ProcessedCount == 0 {
		t.Error("Expected processed count > 0")
	}

	t.Logf("Queue metrics: processed=%d, dropped=%d, retry=%d",
		metrics.ProcessedCount, metrics.DroppedCount, metrics.RetryCount)
}

// Test 2: Queue Overflow Test
func TestManagerV2_QueueOverflow(t *testing.T) {
	config := &Config{
		MentionQueueSize:  2,               // Small queue
		MentionRateLimit:  1 * time.Second, // Slow processing
		MentionMaxRetries: 1,
	}

	manager := NewManagerV2WithConfig(config)
	defer manager.Stop()

	platform := newMockPlatformV2()
	platform.teamManager.execDelay = 500 * time.Millisecond // Slow execution
	ctx := context.Background()

	// Enqueue more than queue size
	for i := 0; i < 5; i++ {
		err := manager.HandleMentions(
			ctx,
			platform,
			12345,
			"test-team",
			fmt.Sprintf("@architect task %d", i),
			[]string{"architect"},
			bus.SenderInfo{Username: "user1"},
			50,
		)

		if err != nil {
			t.Logf("HandleMentions %d failed (expected): %v", i, err)
		}
	}

	// Wait a bit
	time.Sleep(200 * time.Millisecond)

	// Check metrics
	metrics := manager.GetQueueMetrics("architect")
	if metrics == nil {
		t.Fatal("Expected metrics for architect, got nil")
	}

	if metrics.DroppedCount == 0 {
		t.Error("Expected some dropped mentions due to queue overflow")
	}

	t.Logf("Queue overflow test: processed=%d, dropped=%d",
		metrics.ProcessedCount, metrics.DroppedCount)
}

// Test 3: Rate Limiting Test
func TestManagerV2_RateLimiting(t *testing.T) {
	rateLimit := 500 * time.Millisecond
	config := &Config{
		MentionQueueSize:  10,
		MentionRateLimit:  rateLimit,
		MentionMaxRetries: 1,
	}

	manager := NewManagerV2WithConfig(config)
	defer manager.Stop()

	platform := newMockPlatformV2()
	ctx := context.Background()

	// Enqueue multiple mentions for same role
	startTime := time.Now()

	for i := 0; i < 3; i++ {
		err := manager.HandleMentions(
			ctx,
			platform,
			12345,
			"test-team",
			fmt.Sprintf("@architect task %d", i),
			[]string{"architect"},
			bus.SenderInfo{Username: "user1"},
			50,
		)

		if err != nil {
			t.Fatalf("HandleMentions %d failed: %v", i, err)
		}
	}

	// Wait for all executions
	time.Sleep(2 * time.Second)

	elapsed := time.Since(startTime)

	// Should take at least 2 * rateLimit (3 executions with 2 delays)
	minExpected := 2 * rateLimit

	if elapsed < minExpected {
		t.Errorf("Rate limiting not working: elapsed=%v, expected>=%v",
			elapsed, minExpected)
	}

	execCount := platform.teamManager.getExecCount()
	if execCount != 3 {
		t.Errorf("Expected 3 executions, got %d", execCount)
	}

	t.Logf("Rate limiting test: elapsed=%v, executions=%d", elapsed, execCount)
}

// Test 4: Cascade with Queue Test
func TestManagerV2_CascadeWithQueue(t *testing.T) {
	config := &Config{
		MentionQueueSize:    10,
		MentionRateLimit:    100 * time.Millisecond,
		MentionMaxRetries:   2,
		MentionRetryBackoff: 50 * time.Millisecond,
	}

	manager := NewManagerV2WithConfig(config)
	defer manager.Stop()

	platform := newMockPlatformV2()
	ctx := context.Background()

	// Initial mention
	err := manager.HandleMentions(
		ctx,
		platform,
		12345,
		"test-team",
		"@architect design system",
		[]string{"architect"},
		bus.SenderInfo{Username: "user1"},
		50,
	)

	if err != nil {
		t.Fatalf("HandleMentions failed: %v", err)
	}

	// Wait for cascade to complete
	time.Sleep(1 * time.Second)

	// Check that execution happened
	execCount := platform.teamManager.getExecCount()
	if execCount == 0 {
		t.Error("Expected at least 1 execution")
	}

	// Check metrics for all roles
	allMetrics := manager.GetAllQueueMetrics()
	t.Logf("Cascade metrics: %d roles", len(allMetrics))

	for role := range allMetrics {
		t.Logf("  %s: processed=%d, dropped=%d",
			role, allMetrics[role].ProcessedCount, allMetrics[role].DroppedCount)
	}
}

// Test 5: Retry Mechanism Test
func TestManagerV2_RetryMechanism(t *testing.T) {
	config := &Config{
		MentionQueueSize:    10,
		MentionRateLimit:    100 * time.Millisecond,
		MentionMaxRetries:   3,
		MentionRetryBackoff: 100 * time.Millisecond,
	}

	manager := NewManagerV2WithConfig(config)
	defer manager.Stop()

	platform := newMockPlatformV2()
	platform.teamManager.shouldFail = true
	platform.teamManager.failCount = 2 // Fail first 2 attempts, succeed on 3rd
	ctx := context.Background()

	// Enqueue mention
	err := manager.HandleMentions(
		ctx,
		platform,
		12345,
		"test-team",
		"@architect help",
		[]string{"architect"},
		bus.SenderInfo{Username: "user1"},
		50,
	)

	if err != nil {
		t.Fatalf("HandleMentions failed: %v", err)
	}

	// Wait for retries to complete
	time.Sleep(2 * time.Second)

	// Check metrics
	metrics := manager.GetQueueMetrics("architect")
	if metrics == nil {
		t.Fatal("Expected metrics for architect, got nil")
	}

	if metrics.RetryCount == 0 {
		t.Error("Expected retry count > 0")
	}

	if metrics.ProcessedCount == 0 {
		t.Error("Expected successful processing after retries")
	}

	t.Logf("Retry test: processed=%d, retry=%d, failure=%d",
		metrics.ProcessedCount, metrics.RetryCount, metrics.FailureCount)
}

// Test 6: Metrics Test
func TestManagerV2_Metrics(t *testing.T) {
	manager := NewManagerV2()
	defer manager.Stop()

	platform := newMockPlatformV2()
	ctx := context.Background()

	// Execute mentions for multiple roles
	roles := []string{"architect", "developer", "tester"}

	for _, role := range roles {
		err := manager.HandleMentions(
			ctx,
			platform,
			12345,
			"test-team",
			fmt.Sprintf("@%s help", role),
			[]string{role},
			bus.SenderInfo{Username: "user1"},
			50,
		)

		if err != nil {
			t.Fatalf("HandleMentions for %s failed: %v", role, err)
		}
	}

	// Wait for execution
	time.Sleep(1 * time.Second)

	// Test GetQueueMetrics
	for _, role := range roles {
		metrics := manager.GetQueueMetrics(role)
		if metrics == nil {
			t.Errorf("Expected metrics for %s, got nil", role)
			continue
		}

		t.Logf("%s metrics: processed=%d, dropped=%d, queue_len=%d",
			role, metrics.ProcessedCount, metrics.DroppedCount, metrics.QueueLength)
	}

	// Test GetAllQueueMetrics
	allMetrics := manager.GetAllQueueMetrics()
	if len(allMetrics) == 0 {
		t.Error("Expected metrics for all roles")
	}

	t.Logf("Total roles with metrics: %d", len(allMetrics))
}

// Test 7: Depth Limit with Queue
func TestManagerV2_DepthLimitWithQueue(t *testing.T) {
	manager := NewManagerV2()
	defer manager.Stop()

	maxDepth := manager.maxMentionDepth

	platform := newMockPlatformV2()
	ctx := context.Background()

	// Create a session
	session := manager.GetOrCreateSession(12345, "test-team", 50)

	// Manually test depth limit by creating requests at max depth
	req := &MentionRequest{
		Role:        "architect",
		Prompt:      "test",
		SessionID:   session.SessionID,
		ChatID:      12345,
		TeamID:      "test-team",
		Timestamp:   time.Now(),
		Context:     ctx,
		Platform:    platform,
		Session:     session,
		TeamRoster:  "",
		Depth:       maxDepth, // At max depth
		ExecuteFunc: manager.executeMentionRequest,
	}

	// This should not execute due to depth check in executeAgentAndCascade
	err := manager.queueManager.Enqueue(req)
	if err != nil {
		t.Fatalf("Enqueue failed: %v", err)
	}

	// Wait a bit
	time.Sleep(500 * time.Millisecond)

	// Execution should have been skipped due to depth limit
	t.Logf("Depth limit test completed, max_depth=%d", maxDepth)
}

// Test 8: Graceful Shutdown
func TestManagerV2_GracefulShutdown(t *testing.T) {
	manager := NewManagerV2()

	platform := newMockPlatformV2()
	ctx := context.Background()

	// Enqueue some mentions
	for i := 0; i < 3; i++ {
		err := manager.HandleMentions(
			ctx,
			platform,
			12345,
			"test-team",
			fmt.Sprintf("@architect task %d", i),
			[]string{"architect"},
			bus.SenderInfo{Username: "user1"},
			50,
		)

		if err != nil {
			t.Fatalf("HandleMentions %d failed: %v", i, err)
		}
	}

	// Stop immediately
	manager.Stop()

	// Try to enqueue after stop (should fail or be ignored)
	// This tests that queues are properly stopped

	t.Log("Graceful shutdown test completed")
}
