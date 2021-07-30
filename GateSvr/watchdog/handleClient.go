package watchdog

//客户端连接的消息处理
import (
	"Common/log"
	"Common/msg"
	"Common/proto/base"
	"GateSvr/logic"
	"fmt"
	"math"
	"net"
	"runtime/debug"
	"time"

	"github.com/golang/protobuf/proto"
)

//handleClientConnection 客户端连接请求处理
func handleClientConnection(conn net.Conn) {
	defer func() {
		conn.Close()
		if r := recover(); r != nil {
			log.Logger.Error("Recovered in", r, ":", string(debug.Stack()))
		}
	}()
	//连接之后的第一条消息，必须是登入获取userid
	buffer := make([]byte, 2048) //建立一个slice
	conn.SetReadDeadline(time.Now().Add(time.Second * 1))
	n, err := conn.Read(buffer)
	if err != nil {
		log.Logger.Error(conn.RemoteAddr().String(), "read first msg error: ", err)
		return //当远程客户端连接发生错误（断开）后，终止此协程。
	}
	logindata := &base.ClientLogin{}
	err = proto.Unmarshal(buffer[:n], logindata)
	if err != nil {
		log.Logger.Error(conn.RemoteAddr().String(), "encode first msg error: ", err)
		return //当远程客户端连接发生错误（断开）后，终止此协程。
	}
	fmt.Println(logindata)

	//登入验证
	if logic.UserLogin(logindata.Userid, logindata.Token) != 0 {
		//登入验证失败
		return
	}
	//身份验证成功,加入管理队列
	agent := agentmanager.AddAgentClient(logindata.Userid, conn)
	defer agentmanager.RemoveAgentClient(logindata.Userid)

	for {
		//开始通讯
		conn.SetReadDeadline(time.Now().Add(time.Second * 5)) //借此检测心跳包
		n, err := conn.Read(buffer)                           //读取客户端传来的内容
		if err != nil {
			log.Logger.Error(conn.RemoteAddr().String(), "connection error: ", err)
			return //当远程客户端连接发生错误（断开）后，终止此协程。
		}
		//特殊包-心跳包过滤  消息结构[uint8]=200
		if n == 1 && buffer[0] == 200 {
			//log.Logger.Debugln("heart")
			continue
		}

		//log.Logger.Debug(conn.RemoteAddr().String(), "receive data string:\n", string(buffer[:n]))
		if n < msg.GetHeadLength() { //消息大小安全检测
			log.Logger.Error(conn.RemoteAddr().String(), "msg len too samll", logindata.Userid)
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
			//修改消息
			tempbaseH := msg.HeadBase{SignType: msg.Sign_userid, ID: logindata.Userid}
			buf := msg.CreateMsg2(&tempbaseH, buffer[msg.GetHeadbaseLength():n])

			//serverid := binary.LittleEndian.Uint64(buffer[1:9])
			if baseH.ID < math.MaxUint16 { //小于 math.MaxUint16 说明serverid=tid
				//根据tid，转发消息
				if agent.PushData(uint16(baseH.ID), buf) == false { //没用发送成功，说明没用分配给这个客户端对应的服务器
					//分配一个tid对应的服务
					serverid, pSvr := agentmanager.AllocSvr(uint16(baseH.ID))
					if serverid == 0 || pSvr == nil { //分配失败
						agent.SendData(msg.CreateErrorMsg(msg.Err_ServerNoFind))
						break
					}
					//转发送消息给新分配的服务
					if !agentmanager.TransferToServer(serverid, buf, len(buf)) {
						agent.SendData(msg.CreateErrorMsg(msg.Err_MsgSendFail_ServerNoExist))
						break
					}
					//完成分配
					agent.SetSvr(uint16(baseH.ID), pSvr)
				}
				break
			}
			//是指定服务的特殊发送
			if !agentmanager.TransferToServer(baseH.ID, buffer, n) {
				agent.SendData(msg.CreateErrorMsg(msg.Err_MsgSendFail_ServerNoExist))
			}
		case msg.Sign_userid: //后8位是userid
			//userid := binary.BigEndian.Uint64(buffer[1:9])
			log.Logger.Error(conn.RemoteAddr().String(), "客户端不能发消息给其他客户端", logindata.Userid, baseH.ID)
			return
		default: //发现非法协议
			log.Logger.Error(conn.RemoteAddr().String(), "发现客户端的非法协议", logindata.Userid)
			return
		}
	}
}
