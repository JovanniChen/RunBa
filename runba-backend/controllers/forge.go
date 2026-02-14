package controllers

import (
	"steam-backend/constants"
	"steam-backend/models"
	"steam-backend/services"
	"steam-backend/utils"

	"github.com/gin-gonic/gin"
)

// ForgeController 锻刀所控制器
type ForgeController struct {
	forgeService *services.ForgeService
}

// NewForgeController 创建新的锻刀所控制器
func NewForgeController(forgeService *services.ForgeService) *ForgeController {
	return &ForgeController{forgeService: forgeService}
}

// CreateForge 创建锻刀所
func (ctrl *ForgeController) CreateForge(c *gin.Context) {
	var req models.CreateForgeRequest
	if err := utils.BindRequest(c, &req); err != nil {
		utils.ErrorResponseWithCode(c, constants.CodeParamFailure, "参数校验失败", err)
		return
	}

	forge, err := ctrl.forgeService.CreateForge(&req)
	if err != nil {
		utils.ErrorResponseWithCode(c, constants.CodeCreateForgeFailure, "创建锻刀所失败: "+err.Error(), err)
		return
	}

	utils.SuccessResponse(c, constants.CodeSuccess, "创建锻刀所成功", forge)
}

// GetForges 获取锻刀所列表
func (ctrl *ForgeController) GetForges(c *gin.Context) {
	var params models.ForgeListRequest
	if err := c.ShouldBindQuery(&params); err != nil {
		utils.ErrorResponse(c, constants.CodeParamFailure, "参数校验失败", err.Error())
		return
	}

	response, err := ctrl.forgeService.GetForges(&params)
	if err != nil {
		utils.ErrorResponse(c, constants.CodeGetForgesFailure, "获取锻刀所列表失败", err.Error())
		return
	}

	utils.SuccessResponse(c, constants.CodeSuccess, "获取锻刀所列表成功", response)
}

// GetForge 获取锻刀所详情
func (ctrl *ForgeController) GetForge(c *gin.Context) {
	var req models.GetForgeRequest
	if err := utils.BindRequest(c, &req); err != nil {
		utils.ErrorResponseWithCode(c, constants.CodeParamFailure, "参数校验失败", err)
		return
	}

	forge, err := ctrl.forgeService.GetForgeByID(req.ID)
	if err != nil {
		utils.ErrorResponse(c, constants.CodeGetForgeFailure, "获取锻刀所信息失败", err.Error())
		return
	}

	utils.SuccessResponse(c, constants.CodeSuccess, "获取锻刀所信息成功", forge)
}

// UpdateForge 更新锻刀所
func (ctrl *ForgeController) UpdateForge(c *gin.Context) {
	var req models.UpdateForgeRequest
	if err := utils.BindRequest(c, &req); err != nil {
		utils.ErrorResponse(c, constants.CodeParamFailure, "参数校验失败", err.Error())
		return
	}

	forge, err := ctrl.forgeService.UpdateForge(req.ID, &req)
	if err != nil {
		utils.ErrorResponse(c, constants.CodeUpdateForgeFailure, "更新锻刀所失败", err.Error())
		return
	}

	utils.SuccessResponse(c, constants.CodeSuccess, "更新锻刀所成功", forge)
}

// DeleteForge 删除锻刀所
func (ctrl *ForgeController) DeleteForge(c *gin.Context) {
	var req models.DeleteForgeRequest
	if err := utils.BindRequest(c, &req); err != nil {
		utils.ErrorResponseWithCode(c, constants.CodeParamFailure, "参数校验失败", err)
		return
	}

	if err := ctrl.forgeService.DeleteForge(req.ID); err != nil {
		utils.ErrorResponse(c, constants.CodeDeleteForgeFailure, "删除锻刀所失败", err.Error())
		return
	}

	utils.SuccessResponse(c, constants.CodeSuccess, "删除锻刀所成功", nil)
}
