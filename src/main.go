package main

import (
	"fmt"

	"p2r/src/conf"

	_ "embed"

	"github.com/triole/logseal"
)

var (
	lg logseal.Logseal
)

func main() {
	parseArgs()
	lg = logseal.Init(cli.LogLevel, cli.LogFile, cli.LogNoColors, cli.LogJSON)
	conf := conf.Init(
		cli.Config, cli.Action, cli.Cmd.Command,
		cli.DryRun, cli.RsyncDryRun, lg,
	)
	conf.Lg.Info(
		"start "+appName,
		logseal.F{
			"config":    conf.ConfigFile,
			"action":    conf.Action,
			"dry-run":   cli.DryRun,
			"log-level": cli.LogLevel,
		})
	conf.Lg.Debug("read config", logseal.F{"config": fmt.Sprintf("%+v", conf)})

	if conf.Action == "init" {
		conf.InitExample()
	} else {
		runCommands(conf.Commands)
	}
}
