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
	Port       int                `tag:"business" key:"port"`
	DBHost     string             `tag:"mysql" key:"host"`
	DBPort     int                `tag:"mysql" key:"port"`
	DBUserName string             `tag:"mysql" key:"username"`
	DBPwd      string             `tag:"mysql" key:"password"`
	DBName     string             `tag:"mysql" key:"dbname"`
	DBMaxOpen  int                `tag:"mysql" key:"maxopenconns"`
	DBMaxIdle  int                `tag:"mysql" key:"maxidleconns"`

	DBHost_Player     string `tag:"mysql-player" key:"host"`
	DBPort_Player     int    `tag:"mysql-player" key:"port"`
	DBUserName_Player string `tag:"mysql-player" key:"username"`
	DBPwd_Player      string `tag:"mysql-player" key:"password"`
	DBName_Player     string `tag:"mysql-player" key:"dbname"`
	DBMaxOpen_Player  int    `tag:"mysql-player" key:"maxopenconns"`
	DBMaxIdle_Player  int    `tag:"mysql-player" key:"maxidleconns"`

	RedisHost     string `tag:"redis" key:"host"`
	RedisPort     int    `tag:"redis" key:"port"`
	RedisUserName string `tag:"redis" key:"username"`
	RedisPwd      string `tag:"redis" key:"password"`
	RedisNum      int    `tag:"redis" key:"database"`
	RedisMaxOpen  int    `tag:"redis" key:"maxopenconns"`
	RedisMaxIdle  int    `tag:"redis" key:"maxidleconns"`

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
	log.Setup(App.LogLv)
	//敏感词加载
	LoadSensitiveWords()
}
