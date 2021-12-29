package main

import (
	"Common/constant"
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
	//config.Init()
}

func main() {
	// defer func() { //必须要先声明defer，否则不能捕获到panic异常
	// 	if err := recover(); err != nil {
	// 		fmt.Println(err) //这里的err其实就是panic传入的内容
	// 	}
	// }()

	dbmanager.Init()

	//打开socket业务监听
	watchdog.Start()

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

	//启动web服务
	return webmanager.G_WebManager.Start(config.App.WebManagerPort)
}

//服务注册
func registerToDiscovery() bool {
	svrfind.G_ServerRegister.SvrData.ID = constant.GetServerIDName(config.App.TID, config.App.SID)
	svrfind.G_ServerRegister.SvrData.Name = constant.GetServerName(config.App.TID) //本服务的名字
	svrfind.G_ServerRegister.SvrData.Port = config.App.ClientPort
	svrfind.G_ServerRegister.SvrData.Tags = []string{constant.GetServerTag(config.App.TID)}
	svrfind.G_ServerRegister.SvrData.Address = util.GetLocalIP()
	svrfind.G_ServerRegister.SvrData.TaggedAddresses = make(map[string]api.ServiceAddress)
	svrfind.G_ServerRegister.SvrData.TaggedAddresses["client"] = api.ServiceAddress{Address: svrfind.G_ServerRegister.SvrData.Address, Port: config.App.ClientPort}
	svrfind.G_ServerRegister.SvrData.TaggedAddresses["server"] = api.ServiceAddress{Address: svrfind.G_ServerRegister.SvrData.Address, Port: config.App.ServerPort}

	//svritem.SvrData.Check = svritem.CreateAgentServiceCheck(config.App.Base.WebManagerPort)
	return svrfind.G_ServerRegister.Register(config.App.ConsulAddr, config.App.WebManagerPort)

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
