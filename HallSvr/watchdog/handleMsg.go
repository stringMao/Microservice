package watchdog

import (
	"Common/msg"
)

//处理网关服消息
func HandleGateSvrMessage(serverid uint64, mainid, sonid uint32, len uint32, data []byte) bool {
	switch mainid {
	case msg.MID_Gate:
		switch sonid {
		default:
		}
	case msg.MID_Err:

	default:
	}

	return false
}
