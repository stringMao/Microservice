package watchdog

//客户端连接的消息处理
import (
	"Common/log"
	"Common/msg"
	"Common/proto/base"
	"Common/proto/gatesvrproto"
	"Common/try"
	"GateSvr/logic"
	"net"
	"time"

	"github.com/golang/protobuf/proto"
)

//handleClientConnection 客户端连接请求处理
func handleClientConnection(conn net.Conn) {
	defer try.Catch()
	defer conn.Close()

	//fmt.Println(conn.RemoteAddr())
	var err error = nil
	var readLength int = 0       //收到的消息长度
	buffer := make([]byte, 2048) //建立一个slice

	//连接之后的第一条消息，必须是登入获取userid
	conn.SetReadDeadline(time.Now().Add(time.Second * 1))
	readLength, err = conn.Read(buffer)
	if err != nil {
		log.Error(conn.RemoteAddr().String(), "read first msg error: ", err)
		return //当远程客户端连接发生错误（断开）后，终止此协程。
	}
	logindata := &base.ClientLogin{}
	err = proto.Unmarshal(buffer[:readLength], logindata)
	if err != nil {
		log.Errorln(conn.RemoteAddr().String(), "encode first msg error: ", err)
		return //当远程客户端连接发生错误（断开）后，终止此协程。
	}
	//fmt.Println(logindata)

	//登入验证
	if logic.UserLogin(logindata.Userid, logindata.Token) != 0 {
		//登入验证失败
		conn.Write(msg.CreateErrorMsg(msg.Err_Login_AuthenticationFail))
		return
	}
	//顶号处理
	if !agentmanager.ReplaceClient(logindata.Userid) {
		//顶号失败
		conn.Write(msg.CreateErrorMsg(msg.Err_Login_AuthenticationFail))
		return
	}

	//身份验证成功,加入管理队列
	defer agentmanager.RemoveAgentClient(logindata.Userid)
	agent := agentmanager.AddAgentClient(logindata.Userid, conn)

	//加载个人数据
	if !agent.Init() {
		agent.SendData(msg.CreateErrorMsg(msg.Err_Login_InitDataFail))
		return
	} else {
		//发送连接成功消息
		tPro := &gatesvrproto.PlayerInfo{
			NickName: agent.Player.BaseData.NickName,
			Avatar:   agent.Player.BaseData.Avatar,
			Gender:   int32(agent.Player.BaseData.Age),
			Age:      int32(agent.Player.BaseData.Age),
			Score:    agent.Player.CashData.Score,
			Gold:     agent.Player.CashData.Gold,
		}
		dPro, _ := proto.Marshal(tPro)
		agent.SendData(msg.CreateWholeProtoData(msg.MID_Gate, msg.Gate_SendPlayerData, dPro))
	}

	//标记头变量声明
	signhead := &msg.HeadSign{}
	for {
		//开始通讯
		conn.SetReadDeadline(time.Now().Add(time.Second * 5)) //借此检测心跳包
		readLength, err = conn.Read(buffer)                   //读取客户端传来的内容
		if err != nil {
			log.Debug(conn.RemoteAddr().String(), " connection error: ", err)
			return //当远程客户端连接发生错误（断开）后，终止此协程。
		}
		//特殊包-心跳包过滤  消息结构[uint8]=200
		if readLength == 1 && buffer[0] == 200 {
			//log.Logger.Debugln("heart")
			continue
		}
		//消息大小安全检测
		if readLength < msg.GetHeadLength() {
			log.Error(conn.RemoteAddr().String(), "msg len too samll", logindata.Userid)
			return
		}
		//  消息结构  [uint8]+[uint64]+[uint32,uint32,uint32]
		//  [uint8](后面8表示说明，0:后面是serverid 1:后面是userid)
		//  [uint64](userid(8字节)或者0(4字节)+sid(2字节)+tid(2字节))
		//  [uint32,uint32,uint32](mainid+sonid+len)+msg
		signhead.Decode(buffer)

		switch signhead.SignType {
		case msg.Sign_serverid: //后8位是serverid
			//修改消息
			buf := msg.AddSignHead(msg.Sign_userid, logindata.Userid, buffer[msg.GetSignHeadLength():readLength])

			//serverid := binary.LittleEndian.Uint64(buffer[1:9])
			if signhead.Tid == 0 { //发给本服务器
				agent.SendData(msg.CreateErrorMsg(99))
				break
			} else if signhead.Tid != 0 && signhead.Sid == 0 {
				serverid := agent.GetServerId(signhead.Tid)
				if serverid == 0 {
					serverid = agentmanager.AllocSvr(signhead.Tid)
					if serverid == 0 {
						//服务器找不到
						agent.SendData(msg.CreateErrorMsg(msg.Err_ServerNoFind))
						break
					}
					agent.SetServerId(signhead.Tid, serverid)
				}
				if !agentmanager.TransferToServer(serverid, buf) {
					//发送失败
					//切换服务，会有一个号登入多个服的可能，建议断线重登
					agent.SendData(msg.CreateErrorMsg(msg.Err_MsgSendFail_ServerNoExist))
					return
				}
			} else if signhead.Tid != 0 && signhead.Sid != 0 {
				//允不允许客户端指向发送指定的服务，可能存在风险
			}

			//fmt.Print("ss")
		case msg.Sign_userid: //后8位是userid
			//userid := binary.BigEndian.Uint64(buffer[1:9])
			log.Error(conn.RemoteAddr().String(), "客户端不能发消息给其他客户端", logindata.Userid, signhead.SignId)
			return
		default: //发现非法协议
			log.Logger.Error(conn.RemoteAddr().String(), "发现客户端的非法协议", logindata.Userid)
			return
		}

	}
}
