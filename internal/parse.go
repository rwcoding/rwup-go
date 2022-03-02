package internal

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
)

type Doc struct {
	Title    string     `json:"title"`
	Sign     string     `json:"sign"`
	FileSign string     `json:"file_sign"`
	DirSign  string     `json:"dir_sign"`
	Struct   *DocStruct `json:"struct"`
	Content  string     `json:"content"`
	Order    int        `json:"order"`
}

type DocStruct struct {
	Title        string     `json:"title"`
	Route        string     `json:"route"`
	DocDesc      string     `json:"doc_desc"`
	RequestDesc  string     `json:"request_desc"`
	ResponseDesc string     `json:"response_desc"`
	Request      []Property `json:"request"`
	Response     []Property `json:"response"`
}

type Property struct {
	Name     string     `json:"name"`
	Type     string     `json:"type"`
	Length   int        `json:"length"`
	Required bool       `json:"required"`
	Desc     string     `json:"desc"`
	Sample   string     `json:"sample"`
	Order    int        `json:"order"`
	Tree     []Property `json:"tree"`
}

type Dsl struct {
	Title        string                         `toml:"title"`
	Route        string                         `toml:"route"`
	DocDesc      string                         `toml:"doc_desc"`
	RequestDesc  string                         `toml:"request_desc"`
	ResponseDesc string                         `toml:"response_desc"`
	Request      map[string][]string            `toml:"request"`
	Response     map[string][]string            `toml:"response"`
	Component    map[string]map[string][]string `toml:"component"`
}

func parse(file string) (*Doc, error) {
	var err error
	file, err = filepath.Abs(file)
	if err != nil {
		return nil, err
	}

	doc := &Doc{
		Title:    getFileTitleByFile(file),
		Sign:     getFileSignByFile(file, ""),
		FileSign: getFileMd5(file),
		DirSign:  getDirSignByFile(file),
		Order:    10000,
	}

	if filepath.Ext(file) != ".toml" {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}
		doc.Content = string(content)
		doc.Struct = nil
	} else {
		dsl, err := getDsl(file)
		if err != nil {
			return nil, err
		}
		st := DocStruct{
			Title:        dsl.Title,
			Route:        dsl.Route,
			DocDesc:      dsl.DocDesc,
			RequestDesc:  dsl.RequestDesc,
			ResponseDesc: dsl.ResponseDesc,
		}
		req, err := parseProperty(dsl.Request, dsl)
		if err != nil {
			return nil, err
		}
		res, err := parseProperty(dsl.Response, dsl)
		if err != nil {
			return nil, err
		}
		st.Request = req
		st.Response = res
		doc.Struct = &st
	}
	if doc.Title == "" {
		doc.Title = getFileTitleByFile(file)
	}
	return doc, nil
}

func parseProperty(data map[string][]string, dsl *Dsl) ([]Property, error) {
	properties := []Property{}
	for k, v := range data {
		property := Property{Required: true}
		property.Order = 100
		size := len(v)
		property.Name = k
		if size >= 1 {
			property.Type = v[0]
		} else {
			return nil, errors.New("what type for " + k)
		}

		if property.Type == "array" || property.Type == "object" {
			return nil, errors.New("type error for " + k)
		}

		// 数组类型
		if strings.Contains(property.Type, "array[") &&
			property.Type != "array[int]" &&
			property.Type != "array[float]" &&
			property.Type != "array[double]" &&
			property.Type != "array[short]" &&
			property.Type != "array[long]" &&
			property.Type != "array[char]" &&
			property.Type != "array[string]" &&
			property.Type != "array[bool]" {

			typ := strings.ReplaceAll(property.Type, "array[", "")
			typ = strings.ReplaceAll(typ, "]", "")
			if _, ok := dsl.Component[typ]; !ok {
				return nil, errors.New("type error " + property.Type)
			}
			tree, err := parseProperty(dsl.Component[typ], dsl)
			if err != nil {
				return nil, err
			}
			property.Tree = tree
			property.Type = "array[object]"
		}

		// 自定义对象类型
		if !strings.Contains(property.Type, "array[") &&
			property.Type != "int" &&
			property.Type != "float" &&
			property.Type != "double" &&
			property.Type != "short" &&
			property.Type != "long" &&
			property.Type != "char" &&
			property.Type != "string" &&
			property.Type != "bool" {

			if _, ok := dsl.Component[property.Type]; !ok {
				return nil, errors.New("type error " + property.Type)
			}
			tree, err := parseProperty(dsl.Component[property.Type], dsl)
			if err != nil {
				return nil, err
			}
			property.Tree = tree
			property.Type = "object"
		}

		if size >= 2 {
			property.Desc = v[1]
		}

		min := ""
		max := ""
		if size >= 3 {
			for _, vv := range strings.Split(v[2], "|") {
				if vv == "required" {
					property.Required = true
				} else if vv == "ignore" {
					property.Required = false
				}
				if strings.Contains(vv, "=") {
					vvv := strings.Split(vv, "=")
					if vvv[0] == "order" && len(vvv) == 2 {
						ord, err := strconv.Atoi(vvv[1])
						if err == nil {
							property.Order = ord
						}
					}
					if vvv[0] == "min" && len(vvv) == 2 {
						min = vvv[1]
					}
					if vvv[0] == "max" && len(vvv) == 2 {
						max = vvv[1]
						max, err := strconv.Atoi(vvv[1])
						if err == nil {
							if property.Type == "int" {
								property.Length = len(vvv[1])
							} else {
								property.Length = max
							}
						}
					}
				}
			}
		}
		if size >= 4 {
			property.Sample = v[3]
		}

		if min != "" && max != "" {
			property.Desc = property.Desc + "，范围：" + min + " - " + max
		}
		//if min != "" && max == "" {
		//	property.Desc = property.Desc + "，最小：" + min
		//}
		//if min == "" && max != "" {
		//	property.Desc = property.Desc + "，最大：" + max
		//}

		properties = append(properties, property)
	}
	return properties, nil
}
