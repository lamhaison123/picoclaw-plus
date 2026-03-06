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

	"github.com/sipeed/picoclaw/pkg/agent"
	"github.com/sipeed/picoclaw/pkg/bus"
	"github.com/sipeed/picoclaw/pkg/config"
)

// Mock implementations for testing

type mockAgentRegistry struct{}

func (m *mockAgentRegistry) GetAgent(agentID string) (*agent.AgentInstance, bool) {
	return nil, false
}

func (m *mockAgentRegistry) DeregisterAgent(agentID string) error {
	return nil
}

type publishedMessage struct {
	channel string
	data    []byte
}

type subscription struct {
	channel string
	handler func(string, []byte)
}

type mockMessageBus struct {
	published     []publishedMessage
	subscriptions []subscription
	unsubscribed  []string
}

func (m *mockMessageBus) Publish(channel string, data []byte) error {
	m.published = append(m.published, publishedMessage{channel: channel, data: data})
	return nil
}

func (m *mockMessageBus) Subscribe(channel string, handler func(string, []byte)) error {
	m.subscriptions = append(m.subscriptions, subscription{channel: channel, handler: handler})
	return nil
}

func (m *mockMessageBus) Unsubscribe(channel string) error {
	m.unsubscribed = append(m.unsubscribed, channel)
	return nil
}

// createTestTeamManager creates a TeamManager for testing
func createTestTeamManager() *TeamManager {
	cfg := &config.Config{
		Agents: config.AgentsConfig{
			Defaults: config.AgentDefaults{
				Workspace: ".", // Use current directory for tests
			},
		},
	}
	provider := NewMockProvider() // Use MockProvider from test_helpers.go
	registry := agent.NewAgentRegistry(cfg, provider)
	msgBus := bus.NewMessageBus()

	return NewTeamManager(registry, msgBus)
}

// createTestTeamConfig creates a valid team configuration for testing
func createTestTeamConfig() *TeamConfig {
	return &TeamConfig{
		TeamID:      "test-team-001",
		Name:        "Test Team",
		Description: "A test team",
		Pattern:     "sequential",
		Roles: []RoleConfig{
			{
				Name:         "researcher",
				Description:  "Research role",
				Capabilities: []string{"web_search", "document_analysis"},
				Tools:        []string{"webSearch", "webFetch", "readFile"},
				Model:        "gpt-4o",
			},
			{
				Name:         "coder",
				Description:  "Coding role",
				Capabilities: []string{"code_generation", "testing"},
				Tools:        []string{"readCode", "editCode", "fsWrite"},
				Model:        "claude-sonnet-4",
			},
		},
		Coordinator: CoordinatorConfig{
			Role: "researcher",
		},
		Settings: SettingsConfig{
			MaxDelegationDepth:      5,
			AgentTimeoutSeconds:     30,
			FailureThreshold:        3,
			ConsensusTimeoutSeconds: 60,
		},
	}
}

// TestNewTeamManager tests TeamManager creation
func TestNewTeamManager(t *testing.T) {
	tm := createTestTeamManager()

	if tm == nil {
		t.Fatal("Expected non-nil TeamManager")
	}

	if tm.teams == nil {
		t.Error("Expected teams map to be initialized")
	}

	if tm.roleCapabilities == nil {
		t.Error("Expected roleCapabilities map to be initialized")
	}
}

// TestCreateTeamSuccess tests successful team creation
func TestCreateTeamSuccess(t *testing.T) {
	tm := createTestTeamManager()
	ctx := context.Background()
	config := createTestTeamConfig()

	team, err := tm.CreateTeam(ctx, config)
	if err != nil {
		t.Fatalf("CreateTeam failed: %v", err)
	}

	if team == nil {
		t.Fatal("Expected non-nil team")
	}

	if team.ID != config.TeamID {
		t.Errorf("Expected team ID '%s', got '%s'", config.TeamID, team.ID)
	}

	if team.Name != config.Name {
		t.Errorf("Expected team name '%s', got '%s'", config.Name, team.Name)
	}

	if team.Pattern != CollaborationPattern(config.Pattern) {
		t.Errorf("Expected pattern '%s', got '%s'", config.Pattern, team.Pattern)
	}

	if team.Status != TeamStatusActive {
		t.Errorf("Expected status Active, got %s", team.Status)
	}

	if len(team.Agents) != len(config.Roles) {
		t.Errorf("Expected %d agents, got %d", len(config.Roles), len(team.Agents))
	}

	if team.SharedContext == nil {
		t.Error("Expected SharedContext to be initialized")
	}

	if team.CoordinatorID == "" {
		t.Error("Expected coordinator to be set")
	}
}

// TestCreateTeamDuplicateID tests creating team with duplicate ID
func TestCreateTeamDuplicateID(t *testing.T) {
	tm := createTestTeamManager()
	ctx := context.Background()
	config := createTestTeamConfig()

	// Create first team
	_, err := tm.CreateTeam(ctx, config)
	if err != nil {
		t.Fatalf("First CreateTeam failed: %v", err)
	}

	// Try to create team with same ID
	_, err = tm.CreateTeam(ctx, config)
	if err == nil {
		t.Error("Expected error when creating team with duplicate ID")
	}
}

// TestCreateTeamInvalidConfig tests team creation with invalid configuration
func TestCreateTeamInvalidConfig(t *testing.T) {
	tm := createTestTeamManager()
	ctx := context.Background()

	testCases := []struct {
		name   string
		config *TeamConfig
	}{
		{
			name: "missing team ID",
			config: &TeamConfig{
				Name:    "Test",
				Pattern: "sequential",
				Roles:   []RoleConfig{{Name: "test", Capabilities: []string{"test"}}},
			},
		},
		{
			name: "missing name",
			config: &TeamConfig{
				TeamID:  "test-001",
				Pattern: "sequential",
				Roles:   []RoleConfig{{Name: "test", Capabilities: []string{"test"}}},
			},
		},
		{
			name: "missing pattern",
			config: &TeamConfig{
				TeamID: "test-001",
				Name:   "Test",
				Roles:  []RoleConfig{{Name: "test", Capabilities: []string{"test"}}},
			},
		},
		{
			name: "invalid pattern",
			config: &TeamConfig{
				TeamID:  "test-001",
				Name:    "Test",
				Pattern: "invalid",
				Roles:   []RoleConfig{{Name: "test", Capabilities: []string{"test"}}},
			},
		},
		{
			name: "no roles",
			config: &TeamConfig{
				TeamID:  "test-001",
				Name:    "Test",
				Pattern: "sequential",
				Roles:   []RoleConfig{},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tm.CreateTeam(ctx, tc.config)
			if err == nil {
				t.Errorf("Expected error for %s", tc.name)
			}
		})
	}
}

// TestDissolveTeam tests team dissolution
func TestDissolveTeam(t *testing.T) {
	tm := createTestTeamManager()
	ctx := context.Background()
	config := createTestTeamConfig()

	// Create team
	team, err := tm.CreateTeam(ctx, config)
	if err != nil {
		t.Fatalf("CreateTeam failed: %v", err)
	}

	// Dissolve team
	err = tm.DissolveTeam(ctx, team.ID)
	if err != nil {
		t.Fatalf("DissolveTeam failed: %v", err)
	}

	// Verify team is dissolved
	_, err = tm.GetTeamStatus(team.ID)
	if err == nil {
		t.Error("Expected error when getting status of dissolved team")
	}
}

// TestDissolveNonexistentTeam tests dissolving a team that doesn't exist
func TestDissolveNonexistentTeam(t *testing.T) {
	tm := createTestTeamManager()
	ctx := context.Background()

	err := tm.DissolveTeam(ctx, "nonexistent-team")
	if err == nil {
		t.Error("Expected error when dissolving nonexistent team")
	}
}

// TestAddAgent tests adding an agent to a team
func TestAddAgent(t *testing.T) {
	tm := createTestTeamManager()
	ctx := context.Background()
	config := createTestTeamConfig()

	// Create team
	team, err := tm.CreateTeam(ctx, config)
	if err != nil {
		t.Fatalf("CreateTeam failed: %v", err)
	}

	initialAgentCount := len(team.Agents)

	// Add new agent
	agentConfig := AgentConfig{
		AgentID: "new-agent-001",
		Role:    "researcher",
	}

	err = tm.AddAgent(ctx, team.ID, agentConfig)
	if err != nil {
		t.Fatalf("AddAgent failed: %v", err)
	}

	// Verify agent was added
	status, err := tm.GetTeamStatus(team.ID)
	if err != nil {
		t.Fatalf("GetTeamStatus failed: %v", err)
	}

	if status.AgentCount != initialAgentCount+1 {
		t.Errorf("Expected %d agents, got %d", initialAgentCount+1, status.AgentCount)
	}
}

// TestAddAgentInvalidRole tests adding agent with invalid role
func TestAddAgentInvalidRole(t *testing.T) {
	tm := createTestTeamManager()
	ctx := context.Background()
	config := createTestTeamConfig()

	// Create team
	team, err := tm.CreateTeam(ctx, config)
	if err != nil {
		t.Fatalf("CreateTeam failed: %v", err)
	}

	// Try to add agent with invalid role
	agentConfig := AgentConfig{
		AgentID: "new-agent-001",
		Role:    "invalid-role",
	}

	err = tm.AddAgent(ctx, team.ID, agentConfig)
	if err == nil {
		t.Error("Expected error when adding agent with invalid role")
	}
}

// TestRemoveAgent tests removing an agent from a team
func TestRemoveAgent(t *testing.T) {
	tm := createTestTeamManager()
	ctx := context.Background()
	config := createTestTeamConfig()

	// Create team
	team, err := tm.CreateTeam(ctx, config)
	if err != nil {
		t.Fatalf("CreateTeam failed: %v", err)
	}

	// Get first agent ID
	var agentID string
	for id := range team.Agents {
		agentID = id
		break
	}

	initialAgentCount := len(team.Agents)

	// Remove agent
	err = tm.RemoveAgent(ctx, team.ID, agentID)
	if err != nil {
		t.Fatalf("RemoveAgent failed: %v", err)
	}

	// Verify agent was removed
	status, err := tm.GetTeamStatus(team.ID)
	if err != nil {
		t.Fatalf("GetTeamStatus failed: %v", err)
	}

	if status.AgentCount != initialAgentCount-1 {
		t.Errorf("Expected %d agents, got %d", initialAgentCount-1, status.AgentCount)
	}
}

// TestGetTeamStatus tests retrieving team status
func TestGetTeamStatus(t *testing.T) {
	tm := createTestTeamManager()
	ctx := context.Background()
	config := createTestTeamConfig()

	// Create team
	team, err := tm.CreateTeam(ctx, config)
	if err != nil {
		t.Fatalf("CreateTeam failed: %v", err)
	}

	// Get status
	status, err := tm.GetTeamStatus(team.ID)
	if err != nil {
		t.Fatalf("GetTeamStatus failed: %v", err)
	}

	if status.TeamID != team.ID {
		t.Errorf("Expected team ID '%s', got '%s'", team.ID, status.TeamID)
	}

	if status.TeamName != team.Name {
		t.Errorf("Expected team name '%s', got '%s'", team.Name, status.TeamName)
	}

	if status.Status != TeamStatusActive {
		t.Errorf("Expected status Active, got %s", status.Status)
	}

	if status.AgentCount != len(team.Agents) {
		t.Errorf("Expected %d agents, got %d", len(team.Agents), status.AgentCount)
	}

	if status.Uptime < 0 {
		t.Error("Expected non-negative uptime")
	}
}

// TestValidateToolAccess tests tool access validation
func TestValidateToolAccess(t *testing.T) {
	tm := createTestTeamManager()
	ctx := context.Background()
	config := createTestTeamConfig()

	// Create team
	team, err := tm.CreateTeam(ctx, config)
	if err != nil {
		t.Fatalf("CreateTeam failed: %v", err)
	}

	// Get researcher agent ID
	var researcherID string
	for id, agent := range team.Agents {
		if agent.Role == "researcher" {
			researcherID = id
			break
		}
	}

	// Test allowed tool
	if !tm.ValidateToolAccess(researcherID, "webSearch") {
		t.Error("Expected webSearch to be allowed for researcher")
	}

	// Test disallowed tool
	if tm.ValidateToolAccess(researcherID, "editCode") {
		t.Error("Expected editCode to be disallowed for researcher")
	}

	// Test wildcard (if we add file_* pattern)
	// This would require modifying the test config
}

// TestValidateToolAccessWildcard tests wildcard tool permissions
func TestValidateToolAccessWildcard(t *testing.T) {
	tm := createTestTeamManager()

	// Test wildcard matching
	allowedTools := []string{"file_*", "read*"}

	if !tm.checkToolPermission("file_read", allowedTools) {
		t.Error("Expected file_read to match file_*")
	}

	if !tm.checkToolPermission("file_write", allowedTools) {
		t.Error("Expected file_write to match file_*")
	}

	if !tm.checkToolPermission("readFile", allowedTools) {
		t.Error("Expected readFile to match read*")
	}

	if tm.checkToolPermission("writeFile", allowedTools) {
		t.Error("Expected writeFile to not match patterns")
	}
}

// TestConcurrentTeamOperations tests concurrent team operations
func TestConcurrentTeamOperations(t *testing.T) {
	tm := createTestTeamManager()
	ctx := context.Background()

	// Create multiple teams concurrently
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(id int) {
			config := createTestTeamConfig()
			config.TeamID = fmt.Sprintf("team-%d", id)
			_, _ = tm.CreateTeam(ctx, config)
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify teams were created
	tm.mu.RLock()
	teamCount := len(tm.teams)
	tm.mu.RUnlock()

	if teamCount == 0 {
		t.Error("Expected at least some teams to be created")
	}
}

func TestTeamManager_UpdateAgentStatus(t *testing.T) {
	tm := createTestTeamManager()

	config := &TeamConfig{
		TeamID:  "team1",
		Name:    "Test Team",
		Pattern: "sequential",
		Roles: []RoleConfig{
			{Name: "developer", Capabilities: []string{"code"}, Tools: []string{"file_*"}},
		},
		Coordinator: CoordinatorConfig{
			Role: "developer",
		},
		Settings: SettingsConfig{
			MaxDelegationDepth:      5,
			AgentTimeoutSeconds:     30,
			FailureThreshold:        3,
			ConsensusTimeoutSeconds: 30,
		},
	}

	ctx := context.Background()
	team, _ := tm.CreateTeam(ctx, config)

	agentConfig := AgentConfig{
		AgentID: "agent1",
		Role:    "developer",
	}
	tm.AddAgent(ctx, team.ID, agentConfig)

	err := tm.UpdateAgentStatus(team.ID, "agent1", StatusWorking)

	if err != nil {
		t.Fatalf("UpdateAgentStatus failed: %v", err)
	}

	agent, _ := team.Agents["agent1"]
	if agent.Status != StatusWorking {
		t.Errorf("Expected status %s, got %s", StatusWorking, agent.Status)
	}
}

func TestTeamManager_RecordTaskMetrics(t *testing.T) {
	tm := createTestTeamManager()

	config := &TeamConfig{
		TeamID:  "team1",
		Name:    "Test Team",
		Pattern: "sequential",
		Roles: []RoleConfig{
			{Name: "developer", Capabilities: []string{"code"}, Tools: []string{"file_*"}},
		},
	}

	ctx := context.Background()
	team, _ := tm.CreateTeam(ctx, config)

	task := NewTask("Test task", "developer", nil)
	task.MarkAssigned("agent1")
	task.MarkInProgress()
	task.MarkCompleted("result")

	err := tm.RecordTaskMetrics(team.ID, task)

	if err != nil {
		t.Fatalf("RecordTaskMetrics failed: %v", err)
	}

	metricsKey := fmt.Sprintf("metrics_task_%s", task.ID)
	metrics, exists := team.SharedContext.Get(metricsKey)

	if !exists {
		t.Error("Expected metrics to be recorded in shared context")
	}
	if metrics == nil {
		t.Error("Expected non-nil metrics")
	}
}

func TestTeamManager_DetectUnresponsiveAgent(t *testing.T) {
	tm := createTestTeamManager()

	config := &TeamConfig{
		TeamID:  "team1",
		Name:    "Test Team",
		Pattern: "sequential",
		Roles: []RoleConfig{
			{Name: "developer", Capabilities: []string{"code"}, Tools: []string{"file_*"}},
		},
	}

	ctx := context.Background()
	team, _ := tm.CreateTeam(ctx, config)
	tm.AddAgent(ctx, team.ID, AgentConfig{AgentID: "agent1", Role: "developer"})

	// Set last active to past
	agent := team.Agents["agent1"]
	agent.LastActive = time.Now().Add(-2 * time.Minute)

	unresponsive, err := tm.DetectUnresponsiveAgent(team.ID, "agent1", 1*time.Minute)

	if err != nil {
		t.Fatalf("DetectUnresponsiveAgent failed: %v", err)
	}
	if !unresponsive {
		t.Error("Expected agent to be detected as unresponsive")
	}
	if agent.Status != StatusUnresponsive {
		t.Errorf("Expected status %s, got %s", StatusUnresponsive, agent.Status)
	}
}

func TestTeamManager_IncrementAgentFailure(t *testing.T) {
	tm := createTestTeamManager()

	config := &TeamConfig{
		TeamID:  "team1",
		Name:    "Test Team",
		Pattern: "sequential",
		Roles: []RoleConfig{
			{Name: "developer", Capabilities: []string{"code"}, Tools: []string{"file_*"}},
		},
	}

	ctx := context.Background()
	team, _ := tm.CreateTeam(ctx, config)
	tm.AddAgent(ctx, team.ID, AgentConfig{AgentID: "agent1", Role: "developer"})

	// Increment failures
	removed, _ := tm.IncrementAgentFailure(team.ID, "agent1", 3)
	if removed {
		t.Error("Expected agent not to be removed after 1 failure")
	}

	removed, _ = tm.IncrementAgentFailure(team.ID, "agent1", 3)
	if removed {
		t.Error("Expected agent not to be removed after 2 failures")
	}

	removed, _ = tm.IncrementAgentFailure(team.ID, "agent1", 3)
	if !removed {
		t.Error("Expected agent to be removed after 3 failures")
	}

	_, exists := team.Agents["agent1"]
	if exists {
		t.Error("Expected agent to be removed from team")
	}
}

func TestTeamManager_AddAgent_ToActiveTeam(t *testing.T) {
	tm := createTestTeamManager()

	config := &TeamConfig{
		TeamID:  "team1",
		Name:    "Test Team",
		Pattern: "sequential",
		Roles: []RoleConfig{
			{Name: "developer", Capabilities: []string{"code"}, Tools: []string{"file_*"}},
			{Name: "tester", Capabilities: []string{"test"}, Tools: []string{"test_*"}},
		},
	}

	ctx := context.Background()
	team, _ := tm.CreateTeam(ctx, config)
	team.Status = TeamStatusActive // Add agent to active team
	err := tm.AddAgent(ctx, team.ID, AgentConfig{AgentID: "agent2", Role: "tester"})

	if err != nil {
		t.Fatalf("AddAgent to active team failed: %v", err)
	}

	if len(team.Agents) != 3 {
		t.Errorf("Expected 3 agents (2 from config + 1 added), got %d", len(team.Agents))
	}

	// Check notification in shared context
	history := team.SharedContext.GetHistory()
	found := false
	for _, entry := range history {
		if entry.Action == "agent_added" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected agent_added notification in shared context")
	}
}

func TestTeamManager_RemoveAgent_WithTaskReassignment(t *testing.T) {
	tm := createTestTeamManager()

	config := &TeamConfig{
		TeamID:  "team1",
		Name:    "Test Team",
		Pattern: "sequential",
		Roles: []RoleConfig{
			{Name: "developer", Capabilities: []string{"code"}, Tools: []string{"file_*"}},
		},
	}

	ctx := context.Background()
	team, _ := tm.CreateTeam(ctx, config)

	tm.AddAgent(ctx, team.ID, AgentConfig{AgentID: "agent1", Role: "developer"})
	tm.AddAgent(ctx, team.ID, AgentConfig{AgentID: "agent2", Role: "developer"})

	// Simulate active task for agent1
	team.SharedContext.Set("task_active_1", map[string]any{
		"task_id":  "task1",
		"agent_id": "agent1",
		"status":   string(TaskStatusInProgress),
	}, "system")

	// Remove agent1
	err := tm.RemoveAgent(ctx, team.ID, "agent1")

	if err != nil {
		t.Fatalf("RemoveAgent failed: %v", err)
	}

	// Check agent was removed
	if _, exists := team.Agents["agent1"]; exists {
		t.Error("Expected agent1 to be removed")
	}

	// Check notification
	history := team.SharedContext.GetHistory()
	foundRemoval := false
	foundReassignment := false
	for _, entry := range history {
		if entry.Action == "agent_removed" {
			foundRemoval = true
		}
		if entry.Action == "task_reassigned" {
			foundReassignment = true
		}
	}
	if !foundRemoval {
		t.Error("Expected agent_removed notification")
	}
	if !foundReassignment {
		t.Error("Expected task_reassigned notification")
	}
}

func TestTeamManager_RemoveAgent_NoReassignmentNeeded(t *testing.T) {
	tm := createTestTeamManager()

	config := &TeamConfig{
		TeamID:  "team1",
		Name:    "Test Team",
		Pattern: "sequential",
		Roles: []RoleConfig{
			{Name: "developer", Capabilities: []string{"code"}, Tools: []string{"file_*"}},
		},
	}

	ctx := context.Background()
	team, _ := tm.CreateTeam(ctx, config)

	tm.AddAgent(ctx, team.ID, AgentConfig{AgentID: "agent1", Role: "developer"})

	// Remove agent without active tasks
	err := tm.RemoveAgent(ctx, team.ID, "agent1")

	if err != nil {
		t.Fatalf("RemoveAgent failed: %v", err)
	}

	if _, exists := team.Agents["agent1"]; exists {
		t.Error("Expected agent1 to be removed")
	}
}

func TestTeamManager_SendTaskDelegation(t *testing.T) {
	tm := createTestTeamManager()

	task := NewTask("Test task", "developer", nil)
	msg := NewTaskDelegationMessage("team1", "coordinator", "agent1", task, nil)

	err := tm.SendTaskDelegation("team1", msg)

	if err != nil {
		t.Fatalf("SendTaskDelegation failed: %v", err)
	}

	// Message was sent successfully (no error)
}

func TestTeamManager_SendTaskResult(t *testing.T) {
	tm := createTestTeamManager()

	msg := NewTaskResultMessage("team1", "agent1", "coordinator", "task1", TaskStatusCompleted, "result", nil)

	err := tm.SendTaskResult("team1", msg)

	if err != nil {
		t.Fatalf("SendTaskResult failed: %v", err)
	}

	// Message was sent successfully (no error)
}

func TestTeamManager_SendConsensusRequest(t *testing.T) {
	tm := createTestTeamManager()

	msg := NewConsensusRequestMessage("team1", "coordinator", "consensus1", "Question?", []string{"yes", "no"}, VotingRuleMajority, 30, nil)

	err := tm.SendConsensusRequest("team1", msg)

	if err != nil {
		t.Fatalf("SendConsensusRequest failed: %v", err)
	}

	// Message was sent successfully (no error)
}

func TestTeamManager_SubscribeToTeamMessages(t *testing.T) {
	tm := createTestTeamManager()

	handler := func(channel string, data []byte) {
		// Handler would be called when messages arrive
	}

	err := tm.SubscribeToTeamMessages("team1", handler)

	if err != nil {
		t.Fatalf("SubscribeToTeamMessages failed: %v", err)
	}

	// Subscription successful (no error - actual pub/sub not implemented yet)
}

func TestTeamManager_UnsubscribeFromTeamMessages(t *testing.T) {
	tm := createTestTeamManager()

	handler := func(channel string, data []byte) {}
	tm.SubscribeToTeamMessages("team1", handler)

	err := tm.UnsubscribeFromTeamMessages("team1")

	if err != nil {
		t.Fatalf("UnsubscribeFromTeamMessages failed: %v", err)
	}

	// Unsubscription successful (no error - actual pub/sub not implemented yet)
}
