package watchdog

import (
	"Common/kernel/go-scoket/scokets"
	"sync"
)

type Servers struct {
	serverMap sync.Map
}
var ServerList *Servers=nil
func init(){
	ServerList=new(Servers)
}

func (s *Servers) Add(serverId uint64, c *LogicHandler)  {
	s.serverMap.Store(serverId,c)
}
func (s *Servers)Remove(serverId uint64){
	s.serverMap.Delete(serverId)
}
func (s *Servers)Get(serverId uint64)*LogicHandler{
	p,has:=s.serverMap.Load(serverId)
	if has!=true || p==nil{
		return nil
	}
	return p.(*LogicHandler)
}
func (s *Servers)SendData(serverId uint64,buf []byte){
	if v,ok:=s.serverMap.Load(serverId);ok{
		c:=v.(*scokets.Connector)
		c.SendData(buf)
	}
}
func (s *Servers)IsExists(serverId uint64)bool{
	_,has:=s.serverMap.Load(serverId)
	return has
}