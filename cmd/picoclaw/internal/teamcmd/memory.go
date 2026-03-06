// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package teamcmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/sipeed/picoclaw/pkg/team/memory"
)

func newMemoryCommand(memoryFn func() (*memory.TeamMemory, error)) *cobra.Command {
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "memory <team-id>",
		Short: "Display team memory",
		Long:  "Display the persisted memory record of a dissolved team",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			tm, err := memoryFn()
			if err != nil {
				return err
			}

			teamID := args[0]

			// Load team record
			record, err := tm.LoadTeamRecord(teamID)
			if err != nil {
				return fmt.Errorf("failed to load team record: %w", err)
			}

			// Output as JSON if requested
			if jsonOutput {
				data, err := json.MarshalIndent(record, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal record: %w", err)
				}
				fmt.Println(string(data))
				return nil
			}

			// Print formatted output
			fmt.Printf("Team Memory: %s\n", record.TeamName)
			fmt.Printf("ID: %s\n", record.TeamID)
			fmt.Printf("Pattern: %s\n", record.Pattern)
			fmt.Printf("Start Time: %s\n", record.StartTime.Format("2006-01-02 15:04:05"))
			fmt.Printf("End Time: %s\n", record.EndTime.Format("2006-01-02 15:04:05"))
			fmt.Printf("Duration: %s\n", formatDuration(record.EndTime.Sub(record.StartTime)))
			fmt.Printf("Outcome: %s\n", record.Outcome)
			fmt.Println()

			// Print shared context
			fmt.Println("Shared Context:")
			fmt.Println("---------------")
			if len(record.SharedContext) == 0 {
				fmt.Println("  (empty)")
			} else {
				for key, value := range record.SharedContext {
					fmt.Printf("  %s: %v\n", key, value)
				}
			}
			fmt.Println()

			// Print tasks
			fmt.Printf("Tasks: %d\n", len(record.Tasks))
			fmt.Println("------")
			for i, task := range record.Tasks {
				fmt.Printf("  %d. %s\n", i+1, task.Description)
				fmt.Printf("     Role: %s\n", task.Role)
				fmt.Printf("     Status: %s\n", task.Status)
				if task.Error != "" {
					fmt.Printf("     Error: %s\n", task.Error)
				}
			}
			fmt.Println()

			// Print consensus records
			if len(record.Consensus) > 0 {
				fmt.Printf("Consensus: %d\n", len(record.Consensus))
				fmt.Println("----------")
				for i, cons := range record.Consensus {
					fmt.Printf("  %d. %s\n", i+1, cons.Question)
					fmt.Printf("     Outcome: %s\n", cons.Outcome)
					fmt.Printf("     Voting Rule: %s\n", cons.VotingRule)
					fmt.Printf("     Total Votes: %d\n", cons.TotalVotes)
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")

	return cmd
}
