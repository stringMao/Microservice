package msg

//MainID
const (
	MID_Err uint32 = 1 //错误消息

	MID_Gate uint32 = 2 //网关服下发的消息
	MID_Hall uint32 = 3
)

//MID_Gate 网关服的子命令
const (
	Gate_SendPlayerData uint32 = 1 //下发玩家信息
)
