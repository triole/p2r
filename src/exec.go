package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"

	"github.com/triole/logseal"
)

func runCmd(cmdArr []string) ([]byte, int, error) {
	var by []byte
	var err error
	var exitcode int
	var stdBuffer bytes.Buffer
	lg.Info(
		"run command",
		logseal.F{
			"action": cli.Action, "cmd": fmt.Sprintf("%+v", cmdArr),
		},
	)
	if !cli.DryRun {
		cmd := exec.Command(cmdArr[0], cmdArr[1:]...)
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
		by = stdBuffer.Bytes()
		if err != nil {
			lg.IfErrError(
				"exec failed",
				logseal.F{
					"action": cli.Action, "cmd": fmt.Sprintf("%+v", cmdArr),
					"error": err, "output": string(by),
				},
			)
		} else {
			lg.Debug(
				"exec successful",
				logseal.F{
					"action": cli.Action, "cmd": fmt.Sprintf("%+v", cmdArr),
					"error": err, "output": string(by),
				},
			)
		}
	}
	return by, exitcode, err
}
