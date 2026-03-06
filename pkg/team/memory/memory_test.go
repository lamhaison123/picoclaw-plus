// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package memory

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewTeamMemory(t *testing.T) {
	tempDir := t.TempDir()
	tm := NewTeamMemory(tempDir)

	if tm.workspace != tempDir {
		t.Errorf("Expected workspace %s, got %s", tempDir, tm.workspace)
	}

	expectedMemoryDir := filepath.Join(tempDir, "memory", "teams")
	if tm.memoryDir != expectedMemoryDir {
		t.Errorf("Expected memory dir %s, got %s", expectedMemoryDir, tm.memoryDir)
	}

	// Check directory was created
	if _, err := os.Stat(tm.memoryDir); os.IsNotExist(err) {
		t.Error("Expected memory directory to be created")
	}
}

func TestTeamMemory_SaveAndLoadRecord(t *testing.T) {
	tempDir := t.TempDir()
	tm := NewTeamMemory(tempDir)

	record := &TeamMemoryRecord{
		TeamID:    "team1",
		TeamName:  "Test Team",
		Pattern:   "sequential",
		StartTime: time.Now().Add(-1 * time.Hour),
		EndTime:   time.Now(),
		SharedContext: map[string]any{
			"key1": "value1",
			"key2": 123,
		},
		Tasks: []TaskRecord{
			{
				TaskID:      "task1",
				Description: "Test task",
				Role:        "developer",
				AgentID:     "agent1",
				Status:      "completed",
				Result:      "success",
				CreatedAt:   time.Now(),
			},
		},
		Outcome: "success",
	}

	// Save record
	err := tm.SaveTeamRecord(record)
	if err != nil {
		t.Fatalf("SaveTeamRecord failed: %v", err)
	}

	// Load record
	loaded, err := tm.LoadTeamRecord("team1")
	if err != nil {
		t.Fatalf("LoadTeamRecord failed: %v", err)
	}

	if loaded.TeamID != record.TeamID {
		t.Errorf("Expected team ID %s, got %s", record.TeamID, loaded.TeamID)
	}
	if loaded.TeamName != record.TeamName {
		t.Errorf("Expected team name %s, got %s", record.TeamName, loaded.TeamName)
	}
	if loaded.Outcome != record.Outcome {
		t.Errorf("Expected outcome %s, got %s", record.Outcome, loaded.Outcome)
	}
	if len(loaded.Tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(loaded.Tasks))
	}
}

func TestTeamMemory_ListTeamRecords(t *testing.T) {
	tempDir := t.TempDir()
	tm := NewTeamMemory(tempDir)

	// Save multiple records
	record1 := &TeamMemoryRecord{
		TeamID:        "team1",
		TeamName:      "Team 1",
		Pattern:       "sequential",
		StartTime:     time.Now(),
		EndTime:       time.Now(),
		SharedContext: map[string]any{},
		Outcome:       "success",
	}

	record2 := &TeamMemoryRecord{
		TeamID:        "team2",
		TeamName:      "Team 2",
		Pattern:       "parallel",
		StartTime:     time.Now(),
		EndTime:       time.Now(),
		SharedContext: map[string]any{},
		Outcome:       "success",
	}

	tm.SaveTeamRecord(record1)
	tm.SaveTeamRecord(record2)

	// List records
	teamIDs, err := tm.ListTeamRecords()
	if err != nil {
		t.Fatalf("ListTeamRecords failed: %v", err)
	}

	if len(teamIDs) != 2 {
		t.Errorf("Expected 2 team IDs, got %d", len(teamIDs))
	}

	// Check sorted order
	if teamIDs[0] != "team1" || teamIDs[1] != "team2" {
		t.Errorf("Expected sorted team IDs [team1, team2], got %v", teamIDs)
	}
}

func TestTeamMemory_LoadTeamRecord_NotFound(t *testing.T) {
	tempDir := t.TempDir()
	tm := NewTeamMemory(tempDir)

	_, err := tm.LoadTeamRecord("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent team, got nil")
	}
}
