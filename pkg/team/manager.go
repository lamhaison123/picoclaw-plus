// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package team

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/sipeed/picoclaw/pkg/agent"
	"github.com/sipeed/picoclaw/pkg/bus"
	"github.com/sipeed/picoclaw/pkg/config"
	"github.com/sipeed/picoclaw/pkg/logger"
	"github.com/sipeed/picoclaw/pkg/providers"
	"github.com/sipeed/picoclaw/pkg/team/memory"
)

// Error definitions
var (
	ErrTeamNotFound  = errors.New("team not found")
	ErrAgentNotFound = errors.New("agent not found")
)

// TeamManager orchestrates team lifecycle and coordination
type TeamManager struct {
	registry         *agent.AgentRegistry
	bus              *bus.MessageBus
	teams            map[string]*Team
	roleCapabilities map[string][]string
	teamMemory       *memory.TeamMemory
	metrics          *MetricsCollector
	workspace        string        // Workspace path for persistence
	agentExecutor    AgentExecutor // Executor for running tasks with agents
	// Performance optimizations
	agentPool *AgentPool
	roleCache *RoleCache
	// Model selection support
	provider providers.LLMProvider
	cfg      *config.Config
	mu       sync.RWMutex
}

// AgentPool manages reusable agent instances
type AgentPool struct {
	instances map[string]*agent.AgentInstance
	usage     map[string]int
	mu        sync.RWMutex
}

// RoleCache caches role-to-capability mappings
type RoleCache struct {
	mappings map[string][]string
	mu       sync.RWMutex
}

// NewTeamManager creates a new TeamManager instance
func NewTeamManager(registry *agent.AgentRegistry, msgBus *bus.MessageBus) *TeamManager {
	tm := &TeamManager{
		registry:         registry,
		bus:              msgBus,
		teams:            make(map[string]*Team),
		roleCapabilities: make(map[string][]string),
		teamMemory:       nil, // Set via SetTeamMemory
		workspace:        "",  // Set via SetTeamMemory
		metrics:          NewMetricsCollector(),
		agentPool:        NewAgentPool(),
		roleCache:        NewRoleCache(),
		agentExecutor:    nil, // Set via SetAgentExecutor
	}

	// Note: Teams will be loaded when SetTeamMemory is called

	return tm
}

// SetTeamMemory sets the team memory manager
func (tm *TeamManager) SetTeamMemory(teamMemory *memory.TeamMemory) {
	tm.teamMemory = teamMemory
	if teamMemory != nil {
		tm.workspace = teamMemory.GetWorkspace()
		// Load persisted teams now that we have the workspace path
		if err := tm.loadPersistedTeams(); err != nil {
			logger.WarnCF("team", "Failed to load persisted teams after setting memory",
				map[string]any{
					"error": err.Error(),
				})
		}
	}
}

// SetAgentExecutor sets the executor for running tasks with agents
func (tm *TeamManager) SetAgentExecutor(executor AgentExecutor) {
	tm.agentExecutor = executor
}

// SetProvider sets the LLM provider and config for model selection
func (tm *TeamManager) SetProvider(provider providers.LLMProvider, cfg *config.Config) {
	tm.provider = provider
	tm.cfg = cfg
}

// GetMetrics returns the metrics collector
func (tm *TeamManager) GetMetrics() *MetricsCollector {
	return tm.metrics
}

// CreateTeam creates a new team from configuration
// CreateTeam creates a new team from configuration
func (tm *TeamManager) CreateTeam(ctx context.Context, teamConfig *TeamConfig) (*Team, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Validate configuration
	if err := tm.validateConfig(teamConfig); err != nil {
		return nil, fmt.Errorf("invalid team configuration: %w", err)
	}

	// Check if team ID already exists
	if _, exists := tm.teams[teamConfig.TeamID]; exists {
		return nil, fmt.Errorf("team with ID '%s' already exists", teamConfig.TeamID)
	}

	// Create team instance
	team := &Team{
		ID:            teamConfig.TeamID,
		Name:          teamConfig.Name,
		Pattern:       CollaborationPattern(teamConfig.Pattern),
		Agents:        make(map[string]*TeamAgent),
		SharedContext: NewSharedContext(teamConfig.TeamID),
		Status:        TeamStatusInitializing,
		CreatedAt:     time.Now(),
		Config:        teamConfig,
	}

	// Build role-to-capabilities mapping (team-specific, stored in team config)
	// Note: roleCapabilities is shared across teams, so we use composite keys
	// Lock is already held by CreateTeam, so this is safe from race conditions
	for _, roleConfig := range teamConfig.Roles {
		compositeKey := fmt.Sprintf("%s:%s", teamConfig.TeamID, roleConfig.Name)
		tm.roleCapabilities[compositeKey] = roleConfig.Capabilities
	}

	// Register agents for each role
	for _, roleConfig := range teamConfig.Roles {
		agentID := fmt.Sprintf("%s-%s", teamConfig.TeamID, roleConfig.Name)

		teamAgent := &TeamAgent{
			AgentID:      agentID,
			Role:         roleConfig.Name,
			Capabilities: roleConfig.Capabilities,
			Status:       StatusIdle,
			FailureCount: 0,
			LastActive:   time.Now(),
		}

		team.Agents[agentID] = teamAgent

		// Register agent instance with correct model if provider is available
		if tm.provider != nil && tm.cfg != nil {
			// Create agent config for this role
			agentCfg := &config.AgentConfig{
				ID:   agentID,
				Name: fmt.Sprintf("%s (%s)", team.Name, roleConfig.Name),
				Model: &config.AgentModelConfig{
					Primary: roleConfig.Model,
				},
				Workspace: fmt.Sprintf("%s/teams/%s/%s", tm.workspace, teamConfig.TeamID, roleConfig.Name),
			}

			// Create agent instance with role-specific model
			instance := agent.NewAgentInstance(agentCfg, &tm.cfg.Agents.Defaults, tm.cfg, tm.provider)

			// Register in agent registry
			tm.registry.RegisterTeamAgent(agentID, instance)

			logger.InfoCF("team", "Registered team agent with model",
				map[string]any{
					"team_id":   team.ID,
					"agent_id":  agentID,
					"role":      roleConfig.Name,
					"model":     roleConfig.Model,
					"workspace": instance.Workspace,
				})
		} else {
			// Fallback: log warning that model selection is not available
			logger.WarnCF("team", "Model selection not available, agent will use default model",
				map[string]any{
					"team_id":  team.ID,
					"agent_id": agentID,
					"role":     roleConfig.Name,
					"model":    roleConfig.Model,
				})
		}
	}

	// Set coordinator if specified
	if teamConfig.Coordinator.AgentID != "" {
		team.CoordinatorID = teamConfig.Coordinator.AgentID
	} else if teamConfig.Coordinator.Role != "" {
		// Find first agent with coordinator role
		for agentID, agent := range team.Agents {
			if agent.Role == teamConfig.Coordinator.Role {
				team.CoordinatorID = agentID
				break
			}
		}
	}

	// Validate that all required roles are assigned
	if len(team.Agents) == 0 {
		return nil, fmt.Errorf("team must have at least one agent")
	}

	// Set team status to active
	team.Status = TeamStatusActive

	// Store team
	tm.teams[team.ID] = team

	// Record metrics
	tm.metrics.RecordTeamCreation(team.ID)

	logger.InfoCF("team", "Team created successfully",
		map[string]any{
			"team_id":     team.ID,
			"team_name":   team.Name,
			"pattern":     team.Pattern,
			"agent_count": len(team.Agents),
			"coordinator": team.CoordinatorID,
		})

	// Persist team state to disk
	if err := tm.SaveTeamState(team); err != nil {
		// Rollback: remove team from map on persistence failure
		// Note: Already holding tm.mu.Lock() from line 116, no need to lock again
		delete(tm.teams, team.ID)

		logger.ErrorCF("team", "Failed to save team state, rolling back team creation",
			map[string]any{
				"team_id": team.ID,
				"error":   err.Error(),
			})
		return nil, fmt.Errorf("failed to persist team state: %w", err)
	}

	return team, nil
}
func (tm *TeamManager) DissolveTeam(ctx context.Context, teamID string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	team, exists := tm.teams[teamID]
	if !exists {
		return fmt.Errorf("team '%s' not found", teamID)
	}

	// Persist team memory before dissolution
	if tm.teamMemory != nil {
		record := tm.createRecordFromTeam(team, "dissolved")
		if err := tm.teamMemory.SaveTeamRecord(record); err != nil {
			logger.WarnCF("team", "Failed to save team memory",
				map[string]any{
					"team_id": teamID,
					"error":   err.Error(),
				})
		}
	}

	// Cancel any pending tasks by marking them as cancelled
	// Note: We can't directly cancel running tasks, but we mark them
	allData := team.SharedContext.GetAll()
	for _, entry := range allData {
		if taskMap, ok := entry.Value.(map[string]any); ok {
			if status, ok := taskMap["status"].(string); ok {
				if status == string(TaskStatusPending) || status == string(TaskStatusInProgress) {
					taskMap["status"] = string(TaskStatusCancelled)
					team.SharedContext.Set(entry.Key, taskMap, "system")
				}
			}
		}
	}

	// Clear shared context to free memory
	if team.SharedContext != nil {
		team.SharedContext = nil
	}

	// Note: Coordinators are created per-task in ExecuteTask with defer Shutdown(),
	// so no need to shutdown here. They are not stored in the Team struct.

	// Deregister all agents from Agent Registry
	for agentID := range team.Agents {
		// Unregister from registry
		if tm.registry != nil {
			tm.registry.UnregisterAgent(agentID)
		}

		logger.DebugCF("team", "Deregistering team agent",
			map[string]any{
				"team_id":  teamID,
				"agent_id": agentID,
			})
	}

	// Update team status
	team.Status = TeamStatusDissolved

	// Remove from teams map
	delete(tm.teams, teamID)

	// Delete persisted state
	if err := tm.DeleteTeamState(teamID); err != nil {
		logger.WarnCF("team", "Failed to delete team state",
			map[string]any{
				"team_id": teamID,
				"error":   err.Error(),
			})
	}

	logger.InfoCF("team", "Team dissolved",
		map[string]any{
			"team_id":   teamID,
			"team_name": team.Name,
		})

	return nil
}

// AddAgent adds a new agent to an existing team
func (tm *TeamManager) AddAgent(ctx context.Context, teamID string, agentConfig AgentConfig) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	team, exists := tm.teams[teamID]
	if !exists {
		return fmt.Errorf("team '%s' not found", teamID)
	}

	// Validate team is not dissolved
	if team.Status == TeamStatusDissolved {
		return fmt.Errorf("cannot add agent to dissolved team '%s'", teamID)
	}

	// Validate role exists in team configuration
	roleExists := false
	var capabilities []string
	for _, roleConfig := range team.Config.Roles {
		if roleConfig.Name == agentConfig.Role {
			roleExists = true
			capabilities = roleConfig.Capabilities
			break
		}
	}

	if !roleExists {
		return fmt.Errorf("role '%s' not defined in team configuration", agentConfig.Role)
	}

	// Create team agent
	agentID := agentConfig.AgentID
	if agentID == "" {
		agentID = fmt.Sprintf("%s-%s-%d", teamID, agentConfig.Role, len(team.Agents))
	}

	teamAgent := &TeamAgent{
		AgentID:      agentID,
		Role:         agentConfig.Role,
		Capabilities: capabilities,
		Status:       StatusIdle,
		FailureCount: 0,
		LastActive:   time.Now(),
	}

	// Add to team
	team.Agents[agentID] = teamAgent

	// Grant shared context access (already accessible by design)

	// Notify team members of composition change
	team.SharedContext.AddHistoryEntry("system", "agent_added",
		map[string]string{
			"agent_id": agentID,
			"role":     agentConfig.Role,
		})

	logger.InfoCF("team", "Agent added to team",
		map[string]any{
			"team_id":  teamID,
			"agent_id": agentID,
			"role":     agentConfig.Role,
		})

	return nil
}

// RemoveAgent removes an agent from a team
func (tm *TeamManager) RemoveAgent(ctx context.Context, teamID, agentID string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	team, exists := tm.teams[teamID]
	if !exists {
		return fmt.Errorf("team '%s' not found", teamID)
	}

	agent, exists := team.Agents[agentID]
	if !exists {
		return fmt.Errorf("agent '%s' not found in team '%s'", agentID, teamID)
	}

	// Reassign any active tasks from this agent
	if err := tm.reassignAgentTasks(ctx, team, agentID, agent.Role); err != nil {
		logger.WarnCF("team", "Failed to reassign tasks during agent removal",
			map[string]any{
				"team_id":  teamID,
				"agent_id": agentID,
				"error":    err.Error(),
			})
	}

	// Remove from team
	delete(team.Agents, agentID)

	// Notify team members of composition change
	team.SharedContext.AddHistoryEntry("system", "agent_removed",
		map[string]string{
			"agent_id": agentID,
			"role":     agent.Role,
		})

	logger.InfoCF("team", "Agent removed from team",
		map[string]any{
			"team_id":  teamID,
			"agent_id": agentID,
			"role":     agent.Role,
		})

	return nil
}

// reassignAgentTasks reassigns active tasks from a removed agent to other agents with the same role
func (tm *TeamManager) reassignAgentTasks(ctx context.Context, team *Team, removedAgentID, role string) error {
	// Find tasks assigned to this agent in shared context
	allData := team.SharedContext.GetAll()

	tasksToReassign := []string{}
	for _, entry := range allData {
		// Look for task assignments
		if taskMap, ok := entry.Value.(map[string]any); ok {
			if assignedAgent, ok := taskMap["agent_id"].(string); ok && assignedAgent == removedAgentID {
				if status, ok := taskMap["status"].(string); ok {
					if status == string(TaskStatusAssigned) || status == string(TaskStatusInProgress) {
						if taskID, ok := taskMap["task_id"].(string); ok {
							tasksToReassign = append(tasksToReassign, taskID)
						}
					}
				}
			}
		}
	}

	if len(tasksToReassign) == 0 {
		return nil
	}

	// Find another agent with the same role
	var targetAgentID string
	for agentID, agent := range team.Agents {
		if agentID != removedAgentID && agent.Role == role && agent.Status == StatusIdle {
			targetAgentID = agentID
			break
		}
	}

	if targetAgentID == "" {
		return fmt.Errorf("no available agent with role '%s' to reassign tasks", role)
	}

	// Record reassignments in shared context
	for _, taskID := range tasksToReassign {
		team.SharedContext.AddHistoryEntry("system", "task_reassigned",
			map[string]string{
				"task_id":    taskID,
				"from_agent": removedAgentID,
				"to_agent":   targetAgentID,
				"reason":     "agent_removed",
			})
	}

	logger.InfoCF("team", "Reassigned tasks from removed agent",
		map[string]any{
			"team_id":    team.ID,
			"from_agent": removedAgentID,
			"to_agent":   targetAgentID,
			"task_count": len(tasksToReassign),
		})

	return nil
}

// GetTeamStatus returns the current status of a team
func (tm *TeamManager) GetTeamStatus(teamID string) (*TeamStatusInfo, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	team, exists := tm.teams[teamID]
	if !exists {
		return nil, fmt.Errorf("team '%s' not found", teamID)
	}

	// Build agent status map
	agentStatuses := make(map[string]AgentStatus, len(team.Agents))
	for agentID, agent := range team.Agents {
		agentStatuses[agentID] = agent.Status
	}

	statusInfo := &TeamStatusInfo{
		TeamID:        team.ID,
		TeamName:      team.Name,
		Status:        team.Status,
		Pattern:       team.Pattern,
		AgentCount:    len(team.Agents),
		AgentStatuses: agentStatuses,
		CoordinatorID: team.CoordinatorID,
		CreatedAt:     team.CreatedAt,
		Uptime:        time.Since(team.CreatedAt),
	}

	return statusInfo, nil
}

// ValidateToolAccess checks if an agent is allowed to use a tool
func (tm *TeamManager) ValidateToolAccess(agentID, toolName string) bool {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	// Find the team and agent
	var teamAgent *TeamAgent
	for _, team := range tm.teams {
		if agent, exists := team.Agents[agentID]; exists {
			teamAgent = agent
			break
		}
	}

	if teamAgent == nil {
		return false
	}

	// Find role configuration
	for _, team := range tm.teams {
		for _, agent := range team.Agents {
			if agent.AgentID == agentID {
				for _, roleConfig := range team.Config.Roles {
					if roleConfig.Name == agent.Role {
						return tm.checkToolPermission(toolName, roleConfig.Tools)
					}
				}
			}
		}
	}

	return false
}

// checkToolPermission checks if a tool matches any of the allowed patterns
func (tm *TeamManager) checkToolPermission(toolName string, allowedTools []string) bool {
	for _, pattern := range allowedTools {
		if pattern == "*" {
			return true
		}
		if pattern == toolName {
			return true
		}
		// Check wildcard patterns (e.g., "file_*")
		if len(pattern) > 0 && pattern[len(pattern)-1] == '*' {
			prefix := pattern[:len(pattern)-1]
			if len(toolName) >= len(prefix) && toolName[:len(prefix)] == prefix {
				return true
			}
		}
	}
	return false
}

// validateConfig validates team configuration
func (tm *TeamManager) validateConfig(teamConfig *TeamConfig) error {
	if teamConfig.TeamID == "" {
		return fmt.Errorf("team_id is required")
	}
	if teamConfig.Name == "" {
		return fmt.Errorf("name is required")
	}
	if teamConfig.Pattern == "" {
		return fmt.Errorf("pattern is required")
	}

	validPatterns := map[string]bool{
		"sequential":   true,
		"parallel":     true,
		"hierarchical": true,
	}
	if !validPatterns[teamConfig.Pattern] {
		return fmt.Errorf("invalid pattern '%s', must be one of: sequential, parallel, hierarchical", teamConfig.Pattern)
	}

	if len(teamConfig.Roles) == 0 {
		return fmt.Errorf("at least one role is required")
	}

	// Validate each role
	for _, role := range teamConfig.Roles {
		if role.Name == "" {
			return fmt.Errorf("role name is required")
		}
		if len(role.Capabilities) == 0 {
			return fmt.Errorf("role '%s' must have at least one capability", role.Name)
		}
	}

	return nil
}

// AgentConfig represents configuration for adding an agent
type AgentConfig struct {
	AgentID string
	Role    string
}

// GetTeam returns a team by ID
func (tm *TeamManager) GetTeam(teamID string) (*Team, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	team, exists := tm.teams[teamID]
	if !exists {
		return nil, fmt.Errorf("team '%s' not found", teamID)
	}

	return team, nil
}

// TeamStatusInfo contains detailed team status information
type TeamStatusInfo struct {
	TeamID        string
	TeamName      string
	Status        TeamStatus
	Pattern       CollaborationPattern
	AgentCount    int
	AgentStatuses map[string]AgentStatus
	CoordinatorID string
	CreatedAt     time.Time
	Uptime        time.Duration
}

// UpdateAgentStatus updates the status of an agent in a team
func (tm *TeamManager) UpdateAgentStatus(teamID, agentID string, status AgentStatus) error {
	tm.mu.RLock()
	team, exists := tm.teams[teamID]
	tm.mu.RUnlock()

	if !exists {
		return ErrTeamNotFound
	}

	agent, exists := team.Agents[agentID]
	if !exists {
		return ErrAgentNotFound
	}

	agent.Status = status
	agent.LastActive = time.Now()

	return nil
}

// RecordTaskMetrics records metrics for a completed task
func (tm *TeamManager) RecordTaskMetrics(teamID string, task *Task) error {
	tm.mu.RLock()
	team, exists := tm.teams[teamID]
	tm.mu.RUnlock()

	if !exists {
		return ErrTeamNotFound
	}

	// Store metrics in shared context
	metricsKey := fmt.Sprintf("metrics_task_%s", task.ID)
	metrics := map[string]any{
		"task_id":      task.ID,
		"agent_id":     task.AssignedAgentID,
		"status":       task.Status,
		"created_at":   task.CreatedAt,
		"started_at":   task.StartedAt,
		"completed_at": task.CompletedAt,
	}

	if task.StartedAt != nil && task.CompletedAt != nil {
		duration := task.CompletedAt.Sub(*task.StartedAt)
		metrics["duration_ms"] = duration.Milliseconds()
	}

	team.SharedContext.Set(metricsKey, metrics, "system")

	return nil
}

// DetectUnresponsiveAgent checks if an agent has exceeded the heartbeat timeout
func (tm *TeamManager) DetectUnresponsiveAgent(teamID, agentID string, timeout time.Duration) (bool, error) {
	tm.mu.RLock()
	team, exists := tm.teams[teamID]
	tm.mu.RUnlock()

	if !exists {
		return false, ErrTeamNotFound
	}

	agent, exists := team.Agents[agentID]

	if !exists {
		return false, ErrAgentNotFound
	}

	timeSinceActive := time.Since(agent.LastActive)
	if timeSinceActive > timeout {
		// Mark agent as unresponsive
		tm.UpdateAgentStatus(teamID, agentID, StatusUnresponsive)
		return true, nil
	}

	return false, nil
}

// IncrementAgentFailure increments the failure count for an agent
func (tm *TeamManager) IncrementAgentFailure(teamID, agentID string, threshold int) (bool, error) {
	tm.mu.RLock()
	team, exists := tm.teams[teamID]
	tm.mu.RUnlock()

	if !exists {
		return false, ErrTeamNotFound
	}

	agent, exists := team.Agents[agentID]
	if !exists {
		return false, ErrAgentNotFound
	}

	agent.FailureCount++

	// Check if threshold exceeded
	if agent.FailureCount >= threshold {
		// Remove agent from team
		delete(team.Agents, agentID)
		return true, nil
	}

	return false, nil
}

// SendTaskDelegation sends a task delegation message via Message Bus
func (tm *TeamManager) SendTaskDelegation(teamID string, msg *TaskDelegationMessage) error {
	data, err := msg.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to serialize task delegation message: %w", err)
	}

	outboundMsg := bus.OutboundMessage{
		Content: string(data),
	}
	if err := tm.bus.PublishOutbound(context.Background(), outboundMsg); err != nil {
		return fmt.Errorf("failed to publish task delegation: %w", err)
	}

	logger.DebugCF("team", "Task delegation sent",
		map[string]any{
			"team_id":    teamID,
			"task_id":    msg.Task.ID,
			"from_agent": msg.FromAgentID,
			"to_agent":   msg.ToAgentID,
		})

	return nil
}

// SendTaskResult sends a task result message via Message Bus
func (tm *TeamManager) SendTaskResult(teamID string, msg *TaskResultMessage) error {
	data, err := msg.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to serialize task result message: %w", err)
	}

	outboundMsg := bus.OutboundMessage{
		Content: string(data),
	}
	if err := tm.bus.PublishOutbound(context.Background(), outboundMsg); err != nil {
		return fmt.Errorf("failed to publish task result: %w", err)
	}

	logger.DebugCF("team", "Task result sent",
		map[string]any{
			"team_id":    teamID,
			"task_id":    msg.TaskID,
			"from_agent": msg.FromAgentID,
			"to_agent":   msg.ToAgentID,
			"status":     msg.Status,
		})

	return nil
}

// SendConsensusRequest sends a consensus request message via Message Bus
func (tm *TeamManager) SendConsensusRequest(teamID string, msg *ConsensusRequestMessage) error {
	data, err := msg.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to serialize consensus request message: %w", err)
	}

	outboundMsg := bus.OutboundMessage{
		Content: string(data),
	}
	if err := tm.bus.PublishOutbound(context.Background(), outboundMsg); err != nil {
		return fmt.Errorf("failed to publish consensus request: %w", err)
	}

	logger.DebugCF("team", "Consensus request sent",
		map[string]any{
			"team_id":      teamID,
			"consensus_id": msg.ConsensusID,
			"from_agent":   msg.FromAgentID,
			"question":     msg.Question,
		})

	return nil
}

// SubscribeToTeamMessages subscribes to team-specific message channels
// Note: MessageBus uses channel-based communication, not pub/sub pattern
func (tm *TeamManager) SubscribeToTeamMessages(teamID string, handler func(channel string, data []byte)) error {
	// TODO: Implement when MessageBus supports pub/sub pattern
	logger.InfoCF("team", "Team message subscription requested (not yet implemented)",
		map[string]any{
			"team_id": teamID,
		})

	return nil
}

// UnsubscribeFromTeamMessages unsubscribes from team-specific message channels
// Note: MessageBus uses channel-based communication, not pub/sub pattern
func (tm *TeamManager) UnsubscribeFromTeamMessages(teamID string) error {
	// TODO: Implement when MessageBus supports pub/sub pattern
	logger.InfoCF("team", "Team message unsubscription requested (not yet implemented)",
		map[string]any{
			"team_id": teamID,
		})

	return nil
}

// NewAgentPool creates a new agent pool
func NewAgentPool() *AgentPool {
	return &AgentPool{
		instances: make(map[string]*agent.AgentInstance),
		usage:     make(map[string]int),
	}
}

// GetOrCreateInstance gets an existing agent instance or creates a new one
func (ap *AgentPool) GetOrCreateInstance(role string, createFn func() *agent.AgentInstance) *agent.AgentInstance {
	ap.mu.Lock()
	defer ap.mu.Unlock()

	// Check if instance exists for this role
	if instance, exists := ap.instances[role]; exists {
		ap.usage[role]++
		return instance
	}

	// Create new instance
	instance := createFn()
	ap.instances[role] = instance
	ap.usage[role] = 1
	return instance
}

// ReleaseInstance decrements usage count for an agent instance
func (ap *AgentPool) ReleaseInstance(role string) {
	ap.mu.Lock()
	defer ap.mu.Unlock()

	if count, exists := ap.usage[role]; exists {
		if count > 1 {
			ap.usage[role]--
		} else {
			// Remove instance if no longer in use
			delete(ap.instances, role)
			delete(ap.usage, role)
		}
	}
}

// GetUsageCount returns the usage count for a role
func (ap *AgentPool) GetUsageCount(role string) int {
	ap.mu.RLock()
	defer ap.mu.RUnlock()
	return ap.usage[role]
}

// Clear removes all instances from the pool
func (ap *AgentPool) Clear() {
	ap.mu.Lock()
	defer ap.mu.Unlock()
	ap.instances = make(map[string]*agent.AgentInstance)
	ap.usage = make(map[string]int)
}

// NewRoleCache creates a new role cache
func NewRoleCache() *RoleCache {
	return &RoleCache{
		mappings: make(map[string][]string),
	}
}

// Get retrieves capabilities for a role from cache
func (rc *RoleCache) Get(role string) ([]string, bool) {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	caps, exists := rc.mappings[role]
	return caps, exists
}

// Set stores capabilities for a role in cache
func (rc *RoleCache) Set(role string, capabilities []string) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.mappings[role] = capabilities
}

// Invalidate removes a role from cache
func (rc *RoleCache) Invalidate(role string) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	delete(rc.mappings, role)
}

// Clear removes all entries from cache
func (rc *RoleCache) Clear() {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.mappings = make(map[string][]string)
}

// GetCapabilitiesForRole retrieves capabilities for a role with caching
// Note: This method searches across all teams for the role
func (tm *TeamManager) GetCapabilitiesForRole(role string) ([]string, bool) {
	// Try cache first
	if caps, exists := tm.roleCache.Get(role); exists {
		return caps, true
	}

	// Search roleCapabilities map (may have composite keys teamID:role or just role)
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	// Try direct lookup first (backward compatibility)
	if caps, exists := tm.roleCapabilities[role]; exists {
		tm.roleCache.Set(role, caps)
		return caps, true
	}

	// Try composite key lookup (search all teams)
	for key, caps := range tm.roleCapabilities {
		// Check if key ends with :role
		if len(key) > len(role) && key[len(key)-len(role)-1:] == ":"+role {
			tm.roleCache.Set(role, caps)
			return caps, true
		}
	}

	return nil, false
}

// InvalidateRoleCache invalidates the role cache for a specific role
func (tm *TeamManager) InvalidateRoleCache(role string) {
	tm.roleCache.Invalidate(role)
}

// createRecordFromTeam creates a memory record from a team
func (tm *TeamManager) createRecordFromTeam(t *Team, outcome string) *memory.TeamMemoryRecord {
	record := &memory.TeamMemoryRecord{
		TeamID:        t.ID,
		TeamName:      t.Name,
		Pattern:       string(t.Pattern),
		StartTime:     t.CreatedAt,
		EndTime:       time.Now(),
		SharedContext: t.SharedContext.Snapshot(),
		Tasks:         []memory.TaskRecord{},
		Consensus:     []memory.ConsensusRecord{},
		Outcome:       outcome,
	}

	return record
}

// GetAllTeams returns all teams
func (tm *TeamManager) GetAllTeams() map[string]*Team {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	// Return a copy to prevent external modification
	teams := make(map[string]*Team, len(tm.teams))
	for id, team := range tm.teams {
		teams[id] = team
	}
	return teams
}

// ExecuteTask executes a task using a team
func (tm *TeamManager) ExecuteTask(ctx context.Context, teamID string, taskDescription string) (any, error) {
	tm.mu.RLock()
	team, exists := tm.teams[teamID]
	tm.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("team '%s' not found", teamID)
	}

	// Use intelligent router to determine best role for the task
	var requiredRole string

	if tm.agentExecutor != nil {
		// Use LLM-based intelligent routing
		router := NewIntelligentRouter(tm.agentExecutor)
		router.SetTeamID(teamID) // Set team ID for proper agent ID construction
		determinedRole, err := router.DetermineRole(ctx, taskDescription, team.Config.Roles)
		if err != nil {
			logger.WarnCF("team", "Failed to determine role intelligently, using simple routing",
				map[string]any{
					"error": err.Error(),
				})
			// Fallback to simple routing
			requiredRole = router.DetermineRoleSimple(taskDescription, team.Config.Roles)
		} else {
			requiredRole = determinedRole
		}
	} else {
		// Fallback to simple keyword-based routing
		router := NewIntelligentRouter(nil)
		requiredRole = router.DetermineRoleSimple(taskDescription, team.Config.Roles)
	}

	if requiredRole == "" {
		return nil, fmt.Errorf("no suitable role found for task execution: team has no roles configured")
	}

	// Validate role exists in team
	roleExists := false
	for _, roleConfig := range team.Config.Roles {
		if roleConfig.Name == requiredRole {
			roleExists = true
			break
		}
	}
	if !roleExists {
		return nil, fmt.Errorf("selected role '%s' not found in team configuration", requiredRole)
	}

	// Create coordinator
	router := NewDelegationRouter(team.Config.Settings.MaxDelegationDepth)
	coordinator := NewCoordinatorAgent(
		fmt.Sprintf("coordinator-%s", teamID),
		teamID,
		team,
		team.Pattern,
		tm.bus,
		router,
		ctx, // Pass context for goroutine lifecycle
	)

	// BUG FIX: Check if coordinator creation failed
	if coordinator == nil {
		return nil, fmt.Errorf("failed to create coordinator for team '%s'", teamID)
	}

	defer coordinator.Shutdown() // Ensure cleanup of background goroutines

	// BUG FIX #4: Check if agentExecutor is nil before setting
	if tm.agentExecutor == nil {
		return nil, fmt.Errorf("agent executor not configured for team manager")
	}

	// Set executor
	coordinator.Executor = tm.agentExecutor

	// Set timeout from config
	if team.Config.Settings.AgentTimeoutSeconds > 0 {
		coordinator.SetTaskTimeout(time.Duration(team.Config.Settings.AgentTimeoutSeconds) * time.Second)
	}

	// Create task with auto-selected role
	task := NewTask(taskDescription, requiredRole, nil)

	logger.InfoCF("team", "Executing task",
		map[string]any{
			"team_id":       teamID,
			"task_id":       task.ID,
			"required_role": requiredRole,
			"pattern":       team.Pattern,
		})

	// Execute task based on pattern
	var result any
	var err error

	switch team.Pattern {
	case PatternSequential:
		result, err = coordinator.ExecuteSequential(ctx, []*Task{task})
	case PatternParallel:
		result, err = coordinator.ExecuteParallel(ctx, []*Task{task})
	case PatternHierarchical:
		result, err = coordinator.ExecuteHierarchical(ctx, task)
	default:
		return nil, fmt.Errorf("unsupported collaboration pattern: %s", team.Pattern)
	}

	if err != nil {
		logger.ErrorCF("team", "Task execution failed",
			map[string]any{
				"team_id": teamID,
				"task_id": task.ID,
				"error":   err.Error(),
			})
		return nil, fmt.Errorf("task execution failed: %w", err)
	}

	logger.InfoCF("team", "Task completed successfully",
		map[string]any{
			"team_id": teamID,
			"task_id": task.ID,
		})

	return result, nil
}

// ExecuteTaskWithRole executes a task using a team with a specific role
func (tm *TeamManager) ExecuteTaskWithRole(ctx context.Context, teamID string, taskDescription string, requiredRole string) (any, error) {
	tm.mu.RLock()
	team, exists := tm.teams[teamID]
	tm.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("team '%s' not found", teamID)
	}

	// Validate role exists in team
	roleExists := false
	for _, roleConfig := range team.Config.Roles {
		if roleConfig.Name == requiredRole {
			roleExists = true
			break
		}
	}
	if !roleExists {
		return nil, fmt.Errorf("role '%s' not found in team configuration", requiredRole)
	}

	// Create coordinator
	router := NewDelegationRouter(team.Config.Settings.MaxDelegationDepth)
	coordinator := NewCoordinatorAgent(
		fmt.Sprintf("coordinator-%s", teamID),
		teamID,
		team,
		team.Pattern,
		tm.bus,
		router,
		ctx,
	)

	// BUG FIX: Check if coordinator creation failed
	if coordinator == nil {
		return nil, fmt.Errorf("failed to create coordinator for team '%s'", teamID)
	}

	defer coordinator.Shutdown() // Ensure cleanup of background goroutines

	// BUG FIX #4: Check if agentExecutor is nil before setting
	if tm.agentExecutor == nil {
		return nil, fmt.Errorf("agent executor not configured for team manager")
	}

	// Set executor
	coordinator.Executor = tm.agentExecutor

	// Set timeout from config
	if team.Config.Settings.AgentTimeoutSeconds > 0 {
		coordinator.SetTaskTimeout(time.Duration(team.Config.Settings.AgentTimeoutSeconds) * time.Second)
	}

	// Create task with specified role
	task := NewTask(taskDescription, requiredRole, nil)

	logger.InfoCF("team", "Executing task with specific role",
		map[string]any{
			"team_id":       teamID,
			"task_id":       task.ID,
			"required_role": requiredRole,
			"pattern":       team.Pattern,
		})

	// Execute task based on pattern
	var result any
	var err error

	switch team.Pattern {
	case PatternSequential:
		result, err = coordinator.ExecuteSequential(ctx, []*Task{task})
	case PatternParallel:
		result, err = coordinator.ExecuteParallel(ctx, []*Task{task})
	case PatternHierarchical:
		result, err = coordinator.ExecuteHierarchical(ctx, task)
	default:
		return nil, fmt.Errorf("unsupported collaboration pattern: %s", team.Pattern)
	}

	if err != nil {
		logger.ErrorCF("team", "Task execution failed",
			map[string]any{
				"team_id": teamID,
				"task_id": task.ID,
				"error":   err.Error(),
			})
		return nil, fmt.Errorf("task execution failed: %w", err)
	}

	logger.InfoCF("team", "Task completed successfully",
		map[string]any{
			"team_id": teamID,
			"task_id": task.ID,
		})

	return result, nil
}

// GetAllTeamsAsInterface returns all teams as map[string]interface{} for tool compatibility
func (tm *TeamManager) GetAllTeamsAsInterface() map[string]interface{} {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	result := make(map[string]interface{}, len(tm.teams))
	for id, team := range tm.teams {
		result[id] = team
	}
	return result
}

// GetTeamStatusAsInterface returns team status as interface{} for tool compatibility
func (tm *TeamManager) GetTeamStatusAsInterface(teamID string) (interface{}, error) {
	return tm.GetTeamStatus(teamID)
}
