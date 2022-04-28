package watchdog

import (
	"Common/kernel/go-scoket/scokets"
	"Common/log"
	"Common/msg"
	"math/rand"
	"sync"
)

type  ClientManager struct {
    ServerMap 		*sync.Map
	serverList   	SvrListMap //tid对应的sid列表

	PlayerMap        *sync.Map
}

func NewClientManager()*ClientManager{
	return &ClientManager{
		ServerMap:  new(sync.Map),
		serverList: make(SvrListMap, 10),
		PlayerMap:  new(sync.Map),
	}
}

func (c *ClientManager)AddServerClient(p *scokets.Client){
	s:=p.Data.(*ServerData)
	c.ServerMap.Store(s.Serverid,p)
	//新注册的，需要在serverList记录
	c.serverList.Add(s.Tid, s.Sid)
	log.Debugf("新服务加入队列：TID[%d],SID[%d]",s.Tid, s.Sid)
}
func (c *ClientManager)RemoveServerClient(p *scokets.Client){
	s:=p.Data.(*ServerData)
	c.ServerMap.Delete(s.Serverid)
	c.serverList.DeleteValue(s.Tid, s.Sid)
	log.Debugf("服务删除队列：TID[%d],SID[%d]",s.Tid, s.Sid)
}
func (c *ClientManager)ServerIsExists(serverid uint64)bool{
	_,has:=c.ServerMap.Load(serverid)
	return has
}

func (c *ClientManager)AddPlayerClient(p *scokets.Client){
	s:=p.Data.(*PlayerData)
	c.PlayerMap.Store(s.UserId,p)
	log.Debugf("用户加入队列：uid[%d]",s.UserId)
}
func (c *ClientManager)RemovePlayerClient(p *scokets.Client){
	s:=p.Data.(*PlayerData)
	c.PlayerMap.Delete(s.UserId)
	log.Debugf("用户删除队列：uid[%d]",s.UserId)
}
func (c *ClientManager)PlayerIsExists(userid uint64)bool{
	_,has:=c.PlayerMap.Load(userid)
	return has
}
func (c *ClientManager)GetPlayer(userid uint64)*PlayerAgent{
	if p,has:=c.PlayerMap.Load(userid);has{
		if r,ok:=p.(*PlayerAgent);ok{
			return r
		}
	}
	return nil
}

func (c *ClientManager)SendToServer(serverid uint64,msg []byte)bool{
	if svr,has:=c.ServerMap.Load(serverid);has{
		p:=svr.(*scokets.Client)
		return p.SendData(msg)
	}
	return false
}
func (c *ClientManager)SendToPlayer(userid uint64,msg []byte)bool{
	if user,has:=c.PlayerMap.Load(userid);has{
		p:=user.(*PlayerAgent)
		return p.client.Session.Send(msg)
	}
	return false
}
func (c *ClientManager)AllocSvr(tid uint32)uint64{
	sid := c.serverList.RandOneValue(tid)
	if sid > 0 {
		return msg.EncodeServerID(tid, sid)
	}
	return 0
}
//==================================================
type SvrListMap map[uint32][]uint32

func (m SvrListMap) Add(key, value uint32) {
	if len(m[key]) == 0 {
		m[key] = []uint32{value}
	} else {
		m[key] = append(m[key], value)
	}
}

func (m SvrListMap) DeleteValue(key, value uint32) {
	if _, ok := m[key]; ok {
		for i := 0; i < len(m[key]); {
			if m[key][i] == value {
				m[key] = append(m[key][:i], m[key][i+1:]...)
			} else {
				i++
			}
		}
	}
}
func (m SvrListMap) RandOneValue(key uint32) uint32 {
	if length := len(m[key]); length > 0 {
		r := rand.Intn(length)
		return m[key][r]
	}
	return 0
}