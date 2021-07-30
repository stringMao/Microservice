package agent

import (
	"fmt"
	"net"
)

type agentServer struct {
	Serverid uint64
	open     bool
	conn     net.Conn    //客户端连接
	send     chan []byte //需要发送给该服务器的消息
}

func NewAgentServer(serverid uint64, c net.Conn) *agentServer {
	s := &agentServer{
		Serverid: serverid,
		open:     true,
		conn:     c,
		send:     make(chan []byte, 1000),
	}
	go func(*agentServer) {
		for {
			msg := <-s.send
			if msg == nil {
				break
			}
			s.conn.Write(msg)
		}
		fmt.Printf("服务器发送携程关闭:TSid[%d]", s.Serverid)
	}(s)

	return s
}

func (s *agentServer) SendData(msg []byte) {
	s.send <- msg
}

func (s *agentServer) Close() {
	s.send <- nil
	s.conn.Close()
	s.open = false
}
