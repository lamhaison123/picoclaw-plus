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

// HierarchicalWorkflowExample demonstrates a hierarchical analysis workflow
// where a lead analyst decomposes tasks and coordinates specialists
func main() {
	fmt.Println("=== Hierarchical Workflow Example ===\n")

	// Initialize dependencies
	cfg := &config.Config{}
	prov := providers.NewHTTPProvider("", "", "") // Empty provider for demo
	registry := agent.NewAgentRegistry(cfg, prov)
	msgBus := bus.NewMessageBus()
	tm := team.NewTeamManager(registry, msgBus)

	// Define team configuration
	teamConfig := &team.TeamConfig{
		TeamID:  "analysis-team-001",
		Name:    "Analysis Team",
		Pattern: "hierarchical",
		Roles: []team.RoleConfig{
			{
				Name:         "lead_analyst",
				Capabilities: []string{"coordinate", "plan", "synthesize"},
				Tools:        []string{"*"},
			},
			{
				Name:         "code_analyzer",
				Capabilities: []string{"code_analysis", "complexity"},
				Tools:        []string{"file_*", "ast_*"},
			},
			{
				Name:         "performance_analyzer",
				Capabilities: []string{"performance", "profiling"},
				Tools:        []string{"profile_*", "benchmark_*"},
			},
			{
				Name:         "security_analyzer",
				Capabilities: []string{"security", "vulnerability"},
				Tools:        []string{"security_*", "audit_*"},
			},
		},
		Coordinator: team.CoordinatorConfig{
			Role: "lead_analyst",
		},
		Settings: team.SettingsConfig{
			MaxDelegationDepth:      5,
			AgentTimeoutSeconds:     180,
			FailureThreshold:        3,
			ConsensusTimeoutSeconds: 30,
		},
	}

	// Create team
	ctx := context.Background()
	analysisTeam, err := tm.CreateTeam(ctx, teamConfig)
	if err != nil {
		log.Fatalf("Failed to create team: %v", err)
	}
	defer tm.DissolveTeam(ctx, analysisTeam.ID)

	fmt.Printf("✓ Created team: %s (ID: %s)\n", analysisTeam.Name, analysisTeam.ID)
	fmt.Printf("  Pattern: %v\n", analysisTeam.Pattern)
	fmt.Printf("  Agents: %d\n\n", len(analysisTeam.Agents))

	// Create coordinator
	router := team.NewDelegationRouter(5)
	coordinator := team.NewCoordinatorAgent(
		"lead_coordinator",
		analysisTeam.ID,
		analysisTeam,
		team.PatternHierarchical,
		msgBus,
		router,
	, ctx)

	// Define main task (will be decomposed)
	mainTask := team.NewTask(
		"Comprehensive codebase analysis",
		"lead_analyst",
		map[string]interface{}{
			"scope":  "entire codebase",
			"depth":  "detailed",
			"output": "full report",
		},
	)

	fmt.Println("Executing hierarchical workflow...")
	fmt.Printf("Main Task: %s\n", mainTask.Description)
	fmt.Println("(Will be decomposed into subtasks dynamically)\n")

	// Execute hierarchical workflow
	result, err := coordinator.ExecuteHierarchical(ctx, mainTask)
	if err != nil {
		log.Fatalf("Workflow failed: %v", err)
	}

	// Display result
	fmt.Println("✓ Analysis completed!\n")
	fmt.Println("Final Result:")
	fmt.Printf("  %s\n\n", result)

	// Show task decomposition history
	fmt.Println("Task Decomposition History:")
	history := analysisTeam.SharedContext.GetHistory()
	for i, entry := range history {
		fmt.Printf("  %d. [%s] %s: %v\n", i+1, entry.Timestamp.Format("15:04:05"), entry.Action, entry.Data)
	}
	fmt.Println()

	// Show shared context
	fmt.Println("Shared Context:")
	contextData := analysisTeam.SharedContext.GetAll()
	for key, value := range contextData {
		fmt.Printf("  %s: %v\n", key, value)
	}
	fmt.Println()

	// Show metrics
	metrics := tm.GetMetrics().GetMetrics(analysisTeam.ID)
	fmt.Println("Metrics:")
	fmt.Printf("  Tasks delegated: %d\n", metrics.TaskDelegationCount)
	fmt.Printf("  Tasks completed: %d\n", metrics.TaskCompletionCount)
	fmt.Printf("  Tasks failed: %d\n", metrics.TaskFailureCount)
	fmt.Printf("  Average duration: %v\n", metrics.AverageTaskDuration)
	fmt.Println()

	// Show agent statuses
	fmt.Println("Agent Statuses:")
	for agentID, agent := range analysisTeam.Agents {
		fmt.Printf("  %s (%s): %v\n", agentID, agent.Role, agent.Status)
	}
	fmt.Println()

	fmt.Println("=== Example Complete ===")
}
