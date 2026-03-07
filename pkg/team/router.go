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

	"github.com/sipeed/picoclaw/pkg/logger"
)

// DelegationRouter routes tasks between team members and prevents circular delegation
type DelegationRouter struct {
	maxDepth      int
	delegationMap map[string]*DelegationChain
	mu            sync.RWMutex
}

// DelegationChain tracks the delegation path of a task
type DelegationChain struct {
	TaskID    string
	Chain     []string
	Depth     int
	CreatedAt time.Time
}

// NewDelegationRouter creates a new DelegationRouter
func NewDelegationRouter(maxDepth int) *DelegationRouter {
	if maxDepth <= 0 {
		maxDepth = 5 // Default max depth
	}

	return &DelegationRouter{
		maxDepth:      maxDepth,
		delegationMap: make(map[string]*DelegationChain),
	}
}

// RouteTask routes a task to an appropriate agent based on required role
func (dr *DelegationRouter) RouteTask(ctx context.Context, task *Task, team *Team) (string, error) {
	dr.mu.RLock()
	defer dr.mu.RUnlock()

	// Find agents with required role
	var candidates []string
	for agentID, agent := range team.Agents {
		if agent.Role == task.RequiredRole {
			// Check agent status
			if agent.Status == StatusIdle || agent.Status == StatusWorking {
				candidates = append(candidates, agentID)
			}
		}
	}

	if len(candidates) == 0 {
		return "", fmt.Errorf("no available agent with role '%s'", task.RequiredRole)
	}

	// Select first available agent (simple strategy)
	selectedAgent := candidates[0]

	logger.DebugCF("router", "Task routed to agent", map[string]any{
		"task_id":  task.ID,
		"agent_id": selectedAgent,
		"role":     task.RequiredRole,
		"team_id":  team.ID,
	})

	return selectedAgent, nil
}

// ValidateDelegation checks if a delegation is valid (no circular dependencies, within depth limit)
func (dr *DelegationRouter) ValidateDelegation(task *Task, targetAgentID string) error {
	dr.mu.RLock()
	defer dr.mu.RUnlock()

	// Check if target agent is already in delegation chain
	for _, agentID := range task.DelegationChain {
		if agentID == targetAgentID {
			logger.WarnCF("router", "Circular delegation detected", map[string]any{
				"task_id":          task.ID,
				"target_agent":     targetAgentID,
				"delegation_chain": task.DelegationChain,
			})
			return fmt.Errorf("circular delegation detected: agent '%s' already in delegation chain", targetAgentID)
		}
	}

	// Check delegation depth
	if len(task.DelegationChain) >= dr.maxDepth {
		logger.WarnCF("router", "Max delegation depth exceeded", map[string]any{
			"task_id":       task.ID,
			"current_depth": len(task.DelegationChain),
			"max_depth":     dr.maxDepth,
		})
		return fmt.Errorf("max delegation depth %d exceeded (current: %d)", dr.maxDepth, len(task.DelegationChain))
	}

	return nil
}

// RecordDelegation records a delegation in the tracking map
func (dr *DelegationRouter) RecordDelegation(task *Task, targetAgentID string) {
	dr.mu.Lock()
	defer dr.mu.Unlock()

	// Update task delegation chain
	task.DelegationChain = append(task.DelegationChain, targetAgentID)

	// Record in delegation map
	chain := &DelegationChain{
		TaskID:    task.ID,
		Chain:     make([]string, len(task.DelegationChain)),
		Depth:     len(task.DelegationChain),
		CreatedAt: time.Now(),
	}
	copy(chain.Chain, task.DelegationChain)

	dr.delegationMap[task.ID] = chain

	logger.DebugCF("router", "Delegation recorded", map[string]any{
		"task_id":          task.ID,
		"target_agent":     targetAgentID,
		"delegation_depth": chain.Depth,
	})
}

// GetDelegationChain returns the delegation chain for a task
func (dr *DelegationRouter) GetDelegationChain(taskID string) (*DelegationChain, bool) {
	dr.mu.RLock()
	defer dr.mu.RUnlock()

	chain, exists := dr.delegationMap[taskID]
	return chain, exists
}

// ClearDelegationChain removes the delegation chain for a completed task
func (dr *DelegationRouter) ClearDelegationChain(taskID string) {
	dr.mu.Lock()
	defer dr.mu.Unlock()

	delete(dr.delegationMap, taskID)

	logger.DebugCF("router", "Delegation chain cleared", map[string]any{
		"task_id": taskID,
	})
}

// GetStats returns delegation statistics
func (dr *DelegationRouter) GetStats() map[string]any {
	dr.mu.RLock()
	defer dr.mu.RUnlock()

	totalChains := len(dr.delegationMap)
	maxDepthSeen := 0

	for _, chain := range dr.delegationMap {
		if chain.Depth > maxDepthSeen {
			maxDepthSeen = chain.Depth
		}
	}

	return map[string]any{
		"total_active_chains": totalChains,
		"max_depth_seen":      maxDepthSeen,
		"max_depth_limit":     dr.maxDepth,
	}
}
