package gatesvr

import (
	"Common/proto/base"
	"fmt"
	"net"
	"time"

	"github.com/golang/protobuf/proto"
)

type ConnTcp struct {
	conn  net.Conn
	open  bool
	SvrIP string
	send  chan []byte
}

func NewConnTcp(ip string) *ConnTcp {
	return &ConnTcp{
		open:  false,
		SvrIP: ip,
		send:  make(chan []byte, 100),
	}
}

func (c *ConnTcp) Connect(userid uint64, token string) bool {
	var err error
	c.conn, err = net.Dial("tcp", "127.0.0.1:8090") //本机必须用127.0.0.1。直接用ip没法连接
	if err != nil {
		fmt.Println("connect", c.SvrIP, "fail", err.Error())
		return false
	}
	c.open = true

	stSend := &base.ClientLogin{
		Userid: uint64(userid),
		Token:  token,
	}
	pData, err := proto.Marshal(stSend)
	if err != nil {
		panic(err)
	}
	//发送
	c.conn.Write(pData)

	go func(c *ConnTcp) {
		for {
			msg := <-c.send
			c.conn.Write(msg)
		}
	}(c)

	c.OpenHeart()

	buffer := make([]byte, 2048)
	for {
		n, err := c.conn.Read(buffer) //读取客户端传来的内容
		if err != nil {
			fmt.Println("read err:", err, n)
			return false //当远程客户端连接发生错误（断开）后，终止此协程。
		}
		//send <- buffer[:n]
		//protoEncodePrint(buffer, n)
	}
	return true
}

func (c *ConnTcp) OpenHeart() {
	//心跳包
	go func() {
		for {
			if !c.open {
				break
			}
			time.Sleep(time.Second * 1)
			c.send <- []byte{200}
		}
	}()
}

// func protoEncodePrint(buf []byte, n int) {
// 	head := &msg.HeadProto{}
// 	head.Decode(buf)
// 	msgstr := &base.TestMsg{}
// 	if head.Len > 0 {
// 		err := proto.Unmarshal(buf[head.GetHeadLen():n], msgstr)
// 		if err != nil {
// 			fmt.Println("协议解析失败2:", err)
// 			return //当远程客户端连接发生错误（断开）后，终止此协程。
// 		}
// 	}
// 	fmt.Printf("接收消息: mainid[%d] sonid[%d] len[%d] msg:%s \n", head.MainID, head.SonID, head.Len, msgstr.Txt)
// 	return

// }
