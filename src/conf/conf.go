package conf

import (
	"p2r/src/logger"
)

// Config represents the main configuration wrapper.
type Config struct {
	ConfigFile  string
	Action      string
	SubAction   string
	DryRun      bool
	RsyncDryRun bool
	Lg          *logger.Logger
	Parser      *Parser
	Validator   *Validator
	Executor    *Executor
}

// NewConfig creates a new configuration instance.
func NewConfig(
	configFile string,
	action string,
	subAction string,
	dryRun bool,
	rsyncDryRun bool,
	logger *logger.Logger,
) *Config {

	parser := NewParser(configFile, logger)
	validator := NewValidator(logger)
	executor := NewExecutor(logger, false, action) // PrintOnly set later

	return &Config{
		Parser:      parser,
		Validator:   validator,
		Executor:    executor,
		ConfigFile:  configFile,
		Action:      action,
		SubAction:   subAction,
		DryRun:      dryRun,
		RsyncDryRun: rsyncDryRun,
		Lg:          logger,
	}
}

// SetPrintOnly updates executor's print-only mode.
func (c *Config) SetPrintOnly(printOnly bool) {
	c.Executor.PrintOnly = printOnly
	c.Parser.Lg = c.Lg
	c.Validator.Lg = c.Lg
	c.Executor.Lg = c.Lg
}

// Load reads and parses the configuration file.
func (c *Config) Load() (ConfigContent, error) {
	return c.Parser.Parse()
}
