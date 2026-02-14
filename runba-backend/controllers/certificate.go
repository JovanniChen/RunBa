package controllers

import (
	"steam-backend/constants"
	"steam-backend/models"
	"steam-backend/services"
	"steam-backend/utils"

	"github.com/gin-gonic/gin"
)

// CertificateController 证书控制器
// 处理证书相关的HTTP请求
type CertificateController struct {
	certificateService *services.CertificateService
}

// NewCertificateController 创建新的证书控制器实例
func NewCertificateController(certificateService *services.CertificateService) *CertificateController {
	return &CertificateController{
		certificateService: certificateService,
	}
}

// CreateCertificate 创建证书
func (ctrl *CertificateController) CreateCertificate(c *gin.Context) {
	var req models.CreateCertificateRequest
	if err := utils.BindRequest(c, &req); err != nil {
		utils.ErrorResponseWithCode(c, constants.CodeParamFailure, "参数校验失败", err)
		return
	}

	certificate, err := ctrl.certificateService.CreateCertificate(&req)
	if err != nil {
		utils.ErrorResponseWithCode(c, constants.CodeCreateCertificateFailure, "创建证书失败: "+err.Error(), err)
		return
	}

	utils.SuccessResponse(c, constants.CodeSuccess, "创建证书成功", certificate)
}

// GetCertificates 获取证书列表
func (ctrl *CertificateController) GetCertificates(c *gin.Context) {
	var params models.CertificateListRequest
	if err := c.ShouldBindQuery(&params); err != nil {
		utils.ErrorResponse(c, constants.CodeParamFailure, "参数校验失败", err.Error())
		return
	}

	response, err := ctrl.certificateService.GetCertificates(&params)
	if err != nil {
		utils.ErrorResponse(c, constants.CodeGetCertificatesFailure, "获取证书列表失败", err.Error())
		return
	}

	utils.SuccessResponse(c, constants.CodeSuccess, "获取证书列表成功", response)
}

// GetCertificate 获取证书详情
func (ctrl *CertificateController) GetCertificate(c *gin.Context) {
	var req models.GetCertificateRequest
	if err := utils.BindRequest(c, &req); err != nil {
		utils.ErrorResponseWithCode(c, constants.CodeParamFailure, "参数校验失败", err)
		return
	}

	certificate, err := ctrl.certificateService.GetCertificateByID(req.ID)
	if err != nil {
		utils.ErrorResponse(c, constants.CodeGetCertificateFailure, "获取证书信息失败", err.Error())
		return
	}

	utils.SuccessResponse(c, constants.CodeSuccess, "获取证书信息成功", certificate)
}

// GetCertificateByNoPublic 公共查询证书（按证书序号）
func (ctrl *CertificateController) GetCertificateByNoPublic(c *gin.Context) {
	var req models.GetCertificateByNoRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ErrorResponse(c, constants.CodeParamFailure, "参数校验失败", err.Error())
		return
	}

	certificate, err := ctrl.certificateService.GetCertificateByNo(req.No)
	if err != nil {
		utils.ErrorResponse(c, constants.CodeGetCertificateFailure, "获取证书信息失败", err.Error())
		return
	}

	utils.SuccessResponse(c, constants.CodeSuccess, "获取证书信息成功", certificate)
}

// UpdateCertificate 更新证书
func (ctrl *CertificateController) UpdateCertificate(c *gin.Context) {
	var req models.UpdateCertificateRequest
	if err := utils.BindRequest(c, &req); err != nil {
		utils.ErrorResponse(c, constants.CodeParamFailure, "参数校验失败", err.Error())
		return
	}

	certificate, err := ctrl.certificateService.UpdateCertificate(req.ID, &req)
	if err != nil {
		utils.ErrorResponse(c, constants.CodeUpdateCertificateFailure, "更新证书失败", err.Error())
		return
	}

	utils.SuccessResponse(c, constants.CodeSuccess, "更新证书成功", certificate)
}

// DeleteCertificate 删除证书
func (ctrl *CertificateController) DeleteCertificate(c *gin.Context) {
	var req models.DeleteCertificateRequest
	if err := utils.BindRequest(c, &req); err != nil {
		utils.ErrorResponseWithCode(c, constants.CodeParamFailure, "参数校验失败", err)
		return
	}

	if err := ctrl.certificateService.DeleteCertificate(req.ID); err != nil {
		utils.ErrorResponse(c, constants.CodeDeleteCertificateFailure, "删除证书失败", err.Error())
		return
	}

	utils.SuccessResponse(c, constants.CodeSuccess, "删除证书成功", nil)
}
