package controllers

import (
	"steam-backend/constants"
	"steam-backend/models"
	"steam-backend/services"
	"steam-backend/utils"

	"github.com/gin-gonic/gin"
)

// SwordsmithController 刀匠控制器
type SwordsmithController struct {
	swordsmithService *services.SwordsmithService
}

// NewSwordsmithController 创建新的刀匠控制器
func NewSwordsmithController(swordsmithService *services.SwordsmithService) *SwordsmithController {
	return &SwordsmithController{swordsmithService: swordsmithService}
}

// CreateSwordsmith 创建刀匠
func (ctrl *SwordsmithController) CreateSwordsmith(c *gin.Context) {
	var req models.CreateSwordsmithRequest
	if err := utils.BindRequest(c, &req); err != nil {
		utils.ErrorResponseWithCode(c, constants.CodeParamFailure, "参数校验失败", err)
		return
	}

	swordsmith, err := ctrl.swordsmithService.CreateSwordsmith(&req)
	if err != nil {
		utils.ErrorResponseWithCode(c, constants.CodeCreateSwordsmithFailure, "创建刀匠失败: "+err.Error(), err)
		return
	}

	utils.SuccessResponse(c, constants.CodeSuccess, "创建刀匠成功", swordsmith)
}

// GetSwordsmiths 获取刀匠列表
func (ctrl *SwordsmithController) GetSwordsmiths(c *gin.Context) {
	var params models.SwordsmithListRequest
	if err := c.ShouldBindQuery(&params); err != nil {
		utils.ErrorResponse(c, constants.CodeParamFailure, "参数校验失败", err.Error())
		return
	}

	response, err := ctrl.swordsmithService.GetSwordsmiths(&params)
	if err != nil {
		utils.ErrorResponse(c, constants.CodeGetSwordsmithsFailure, "获取刀匠列表失败", err.Error())
		return
	}

	utils.SuccessResponse(c, constants.CodeSuccess, "获取刀匠列表成功", response)
}

// GetSwordsmith 获取刀匠详情
func (ctrl *SwordsmithController) GetSwordsmith(c *gin.Context) {
	var req models.GetSwordsmithRequest
	if err := utils.BindRequest(c, &req); err != nil {
		utils.ErrorResponseWithCode(c, constants.CodeParamFailure, "参数校验失败", err)
		return
	}

	swordsmith, err := ctrl.swordsmithService.GetSwordsmithByID(req.ID)
	if err != nil {
		utils.ErrorResponse(c, constants.CodeGetSwordsmithFailure, "获取刀匠信息失败", err.Error())
		return
	}

	utils.SuccessResponse(c, constants.CodeSuccess, "获取刀匠信息成功", swordsmith)
}

// UpdateSwordsmith 更新刀匠
func (ctrl *SwordsmithController) UpdateSwordsmith(c *gin.Context) {
	var req models.UpdateSwordsmithRequest
	if err := utils.BindRequest(c, &req); err != nil {
		utils.ErrorResponse(c, constants.CodeParamFailure, "参数校验失败", err.Error())
		return
	}

	swordsmith, err := ctrl.swordsmithService.UpdateSwordsmith(req.ID, &req)
	if err != nil {
		utils.ErrorResponse(c, constants.CodeUpdateSwordsmithFailure, "更新刀匠失败", err.Error())
		return
	}

	utils.SuccessResponse(c, constants.CodeSuccess, "更新刀匠成功", swordsmith)
}

// DeleteSwordsmith 删除刀匠
func (ctrl *SwordsmithController) DeleteSwordsmith(c *gin.Context) {
	var req models.DeleteSwordsmithRequest
	if err := utils.BindRequest(c, &req); err != nil {
		utils.ErrorResponseWithCode(c, constants.CodeParamFailure, "参数校验失败", err)
		return
	}

	if err := ctrl.swordsmithService.DeleteSwordsmith(req.ID); err != nil {
		utils.ErrorResponse(c, constants.CodeDeleteSwordsmithFailure, "删除刀匠失败", err.Error())
		return
	}

	utils.SuccessResponse(c, constants.CodeSuccess, "删除刀匠成功", nil)
}
