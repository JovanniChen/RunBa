package api

import (
	"steam-backend/routes/constants"
	"steam-backend/routes/controllers"

	"github.com/gin-gonic/gin"
)

// SwordsmithsRoutes 刀匠管理路由
type SwordsmithsRoutes struct{}

// NewSwordsmithsRoutes 创建刀匠管理路由实例
func NewSwordsmithsRoutes() *SwordsmithsRoutes {
	return &SwordsmithsRoutes{}
}

// RegisterRoutes 注册刀匠管理路由（需要认证）
func (r *SwordsmithsRoutes) RegisterRoutes(protected *gin.RouterGroup, ctrls *controllers.Controllers) {
	swordsmiths := protected.Group(constants.SwordsmithsGroup)
	{
		swordsmiths.POST(constants.CreateSwordsmith, ctrls.Swordsmith.CreateSwordsmith)
		swordsmiths.GET(constants.GetSwordsmiths, ctrls.Swordsmith.GetSwordsmiths)
		swordsmiths.GET(constants.GetSwordsmithByID, ctrls.Swordsmith.GetSwordsmith)
		swordsmiths.POST(constants.UpdateSwordsmith, ctrls.Swordsmith.UpdateSwordsmith)
		swordsmiths.GET(constants.DeleteSwordsmithByID, ctrls.Swordsmith.DeleteSwordsmith)
	}
}
