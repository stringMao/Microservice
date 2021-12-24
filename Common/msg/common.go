package msg

//CreateToClientMsg 生成一条发送给客户端得消息
func CreateToClientMsg(userid uint64, mainid, sonid uint32, data []byte) []byte {
	return append(CreateSignAndProtoHead(Sign_userid, userid, mainid, sonid, uint32(len(data))), data...)
}

//CreateToClientMsg 生成一条发送给客户端得消息
func CreateToSvrMsg(serverid uint64, mainid, sonid uint32, data []byte) []byte {
	return append(CreateSignAndProtoHead(Sign_serverid, serverid, mainid, sonid, uint32(len(data))), data...)
}
