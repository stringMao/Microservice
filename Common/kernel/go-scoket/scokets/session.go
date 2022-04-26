package scokets

import (
	"Common/kernel/module/bytepool"
	"Common/try"
	"net"
)

type Session struct {
	sessionID   uint64
	NetWork    string  //"socket","websocket"
	Conn   		net.Conn
	//WebConn
	writeChan  	chan []byte
	closeChan  	chan bool
	ReadBuf    	[]byte
	BytePool    *bytepool.BytePoolCap

	Addr        string
}
func NewSession(fd uint64,conn net.Conn,writeCap int)*Session{
	s:= &Session{
		sessionID: fd,
		NetWork:"socket",
		Conn: conn,
		//WebConn:nil,
		writeChan:make(chan []byte,writeCap),
		closeChan: make(chan bool,1),
		ReadBuf: make([]byte,messageMaxLen),
		BytePool:bytePool,
	}
	if conn!=nil{
		s.Addr= conn.RemoteAddr().String()
	}
	return s
}

// NewWebSocketSession TODO:websocket待实现
func NewWebSocketSession(fd uint64)*Session{
	return &Session{
		sessionID: fd,
		NetWork:"websocket",
	}
}


func (s *Session)SetID(id uint64){
	s.sessionID=id
}
func (s *Session)SetBytePool(bp *bytepool.BytePoolCap){
	s.BytePool=bp
}

func (s *Session)Send(msg []byte)(suc bool) {
	defer func() {
		if recover() != nil {
			suc = false //发送失败
		}
	}()
	s.writeChan <- msg //a.send chan在被close之后，插入数据会异常
	return true
}

func (s *Session)SyncSend(msg []byte)bool{
	_,err:=s.Conn.Write(msg)
	return err==nil
}
func (s *Session)Close(){
	defer try.Catch()
	s.closeChan<-true
	close(s.writeChan)

	if s.Conn!=nil{
		s.Conn.Close()
	}
}


