// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package team

import (
	"testing"
	"time"
)

func TestMetricsCollector_RecordTeamCreation(t *testing.T) {
	mc := NewMetricsCollector()

	mc.RecordTeamCreation("team1")
	mc.RecordTeamCreation("team1")

	metrics := mc.GetMetrics("team1")
	if metrics == nil {
		t.Fatal("Expected metrics to exist")
	}

	if metrics.CreationCount != 2 {
		t.Errorf("Expected creation count 2, got %d", metrics.CreationCount)
	}
}

func TestMetricsCollector_RecordTaskCompletion(t *testing.T) {
	mc := NewMetricsCollector()
	mc.RecordTeamCreation("team1")

	mc.RecordTaskCompletion("team1", 100*time.Millisecond)
	mc.RecordTaskCompletion("team1", 200*time.Millisecond)

	metrics := mc.GetMetrics("team1")
	if metrics.TaskCompletionCount != 2 {
		t.Errorf("Expected completion count 2, got %d", metrics.TaskCompletionCount)
	}

	if metrics.AverageTaskDuration == 0 {
		t.Error("Expected non-zero average duration")
	}
}

func TestMetricsCollector_GetAllMetrics(t *testing.T) {
	mc := NewMetricsCollector()

	mc.RecordTeamCreation("team1")
	mc.RecordTeamCreation("team2")

	allMetrics := mc.GetAllMetrics()

	if len(allMetrics) != 2 {
		t.Errorf("Expected 2 teams, got %d", len(allMetrics))
	}

	if _, exists := allMetrics["team1"]; !exists {
		t.Error("Expected team1 metrics")
	}

	if _, exists := allMetrics["team2"]; !exists {
		t.Error("Expected team2 metrics")
	}
}

func TestHealthChecker_CheckAgentHeartbeat(t *testing.T) {
	hc := NewHealthChecker()

	team := &Team{
		ID:     "team1",
		Name:   "Test Team",
		Status: TeamStatusActive,
		Agents: map[string]*TeamAgent{
			"agent1": {
				AgentID:    "agent1",
				Role:       "developer",
				Status:     StatusIdle,
				LastActive: time.Now(),
			},
			"agent2": {
				AgentID:    "agent2",
				Role:       "tester",
				Status:     StatusIdle,
				LastActive: time.Now().Add(-2 * time.Minute),
			},
		},
	}

	check := hc.CheckAgentHeartbeat(team, 1*time.Minute)

	if check.Component != "agent_heartbeat" {
		t.Errorf("Expected component agent_heartbeat, got %s", check.Component)
	}

	if check.Status != "degraded" {
		t.Errorf("Expected status degraded, got %s", check.Status)
	}
}

func TestHealthChecker_CheckMessageBusConnectivity(t *testing.T) {
	hc := NewHealthChecker()

	// Test with valid bus
	bus := &mockMessageBus{}
	check := hc.CheckMessageBusConnectivity(bus)

	if check.Status != "healthy" {
		t.Errorf("Expected status healthy, got %s", check.Status)
	}

	// Test with nil bus
	check = hc.CheckMessageBusConnectivity(nil)

	if check.Status != "unhealthy" {
		t.Errorf("Expected status unhealthy, got %s", check.Status)
	}
}

func TestHealthChecker_CheckSharedContextAccessibility(t *testing.T) {
	hc := NewHealthChecker()

	ctx := NewSharedContext("test-team")
	check := hc.CheckSharedContextAccessibility(ctx)

	if check.Status != "healthy" {
		t.Errorf("Expected status healthy, got %s", check.Status)
	}

	if check.Component != "shared_context" {
		t.Errorf("Expected component shared_context, got %s", check.Component)
	}
}

func TestHealthChecker_GetHealthChecks(t *testing.T) {
	hc := NewHealthChecker()

	ctx := NewSharedContext("test-team")
	hc.CheckSharedContextAccessibility(ctx)
	hc.CheckMessageBusConnectivity(&mockMessageBus{})

	checks := hc.GetHealthChecks()

	if len(checks) != 2 {
		t.Errorf("Expected 2 health checks, got %d", len(checks))
	}
}

