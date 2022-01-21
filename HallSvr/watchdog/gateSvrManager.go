package watchdog

//管家
import (
	"Common/constant"
	"Common/log"
	"Common/proto/base"
	"Common/svrfind"
	"HallSvr/config"
	"HallSvr/kernel"
	"fmt"
	"strconv"
	"time"

	"github.com/golang/protobuf/proto"
)

//连接的临时存储容器
var m_tempConnMap map[uint64]*kernel.ConnTcp = make(map[uint64]*kernel.ConnTcp, 10)

//连接网关服
func ConnectGateSvrs() {
	//登入验证
	logindata := &base.ServerLogin{
		Tid:      uint32(config.App.TID),
		Sid:      uint32(config.App.SID),
		Password: "test",
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

				if kernel.GetManagerSvrs().IsExist(serverid) {
					continue
				}

				pData, _ := proto.Marshal(logindata)
				conn := kernel.NewConnTcp(serverid, fmt.Sprintf("%s:%d", addr.Address, addr.Port), pData, DistributeMessage, HandleErr)
				if conn.Connect() {
					log.Infof("网关服连接[addr:%s]成功", conn.Addr)
					m_tempConnMap[uint64(serverid)] = conn
				}
			}
		}
		time.Sleep(10 * time.Second)
	}
}
