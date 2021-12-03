package main

import (
	"Common/msg"
	"Common/proto/base"
	"TestClient/gatesvr"
	"TestClient/loginsvr"
	"fmt"
	"os"
	"time"

	"github.com/golang/protobuf/proto"
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
				msgData := &base.TestMsg{
					Txt: "这是一条测试消息",
				}
				dPro, _ := proto.Marshal(msgData)
				testmsg := msg.CreateWholeMsgData(msg.Sign_serverid, 3, msg.MID_Hall, msg.Hall_TestMsg, dPro)
				gate.TestSend(testmsg)
				//gate.TestSend(msg.CreateWholeMsgData(msg.Sign_serverid, 0, msg.MID_Gate, 1, []byte("twes")))
			}
		}

	}()

	c := make(chan os.Signal)
	<-c
}
