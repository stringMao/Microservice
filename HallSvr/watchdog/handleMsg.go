package watchdog

import (
	"Common/constant"
	"Common/log"
	"Common/msg"
	"Common/proto/base"
	"HallSvr/config"
	"HallSvr/core/send"
	"HallSvr/kernel"
	"fmt"

	"github.com/golang/protobuf/proto"
)

type ClientMsgData struct {
	UserId  uint64
	GateSvr uint64
	MsgHead *msg.HeadProto
	MsgData []byte
}
type SvrMsgData struct {
	ServerId uint64
	MsgHead  *msg.HeadProto
	MsgData  []byte
}
type SelfMsgData struct {
	MsgHead *msg.HeadProto
	MsgData []byte
}

var G_SvrMsg chan SvrMsgData = make(chan SvrMsgData, 1000) //服务器消息队列
var G_SvrMsgSwitch chan int = make(chan int, 1)            //开关

var G_ClientMsg chan ClientMsgData = make(chan ClientMsgData, 500) //玩家消息队列
var G_ClientMsgSwitch chan int = make(chan int, 1)

var G_SelfMsg chan SelfMsgData = make(chan SelfMsgData, 10) //自己的消息队列
var G_SelfMsgSwitch chan int = make(chan int, 1)            //开关

//消息总分发入口
func DistributeMessage(id uint64, data []byte, len uint32) {
	sign := msg.GetSign(data)
	switch sign.SignType {
	case msg.Sign_serverid: //后8位是serverid
		head := msg.GetHead(data)

		data := SvrMsgData{
			ServerId: sign.SignId,
			MsgHead:  head,
			MsgData:  data[msg.GetHeadLength():],
		}
		G_SvrMsg <- data
	case msg.Sign_userid: //后8位是userid
		head := msg.GetHead(data)

		data := ClientMsgData{
			UserId:  sign.SignId,
			GateSvr: id,
			MsgHead: head,
			MsgData: data[msg.GetHeadLength():],
		}
		G_ClientMsg <- data
	case msg.Sign_Self:
		head := msg.GetHead(data)

		data := SelfMsgData{
			MsgHead: head,
			MsgData: data[msg.GetHeadLength():],
		}
		G_SelfMsg <- data
	default:
	}
}

//服务器消息单独携程处理
func HandleSvrMsg() {
	for {
		select {
		case msg, ok := <-G_SvrMsg:
			if ok {
				OnEventSvrMessage(msg.ServerId, msg.MsgHead.MainID, msg.MsgHead.SonID, msg.MsgHead.Len, msg.MsgData)
			}
		case <-G_SvrMsgSwitch:
			goto stop
		}
	}
stop:
}

//客户端的消息单独携程处理
func HandlClientMsg() {
	for {
		select {
		case msg, ok := <-G_ClientMsg:
			if ok {
				OnEventClientMessage(msg.UserId, msg.GateSvr, msg.MsgHead.MainID, msg.MsgHead.SonID, msg.MsgHead.Len, msg.MsgData)
			}
		case <-G_ClientMsgSwitch:
			goto stop
		}
	}
stop:
}

//自己的消息单独携程处理
func HandleSelfMsg() {
	for {
		select {
		case msg, ok := <-G_SelfMsg:
			if ok {
				OnEventSelfMessage(msg.MsgHead.MainID, msg.MsgHead.SonID, msg.MsgHead.Len, msg.MsgData)
			}
		case <-G_SelfMsgSwitch:
			goto stop
		}
	}
stop:
}

//处理服务器消息服消息  serverid=来源svr
func OnEventSvrMessage(serverid uint64, mainid, sonid uint32, len uint32, data []byte) bool {
	log.Debugf("OnEventSvrMessage:serverid[%d],mainid[%d],sonid[%d]", serverid, mainid, sonid)
	switch mainid {
	case msg.MID_Gate:
		switch sonid {
		case msg.Gate_SS_SvrLoginResult: //网关服登入结果
			DoLoginGateSvr(serverid, len, data)
		case msg.Gate_SS_ClientJionReq: //有用户请求加入本服
			DoClientJionReq(serverid, len, data)
		case msg.Gate_SS_ClientOffline: //用户离线

		default:
		}
	case msg.MID_Err:

	default:
	}

	return false
}

//处理客户端消息
func OnEventClientMessage(userid uint64, gatesvrid uint64, mainid, sonid uint32, len uint32, data []byte) bool {
	log.Debugf("OnEventClientMessage:userid[%d],gatesvrid[%d],mainid[%d],sonid[%d]", userid, gatesvrid, mainid, sonid)
	switch mainid {
	case msg.MID_Test:
		switch sonid {
		case msg.Test_1:
			DoClientTestMessage(gatesvrid, userid, mainid, sonid, len, data)
		default:
		}
	case msg.MID_Err:

	default:
	}

	return false
}

func OnEventSelfMessage(mainid, sonid uint32, len uint32, data []byte) {
	log.Debugf("OnEventSelfMessage:mainid[%d],sonid[%d]", mainid, sonid)
}

func HandleErr(serverid uint64, codeid int) {
	if codeid == 1 { //断线
		kernel.GetManagerSvrs().DeleteGateSvr(serverid)
	}
}

//=具体业务消息逻辑处理===========================

//网关服登入结果
func DoLoginGateSvr(id uint64, len uint32, data []byte) bool {
	c, ok := m_tempConnMap[id]
	if !ok {
		return false
	}
	defer delete(m_tempConnMap, id)

	loginResult := &base.LoginResult{}
	err := proto.Unmarshal(data[:len], loginResult)
	if err != nil || loginResult.Code != 0 {
		c.CloseConnet()
		//log.Errorf("GateSvr Login is fail. err:%s", err.Error())
		//log.Errorf("GateSvr Login is fail. code[%d] txt[%s]", loginResult.Code, loginResult.Msg)
		return false
	}

	agent := kernel.NewSeverAgent(c)
	agent.BaseData.TID = loginResult.Tid
	agent.BaseData.SID = loginResult.Sid
	agent.BaseData.ServerId = msg.EncodeServerID(agent.BaseData.TID, agent.BaseData.SID)
	agent.BaseData.Name = loginResult.Name

	log.Infof("网关服登入成功 serverid:%d", agent.BaseData.ServerId)
	kernel.GetManagerSvrs().AddGateSvr(agent.BaseData.ServerId, agent)
	agent.StartWork()

	return true
}

func DoClientTestMessage(key uint64, userid uint64, mainid, sonid uint32, len uint32, data []byte) {

	msgData := &base.TestMsg{}
	err := proto.Unmarshal(data[:len], msgData)
	if err != nil {
		fmt.Println("协议解析失败2:", err)
		return //当远程客户端连接发生错误（断开）后，终止此协程。
	}
	fmt.Printf("userid:%d,mainid:%d,sonid:%d,str:%s\n", userid, mainid, sonid, msgData.Txt)

	//回复
	msgData = &base.TestMsg{
		Txt: fmt.Sprintf("收到测试消息,我是[%s]", constant.GetServerIDName(config.App.TID, config.App.SID)),
	}
	dPro, _ := proto.Marshal(msgData)
	//testmsg := msg.CreateWholeMsgData(msg.Sign_userid, userid, msg.MID_Test, msg.Test_1, dPro)
	kernel.GetManagerSvrs().SendData(key, send.CreateMsgToClient(userid, mainid, sonid, dPro))
	//SendMsgToClient(key, userid, mainid, sonid, dPro)
}

func DoClientJionReq(id uint64, len uint32, data []byte) bool {
	pData := &base.NotifyJionServerReq{}
	if proto.Unmarshal(data[:len], pData) != nil {
		fmt.Println("DoClientJionReq 协议解析失败:")
		return false
	}

	rSt := &base.NotifyJionServerResult{
		Userid: pData.Userid,
		Codeid: 0,
	}
	dPro, _ := proto.Marshal(rSt)
	kernel.GetManagerSvrs().SendData(id, send.CreateMsgToServerID(id, msg.MID_Gate, msg.Gate_SS_ClientJionResult, dPro))

	return true
}

func SendMsgToClient(key uint64, userid uint64, mainid, sonid uint32, data []byte) {
	// msgData := &base.TestMsg{
	// 	Txt: fmt.Sprintf("收到测试消息,我是[%s]", constant.GetServerIDName(config.App.TID, config.App.SID)),
	// }
	// dPro, _ := proto.Marshal(msgData)
	testmsg := msg.CreateWholeMsgData(msg.Sign_userid, userid, msg.MID_Test, msg.Test_1, data)

	kernel.GetManagerSvrs().SendData(key, testmsg)
}
