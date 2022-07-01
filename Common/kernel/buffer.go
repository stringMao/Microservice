package kernel

import (
	"encoding/binary"
	"io"
)

//buffer.go
type ByteBuf struct {
	Data []byte
	readIndex int
	bytePool *BytePoolCap
}
func NewByteBuf(pool *BytePoolCap)*ByteBuf{
	return &ByteBuf{
		Data:pool.Get(),
		readIndex:0,
		bytePool: pool,
	}
}
func (bb *ByteBuf)InitBuf(size int){
	if cap(bb.Data) < size{
		bb.Free()
		bb.Data=make([]byte,0,size)
	}else{
		bb.Data = bb.Data[0:0]
	}
	bb.readIndex=0
}

func (bb *ByteBuf)Free(){
	bb.bytePool.Put(bb.Data)
	bb.Data=nil
	bb.bytePool=nil
	bb.readIndex=0
}
func (bb *ByteBuf)ResetReadIndex(n int){
	bb.readIndex=n
}
func (bb *ByteBuf) Slice(n int) []byte {
	r := bb.Data[bb.readIndex : bb.readIndex+n]
	bb.readIndex += n
	return r
}
func (bb *ByteBuf) Append(p ...byte) {
	bb.Data = append(bb.Data, p...)
}

func (bb *ByteBuf) WriteUint32BE(v uint32) {
	bb.Append(byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}
func (bb *ByteBuf) ReadUint32BE() uint32 {
	return binary.BigEndian.Uint32(bb.Slice(4))
}

func (bb *ByteBuf) Read(by []byte) (int, error) {
	if bb.readIndex == len(bb.Data) {
		return 0, io.EOF
	}
	n := len(by)
	if n+bb.readIndex > len(bb.Data) {
		n = len(bb.Data) - bb.readIndex
	}
	copy(by, bb.Data[bb.readIndex:])
	bb.readIndex += n
	return n, nil
}
func (bb *ByteBuf) Write(p []byte) (int, error) {
	bb.Data = append(bb.Data, p...)
	return len(p), nil
}

