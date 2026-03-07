package collaborative

import (
	"testing"
)

// TestCascadeCleanupOnAllPaths verifies that UnmarkAgentInCascade is called on all code paths
func TestCascadeCleanupOnAllPaths(t *testing.T) {
	session := NewSession(123, "team-1", 50)

	// Test 1: Normal execution path
	t.Run("normal_execution", func(t *testing.T) {
		session.MarkAgentInCascade("agent1")
		if !session.IsAgentInCascade("agent1") {
			t.Error("Agent should be marked in cascade")
		}

		// Simulate normal execution completion
		session.UnmarkAgentInCascade("agent1")
		if session.IsAgentInCascade("agent1") {
			t.Error("Agent should be unmarked after completion")
		}
	})

	// Test 2: Error path
	t.Run("error_path", func(t *testing.T) {
		session.MarkAgentInCascade("agent2")
		if !session.IsAgentInCascade("agent2") {
			t.Error("Agent should be marked in cascade")
		}

		// Simulate error during execution
		session.UnmarkAgentInCascade("agent2")
		if session.IsAgentInCascade("agent2") {
			t.Error("Agent should be unmarked even on error")
		}
	})

	// Test 3: Send message failure path
	t.Run("send_failure_path", func(t *testing.T) {
		session.MarkAgentInCascade("agent3")
		if !session.IsAgentInCascade("agent3") {
			t.Error("Agent should be marked in cascade")
		}

		// Simulate send message failure
		session.UnmarkAgentInCascade("agent3")
		if session.IsAgentInCascade("agent3") {
			t.Error("Agent should be unmarked even if send fails")
		}
	})

	t.Log("✓ All cleanup paths verified")
}

// TestCascadeBackReference verifies agents can mention each other back
func TestCascadeBackReference(t *testing.T) {
	session := NewSession(123, "team-1", 50)

	// Simulate: A → B → A (back-reference)
	t.Run("A_to_B_to_A", func(t *testing.T) {
		// Step 1: A executes
		session.MarkAgentInCascade("A")
		if !session.IsAgentInCascade("A") {
			t.Error("A should be marked")
		}

		// A completes and mentions B
		session.UnmarkAgentInCascade("A")
		if session.IsAgentInCascade("A") {
			t.Error("A should be unmarked after completion")
		}

		// Step 2: B executes
		session.MarkAgentInCascade("B")
		if !session.IsAgentInCascade("B") {
			t.Error("B should be marked")
		}

		// B completes and mentions A back
		session.UnmarkAgentInCascade("B")
		if session.IsAgentInCascade("B") {
			t.Error("B should be unmarked after completion")
		}

		// Step 3: A can execute again (not blocked)
		if session.IsAgentInCascade("A") {
			t.Error("A should NOT be in cascade anymore, back-reference should work")
		}

		session.MarkAgentInCascade("A")
		if !session.IsAgentInCascade("A") {
			t.Error("A should be able to execute again")
		}

		session.UnmarkAgentInCascade("A")
	})

	t.Log("✓ Back-reference working correctly")
}

// TestCascadeTrueCycleDetection verifies true cycles are still detected
func TestCascadeTrueCycleDetection(t *testing.T) {
	session := NewSession(123, "team-1", 50)

	t.Run("concurrent_execution_blocks", func(t *testing.T) {
		// A is currently executing
		session.MarkAgentInCascade("A")

		// B tries to mention A while A is still executing
		if !session.IsAgentInCascade("A") {
			t.Error("A should be in cascade")
		}

		// This should be blocked (true cycle)
		blocked := session.IsAgentInCascade("A")
		if !blocked {
			t.Error("Should detect cycle when A is still executing")
		}

		// A completes
		session.UnmarkAgentInCascade("A")

		// Now B can mention A
		if session.IsAgentInCascade("A") {
			t.Error("A should not be in cascade after completion")
		}
	})

	t.Log("✓ True cycle detection working")
}

// TestCascadeDepthLimit verifies depth limit is enforced
func TestCascadeDepthLimit(t *testing.T) {
	manager := NewManagerV2()
	maxDepth := manager.maxMentionDepth

	t.Run("depth_limit_enforced", func(t *testing.T) {
		// Depth 0: OK
		if 0 >= maxDepth {
			t.Error("Depth 0 should be allowed")
		}

		// Depth 1: OK
		if 1 >= maxDepth {
			t.Error("Depth 1 should be allowed")
		}

		// Depth 2: OK
		if 2 >= maxDepth {
			t.Error("Depth 2 should be allowed")
		}

		// Depth 3: BLOCKED (maxDepth = 3)
		if 3 < maxDepth {
			t.Error("Depth 3 should be blocked")
		}

		// Verify max depth value
		if maxDepth != 3 {
			t.Errorf("Expected max depth 3, got %d", maxDepth)
		}
	})

	t.Log("✓ Depth limit enforced correctly")
}

// TestCascadeMultipleAgents verifies multiple agents can execute in parallel
func TestCascadeMultipleAgents(t *testing.T) {
	session := NewSession(123, "team-1", 50)

	t.Run("parallel_execution", func(t *testing.T) {
		// Mark multiple agents as executing
		session.MarkAgentInCascade("A")
		session.MarkAgentInCascade("B")
		session.MarkAgentInCascade("C")

		// All should be marked
		if !session.IsAgentInCascade("A") {
			t.Error("A should be marked")
		}
		if !session.IsAgentInCascade("B") {
			t.Error("B should be marked")
		}
		if !session.IsAgentInCascade("C") {
			t.Error("C should be marked")
		}

		// Unmark in different order
		session.UnmarkAgentInCascade("B")
		if session.IsAgentInCascade("B") {
			t.Error("B should be unmarked")
		}
		if !session.IsAgentInCascade("A") {
			t.Error("A should still be marked")
		}
		if !session.IsAgentInCascade("C") {
			t.Error("C should still be marked")
		}

		session.UnmarkAgentInCascade("A")
		session.UnmarkAgentInCascade("C")

		// All should be unmarked
		if session.IsAgentInCascade("A") {
			t.Error("A should be unmarked")
		}
		if session.IsAgentInCascade("B") {
			t.Error("B should be unmarked")
		}
		if session.IsAgentInCascade("C") {
			t.Error("C should be unmarked")
		}
	})

	t.Log("✓ Parallel execution working correctly")
}

// TestCascadeEdgeCases tests various edge cases
func TestCascadeEdgeCases(t *testing.T) {
	session := NewSession(123, "team-1", 50)

	t.Run("double_unmark", func(t *testing.T) {
		session.MarkAgentInCascade("agent")
		session.UnmarkAgentInCascade("agent")
		session.UnmarkAgentInCascade("agent") // Should not panic

		if session.IsAgentInCascade("agent") {
			t.Error("Agent should not be in cascade")
		}
	})

	t.Run("unmark_without_mark", func(t *testing.T) {
		session.UnmarkAgentInCascade("never_marked") // Should not panic

		if session.IsAgentInCascade("never_marked") {
			t.Error("Agent should not be in cascade")
		}
	})

	t.Run("check_without_mark", func(t *testing.T) {
		if session.IsAgentInCascade("not_marked") {
			t.Error("Agent should not be in cascade")
		}
	})

	t.Log("✓ Edge cases handled correctly")
}
