package msg

import (
	"Common/proto/base"

	"github.com/gogo/protobuf/proto"
)


//错误码
const (
	// ErrCode_Undefined  未定义
	ErrCode_Undefined  int32 = 1

	// ErrCode_LoginGateSvrAuthFail 网关服登入失败-身份验证未通过
	ErrCode_LoginGateSvrAuthFail int32=1000

)

//错误消息的子命令定义
//MainID_Err => SonID
const (
	//一般错误
	Err_Nomoal uint32 = 1

	//身份认证失败
	Err_Login_AuthenticationFail uint32 = 1000
	Err_Login_InitDataFail       uint32 = 1001
	Err_Login_ReplaceClientFail  uint32 = 1002

	//
	//Err_ServerNoFind 服务未找到
	Err_ServerNoFind uint32 = 1100
	//Err_MsgSendFail 消息转发失败
	Err_MsgSendFail uint32 = 1101
	//Err_MsgSendFail_ServerNoExist 消息发送失败，转发的服务器不存在
	Err_MsgSendFail_ServerNoExist uint32 = 1102
)

//CreateErrorMsg 创建一条错误码消息
func CreateErrorMsg(errid uint32) []byte {
	head := &HeadProto{
		MainID: MID_Err,
		SonID:  errid,
		Len:    0,
	}
	return head.Encode()
}

//创建一条发给服务器的错误消息
func CreateErrSvrMsg(serverid uint64, errid uint32) []byte {
	return append(CreateSignHead(Sign_serverid, serverid), CreateErrorMsg(errid)...)
}

//创建 带有错误描述的消息
func CreateErrorMsgData(errid uint32, txt string) []byte {
	pro := &base.Txt{
		Txt: txt,
	}
	pData, _ := proto.Marshal(pro)
	return CreateWholeProtoData(MID_Err, errid, pData)
}

func CreateErrSvrMsgData(serverid uint64, errid uint32, txt string) []byte {
	return append(CreateSignHead(Sign_serverid, serverid), CreateErrorMsgData(errid, txt)...)
}

// func CreateErrorMsgData(errid uint32, m string) {
// 	data := []byte(m)
// 	head := &HeadProto{
// 		MainID: MID_Err,
// 		SonID:  errid,
// 		Len: func(d []byte) uint32 {
// 			if d == nil {
// 				return 0
// 			} else {
// 				return uint32(len(d))
// 			}
// 		}(data),
// 	}
// }

//func CreateErrMsg
