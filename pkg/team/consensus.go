// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package team

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ConsensusRequest represents a request for consensus voting
type ConsensusRequest struct {
	ID         string
	Question   string
	Options    []string
	VotingRule VotingRule
	Voters     []string
	Timeout    time.Duration
	Context    map[string]any
	CreatedAt  time.Time
}

// ConsensusVote represents a vote from an agent
type ConsensusVote struct {
	ConsensusID string
	VoterID     string
	Vote        string
	Weight      float64
	Rationale   string
	Timestamp   time.Time
}

// ConsensusResult represents the outcome of a consensus
type ConsensusResult struct {
	ConsensusID string
	Question    string
	Outcome     string
	Votes       []*ConsensusVote
	TotalVotes  int
	Rule        VotingRule
	Timestamp   time.Time
}

// ConsensusManager manages consensus voting within a team
type ConsensusManager struct {
	requests map[string]*ConsensusRequest
	votes    map[string][]*ConsensusVote
	results  map[string]*ConsensusResult
	mu       sync.RWMutex
}

// NewConsensusManager creates a new consensus manager
func NewConsensusManager() *ConsensusManager {
	return &ConsensusManager{
		requests: make(map[string]*ConsensusRequest),
		votes:    make(map[string][]*ConsensusVote),
		results:  make(map[string]*ConsensusResult),
	}
}

// InitiateConsensus starts a new consensus voting process
func (cm *ConsensusManager) InitiateConsensus(question string, options []string, voters []string, rule VotingRule, timeout time.Duration, context map[string]any) (*ConsensusRequest, error) {
	if len(options) == 0 {
		return nil, fmt.Errorf("consensus must have at least one option")
	}
	if len(voters) == 0 {
		return nil, fmt.Errorf("consensus must have at least one voter")
	}

	request := &ConsensusRequest{
		ID:         generateConsensusID(),
		Question:   question,
		Options:    options,
		VotingRule: rule,
		Voters:     voters,
		Timeout:    timeout,
		Context:    context,
		CreatedAt:  time.Now(),
	}

	cm.mu.Lock()
	cm.requests[request.ID] = request
	cm.votes[request.ID] = []*ConsensusVote{}
	cm.mu.Unlock()

	return request, nil
}

// SubmitVote records a vote for a consensus
func (cm *ConsensusManager) SubmitVote(consensusID, voterID, vote string, weight float64, rationale string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	request, exists := cm.requests[consensusID]
	if !exists {
		return fmt.Errorf("consensus %s not found", consensusID)
	}

	// Validate voter
	validVoter := false
	for _, v := range request.Voters {
		if v == voterID {
			validVoter = true
			break
		}
	}
	if !validVoter {
		return fmt.Errorf("voter %s not authorized for consensus %s", voterID, consensusID)
	}

	// Check if already voted
	for _, v := range cm.votes[consensusID] {
		if v.VoterID == voterID {
			return fmt.Errorf("voter %s has already voted", voterID)
		}
	}

	// Validate vote option
	validOption := false
	for _, opt := range request.Options {
		if opt == vote {
			validOption = true
			break
		}
	}
	if !validOption {
		return fmt.Errorf("invalid vote option: %s", vote)
	}

	// Record vote
	voteRecord := &ConsensusVote{
		ConsensusID: consensusID,
		VoterID:     voterID,
		Vote:        vote,
		Weight:      weight,
		Rationale:   rationale,
		Timestamp:   time.Now(),
	}
	cm.votes[consensusID] = append(cm.votes[consensusID], voteRecord)

	return nil
}

// DetermineOutcome calculates the consensus outcome based on voting rule
func (cm *ConsensusManager) DetermineOutcome(consensusID string) (*ConsensusResult, error) {
	cm.mu.RLock()
	request, exists := cm.requests[consensusID]
	if !exists {
		cm.mu.RUnlock()
		return nil, fmt.Errorf("consensus %s not found", consensusID)
	}

	votes := cm.votes[consensusID]
	cm.mu.RUnlock()

	var outcome string
	switch request.VotingRule {
	case VotingRuleMajority:
		outcome = cm.determineMajority(votes)
	case VotingRuleUnanimous:
		outcome = cm.determineUnanimous(votes, request.Voters)
	case VotingRuleWeighted:
		outcome = cm.determineWeighted(votes)
	default:
		return nil, fmt.Errorf("unknown voting rule: %s", request.VotingRule)
	}

	result := &ConsensusResult{
		ConsensusID: consensusID,
		Question:    request.Question,
		Outcome:     outcome,
		Votes:       votes,
		TotalVotes:  len(votes),
		Rule:        request.VotingRule,
		Timestamp:   time.Now(),
	}

	cm.mu.Lock()
	cm.results[consensusID] = result
	cm.mu.Unlock()

	return result, nil
}

// determineMajority finds the option with the most votes
func (cm *ConsensusManager) determineMajority(votes []*ConsensusVote) string {
	if len(votes) == 0 {
		return ""
	}

	voteCounts := make(map[string]int)
	for _, vote := range votes {
		voteCounts[vote.Vote]++
	}

	maxVotes := 0
	var outcome string
	for option, count := range voteCounts {
		if count > maxVotes {
			maxVotes = count
			outcome = option
		}
	}

	return outcome
}

// determineUnanimous checks if all voters agree
func (cm *ConsensusManager) determineUnanimous(votes []*ConsensusVote, voters []string) string {
	if len(votes) != len(voters) {
		return "" // Not all voters have voted
	}

	if len(votes) == 0 {
		return ""
	}

	firstVote := votes[0].Vote
	for _, vote := range votes {
		if vote.Vote != firstVote {
			return "" // No unanimous agreement
		}
	}

	return firstVote
}

// determineWeighted calculates weighted vote outcome
func (cm *ConsensusManager) determineWeighted(votes []*ConsensusVote) string {
	if len(votes) == 0 {
		return ""
	}

	weightedCounts := make(map[string]float64)
	for _, vote := range votes {
		weightedCounts[vote.Vote] += vote.Weight
	}

	maxWeight := 0.0
	var outcome string
	for option, weight := range weightedCounts {
		if weight > maxWeight {
			maxWeight = weight
			outcome = option
		}
	}

	return outcome
}

// GetResult retrieves a consensus result
func (cm *ConsensusManager) GetResult(consensusID string) (*ConsensusResult, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	result, exists := cm.results[consensusID]
	return result, exists
}

// GetVotes retrieves all votes for a consensus
func (cm *ConsensusManager) GetVotes(consensusID string) ([]*ConsensusVote, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	votes, exists := cm.votes[consensusID]
	if !exists {
		return nil, fmt.Errorf("consensus %s not found", consensusID)
	}

	return votes, nil
}

// WaitForConsensus waits for consensus to complete or timeout
func (cm *ConsensusManager) WaitForConsensus(ctx context.Context, consensusID string) (*ConsensusResult, error) {
	cm.mu.RLock()
	request, exists := cm.requests[consensusID]
	cm.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("consensus %s not found", consensusID)
	}

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	timeout := time.After(request.Timeout)

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-timeout:
			// Timeout reached, determine outcome with partial votes
			result, err := cm.DetermineOutcome(consensusID)
			if err != nil {
				return nil, fmt.Errorf("timeout: %w", err)
			}
			// Mark result as incomplete if not all voters participated
			cm.mu.RLock()
			votes := cm.votes[consensusID]
			voterCount := len(request.Voters)
			cm.mu.RUnlock()

			if len(votes) < voterCount {
				result.Outcome = fmt.Sprintf("TIMEOUT_PARTIAL: %s (votes: %d/%d)", result.Outcome, len(votes), voterCount)
			}
			return result, nil
		case <-ticker.C:
			// Check votes without holding lock for too long
			cm.mu.RLock()
			votes := cm.votes[consensusID]
			voterCount := len(request.Voters)
			cm.mu.RUnlock()

			if len(votes) == voterCount {
				return cm.DetermineOutcome(consensusID)
			}
		}
	}
}

// generateConsensusID generates a unique consensus ID
func generateConsensusID() string {
	return "consensus_" + time.Now().Format("20060102150405") + "_" + randomString(8)
}
