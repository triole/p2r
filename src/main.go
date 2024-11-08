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
	lg.Info(
		"start "+appName,
		logseal.F{
			"command":   cli.Action,
			"config":    cli.Config,
			"dry-run":   cli.DryRun,
			"log-level": cli.LogLevel,
		})
	conf := readConfig(cli.Config)
	lg.Debug("read config", logseal.F{"config": fmt.Sprintf("%+v", conf)})

	switch cli.Action {
	case "tun":
		fmt.Printf("%+v\n", "to be implemented")
	case "pull", "push":
		runSync(conf.SyncSteps)
	default:
		lg.Info("display read config file")
		displayInfo(conf)
	}
}

func displayInfo(conf tConf) {
	s, _ := json.MarshalIndent(conf, "", "  ")
	fmt.Println(string(s))
}
