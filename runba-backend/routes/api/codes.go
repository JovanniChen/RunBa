package api

import (
	"steam-backend/routes/constants"
	"steam-backend/routes/controllers"

	"github.com/gin-gonic/gin"
)

// CodesRoutes 激活码管理路由
type CodesRoutes struct{}

// NewCodesRoutes 创建激活码管理路由实例
func NewCodesRoutes() *CodesRoutes {
	return &CodesRoutes{}
}

// RegisterRoutes 注册激活码管理路由（需要认证）
func (r *CodesRoutes) RegisterRoutes(protected *gin.RouterGroup, ctrls *controllers.Controllers) {
	// 激活码管理相关操作
	activationCodes := protected.Group(constants.ActivationCodesGroup)
	{
		// 基本CRUD操作
		activationCodes.POST(constants.CreateBatchActivationCodes, ctrls.ActivationCode.CreateBatchActivationCodes) // 批量创建激活码
		activationCodes.GET(constants.GetActivationCodes, ctrls.ActivationCode.GetActivationCodes)                  // 获取激活码列表
		activationCodes.GET(constants.GetActivationCode, ctrls.ActivationCode.GetActivationCode)                    // 获取单个激活码信息
		activationCodes.POST(constants.UpdateActivationCode, ctrls.ActivationCode.UpdateActivationCode)             // 更新激活码
		activationCodes.GET(constants.DeleteActivationCode, ctrls.ActivationCode.DeleteActivationCode)              // 删除激活码（软删除）
	}
}
