package middleware

import (
	"steam-backend/middleware"

	"github.com/gin-gonic/gin"
)

// MiddlewareSetup 中间件配置
type MiddlewareSetup struct{}

// NewMiddlewareSetup 创建中间件配置实例
func NewMiddlewareSetup() *MiddlewareSetup {
	return &MiddlewareSetup{}
}

// SetupGlobal 设置全局中间件
func (m *MiddlewareSetup) SetupGlobal(router *gin.Engine) {
	// 添加全局中间件
	router.Use(middleware.LoggerMiddleware()) // 日志中间件
	router.Use(gin.Recovery())                // 恢复中间件，处理panic
	router.Use(middleware.CORSMiddleware())   // 跨域中间件
}

// SetupProtected 设置需要认证的中间件
func (m *MiddlewareSetup) SetupProtected(group *gin.RouterGroup) {
	group.Use(middleware.RequireJWT()) // 应用JWT认证中间件
}

// SetupAdmin 设置管理员权限中间件（可选，未来扩展）
func (m *MiddlewareSetup) SetupAdmin(group *gin.RouterGroup) {
	// 可以在这里添加管理员权限验证中间件
	// group.Use(middleware.RequireAdmin())
}
