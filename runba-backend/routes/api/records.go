package api

import (
	"steam-backend/routes/constants"
	"steam-backend/routes/controllers"

	"github.com/gin-gonic/gin"
)

// RecordsRoutes 兑换记录管理路由
type RecordsRoutes struct{}

// NewRecordsRoutes 创建兑换记录管理路由实例
func NewRecordsRoutes() *RecordsRoutes {
	return &RecordsRoutes{}
}

// RegisterRoutes 注册兑换记录管理路由（需要认证）
func (r *RecordsRoutes) RegisterRoutes(protected *gin.RouterGroup, ctrls *controllers.Controllers) {
	// 兑换记录管理相关操作
	records := protected.Group(constants.RecordsGroup)
	{
		records.POST(constants.CreateExchangeRecord, ctrls.Record.CreateExchangeRecord) // 创建兑换记录
		records.GET(constants.GetExchangeRecords, ctrls.Record.GetExchangeRecords)      // 获取兑换记录列表
		records.GET(constants.GetExchangeRecordByID, ctrls.Record.GetExchangeRecord)    // 获取单个兑换记录信息
	}
}
