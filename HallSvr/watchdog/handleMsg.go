package watchdog

import (
	"Common/kernel/go-scoket/sockets"
	"Common/log"
	"Common/msg"
	"Common/proto/codes"
	"Common/proto/gateProto"
	"github.com/golang/protobuf/proto"
)

func NewMsgToClient(userid uint64,id uint32,pb proto.Message)[]byte{
	return msg.NewMessage(id,userid,0,pb)
}
func NewMsgToSvr(serverid uint64,id uint32,pb proto.Message)[]byte{
	return msg.NewMessage(id,0,serverid,pb)
}
func NewMsgToSvr2(engine *sockets.Engine,id uint32,pb proto.Message)[]byte{
	data:=engine.Data.(*ServerData)
	return msg.NewMessage(id,0,data.Serverid,pb)
}

//=具体业务消息逻辑处理===========================

// RegisterGateResult 网关服登入结果
func RegisterGateResult(engine *sockets.Engine,buf []byte){
	pData:=&gateProto.SvrRegisterResult{}
	pForward:=msg.ProtoDecode(buf,pData)
	if pForward==nil{
		return
	}
	log.Debugf("向网关注册返回 gate-sid[%d] code[%d]",pData.Sid,pData.Code)

	if pData.Code==codes.Code_Success{
		serverId:=msg.EncodeServerID(pData.Tid,pData.Sid)
		engine.Data=&ServerData{
			Tid:      pData.Tid,
			Sid:      pData.Sid,
			Serverid: serverId,
			SvrType:  0,
		}
		//ServerList.Add(msg.EncodeServerID(pData.Tid,pData.Sid),engine)
	}
}
func UserJoin(engine *sockets.Engine,buf []byte){
	pData:=&gateProto.UserJoinQuit{}
	pForward:=msg.ProtoDecode(buf,pData)
	if pForward==nil{
		return
	}
	log.Debugf("UserJoin userid[%d]",pData.Userid)
    pObj:=&gateProto.UserJoinQuitResult{
		Code:   codes.Code_Success,
		Userid: pData.Userid,
	}

	engine.SendData(NewMsgToSvr2(engine,msg.ToGateSvr_UserJoinSvrResult,pObj))
}
func UserQuit(engine *sockets.Engine,buf []byte){
	pData:=&gateProto.UserJoinQuit{}
	pForward:=msg.ProtoDecode(buf,pData)
	if pForward==nil{
		return
	}
	log.Debugf("UserQuit userid[%d]",pData.Userid)
	pObj:=&gateProto.UserJoinQuitResult{
		Code:   codes.Code_Success,
		Userid: pData.Userid,
	}
	engine.SendData(NewMsgToSvr2(engine,msg.ToGateSvr_UserQuitSvrResult,pObj))
}
func UserOffline(engine *sockets.Engine,buf []byte){
	pData := &gateProto.UserJoinQuit{}
	pForward:=msg.ProtoDecode(buf,pData)
	if pForward==nil{
		return
	}
	log.Debugf("UserOffline userid[%d]",pData.Userid)
}

//func DoPlayerTestMsg(srcType uint8,srcId uint64,s *scokets.Connector,buf []byte){
//	msgData := &base.TestMsg{}
//	err := proto.Unmarshal(buf, msgData)
//	if err != nil {
//		log.Errorln("协议解析失败2:", err)
//		return //当远程客户端连接发生错误（断开）后，终止此协程。
//	}
//	fmt.Printf("userid:%d,str:%s\n", srcId, msgData.Txt)
//
//	//回复
//	msgData = &base.TestMsg{
//		Txt: fmt.Sprintf("收到测试消息,我是[%s]", constant.GetServerIDName(config.App.TID, config.App.SID)),
//	}
//	dPro, _ := proto.Marshal(msgData)
//	//testmsg := msg.CreateWholeMsgData(msg.Sign_userid, userid, msg.MID_Test, msg.Test_1, dPro)
//	s.SendData(send.CreateMsgToClient(srcId, msg.MID_Test, msg.Test_1, dPro))
//}

