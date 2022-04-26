package scokets

import (
	"Common/kernel/module/bytepool"
	"Common/log"
	"Common/try"
	"time"
)

type IHandler interface {
	BeforeHandle(*Client,int,[]byte)
	CloseHandle(*Client)
}

const(
	Type_Scoket =iota+1
	Type_WebScoket
)
const messageMaxLen = 1024

type Handler struct {
	handleFuncMap  map[uint32]func(s *Client,len int,buf []byte)
}
func (s *Handler)GetHandleFunc(id uint32)func(s *Client,len int,buf []byte){
	if s.handleFuncMap==nil{
		return nil
	}
	fn, has := s.handleFuncMap[id]
	if has {
		return fn
	}
	return nil
}
func (s *Handler)AddHandleFuc(id uint32,fn func(s *Client,len int,buf []byte)){
	if s.handleFuncMap==nil{
		s.handleFuncMap= make(map[uint32]func(s *Client,len int,buf []byte),10)
	}
	s.handleFuncMap[id]=fn
}

//==字节池==
var bytePool  *bytepool.BytePoolCap=nil
func InitBytePool(maxsize,len,cap int){
	bytePool=bytepool.NewBytePoolCap(maxsize,len,cap)
}
func GetByteFormPool()[]byte{
	if bytePool==nil{
		InitBytePool(1000,messageMaxLen,messageMaxLen)
	}

	return bytePool.Get()
}


//===
const ServerConnect = 1
const ClientConnect = 2

const (
	STATUS_INIT = iota
	STATUS_CONNECT
	STATUS_CLOSE
)

type Client struct {
	role       int  //角色
	status     int
	Session  		*Session
	Data     		interface{}
	handleClient 	IHandler
	Heart     	time.Duration //单位10的-6次
}
func NewClient(role int,session *Session,data interface{},hearttime time.Duration)*Client{
	s:=&Client{
		role: role,
		status: STATUS_INIT,
		Session: session,
		Data:data,
		Heart:hearttime,
	}
	return s
}


func (c *Client)Close(){
	if c.status==STATUS_CONNECT{
		c.status=STATUS_CLOSE
		c.Session.Close()
	}
}

func (c *Client)Start(){
	defer try.Catch()
	defer func() {
		if c.status==STATUS_CONNECT{//异常退出
			c.Session.Close()//确保写携程能退出
		}
		c.handleClient.CloseHandle(c)
	}()

	c.status=STATUS_CONNECT
	go c.WriteMessage()
	for{
		//开始通讯
		if c.role==ServerConnect && c.Heart > 0{
			c.Session.Conn.SetReadDeadline(time.Now().Add(c.Heart)) //借此检测心跳包
		}
		n,err := c.Session.Conn.Read(c.Session.ReadBuf)            //读取客户端传来的内容
		if err != nil {
			log.Debug(c.Session.Addr, " connet is close: ", err)
			return //当远程客户端连接发生错误（断开）后，终止此协程。
		}
		//特殊包-心跳包过滤  消息结构[uint8]=200
		if n == 1 && c.Session.ReadBuf[0] == 200 && c.role==ServerConnect{
			//log.Logger.Debugln("heart")
			continue
		}
		if c.handleClient!=nil{
			c.handleClient.BeforeHandle(c,n,c.Session.ReadBuf)
		}
	}
}



func (c *Client)WriteMessage(){
	defer try.Catch()
	defer func() {
		c.Session.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.Session.writeChan:
			if !ok{
				goto stop
			}
			c.Session.Conn.Write(msg)
			if c.Session.BytePool!=nil{
				c.Session.BytePool.Put(msg)
			}
		case <-c.heartTicker():
				 _, err := c.Session.Conn.Write([]byte{200})
				 if err != nil {
					 log.Errorln("Connector Heart err:", err)
				 }
		 case <-c.Session.closeChan:
				 goto stop
			 }
		 }

stop:
	log.Debugln("client WriteMessage is stop")
}

func (c *Client)SetHandler(h IHandler){
	c.handleClient=h
}

func (c *Client)heartTicker()<-chan time.Time{
	if c.role==ServerConnect{
		return  nil
	}
	return time.Tick(c.Heart)
}

func (c *Client)SendData(data []byte)bool{
	return c.Session.Send(data)
}
