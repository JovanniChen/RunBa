package controllers

import (
	"steam-backend/constants"
	"steam-backend/models"
	"steam-backend/services"
	"steam-backend/utils"

	"github.com/gin-gonic/gin"
)

// ConfController 配置控制器结构体
type ConfController struct {
	confService *services.ConfService
}

// NewConfController 创建新的配置控制器实例
func NewConfController() *ConfController {
	return &ConfController{
		confService: services.NewConfService(),
	}
}

func (c *ConfController) CreateOrUpdateConf(ctx *gin.Context) {
	var req models.CreateConfRequest

	// 绑定请求参数
	if err := utils.BindRequest(ctx, &req); err != nil {
		utils.ErrorResponse(ctx, constants.CodeParamFailure, "请求参数错误", err.Error())
		return
	}

	// 创建或更新配置
	conf, err := c.confService.CreateOrUpdateConf(&req)
	if err != nil {
		utils.ErrorResponse(ctx, constants.CodeCreateOrUpdateConfFailure, "配置操作失败", err.Error())
		return
	}

	utils.SuccessResponse(ctx, constants.CodeSuccess, "配置操作成功", conf)
}

func (c *ConfController) GetLatestConf(ctx *gin.Context) {
	// 获取最新配置
	conf, err := c.confService.GetLatestConf()
	if err != nil {
		utils.ErrorResponse(ctx, constants.CodeConfNotExists, "获取最新配置失败", err.Error())
		return
	}

	utils.SuccessResponse(ctx, constants.CodeSuccess, "获取最新配置成功", conf)
}
