package kernel

import (
	"Common/log"
	"fmt"
	"net"
	"runtime/debug"
	"sync"
	"time"
)

type ConnTcp struct {
	conn         net.Conn
	open         bool
	SvrIP        string
	send         chan []byte //消息发送通道
	wg           sync.WaitGroup
	callBackFunc func([]byte, int) //消息回调接口
}

func NewConnTcp(ip string, port int, f func([]byte, int)) *ConnTcp {
	if f == nil {
		log.Fatalln("NewConnTcp callBackFunc is nil")
		return nil
	}
	return &ConnTcp{
		open:  false,
		SvrIP: fmt.Sprintf("%s:%d", ip, port), //
		send:  make(chan []byte, 1000),
		//rev:   make(chan []byte, 1000),

		callBackFunc: f,
	}
}

func (c *ConnTcp) Connect() bool {
	var err error
	c.conn, err = net.Dial("tcp", c.SvrIP) //本机必须用127.0.0.1。直接用ip没法连接
	if err != nil {
		log.Errorln("connect", c.SvrIP, "fail", err.Error())
		return false
	}
	c.open = true

	//发送消息携程打开
	c.openWrite()

	//接收消息携程打开
	c.openRead()

	//心跳包发送携程打开
	c.openHeart()
	return true
}

//主动关闭scoket连接的流程
func (c *ConnTcp) CloseConnet() {
	//心跳 被关闭
	c.open = false

	//close通道之后，新的发送消息无法再进入，发送协程在发送完缓存区的消息后被关闭
	close(c.send)

	//等待所有协程关闭
	c.wg.Wait()

	//关闭socket
	c.conn.Close()

	//新消息读取协程 会捕捉到err退出

}

//scoket连接被中断
func (c *ConnTcp) DisConnet() {

}

func (c *ConnTcp) openRead() {
	//c.wg.Add(1)

	go func() {
		defer func() {
			if c.open { //socket 中断时，会进入
				c.DisConnet()
			}
		}()
		//defer c.wg.Done()
		buffer := make([]byte, 2048)
		for c.open {
			n, err := c.conn.Read(buffer) //读取scoket传来的内容
			if err != nil {
				log.Errorln("ConnTcp openRead err:", err, n)
				return //当远程客户端连接发生错误（断开）后，终止此协程。
			}
			//send <- buffer[:n]
			msg := make([]byte, n)
			copy(msg, buffer[:n])
			c.callBackFunc(msg, n)
		}
	}()
}
func (c *ConnTcp) openWrite() {
	c.wg.Add(1)

	go func() {
		defer c.wg.Done()

		defer func() {
			if r := recover(); r != nil {
				log.PrintPanicStack(r, string(debug.Stack()))
				c.conn.Close()
			}
		}()

		for {
			//chan被close之后，缓存区还有数据，还是会返回ok=true，直到缓冲区清空
			if msg, ok := <-c.send; ok {
				_, err := c.conn.Write(msg)
				if err != nil {
					log.Errorln("ConnTcp openWrite err:", err)
					return
				}
			} else {
				break
			}
		}

		// if c.open { //socket 中断时，会进入
		// 	c.DisConnet()
		// }

	}()
}

func (c *ConnTcp) openHeart() {
	c.wg.Add(1)
	//心跳包
	go func() {
		defer c.wg.Done()

		for c.open {
			time.Sleep(time.Second * 1)
			//c.send <- []byte{200}
			if !c.SendData([]byte{200}) {
				return
			}
		}
	}()
}

func (c *ConnTcp) SendData(msg []byte) (suc bool) {
	defer func() {
		if recover() != nil {
			suc = false //发送失败
		}
	}()
	c.send <- msg //s.send chan在被close之后，插入数据会异常
	return true
}
