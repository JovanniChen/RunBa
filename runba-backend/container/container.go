// 服务容器 - 统一依赖注入管理
//
// 本包实现了依赖注入容器模式，用于统一管理应用程序的所有依赖关系。
// 容器负责创建、配置和管理服务实例的生命周期，提供了以下功能：
//
// 1. 基础依赖管理：配置、数据库、日志系统
// 2. 服务层管理：用户、账户、代理等核心业务服务
// 3. 中间件管理：认证中间件、限流中间件等
// 4. 资源生命周期管理：优雅启动和关闭
//
// 使用依赖注入容器的优势：
// - 解耦：服务之间的依赖关系通过容器管理，而非硬编码
// - 可测试性：便于单元测试时模拟依赖
// - 配置集中化：所有依赖的创建逻辑集中在一处
// - 生命周期管理：统一管理资源的创建和清理
//
// Example:
//
//	// 创建并初始化容器
//	container, err := NewContainer()
//	if err != nil {
//	    log.Fatal("容器初始化失败:", err)
//	}
//	defer container.Close()
//
//	// 获取服务实例
//	userService := container.GetUserService()
package container

import (
	"fmt"
	"steam-backend/config"
	"steam-backend/logger"
	"steam-backend/middleware"
	"steam-backend/services"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Container 服务容器
//
// Container是应用程序的依赖注入容器，负责管理所有服务实例的生命周期。
// 它采用分层架构设计，包含以下几个层次：
//
// 基础设施层：
//   - Config: 配置管理器，处理应用配置
//   - DB: 数据库连接实例，提供数据持久化
//   - Logger: 结构化日志记录器，用于应用日志输出
//
// 业务服务层：
//   - UserService: 用户管理服务，处理用户注册、登录、信息管理等
//   - AccountService: 账户服务，处理账户信息、状态与积分
//   - CertificateService: 证书服务，处理证书管理
//   - ForgeService: 锻刀所服务，处理锻刀所管理
//   - SwordsmithService: 刀匠服务，处理刀匠管理
//   - ProxyService: 代理服务，管理可用的代理节点
//
// 中间件层：
//   - AuthMiddleware: 统一认证中间件
//
// 容器采用懒加载模式，只有在需要时才创建服务实例，
// 并确保每个服务只创建一次（单例模式）。
type Container struct {
	// 基础依赖
	Config *config.AppConfig // 配置管理器，使用统一配置
	DB     *gorm.DB          // GORM数据库连接，提供ORM功能
	Logger *zap.Logger       // Zap结构化日志记录器
	// 服务层
	UserService        services.UserServiceInterface // 用户管理服务接口
	AccountService     *services.AccountService      // 账户服务
	CertificateService *services.CertificateService  // 证书服务
	ForgeService       *services.ForgeService        // 锻刀所服务
	SwordsmithService  *services.SwordsmithService   // 刀匠服务
	ProxyService       *services.ProxyService        // 代理管理服务
	// 中间件
	AuthMiddleware *middleware.AuthMiddleware // 统一认证中间件
}

// NewContainer 创建新的服务容器
//
// NewContainer是容器的工厂函数，负责创建并初始化服务容器。
// 它按照以下顺序进行初始化：
//  1. 创建空的容器实例
//  2. 初始化基础依赖（配置、数据库、日志）
//  3. 初始化业务服务（用户、账户、代理等）
//  4. 初始化中间件（认证中间件）
//
// 如果任何步骤失败，函数将返回错误，确保容器处于完全可用状态。
//
// Returns:
//   - *Container: 完全初始化的服务容器实例
//   - error: 初始化过程中的任何错误
//
// Example:
//
//	container, err := NewContainer()
//	if err != nil {
//	    log.Fatal("容器初始化失败:", err)
//	}
//	defer container.Close()
func NewContainer() (*Container, error) {
	container := &Container{}

	// 初始化基础依赖
	if err := container.initBaseDependencies(); err != nil {
		return nil, err
	}

	// 初始化服务层
	if err := container.initServices(); err != nil {
		return nil, err
	}

	// 初始化中间件
	container.initMiddleware()

	return container, nil
}

// initBaseDependencies 初始化基础依赖
func (c *Container) initBaseDependencies() error {
	// 初始化配置
	c.Config = config.LoadConfig()

	// 初始化数据库
	if err := c.Config.InitDatabase(); err != nil {
		return err
	}
	c.DB = c.Config.GetDB()

	// 初始化日志系统
	logger, err := c.initLogger()
	if err != nil {
		return err
	}
	c.Logger = logger

	return nil
}

// initServices 初始化服务层
func (c *Container) initServices() error {
	logger.Info("开始初始化服务层...")

	// 初始化代理服务
	logger.Info("正在初始化代理服务...")
	c.ProxyService = services.NewProxyService()
	if c.ProxyService == nil {
		logger.Error("代理服务初始化失败")
		return fmt.Errorf("代理服务初始化失败")
	}
	logger.Info("代理服务初始化成功")

	// 初始化账户服务
	logger.Info("正在初始化账户服务...")
	c.AccountService = services.NewAccountService(c.ProxyService)
	if c.AccountService == nil {
		logger.Error("账户服务初始化失败")
		return fmt.Errorf("账户服务初始化失败")
	}
	logger.Info("账户服务初始化成功")

	// 初始化证书服务
	logger.Info("正在初始化证书服务...")
	c.CertificateService = services.NewCertificateService()
	if c.CertificateService == nil {
		logger.Error("证书服务初始化失败")
		return fmt.Errorf("证书服务初始化失败")
	}
	logger.Info("证书服务初始化成功")

	// 初始化锻刀所服务
	logger.Info("正在初始化锻刀所服务...")
	c.ForgeService = services.NewForgeService()
	if c.ForgeService == nil {
		logger.Error("锻刀所服务初始化失败")
		return fmt.Errorf("锻刀所服务初始化失败")
	}
	logger.Info("锻刀所服务初始化成功")

	// 初始化刀匠服务
	logger.Info("正在初始化刀匠服务...")
	c.SwordsmithService = services.NewSwordsmithService()
	if c.SwordsmithService == nil {
		logger.Error("刀匠服务初始化失败")
		return fmt.Errorf("刀匠服务初始化失败")
	}
	logger.Info("刀匠服务初始化成功")

	// 初始化用户服务
	logger.Info("正在初始化用户服务...")
	c.UserService = services.NewUserService()
	if c.UserService == nil {
		logger.Error("用户服务初始化失败")
		return fmt.Errorf("用户服务初始化失败")
	}
	logger.Info("用户服务初始化成功")

	logger.Info("服务层初始化完成")
	return nil
}

// initMiddleware 初始化中间件
func (c *Container) initMiddleware() {
	// 统一认证中间件
	c.AuthMiddleware = middleware.NewAuthMiddleware(c.Logger)
}

// initLogger 初始化日志系统
func (c *Container) initLogger() (*zap.Logger, error) {
	// 根据环境配置日志
	var logger *zap.Logger
	var err error

	if c.Config.GetString("app.environment") == "production" {
		// 生产环境配置
		loggerConfig := zap.NewProductionConfig()
		loggerConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		logger, err = loggerConfig.Build()
	} else {
		// 开发环境配置
		loggerConfig := zap.NewDevelopmentConfig()
		loggerConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		logger, err = loggerConfig.Build()
	}

	if err != nil {
		return nil, err
	}

	return logger, nil
}

// Close 关闭容器，清理资源
func (c *Container) Close() error {
	// 同步日志
	if c.Logger != nil {
		c.Logger.Sync()
	}

	// 关闭数据库连接
	if c.Config != nil {
		return c.Config.CloseDatabase()
	}

	return nil
}

// GetAccountService 获取账户服务
func (c *Container) GetAccountService() *services.AccountService {
	return c.AccountService
}

// GetCertificateService 获取证书服务
func (c *Container) GetCertificateService() *services.CertificateService {
	return c.CertificateService
}

// GetForgeService 获取锻刀所服务
func (c *Container) GetForgeService() *services.ForgeService {
	return c.ForgeService
}

// GetSwordsmithService 获取刀匠服务
func (c *Container) GetSwordsmithService() *services.SwordsmithService {
	return c.SwordsmithService
}

// GetUserService 获取用户服务
func (c *Container) GetUserService() services.UserServiceInterface {
	return c.UserService
}

// GetProxyService 获取代理服务
func (c *Container) GetProxyService() *services.ProxyService {
	return c.ProxyService
}

// GetAuthMiddleware 获取统一认证中间件
func (c *Container) GetAuthMiddleware() *middleware.AuthMiddleware {
	return c.AuthMiddleware
}

// GetDB 获取数据库实例
func (c *Container) GetDB() *gorm.DB {
	return c.DB
}

// GetLogger 获取日志实例
func (c *Container) GetLogger() *zap.Logger {
	return c.Logger
}

// GetConfig 获取配置实例
func (c *Container) GetConfig() *config.AppConfig {
	return c.Config
}
