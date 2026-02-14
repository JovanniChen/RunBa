package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"steam-backend/config"
	"steam-backend/container"
	"steam-backend/logger"
	"steam-backend/routes"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

// main 应用程序主入口函数
func main() {
	// 添加全局 panic 恢复机制
	defer func() {
		if r := recover(); r != nil {
			logger.Fatal("应用程序发生严重错误，即将退出: %v", r)
		}
	}()

	// 加载配置以获取日志级别设置
	cfg := config.LoadConfig()
	// 使用配置文件中的日志级别
	logConfig := logger.DefaultConfig
	logConfig.Level = logger.ParseLogLevel(cfg.GetString("log.level"))
	logger.InitGlobalLogger("steam-backend", logConfig)

	// 打印应用启动信息
	logger.Info("=================================")
	logger.Info("润爸刀剑后端API v1.0.0正在启动...")
	logger.Info("=================================")

	// 初始化服务容器
	logger.Info("正在初始化服务容器...")
	appContainer, err := container.NewContainer()
	if err != nil {
		logger.Fatal("服务容器初始化失败: %v", err)
	}
	defer func() {
		if err := appContainer.Close(); err != nil {
			logger.Error("关闭服务容器失败: %v", err)
		}
	}()

	// 获取配置
	appConfig := appContainer.GetConfig()
	logger.Info("应用名称: %s", appConfig.GetString("app.name"))
	logger.Info("应用版本: %s", appConfig.GetString("app.version"))

	// 设置Gin运行模式
	mode := appConfig.GetString("server.mode")
	gin.SetMode(mode)
	logger.Info("Gin运行模式: %s", mode)

	// 设置路由，传入容器实例
	logger.Info("正在设置路由...")
	router := routes.SetupRoutes(appContainer)

	// 获取服务器地址
	serverAddr := config.GetServerAddress()

	// 启动HTTP服务器
	logger.Info("服务器启动成功! 监听地址: %s", serverAddr)
	// printAPIEndpoints()

	// 优雅关闭设置
	setupGracefulShutdown(router, serverAddr)
}

// setupGracefulShutdown 设置优雅关闭
func setupGracefulShutdown(router *gin.Engine, serverAddr string) {
	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	// 在goroutine中启动服务器
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("服务器启动失败: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("正在关闭服务器...")

	// 给服务器5秒时间完成当前请求
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("服务器强制关闭: %v", err)
	}

	logger.Info("服务器已退出")
}
