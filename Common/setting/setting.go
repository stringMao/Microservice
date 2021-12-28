package setting

import (
	"Common/constant"
	"Common/log"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/Unknwon/goconfig"
)

//BaseConfig 基础配置，必须要读取到的
type BaseConfig struct {
	TID            int
	SID            int    `tag:"server" key:"sid"`
	LogLv          string `tag:"log" key:"level"` //日志等级
	WebManagerPort int    `tag:"webmanager" key:"port"`
	ConsulAddr     string `tag:"webmanager" key:"consuladdr"`
}

//AppCfg 系统配置的全局变量
var appCfg *goconfig.ConfigFile

func LoadAppConfig(v interface{}) {
	var err error
	//dirPath, err := filepath.Abs(filepath.Dir(os.Args[0])) //Getwd()
	dirPath, err := os.Getwd()
	if err != nil {
		log.Logger.Fatal("[app.ini]路径未找到 err:", err)
	}
	confPath, err := filepath.Abs(dirPath + "/app.ini")
	if err != nil {
		log.Logger.Fatal("[app.ini]文件未找到：", err)
	}
	appCfg, err = goconfig.LoadConfigFile(confPath)
	if err != nil {
		log.Logger.Fatal("app.ini read err:", err)
	}
	val := reflect.ValueOf(v).Elem()
	Parsing(val)
}

func Parsing(val reflect.Value) {
	//typ := reflect.TypeOf(v)
	//val := reflect.ValueOf(v).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		if field.Tag.Get("base") == "true" {
			Parsing(val.FieldByName(field.Name))
		}

		tag := field.Tag.Get("tag")
		key := field.Tag.Get("key")
		if tag != "" && key != "" {
			switch field.Type.Name() {
			case "int":
				value, err := appCfg.Int(tag, key)
				if err != nil {
					log.Logger.Fatalf("read app.ini of %s-%s is err:%v", tag, key, err)
					return
				}
				val.FieldByName(field.Name).Set(reflect.ValueOf(value))
			case "string":
				value, err := appCfg.GetValue(tag, key)
				if err != nil {
					log.Logger.Fatalf("read app.ini of %s-%s is err:%v", tag, key, err)
					return
				}
				val.FieldByName(field.Name).Set(reflect.ValueOf(value))
			case "bool":
				value, err := appCfg.Bool(tag, key)
				if err != nil {
					log.Logger.Fatalf("read app.ini of %s-%s is err:%v", tag, key, err)
					return
				}
				val.FieldByName(field.Name).Set(reflect.ValueOf(value))
			}

		}
	}
}

type ServerBase interface {
	GetServerID() uint64
	GetServerIDStr() string
	GetServerName() string
	GetServerTag() string
}

func (s *BaseConfig) GetServerID() uint64 {
	return uint64(s.SID)<<32 + uint64(s.TID)
}
func (s *BaseConfig) GetServerIDStr() string {
	return fmt.Sprintf("TID:%d_SID:%d", s.TID, s.SID)
}

func (s *BaseConfig) GetServerName() string {
	return GetServerName(s.TID) + s.GetServerIDStr()
}

func (s *BaseConfig) GetServerTag() string {
	return GetServerTag(s.TID)
}

func GetServerName(tid int) string {
	switch tid {
	case constant.TID_LoginSvr:
		return "登入服"
	case constant.TID_GateSvr:
		return "网关服"
	case constant.TID_HallSvr:
		return "大厅服"
	default:
		return "未命名"
	}
}

func GetServerTag(tid int) string {
	switch tid {
	case constant.TID_LoginSvr:
		return fmt.Sprintf("登入服")
	case constant.TID_GateSvr:
		return fmt.Sprintf("网关服")
	case constant.TID_HallSvr:
		return fmt.Sprintf("大厅服")
	default:
		return fmt.Sprintf("未命名")
	}
}
