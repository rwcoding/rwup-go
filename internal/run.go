package internal

import (
	"fmt"
	"log"
	"path"
	"strconv"
)

func Run() {
	dirs := recurseDirName()
	files := recurseFileMd5()

	result := apiCheck(dirs, files)
	if result.Code != 10000 {
		log.Fatal(result.Msg)
	}
	num := 0
	var docs []interface{}
	for _, v := range result.Data.Files {
		fmt.Println(" -- sync " + v + " ......")
		if v[:1] == "$" {
			v = v[1:]
		}
		parsed, err := parse(path.Join(root(), v))
		if err != nil {
			fmt.Println(err)
			continue
		}
		docs = append(docs, parsed)
		if len(docs) >= 5 {
			num += syncDocs(docs)
			docs = docs[0:0]
		}
		fmt.Println(v)
	}
	if len(docs) > 0 {
		num += syncDocs(docs)
	}

	fmt.Println("[over] sync document " + strconv.Itoa(num) + ".")
}

func syncDocs(docs []interface{}) int {
	result := apiSync(docs)
	if result.Code != 10000 {
		log.Println("api err: " + result.Msg)
		return 0
	} else {
		return result.Data.Num
	}
}
