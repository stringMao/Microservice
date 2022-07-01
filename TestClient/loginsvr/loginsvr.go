package loginsvr

import (
	"Common/log"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type RespDataResult struct {
	Status  int           `json:"status"`
	Message string        `json:"message"`
	Data    replyAccLogin `json:"data"`
}

type ServerIPInfo struct {
	Address string //ip:port
	Svrname string
}
type replyAccLogin struct {
	Userid  uint64         `json:"userid"`
	Token   string         `json:"token"`
	Svrlist []ServerIPInfo `json:"svrlist"`
}

func Signin() (userid uint64, token string, ip string) {

	m := make(map[string]interface{})
	m["username"] = url.QueryEscape("test001")
	m["passwd"] = url.QueryEscape("123")
	m["channel"] = 1
	bytesData, err := json.Marshal(m)
	if err != nil {
		log.Logger.Println(err)
		return
	}
	reader := bytes.NewReader(bytesData)

	request, err := http.NewRequest("POST", "http://192.168.32.36:8202/mloginsvr/client/signin", reader)
	if err != nil {
		log.Logger.Println(err)
		return
	}
	request.Header.Set("sign", "test")

	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Logger.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Logger.Println(err)
		return
	}

	var ret RespDataResult
	err = json.Unmarshal(body, &ret)
	if err != nil {
		log.Logger.Println(err)
		return
	}

	if ret.Status == 0 {
		return uint64(ret.Data.Userid), ret.Data.Token, ret.Data.Svrlist[0].Address
	}

	return 0, "", ""
}
