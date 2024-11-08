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
	lg = logseal.Init(cli.LogLevel, cli.LogFile, cli.LogNoColors, cli.LogJSON)
	lg.Info(
		"start "+appName,
		logseal.F{
			"command":   cli.Cmd,
			"config":    cli.Config,
			"dry-run":   cli.DryRun,
			"log-level": cli.LogLevel,
		})
	conf := readConfig(cli.Config)
	lg.Debug("read config", logseal.F{"config": fmt.Sprintf("%+v", conf)})

	switch cli.Cmd {
	case "tun":
		fmt.Printf("%+v\n", "TUNNEL")
	default:
		runSync(conf.SyncSteps)
	}
}
