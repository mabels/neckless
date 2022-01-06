package neckless

import (
	"os"
	"path"
	"strings"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func findFile(plainFname string, dirs ...string) string {
	wd, err := os.Getwd()
	if err == nil {
		var last string
		for ok := true; ok; ok = strings.Compare(last, wd) != 0 {
			tName := path.Join(wd, plainFname)
			if fileExists(tName) {
				return tName
			}
			last = wd
			wd = path.Dir(wd)
		}
	}
	for i := range dirs {
		if fileExists(dirs[i]) {
			return dirs[i]
		}
		tName := path.Join(path.Dir(dirs[i]), plainFname)
		if fileExists(tName) {
			return tName
		}
		tName = path.Join(dirs[i], plainFname)
		if fileExists(tName) {
			return tName
		}
	}
	return plainFname
}
