package msg

//网关服使用的消息=========================================

const (
	//“网关服”与“其他服务器”的消息
	SS_SvrRegisterGateSvr uint32 = iota + 10000  //服务器像gate注册自己
	CS_PlayerLoginGateSvr    //玩家登入网关
	SC_PlayerLoginGateSvr    //网关返回用户登入结果
	Gate_SS_SvrLoginResult  //服务器登入网关服的结果
	Gate_SS_ClientJionReq                        //用户加入请求
	Gate_SS_ClientJionResult
	Gate_SS_ClientLeaveReq    //用户离开请求
	Gate_SS_ClientLeaveResult //用户离开结果
	Gate_SS_ClientOffline     //用户离线

	//网关服发送给客户端的消息
	Gate_SC_SendPlayerData uint32 = iota + 20000 //网关服下发玩家的信息
	Gate_SC_ClientJionResult
	Gate_SC_ClientLeaveResult

	//客户端发送给网关的消息
	Gate_CS_JionServerReq  uint32 = 30000 //客户端请求加入服务器
	Gate_CS_LeaveServerReq uint32 = 30001 //客户端请求退出服务器
)

//创建
