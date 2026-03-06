// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package teamcmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/sipeed/picoclaw/pkg/team"
)

func newStatusCommand(managerFn func() (*team.TeamManager, error)) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status <team-id>",
		Short: "Show detailed team status",
		Long:  "Display detailed information about a team including agent statuses and metrics",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			tm, err := managerFn()
			if err != nil {
				return err
			}

			teamID := args[0]

			// Get team status
			status, err := tm.GetTeamStatus(teamID)
			if err != nil {
				return fmt.Errorf("failed to get team status: %w", err)
			}

			// Print team information
			fmt.Printf("Team: %s\n", status.TeamName)
			fmt.Printf("ID: %s\n", status.TeamID)
			fmt.Printf("Status: %s\n", status.Status)
			fmt.Printf("Pattern: %s\n", status.Pattern)
			fmt.Printf("Agent Count: %d\n", status.AgentCount)
			fmt.Printf("Uptime: %s\n", formatDuration(status.Uptime))
			fmt.Println()

			// Print agent statuses
			fmt.Println("Agents:")
			fmt.Println("-------")
			for agentID, agentStatus := range status.AgentStatuses {
				fmt.Printf("  %s: %s\n", agentID, agentStatus)
			}
			fmt.Println()

			// Get and print metrics
			metrics := tm.GetMetrics().GetMetrics(teamID)
			fmt.Println("Metrics:")
			fmt.Println("--------")
			if metrics != nil {
				fmt.Printf("  Tasks Delegated: %d\n", metrics.TaskDelegationCount)
				fmt.Printf("  Tasks Completed: %d\n", metrics.TaskCompletionCount)
				fmt.Printf("  Tasks Failed: %d\n", metrics.TaskFailureCount)
				fmt.Printf("  Consensus Count: %d\n", metrics.ConsensusCount)
				fmt.Printf("  Average Task Duration: %s\n", formatDuration(metrics.AverageTaskDuration))
			} else {
				fmt.Println("  No metrics available yet")
			}

			return nil
		},
	}

	return cmd
}

func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%.1fm", d.Minutes())
	}
	return fmt.Sprintf("%.1fh", d.Hours())
}
