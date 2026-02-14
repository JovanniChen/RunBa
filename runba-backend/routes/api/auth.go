package api

import (
	"steam-backend/routes/constants"
	"steam-backend/routes/controllers"

	"github.com/gin-gonic/gin"
)

// AuthRoutes 认证相关路由
type AuthRoutes struct{}

// NewAuthRoutes 创建认证路由实例
func NewAuthRoutes() *AuthRoutes {
	return &AuthRoutes{}
}

// RegisterRoutes 注册认证相关路由
func (r *AuthRoutes) RegisterRoutes(v1 *gin.RouterGroup, ctrls *controllers.Controllers) {
	// 公开认证路由（不需要JWT验证）
	auth := v1.Group(constants.AuthGroup)
	{
		auth.POST(constants.AuthLogin, ctrls.User.Login) // 用户登录
	}
}

// RegisterProtectedRoutes 注册需要认证的用户相关路由
func (r *AuthRoutes) RegisterProtectedRoutes(protected *gin.RouterGroup, ctrls *controllers.Controllers) {
	// 当前用户相关操作
	protected.GET(constants.AuthGroup+constants.AuthProfile, ctrls.User.GetProfile)             // 获取当前用户信息
	protected.PUT(constants.AuthGroup+constants.AuthProfile, ctrls.User.UpdateProfile)          // 更新当前用户信息
	protected.POST(constants.AuthGroup+constants.AuthChangePassword, ctrls.User.ChangePassword) // 修改密码
}
