package api

import (
	"steam-backend/routes/constants"
	"steam-backend/routes/controllers"

	"github.com/gin-gonic/gin"
)

// ForgesRoutes 锻刀所管理路由
type ForgesRoutes struct{}

// NewForgesRoutes 创建锻刀所管理路由实例
func NewForgesRoutes() *ForgesRoutes {
	return &ForgesRoutes{}
}

// RegisterRoutes 注册锻刀所管理路由（需要认证）
func (r *ForgesRoutes) RegisterRoutes(protected *gin.RouterGroup, ctrls *controllers.Controllers) {
	forges := protected.Group(constants.ForgesGroup)
	{
		forges.POST(constants.CreateForge, ctrls.Forge.CreateForge)
		forges.GET(constants.GetForges, ctrls.Forge.GetForges)
		forges.GET(constants.GetForgeByID, ctrls.Forge.GetForge)
		forges.POST(constants.UpdateForge, ctrls.Forge.UpdateForge)
		forges.GET(constants.DeleteForgeByID, ctrls.Forge.DeleteForge)
	}
}
