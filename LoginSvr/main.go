package main

import (
	"Common/constant"
	"Common/log"
	"Common/setting"
	"Common/svrfind"
	"Common/util"
	"Common/webmanager"
	"LoginSvr/config"
	"LoginSvr/dbmanager"
	"LoginSvr/global"
	"LoginSvr/router"
	"LoginSvr/svrbalanced"
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
	//defer ServerEnd()
	//mysql连接
	dbmanager.ConnectDB_LoginSvr()
	dbmanager.ConnectDB_Player()
	//
	dbmanager.ConnectRedis()
	//6.
	global.Init()

	//web后台功能开启
	registerWebManagerRoute()
	//7.业务路由
	router.Init()
	go router.Start(config.App.Port)

	//服务注册及服务发现
	registerToDiscovery()

	c := make(chan os.Signal)
	<-c
}

func ServerEnd() {
	//通知consul注销服务
	//svrfind.Deregister(config.App.Base.GetServerIDStr())
}

//注册后台路由，启动路由监听
func registerWebManagerRoute() {
	webagent := webmanager.CreateRouterAgent()

	//供consul健康检查接口
	consulCheckHealth := webmanager.RouterHelper{
		Type:   webmanager.RouterType_consul,
		Path:   "/check",
		Method: "GET",
		Handlers: []gin.HandlerFunc{
			svrfind.CheckHealth,
		},
	}
	webagent.RegisterRouter(consulCheckHealth)

	go webagent.Start(config.App.Base.WebManagerPort)
}

//服务注册
func registerToDiscovery() {
	svritem := svrfind.NewServerItem(config.App.Base.ConsulAddr)
	svritem.SvrData.ID = config.App.Base.GetServerIDStr()
	svritem.SvrData.Name = config.App.Base.GetServerName() //本服务的名字
	svritem.SvrData.Port = config.App.Port
	svritem.SvrData.Tags = []string{config.App.Base.GetServerTag()}
	svritem.SvrData.Address = util.GetLocalIP()
	//svritem.SvrData.Check = svritem.CreateAgentServiceCheck(config.App.Base.WebManagerPort)
	if svritem.Register(config.App.Base.WebManagerPort) {
		log.Infoln("服务注册成功!")
	}

	//拉取该服务需要知道的其他服务状态
	go svrbalanced.RefreshSvrList(svritem, setting.GetServerName(constant.TID_GateSvr), "")
}
