// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

//go:build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sipeed/picoclaw/pkg/team"
)

// ConsensusExample demonstrates consensus voting mechanisms
func main() {
	fmt.Println("=== Consensus Voting Example ===\n")

	// Create consensus manager
	cm := team.NewConsensusManager()

	// Example 1: Majority Voting
	fmt.Println("Example 1: Majority Voting")
	fmt.Println("Question: Should we proceed with the new architecture?")
	fmt.Println("Voters: architect, developer1, developer2, tester")
	fmt.Println()

	request1, err := cm.InitiateConsensus(
		"Should we proceed with the new architecture?",
		[]string{"yes", "no"},
		[]string{"architect", "developer1", "developer2", "tester"},
		team.VotingRuleMajority,
		30*time.Second,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to initiate consensus: %v", err)
	}

	// Submit votes
	cm.SubmitVote(request1.ID, "architect", "yes", 1.0, "Architecture looks solid")
	cm.SubmitVote(request1.ID, "developer1", "yes", 1.0, "Implementation is feasible")
	cm.SubmitVote(request1.ID, "developer2", "no", 1.0, "Concerns about complexity")
	cm.SubmitVote(request1.ID, "tester", "yes", 1.0, "Testable design")

	// Determine outcome
	result1, err := cm.DetermineOutcome(request1.ID)
	if err != nil {
		log.Fatalf("Failed to determine outcome: %v", err)
	}

	fmt.Printf("Votes:\n")
	fmt.Printf("  yes: 3\n")
	fmt.Printf("  no: 1\n")
	fmt.Printf("✓ Outcome: %s (majority wins)\n\n", result1.Outcome)

	// Example 2: Unanimous Voting
	fmt.Println("Example 2: Unanimous Voting")
	fmt.Println("Question: Should we deploy to production?")
	fmt.Println("Voters: lead, developer, tester, security")
	fmt.Println()

	request2, err := cm.InitiateConsensus(
		"Should we deploy to production?",
		[]string{"approve", "reject"},
		[]string{"lead", "developer", "tester", "security"},
		team.VotingRuleUnanimous,
		30*time.Second,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to initiate consensus: %v", err)
	}

	// All must approve
	cm.SubmitVote(request2.ID, "lead", "approve", 1.0, "All requirements met")
	cm.SubmitVote(request2.ID, "developer", "approve", 1.0, "Code is production ready")
	cm.SubmitVote(request2.ID, "tester", "approve", 1.0, "All tests passing")
	cm.SubmitVote(request2.ID, "security", "approve", 1.0, "No security issues found")

	result2, err := cm.DetermineOutcome(request2.ID)
	if err != nil {
		log.Fatalf("Failed to determine outcome: %v", err)
	}

	fmt.Printf("Votes:\n")
	fmt.Printf("  approve: 4\n")
	fmt.Printf("  reject: 0\n")
	fmt.Printf("✓ Outcome: %s (unanimous)\n\n", result2.Outcome)

	// Example 3: Weighted Voting
	fmt.Println("Example 3: Weighted Voting")
	fmt.Println("Question: Which technology stack should we use?")
	fmt.Println("Voters: senior_architect (weight: 3.0), developer1 (weight: 1.0), developer2 (weight: 1.0)")
	fmt.Println()

	request3, err := cm.InitiateConsensus(
		"Which technology stack should we use?",
		[]string{"stack_a", "stack_b"},
		[]string{"senior_architect", "developer1", "developer2"},
		team.VotingRuleWeighted,
		30*time.Second,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to initiate consensus: %v", err)
	}

	// Senior architect has more weight
	cm.SubmitVote(request3.ID, "senior_architect", "stack_a", 3.0, "Better long-term scalability")
	cm.SubmitVote(request3.ID, "developer1", "stack_b", 1.0, "Easier to learn")
	cm.SubmitVote(request3.ID, "developer2", "stack_b", 1.0, "More community support")

	result3, err := cm.DetermineOutcome(request3.ID)
	if err != nil {
		log.Fatalf("Failed to determine outcome: %v", err)
	}

	fmt.Printf("Votes:\n")
	fmt.Printf("  stack_a: 3.0 (senior_architect)\n")
	fmt.Printf("  stack_b: 2.0 (developer1 + developer2)\n")
	fmt.Printf("✓ Outcome: %s (weighted score wins)\n\n", result3.Outcome)

	// Example 4: Consensus with Timeout
	fmt.Println("Example 4: Consensus with Timeout")
	fmt.Println("Question: Should we refactor the authentication module?")
	fmt.Println("Timeout: 2 seconds")
	fmt.Println()

	request4, err := cm.InitiateConsensus(
		"Should we refactor the authentication module?",
		[]string{"yes", "no"},
		[]string{"agent1", "agent2", "agent3"},
		team.VotingRuleMajority,
		2*time.Second,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to initiate consensus: %v", err)
	}

	// Only partial votes
	cm.SubmitVote(request4.ID, "agent1", "yes", 1.0, "Code needs refactoring")
	cm.SubmitVote(request4.ID, "agent2", "yes", 1.0, "Will improve maintainability")
	// agent3 doesn't vote

	// Wait for consensus with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result4, err := cm.WaitForConsensus(ctx, request4.ID)
	if err != nil {
		log.Printf("Consensus timed out: %v\n", err)
		// Determine outcome with partial votes
		result4, _ = cm.DetermineOutcome(request4.ID)
	}

	fmt.Printf("Votes received: 2/3\n")
	fmt.Printf("  yes: 2\n")
	fmt.Printf("  no: 0\n")
	fmt.Printf("✓ Outcome: %s (partial votes, majority)\n\n", result4.Outcome)

	// Example 5: Vote Validation
	fmt.Println("Example 5: Vote Validation")
	fmt.Println("Demonstrating vote validation rules")
	fmt.Println()

	request5, err := cm.InitiateConsensus(
		"Test question",
		[]string{"option1", "option2"},
		[]string{"voter1", "voter2"},
		team.VotingRuleMajority,
		30*time.Second,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to initiate consensus: %v", err)
	}

	// Valid vote
	err = cm.SubmitVote(request5.ID, "voter1", "option1", 1.0, "I prefer option1")
	if err != nil {
		fmt.Printf("✗ Valid vote rejected: %v\n", err)
	} else {
		fmt.Printf("✓ Valid vote accepted\n")
	}

	// Duplicate vote (should fail)
	err = cm.SubmitVote(request5.ID, "voter1", "option2", 1.0, "Changed my mind")
	if err != nil {
		fmt.Printf("✓ Duplicate vote rejected: %v\n", err)
	} else {
		fmt.Printf("✗ Duplicate vote accepted (should have been rejected)\n")
	}

	// Unauthorized voter (should fail)
	err = cm.SubmitVote(request5.ID, "unauthorized", "option1", 1.0, "I want to vote")
	if err != nil {
		fmt.Printf("✓ Unauthorized voter rejected: %v\n", err)
	} else {
		fmt.Printf("✗ Unauthorized voter accepted (should have been rejected)\n")
	}

	// Invalid option (should fail)
	err = cm.SubmitVote(request5.ID, "voter2", "invalid_option", 1.0, "This option doesn't exist")
	if err != nil {
		fmt.Printf("✓ Invalid option rejected: %v\n", err)
	} else {
		fmt.Printf("✗ Invalid option accepted (should have been rejected)\n")
	}

	fmt.Println()
	fmt.Println("=== Example Complete ===")
}
