package send

import (
	"Common/msg"
)

//给某个服务器发送消息
func CreateMsgToServerID(targetSvrid uint64, mainid, sonid uint32, data []byte) []byte {
	return msg.CreateWholeMsgData(msg.Sign_serverid, targetSvrid, mainid, sonid, data)
}
func CreateMsgToTID(targetTid uint32, mainid, sonid uint32, data []byte) []byte {
	return msg.CreateWholeMsgData(msg.Sign_serverid, msg.EncodeServerID(targetTid, 0), mainid, sonid, data)
}

//给客户端的消息
func CreateMsgToClient(userid uint64, mainid, sonid uint32, data []byte) []byte {
	return msg.CreateWholeMsgData(msg.Sign_userid, userid, mainid, sonid, data)
}
