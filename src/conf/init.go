package conf

import (
	"os"

	"p2r/src/logger"

	_ "embed"
)

//go:embed example.yaml
var content string

// InitExample creates an example configuration file.
func InitExample(configFile string, logger *logger.Logger) error {
	exists, _ := pathExists(configFile)
	if !exists {
		logger.Info("creating example config", map[string]interface{}{"path": configFile})
		err := os.WriteFile(configFile, []byte(content), 0600)
		if err != nil {
			logger.Error("failed to create example config", map[string]interface{}{"error": err})
			return err
		}
	} else {
		logger.Info("config file already exists, skipping creation", map[string]interface{}{"path": configFile})
	}

	return nil
}

// pathExists checks if a file exists.
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// ListAvailableCommands displays all available commands from config.
func (cfg *Config) ListAvailableCommands() {
	configContent, err := cfg.Load()
	if err != nil {
		cfg.Lg.Error("failed to load config", map[string]interface{}{"error": err})
		return
	}
	cfg.Executor.ListCommands(configContent)
}
