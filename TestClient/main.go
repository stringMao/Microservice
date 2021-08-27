package main

import (
	"Common/msg"
	"TestClient/gatesvr"
	"TestClient/loginsvr"
	"fmt"
	"os"
	"time"
)

func main() {

	// var test map[int]string
	// test = make(map[int]string)
	// test[1] = "mao"
	// test[2] = "ling"
	// test[2] = "jia"
	// test[1] = "jia"
	// delete(test, 2)
	// delete(test, 2)

	fmt.Println("111")
	userid, token, ip := loginsvr.Signin()

	gate := gatesvr.NewConnTcp(ip)

	go gate.Connect(userid, token)

	go func() {
		for {
			time.Sleep(time.Second * 1)
			if gatesvr.ConnSucc {
				gate.TestSend(msg.CreateWholeMsgData(msg.Sign_serverid, 0, msg.MID_Gate, 1, []byte("twes")))
			}
		}

	}()

	c := make(chan os.Signal)
	<-c
}
