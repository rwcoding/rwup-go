package internal

import (
	"io/fs"
	"log"
	"path/filepath"
	"strings"
)

type ReqFile struct {
	Md5   string `json:"md5"`
	File  string `json:"file"`
	Title string `json:"title"`
}

func recurseFileMd5() map[string]ReqFile {
	dataMap := make(map[string]ReqFile)
	_ = filepath.WalkDir(root(), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
			return err
		}
		if !d.IsDir() {
			first := d.Name()[:1]
			if first == "." || first == "_" || !allowFileExt(path) {
				return err
			}
			dataMap[getFileSignByFile(path, "")] = ReqFile{
				Md5:   getFileMd5(path),
				File:  getRelativePathByFile(path),
				Title: getFileTitleByFile(path),
			}
		}
		return err
	})
	return dataMap
}

func allowFileExt(file string) bool {
	ext := strings.ReplaceAll(filepath.Ext(file), ".", "")
	allow := strings.Split(config.Parse, ",")
	if len(allow) == 0 {
		return true
	}
	for _, v := range allow {
		if v == ext {
			return true
		}
	}
	return false
}

func getFileSignByFile(file, route string) string {
	if route != "" {
		return "$" + route
	}
	if strings.Contains(file, ".toml") {
		dsl, err := getDsl(file)
		if err == nil && dsl.Route != "" {
			return "$" + dsl.Route
		}
	}
	return "$" + strings.ReplaceAll(strings.ReplaceAll(file, root(), ""), "\\", "/")
}

func getRelativePathByFile(file string) string {
	return strings.ReplaceAll(strings.ReplaceAll(file, root(), ""), "\\", "/")
}

func getFileTitleByFile(file string) string {
	sign := getFileSignByFile(file, "")
	if name, ok := mapFileTitles[sign]; ok {
		return name
	}
	if strings.Contains(file, ".toml") {
		dsl, err := getDsl(file)
		if err == nil && dsl.Title != "" {
			return dsl.Title
		}
	}
	return strings.ReplaceAll(filepath.Base(sign), filepath.Ext(sign), "")
}
