package main

import (
	"TestClient/gatesvr"
	"TestClient/loginsvr"
	"fmt"
	"os"
)

func main() {
	fmt.Println("111")
	userid, token, ip := loginsvr.Signin()

	gate := gatesvr.NewConnTcp(ip)

	go gate.Connect(userid, token)

	c := make(chan os.Signal)
	<-c
}
