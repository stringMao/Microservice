package msg

import (
	"encoding/binary"
	"errors"
)

//消息协议头
type HeadProto struct {
	MainID uint32
	SonID  uint32
	Len    uint32
}

var size_head int = binary.Size(HeadProto{})

func (h *HeadProto) GetHeadLen() int {
	return size_head
}

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

//解码
func (h *HeadProto) Decode(b []byte) {
	h.MainID = binary.LittleEndian.Uint32(b[:4])
	h.SonID = binary.LittleEndian.Uint32(b[4:8])
	h.Len = binary.LittleEndian.Uint32(b[8:12])
}

//HeadBase 消息的基础“指向信息头”===================================================================
type HeadBase struct {
	SignType uint8  //sign_serverid sign_userid
	ID       uint64 //SignType==sign_serverid?serverid:userid
}

const Sign_serverid uint8 = 0
const Sign_userid uint8 = 1

var size_base int = binary.Size(HeadBase{})

//编码
func (h *HeadBase) Encode() []byte {
	b := make([]byte, size_base)
	b[0] = byte(h.SignType)
	b[1] = byte(h.ID)
	b[2] = byte(h.ID >> 8)
	b[3] = byte(h.ID >> 16)
	b[4] = byte(h.ID >> 24)
	b[5] = byte(h.ID >> 32)
	b[6] = byte(h.ID >> 40)
	b[7] = byte(h.ID >> 48)
	b[8] = byte(h.ID >> 56)
	return b
}

//解码
func (h *HeadBase) Decode(b []byte) {
	//_ = b[size_base-1] // bounds check hint to compiler;
	h.SignType = b[0]
	h.ID = binary.LittleEndian.Uint64(b[1:size_base])
}

//
func (h *HeadBase) ReplaceHeadBase(b []byte) {
	_ = b[8] // early bounds check to guarantee safety of writes below
	b[0] = byte(h.SignType)
	b[1] = byte(h.ID)
	b[2] = byte(h.ID >> 8)
	b[3] = byte(h.ID >> 16)
	b[4] = byte(h.ID >> 24)
	b[5] = byte(h.ID >> 32)
	b[6] = byte(h.ID >> 40)
	b[7] = byte(h.ID >> 48)
	b[8] = byte(h.ID >> 56)
}

//================================================================================
func CreateMsgHead(base *HeadBase, head *HeadProto) []byte {
	b := make([]byte, size_base+size_head)
	copy(b, base.Encode())
	copy(b[size_base:], head.Encode())
	return b
}

func DecodeMsgHead(b []byte, n int) (*HeadBase, *HeadProto, error) {
	if n < (size_base + size_head) {
		return nil, nil, errors.New("长度小于head")
	}
	base := &HeadBase{}
	base.Decode(b[:size_base])
	//base.SignType = b[0]
	//base.ID = binary.LittleEndian.Uint64(b[1:9])

	head := &HeadProto{}
	head.Decode(b[size_base:(size_base + size_head)])
	//head.MainID = binary.LittleEndian.Uint32(b[9:13])
	//head.SonID = binary.LittleEndian.Uint32(b[13:17])
	//head.Len = binary.LittleEndian.Uint32(b[17:21])

	return base, head, nil
}

func CreateMsgData(base *HeadBase, head *HeadProto, data []byte) []byte {
	b := make([]byte, size_base+size_head+int(head.Len))
	copy(b, base.Encode())
	copy(b[size_base:], head.Encode())
	if head.Len > 0 {
		copy(b[size_base+size_head:], data)
	}
	return b
}

func CreateMsg2(base *HeadBase, data []byte) []byte {
	b := make([]byte, size_base+len(data))
	copy(b, base.Encode())
	copy(b[size_base:], data)
	return b
}

func GetHeadLength() int {
	return size_base + size_head
}

func GetHeadbaseLength() int {
	return size_base
}
