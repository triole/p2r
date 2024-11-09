package main

import (
	"encoding/json"
	"fmt"

	"github.com/triole/logseal"
)

var (
	lg logseal.Logseal
)

func main() {
	parseArgs()
	lg = logseal.Init(cli.LogLevel, cli.LogFile, cli.LogNoColors, cli.LogJSON)
	confFile := pabs(cli.Config)
	lg.Info(
		"start "+appName,
		logseal.F{
			"command":   cli.Action,
			"config":    confFile,
			"dry-run":   cli.DryRun,
			"log-level": cli.LogLevel,
		})
	conf := readConfig(confFile)
	lg.Debug("read config", logseal.F{"config": fmt.Sprintf("%+v", conf)})

	switch cli.Action {
	case "pull", "push":
		runSync(conf.SyncSteps)
	case "list":
		list(conf.SyncSteps)
	default:
		lg.Info("display read config file")
		displayInfo(conf)
	}
}

func displayInfo(conf tConf) {
	s, _ := json.MarshalIndent(conf, "", "  ")
	fmt.Println(string(s))
}
