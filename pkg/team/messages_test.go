// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package team

import (
	"errors"
	"testing"
	"time"
)

func TestNewTask(t *testing.T) {
	description := "Test task"
	role := "developer"
	context := map[string]any{"key": "value"}

	task := NewTask(description, role, context)

	if task.Description != description {
		t.Errorf("Expected description %s, got %s", description, task.Description)
	}
	if task.RequiredRole != role {
		t.Errorf("Expected role %s, got %s", role, task.RequiredRole)
	}
	if task.Status != TaskStatusPending {
		t.Errorf("Expected status %s, got %s", TaskStatusPending, task.Status)
	}
	if len(task.DelegationChain) != 0 {
		t.Errorf("Expected empty delegation chain, got %d items", len(task.DelegationChain))
	}
	if task.ID == "" {
		t.Error("Expected non-empty task ID")
	}
}

func TestTaskValidation(t *testing.T) {
	tests := []struct {
		name        string
		description string
		role        string
		expectError bool
	}{
		{"Valid task", "Test task", "developer", false},
		{"Missing description", "", "developer", true},
		{"Missing role", "Test task", "", true},
		{"Missing both", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := NewTask(tt.description, tt.role, nil)
			err := task.Validate()
			if tt.expectError && err == nil {
				t.Error("Expected validation error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		})
	}
}

func TestTaskDelegationChain(t *testing.T) {
	task := NewTask("Test task", "developer", nil)

	// Test adding to chain
	task.AddToDelegationChain("agent1")
	task.AddToDelegationChain("agent2")

	if len(task.DelegationChain) != 2 {
		t.Errorf("Expected 2 agents in chain, got %d", len(task.DelegationChain))
	}

	// Test checking if agent is in chain
	if !task.IsInDelegationChain("agent1") {
		t.Error("Expected agent1 to be in chain")
	}
	if !task.IsInDelegationChain("agent2") {
		t.Error("Expected agent2 to be in chain")
	}
	if task.IsInDelegationChain("agent3") {
		t.Error("Expected agent3 not to be in chain")
	}
}

func TestTaskStatusTransitions(t *testing.T) {
	task := NewTask("Test task", "developer", nil)

	// Test assigned
	task.MarkAssigned("agent1")
	if task.Status != TaskStatusAssigned {
		t.Errorf("Expected status %s, got %s", TaskStatusAssigned, task.Status)
	}
	if task.AssignedAgentID != "agent1" {
		t.Errorf("Expected assigned agent agent1, got %s", task.AssignedAgentID)
	}

	// Test in progress
	task.MarkInProgress()
	if task.Status != TaskStatusInProgress {
		t.Errorf("Expected status %s, got %s", TaskStatusInProgress, task.Status)
	}
	if task.StartedAt == nil {
		t.Error("Expected StartedAt to be set")
	}

	// Test completed
	result := "test result"
	task.MarkCompleted(result)
	if task.Status != TaskStatusCompleted {
		t.Errorf("Expected status %s, got %s", TaskStatusCompleted, task.Status)
	}
	if task.Result != result {
		t.Errorf("Expected result %v, got %v", result, task.Result)
	}
	if task.CompletedAt == nil {
		t.Error("Expected CompletedAt to be set")
	}
}

func TestTaskFailure(t *testing.T) {
	task := NewTask("Test task", "developer", nil)
	testError := errors.New("test error")

	task.MarkFailed(testError)

	if task.Status != TaskStatusFailed {
		t.Errorf("Expected status %s, got %s", TaskStatusFailed, task.Status)
	}
	if task.Error != testError {
		t.Errorf("Expected error %v, got %v", testError, task.Error)
	}
	if task.CompletedAt == nil {
		t.Error("Expected CompletedAt to be set")
	}
}

func TestTaskCancellation(t *testing.T) {
	task := NewTask("Test task", "developer", nil)

	task.MarkCancelled()

	if task.Status != TaskStatusCancelled {
		t.Errorf("Expected status %s, got %s", TaskStatusCancelled, task.Status)
	}
	if task.CompletedAt == nil {
		t.Error("Expected CompletedAt to be set")
	}
}

func TestTaskDelegationMessage(t *testing.T) {
	task := NewTask("Test task", "developer", nil)
	context := map[string]any{"key": "value"}

	msg := NewTaskDelegationMessage("team1", "agent1", "agent2", task, context)

	if msg.TeamID != "team1" {
		t.Errorf("Expected team ID team1, got %s", msg.TeamID)
	}
	if msg.FromAgentID != "agent1" {
		t.Errorf("Expected from agent agent1, got %s", msg.FromAgentID)
	}
	if msg.ToAgentID != "agent2" {
		t.Errorf("Expected to agent agent2, got %s", msg.ToAgentID)
	}
	if msg.Task != task {
		t.Error("Expected task to match")
	}
	if msg.MessageID == "" {
		t.Error("Expected non-empty message ID")
	}
}

func TestTaskDelegationMessageSerialization(t *testing.T) {
	task := NewTask("Test task", "developer", nil)
	msg := NewTaskDelegationMessage("team1", "agent1", "agent2", task, nil)

	// Serialize
	data, err := msg.ToJSON()
	if err != nil {
		t.Fatalf("Failed to serialize message: %v", err)
	}

	// Deserialize
	decoded, err := TaskDelegationMessageFromJSON(data)
	if err != nil {
		t.Fatalf("Failed to deserialize message: %v", err)
	}

	if decoded.MessageID != msg.MessageID {
		t.Errorf("Expected message ID %s, got %s", msg.MessageID, decoded.MessageID)
	}
	if decoded.TeamID != msg.TeamID {
		t.Errorf("Expected team ID %s, got %s", msg.TeamID, decoded.TeamID)
	}
	if decoded.Task.ID != task.ID {
		t.Errorf("Expected task ID %s, got %s", task.ID, decoded.Task.ID)
	}
}

func TestTaskResultMessage(t *testing.T) {
	result := "test result"
	testError := errors.New("test error")

	// Test success message
	msg := NewTaskResultMessage("team1", "agent1", "agent2", "task1", TaskStatusCompleted, result, nil)

	if msg.TeamID != "team1" {
		t.Errorf("Expected team ID team1, got %s", msg.TeamID)
	}
	if msg.TaskID != "task1" {
		t.Errorf("Expected task ID task1, got %s", msg.TaskID)
	}
	if msg.Status != TaskStatusCompleted {
		t.Errorf("Expected status %s, got %s", TaskStatusCompleted, msg.Status)
	}
	if msg.Result != result {
		t.Errorf("Expected result %v, got %v", result, msg.Result)
	}
	if msg.Error != "" {
		t.Errorf("Expected no error, got %s", msg.Error)
	}

	// Test error message
	msgErr := NewTaskResultMessage("team1", "agent1", "agent2", "task1", TaskStatusFailed, nil, testError)
	if msgErr.Error != testError.Error() {
		t.Errorf("Expected error %s, got %s", testError.Error(), msgErr.Error)
	}
}

func TestTaskResultMessageSerialization(t *testing.T) {
	msg := NewTaskResultMessage("team1", "agent1", "agent2", "task1", TaskStatusCompleted, "result", nil)

	// Serialize
	data, err := msg.ToJSON()
	if err != nil {
		t.Fatalf("Failed to serialize message: %v", err)
	}

	// Deserialize
	decoded, err := TaskResultMessageFromJSON(data)
	if err != nil {
		t.Fatalf("Failed to deserialize message: %v", err)
	}

	if decoded.MessageID != msg.MessageID {
		t.Errorf("Expected message ID %s, got %s", msg.MessageID, decoded.MessageID)
	}
	if decoded.TaskID != msg.TaskID {
		t.Errorf("Expected task ID %s, got %s", msg.TaskID, decoded.TaskID)
	}
}

func TestConsensusRequestMessage(t *testing.T) {
	options := []string{"option1", "option2", "option3"}
	context := map[string]any{"key": "value"}

	msg := NewConsensusRequestMessage("team1", "agent1", "consensus1", "Test question?", options, VotingRuleMajority, 30, context)

	if msg.TeamID != "team1" {
		t.Errorf("Expected team ID team1, got %s", msg.TeamID)
	}
	if msg.ConsensusID != "consensus1" {
		t.Errorf("Expected consensus ID consensus1, got %s", msg.ConsensusID)
	}
	if msg.Question != "Test question?" {
		t.Errorf("Expected question 'Test question?', got %s", msg.Question)
	}
	if len(msg.Options) != 3 {
		t.Errorf("Expected 3 options, got %d", len(msg.Options))
	}
	if msg.VotingRule != VotingRuleMajority {
		t.Errorf("Expected voting rule %s, got %s", VotingRuleMajority, msg.VotingRule)
	}
	if msg.Timeout != 30 {
		t.Errorf("Expected timeout 30, got %d", msg.Timeout)
	}
}

func TestConsensusRequestMessageSerialization(t *testing.T) {
	options := []string{"option1", "option2"}
	msg := NewConsensusRequestMessage("team1", "agent1", "consensus1", "Test?", options, VotingRuleMajority, 30, nil)

	// Serialize
	data, err := msg.ToJSON()
	if err != nil {
		t.Fatalf("Failed to serialize message: %v", err)
	}

	// Deserialize
	decoded, err := ConsensusRequestMessageFromJSON(data)
	if err != nil {
		t.Fatalf("Failed to deserialize message: %v", err)
	}

	if decoded.ConsensusID != msg.ConsensusID {
		t.Errorf("Expected consensus ID %s, got %s", msg.ConsensusID, decoded.ConsensusID)
	}
	if len(decoded.Options) != len(msg.Options) {
		t.Errorf("Expected %d options, got %d", len(msg.Options), len(decoded.Options))
	}
}

func TestConsensusVoteMessage(t *testing.T) {
	context := map[string]any{"key": "value"}

	msg := NewConsensusVoteMessage("team1", "agent1", "agent2", "consensus1", "option1", 1.0, "My rationale", context)

	if msg.TeamID != "team1" {
		t.Errorf("Expected team ID team1, got %s", msg.TeamID)
	}
	if msg.ConsensusID != "consensus1" {
		t.Errorf("Expected consensus ID consensus1, got %s", msg.ConsensusID)
	}
	if msg.Vote != "option1" {
		t.Errorf("Expected vote option1, got %s", msg.Vote)
	}
	if msg.Weight != 1.0 {
		t.Errorf("Expected weight 1.0, got %f", msg.Weight)
	}
	if msg.Rationale != "My rationale" {
		t.Errorf("Expected rationale 'My rationale', got %s", msg.Rationale)
	}
}

func TestConsensusVoteMessageSerialization(t *testing.T) {
	msg := NewConsensusVoteMessage("team1", "agent1", "agent2", "consensus1", "option1", 1.0, "rationale", nil)

	// Serialize
	data, err := msg.ToJSON()
	if err != nil {
		t.Fatalf("Failed to serialize message: %v", err)
	}

	// Deserialize
	decoded, err := ConsensusVoteMessageFromJSON(data)
	if err != nil {
		t.Fatalf("Failed to deserialize message: %v", err)
	}

	if decoded.ConsensusID != msg.ConsensusID {
		t.Errorf("Expected consensus ID %s, got %s", msg.ConsensusID, decoded.ConsensusID)
	}
	if decoded.Vote != msg.Vote {
		t.Errorf("Expected vote %s, got %s", msg.Vote, decoded.Vote)
	}
	if decoded.Weight != msg.Weight {
		t.Errorf("Expected weight %f, got %f", msg.Weight, decoded.Weight)
	}
}

func TestMessageTimestamps(t *testing.T) {
	before := time.Now()
	time.Sleep(1 * time.Millisecond)

	task := NewTask("Test", "dev", nil)
	msg1 := NewTaskDelegationMessage("team1", "agent1", "agent2", task, nil)
	msg2 := NewTaskResultMessage("team1", "agent1", "agent2", "task1", TaskStatusCompleted, nil, nil)
	msg3 := NewConsensusRequestMessage("team1", "agent1", "c1", "Q?", []string{"a"}, VotingRuleMajority, 30, nil)
	msg4 := NewConsensusVoteMessage("team1", "agent1", "agent2", "c1", "a", 1.0, "", nil)

	time.Sleep(1 * time.Millisecond)
	after := time.Now()

	// Check all timestamps are within range
	if msg1.Timestamp.Before(before) || msg1.Timestamp.After(after) {
		t.Error("TaskDelegationMessage timestamp out of range")
	}
	if msg2.Timestamp.Before(before) || msg2.Timestamp.After(after) {
		t.Error("TaskResultMessage timestamp out of range")
	}
	if msg3.Timestamp.Before(before) || msg3.Timestamp.After(after) {
		t.Error("ConsensusRequestMessage timestamp out of range")
	}
	if msg4.Timestamp.Before(before) || msg4.Timestamp.After(after) {
		t.Error("ConsensusVoteMessage timestamp out of range")
	}
}
