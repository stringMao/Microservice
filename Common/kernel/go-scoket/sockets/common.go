package sockets

import "Common/kernel"

const messageMaxLen = 1024
const message_head_size = 4
//==字节池=============================
var bytePool  *kernel.BytePoolCap =nil

const (
	size_d = 1000
	len_d = messageMaxLen
	cap_d = messageMaxLen
)

func GetBytePool()*kernel.BytePoolCap{
	if bytePool==nil{
		InitBytePool(size_d,len_d,cap_d)
	}
	return bytePool
}
func InitBytePool(maxsize,len,cap int){
	bytePool= kernel.NewBytePoolCap(maxsize,len,cap)
}
func GetByteFormPool()[]byte{
	return GetBytePool().Get()
}
func PutByteToPool(buf []byte){
	GetBytePool().Put(buf)
}
