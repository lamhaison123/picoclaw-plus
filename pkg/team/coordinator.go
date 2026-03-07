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

	"github.com/sipeed/picoclaw/pkg/bus"
	"github.com/sipeed/picoclaw/pkg/logger"
)

// CoordinatorAgent coordinates task delegation and execution within a team
type CoordinatorAgent struct {
	AgentID          string
	TeamID           string
	Team             *Team
	Pattern          CollaborationPattern
	Bus              *bus.MessageBus
	Router           *DelegationRouter
	ConsensusManager *ConsensusManager
	Executor         AgentExecutor // Executor for running tasks
	pendingTasks     map[string]*Task
	taskResults      map[string]*TaskResultMessage
	mu               sync.RWMutex
	taskTimeout      time.Duration
	shutdownCtx      context.Context    // Context for shutdown
	shutdownCancel   context.CancelFunc // Cancel function for shutdown
}

// NewCoordinatorAgent creates a new coordinator agent for a team
func NewCoordinatorAgent(agentID, teamID string, team *Team, pattern CollaborationPattern, bus *bus.MessageBus, delegationRouter *DelegationRouter, ctx context.Context) *CoordinatorAgent {
	// Create shutdown context that can be cancelled independently
	shutdownCtx, shutdownCancel := context.WithCancel(context.Background())

	c := &CoordinatorAgent{
		AgentID:          agentID,
		TeamID:           teamID,
		Team:             team,
		Pattern:          pattern,
		Bus:              bus,
		Router:           delegationRouter,
		ConsensusManager: NewConsensusManager(),
		pendingTasks:     make(map[string]*Task),
		taskResults:      make(map[string]*TaskResultMessage),
		taskTimeout:      270 * time.Second, // Default 270 seconds (4.5 min) - increased per architecture review to prevent parent timeout before child operations complete
		shutdownCtx:      shutdownCtx,
		shutdownCancel:   shutdownCancel,
	}

	// Start cleanup goroutine with shutdown context
	go c.cleanupTaskResults(shutdownCtx)

	return c
}

// Shutdown gracefully shuts down the coordinator and stops all background goroutines
func (c *CoordinatorAgent) Shutdown() {
	if c.shutdownCancel != nil {
		c.shutdownCancel()
	}
}

// cleanupTaskResults periodically removes old task results to prevent memory leak
func (c *CoordinatorAgent) cleanupTaskResults(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			// Context cancelled, exit goroutine
			return
		case <-ticker.C:
			c.mu.Lock()
			now := time.Now()
			for taskID, result := range c.taskResults {
				// Remove results older than 10 minutes
				if result.Timestamp.Add(10 * time.Minute).Before(now) {
					delete(c.taskResults, taskID)
				}
			}
			c.mu.Unlock()
		}
	}
}

// DelegateTask delegates a task to an appropriate agent
func (c *CoordinatorAgent) DelegateTask(ctx context.Context, task *Task) error {
	if err := task.Validate(); err != nil {
		return fmt.Errorf("invalid task: %w", err)
	}

	// Route task to appropriate agent
	targetAgentID, err := c.Router.RouteTask(ctx, task, c.Team)
	if err != nil {
		return fmt.Errorf("failed to route task: %w", err)
	}

	// Validate delegation (check for circular dependencies)
	if err := c.Router.ValidateDelegation(task, targetAgentID); err != nil {
		return fmt.Errorf("delegation validation failed: %w", err)
	}

	// Record delegation
	c.Router.RecordDelegation(task, targetAgentID)

	// Mark task as assigned
	task.MarkAssigned(targetAgentID)

	// Store pending task
	c.mu.Lock()
	c.pendingTasks[task.ID] = task
	c.mu.Unlock()

	// Create delegation message
	msg := NewTaskDelegationMessage(c.TeamID, c.AgentID, targetAgentID, task, task.Context)

	// Send message via bus
	data, err := msg.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to serialize delegation message: %w", err)
	}

	// Publish to team-specific channel via outbound message
	outboundMsg := bus.OutboundMessage{
		Content: string(data),
	}
	if err := c.Bus.PublishOutbound(ctx, outboundMsg); err != nil {
		return fmt.Errorf("failed to publish delegation message: %w", err)
	}

	return nil
}

// DelegateAndWait delegates a task and waits for the result
func (c *CoordinatorAgent) DelegateAndWait(ctx context.Context, task *Task) (*TaskResultMessage, error) {
	// Validate task
	if err := task.Validate(); err != nil {
		return nil, fmt.Errorf("invalid task: %w", err)
	}

	// Validate team is not nil (check early before any operations)
	if c.Team == nil {
		return nil, fmt.Errorf("coordinator team is nil")
	}

	// Route task to appropriate agent
	targetAgentID, err := c.Router.RouteTask(ctx, task, c.Team)
	if err != nil {
		return nil, fmt.Errorf("failed to route task: %w", err)
	}

	// Validate delegation
	if err := c.Router.ValidateDelegation(task, targetAgentID); err != nil {
		return nil, fmt.Errorf("delegation validation failed: %w", err)
	}

	// Record delegation
	c.Router.RecordDelegation(task, targetAgentID)

	// Mark task as assigned
	task.MarkAssigned(targetAgentID)

	// Store pending task
	c.mu.Lock()
	c.pendingTasks[task.ID] = task
	c.mu.Unlock()

	// Execute task using executor if available
	if c.Executor != nil {
		// Create context with timeout
		execCtx, cancel := context.WithTimeout(ctx, c.taskTimeout)
		defer cancel()

		// Execute task
		result, err := c.Executor.Execute(execCtx, targetAgentID, task)

		// Clean up
		c.mu.Lock()
		delete(c.pendingTasks, task.ID)
		c.mu.Unlock()

		// Build result message
		if err != nil {
			task.MarkFailed(err)
			return NewTaskResultMessage(
				c.TeamID,
				targetAgentID,
				c.AgentID,
				task.ID,
				TaskStatusFailed,
				nil,
				err,
			), nil
		}

		task.MarkCompleted(result)
		return NewTaskResultMessage(
			c.TeamID,
			targetAgentID,
			c.AgentID,
			task.ID,
			TaskStatusCompleted,
			result,
			nil,
		), nil
	}

	// Fallback: old message bus approach (for backward compatibility)
	// Create delegation message
	msg := NewTaskDelegationMessage(c.TeamID, c.AgentID, targetAgentID, task, task.Context)

	// Send message via bus
	data, err := msg.ToJSON()
	if err != nil {
		c.mu.Lock()
		delete(c.pendingTasks, task.ID)
		c.mu.Unlock()
		return nil, fmt.Errorf("failed to serialize delegation message: %w", err)
	}

	outboundMsg := bus.OutboundMessage{
		Content: string(data),
	}
	if err := c.Bus.PublishOutbound(ctx, outboundMsg); err != nil {
		c.mu.Lock()
		delete(c.pendingTasks, task.ID)
		c.mu.Unlock()
		return nil, fmt.Errorf("failed to publish delegation message: %w", err)
	}

	// Wait for result with timeout
	resultChan := make(chan *TaskResultMessage, 1)
	errorChan := make(chan error, 1)
	
	goroutineCtx, goroutineCancel := context.WithCancel(context.Background())
	defer goroutineCancel() // Ensure goroutine exits
	
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		
		timeout := time.After(c.taskTimeout)
		
		for {
			select {
			case <-goroutineCtx.Done():
				return
			case <-ctx.Done():
				select {
				case errorChan <- ctx.Err():
				case <-goroutineCtx.Done():
				}
				return
			case <-timeout:
				select {
				case errorChan <- errors.New("task execution timeout"):
				case <-goroutineCtx.Done():
				}
				return
			case <-ticker.C:
				c.mu.RLock()
				result, ok := c.taskResults[task.ID]
				c.mu.RUnlock()
				
				if ok {
					select {
					case resultChan <- result:
					case <-goroutineCtx.Done():
					}
					return
				}
			}
		}
	}()
	
	var result *TaskResultMessage
	var waitErr error
	
	select {
	case result = <-resultChan:
		waitErr = nil
	case waitErr = <-errorChan:
		result = nil
	}

	// Clean up
	c.mu.Lock()
	delete(c.pendingTasks, task.ID)
	if waitErr == nil {
		delete(c.taskResults, task.ID)
	}
	c.mu.Unlock()

	return result, waitErr
}

// HandleTaskResult handles a task result message from an agent
func (c *CoordinatorAgent) HandleTaskResult(msg *TaskResultMessage) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Validate team is not nil
	if c.Team == nil {
		return fmt.Errorf("coordinator team is nil")
	}

	// Check if we're waiting for this task
	task, ok := c.pendingTasks[msg.TaskID]
	if !ok {
		return fmt.Errorf("received result for unknown task: %s", msg.TaskID)
	}

	// Validate status transition
	if task.Status == TaskStatusCompleted || task.Status == TaskStatusFailed {
		return fmt.Errorf("task %s already in terminal state: %s", msg.TaskID, task.Status)
	}

	// Check if result already exists (duplicate)
	if _, exists := c.taskResults[msg.TaskID]; exists {
		return fmt.Errorf("task %s result already received", msg.TaskID)
	}

	// Validate result is from authorized agent
	if _, exists := c.Team.Agents[msg.FromAgentID]; !exists {
		return fmt.Errorf("received result from unauthorized agent: %s", msg.FromAgentID)
	}

	// Update task status
	if msg.Status == TaskStatusCompleted {
		task.MarkCompleted(msg.Result)
	} else if msg.Status == TaskStatusFailed {
		task.MarkFailed(errors.New(msg.Error))
	} else {
		return fmt.Errorf("invalid task status in result: %s", msg.Status)
	}

	// Store result
	c.taskResults[msg.TaskID] = msg

	return nil
}

// GetPendingTasks returns all pending tasks
func (c *CoordinatorAgent) GetPendingTasks() []*Task {
	c.mu.RLock()
	defer c.mu.RUnlock()

	tasks := make([]*Task, 0, len(c.pendingTasks))
	for _, task := range c.pendingTasks {
		tasks = append(tasks, task)
	}
	return tasks
}

// SetTaskTimeout sets the timeout for task execution
func (c *CoordinatorAgent) SetTaskTimeout(timeout time.Duration) {
	c.taskTimeout = timeout
}

// ExecuteSequential executes tasks in sequential order, passing results between tasks
func (c *CoordinatorAgent) ExecuteSequential(ctx context.Context, tasks []*Task) ([]any, error) {
	if len(tasks) == 0 {
		return []any{}, nil
	}

	results := make([]any, 0, len(tasks))
	var previousResult any

	for i, task := range tasks {
		// Check if context is cancelled
		if ctx.Err() != nil {
			return results, ctx.Err()
		}

		// Pass previous result as context to next task
		if i > 0 && previousResult != nil {
			if task.Context == nil {
				task.Context = make(map[string]any)
			}
			task.Context["previous_result"] = previousResult
		}

		// Delegate task and wait for completion
		resultMsg, err := c.DelegateAndWait(ctx, task)
		if err != nil {
			return results, fmt.Errorf("task %d failed: %w", i, err)
		}

		// Check if task failed
		if resultMsg.Status == TaskStatusFailed {
			return results, fmt.Errorf("task %d failed: %s", i, resultMsg.Error)
		}

		// Store result
		results = append(results, resultMsg.Result)
		previousResult = resultMsg.Result

		// Update shared context with result
		c.Team.SharedContext.Set(fmt.Sprintf("task_%d_result", i), resultMsg.Result, c.AgentID)
	}

	return results, nil
}

// ExecuteParallel executes tasks in parallel and collects all results
func (c *CoordinatorAgent) ExecuteParallel(ctx context.Context, tasks []*Task) ([]any, error) {
	if len(tasks) == 0 {
		return []any{}, nil
	}

	// Create channels for results and errors
	type taskResult struct {
		index  int
		result any
		err    error
	}
	resultChan := make(chan taskResult, len(tasks))

	// Create cancellable context for child goroutines
	childCtx, cancel := context.WithCancel(ctx)
	defer cancel() // Ensure all child goroutines are cancelled

	// Launch all tasks in parallel
	var wg sync.WaitGroup
	for i, task := range tasks {
		wg.Add(1)
		go func(idx int, t *Task) {
			defer wg.Done()

			// Check if context is cancelled before starting
			select {
			case <-childCtx.Done():
				resultChan <- taskResult{index: idx, err: childCtx.Err()}
				return
			default:
			}

			resultMsg, err := c.DelegateAndWait(childCtx, t)
			if err != nil {
				resultChan <- taskResult{index: idx, err: err}
				return
			}

			if resultMsg.Status == TaskStatusFailed {
				resultChan <- taskResult{index: idx, err: fmt.Errorf("task failed: %s", resultMsg.Error)}
				return
			}

			resultChan <- taskResult{index: idx, result: resultMsg.Result}
		}(i, task)
	}

	// Wait for all tasks to complete with timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(resultChan)
		close(done)
	}()

	// Wait for completion or context cancellation
	select {
	case <-done:
		// All tasks completed normally
	case <-ctx.Done():
		// Context cancelled, cancel all child goroutines
		cancel()
		// Wait for goroutines to exit with timeout
		select {
		case <-done:
			// Goroutines exited cleanly
		case <-time.After(5 * time.Second):
			// Timeout waiting for goroutines - log warning
			logger.WarnCF("team", "Timeout waiting for parallel tasks to exit",
				map[string]any{
					"team_id": c.TeamID,
				})
		}
		return nil, ctx.Err()
	}

	// Collect results safely
	results := make([]any, len(tasks))
	var errors []error

	// Drain resultChan safely - it's already closed
	for res := range resultChan {
		if res.err != nil {
			errors = append(errors, fmt.Errorf("task %d: %w", res.index, res.err))
			// Leave results[res.index] as nil to indicate failure
		} else {
			results[res.index] = res.result
			// Update shared context
			c.Team.SharedContext.Set(fmt.Sprintf("task_%d_result", res.index), res.result, c.AgentID)
		}
	}

	// Return partial results with error if any task failed
	// Note: results array may contain nil values for failed tasks
	if len(errors) > 0 {
		return results, fmt.Errorf("parallel execution had %d failures: %v", len(errors), errors[0])
	}

	return results, nil
}

// DecomposeTask breaks down a main task into subtasks
func (c *CoordinatorAgent) DecomposeTask(mainTask *Task) ([]*Task, error) {
	// Simple decomposition: create subtasks based on task description
	// In a real implementation, this would use LLM to analyze and decompose
	subtasks := []*Task{
		NewTask(fmt.Sprintf("Subtask 1 of %s", mainTask.Description), mainTask.RequiredRole, mainTask.Context),
		NewTask(fmt.Sprintf("Subtask 2 of %s", mainTask.Description), mainTask.RequiredRole, mainTask.Context),
	}

	// Set parent task ID
	for _, subtask := range subtasks {
		subtask.ParentTaskID = mainTask.ID
	}

	return subtasks, nil
}

// ExecuteHierarchical executes a task by decomposing it into subtasks
func (c *CoordinatorAgent) ExecuteHierarchical(ctx context.Context, mainTask *Task) (any, error) {
	// Decompose main task
	subtasks, err := c.DecomposeTask(mainTask)
	if err != nil {
		return nil, fmt.Errorf("failed to decompose task: %w", err)
	}

	// Execute subtasks sequentially (can be parallel in advanced implementation)
	subtaskResults := make([]any, 0, len(subtasks))
	for i, subtask := range subtasks {
		resultMsg, err := c.DelegateAndWait(ctx, subtask)
		if err != nil {
			return nil, fmt.Errorf("subtask %d failed: %w", i, err)
		}

		if resultMsg.Status == TaskStatusFailed {
			return nil, fmt.Errorf("subtask %d failed: %s", i, resultMsg.Error)
		}

		subtaskResults = append(subtaskResults, resultMsg.Result)

		// Analyze intermediate results and potentially create more subtasks
		if c.shouldCreateMoreSubtasks(subtaskResults) {
			additionalSubtask := NewTask(
				fmt.Sprintf("Additional subtask based on result %d", i),
				mainTask.RequiredRole,
				map[string]any{"previous_results": subtaskResults},
			)
			additionalSubtask.ParentTaskID = mainTask.ID
			subtasks = append(subtasks, additionalSubtask)
		}
	}

	// Integrate all subtask results
	integratedResult := c.integrateResults(subtaskResults)

	// Update shared context
	c.Team.SharedContext.Set(fmt.Sprintf("task_%s_result", mainTask.ID), integratedResult, c.AgentID)

	return integratedResult, nil
}

// shouldCreateMoreSubtasks analyzes intermediate results to decide if more subtasks are needed
func (c *CoordinatorAgent) shouldCreateMoreSubtasks(results []any) bool {
	// Simple heuristic: don't create more subtasks for now
	// In real implementation, this would use LLM to analyze results
	return false
}

// integrateResults combines subtask results into a final result
func (c *CoordinatorAgent) integrateResults(results []any) any {
	// Simple integration: concatenate results
	// In real implementation, this would use LLM to synthesize results
	integrated := make(map[string]any)
	for i, result := range results {
		integrated[fmt.Sprintf("subtask_%d", i)] = result
	}
	return integrated
}

// ReassignTask reassigns a task from one agent to another
func (c *CoordinatorAgent) ReassignTask(ctx context.Context, task *Task, newAgentID string) error {
	// Validate task state
	if task.Status == TaskStatusCompleted {
		return fmt.Errorf("cannot reassign completed task")
	}
	if task.Status == TaskStatusCancelled {
		return fmt.Errorf("cannot reassign cancelled task")
	}

	// Validate new agent exists
	_, exists := c.Team.Agents[newAgentID]
	if !exists {
		return fmt.Errorf("agent %s not found in team", newAgentID)
	}

	// Cancel current assignment (mark as cancelled)
	task.MarkCancelled()

	// Create new task with same description (but fresh delegation chain)
	newTask := NewTask(task.Description, task.RequiredRole, task.Context)
	newTask.ParentTaskID = task.ParentTaskID
	// Don't copy delegation chain - start fresh for reassignment

	// Validate delegation with new task (fresh chain)
	if err := c.Router.ValidateDelegation(newTask, newAgentID); err != nil {
		return fmt.Errorf("reassignment validation failed: %w", err)
	}

	// Delegate to new agent
	if err := c.DelegateTask(ctx, newTask); err != nil {
		return fmt.Errorf("failed to reassign task: %w", err)
	}

	return nil
}

// AggregationStrategy defines how results are combined
type AggregationStrategy string

const (
	StrategyConcat    AggregationStrategy = "concatenate"
	StrategyMerge     AggregationStrategy = "merge"
	StrategyIntegrate AggregationStrategy = "integrate"
)

// AggregateResults combines results from multiple tasks with attribution
func (c *CoordinatorAgent) AggregateResults(results []any, agentIDs []string, strategy AggregationStrategy) (any, error) {
	if len(results) == 0 || len(agentIDs) == 0 {
		return nil, fmt.Errorf("results and agentIDs cannot be empty")
	}
	if len(results) != len(agentIDs) {
		return nil, fmt.Errorf("results and agentIDs length mismatch: %d vs %d", len(results), len(agentIDs))
	}

	switch strategy {
	case StrategyConcat:
		return c.aggregateConcatenate(results, agentIDs)
	case StrategyMerge:
		return c.aggregateMerge(results, agentIDs)
	case StrategyIntegrate:
		return c.aggregateIntegrate(results, agentIDs)
	default:
		return nil, fmt.Errorf("unknown aggregation strategy: %s", strategy)
	}
}

// aggregateConcatenate concatenates results in order (for sequential)
func (c *CoordinatorAgent) aggregateConcatenate(results []any, agentIDs []string) (any, error) {
	aggregated := make([]map[string]any, 0, len(results))

	for i, result := range results {
		entry := map[string]any{
			"agent_id": agentIDs[i],
			"result":   result,
			"index":    i,
		}
		aggregated = append(aggregated, entry)
	}

	return map[string]any{
		"strategy": "concatenate",
		"results":  aggregated,
		"count":    len(results),
	}, nil
}

// aggregateMerge merges results into a single structure (for parallel)
func (c *CoordinatorAgent) aggregateMerge(results []any, agentIDs []string) (any, error) {
	merged := make(map[string]any)

	for i, result := range results {
		agentKey := fmt.Sprintf("agent_%s", agentIDs[i])
		merged[agentKey] = result
	}

	merged["strategy"] = "merge"
	merged["agent_count"] = len(agentIDs)

	return merged, nil
}

// aggregateIntegrate integrates results hierarchically (for hierarchical)
func (c *CoordinatorAgent) aggregateIntegrate(results []any, agentIDs []string) (any, error) {
	integrated := map[string]any{
		"strategy":     "integrate",
		"main_result":  nil,
		"sub_results":  []map[string]any{},
		"contributors": agentIDs,
	}

	// First result is main result
	if len(results) > 0 {
		integrated["main_result"] = results[0]
	}

	// Rest are sub-results
	for i := 1; i < len(results); i++ {
		subResult := map[string]any{
			"agent_id": agentIDs[i],
			"result":   results[i],
		}
		integrated["sub_results"] = append(integrated["sub_results"].([]map[string]any), subResult)
	}

	return integrated, nil
}

// ResolveConflicts detects and resolves contradictory results
func (c *CoordinatorAgent) ResolveConflicts(results []any, agentIDs []string, resolution string) (any, error) {
	// Simple conflict detection: check if results are different
	conflicts := c.detectConflicts(results)

	if len(conflicts) == 0 {
		// No conflicts, return first result
		return map[string]any{
			"resolved":      results[0],
			"had_conflicts": false,
			"resolution":    "none_needed",
		}, nil
	}

	// Resolve based on strategy
	var resolved any
	switch resolution {
	case "voting":
		resolved = c.resolveByVoting(results, agentIDs)
	case "priority":
		resolved = c.resolveByPriority(results, agentIDs)
	case "consensus":
		resolved = c.resolveByConsensus(results, agentIDs)
	default:
		return nil, fmt.Errorf("unknown resolution strategy: %s", resolution)
	}

	return map[string]any{
		"resolved":      resolved,
		"had_conflicts": true,
		"conflicts":     conflicts,
		"resolution":    resolution,
		"contributors":  agentIDs,
	}, nil
}

// detectConflicts identifies contradictory results
func (c *CoordinatorAgent) detectConflicts(results []any) []string {
	conflicts := []string{}

	// Simple check: if results are different, mark as conflict
	if len(results) < 2 {
		return conflicts
	}

	firstResult := fmt.Sprintf("%v", results[0])
	for i := 1; i < len(results); i++ {
		if fmt.Sprintf("%v", results[i]) != firstResult {
			conflicts = append(conflicts, fmt.Sprintf("result_%d_differs", i))
		}
	}

	return conflicts
}

// resolveByVoting uses majority voting
func (c *CoordinatorAgent) resolveByVoting(results []any, agentIDs []string) any {
	votes := make(map[string]int)
	resultMap := make(map[string]any)

	for _, result := range results {
		key := fmt.Sprintf("%v", result)
		votes[key]++
		resultMap[key] = result
	}

	// Find result with most votes
	maxVotes := 0
	var winner string
	for key, count := range votes {
		if count > maxVotes {
			maxVotes = count
			winner = key
		}
	}

	return resultMap[winner]
}

// resolveByPriority uses first agent's result (priority order)
func (c *CoordinatorAgent) resolveByPriority(results []any, agentIDs []string) any {
	if len(results) > 0 {
		return results[0]
	}
	return nil
}

// resolveByConsensus requires all agents to agree
func (c *CoordinatorAgent) resolveByConsensus(results []any, agentIDs []string) any {
	if len(results) == 0 {
		return nil
	}

	// Check if all results are the same
	firstResult := fmt.Sprintf("%v", results[0])
	for _, result := range results {
		if fmt.Sprintf("%v", result) != firstResult {
			return map[string]any{
				"consensus": false,
				"message":   "No consensus reached",
			}
		}
	}

	return results[0]
}

// ShouldRetry determines if a failed task should be retried
func (c *CoordinatorAgent) ShouldRetry(task *Task, attemptCount int, maxAttempts int) bool {
	if attemptCount >= maxAttempts {
		return false
	}

	// Don't retry if task was cancelled
	if task.Status == TaskStatusCancelled {
		return false
	}

	// Retry if task failed
	if task.Status == TaskStatusFailed {
		return true
	}

	return false
}

// RetryTask retries a failed task with exponential backoff
func (c *CoordinatorAgent) RetryTask(ctx context.Context, task *Task, maxAttempts int) (*TaskResultMessage, error) {
	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		// Log retry attempt
		c.Team.SharedContext.AddHistoryEntry("system", "task_retry",
			map[string]string{
				"task_id": task.ID,
				"attempt": fmt.Sprintf("%d", attempt),
			})

		// Reset task status
		task.Status = TaskStatusPending
		task.Error = nil

		// Delegate and wait
		result, err := c.DelegateAndWait(ctx, task)
		if err == nil && result.Status == TaskStatusCompleted {
			return result, nil
		}

		lastErr = err

		// Exponential backoff
		if attempt < maxAttempts {
			backoff := time.Duration(attempt*attempt) * 100 * time.Millisecond
			time.Sleep(backoff)
		}
	}

	return nil, fmt.Errorf("task failed after %d attempts: %w", maxAttempts, lastErr)
}

// ReassignFailedTask reassigns a failed task to a different agent
func (c *CoordinatorAgent) ReassignFailedTask(ctx context.Context, task *Task) error {
	// Find another agent with the same role
	var targetAgentID string
	for agentID, agent := range c.Team.Agents {
		if agentID != task.AssignedAgentID && agent.Role == task.RequiredRole && agent.Status == StatusIdle {
			targetAgentID = agentID
			break
		}
	}

	if targetAgentID == "" {
		return fmt.Errorf("no available agent with role '%s' for reassignment", task.RequiredRole)
	}

	// Use ReassignTask method
	return c.ReassignTask(ctx, task, targetAgentID)
}

// AbortWorkflow aborts the entire workflow on critical failure
func (c *CoordinatorAgent) AbortWorkflow(reason string) error {
	// Update team status
	c.Team.Status = TeamStatusPaused

	// Log abort
	c.Team.SharedContext.AddHistoryEntry("system", "workflow_aborted",
		map[string]string{
			"reason":      reason,
			"coordinator": c.AgentID,
		})

	// Cancel all pending tasks
	c.mu.Lock()
	for taskID, task := range c.pendingTasks {
		task.MarkCancelled()
		delete(c.pendingTasks, taskID)
	}
	c.mu.Unlock()

	return fmt.Errorf("workflow aborted: %s", reason)
}

// LogFailure logs a failure to team memory
func (c *CoordinatorAgent) LogFailure(task *Task, agentID string, err error) {
	failureRecord := map[string]any{
		"task_id":   task.ID,
		"agent_id":  agentID,
		"error":     err.Error(),
		"timestamp": time.Now().Format(time.RFC3339),
		"status":    string(task.Status),
	}

	// Store in shared context
	failureKey := fmt.Sprintf("failure_%s_%d", task.ID, time.Now().Unix())
	c.Team.SharedContext.Set(failureKey, failureRecord, c.AgentID)

	// Add to history
	c.Team.SharedContext.AddHistoryEntry("system", "task_failure",
		map[string]string{
			"task_id":  task.ID,
			"agent_id": agentID,
			"error":    err.Error(),
		})
}
