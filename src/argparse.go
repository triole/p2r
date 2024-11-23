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
	BUILDTAGS      string
	appName        = "p2r"
	appDescription = "push or pull remote file systems"
	appMainversion = "0.1"
)

var cli struct {
	Action      string `kong:"-" enum:"pull,push,cmd,info" default:"info"`
	Config      string `help:"config file" default:"${configFile}" short:"c"`
	DryRun      bool   `help:"only print commands what would have been executed" short:"n"`
	RsyncDryRun bool   `help:"enable rsync dry runs" short:"m"`
	LogFile     string `help:"log file" default:"/dev/stdout"`
	LogLevel    string `help:"log level, can be: [${enum}]" default:"info" enum:"trace,debug,info,error"`
	LogNoColors bool   `help:"disable output colours, print plain text"`
	LogJSON     bool   `help:"enable json log, instead of text one"`
	VersionFlag bool   `help:"display version" short:"V"`

	Pull struct {
		Plain bool `help:"print plain list, file names only" short:"p"`
	} `cmd:"" help:"list files matching the criteria"`

	Push struct {
		Plain bool `help:"print plain list, file names only" short:"p"`
	} `cmd:"" help:"list files matching the criteria"`

	Cmd struct {
		Command string `help:"run a command defined in the config yaml" arg:""`
	} `cmd:"" help:"run a command defined in the config yaml"`

	Init struct {
		Plain bool `help:"init a config template to customise" short:"i"`
	} `cmd:"" help:"init a config template to customise" short:"i"`

	Version struct{} `cmd:"" help:"display version"`
}

func parseArgs() {
	curdir := pwd()
	ctx := kong.Parse(&cli,
		kong.Name(appName),
		kong.Description(appDescription),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact:      true,
			Summary:      true,
			NoAppSummary: true,
			FlagsLast:    false,
		}),
		kong.Vars{
			"configFile": filepath.Join(curdir, "p2r.yaml"),
		},
	)
	_ = ctx.Run()
	// ctx.FatalIfErrorf(err)
	cli.Action = strings.Split(ctx.Command(), " ")[0]
	if cli.Action == "version" {
		printBuildTags(BUILDTAGS)
		os.Exit(0)
	}
}

func printBuildTags(buildtags string) {
	regexp, _ := regexp.Compile(`({|}|,)`)
	s := regexp.ReplaceAllString(buildtags, "\n")
	s = strings.Replace(s, "_subversion: ", "Version: "+appMainversion+".", -1)
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
