package svrbalanced

//负载均衡相关
import (
	"Common/constant"
	"Common/svrfind"
	"LoginSvr/global"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var rwlock *sync.RWMutex = new(sync.RWMutex)
var gatesvrlist []global.ServerIPInfo

func RefreshGateSvrList() {
	go func() {
		for {
			secSvrEntry := svrfind.G_ServerRegister.GetSvr(constant.GetServerName(constant.TID_GateSvr), constant.GetServerTag(constant.TID_GateSvr))
			rwlock.Lock()
			gatesvrlist = gatesvrlist[0:0] //清空
			for _, v := range secSvrEntry {
				//fmt.Println(k, "  ", v.Service.Address)
				//fmt.Println(k, "  ", v.Service.Port)
				//fmt.Printf(" %d:%+v \n", k, v.Service)
				gatesvrlist = append(gatesvrlist, global.ServerIPInfo{Address: fmt.Sprintf("%s:%d", v.Service.Address, v.Service.Port)})
			}
			rwlock.Unlock()
			time.Sleep(time.Second * 10)
		}
	}()

}

//随机负载均衡
func GetSvrAddr() []global.ServerIPInfo {
	rwlock.RLock()
	defer rwlock.RUnlock()
	n := len(gatesvrlist)
	if n == 0 {
		return nil
	} else if n == 1 {
		return []global.ServerIPInfo{gatesvrlist[0]}
	} else {
		return []global.ServerIPInfo{gatesvrlist[rand.Intn(n)]}
	}
}
