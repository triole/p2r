package conf

import (
	"bytes"
	"os/user"
	"path/filepath"
	"strings"
	"text/template"
)

func isLocalPath(path string) bool {
	return !strings.Contains(path, ":")
}

func (conf Conf) assembleCommand(step SyncStep) (cmdArr []string, errArr []error) {
	var source, target string
	cmdArr = step.Cmd
	if cmdArr[0] == "rsync" && conf.RsyncDryRun {
		cmdArr = append(cmdArr, "-n")
	}
	if conf.Action == "pull" {
		source = step.Set.Remote.FullPath
		target = step.Set.Local.FullPath
	}
	if conf.Action == "push" {
		source = step.Set.Local.FullPath
		target = step.Set.Remote.FullPath
	}
	cmdArr = append(cmdArr, source)
	cmdArr = append(cmdArr, target)
	errArr = step.Set.Local.Errors
	errArr = append(errArr, step.Set.Remote.Errors...)
	return
}

func (conf Conf) templateFile(str string) (by []byte, err error) {
	ud := getUserdataMap()
	buf := &bytes.Buffer{}
	templ, err := template.New("conf").Parse(str)
	if err == nil {
		templ.Execute(buf, map[string]interface{}{
			"confdir": filepath.Dir(conf.pabs(conf.ConfigFile)),
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
