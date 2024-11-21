package conf

import "github.com/triole/logseal"

type Conf struct {
	ConfigFile  string
	Action      string
	DryRun      bool
	RsyncDryRun bool
	Lg          logseal.Logseal
	SyncSteps   tSyncSteps `yaml:"sync_steps"`
	Commands    tCommands  `yaml:"commands"`
}

type tSyncSteps []tSyncStep

type tSyncStep struct {
	Cmd    []string `yaml:"cmd"`
	Local  string   `yaml:"local"`
	Remote string   `yaml:"remote"`
	Set    tSet
}

type tSet struct {
	Local   tPath
	Remote  tPath
	Command []string
	Errors  []error
}

type tPath struct {
	IsLocal   bool
	Machine   string
	Path      string
	FullPath  string
	IsFolder  interface{}
	IsEmpty   interface{}
	IsHealthy bool
	Errors    []error
}

type tCommands map[string][][]string
