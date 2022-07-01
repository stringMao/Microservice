package watchdog

import (
	"Common/constant"
	"Common/kernel/go-scoket/sockets"
	"Common/log"
	"Common/msg"
	"Common/proto/codes"
	"Common/proto/gateProto"
	"GateSvr/config"
	"fmt"
	"github.com/golang/protobuf/proto"
)

var G_ClientManager *ClientManager=nil
//代理管理对象创建
var ServerListener *sockets.Listener =nil
var PlayerListener *sockets.Listener =nil



//业务启动
func Start() {
	//创建代理管理器
	G_ClientManager=NewClientManager()

	ServerListener = sockets.NewListener(sockets.Type_Socket, fmt.Sprintf("%s:%d",config.App.WebManagerIP,config.App.ServerPort) ,GetServerConnection)
	ServerListener.AddHandleFuc(msg.ToGateSvr_UserJoinSvrResult, JoinServerResult)
	ServerListener.AddHandleFuc(msg.ToGateSvr_UserQuitSvrResult, QuitServerResult)
	ServerListener.AddHandleFuc(msg.ToGateSvr_SvrRegister,SvrRegister)
	ServerListener.StartListen()

	PlayerListener = sockets.NewListener(sockets.Type_Socket,fmt.Sprintf("%s:%d",config.App.WebManagerIP,config.App.ClientPort),GetPlayerConnection)
	PlayerListener.AddHandleFuc(msg.ToGateSvr_UserLogin,UserLogin)
	PlayerListener.AddHandleFuc(msg.ToGateSvr_UserJoinSvrReq, UserJoinSvrReq)
	PlayerListener.AddHandleFuc(msg.ToGateSvr_UserQuitSvrReq, UserQuitSvrReq)
	PlayerListener.StartListen()
}




//====================================================================================================

// UserJoinSvrReq 玩家消息
func UserJoinSvrReq(engine *sockets.Engine,buf []byte){
	pData := &gateProto.JoinQuitServerReq{}
	err := proto.Unmarshal(buf, pData)
	if err != nil {
		return
	}
	if pData.Tid == constant.TID_GateSvr {
		return
	}
	pUserData:= engine.Data.(*PlayerData)
	log.Debugf("UserJoinSvrReq uid[%d] tid[%d]",pUserData.UserId,pData.Tid)
	if _,has:=pUserData.SvrList[pData.Tid];has{
		//已经加入该type的服务器
		pObj:=&gateProto.JoinQuitServerResult{
			Code:   codes.Code_JoinSvr_AgainJoin, //0成功 1服务器找不到 2重复加入
			Tid:  	pData.Tid,
		}
		engine.SendData(msg.NewClientMessage(msg.ToUser_JoinSvrResult,pObj))
		return
	}
	serverId :=G_ClientManager.AllocSvr(pData.Tid)
	if serverId == 0 {
		//服务器找不到
		pObj:=&gateProto.JoinQuitServerResult{
			Code:   codes.Code_JoinSvr_SvrNoFind,
			Tid:  pData.Tid,
		}
		engine.SendData(msg.NewClientMessage(msg.ToUser_JoinSvrResult,pObj))
		return
	}

	pObj := &gateProto.UserJoinQuit{Userid: pUserData.UserId}
	//发送加入请求给指定服务器
	G_ClientManager.SendToServer2(serverId,
		msg.NewMessage(msg.CommonSvrMsg_UserJoin,0,config.App.ServerID,pObj))

}

func UserQuitSvrReq(engine *sockets.Engine,buf []byte){
	pData := &gateProto.JoinQuitServerReq{}
	err := proto.Unmarshal(buf, pData)
	if err != nil {
		return
	}
	if pData.Tid == constant.TID_GateSvr {
		return
	}

	pUserData:= engine.Data.(*PlayerData)
	log.Debugf("UserQuitSvrReq uid[%d] tid[%d]",pUserData.UserId,pData.Tid)
	if sid,has:=pUserData.SvrList[pData.Tid];has{
		serverId:=msg.EncodeServerID(pData.Tid, sid)
		//转发离开请求给指定服务器
		pObj := &gateProto.UserJoinQuit{Userid: pUserData.UserId}
		G_ClientManager.SendToServer2(serverId,
			msg.NewMessage(msg.CommonSvrMsg_UserQuit,0,config.App.ServerID,pObj))
	}else{
		//玩家本来就没有加入该服务，那就直接告知退出成功即可
		pObj:=&gateProto.JoinQuitServerResult{
			Code:   codes.Code_Success, //0成功
			Tid:  	pData.Tid,
		}
		engine.SendData(msg.NewClientMessage(msg.ToUser_QuitSvrResult,pObj))
	}

}

// JoinServerResult 服务器返回的用户join结果消息
func JoinServerResult(engine *sockets.Engine,buf []byte){
	pData := &gateProto.UserJoinQuitResult{}
	if err := proto.Unmarshal(buf, pData); err != nil {
		//协议解析错误
		return
	}
	pServerData:= engine.Data.(*ServerData)
	log.Debugf("JoinServerResult code[%d],tid[%d],sid[%d],uid[%d]",pData.Code,pServerData.Tid,pServerData.Sid,pData.Userid)
	if pData.Code== codes.Code_Success{
		pPlayerEngine:=G_ClientManager.GetPlayer(pData.Userid)
		if pPlayerEngine==nil{
			return
		}
		//保存join结果
		data:=pPlayerEngine.Data.(*PlayerData)
		if data!=nil && data.SvrList!=nil {
			data.SvrList[pServerData.Tid] = pServerData.Sid
		}
	}
	//成功通知客户端
	pObj := &gateProto.JoinQuitServerResult{
		Code:   pData.Code, //0成功
		Tid: 	pServerData.Tid,
	}
	G_ClientManager.SendToPlayer2(pData.Userid,msg.NewClientMessage(msg.ToUser_JoinSvrResult,pObj))
}

// QuitServerResult 服务器返回的用户Quit结果消息
func QuitServerResult(engine *sockets.Engine,buf []byte){
	pData := &gateProto.UserJoinQuitResult{}
	if err := proto.Unmarshal(buf, pData); err != nil {
		//协议解析错误
		return
	}
	pServerData:= engine.Data.(*ServerData)
	log.Debugf("QuitServerResult code[%d],tid[%d],sid[%d],uid[%d]",pData.Code,pServerData.Tid,pServerData.Sid,pData.Userid)
	if pData.Code== codes.Code_Success{
		pPlayerEngine:=G_ClientManager.GetPlayer(pData.Userid)
		if pPlayerEngine==nil{
			return
		}

		//保存quit结果
		data:=pPlayerEngine.Data.(*PlayerData)
		if data!=nil && data.SvrList!=nil {
			delete(data.SvrList, pServerData.Tid)
		}

	}
	//通知客户端
	pObj := &gateProto.JoinQuitServerResult{
		Code:   pData.Code, //0成功 1服务器找不到 2重复加入
		Tid: pServerData.Tid,
	}
	G_ClientManager.SendToPlayer2(pData.Userid,msg.NewClientMessage(msg.ToUser_QuitSvrResult,pObj))
}