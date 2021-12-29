package constant

import "fmt"

//TID 分配
const (
	TID_LoginSvr = 1 //登入服tid
	TID_GateSvr  = 2 //网关服tid
	TID_HallSvr  = 3 //大厅服tid
)

//GetServerName 获得服务name
func GetServerName(tid int) string {
	switch tid {
	case TID_LoginSvr:
		return "登入服"
	case TID_GateSvr:
		return "网关服"
	case TID_HallSvr:
		return "大厅服"
	default:
		return "未命名"
	}
}

//GetServerTag 获得服务标签名
func GetServerTag(tid int) string {
	switch tid {
	case TID_LoginSvr:
		return "登入服"
	case TID_GateSvr:
		return "网关服"
	case TID_HallSvr:
		return "大厅服"
	default:
		return "未命名"
	}
}

func GetServerID(tid, sid int) uint64 {
	return uint64(sid)<<32 + uint64(tid)
}

func GetServerIDName(tid, sid int) string {
	return fmt.Sprintf("TID:%d_SID:%d_ServerId:%d", tid, sid, GetServerID(tid, sid))
}
