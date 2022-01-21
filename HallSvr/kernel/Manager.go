package kernel

import "sync"

type ManagerSvrs struct {
	GateSvrMap map[uint64]*ServerAgent
	rwlock     *sync.RWMutex
}

//单例模式==》饿汉模式
var instance *ManagerSvrs

func init() {
	instance = new(ManagerSvrs)
	instance.GateSvrMap = make(map[uint64]*ServerAgent)
	instance.rwlock = new(sync.RWMutex)
}
func GetManagerSvrs() *ManagerSvrs {
	return instance
}

func (m *ManagerSvrs) AddGateSvr(serverid uint64, agent *ServerAgent) {
	agent.ConnObject.Id = serverid
	m.rwlock.Lock()
	defer m.rwlock.Unlock()
	m.GateSvrMap[serverid] = agent
}
func (m *ManagerSvrs) DeleteGateSvr(serverid uint64) {
	m.rwlock.Lock()
	defer m.rwlock.Unlock()
	if p, ok := m.GateSvrMap[serverid]; ok {
		p.Close()
		delete(m.GateSvrMap, serverid)
	}
}

func (m *ManagerSvrs) SendData(key uint64, data []byte) bool {
	m.rwlock.RLock()
	defer m.rwlock.RUnlock()
	if p, ok := m.GateSvrMap[key]; ok {
		return p.ConnObject.SendData(data)
	}
	return false
}

func (m *ManagerSvrs) IsExist(serverid uint64) bool {
	m.rwlock.RLock()
	defer m.rwlock.RUnlock()
	if _, ok := m.GateSvrMap[serverid]; ok {
		return true
	}
	return false
}

func SendData(key uint64, data []byte) bool {
	return instance.SendData(key, data)
}
