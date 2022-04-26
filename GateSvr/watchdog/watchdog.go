package watchdog

import (
	"Common/constant"
	"Common/kernel/go-scoket/scokets"
	"Common/msg"
	"Common/proto/base"
	"GateSvr/config"
	"GateSvr/core/send"
	"GateSvr/util/msgbody"
	"fmt"
	"github.com/golang/protobuf/proto"
)

var G_ClientManager *ClientManager=nil
//代理管理对象创建
var ServerListener *scokets.Listener =nil
var PlayerListener *scokets.Listener =nil

//业务启动
func Start() {
	//创建代理管理器
	G_ClientManager=NewClientManager()

	scokets.InitBytePool(1000,1024,1024)

	ServerListener =scokets.NewListener(scokets.Type_Scoket, fmt.Sprintf("%s:%d",config.App.WebManagerIP,config.App.ServerPort) ,GetServerConnection)
	ServerListener.AddHandleFuc(msg.Gate_SS_ClientJionResult,JionServerResult)
	ServerListener.AddHandleFuc(msg.SCS_SvrRegisterGateSvr,SvrRegister)
	ServerListener.StartListen()

	//PlayerListener =scokets.NewListener(scokets.Type_Scoket,fmt.Sprintf("%s:%d",config.App.WebManagerIP,config.App.ClientPort),GetPlayerConnection)
	//PlayerListener.AddHandleFuc(msg.Gate_CS_JionServerReq,ReqJionServer)
	//PlayerListener.AddHandleFuc(msg.Gate_CS_LeaveServerReq,ReqLeaveServer)
	//PlayerListener.StartListen()
}

//====================================================================================================
//玩家消息
func ReqJionServer(client *scokets.Client,buf []byte){
	obj := &base.ClientJionServerReq{}
	err := proto.Unmarshal(buf, obj)
	if err != nil {
		return
	}
	if obj.Tid == constant.TID_GateSvr {
		return
	}

	pd:=client.Data.(PlayerData)
	if _,has:=pd.SvrList[obj.Tid];has{
		//已经加入该type的服务器
		client.Session.Send(send.CreateMsgToClient(msg.Gate_SC_ClientJionResult,
			msgbody.MakeToClientJionServerResult(2, uint64(obj.Tid))))
		return
	}
	serverid :=G_ClientManager.AllocSvr(obj.Tid)
	if serverid == 0 {
		//服务器找不到
		client.Session.Send(send.CreateMsgToClient(msg.Gate_SC_ClientJionResult,
			msgbody.MakeToClientJionServerResult(1, msg.EncodeServerID(obj.Tid, 0))))
		return
	}
	//转发加入请求给指定服务器
	tPro := &base.NotifyJionServerReq{Userid: pd.UserId}
	dPro, _ := proto.Marshal(tPro)
	G_ClientManager.SendToServer(serverid, send.CreateMsgToSvr(msg.Gate_SS_ClientJionReq, dPro))
}

func ReqLeaveServer(client *scokets.Client,buf []byte){
	obj := &base.ClientLeaveServerReq{}
	err := proto.Unmarshal(buf, obj)
	if err != nil {
		return
	}
	if obj.Tid == constant.TID_GateSvr {
		return
	}

	pd:=client.Data.(PlayerData)
	if sid,has:=pd.SvrList[obj.Tid];has{
		serverid:=msg.EncodeServerID(obj.Tid, sid)
		//转发加入请求给指定服务器
		tPro := &base.NotifyLeaveServerReq{Userid: pd.UserId}
		dPro, _ := proto.Marshal(tPro)
		G_ClientManager.SendToServer(serverid, send.CreateMsgToSvr(msg.Gate_SS_ClientLeaveReq, dPro))
	}

}

//服务器消息
func JionServerResult(client *scokets.Client,len int,buf []byte){
	obj := &base.NotifyJionServerResult{}
	if err := proto.Unmarshal(buf[:len], obj); err != nil {
		//协议解析错误
		return
	}
	cdata:=client.Data.(ServerData)
	if obj.Codeid != 0 {
		//加入失败
		G_ClientManager.SendToPlayer(obj.Userid, send.CreateMsgToClient(msg.Gate_SC_ClientJionResult,
			msgbody.MakeToClientJionServerResult(int(obj.Codeid), cdata.Serverid)))
		return
	}
	player:=G_ClientManager.GetPlayer(obj.Userid)
	if player==nil{
		return
	}
	
	//保存jion结果
	player.JionServer(cdata.Tid,cdata.Sid)

	//成功通知客户端
	G_ClientManager.SendToPlayer(obj.Userid,send.CreateMsgToClient(msg.Gate_SC_ClientJionResult,
		msgbody.MakeToClientJionServerResult(int(obj.Codeid), cdata.Serverid)))

}
