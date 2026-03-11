// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package tools

import (
	"context"
	"fmt"
)

// TeamManagerAdapter adapts team.TeamManager to tools.TeamManager interface
type TeamManagerAdapter struct {
	impl interface {
		ExecuteTask(ctx context.Context, teamID string, taskDescription string) (any, error)
		ExecuteTaskWithRole(ctx context.Context, teamID string, taskDescription string, role string) (any, error)
		GetTeamStatusAsInterface(teamID string) (interface{}, error)
		GetAllTeamsAsInterface() map[string]interface{}
	}
}

// NewTeamManagerAdapter creates a new adapter
func NewTeamManagerAdapter(impl interface {
	ExecuteTask(ctx context.Context, teamID string, taskDescription string) (any, error)
	ExecuteTaskWithRole(ctx context.Context, teamID string, taskDescription string, role string) (any, error)
	GetTeamStatusAsInterface(teamID string) (interface{}, error)
	GetAllTeamsAsInterface() map[string]interface{}
}) TeamManager {
	// BUG FIX: Validate implementation is not nil
	if impl == nil {
		panic("NewTeamManagerAdapter: impl cannot be nil")
	}
	return &TeamManagerAdapter{impl: impl}
}

func (a *TeamManagerAdapter) ExecuteTask(ctx context.Context, teamID string, taskDescription string) (any, error) {
	// BUG FIX: Check if impl is nil
	if a.impl == nil {
		return nil, fmt.Errorf("team manager implementation is nil")
	}
	return a.impl.ExecuteTask(ctx, teamID, taskDescription)
}

func (a *TeamManagerAdapter) ExecuteTaskWithRole(ctx context.Context, teamID string, taskDescription string, role string) (any, error) {
	// BUG FIX: Check if impl is nil
	if a.impl == nil {
		return nil, fmt.Errorf("team manager implementation is nil")
	}
	return a.impl.ExecuteTaskWithRole(ctx, teamID, taskDescription, role)
}

func (a *TeamManagerAdapter) GetTeamStatus(teamID string) (interface{}, error) {
	// BUG FIX: Check if impl is nil
	if a.impl == nil {
		return nil, fmt.Errorf("team manager implementation is nil")
	}
	return a.impl.GetTeamStatusAsInterface(teamID)
}

func (a *TeamManagerAdapter) GetAllTeams() map[string]interface{} {
	// BUG FIX: Check if impl is nil
	if a.impl == nil {
		return make(map[string]interface{})
	}
	return a.impl.GetAllTeamsAsInterface()
}
