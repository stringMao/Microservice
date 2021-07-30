package agent

import (
	"fmt"
	"math"
	"math/rand"
	"net"
	"sync"
)

//var m_conManagerMap map[int64]agentManager.Agent = make(map[int64]agentManager.Agent, 1)
type AgentManager struct {
	AgentClients map[uint64]*agentClient

	rwSvrlock    *sync.RWMutex       //服务列表的读写锁
	serverList   map[uint16][]uint16 //tid对应的sid列表
	AgentServers map[uint64]*agentServer
}

func NewAgentManager() *AgentManager {
	return &AgentManager{
		AgentClients: make(map[uint64]*agentClient, 1),
		rwSvrlock:    new(sync.RWMutex),
		serverList:   make(map[uint16][]uint16, 10),
		AgentServers: make(map[uint64]*agentServer, 1),
	}
}

//添加一个客户端代理对象
func (m *AgentManager) AddAgentClient(userid uint64, c net.Conn) *agentClient {

	if old, ok := m.AgentClients[userid]; ok { //水平扩展下，顶号如何判断?
		//已有连接，原连接断开
		old.Close()
		//return old
	}
	a := NewAgentClient(userid, c)
	fmt.Println(len(m.AgentClients))
	m.AgentClients[userid] = a
	return a
}

func (m *AgentManager) RemoveAgentClient(userid uint64) {
	if a, ok := m.AgentClients[userid]; ok {
		a.Close()
		delete(m.AgentClients, userid)
	}
}

//服务器注册
func (m *AgentManager) AddAgentServer(serverid uint64, c net.Conn) *agentServer {
	s := NewAgentServer(serverid, c)
	m.rwSvrlock.Lock()
	defer m.rwSvrlock.Unlock()
	//重复注册的，关闭前面的。
	if old, ok := m.AgentServers[serverid]; ok {
		old.Close()
	} else {
		//新注册的，需要在serverList记录
		tid, sid := decodeServerID(serverid)
		if l, ok := m.serverList[tid]; ok {
			l = append(l, sid)
		} else {
			l = make([]uint16, 1)
			l = append(l, sid)
		}
	}
	m.AgentServers[serverid] = s

	return s
}

//移除一个服务器连接
func (m *AgentManager) RemoveAgentServer(serverid uint64) {
	m.rwSvrlock.Lock()
	defer m.rwSvrlock.Unlock()
	if a, ok := m.AgentServers[serverid]; ok {
		a.Close()
		delete(m.AgentServers, serverid)
	}

	tid, sid := decodeServerID(serverid)
	if l, ok := m.serverList[tid]; ok {
		for i := 0; i < len(l); {
			if l[i] == sid {
				l = append(l[:i], l[i+1:]...)
			} else {
				i++
			}
		}
	}
}

//客户端消息转发给指定的服务器
func (m *AgentManager) TransferToServer(serverid uint64, msg []byte, len int) bool {
	m.rwSvrlock.RLock()
	defer m.rwSvrlock.RUnlock()
	if s, ok := m.AgentServers[serverid]; ok && s.open {
		s.send <- msg
		return true
	}
	return false
}

func (m *AgentManager) TransferToClient(userid uint64, msg []byte) bool {
	//大端拼接
	if c, ok := m.AgentClients[userid]; ok {
		c.send <- msg
		return true
	}
	return false
}

//分配一个指定类型的服务
func (m *AgentManager) AllocSvr(tid uint16) (uint64, *agentServer) {
	m.rwSvrlock.RLock()
	defer m.rwSvrlock.RUnlock()

	if l, ok := m.serverList[tid]; ok && l != nil && len(l) > 0 {
		r := rand.Intn(len(l)) //随机负载均衡
		sid := l[r]
		serverid := encodeServerID(tid, sid)

		if pSvr, ok2 := m.AgentServers[serverid]; ok2 {
			return serverid, pSvr
		}
	}
	return 0, nil
}

//解码serverid to tid sid
func decodeServerID(serverid uint64) (tid, sid uint16) {
	sid = uint16(serverid & (math.MaxUint16 << 16))
	tid = uint16(serverid & math.MaxUint16)
	return tid, sid
}
func encodeServerID(tid, sid uint16) uint64 {
	return uint64(sid)<<16 + uint64(tid)
}
