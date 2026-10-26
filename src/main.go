package main

import (
	"fmt"

	"p2r/src/conf"
	"p2r/src/logger"
)

var (
	lg *logger.Logger
)

func main() {
	parseArgs()

	// Initialize logger with zero dependencies
	lg = logger.NewLogger(cli.LogLevel, cli.LogFile, cli.LogNoColors)
	defer lg.Close() // Ensure file handle is closed on exit

	// Build configuration
	cfg := conf.NewConfig(
		cli.Config,
		cli.Action,
		cli.Cmd.Command,
		cli.PrintOnly,
		cli.RsyncDryRun,
		lg,
	)

	// Log startup information
	lg.Info("starting "+AppName, map[string]interface{}{
		"config":    cfg.ConfigFile,
		"action":    cfg.Action,
		"dry-run":   cli.PrintOnly,
		"log-level": cli.LogLevel,
	})

	// Handle different command modes
	if cli.Cmd.List {
		cfg.ListAvailableCommands()
		return
	}

	switch cfg.Action {
	case ActionInit:
		conf.InitExample(cli.Config, lg)
		return

	case ActionVersion:
		printBuildTags(BUILDTAGS)
		return

	default:
		// Execute the action
		if err := executeAction(cfg); err != nil {
			cfg.Lg.Fatal("execution failed", map[string]interface{}{"error": err})
		}
	}
}

// executeAction handles the main execution flow for push/pull/cmd actions.
func executeAction(cfg *conf.Config) error {
	// Load configuration
	configContent, err := cfg.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	cfg.Lg.Debug("config loaded", map[string]interface{}{"config": fmt.Sprintf("%+v", configContent)})

	// Build commands based on action
	var commands []conf.ExecCommand

	switch cfg.Action {
	case ActionPull, ActionPush:
		commands = buildSyncCommands(configContent, cfg.Action, cfg.RsyncDryRun)

	case ActionCmd:
		if cfg.SubAction == "" {
			return fmt.Errorf("no sub-action specified for cmd")
		}
		if cmdList, ok := configContent.Commands[cfg.SubAction]; ok {
			for _, cmd := range cmdList {
				commands = append(commands, conf.ExecCommand{Cmd: cmd})
			}
		} else {
			return fmt.Errorf("command '%s' not found in config", cfg.SubAction)
		}

	default:
		return fmt.Errorf("unknown action: %s", cfg.Action)
	}

	// Execute commands
	err = cfg.Executor.ExecuteCommands(commands, cfg.DryRun, cfg.RsyncDryRun)
	return err
}

// buildSyncCommands builds the command list for pull/push actions.
func buildSyncCommands(configContent conf.ConfigContent, action string, rsyncDryRun bool) []conf.ExecCommand {
	var commands []conf.ExecCommand

	for _, step := range configContent.SyncSteps {
		sources := step.Sources
		targets := step.Targets

		// Swap sources/targets for pull action
		if action == ActionPull {
			sources, targets = targets, sources
		}

		// Build command base
		cmdBase := make([]string, len(step.Cmd))
		copy(cmdBase, step.Cmd)

		// Add --dry-run flag for rsync if enabled
		if rsyncDryRun && len(cmdBase) > 0 && cmdBase[0] == "rsync" {
			cmdBase = append(cmdBase, "--dry-run")
		}

		// Create cross-product of sources × targets
		for _, source := range sources {
			for _, target := range targets {
				cmd := conf.ExecCommand{
					Cmd: append(append(cmdBase, source), target),
				}
				commands = append(commands, cmd)
			}
		}
	}

	return commands
}
