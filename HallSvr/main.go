package main

import (
	"Common/log"
	"Common/svrfind"
	"Common/util"
	"Common/webmanager"
	"HallSvr/config"
	"HallSvr/watchdog"
	"math/rand"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	// 设置随机数种子
	rand.Seed(time.Now().Unix())
	//读取系统配置app.ini
	config.Init()
}

func main() {

	//web后台功能开启
	registerWebManagerRoute()
	//服务注册
	svritem := registerToDiscovery()

	//
	watchdog.ConnectGateSvr(svritem)

	c := make(chan os.Signal)
	<-c
}

//注册后台路由，启动路由监听
func registerWebManagerRoute() {
	webagent := webmanager.CreateRouterAgent()

	//健康检查接口
	consulCheckHealth := webmanager.RouterHelper{
		Type:   webmanager.RouterType_consul,
		Path:   "/check",
		Method: "GET",
		Handlers: []gin.HandlerFunc{
			svrfind.CheckHealth,
		},
	}
	webagent.RegisterRouter(consulCheckHealth)

	go webagent.Start(config.App.WebManagerPort)
}

//服务注册
func registerToDiscovery() *svrfind.ServerItem {
	svritem := svrfind.NewServerItem(config.App.ConsulAddr)
	svritem.SvrData.ID = config.App.GetServerIDStr()
	svritem.SvrData.Name = config.App.GetServerName() //本服务的名字
	svritem.SvrData.Port = config.App.Port
	svritem.SvrData.Tags = []string{config.App.GetServerTag()}
	svritem.SvrData.Address = util.GetLocalIP()
	//svritem.SvrData.Check = svritem.CreateAgentServiceCheck(config.App.Base.WebManagerPort)
	if svritem.Register(config.App.WebManagerPort) {
		log.Infoln("服务注册成功!")
	}
	return svritem
}
