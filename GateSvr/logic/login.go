package logic

import (
	"GateSvr/dbmanager"
	"fmt"
	"strings"

	"github.com/gomodule/redigo/redis"
)

func UserLogin(userid uint64, token string) byte {
	redisCon := dbmanager.GetRedis()
	defer redisCon.Close()

	redistoken, err := redis.String(redisCon.Do("GET", fmt.Sprintf("token:%d", userid)))
	if err != nil {
		return 1
	}
	if strings.EqualFold(token, redistoken) {
		return 0
	}

	return 2
}

//服务器连接的身份验证
func ServerLoginAuthentication(serverid uint64, password string) bool {

	return true
}
