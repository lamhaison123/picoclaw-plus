// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package team

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// Error definitions
var (
	ErrTaskDescriptionRequired = errors.New("task description is required")
	ErrTaskRoleRequired        = errors.New("task required role is required")
)

// CollaborationPattern defines how agents interact within a team
type CollaborationPattern string

const (
	PatternSequential   CollaborationPattern = "sequential"
	PatternParallel     CollaborationPattern = "parallel"
	PatternHierarchical CollaborationPattern = "hierarchical"
)

// AgentStatus represents the current state of an agent
type AgentStatus string

const (
	StatusIdle         AgentStatus = "idle"
	StatusWorking      AgentStatus = "working"
	StatusWaiting      AgentStatus = "waiting"
	StatusFailed       AgentStatus = "failed"
	StatusUnresponsive AgentStatus = "unresponsive"
)

// TaskStatus represents the current state of a task
type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusAssigned   TaskStatus = "assigned"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusFailed     TaskStatus = "failed"
	TaskStatusCancelled  TaskStatus = "cancelled"
)

// TeamStatus represents the current state of a team
type TeamStatus string

const (
	TeamStatusInitializing TeamStatus = "initializing"
	TeamStatusActive       TeamStatus = "active"
	TeamStatusPaused       TeamStatus = "paused"
	TeamStatusDissolved    TeamStatus = "dissolved"
)

// VotingRule defines how consensus is reached
type VotingRule string

const (
	VotingRuleMajority  VotingRule = "majority"
	VotingRuleUnanimous VotingRule = "unanimous"
	VotingRuleWeighted  VotingRule = "weighted"
)

// Team represents a collection of agents working together
type Team struct {
	ID            string
	Name          string
	Pattern       CollaborationPattern
	Agents        map[string]*TeamAgent
	CoordinatorID string
	SharedContext *SharedContext
	Status        TeamStatus
	CreatedAt     time.Time
	Config        *TeamConfig
}

// TeamAgent represents an agent within a team
type TeamAgent struct {
	AgentID      string
	Role         string
	Capabilities []string
	Status       AgentStatus
	FailureCount int
	LastActive   time.Time
}

// Task represents a unit of work assigned to an agent
type Task struct {
	ID              string
	Description     string
	RequiredRole    string
	Context         map[string]any
	ParentTaskID    string
	DelegationChain []string
	Status          TaskStatus
	Result          any
	Error           error
	AssignedAgentID string
	CreatedAt       time.Time
	StartedAt       *time.Time
	CompletedAt     *time.Time
}

// TeamConfig defines the configuration for a team
type TeamConfig struct {
	TeamID      string            `json:"team_id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Pattern     string            `json:"pattern"`
	Roles       []RoleConfig      `json:"roles"`
	Coordinator CoordinatorConfig `json:"coordinator"`
	Settings    SettingsConfig    `json:"settings"`
}

// RoleConfig defines a role within a team
type RoleConfig struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Capabilities []string `json:"capabilities"`
	Tools        []string `json:"tools"`
	Model        string   `json:"model"`
	MaxTokens    int      `json:"max_tokens,omitempty"`
	Temperature  float64  `json:"temperature,omitempty"`
}

// CoordinatorConfig defines the coordinator agent configuration
type CoordinatorConfig struct {
	Role    string `json:"role"`
	AgentID string `json:"agent_id"`
}

// SettingsConfig defines team-wide settings
type SettingsConfig struct {
	MaxDelegationDepth      int `json:"max_delegation_depth"`
	AgentTimeoutSeconds     int `json:"agent_timeout_seconds"`
	FailureThreshold        int `json:"failure_threshold"`
	ConsensusTimeoutSeconds int `json:"consensus_timeout_seconds"`
	MaxSpawnedAgents        int `json:"max_spawned_agents,omitempty"`
}

// Helper functions

// generateTaskID generates a unique task ID
func generateTaskID() string {
	return "task_" + time.Now().Format("20060102150405") + "_" + randomString(8)
}

// randomString generates a random string of the given length
func randomString(length int) string {
	// Validate length
	if length <= 0 {
		return ""
	}

	bytes := make([]byte, length/2+1)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp-based ID if random fails
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	encoded := hex.EncodeToString(bytes)
	if len(encoded) < length {
		return encoded
	}
	return encoded[:length]
}

// Task methods

// NewTask creates a new task with the given description and required role
func NewTask(description, requiredRole string, context map[string]any) *Task {
	return &Task{
		ID:              generateTaskID(),
		Description:     description,
		RequiredRole:    requiredRole,
		Context:         context,
		DelegationChain: []string{},
		Status:          TaskStatusPending,
		CreatedAt:       time.Now(),
	}
}

// AddToDelegationChain adds an agent to the delegation chain
func (t *Task) AddToDelegationChain(agentID string) {
	t.DelegationChain = append(t.DelegationChain, agentID)
}

// IsInDelegationChain checks if an agent is in the delegation chain
func (t *Task) IsInDelegationChain(agentID string) bool {
	for _, id := range t.DelegationChain {
		if id == agentID {
			return true
		}
	}
	return false
}

// MarkAssigned marks the task as assigned to an agent
func (t *Task) MarkAssigned(agentID string) {
	t.Status = TaskStatusAssigned
	t.AssignedAgentID = agentID
}

// MarkInProgress marks the task as in progress
func (t *Task) MarkInProgress() {
	t.Status = TaskStatusInProgress
	now := time.Now()
	t.StartedAt = &now
}

// MarkCompleted marks the task as completed with a result
func (t *Task) MarkCompleted(result any) {
	t.Status = TaskStatusCompleted
	t.Result = result
	now := time.Now()
	t.CompletedAt = &now
}

// MarkFailed marks the task as failed with an error
func (t *Task) MarkFailed(err error) {
	t.Status = TaskStatusFailed
	t.Error = err
	now := time.Now()
	t.CompletedAt = &now
}

// MarkCancelled marks the task as cancelled
func (t *Task) MarkCancelled() {
	t.Status = TaskStatusCancelled
	now := time.Now()
	t.CompletedAt = &now
}

// Validate checks if the task has all required fields
func (t *Task) Validate() error {
	if t.Description == "" {
		return ErrTaskDescriptionRequired
	}
	if t.RequiredRole == "" {
		return ErrTaskRoleRequired
	}
	return nil
}
