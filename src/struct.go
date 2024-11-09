package main

type tConf struct {
	SyncSteps tSyncSteps `yaml:"sync_steps"`
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

func (pth tPath) isFolder() (b bool) {
	b = false
	switch val := pth.IsFolder.(type) {
	case bool:
		if val {
			b = true
		}
	}
	return
}

func (pth tPath) isEmpty() (b bool) {
	b = false
	switch val := pth.IsEmpty.(type) {
	case bool:
		if val {
			b = true
		}
	}
	return
}
