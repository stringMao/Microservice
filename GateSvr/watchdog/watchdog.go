package watchdog

import (
	"Common/log"
	"Common/util"
	"GateSvr/agent"
	"GateSvr/config"
	"fmt"
	"net"
)

const (
	//服务器用户监听
	listenType_Server = iota
	// 客户端用户监听
	listenType_Client
)

//代理管理对象创建
var agentmanager *agent.AgentManager

//业务启动
func Start() {
	//创建代理管理器
	agentmanager = agent.NewAgentManager()
	//开启服务器连接的监听
	go StartTCPListen(config.App.ServerPort, listenType_Server)
	//开启客户端连接的监听
	go StartTCPListen(config.App.ClientPort, listenType_Client)
}

//Scoket开始监听
func StartTCPListen(port int, typeid int) {
	netListen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", util.GetLocalIP(), port))
	if err != nil {
		log.Logger.Fatal(err)
		return
	}
	defer netListen.Close()

	for {
		conn, err := netListen.Accept()
		if err != nil {
			log.Logger.Error(err)
			continue
		}
		if typeid == listenType_Client { //监听客户端连接
			go handleClientConnection(conn)
		} else if typeid == listenType_Server {
			go handleServerConnection(conn)
		}

	}
}
