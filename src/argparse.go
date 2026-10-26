package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alecthomas/kong"
)

var (
	BUILDTAGS string
)

var cli struct {
	Action      string `kong:"-" enum:"pull,push,cmd,info" default:"info"`
	Config      string `help:"config file" default:"${configFile}" short:"c"`
	PrintOnly   bool   `help:"only print commands what would have been executed" short:"p"`
	RsyncDryRun bool   `help:"enable rsync dry runs" short:"n"`
	LogFile     string `help:"log file" default:"/dev/stdout"`
	LogLevel    string `help:"log level, can be: [${enum}]" default:"info" enum:"trace,debug,info,error"`
	LogNoColors bool   `help:"disable output colours, print plain text"`
	VersionFlag bool   `help:"display version" short:"V"`

	Pull struct {
	} `cmd:"" help:"list files matching the criteria"`

	Push struct {
	} `cmd:"" help:"list files matching the criteria"`

	Cmd struct {
		Command string `help:"run a command defined in the config yaml" optional:"" arg:""`
		List    bool   `help:"list available commands" short:"l"`
	} `cmd:"" help:"run a command defined in the config yaml"`

	Init struct {
	} `cmd:"" help:"init a config template to customise" short:"i"`

	Version struct{} `cmd:"" help:"display version"`
}

func parseArgs() {
	curdir := pwd()
	ctx := kong.Parse(&cli,
		kong.Name(AppName),
		kong.Description(AppDescription),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact:      true,
			Summary:      true,
			NoAppSummary: true,
			FlagsLast:    false,
		}),
		kong.Vars{
			"configFile": filepath.Join(curdir, DefaultConfig),
			"enum":       LogLevelEnum,
		},
	)
	_ = ctx.Run()
	// ctx.FatalIfErrorf(err)
	cli.Action = strings.Split(ctx.Command(), " ")[0]
	if cli.Action == ActionVersion {
		printBuildTags(BUILDTAGS)
		os.Exit(0)
	}
}

func printBuildTags(buildtags string) {
	regexp, _ := regexp.Compile(`({|}|,)`)
	s := regexp.ReplaceAllString(buildtags, "\n")
	s = strings.Replace(s, "_subversion: ", "Version: "+AppVersion+".", -1)
	fmt.Printf("%s\n", s)
}

func pwd() string {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return pwd
}
