package watchdog

//服务器连接的消息处理
import (
	"Common/log"
	"Common/msg"
	"Common/proto/base"
	"fmt"
	"net"
	"runtime/debug"
	"time"

	"github.com/golang/protobuf/proto"
)

//服务器连接请求处理
func handleServerConnection(conn net.Conn) {
	defer func() {
		conn.Close()
		if r := recover(); r != nil {
			log.Logger.Error("Recovered in", r, ":", string(debug.Stack()))
		}
	}()
	//conn.SetReadDeadline(time.Now().Add(time.Second * 30))
	//连接之后的第一条消息，必须是验证身份，并且获得tid，sid
	buffer := make([]byte, 2048) //建立一个slice
	n, err := conn.Read(buffer)
	if err != nil {
		log.Logger.Error(conn.RemoteAddr().String(), "read server first msg error: ", err)
		return //当远程客户端连接发生错误（断开）后，终止此协程。
	}
	//buf := buffer[:n]
	logindata := &base.ServerLogin{}
	err = proto.Unmarshal(buffer[:n], logindata)
	if err != nil {
		log.Logger.Error(conn.RemoteAddr().String(), "read first msg error4: ", err)
		return //当远程客户端连接发生错误（断开）后，终止此协程。
	}
	fmt.Println(logindata)

	//登入验证
	//logindata.Password

	//身份验证成功,加入管理队列
	agent := agentmanager.AddAgentServer(logindata.Serverid, conn)
	defer agentmanager.RemoveAgentServer(logindata.Serverid)
	for {
		//开始通讯
		conn.SetReadDeadline(time.Now().Add(time.Second * 10)) //借此检测心跳包
		n, err := conn.Read(buffer)                            //读取客户端传来的内容
		if err != nil {
			log.Logger.Error(conn.RemoteAddr().String(), "server connection error: ", err)
			return //当远程客户端连接发生错误（断开）后，终止此协程。
		}
		//特殊包-心跳包过滤  消息结构[uint8]=200
		if n == 1 && buffer[0] == 200 {
			//log.Logger.Debugln("heart")
			continue
		}

		if n < msg.GetHeadLength() { //消息大小安全检测
			log.Logger.Error(conn.RemoteAddr().String(), "server msg len too samll", logindata.Serverid)
			return
		}
		//  消息结构  [uint8]+[uint64]+[uint32,uint32,uint32]
		//  [uint8](后面8表示说明，0:后面是serverid 1:后面是userid)
		//  [uint64](userid(8字节)或者0(4字节)+sid(2字节)+tid(2字节))
		//  [uint32,uint32,uint32](mainid+sonid+len)+msg
		baseH := &msg.HeadBase{}
		baseH.Decode(buffer)

		switch baseH.SignType {
		case msg.Sign_serverid: //后8位是serverid
			//serverid := binary.LittleEndian.Uint64(buffer[1:9])
			tempbaseH := msg.HeadBase{SignType: msg.Sign_serverid, ID: logindata.Serverid}
			buf := msg.CreateMsg2(&tempbaseH, buffer[msg.GetHeadbaseLength():n])

			if !agentmanager.TransferToServer(baseH.ID, buf, len(buf)) {
				agent.SendData(msg.CreateErrorMsg(msg.Err_MsgSendFail_ServerNoExist))
				break
			}
		case msg.Sign_userid: //后8位是userid
			//消息转发给客户端
			//userid := binary.BigEndian.Uint64(buffer[1:9])
			agentmanager.TransferToClient(baseH.ID, buffer[9:n])
		default: //发现非法协议
			log.Logger.Error(conn.RemoteAddr().String(), "发现客户端的非法协议", logindata.Serverid)
			return
		}

	}
}
