package main

import (
	"Common/svrfind"
	"Common/util"
	"Common/webmanager"
	"GateSvr/config"
	"GateSvr/watchdog"
	"math/rand"
	"os"
	"strconv"
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

	go webagent.Start(config.App.Base.WebManagerPort)
}

//服务注册
func registerToDiscovery() {
	registration := svrfind.CreateRegistration()
	registration.ID = config.App.Base.GetServerIDStr()
	registration.Name = config.App.Base.GetServerName()
	registration.Port = config.App.ClientPort

	registration.Tags = []string{config.App.Base.GetServerTag()}
	registration.Address = util.GetLocalIP()
	registration.TaggedAddresses = make(map[string]api.ServiceAddress)
	registration.TaggedAddresses["client"] = api.ServiceAddress{Address: registration.Address, Port: config.App.ClientPort}
	registration.TaggedAddresses["server"] = api.ServiceAddress{Address: registration.Address, Port: config.App.ServerPort}

	registration.Meta = make(map[string]string)
	registration.Meta["online"] = strconv.Itoa(99)
	registration.Meta["test"] = "test"

	registration.Check = svrfind.CreateCheck(registration.Address, config.App.Base.WebManagerPort)
	svrfind.Register(registration)
}
