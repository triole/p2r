package main

import (
	"strings"

	"github.com/triole/logseal"
)

func runSync(steps tSyncSteps) {
	for _, step := range steps {
		cmdBase, cmdArgs := assembleCommand(step)
		lg.Info("sync", logseal.F{"action": cli.Action, "cmd_base": cmdBase, "cmd_args": cmdArgs})
		if !cli.Print {
			_, _, err := runCmd(cmdBase, cmdArgs)
			lg.IfErrFatal("sync failed", logseal.F{"error": err})
		}
	}
}

func assembleCommand(step tSyncStep) (cmdBase string, cmdArgs []string) {
	cmdBase = step.Cmd.Base
	if cli.DryRun {
		cmdArgs = append(cmdArgs, step.Cmd.DebugArg)
	}
	cmdArgs = append(cmdArgs, step.Cmd.Args...)
	if cli.Action == "pull" {
		cmdArgs = append(cmdArgs, makeRemotePath(step.Remote))
		cmdArgs = append(cmdArgs, step.Local.Folder)
	}
	if cli.Action == "push" {
		cmdArgs = append(cmdArgs, step.Local.Folder)
		cmdArgs = append(cmdArgs, makeRemotePath(step.Remote))
	}
	return
}

func makeRemotePath(remote tRemote) string {
	return strings.Join([]string{remote.Host, remote.Folder}, ":")
}
