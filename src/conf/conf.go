package conf

import (
	"bytes"
	"os/user"
	"path/filepath"
	"text/template"
)

func (conf Conf) templateFile(str string) (by []byte, err error) {
	ud := conf.getUserdataMap()
	buf := &bytes.Buffer{}
	templ, err := template.New("conf").Parse(str)
	if err == nil {
		templ.Execute(buf, map[string]interface{}{
			"confdir": filepath.Dir(conf.pabs(conf.ConfigFile)),
			"CONFDIR": filepath.Dir(conf.pabs(conf.ConfigFile)),
			"workdir": pwd(),
			"WORKDIR": pwd(),
			"home":    ud["home"],
			"HOME":    ud["home"],
			"uid":     ud["uid"],
			"UID":     ud["uid"],
			"gid":     ud["gid"],
			"GID":     ud["gid"],
			"user":    ud["username"],
			"USER":    ud["username"],
		})
		by = buf.Bytes()
	}
	return
}

func (conf Conf) getUserdataMap() map[string]string {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	m := make(map[string]string)
	m["home"] = user.HomeDir
	m["uid"] = user.Uid
	m["gid"] = user.Gid
	m["username"] = user.Username
	m["name"] = user.Name
	return m
}
