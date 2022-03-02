package internal

import (
	"github.com/BurntSushi/toml"
	"log"
	"path"
)

var mapFileTitles = map[string]string{}

func initFileTitles() {
	file := path.Join(root(), "_titles.toml")
	dataMap := map[string]string{}
	if fileExists(file) {
		if _, err := toml.DecodeFile(file, &dataMap); err != nil {
			log.Fatal(err)
		}
	}
	for k, v := range dataMap {
		if k[:1] == "$" {
			mapFileTitles[k] = v
		} else {
			mapFileTitles["$"+k] = v
		}
	}
}
