package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

var username string
var passwd string
var separated = "++++++++++++++++++++++++++++++"

func init() {
	flag.StringVar(&username, "username", "0", "username")
	flag.StringVar(&passwd, "passwd", "", "password")
	flag.Parse()
}

func main() {
	cookies := getCookies(username, passwd)
	resp := fuckUpcCOVID(cookies)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("%v\nresp: %+v\n", separated, string(body))
	_ = resp.Body.Close()
}

func fuckUpcCOVID(cookies []*http.Cookie) *http.Response {
	resp := request("https://app.upc.edu.cn/ncov/wap/default/index", "GET", nil, cookies)
	defer func() {
		_ = resp.Body.Close()
	}()
	body, _ := ioutil.ReadAll(resp.Body)
	m := make(map[string]interface{}, 0)
	err := json.Unmarshal(body, &m)
	if err != nil {
		panic(err)
	}
	info := m["d"].(map[string]interface{})["info"].(map[string]interface{})
	oldInfo := m["d"].(map[string]interface{})["oldInfo"].(map[string]interface{})
	oldInfo["date"] = time.Now().Format("20060102")
	oldInfo["created"] = info["created"]
	oldInfo["id"] = info["id"]
	oldInfo["uid"] = info["uid"]
	fmt.Printf("info: %+v\n%v\noldInfo: %+v\n", info, separated, oldInfo)
	data := url.Values{}
	for k, v := range oldInfo {
		data.Set(k, StringValue(v))
	}
	return request("https://app.upc.edu.cn/ncov/wap/default/save", "POST", data, cookies)
}
