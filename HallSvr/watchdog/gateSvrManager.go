package watchdog

//管家
import (
	"Common/constant"
	"Common/log"
	"Common/msg"
	"Common/proto/base"
	"Common/setting"
	"Common/svrfind"
	"HallSvr/config"
	"HallSvr/kernel"
	"fmt"

	"github.com/golang/protobuf/proto"
)

//网关服连接列表
var GateSvrMap map[int]*kernel.ConnTcp = make(map[int]*kernel.ConnTcp)

func ConnectGateSvr(svritem *svrfind.ServerItem) {
	//连接GateSvr
	gatelist := svritem.GetSvr(setting.GetServerName(constant.TID_GateSvr), setting.GetServerTag(constant.TID_GateSvr))
	for k, v := range gatelist {
		fmt.Println(k, "  ", v.Service.Address)
		fmt.Println(k, "  ", v.Service.Port)
		fmt.Printf(" %d:%+v \n", k, v.Service)

		if addr, ok := v.Service.TaggedAddresses["server"]; ok {
			agentGateSvr := kernel.NewConnTcp(addr.Address, addr.Port, HandleMessage)
			if !agentGateSvr.Connect() {
				log.Errorf("网关服连接失败 addr:%s", agentGateSvr.SvrIP)
			} else {
				log.Infof("网关服连接成功 addr:%s", agentGateSvr.SvrIP)
				GateSvrMap[v.Service.Port] = agentGateSvr

				//登入验证
				logindata := &base.ServerLogin{
					Tid:      uint32(config.App.TID),
					Sid:      uint32(config.App.SID),
					Password: "test",
				}
				pData, err := proto.Marshal(logindata)
				if err != nil {
					log.Errorf("ConnectGateSvr proto.Marshal err:%s", err)
				}
				agentGateSvr.SendData(pData)
			}
		}
	}
}

func HandleMessage(data []byte, len int) {
	signhead := &msg.HeadSign{}
	signhead.Decode(data)

	switch signhead.SignType {
	case msg.Sign_serverid: //后8位是serverid

	case msg.Sign_userid: //后8位是userid
		head := &msg.HeadProto{}
		head.Decode(data[msg.GetSignHeadLength():])

		HandleClientMessage(head.MainID, head.SonID, head.Len, data[msg.GetHeadLength():])

	default:
	}
}

func HandleClientMessage(mainid, sonid uint32, len uint32, data []byte) {
	msgData := &base.TestMsg{}
	err := proto.Unmarshal(data, msgData)
	if err != nil {
		fmt.Println("协议解析失败2:", err)
		return //当远程客户端连接发生错误（断开）后，终止此协程。
	}
	fmt.Printf("mainid:%d,sonid:%d,str:%s\n", mainid, sonid, msgData.Txt)
}
