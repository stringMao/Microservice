package main

import (
	"Common/constant"
	"Common/log"
	"Common/svrfind"
	"Common/webmanager"
	"LoginSvr/config"
	"LoginSvr/dbmanager"
	"LoginSvr/global"
	"LoginSvr/router"
	"LoginSvr/svrbalanced"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
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
	go func() {
		log.Debugln(http.ListenAndServe(":8060", nil))
	}()
	//defer ServerEnd()
	//log.Logger.Debugln("1111111")
	//log.Debugln("22222222")
	//mysql连接
	dbmanager.ConnectDB_LoginSvr()
	dbmanager.ConnectDB_Player()
	//
	dbmanager.ConnectRedis()
	//6.
	global.Init()

	//web后台功能开启
	if !registerWebManagerRoute() {
		log.Fatalln("web服务启动失败!")
	}
	log.Infoln("web服务启动成功!")

	//7.业务路由
	router.Init()
	go router.Start(config.App.Port)

	//服务注册
	if !registerToDiscovery() {
		log.Fatalln("服务注册失败!")
	}
	log.Infoln("服务注册成功!")

	//拉取该服务需要知道的其他服务状态
	svrbalanced.RefreshGateSvrList()

	c := make(chan os.Signal)
	<-c
}

func ServerEnd() {
	//通知consul注销服务
	//svrfind.Deregister(config.App.Base.GetServerIDStr())
}

//注册后台路由，启动路由监听
func registerWebManagerRoute() bool {
	//供consul健康检查接口
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
	svrfind.G_ServerRegister.SvrData.Address = config.App.WebManagerIP
	//svritem.SvrData.Check = svritem.CreateAgentServiceCheck(config.App.Base.WebManagerPort)
	return svrfind.G_ServerRegister.Register(config.App.ConsulAddr, config.App.WebManagerPort)

}
