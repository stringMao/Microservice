package models

import (
	"Common/log"
	"LoginSvr/dbmanager"
)

//UserRealInfo 用户实名信息
type UserRealInfo struct {
	Userid   int64 `xorm:"pk autoincr notnull"`
	Name     string
	Identity string
	Gender   int
	Addr     string
}

//TableName ..
func (*UserRealInfo) TableName() string {
	return "userRealInfo"
}

//GetByUserid 通过id查询
func (a *UserRealInfo) GetByUserid(id int64) bool {
	has, err := dbmanager.Get_LoginSvr().Where("userid=?", id).Get(a)
	if err != nil {
		log.WithFields(log.Fields{
			"has":    has,
			"err":    err,
			"userid": id,
		}).Error("UserRealInfo [GetByUserid] is err")
		return false
	}

	return has
}

//UpdateOrInsert 更新或者插入实名信息
func (a *UserRealInfo) UpdateOrInsert() bool {
	has, err := dbmanager.Get_LoginSvr().Exist(&UserRealInfo{Userid: a.Userid})
	if err != nil {
		log.WithFields(log.Fields{
			"has": has,
			"err": err,
			"id":  a.Userid,
		}).Error("UserRealInfo.go [UpdateOrInsert]-1 is err")
		return false
	}
	if has {
		//存在则更新
		affected, err := dbmanager.Get_LoginSvr().Update(a, &UserRealInfo{Userid: a.Userid})
		if err != nil {
			log.WithFields(log.Fields{
				"affected": affected,
				"err":      err,
				"id":       a.Userid,
			}).Error("UserRealInfo.go [UpdateOrInsert]-2 is err")
			return false
		}
	} else {
		//不存在，则插入
		affected, err := dbmanager.Get_LoginSvr().Insert(a)
		if err != nil || affected != 1 {
			log.WithFields(log.Fields{
				"affected": affected,
				"err":      err,
				"id":       a.Userid,
			}).Error("UserRealInfo.go [UpdateOrInsert]-3 is err")
			return false
		}
	}
	return true
}
