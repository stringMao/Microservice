package msg

import "Common/constant"

//MainID定义
const (
	MID_Err   uint32 = 1  //错误消息
	MID_Fatal uint32 = 2  //异常消息
	MID_Test  uint32 = 10 //测试用消息

	MID_Gate uint32 = 100 //网关服下发的消息
	MID_Hall uint32 = 101 //
)


//每个服务分配的消息id区间 [10000*TID,10000*TID+9999]
const msgInterval = 10000
func ParseMessageID(id uint32)int{
	return int(id/msgInterval)
}
//玩家消息[0,10000)
const(
	ToUser uint32=iota+1
	ToUser_Err
	ToUser_Test
	ToUser_GateLoginResult
	ToUser_JoinSvrResult
	ToUser_QuitSvrResult
)

//公共消息[10000,20000)
const (
	CommonSvrMsg  uint32=iota+constant.TID_Common*msgInterval
	CommonSvrMsg_SvrRegisterResult //服务器像网关服注册结果
	CommonSvrMsg_UserJoin      //用户请求加入你
	CommonSvrMsg_UserQuit      //用户请求离开你
	CommonSvrMsg_UserOffline   //用户离线
)

//TID_GateSvr  = 11
//发送给网关服的消息 [110000,120000)
const (
	ToGateSvr 	uint32=iota+constant.TID_GateSvr*msgInterval
	ToGateSvr_UserLogin 	   //登入网关服
	ToGateSvr_SvrRegister 			//服务器向网关注册
	ToGateSvr_UserJoinSvrReq  		//用户加入服务器的请求
	ToGateSvr_UserJoinSvrResult
	ToGateSvr_UserQuitSvrReq  		//用户退出服务器的请求
	ToGateSvr_UserQuitSvrResult
	ToGateSvr_UserKickOut       	//将用户强制踢出
)

//TID_HallSvr  = 12 //大厅服tid
// 发送给大厅服的消息 [120000,130000)
const (
	ToHallSvr  uint32=iota+constant.TID_HallSvr*msgInterval
	ToHallSvr_Test
)

