package kernel

import (
	"Common/log"
	"net"
	"runtime/debug"
	"sync"
	"time"
)

const (
	STATUS_INIT = iota
	STATUS_CONNECT
	STATUS_CLOSE
)

type ConnTcp struct {
	Id            uint64
	Addr          string
	conn          net.Conn
	status        int         //连接状态码
	writeChan     chan []byte //消息发送通道
	HeartSwitch   bool        //心跳开关
	stopWriteChan chan int
	wg            sync.WaitGroup
	callBackFunc  func(uint64, []byte, uint32) //消息回调接口
	errFunc       func(uint64, int)
	loginData     []byte //登入需要发送的消息体
}

func NewConnTcp(id uint64, addr string, logindata []byte, callback func(uint64, []byte, uint32), errhandle func(uint64, int)) *ConnTcp {
	return &ConnTcp{
		Id:            id,
		status:        STATUS_INIT,
		Addr:          addr, //
		writeChan:     make(chan []byte, 1000),
		stopWriteChan: make(chan int, 1),
		HeartSwitch:   false,
		callBackFunc:  callback,
		errFunc:       errhandle,
		loginData:     logindata,
	}
}

func (c *ConnTcp) Connect() bool {
	var err error
	c.conn, err = net.Dial("tcp", c.Addr) //
	if err != nil {
		log.Errorln("connect", c.Addr, "fail", err.Error())
		return false
	}
	_, err = c.conn.Write(c.loginData)
	if err != nil {
		c.conn.Close()
		return false
	}
	c.status = STATUS_CONNECT

	//发送消息携程打开
	go c.openWrite()

	//接收消息携程打开
	go c.openRead()

	//心跳包发送打开
	c.HeartSwitch = true
	return true
}
func (c *ConnTcp) ReConnect(maxcount int) bool {
	log.Errorln("开始重连===Addr:[%s]", c.Addr)
	count := 0
	for {
		if c.Connect() {
			log.Errorln("重连成功===Addr:[%s]", c.Addr)
			return true
		}
		count++
		if count >= maxcount {
			log.Errorln("重连失败===Addr:[%s]", c.Addr)
			return false
		}
		time.Sleep(1 * time.Second)
	}
}

//主动关闭scoket连接的流程
func (c *ConnTcp) CloseConnet() {
	if c.status == STATUS_CONNECT {
		c.status = STATUS_CLOSE
		c.HeartSwitch = false
		c.stopWriteChan <- 1

		//等待所有协程关闭
		c.wg.Wait()
		//close通道之后，新的发送消息无法再进入，发送协程在发送完缓存区的消息后被关闭
		//close(c.writeChan)
		//close(c.stopWriteChan)

		//关闭socket
		c.conn.Close()
	}

}

//scoket连接异常中断，需要重连
func (c *ConnTcp) Disconnet() {
	c.CloseConnet()
	//尝试重连
	if !c.ReConnect(2) {
		//重连失败后向上通知
		c.errFunc(c.Id, 1)
	}
}

func (c *ConnTcp) openRead() {
	defer func() {
		if c.status == STATUS_CONNECT { //socket 意外中断时，进入重连
			c.Disconnet()
		}
	}()

	for {
		if c.status != STATUS_CONNECT {
			break
		}
		buffer := make([]byte, 2048)
		n, err := c.conn.Read(buffer) //读取scoket传来的内容
		if err != nil {
			log.Errorln("ConnTcp openRead err:", err, n)
			return //当远程客户端连接发生错误（断开）后，终止此协程。
		}
		if n > 0 && n < 2048 {
			c.callBackFunc(c.Id, buffer, uint32(n))
		}

	}

}

func (c *ConnTcp) openWrite() {
	c.wg.Add(1)
	defer c.wg.Done()

	defer func() {
		if r := recover(); r != nil {
			log.PrintPanicStack(r, string(debug.Stack()))
			c.conn.Close()
		}
	}()
	for {
		select {
		case msg := <-c.writeChan:
			_, err := c.conn.Write(msg)
			if err != nil {
				log.Errorln("ConnTcp openWrite Write err:", err)
				//goto stop
			}
		case <-time.Tick(1 * time.Second):
			if c.HeartSwitch {
				_, err := c.conn.Write([]byte{200})
				if err != nil {
					log.Errorln("ConnTcp openWrite Heart err:", err)
					//goto stop
					//心跳发生失败++
				}
			}
		case <-c.stopWriteChan:
			goto stop
		}
	}
stop:
	log.Infoln("openWrite is stop:")
}

func (c *ConnTcp) SendData(msg []byte) (suc bool) {
	defer func() {
		if recover() != nil {
			suc = false //发送失败
		}
	}()
	c.writeChan <- msg //s.send chan在被close之后，插入数据会异常
	return true
}
