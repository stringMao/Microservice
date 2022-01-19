package watchdog

func StartWork() {
	//开启消息处理协程
	go HandleSvrMsg()
	go HandlClientMsg()
	go HandleSelfMsg()

	//开启网关服连接和检查协程
	go ConnectGateSvrs()
}
