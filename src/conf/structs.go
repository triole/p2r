package conf

// ConfigContent represents the parsed YAML configuration.
type ConfigContent struct {
	SyncSteps SyncSteps             `yaml:"sync_steps"`
	Commands  map[string][][]string `yaml:"commands"`
}

// SyncSteps is a list of synchronization steps.
type SyncSteps []SyncStep

// SyncStep defines a single synchronization step.
type SyncStep struct {
	Cmd     []string `yaml:"cmd"`
	Sources []string `yaml:"sources"`
	Targets []string `yaml:"targets"`
}

// Commands is a list of commands.
type Commands []ExecCommand

// ExecCommand represents a command to be executed.
type ExecCommand struct {
	Cmd []string
	Err []error
}

// ParsedPath represents a validated path.
type ParsedPath struct {
	IsLocal   bool
	Machine   string
	Path      string
	FullPath  string
	Exists    bool
	IsFolder  bool
	IsEmpty   bool
	IsHealthy bool
	Errors    []error
}
