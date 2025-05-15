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
	by, err = conf.templateFile(string(by))
	err = yaml.Unmarshal(by, &configContent)
	conf.Lg.IfErrFatal(
		"can not unmarshal config", logseal.F{"path": conf.ConfigFile, "error": err},
	)
	return
}

func (conf *Conf) assembleCommands(configContent ConfigContent) (commands Commands) {
	switch conf.Action {
	case "pull", "push":
		for _, step := range configContent.SyncSteps {
			conf.Commands = append(conf.Commands, conf.assembleSyncCommand(step))
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

func (conf *Conf) assembleSyncCommand(step SyncStep) (cmd Command) {
	local := conf.parsePath(step.Local)
	remote := conf.parsePath(step.Remote)
	if step.Cmd[0] == "rsync" && conf.RsyncDryRun {
		step.Cmd = append(step.Cmd, "--dry-run")
	}
	switch conf.Action {
	case "pull":
		step.Cmd = append(step.Cmd, remote.FullPath)
		cmd.Cmd = append(step.Cmd, local.FullPath)
		cmd.Err = append(cmd.Err, remote.Errors...)
		cmd.Err = append(cmd.Err, local.Errors...)
	case "push":
		step.Cmd = append(step.Cmd, local.FullPath)
		cmd.Cmd = append(step.Cmd, remote.FullPath)
		cmd.Err = append(cmd.Err, local.Errors...)
		cmd.Err = append(cmd.Err, remote.Errors...)
	}
	return
}
