package config

import (
	"Common/constant"
	"Common/log"
	"Common/setting"
)

//AppCfg 系统配置的全局变量
//var AppCfg *goconfig.ConfigFile

type ServerConfig struct {
	// TID            int
	// SID            int    `tag:"server" key:"sid"`
	// LogLv          string `tag:"log" key:"level"` //日志等级
	// WebManagerPort int    `tag:"webmanager" key:"port"`
	//Base       setting.BaseConfig `base:"true"`
	setting.BaseConfig `base:"true"`
	ClientPort         int `tag:"business" key:"clientPort"`
	ServerPort         int `tag:"business" key:"serverPort"`

	RedisHost     string `tag:"redis" key:"host"`
	RedisPort     int    `tag:"redis" key:"port"`
	RedisUserName string `tag:"redis" key:"username"`
	RedisPwd      string `tag:"redis" key:"password"`
	RedisNum      int    `tag:"redis" key:"database"`
	RedisMaxOpen  int    `tag:"redis" key:"maxopenconns"`
	RedisMaxIdle  int    `tag:"redis" key:"maxidleconns"`

	DBHost     string `tag:"mysql" key:"host"`
	DBPort     int    `tag:"mysql" key:"port"`
	DBUserName string `tag:"mysql" key:"username"`
	DBPwd      string `tag:"mysql" key:"password"`
	DBName     string `tag:"mysql" key:"dbname"`
	DBMaxOpen  int    `tag:"mysql" key:"maxopenconns"`
	DBMaxIdle  int    `tag:"mysql" key:"maxidleconns"`
}

var App *ServerConfig

func init() {
	App = new(ServerConfig)
	App.TID = constant.TID_GateSvr
	setting.LoadAppConfig(App)
	//日志等级设置
	log.Setup(App.LogLv)
}
