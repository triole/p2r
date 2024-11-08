package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
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

func fileExists(path string) bool {
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
	if err != nil {
		log.Fatal(err)
	}
	b = fi.IsDir()
	return
}

func isFolderEmpty(name string) (r bool) {
	r = true
	arr, _ := filepath.Glob(name + "*")
	if len(arr) > 0 {
		r = false
	}
	return
}

func pwd() string {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return pwd
}
