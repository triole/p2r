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
	conf.expand()
	return conf
}

func (conf Conf) expand() {
	configFile := conf.pabs(conf.ConfigFile)
	by, err := os.ReadFile(configFile)
	conf.Lg.IfErrFatal(
		"can not read file", logseal.F{"path": configFile, "error": err},
	)
	by, err = conf.templateFile(string(by))
	err = yaml.Unmarshal(by, &conf)
	conf.Lg.IfErrFatal(
		"can not unmarshal config", logseal.F{"path": conf.ConfigFile, "error": err},
	)
	for idx, step := range conf.SyncSteps {
		conf.SyncSteps[idx].Set.Local = conf.parsePath(step.Local)
		conf.SyncSteps[idx].Set.Remote = conf.parsePath(step.Remote)
		conf.SyncSteps[idx].Set.Command, conf.SyncSteps[idx].Set.Errors = conf.assembleCommand(step)
	}
	return
}
