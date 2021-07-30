package pub

//一些逻辑业务操作整合接口

import (
	"LoginSvr/dbmanager"
	"LoginSvr/global"
	"LoginSvr/models"
	"fmt"
)

//InsertTokenToRedis 将登入token插入redis
func InsertTokenToRedis(token string, acc models.Account) bool {
	redisCon := dbmanager.GetRedis()
	defer redisCon.Close()

	redisCon.Do("SET", fmt.Sprintf("token:%d", acc.Userid), token, "EX", global.TokenActiveTime)

	var key = fmt.Sprintf("accountdata:%d", acc.Userid)
	redisCon.Do("HSET", key, "nickname", acc.Nickname)
	redisCon.Do("HSET", key, "accounttype", acc.Accounttype)
	var userinfo models.UserRealInfo
	userinfo.GetByUserid(acc.Userid)
	redisCon.Do("HSET", key, "age", global.GetCitizenAge([]byte(userinfo.Identity), true)) //年龄-1表示未实名认证
	redisCon.Do("HSET", key, "gender", userinfo.Gender)
	redisCon.Do("EXPIRE", key, global.TokenActiveTime)
	return true

	// u := global.RedisUserInfo{Gender: 0, CompleteRealname: false, Age: 0}
	// u.Token = token
	// u.Nickname = acc.Nickname
	// u.Accounttype = acc.Accounttype
	// //u.Indentity = acc.Identity
	// //获得实名信息
	// var userinfo models.UserRealInfo
	// if userinfo.GetByUserid(acc.Userid) {
	// 	u.CompleteRealname = true
	// 	u.Gender = userinfo.Gender
	// 	u.Age = global.GetCitizenAge([]byte(userinfo.Identity), false)
	// }

	// var key = fmt.Sprintf("token:%d", acc.Userid)
	// if val, err := json.Marshal(u); err == nil {
	// 	redisCon.Do("SET", key, string(val), "EX", global.TokenActiveTime)
	// 	return true
	// }
	// return false
}
