package api

import (
	"steam-backend/routes/constants"
	"steam-backend/routes/controllers"

	"github.com/gin-gonic/gin"
)

// UsersRoutes 用户管理路由
type UsersRoutes struct{}

// NewUsersRoutes 创建用户管理路由实例
func NewUsersRoutes() *UsersRoutes {
	return &UsersRoutes{}
}

// RegisterRoutes 注册用户管理路由（需要认证）
func (r *UsersRoutes) RegisterRoutes(protected *gin.RouterGroup, ctrls *controllers.Controllers) {
	// 用户管理相关操作（管理员功能）
	users := protected.Group(constants.UsersGroup)
	{
		users.GET(constants.GetUsers, ctrls.User.GetUsers)         // 获取用户列表
		users.GET(constants.GetUserByID, ctrls.User.GetUser)       // 获取单个用户信息
		users.POST(constants.UpdateUser, ctrls.User.UpdateUser)    // 更新用户信息
		users.GET(constants.DeleteUserByID, ctrls.User.DeleteUser) // 删除用户（软删除）
		// 可选功能，根据需要启用
		// users.POST(constants.UpdateUserStatus, ctrls.User.UpdateUserStatus) // 更新用户状态
	}
}
