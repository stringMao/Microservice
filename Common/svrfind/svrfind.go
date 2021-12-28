package svrfind

//服务发现
//consul API文档  https://www.consul.io/api-docs/config
//golang调用API文档  https://pkg.go.dev/github.com/hashicorp/consul/api
import (
	"Common/log"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	consulapi "github.com/hashicorp/consul/api"
)

type ServerItem struct {
	//ConsulAddr string //consul服务的地址
	Config *consulapi.Config
	Client *consulapi.Client
	//SvrIp   string //本服务ip
	SvrData *consulapi.AgentServiceRegistration
}

var G_ServerRegister *ServerItem

func init() {
	G_ServerRegister = new(ServerItem)
	G_ServerRegister.Config = consulapi.DefaultConfig()
	G_ServerRegister.Client = nil
	G_ServerRegister.SvrData = new(consulapi.AgentServiceRegistration)
}

//向consul注册服务
func (s *ServerItem) Register(consuladdr string, checkPort int) bool {
	G_ServerRegister.Config.Address = consuladdr
	var err error
	G_ServerRegister.Client, err = consulapi.NewClient(G_ServerRegister.Config)
	if err != nil {
		log.Errorln("consul new client error : ", err)
		return false
	}

	//健康检查接口信息
	s.SvrData.Check = &consulapi.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s:%d/consul%s", s.SvrData.Address, checkPort, "/check"),
		Timeout:                        "5s",  //设置超时 5s。
		Interval:                       "5s",  //设置间隔 5s。
		DeregisterCriticalServiceAfter: "30s", //check失败后30秒删除本服务
	}
	//注册
	err = s.Client.Agent().ServiceRegister(s.SvrData)
	if err != nil {
		log.Errorln("ServerItem register server error : ", err)
		return false
	}
	return true
}

func (s *ServerItem) Deregister() {
	s.Client.Agent().ServiceDeregister(s.SvrData.ID)
}

func (s *ServerItem) GetSvr(svrName string, tag string) []*consulapi.ServiceEntry {
	defer func() { //必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			log.Logger.Errorln(err) //这里的err其实就是panic传入的内容
		}
	}()

	serviceEntry, _, _ := s.Client.Health().Service(svrName, tag, true, &consulapi.QueryOptions{})

	//fmt.Println(serviceEntry[0].Service.Address)
	// for k, v := range serviceEntry {
	// 	//fmt.Println(k, "  ", v.Service.Address)
	// 	//fmt.Println(k, "  ", v.Service.Port)
	// 	fmt.Printf(" %d:%+v \n", k, v.Service)
	// }

	return serviceEntry
}

// func CreateRegistration() *consulapi.AgentServiceRegistration {
// 	return new(consulapi.AgentServiceRegistration)
// }
// func CreateCheck(address string, port int) *consulapi.AgentServiceCheck {
// 	return &consulapi.AgentServiceCheck{
// 		HTTP:                           fmt.Sprintf("http://%s:%d/consul%s", address, port, "/check"),
// 		Timeout:                        "5s",  //设置超时 5s。
// 		Interval:                       "5s",  //设置间隔 5s。
// 		DeregisterCriticalServiceAfter: "30s", //check失败后30秒删除本服务
// 	}
// }

// //服务注册  serverid格式"tid-sid"  tag例如：loginsvr、gatesvr
// func Register(consulAddr string, r *consulapi.AgentServiceRegistration) {
// 	config := consulapi.DefaultConfig()
// 	config.Address = consulAddr
// 	client, err := consulapi.NewClient(config)
// 	if err != nil {
// 		log.Logger.Fatal("consul client error : ", err)
// 	}
// 	err = client.Agent().ServiceRegister(r)
// 	//defer client.Agent().ServiceDeregister(serverid)

// 	if err != nil {
// 		log.Logger.Fatal("svrdiscovery register server error : ", err)
// 	}

// 	// http.HandleFunc("/check", consulCheck)
// 	// go http.ListenAndServe(fmt.Sprintf(":%d", checkPort), nil)
// }

//服务注销
// func Deregister(svrId string) {
// 	config := consulapi.DefaultConfig()
// 	client, err := consulapi.NewClient(config)
// 	if err != nil {
// 		log.Logger.Errorln("consul Deregister client error : ", err)
// 	}
// 	client.Agent().ServiceDeregister(svrId)
// }

// //健康检查
// func consulCheck(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintln(w, "consulCheck")
// }

//服务发现的健康检查接口
func CheckHealth(c *gin.Context) {
	c.String(http.StatusOK, "consulCheck")
}

//
// func GetSvr(svrName string, tag string) []*consulapi.ServiceEntry {
// 	defer func() { //必须要先声明defer，否则不能捕获到panic异常
// 		if err := recover(); err != nil {
// 			log.Logger.Errorln(err) //这里的err其实就是panic传入的内容
// 		}
// 	}()

// 	client, err := consulapi.NewClient(consulapi.DefaultConfig())
// 	if err != nil {
// 		log.Logger.Errorln("consul client error : ", err)
// 	}
// 	serviceEntry, _, _ := client.Health().Service(svrName, tag, true, &consulapi.QueryOptions{})

// 	//fmt.Println(serviceEntry[0].Service.Address)
// 	// for k, v := range serviceEntry {
// 	// 	//fmt.Println(k, "  ", v.Service.Address)
// 	// 	//fmt.Println(k, "  ", v.Service.Port)
// 	// 	fmt.Printf(" %d:%+v \n", k, v.Service)
// 	// }

// 	return serviceEntry
// }
