// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package teamcmd

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/sipeed/picoclaw/pkg/team"
)

func newCreateCommand(
	managerFn func() (*team.TeamManager, error),
	workspaceFn func() (string, error),
) *cobra.Command {
	var configFile string

	cmd := &cobra.Command{
		Use:   "create [config-file]",
		Short: "Create a new team from configuration",
		Long:  "Create a new multi-agent team from a JSON configuration file",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			tm, err := managerFn()
			if err != nil {
				return err
			}

			workspace, err := workspaceFn()
			if err != nil {
				return err
			}

			// Determine config file path
			if len(args) > 0 {
				configFile = args[0]
			}

			if configFile == "" {
				return fmt.Errorf("config file is required")
			}

			// Make path absolute if relative
			if !filepath.IsAbs(configFile) {
				configFile = filepath.Join(workspace, configFile)
			}

			// Load configuration
			config, err := team.LoadTeamConfig(configFile)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			// Create team
			ctx := context.Background()
			createdTeam, err := tm.CreateTeam(ctx, config)
			if err != nil {
				return fmt.Errorf("failed to create team: %w", err)
			}

			// Print success message
			fmt.Printf("✓ Team created successfully!\n")
			fmt.Printf("  ID: %s\n", createdTeam.ID)
			fmt.Printf("  Name: %s\n", createdTeam.Name)
			fmt.Printf("  Pattern: %s\n", createdTeam.Pattern)
			fmt.Printf("  Agents: %d\n", len(createdTeam.Agents))
			fmt.Printf("  Status: %s\n", createdTeam.Status)

			return nil
		},
	}

	cmd.Flags().StringVarP(&configFile, "config", "c", "", "Path to team configuration file")

	return cmd
}
