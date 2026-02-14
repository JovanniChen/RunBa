package controllers

import (
	"steam-backend/constants"
	"steam-backend/models"
	"steam-backend/services"
	"steam-backend/utils"

	"github.com/gin-gonic/gin"
)

// ProxyController 代理控制器结构体
type ProxyController struct {
	proxyService *services.ProxyService
}

// NewProxyController 创建新的代理控制器实例
func NewProxyController() *ProxyController {
	return &ProxyController{
		proxyService: services.NewProxyService(),
	}
}

func (c *ProxyController) CreateProxy(ctx *gin.Context) {
	var req models.CreateProxyRequest

	// 绑定请求参数
	if err := utils.BindRequest(ctx, &req); err != nil {
		utils.ErrorResponse(ctx, constants.CodeParamFailure, "请求参数错误", err.Error())
		return
	}

	// 创建代理
	proxy, err := c.proxyService.CreateProxy(&req)
	if err != nil {
		utils.ErrorResponse(ctx, constants.CodeCreateProxyFailure, "创建代理失败", err.Error())
		return
	}

	utils.SuccessResponse(ctx, constants.CodeSuccess, "代理创建成功", proxy)
}

func (c *ProxyController) GetProxy(ctx *gin.Context) {
	var req models.GetAccountRequest

	// 绑定请求参数
	if err := utils.BindRequest(ctx, &req); err != nil {
		utils.ErrorResponseWithCode(ctx, constants.CodeParamFailure, "参数校验失败", err)
		return
	}

	// 获取代理
	proxy, err := c.proxyService.GetProxyByID(req.ID)
	if err != nil {
		utils.ErrorResponse(ctx, constants.CodeProxyNotExists, "代理不存在", err.Error())
		return
	}

	utils.SuccessResponse(ctx, constants.CodeSuccess, "获取代理成功", proxy)
}

func (c *ProxyController) GetProxies(ctx *gin.Context) {
	var req models.ProxyListRequest

	// 绑定查询参数
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.ErrorResponse(ctx, constants.CodeParamFailure, "请求参数错误", err.Error())
		return
	}

	// 获取代理列表
	response, err := c.proxyService.GetProxyList(&req)
	if err != nil {
		utils.ErrorResponse(ctx, constants.CodeGetProxiesFailure, "获取代理列表失败", err.Error())
		return
	}

	utils.SuccessResponse(ctx, constants.CodeSuccess, "获取代理列表成功", response)
}

func (c *ProxyController) UpdateProxy(ctx *gin.Context) {
	var req models.UpdateProxyRequest

	// 绑定请求参数
	if err := utils.BindRequest(ctx, &req); err != nil {
		utils.ErrorResponse(ctx, constants.CodeParamFailure, "请求参数错误", err.Error())
		return
	}

	// 更新代理
	proxy, err := c.proxyService.UpdateProxy(req.ID, &req)
	if err != nil {
		utils.ErrorResponse(ctx, constants.CodeUpdateProxyFailure, "更新代理失败", err.Error())
		return
	}

	utils.SuccessResponse(ctx, constants.CodeSuccess, "代理更新成功", proxy)
}

func (c *ProxyController) DeleteProxy(ctx *gin.Context) {
	var req models.DeleteProxyRequest

	// 绑定请求参数
	if err := utils.BindRequest(ctx, &req); err != nil {
		utils.ErrorResponse(ctx, constants.CodeParamFailure, "请求参数错误", err.Error())
		return
	}

	// 删除代理
	err := c.proxyService.DeleteProxy(req.ID)
	if err != nil {
		utils.ErrorResponse(ctx, constants.CodeDeleteProxyFailure, "删除代理失败", err.Error())
		return
	}

	utils.SuccessResponse(ctx, constants.CodeSuccess, "代理删除成功", nil)
}

func (c *ProxyController) UpdateProxyStatus(ctx *gin.Context) {
	var req models.UpdateProxyStatusRequest

	// 绑定请求参数
	if err := utils.BindRequest(ctx, &req); err != nil {
		utils.ErrorResponse(ctx, constants.CodeParamFailure, "请求参数错误", err.Error())
		return
	}

	// 更新代理状态
	err := c.proxyService.UpdateProxyStatus(req.ID, req.Status)
	if err != nil {
		utils.ErrorResponse(ctx, constants.CodeUpdateProxyStatusFailure, "更新代理状态失败", err.Error())
		return
	}

	utils.SuccessResponse(ctx, constants.CodeSuccess, "代理状态更新成功", nil)
}
