package db

import (
	"Common/log"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

//MasterDB ..
//var MasterDB *xorm.Engine

type MysqlData struct {
	Host         string
	Port         int
	UserName     string
	PassWord     string
	DBName       string
	MaxIdleConns int //连接池的空闲数大小
	MaxOpenConns int //连接池最大打开连接数
	engine       *xorm.Engine
}

func (db *MysqlData) Open() bool {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", db.UserName, db.PassWord, db.Host, db.Port, db.DBName)
	var err error
	db.engine, err = xorm.NewEngine("mysql", dataSourceName)
	if err != nil {
		log.Logger.Fatalln("mysql connet is err:", err)
		return false
	}
	//engine.SetMapper(names.SameMapper{})//"xorm.io/xorm/names"

	db.engine.SetMaxIdleConns(db.MaxIdleConns) //连接池的空闲数大小
	db.engine.SetMaxOpenConns(db.MaxOpenConns) //最大打开连接数
	return true
}

func (db *MysqlData) GetMysql() *xorm.Engine {
	return db.engine
}
