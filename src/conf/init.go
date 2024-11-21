package conf

import (
	"os"

	"github.com/triole/logseal"
	yaml "gopkg.in/yaml.v3"
)

func Init(filename, action string, dryRun, rsyncDryRun bool, lg logseal.Logseal) (conf Conf) {
	conf.ConfigFile = conf.pabs(filename)
	conf.Action = action
	conf.DryRun = dryRun
	conf.RsyncDryRun = rsyncDryRun
	conf.Lg = lg
	configContent := conf.expand()
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
	}
	return
}

func (conf *Conf) assembleSyncCommand(step SyncStep) (cmd Command) {
	local := conf.parsePath(step.Local)
	remote := conf.parsePath(step.Remote)
	switch conf.Action {
	case "pull":
		cmd.Cmd = append(step.Cmd, remote.FullPath)
		cmd.Cmd = append(step.Cmd, local.FullPath)
		cmd.Err = append(cmd.Err, remote.Errors...)
		cmd.Err = append(cmd.Err, local.Errors...)
	case "push":
		cmd.Cmd = append(step.Cmd, local.FullPath)
		cmd.Cmd = append(step.Cmd, remote.FullPath)
		cmd.Err = append(cmd.Err, local.Errors...)
		cmd.Err = append(cmd.Err, remote.Errors...)
	}
	return
}
