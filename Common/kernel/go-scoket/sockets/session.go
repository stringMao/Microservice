package sockets

import (
	"Common/kernel"
	"Common/try"
	"net"
)

type SessionOption func(s *Session)

type Session struct {
	sessionID   uint64
	Addr        string
	NetWork     string  //"socket","websocket"
	Conn   		net.Conn
	//WebConn
	writeChan  	chan []byte
	//ReadBuf    	[]byte
	ReadBuf     *kernel.ByteBuf
	OutBuf      *kernel.ByteBuf
	TempBuf     []byte
	HeartBuf    []byte
}
func WithSessionId(id uint64)SessionOption{
	return func(s *Session) {
		s.sessionID=id
	}
}
func WithNetWork(str string)SessionOption{
	return func(s *Session) {
		s.NetWork=str
	}
}
func NewSession(conn net.Conn,writeCap int,options ...SessionOption)*Session{
	s:= &Session{
		sessionID: 0,
		NetWork:"socket",
		Conn: conn,
		//WebConn:nil,
		writeChan:make(chan []byte,writeCap),
		//closeChan: make(chan bool,1),
		//ReadBuf: make([]byte,messageMaxLen),
		ReadBuf:kernel.NewByteBuf(GetBytePool()),
		OutBuf:kernel.NewByteBuf(GetBytePool()),
		TempBuf:nil,
		HeartBuf: make([]byte,1,1),
	}
	s.HeartBuf[0]=200
	for _,op:=range options{
		op(s)
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
	var err error
	if s.Conn!=nil{
		err = s.Conn.Close() //重复close 会err
	}

 	if err==nil{
		s.ReadBuf.Free()
		s.OutBuf.Free()
		//s.closeChan<-true
		close(s.writeChan)
	}

}


