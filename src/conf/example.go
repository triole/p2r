package conf

import (
	_ "embed"
	"os"

	"github.com/triole/logseal"
)

//go:embed p2r.yaml
var exampleConf string

func (conf Conf) InitExample() {
	b, _ := conf.exists(conf.ConfigFile)
	if !b {
		err := os.WriteFile(conf.ConfigFile, []byte(exampleConf), 0644)
		if err != nil {
			conf.Lg.Error(
				"could not create example conf",
				logseal.F{"error": err, "path": conf.ConfigFile},
			)
		}
	} else {
		conf.Lg.Info(
			"did not create example conf, file exists",
			logseal.F{"path": conf.ConfigFile},
		)
	}
}
