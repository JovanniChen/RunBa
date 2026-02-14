package config

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"steam-backend/logger"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// AppConfig 配置管理器
type AppConfig struct {
	viper     *viper.Viper
	db        *gorm.DB
	validator *ConfigValidator
}

var globalConfig *AppConfig

// Config 应用配置结构
type Config struct {
	App      ApplicationConfig `json:"app"`
	Server   ServerConfig      `json:"server"`
	Database DatabaseConfig    `json:"database"`
	JWT      JWTConfig         `json:"jwt"`
	Redis    RedisConfig       `json:"redis"`
	Log      LogConfig         `json:"log"`
	Jobs     JobsConfig        `json:"jobs"`
}

// ApplicationConfig 应用基本配置
type ApplicationConfig struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Environment string `json:"environment"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host         string `json:"host"`
	Port         string `json:"port"`
	Mode         string `json:"mode"`
	ReadTimeout  int    `json:"read_timeout"`
	WriteTimeout int    `json:"write_timeout"`
	IdleTimeout  int    `json:"idle_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host            string `json:"host"`
	Port            string `json:"port"`
	User            string `json:"user"`
	Password        string `json:"password"`
	DBName          string `json:"db_name"`
	Charset         string `json:"charset"`
	MaxIdleConns    int    `json:"max_idle_conns"`
	MaxOpenConns    int    `json:"max_open_conns"`
	ConnMaxLifetime int    `json:"conn_max_lifetime"`
	LogLevel        string `json:"log_level"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret            string `json:"secret"`
	Expiration        int    `json:"expiration"`         // 小时
	RefreshExpiration int    `json:"refresh_expiration"` // 小时
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host         string `json:"host"`
	Port         string `json:"port"`
	Password     string `json:"password"`
	DB           int    `json:"db"`
	PoolSize     int    `json:"pool_size"`
	DialTimeout  int    `json:"dial_timeout"`  // 连接超时(秒)
	ReadTimeout  int    `json:"read_timeout"`  // 读取超时(秒)
	WriteTimeout int    `json:"write_timeout"` // 写入超时(秒)
}

// LogConfig 日志配置
type LogConfig struct {
	Level  string `json:"level"`
	Format string `json:"format"`
	Output string `json:"output"`
}

// JobsConfig 任务调度器配置
type JobsConfig struct {
	Enabled bool `json:"enabled"` // 是否启用任务调度器
}

// ConfigValidator 配置验证器
type ConfigValidator struct{}

// NewAppConfig 创建新的配置管理器
func NewAppConfig() *AppConfig {
	v := viper.New()

	// 配置文件设置 - 只使用config.yaml
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	// 设置默认值
	setDefaults(v)

	return &AppConfig{
		viper:     v,
		validator: &ConfigValidator{},
	}
}

// setDefaults 设置默认配置值
func setDefaults(v *viper.Viper) {
	// 应用配置
	v.SetDefault("app.name", "RunBa Backend API")
	v.SetDefault("app.version", "1.0.0")
	v.SetDefault("app.environment", "development")

	// 服务器配置
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", "8080")
	v.SetDefault("server.mode", "debug")
	v.SetDefault("server.read_timeout", 30)
	v.SetDefault("server.write_timeout", 30)
	v.SetDefault("server.idle_timeout", 60)

	// 数据库配置
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", "3306")
	v.SetDefault("database.user", "root")
	v.SetDefault("database.password", "")
	v.SetDefault("database.name", "runba_backend")
	v.SetDefault("database.charset", "utf8mb4")
	v.SetDefault("database.max_idle_conns", 10)
	v.SetDefault("database.max_open_conns", 100)
	v.SetDefault("database.conn_max_lifetime", 3600)
	v.SetDefault("database.log_level", "info")

	// JWT配置
	v.SetDefault("jwt.secret", "your-super-secret-jwt-key-change-this-in-production")
	v.SetDefault("jwt.expiration", 24)
	v.SetDefault("jwt.refresh_expiration", 168)

	// Redis配置
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", "6379")
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)
	v.SetDefault("redis.pool_size", 10)
	v.SetDefault("redis.dial_timeout", 5)
	v.SetDefault("redis.read_timeout", 3)
	v.SetDefault("redis.write_timeout", 3)

	// 日志配置
	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "json")
	v.SetDefault("log.output", "stdout")

	// 任务调度器配置
	v.SetDefault("jobs.enabled", true)
}

// LoadConfig 加载配置
func (ac *AppConfig) LoadConfig() error {
	if err := ac.viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("读取配置文件失败: %w", err)
		}
		logger.Info("警告: 未找到config.yaml文件，使用默认配置")
	} else {
		logger.Info("成功加载配置文件: config.yaml")
	}
	return nil
}

// GetConfig 获取完整配置
func (ac *AppConfig) GetConfig() *Config {
	return &Config{
		App:      ac.GetAppConfig(),
		Server:   ac.GetServerConfig(),
		Database: ac.GetDatabaseConfig(),
		JWT:      ac.GetJWTConfig(),
		Redis:    ac.GetRedisConfig(),
		Log:      ac.GetLogConfig(),
		Jobs:     ac.GetJobsConfig(),
	}
}

// GetAppConfig 获取应用配置
func (ac *AppConfig) GetAppConfig() ApplicationConfig {
	return ApplicationConfig{
		Name:        ac.GetString("app.name"),
		Version:     ac.GetString("app.version"),
		Environment: ac.GetString("app.environment"),
	}
}

// GetServerConfig 获取服务器配置
func (ac *AppConfig) GetServerConfig() ServerConfig {
	return ServerConfig{
		Host:         ac.GetString("server.host"),
		Port:         ac.GetString("server.port"),
		Mode:         ac.GetString("server.mode"),
		ReadTimeout:  ac.GetInt("server.read_timeout"),
		WriteTimeout: ac.GetInt("server.write_timeout"),
		IdleTimeout:  ac.GetInt("server.idle_timeout"),
	}
}

// GetDatabaseConfig 获取数据库配置
func (ac *AppConfig) GetDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Host:            ac.GetString("database.host"),
		Port:            ac.GetString("database.port"),
		User:            ac.GetString("database.user"),
		Password:        ac.GetString("database.password"),
		DBName:          ac.GetString("database.name"),
		Charset:         ac.GetString("database.charset"),
		MaxIdleConns:    ac.GetInt("database.max_idle_conns"),
		MaxOpenConns:    ac.GetInt("database.max_open_conns"),
		ConnMaxLifetime: ac.GetInt("database.conn_max_lifetime"),
		LogLevel:        ac.GetString("database.log_level"),
	}
}

// GetJWTConfig 获取JWT配置
func (ac *AppConfig) GetJWTConfig() JWTConfig {
	return JWTConfig{
		Secret:            ac.GetString("jwt.secret"),
		Expiration:        ac.GetInt("jwt.expiration"),
		RefreshExpiration: ac.GetInt("jwt.refresh_expiration"),
	}
}

// GetRedisConfig 获取Redis配置
func (ac *AppConfig) GetRedisConfig() RedisConfig {
	return RedisConfig{
		Host:         ac.GetString("redis.host"),
		Port:         ac.GetString("redis.port"),
		Password:     ac.GetString("redis.password"),
		DB:           ac.GetInt("redis.db"),
		PoolSize:     ac.GetInt("redis.pool_size"),
		DialTimeout:  ac.GetInt("redis.dial_timeout"),
		ReadTimeout:  ac.GetInt("redis.read_timeout"),
		WriteTimeout: ac.GetInt("redis.write_timeout"),
	}
}

// GetLogConfig 获取日志配置
func (ac *AppConfig) GetLogConfig() LogConfig {
	return LogConfig{
		Level:  ac.GetString("log.level"),
		Format: ac.GetString("log.format"),
		Output: ac.GetString("log.output"),
	}
}

// GetJobsConfig 获取任务调度器配置
func (ac *AppConfig) GetJobsConfig() JobsConfig {
	return JobsConfig{
		Enabled: ac.GetBool("jobs.enabled"),
	}
}

// InitDatabase 初始化数据库连接
func (ac *AppConfig) InitDatabase() error {
	dbConfig := ac.GetDatabaseConfig()

	// 构建DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DBName,
		dbConfig.Charset,
	)

	// 配置GORM
	var logLevel gormlogger.LogLevel
	switch dbConfig.LogLevel {
	case "silent":
		logLevel = gormlogger.Silent
	case "error":
		logLevel = gormlogger.Error
	case "warn":
		logLevel = gormlogger.Warn
	case "info":
		logLevel = gormlogger.Info
	default:
		logLevel = gormlogger.Info
	}

	gormConfig := &gorm.Config{
		Logger:                                   gormlogger.Default.LogMode(logLevel),
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}

	// 连接数据库
	var err error
	ac.db, err = gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return fmt.Errorf("数据库连接失败: %w", err)
	}

	// 配置连接池
	sqlDB, err := ac.db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %w", err)
	}

	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(dbConfig.ConnMaxLifetime) * time.Second)

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("数据库连接测试失败: %w", err)
	}

	logger.Info("数据库连接成功! 连接池配置: MaxIdle=%d, MaxOpen=%d, MaxLifetime=%ds",
		dbConfig.MaxIdleConns,
		dbConfig.MaxOpenConns,
		dbConfig.ConnMaxLifetime,
	)

	// 启动监控
	go ac.monitorConnectionPool(sqlDB)

	return nil
}

// monitorConnectionPool 监控连接池
func (ac *AppConfig) monitorConnectionPool(sqlDB *sql.DB) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		stats := sqlDB.Stats()
		logger.Info("数据库连接池状态: Open=%d, InUse=%d, Idle=%d, WaitCount=%d, WaitDuration=%v",
			stats.OpenConnections,
			stats.InUse,
			stats.Idle,
			stats.WaitCount,
			stats.WaitDuration,
		)

		if stats.WaitCount > 100 {
			logger.Info("警告: 数据库连接池等待次数过多: %d", stats.WaitCount)
		}
	}
}

// CloseDatabase 关闭数据库连接
func (ac *AppConfig) CloseDatabase() error {
	if ac.db != nil {
		sqlDB, err := ac.db.DB()
		if err == nil {
			err = sqlDB.Close()
			if err != nil {
				return err
			}
			logger.Info("数据库连接已关闭")
		}
	}
	return nil
}

// ValidateConfig 验证配置
func (ac *AppConfig) ValidateConfig() error {
	config := ac.GetConfig()
	return ac.validator.ValidateAll(config)
}

// 基础访问方法
func (ac *AppConfig) GetString(key string) string {
	return ac.viper.GetString(key)
}

func (ac *AppConfig) GetInt(key string) int {
	return ac.viper.GetInt(key)
}

func (ac *AppConfig) GetBool(key string) bool {
	return ac.viper.GetBool(key)
}

func (ac *AppConfig) GetDuration(key string) time.Duration {
	return ac.viper.GetDuration(key)
}

func (ac *AppConfig) IsSet(key string) bool {
	return ac.viper.IsSet(key)
}

func (ac *AppConfig) GetDB() *gorm.DB {
	return ac.db
}

// 配置验证器方法
func (cv *ConfigValidator) ValidateAll(config *Config) error {
	if err := cv.ValidateDatabaseConfig(config.Database); err != nil {
		return fmt.Errorf("数据库配置错误: %w", err)
	}

	if err := cv.ValidateJWTConfig(config.JWT); err != nil {
		return fmt.Errorf("JWT配置错误: %w", err)
	}

	if err := cv.ValidateServerConfig(config.Server); err != nil {
		return fmt.Errorf("服务器配置错误: %w", err)
	}

	return nil
}

func (cv *ConfigValidator) ValidateDatabaseConfig(config DatabaseConfig) error {
	if config.Host == "" {
		return fmt.Errorf("数据库主机地址不能为空")
	}

	if config.Port == "" {
		return fmt.Errorf("数据库端口不能为空")
	}

	if _, err := strconv.Atoi(config.Port); err != nil {
		return fmt.Errorf("数据库端口格式错误: %s", config.Port)
	}

	if config.User == "" {
		return fmt.Errorf("数据库用户名不能为空")
	}

	if config.DBName == "" {
		return fmt.Errorf("数据库名称不能为空")
	}

	return nil
}

func (cv *ConfigValidator) ValidateJWTConfig(config JWTConfig) error {
	if config.Secret == "" {
		return fmt.Errorf("JWT密钥不能为空")
	}

	if len(config.Secret) < 32 {
		return fmt.Errorf("JWT密钥长度不足，建议至少32个字符")
	}

	if config.Secret == "your-super-secret-jwt-key-change-this-in-production" {
		return fmt.Errorf("请修改默认的JWT密钥")
	}

	return nil
}

func (cv *ConfigValidator) ValidateServerConfig(config ServerConfig) error {
	if config.Port == "" {
		return fmt.Errorf("服务器端口不能为空")
	}

	port, err := strconv.Atoi(config.Port)
	if err != nil {
		return fmt.Errorf("服务器端口格式错误: %s", config.Port)
	}

	if port < 1 || port > 65535 {
		return fmt.Errorf("服务器端口超出范围: %d", port)
	}

	validModes := map[string]bool{
		"debug":   true,
		"release": true,
		"test":    true,
	}
	if !validModes[config.Mode] {
		return fmt.Errorf("无效的Gin模式: %s", config.Mode)
	}

	return nil
}

// 全局配置管理
func LoadConfig() *AppConfig {
	if globalConfig == nil {
		globalConfig = NewAppConfig()
		if err := globalConfig.LoadConfig(); err != nil {
			logger.Info("警告: 无法加载配置文件: %v", err)
		}
	}
	return globalConfig
}

func GetDB() *gorm.DB {
	config := LoadConfig()
	return config.GetDB()
}

func GetJWTSecret() string {
	config := LoadConfig()
	return config.GetString("jwt.secret")
}

func GetServerAddress() string {
	config := LoadConfig()
	port := config.GetString("server.port")
	return ":" + port
}

// GetJobsConfig 获取任务调度器配置
func GetJobsConfig() JobsConfig {
	config := LoadConfig()
	return config.GetJobsConfig()
}

// NewRedisClient 创建Redis客户端
func (ac *AppConfig) NewRedisClient() (*redis.Client, error) {
	redisConfig := ac.GetRedisConfig()

	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", redisConfig.Host, redisConfig.Port),
		Password:     redisConfig.Password,
		DB:           redisConfig.DB,
		PoolSize:     redisConfig.PoolSize,
		DialTimeout:  time.Duration(redisConfig.DialTimeout) * time.Second,
		ReadTimeout:  time.Duration(redisConfig.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(redisConfig.WriteTimeout) * time.Second,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("Redis连接失败: %v", err)
	}

	logger.Info("Redis连接成功: %s:%s", redisConfig.Host, redisConfig.Port)
	return client, nil
}
