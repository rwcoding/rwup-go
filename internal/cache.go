package internal

import (
	"github.com/BurntSushi/toml"
)

var dslList = map[string]*Dsl{}

func getDsl(file string) (*Dsl, error) {
	if v, ok := dslList[file]; ok {
		return v, nil
	}
	dsl := new(Dsl)
	_, err := toml.DecodeFile(file, dsl)
	if err != nil {
		// panic(err)
		return nil, err
	}
	dslList[file] = dsl
	return dsl, nil
}
