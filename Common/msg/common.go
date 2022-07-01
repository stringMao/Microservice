package msg

import (
	"Common/log"
	"Common/proto/base"
	"github.com/golang/protobuf/proto"
)


func ProtoMarshal(pb proto.Message) []byte {
	bytes, err := proto.Marshal(pb)
	if err != nil {
		log.Error("marshalProto Marshal error:", err)
		return nil
	}
	return bytes
}
func ProtoDecode(buf []byte,pb proto.Message)(pForward *base.Forward){
	pForward=&base.Forward{}
	err:=proto.Unmarshal(buf, pForward)
	if err!= nil{
		log.Error("ProtoDecode error:", err)
		pForward=nil
		return
	}
	if pb!=nil{
		err=proto.Unmarshal(pForward.Body, pb)
		if err!=nil{
			log.Error("ProtoDecode error2:", err)
			pForward=nil
			return
		}
	}
	return
}

//func CreateBaseMessage(msgId uint32,pb proto.Message)[]byte{
//	return ProtoMarshal(&base.Message{
//		MessageId: msgId,
//		Body: ProtoMarshal(pb),
//	})
//}


func CreateForward(userid,serverid uint64,data []byte)[]byte{
	return ProtoMarshal(&base.Forward{
		UserId: userid,
		ServerId: serverid,
		Body: data,
	})
}
func ForwardFromClient(userid uint64,data []byte)[]byte{
	return CreateForward(userid,0,data)
}

//ForwardFromServer 。。
func ForwardFromServer(serverid uint64,data []byte)[]byte{
	return CreateForward(0,serverid,data)
}

func NewMessage(msgId uint32,userid,serverid uint64,pb proto.Message)[]byte{
	return ProtoMarshal(&base.Message{
		MessageId: msgId,
		Body: CreateForward(userid,serverid,ProtoMarshal(pb)),
	})
}
func NewClientMessage(msgId uint32,pb proto.Message)[]byte{
	return ProtoMarshal(&base.Message{
		MessageId: msgId,
		Body: ProtoMarshal(pb),
	})
}

