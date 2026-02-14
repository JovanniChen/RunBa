package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTSecret 全局JWT签名密钥
// 在应用启动时从配置中加载
var JWTSecret []byte

// Claims JWT声明结构体
//
// 定义JWT令牌中包含的用户信息
// 包含标准JWT声明和自定义用户信息
type Claims struct {
	UserID               uint   `json:"user_id"`  // 用户ID
	Username             string `json:"username"` // 用户名
	jwt.RegisteredClaims        // 标准JWT声明（过期时间、签发时间等）
}

// GenerateToken 生成JWT令牌
//
// 功能：
// - 根据用户信息生成JWT令牌
// - 设置令牌的过期时间
// - 使用全局密钥进行签名
//
// 参数：
// - userID: 用户ID
// - username: 用户名
// - expirationHours: 过期时间（小时）
//
// 返回值：
// - string: 生成的JWT令牌字符串
// - error: 生成过程中的错误
//
// 安全考虑：
// - 令牌包含用户ID和用户名信息
// - 默认过期时间为24小时
// - 使用强密钥进行签名
func GenerateToken(userID uint, username string, expirationHours int) (string, error) {
	// 计算过期时间
	expirationTime := time.Now().Add(time.Duration(expirationHours) * time.Hour)

	// 创建JWT声明
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),     // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),     // 生效时间
		},
	}

	// 创建JWT令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥签名
	tokenString, err := token.SignedString(JWTSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken 解析JWT令牌
//
// 功能：
// - 验证JWT令牌的有效性
// - 解析令牌中的用户信息
// - 检查令牌是否过期
//
// 参数：
// - tokenString: JWT令牌字符串
//
// 返回值：
// - *Claims: 解析出的JWT声明
// - error: 解析过程中的错误
//
// 错误类型：
// - TokenExpiredError: 令牌已过期
// - TokenNotValidYet: 令牌尚未生效
// - TokenMalformed: 令牌格式错误
// - TokenInvalidSignature: 令牌签名无效
func ParseToken(tokenString string) (*Claims, error) {
	// 解析令牌
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("无效的签名方法")
		}
		return JWTSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// 验证令牌并提取声明
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的令牌")
}

// RefreshToken 刷新JWT令牌
//
// 功能：
// - 验证当前令牌的有效性
// - 生成新的令牌（延长过期时间）
// - 保持用户信息不变
//
// 参数：
// - tokenString: 当前JWT令牌字符串
// - newExpirationHours: 新令牌的过期时间（小时）
//
// 返回值：
// - string: 新的JWT令牌字符串
// - error: 刷新过程中的错误
//
// 使用场景：
// - 用户活跃时自动刷新令牌
// - 延长用户会话时间
func RefreshToken(tokenString string, newExpirationHours int) (string, error) {
	// 解析当前令牌
	claims, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}

	// 生成新令牌
	return GenerateToken(claims.UserID, claims.Username, newExpirationHours)
}

// ValidateToken 验证JWT令牌（不返回声明）
//
// 功能：
// - 仅验证令牌的有效性
// - 不返回令牌中的用户信息
// - 用于简单的令牌验证场景
//
// 参数：
// - tokenString: JWT令牌字符串
//
// 返回值：
// - bool: 令牌是否有效
// - error: 验证过程中的错误
//
// 使用场景：
// - 简单的令牌有效性检查
// - 不需要用户信息的验证场景
func ValidateToken(tokenString string) (bool, error) {
	_, err := ParseToken(tokenString)
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetTokenExpiration 获取令牌过期时间
//
// 功能：
// - 解析令牌并返回过期时间
// - 不验证令牌的有效性
//
// 参数：
// - tokenString: JWT令牌字符串
//
// 返回值：
// - time.Time: 令牌过期时间
// - error: 解析过程中的错误
//
// 使用场景：
// - 检查令牌剩余有效期
// - 提前提醒用户令牌即将过期
func GetTokenExpiration(tokenString string) (time.Time, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return time.Time{}, err
	}

	return claims.ExpiresAt.Time, nil
}
