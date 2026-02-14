package controllers

import (
	"steam-backend/constants"
	"steam-backend/middleware"
	"steam-backend/models"
	"steam-backend/services"
	"steam-backend/utils"

	"github.com/gin-gonic/gin"
)

// UserController 统一用户控制器
// 包含用户认证和用户管理的所有功能
type UserController struct {
	userService *services.UserService
}

// NewUserController 创建新的用户控制器实例
func NewUserController() *UserController {
	return &UserController{
		userService: services.NewUserService(),
	}
}

// ======================
// 认证相关功能
// ======================

// Login 用户登录接口
// POST /api/v1/auth/login
func (ctrl *UserController) Login(c *gin.Context) {
	var loginRequest models.LoginRequest

	// 绑定请求数据 - 支持 JSON 和 urlencoded 格式
	if err := utils.BindRequest(c, &loginRequest); err != nil {
		utils.ErrorResponseWithCode(c, constants.CodeParamFailure, "请求数据格式错误", err)
		return
	}

	// 调用服务层进行登录
	loginResponse, err := ctrl.userService.LoginUser(&loginRequest)
	if err != nil {
		// 根据错误类型返回不同的状态码
		if err.Error() == "用户不存在" || err.Error() == "密码错误" {
			utils.ErrorResponse(c, constants.CodeUsernameOrPasswordFailure, "用户名或密码错误", "")
		} else {
			utils.ErrorResponseWithCode(c, constants.CodeLoginFailure, "登录失败", err)
		}
		return
	}

	// 登录成功
	utils.SuccessResponse(c, constants.CodeSuccess, "登录成功", loginResponse)
}

// GetProfile 获取当前用户信息
// GET /api/v1/auth/profile
func (ctrl *UserController) GetProfile(c *gin.Context) {
	// 从JWT中间件获取用户ID
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, constants.CodeUnauthorized, "未找到用户信息", "")
		return
	}

	// 调用服务层获取用户信息
	userInfoResponse, err := ctrl.userService.GetUserProfile(userID)
	if err != nil {
		if err.Error() == "用户不存在" {
			utils.ErrorResponse(c, constants.CodeUserNoExists, err.Error(), "")
		} else {
			utils.ErrorResponseWithCode(c, constants.CodeGetUserFailure, "获取用户信息失败", err)
		}
		return
	}

	// 返回用户信息
	utils.SuccessResponse(c, constants.CodeSuccess, "获取用户信息成功", userInfoResponse)
}

// UpdateProfile 更新当前用户信息
// PUT /api/v1/auth/profile
func (ctrl *UserController) UpdateProfile(c *gin.Context) {
	// 从JWT中间件获取用户ID
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, constants.CodeUnauthorized, "未找到用户信息", "")
		return
	}

	var updateRequest models.UpdateUserRequest

	// 绑定请求数据
	if err := utils.BindRequest(c, &updateRequest); err != nil {
		utils.ErrorResponseWithCode(c, constants.CodeParamFailure, "请求数据格式错误", err)
		return
	}

	// 调用服务层更新用户信息
	userInfoResponse, err := ctrl.userService.UpdateUserProfile(userID, &updateRequest)
	if err != nil {
		if err.Error() == "用户不存在" {
			utils.ErrorResponse(c, constants.CodeUserNoExists, err.Error(), "")
		} else {
			utils.ErrorResponseWithCode(c, constants.CodeUpdateUserFailure, "更新用户信息失败", err)
		}
		return
	}

	// 返回更新后的用户信息
	utils.SuccessResponse(c, constants.CodeSuccess, "用户信息更新成功", userInfoResponse)
}

// ChangePassword 修改密码
// POST /api/v1/auth/change-password
func (ctrl *UserController) ChangePassword(c *gin.Context) {
	// 从JWT中间件获取用户ID
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, constants.CodeUnauthorized, "未找到用户信息", "")
		return
	}

	var changePasswordRequest models.ChangePasswordRequest

	// 绑定请求数据
	if err := utils.BindRequest(c, &changePasswordRequest); err != nil {
		utils.ErrorResponseWithCode(c, constants.CodeParamFailure, "请求数据格式错误", err)
		return
	}

	// 验证新密码确认
	if changePasswordRequest.NewPassword != changePasswordRequest.ConfirmNewPassword {
		utils.ErrorResponse(c, constants.CodeRepeatPasswordFailure, "新密码确认不匹配", "")
		return
	}

	// 调用服务层修改密码
	err := ctrl.userService.ChangePassword(userID, &changePasswordRequest)
	if err != nil {
		if err.Error() == "原密码错误" {
			utils.ErrorResponse(c, constants.CodePasswordFailure, err.Error(), "")
		} else if err.Error() == "用户不存在" {
			utils.ErrorResponse(c, constants.CodeUserNoExists, err.Error(), "")
		} else {
			utils.ErrorResponseWithCode(c, constants.CodeModifyPasswordFailure, "密码修改失败", err)
		}
		return
	}

	// 密码修改成功
	utils.SuccessResponse(c, constants.CodeSuccess, "密码修改成功", nil)
}

// ======================
// 用户管理功能（管理员）
// ======================

// GetUsers 获取用户列表
// GET /api/v1/users
func (ctrl *UserController) GetUsers(c *gin.Context) {
	var params models.UserListRequest

	// 绑定查询参数
	if err := c.ShouldBindQuery(&params); err != nil {
		utils.ErrorResponseWithCode(c, constants.CodeParamFailure, "查询参数错误", err)
		return
	}

	// 调用Service层获取用户列表
	userListResponse, err := ctrl.userService.GetUsers(&params)
	if err != nil {
		utils.ErrorResponse(c, constants.CodeGetUsersFailure, err.Error(), "获取用户列表失败")
		return
	}

	// 返回响应
	utils.SuccessResponse(c, constants.CodeSuccess, "获取用户列表成功", userListResponse)
}

// GetUser 获取单个用户信息
// GET /api/v1/users/:id
func (ctrl *UserController) GetUser(c *gin.Context) {
	var req models.GetUserRequest

	// 绑定请求参数
	if err := utils.BindRequest(c, &req); err != nil {
		utils.ErrorResponse(c, constants.CodeParamFailure, "请求参数错误", err.Error())
		return
	}

	// 调用Service层获取用户信息
	userInfoResponse, err := ctrl.userService.GetUserByID(req.ID)
	if err != nil {
		if err.Error() == "用户不存在" {
			utils.ErrorResponse(c, constants.CodeUserNoExists, err.Error(), "")
		} else {
			utils.ErrorResponse(c, constants.CodeGetUserFailure, err.Error(), "获取用户信息失败")
		}
		return
	}

	// 返回响应
	utils.SuccessResponse(c, constants.CodeSuccess, "获取用户信息成功", userInfoResponse)
}

// UpdateUser 更新用户信息（管理员功能）
// PUT /api/v1/users/:id
func (ctrl *UserController) UpdateUser(c *gin.Context) {
	var req models.UpdateUserRequest

	// 绑定请求数据
	if err := utils.BindRequest(c, &req); err != nil {
		utils.ErrorResponseWithCode(c, constants.CodeParamFailure, "请求数据格式错误", err)
		return
	}

	// 调用Service层更新用户信息
	userInfoResponse, err := ctrl.userService.UpdateUser(req.ID, &req)
	if err != nil {
		if err.Error() == "用户不存在" {
			utils.ErrorResponse(c, constants.CodeUserNoExists, err.Error(), "")
		} else {
			utils.ErrorResponse(c, constants.CodeUpdateUserFailure, err.Error(), "更新用户信息失败")
		}
		return
	}

	// 返回响应
	utils.SuccessResponse(c, constants.CodeSuccess, "用户信息更新成功", userInfoResponse)
}

// DeleteUser 删除用户（软删除）
// DELETE /api/v1/users/:id
func (ctrl *UserController) DeleteUser(c *gin.Context) {
	var req models.DeleteUserRequest

	// 绑定请求参数
	if err := utils.BindRequest(c, &req); err != nil {
		utils.ErrorResponse(c, constants.CodeParamFailure, "请求参数错误", err.Error())
		return
	}

	// 调用Service层删除用户
	err := ctrl.userService.DeleteUser(uint(req.ID))
	if err != nil {
		if err.Error() == "用户不存在" {
			utils.ErrorResponse(c, constants.CodeUserNoExists, err.Error(), "")
		} else {
			utils.ErrorResponse(c, constants.CodeDeleteUserFailure, err.Error(), "删除用户失败")
		}
		return
	}

	// 返回响应
	utils.SuccessResponse(c, constants.CodeSuccess, "用户删除成功", nil)
}

// UpdateUserStatus 更新用户状态
// PATCH /api/v1/users/:id/status
func (ctrl *UserController) UpdateUserStatus(c *gin.Context) {
	var req models.UpdateUserStatusRequest

	// 绑定请求数据
	if err := utils.BindRequest(c, &req); err != nil {
		utils.ErrorResponseWithCode(c, constants.CodeParamFailure, "请求数据格式错误", err)
		return
	}

	// 调用Service层更新用户状态
	userInfo, err := ctrl.userService.UpdateUserStatus(uint(req.ID), req.Status)
	if err != nil {
		if err.Error() == "用户不存在" {
			utils.ErrorResponse(c, constants.CodeUserNoExists, err.Error(), "")
		} else {
			utils.ErrorResponse(c, constants.CodeUpdateUserStatusFailure, err.Error(), "更新用户状态失败")
		}
		return
	}

	// 返回响应
	utils.SuccessResponse(c, constants.CodeSuccess, "用户状态更新成功", userInfo)
}
