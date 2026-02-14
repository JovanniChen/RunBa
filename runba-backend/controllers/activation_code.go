package controllers

import (
	"steam-backend/constants"
	"steam-backend/models"
	"steam-backend/services"
	"steam-backend/utils"

	"github.com/gin-gonic/gin"
)

// ActivationCodeController 激活码控制器结构体
type ActivationCodeController struct {
	activationCodeService *services.ActivationCodeService
}

// NewActivationCodeController 创建新的激活码控制器实例
func NewActivationCodeController() *ActivationCodeController {
	return &ActivationCodeController{
		activationCodeService: services.NewActivationCodeService(),
	}
}

func (c *ActivationCodeController) CreateBatchActivationCodes(ctx *gin.Context) {
	var req models.CreateBatchActivationCodeRequest

	// 绑定请求参数
	if err := utils.BindRequest(ctx, &req); err != nil {
		utils.ErrorResponse(ctx, constants.CodeParamFailure, "请求参数错误", err.Error())
		return
	}

	// 批量创建激活码
	activationCodes, err := c.activationCodeService.CreateBatchActivationCodes(&req)
	if err != nil {
		utils.ErrorResponse(ctx, constants.CodeCreateBatchActivationCodesFailure, "批量创建激活码失败", err.Error())
		return
	}

	utils.SuccessResponse(ctx, constants.CodeSuccess, "批量创建激活码成功", activationCodes)
}

func (c *ActivationCodeController) GetActivationCodes(ctx *gin.Context) {
	var req models.ActivationCodeListRequest

	// 绑定查询参数
	if err := utils.BindRequest(ctx, &req); err != nil {
		utils.ErrorResponse(ctx, constants.CodeParamFailure, "请求参数错误", err.Error())
		return
	}

	// 获取激活码列表
	result, err := c.activationCodeService.GetActivationCodes(&req)
	if err != nil {
		utils.ErrorResponse(ctx, constants.CodeGetActivationCodesFailure, "获取激活码列表失败", err.Error())
		return
	}

	utils.SuccessResponse(ctx, constants.CodeSuccess, "获取激活码列表成功", result)
}

func (c *ActivationCodeController) GetActivationCode(ctx *gin.Context) {
	var req models.GetActivationCodeRequest
	// 绑定查询参数
	if err := utils.BindRequest(ctx, &req); err != nil {
		utils.ErrorResponse(ctx, constants.CodeParamFailure, "请求参数错误", err.Error())
		return
	}

	// 获取激活码信息
	activationCode, err := c.activationCodeService.GetActivationCodeByID(uint(req.ID))
	if err != nil {
		if err.Error() == "激活码不存在" {
			utils.ErrorResponse(ctx, constants.CodeActivationCodeNotExists, "激活码不存在", err.Error())
		} else {
			utils.ErrorResponse(ctx, constants.CodeGetActivationCodeFailure, "获取激活码失败", err.Error())
		}
		return
	}

	utils.SuccessResponse(ctx, constants.CodeSuccess, "获取激活码成功", activationCode)
}

func (c *ActivationCodeController) UpdateActivationCode(ctx *gin.Context) {
	var req models.UpdateActivationCodeRequest

	// 绑定请求参数
	if err := utils.BindRequest(ctx, &req); err != nil {
		utils.ErrorResponse(ctx, constants.CodeParamFailure, "请求参数错误", err.Error())
		return
	}

	// 更新激活码
	activationCode, err := c.activationCodeService.UpdateActivationCode(uint(req.ID), &req)
	if err != nil {
		if err.Error() == "激活码不存在" {
			utils.ErrorResponse(ctx, constants.CodeActivationCodeNotExists, "激活码不存在", err.Error())
		} else if err.Error() == "激活码已存在" {
			utils.ErrorResponse(ctx, constants.CodeGetActivationCodeExists, "激活码已存在", err.Error())
		} else {
			utils.ErrorResponse(ctx, constants.CodeUpdateActivationCodeFailure, "更新激活码失败", err.Error())
		}
		return
	}

	utils.SuccessResponse(ctx, constants.CodeSuccess, "激活码更新成功", activationCode)
}

func (c *ActivationCodeController) DeleteActivationCode(ctx *gin.Context) {
	var req models.DeleteActivationCodeRequest

	// 绑定请求参数
	if err := utils.BindRequest(ctx, &req); err != nil {
		utils.ErrorResponse(ctx, constants.CodeParamFailure, "请求参数错误", err.Error())
		return
	}

	// 删除激活码
	err := c.activationCodeService.DeleteActivationCode(uint(req.ID))
	if err != nil {
		if err.Error() == "激活码不存在" {
			utils.ErrorResponse(ctx, constants.CodeActivationCodeNotExists, "激活码不存在", err.Error())
		} else {
			utils.ErrorResponse(ctx, constants.CodeDeleteActivationCodeFailure, "删除激活码失败", err.Error())
		}
		return
	}

	utils.SuccessResponse(ctx, constants.CodeSuccess, "激活码删除成功", nil)
}
