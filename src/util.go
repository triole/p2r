package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/triole/logseal"
)

func cleanPath(s string) (r string) {
	r = path.Clean(s)
	r = pabs(r)
	if !strings.HasSuffix(r, "/") {
		r += "/"
	}
	return
}

func pabs(pathstring string) string {
	r, err := filepath.Abs(pathstring)
	if err != nil {
		fmt.Printf("Unable to make absolute path. %s\n", err)
		os.Exit(1)
	}
	return r
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func isFolder(fn string) (b bool) {
	fi, err := os.Stat(fn)
	lg.IfErrFatal("path does not exist", logseal.F{"error": err})
	b = fi.IsDir()
	return
}

func isEmpty(name string) (bool, error) {
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

func isLocalPath(path string) bool {
	return !strings.Contains(path, ":")
}

func pwd() string {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return pwd
}

func rxMatch(rx string, str string) (b bool) {
	re, _ := regexp.Compile(rx)
	b = re.MatchString(str)
	return
}
