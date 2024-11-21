package main

import (
	"fmt"

	"p2r/src/conf"

	"github.com/triole/logseal"
)

var (
	lg logseal.Logseal
)

func main() {
	parseArgs()
	lg = logseal.Init(cli.LogLevel, cli.LogFile, cli.LogNoColors, cli.LogJSON)
	conf := conf.Init(cli.Config, cli.Action, cli.DryRun, cli.RsyncDryRun, lg)
	fmt.Printf("%+v\n", cli.Action)
	lg.Info(
		"start "+appName,
		logseal.F{
			"config":    conf.ConfigFile,
			"action":    conf.Action,
			"dry-run":   cli.DryRun,
			"log-level": cli.LogLevel,
		})
	fmt.Printf("%+v\n", conf)
	// lg.Debug("read config", logseal.F{"config": fmt.Sprintf("%+v", conf)})

	// switch cli.Action {
	// case "pull", "push":
	// 	runSync(conf.SyncSteps)
	// case "list":
	// 	list(conf.SyncSteps)
	// case "cmd":
	// 	fmt.Printf("%+v\n", conf)
	// 	// runCommand(conf.Commands)
	// default:
	// 	lg.Info("display read config file")
	// 	displayInfo(conf)
	// }
}

// func runCommand(commands [][]string) {
// 	for _, el := range commands {
// 		fmt.Printf("%+v\n", el)
// 	}
// }

// func displayInfo(conf tConf) {
// 	s, _ := json.MarshalIndent(conf, "", "  ")
// 	fmt.Println(string(s))
// }
