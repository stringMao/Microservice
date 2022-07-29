package watchdog

//管家
import (
	"Common/constant"
	"Common/kernel/go-scoket/sockets"
	"Common/log"
	"Common/msg"
	"Common/proto/base"
	"Common/proto/gateProto"
	"Common/svrfind"
	"Common/try"
	"HallSvr/config"
	"fmt"
	"strconv"
	"time"

	"github.com/golang/protobuf/proto"
)

func InitHandleFunc(p *sockets.Connector){
	p.AddHandleFuc(msg.CommonSvrMsg_SvrRegisterResult, RegisterGateResult)
	p.AddHandleFuc(msg.CommonSvrMsg_UserJoin, UserJoin)
	p.AddHandleFuc(msg.CommonSvrMsg_UserQuit, UserQuit)
	p.AddHandleFuc(msg.CommonSvrMsg_UserOffline, UserOffline)
	p.AddHandleFuc(msg.ToHallSvr_Test, TestMessage)


}

type ServerData struct {
	Tid      uint32 //服务类型id
	Sid      uint32 //同服务类型下的唯一标识id
	Serverid uint64 //0+0+sid+tid 位运算获得
	SvrType  uint32
}
type LogicAgent struct {
	connector  *sockets.Connector
}

func NewLogicHandler()*LogicAgent {
	return &LogicAgent{
		connector:nil,
	}
}

func (c *LogicAgent)BeforeHandle(engine *sockets.Engine,len int,buffer []byte){
	defer try.Catch()
	pMessage:=&base.Message{}
	err := proto.Unmarshal(buffer, pMessage)
	if err != nil {
		log.Warnln(" BeforeHandle is err:",err)
		return
	}
    fn:=c.connector.Handler.GetHandleFunc(pMessage.MessageId)
    if fn==nil{
		log.Errorln(" BeforeHandle messsgeId  no found:",pMessage.MessageId)
		return
	}
	fn(engine,pMessage.Body)
}

func (c *LogicAgent)CloseHandle(engine *sockets.Engine){
	if engine==nil || engine.Data==nil{
		return
	}
	pServerData:=engine.Data.(*ServerData)
	ServerList.Remove(pServerData.Serverid)
	//遍历在线用户列表，来自该gate服的都标记成短线或者直接清理
}

// ConnectGateSvrs 连接网关服
func ConnectGateSvrs() {
	for {
		//先清理没完成注册的连接
		ServerList.ClearRegisterTimeOut()


		//检查没有连接的gate，并且连接注册
		gateList := svrfind.G_ServerRegister.GetSvr(constant.GetServerName(constant.TID_GateSvr), constant.GetServerTag(constant.TID_GateSvr))
		for _, v := range gateList {
			//fmt.Println(k, "  ", v.Service.Address)
			//fmt.Println(k, "  ", v.Service.Port)
			//fmt.Printf(" %d:%+v \n", k, v.Service)

			if addr, ok := v.Service.TaggedAddresses["server"]; ok {
				str, ok := v.Service.Meta["ServerID"]
				if !ok {
					log.Errorf("网关服未正确注册 serverid")
					continue
				}
				serverid, err := strconv.ParseUint(str, 10, 64) //strconv.Atoi(strserverid)
				if err != nil {
					log.Errorf("网关服未正确注册2 serverid ")
					continue
				}

				if ServerList.IsExists(serverid) {
					//如果登入失败的，则删除
					continue
				}

				pLogic:=NewLogicHandler()
				connector:=sockets.NewConnector(fmt.Sprintf("%s:%d", addr.Address, addr.Port),
					pLogic,time.Second)
				if connector.StartConnect(){
					log.Infof("网关服连接[addr:%s]成功\n", connector.GetAddrInfo())
					pLogic.connector=connector
					ServerList.Add(serverid,connector.Engine)
					InitHandleFunc(connector)

					pObj:=&gateProto.SvrRegisterReq{
						Tid:      uint32(config.App.TID),
						Sid:      uint32(config.App.SID),
						Password: "123",
						SvrType:  11,
					}
					connector.SendData(NewMsgToSvr(serverid,msg.ToGateSvr_SvrRegister,pObj))
				}
			}
		}
		time.Sleep(10 * time.Second)
	}
}
