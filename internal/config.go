package internal

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"gopkg.in/ini.v1"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type configDS struct {
	Username string `ini:"username"`
	Password string `ini:"password"`
	Project  string `ini:"project"`
	Root     string `ini:"root"`
	Url      string `ini:"url"`
	Parse    string `ini:"parse"`
}

var config = configDS{}

func init() {
	initConf()
	initFileTitles()
	initDirs()
}

func initConf() {
	var conf string
	flag.StringVar(&conf, "conf", "", "")
	flag.Parse()

	if conf == "" {
		conf, _ = getCurrentDirConf()
	}

	if conf == "" {
		log.Fatal(".docg not found")
	}

	cfg, err := ini.Load(conf)
	if err != nil {
		log.Fatal(err)
	}

	if err := cfg.MapTo(&config); err != nil {
		log.Fatal(err)
	}
}

func root() string {
	dir, _ := os.Getwd()
	root := ""
	if config.Root[:1] == "/" || strings.Contains(config.Root, ":") {
		root = config.Root
	} else {
		root = path.Join(dir, config.Root)
	}
	s, err := filepath.Abs(root)
	if err != nil {
		log.Fatal(err)
	}
	return s
}

func fileExists(file string) bool {
	if _, err := os.Stat(file); err != nil {
		return false
	}
	return true
}

func getCurrentDirConf() (file string, err error) {
	file = ""
	dir, err := os.Getwd()
	if err == nil {
		file = dir + string(os.PathSeparator) + ".docg"
	}
	return
}

func getFileMd5(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	hash := md5.New()
	_, _ = io.Copy(hash, file)
	return hex.EncodeToString(hash.Sum(nil))
}
