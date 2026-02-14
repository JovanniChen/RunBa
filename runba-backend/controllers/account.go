package controllers

import (
	"fmt"
	"steam-backend/constants"
	"steam-backend/models"
	"steam-backend/services"
	"steam-backend/utils"

	"github.com/gin-gonic/gin"
)

// AccountController Steam账户控制器结构体
// 处理Steam账户相关的HTTP请求
type AccountController struct {
	accountService *services.AccountService // 直接使用具体类型
}

// NewAccountController 创建新的账户控制器实例
func NewAccountController(accountService *services.AccountService) *AccountController {
	return &AccountController{
		accountService: accountService,
	}
}

func (ctrl *AccountController) CreateAccount(c *gin.Context) {
	var req models.CreateAccountRequest

	// 绑定请求参数
	if err := utils.BindRequest(c, &req); err != nil {
		utils.ErrorResponseWithCode(c, constants.CodeParamFailure, "参数校验失败", err)
		return
	}

	// 调用服务层创建账户
	account, err := ctrl.accountService.CreateAccount(&req)
	if err != nil {
		utils.ErrorResponseWithCode(c, constants.CodeCreateAccountFailure, "创建钱包失败:,"+err.Error(), err)
		return
	}

	utils.SuccessResponse(c, constants.CodeSuccess, "创建钱包成功", account)
}

func (ctrl *AccountController) GetAccounts(c *gin.Context) {
	var params models.AccountListRequest

	// 绑定查询参数
	if err := c.ShouldBindQuery(&params); err != nil {
		utils.ErrorResponse(c, constants.CodeParamFailure, "参数校验失败", err.Error())
		return
	}

	// 调用服务层获取账户列表
	response, err := ctrl.accountService.GetAccounts(&params)
	if err != nil {
		utils.ErrorResponse(c, constants.CodeGetAccountsFailure, "获取钱包列表失败", err.Error())
		return
	}

	utils.SuccessResponse(c, constants.CodeSuccess, "获取钱包列表成功", response)
}

func (ctrl *AccountController) GetAccount(c *gin.Context) {
	var req models.GetAccountRequest

	// 绑定请求参数
	if err := utils.BindRequest(c, &req); err != nil {
		utils.ErrorResponseWithCode(c, constants.CodeParamFailure, "参数校验失败", err)
		return
	}

	// 调用服务层获取账户信息
	account, err := ctrl.accountService.GetAccountByID(req.ID)
	if err != nil {
		utils.ErrorResponse(c, constants.CodeGetAccountFailure, "获取钱包信息失败", err.Error())
		return
	}

	utils.SuccessResponse(c, constants.CodeSuccess, "获取Steam账户信息成功", account)
}

func (ctrl *AccountController) UpdateAccount(c *gin.Context) {
	var req models.UpdateAccountRequest
	// 绑定请求参数
	if err := utils.BindRequest(c, &req); err != nil {
		utils.ErrorResponse(c, constants.CodeParamFailure, "参数校验失败", err.Error())
		return
	}

	// 调用服务层更新账户信息
	account, err := ctrl.accountService.UpdateAccount(req.ID, &req)
	if err != nil {
		utils.ErrorResponse(c, constants.CodeUpdateAccountFailure, "更新钱包信息失败", err.Error())
		return
	}

	utils.SuccessResponse(c, constants.CodeSuccess, "更新Steam账户信息成功", account)
}

func (ctrl *AccountController) DeleteAccount(c *gin.Context) {
	var req models.DeleteAccountRequest

	// 绑定请求参数
	if err := utils.BindRequest(c, &req); err != nil {
		utils.ErrorResponseWithCode(c, constants.CodeParamFailure, "参数校验失败", err)
		return
	}

	// 调用服务层删除账户
	err := ctrl.accountService.DeleteAccount(req.ID)
	if err != nil {
		utils.ErrorResponse(c, constants.CodeDeleteAccountFailure, "删除Steam账户失败", err.Error())
		return
	}

	utils.SuccessResponse(c, constants.CodeSuccess, "删除Steam账户成功", nil)
}

func (ctrl *AccountController) UpdateAccountStatus(c *gin.Context) {
	var req models.UpdateAccountStatusRequest
	// 绑定请求参数
	if err := utils.BindRequest(c, &req); err != nil {
		utils.ErrorResponse(c, constants.CodeParamFailure, "参数验证失败", err.Error())
		return
	}

	fmt.Println("req.Status", req.Status)

	// 调用服务层更新账户状态（直接设置，覆盖所有状态）
	account, err := ctrl.accountService.SetAccountStatusSafe(req.ID, req.Status)
	if err != nil {
		utils.ErrorResponse(c, constants.CodeUpdateAccountStatusFailure, "更新账户状态失败", err.Error())
		return
	}

	utils.SuccessResponse(c, constants.CodeSuccess, "更新账户状态成功", account)
}

// AddAccountStatus 添加账户状态
func (ctrl *AccountController) AddAccountStatus(c *gin.Context) {
	var req models.AddAccountStatusRequest
	// 绑定请求参数
	if err := utils.BindRequest(c, &req); err != nil {
		utils.ErrorResponse(c, constants.CodeParamFailure, "参数验证失败", err.Error())
		return
	}

	// 调用服务层添加账户状态
	account, err := ctrl.accountService.AddAccountStatus(req.ID, req.Status)
	if err != nil {
		utils.ErrorResponse(c, constants.CodeUpdateAccountStatusFailure, "添加账户状态失败", err.Error())
		return
	}

	utils.SuccessResponse(c, constants.CodeSuccess, "添加账户状态成功", account)
}

// RemoveAccountStatus 移除账户状态
func (ctrl *AccountController) RemoveAccountStatus(c *gin.Context) {
	var req models.RemoveAccountStatusRequest
	// 绑定请求参数
	if err := utils.BindRequest(c, &req); err != nil {
		utils.ErrorResponse(c, constants.CodeParamFailure, "参数验证失败", err.Error())
		return
	}

	// 调用服务层移除账户状态
	account, err := ctrl.accountService.RemoveAccountStatus(req.ID, req.Status)
	if err != nil {
		utils.ErrorResponse(c, constants.CodeUpdateAccountStatusFailure, "移除账户状态失败", err.Error())
		return
	}

	utils.SuccessResponse(c, constants.CodeSuccess, "移除账户状态成功", account)
}

// SetAccountStatusSafe 安全设置账户状态
func (ctrl *AccountController) SetAccountStatusSafe(c *gin.Context) {
	var req models.UpdateAccountStatusRequest
	// 绑定请求参数
	if err := utils.BindRequest(c, &req); err != nil {
		utils.ErrorResponse(c, constants.CodeParamFailure, "参数验证失败", err.Error())
		return
	}

	// 调用服务层安全设置账户状态
	account, err := ctrl.accountService.SetAccountStatusSafe(req.ID, req.Status)
	if err != nil {
		utils.ErrorResponse(c, constants.CodeUpdateAccountStatusFailure, "设置账户状态失败", err.Error())
		return
	}

	utils.SuccessResponse(c, constants.CodeSuccess, "设置账户状态成功", account)
}
