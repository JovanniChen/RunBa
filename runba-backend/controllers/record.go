package controllers

import (
	"steam-backend/constants"
	"steam-backend/models"
	"steam-backend/services"
	"steam-backend/utils"

	"github.com/gin-gonic/gin"
)

// RecordController 兑换记录控制器结构体
type RecordController struct {
	recordService         services.RecordServiceInterface
	activationCodeService services.ActivationCodeServiceInterface
}

// NewRecordController 创建新的兑换记录控制器实例
func NewRecordController() *RecordController {
	return &RecordController{
		recordService:         services.NewRecordService(),
		activationCodeService: services.NewActivationCodeService(),
	}
}

func (ctrl *RecordController) CreateExchangeRecord(c *gin.Context) {
	var req models.CreateExchangeRecordRequest

	// 绑定请求参数
	if err := utils.BindRequest(c, &req); err != nil {
		utils.ErrorResponseWithCode(c, constants.CodeParamFailure, "参数验证失败", err)
		return
	}

	// 调用服务层创建兑换记录
	record, err := ctrl.recordService.CreateExchangeRecord(&req)
	if err != nil {
		utils.ErrorResponseWithCode(c, constants.CodeCreateExchangeRecordFailure, "创建兑换记录失败", err)
		return
	}

	utils.SuccessResponse(c, constants.CodeSuccess, "兑换请求已提交，正在处理中", record)
}

func (ctrl *RecordController) GetExchangeRecords(c *gin.Context) {
	var params models.ExchangeRecordListRequest

	// 绑定查询参数
	if err := c.ShouldBindQuery(&params); err != nil {
		utils.ErrorResponse(c, constants.CodeParamFailure, "参数验证失败", err.Error())
		return
	}

	// 调用服务层获取兑换记录列表
	response, err := ctrl.recordService.GetExchangeRecords(&params)
	if err != nil {
		utils.ErrorResponse(c, constants.CodeGetExchangeRecordsFailure, "获取兑换记录列表失败", err.Error())
		return
	}

	utils.SuccessResponse(c, constants.CodeSuccess, "获取兑换记录列表成功", response)
}

func (ctrl *RecordController) GetExchangeRecord(c *gin.Context) {
	var req models.GetExchangeRecord

	// 绑定请求参数
	if err := utils.BindRequest(c, &req); err != nil {
		utils.ErrorResponse(c, constants.CodeParamFailure, "请求参数错误", err.Error())
		return
	}

	// 调用服务层获取兑换记录信息（默认包含礼物记录）
	record, err := ctrl.recordService.GetExchangeRecordByID(uint(req.ID))
	if err != nil {
		utils.ErrorResponse(c, constants.CodeGetExchangeRecordFailure, "获取兑换记录信息失败", err.Error())
		return
	}

	utils.SuccessResponse(c, constants.CodeSuccess, "获取兑换记录信息成功", record)
}
