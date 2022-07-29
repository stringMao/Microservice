package watchdog

import (
	"Common/constant"
	"Common/kernel/go-scoket/sockets"
	"Common/log"
	"Common/msg"
	"Common/proto/base"
	"Common/proto/codes"
	"Common/proto/gateProto"
	"Common/try"
	"GateSvr/config"
	"GateSvr/logic"
	"github.com/golang/protobuf/proto"
)

type PlayerData struct{
	UserId      uint64
	//TODO  成员改成结构体，保存tid，sid，serverid等详细信息
	SvrList     map[uint32]uint32
}
type PlayerAgent struct {
	//SignHead *msg.HeadSign
	//MsgHead  *msg.HeadProto
	engine *sockets.Engine
}
//func (s *PlayerAgent) JoinServer(tid,sid uint32){
//	data:=s.engine.Data.(*PlayerData)
//	if data!=nil && data.SvrList!=nil {
//		data.SvrList[tid] = sid
//	}
//}
//func (s *PlayerAgent) QuitServer(tid uint32){
//	data:=s.engine.Data.(*PlayerData)
//	if data!=nil && data.SvrList!=nil {
//		delete(data.SvrList, tid)
//	}
//}


//GetPlayerConnection 收到一个客户端连接
func GetPlayerConnection(engine *sockets.Engine) {
	pAgent:= &PlayerAgent{
		engine: engine,
		//SignHead: new(msg.HeadSign),
		//MsgHead:  new(msg.HeadProto),
	}
	pAgent.engine.SetHandler(pAgent)
	pAgent.engine.Start()
}
func (s *PlayerAgent)CloseHandle(engine *sockets.Engine){
	if engine.Data==nil{
		return
	}
	//通知该玩家所在的 服务器 用户离开
	pPlayerData:=engine.Data.(*PlayerData)
	for k,v:=range pPlayerData.SvrList{
		pObj := &gateProto.UserJoinQuit{Userid: pPlayerData.UserId}
		G_ClientManager.SendToServer2(msg.EncodeServerID(k,v),
			msg.NewMessage(msg.CommonSvrMsg_UserOffline,0,config.App.ServerID,pObj))
	}
	//用户在线列表里删除用户
	G_ClientManager.RemovePlayer(engine)
}

func (s *PlayerAgent)BeforeHandle(client *sockets.Engine,len int,buffer []byte){
	//====这里需要捕获异常
	defer try.Catch()
	//
	pMessage := &base.Message{}
	if err:=proto.Unmarshal(buffer, pMessage);err != nil {
		log.Warnf("PlayerAgent BeforeHandle len[%d] is err:%s",len,err.Error())
		return
	}
    msgType:=msg.ParseMessageID(pMessage.MessageId)

	switch  msgType{
	case constant.TID_GateSvr:
		//同一个携程的同步处理，所以不需要new新的buf拷贝消息内容
		fn:= PlayerListener.GetHandleFunc(pMessage.MessageId)
		if fn!=nil{
			fn(s.engine, pMessage.Body)
		}
	default:
		if msgType==0 {
			return
		}
		playerData:=s.engine.Data.(*PlayerData)
		if sid,ok:=playerData.SvrList[uint32(msgType)];ok{
			//协议再包一层
			serverId:=msg.EncodeServerID(uint32(msgType), sid)
			pb:=&base.Forward{
				UserId:   playerData.UserId,
				ServerId: 0,
				Body:     pMessage.Body,
			}
			pMessage.Body=msg.ProtoMarshal(pb)
			if G_ClientManager.SendToServer(serverId, pMessage) {
				break
			}
		}
	}

}


func UserLogin(engine *sockets.Engine,buf []byte) {
	pData := &gateProto.UserLoginReq{}
	err := proto.Unmarshal(buf, pData)
	if err != nil {
		return
	}
	log.Infof("收到玩家登入消息 userid[%d] token[%s] ",pData.Userid,pData.Token)
	//登入验证
	if logic.UserLogin(pData.Userid, pData.Token) != 0 {
		//登入验证失败
		log.Debugf("玩家登入验证失败 uid[%d],token[%s]", pData.Userid, pData.Token)
		pObj:=&gateProto.UserLoginResult{
			Code: codes.Code_LoginGateSvrAuthFail,
		}
		engine.SyncSendData(msg.NewClientMessage(msg.ToUser_GateLoginResult,pObj))
		engine.Close()
		return
	}
	//顶号
	if G_ClientManager.PlayerIsExists(pData.Userid) {
		//顶号处理

		//顶号失败
		pObj:=&gateProto.UserLoginResult{
			Code: codes.Code_LoginGateSvrReplaceFail,
		}
		engine.SyncSendData(msg.NewClientMessage(msg.ToUser_GateLoginResult,pObj))
		engine.Close()
		return
	}
	//用户数据保存
	engine.Data= &PlayerData{
		UserId:  pData.Userid,
		SvrList: make(map[uint32]uint32,10),
	}
	G_ClientManager.AddPlayer(engine)

	//返回结果
	pObj:=&gateProto.UserLoginResult{
		Code: codes.Code_Success,
	}
	engine.SendData(msg.NewClientMessage(msg.ToUser_GateLoginResult,pObj))
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
