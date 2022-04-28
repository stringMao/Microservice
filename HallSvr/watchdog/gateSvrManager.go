package watchdog

//管家
import (
	"Common/constant"
	"Common/kernel/go-scoket/scokets"
	"Common/log"
	"Common/msg"
	"Common/proto/base"
	"Common/svrfind"
	"HallSvr/config"
	"HallSvr/core/send"
	"fmt"
	"strconv"
	"time"

	"github.com/golang/protobuf/proto"
)



type HandleFunction struct {
	HandleFuncMap  map[uint64]func(srcType uint8,srcId uint64,s *scokets.Connector,buf []byte)
}
var m_handleFunction *HandleFunction=new(HandleFunction)
func (c *HandleFunction)GetHandleFunc(id uint64)func(srcType uint8,srcId uint64,s *scokets.Connector,buf []byte){
	if c.HandleFuncMap==nil{
		return nil
	}
	fn, has := c.HandleFuncMap[id]
	if has {
		return fn
	}
	return nil
}
func (c *HandleFunction)AddHandleFuc(id uint64,fn func(srcType uint8,srcId uint64,s *scokets.Connector,buf []byte)){
	if c.HandleFuncMap==nil{
		c.HandleFuncMap= make(map[uint64]func(srcType uint8,srcId uint64,s *scokets.Connector,buf []byte))
	}
	c.HandleFuncMap[id]=fn
}

func InitHandleFunc(){
	m_handleFunction.AddHandleFuc(msg.MergeMsgID(msg.MID_Gate,msg.Gate_SS_SvrLoginResult),DoLoginGateSvr)
    m_handleFunction.AddHandleFuc(msg.MergeMsgID(msg.MID_Gate,msg.Gate_SS_ClientJionReq),DoPlayerJionReq)
	m_handleFunction.AddHandleFuc(msg.MergeMsgID(msg.MID_Gate,msg.Gate_SS_ClientLeaveReq),DoPlayerLeaveReq)

	m_handleFunction.AddHandleFuc(msg.MergeMsgID(msg.MID_Test,msg.Test_1),DoPlayerTestMsg)

}

type LogicHandler struct {
	HandleFunc  *HandleFunction
	connector  *scokets.Connector
	SignHead *msg.HeadSign
	MsgHead  *msg.HeadProto
	bLogin   bool
	loginTick *time.Timer
}

func NewLogicHandler()*LogicHandler{
	return &LogicHandler{
		HandleFunc:m_handleFunction,
		connector:nil,
		SignHead: new(msg.HeadSign),
		MsgHead: new(msg.HeadProto),
		bLogin:false,
		loginTick:nil,
	}
}

func (c *LogicHandler)BeforeHandle(client *scokets.Client,len int,buffer []byte){
	if len < msg.GetHeadLength() { //消息大小安全检测
		log.Error("msg len too samll")
		return
	}
	if msg.ParseSign(c.SignHead,buffer) && msg.ParseHead(c.MsgHead,buffer){
		if  fn:=c.HandleFunc.GetHandleFunc(msg.MergeMsgID(c.MsgHead.MainID,c.MsgHead.SonID));fn!=nil{
			buf:=scokets.GetByteFormPool()
			copy(buf, buffer[msg.GetHeadLength():len])
			fn(c.SignHead.SignType,c.SignHead.SignId,c.connector,buf)
		}
	}
}

func (c *LogicHandler)CloseHandle(client *scokets.Client){

}

//连接网关服
func ConnectGateSvrs() {
	//登入验证
	logindata := &base.ServerLogin{
		Tid:      uint32(config.App.TID),
		Sid:      uint32(config.App.SID),
		Password: "123",
	}

	
	for {
		gatelist := svrfind.G_ServerRegister.GetSvr(constant.GetServerName(constant.TID_GateSvr), constant.GetServerTag(constant.TID_GateSvr))
		for _, v := range gatelist {
			//fmt.Println(k, "  ", v.Service.Address)
			//fmt.Println(k, "  ", v.Service.Port)
			//fmt.Printf(" %d:%+v \n", k, v.Service)

			if addr, ok := v.Service.TaggedAddresses["server"]; ok {
				strserverid, ok := v.Service.Meta["ServerID"]
				if !ok {
					log.Errorf("网关服未正确注册 serverid")
					continue
				}

				serverid, err := strconv.ParseUint(strserverid, 10, 64) //strconv.Atoi(strserverid)
				if err != nil {
					log.Errorf("网关服未正确注册2 serverid ")
					continue
				}

				if ServerList.IsExists(serverid) {
					//如果登入失败的，则删除

					continue
				}

				pData, _ := proto.Marshal(logindata)

				logicHandler:=NewLogicHandler()
				connector:=scokets.NewScoketConnector(fmt.Sprintf("%s:%d", addr.Address, addr.Port),
					logicHandler,time.Second)
				if connector.StartConnect(){
					log.Infof("网关服连接[addr:%s]成功", connector.GetAddrInfo())
					logicHandler.connector=connector

					//agenter保存到队列
					ServerList.Add(serverid,logicHandler)
					connector.SendData(send.CreateMsgToServerID(serverid, msg.MID_Gate, msg.SS_SvrRegisterGateSvr, pData))
					//connector.SendData(pData)
				}
			}
		}
		time.Sleep(30 * time.Second)
	}
}
