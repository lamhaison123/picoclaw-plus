// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package memory

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// TeamMemory manages persistence of team records
type TeamMemory struct {
	workspace string
	memoryDir string
}

// TeamMemoryRecord represents a complete record of a team's execution
type TeamMemoryRecord struct {
	TeamID        string            `json:"team_id"`
	TeamName      string            `json:"team_name"`
	Pattern       string            `json:"pattern"`
	StartTime     time.Time         `json:"start_time"`
	EndTime       time.Time         `json:"end_time"`
	SharedContext map[string]any    `json:"shared_context"`
	Tasks         []TaskRecord      `json:"tasks"`
	Consensus     []ConsensusRecord `json:"consensus"`
	Outcome       string            `json:"outcome"`
}

// TaskRecord represents a task execution record
type TaskRecord struct {
	TaskID      string         `json:"task_id"`
	Description string         `json:"description"`
	Role        string         `json:"role"`
	AgentID     string         `json:"agent_id"`
	Status      string         `json:"status"`
	Result      any            `json:"result,omitempty"`
	Error       string         `json:"error,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	StartedAt   *time.Time     `json:"started_at,omitempty"`
	CompletedAt *time.Time     `json:"completed_at,omitempty"`
	Context     map[string]any `json:"context,omitempty"`
}

// ConsensusRecord represents a consensus voting record
type ConsensusRecord struct {
	ConsensusID string         `json:"consensus_id"`
	Question    string         `json:"question"`
	Options     []string       `json:"options"`
	Outcome     string         `json:"outcome"`
	VotingRule  string         `json:"voting_rule"`
	TotalVotes  int            `json:"total_votes"`
	Timestamp   time.Time      `json:"timestamp"`
	Context     map[string]any `json:"context,omitempty"`
}

// NewTeamMemory creates a new team memory manager
func NewTeamMemory(workspace string) *TeamMemory {
	memoryDir := filepath.Join(workspace, "memory", "teams")
	os.MkdirAll(memoryDir, 0o755)

	tm := &TeamMemory{
		workspace: workspace,
		memoryDir: memoryDir,
	}

	// Clean up stale temp files on startup
	tm.cleanupStaleTempFiles()

	return tm
}

// cleanupStaleTempFiles removes temporary files older than 1 hour
func (tm *TeamMemory) cleanupStaleTempFiles() {
	files, err := os.ReadDir(tm.memoryDir)
	if err != nil {
		return // Silently fail if directory doesn't exist yet
	}

	now := time.Now()
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".tmp" {
			filePath := filepath.Join(tm.memoryDir, file.Name())
			info, err := file.Info()
			if err != nil {
				continue
			}

			// Remove temp files older than 1 hour
			if now.Sub(info.ModTime()) > time.Hour {
				os.Remove(filePath)
			}
		}
	}
}

// SaveTeamRecord saves a team record to disk
func (tm *TeamMemory) SaveTeamRecord(record *TeamMemoryRecord) error {
	// Generate filename
	timestamp := record.EndTime.Format("20060102_150405")
	filename := fmt.Sprintf("%s_%s.json", record.TeamID, timestamp)
	filepath := filepath.Join(tm.memoryDir, filename)

	// Serialize to JSON
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal team record: %w", err)
	}

	// Write atomically
	tempFile := filepath + ".tmp"
	if err := os.WriteFile(tempFile, data, 0o644); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	if err := os.Rename(tempFile, filepath); err != nil {
		os.Remove(tempFile)
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}

// LoadTeamRecord loads a team record by team ID (latest record)
func (tm *TeamMemory) LoadTeamRecord(teamID string) (*TeamMemoryRecord, error) {
	records, err := tm.ListTeamRecords()
	if err != nil {
		return nil, err
	}

	// Find latest record for this team
	for i := len(records) - 1; i >= 0; i-- {
		if records[i] == teamID {
			// Load the file
			files, err := os.ReadDir(tm.memoryDir)
			if err != nil {
				return nil, fmt.Errorf("failed to read memory directory: %w", err)
			}

			for _, file := range files {
				if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
					// Check if filename starts with teamID
					if len(file.Name()) > len(teamID) && file.Name()[:len(teamID)] == teamID {
						data, err := os.ReadFile(filepath.Join(tm.memoryDir, file.Name()))
						if err != nil {
							continue
						}

						var record TeamMemoryRecord
						if err := json.Unmarshal(data, &record); err != nil {
							continue
						}

						if record.TeamID == teamID {
							return &record, nil
						}
					}
				}
			}
		}
	}

	return nil, fmt.Errorf("team record not found: %s", teamID)
}

// ListTeamRecords lists all team IDs with saved records
func (tm *TeamMemory) ListTeamRecords() ([]string, error) {
	files, err := os.ReadDir(tm.memoryDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read memory directory: %w", err)
	}

	teamIDs := make(map[string]bool)
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			// Extract team ID from filename (format: teamID_timestamp.json)
			name := file.Name()
			underscoreIdx := -1
			for i, c := range name {
				if c == '_' {
					underscoreIdx = i
					break
				}
			}
			if underscoreIdx > 0 {
				teamID := name[:underscoreIdx]
				teamIDs[teamID] = true
			}
		}
	}

	// Convert to sorted slice
	result := make([]string, 0, len(teamIDs))
	for teamID := range teamIDs {
		result = append(result, teamID)
	}
	sort.Strings(result)

	return result, nil
}

// GetWorkspace returns the workspace path
func (tm *TeamMemory) GetWorkspace() string {
	return tm.workspace
}
