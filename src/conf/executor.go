package conf

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"p2r/src/logger"
)

// Executor handles command execution.
type Executor struct {
	Lg        *logger.Logger
	PrintOnly bool
	Action    string
}

// NewExecutor creates a new Executor instance.
func NewExecutor(lg *logger.Logger, printOnly bool, action string) *Executor {
	return &Executor{
		Lg:        lg,
		PrintOnly: printOnly,
		Action:    action,
	}
}

// ExecuteCommands runs a list of commands.
func (e *Executor) ExecuteCommands(commands []ExecCommand, dryRun bool, rsyncDryRun bool) error {
	if len(commands) == 0 {
		e.Lg.Error("no commands to execute", nil)
		return fmt.Errorf("no commands provided")
	}

	var lastErr error
	for i, cmd := range commands {
		// Apply dry-run flag for rsync if enabled
		cmdCopy := make([]string, len(cmd.Cmd))
		copy(cmdCopy, cmd.Cmd)

		if rsyncDryRun && len(cmdCopy) > 0 && cmdCopy[0] == "rsync" {
			cmdCopy = append(cmdCopy, "--dry-run")
		}

		if len(cmd.Err) > 0 {
			for _, err := range cmd.Err {
				e.Lg.Warn("skipping command with errors", map[string]interface{}{
					"index": i,
					"cmd":   fmt.Sprintf("%+v", cmd.Cmd),
					"error": err,
				})
			}
			continue
		}

		if _, exitCode, err := e.runCommand(cmdCopy); err != nil {
			e.Lg.Error("command execution failed", map[string]interface{}{
				"index":     i,
				"cmd":       fmt.Sprintf("%+v", cmdCopy),
				"exit_code": exitCode,
				"error":     err,
				"action":    e.Action,
				"dry_run":   dryRun,
			})
			lastErr = err
		}
	}

	return lastErr
}

// runCommand executes a single command.
// Output streams directly to stdout; no buffering to avoid memory issues with large outputs.
func (e *Executor) runCommand(cmdArr []string) ([]byte, int, error) {
	if len(cmdArr) == 0 {
		return nil, 0, fmt.Errorf("empty command array")
	}

	e.Lg.Info("executing command", map[string]interface{}{
		"action": e.Action,
		"cmd":    fmt.Sprintf("%+v", cmdArr),
	})

	if e.PrintOnly {
		output := fmt.Sprintf("WOULD EXECUTE: %s %s\n",
			cmdArr[0],
			strings.Join(cmdArr[1:], " "))
		return []byte(output), 0, nil
	}

	cmd := exec.Command(cmdArr[0], cmdArr[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	var exitCode int
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				exitCode = status.ExitStatus()
			}
		}
		e.Lg.Debug("command failed", map[string]interface{}{
			"action": e.Action,
			"cmd":    fmt.Sprintf("%+v", cmdArr),
			"error":  err,
			"exit":   exitCode,
		})
		return nil, exitCode, err
	}

	e.Lg.Debug("command completed successfully", map[string]interface{}{
		"action": e.Action,
		"cmd":    fmt.Sprintf("%+v", cmdArr),
		"exit":   exitCode,
	})

	return nil, exitCode, nil
}

// ListCommands displays all available commands from config.
func (e *Executor) ListCommands(configContent ConfigContent) {
	e.Lg.Info("available commands:", nil)
	for key, commands := range configContent.Commands {
		fmt.Printf("%s:\n", key)
		for _, cmd := range commands {
			fmt.Printf("  %v\n", cmd)
		}
	}
}
