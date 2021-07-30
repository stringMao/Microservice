package dbmanager

import (
	"Common/db"
	"Common/log"
	"LoginSvr/config"

	"github.com/go-xorm/xorm"
	"github.com/gomodule/redigo/redis"
)

var db_LoginSvr = &db.MysqlData{}

func ConnectDB_LoginSvr() {
	db_LoginSvr.Host = config.App.DBHost
	db_LoginSvr.Port = config.App.DBPort
	db_LoginSvr.UserName = config.App.DBUserName
	db_LoginSvr.PassWord = config.App.DBPwd
	db_LoginSvr.DBName = config.App.DBName
	db_LoginSvr.MaxIdleConns = config.App.DBMaxIdle
	db_LoginSvr.MaxOpenConns = config.App.DBMaxOpen
	if !db_LoginSvr.Open() {
		log.Logger.Fatalln("mysql init fail")
	}
	log.Logger.Info("mysql init success")
}
func Get_LoginSvr() *xorm.Engine {
	return db_LoginSvr.GetMysql()
}

var redisData = &db.RedisData{}

func ConnectRedis() {
	redisData.Host = config.App.RedisHost
	redisData.Port = config.App.RedisPort
	redisData.PassWord = config.App.RedisPwd
	redisData.DataBase = config.App.RedisNum
	redisData.MaxOpenConns = config.App.RedisMaxOpen
	redisData.MaxIdleConns = config.App.RedisMaxIdle

	if !redisData.Open() {
		log.Logger.Fatalln("redis init fail")
	}
	log.Logger.Info("redis init success")
}

func GetRedis() redis.Conn {
	return redisData.GetRedis()
}
