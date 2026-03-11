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
	"time"

	"github.com/sipeed/picoclaw/pkg/agent"
	"github.com/sipeed/picoclaw/pkg/config"
	"github.com/sipeed/picoclaw/pkg/logger"
)

// TeamState represents the persisted state of an active team
type TeamState struct {
	ID            string                `json:"id"`
	Name          string                `json:"name"`
	Pattern       CollaborationPattern  `json:"pattern"`
	Status        TeamStatus            `json:"status"`
	CoordinatorID string                `json:"coordinator_id"`
	CreatedAt     time.Time             `json:"created_at"`
	Config        *TeamConfig           `json:"config"`
	SharedContext map[string]any        `json:"shared_context"`
	Agents        map[string]*TeamAgent `json:"agents"`
}

// SaveTeamState saves the current state of a team to disk
func (tm *TeamManager) SaveTeamState(team *Team) error {
	stateDir := filepath.Join(tm.getStateDir(), "active")
	if err := os.MkdirAll(stateDir, 0o755); err != nil {
		return fmt.Errorf("failed to create state directory: %w", err)
	}

	state := &TeamState{
		ID:            team.ID,
		Name:          team.Name,
		Pattern:       team.Pattern,
		Status:        team.Status,
		CoordinatorID: team.CoordinatorID,
		CreatedAt:     team.CreatedAt,
		Config:        team.Config,
		SharedContext: team.SharedContext.Snapshot(),
		Agents:        team.Agents,
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal team state: %w", err)
	}

	statePath := filepath.Join(stateDir, fmt.Sprintf("%s.json", team.ID))
	tempPath := statePath + ".tmp"

	logger.InfoCF("team", "Saving team state",
		map[string]any{
			"team_id":    team.ID,
			"state_path": statePath,
		})

	if err := os.WriteFile(tempPath, data, 0o644); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	if err := os.Rename(tempPath, statePath); err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	logger.InfoCF("team", "Team state saved successfully",
		map[string]any{
			"team_id":    team.ID,
			"state_path": statePath,
		})

	return nil
}

// LoadTeamState loads a team state from disk
func (tm *TeamManager) LoadTeamState(teamID string) (*TeamState, error) {
	stateDir := filepath.Join(tm.getStateDir(), "active")
	statePath := filepath.Join(stateDir, fmt.Sprintf("%s.json", teamID))

	data, err := os.ReadFile(statePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("team state not found: %s", teamID)
		}
		return nil, fmt.Errorf("failed to read team state: %w", err)
	}

	var state TeamState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to unmarshal team state: %w", err)
	}

	return &state, nil
}

// LoadAllTeamStates loads all active team states from disk
func (tm *TeamManager) LoadAllTeamStates() ([]*TeamState, error) {
	stateDir := filepath.Join(tm.getStateDir(), "active")

	logger.InfoCF("team", "Loading team states",
		map[string]any{
			"state_dir": stateDir,
		})

	if err := os.MkdirAll(stateDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create state directory: %w", err)
	}

	files, err := os.ReadDir(stateDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read state directory: %w", err)
	}

	logger.InfoCF("team", "Found state files",
		map[string]any{
			"count": len(files),
		})

	var states []*TeamState
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

		data, err := os.ReadFile(filepath.Join(stateDir, file.Name()))
		if err != nil {
			logger.WarnCF("team", "Failed to read state file",
				map[string]any{
					"file":  file.Name(),
					"error": err.Error(),
				})
			continue
		}

		var state TeamState
		if err := json.Unmarshal(data, &state); err != nil {
			logger.WarnCF("team", "Failed to unmarshal state file",
				map[string]any{
					"file":  file.Name(),
					"error": err.Error(),
				})
			continue
		}

		states = append(states, &state)
	}

	logger.InfoCF("team", "Loaded team states",
		map[string]any{
			"count": len(states),
		})

	return states, nil
}

// DeleteTeamState removes a team state file from disk
func (tm *TeamManager) DeleteTeamState(teamID string) error {
	stateDir := filepath.Join(tm.getStateDir(), "active")
	statePath := filepath.Join(stateDir, fmt.Sprintf("%s.json", teamID))

	if err := os.Remove(statePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete team state: %w", err)
	}

	return nil
}

// RestoreTeamFromState restores a team from saved state
func (tm *TeamManager) RestoreTeamFromState(state *TeamState) (*Team, error) {
	team := &Team{
		ID:            state.ID,
		Name:          state.Name,
		Pattern:       state.Pattern,
		Status:        state.Status,
		CoordinatorID: state.CoordinatorID,
		CreatedAt:     state.CreatedAt,
		Config:        state.Config,
		Agents:        state.Agents,
		SharedContext: NewSharedContext(state.ID),
	}

	// Restore shared context
	for key, value := range state.SharedContext {
		team.SharedContext.Set(key, value, "system")
	}

	// Restore role capabilities
	if state.Config != nil {
		for _, roleConfig := range state.Config.Roles {
			compositeKey := fmt.Sprintf("%s:%s", state.ID, roleConfig.Name)
			tm.roleCapabilities[compositeKey] = roleConfig.Capabilities
		}
	}

	// Register agents for each role (same as CreateTeam)
	if tm.provider != nil && tm.cfg != nil && state.Config != nil {
		for _, roleConfig := range state.Config.Roles {
			agentID := fmt.Sprintf("%s-%s", state.ID, roleConfig.Name)

			// Create agent config for this role
			agentCfg := &config.AgentConfig{
				ID:   agentID,
				Name: fmt.Sprintf("%s (%s)", state.Name, roleConfig.Name),
				Model: &config.AgentModelConfig{
					Primary: roleConfig.Model,
				},
				Workspace: fmt.Sprintf("%s/teams/%s/%s", tm.workspace, state.ID, roleConfig.Name),
			}

			// Create agent instance with role-specific model
			instance := agent.NewAgentInstance(agentCfg, &tm.cfg.Agents.Defaults, tm.cfg, tm.provider)

			// Register in agent registry
			tm.registry.RegisterTeamAgent(agentID, instance)

			logger.InfoCF("team", "Registered team agent on restore",
				map[string]any{
					"team_id":   state.ID,
					"agent_id":  agentID,
					"role":      roleConfig.Name,
					"model":     roleConfig.Model,
					"workspace": instance.Workspace,
				})
		}
	} else {
		logger.WarnCF("team", "Cannot register agents on restore: provider or config not set",
			map[string]any{
				"team_id":      state.ID,
				"has_provider": tm.provider != nil,
				"has_cfg":      tm.cfg != nil,
				"has_config":   state.Config != nil,
			})
	}

	// Initialize metrics for restored team
	tm.metrics.RecordTeamCreation(state.ID)

	return team, nil
}

// getStateDir returns the directory for storing team states
func (tm *TeamManager) getStateDir() string {
	if tm.workspace != "" {
		return filepath.Join(tm.workspace, "teams")
	}
	// v0.2.1: Use PICOCLAW_HOME if set
	if picoclawHome := os.Getenv("PICOCLAW_HOME"); picoclawHome != "" {
		return filepath.Join(picoclawHome, "teams")
	}
	// Fallback to HOME directory
	if home := os.Getenv("HOME"); home != "" {
		return filepath.Join(home, ".picoclaw", "teams")
	}
	// Last resort fallback
	return filepath.Join(".", ".picoclaw", "teams")
}

// loadPersistedTeams loads all persisted teams into memory
func (tm *TeamManager) loadPersistedTeams() error {
	states, err := tm.LoadAllTeamStates()
	if err != nil {
		return err
	}

	tm.mu.Lock()
	defer tm.mu.Unlock()

	for _, state := range states {
		team, err := tm.RestoreTeamFromState(state)
		if err != nil {
			logger.WarnCF("team", "Failed to restore team from state",
				map[string]any{
					"team_id": state.ID,
					"error":   err.Error(),
				})
			continue
		}

		tm.teams[team.ID] = team

		logger.InfoCF("team", "Restored team from disk",
			map[string]any{
				"team_id":   team.ID,
				"team_name": team.Name,
			})
	}

	return nil
}
