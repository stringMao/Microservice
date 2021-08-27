package player

import (
	"fmt"
	"time"
)

//道具结构
type ToolData struct {
	ID      int       //道具ID
	Typ     int       //道具类型  0一般道具  1计时道具
	Count   int       //道具数量
	EndTime time.Time //到期时间
}

type PlayerTool struct {
	ToolList map[int]ToolData //道具列表
	redisKey string           //redis里的缓存key
}

func newPlayerTool(userid uint64) *PlayerTool {
	return &PlayerTool{
		ToolList: make(map[int]ToolData),
		redisKey: fmt.Sprintf("PlayerTool:%d", userid),
	}
}

//一般道具设置数量
func (p *PlayerTool) Set(id, count int) {
	if tooldata, ok := p.ToolList[id]; ok {
		tooldata.Count = count
	} else {
		t := ToolData{
			ID:    id,
			Typ:   0,
			Count: count,
		}
		p.ToolList[id] = t
	}
}

//更新道具数量(count 增量)
func (p *PlayerTool) UpdateCount(id, count int) {
	if tooldata, ok := p.ToolList[id]; ok {
		tooldata.Count += count
	} else {
		p.Set(id, count)
	}
}
