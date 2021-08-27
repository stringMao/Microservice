package agent

import (
	"Common/log"
	"Common/msg"
	"net"
	"runtime/debug"
)

type agentServer struct {
	Tid      uint32 //服务类型id
	Sid      uint32 //同服务类型下的唯一标识id
	Serverid uint64 //0+0+sid+tid 位运算获得

	conn    net.Conn    //
	send    chan []byte //需要发送给该服务器的消息
	sendEnd chan bool   //发送消息的协程开关
}

func NewAgentServer(tid, sid uint32, c net.Conn) *agentServer {
	s := &agentServer{
		Tid:      tid,
		Sid:      sid,
		Serverid: msg.EncodeServerID(tid, sid),

		conn:    c,
		send:    make(chan []byte, 1000),
		sendEnd: make(chan bool),
	}
	go func(s *agentServer) {
		defer func() {
			s.sendEnd <- true
			log.Debugf("服务器发送携程关闭:serverid[%d]", s.Serverid)
		}()
		defer func() {
			if r := recover(); r != nil {
				log.PrintPanicStack(r, string(debug.Stack()))
				s.conn.Close()
			}
		}()

		for {
			if msg, ok := <-s.send; ok {
				s.conn.Write(msg)
			} else {
				break
			}
		}
	}(s)

	return s
}

//发送数据
func (s *agentServer) SendData(msg []byte) (suc bool) {
	defer func() {
		if recover() != nil {
			suc = false //发送失败
		}
	}()
	s.send <- msg //s.send chan在被close之后，插入数据会异常
	return true
}

func (s *agentServer) Close() {
	//关闭子携程 close要保证不可以重复
	close(s.send) //close之后，缓存区还有数据，还是会返回ok=true，直到缓冲区清空

	//确保发送协程结束
	<-s.sendEnd

	s.conn.Close()
}
