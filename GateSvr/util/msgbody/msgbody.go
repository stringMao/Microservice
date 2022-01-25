package msgbody

import (
	"Common/proto/base"

	"github.com/golang/protobuf/proto"
)

//协议体构建============================================================

func MakeToClientJionServerResult(codeid int, serverid uint64) []byte {
	pStruct := &base.ToClientJionServerResult{
		Codeid:   int32(codeid), //0成功 1服务器找不到 2重复加入
		Serverid: serverid,
	}
	dStruct, _ := proto.Marshal(pStruct)
	return dStruct
}
func MakeToClientLeaveServerResult(codeid int, serverid uint64) []byte {
	pStruct := &base.ToClientLeaveServerResult{
		Codeid:   int32(codeid), //0成功 1服务器找不到 2重复加入
		Serverid: serverid,
	}
	dStruct, _ := proto.Marshal(pStruct)
	return dStruct
}
