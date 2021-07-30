package agent

import (
	"fmt"
	"net"
)

type agentClient struct {
	Userid uint64
	conn   net.Conn //客户端连接
	send   chan []byte
	//svrlist map[uint16]uint64 //该用户连接的服务器列表map[tid]serverid（每种服务器只能连一个 serverid=0+0+sid+tid）

	svrMap map[uint16]*agentServer
}

func NewAgentClient(userid uint64, c net.Conn) *agentClient {
	a := &agentClient{
		Userid: userid,
		conn:   c,
		send:   make(chan []byte, 100),
		//svrlist: make(map[uint16]uint64, 1),
		svrMap: make(map[uint16]*agentServer, 5),
	}

	go func(a *agentClient) {
		for {
			msg := <-a.send
			if msg == nil { //nil表示关闭
				break
			}
			a.conn.Write(msg)
		}
		fmt.Printf("客户端发送携程关闭:userid[%d]", a.Userid)
	}(a)

	return a
}

func (a *agentClient) Close() {
	a.send <- nil
	a.conn.Close()
}

//给该链接发送消息
func (a *agentClient) SendData(msg []byte) {
	a.send <- msg
}

//上传给服务器消息
func (a *agentClient) PushData(tid uint16, msg []byte) bool {
	if s, ok := a.svrMap[tid]; ok {
		if s != nil {
			if s.open {
				s.SendData(msg)
				return true
			} else {
				delete(a.svrMap, tid)
			}
		}
	}
	return false
}

func (a *agentClient) SetSvr(tid uint16, psvr *agentServer) {
	a.svrMap[tid] = psvr
}

//GetSvr
// func (a *agentClient) GetSvr(tid uint16) uint64 {
// 	if s, ok := a.svrlist[tid]; ok {
// 		return s
// 	}
// 	//如果原来没有分配具体的sid。则负载均衡一个
// 	serverid := balanceManager.AllocServerid(tid)
// 	if serverid != 0 {
// 		a.SetSvr(tid, serverid)
// 	}
// 	return serverid
// }
