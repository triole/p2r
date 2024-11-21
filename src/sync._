package main

import (
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
