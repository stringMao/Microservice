package main

import (
	"Common/constant"
	"Common/setting"
	"Common/svrfind"
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
	//mysql连接
	dbmanager.ConnectDB_LoginSvr()
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

	go webagent.Start(config.App.Base.WebManagerPort)
}

//服务注册
func registerToDiscovery() {
	registration := svrfind.CreateRegistration()
	registration.ID = config.App.Base.GetServerIDStr()
	registration.Name = config.App.Base.GetServerName()
	registration.Port = config.App.Port
	//registration.TaggedAddresses["Client"]=
	registration.Tags = []string{config.App.Base.GetServerTag()}
	registration.Address = "127.0.0.1"
	registration.Check = svrfind.CreateCheck(registration.Address, config.App.Base.WebManagerPort)
	svrfind.Register(registration)

	//拉取该服务需要知道的其他服务状态
	go svrbalanced.RefreshSvrList(setting.GetServerName(constant.TID_GateSvr), "")
}
