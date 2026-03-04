// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package team

import (
	"encoding/json"
	"time"
)

// MessageType defines the type of team communication message
type MessageType string

const (
	MessageTypeTaskDelegation   MessageType = "task_delegation"
	MessageTypeTaskResult       MessageType = "task_result"
	MessageTypeConsensusRequest MessageType = "consensus_request"
	MessageTypeConsensusVote    MessageType = "consensus_vote"
)

// TaskDelegationMessage represents a task delegation from coordinator to agent
type TaskDelegationMessage struct {
	MessageID   string                 `json:"message_id"`
	TeamID      string                 `json:"team_id"`
	FromAgentID string                 `json:"from_agent_id"`
	ToAgentID   string                 `json:"to_agent_id"`
	Task        *Task                  `json:"task"`
	Context     map[string]any `json:"context,omitempty"`
	Timestamp   time.Time              `json:"timestamp"`
}

// NewTaskDelegationMessage creates a new task delegation message
func NewTaskDelegationMessage(teamID, fromAgentID, toAgentID string, task *Task, context map[string]any) *TaskDelegationMessage {
	return &TaskDelegationMessage{
		MessageID:   generateMessageID(),
		TeamID:      teamID,
		FromAgentID: fromAgentID,
		ToAgentID:   toAgentID,
		Task:        task,
		Context:     context,
		Timestamp:   time.Now(),
	}
}

// ToJSON serializes the message to JSON
func (m *TaskDelegationMessage) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

// TaskDelegationMessageFromJSON deserializes a task delegation message from JSON
func TaskDelegationMessageFromJSON(data []byte) (*TaskDelegationMessage, error) {
	var msg TaskDelegationMessage
	err := json.Unmarshal(data, &msg)
	return &msg, err
}

// TaskResultMessage represents the result of a completed task
type TaskResultMessage struct {
	MessageID   string      `json:"message_id"`
	TeamID      string      `json:"team_id"`
	FromAgentID string      `json:"from_agent_id"`
	ToAgentID   string      `json:"to_agent_id"`
	TaskID      string      `json:"task_id"`
	Status      TaskStatus  `json:"status"`
	Result      any `json:"result,omitempty"`
	Error       string      `json:"error,omitempty"`
	Timestamp   time.Time   `json:"timestamp"`
}

// NewTaskResultMessage creates a new task result message
func NewTaskResultMessage(teamID, fromAgentID, toAgentID, taskID string, status TaskStatus, result any, err error) *TaskResultMessage {
	msg := &TaskResultMessage{
		MessageID:   generateMessageID(),
		TeamID:      teamID,
		FromAgentID: fromAgentID,
		ToAgentID:   toAgentID,
		TaskID:      taskID,
		Status:      status,
		Result:      result,
		Timestamp:   time.Now(),
	}
	if err != nil {
		msg.Error = err.Error()
	}
	return msg
}

// ToJSON serializes the message to JSON
func (m *TaskResultMessage) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

// TaskResultMessageFromJSON deserializes a task result message from JSON
func TaskResultMessageFromJSON(data []byte) (*TaskResultMessage, error) {
	var msg TaskResultMessage
	err := json.Unmarshal(data, &msg)
	return &msg, err
}

// ConsensusRequestMessage represents a request for consensus voting
type ConsensusRequestMessage struct {
	MessageID   string                 `json:"message_id"`
	TeamID      string                 `json:"team_id"`
	FromAgentID string                 `json:"from_agent_id"`
	ConsensusID string                 `json:"consensus_id"`
	Question    string                 `json:"question"`
	Options     []string               `json:"options"`
	VotingRule  VotingRule             `json:"voting_rule"`
	Context     map[string]any `json:"context,omitempty"`
	Timeout     int                    `json:"timeout_seconds"`
	Timestamp   time.Time              `json:"timestamp"`
}

// NewConsensusRequestMessage creates a new consensus request message
func NewConsensusRequestMessage(teamID, fromAgentID, consensusID, question string, options []string, rule VotingRule, timeout int, context map[string]any) *ConsensusRequestMessage {
	return &ConsensusRequestMessage{
		MessageID:   generateMessageID(),
		TeamID:      teamID,
		FromAgentID: fromAgentID,
		ConsensusID: consensusID,
		Question:    question,
		Options:     options,
		VotingRule:  rule,
		Context:     context,
		Timeout:     timeout,
		Timestamp:   time.Now(),
	}
}

// ToJSON serializes the message to JSON
func (m *ConsensusRequestMessage) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

// ConsensusRequestMessageFromJSON deserializes a consensus request message from JSON
func ConsensusRequestMessageFromJSON(data []byte) (*ConsensusRequestMessage, error) {
	var msg ConsensusRequestMessage
	err := json.Unmarshal(data, &msg)
	return &msg, err
}

// ConsensusVoteMessage represents a vote in response to a consensus request
type ConsensusVoteMessage struct {
	MessageID   string                 `json:"message_id"`
	TeamID      string                 `json:"team_id"`
	FromAgentID string                 `json:"from_agent_id"`
	ToAgentID   string                 `json:"to_agent_id"`
	ConsensusID string                 `json:"consensus_id"`
	Vote        string                 `json:"vote"`
	Weight      float64                `json:"weight,omitempty"`
	Rationale   string                 `json:"rationale,omitempty"`
	Context     map[string]any `json:"context,omitempty"`
	Timestamp   time.Time              `json:"timestamp"`
}

// NewConsensusVoteMessage creates a new consensus vote message
func NewConsensusVoteMessage(teamID, fromAgentID, toAgentID, consensusID, vote string, weight float64, rationale string, context map[string]any) *ConsensusVoteMessage {
	return &ConsensusVoteMessage{
		MessageID:   generateMessageID(),
		TeamID:      teamID,
		FromAgentID: fromAgentID,
		ToAgentID:   toAgentID,
		ConsensusID: consensusID,
		Vote:        vote,
		Weight:      weight,
		Rationale:   rationale,
		Context:     context,
		Timestamp:   time.Now(),
	}
}

// ToJSON serializes the message to JSON
func (m *ConsensusVoteMessage) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

// ConsensusVoteMessageFromJSON deserializes a consensus vote message from JSON
func ConsensusVoteMessageFromJSON(data []byte) (*ConsensusVoteMessage, error) {
	var msg ConsensusVoteMessage
	err := json.Unmarshal(data, &msg)
	return &msg, err
}

// generateMessageID generates a unique message ID
func generateMessageID() string {
	return "msg_" + time.Now().Format("20060102150405") + "_" + randomString(8)
}
