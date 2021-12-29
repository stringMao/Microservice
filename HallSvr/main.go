package main

import (
	"Common/constant"
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
	//config.Init()
}

func main() {

	//web后台功能开启
	if !registerWebManagerRoute() {
		log.Fatalln("web服务启动失败!")
	}
	log.Infoln("web服务启动成功!")

	//服务注册
	if !registerToDiscovery() {
		log.Fatalln("服务注册失败!")
	}
	log.Infoln("服务注册成功!")

	//
	watchdog.ConnectGateSvrs()

	c := make(chan os.Signal)
	<-c
}

//注册后台路由，启动路由监听
func registerWebManagerRoute() bool {

	//健康检查接口
	consulCheckHealth := webmanager.RouterHelper{
		Type:   webmanager.RouterType_consul,
		Path:   "/check",
		Method: "GET",
		Handlers: []gin.HandlerFunc{
			svrfind.CheckHealth,
		},
	}
	webmanager.G_WebManager.RegisterRouter(consulCheckHealth)

	return webmanager.G_WebManager.Start(config.App.WebManagerPort)
}

//服务注册
func registerToDiscovery() bool {
	svrfind.G_ServerRegister.SvrData.ID = constant.GetServerIDName(config.App.TID, config.App.SID)
	svrfind.G_ServerRegister.SvrData.Name = constant.GetServerName(config.App.TID) //本服务的名字
	svrfind.G_ServerRegister.SvrData.Port = config.App.Port
	svrfind.G_ServerRegister.SvrData.Tags = []string{constant.GetServerTag(config.App.TID)}
	svrfind.G_ServerRegister.SvrData.Address = util.GetLocalIP()
	//svritem.SvrData.Check = svritem.CreateAgentServiceCheck(config.App.Base.WebManagerPort)
	return svrfind.G_ServerRegister.Register(config.App.ConsulAddr, config.App.WebManagerPort)
}
