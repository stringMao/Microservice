package config

import (
	"Common/constant"
	"Common/log"
	"Common/setting"
)

//AppCfg 系统配置的全局变量
//var AppCfg *goconfig.ConfigFile

type ServerConfig struct {
	setting.BaseConfig `base:"true"`
	Port       int                `tag:"business" key:"port" binding:"required"`
	DBHost     string             `tag:"mysql" key:"host" binding:"required"`
	DBPort     int                `tag:"mysql" key:"port" binding:"required"`
	DBUserName string             `tag:"mysql" key:"username" binding:"required"`
	DBPwd      string             `tag:"mysql" key:"password" binding:"required"`
	DBName     string             `tag:"mysql" key:"dbname" binding:"required"`
	DBMaxOpen  int                `tag:"mysql" key:"maxopenconns" binding:"required"`
	DBMaxIdle  int                `tag:"mysql" key:"maxidleconns" binding:"required"`

	DBHost_Player     string `tag:"mysql-player" key:"host" binding:"required"`
	DBPort_Player     int    `tag:"mysql-player" key:"port" binding:"required"`
	DBUserName_Player string `tag:"mysql-player" key:"username" binding:"required"`
	DBPwd_Player      string `tag:"mysql-player" key:"password" binding:"required"`
	DBName_Player     string `tag:"mysql-player" key:"dbname" binding:"required"`
	DBMaxOpen_Player  int    `tag:"mysql-player" key:"maxopenconns" binding:"required"`
	DBMaxIdle_Player  int    `tag:"mysql-player" key:"maxidleconns" binding:"required"`

	RedisHost     string `tag:"redis" key:"host" binding:"required"`
	RedisPort     int    `tag:"redis" key:"port" binding:"required"`
	RedisUserName string `tag:"redis" key:"username" binding:"required"`
	RedisPwd      string `tag:"redis" key:"password" binding:"required"`
	RedisNum      int    `tag:"redis" key:"database" binding:"required"`
	RedisMaxOpen  int    `tag:"redis" key:"maxopenconns" binding:"required"`
	RedisMaxIdle  int    `tag:"redis" key:"maxidleconns" binding:"required"`

	WordsPath string `tag:"business" key:"wordsfile"` //敏感词文件
}

var App *ServerConfig

//Init 系统配置读取
func init() {
	App = new(ServerConfig)
	setting.LoadAppConfig(App)
	App.TID = constant.TID_LoginSvr
	App.ServerID = constant.GetServerID(App.TID, App.SID)

	//日志等级设置
	log.Reset(App.LogLv,App.LogWithFunc,App.LogWithFile)
	//敏感词加载
	LoadSensitiveWords()
}
