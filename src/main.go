package main

import (
	"fmt"

	"github.com/triole/logseal"
)

var (
	lg logseal.Logseal
)

func main() {
	parseArgs()
	lg = logseal.Init(CLI.LogLevel, CLI.LogFile, CLI.LogNoColors, CLI.LogJSON)
	lg.Info(
		"start "+appName,
		logseal.F{
			"command":   CLI.Cmd,
			"config":    CLI.Config,
			"dry-run":   CLI.DryRun,
			"log-level": CLI.LogLevel,
		})
	conf := readConfig(CLI.Config)
	lg.Debug("read config", logseal.F{"config": fmt.Sprintf("%+v", conf)})
}
