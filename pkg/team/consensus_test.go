// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package team

import (
	"context"
	"testing"
	"time"
)

func TestConsensusManager_InitiateConsensus(t *testing.T) {
	cm := NewConsensusManager()

	options := []string{"option1", "option2", "option3"}
	voters := []string{"agent1", "agent2", "agent3"}

	request, err := cm.InitiateConsensus("Test question?", options, voters, VotingRuleMajority, 30*time.Second, nil)

	if err != nil {
		t.Fatalf("InitiateConsensus failed: %v", err)
	}
	if request.ID == "" {
		t.Error("Expected non-empty consensus ID")
	}
	if request.Question != "Test question?" {
		t.Errorf("Expected question 'Test question?', got %s", request.Question)
	}
	if len(request.Options) != 3 {
		t.Errorf("Expected 3 options, got %d", len(request.Options))
	}
}

func TestConsensusManager_InitiateConsensus_NoOptions(t *testing.T) {
	cm := NewConsensusManager()

	_, err := cm.InitiateConsensus("Test?", []string{}, []string{"agent1"}, VotingRuleMajority, 30*time.Second, nil)

	if err == nil {
		t.Error("Expected error for no options, got nil")
	}
}

func TestConsensusManager_InitiateConsensus_NoVoters(t *testing.T) {
	cm := NewConsensusManager()

	_, err := cm.InitiateConsensus("Test?", []string{"option1"}, []string{}, VotingRuleMajority, 30*time.Second, nil)

	if err == nil {
		t.Error("Expected error for no voters, got nil")
	}
}

func TestConsensusManager_SubmitVote_Success(t *testing.T) {
	cm := NewConsensusManager()

	options := []string{"option1", "option2"}
	voters := []string{"agent1", "agent2"}
	request, _ := cm.InitiateConsensus("Test?", options, voters, VotingRuleMajority, 30*time.Second, nil)

	err := cm.SubmitVote(request.ID, "agent1", "option1", 1.0, "My rationale")

	if err != nil {
		t.Fatalf("SubmitVote failed: %v", err)
	}

	votes, _ := cm.GetVotes(request.ID)
	if len(votes) != 1 {
		t.Errorf("Expected 1 vote, got %d", len(votes))
	}
	if votes[0].Vote != "option1" {
		t.Errorf("Expected vote option1, got %s", votes[0].Vote)
	}
}

func TestConsensusManager_SubmitVote_UnauthorizedVoter(t *testing.T) {
	cm := NewConsensusManager()

	options := []string{"option1", "option2"}
	voters := []string{"agent1", "agent2"}
	request, _ := cm.InitiateConsensus("Test?", options, voters, VotingRuleMajority, 30*time.Second, nil)

	err := cm.SubmitVote(request.ID, "agent3", "option1", 1.0, "")

	if err == nil {
		t.Error("Expected error for unauthorized voter, got nil")
	}
}

func TestConsensusManager_SubmitVote_DuplicateVote(t *testing.T) {
	cm := NewConsensusManager()

	options := []string{"option1", "option2"}
	voters := []string{"agent1", "agent2"}
	request, _ := cm.InitiateConsensus("Test?", options, voters, VotingRuleMajority, 30*time.Second, nil)

	cm.SubmitVote(request.ID, "agent1", "option1", 1.0, "")
	err := cm.SubmitVote(request.ID, "agent1", "option2", 1.0, "")

	if err == nil {
		t.Error("Expected error for duplicate vote, got nil")
	}
}

func TestConsensusManager_SubmitVote_InvalidOption(t *testing.T) {
	cm := NewConsensusManager()

	options := []string{"option1", "option2"}
	voters := []string{"agent1", "agent2"}
	request, _ := cm.InitiateConsensus("Test?", options, voters, VotingRuleMajority, 30*time.Second, nil)

	err := cm.SubmitVote(request.ID, "agent1", "invalid", 1.0, "")

	if err == nil {
		t.Error("Expected error for invalid option, got nil")
	}
}

func TestConsensusManager_DetermineOutcome_Majority(t *testing.T) {
	cm := NewConsensusManager()

	options := []string{"option1", "option2"}
	voters := []string{"agent1", "agent2", "agent3"}
	request, _ := cm.InitiateConsensus("Test?", options, voters, VotingRuleMajority, 30*time.Second, nil)

	cm.SubmitVote(request.ID, "agent1", "option1", 1.0, "")
	cm.SubmitVote(request.ID, "agent2", "option1", 1.0, "")
	cm.SubmitVote(request.ID, "agent3", "option2", 1.0, "")

	result, err := cm.DetermineOutcome(request.ID)

	if err != nil {
		t.Fatalf("DetermineOutcome failed: %v", err)
	}
	if result.Outcome != "option1" {
		t.Errorf("Expected outcome option1, got %s", result.Outcome)
	}
	if result.TotalVotes != 3 {
		t.Errorf("Expected 3 total votes, got %d", result.TotalVotes)
	}
}

func TestConsensusManager_DetermineOutcome_Unanimous(t *testing.T) {
	cm := NewConsensusManager()

	options := []string{"option1", "option2"}
	voters := []string{"agent1", "agent2", "agent3"}
	request, _ := cm.InitiateConsensus("Test?", options, voters, VotingRuleUnanimous, 30*time.Second, nil)

	cm.SubmitVote(request.ID, "agent1", "option1", 1.0, "")
	cm.SubmitVote(request.ID, "agent2", "option1", 1.0, "")
	cm.SubmitVote(request.ID, "agent3", "option1", 1.0, "")

	result, err := cm.DetermineOutcome(request.ID)

	if err != nil {
		t.Fatalf("DetermineOutcome failed: %v", err)
	}
	if result.Outcome != "option1" {
		t.Errorf("Expected unanimous outcome option1, got %s", result.Outcome)
	}
}

func TestConsensusManager_DetermineOutcome_Unanimous_NoAgreement(t *testing.T) {
	cm := NewConsensusManager()

	options := []string{"option1", "option2"}
	voters := []string{"agent1", "agent2", "agent3"}
	request, _ := cm.InitiateConsensus("Test?", options, voters, VotingRuleUnanimous, 30*time.Second, nil)

	cm.SubmitVote(request.ID, "agent1", "option1", 1.0, "")
	cm.SubmitVote(request.ID, "agent2", "option1", 1.0, "")
	cm.SubmitVote(request.ID, "agent3", "option2", 1.0, "")

	result, err := cm.DetermineOutcome(request.ID)

	if err != nil {
		t.Fatalf("DetermineOutcome failed: %v", err)
	}
	if result.Outcome != "" {
		t.Errorf("Expected no unanimous outcome, got %s", result.Outcome)
	}
}

func TestConsensusManager_DetermineOutcome_Weighted(t *testing.T) {
	cm := NewConsensusManager()

	options := []string{"option1", "option2"}
	voters := []string{"agent1", "agent2", "agent3"}
	request, _ := cm.InitiateConsensus("Test?", options, voters, VotingRuleWeighted, 30*time.Second, nil)

	cm.SubmitVote(request.ID, "agent1", "option1", 3.0, "")
	cm.SubmitVote(request.ID, "agent2", "option2", 1.0, "")
	cm.SubmitVote(request.ID, "agent3", "option2", 1.0, "")

	result, err := cm.DetermineOutcome(request.ID)

	if err != nil {
		t.Fatalf("DetermineOutcome failed: %v", err)
	}
	if result.Outcome != "option1" {
		t.Errorf("Expected weighted outcome option1, got %s", result.Outcome)
	}
}

func TestConsensusManager_WaitForConsensus_AllVoted(t *testing.T) {
	cm := NewConsensusManager()

	options := []string{"option1", "option2"}
	voters := []string{"agent1", "agent2"}
	request, _ := cm.InitiateConsensus("Test?", options, voters, VotingRuleMajority, 5*time.Second, nil)

	// Submit votes in background
	go func() {
		time.Sleep(100 * time.Millisecond)
		cm.SubmitVote(request.ID, "agent1", "option1", 1.0, "")
		time.Sleep(100 * time.Millisecond)
		cm.SubmitVote(request.ID, "agent2", "option1", 1.0, "")
	}()

	ctx := context.Background()
	result, err := cm.WaitForConsensus(ctx, request.ID)

	if err != nil {
		t.Fatalf("WaitForConsensus failed: %v", err)
	}
	if result.Outcome != "option1" {
		t.Errorf("Expected outcome option1, got %s", result.Outcome)
	}
}

func TestConsensusManager_WaitForConsensus_Timeout(t *testing.T) {
	cm := NewConsensusManager()

	options := []string{"option1", "option2"}
	voters := []string{"agent1", "agent2", "agent3"}
	request, _ := cm.InitiateConsensus("Test?", options, voters, VotingRuleMajority, 200*time.Millisecond, nil)

	// Only submit partial votes
	cm.SubmitVote(request.ID, "agent1", "option1", 1.0, "")

	ctx := context.Background()
	result, err := cm.WaitForConsensus(ctx, request.ID)

	if err != nil {
		t.Fatalf("WaitForConsensus failed: %v", err)
	}
	// Should determine outcome with partial votes
	if result.TotalVotes != 1 {
		t.Errorf("Expected 1 vote at timeout, got %d", result.TotalVotes)
	}
}

func TestConsensusManager_WaitForConsensus_ContextCancelled(t *testing.T) {
	cm := NewConsensusManager()

	options := []string{"option1", "option2"}
	voters := []string{"agent1", "agent2"}
	request, _ := cm.InitiateConsensus("Test?", options, voters, VotingRuleMajority, 5*time.Second, nil)

	ctx, cancel := context.WithCancel(context.Background())

	// Cancel context after 100ms
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	_, err := cm.WaitForConsensus(ctx, request.ID)

	if err == nil {
		t.Error("Expected context cancelled error, got nil")
	}
	if err != context.Canceled {
		t.Errorf("Expected context.Canceled, got %v", err)
	}
}

func TestConsensusManager_GetResult(t *testing.T) {
	cm := NewConsensusManager()

	options := []string{"option1", "option2"}
	voters := []string{"agent1", "agent2"}
	request, _ := cm.InitiateConsensus("Test?", options, voters, VotingRuleMajority, 30*time.Second, nil)

	cm.SubmitVote(request.ID, "agent1", "option1", 1.0, "")
	cm.SubmitVote(request.ID, "agent2", "option1", 1.0, "")
	cm.DetermineOutcome(request.ID)

	result, exists := cm.GetResult(request.ID)

	if !exists {
		t.Error("Expected result to exist")
	}
	if result.Outcome != "option1" {
		t.Errorf("Expected outcome option1, got %s", result.Outcome)
	}
}
