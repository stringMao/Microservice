package setting

import (
	"Common/log"
	"Common/util"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/Unknwon/goconfig"
)

//BaseConfig 基础配置，必须要读取到的
type BaseConfig struct {
	TID            int
	SID            int `tag:"server" key:"sid" binding:"required"`
	ServerID       uint64
	LogLv          string `tag:"log" key:"level"` //日志等级
	LogWithFunc    bool `tag:"log" key:"withfunc"` //日志打印函数名
	LogWithFile    bool `tag:"log" key:"withfile"` //日志打印文件名及行号
	WebManagerIP   string `tag:"webmanager" key:"ip" binding:"required"`
	WebManagerPort int    `tag:"webmanager" key:"port" binding:"required"`
	ConsulAddr     string `tag:"webmanager" key:"consuladdr" binding:"required"`
}

//AppCfg 系统配置的全局变量
var appCfg *goconfig.ConfigFile

func LoadAppConfig(v interface{}) {
	var err error
	//dirPath, err := filepath.Abs(filepath.Dir(os.Args[0])) //Getwd()
	dirPath, err :=util.GetCurrentPath()
	//dirPath, err := os.Getwd()
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
			//Parsing(val.FieldByName(field.Name))
			Parsing(val.Field(i))
		}


		tag := field.Tag.Get("tag")
		key := field.Tag.Get("key")
		bind:=field.Tag.Get("binding")
		secBind:=strings.Split(bind,",")
		if tag != "" && key != "" {
			switch field.Type.Name() {
			case "int":
				value, err := appCfg.Int(tag, key)
				if err != nil && isHaveString(secBind,"required"){
					log.Logger.Fatalf("read app.ini of %s-%s is err:%v", tag, key, err)
					return
				}
				val.FieldByName(field.Name).Set(reflect.ValueOf(value))
			case "string":
				value, err := appCfg.GetValue(tag, key)
				if err != nil && isHaveString(secBind,"required") {
					log.Logger.Fatalf("read app.ini of %s-%s is err:%v", tag, key, err)
					return
				}
				val.FieldByName(field.Name).Set(reflect.ValueOf(value))
			case "bool":
				value, err := appCfg.Bool(tag, key)
				if err != nil && isHaveString(secBind,"required") {
					log.Logger.Fatalf("read app.ini of %s-%s is err:%v", tag, key, err)
					return
				}
				val.FieldByName(field.Name).Set(reflect.ValueOf(value))
			}

		}
	}
}

func isHaveString(src []string,s string)bool{
	for _, v := range src {
		if strings.Compare(v,s)==0 {
			return true
		}
	}
	return false
}
