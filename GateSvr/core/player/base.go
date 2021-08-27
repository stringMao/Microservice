package player

import (
	"Common/try"
	"GateSvr/dbmanager"
	"fmt"
	"strconv"

	"github.com/gomodule/redigo/redis"
)

//玩家的账号基础信息
type PlayerBase struct {
	NickName string //昵称
	Avatar   string //头像
	Gender   int    //性别
	Age      int
	redisKey string //redis里的缓存key
}

func newPlyerBase(userid uint64) *PlayerBase {
	return &PlayerBase{
		NickName: "未知",
		Avatar:   "默认头像url",
		Gender:   0,
		Age:      0,
		redisKey: fmt.Sprintf("PlayerBase:%d", userid),
	}
}

func (p *PlayerBase) LoadRedisData() bool {
	defer try.Catch()
	redisCon := dbmanager.GetRedis()
	defer redisCon.Close()
	if b, err := redis.Bool(redisCon.Do("EXISTS", p.redisKey)); b && err == nil {
		m, _ := redis.StringMap(redisCon.Do("HGETALL", p.redisKey))

		p.NickName = m["nickname"]
		//p.Avatar
		p.Gender, _ = strconv.Atoi(m["gender"])
		p.Age, _ = strconv.Atoi(m["age"])

		return true
	}
	return false
}
