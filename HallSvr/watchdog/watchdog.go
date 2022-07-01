package watchdog



func StartWork() {

	//开启网关服连接和检查协程
	go ConnectGateSvrs()
}

