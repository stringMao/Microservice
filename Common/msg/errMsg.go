package msg

//MainID_Err => SonID
const (
	//Err_ServerNoFind 服务未找到
	Err_ServerNoFind uint32 = 1
	//Err_MsgSendFail 消息转发失败
	Err_MsgSendFail uint32 = 2
	//Err_MsgSendFail_ServerNoExist 消息发送失败，转发的服务器不存在
	Err_MsgSendFail_ServerNoExist = 3
)

//CreateErrorMsg 创建一条错误码消息
func CreateErrorMsg(errid uint32) []byte {
	head := &HeadProto{
		MainID: MainID_Err,
		SonID:  errid,
		Len:    0,
	}
	return head.Encode()
}
