package config

import (
	"Common/constant"
	"Common/log"
	"Common/setting"
)

type ServerConfig struct {
	setting.BaseConfig `base:"true"`
	Port               int `tag:"business" key:"port"`
}

var App *ServerConfig

func init() {
	App = new(ServerConfig)
	App.TID = constant.TID_HallSvr
	setting.LoadAppConfig(App)
	//日志等级设置
	log.Setup(App.LogLv)
}
