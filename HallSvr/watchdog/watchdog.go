package watchdog



func StartWork() {
	Init()
	
	//开启网关服连接和检查协程
	go ConnectGateSvrs()
}

func Init(){
	InitHandleFunc()
}