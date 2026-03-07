// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package team

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// createTestTeam creates a test team
func createTestTeamForRouter() *Team {
	return &Team{
		ID:      "test-team-001",
		Name:    "Test Team",
		Pattern: PatternSequential,
		Agents: map[string]*TeamAgent{
			"agent-001": {
				AgentID:      "agent-001",
				Role:         "researcher",
				Capabilities: []string{"web_search"},
				Status:       StatusIdle,
			},
			"agent-002": {
				AgentID:      "agent-002",
				Role:         "coder",
				Capabilities: []string{"code_generation"},
				Status:       StatusIdle,
			},
			"agent-003": {
				AgentID:      "agent-003",
				Role:         "researcher",
				Capabilities: []string{"web_search"},
				Status:       StatusWorking,
			},
		},
		SharedContext: NewSharedContext("test-team-001"),
		Status:        TeamStatusActive,
		CreatedAt:     time.Now(),
	}
}

// TestNewDelegationRouter tests router creation
func TestNewDelegationRouter(t *testing.T) {
	router := NewDelegationRouter(5)

	if router == nil {
		t.Fatal("Expected non-nil router")
	}

	if router.maxDepth != 5 {
		t.Errorf("Expected maxDepth 5, got %d", router.maxDepth)
	}

	if router.delegationMap == nil {
		t.Error("Expected delegationMap to be initialized")
	}
}

// TestNewDelegationRouterDefaultDepth tests default max depth
func TestNewDelegationRouterDefaultDepth(t *testing.T) {
	router := NewDelegationRouter(0)

	if router.maxDepth != 5 {
		t.Errorf("Expected default maxDepth 5, got %d", router.maxDepth)
	}
}

// TestRouteTaskSuccess tests successful task routing
func TestRouteTaskSuccess(t *testing.T) {
	testTeam := createTestTeamForRouter()
	router := NewDelegationRouter(5)
	ctx := context.Background()

	task := &Task{
		ID:           "task-001",
		Description:  "Test task",
		RequiredRole: "researcher",
	}

	agentID, err := router.RouteTask(ctx, task, testTeam)
	if err != nil {
		t.Fatalf("RouteTask failed: %v", err)
	}

	if agentID == "" {
		t.Error("Expected non-empty agent ID")
	}

	// Verify agent has correct role
	agent := testTeam.Agents[agentID]
	if agent.Role != "researcher" {
		t.Errorf("Expected role 'researcher', got '%s'", agent.Role)
	}
}

// TestRouteTaskNoAvailableAgent tests routing when no agent is available
func TestRouteTaskNoAvailableAgent(t *testing.T) {
	testTeam := createTestTeamForRouter()
	router := NewDelegationRouter(5)
	ctx := context.Background()

	task := &Task{
		ID:           "task-001",
		Description:  "Test task",
		RequiredRole: "nonexistent-role",
	}

	_, err := router.RouteTask(ctx, task, testTeam)
	if err == nil {
		t.Error("Expected error when no agent with required role")
	}
}

// TestValidateDelegationSuccess tests successful delegation validation
func TestValidateDelegationSuccess(t *testing.T) {
	router := NewDelegationRouter(5)

	task := &Task{
		ID:              "task-001",
		DelegationChain: []string{"agent-001", "agent-002"},
	}

	err := router.ValidateDelegation(task, "agent-003")
	if err != nil {
		t.Errorf("ValidateDelegation failed: %v", err)
	}
}

// TestValidateDelegationCircular tests circular delegation detection
func TestValidateDelegationCircular(t *testing.T) {
	router := NewDelegationRouter(5)

	task := &Task{
		ID:              "task-001",
		DelegationChain: []string{"agent-001", "agent-002", "agent-003"},
	}

	// Try to delegate back to agent-002 (circular)
	err := router.ValidateDelegation(task, "agent-002")
	if err == nil {
		t.Error("Expected error for circular delegation")
	}
}

// TestValidateDelegationMaxDepth tests max depth enforcement
func TestValidateDelegationMaxDepth(t *testing.T) {
	router := NewDelegationRouter(3)

	task := &Task{
		ID:              "task-001",
		DelegationChain: []string{"agent-001", "agent-002", "agent-003"},
	}

	// Chain length is 3, which equals maxDepth, so next delegation should fail
	err := router.ValidateDelegation(task, "agent-004")
	if err == nil {
		t.Error("Expected error for max depth exceeded")
	}
}

// TestRecordDelegation tests delegation recording
func TestRecordDelegation(t *testing.T) {
	router := NewDelegationRouter(5)

	task := &Task{
		ID:              "task-001",
		DelegationChain: []string{"agent-001"},
	}

	router.RecordDelegation(task, "agent-002")

	// Verify task delegation chain was updated
	if len(task.DelegationChain) != 2 {
		t.Errorf("Expected delegation chain length 2, got %d", len(task.DelegationChain))
	}

	if task.DelegationChain[1] != "agent-002" {
		t.Errorf("Expected last agent 'agent-002', got '%s'", task.DelegationChain[1])
	}

	// Verify delegation was recorded in map
	chain, exists := router.GetDelegationChain("task-001")
	if !exists {
		t.Error("Expected delegation chain to be recorded")
	}

	if chain.Depth != 2 {
		t.Errorf("Expected depth 2, got %d", chain.Depth)
	}
}

// TestGetDelegationChain tests retrieving delegation chain
func TestGetDelegationChain(t *testing.T) {
	router := NewDelegationRouter(5)

	task := &Task{
		ID:              "task-001",
		DelegationChain: []string{},
	}

	router.RecordDelegation(task, "agent-001")
	router.RecordDelegation(task, "agent-002")

	chain, exists := router.GetDelegationChain("task-001")
	if !exists {
		t.Fatal("Expected delegation chain to exist")
	}

	if chain.TaskID != "task-001" {
		t.Errorf("Expected task ID 'task-001', got '%s'", chain.TaskID)
	}

	if chain.Depth != 2 {
		t.Errorf("Expected depth 2, got %d", chain.Depth)
	}

	if len(chain.Chain) != 2 {
		t.Errorf("Expected chain length 2, got %d", len(chain.Chain))
	}
}

// TestClearDelegationChain tests clearing delegation chain
func TestClearDelegationChain(t *testing.T) {
	router := NewDelegationRouter(5)

	task := &Task{
		ID:              "task-001",
		DelegationChain: []string{},
	}

	router.RecordDelegation(task, "agent-001")

	// Verify chain exists
	_, exists := router.GetDelegationChain("task-001")
	if !exists {
		t.Error("Expected delegation chain to exist before clearing")
	}

	// Clear chain
	router.ClearDelegationChain("task-001")

	// Verify chain is cleared
	_, exists = router.GetDelegationChain("task-001")
	if exists {
		t.Error("Expected delegation chain to be cleared")
	}
}

// TestGetStats tests delegation statistics
func TestGetStats(t *testing.T) {
	router := NewDelegationRouter(5)

	// Record some delegations
	task1 := &Task{ID: "task-001", DelegationChain: []string{}}
	task2 := &Task{ID: "task-002", DelegationChain: []string{}}

	router.RecordDelegation(task1, "agent-001")
	router.RecordDelegation(task1, "agent-002")
	router.RecordDelegation(task2, "agent-003")

	stats := router.GetStats()

	totalChains := stats["total_active_chains"].(int)
	if totalChains != 2 {
		t.Errorf("Expected 2 active chains, got %d", totalChains)
	}

	maxDepthSeen := stats["max_depth_seen"].(int)
	if maxDepthSeen != 2 {
		t.Errorf("Expected max depth seen 2, got %d", maxDepthSeen)
	}

	maxDepthLimit := stats["max_depth_limit"].(int)
	if maxDepthLimit != 5 {
		t.Errorf("Expected max depth limit 5, got %d", maxDepthLimit)
	}
}

// TestConcurrentDelegationRecording tests concurrent delegation operations
func TestConcurrentDelegationRecording(t *testing.T) {
	router := NewDelegationRouter(10)

	done := make(chan bool)

	// Record delegations concurrently
	for i := 0; i < 10; i++ {
		go func(id int) {
			task := &Task{
				ID:              fmt.Sprintf("task-%d", id),
				DelegationChain: []string{},
			}
			router.RecordDelegation(task, "agent-001")
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	stats := router.GetStats()
	totalChains := stats["total_active_chains"].(int)

	if totalChains != 10 {
		t.Errorf("Expected 10 active chains, got %d", totalChains)
	}
}
