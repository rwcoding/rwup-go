package internal

import (
	"github.com/BurntSushi/toml"
	"io/fs"
	"log"
	"path"
	"path/filepath"
	"strings"
)

var mapDirNames = map[string]string{}

func initDirs() {
	file := path.Join(root(), "_dirs.toml")
	dataMap := map[string]string{}
	if fileExists(file) {
		if _, err := toml.DecodeFile(file, &dataMap); err != nil {
			log.Fatal(err)
		}
	}
	for k, v := range dataMap {
		if k[:1] == "$" {
			mapDirNames[k] = v
		} else {
			mapDirNames["$"+k] = v
		}
	}
}

func getDirSignByDir(path string) string {
	return "$" + strings.ReplaceAll(strings.ReplaceAll(path, root(), ""), "\\", "/")
}

func getDirNameByDir(path string) string {
	sign := getDirSignByDir(path)
	if name, ok := mapDirNames[sign]; ok {
		return name
	}
	return strings.ReplaceAll(filepath.Base(sign), filepath.Ext(sign), "")
}

func recurseDirName() map[string]string {
	dataMap := make(map[string]string)
	_ = filepath.WalkDir(root(), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
			return err
		}
		if d.IsDir() {
			first := d.Name()[:1]
			if first == "." || first == "_" {
				return err
			}
			sign := getDirSignByDir(path)
			if sign == "$" {
				return err
			}
			dataMap[sign] = getDirNameByDir(path)
		}
		return err
	})
	return dataMap
}

func getDirSignByFile(file string) string {
	return getDirSignByDir(filepath.Dir(file))
}
