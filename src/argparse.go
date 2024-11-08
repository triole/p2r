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
	appDescription = "push, pull or tun to remote file systems"
	appMainversion = "0.1"
)

var cli struct {
	Cmd         string `help:"command to run, can be: [${enum}]" arg:"" enum:"info,push,pull,tun" default:"info"`
	Config      string `help:"config file" default:"${configFile}" short:"f"`
	Print       bool   `help:"only print commands that would have been executed" short:"p"`
	DryRun      bool   `help:"execute sync dry run syncs" short:"n"`
	LogFile     string `help:"log file" default:"/dev/stdout"`
	LogLevel    string `help:"log level, can be: [${enum}]" default:"info" enum:"trace,debug,info,error"`
	LogNoColors bool   `help:"disable output colours, print plain text"`
	LogJSON     bool   `help:"enable json log, instead of text one"`
	VersionFlag bool   `help:"display version" short:"V"`
}

func parseArgs() {
	curdir := pwd()
	_ = kong.Parse(&cli,
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
	// err := ctx.Run()
	// ctx.FatalIfErrorf(err)

	if cli.VersionFlag {
		printBuildTags(BUILDTAGS)
		os.Exit(0)
	}
	// ctx.FatalIfErrorf(err)
}

func printBuildTags(buildtags string) {
	regexp, _ := regexp.Compile(`({|}|,)`)
	s := regexp.ReplaceAllString(buildtags, "\n")
	s = strings.Replace(s, "_subversion: ", "Version: "+appMainversion+".", -1)
	fmt.Printf("%s\n", s)
}
