// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package team

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/sipeed/picoclaw/pkg/team/memory"
)

func TestTeamPersistence(t *testing.T) {
	// Create temp workspace
	tmpDir := t.TempDir()

	// Use helper to create test team manager
	tm1 := createTestTeamManager()
	teamMemory := memory.NewTeamMemory(tmpDir)
	tm1.SetTeamMemory(teamMemory)

	// Create a team
	teamConfig := &TeamConfig{
		TeamID:      "test-team-001",
		Name:        "Test Team",
		Description: "Test team for persistence",
		Pattern:     "sequential",
		Roles: []RoleConfig{
			{
				Name:         "developer",
				Description:  "Developer role",
				Capabilities: []string{"coding"},
				Tools:        []string{"editCode"},
				Model:        "test-model",
				MaxTokens:    1000,
				Temperature:  0.7,
			},
		},
		Coordinator: CoordinatorConfig{
			Role: "developer",
		},
		Settings: SettingsConfig{
			MaxDelegationDepth:      3,
			AgentTimeoutSeconds:     60,
			FailureThreshold:        2,
			ConsensusTimeoutSeconds: 30,
			MaxSpawnedAgents:        5,
		},
	}

	ctx := context.Background()
	team, err := tm1.CreateTeam(ctx, teamConfig)
	if err != nil {
		t.Fatalf("Failed to create team: %v", err)
	}

	// Verify state file was created
	statePath := filepath.Join(tmpDir, "teams", "active", "test-team-001.json")
	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		t.Fatalf("State file was not created: %s", statePath)
	}

	// Create second team manager (simulating a new process)
	tm2 := createTestTeamManager()
	teamMemory2 := memory.NewTeamMemory(tmpDir)
	tm2.SetTeamMemory(teamMemory2)

	// Verify team was loaded
	loadedTeam, err := tm2.GetTeam("test-team-001")
	if err != nil {
		t.Fatalf("Failed to get loaded team: %v", err)
	}

	if loadedTeam.ID != team.ID {
		t.Errorf("Team ID mismatch: got %s, want %s", loadedTeam.ID, team.ID)
	}

	if loadedTeam.Name != team.Name {
		t.Errorf("Team name mismatch: got %s, want %s", loadedTeam.Name, team.Name)
	}

	if len(loadedTeam.Agents) != len(team.Agents) {
		t.Errorf("Agent count mismatch: got %d, want %d", len(loadedTeam.Agents), len(team.Agents))
	}

	// Dissolve team
	if err := tm2.DissolveTeam(ctx, "test-team-001"); err != nil {
		t.Fatalf("Failed to dissolve team: %v", err)
	}

	// Verify state file was deleted
	if _, err := os.Stat(statePath); !os.IsNotExist(err) {
		t.Errorf("State file was not deleted after dissolution")
	}
}
