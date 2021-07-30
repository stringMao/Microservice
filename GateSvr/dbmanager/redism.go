package dbmanager

import (
	"Common/db"
	"Common/log"
	"GateSvr/config"

	"github.com/gomodule/redigo/redis"
)

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
