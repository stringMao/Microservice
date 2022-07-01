package watchdog

import (
	"Common/kernel/go-scoket/sockets"
	"Common/log"
	"Common/msg"
	"github.com/golang/protobuf/proto"
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

func (c *ClientManager) AddServer(p *sockets.Engine){
	s:=p.Data.(*ServerData)
	c.ServerMap.Store(s.Serverid,p)
	//新注册的，需要在serverList记录
	c.serverList.Add(s.Tid, s.Serverid)
	log.Debugf("新服务加入队列：TID[%d],SID[%d]",s.Tid, s.Sid)
}
func (c *ClientManager) RemoveServer(p *sockets.Engine){
	s:=p.Data.(*ServerData)
	c.serverList.DeleteValue(s.Tid, s.Serverid)
	c.ServerMap.Delete(s.Serverid)

	log.Debugf("服务删除队列：TID[%d],SID[%d]",s.Tid, s.Sid)
}
func (c *ClientManager)ServerIsExists(serverid uint64)bool{
	_,has:=c.ServerMap.Load(serverid)
	return has
}

func (c *ClientManager) AddPlayer(p *sockets.Engine){
	s:=p.Data.(*PlayerData)
	c.PlayerMap.Store(s.UserId,p)
	log.Debugf("uid[%d]加入队列",s.UserId)
}
func (c *ClientManager) RemovePlayer(p *sockets.Engine){
	s:=p.Data.(*PlayerData)
	c.PlayerMap.Delete(s.UserId)
	log.Debugf("uid[%d]删除队列",s.UserId)
}
func (c *ClientManager)PlayerIsExists(userid uint64)bool{
	_,has:=c.PlayerMap.Load(userid)
	return has
}
func (c *ClientManager)GetPlayer(userid uint64)*sockets.Engine{
	if p,has:=c.PlayerMap.Load(userid);has{
		if r,ok:=p.(*sockets.Engine);ok{
			return r
		}
	}
	return nil
}

// AllocSvr 负载均衡的分配一个服务器
func (c *ClientManager)AllocSvr(tid uint32)uint64{
	return c.serverList.RandOneValue(tid)
}


//消息发送====================================================================

func (c *ClientManager)SendToServer(serverId uint64,pb proto.Message)bool {
	buf:=msg.ProtoMarshal(pb)
	return c.SendToServer2(serverId,buf)
}
func (c *ClientManager)SendToServer2(serverid uint64,msg []byte)bool{
	if svr,has:=c.ServerMap.Load(serverid);has{
		p:=svr.(*sockets.Engine)
		return p.SendData(msg)
	}
	return false
}

// BroadcastToTid 广播消息给指定类型TID的服务
func (c *ClientManager) BroadcastToTid(tid uint32,pb proto.Message){
	//TODO 虽然buf指针被多次复制，但是没关系，他指向的内存不会被修改
	buf:=msg.ProtoMarshal(pb)
	if list,has:=c.serverList[tid];has{
		for _, serverId := range list {
			c.SendToServer2(serverId,buf)
		}
	}
}

// BroadcastToServers 广播给所有服务
func  (c *ClientManager)BroadcastToServers(pList []uint64,pb proto.Message){
	buf:=msg.ProtoMarshal(pb)
	if pList==nil{
		c.ServerMap.Range(func(key, value interface{}) bool {
			p:=value.(*ServerAgent)
			p.engine.Session.Send(buf)
			return  true
		})
	}else{
		for _,v:=range pList {
			c.SendToServer2(v,buf)
		}
	}
}

func (c *ClientManager)SendToPlayer(userid uint64,pb proto.Message)bool{
	buf:=msg.ProtoMarshal(pb)
	return c.SendToPlayer2(userid,buf)
}
func (c *ClientManager)SendToPlayer2(userid uint64,msg []byte)bool{
	if user,has:=c.PlayerMap.Load(userid);has{
		p:=user.(*sockets.Engine)
		return p.Session.Send(msg)
	}
	return false
}

// BroadcastToPlayers 广播给玩家
func  (c *ClientManager)BroadcastToPlayers(pList []uint64,pb proto.Message){
	buf:=msg.ProtoMarshal(pb)
	if pList==nil{
		c.PlayerMap.Range(func(key, value interface{}) bool {
			p:=value.(*sockets.Engine)
			p.Session.Send(buf)
			return  true
		})
	}else{
		for _,v:=range pList {
			c.SendToPlayer2(v,buf)
		}
	}
}

func (c *ClientManager)SendToPlayerOfTid(userid uint64,tid uint32,pb proto.Message)bool{
	buf:=msg.ProtoMarshal(pb)
	return c.SendToPlayerOfTid2(userid,tid,buf)
}
func (c *ClientManager)SendToPlayerOfTid2(userid uint64,tid uint32 ,buf []byte)bool{
	if user,has:=c.PlayerMap.Load(userid);has{
		p:=user.(*sockets.Engine)
		pData :=p.Data.(*PlayerData)
		if sid,ok:= pData.SvrList[tid];ok{
			return c.SendToServer2(msg.EncodeServerID(tid,sid),buf)
		}
	}
	return false
}



//==================================================

type SvrListMap map[uint32][]uint64

func (m SvrListMap) Add(key uint32, value uint64) {
	if m[key] == nil {
		m[key] = []uint64{value}
	} else {
		m[key] = append(m[key], value)
	}
}

func (m SvrListMap) DeleteValue(key uint32, value uint64) {
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
func (m SvrListMap) RandOneValue(key uint32) uint64 {
	if length := len(m[key]); length > 0 {
		r := rand.Intn(length)
		return m[key][r]
	}
	return 0
}


