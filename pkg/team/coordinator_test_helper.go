// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package team

import (
	"context"

	"github.com/sipeed/picoclaw/pkg/bus"
)

// Helper function to create coordinator for tests
func newTestCoordinator(agentID, teamID string, team *Team, pattern CollaborationPattern) *CoordinatorAgent {
	messageBus := bus.NewMessageBus()
	delegationRouter := NewDelegationRouter(5)
	ctx := context.Background()
	return NewCoordinatorAgent(agentID, teamID, team, pattern, messageBus, delegationRouter, ctx)
}
