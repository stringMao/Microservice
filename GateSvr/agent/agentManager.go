package agent

import (
	"Common/log"
	"Common/msg"
	"math/rand"
	"net"
	"sync"
	"time"
)

type ServerObj map[uint16]*agentServer //uint16=sid

//var m_conManagerMap map[int64]agentManager.Agent = make(map[int64]agentManager.Agent, 1)
type AgentManager struct {
	AgentClients map[uint64]*agentClient //客户端列表
	rwClientLock *sync.RWMutex

	serverList   map[uint32][]uint32 //tid对应的sid列表
	AgentServers map[uint64]*agentServer
	rwSvrlock    *sync.RWMutex //服务列表的读写锁

}

func NewAgentManager() *AgentManager {
	return &AgentManager{
		AgentClients: make(map[uint64]*agentClient, 1),
		rwClientLock: new(sync.RWMutex),
		serverList:   make(map[uint32][]uint32, 10),
		AgentServers: make(map[uint64]*agentServer, 100),
		rwSvrlock:    new(sync.RWMutex),
	}
}

//客户端==start======================================================================================

//客户端的顶号处理  //水平扩展下，顶号如何判断?？？？？？？
func (m *AgentManager) ReplaceClient(userid uint64) bool {
	var applyclose = false
	for i := 0; i < 3; i++ {
		m.rwClientLock.RLock()
		old, ok := m.AgentClients[userid]
		m.rwClientLock.RUnlock()
		if ok {
			if !applyclose {
				old.conn.Close()  //原连接断开
				applyclose = true //保证值申请一次关闭，避免重复
			}
			//被顶处，安全退出及数据保存需要时间
			time.Sleep(1000)
		} else {
			return true //已经没有就连接存在
		}
	}
	return false //尝试多次，旧连接依然存在，则线登入失败，让客户端重新登入
}

//添加一个客户端代理对象
func (m *AgentManager) AddAgentClient(userid uint64, c net.Conn) *agentClient {
	//========================
	a := NewAgentClient(userid, c)

	log.Debug("AgentClients 数量：", len(m.AgentClients))

	m.rwClientLock.Lock()
	m.AgentClients[userid] = a
	m.rwClientLock.Unlock()
	return a
}

//移除客户端对象
func (m *AgentManager) RemoveAgentClient(userid uint64) {
	m.rwClientLock.RLock()
	agent, ok := m.AgentClients[userid]
	m.rwClientLock.RUnlock()
	if ok {
		//先关闭代理对象
		agent.Close()
		//再清理队列里的对象指针(确保被清理的对象已经完成内存保存)
		m.rwClientLock.Lock()
		delete(m.AgentClients, userid)
		m.rwClientLock.Unlock()
	}
}

//客户端==end======================================================================================

//重复注册检查
func (m *AgentManager) ReplaceServer(tid, sid uint32) bool {
	serverid := msg.EncodeServerID(tid, sid)
	m.rwSvrlock.RLock()
	defer m.rwSvrlock.RUnlock()
	if _, ok := m.AgentServers[serverid]; ok {
		return false //重复注册，则无效
	}
	return true
}

//服务器注册
func (m *AgentManager) AddAgentServer(tid, sid uint32, c net.Conn) *agentServer {
	s := NewAgentServer(tid, sid, c)
	m.rwSvrlock.Lock()
	defer m.rwSvrlock.Unlock()

	//新注册的，需要在serverList记录
	sidlist, ok := m.serverList[tid]
	if !ok {
		sidlist = make([]uint32, 10)
		m.serverList[tid] = sidlist
	}
	sidlist = append(sidlist, sid)
	m.AgentServers[s.Serverid] = s

	return s
}

//移除一个服务器连接
func (m *AgentManager) RemoveAgentServer(tid, sid uint32) {
	m.rwSvrlock.Lock()
	defer m.rwSvrlock.Unlock()

	if l, ok := m.serverList[tid]; ok {
		for i := 0; i < len(l); {
			if l[i] == sid {
				l = append(l[:i], l[i+1:]...)
			} else {
				i++
			}
		}
	}

	serverid := msg.EncodeServerID(tid, sid)
	if a, ok := m.AgentServers[serverid]; ok {
		a.Close()
		delete(m.AgentServers, serverid)
	}

}

//客户端消息转发给指定的服务器
func (m *AgentManager) TransferToServer(serverid uint64, msg []byte) bool {
	m.rwSvrlock.RLock()
	s, ok := m.AgentServers[serverid]
	m.rwSvrlock.RUnlock()
	if ok {
		s.send <- msg
		return true
	}

	return false
}

func (m *AgentManager) TransferToClient(userid uint64, msg []byte) bool {
	m.rwClientLock.RLock()
	c, ok := m.AgentClients[userid]
	m.rwClientLock.RUnlock()
	if ok {
		c.send <- msg
		return true
	}
	return false
}

//分配一个指定类型的服务
func (m *AgentManager) AllocSvr(tid uint32) uint64 {
	m.rwSvrlock.RLock()
	defer m.rwSvrlock.RUnlock()

	if l, ok := m.serverList[tid]; ok && l != nil && len(l) > 0 {
		r := rand.Intn(len(l)) //随机负载均衡
		sid := l[r]
		serverid := msg.EncodeServerID(tid, sid)

		return serverid
	}
	return 0
}
