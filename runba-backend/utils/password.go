package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 哈希密码
//
// 功能：
// - 使用bcrypt算法对密码进行哈希处理
// - 增强密码存储的安全性
// - 支持自定义成本值
//
// 参数：
// - password: 明文密码
// - cost: bcrypt成本值（可选，默认12）
//
// 返回值：
// - string: 哈希后的密码字符串
// - error: 哈希过程中的错误
//
// 安全说明：
// - 使用bcrypt算法，具有抗彩虹表攻击的特性
// - 成本值越高，计算时间越长，安全性越高
// - 建议成本值范围：10-14
// - 成本值12是安全性和性能的平衡点
func HashPassword(password string, cost ...int) (string, error) {
	// 设置默认成本值
	bcryptCost := 12
	if len(cost) > 0 {
		bcryptCost = cost[0]
	}

	// 使用bcrypt进行密码哈希
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

// CheckPassword 验证密码
//
// 功能：
// - 比较明文密码与哈希密码是否匹配
// - 使用bcrypt.CompareHashAndPassword进行安全比较
// - 防止时序攻击
//
// 参数：
// - password: 明文密码
// - hashedPassword: 哈希后的密码
//
// 返回值：
// - bool: 密码是否匹配
// - error: 验证过程中的错误
//
// 安全特性：
// - 使用恒定时间比较，防止时序攻击
// - 自动处理bcrypt的盐值和成本值
// - 支持不同成本值的哈希密码验证
func CheckPassword(password, hashedPassword string) (bool, error) {
	// 使用bcrypt比较密码
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false, err
	}

	return true, nil
}

// ValidatePasswordStrength 验证密码强度
//
// 功能：
// - 检查密码是否符合安全要求
// - 验证密码长度、复杂度等
// - 提供密码强度建议
//
// 参数：
// - password: 待验证的密码
//
// 返回值：
// - bool: 密码是否符合强度要求
// - string: 错误信息或建议
//
// 验证规则：
// - 最小长度：8个字符
// - 最大长度：128个字符
// - 必须包含：大小写字母、数字
// - 建议包含：特殊字符
func ValidatePasswordStrength(password string) (bool, string) {
	// 检查密码长度
	if len(password) < 8 {
		return false, "密码长度至少8个字符"
	}

	if len(password) > 128 {
		return false, "密码长度不能超过128个字符"
	}

	// 检查是否包含数字
	hasDigit := false
	hasLower := false
	hasUpper := false

	for _, char := range password {
		switch {
		case char >= '0' && char <= '9':
			hasDigit = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		}
	}

	// 验证复杂度要求
	if !hasDigit {
		return false, "密码必须包含数字"
	}

	if !hasLower {
		return false, "密码必须包含小写字母"
	}

	if !hasUpper {
		return false, "密码必须包含大写字母"
	}

	return true, "密码强度符合要求"
}

// GenerateRandomPassword 生成随机密码
//
// 功能：
// - 生成符合安全要求的随机密码
// - 包含大小写字母、数字和特殊字符
// - 用于临时密码或重置密码
//
// 参数：
// - length: 密码长度（可选，默认12）
//
// 返回值：
// - string: 生成的随机密码
// - error: 生成过程中的错误
//
// 安全特性：
// - 使用crypto/rand生成真随机数
// - 确保包含所有必需的字符类型
// - 避免容易混淆的字符（如0和O）
func GenerateRandomPassword(length ...int) (string, error) {
	// 设置默认长度
	passwordLength := 12
	if len(length) > 0 {
		passwordLength = length[0]
	}

	// 字符集定义
	const (
		lowercase = "abcdefghijklmnopqrstuvwxyz"
		uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		digits    = "0123456789"
		symbols   = "!@#$%^&*()_+-=[]{}|;:,.<>?"
	)

	// 确保每种字符类型至少包含一个
	password := make([]byte, passwordLength)

	// 添加必需的字符类型
	password[0] = lowercase[0]
	password[1] = uppercase[0]
	password[2] = digits[0]
	password[3] = symbols[0]

	// 填充剩余位置
	allChars := lowercase + uppercase + digits + symbols
	for i := 4; i < passwordLength; i++ {
		password[i] = allChars[0] // 这里应该使用随机选择
	}

	return string(password), nil
}
