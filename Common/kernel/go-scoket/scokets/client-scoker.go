package scokets

import (
	"Common/log"
	"Common/try"
	"net"
	"time"
)



type Connector struct {
	Handler
	Client   *Client
	NetType   		int  //1=scoket  2=webscoket
	reconnect    	time.Duration  //重连间隔
}

func NewConnector(netType int,addr string,handleAgent IHandler,hearttime time.Duration)*Connector{
	session:=NewSession(0,nil,20)
	session.Addr=addr
	client:=NewClient(ClientConnect,session,nil,hearttime)
	client.SetHandler(handleAgent)
	c:= &Connector{
		NetType: netType,
		Client: client,
		reconnect: 0,
	}
	return c
}
func NewScoketConnector(addr string,handleAgent IHandler,hearttime time.Duration )*Connector {
	session:=NewSession(0,nil,20)
	session.Addr=addr
	client:=NewClient(ClientConnect,session,nil,hearttime)
	client.SetHandler(handleAgent)
	c:= &Connector{
		NetType: Type_Scoket,
		Client: client,
		reconnect: 0,
	}
	return c
}

func (c *Connector)StartConnect()bool{
	defer try.Catch()

	switch c.NetType{
	case Type_Scoket:
		//var err error
		conn, err := net.Dial("tcp", c.Client.Session.Addr) //
		if err != nil {
			log.Errorf("connect[%s] is fail,err:%s", c.Client.Session.Addr, err.Error())
			return false
		}
		c.Client.Session.Conn=conn

		go c.Client.Start()

		return true
	case Type_WebScoket:
		log.Fatalln("webscoket 连接未实现")
	default:
		log.Fatalln("未知连接方式")
	}
	return false
}
func (c *Connector)ReConnect(){
	if c.reconnect>0{

	}
}

func (c *Connector)SendData(msg []byte){
	defer try.Catch()
	c.Client.Session.Send(msg)
}

func (c *Connector)GetAddrInfo()string{
	if c.Client!=nil && c.Client.Session!=nil{
		return c.Client.Session.Addr
	}
	return ""
}

//func (c *Connector)Read(){
//	defer try.Catch()
//	defer func() {
//		if c.status == STATUS_CONNECT { //socket 意外中断时，进入重连
//			c.ReConnect()
//		}
//	}()
//
//	for {
//		n, err := c.Session.Conn.Read(c.Session.ReadBuf) //读取scoket传来的内容
//		if err != nil {
//			log.Infof("[%s] Read err:", c.Session.Addr,err.Error())
//			return //当远程客户端连接发生错误（断开）后，终止此协程。
//		}
//		if n > 0 && n <= messageMaxLen &&c.CallBack!=nil {
//			c.CallBack(n,c.Session.ReadBuf)
//		}
//
//	}
//}
//
//func (c *Connector)Write(){
//	defer try.Catch()
//	defer c.Close()
//
//	for {
//		select {
//		case msg := <-c.Session.writeChan:
//			_, err := c.Session.Conn.Write(msg)
//			if err != nil {
//				log.Errorln("Connector Write err:", err)
//				//goto stop
//			}
//		case <-time.Tick(c.heart):
//			_, err := c.Session.Conn.Write([]byte{200})
//			if err != nil {
//				log.Errorln("Connector Heart err:", err)
//				//goto stop
//				//心跳发生失败++
//			}
//		case <-c.Session.closeChan:
//			goto stop
//		}
//	}
//stop:
//	log.Infoln("Connector Write is stop:")
//}
//
//func (c *Connector)Close(){
//	if c.status == STATUS_CONNECT {
//		c.status = STATUS_CLOSE
//		//关闭socket
//		c.Session.Close()
//	}
//}