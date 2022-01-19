package watchdog

import (
	"Common/log"
	"Common/try"
	"net"
	"strconv"
	"time"
)

func StartWork() {
	//开启消息处理协程
	go HandleSvrMsg()
	go HandlClientMsg()
	go HandleSelfMsg()

	//开启网关服连接和检查协程
	go ConnectGateSvrs()
}

//Scoket开始监听
func StartTCPListen(port int) {
	netListen, err := net.Listen("tcp", "localhost:"+strconv.Itoa(port))
	if err != nil {
		log.Logger.Fatal(err)
		return
	}
	defer netListen.Close()

	for {
		conn, err := netListen.Accept()
		if err != nil {
			log.Logger.Error(err)
			continue
		}
		//监听客户端连接
		go handleConnection(conn)

	}
}

func handleConnection(conn net.Conn) {
	defer try.Catch()
	defer conn.Close()
	//conn.SetReadDeadline(time.Now().Add(time.Second * 30))
	//连接之后的第一条消息，必须是验证身份，并且获得tid，sid
	buffer := make([]byte, 2048) //建立一个slice
	_, err := conn.Read(buffer)
	if err != nil {
		log.Logger.Error(conn.RemoteAddr().String(), " read server first msg error: ", err)
		return //当远程客户端连接发生错误（断开）后，终止此协程。
	}
	//buf := buffer[:n]

	for {
		//开始通讯
		conn.SetReadDeadline(time.Now().Add(time.Second * 10)) //借此检测心跳包
		n, err := conn.Read(buffer)                            //读取客户端传来的内容
		if err != nil {
			log.Logger.Debug(conn.RemoteAddr().String(), " server connection error: ", err)
			return //当远程客户端连接发生错误（断开）后，终止此协程。
		}
		//特殊包-心跳包过滤  消息结构[uint8]=200
		if n == 1 && buffer[0] == 200 {
			//log.Logger.Debugln("heart")
			continue
		}

	}
}
