package main

import (
	"os"
	"os/exec"
)

// runCommand is the execution function.
// The actual command execution is now handled by conf.Executor.
// This function is kept for backward compatibility but should be removed.
func runCommand(cmdArr []string, printOnly bool) ([]byte, int, error) {
	if printOnly {
		return nil, 0, nil
	}

	cmd := exec.Command(cmdArr[0], cmdArr[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return nil, 0, cmd.Run()
}
