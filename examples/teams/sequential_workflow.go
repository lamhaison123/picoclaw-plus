// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

//go:build ignore

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/sipeed/picoclaw/pkg/agent"
	"github.com/sipeed/picoclaw/pkg/bus"
	"github.com/sipeed/picoclaw/pkg/config"
	"github.com/sipeed/picoclaw/pkg/providers"
	"github.com/sipeed/picoclaw/pkg/team"
)

// SequentialWorkflowExample demonstrates a sequential development workflow
// where tasks execute in order: design → implement → test → review
func main() {
	fmt.Println("=== Sequential Workflow Example ===\n")

	// Initialize dependencies
	cfg := &config.Config{}
	prov := providers.NewHTTPProvider("", "", "") // Empty provider for demo
	registry := agent.NewAgentRegistry(cfg, prov)
	msgBus := bus.NewMessageBus()
	tm := team.NewTeamManager(registry, msgBus)

	// Define team configuration
	teamConfig := &team.TeamConfig{
		TeamID:  "dev-team-001",
		Name:    "Development Team",
		Pattern: "sequential",
		Roles: []team.RoleConfig{
			{
				Name:         "architect",
				Capabilities: []string{"design", "architecture"},
				Tools:        []string{"file_read", "diagram_*"},
			},
			{
				Name:         "developer",
				Capabilities: []string{"code", "implement"},
				Tools:        []string{"file_*", "git_*"},
			},
			{
				Name:         "tester",
				Capabilities: []string{"test", "qa"},
				Tools:        []string{"test_*", "file_read"},
			},
			{
				Name:         "reviewer",
				Capabilities: []string{"review", "approve"},
				Tools:        []string{"file_read", "git_*"},
			},
		},
		Coordinator: team.CoordinatorConfig{
			Role: "architect",
		},
		Settings: team.SettingsConfig{
			MaxDelegationDepth:      5,
			AgentTimeoutSeconds:     60,
			FailureThreshold:        3,
			ConsensusTimeoutSeconds: 30,
		},
	}

	// Create team
	ctx := context.Background()
	devTeam, err := tm.CreateTeam(ctx, teamConfig)
	if err != nil {
		log.Fatalf("Failed to create team: %v", err)
	}
	defer tm.DissolveTeam(ctx, devTeam.ID)

	fmt.Printf("✓ Created team: %s (ID: %s)\n", devTeam.Name, devTeam.ID)
	fmt.Printf("  Pattern: %v\n", devTeam.Pattern)
	fmt.Printf("  Agents: %d\n\n", len(devTeam.Agents))

	// Create coordinator
	router := team.NewDelegationRouter(5)
	coordinator := team.NewCoordinatorAgent(
		"coordinator",
		devTeam.ID,
		devTeam,
		team.PatternSequential,
		msgBus,
		router,
	, ctx)

	// Define sequential tasks
	tasks := []*team.Task{
		team.NewTask(
			"Design user authentication system",
			"architect",
			map[string]interface{}{
				"requirements": "OAuth2, JWT tokens, role-based access",
			},
		),
		team.NewTask(
			"Implement authentication endpoints",
			"developer",
			nil, // Will receive design from previous task
		),
		team.NewTask(
			"Write and run integration tests",
			"tester",
			nil, // Will receive implementation details
		),
		team.NewTask(
			"Review code and approve",
			"reviewer",
			nil, // Will receive test results
		),
	}

	fmt.Println("Executing sequential workflow...")
	fmt.Println("Tasks:")
	for i, task := range tasks {
		fmt.Printf("  %d. %s (role: %s)\n", i+1, task.Description, task.RequiredRole)
	}
	fmt.Println()

	// Execute sequential workflow
	results, err := coordinator.ExecuteSequential(ctx, tasks)
	if err != nil {
		log.Fatalf("Workflow failed: %v", err)
	}

	// Display results
	fmt.Println("✓ Workflow completed successfully!\n")
	fmt.Println("Results:")
	for i, result := range results {
		fmt.Printf("  Task %d: %s\n", i+1, result)
	}
	fmt.Println()

	// Show shared context
	fmt.Println("Shared Context:")
	contextData := devTeam.SharedContext.GetAll()
	for key, value := range contextData {
		fmt.Printf("  %s: %v\n", key, value)
	}
	fmt.Println()

	// Show metrics
	metrics := tm.GetMetrics().GetMetrics(devTeam.ID)
	fmt.Println("Metrics:")
	fmt.Printf("  Tasks delegated: %d\n", metrics.TaskDelegationCount)
	fmt.Printf("  Tasks completed: %d\n", metrics.TaskCompletionCount)
	fmt.Printf("  Average duration: %v\n", metrics.AverageTaskDuration)
	fmt.Println()

	fmt.Println("=== Example Complete ===")
}
