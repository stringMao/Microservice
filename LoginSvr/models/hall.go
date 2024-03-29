package models

import (
	"Common/log"
	"LoginSvr/dbmanager"
)

//Hall 大厅服务器信息表
type Hall struct {
	Serverid   int `xorm:"pk notnull"`
	Servername string
	Address    string
	Channel    int
	Status     int
}

//TableName ..
func (*Hall) TableName() string {
	return "halls"
}

//LoadServerInfo 加载大厅服务器组信息
func LoadServerInfo() []Hall {
	var halls []Hall
	err := dbmanager.Get_LoginSvr().Where("status=?", 0).Find(&halls)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("Hall [LoadServerInfo] is err")
	}
	return halls
}
