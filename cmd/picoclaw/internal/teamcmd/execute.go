// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package teamcmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/sipeed/picoclaw/pkg/team"
)

func newExecuteCommand(managerFn func() (*team.TeamManager, error)) *cobra.Command {
	var (
		taskDescription string
		requiredRole    string
		timeout         int
	)

	cmd := &cobra.Command{
		Use:   "execute <team-id>",
		Short: "Execute a task using a team",
		Long:  "Delegate a task to a team for execution using the configured collaboration pattern",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			teamID := args[0]

			if taskDescription == "" {
				return fmt.Errorf("task description is required (use --task flag)")
			}

			tm, err := managerFn()
			if err != nil {
				return err
			}

			// Verify team exists
			_, err = tm.GetTeam(teamID)
			if err != nil {
				return fmt.Errorf("failed to get team: %w", err)
			}

			fmt.Printf("Executing task on team '%s'...\n", teamID)
			fmt.Printf("Task: %s\n", taskDescription)
			if requiredRole != "" {
				fmt.Printf("Role: %s\n", requiredRole)
			}
			fmt.Println()

			// Execute task with timeout
			ctx := context.Background()
			if timeout > 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
				defer cancel()
			}

			var result interface{}
			if requiredRole != "" {
				// Execute with specific role
				result, err = tm.ExecuteTaskWithRole(ctx, teamID, taskDescription, requiredRole)
			} else {
				// Auto-select role
				result, err = tm.ExecuteTask(ctx, teamID, taskDescription)
			}

			if err != nil {
				return fmt.Errorf("task execution failed: %w", err)
			}

			// Print result
			fmt.Println("✓ Task completed successfully!")
			fmt.Println()
			fmt.Println("Result:")
			fmt.Println("-------")
			fmt.Printf("%v\n", result)

			return nil
		},
	}

	cmd.Flags().StringVarP(&taskDescription, "task", "t", "", "Task description (required)")
	cmd.Flags().StringVarP(&requiredRole, "role", "r", "", "Required role for task execution (optional, auto-select if not specified)")
	cmd.Flags().IntVar(&timeout, "timeout", 300, "Task timeout in seconds (0 for no timeout)")

	cmd.MarkFlagRequired("task")

	return cmd
}
