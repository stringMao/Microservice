package webmanager

//web后台
import (
	"Common/log"
	"fmt"

	"github.com/gin-gonic/gin"
)

type RouterHelper struct {
	Type     int //路由组类型
	Path     string
	Method   string
	Handlers []gin.HandlerFunc
}

//路由组类型
const (
	RouterType_consul = iota //服务发现功能路由
	RouterType_gm            //gm功能路由
)

type RouterAgent struct {
	Router       *gin.Engine
	Group_consul *gin.RouterGroup
	Group_gm     *gin.RouterGroup
}

var G_WebManager *RouterAgent

func init() {
	gin.SetMode(gin.ReleaseMode)
	G_WebManager = new(RouterAgent)
	G_WebManager.Router = gin.Default() //gin.New()然后自定义日志输出
	G_WebManager.Group_consul = nil
	G_WebManager.Group_gm = nil
}

// func CreateRouterAgent() *RouterAgent {
// 	//gin.New()
// 	gin.SetMode(gin.ReleaseMode)
// 	return &RouterAgent{
// 		Router:       gin.Default(),
// 		Group_consul: nil,
// 		Group_gm:     nil,
// 	}
// }

//路由注册
func (agent *RouterAgent) RegisterRouter(r RouterHelper) bool {
	var rg *gin.RouterGroup = nil

	if r.Type == RouterType_consul { //服务发现相关路由
		if agent.Group_consul == nil {
			agent.Group_consul = agent.Router.Group("/consul")
		}
		rg = agent.Group_consul
	} else if r.Type == RouterType_gm { //gm功能相关路由
		if agent.Group_gm == nil {
			agent.Group_gm = agent.Router.Group("/gm")
			//添加身份验证的中间件
		}
		rg = agent.Group_gm
	}
	if rg == nil {
		return false
	}

	switch r.Method {
	case "GET":
		rg.GET(r.Path, r.Handlers...)
	case "POST":
		rg.POST(r.Path, r.Handlers...)
	default:
		return false
	}
	return true
}

func (agent *RouterAgent) Start(webManagerPort int) bool {
	go func() {
		err := agent.Router.Run(fmt.Sprintf("0.0.0.0:%d", webManagerPort))
		if err != nil {
			log.Errorln("webManager Run is err:", err)
		}
	}()
	return true
}

