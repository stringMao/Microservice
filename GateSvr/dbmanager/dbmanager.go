package dbmanager

import (
	"Common/db"
	"Common/log"
	"GateSvr/config"

	"github.com/garyburd/redigo/redis"
	"github.com/go-xorm/xorm"
)

type DBAgent struct {
	player db.MysqlData
	redis1 db.RedisData
}

func GetDB_Player() *xorm.Engine {
	return dbmanager.player.GetMysql()
}
func GetRedis() redis.Conn {
	return dbmanager.redis1.GetRedis()
}

var dbmanager DBAgent

func Init() {
	dbmanager.player.Host = config.App.DBHost
	dbmanager.player.Port = config.App.DBPort
	dbmanager.player.UserName = config.App.DBUserName
	dbmanager.player.PassWord = config.App.DBPwd
	dbmanager.player.DBName = config.App.DBName
	dbmanager.player.MaxIdleConns = config.App.DBMaxIdle
	dbmanager.player.MaxOpenConns = config.App.DBMaxOpen
	if !dbmanager.player.Open() {
		log.Logger.Fatalln("mysql [player] init fail")
	}
	log.Logger.Info("mysql [player]  init success")

	dbmanager.redis1.Host = config.App.RedisHost
	dbmanager.redis1.Port = config.App.RedisPort
	dbmanager.redis1.PassWord = config.App.RedisPwd
	dbmanager.redis1.DataBase = config.App.RedisNum
	dbmanager.redis1.MaxOpenConns = config.App.RedisMaxOpen
	dbmanager.redis1.MaxIdleConns = config.App.RedisMaxIdle

	if !dbmanager.redis1.Open() {
		log.Logger.Fatalln("redis [1]init fail")
	}
	log.Logger.Info("redis [1] init success")
}
