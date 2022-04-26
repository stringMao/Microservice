package scokets

import (
	"Common/log"
	"net"
	//uuid "github.com/satori/go.uuid"
	"time"
)

type Listener struct {
	Handler
	NetType   		int  //1=scoket  2=webscoket
	Addr       		string
	CallBack        func(client *Client)
}

func NewListener(netType int,addr string,callback func(client *Client))*Listener{
	return &Listener{
		NetType: netType,
		Addr: addr,
		CallBack: callback,
	}
}

func (l Listener)StartListen(){
	go func(){
		//defer try.Catch()
		if l.NetType==Type_Scoket{
			netListen, err := net.Listen("tcp", l.Addr)
			if err != nil {
				log.Fatalln(err)
				return
			}
			defer netListen.Close()
			log.Debugf("[%s]监听成功",l.Addr)

			for {
				if conn, err := netListen.Accept();err==nil{
					session:=NewSession(0,conn,100)
					cli:=NewClient(ServerConnect,session,nil,5*time.Second)
					//go cli.Start()
					if l.CallBack!=nil{
						go l.CallBack(cli)
					}
				}
			}
		}else if l.NetType==Type_WebScoket{
			log.Fatalln("webscoket 未实现")
			return
		}
		log.Fatalln("StartListen 未知类型")
		return
	}()
}








