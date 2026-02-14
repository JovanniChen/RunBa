package middleware

import (
	"fmt"
	"net/http"
	"steam-backend/models"
	"steam-backend/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AuthConfig 认证配置
type AuthConfig struct {
	RequireJWT bool // 是否要求JWT认证
	Optional   bool // 是否为可选认证
}

// AuthMiddleware 统一认证中间件
type AuthMiddleware struct {
	logger *zap.Logger
}

// NewAuthMiddleware 创建认证中间件
func NewAuthMiddleware(logger *zap.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		logger: logger,
	}
}

// Authenticate 统一认证方法
func (m *AuthMiddleware) Authenticate(config AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// JWT认证处理
		if config.RequireJWT || !config.Optional {
			if !m.handleJWTAuth(c, config.Optional) {
				return
			}
		}

		c.Next()
	}
}

// handleJWTAuth 处理JWT认证
func (m *AuthMiddleware) handleJWTAuth(c *gin.Context, optional bool) bool {
	authHeader := c.GetHeader("token")

	if authHeader == "" {
		if !optional {
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse(
				401,
				"缺少token头",
				nil,
			))
			c.Abort()
			return false
		}
		return true
	}

	if !strings.HasPrefix(authHeader, "") {
		if !optional {
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse(
				401,
				"token头格式错误，应为'<token>'",
				nil,
			))
			c.Abort()
			return false
		}
		return true
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := utils.ParseToken(tokenString)

	if err != nil {
		if !optional {
			c.JSON(http.StatusUnauthorized, models.NewErrorResponse(
				401,
				"无效的token",
				err,
			))
			c.Abort()
			return false
		}
		return true
	}

	// JWT认证成功，存储用户信息
	c.Set("user_id", claims.UserID)
	c.Set("username", claims.Username)
	return true
}

// RequireJWT 要求JWT认证 (向后兼容)
func (m *AuthMiddleware) RequireJWT() gin.HandlerFunc {
	return m.Authenticate(AuthConfig{
		RequireJWT: true,
		Optional:   false,
	})
}

// OptionalAuth 可选认证 (向后兼容)
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return m.Authenticate(AuthConfig{
		RequireJWT: false,
		Optional:   true,
	})
}

// 向后兼容的函数
func RequireJWT() gin.HandlerFunc {
	middleware := &AuthMiddleware{}
	return middleware.RequireJWT()
}

func OptionalAuth() gin.HandlerFunc {
	middleware := &AuthMiddleware{}
	return middleware.OptionalAuth()
}

// CORSMiddleware 跨域资源共享中间件
//
// 功能：
// - 处理跨域请求，允许前端应用访问API
// - 设置CORS相关的HTTP头
// - 处理预检请求（OPTIONS）
//
// 安全考虑：
// - 生产环境应该配置具体的允许域名
// - 当前设置为允许所有源（*），仅用于开发环境
//
// 配置的CORS头：
// - Access-Control-Allow-Origin: 允许的源
// - Access-Control-Allow-Methods: 允许的HTTP方法
// - Access-Control-Allow-Headers: 允许的请求头
// - Access-Control-Allow-Credentials: 允许携带凭证
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置允许的源，生产环境中应该配置具体的域名
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		// 设置允许的HTTP方法
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// 设置允许的请求头
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, token")

		// 设置允许携带凭证
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// LoggerMiddleware 自定义日志中间件
//
// 功能：
// - 记录请求的详细信息，便于调试和监控
// - 自定义日志格式，包含时间戳、方法、路径、状态码等
// - 用于性能分析和问题排查
//
// 日志格式：
// [时间戳] 方法 路径 状态码 响应时间 客户端IP
//
// 示例输出：
// [2024-01-15 10:30:45] GET /api/v1/users 200 15.2ms 192.168.1.100
func LoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// 自定义日志格式
		return fmt.Sprintf("[%s] %s %s %d %s %s\n",
			param.TimeStamp.Format("2006-01-02 15:04:05"),
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
			param.ClientIP,
		)
	})
}

// GetUserID 从Gin上下文中获取当前用户ID
//
// 功能：
// - 辅助函数，用于在控制器中快速获取当前登录用户的ID
// - 从中间件设置的上下文中提取用户ID
//
// 参数：
// - c: Gin上下文对象
//
// 返回值：
// - uint: 用户ID
// - bool: 是否成功获取到用户ID
//
// 使用示例：
// userID, exists := middleware.GetUserID(c)
//
//	if !exists {
//	    // 处理未找到用户ID的情况
//	}
func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	if id, ok := userID.(uint); ok {
		return id, true
	}

	return 0, false
}

// GetUsername 从Gin上下文中获取当前用户名
//
// 功能：
// - 辅助函数，用于在控制器中快速获取当前登录用户的用户名
// - 从中间件设置的上下文中提取用户名
//
// 参数：
// - c: Gin上下文对象
//
// 返回值：
// - string: 用户名
// - bool: 是否成功获取到用户名
//
// 使用示例：
// username, exists := middleware.GetUsername(c)
//
//	if !exists {
//	    // 处理未找到用户名的情况
//	}
func GetUsername(c *gin.Context) (string, bool) {
	username, exists := c.Get("username")
	if !exists {
		return "", false
	}

	if name, ok := username.(string); ok {
		return name, true
	}

	return "", false
}
