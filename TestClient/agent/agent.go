package agent

import (
	"Common/kernel/go-scoket/sockets"
	"Common/log"
	"Common/proto/base"
	"github.com/golang/protobuf/proto"
)

type LogicAgent struct {
	Connector *sockets.Connector
}

func NewLogicAgent(connector *sockets.Connector) *LogicAgent {
	return &LogicAgent{Connector: connector}
}


func (c *LogicAgent)BeforeHandle(engine *sockets.Engine,len int,buffer []byte){
	pMessage := &base.Message{}
	if proto.Unmarshal(buffer, pMessage) != nil {
		log.Warnln("ServerAgent BeforeHandle is err")
		return
	}
	fn:=c.Connector.GetHandleFunc(pMessage.MessageId)
	if fn==nil{
		log.Errorf("收到未知消息 msgid[%d]",pMessage.MessageId)
		return
	}
	fn(engine,pMessage.Body)
}

func (c *LogicAgent)CloseHandle(engine *sockets.Engine){
	log.Debugln("socket is close")
}