package agent

import (
	"Common/log"
	"GateSvr/core/player"
	"net"
	"runtime/debug"
	"sync"
)

//客户端对象
type AgentClient struct {
	Userid        uint64
	conn          net.Conn //客户端连接
	send          chan []byte
	sendEnd       chan bool         //发送消息的协程开关
	svrlist       map[uint32]uint64 //该用户连接的服务器列表map[tid]serverid（每种服务器只能连一个 serverid=0+0+sid+tid）
	svrlistRWLock *sync.RWMutex     //对服务器列表的读写锁
	//svrMap map[uint16]*AgentServer //服务器列表

	Player *player.Player //玩家数据
}

func NewAgentClient(userid uint64, c net.Conn) *AgentClient {
	a := &AgentClient{
		Userid:        userid,
		conn:          c,
		send:          make(chan []byte, 100),
		sendEnd:       make(chan bool),
		svrlist:       make(map[uint32]uint64, 5),
		svrlistRWLock: new(sync.RWMutex),
		//svrMap: make(map[uint16]*AgentServer, 5),
	}

	go func(a *AgentClient) {
		defer func() {
			a.sendEnd <- true
		}()
		defer func() {
			if r := recover(); r != nil {
				log.PrintPanicStack(r, string(debug.Stack()))
				a.conn.Close()
			}
		}()

		for {
			if msg, ok := <-a.send; ok {
				a.conn.Write(msg)
			} else {
				break
			}
		}
		//fmt.Printf("客户端发送携程关闭:userid[%d]\n", a.Userid)
	}(a)

	return a
}

//
func (a *AgentClient) Close() {
	//通知其他业务服务器该用户离线

	//缓存同步
	a.Save()

	//关闭子携程 close要保证不可以重复
	close(a.send) //close之后，缓存区还有数据，还是会返回ok=true，直到缓冲区清空

	//确保发送协程清空并且结束
	<-a.sendEnd
	//关闭连接
	a.conn.Close()
}

//给该链接发送消息
func (a *AgentClient) SendData(msg []byte) (suc bool) {
	defer func() {
		if recover() != nil {
			suc = false //发送失败
		}
	}()
	a.send <- msg //a.send chan在被close之后，插入数据会异常
	return true
}

func (a *AgentClient) GetServerId(tid uint32) uint64 {
	a.svrlistRWLock.RLock()
	defer a.svrlistRWLock.RUnlock()
	if s, ok := a.svrlist[tid]; ok {
		return s
	}
	return 0
}

func (a *AgentClient) SetServerId(tid uint32, serverid uint64) {
	a.svrlistRWLock.Lock()
	defer a.svrlistRWLock.Unlock()
	a.svrlist[tid] = serverid
}

func (a *AgentClient) GetAllSvrList() []uint64 {
	r := []uint64{}
	for _, val := range a.svrlist {
		r = append(r, val)
	}
	return r
}
