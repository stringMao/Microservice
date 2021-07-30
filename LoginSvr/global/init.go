package global

import (
	"Common/log"
	"LoginSvr/models"
)

//HallList 大厅服务器列表信息
var HallList []models.Hall

//Init global初始化
func Init() {
	HallList = models.LoadServerInfo()
	log.Logger.Info("global init success")
}
