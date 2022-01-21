package main

import (
	"Common/log"
	"Common/msg"
	"Common/proto/base"
	"Common/proto/gatesvrproto"
	"TestClient/kernel"
	"TestClient/loginsvr"
	"fmt"
	"os"
	"time"

	"github.com/golang/protobuf/proto"
)

var G_GateSvr *kernel.ConnTcp = nil

var bjionHall = false

func main() {

	fmt.Println("111")
	userid, token, ip := loginsvr.Signin()

	//登入消息
	stSend := &base.ClientLogin{
		Userid: uint64(userid),
		Token:  token,
	}
	pData, _ := proto.Marshal(stSend)

	G_GateSvr = kernel.NewConnTcp(1, ip, pData, DistributeMessage, HandleErr)
	if !G_GateSvr.Connect() {
		return
	}
	log.Infof("网关服连接[addr:%s]成功", G_GateSvr.Addr)
	G_GateSvr.HeartSwitch = true
	//gate := gatesvr.NewConnTcp(ip)
	//go gate.Connect(userid, token)

	//===========
	jionhalldata := &base.ClientJionServerReq{
		Tid: 3,
	}
	djionhalldata, _ := proto.Marshal(jionhalldata)
	pjionhalldata := msg.CreateWholeMsgData(msg.Sign_serverid, 2, msg.MID_Gate, msg.Gate_CS_JionServerReq, djionhalldata)
	G_GateSvr.SendData(pjionhalldata)
	//===================

	go func() {
		for {
			if bjionHall {
				time.Sleep(time.Second * 1)
				//if gatesvr.ConnSucc {
				msgData := &base.TestMsg{
					Txt: "这是一条测试消息",
				}
				dPro, _ := proto.Marshal(msgData)
				testmsg := msg.CreateWholeMsgData(msg.Sign_serverid, 3, msg.MID_Test, msg.Test_1, dPro)
				//gate.TestSend(testmsg)
				G_GateSvr.SendData(testmsg)
				//gate.TestSend(msg.CreateWholeMsgData(msg.Sign_serverid, 0, msg.MID_Gate, 1, []byte("twes")))
				//}
			}
		}
	}()

	c := make(chan os.Signal)
	<-c
}

func HandleErr(serverid uint64, codeid int) {
	if codeid == 1 { //断线
		//kernel.GetManagerSvrs().DeleteGateSvr(serverid)
		//c <- os.Interrupt{}
	}
}

//消息总分发入口
func DistributeMessage(id uint64, data []byte, len uint32) {
	protoEncodePrint(data, len)

	head := &msg.HeadProto{}
	head.Decode(data)
	switch head.MainID {
	case msg.MID_Gate:
		switch head.SonID {
		case msg.Gate_SC_ClientJionResult: //
			jionResult := &base.ToClientJionServerResult{}
			if err := proto.Unmarshal(data[msg.GetProtoHeadLength():len], jionResult); err != nil {
				log.Debug("协议解析失败")
				return
			}
			if jionResult.Codeid == 0 {
				bjionHall = true
			} else {
				log.Debug("jion hall fail")
			}

		default:
		}
	default:
	}

}

func protoEncodePrint(buf []byte, n uint32) {
	head := &msg.HeadProto{}
	head.Decode(buf)

	fmt.Printf("接收消息: mainid[%d] sonid[%d] len[%d] \n", head.MainID, head.SonID, head.Len)
	if head.Len > 0 {

		switch head.MainID {
		case msg.MID_Err:
		case msg.MID_Test:
			switch head.SonID {
			case msg.Test_1:
				msgData := &base.TestMsg{}
				err := proto.Unmarshal(buf[msg.GetProtoHeadLength():n], msgData)
				if err != nil {
					fmt.Println("协议解析失败2:", err)
					return //当远程客户端连接发生错误（断开）后，终止此协程。
				}
				fmt.Printf("mainid:%d,sonid:%d,str:%s\n", head.MainID, head.SonID, msgData.Txt)
			default:
			}
		case msg.MID_Hall:
		case msg.MID_Gate:
			switch head.SonID {
			case msg.Gate_SC_SendPlayerData:
				msgstr := &gatesvrproto.PlayerInfo{}
				err := proto.Unmarshal(buf[head.GetHeadLen():n], msgstr)
				if err != nil {
					fmt.Println("协议解析失败2:", err)
					return //当远程客户端连接发生错误（断开）后，终止此协程。
				}
				fmt.Printf("%+v\n", msgstr)
				//ConnSucc = true
			}

		}

	}
	//fmt.Printf("接收消息: mainid[%d] sonid[%d] len[%d] msg:%s \n", head.MainID, head.SonID, head.Len, msgstr.Txt)
	return

}
