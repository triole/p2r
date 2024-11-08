package main

import (
	"bytes"
	"os"
	"os/user"
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
		conf.SyncSteps[idx].Local.Folder = cleanPath(step.Local.Folder)
		conf.SyncSteps[idx].Remote.Folder = cleanPath(step.Remote.Folder)
	}
	return
}

func templateFile(conf string) (by []byte, err error) {
	ud := getUserdataMap()
	buf := &bytes.Buffer{}
	templ, err := template.New("conf").Parse(conf)
	if err == nil {
		templ.Execute(buf, map[string]interface{}{
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
