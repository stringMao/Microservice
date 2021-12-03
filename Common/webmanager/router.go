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

func CreateRouterAgent() *RouterAgent {
	//gin.New()
	gin.SetMode(gin.ReleaseMode)
	return &RouterAgent{
		Router:       gin.Default(),
		Group_consul: nil,
		Group_gm:     nil,
	}
}

func CreateRouterAgent2() *RouterAgent {
	gin.SetMode(gin.ReleaseMode)
	r := &RouterAgent{
		Router:       gin.New(),
		Group_consul: nil,
		Group_gm:     nil,
	}
	return r
}

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

func (agent *RouterAgent) Start(webManagerPort int) {
	err := agent.Router.Run(fmt.Sprintf(":%d", webManagerPort))
	if err != nil {
		log.Logger.Fatalln("webManager start is fail, err:", err)
	}
}
