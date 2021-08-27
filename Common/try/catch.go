package try

import (
	"Common/log"
	"runtime/debug"
)

func Catch() {
	if r := recover(); r != nil {
		// log.Logger.Error("catch Exception:", r)
		// log.Logger.Error("Stack Info start ============")
		// log.Logger.Error(string(debug.Stack()))
		// log.Logger.Error("Stack Info end ==============")

		log.PrintPanicStack(r, string(debug.Stack()))
	}
}
