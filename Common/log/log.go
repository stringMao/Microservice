package log

import (
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

//Logger 全局日志对象
var Logger = logrus.New()

func init() {
	//Logger.SetFormatter(&logrus.JSONFormatter{})
	Logger.AddHook(newLfsHook("logs/"))
}

//Setup ..
func Setup(lv string) {
	//更新日志设置
	logrusLogLevel, err := logrus.ParseLevel(lv)
	if err != nil {
		Logger.Fatalln("app.ini of log-lvevl is err:", err)
	}
	Logger.SetLevel(logrusLogLevel) //设置等级
}

// 设置日志文件切割及软链接
func newLfsHook(filepath string) logrus.Hook {
	var err error
	//===debuglog======================================
	logpath := filepath + "debug/"
	writerDebug, err := rotatelogs.New(
		logpath+"%Y%m%d%H%M",
		rotatelogs.WithLinkName(logpath),       // 生成软链，指向最新日志文件
		rotatelogs.WithRotationTime(time.Hour), //设置日志分割的时间，这里设置为一小时分割一次
		//WithMaxAge和WithRotationCount二者只能设置一个，
		rotatelogs.WithMaxAge(time.Hour*24*5), // 文件最大保存时间
	)
	if err != nil {
		logrus.Errorf("writerDebug logger error. %+v", errors.WithStack(err))
	}

	//====infolog===================================
	logpath = filepath + "info/"
	writerInfo, err := rotatelogs.New(
		logpath+"%Y%m%d%H%M",
		rotatelogs.WithLinkName(logpath),
		rotatelogs.WithRotationTime(time.Hour),
		rotatelogs.WithMaxAge(time.Hour*24*30),
	)
	if err != nil {
		logrus.Errorf("writerInfo logger error. %+v", errors.WithStack(err))
	}
	//===warn log===============================================
	logpath = filepath + "warn/"
	writerWarn, err := rotatelogs.New(
		logpath+"%Y%m%d%H%M",
		rotatelogs.WithLinkName(logpath),
		rotatelogs.WithRotationTime(time.Hour),
		rotatelogs.WithMaxAge(time.Hour*24*30),
	)
	if err != nil {
		logrus.Errorf("writerWarn logger error. %+v", errors.WithStack(err))
	}

	//====Errlog===================================
	logpath = filepath + "error/"
	writerErr, err := rotatelogs.New(
		logpath+"%Y%m%d%H%M",
		rotatelogs.WithLinkName(logpath),
		rotatelogs.WithRotationTime(time.Hour*24),
		rotatelogs.WithMaxAge(time.Hour*24*30),
	)
	if err != nil {
		logrus.Errorf("writerErr logger error. %+v", errors.WithStack(err))
	}
	//==Fatal log=========================================
	logpath = filepath + "fatal/"
	writerFatal, err := rotatelogs.New(
		logpath+"%Y%m%d%H%M",
		rotatelogs.WithLinkName(logpath),
		rotatelogs.WithRotationTime(time.Hour*24),
		rotatelogs.WithMaxAge(time.Hour*24*30),
	)
	if err != nil {
		logrus.Errorf("writerFatal logger error. %+v", errors.WithStack(err))
	}

	//===Panic log===============================================
	logpath = filepath + "panic/"
	writerPanic, err := rotatelogs.New(
		logpath+"%Y%m%d%H%M",
		rotatelogs.WithLinkName(logpath),
		rotatelogs.WithRotationTime(time.Hour*24),
		rotatelogs.WithMaxAge(time.Hour*24*30),
	)
	if err != nil {
		logrus.Errorf("writerPanic logger error. %+v", errors.WithStack(err))
	}

	/*
		logrus 拥有六种日志级别：debug、info、warn、error、fatal 和 panic，
		logrus.Debug(“Useful debugging information.”)
		logrus.Info(“Something noteworthy happened!”)
		logrus.Warn(“You should probably take a look at this.”)
		logrus.Error(“Something failed but I'm not quitting.”)
		logrus.Fatal(“Bye.”) //log之后会调用os.Exit(1)
		logrus.Panic(“I'm bailing.”) //log之后会panic()
	*/
	//设置默认等级
	logrusLogLevel, _ := logrus.ParseLevel("debug")
	Logger.SetLevel(logrusLogLevel) //设置等级

	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writerDebug,
		logrus.InfoLevel:  writerInfo,
		logrus.WarnLevel:  writerWarn,
		logrus.ErrorLevel: writerErr,
		logrus.FatalLevel: writerFatal,
		logrus.PanicLevel: writerPanic,
	}, &logrus.TextFormatter{DisableColors: true})

	return lfsHook
}

//Fields ..
type Fields map[string]interface{}

//WithFields 重写此函数，便于使用
func WithFields(fields Fields) *logrus.Entry {
	return Logger.WithFields(logrus.Fields(fields))
}

//WithField 重写此函数，便于使用
func WithField(key string, value interface{}) *logrus.Entry {
	return Logger.WithField(key, value)
}

func Debug(args ...interface{}) {
	Logger.Debug(args...)
}

func Debugln(args ...interface{}) {
	Logger.Debugln(args...)
}

func Debugf(format string, args ...interface{}) {
	Logger.Debugf(format, args...)
}

func Info(args ...interface{}) {
	Logger.Info(args...)
}
func Infoln(args ...interface{}) {
	Logger.Infoln(args...)
}
func Infof(format string, args ...interface{}) {
	Logger.Infof(format, args...)
}

func Warn(args ...interface{}) {
	Logger.Warn(args...)
}
func Warnln(args ...interface{}) {
	Logger.Warnln(args...)
}
func Warnf(format string, args ...interface{}) {
	Logger.Warnf(format, args...)
}

func Error(args ...interface{}) {
	Logger.Error(args...)
}
func Errorln(args ...interface{}) {
	Logger.Errorln(args...)
}
func Errorf(format string, args ...interface{}) {
	Logger.Errorf(format, args...)
}
func Fatal(args ...interface{}) {
	Logger.Fatal(args...)
}
func Fatalln(args ...interface{}) {
	Logger.Fatalln(args...)
}
func Fatalf(format string, args ...interface{}) {
	Logger.Fatalf(format, args...)
}

func Panic(args ...interface{}) {
	Logger.Panic(args...)
}
func Panicln(args ...interface{}) {
	Logger.Panicln(args...)
}
func Panicf(format string, args ...interface{}) {
	Logger.Panicf(format, args...)
}

//打印异常 及堆栈信息
func PrintPanicStack(err interface{}, stack string) {
	Logger.Error("catch Exception:", err)
	Logger.Error("Stack Info start ============")
	Logger.Error(stack)
	Logger.Error("Stack Info end ==============")
}
