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
	"github.com/sipeed/picoclaw/pkg/team/memory"
)

// TestEndToEndSequentialWorkflow tests complete sequential workflow
func TestEndToEndSequentialWorkflow(t *testing.T) {
	// Setup
	cfg := &config.Config{
		Agents: config.AgentsConfig{
			Defaults: config.AgentDefaults{
				Workspace: t.TempDir(), // Use temp directory for test
			},
		},
	}
	prov := NewMockProvider()
	registry := agent.NewAgentRegistry(cfg, prov)
	msgBus := bus.NewMessageBus()
	tm := NewTeamManager(registry, msgBus)

	teamConfig := &TeamConfig{
		TeamID:  "seq-team",
		Name:    "Sequential Team",
		Pattern: "sequential",
		Roles: []RoleConfig{
			{Name: "coordinator", Capabilities: []string{"coordinate"}, Tools: []string{"*"}},
			{Name: "agent1", Capabilities: []string{"task1"}, Tools: []string{"tool1"}},
			{Name: "agent2", Capabilities: []string{"task2"}, Tools: []string{"tool2"}},
			{Name: "agent3", Capabilities: []string{"task3"}, Tools: []string{"tool3"}},
		},
		Coordinator: CoordinatorConfig{
			Role: "coordinator",
		},
		Settings: SettingsConfig{
			MaxDelegationDepth:      5,
			AgentTimeoutSeconds:     30,
			ConsensusTimeoutSeconds: 30, FailureThreshold: 3,
		},
	}

	ctx := context.Background()
	team, err := tm.CreateTeam(ctx, teamConfig)
	if err != nil {
		t.Fatalf("Failed to create team: %v", err)
	}
	defer tm.DissolveTeam(ctx, team.ID)

	// Verify team created
	if team.Status != TeamStatusActive {
		t.Errorf("Expected team status Active, got %v", team.Status)
	}

	// Verify agents registered
	if len(team.Agents) != 4 {
		t.Errorf("Expected 4 agents, got %d", len(team.Agents))
	}

	// Verify shared context initialized
	if team.SharedContext == nil {
		t.Fatal("Expected shared context to be initialized")
	}

	// Test shared context operations
	team.SharedContext.Set("test_key", "test_value", "system")
	value, exists := team.SharedContext.Get("test_key")
	if !exists || value != "test_value" {
		t.Error("Shared context not working correctly")
	}

	// Verify history tracking
	history := team.SharedContext.GetHistory()
	if len(history) == 0 {
		t.Error("Expected history entries")
	}
}

// TestEndToEndParallelWorkflow tests complete parallel workflow
func TestEndToEndParallelWorkflow(t *testing.T) {
	// Setup
	cfg := &config.Config{
		Agents: config.AgentsConfig{
			Defaults: config.AgentDefaults{
				Workspace: t.TempDir(), // Use temp directory for test
			},
		},
	}
	prov := NewMockProvider()
	registry := agent.NewAgentRegistry(cfg, prov)
	msgBus := bus.NewMessageBus()
	tm := NewTeamManager(registry, msgBus)

	teamConfig := &TeamConfig{
		TeamID:  "par-team",
		Name:    "Parallel Team",
		Pattern: "parallel",
		Roles: []RoleConfig{
			{Name: "coordinator", Capabilities: []string{"coordinate"}, Tools: []string{"*"}},
			{Name: "worker1", Capabilities: []string{"work"}, Tools: []string{"tool1"}},
			{Name: "worker2", Capabilities: []string{"work"}, Tools: []string{"tool2"}},
			{Name: "worker3", Capabilities: []string{"work"}, Tools: []string{"tool3"}},
		},
		Coordinator: CoordinatorConfig{
			Role: "coordinator",
		},
		Settings: SettingsConfig{
			MaxDelegationDepth:      5,
			AgentTimeoutSeconds:     30,
			ConsensusTimeoutSeconds: 30, FailureThreshold: 3,
		},
	}

	ctx := context.Background()
	team, err := tm.CreateTeam(ctx, teamConfig)
	if err != nil {
		t.Fatalf("Failed to create team: %v", err)
	}
	defer tm.DissolveTeam(ctx, team.ID)

	// Verify team created with parallel pattern
	if team.Pattern != PatternParallel {
		t.Errorf("Expected parallel pattern, got %v", team.Pattern)
	}

	// Test concurrent shared context access
	done := make(chan bool)
	for i := 0; i < 3; i++ {
		go func(id int) {
			team.SharedContext.Set("key", "value", fmt.Sprintf("agent-%d", id))
			_, exists := team.SharedContext.Get("key")
			if !exists {
				t.Error("Expected key to exist")
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 3; i++ {
		<-done
	}
}

// TestEndToEndHierarchicalWorkflow tests complete hierarchical workflow
func TestEndToEndHierarchicalWorkflow(t *testing.T) {
	// Setup
	cfg := &config.Config{
		Agents: config.AgentsConfig{
			Defaults: config.AgentDefaults{
				Workspace: t.TempDir(), // Use temp directory for test
			},
		},
	}
	prov := NewMockProvider()
	registry := agent.NewAgentRegistry(cfg, prov)
	msgBus := bus.NewMessageBus()
	tm := NewTeamManager(registry, msgBus)

	teamConfig := &TeamConfig{
		TeamID:  "hier-team",
		Name:    "Hierarchical Team",
		Pattern: "hierarchical",
		Roles: []RoleConfig{
			{Name: "lead", Capabilities: []string{"coordinate", "plan"}, Tools: []string{"*"}},
			{Name: "specialist1", Capabilities: []string{"spec1"}, Tools: []string{"tool1"}},
			{Name: "specialist2", Capabilities: []string{"spec2"}, Tools: []string{"tool2"}},
		},
		Coordinator: CoordinatorConfig{
			Role: "lead",
		},
		Settings: SettingsConfig{
			MaxDelegationDepth:      5,
			AgentTimeoutSeconds:     30,
			ConsensusTimeoutSeconds: 30, FailureThreshold: 3,
		},
	}

	ctx := context.Background()
	team, err := tm.CreateTeam(ctx, teamConfig)
	if err != nil {
		t.Fatalf("Failed to create team: %v", err)
	}
	defer tm.DissolveTeam(ctx, team.ID)

	// Verify hierarchical pattern
	if team.Pattern != PatternHierarchical {
		t.Errorf("Expected hierarchical pattern, got %v", team.Pattern)
	}

	// Create coordinator and router
	router := NewDelegationRouter(5)
	coordinator := NewCoordinatorAgent("lead", team.ID, team, PatternHierarchical, msgBus, router, context.Background())

	// Test task decomposition
	mainTask := NewTask("Main task", "lead", nil)

	// Verify coordinator can delegate
	err = coordinator.DelegateTask(ctx, mainTask)
	if err != nil {
		t.Errorf("Failed to delegate task: %v", err)
	}
}

// TestConsensusProtocolIntegration tests consensus protocol end-to-end
func TestConsensusProtocolIntegration(t *testing.T) {
	cm := NewConsensusManager()

	voters := []string{"agent1", "agent2", "agent3"}
	options := []string{"approve", "reject"}

	// Initiate consensus
	request, err := cm.InitiateConsensus(
		"Should we proceed?",
		options,
		voters,
		VotingRuleMajority,
		30*time.Second,
		nil,
	)
	if err != nil {
		t.Fatalf("Failed to initiate consensus: %v", err)
	}

	// Submit votes
	cm.SubmitVote(request.ID, "agent1", "approve", 1.0, "I approve")
	cm.SubmitVote(request.ID, "agent2", "approve", 1.0, "I approve")
	cm.SubmitVote(request.ID, "agent3", "reject", 1.0, "I reject")

	// Determine outcome
	result, err := cm.DetermineOutcome(request.ID)
	if err != nil {
		t.Fatalf("Failed to determine outcome: %v", err)
	}

	// Verify majority wins
	if result.Outcome != "approve" {
		t.Errorf("Expected outcome 'approve', got '%s'", result.Outcome)
	}

	// Verify result details
	if result.TotalVotes != 3 {
		t.Errorf("Expected 3 total votes, got %d", result.TotalVotes)
	}
}

// TestDynamicCompositionIntegration tests dynamic team composition
func TestDynamicCompositionIntegration(t *testing.T) {
	// Setup
	cfg := &config.Config{
		Agents: config.AgentsConfig{
			Defaults: config.AgentDefaults{
				Workspace: t.TempDir(), // Use temp directory for test
			},
		},
	}
	prov := NewMockProvider()
	registry := agent.NewAgentRegistry(cfg, prov)
	msgBus := bus.NewMessageBus()
	tm := NewTeamManager(registry, msgBus)

	teamConfig := &TeamConfig{
		TeamID:  "dynamic-team",
		Name:    "Dynamic Team",
		Pattern: "sequential",
		Roles: []RoleConfig{
			{Name: "coordinator", Capabilities: []string{"coordinate"}, Tools: []string{"*"}},
			{Name: "worker", Capabilities: []string{"work"}, Tools: []string{"tool1"}},
		},
		Coordinator: CoordinatorConfig{
			Role: "coordinator",
		},
		Settings: SettingsConfig{
			MaxDelegationDepth:      5,
			AgentTimeoutSeconds:     30,
			ConsensusTimeoutSeconds: 30, FailureThreshold: 3,
		},
	}

	ctx := context.Background()
	team, err := tm.CreateTeam(ctx, teamConfig)
	if err != nil {
		t.Fatalf("Failed to create team: %v", err)
	}
	defer tm.DissolveTeam(ctx, team.ID)

	// Add agent during execution
	newAgentConfig := AgentConfig{
		AgentID: "new-worker",
		Role:    "worker",
	}

	err = tm.AddAgent(ctx, team.ID, newAgentConfig)
	if err != nil {
		t.Fatalf("Failed to add agent: %v", err)
	}

	// Verify agent added
	updatedTeam, _ := tm.GetTeam(team.ID)
	if len(updatedTeam.Agents) != 3 {
		t.Errorf("Expected 3 agents after addition, got %d", len(updatedTeam.Agents))
	}

	// Remove agent
	err = tm.RemoveAgent(ctx, team.ID, "new-worker")
	if err != nil {
		t.Fatalf("Failed to remove agent: %v", err)
	}

	// Verify agent removed
	updatedTeam, _ = tm.GetTeam(team.ID)
	if len(updatedTeam.Agents) != 2 {
		t.Errorf("Expected 2 agents after removal, got %d", len(updatedTeam.Agents))
	}
}

// TestFailureRecoveryIntegration tests failure detection and recovery
func TestFailureRecoveryIntegration(t *testing.T) {
	// Setup
	cfg := &config.Config{
		Agents: config.AgentsConfig{
			Defaults: config.AgentDefaults{
				Workspace: t.TempDir(), // Use temp directory for test
			},
		},
	}
	prov := NewMockProvider()
	registry := agent.NewAgentRegistry(cfg, prov)
	msgBus := bus.NewMessageBus()
	tm := NewTeamManager(registry, msgBus)

	teamConfig := &TeamConfig{
		TeamID:  "recovery-team",
		Name:    "Recovery Team",
		Pattern: "sequential",
		Roles: []RoleConfig{
			{Name: "coordinator", Capabilities: []string{"coordinate"}, Tools: []string{"*"}},
			{Name: "worker", Capabilities: []string{"work"}, Tools: []string{"tool1"}},
		},
		Coordinator: CoordinatorConfig{
			Role: "coordinator",
		},
		Settings: SettingsConfig{
			MaxDelegationDepth:      5,
			AgentTimeoutSeconds:     30,
			ConsensusTimeoutSeconds: 30, FailureThreshold: 3,
		},
	}

	ctx := context.Background()
	team, err := tm.CreateTeam(ctx, teamConfig)
	if err != nil {
		t.Fatalf("Failed to create team: %v", err)
	}
	defer tm.DissolveTeam(ctx, team.ID)

	// Simulate agent failure
	agentID := "recovery-team-worker"
	err = tm.UpdateAgentStatus(team.ID, agentID, StatusFailed)
	if err != nil {
		t.Fatalf("Failed to update agent status: %v", err)
	}

	// Verify status updated
	updatedTeam, _ := tm.GetTeam(team.ID)
	if updatedTeam.Agents[agentID].Status != StatusFailed {
		t.Error("Expected agent status to be Failed")
	}

	// Test failure threshold
	removed, err := tm.IncrementAgentFailure(team.ID, agentID, 3)
	if err != nil {
		t.Fatalf("Failed to increment failure count: %v", err)
	}

	if updatedTeam.Agents[agentID].FailureCount != 1 {
		t.Errorf("Expected failure count 1, got %d", updatedTeam.Agents[agentID].FailureCount)
	}

	// Increment to threshold
	tm.IncrementAgentFailure(team.ID, agentID, 3)
	removed, _ = tm.IncrementAgentFailure(team.ID, agentID, 3)

	if !removed {
		t.Error("Expected agent to be removed after reaching threshold")
	}
}

// TestAgentRegistryIntegration tests Agent Registry integration
func TestAgentRegistryIntegration(t *testing.T) {
	// Setup
	cfg := &config.Config{
		Agents: config.AgentsConfig{
			Defaults: config.AgentDefaults{
				Workspace: t.TempDir(), // Use temp directory for test
			},
		},
	}
	prov := NewMockProvider()
	registry := agent.NewAgentRegistry(cfg, prov)
	msgBus := bus.NewMessageBus()
	tm := NewTeamManager(registry, msgBus)

	teamConfig := &TeamConfig{
		TeamID:  "registry-team",
		Name:    "Registry Team",
		Pattern: "sequential",
		Roles: []RoleConfig{
			{Name: "coordinator", Capabilities: []string{"coordinate"}, Tools: []string{"*"}},
			{Name: "worker", Capabilities: []string{"work"}, Tools: []string{"tool1"}},
		},
		Coordinator: CoordinatorConfig{
			Role: "coordinator",
		},
		Settings: SettingsConfig{
			MaxDelegationDepth:      5,
			AgentTimeoutSeconds:     30,
			ConsensusTimeoutSeconds: 30, FailureThreshold: 3,
		},
	}

	ctx := context.Background()
	team, err := tm.CreateTeam(ctx, teamConfig)
	if err != nil {
		t.Fatalf("Failed to create team: %v", err)
	}

	// Verify agents registered
	// Note: This would require Agent Registry to expose a method to check registration
	// For now, we verify team creation succeeded

	// Dissolve team
	err = tm.DissolveTeam(ctx, team.ID)
	if err != nil {
		t.Fatalf("Failed to dissolve team: %v", err)
	}

	// Verify team dissolved
	_, err = tm.GetTeam(team.ID)
	if err == nil {
		t.Error("Expected error when getting dissolved team")
	}
}

// TestMessageBusIntegration tests Message Bus integration
func TestMessageBusIntegration(t *testing.T) {
	t.Skip("Skipping: MessageBus pub/sub pattern not yet implemented")

	msgBus := bus.NewMessageBus()
	tm := &TeamManager{
		bus: msgBus,
	}

	teamID := "msg-team"

	// Subscribe to team messages
	received := make(chan bool, 1)
	err := tm.SubscribeToTeamMessages(teamID, func(channel string, data []byte) {
		received <- true
	})
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	// Send task delegation message
	msg := &TaskDelegationMessage{
		MessageID:   "msg-1",
		TeamID:      teamID,
		FromAgentID: "agent-1",
		ToAgentID:   "agent-2",
		Task:        NewTask("Test task", "worker", nil),
		Context:     map[string]any{"key": "value"},
		Timestamp:   time.Now(),
	}

	err = tm.SendTaskDelegation(teamID, msg)
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	// Wait for message delivery
	select {
	case <-received:
		// Success
	case <-time.After(1 * time.Second):
		t.Error("Message not received within timeout")
	}

	// Unsubscribe
	err = tm.UnsubscribeFromTeamMessages(teamID)
	if err != nil {
		t.Fatalf("Failed to unsubscribe: %v", err)
	}
}

// TestTeamMemoryIntegration tests Team Memory persistence
func TestTeamMemoryIntegration(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()

	// Setup
	cfg := &config.Config{
		Agents: config.AgentsConfig{
			Defaults: config.AgentDefaults{
				Workspace: tmpDir,
			},
		},
	}
	prov := NewMockProvider()
	registry := agent.NewAgentRegistry(cfg, prov)
	msgBus := bus.NewMessageBus()
	tm := NewTeamManager(registry, msgBus)

	teamMemory := memory.NewTeamMemory(tmpDir)
	tm.SetTeamMemory(teamMemory)

	teamConfig := &TeamConfig{
		TeamID:  "memory-team",
		Name:    "Memory Team",
		Pattern: "sequential",
		Roles: []RoleConfig{
			{Name: "coordinator", Capabilities: []string{"coordinate"}, Tools: []string{"*"}},
			{Name: "worker", Capabilities: []string{"work"}, Tools: []string{"tool1"}},
		},
		Coordinator: CoordinatorConfig{
			Role: "coordinator",
		},
		Settings: SettingsConfig{
			MaxDelegationDepth:      5,
			AgentTimeoutSeconds:     30,
			ConsensusTimeoutSeconds: 30, FailureThreshold: 3,
		},
	}

	ctx := context.Background()
	team, err := tm.CreateTeam(ctx, teamConfig)
	if err != nil {
		t.Fatalf("Failed to create team: %v", err)
	}

	// Add some data to shared context
	team.SharedContext.Set("test_key", "test_value", "system")

	// Dissolve team (should persist memory)
	err = tm.DissolveTeam(ctx, team.ID)
	if err != nil {
		t.Fatalf("Failed to dissolve team: %v", err)
	}

	// Load team memory
	record, err := teamMemory.LoadTeamRecord(team.ID)
	if err != nil {
		t.Fatalf("Failed to load team record: %v", err)
	}

	// Verify record
	if record.TeamID != team.ID {
		t.Errorf("Expected team ID %s, got %s", team.ID, record.TeamID)
	}

	if record.TeamName != team.Name {
		t.Errorf("Expected team name %s, got %s", team.Name, record.TeamName)
	}

	// Verify shared context persisted
	if value, exists := record.SharedContext["test_key"]; !exists || value != "test_value" {
		t.Error("Shared context not persisted correctly")
	}
}
