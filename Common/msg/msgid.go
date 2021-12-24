package msg

//MainID定义
const (
	MID_Err   uint32 = 1  //错误消息
	MID_Fatal uint32 = 2  //异常消息
	MID_Test  uint32 = 10 //测试用消息

	MID_Gate uint32 = 100 //网关服下发的消息
	MID_Hall uint32 = 101 //
)

const (
	Test_1 uint32 = 1
)

//MID_Gate 网关服的子命令
const (
	Gate_SendPlayerData uint32 = 1 //下发玩家信息
)

const (
	Hall_TestMsg uint32 = 1 //下发玩家信息
)
