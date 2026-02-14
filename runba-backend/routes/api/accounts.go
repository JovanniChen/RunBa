package api

import (
	"steam-backend/routes/constants"
	"steam-backend/routes/controllers"

	"github.com/gin-gonic/gin"
)

// AccountsRoutes Steam账号管理路由
type AccountsRoutes struct{}

// NewAccountsRoutes 创建账号管理路由实例
func NewAccountsRoutes() *AccountsRoutes {
	return &AccountsRoutes{}
}

// RegisterRoutes 注册Steam账号管理路由（需要认证）
func (r *AccountsRoutes) RegisterRoutes(protected *gin.RouterGroup, ctrls *controllers.Controllers) {
	// Steam 账号管理相关操作
	accounts := protected.Group(constants.AccountsGroup)
	{
		// 基本CRUD操作
		accounts.POST(constants.CreateAccount, ctrls.Account.CreateAccount)             // 创建账号
		accounts.GET(constants.GetAccounts, ctrls.Account.GetAccounts)                  // 获取账号列表
		accounts.GET(constants.GetAccountByID, ctrls.Account.GetAccount)                // 获取单个账号信息
		accounts.POST(constants.UpdateAccountStatus, ctrls.Account.UpdateAccountStatus) // 更新账号状态
		accounts.GET(constants.DeleteAccountByID, ctrls.Account.DeleteAccount)          // 删除账号（软删除）
	}
}
