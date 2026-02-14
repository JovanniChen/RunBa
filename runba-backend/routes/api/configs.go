package api

import (
	"steam-backend/routes/constants"
	"steam-backend/routes/controllers"

	"github.com/gin-gonic/gin"
)

// ConfigsRoutes 系统配置管理路由
type ConfigsRoutes struct{}

// NewConfigsRoutes 创建配置管理路由实例
func NewConfigsRoutes() *ConfigsRoutes {
	return &ConfigsRoutes{}
}

// RegisterRoutes 注册系统配置管理路由（需要认证）
func (r *ConfigsRoutes) RegisterRoutes(protected *gin.RouterGroup, ctrls *controllers.Controllers) {
	// 系统配置管理相关操作
	confs := protected.Group(constants.ConfsGroup)
	{
		// 推荐使用的创建或更新配置接口
		confs.POST(constants.CreateOrUpdateConf, ctrls.Conf.CreateOrUpdateConf) // 创建或更新配置（推荐使用）
		confs.GET(constants.GetConf, ctrls.Conf.GetLatestConf)                  // 获取最新配置
	}
}
