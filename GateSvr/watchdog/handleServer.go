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

type ServerData struct {
	Tid      uint32 //服务类型id
	Sid      uint32 //同服务类型下的唯一标识id
	Serverid uint64 //0+0+sid+tid 位运算获得
	SvrType  uint32
}

type ServerAgent struct {
	//SignHead *msg.HeadSign
	//MsgHead  *msg.HeadProto
	engine *sockets.Engine
}

// GetServerConnection 服务器连接请求处理
func GetServerConnection(engine *sockets.Engine){
	log.Debugln("收到新的服务器连接")
	pAgent:= &ServerAgent{
		engine: engine,
		//SignHead: new(msg.HeadSign),
		//MsgHead:  new(msg.HeadProto),
	}
	pAgent.engine.SetHandler(pAgent)
	pAgent.engine.Start()
}
func (s *ServerAgent)CloseHandle(engine *sockets.Engine){
	log.Debugf("server CloseHandle")
	if engine.Data==nil{
		return
	}

	pSvrData:=engine.Data.(*ServerData)
	log.Debugf("server CloseHandle tid[%d] sid[%d]",pSvrData.Tid,pSvrData.Sid)
	//检查所有玩家，有连接再该服务器的都通知离开
	G_ClientManager.PlayerMap.Range(func (k,v interface{})bool{
		pClientEngine:=v.(*sockets.Engine)
		pPlayerData :=pClientEngine.Data.(*PlayerData)
		if sid,ok:= pPlayerData.SvrList[pSvrData.Tid];ok && sid==pSvrData.Sid{
			delete(pPlayerData.SvrList,pSvrData.Tid)
			pObj:=&gateProto.JoinQuitServerResult{
				Code:   codes.Code_Success, //0成功
				Tid:  	pSvrData.Tid,
			}
			pClientEngine.SendData(msg.NewClientMessage(msg.ToUser_QuitSvrResult,pObj))
		}
		return true
	})
	//最后"服务在线列表"移除该服务器
	G_ClientManager.RemoveServer(engine)
}

func (s *ServerAgent)BeforeHandle(client *sockets.Engine,len int,buffer []byte){
	//需要捕获异常
	defer try.Catch()
	//log.Debugln("收到一条消息")

	pMessage:=&base.Message{}
	err := proto.Unmarshal(buffer, pMessage)
	if err != nil {
		log.Warnf("ServerAgent BeforeHandle len[%d] is err:%s",len,err.Error())
		return
	}

	pForward := &base.Forward{}
	if err = proto.Unmarshal(pMessage.Body,pForward );err!=nil{
		log.Warnln("ServerAgent BeforeHandle is err2:",err)
		return
	}

	if pForward.UserId>0 && pForward.ServerId==0{
		//转发给指定客户端
		pMessage.Body=pForward.Body
		G_ClientManager.SendToPlayer(pForward.UserId, pMessage)
	}else if pForward.UserId>0 && pForward.ServerId>0{
		//发送给该用户连接着的某个svr，ServerId=tid，找该用户连接的tid
		tid, _ := msg.DecodeServerID(pForward.ServerId)
		uid := pForward.UserId
		pForward.UserId = 0
		pForward.ServerId = client.Data.(*ServerData).Serverid
		pMessage.Body=msg.ProtoMarshal(pForward)
		G_ClientManager.SendToPlayerOfTid(uid, tid,pMessage)
	}else if pForward.UserId==0 && pForward.ServerId>0{
		if pForward.ServerId==config.App.ServerID{ //发生给本服务的消息
			fn:=ServerListener.GetHandleFunc(pMessage.MessageId)
			if fn==nil{
				log.Errorf("收到无效消息，msgId[%d]",pMessage.MessageId)
				return
			}
			fn(s.engine,pForward.Body)
		}else{
			//发生消息给指定服务器，注意检验身份。解析出tid和sid，
			tid,sid:=msg.DecodeServerID(pForward.ServerId)
			targetId:=pForward.ServerId
			pServerData:=client.Data.(*ServerData)
			pForward.ServerId=pServerData.Serverid
			pForward.UserId=0
			pMessage.Body=msg.ProtoMarshal(pForward)

			if tid>0 && sid==0{//如果tid>0,sid=0,则发送给所有给类型服务器
				G_ClientManager.BroadcastToTid(tid,pMessage)
			}else if tid>0 && sid>0{//如果tid>0,sid>0,发送给特定的一个服务器
				G_ClientManager.SendToServer(targetId,pMessage)
			}
		}
	}else if pForward.UserId==0 && pForward.ServerId==0{ //发送给所有用户
		pMessage.Body=pForward.Body
		G_ClientManager.BroadcastToPlayers(nil,pMessage)
	}
}

func SvrRegister(engine *sockets.Engine,buf []byte){
	pData := &gateProto.SvrRegisterReq{}
	err := proto.Unmarshal(buf, pData)
	if err != nil {
		engine.Close()
		return
	}
	log.Infof("收到服务注册消息tid[%d] sid[%d] svrType[%d]",pData.Tid,pData.Sid,pData.SvrType)

	//登入验证
	if !logic.ServerLoginAuthentication(uint64(pData.Tid), pData.Password) {
		log.Infof("server 身份验证失败 tid:%d,sid:%d", pData.Tid, pData.Sid)

		pObj:=&gateProto.SvrRegisterResult{
			Code: codes.Code_SvrRegister_AuthFail,
			Tid:  uint32(config.App.TID),
			Sid:  uint32(config.App.SID),
		}
		engine.SyncSendData(msg.NewMessage(msg.CommonSvrMsg_SvrRegisterResult,0,config.App.ServerID,pObj))
		engine.Close()
		return
	}

	serverId:=msg.EncodeServerID(pData.Tid, pData.Sid)
	//重复注册判断
	if G_ClientManager.ServerIsExists(serverId) {
		log.Infof("server 重复注册 tid:%d,sid:%d", pData.Tid, pData.Sid)
		pObj:=&gateProto.SvrRegisterResult{
			Code: codes.Code_SvrRegister_ExistedFail,
			Tid:  uint32(config.App.TID),
			Sid:  uint32(config.App.SID),
		}
		engine.SyncSendData(msg.NewMessage(msg.CommonSvrMsg_SvrRegisterResult,0,config.App.ServerID,pObj))
		engine.Close()
		return
	}

	//身份验证成功========
	engine.Data=&ServerData{
		Tid:      pData.Tid,
		Sid:      pData.Sid,
		Serverid: serverId,
	}
	G_ClientManager.AddServer(engine)

	pObj:=&gateProto.SvrRegisterResult{
		Code: codes.Code_Success,
		Tid:  uint32(config.App.TID),
		Sid:  uint32(config.App.SID),
		Name: constant.GetServerIDName(config.App.TID, config.App.SID),
	}

	engine.SendData(msg.NewMessage(msg.CommonSvrMsg_SvrRegisterResult,0,config.App.ServerID,pObj))
	log.Infof("服务注册成功，ServerID[%d] TID[%d] SID[%d]", serverId, pData.Tid, pData.Sid)
	return
}












//
////服务器连接的消息处理
//import (
//	"Common/constant"
//	"Common/log"
//	"Common/msg"
//	"Common/proto/base"
//	"Common/try"
//	"GateSvr/agent"
//	"GateSvr/config"
//	"GateSvr/core/send"
//	"GateSvr/logic"
//	"GateSvr/util/msgbody"
//	"fmt"
//	"net"
//	"time"
//
//	"github.com/golang/protobuf/proto"
//)
//
////服务器连接请求处理
//func HandleServerConnection(conn net.Conn) {
//	defer try.Catch()
//	defer conn.Close()
//	//conn.SetReadDeadline(time.Now().Add(time.Second * 30))
//	//连接之后的第一条消息，必须是验证身份，并且获得tid，sid
//	buffer := make([]byte, 2048) //建立一个slice
//	n, err := conn.Read(buffer)
//	if err != nil {
//		log.Logger.Error(conn.RemoteAddr().String(), " read server first msg error: ", err)
//		return //当远程客户端连接发生错误（断开）后，终止此协程。
//	}
//	//buf := buffer[:n]
//	logindata := &base.ServerLogin{}
//	err = proto.Unmarshal(buffer[:n], logindata)
//	if err != nil {
//		log.Logger.Error(conn.RemoteAddr().String(), "handleServerConnection proto.Unmarshal ServerLogin error: ", err)
//		return //当远程客户端连接发生错误（断开）后，终止此协程。
//	}
//	fmt.Println(logindata)
//
//	//登入验证
//	if !logic.ServerLoginAuthentication(uint64(logindata.Tid), logindata.Password) {
//		conn.Write(msg.CreateErrSvrMsgData(0, msg.Err_Nomoal, "服务注册认证失败"))
//		return
//	}
//	//重复注册判断
//	if !agentmanager.ReplaceServer(logindata.Tid, logindata.Sid) {
//		conn.Write(msg.CreateErrSvrMsgData(0, msg.Err_Nomoal, "服务重复注册"))
//		return
//	}
//
//	defer agentmanager.RemoveAgentServer(logindata.Tid, logindata.Sid)
//	//身份验证成功,加入管理队列
//	agenter := agentmanager.AddAgentServer(logindata.Tid, logindata.Sid, conn)
//
//	loginResult := &base.LoginResult{}
//	if agenter == nil { //注册失败
//		loginResult.Code = 1
//		ploginResult, _ := proto.Marshal(loginResult)
//		conn.Write(send.CreateMsgToSvr(msg.Gate_SS_SvrLoginResult, ploginResult))
//
//		log.Infof("服务器登入失败，ServerID[%d] TID[%d] SID[%d]", agenter.Serverid, logindata.Tid, logindata.Sid)
//		return
//	} else {
//		//conn.Write(msg.CreateErrSvrMsgData(0, msg.Err_Nomoal, "服务与网关服注册成功"))
//		loginResult.Code = 0
//		loginResult.Tid = uint32(config.App.TID)
//		loginResult.Sid = uint32(config.App.SID)
//		loginResult.Name = constant.GetServerIDName(config.App.TID, config.App.SID)
//		ploginResult, _ := proto.Marshal(loginResult)
//		conn.Write(send.CreateMsgToSvr(msg.Gate_SS_SvrLoginResult, ploginResult))
//
//		log.Infof("服务器登入成功，ServerID[%d] TID[%d] SID[%d]", agenter.Serverid, logindata.Tid, logindata.Sid)
//	}
//
//	signhead := &msg.HeadSign{}
//	for {
//		//开始通讯
//		conn.SetReadDeadline(time.Now().Add(time.Second * 10)) //借此检测心跳包
//		n, err := conn.Read(buffer)                            //读取客户端传来的内容
//		if err != nil {
//			log.Logger.Debug(conn.RemoteAddr().String(), " server handleServerConnection error: ", err)
//			return //当远程客户端连接发生错误（断开）后，终止此协程。
//		}
//		//特殊包-心跳包过滤  消息结构[uint8]=200
//		if n == 1 && buffer[0] == 200 {
//			//log.Logger.Debugln("heart")
//			continue
//		}
//
//		if n < msg.GetHeadLength() { //消息大小安全检测
//			log.Logger.Error(conn.RemoteAddr().String(), "server msg len too samll", logindata.Tid, logindata.Sid)
//			return
//		}
//		//  消息结构  [uint8]+[uint64]+[uint32,uint32,uint32]
//		//  [uint8](后面8表示说明，0:后面是serverid 1:后面是userid)
//		//  [uint64](userid(8字节)或者0(4字节)+sid(2字节)+tid(2字节))
//		//  [uint32,uint32,uint32](mainid+sonid+len)+msg
//		signhead.Decode(buffer)
//
//		//拷贝消息切片
//		buf := make([]byte, n)
//		copy(buf, buffer[:n])
//
//		switch signhead.SignType {
//		case msg.Sign_serverid: //后8位是serverid
//			//serverid := binary.LittleEndian.Uint64(buffer[1:9])
//			if signhead.Tid == constant.TID_GateSvr {
//				HandleServerMessage(agenter, buf)
//			} else if signhead.Tid != 0 && signhead.Sid != 0 {
//				msg.ChangeSignHead(msg.Sign_serverid, agenter.Serverid, buf)
//
//				if !agentmanager.TransferToServer(signhead.SignId, buf) {
//					agenter.SendData(msg.CreateErrorMsg(msg.Err_MsgSendFail_ServerNoExist))
//				}
//			} else if signhead.Tid != 0 && signhead.Sid == 0 {
//				//查找给 该玩家分配的服务器是哪一台，然后转发过去
//			}
//		case msg.Sign_userid: //后8位是userid
//			//消息转发给客户端
//			//userid := binary.BigEndian.Uint64(buffer[1:9])
//			agentmanager.TransferToClient(signhead.SignId, buffer[msg.GetSignHeadLength():n])
//		default: //发现非法协议
//			log.Errorln(conn.RemoteAddr().String(), "发现服务器的非法协议", logindata.Tid)
//			return
//		}
//
//	}
//}
//
//func HandleServerMessage(p *agent.AgentServer, data []byte) {
//	head := msg.GetHead(data)
//
//	if head.MainID != msg.MID_Gate {
//		return
//	}
//
//	switch head.SonID {
//	case msg.Gate_SS_ClientJionResult: //用户加入业务服务器的结果返回
//		jionResult := &base.NotifyJionServerResult{}
//		if err := proto.Unmarshal(data[msg.GetHeadLength():], jionResult); err != nil {
//			//协议解析错误
//			return
//		}
//
//		if jionResult.Codeid != 0 {
//			//加入失败
//			agentmanager.TransferToClient(jionResult.Userid, send.CreateMsgToClient(msg.Gate_SC_ClientJionResult,
//				msgbody.MakeToClientJionServerResult(int(jionResult.Codeid), msg.EncodeServerID(p.Tid, 0))))
//			return
//		}
//		cAgenter := agentmanager.GetAgentClient(jionResult.Userid)
//		if cAgenter == nil {
//			return
//		}
//		//保存jion结果
//		cAgenter.SetServerId(p.Tid, p.Serverid)
//
//		//成功通知客户端
//		agentmanager.TransferToClient(jionResult.Userid, send.CreateMsgToClient(msg.Gate_SC_ClientJionResult,
//			msgbody.MakeToClientJionServerResult(int(jionResult.Codeid), p.Serverid)))
//	case msg.Gate_SS_ClientLeaveResult: //用户离开业务服务器结果
//		pData := &base.NotifyLeaveServerResult{}
//		if err := proto.Unmarshal(data[msg.GetHeadLength():], pData); err != nil {
//			//协议解析错误
//			return
//		}
//
//		if pData.Codeid != 0 {
//			//离开失败
//			agentmanager.TransferToClient(pData.Userid, send.CreateMsgToClient(msg.Gate_SC_ClientLeaveResult,
//				msgbody.MakeToClientLeaveServerResult(int(pData.Codeid), p.Serverid)))
//			return
//		}
//		cAgenter := agentmanager.GetAgentClient(pData.Userid)
//		if cAgenter == nil {
//			return
//		}
//		//保存leave结果
//		cAgenter.SetServerId(p.Tid, 0)
//
//		//成功通知客户端
//		agentmanager.TransferToClient(pData.Userid, send.CreateMsgToClient(msg.Gate_SC_ClientLeaveResult,
//			msgbody.MakeToClientJionServerResult(int(pData.Codeid), p.Serverid)))
//	default:
//		return
//	}
//}
