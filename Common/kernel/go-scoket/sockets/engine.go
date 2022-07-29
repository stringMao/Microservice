package sockets

import (
	//"Common/kernel"
	"Common/log"
	"Common/try"
	"errors"
	"time"
)
func init(){
}

type IHandler interface {
	BeforeHandle(*Engine,int,[]byte)
	CloseHandle(*Engine)
}

const(
	Type_Socket =iota+1
	Type_WebSocket
)


type Handler struct {
	handleFuncMap  map[uint32]func(s *Engine,buf []byte)
}
func (s *Handler)GetHandleFunc(id uint32)func(s *Engine,buf []byte){
	if s.handleFuncMap==nil{
		return nil
	}
	fn, has := s.handleFuncMap[id]
	if has {
		return fn
	}
	return nil
}
func (s *Handler)AddHandleFuc(id uint32,fn func(s *Engine,buf []byte)){
	if s.handleFuncMap==nil{
		s.handleFuncMap= make(map[uint32]func(s *Engine,buf []byte),10)
	}
	s.handleFuncMap[id]=fn
}


const ServerConnect = 1 //engine对像是在服务端创建的
const ClientConnect = 2 //engine对象创建是为了作为客户端去连接服务器的


type engineOption func(c *Engine)
type Engine struct {
	role       int  //角色
	Session  		*Session
	Data     		interface{}
	handleClient 	IHandler
	Heart     	time.Duration //单位10的-6次
}

func WithData(data interface{}) engineOption {
	return func(c *Engine) {
		c.Data=data
	}
}
func WithHeart(t time.Duration) engineOption {
	return func(c *Engine) {
		c.Heart=t
	}
}

func NewEngine(role int,session *Session,options ...engineOption)*Engine {
	s:=&Engine{
		role: role,
		//status: STATUS_INIT,
		Session: session,
		Data:nil,
		Heart:0,
	}
	for _,op:=range options {
		op(s)
	}
	return s
}


func (c *Engine)Close(){
	c.Session.Close()
}

func (c *Engine)Start(){
	defer try.Catch()
	defer func() {
		c.Session.Close()
		c.handleClient.CloseHandle(c)
	}()

	//c.status=STATUS_CONNECT
	go c.WriteMessage()
	for{
		//开始通讯
		if c.role==ServerConnect && c.Heart > 0{
			c.Session.Conn.SetReadDeadline(time.Now().Add(c.Heart)) //借此检测心跳包
		}
		//c.Session.ReadBuf.Data=c.Session.ReadBuf.Data[0:0]
		n,err := c.Session.Conn.Read(c.Session.ReadBuf.Data)            //读取客户端传来的内容
		if err != nil {
			log.Debug(c.Session.Addr, " connect is close: ", err)
			return //当远程客户端连接发生错误（断开）后，终止此协程。
		}
		//特殊包-心跳包过滤  消息结构[uint8]=200
		if n == 1 && c.Session.ReadBuf.Data[0] == 200 && c.role==ServerConnect{
			//log.Logger.Debugln("收到 heart")
			continue
		}
		// 协议包格式 协议包长度（4字节,表示后面字节流长度，不包含本身的4字节）+[]byte
		if c.handleClient!=nil{
			//协议解包和粘包
			if c.Session.TempBuf!=nil{
				c.Session.ReadBuf.Data=append(c.Session.TempBuf,c.Session.ReadBuf.Data[:n]...)
				n=len(c.Session.ReadBuf.Data)
				c.Session.TempBuf = nil
			}
			c.unpack(n)
		}
	}
}


func (c *Engine)unpack(length int)error {
	if length < message_head_size{
		return errors.New("协议太小")
	}
	index:=0
	for index<length{
		c.Session.ReadBuf.ResetReadIndex(index)
		msgLen:=int(c.Session.ReadBuf.ReadUint32BE())

		if msgLen+message_head_size > length-index{
			//包长度不够，说明包不完整，需要等待下一条消息粘包
			c.Session.TempBuf=make([]byte,length-index,2048)
			copy(c.Session.TempBuf,c.Session.ReadBuf.Data[index:length])
			break
		}else{
			index=index+message_head_size
			c.handleClient.BeforeHandle(c,msgLen,c.Session.ReadBuf.Data[index:index+msgLen])
			index=index+msgLen
		}
	}
	return nil
}

func (c *Engine)WriteMessage(){
	defer try.Catch()
	for {
		select {
		case msg, ok := <-c.Session.writeChan:
			if !ok{
				goto stop
			}
			//协议组合
			msgLen:=len(msg)
			c.Session.OutBuf.InitBuf(msgLen+message_head_size)
			c.Session.OutBuf.WriteUint32BE(uint32(msgLen))
			c.Session.OutBuf.Write(msg)

			//发送
			c.Session.Conn.Write(c.Session.OutBuf.Data)
		case <-c.heartTicker():
			 _, err := c.Session.Conn.Write(c.Session.HeartBuf)
			 if err != nil {
				 log.Errorln("Connector Heart err:", err)
			 }
		 //case <-c.Session.closeChan:
			//	 goto stop
			 }
		 }

stop:
	log.Debugln("WriteMessage goroutine is stop")
}

func (c *Engine)SetHandler(h IHandler){
	c.handleClient=h
}

func (c *Engine)heartTicker()<-chan time.Time{
	if c.role==ServerConnect{
		return  nil
	}
	return time.Tick(c.Heart)
}

func (c *Engine)SendData(data []byte)bool{
	if data==nil{
		return false
	}
	return c.Session.Send(data)
}
func (c *Engine)SyncSendData(data []byte)bool{
	if data==nil{
		return false
	}
	return c.Session.SyncSend(data)
}