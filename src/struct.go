package main

type tConf struct {
	SyncSteps tSyncSteps `yaml:"sync_steps"`
}

type tSyncSteps []tSyncStep

type tSyncStep struct {
	Cmd    []string `yaml:"cmd"`
	Local  string   `yaml:"local"`
	Remote string   `yaml:"remote"`
}
