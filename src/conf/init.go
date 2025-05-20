package conf

import (
	"os"

	"github.com/triole/logseal"
	yaml "gopkg.in/yaml.v3"
)

func Init(filename, action, subaction string, dryRun, rsyncDryRun bool, lg logseal.Logseal) (conf Conf) {
	conf.ConfigFile = conf.pabs(filename)
	conf.Action = action
	conf.SubAction = subaction
	conf.DryRun = dryRun
	conf.RsyncDryRun = rsyncDryRun
	conf.Lg = lg
	var configContent ConfigContent
	if conf.Action == "init" {
		conf.InitExample()
		os.Exit(0)
	} else {
		configContent = conf.expand()
	}
	conf.assembleCommands(configContent)
	return conf
}

func (conf *Conf) expand() (configContent ConfigContent) {
	by, err := os.ReadFile(conf.ConfigFile)
	conf.Lg.IfErrFatal(
		"can not read file", logseal.F{"path": conf.ConfigFile, "error": err},
	)
	by, _ = conf.templateFile(string(by))
	err = yaml.Unmarshal(by, &configContent)
	conf.Lg.IfErrFatal(
		"can not unmarshal config", logseal.F{"path": conf.ConfigFile, "error": err},
	)
	conf.Lg.Debug(configContent)
	return
}

func (conf *Conf) assembleCommands(configContent ConfigContent) (commands Commands) {
	switch conf.Action {
	case "pull", "push":
		for _, step := range configContent.SyncSteps {
			conf.Commands = append(conf.Commands, conf.assembleSyncCommands(step)...)
		}
	case "cmd":
		if commands, ok := configContent.Commands[conf.SubAction]; ok {
			for _, command := range commands {
				var cmd Command
				cmd.Cmd = command
				conf.Commands = append(conf.Commands, cmd)
			}
		}
	}
	return
}

func (conf *Conf) assembleSyncCommands(step SyncStep) (cmds []Command) {
	sources := conf.parsePathList(step.Sources)
	targets := conf.parsePathList(step.Targets)

	if conf.Action == "pull" {
		temp := sources
		sources = targets
		targets = temp
	}

	if step.Cmd[0] == "rsync" && conf.RsyncDryRun {
		step.Cmd = append(step.Cmd, "--dry-run")
	}
	for _, source := range sources {
		for _, target := range targets {
			cmd := Command{}
			cmd.Cmd = append(cmd.Cmd, step.Cmd...)
			cmd.Cmd = append(cmd.Cmd, source.FullPath)
			cmd.Cmd = append(cmd.Cmd, target.FullPath)
			cmd.Err = append(cmd.Err, source.Errors...)
			cmd.Err = append(cmd.Err, target.Errors...)
			cmds = append(cmds, cmd)
		}
	}
	return
}
