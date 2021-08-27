package orm

import (
	"Common/log"
	"GateSvr/dbmanager"
)

type PlayerCash struct {
	Userid uint64
	Score  int64
	Gold   int64
}

//TableName ..
func (*PlayerCash) TableName() string {
	return "playercash"
}

func (p *PlayerCash) Get(userid uint64) bool {
	has, err := dbmanager.GetDB_Player().Where("userid=?", userid).Get(p)
	if !has { //用户数据不存在，说明是第一次注册的
		//初始化数据
		p.Userid = userid
		p.Gold = 0
		p.Score = 0
		affected, err2 := dbmanager.GetDB_Player().Insert(p)
		if err2 != nil || affected != 1 {
			log.WithFields(log.Fields{
				"affected": affected,
				"err":      err2,
				"userid":   userid,
			}).Error("PlayerCash [Get] insert is err")
			return false
		}
		return true
	}
	if err != nil {
		log.WithFields(log.Fields{
			"has":    has,
			"err":    err,
			"userid": userid,
		}).Error("PlayerCash [Get] is err")
		return false
	}
	return true
}

func (p *PlayerCash) Update(columns ...string) bool {
	_, err := dbmanager.GetDB_Player().Where("userid=?", p.Userid).Cols(columns...).Update(p)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("PlayerCash.go [Update] is err")
		return false
	}
	return true
}
