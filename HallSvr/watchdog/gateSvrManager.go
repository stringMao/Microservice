package watchdog

//管家
import (
	"Common/constant"
	"Common/log"
	"Common/msg"
	"Common/proto/base"
	"Common/svrfind"
	"HallSvr/config"
	"HallSvr/kernel"
	"fmt"
	"strconv"
	"time"

	"github.com/golang/protobuf/proto"
)

//连接的临时存储容器
var m_tempConnMap map[uint64]*kernel.ConnTcp = make(map[uint64]*kernel.ConnTcp, 10)

//bool
func ConnectGateSvrs() {
	//登入验证
	logindata := &base.ServerLogin{
		Tid:      uint32(config.App.TID),
		Sid:      uint32(config.App.SID),
		Password: "test",
	}
	for {
		gatelist := svrfind.G_ServerRegister.GetSvr(constant.GetServerName(constant.TID_GateSvr), constant.GetServerTag(constant.TID_GateSvr))
		for _, v := range gatelist {
			//fmt.Println(k, "  ", v.Service.Address)
			//fmt.Println(k, "  ", v.Service.Port)
			//fmt.Printf(" %d:%+v \n", k, v.Service)

			if addr, ok := v.Service.TaggedAddresses["server"]; ok {
				strserverid, ok := v.Service.Meta["ServerID"]
				if !ok {
					log.Errorf("网关服未正确注册 serverid")
					continue
				}

				serverid, err := strconv.ParseUint(strserverid, 10, 64) //strconv.Atoi(strserverid)
				if err != nil {
					log.Errorf("网关服未正确注册2 serverid ")
					continue
				}

				if kernel.GetManagerSvrs().IsExist(serverid) {
					continue
				}

				pData, _ := proto.Marshal(logindata)
				conn := kernel.NewConnTcp(serverid, fmt.Sprintf("%s:%d", addr.Address, addr.Port), pData, DistributeMessage, HandleErr)
				if conn.Connect() {
					log.Infof("网关服连接[addr:%s]成功", conn.Addr)
					m_tempConnMap[uint64(serverid)] = conn
				}
			}
		}
		time.Sleep(10 * time.Second)
	}
}

func HandleErr(serverid uint64, codeid int) {
	if codeid == 1 { //断线
		kernel.GetManagerSvrs().DeleteGateSvr(serverid)
	}
}

//消息总分发入口
func DistributeMessage(id uint64, data []byte, len uint32) {
	signhead := &msg.HeadSign{}
	signhead.Decode(data)

	switch signhead.SignType {
	case msg.Sign_serverid: //后8位是serverid
		head := &msg.HeadProto{}
		head.Decode(data[msg.GetSignHeadLength():])

		data := SvrMsgData{
			ServerId: signhead.SignId,
			MsgHead:  head,
			MsgData:  data[msg.GetHeadLength():],
		}
		G_SvrMsg <- data

		// switch signhead.Tid {
		// case constant.TID_GateSvr:
		// 	if head.MainID == msg.MID_Gate && head.SonID == msg.GateSvr_SvrLoginResult {
		// 		OnEventLoginGateSvr(c, len, data[msg.GetHeadLength():])
		// 	} else {
		// 		HandleGateSvrMessage(signhead.SignId, head.MainID, head.SonID, head.Len, data[msg.GetHeadLength():])
		// 	}

		// case constant.TID_LoginSvr:
		// case constant.TID_HallSvr:
		// default:
		// }
	case msg.Sign_userid: //后8位是userid
		head := &msg.HeadProto{}
		head.Decode(data[msg.GetSignHeadLength():])

		data := ClientMsgData{
			UserId:  signhead.SignId,
			GateSvr: id,
			MsgHead: head,
			MsgData: data[msg.GetHeadLength():],
		}
		G_ClientMsg <- data

		//HandleClientMessage(key, signhead.SignId, head.MainID, head.SonID, head.Len, data[msg.GetHeadLength():])
	case msg.Sign_Self:
		head := &msg.HeadProto{}
		head.Decode(data[msg.GetSignHeadLength():])

		data := SelfMsgData{
			MsgHead: head,
			MsgData: data[msg.GetHeadLength():],
		}
		G_SelfMsg <- data
	default:
	}
}

// //网关服登入结果
// func OnEventLoginGateSvr(c *kernel.ConnTcp, len uint32, data []byte) bool {
// 	loginResult := &base.LoginResult{}
// 	err := proto.Unmarshal(data, loginResult)
// 	if err != nil {
// 		c.CloseConnet()
// 		log.Errorf("GateSvr Login is fail. err:%s", err.Error())
// 		return false
// 	}
// 	if loginResult.Code != 0 {
// 		c.CloseConnet()
// 		log.Errorf("GateSvr Login is fail. code[%d] txt[%s]", loginResult.Code, loginResult.Msg)
// 		return false
// 	}
// 	agent := kernel.NewSeverAgent(c)
// 	agent.BaseData.TID = loginResult.Tid
// 	agent.BaseData.SID = loginResult.Sid
// 	agent.BaseData.ServerId = msg.EncodeServerID(agent.BaseData.TID, agent.BaseData.SID)
// 	agent.BaseData.Name = loginResult.Name
// 	agent.StartWork()
// 	log.Infof("网关服登入成功 serverid:%d", agent.BaseData.ServerId)
// 	kernel.GetManagerSvrs().AddGateSvr(agent.BaseData.ServerId, agent)

// 	return true
// }

func HandleClientMessage(key uint64, userid uint64, mainid, sonid uint32, len uint32, data []byte) {

	msgData := &base.TestMsg{}
	err := proto.Unmarshal(data[:len], msgData)
	if err != nil {
		fmt.Println("协议解析失败2:", err)
		return //当远程客户端连接发生错误（断开）后，终止此协程。
	}
	fmt.Printf("userid:%d,mainid:%d,sonid:%d,str:%s\n", userid, mainid, sonid, msgData.Txt)
	//回复
	SendMsgToClient(key, userid, mainid, sonid, 0, nil)
}

func SendMsgToClient(key uint64, userid uint64, mainid, sonid uint32, len uint32, data []byte) {
	msgData := &base.TestMsg{
		Txt: fmt.Sprintf("收到测试消息,我是[%s]", constant.GetServerIDName(config.App.TID, config.App.SID)),
	}
	dPro, _ := proto.Marshal(msgData)
	testmsg := msg.CreateWholeMsgData(msg.Sign_userid, userid, msg.MID_Test, msg.Test_1, dPro)

	kernel.GetManagerSvrs().SendData(key, testmsg)
}
