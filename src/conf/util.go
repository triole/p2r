package conf

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/triole/logseal"
)

func (conf Conf) isHealthy(pth tPath) (b bool, errArr []error) {
	if conf.rxMatch(pth.FullPath, "^[:/-_]$") {
		errArr = append(
			errArr, errors.New("path seems short: "+pth.FullPath),
		)
	}
	if pth.isEmpty() {
		errArr = append(
			errArr, errors.New("folder empty: "+pth.FullPath),
		)
	}
	return len(errArr) == 0, errArr
}

func (conf Conf) isFolder(fn string) (b bool) {
	fi, err := os.Stat(fn)
	conf.Lg.IfErrFatal("path does not exist", logseal.F{"error": err})
	b = fi.IsDir()
	return
}

func (conf Conf) isEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()
	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

func (conf Conf) rxMatch(rx string, str string) (b bool) {
	re, _ := regexp.Compile(rx)
	b = re.MatchString(str)
	return
}

func (conf Conf) pabs(pathstring string) string {
	r, err := filepath.Abs(pathstring)
	if err != nil {
		fmt.Printf("Unable to make absolute path. %s\n", err)
		os.Exit(1)
	}
	return r
}

func pwd() string {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return pwd
}
