package send

import (
	"Common/msg"
	"GateSvr/config"
)

func CreateMsgToSvr(sonid uint32, data []byte) []byte {
	return msg.CreateWholeMsgData(msg.Sign_serverid, config.App.ServerID, msg.MID_Gate, sonid, data)
}

func CreateMsgToClient(sonid uint32, data []byte) []byte {
	return msg.CreateWholeProtoData(msg.MID_Gate, sonid, data)
}
