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
	ClientPort         int `tag:"business" key:"clientPort" binding:"required"`
	ServerPort         int `tag:"business" key:"serverPort" binding:"required"`

	RedisHost     string `tag:"redis" key:"host" binding:"required"`
	RedisPort     int    `tag:"redis" key:"port" binding:"required"`
	RedisUserName string `tag:"redis" key:"username" binding:"required"`
	RedisPwd      string `tag:"redis" key:"password" binding:"required"`
	RedisNum      int    `tag:"redis" key:"database" binding:"required"`
	RedisMaxOpen  int    `tag:"redis" key:"maxopenconns" binding:"required"`
	RedisMaxIdle  int    `tag:"redis" key:"maxidleconns" binding:"required"`

	DBHost     string `tag:"mysql" key:"host" binding:"required"`
	DBPort     int    `tag:"mysql" key:"port" binding:"required"`
	DBUserName string `tag:"mysql" key:"username" binding:"required"`
	DBPwd      string `tag:"mysql" key:"password" binding:"required"`
	DBName     string `tag:"mysql" key:"dbname" binding:"required"`
	DBMaxOpen  int    `tag:"mysql" key:"maxopenconns" binding:"required"`
	DBMaxIdle  int    `tag:"mysql" key:"maxidleconns" binding:"required"`
}

var App *ServerConfig

func init() {
	App = new(ServerConfig)

	setting.LoadAppConfig(App)
	App.TID = constant.TID_GateSvr
	App.ServerID = constant.GetServerID(App.TID, App.SID)
	//日志等级设置
	log.Reset(App.LogLv,App.LogWithFunc,App.LogWithFile)
}
