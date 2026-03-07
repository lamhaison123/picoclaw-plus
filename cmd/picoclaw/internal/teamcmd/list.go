// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package teamcmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/sipeed/picoclaw/pkg/team"
)

func newListCommand(managerFn func() (*team.TeamManager, error)) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all active teams",
		Long:  "Display a list of all currently active teams with their status",
		RunE: func(cmd *cobra.Command, _ []string) error {
			tm, err := managerFn()
			if err != nil {
				return err
			}

			// Get all teams
			teams := tm.GetAllTeams()

			if len(teams) == 0 {
				fmt.Println("No active teams")
				fmt.Println()
				fmt.Println("Use 'picoclaw team create' to create a new team")
				return nil
			}

			fmt.Printf("Active Teams (%d):\n", len(teams))
			fmt.Println("------------------")

			for _, t := range teams {
				fmt.Printf("  %s (%s)\n", t.Name, t.ID)
				fmt.Printf("    Pattern: %s\n", t.Pattern)
				fmt.Printf("    Status: %s\n", t.Status)
				fmt.Printf("    Agents: %d\n", len(t.Agents))
				fmt.Println()
			}

			return nil
		},
	}

	return cmd
}
