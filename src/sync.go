package main

import (
	"errors"
	"fmt"

	"github.com/triole/logseal"
)

func runSync(steps tSyncSteps) {
	for _, step := range steps {
		cmdArr, err := assembleCommand(step)
		if len(err) == 0 {
			_, _, err := runCmd(cmdArr)
			lg.IfErrFatal("sync failed", logseal.F{"error": err})
		} else {
			for _, el := range err {
				lg.Error(
					"skip command", logseal.F{
						"cmd": fmt.Sprintf("%+v", cmdArr), "error": el,
					},
				)
			}
		}
	}
}

func assembleCommand(step tSyncStep) (cmdArr []string, errArr []error) {
	var source, target string
	cmdArr = step.Cmd
	if cmdArr[0] == "rsync" && cli.RsyncDryRun {
		cmdArr = append(cmdArr, "-n")
	}
	if cli.Action == "pull" {
		source = step.Remote
		target = step.Local
	}
	if cli.Action == "push" {
		source = step.Local
		target = step.Remote
	}
	errArr = isHealthy(source, "source", cli.Action, errArr)
	errArr = isHealthy(target, "target", cli.Action, errArr)
	cmdArr = append(cmdArr, source)
	cmdArr = append(cmdArr, target)
	return
}

func isHealthy(path, typ, action string, errArr []error) []error {
	if rxMatch(path, "^[:/-_]$") {
		errArr = append(
			errArr, errors.New(action+" "+typ+" path seems short: "+path),
		)
	}
	if isLocalPath(path) {
		if isFolder(path) {
			b, _ := isEmpty(path)
			if b {
				errArr = append(
					errArr, errors.New(action+" "+typ+" folder empty: "+path),
				)
			}
		}
	}
	return errArr
}
