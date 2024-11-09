package main

import (
	"bytes"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/triole/logseal"
	yaml "gopkg.in/yaml.v3"
)

func readConfig(filename string) (conf tConf) {
	by, err := os.ReadFile(filename)
	lg.IfErrFatal(
		"can not read file", logseal.F{"path": filename, "error": err},
	)
	by, err = templateFile(string(by))
	err = yaml.Unmarshal(by, &conf)
	lg.IfErrFatal(
		"can not unmarshal config", logseal.F{"path": filename, "error": err},
	)
	for idx, step := range conf.SyncSteps {
		conf.SyncSteps[idx].Set.Local = parsePath(step.Local)
		conf.SyncSteps[idx].Set.Remote = parsePath(step.Remote)
		conf.SyncSteps[idx].Set.Command, conf.SyncSteps[idx].Set.Errors = assembleCommand(step)
	}
	return
}

func parsePath(pth string) (p tPath) {
	p.FullPath = pth
	p.Path = pth
	p.IsFolder = nil
	p.IsLocal = isLocalPath(p.FullPath)
	p.IsHealthy, p.Errors = isHealthy(p)
	if p.IsLocal {
		p.IsFolder = isFolder(p.FullPath)
		p.IsEmpty, _ = isEmpty(p.FullPath)
	} else {
		p.IsLocal = false
		arr := strings.Split(p.FullPath, ":")
		p.Machine = arr[0]
		p.Path = arr[1]
	}
	return
}

func assembleCommand(step tSyncStep) (cmdArr []string, errArr []error) {
	var source, target string
	cmdArr = step.Cmd
	if cmdArr[0] == "rsync" && cli.RsyncDryRun {
		cmdArr = append(cmdArr, "-n")
	}
	if cli.Action == "pull" {
		source = step.Set.Remote.FullPath
		target = step.Set.Local.FullPath
	}
	if cli.Action == "push" {
		source = step.Set.Local.FullPath
		target = step.Set.Remote.FullPath
	}
	cmdArr = append(cmdArr, source)
	cmdArr = append(cmdArr, target)
	errArr = step.Set.Local.Errors
	errArr = append(errArr, step.Set.Remote.Errors...)
	return
}

func templateFile(conf string) (by []byte, err error) {
	ud := getUserdataMap()
	buf := &bytes.Buffer{}
	templ, err := template.New("conf").Parse(conf)
	if err == nil {
		templ.Execute(buf, map[string]interface{}{
			"confdir": filepath.Dir(pabs(cli.Config)),
			"workdir": pwd(),
			"home":    ud["home"],
			"uid":     ud["uid"],
			"gid":     ud["gid"],
			"user":    ud["username"],
		})
		by = buf.Bytes()
	}
	return
}

func getUserdataMap() map[string]string {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	m := make(map[string]string)
	m["home"] = user.HomeDir + "/"
	m["uid"] = user.Uid
	m["gid"] = user.Gid
	m["username"] = user.Username
	m["name"] = user.Name
	return m
}
