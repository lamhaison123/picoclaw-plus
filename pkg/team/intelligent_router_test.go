// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package team

import (
	"context"
	"testing"
)

// MockExecutor for testing
type MockExecutor struct {
	executeFunc func(ctx context.Context, agentID string, task *Task) (any, error)
}

func (m *MockExecutor) Execute(ctx context.Context, agentID string, task *Task) (any, error) {
	if m.executeFunc != nil {
		return m.executeFunc(ctx, agentID, task)
	}
	return nil, nil
}

func TestIntelligentRouter_DetermineRole_ParsesMapResult(t *testing.T) {
	// Create mock executor that returns map structure (like DirectAgentExecutor)
	mockExecutor := &MockExecutor{
		executeFunc: func(ctx context.Context, agentID string, task *Task) (any, error) {
			// Simulate DirectAgentExecutor returning a map
			return map[string]any{
				"status":      "completed",
				"agent_id":    agentID,
				"task_id":     task.ID,
				"role":        "manager",
				"description": task.Description,
				"result":      "developer", // The actual LLM response
				"timestamp":   "2026-03-04T10:00:00Z",
			}, nil
		},
	}

	router := NewIntelligentRouter(mockExecutor)
	router.SetTeamID("test-team")

	roles := []RoleConfig{
		{Name: "developer", Description: "Writes code", Capabilities: []string{"coding"}},
		{Name: "tester", Description: "Tests code", Capabilities: []string{"testing"}},
	}

	role, err := router.DetermineRole(context.Background(), "write code for login", roles)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if role != "developer" {
		t.Errorf("Expected role 'developer', got '%s'", role)
	}
}

func TestIntelligentRouter_DetermineRole_UsesTeamID(t *testing.T) {
	var capturedAgentID string

	mockExecutor := &MockExecutor{
		executeFunc: func(ctx context.Context, agentID string, task *Task) (any, error) {
			capturedAgentID = agentID
			return map[string]any{
				"result": "developer",
			}, nil
		},
	}

	router := NewIntelligentRouter(mockExecutor)
	router.SetTeamID("my-team-123")

	roles := []RoleConfig{
		{Name: "developer", Description: "Writes code", Capabilities: []string{"coding"}},
	}

	_, err := router.DetermineRole(context.Background(), "write code", roles)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expectedAgentID := "my-team-123-manager"
	if capturedAgentID != expectedAgentID {
		t.Errorf("Expected agent ID '%s', got '%s'", expectedAgentID, capturedAgentID)
	}
}

func TestIntelligentRouter_DetermineRole_FallbackToFirstRole(t *testing.T) {
	mockExecutor := &MockExecutor{
		executeFunc: func(ctx context.Context, agentID string, task *Task) (any, error) {
			return map[string]any{
				"result": "unknown-role", // Invalid role
			}, nil
		},
	}

	router := NewIntelligentRouter(mockExecutor)

	roles := []RoleConfig{
		{Name: "developer", Description: "Writes code", Capabilities: []string{"coding"}},
		{Name: "tester", Description: "Tests code", Capabilities: []string{"testing"}},
	}

	role, err := router.DetermineRole(context.Background(), "do something", roles)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Should fallback to first role
	if role != "developer" {
		t.Errorf("Expected fallback to first role 'developer', got '%s'", role)
	}
}

func TestIntelligentRouter_DetermineRoleSimple(t *testing.T) {
	router := NewIntelligentRouter(nil)

	roles := []RoleConfig{
		{Name: "architect", Description: "Designs systems", Capabilities: []string{"design"}},
		{Name: "developer", Description: "Writes code", Capabilities: []string{"coding"}},
		{Name: "tester", Description: "Tests code", Capabilities: []string{"testing"}},
	}

	tests := []struct {
		task         string
		expectedRole string
	}{
		{"write code for authentication", "developer"},
		{"implement user login", "developer"},
		{"design system architecture", "architect"},
		{"test the login feature", "tester"},
		{"run tests and validate", "tester"},
	}

	for _, tt := range tests {
		t.Run(tt.task, func(t *testing.T) {
			role := router.DetermineRoleSimple(tt.task, roles)
			if role != tt.expectedRole {
				t.Errorf("For task '%s', expected role '%s', got '%s'", tt.task, tt.expectedRole, role)
			}
		})
	}
}
