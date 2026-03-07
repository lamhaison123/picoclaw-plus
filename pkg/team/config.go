// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package team

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LoadTeamConfig loads a team configuration from a JSON file
func LoadTeamConfig(configPath string) (*TeamConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config TeamConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config JSON: %w", err)
	}

	// Validate configuration
	if err := ValidateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// ValidateConfig validates a team configuration
func ValidateConfig(config *TeamConfig) error {
	if config.TeamID == "" {
		return fmt.Errorf("team_id is required")
	}

	if config.Name == "" {
		return fmt.Errorf("name is required")
	}

	if config.Pattern == "" {
		return fmt.Errorf("pattern is required")
	}

	validPatterns := map[string]bool{
		"sequential":   true,
		"parallel":     true,
		"hierarchical": true,
	}
	if !validPatterns[config.Pattern] {
		return fmt.Errorf("invalid pattern '%s', must be one of: sequential, parallel, hierarchical", config.Pattern)
	}

	if len(config.Roles) == 0 {
		return fmt.Errorf("at least one role is required")
	}

	// Validate roles
	roleNames := make(map[string]bool)
	for i, role := range config.Roles {
		if role.Name == "" {
			return fmt.Errorf("role %d: name is required", i)
		}

		if roleNames[role.Name] {
			return fmt.Errorf("duplicate role name: %s", role.Name)
		}
		roleNames[role.Name] = true

		if len(role.Capabilities) == 0 {
			return fmt.Errorf("role '%s': at least one capability is required", role.Name)
		}

		if len(role.Tools) == 0 {
			return fmt.Errorf("role '%s': at least one tool is required", role.Name)
		}
	}

	// Validate coordinator role exists
	if config.Coordinator.Role != "" {
		if !roleNames[config.Coordinator.Role] {
			return fmt.Errorf("coordinator role '%s' not defined in roles", config.Coordinator.Role)
		}
	}

	// Validate settings
	if config.Settings.MaxDelegationDepth < 1 {
		return fmt.Errorf("max_delegation_depth must be at least 1")
	}

	if config.Settings.AgentTimeoutSeconds < 1 {
		return fmt.Errorf("agent_timeout_seconds must be at least 1")
	}

	if config.Settings.FailureThreshold < 1 {
		return fmt.Errorf("failure_threshold must be at least 1")
	}

	if config.Settings.ConsensusTimeoutSeconds < 1 {
		return fmt.Errorf("consensus_timeout_seconds must be at least 1")
	}

	return nil
}

// SubstituteVariables replaces template variables in configuration
func SubstituteVariables(config *TeamConfig, variables map[string]string) error {
	// Validate variable names (alphanumeric + underscore only)
	for key := range variables {
		if !isValidVariableName(key) {
			return fmt.Errorf("invalid variable name '%s': must contain only alphanumeric characters and underscores", key)
		}
	}

	// Substitute in team name
	config.Name = substituteString(config.Name, variables)
	config.Description = substituteString(config.Description, variables)

	// Substitute in roles
	for i := range config.Roles {
		config.Roles[i].Name = substituteString(config.Roles[i].Name, variables)
		config.Roles[i].Description = substituteString(config.Roles[i].Description, variables)
		config.Roles[i].Model = substituteString(config.Roles[i].Model, variables)

		for j := range config.Roles[i].Capabilities {
			config.Roles[i].Capabilities[j] = substituteString(config.Roles[i].Capabilities[j], variables)
		}

		for j := range config.Roles[i].Tools {
			config.Roles[i].Tools[j] = substituteString(config.Roles[i].Tools[j], variables)
		}
	}

	// Substitute in coordinator
	config.Coordinator.Role = substituteString(config.Coordinator.Role, variables)
	config.Coordinator.AgentID = substituteString(config.Coordinator.AgentID, variables)

	return nil
}

// substituteString replaces ${VAR} or $VAR with values from variables map
func substituteString(s string, variables map[string]string) string {
	result := s

	// Replace ${VAR} format
	for key, value := range variables {
		result = strings.ReplaceAll(result, fmt.Sprintf("${%s}", key), value)
		result = strings.ReplaceAll(result, fmt.Sprintf("$%s", key), value)
	}

	// Replace environment variables
	result = os.ExpandEnv(result)

	return result
}

// isValidVariableName checks if a variable name is valid (alphanumeric + underscore)
func isValidVariableName(name string) bool {
	if name == "" {
		return false
	}
	for _, ch := range name {
		if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '_') {
			return false
		}
	}
	return true
}

// SaveTeamConfig saves a team configuration to a JSON file
func SaveTeamConfig(config *TeamConfig, configPath string) error {
	// Validate before saving
	if err := ValidateConfig(config); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Ensure directory exists
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(configPath, data, 0o644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// ListTeamConfigs lists all team configuration files in a directory
func ListTeamConfigs(configDir string) ([]string, error) {
	entries, err := os.ReadDir(configDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read config directory: %w", err)
	}

	configs := []string{}
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			configs = append(configs, filepath.Join(configDir, entry.Name()))
		}
	}

	return configs, nil
}
