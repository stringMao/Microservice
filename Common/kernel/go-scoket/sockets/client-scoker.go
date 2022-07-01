package sockets

import (
	"Common/log"
	"Common/try"
	"net"
	"time"
)


type connectorOption func(c *Connector)
type Connector struct {
	Handler
	Engine    *Engine
	NetType   int  //1=scoket  2=webscoket
	reconnect time.Duration  //重连间隔
}
func WithNetType(netType int)connectorOption{
	return func(c *Connector){
		c.NetType=netType
	}
}
func WithReConnect(d time.Duration)connectorOption{
	return func(c *Connector) {
		c.reconnect=d
	}
}


func NewConnector(addr string, handler IHandler, heartTime time.Duration,options ...connectorOption)*Connector{
	session:=NewSession(nil,100)
	session.Addr=addr
	engine:=NewEngine(ClientConnect,session,WithHeart(heartTime))
	engine.SetHandler(handler)
	c:= &Connector{
		NetType:   Type_Socket,
		Engine:engine,
		reconnect: -1,
	}
	for _,op:=range options {
		op(c)
	}
	return c
}


func (c *Connector)StartConnect()bool{
	defer try.Catch()

	switch c.NetType{
	case Type_Socket:
		//var err error
		conn, err := net.Dial("tcp", c.Engine.Session.Addr) //
		if err != nil {
			log.Errorf("connect[%s] is fail,err:%s", c.Engine.Session.Addr, err.Error())
			return false
		}
		c.Engine.Session.Conn=conn

		go c.Engine.Start()

		return true
	case Type_WebSocket:
		log.Fatalln("webscoket 连接未实现")
	default:
		log.Fatalln("未知连接方式")
	}
	return false
}
func (c *Connector)ReConnect(){
	if c.reconnect>=0{

	}
}

func (c *Connector)SendData(msg []byte){
	defer try.Catch()
	c.Engine.Session.Send(msg)
}

func (c *Connector)GetAddrInfo()string{
	if c.Engine !=nil && c.Engine.Session!=nil{
		return c.Engine.Session.Addr
	}
	return ""
}

