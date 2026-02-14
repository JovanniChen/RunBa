package controllers

import (
	"steam-backend/container"
	"steam-backend/controllers"
)

// Controllers 控制器集合
// 统一管理所有控制器实例
type Controllers struct {
	User           *controllers.UserController
	Account        *controllers.AccountController
	Certificate    *controllers.CertificateController
	Forge          *controllers.ForgeController
	Swordsmith     *controllers.SwordsmithController
	Record         *controllers.RecordController
	ActivationCode *controllers.ActivationCodeController
	Conf           *controllers.ConfController
	Proxy          *controllers.ProxyController
}

// NewControllers 创建控制器集合
// 从容器中获取依赖并创建所有控制器实例
func NewControllers(appContainer *container.Container) *Controllers {
	return &Controllers{
		User:           controllers.NewUserController(),
		Account:        controllers.NewAccountController(appContainer.GetAccountService()),
		Certificate:    controllers.NewCertificateController(appContainer.GetCertificateService()),
		Forge:          controllers.NewForgeController(appContainer.GetForgeService()),
		Swordsmith:     controllers.NewSwordsmithController(appContainer.GetSwordsmithService()),
		ActivationCode: controllers.NewActivationCodeController(),
		Conf:           controllers.NewConfController(),
		Proxy:          controllers.NewProxyController(),
	}
}

// RouteRegistrar 路由注册器接口
// 定义路由模块必须实现的接口
type RouteRegistrar interface {
	RegisterRoutes(group interface{}, controllers *Controllers)
}
