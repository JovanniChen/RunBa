package api

import (
	"steam-backend/routes/constants"
	"steam-backend/routes/controllers"

	"github.com/gin-gonic/gin"
)

// ProxiesRoutes 代理管理路由
type ProxiesRoutes struct{}

// NewProxiesRoutes 创建代理管理路由实例
func NewProxiesRoutes() *ProxiesRoutes {
	return &ProxiesRoutes{}
}

// RegisterRoutes 注册代理管理路由（需要认证）
func (r *ProxiesRoutes) RegisterRoutes(protected *gin.RouterGroup, ctrls *controllers.Controllers) {
	// 代理管理相关操作
	proxies := protected.Group(constants.ProxiesGroup)
	{
		// 基本CRUD操作
		proxies.POST(constants.CreateProxy, ctrls.Proxy.CreateProxy) // 创建代理
		proxies.GET(constants.GetProxies, ctrls.Proxy.GetProxies)    // 获取代理列表
		proxies.GET(constants.GetProxyByID, ctrls.Proxy.GetProxy)    // 获取单个代理信息
		proxies.POST(constants.UpdateProxy, ctrls.Proxy.UpdateProxy) // 更新代理信息
		proxies.GET(constants.DeleteProxy, ctrls.Proxy.DeleteProxy)  // 删除代理（软删除）
	}
}
