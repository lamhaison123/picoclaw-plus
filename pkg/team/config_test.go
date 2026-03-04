// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package team

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateConfig_Valid(t *testing.T) {
	config := &TeamConfig{
		TeamID:  "team1",
		Name:    "Test Team",
		Pattern: "sequential",
		Roles: []RoleConfig{
			{
				Name:         "developer",
				Capabilities: []string{"code"},
				Tools:        []string{"file_*"},
			},
		},
		Coordinator: CoordinatorConfig{
			Role: "developer",
		},
		Settings: SettingsConfig{
			MaxDelegationDepth:      5,
			AgentTimeoutSeconds:     30,
			FailureThreshold:        3,
			ConsensusTimeoutSeconds: 60,
		},
	}

	err := ValidateConfig(config)
	if err != nil {
		t.Errorf("Expected valid config, got error: %v", err)
	}
}

func TestValidateConfig_MissingTeamID(t *testing.T) {
	config := &TeamConfig{
		Name:    "Test Team",
		Pattern: "sequential",
		Roles: []RoleConfig{
			{Name: "developer", Capabilities: []string{"code"}, Tools: []string{"file_*"}},
		},
		Settings: SettingsConfig{
			MaxDelegationDepth:      5,
			AgentTimeoutSeconds:     30,
			FailureThreshold:        3,
			ConsensusTimeoutSeconds: 60,
		},
	}

	err := ValidateConfig(config)
	if err == nil {
		t.Error("Expected error for missing team_id, got nil")
	}
}

func TestValidateConfig_InvalidPattern(t *testing.T) {
	config := &TeamConfig{
		TeamID:  "team1",
		Name:    "Test Team",
		Pattern: "invalid",
		Roles: []RoleConfig{
			{Name: "developer", Capabilities: []string{"code"}, Tools: []string{"file_*"}},
		},
		Settings: SettingsConfig{
			MaxDelegationDepth:      5,
			AgentTimeoutSeconds:     30,
			FailureThreshold:        3,
			ConsensusTimeoutSeconds: 60,
		},
	}

	err := ValidateConfig(config)
	if err == nil {
		t.Error("Expected error for invalid pattern, got nil")
	}
}

func TestValidateConfig_DuplicateRole(t *testing.T) {
	config := &TeamConfig{
		TeamID:  "team1",
		Name:    "Test Team",
		Pattern: "sequential",
		Roles: []RoleConfig{
			{Name: "developer", Capabilities: []string{"code"}, Tools: []string{"file_*"}},
			{Name: "developer", Capabilities: []string{"test"}, Tools: []string{"test_*"}},
		},
		Settings: SettingsConfig{
			MaxDelegationDepth:      5,
			AgentTimeoutSeconds:     30,
			FailureThreshold:        3,
			ConsensusTimeoutSeconds: 60,
		},
	}

	err := ValidateConfig(config)
	if err == nil {
		t.Error("Expected error for duplicate role, got nil")
	}
}

func TestValidateConfig_InvalidCoordinatorRole(t *testing.T) {
	config := &TeamConfig{
		TeamID:  "team1",
		Name:    "Test Team",
		Pattern: "sequential",
		Roles: []RoleConfig{
			{Name: "developer", Capabilities: []string{"code"}, Tools: []string{"file_*"}},
		},
		Coordinator: CoordinatorConfig{
			Role: "nonexistent",
		},
		Settings: SettingsConfig{
			MaxDelegationDepth:      5,
			AgentTimeoutSeconds:     30,
			FailureThreshold:        3,
			ConsensusTimeoutSeconds: 60,
		},
	}

	err := ValidateConfig(config)
	if err == nil {
		t.Error("Expected error for invalid coordinator role, got nil")
	}
}

func TestSubstituteVariables(t *testing.T) {
	config := &TeamConfig{
		TeamID:      "team1",
		Name:        "${TEAM_NAME}",
		Description: "Team for $PROJECT",
		Pattern:     "sequential",
		Roles: []RoleConfig{
			{
				Name:         "${ROLE_NAME}",
				Description:  "Role for $PROJECT",
				Capabilities: []string{"${CAPABILITY}"},
				Tools:        []string{"$TOOL"},
				Model:        "${MODEL}",
			},
		},
		Settings: SettingsConfig{
			MaxDelegationDepth:      5,
			AgentTimeoutSeconds:     30,
			FailureThreshold:        3,
			ConsensusTimeoutSeconds: 60,
		},
	}

	variables := map[string]string{
		"TEAM_NAME":  "My Team",
		"PROJECT":    "MyProject",
		"ROLE_NAME":  "developer",
		"CAPABILITY": "code",
		"TOOL":       "file_*",
		"MODEL":      "gpt-4",
	}

	err := SubstituteVariables(config, variables)
	if err != nil {
		t.Fatalf("SubstituteVariables failed: %v", err)
	}

	if config.Name != "My Team" {
		t.Errorf("Expected name 'My Team', got '%s'", config.Name)
	}
	if config.Description != "Team for MyProject" {
		t.Errorf("Expected description 'Team for MyProject', got '%s'", config.Description)
	}
	if config.Roles[0].Name != "developer" {
		t.Errorf("Expected role name 'developer', got '%s'", config.Roles[0].Name)
	}
	if config.Roles[0].Model != "gpt-4" {
		t.Errorf("Expected model 'gpt-4', got '%s'", config.Roles[0].Model)
	}
}

func TestSaveAndLoadTeamConfig(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "team.json")

	config := &TeamConfig{
		TeamID:  "team1",
		Name:    "Test Team",
		Pattern: "sequential",
		Roles: []RoleConfig{
			{
				Name:         "developer",
				Capabilities: []string{"code"},
				Tools:        []string{"file_*"},
			},
		},
		Settings: SettingsConfig{
			MaxDelegationDepth:      5,
			AgentTimeoutSeconds:     30,
			FailureThreshold:        3,
			ConsensusTimeoutSeconds: 60,
		},
	}

	// Save config
	err := SaveTeamConfig(config, configPath)
	if err != nil {
		t.Fatalf("SaveTeamConfig failed: %v", err)
	}

	// Check file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Expected config file to be created")
	}

	// Load config
	loaded, err := LoadTeamConfig(configPath)
	if err != nil {
		t.Fatalf("LoadTeamConfig failed: %v", err)
	}

	if loaded.TeamID != config.TeamID {
		t.Errorf("Expected team ID %s, got %s", config.TeamID, loaded.TeamID)
	}
	if loaded.Name != config.Name {
		t.Errorf("Expected name %s, got %s", config.Name, loaded.Name)
	}
	if len(loaded.Roles) != len(config.Roles) {
		t.Errorf("Expected %d roles, got %d", len(config.Roles), len(loaded.Roles))
	}
}

func TestLoadTeamConfig_InvalidJSON(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "invalid.json")

	// Write invalid JSON
	os.WriteFile(configPath, []byte("invalid json"), 0o644)

	_, err := LoadTeamConfig(configPath)
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestLoadTeamConfig_FileNotFound(t *testing.T) {
	_, err := LoadTeamConfig("/nonexistent/config.json")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}

func TestListTeamConfigs(t *testing.T) {
	tempDir := t.TempDir()

	// Create test config files
	config1 := &TeamConfig{
		TeamID:  "team1",
		Name:    "Team 1",
		Pattern: "sequential",
		Roles: []RoleConfig{
			{Name: "developer", Capabilities: []string{"code"}, Tools: []string{"file_*"}},
		},
		Settings: SettingsConfig{
			MaxDelegationDepth:      5,
			AgentTimeoutSeconds:     30,
			FailureThreshold:        3,
			ConsensusTimeoutSeconds: 60,
		},
	}

	config2 := &TeamConfig{
		TeamID:  "team2",
		Name:    "Team 2",
		Pattern: "parallel",
		Roles: []RoleConfig{
			{Name: "tester", Capabilities: []string{"test"}, Tools: []string{"test_*"}},
		},
		Settings: SettingsConfig{
			MaxDelegationDepth:      5,
			AgentTimeoutSeconds:     30,
			FailureThreshold:        3,
			ConsensusTimeoutSeconds: 60,
		},
	}

	SaveTeamConfig(config1, filepath.Join(tempDir, "team1.json"))
	SaveTeamConfig(config2, filepath.Join(tempDir, "team2.json"))

	// Create non-JSON file
	os.WriteFile(filepath.Join(tempDir, "readme.txt"), []byte("test"), 0o644)

	// List configs
	configs, err := ListTeamConfigs(tempDir)
	if err != nil {
		t.Fatalf("ListTeamConfigs failed: %v", err)
	}

	if len(configs) != 2 {
		t.Errorf("Expected 2 config files, got %d", len(configs))
	}
}

func TestListTeamConfigs_EmptyDirectory(t *testing.T) {
	tempDir := t.TempDir()

	configs, err := ListTeamConfigs(tempDir)
	if err != nil {
		t.Fatalf("ListTeamConfigs failed: %v", err)
	}

	if len(configs) != 0 {
		t.Errorf("Expected 0 config files, got %d", len(configs))
	}
}

func TestListTeamConfigs_NonExistentDirectory(t *testing.T) {
	configs, err := ListTeamConfigs("/nonexistent/directory")
	if err != nil {
		t.Fatalf("Expected no error for non-existent directory, got %v", err)
	}

	if len(configs) != 0 {
		t.Errorf("Expected 0 config files, got %d", len(configs))
	}
}

func TestSubstituteVariables_ValidNames(t *testing.T) {
	config := &TeamConfig{
		TeamID:      "team1",
		Name:        "${TEAM_NAME}",
		Description: "$PROJECT_DESC",
		Pattern:     "sequential",
		Roles: []RoleConfig{
			{
				Name:         "${ROLE_NAME}",
				Description:  "Developer role",
				Capabilities: []string{"code"},
				Tools:        []string{"file_*"},
			},
		},
		Settings: SettingsConfig{
			MaxDelegationDepth:      5,
			AgentTimeoutSeconds:     30,
			FailureThreshold:        3,
			ConsensusTimeoutSeconds: 60,
		},
	}

	variables := map[string]string{
		"TEAM_NAME":    "My Team",
		"PROJECT_DESC": "Test Project",
		"ROLE_NAME":    "developer",
	}

	err := SubstituteVariables(config, variables)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if config.Name != "My Team" {
		t.Errorf("Expected name 'My Team', got '%s'", config.Name)
	}

	if config.Description != "Test Project" {
		t.Errorf("Expected description 'Test Project', got '%s'", config.Description)
	}

	if config.Roles[0].Name != "developer" {
		t.Errorf("Expected role name 'developer', got '%s'", config.Roles[0].Name)
	}
}

func TestSubstituteVariables_InvalidNames(t *testing.T) {
	config := &TeamConfig{
		TeamID:  "team1",
		Name:    "Test Team",
		Pattern: "sequential",
		Roles: []RoleConfig{
			{
				Name:         "developer",
				Capabilities: []string{"code"},
				Tools:        []string{"file_*"},
			},
		},
		Settings: SettingsConfig{
			MaxDelegationDepth:      5,
			AgentTimeoutSeconds:     30,
			FailureThreshold:        3,
			ConsensusTimeoutSeconds: 60,
		},
	}

	testCases := []struct {
		name      string
		variables map[string]string
		wantError bool
	}{
		{
			name: "invalid dash",
			variables: map[string]string{
				"TEAM-NAME": "My Team",
			},
			wantError: true,
		},
		{
			name: "invalid space",
			variables: map[string]string{
				"TEAM NAME": "My Team",
			},
			wantError: true,
		},
		{
			name: "invalid special char",
			variables: map[string]string{
				"TEAM@NAME": "My Team",
			},
			wantError: true,
		},
		{
			name: "empty name",
			variables: map[string]string{
				"": "My Team",
			},
			wantError: true,
		},
		{
			name: "valid underscore",
			variables: map[string]string{
				"TEAM_NAME": "My Team",
			},
			wantError: false,
		},
		{
			name: "valid alphanumeric",
			variables: map[string]string{
				"TEAM123": "My Team",
			},
			wantError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := SubstituteVariables(config, tc.variables)
			if tc.wantError && err == nil {
				t.Errorf("Expected error for invalid variable name, got nil")
			}
			if !tc.wantError && err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}
		})
	}
}
