package try

import (
	"Common/log"
	"fmt"
	"path"
	"runtime"
	"runtime/debug"
	"strings"
)

func Catch() {
	if r := recover(); r != nil {
		// log.Logger.Error("catch Exception:", r)
		// log.Logger.Error("Stack Info start ============")
		// log.Logger.Error(string(debug.Stack()))
		// log.Logger.Error("Stack Info end ==============")

		pc, file, lineNo, ok := runtime.Caller(4)
		if ok{
			//打印异常捕捉处
			log.Errorf("异常 func=%s  file=%s",runtime.FuncForPC(pc).Name(),fmt.Sprintf("%s:%d",path.Base(file),lineNo))
		}

		str:=string(debug.Stack())
		str=strings.Replace(str,"\t","",-1)
		sec:=strings.Split(str,"\n")
		log.Error("catch Exception: ", r)
		log.Error("Stack Info start =======================")
		for _,v:=range sec{
			log.Error(v)
		}

		log.Error("Stack Info end =========================")
		//log.PrintPanicStack(r, string(debug.Stack()))
	}
}
