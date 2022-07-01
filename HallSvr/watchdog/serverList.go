package watchdog

import (
	"Common/kernel/go-scoket/sockets"
	"sync"
)

type Servers struct {
	serverMap sync.Map
}
var ServerList *Servers=nil
func init(){
	ServerList=new(Servers)
}

func (s *Servers) Add(serverId uint64, c *sockets.Engine)  {
	s.serverMap.Store(serverId,c)
}
func (s *Servers)Remove(serverId uint64){
	s.serverMap.Delete(serverId)
}
func (s *Servers)Get(serverId uint64)*sockets.Engine {
	p,has:=s.serverMap.Load(serverId)
	if has!=true || p==nil{
		return nil
	}
	return p.(*sockets.Engine)
}
func (s *Servers)SendData(serverId uint64,buf []byte){
	if v,ok:=s.serverMap.Load(serverId);ok{
		c:=v.(*sockets.Engine)
		c.SendData(buf)
	}
}
func (s *Servers)IsExists(serverId uint64)bool{
	_,has:=s.serverMap.Load(serverId)
	return has
}

func (s *Servers) ClearRegisterTimeOut(){
	s.serverMap.Range(func(key, value interface{}) bool {
		p:=value.(*sockets.Engine)
		//p.Data==nil 说明注册消息没有返回成功
		if p.Data==nil{
			s.serverMap.Delete(key)
			p.Close()
		}
		return true
	})
}