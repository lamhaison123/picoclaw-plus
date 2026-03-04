// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package teamcmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/sipeed/picoclaw/cmd/picoclaw/internal"
	"github.com/sipeed/picoclaw/pkg/agent"
	"github.com/sipeed/picoclaw/pkg/bus"
	"github.com/sipeed/picoclaw/pkg/config"
	"github.com/sipeed/picoclaw/pkg/providers"
	"github.com/sipeed/picoclaw/pkg/team"
	"github.com/sipeed/picoclaw/pkg/team/memory"
)

type deps struct {
	workspace   string
	cfg         *config.Config
	teamManager *team.TeamManager
	teamMemory  *memory.TeamMemory
}

func NewTeamCommand() *cobra.Command {
	var d deps

	cmd := &cobra.Command{
		Use:   "team",
		Short: "Manage multi-agent teams",
		Long:  "Create, manage, and monitor multi-agent collaboration teams",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := internal.LoadConfig()
			if err != nil {
				return fmt.Errorf("error loading config: %w", err)
			}

			d.cfg = cfg
			d.workspace = cfg.WorkspacePath()

			// Initialize team memory
			d.teamMemory = memory.NewTeamMemory(d.workspace)

			// Initialize team manager
		// Initialize team manager
		prov, _, err := providers.CreateProvider(cfg)
		if err != nil {
			return fmt.Errorf("error creating provider: %w", err)
		}
		msgBus := bus.NewMessageBus()

		// Create agent loop for task execution
		agentLoop := agent.NewAgentLoop(cfg, msgBus, prov)

		// Get registry from agent loop (they must share the same registry!)
		registry := agentLoop.GetRegistry()

		// Create executor that uses agent loop
		executor := team.NewDirectAgentExecutor(agentLoop)

		d.teamManager = team.NewTeamManager(registry, msgBus)
		d.teamManager.SetProvider(prov, cfg) // Set provider FIRST
		d.teamManager.SetAgentExecutor(executor)
		d.teamManager.SetTeamMemory(d.teamMemory) // Load teams AFTER provider is set

		return nil

		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	managerFn := func() (*team.TeamManager, error) {
		if d.teamManager == nil {
			return nil, fmt.Errorf("team manager is not initialized")
		}
		return d.teamManager, nil
	}

	memoryFn := func() (*memory.TeamMemory, error) {
		if d.teamMemory == nil {
			return nil, fmt.Errorf("team memory is not initialized")
		}
		return d.teamMemory, nil
	}

	workspaceFn := func() (string, error) {
		if d.workspace == "" {
			return "", fmt.Errorf("workspace is not initialized")
		}
		return d.workspace, nil
	}

	cmd.AddCommand(
		newCreateCommand(managerFn, workspaceFn),
		newListCommand(managerFn),
		newStatusCommand(managerFn),
		newDissolveCommand(managerFn),
		newMemoryCommand(memoryFn),
		newExecuteCommand(managerFn),
	)

	return cmd
}

