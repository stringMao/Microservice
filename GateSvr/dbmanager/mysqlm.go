package dbmanager

import (
	"Common/db"
	"Common/log"
	"GateSvr/config"

	"github.com/go-xorm/xorm"
)

var db_Player = &db.MysqlData{}

func ConnectDB() {
	db_Player.Host = config.App.DBHost
	db_Player.Port = config.App.DBPort
	db_Player.UserName = config.App.DBUserName
	db_Player.PassWord = config.App.DBPwd
	db_Player.DBName = config.App.DBName
	db_Player.MaxIdleConns = config.App.DBMaxIdle
	db_Player.MaxOpenConns = config.App.DBMaxOpen
	if !db_Player.Open() {
		log.Logger.Fatalln("mysql init fail")
	}

	log.Logger.Info("mysql init success")
}

func GetDB_Player() *xorm.Engine {
	return db_Player.GetMysql()
}
