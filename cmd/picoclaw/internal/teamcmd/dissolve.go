// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package teamcmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/sipeed/picoclaw/pkg/team"
)

func newDissolveCommand(managerFn func() (*team.TeamManager, error)) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "dissolve <team-id>",
		Short: "Dissolve a team",
		Long:  "Dissolve a team and persist its memory. All agents will be deregistered.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			tm, err := managerFn()
			if err != nil {
				return err
			}

			teamID := args[0]

			// Get team info before dissolution
			status, err := tm.GetTeamStatus(teamID)
			if err != nil {
				return fmt.Errorf("failed to get team status: %w", err)
			}

			// Confirm dissolution if not forced
			if !force {
				fmt.Printf("Are you sure you want to dissolve team '%s' (%s)? [y/N]: ", status.TeamName, teamID)
				var response string
				fmt.Scanln(&response)
				if response != "y" && response != "Y" {
					fmt.Println("Dissolution cancelled")
					return nil
				}
			}

			// Dissolve team
			ctx := context.Background()
			if err := tm.DissolveTeam(ctx, teamID); err != nil {
				return fmt.Errorf("failed to dissolve team: %w", err)
			}

			fmt.Printf("✓ Team '%s' dissolved successfully\n", status.TeamName)
			fmt.Println("  Team memory has been persisted")
			fmt.Printf("  %d agents deregistered\n", len(status.AgentStatuses))

			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force dissolution without confirmation")

	return cmd
}
