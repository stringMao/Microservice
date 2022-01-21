package msg

import (
	"encoding/binary"
	"errors"
	"math"
)

//协议头长度
var size_head int = 0

//标记头长度
var size_sign int = 0

func init() {
	size_head = binary.Size(HeadProto{})
	size_sign = binary.Size(HeadSign{})
}

//协议头
type HeadProto struct {
	MainID uint32
	SonID  uint32
	Len    uint32
}

func (h *HeadProto) GetHeadLen() int {
	return size_head
}

//协议头（小端） 编码成二进制切片
func (h *HeadProto) Encode() []byte {
	b := make([]byte, size_head)
	b[0] = byte(h.MainID)
	b[1] = byte(h.MainID >> 8)
	b[2] = byte(h.MainID >> 16)
	b[3] = byte(h.MainID >> 24)

	b[4] = byte(h.SonID)
	b[5] = byte(h.SonID >> 8)
	b[6] = byte(h.SonID >> 16)
	b[7] = byte(h.SonID >> 24)

	b[8] = byte(h.Len)
	b[9] = byte(h.Len >> 8)
	b[10] = byte(h.Len >> 16)
	b[11] = byte(h.Len >> 24)

	return b
}

//二进制切片解码 成协议头
func (h *HeadProto) Decode(b []byte) {
	h.MainID = binary.LittleEndian.Uint32(b[:4])
	h.SonID = binary.LittleEndian.Uint32(b[4:8])
	h.Len = binary.LittleEndian.Uint32(b[8:12])
}

//HeadSign 消息的“标记头” 消息来自哪里或者消息去向哪里===================================================================

const Sign_serverid uint8 = 1 //表示 serverid
const Sign_userid uint8 = 2   //表示 userid
const Sign_Self uint8 = 3     //自己内部消息

type HeadSign struct {
	SignType uint8  //sign_serverid sign_userid
	SignId   uint64 //SignType==sign_serverid?serverid:userid
	Tid      uint32
	Sid      uint32
}

//“标记头” 编码
func (h *HeadSign) Encode() []byte {
	b := make([]byte, size_sign)
	b[0] = byte(h.SignType)
	b[1] = byte(h.SignId)
	b[2] = byte(h.SignId >> 8)
	b[3] = byte(h.SignId >> 16)
	b[4] = byte(h.SignId >> 24)
	b[5] = byte(h.SignId >> 32)
	b[6] = byte(h.SignId >> 40)
	b[7] = byte(h.SignId >> 48)
	b[8] = byte(h.SignId >> 56)
	return b
}

//解码成“标记头”
func (h *HeadSign) Decode(b []byte) error {
	//_ = b[size_sign-1] // bounds check hint to compiler;
	h.SignType = b[0]
	if h.SignType != Sign_serverid && h.SignType != Sign_userid {
		return errors.New("消息头的SignType不存在")
	}
	h.SignId = binary.LittleEndian.Uint64(b[1:size_sign])
	if h.SignType == Sign_serverid {
		h.Tid, h.Sid = DecodeServerID(h.SignId)
	}
	return nil
}

//将二进制流头部的“标记头”重置
// func (h *HeadSign) ReplaceHeadBase(b []byte) {
// 	_ = b[8] // early bounds check to guarantee safety of writes below
// 	b[0] = byte(h.SignType)
// 	b[1] = byte(h.ID)
// 	b[2] = byte(h.ID >> 8)
// 	b[3] = byte(h.ID >> 16)
// 	b[4] = byte(h.ID >> 24)
// 	b[5] = byte(h.ID >> 32)
// 	b[6] = byte(h.ID >> 40)
// 	b[7] = byte(h.ID >> 48)
// 	b[8] = byte(h.ID >> 56)
// }

//================================================================================

//获得“标记头”+“协议头”的总长度
func GetHeadLength() int {
	return size_sign + size_head
}

//获得“标记头”的长度
func GetSignHeadLength() int {
	return size_sign
}

//获得“协议头”的总长度
func GetProtoHeadLength() int {
	return size_head
}

//解码serverid to tid sid
func DecodeServerID(serverid uint64) (tid, sid uint32) {
	sid = uint32(serverid & (math.MaxUint32 << 32))
	tid = uint32(serverid & math.MaxUint32)
	return tid, sid
}

//将tid 和sid组合成serverid
func EncodeServerID(tid, sid uint32) uint64 {
	return uint64(sid)<<32 + uint64(tid)
}

//创建“标记头”的二进制流
func CreateSignHead(signtype uint8, id uint64) []byte {
	return (&HeadSign{
		SignType: signtype,
		SignId:   id,
	}).Encode()
}

//创建“协议头”的二进制流
func CreateProtoHead(mainid, sonid, lenght uint32) []byte {
	return (&HeadProto{
		MainID: mainid,
		SonID:  sonid,
		Len:    lenght,
	}).Encode()
}

//创建“标记头”+“协议头”的二进制流
func CreateSignAndProtoHead(signtype uint8, id uint64, mainid, sonid, lenght uint32) []byte {
	return append(CreateSignHead(signtype, id), CreateProtoHead(mainid, sonid, lenght)...)
}

//创建完整的消息流 “标记头”+“协议头”+“协议体”
func CreateWholeMsgData(signtype uint8, id uint64, mainid, sonid uint32, data []byte) []byte {
	return append(CreateSignAndProtoHead(signtype, id, mainid, sonid, uint32(len(data))), data...)
}

//创建完整的协议流 “协议头”+“协议体”
func CreateWholeProtoData(mainid, sonid uint32, data []byte) []byte {
	return append(CreateProtoHead(mainid, sonid, uint32(len(data))), data...)
}

//为协议流添加“标记头”
func AddSignHead(signtype uint8, id uint64, data []byte) []byte {
	return append(CreateSignHead(signtype, id), data...)
}

//修改“标记头”字段
func ChangeSignHead(signtype uint8, id uint64, data []byte) {
	data[0] = byte(signtype)
	data[1] = byte(id)
	data[2] = byte(id >> 8)
	data[3] = byte(id >> 16)
	data[4] = byte(id >> 24)
	data[5] = byte(id >> 32)
	data[6] = byte(id >> 40)
	data[7] = byte(id >> 48)
	data[8] = byte(id >> 56)
}

//完整的消息结构
// type WholeProtoStruct struct {
// 	Sign HeadSign
// 	Head HeadProto
// 	Data []byte
// }

//标准的协议结构
// type ProtoStruct struct {
// 	Head HeadProto
// 	Data []byte
// }

//获得”标记头“
func GetSign(b []byte) *HeadSign {
	h := &HeadSign{}
	h.SignType = b[0]
	// if h.SignType != Sign_serverid && h.SignType != Sign_userid {
	// 	return nil
	// }
	h.SignId = binary.LittleEndian.Uint64(b[1:size_sign])
	if h.SignType == Sign_serverid {
		h.Tid, h.Sid = DecodeServerID(h.SignId)
	}
	return h
}

//从完整包里获得”消息头“
func GetHead(b []byte) *HeadProto {
	idex := size_sign
	h := &HeadProto{
		MainID: binary.LittleEndian.Uint32(b[idex : idex+4]),
		SonID:  binary.LittleEndian.Uint32(b[idex+4 : idex+8]),
		Len:    binary.LittleEndian.Uint32(b[idex+8 : idex+12]),
	}
	return h
}

//从只有“消息头和消息体”的包中获得”消息头“
func GetHeadOfProto(b []byte) *HeadProto {
	h := &HeadProto{
		MainID: binary.LittleEndian.Uint32(b[:+4]),
		SonID:  binary.LittleEndian.Uint32(b[+4:+8]),
		Len:    binary.LittleEndian.Uint32(b[+8:+12]),
	}
	return h
}
