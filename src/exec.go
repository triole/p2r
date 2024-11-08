package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"syscall"

	"github.com/triole/logseal"
)

func runCmd(cmdBase string, cmdArgs []string) ([]byte, int, error) {
	var err error
	var exitcode int
	var stdBuffer bytes.Buffer

	cmd := exec.Command(cmdBase, cmdArgs...)
	// mw := io.MultiWriter(&stdBuffer)
	mw := io.MultiWriter(os.Stdout, &stdBuffer)

	cmd.Stdout = mw
	cmd.Stderr = mw
	if err = cmd.Run(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// the program has exited with an exit code != 0
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				exitcode = status.ExitStatus()
			}
		}
	}
	by := stdBuffer.Bytes()
	if err != nil {
		lg.IfErrError(
			"exec failed",
			logseal.F{
				"cmd_base": cmdBase, "cmd_args": cmdArgs,
				"error": err, "output": string(by),
			},
		)
	} else {
		lg.Debug(
			"exec successful",
			logseal.F{
				"cmd_base": cmdBase, "cmd_args": cmdArgs,
				"error": err, "output": string(by),
			},
		)
	}
	return by, exitcode, err
}
