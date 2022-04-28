package watchdog

import (
	"Common/constant"
	"Common/kernel/go-scoket/scokets"
	"Common/log"
	"Common/msg"
	"Common/proto/base"
	"GateSvr/config"
	"GateSvr/core/send"
	"GateSvr/logic"
	"github.com/golang/protobuf/proto"
	"time"
)

type ServerData struct {
	Tid      uint32 //服务类型id
	Sid      uint32 //同服务类型下的唯一标识id
	Serverid uint64 //0+0+sid+tid 位运算获得
}

type ServerAgent struct {
	SignHead *msg.HeadSign
	MsgHead  *msg.HeadProto
	client *scokets.Client
}

func NewServerAgent(session *scokets.Session,tid,sid uint32)*ServerAgent{
	serverid:=msg.EncodeServerID(tid, sid)
	svrdata:=&ServerData{
		Tid: tid,
		Sid: sid,
		Serverid: serverid,
	}
	session.SetID(serverid)

	agent:= &ServerAgent{
		client: scokets.NewClient(scokets.ServerConnect,session,svrdata,4*time.Second),
		SignHead: new(msg.HeadSign),
		MsgHead: new(msg.HeadProto),
	}
	agent.client.SetHandler(agent)

	return agent
}

//服务器连接请求处理
func GetServerConnection(client *scokets.Client){
	log.Debugln("收到新的服务器连接")
	pAgenter:= &ServerAgent{
		client: client,
		SignHead: new(msg.HeadSign),
		MsgHead: new(msg.HeadProto),
	}
	pAgenter.client.SetHandler(pAgenter)
	//G_ClientManager.AddServerClient(pAgenter)
	//defer G_ClientManager.RemoveServerClient(pAgenter)

	pAgenter.client.Start()

}
//func GetServerConnection2(conn net.Conn) {
//	defer try.Catch()
//	session:=scokets.NewSession( 0,conn,200)
//	defer session.Close()
//
//	conn.SetReadDeadline(time.Now().Add(time.Millisecond))
//	//连接之后的第一条消息，必须是验证身份，并且获得tid，sid
//	n, err := conn.Read(session.ReadBuf)
//	if err != nil {
//		log.Logger.Error(conn.RemoteAddr().String(), " read server first msg error: ", err)
//		return //当远程客户端连接发生错误（断开）后，终止此协程。
//	}
//	//buf := buffer[:n]
//	logindata := &base.ServerLogin{}
//	err = proto.Unmarshal(session.ReadBuf[:n], logindata)
//	if err != nil {
//		log.Logger.Error(conn.RemoteAddr().String(), "handleServerConnection proto.Unmarshal ServerLogin error: ", err)
//		return //当远程客户端连接发生错误（断开）后，终止此协程。
//	}
//	//fmt.Println(logindata)
//
//	//登入验证
//	if !logic.ServerLoginAuthentication(uint64(logindata.Tid), logindata.Password) {
//		conn.Write(msg.CreateErrSvrMsgData(0, msg.Err_Nomoal, "服务注册认证失败"))
//		log.Infof("server 身份验证失败 tid:%d,sid:%d",logindata.Tid,logindata.Sid)
//		return
//	}
//
//	serverid:=msg.EncodeServerID(logindata.Tid, logindata.Sid)
//	//重复注册判断
//	if G_ClientManager.ServerIsExists(serverid) {
//		conn.Write(msg.CreateErrSvrMsgData(0, msg.Err_Nomoal, "服务重复注册"))
//		log.Infof("server 重复注册 tid:%d,sid:%d",logindata.Tid,logindata.Sid)
//		return
//	}
//
//	//身份验证成功
//	pAgenter:=NewServerAgent(session,logindata.Tid, logindata.Sid)
//
//	G_ClientManager.AddServerClient(pAgenter)
//	defer G_ClientManager.RemoveServerClient(pAgenter)
//
//	loginResult := &base.LoginResult{}
//	loginResult.Code = 0
//	loginResult.Tid = uint32(config.App.TID)
//	loginResult.Sid = uint32(config.App.SID)
//	loginResult.Name = constant.GetServerIDName(config.App.TID, config.App.SID)
//	ploginResult, _ := proto.Marshal(loginResult)
//	conn.Write(send.CreateMsgToSvr(msg.Gate_SS_SvrLoginResult, ploginResult))
//	log.Infof("服务注册成功，ServerID[%d] TID[%d] SID[%d]", serverid, logindata.Tid, logindata.Sid)
//
//	pAgenter.client.Start()
//
//}
//

func (s *ServerAgent)BeforeHandle(client *scokets.Client,len int,buffer []byte){
	if len < msg.GetHeadLength() { //消息大小安全检测
		log.Error("server msg len too samll")
		return
	}
	//  消息结构  [uint8]+[uint64]+[uint32,uint32,uint32]
	//  [uint8](后面8表示说明，0:后面是serverid 1:后面是userid)
	//  [uint64](userid(8字节)或者0(4字节)+sid(2字节)+tid(2字节))
	//  [uint32,uint32,uint32](mainid+sonid+len)+msg
	s.SignHead.Decode(buffer)

	switch s.SignHead.SignType {
	case msg.Sign_serverid: //后8位是serverid
		if s.SignHead.Tid == constant.TID_GateSvr {
			msg.ParseHead(s.MsgHead,buffer)
			fn:=ServerListener.GetHandleFunc(s.MsgHead.SonID)
			if fn==nil{
				log.Errorf("收到无效消息，sonid[%d]",s.MsgHead.SonID)
				break
			}
			buf :=scokets.GetByteFormPool()
			//拷贝消息切片
			copy(buf, buffer[msg.GetHeadLength():len])
			fn(s.client,int(s.MsgHead.Len),buf)
		} else if s.SignHead.Tid != 0 && s.SignHead.Sid != 0 {
			sd:=s.client.Data.(*ServerData)
			buf :=scokets.GetByteFormPool()
			copy(buf, buffer[:len])
			msg.ChangeSignHead(msg.Sign_serverid, sd.Serverid, buf)//将标记头改成来源

			if !G_ClientManager.SendToServer(s.SignHead.SignId, buf[:len]) {
				s.client.Session.Send(msg.CreateErrorMsg(msg.Err_MsgSendFail_ServerNoExist))
			}
		} else if s.SignHead.Tid != 0 && s.SignHead.Sid == 0 {
			//给该类型服务器群发本消息

		}
	case msg.Sign_userid: //后8位是userid
		//消息转发给客户端
		buf :=scokets.GetByteFormPool()
		n:=copy(buf, buffer[msg.GetSignHeadLength():len])
		G_ClientManager.SendToPlayer(s.SignHead.SignId, buf[:n])
	default: //发现非法协议
		log.Errorln( "发现服务器的非法协议")
		return
	}
}

func (s *ServerAgent)CloseHandle(client *scokets.Client){
	if client.Data==nil{
		return
	}
	G_ClientManager.RemoveServerClient(client)
}


func SvrRegister(client *scokets.Client,len int,buf []byte){
	log.Debugln("收到服务注册消息")
	obj := &base.ServerLogin{}
	err := proto.Unmarshal(buf[:len], obj)
	if err != nil {
		client.Close()
		return
	}

	//登入验证
	if !logic.ServerLoginAuthentication(uint64(obj.Tid), obj.Password) {
		client.SendData(msg.CreateErrSvrMsgData(0, msg.Err_Nomoal, "服务注册认证失败"))
		log.Infof("server 身份验证失败 tid:%d,sid:%d",obj.Tid,obj.Sid)
		client.Close()
		return
	}

	serverId:=msg.EncodeServerID(obj.Tid, obj.Sid)
	//重复注册判断
	if G_ClientManager.ServerIsExists(serverId) {
		client.SendData(msg.CreateErrSvrMsgData(0, msg.Err_Nomoal, "服务重复注册"))
		log.Infof("server 重复注册 tid:%d,sid:%d",obj.Tid,obj.Sid)
		client.Close()
		return
	}

	//身份验证成功========
	svrData:=&ServerData{
		Tid: obj.Tid,
		Sid: obj.Sid,
		Serverid: serverId,
	}
	client.Data=svrData

	G_ClientManager.AddServerClient(client)


	loginResult := &base.LoginResult{}
	loginResult.Code = 0
	loginResult.Tid = uint32(config.App.TID)
	loginResult.Sid = uint32(config.App.SID)
	loginResult.Name = constant.GetServerIDName(config.App.TID, config.App.SID)
	ploginResult, _ := proto.Marshal(loginResult)
	client.SendData(send.CreateMsgToSvr(msg.Gate_SS_SvrLoginResult, ploginResult))
	log.Infof("服务注册成功，ServerID[%d] TID[%d] SID[%d]", serverId, obj.Tid, obj.Sid)

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
