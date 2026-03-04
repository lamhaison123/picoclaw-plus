// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package teamcmd

import (
	"testing"
)

func TestNewTeamCommand(t *testing.T) {
	cmd := NewTeamCommand()

	if cmd == nil {
		t.Fatal("Expected command, got nil")
	}

	if cmd.Use != "team" {
		t.Errorf("Expected use 'team', got '%s'", cmd.Use)
	}

	// Check subcommands
	subcommands := cmd.Commands()
	expectedSubcommands := []string{"create", "list", "status", "dissolve", "memory"}

	if len(subcommands) != len(expectedSubcommands) {
		t.Errorf("Expected %d subcommands, got %d", len(expectedSubcommands), len(subcommands))
	}

	for _, expected := range expectedSubcommands {
		found := false
		for _, sub := range subcommands {
			if sub.Use == expected || sub.Use[:len(expected)] == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected subcommand '%s' not found", expected)
		}
	}
}
