package kernel

type ServerAgent struct {
	BaseData   *ServerBaseData
	ConnObject *ConnTcp
}

func NewSeverAgent(c *ConnTcp) *ServerAgent {
	agent := &ServerAgent{
		BaseData:   &ServerBaseData{},
		ConnObject: c,
	}
	return agent
}

func (agent *ServerAgent) Close() {
	if agent.ConnObject != nil {
		agent.ConnObject.CloseConnet()
	}
}

func (agent *ServerAgent) StartWork() {
	if agent.ConnObject != nil {
		//agent.ConnObject.HeartSwitch = true
	}
}

//SetConnectData 设置连接的基本信息
// func (agent *ServerAgent) SetConnectData(ip string, port int, f func([]byte, int)) {
// 	agent.ConnObject.SvrIP = fmt.Sprintf("%s:%d", ip, port)
// 	agent.ConnObject.CallBackFunc = f
// }

// func (agent *ServerAgent) LoginGateSvr() bool {
// 	if agent.ConnObject.Connect() {
// 		//发送登入消息
// 		//登入验证
// 		logindata := &base.ServerLogin{
// 			Tid:      uint32(config.App.TID),
// 			Sid:      uint32(config.App.SID),
// 			Password: "test",
// 		}
// 		pData, err := proto.Marshal(logindata)
// 		if err != nil {
// 			log.Errorf("ConnectGateSvr proto.Marshal err:%s", err)
// 		}
// 		//agentGateSvr.SendData(pData)
// 		agent.ConnObject.SendData(pData)
// 		return true
// 	}
// 	return false
// }
