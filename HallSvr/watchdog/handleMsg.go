package watchdog

import (
	"Common/constant"
	"Common/kernel/go-scoket/scokets"
	"Common/log"
	"Common/msg"
	"Common/proto/base"
	"HallSvr/config"
	"HallSvr/core/send"
	"fmt"

	"github.com/golang/protobuf/proto"
)


//=具体业务消息逻辑处理===========================
//网关服登入结果
func DoLoginGateSvr(srcType uint8,srcId uint64,s *scokets.Connector,buf []byte){
	serverer:=ServerList.Get(srcId)
	if serverer.bLogin{
		log.Errorf("GateSvr[%d] 已经登入过了", srcId)
		return
	}
	loginResult := &base.LoginResult{}
	err := proto.Unmarshal(buf, loginResult)
	if err != nil || loginResult.Code != 0 {
		ServerList.Remove(srcId)
		//log.Errorf("GateSvr Login is fail. err:%s", err.Error())
		//log.Errorf("GateSvr Login is fail. code[%d] txt[%s]", loginResult.Code, loginResult.Msg)
		return
	}
	serverer.bLogin=true
	log.Infof("网关服登入成功 serverid:%d", srcId)
}
func DoPlayerJionReq(srcType uint8,srcId uint64,s *scokets.Connector,buf []byte){
	pData := &base.NotifyJionServerReq{}
	if proto.Unmarshal(buf, pData) != nil {
		fmt.Println("DoPlayerJionReq 协议解析失败:")
		return
	}
	rSt := &base.NotifyJionServerResult{
		Userid: pData.Userid,
		Codeid: 0,
	}
	dPro, _ := proto.Marshal(rSt)
	s.SendData(send.CreateMsgToServerID(srcId, msg.MID_Gate, msg.Gate_SS_ClientJionResult, dPro))
}
func DoPlayerLeaveReq(srcType uint8,srcId uint64,s *scokets.Connector,buf []byte){
	pData := &base.NotifyLeaveServerReq{}
	if proto.Unmarshal(buf, pData) != nil {
		log.Errorln("DoPlayerLeave 协议解析失败:")
		return
	}
	rSt := &base.NotifyLeaveServerResult{
		Userid: pData.Userid,
		Codeid: 0,
	}
	dPro, _ := proto.Marshal(rSt)
	s.SendData(send.CreateMsgToServerID(srcId, msg.MID_Gate, msg.Gate_SS_ClientLeaveResult, dPro))

	return
}
func DoPlayerTestMsg(srcType uint8,srcId uint64,s *scokets.Connector,buf []byte){
	msgData := &base.TestMsg{}
	err := proto.Unmarshal(buf, msgData)
	if err != nil {
		log.Errorln("协议解析失败2:", err)
		return //当远程客户端连接发生错误（断开）后，终止此协程。
	}
	fmt.Printf("userid:%d,str:%s\n", srcId, msgData.Txt)

	//回复
	msgData = &base.TestMsg{
		Txt: fmt.Sprintf("收到测试消息,我是[%s]", constant.GetServerIDName(config.App.TID, config.App.SID)),
	}
	dPro, _ := proto.Marshal(msgData)
	//testmsg := msg.CreateWholeMsgData(msg.Sign_userid, userid, msg.MID_Test, msg.Test_1, dPro)
	s.SendData(send.CreateMsgToClient(srcId, msg.MID_Test, msg.Test_1, dPro))
}

