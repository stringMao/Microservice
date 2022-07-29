package main

import (
	"Common/constant"
	"Common/kernel/go-scoket/sockets"
	"Common/log"
	"Common/msg"
	"Common/proto/base"
	"Common/proto/codes"
	"Common/proto/gateProto"
	"TestClient/agent"
	"TestClient/loginsvr"
	"github.com/golang/protobuf/proto"
	"os"
	"time"
)


func main() {
	userid, token, ip := loginsvr.Signin()

	connector:=sockets.NewConnector(ip, nil,time.Second)
	if !connector.StartConnect(){
		log.Errorln("gate 连接 失败")
		return
	}
	logicAgent:=agent.NewLogicAgent(connector)
	connector.Engine.SetHandler(logicAgent)
	//添加业务消息处理函数
	connector.AddHandleFuc(msg.ToUser_GateLoginResult,GateJoinResult)
	connector.AddHandleFuc(msg.ToUser_JoinSvrResult,SvrJoinResult)
	connector.AddHandleFuc(msg.ToUser_QuitSvrResult,SvrQuitResult)
	//connector.AddHandleFuc(msg.ToUser_QuitSvrResult,GateJoinResult)
	connector.AddHandleFuc(msg.ToUser_Test,SvrTestResult)

	//发送登入消息
	pObj := &gateProto.UserLoginReq{
		Userid: uint64(userid),
		Token:  token,
	}
	connector.SendData(msg.NewClientMessage(msg.ToGateSvr_UserLogin,pObj))

	//
	c := make(chan os.Signal)
	<-c
}

func GateJoinResult(engine *sockets.Engine,buf []byte)  {
	pData:=&gateProto.UserLoginResult{}
	err := proto.Unmarshal(buf, pData)
	if err != nil {
		return
	}
	if pData.Code==codes.Code_Success{
		log.Debugf("gate login success ")
		//登入大厅请求
		pObj := &gateProto.JoinQuitServerReq{
			Tid:    constant.TID_HallSvr,
		}
		engine.SendData(msg.NewClientMessage(msg.ToGateSvr_UserJoinSvrReq,pObj))
		
	}else{
		log.Debugf("gate login fail code[%d]",pData.Code)
	}
}

func SvrJoinResult(engine *sockets.Engine,buf []byte)  {
	pData:=&gateProto.JoinQuitServerResult{}
	err := proto.Unmarshal(buf, pData)
	if err != nil {
		return
	}
	if pData.Code!=codes.Code_Success{
		log.Debugf("svr login fail code[%d]",pData.Code)
		return
	}
	log.Debugf("svr[%d] login success ",pData.Tid)

	//
	//登入大厅请求
	//pObj := &gateProto.JoinQuitServerReq{
	//	Tid:    constant.TID_HallSvr,
	//}
	//engine.SendData(msg.NewClientMessage(msg.ToGateSvr_UserQuitSvrReq,pObj))
    i:=0
	for i<1000  {
		i++
		pObj := &base.ReplyResult{
			Code: int32(i),
			Txt:  "测试消息",
		}
		engine.SendData(msg.NewClientMessage(msg.ToHallSvr_Test,pObj))
		time.Sleep(10*time.Microsecond)
	}


}

func SvrQuitResult(engine *sockets.Engine,buf []byte)  {
	pData:=&gateProto.JoinQuitServerResult{}
	err := proto.Unmarshal(buf, pData)
	if err != nil {
		return
	}
	if pData.Code!=codes.Code_Success{
		log.Debugf("svr Quit fail code[%d]",pData.Code)
		return
	}
	log.Debugf("svr[%d] Quit success ",pData.Tid)
}

func SvrTestResult(engine *sockets.Engine,buf []byte){

	pData:=&base.ReplyResult{}
	err := proto.Unmarshal(buf, pData)
	if err != nil {
		return
	}
	log.Debugf("TestMessage； %s:%d",pData.Txt,pData.Code)
	return
	pData.Code++
	engine.SendData(msg.NewClientMessage(msg.ToHallSvr_Test,pData))
}

//func protoEncodePrint(buf []byte, n uint32) {
//	head := &msg.HeadProto{}
//	head.Decode(buf)
//
//	fmt.Printf("接收消息: mainid[%d] sonid[%d] len[%d] \n", head.MainID, head.SonID, head.Len)
//	if head.Len > 0 {
//
//		switch head.MainID {
//		case msg.MID_Err:
//		case msg.MID_Test:
//			switch head.SonID {
//			case msg.Test_1:
//				msgData := &base.TestMsg{}
//				err := proto.Unmarshal(buf[msg.GetProtoHeadLength():n], msgData)
//				if err != nil {
//					fmt.Println("协议解析失败2:", err)
//					return //当远程客户端连接发生错误（断开）后，终止此协程。
//				}
//				fmt.Printf("mainid:%d,sonid:%d,str:%s\n", head.MainID, head.SonID, msgData.Txt)
//			default:
//			}
//		case msg.MID_Hall:
//		case msg.MID_Gate:
//			switch head.SonID {
//			case msg.Gate_SC_SendPlayerData:
//				msgstr := &gatesvrproto.PlayerInfo{}
//				err := proto.Unmarshal(buf[head.GetHeadLen():n], msgstr)
//				if err != nil {
//					fmt.Println("协议解析失败2:", err)
//					return //当远程客户端连接发生错误（断开）后，终止此协程。
//				}
//				fmt.Printf("%+v\n", msgstr)
//				//ConnSucc = true
//			}
//
//		}
//
//	}
//	//fmt.Printf("接收消息: mainid[%d] sonid[%d] len[%d] msg:%s \n", head.MainID, head.SonID, head.Len, msgstr.Txt)
//	return
//
//}
