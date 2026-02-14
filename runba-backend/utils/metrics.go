package utils

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// HTTP请求计数器
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "HTTP请求总数",
		},
		[]string{"method", "endpoint", "status"},
	)

	// HTTP请求持续时间
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP请求持续时间",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// 活跃用户数
	activeUsers = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_users",
			Help: "当前活跃用户数",
		},
	)

	// 数据库连接数
	dbConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_connections",
			Help: "当前数据库连接数",
		},
	)
)

// InitMetrics 初始化监控指标
func InitMetrics() {
	// 注册指标
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
	prometheus.MustRegister(activeUsers)
	prometheus.MustRegister(dbConnections)
}

// MetricsMiddleware 监控中间件
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 处理请求
		c.Next()

		// 记录请求指标
		duration := time.Since(start).Seconds()
		status := c.Writer.Status()
		method := c.Request.Method
		endpoint := c.FullPath()

		// 增加请求计数
		httpRequestsTotal.WithLabelValues(method, endpoint, string(rune(status))).Inc()

		// 记录请求持续时间
		httpRequestDuration.WithLabelValues(method, endpoint).Observe(duration)
	}
}

// SetupMetrics 设置监控路由
func SetupMetrics(router *gin.Engine) {
	// 添加监控中间件
	router.Use(MetricsMiddleware())

	// 添加Prometheus指标端点
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}

// UpdateActiveUsers 更新活跃用户数
func UpdateActiveUsers(count int) {
	activeUsers.Set(float64(count))
}

// UpdateDBConnections 更新数据库连接数
func UpdateDBConnections(count int) {
	dbConnections.Set(float64(count))
}

// GetMetrics 获取当前指标
func GetMetrics() map[string]interface{} {
	return map[string]interface{}{
		"http_requests_total":   httpRequestsTotal,
		"http_request_duration": httpRequestDuration,
		"active_users":          activeUsers,
		"db_connections":        dbConnections,
	}
}
