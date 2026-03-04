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

// BenchmarkTeamCreation benchmarks team creation with 10 agents
// Target: < 100ms
func BenchmarkTeamCreation(b *testing.B) {
	// Setup
	cfg := &config.Config{}
	prov := NewMockProvider()
	registry := agent.NewAgentRegistry(cfg, prov)
	msgBus := bus.NewMessageBus()
	tm := NewTeamManager(registry, msgBus)

	teamConfig := &TeamConfig{
		TeamID:  "bench-team",
		Name:    "Benchmark Team",
		Pattern: "sequential",
		Roles: []RoleConfig{
			{Name: "role1", Capabilities: []string{"cap1"}, Tools: []string{"tool1"}},
			{Name: "role2", Capabilities: []string{"cap2"}, Tools: []string{"tool2"}},
			{Name: "role3", Capabilities: []string{"cap3"}, Tools: []string{"tool3"}},
			{Name: "role4", Capabilities: []string{"cap4"}, Tools: []string{"tool4"}},
			{Name: "role5", Capabilities: []string{"cap5"}, Tools: []string{"tool5"}},
			{Name: "role6", Capabilities: []string{"cap6"}, Tools: []string{"tool6"}},
			{Name: "role7", Capabilities: []string{"cap7"}, Tools: []string{"tool7"}},
			{Name: "role8", Capabilities: []string{"cap8"}, Tools: []string{"tool8"}},
			{Name: "role9", Capabilities: []string{"cap9"}, Tools: []string{"tool9"}},
			{Name: "role10", Capabilities: []string{"cap10"}, Tools: []string{"tool10"}},
		},
		Coordinator: CoordinatorConfig{
			Role: "role1",
		},
		Settings: SettingsConfig{
			MaxDelegationDepth:      5,
			AgentTimeoutSeconds:     30,
			FailureThreshold:        3,
			ConsensusTimeoutSeconds: 30,
		},
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Use unique team ID for each iteration
		teamConfig.TeamID = fmt.Sprintf("bench-team-%d", i)

		start := time.Now()
		team, err := tm.CreateTeam(ctx, teamConfig)
		elapsed := time.Since(start)

		if err != nil {
			b.Fatalf("Failed to create team: %v", err)
		}

		// Check if within target (100ms)
		if elapsed > 100*time.Millisecond {
			b.Logf("Warning: Team creation took %v (target: <100ms)", elapsed)
		}

		// Cleanup
		tm.DissolveTeam(ctx, team.ID)
	}
}

// BenchmarkMessageDelivery benchmarks task delegation message delivery
// Target: < 10ms
func BenchmarkMessageDelivery(b *testing.B) {
	msgBus := bus.NewMessageBus()
	teamID := "bench-team"

	msg := &TaskDelegationMessage{
		MessageID:   "msg-1",
		TeamID:      teamID,
		FromAgentID: "agent-1",
		ToAgentID:   "agent-2",
		Task:        NewTask("Test task", "developer", nil),
		Context:     map[string]any{"key": "value"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		start := time.Now()

		// Serialize and publish message
		data, _ := msg.ToJSON()
		outboundMsg := bus.OutboundMessage{
			Content: string(data),
		}
		_ = msgBus.PublishOutbound(context.Background(), outboundMsg)

		elapsed := time.Since(start)

		// Check if within target (10ms)
		if elapsed > 10*time.Millisecond {
			b.Logf("Warning: Message delivery took %v (target: <10ms)", elapsed)
		}
	}
}

// BenchmarkSharedContextRead benchmarks concurrent read operations
// Target: < 1ms per read
func BenchmarkSharedContextRead(b *testing.B) {
	ctx := NewSharedContext("bench-team")

	// Populate context with data
	for i := 0; i < 100; i++ {
		ctx.Set(fmt.Sprintf("key-%d", i), fmt.Sprintf("value-%d", i), "bench-agent")
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			start := time.Now()

			// Read operation
			_, exists := ctx.Get("key-50")

			elapsed := time.Since(start)

			if !exists {
				b.Error("Expected key to exist")
			}

			// Check if within target (1ms)
			if elapsed > time.Millisecond {
				b.Logf("Warning: Context read took %v (target: <1ms)", elapsed)
			}
		}
	})
}

// BenchmarkSharedContextWrite benchmarks write operations
func BenchmarkSharedContextWrite(b *testing.B) {
	ctx := NewSharedContext("bench-team")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx.Set(fmt.Sprintf("key-%d", i), fmt.Sprintf("value-%d", i), "bench-agent")
	}
}

// BenchmarkAgentPoolReuse benchmarks agent instance reuse
func BenchmarkAgentPoolReuse(b *testing.B) {
	pool := NewAgentPool()
	createFn := func() *agent.AgentInstance {
		return &agent.AgentInstance{}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		instance := pool.GetOrCreateInstance("developer", createFn)
		if instance == nil {
			b.Fatal("Expected instance, got nil")
		}
	}
}

// BenchmarkRoleCacheLookup benchmarks role cache lookups
func BenchmarkRoleCacheLookup(b *testing.B) {
	cache := NewRoleCache()
	cache.Set("developer", []string{"code", "test", "review"})

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, exists := cache.Get("developer")
			if !exists {
				b.Error("Expected role to exist in cache")
			}
		}
	})
}

// BenchmarkTaskDelegation benchmarks task delegation workflow
func BenchmarkTaskDelegation(b *testing.B) {
	// Setup
	cfg := &config.Config{}
	prov := NewMockProvider()
	registry := agent.NewAgentRegistry(cfg, prov)
	msgBus := bus.NewMessageBus()
	tm := NewTeamManager(registry, msgBus)

	teamConfig := &TeamConfig{
		TeamID:  "bench-team",
		Name:    "Benchmark Team",
		Pattern: "sequential",
		Roles: []RoleConfig{
			{Name: "coordinator", Capabilities: []string{"coordinate"}, Tools: []string{"*"}},
			{Name: "developer", Capabilities: []string{"code"}, Tools: []string{"file_*"}},
		},
		Coordinator: CoordinatorConfig{
			Role: "coordinator",
		},
		Settings: SettingsConfig{
			MaxDelegationDepth:      5,
			AgentTimeoutSeconds:     30,
			FailureThreshold:        3,
			ConsensusTimeoutSeconds: 30,
		},
	}

	ctx := context.Background()
	team, err := tm.CreateTeam(ctx, teamConfig)
	if err != nil {
		b.Fatalf("Failed to create team: %v", err)
	}
	defer tm.DissolveTeam(ctx, team.ID)

	router := NewDelegationRouter(5)
	coordinator := NewCoordinatorAgent("coordinator", team.ID, team, PatternSequential, msgBus, router, context.Background())

	task := NewTask("Benchmark task", "developer", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := coordinator.DelegateTask(ctx, task)
		if err != nil {
			b.Fatalf("Failed to delegate task: %v", err)
		}
	}
}

// BenchmarkConsensusVoting benchmarks consensus voting
func BenchmarkConsensusVoting(b *testing.B) {
	cm := NewConsensusManager()

	voters := []string{"agent1", "agent2", "agent3", "agent4", "agent5"}
	options := []string{"yes", "no"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		request, err := cm.InitiateConsensus(
			"Benchmark question",
			options,
			voters,
			VotingRuleMajority,
			30*time.Second,
			nil,
		)
		if err != nil {
			b.Fatalf("Failed to initiate consensus: %v", err)
		}

		// Submit votes
		for _, voter := range voters {
			cm.SubmitVote(request.ID, voter, "yes", 1.0, "benchmark vote")
		}

		// Determine outcome
		cm.DetermineOutcome(request.ID)
	}
}

// BenchmarkMemoryPersistence benchmarks team memory save operations
func BenchmarkMemoryPersistence(b *testing.B) {
	// Create temporary directory for benchmarks
	tmpDir := b.TempDir()

	teamMemory := memory.NewTeamMemory(tmpDir)

	record := &memory.TeamMemoryRecord{
		TeamID:    "bench-team",
		TeamName:  "Benchmark Team",
		Pattern:   "sequential",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Hour),
		SharedContext: map[string]any{
			"key1": "value1",
			"key2": "value2",
		},
		Tasks: []memory.TaskRecord{
			{
				TaskID:      "task-1",
				Description: "Test task",
				Role:        "developer",
				Status:      "completed",
				Result:      "success",
			},
		},
		Outcome: "success",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		record.TeamID = fmt.Sprintf("bench-team-%d", i)
		err := teamMemory.SaveTeamRecord(record)
		if err != nil {
			b.Fatalf("Failed to save team record: %v", err)
		}
	}
}
