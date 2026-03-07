// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package team

import (
	"fmt"
	"sync"
	"time"
)

// TeamMetrics tracks team-level metrics
type TeamMetrics struct {
	TeamID               string
	CreationCount        int64
	DissolutionCount     int64
	TaskDelegationCount  int64
	TaskCompletionCount  int64
	TaskFailureCount     int64
	ConsensusCount       int64
	AgentAdditionCount   int64
	AgentRemovalCount    int64
	AverageTaskDuration  time.Duration
	MessageBusThroughput int64
	SharedContextOps     int64
	mu                   sync.RWMutex
}

// MetricsCollector collects and aggregates team metrics
type MetricsCollector struct {
	teamMetrics map[string]*TeamMetrics
	mu          sync.RWMutex
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		teamMetrics: make(map[string]*TeamMetrics),
	}
}

// RecordTeamCreation records a team creation event
func (mc *MetricsCollector) RecordTeamCreation(teamID string) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	metrics, exists := mc.teamMetrics[teamID]
	if !exists {
		metrics = &TeamMetrics{TeamID: teamID}
		mc.teamMetrics[teamID] = metrics
	}

	metrics.mu.Lock()
	metrics.CreationCount++
	metrics.mu.Unlock()
}

// RecordTeamDissolution records a team dissolution event
func (mc *MetricsCollector) RecordTeamDissolution(teamID string) {
	mc.mu.RLock()
	metrics, exists := mc.teamMetrics[teamID]
	mc.mu.RUnlock()

	if !exists {
		return
	}

	metrics.mu.Lock()
	metrics.DissolutionCount++
	metrics.mu.Unlock()
}

// RecordTaskDelegation records a task delegation event
func (mc *MetricsCollector) RecordTaskDelegation(teamID string) {
	mc.mu.RLock()
	metrics, exists := mc.teamMetrics[teamID]
	mc.mu.RUnlock()

	if !exists {
		return
	}

	metrics.mu.Lock()
	metrics.TaskDelegationCount++
	metrics.mu.Unlock()
}

// RecordTaskCompletion records a task completion event
func (mc *MetricsCollector) RecordTaskCompletion(teamID string, duration time.Duration) {
	mc.mu.RLock()
	metrics, exists := mc.teamMetrics[teamID]
	mc.mu.RUnlock()

	if !exists {
		return
	}

	metrics.mu.Lock()
	metrics.TaskCompletionCount++

	// Update average duration
	if metrics.AverageTaskDuration == 0 {
		metrics.AverageTaskDuration = duration
	} else {
		metrics.AverageTaskDuration = (metrics.AverageTaskDuration + duration) / 2
	}
	metrics.mu.Unlock()
}

// RecordTaskFailure records a task failure event
func (mc *MetricsCollector) RecordTaskFailure(teamID string) {
	mc.mu.RLock()
	metrics, exists := mc.teamMetrics[teamID]
	mc.mu.RUnlock()

	if !exists {
		return
	}

	metrics.mu.Lock()
	metrics.TaskFailureCount++
	metrics.mu.Unlock()
}

// RecordConsensus records a consensus event
func (mc *MetricsCollector) RecordConsensus(teamID string) {
	mc.mu.RLock()
	metrics, exists := mc.teamMetrics[teamID]
	mc.mu.RUnlock()

	if !exists {
		return
	}

	metrics.mu.Lock()
	metrics.ConsensusCount++
	metrics.mu.Unlock()
}

// RecordAgentAddition records an agent addition event
func (mc *MetricsCollector) RecordAgentAddition(teamID string) {
	mc.mu.RLock()
	metrics, exists := mc.teamMetrics[teamID]
	mc.mu.RUnlock()

	if !exists {
		return
	}

	metrics.mu.Lock()
	metrics.AgentAdditionCount++
	metrics.mu.Unlock()
}

// RecordAgentRemoval records an agent removal event
func (mc *MetricsCollector) RecordAgentRemoval(teamID string) {
	mc.mu.RLock()
	metrics, exists := mc.teamMetrics[teamID]
	mc.mu.RUnlock()

	if !exists {
		return
	}

	metrics.mu.Lock()
	metrics.AgentRemovalCount++
	metrics.mu.Unlock()
}

// RecordMessageBusOperation records a message bus operation
func (mc *MetricsCollector) RecordMessageBusOperation(teamID string) {
	mc.mu.RLock()
	metrics, exists := mc.teamMetrics[teamID]
	mc.mu.RUnlock()

	if !exists {
		return
	}

	metrics.mu.Lock()
	metrics.MessageBusThroughput++
	metrics.mu.Unlock()
}

// RecordSharedContextOperation records a shared context operation
func (mc *MetricsCollector) RecordSharedContextOperation(teamID string) {
	mc.mu.RLock()
	metrics, exists := mc.teamMetrics[teamID]
	mc.mu.RUnlock()

	if !exists {
		return
	}

	metrics.mu.Lock()
	metrics.SharedContextOps++
	metrics.mu.Unlock()
}

// GetMetrics returns metrics for a team
func (mc *MetricsCollector) GetMetrics(teamID string) *TeamMetrics {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	metrics, exists := mc.teamMetrics[teamID]
	if !exists {
		return nil
	}

	// Return a copy
	metrics.mu.RLock()
	defer metrics.mu.RUnlock()

	return &TeamMetrics{
		TeamID:               metrics.TeamID,
		CreationCount:        metrics.CreationCount,
		DissolutionCount:     metrics.DissolutionCount,
		TaskDelegationCount:  metrics.TaskDelegationCount,
		TaskCompletionCount:  metrics.TaskCompletionCount,
		TaskFailureCount:     metrics.TaskFailureCount,
		ConsensusCount:       metrics.ConsensusCount,
		AgentAdditionCount:   metrics.AgentAdditionCount,
		AgentRemovalCount:    metrics.AgentRemovalCount,
		AverageTaskDuration:  metrics.AverageTaskDuration,
		MessageBusThroughput: metrics.MessageBusThroughput,
		SharedContextOps:     metrics.SharedContextOps,
	}
}

// GetAllMetrics returns metrics for all teams
func (mc *MetricsCollector) GetAllMetrics() map[string]*TeamMetrics {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	result := make(map[string]*TeamMetrics)
	for teamID := range mc.teamMetrics {
		result[teamID] = mc.GetMetrics(teamID)
	}

	return result
}

// HealthCheck performs health checks on team components
type HealthCheck struct {
	Component string
	Status    string
	Message   string
	Timestamp time.Time
}

// HealthChecker performs health checks
type HealthChecker struct {
	checks []HealthCheck
	mu     sync.RWMutex
}

// NewHealthChecker creates a new health checker
func NewHealthChecker() *HealthChecker {
	return &HealthChecker{
		checks: []HealthCheck{},
	}
}

// CheckAgentHeartbeat checks if agents are responsive
func (hc *HealthChecker) CheckAgentHeartbeat(team *Team, timeout time.Duration) HealthCheck {
	unresponsiveCount := 0
	for _, agent := range team.Agents {
		if time.Since(agent.LastActive) > timeout {
			unresponsiveCount++
		}
	}

	status := "healthy"
	message := "All agents responsive"
	if unresponsiveCount > 0 {
		status = "degraded"
		message = fmt.Sprintf("%d agents unresponsive", unresponsiveCount)
	}

	check := HealthCheck{
		Component: "agent_heartbeat",
		Status:    status,
		Message:   message,
		Timestamp: time.Now(),
	}

	hc.mu.Lock()
	hc.checks = append(hc.checks, check)
	hc.mu.Unlock()

	return check
}

// CheckMessageBusConnectivity checks message bus connectivity
func (hc *HealthChecker) CheckMessageBusConnectivity(bus any) HealthCheck {
	// Simple connectivity check
	status := "healthy"
	message := "Message bus connected"

	if bus == nil {
		status = "unhealthy"
		message = "Message bus not available"
	}

	check := HealthCheck{
		Component: "message_bus",
		Status:    status,
		Message:   message,
		Timestamp: time.Now(),
	}

	hc.mu.Lock()
	hc.checks = append(hc.checks, check)
	hc.mu.Unlock()

	return check
}

// CheckSharedContextAccessibility checks shared context accessibility
func (hc *HealthChecker) CheckSharedContextAccessibility(ctx *SharedContext) HealthCheck {
	status := "healthy"
	message := "Shared context accessible"

	// Try to perform a read operation
	defer func() {
		if r := recover(); r != nil {
			status = "unhealthy"
			message = fmt.Sprintf("Shared context panic: %v", r)
		}
	}()

	_ = ctx.GetAll()

	check := HealthCheck{
		Component: "shared_context",
		Status:    status,
		Message:   message,
		Timestamp: time.Now(),
	}

	hc.mu.Lock()
	hc.checks = append(hc.checks, check)
	hc.mu.Unlock()

	return check
}

// CheckTeamManagerResponsiveness checks team manager responsiveness
func (hc *HealthChecker) CheckTeamManagerResponsiveness(tm *TeamManager) HealthCheck {
	status := "healthy"
	message := "Team manager responsive"

	// Try to get team status
	start := time.Now()
	_, _ = tm.GetTeamStatus("test")
	duration := time.Since(start)

	if duration > 100*time.Millisecond {
		status = "degraded"
		message = fmt.Sprintf("Team manager slow: %v", duration)
	}

	check := HealthCheck{
		Component: "team_manager",
		Status:    status,
		Message:   message,
		Timestamp: time.Now(),
	}

	hc.mu.Lock()
	hc.checks = append(hc.checks, check)
	hc.mu.Unlock()

	return check
}

// GetHealthChecks returns all health checks
func (hc *HealthChecker) GetHealthChecks() []HealthCheck {
	hc.mu.RLock()
	defer hc.mu.RUnlock()

	result := make([]HealthCheck, len(hc.checks))
	copy(result, hc.checks)
	return result
}
