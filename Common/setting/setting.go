package setting

import (
	"Common/log"
	"os"
	"path/filepath"
	"reflect"

	"github.com/Unknwon/goconfig"
)

//BaseConfig 基础配置，必须要读取到的
type BaseConfig struct {
	TID            int
	SID            int `tag:"server" key:"sid"`
	ServerID       uint64
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
