package internal

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type responseCheckItem struct {
	Files []string `json:"files"`
}

type responseCheck struct {
	Code int               `json:"code"`
	Msg  string            `json:"msg"`
	Data responseCheckItem `json:"data"`
}

func apiCheck(dirs map[string]string, files map[string]ReqFile) *responseCheck {
	params := map[string]interface{}{}
	params["dirs"] = dirs
	params["files"] = files
	result := request("/api/open/sync/check", params)
	res := new(responseCheck)
	if err := json.Unmarshal(result, res); err != nil {
		log.Fatal()
	}
	return res
}

type responseSyncItem struct {
	Num int `json:"num"`
}
type responseSync struct {
	Code int              `json:"code"`
	Msg  string           `json:"msg"`
	Data responseSyncItem `json:"data"`
}

func apiSync(docs []interface{}) *responseSync {

	params := map[string]interface{}{}
	params["docs"] = docs

	result := request("/api/open/sync", params)
	res := new(responseSync)
	if err := json.Unmarshal(result, res); err != nil {
		log.Fatal()
	}
	return res
}

func request(route string, params map[string]interface{}) []byte {
	params["username"] = config.Username
	params["password"] = config.Password
	params["project"] = config.Project

	jsonByte, err := json.Marshal(params)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	url := config.Url + route
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonByte))
	if err != nil {
		log.Fatal(err)
		return nil
	}

	ut := strconv.FormatInt(time.Now().Unix(), 10)
	token := "-"
	key := "-"
	str := route + ut + token + string(jsonByte) + key
	m := md5.New()
	m.Write([]byte(str))
	sign := hex.EncodeToString(m.Sum(nil))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Time", strconv.FormatInt(time.Now().Unix(), 10))
	req.Header.Set("X-Token", token)
	req.Header.Set("X-Sign", sign)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	return body
}
