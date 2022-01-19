package watchdog

import (
	"Common/log"
	"Common/msg"
	"Common/proto/base"
	"HallSvr/kernel"

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

//处理服务器消息服消息
func OnEventSvrMessage(serverid uint64, mainid, sonid uint32, len uint32, data []byte) bool {
	log.Debugf("OnEventSvrMessage:serverid[%d],mainid[%d],sonid[%d]", serverid, mainid, sonid)
	switch mainid {
	case msg.MID_Gate:
		switch sonid {
		case msg.GateSvr_SvrLoginResult: //网关服登入结果
			OnEventLoginGateSvr(serverid, len, data)
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
			HandleClientMessage(gatesvrid, userid, mainid, sonid, len, data)
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

//=具体业务消息逻辑处理===========================

//网关服登入结果
func OnEventLoginGateSvr(id uint64, len uint32, data []byte) bool {
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
