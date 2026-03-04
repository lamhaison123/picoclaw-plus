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

// ParallelWorkflowExample demonstrates a parallel research workflow
// where multiple researchers work simultaneously on different aspects
func main() {
	fmt.Println("=== Parallel Workflow Example ===\n")

	// Initialize dependencies
	cfg := &config.Config{}
	prov := providers.NewHTTPProvider("", "", "") // Empty provider for demo
	registry := agent.NewAgentRegistry(cfg, prov)
	msgBus := bus.NewMessageBus()
	tm := team.NewTeamManager(registry, msgBus)

	// Define team configuration
	teamConfig := &team.TeamConfig{
		TeamID:  "research-team-001",
		Name:    "Research Team",
		Pattern: "parallel",
		Roles: []team.RoleConfig{
			{
				Name:         "coordinator",
				Capabilities: []string{"coordinate", "synthesize"},
				Tools:        []string{"*"},
			},
			{
				Name:         "doc_researcher",
				Capabilities: []string{"search", "documentation"},
				Tools:        []string{"web_*", "file_read"},
			},
			{
				Name:         "code_analyst",
				Capabilities: []string{"analyze", "code_review"},
				Tools:        []string{"file_*", "git_*"},
			},
			{
				Name:         "security_analyst",
				Capabilities: []string{"security", "vulnerability"},
				Tools:        []string{"security_*", "file_read"},
			},
		},
		Coordinator: team.CoordinatorConfig{
			Role: "coordinator",
		},
		Settings: team.SettingsConfig{
			MaxDelegationDepth:      5,
			AgentTimeoutSeconds:     120,
			FailureThreshold:        3,
			ConsensusTimeoutSeconds: 30,
		},
	}

	// Create team
	ctx := context.Background()
	researchTeam, err := tm.CreateTeam(ctx, teamConfig)
	if err != nil {
		log.Fatalf("Failed to create team: %v", err)
	}
	defer tm.DissolveTeam(ctx, researchTeam.ID)

	fmt.Printf("✓ Created team: %s (ID: %s)\n", researchTeam.Name, researchTeam.ID)
	fmt.Printf("  Pattern: %v\n", researchTeam.Pattern)
	fmt.Printf("  Agents: %d\n\n", len(researchTeam.Agents))

	// Create coordinator
	router := team.NewDelegationRouter(5)
	coordinator := team.NewCoordinatorAgent(
		"coordinator",
		researchTeam.ID,
		researchTeam,
		team.PatternParallel,
		msgBus,
		router,
	, ctx)

	// Define parallel tasks
	tasks := []*team.Task{
		team.NewTask(
			"Research authentication best practices",
			"doc_researcher",
			map[string]interface{}{
				"topic":   "OAuth2 and JWT",
				"sources": []string{"official docs", "security guides"},
			},
		),
		team.NewTask(
			"Analyze existing authentication code",
			"code_analyst",
			map[string]interface{}{
				"files": []string{"auth/*.go", "middleware/*.go"},
			},
		),
		team.NewTask(
			"Identify security vulnerabilities",
			"security_analyst",
			map[string]interface{}{
				"focus": []string{"injection", "XSS", "CSRF"},
			},
		),
	}

	fmt.Println("Executing parallel workflow...")
	fmt.Println("Tasks (will run simultaneously):")
	for i, task := range tasks {
		fmt.Printf("  %d. %s (role: %s)\n", i+1, task.Description, task.RequiredRole)
	}
	fmt.Println()

	// Execute parallel workflow
	results, err := coordinator.ExecuteParallel(ctx, tasks)
	if err != nil {
		log.Fatalf("Workflow failed: %v", err)
	}

	// Display results
	fmt.Println("✓ All tasks completed!\n")
	fmt.Println("Results:")
	for i, result := range results {
		fmt.Printf("  Task %d: %s\n", i+1, result)
	}
	fmt.Println()

	// Aggregate results
	agentIDs := []string{
		"research-team-001-doc_researcher",
		"research-team-001-code_analyst",
		"research-team-001-security_analyst",
	}

	aggregated, err := coordinator.AggregateResults(results, agentIDs, team.StrategyMerge)
	if err != nil {
		log.Printf("Warning: Failed to aggregate results: %v", err)
	} else {
		fmt.Println("Aggregated Results:")
		fmt.Printf("  %s\n\n", aggregated)
	}

	// Show shared context
	fmt.Println("Shared Context:")
	contextData := researchTeam.SharedContext.GetAll()
	for key, value := range contextData {
		fmt.Printf("  %s: %v\n", key, value)
	}
	fmt.Println()

	// Show metrics
	metrics := tm.GetMetrics().GetMetrics(researchTeam.ID)
	fmt.Println("Metrics:")
	fmt.Printf("  Tasks delegated: %d\n", metrics.TaskDelegationCount)
	fmt.Printf("  Tasks completed: %d\n", metrics.TaskCompletionCount)
	fmt.Printf("  Tasks failed: %d\n", metrics.TaskFailureCount)
	fmt.Printf("  Average duration: %v\n", metrics.AverageTaskDuration)
	fmt.Println()

	fmt.Println("=== Example Complete ===")
}
