package watchdog

//服务器连接的消息处理
import (
	"Common/log"
	"Common/msg"
	"Common/proto/base"
	"Common/try"
	"GateSvr/logic"
	"fmt"
	"net"
	"time"

	"github.com/golang/protobuf/proto"
)

//服务器连接请求处理
func handleServerConnection(conn net.Conn) {
	defer try.Catch()
	defer conn.Close()
	//conn.SetReadDeadline(time.Now().Add(time.Second * 30))
	//连接之后的第一条消息，必须是验证身份，并且获得tid，sid
	buffer := make([]byte, 2048) //建立一个slice
	n, err := conn.Read(buffer)
	if err != nil {
		log.Logger.Error(conn.RemoteAddr().String(), " read server first msg error: ", err)
		return //当远程客户端连接发生错误（断开）后，终止此协程。
	}
	//buf := buffer[:n]
	logindata := &base.ServerLogin{}
	err = proto.Unmarshal(buffer[:n], logindata)
	if err != nil {
		log.Logger.Error(conn.RemoteAddr().String(), " read server first msg error: ", err)
		return //当远程客户端连接发生错误（断开）后，终止此协程。
	}
	fmt.Println(logindata)

	//登入验证
	if !logic.ServerLoginAuthentication(uint64(logindata.Tid), logindata.Password) {
		conn.Write(msg.CreateErrSvrMsgData(0, msg.Err_Nomoal, "服务注册认证失败"))
		return
	}
	//重复注册判断
	if !agentmanager.ReplaceServer(logindata.Tid, logindata.Sid) {
		conn.Write(msg.CreateErrSvrMsgData(0, msg.Err_Nomoal, "服务重复注册"))
		return
	}

	defer agentmanager.RemoveAgentServer(logindata.Tid, logindata.Sid)
	//身份验证成功,加入管理队列
	agent := agentmanager.AddAgentServer(logindata.Tid, logindata.Sid, conn)
	if agent == nil { //注册失败
		return
	}

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

		if n < msg.GetHeadLength() { //消息大小安全检测
			log.Logger.Error(conn.RemoteAddr().String(), "server msg len too samll", logindata.Tid, logindata.Sid)
			return
		}
		//  消息结构  [uint8]+[uint64]+[uint32,uint32,uint32]
		//  [uint8](后面8表示说明，0:后面是serverid 1:后面是userid)
		//  [uint64](userid(8字节)或者0(4字节)+sid(2字节)+tid(2字节))
		//  [uint32,uint32,uint32](mainid+sonid+len)+msg
		signhead := &msg.HeadSign{}
		signhead.Decode(buffer)

		switch signhead.SignType {
		case msg.Sign_serverid: //后8位是serverid
			//serverid := binary.LittleEndian.Uint64(buffer[1:9])
			buf := msg.AddSignHead(msg.Sign_serverid, agent.Serverid, buffer[msg.GetSignHeadLength():n])

			if !agentmanager.TransferToServer(signhead.SignId, buf) {
				agent.SendData(msg.CreateErrorMsg(msg.Err_MsgSendFail_ServerNoExist))
				break
			}
		case msg.Sign_userid: //后8位是userid
			//消息转发给客户端
			//userid := binary.BigEndian.Uint64(buffer[1:9])
			agentmanager.TransferToClient(signhead.SignId, buffer[9:n])
		default: //发现非法协议
			log.Errorln(conn.RemoteAddr().String(), "发现服务器的非法协议", logindata.Tid)
			return
		}

	}
}
