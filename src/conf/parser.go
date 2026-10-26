package conf

import (
	"bytes"
	"os"
	"os/user"
	"path/filepath"
	"text/template"

	"p2r/src/logger"

	yaml "gopkg.in/yaml.v3"
)

// Parser handles configuration file parsing and template expansion.
type Parser struct {
	ConfigFile string
	Lg         *logger.Logger
}

// NewParser creates a new Parser instance.
func NewParser(configFile string, lg *logger.Logger) *Parser {
	return &Parser{
		ConfigFile: pabs(configFile),
		Lg:         lg,
	}
}

// Parse reads and parses the YAML configuration file with template expansion.
func (p *Parser) Parse() (ConfigContent, error) {
	// Read config file
	content, err := os.ReadFile(p.ConfigFile)
	if err != nil {
		return ConfigContent{}, err
	}

	// Expand templates
	expanded, err := p.expandTemplate(string(content))
	if err != nil {
		return ConfigContent{}, err
	}

	// Unmarshal YAML
	var config ConfigContent
	err = yaml.Unmarshal(expanded, &config)
	if err != nil {
		return ConfigContent{}, err
	}

	p.Lg.Debug("config parsed", map[string]interface{}{"content": config})
	return config, nil
}

// expandTemplate processes template variables in the config content.
func (p *Parser) expandTemplate(content string) ([]byte, error) {
	userData := p.getUserDataMap()

	tmpl, err := template.New("conf").Parse(content)
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, map[string]interface{}{
		"confdir": filepath.Dir(p.ConfigFile),
		"CONFDIR": filepath.Dir(p.ConfigFile),
		"workdir": currentDir(),
		"WORKDIR": currentDir(),
		"home":    userData["home"],
		"HOME":    userData["home"],
		"uid":     userData["uid"],
		"UID":     userData["uid"],
		"gid":     userData["gid"],
		"GID":     userData["gid"],
		"user":    userData["username"],
		"USER":    userData["username"],
	})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// getUserDataMap retrieves current user information.
func (p *Parser) getUserDataMap() map[string]string {
	currentUser, err := user.Current()
	if err != nil {
		// Return defaults on error instead of panic
		p.Lg.Error("failed to get current user", map[string]interface{}{"error": err})
		return map[string]string{
			"home":     "",
			"uid":      "",
			"gid":      "",
			"username": "",
			"name":     "",
		}
	}

	return map[string]string{
		"home":     currentUser.HomeDir,
		"uid":      currentUser.Uid,
		"gid":      currentUser.Gid,
		"username": currentUser.Username,
		"name":     currentUser.Name,
	}
}

// currentDir returns the current working directory.
func currentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return dir
}

// pabs returns the absolute path for a given path string.
func pabs(pathstring string) string {
	abs, err := filepath.Abs(pathstring)
	if err != nil {
		return pathstring // Return original path on error
	}
	return abs
}
