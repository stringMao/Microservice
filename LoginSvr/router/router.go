package router

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"

	"Common/log"
	"LoginSvr/global"
	"LoginSvr/logic/realname"
	"LoginSvr/logic/signin"
	"LoginSvr/logic/signup"
	"LoginSvr/middle"
	"LoginSvr/sdk/wechat"
)

var router = gin.Default()

//Init ..
func Init() {

	//gin.SetMode(gin.ReleaseMode)
	//路由分组
	//大厅服务器路由
	svrRG := router.Group(strings.ToLower("/mloginsvr/hall"))
	svrRG.Use(middle.Respone(), middle.AuthCheckSign(global.HallSignKey), middle.Authentication()) //中间件设置
	svrRGRouter(svrRG)

	//客户端路由
	clientRG := router.Group(strings.ToLower("/mloginsvr/client"))
	clientRG.Use(middle.Respone(), middle.AuthCheckSign(global.ClientSignKey)) //中间件设置
	clientRGRouter(clientRG)

	//管理员路由

}

//Start webapi启动
func Start(port int) {
	err := router.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Logger.Fatalln("router start is err:", err)
	}

	//router.RunTLS(conf.JSONConf.Port, crtPath, keyPath)
}

//svrRGRouter 大厅服务器访问接口(需要身份认证)
func svrRGRouter(group *gin.RouterGroup) {
	//修改昵称
	group.POST(strings.ToLower("/modifynickname"), signup.ModifyNickname)
	//实名认证
	group.POST(strings.ToLower("/realnameverify"), realname.RealNameVerify)
}

//clientRGRouter 客户端请求接口
func clientRGRouter(group *gin.RouterGroup) {
	//客户端账号登入
	group.POST(strings.ToLower("/signin"), signin.AccountLogin)
	//账号注册
	group.POST(strings.ToLower("/signup"), signup.RegisterAccount)
	//请求短信验证码
	group.POST(strings.ToLower("/applysms"), signup.ApplySMSVerificationCode)
	//忘记密码重置密码
	group.POST(strings.ToLower("/resetpasswd"), signup.LostPasswd)
	//微信登入-code
	group.POST(strings.ToLower("/wechat/codelogin"), wechat.GetAccessTokenByCode)
	//微信登入-access_token
	group.POST(strings.ToLower("/wechat/accesstokenlogin"), wechat.LoginByAccessToken)

}
