package watchdog

import (
	"Common/constant"
	"Common/kernel/go-scoket/scokets"
	"Common/log"
	"Common/msg"
	"Common/proto/base"
	"Common/proto/gatesvrproto"
	"Common/try"
	"GateSvr/logic"
	"github.com/golang/protobuf/proto"
	"net"
	"time"
)

type PlayerData struct{
	UserId      uint64
	SvrList     map[uint32]uint32
}
type PlayerAgent struct {
	SignHead *msg.HeadSign
	MsgHead  *msg.HeadProto
	client *scokets.Client
}

func NewPlayerAgent(session *scokets.Session,userid uint64)*PlayerAgent{
	playerdata:=&PlayerData{
		UserId: userid,
		SvrList: make(map[uint32]uint32,10),
	}
	session.SetID(userid)

	agent:= &PlayerAgent{
		client: scokets.NewClient(scokets.ServerConnect,session,playerdata,4*time.Second),
		SignHead: new(msg.HeadSign),
		MsgHead: new(msg.HeadProto),
	}
	agent.client.SetHandler(agent)

	return agent
}

//handleClientConnection 客户端连接请求处理
func GetPlayerConnection(conn net.Conn) {
	defer try.Catch()
	session:=scokets.NewSession( 0,conn,100)
	defer session.Close()
	log.Debug("handleClientConnection===================")
	//fmt.Println(conn.RemoteAddr())


	//连接之后的第一条消息，必须是登入获取userid
	conn.SetReadDeadline(time.Now().Add(time.Second * 1))
	n, err := conn.Read(session.ReadBuf)
	if err != nil {
		log.Error(conn.RemoteAddr().String(), "read first msg error: ", err)
		return //当远程客户端连接发生错误（断开）后，终止此协程。
	}
	logindata := &base.ClientLogin{}
	if err = proto.Unmarshal(session.ReadBuf[:n], logindata);err != nil{
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
	if !G_ClientManager.PlayerIsExists(logindata.Userid) {
		//顶号失败
		conn.Write(msg.CreateErrorMsg(msg.Err_Login_AuthenticationFail))
		return
	}
	agenter:=NewPlayerAgent(session,logindata.Userid)
	//身份验证成功,加入管理队列
	G_ClientManager.AddPlayerClient(agenter)
	defer G_ClientManager.RemovePlayerClient(agenter)


	//发送连接成功消息
	tPro := &gatesvrproto.PlayerInfo{
	}
	dPro, _ := proto.Marshal(tPro)
	conn.Write(msg.CreateWholeProtoData(msg.MID_Gate, msg.Gate_SC_SendPlayerData, dPro))
	log.Debugf("客户端登入成功 Userid[%d]", logindata.Userid)

	agenter.client.Start()
}

func (s *PlayerAgent)BeforeHandle(client *scokets.Client,len int,buffer []byte){
	if len < msg.GetHeadLength() { //消息大小安全检测
		log.Error("server msg len too samll")
		return
	}

	s.SignHead.Decode(buffer)

	switch s.SignHead.SignType {
	case msg.Sign_serverid: //后8位是serverid
		if s.SignHead.Tid == constant.TID_GateSvr { //发给本服务器
			msg.ParseHead(s.MsgHead,buffer)
			fn:= PlayerListener.GetHandleFunc(s.MsgHead.SonID)
			if fn!=nil{
				buf :=scokets.GetByteFormPool()
				//拷贝消息切片
				copy(buf, buffer[msg.GetHeadLength():len])
				fn(s.client,int(s.MsgHead.Len),buf)
			}
		} else if s.SignHead.Tid != 0 && s.SignHead.Sid == 0 {
			sd:=s.client.Data.(*PlayerData)
			if sid,ok:=sd.SvrList[s.SignHead.Tid];ok{
				buf :=scokets.GetByteFormPool()
				copy(buf, buffer[:len])
				msg.ChangeSignHead(msg.Sign_userid, sd.UserId, buf)//将标记头改成来源
				serverid:=msg.EncodeServerID(s.SignHead.Tid, sid)
				if G_ClientManager.SendToServer(serverid, buf) {
					break
				}
			}
			s.client.Session.Send(msg.CreateErrorMsg(msg.Err_MsgSendFail_ServerNoExist))

		} else if s.SignHead.Tid != 0 && s.SignHead.Sid != 0 {
			//允不允许客户端指向发送指定的服务，可能存在风险
		}

		//fmt.Print("ss")
	case msg.Sign_userid: //后8位是userid
		//userid := binary.BigEndian.Uint64(buffer[1:9])
		log.Error( "客户端不能发消息给其他客户端")
		return
	default: //发现非法协议
		log.Logger.Error("发现客户端的非法协议")
		return
	}
}
func (s *PlayerAgent)CloseHandle(clisent *scokets.Client){

}



func (s *PlayerAgent)JionServer(tid,sid uint32){
	data:=s.client.Data.(PlayerData)
	data.SvrList[tid]=sid
}










//
////客户端连接的消息处理
//import (
//	"Common/constant"
//	"Common/log"
//	"Common/msg"
//	"Common/proto/base"
//	"Common/proto/gatesvrproto"
//	"Common/try"
//	"GateSvr/agent"
//	"GateSvr/core/send"
//	"GateSvr/logic"
//	"GateSvr/util/msgbody"
//	"net"
//	"time"
//
//	"github.com/golang/protobuf/proto"
//)
//
////handleClientConnection 客户端连接请求处理
//func handleClientConnection(conn net.Conn) {
//	defer try.Catch()
//	defer conn.Close()
//
//	log.Debug("handleClientConnection===================")
//	//fmt.Println(conn.RemoteAddr())
//	var err error = nil
//	var readLength int = 0       //收到的消息长度
//	buffer := make([]byte, 2048) //建立一个slice
//
//	//连接之后的第一条消息，必须是登入获取userid
//	conn.SetReadDeadline(time.Now().Add(time.Second * 1))
//	readLength, err = conn.Read(buffer)
//	if err != nil {
//		log.Error(conn.RemoteAddr().String(), "read first msg error: ", err)
//		return //当远程客户端连接发生错误（断开）后，终止此协程。
//	}
//	logindata := &base.ClientLogin{}
//	err = proto.Unmarshal(buffer[:readLength], logindata)
//	if err != nil {
//		log.Errorln(conn.RemoteAddr().String(), "encode first msg error: ", err)
//		return //当远程客户端连接发生错误（断开）后，终止此协程。
//	}
//	//fmt.Println(logindata)
//
//	//登入验证
//	if logic.UserLogin(logindata.Userid, logindata.Token) != 0 {
//		//登入验证失败
//		conn.Write(msg.CreateErrorMsg(msg.Err_Login_AuthenticationFail))
//		return
//	}
//	//顶号处理
//	if !agentmanager.ReplaceClient(logindata.Userid) {
//		//顶号失败
//		conn.Write(msg.CreateErrorMsg(msg.Err_Login_AuthenticationFail))
//		return
//	}
//
//	//身份验证成功,加入管理队列
//	agenter := agentmanager.AddAgentClient(logindata.Userid, conn)
//	defer agentmanager.RemoveAgentClient(logindata.Userid)
//
//	//加载个人数据
//	if !agenter.Init() {
//		agenter.SendData(msg.CreateErrorMsg(msg.Err_Login_InitDataFail))
//		return
//	} else {
//		//发送连接成功消息
//		tPro := &gatesvrproto.PlayerInfo{
//			NickName: agenter.Player.BaseData.NickName,
//			Avatar:   agenter.Player.BaseData.Avatar,
//			Gender:   int32(agenter.Player.BaseData.Age),
//			Age:      int32(agenter.Player.BaseData.Age),
//			Score:    agenter.Player.CashData.Score,
//			Gold:     agenter.Player.CashData.Gold,
//		}
//		dPro, _ := proto.Marshal(tPro)
//		agenter.SendData(msg.CreateWholeProtoData(msg.MID_Gate, msg.Gate_SC_SendPlayerData, dPro))
//
//		log.Debugf("客户端登入成功 Userid[%d]", logindata.Userid)
//	}
//
//	//标记头变量声明
//	signhead := &msg.HeadSign{}
//	for {
//		//开始通讯
//		conn.SetReadDeadline(time.Now().Add(time.Second * 5)) //借此检测心跳包
//		readLength, err = conn.Read(buffer)                   //读取客户端传来的内容
//		if err != nil {
//			log.Debug(conn.RemoteAddr().String(), " connection error: ", err)
//			return //当远程客户端连接发生错误（断开）后，终止此协程。
//		}
//		//特殊包-心跳包过滤  消息结构[uint8]=200
//		if readLength == 1 && buffer[0] == 200 {
//			//log.Logger.Debugln("heart")
//			continue
//		}
//		//消息大小安全检测
//		if readLength < msg.GetHeadLength() {
//			log.Error(conn.RemoteAddr().String(), "msg len too samll", logindata.Userid)
//			return
//		}
//		//  消息结构  [uint8]+[uint64]+[uint32,uint32,uint32]
//		//  [uint8](后面8表示说明，0:后面是serverid 1:后面是userid)
//		//  [uint64](userid(8字节)或者0(4字节)+sid(2字节)+tid(2字节))
//		//  [uint32,uint32,uint32](mainid+sonid+len)+msg
//		signhead.Decode(buffer)
//
//		//拷贝消息切片
//		buf := make([]byte, readLength)
//		copy(buf, buffer[:readLength])
//
//		switch signhead.SignType {
//		case msg.Sign_serverid: //后8位是serverid
//			//serverid := binary.LittleEndian.Uint64(buffer[1:9])
//			if signhead.Tid == constant.TID_GateSvr { //发给本服务器
//				HandleClientMessage(agenter, buf)
//			} else if signhead.Tid != 0 && signhead.Sid == 0 {
//				serverid := agenter.GetServerId(signhead.Tid)
//				if serverid == 0 {
//					//未加入该服务器
//					agenter.SendData(msg.CreateErrorMsg(msg.Err_ServerNoFind))
//					break
//				}
//				//修改消息
//				msg.ChangeSignHead(msg.Sign_userid, logindata.Userid, buf)
//				if !agentmanager.TransferToServer(serverid, buf) {
//					//发送失败
//					//切换服务，会有一个号登入多个服的可能，建议断线重登
//					agenter.SendData(msg.CreateErrorMsg(msg.Err_MsgSendFail_ServerNoExist))
//					return
//				}
//			} else if signhead.Tid != 0 && signhead.Sid != 0 {
//				//允不允许客户端指向发送指定的服务，可能存在风险
//			}
//
//			//fmt.Print("ss")
//		case msg.Sign_userid: //后8位是userid
//			//userid := binary.BigEndian.Uint64(buffer[1:9])
//			log.Error(conn.RemoteAddr().String(), "客户端不能发消息给其他客户端", logindata.Userid, signhead.SignId)
//			return
//		default: //发现非法协议
//			log.Logger.Error(conn.RemoteAddr().String(), "发现客户端的非法协议", logindata.Userid)
//			return
//		}
//
//	}
//}
//
//func HandleClientMessage(p *agent.AgentClient, data []byte) {
//	head := msg.GetHead(data)
//
//	if head.MainID != msg.MID_Gate {
//		return
//	}
//
//	switch head.SonID {
//	case msg.Gate_CS_JionServerReq: //请求加入某个业务服务器
//		loginreq := &base.ClientJionServerReq{}
//		err := proto.Unmarshal(data[msg.GetHeadLength():], loginreq)
//		if err != nil {
//			return
//		}
//		if loginreq.Tid == constant.TID_GateSvr {
//			return
//		}
//
//		if s := p.GetServerId(loginreq.Tid); s != 0 {
//			//已经加入该type的服务器
//			p.SendData(send.CreateMsgToClient(msg.Gate_SC_ClientJionResult,
//				msgbody.MakeToClientJionServerResult(2, s)))
//			return
//		}
//		serverid := agentmanager.AllocSvr(loginreq.Tid)
//		if serverid == 0 {
//			//服务器找不到
//			p.SendData(send.CreateMsgToClient(msg.Gate_SC_ClientJionResult,
//				msgbody.MakeToClientJionServerResult(1, msg.EncodeServerID(loginreq.Tid, 0))))
//			return
//		}
//		//转发加入请求给指定服务器
//		tPro := &base.NotifyJionServerReq{Userid: p.Userid}
//		dPro, _ := proto.Marshal(tPro)
//		agentmanager.TransferToServer(serverid, send.CreateMsgToSvr(msg.Gate_SS_ClientJionReq, dPro))
//	case msg.Gate_CS_LeaveServerReq: //请求离开某个业务服务器
//		leaveReq := &base.ClientLeaveServerReq{}
//		err := proto.Unmarshal(data[msg.GetHeadLength():], leaveReq)
//		if err != nil {
//			return
//		}
//		if leaveReq.Tid == constant.TID_GateSvr {
//			return
//		}
//
//		serverid := p.GetServerId(leaveReq.Tid)
//		if serverid == 0 {
//			//没有可离开的服务器
//			p.SendData(send.CreateMsgToClient(msg.Gate_SC_ClientLeaveResult,
//				msgbody.MakeToClientLeaveServerResult(1, serverid)))
//			return
//		}
//
//		//转发加入请求给指定服务器
//		tPro := &base.NotifyLeaveServerReq{Userid: p.Userid}
//		dPro, _ := proto.Marshal(tPro)
//		agentmanager.TransferToServer(serverid, send.CreateMsgToSvr(msg.Gate_SS_ClientLeaveReq, dPro))
//	default:
//		return
//	}
//}
