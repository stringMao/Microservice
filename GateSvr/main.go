package main

import (
	"Common/log"
	"Common/svrfind"
	"Common/util"
	"Common/webmanager"
	"GateSvr/config"
	"GateSvr/dbmanager"
	"GateSvr/watchdog"
	"math/rand"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
)

func init() {
	// 设置随机数种子
	rand.Seed(time.Now().Unix())
	//读取系统配置app.ini
	config.Init()
}

func main() {
	// defer func() { //必须要先声明defer，否则不能捕获到panic异常
	// 	if err := recover(); err != nil {
	// 		fmt.Println(err) //这里的err其实就是panic传入的内容
	// 	}
	// }()
	dbmanager.Init()

	//打开socket监听
	watchdog.Start()

	//web后台功能开启
	registerWebManagerRoute()

	//服务注册
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

	go webagent.Start(config.App.WebManagerPort)
}

//服务注册
func registerToDiscovery() {
	svritem := svrfind.NewServerItem(config.App.ConsulAddr)
	svritem.SvrData.ID = config.App.GetServerIDStr()
	svritem.SvrData.Name = config.App.GetServerName() //本服务的名字
	svritem.SvrData.Port = config.App.ClientPort
	svritem.SvrData.Tags = []string{config.App.GetServerTag()}
	svritem.SvrData.Address = util.GetLocalIP()
	svritem.SvrData.TaggedAddresses = make(map[string]api.ServiceAddress)
	svritem.SvrData.TaggedAddresses["client"] = api.ServiceAddress{Address: svritem.SvrData.Address, Port: config.App.ClientPort}
	svritem.SvrData.TaggedAddresses["server"] = api.ServiceAddress{Address: svritem.SvrData.Address, Port: config.App.ServerPort}

	//svritem.SvrData.Check = svritem.CreateAgentServiceCheck(config.App.Base.WebManagerPort)
	if svritem.Register(config.App.WebManagerPort) {
		log.Infoln("服务注册成功!")
	}

	// registration := svrfind.CreateRegistration()
	// registration.ID = config.App.GetServerIDStr()
	// registration.Name = config.App.GetServerName()
	// registration.Port = config.App.ClientPort

	// registration.Tags = []string{config.App.GetServerTag()}
	// registration.Address = util.GetLocalIP()
	// registration.TaggedAddresses = make(map[string]api.ServiceAddress)
	// registration.TaggedAddresses["client"] = api.ServiceAddress{Address: registration.Address, Port: config.App.ClientPort}
	// registration.TaggedAddresses["server"] = api.ServiceAddress{Address: registration.Address, Port: config.App.ServerPort}

	// //registration.Meta = make(map[string]string)
	// //registration.Meta["online"] = strconv.Itoa(99)
	// //registration.Meta["test"] = "test"

	// registration.Check = svrfind.CreateCheck(registration.Address, config.App.WebManagerPort)
	// svrfind.Register(config.App.ConsulAddr, registration)
}
