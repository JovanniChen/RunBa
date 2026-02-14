package routes

import (
	"steam-backend/container"
	"steam-backend/routes/api"
	"steam-backend/routes/constants"
	"steam-backend/routes/controllers"
	"steam-backend/routes/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置所有路由
// 配置应用的所有API路由和中间件
func SetupRoutes(appContainer *container.Container) *gin.Engine {
	// 创建Gin路由器
	router := gin.New()

	// 设置全局中间件
	middlewareSetup := middleware.NewMiddlewareSetup()
	middlewareSetup.SetupGlobal(router)

	// 创建控制器实例
	ctrls := controllers.NewControllers(appContainer)

	// 设置系统路由
	setupSystemRoutes(router)

	// 设置API路由
	setupAPIRoutes(router, ctrls, middlewareSetup)

	return router
}

// setupSystemRoutes 设置系统路由
func setupSystemRoutes(router *gin.Engine) {
	// 健康检查路由（无需认证）
	router.GET(constants.Health, func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Steam Backend API is running",
		})
	})

	// API根路径信息
	router.GET(constants.Root, func(c *gin.Context) {
		c.JSON(200, gin.H{
			"name":        "Steam Backend API",
			"version":     "1.0.0",
			"description": "基于Gin + GORM + MySQL的Steam后端API",
			"health":      "GET /health",
			"api":         "GET /api/v1/*",
		})
	})
}

// setupAPIRoutes 设置API路由
func setupAPIRoutes(router *gin.Engine, ctrls *controllers.Controllers, middlewareSetup *middleware.MiddlewareSetup) {
	// API版本1路由组
	v1 := router.Group(constants.APIVersion1)

	// 注册认证相关路由（公开路由）
	authRoutes := api.NewAuthRoutes()
	authRoutes.RegisterRoutes(v1, ctrls)

	// 公开证书查询路由
	publicCertificates := v1.Group(constants.CertificatesGroup)
	publicCertificates.GET(constants.GetCertificateByNo, ctrls.Certificate.GetCertificateByNoPublic)

	// 设置需要JWT认证的路由组
	protected := v1.Group("/")
	middlewareSetup.SetupProtected(protected)

	authRoutes.RegisterProtectedRoutes(protected, ctrls)

	usersRoutes := api.NewUsersRoutes()
	usersRoutes.RegisterRoutes(protected, ctrls)

	certificatesRoutes := api.NewCertificatesRoutes()
	certificatesRoutes.RegisterRoutes(protected, ctrls)

	forgesRoutes := api.NewForgesRoutes()
	forgesRoutes.RegisterRoutes(protected, ctrls)

	swordsmithsRoutes := api.NewSwordsmithsRoutes()
	swordsmithsRoutes.RegisterRoutes(protected, ctrls)
}
