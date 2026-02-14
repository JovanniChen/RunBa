package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter 限流器结构体
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rps      float64 // 每秒请求数
	burst    int     // 突发请求数
}

// NewRateLimiter 创建限流器
func NewRateLimiter(rps float64, burst int) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rps:      rps,
		burst:    burst,
	}
}

// getLimiter 获取指定IP的限流器
func (rl *RateLimiter) getLimiter(key string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[key]
	if !exists {
		limiter = rate.NewLimiter(rate.Limit(rl.rps), rl.burst)
		rl.limiters[key] = limiter
	}

	return limiter
}

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware(rps float64, burst int) gin.HandlerFunc {
	limiter := NewRateLimiter(rps, burst)

	return func(c *gin.Context) {
		// 获取客户端IP
		clientIP := c.ClientIP()
		if clientIP == "" {
			clientIP = "unknown"
		}

		// 获取限流器
		limiter := limiter.getLimiter(clientIP)

		// 检查是否允许请求
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":        429,
				"message":     "请求过于频繁，请稍后再试",
				"retry_after": time.Now().Add(time.Second).Unix(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// IPWhitelistMiddleware IP白名单中间件
func IPWhitelistMiddleware(whitelist []string) gin.HandlerFunc {
	whitelistMap := make(map[string]bool)
	for _, ip := range whitelist {
		whitelistMap[ip] = true
	}

	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		// 检查是否在白名单中
		if !whitelistMap[clientIP] {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "IP地址不在白名单中",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// UserRateLimitMiddleware 基于用户的限流中间件
func UserRateLimitMiddleware(rps float64, burst int) gin.HandlerFunc {
	limiter := NewRateLimiter(rps, burst)

	return func(c *gin.Context) {
		// 从JWT中获取用户ID
		userID, exists := GetUserID(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未授权访问",
			})
			c.Abort()
			return
		}

		// 使用用户ID作为限流键
		key := fmt.Sprintf("user:%d", userID)
		limiter := limiter.getLimiter(key)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// CleanupExpiredLimiters 清理过期的限流器
func (rl *RateLimiter) CleanupExpiredLimiters() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
}
