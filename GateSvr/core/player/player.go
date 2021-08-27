package player

type Player struct {
	BaseData *PlayerBase //基础信息
	CashData *PlayerCash //货币数据
	ToolData *PlayerTool //道具数据
}

func NewPlayer(userid uint64) *Player {
	p := &Player{
		BaseData: newPlyerBase(userid),
		CashData: newPlayerCash(userid),
		ToolData: newPlayerTool(userid),
	}
	return p
}

func (p *Player) Init(userid uint64) bool {
	//初始基础信息
	p.BaseData.LoadRedisData()

	//初始货币数据
	if !p.CashData.LoadRedisData(true) {
		if !p.CashData.LoadDBData(userid) {
			return false
		}
		p.CashData.SetRedisData()
	}
	//初始道具信息

	return true
}

//保存玩家的内存数据
func (p *Player) Save(userid uint64) {
	p.CashData.Save(userid)
	//p.ToolData.Save
}
