package conf

import "github.com/triole/logseal"

type Conf struct {
	ConfigFile  string
	Action      string
	SubAction   string
	DryRun      bool
	RsyncDryRun bool
	Lg          logseal.Logseal
	Commands    Commands
}

type Commands []Command

type Command struct {
	Cmd []string
	Err []error
}

type ConfigContent struct {
	SyncSteps SyncSteps             `yaml:"sync_steps"`
	Commands  map[string][][]string `yaml:"commands"`
}
type SyncSteps []SyncStep

type SyncStep struct {
	Cmd     []string `yaml:"cmd"`
	Sources []string `yaml:"sources"`
	Targets []string `yaml:"targets"`
	Set     Set
}

type Set struct {
	Local   Path
	Remote  Path
	Command []string
	Errors  []error
}

type Path struct {
	IsLocal   bool
	Machine   string
	Path      string
	FullPath  string
	Exists    bool
	IsFolder  interface{}
	IsEmpty   interface{}
	IsHealthy bool
	Errors    []error
}
