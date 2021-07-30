package balanceManager

//服务分配时的负载均衡策略

//分配一个tid类型的服务serverid
func AllocServerid(tid uint16) uint64 {
	return uint64(1)<<16 + uint64(tid)
}
