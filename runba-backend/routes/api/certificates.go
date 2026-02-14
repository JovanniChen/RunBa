package api

import (
	"steam-backend/routes/constants"
	"steam-backend/routes/controllers"

	"github.com/gin-gonic/gin"
)

// CertificatesRoutes 证书管理路由
//
type CertificatesRoutes struct{}

// NewCertificatesRoutes 创建证书管理路由实例
func NewCertificatesRoutes() *CertificatesRoutes {
	return &CertificatesRoutes{}
}

// RegisterRoutes 注册证书管理路由（需要认证）
func (r *CertificatesRoutes) RegisterRoutes(protected *gin.RouterGroup, ctrls *controllers.Controllers) {
	certificates := protected.Group(constants.CertificatesGroup)
	{
		certificates.POST(constants.CreateCertificate, ctrls.Certificate.CreateCertificate)
		certificates.GET(constants.GetCertificates, ctrls.Certificate.GetCertificates)
		certificates.GET(constants.GetCertificateByID, ctrls.Certificate.GetCertificate)
		certificates.POST(constants.UpdateCertificate, ctrls.Certificate.UpdateCertificate)
		certificates.GET(constants.DeleteCertificateByID, ctrls.Certificate.DeleteCertificate)
	}
}
