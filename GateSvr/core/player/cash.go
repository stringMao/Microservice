package player

import (
	"GateSvr/dbmanager"

	"GateSvr/orm"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

const (
	mark_score = iota
	mark_gold
	mark_max
)

const rediskey_ttl int = 60 * 60 * 24 * 30 //30天
//玩家货币
type PlayerCash struct {
	Score     int64        //积分
	Gold      int64        //金币
	redisMark map[int]bool //内存变更标记
	dbMark    map[int]bool
	redisKey  string //redis里的缓存key
}

//
func newPlayerCash(userid uint64) *PlayerCash {
	return &PlayerCash{
		Score:     0,
		Gold:      0,
		redisMark: make(map[int]bool, mark_max),
		dbMark:    make(map[int]bool, mark_max),
		redisKey:  fmt.Sprintf("PlayerCash:%d", userid),
	}
}

//积分内存变更
func (p *PlayerCash) UpdateScore(score int) {
	p.Score += int64(score)
	p.redisMark[mark_score] = true //标记内存的变更
	p.dbMark[mark_score] = true
}

//加载redis里的数据
func (p *PlayerCash) LoadRedisData(updaTTL bool) bool {
	redisCon := dbmanager.GetRedis()
	defer redisCon.Close()
	if b, err := redis.Bool(redisCon.Do("EXISTS", p.redisKey)); b && err == nil {
		m, _ := redis.Int64Map(redisCon.Do("HGETALL", p.redisKey))

		p.Score = m["Score"]
		p.Gold = m["Gold"]

		if updaTTL {
			redisCon.Do("EXPIRE", p.redisKey, rediskey_ttl)
		}
		return true
	}
	return false
}

//设置redis缓存
func (p *PlayerCash) SetRedisData() {
	redisCon := dbmanager.GetRedis()
	defer redisCon.Close()

	redisCon.Do("HMSET", p.redisKey, "Score", p.Score, "Gold", p.Gold)
	redisCon.Do("EXPIRE", p.redisKey, rediskey_ttl)
}

//加载DB数据
func (p *PlayerCash) LoadDBData(userid uint64) bool {
	var pc orm.PlayerCash
	if !pc.Get(userid) {
		return false
	}
	p.Score = pc.Score
	p.Gold = pc.Gold
	return true
}

func (p *PlayerCash) Save(userid uint64) {
	redisCon := dbmanager.GetRedis()
	defer redisCon.Close()

	//redis缓存同步
	for i := mark_score; i < mark_max; i++ {
		if v, ok := p.redisMark[i]; ok && v { //内存和redis缓存不一致
			switch i {
			case mark_score:
				redisCon.Do("HSET", p.redisKey, "Score", p.Score)
			case mark_gold:
				redisCon.Do("HSET", p.redisKey, "Gold", p.Gold)
			}
			p.redisMark[i] = false
		}
	}
	//db同步
	var pc orm.PlayerCash
	pc.Userid = userid
	var columns []string = make([]string, 0)
	for i := mark_score; i < mark_max; i++ {
		if v, ok := p.dbMark[i]; ok && v {
			switch i {
			case mark_score:
				pc.Score = p.Score
				columns = append(columns, "score") //db里的字段名
			case mark_gold:
				pc.Gold = p.Gold
				columns = append(columns, "gold")
			}
			p.dbMark[i] = false
		}
	}
	if len(columns) > 0 {
		pc.Update(columns...)
	}

}
