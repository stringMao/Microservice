package kernel

type ManagerSvrs struct {
	GateSvrMap map[uint64]*ServerAgent
}

//单例模式==》饿汉模式
var instance *ManagerSvrs

func init() {
	instance = new(ManagerSvrs)
	instance.GateSvrMap = make(map[uint64]*ServerAgent)
}
func GetManagerSvrs() *ManagerSvrs {
	return instance
}

func (m *ManagerSvrs) AddGateSvr(serverid uint64, agent *ServerAgent) {
	agent.ConnObject.KeyMap = serverid
	m.GateSvrMap[serverid] = agent
}

func (m *ManagerSvrs) SendData(key uint64, data []byte) bool {
	if p, ok := m.GateSvrMap[key]; ok {
		return p.ConnObject.SendData(data)
	}
	return false
}
