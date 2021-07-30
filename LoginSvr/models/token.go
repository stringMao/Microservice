package models

import (
	"Common/log"
	"LoginSvr/dbmanager"
	"time"
)

//Token ..
type Token struct {
	Userid     int64
	Token      string
	Updatetime time.Time `xorm:"updated"`
}

//TableName ..
func (*Token) TableName() string {
	return "tokens"
}

//GetByUserid ..
func (t *Token) GetByUserid(id int64) string {
	has, err := dbmanager.Get_LoginSvr().Where("userid=?", id).Get(t)
	if err != nil {
		log.WithFields(log.Fields{
			"has": has,
			"err": err,
			"id":  id,
		}).Error("Token.go [GetByUserid] is err")
		return ""
	}
	return t.Token
}

//InsertOrUpdate 更新或插入token
func (t *Token) InsertOrUpdate() bool {
	has, err := dbmanager.Get_LoginSvr().Where("userid = ?", t.Userid).Exist(&Token{})
	if err != nil {
		log.WithFields(log.Fields{
			"has": has,
			"err": err,
			"id":  t.Userid,
		}).Error("Token.go [InsertOrUpdate]-1 is err")
		return false
	}
	if has {
		//存在则更新
		affected, err := dbmanager.Get_LoginSvr().Where("userid = ?", t.Userid).Cols("token").Update(t)
		if err != nil || affected != 1 {
			log.WithFields(log.Fields{
				"affected": affected,
				"err":      err,
				"id":       t.Userid,
			}).Error("Token.go [InsertOrUpdate]-2 is err")
			return false
		}
	} else {
		//不存在，则插入
		affected, err := dbmanager.Get_LoginSvr().Insert(t)
		if err != nil || affected != 1 {
			log.WithFields(log.Fields{
				"affected": affected,
				"err":      err,
				"id":       t.Userid,
			}).Error("Token.go [InsertOrUpdate]-3 is err")
			return false
		}
	}

	return true
}
