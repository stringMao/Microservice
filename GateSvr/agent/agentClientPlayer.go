package agent

import (
	"GateSvr/core/player"
)

//加载玩家数据
func (a *agentClient) Init() bool {
	if a.Player == nil {
		a.Player = player.NewPlayer(a.Userid)
	}

	if !a.Player.Init(a.Userid) {
		return false
	}

	return true
}

//
func (a *agentClient) Save() {
	if a.Player == nil {
		return
	}
	a.Player.Save(a.Userid)
}

//

//积分变更
func (a *agentClient) UpdateScore(score int) {
	a.Player.CashData.UpdateScore(score)
}
