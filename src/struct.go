package main

type tConfig struct {
	SyncSteps []tSyncSteps `yaml:"sync_steps"`
	Tunnel    tTunnel      `yaml:"tunnel"`
}

type tSyncSteps struct {
	Cmd    tCmd    `yaml:"cmd"`
	Local  tLocal  `yaml:"local"`
	Remote tRemote `yaml:"remote"`
}

type tCmd struct {
	Base     string   `yaml:"base"`
	Args     []string `yaml:"args"`
	DebugArg string   `yaml:"debug_arg"`
}

type tLocal struct {
	Folder string `yaml:"folder"`
}

type tRemote struct {
	Host   string `yaml:"host"`
	Folder string `yaml:"folder"`
	User   string `yaml:"user"`
	Group  string `yaml:"group"`
}

type tTunnel struct {
	LocalPort  int `yaml:"local_port"`
	RemotePort int `yaml:"remote_port"`
}
