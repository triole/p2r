package conf

import "fmt"

func (conf Conf) ListAvailableCommands() {
	conf.Lg.Info("the following commands are available")
	configContent := conf.expand()
	for key, vals := range configContent.Commands {
		fmt.Printf("%v\n", key)
		for _, val := range vals {
			fmt.Printf("  %v\n", val)
		}
	}
}
