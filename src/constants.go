package main

// Application constants.
const (
	AppName        = "p2r"
	AppDescription = "push or pull remote file systems"
	AppVersion     = "0.1"
	DefaultConfig  = "p2r.yaml"
)

// Action constants.
const (
	ActionPull    = "pull"
	ActionPush    = "push"
	ActionCmd     = "cmd"
	ActionInfo    = "info"
	ActionInit    = "init"
	ActionVersion = "version"
)

// Log level constants.
const (
	LogLevelTrace = "trace"
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelError = "error"
)

// Log level validation message.
const LogLevelEnum = LogLevelTrace + "," + LogLevelDebug + "," + LogLevelInfo + "," + LogLevelError
