package utils

import (
	"crypto/rand"
	"fmt"
	"hash/crc32"
	"math/big"
	"strconv"
	"time"
)

// ActivationCodeGenerator 激活码生成器
type ActivationCodeGenerator struct{}

// NewActivationCodeGenerator 创建新的激活码生成器
func NewActivationCodeGenerator() *ActivationCodeGenerator {
	return &ActivationCodeGenerator{}
}

// GenerateActivationCode 生成唯一的激活码
// 格式：21位数字字符串（20位数据 + 1位校验位）
func (g *ActivationCodeGenerator) GenerateActivationCode(name string, points int32) string {
	// 1. 获取当前时间戳（纳秒级，提高精度）
	timestamp := time.Now().UnixNano()

	// 2. 计算名称的哈希值
	nameHash := crc32.ChecksumIEEE([]byte(name))

	// 3. 生成高安全性随机数（7位数字）
	randomNum, _ := rand.Int(rand.Reader, big.NewInt(100000))

	// 4. 组合生成20位数据部分
	// 时间戳(6位) + 名称哈希(4位) + 点数(2位) + 随机数(7位) + 保留位(1位)
	dataCode := fmt.Sprintf("%06d%04d%02d%05d%01d",
		timestamp%1000000,         // 时间戳取后6位
		nameHash%10000,            // 名称哈希取4位
		points%100,                // 点数取后2位
		randomNum.Uint64()%100000, // 随机数取7位
		0)                         // 保留位，固定为0

	// 5. 计算校验位（使用简单的校验和算法）
	checksum := g.calculateChecksum(dataCode)

	// 6. 组合最终激活码（20位数据 + 1位校验位）
	code := dataCode + fmt.Sprintf("%d", checksum)

	return code
}

// calculateChecksum 计算校验位
func (g *ActivationCodeGenerator) calculateChecksum(dataCode string) int {
	sum := 0
	for i, char := range dataCode {
		digit := int(char - '0')
		// 使用加权算法：奇数位权重为1，偶数位权重为3
		if i%2 == 0 {
			sum += digit
		} else {
			sum += digit * 3
		}
	}
	// 校验位 = (10 - (总和 % 10)) % 10
	checksum := (10 - (sum % 10)) % 10
	return checksum
}

// GenerateBatchActivationCodes 批量生成激活码
func (g *ActivationCodeGenerator) GenerateBatchActivationCodes(name string, points int32, count int) []string {
	codes := make([]string, count)
	usedCodes := make(map[string]bool) // 用于检测重复

	for i := 0; i < count; i++ {
		var code string
		attempts := 0
		maxAttempts := 100 // 最大尝试次数，防止无限循环

		// 生成唯一激活码，避免重复
		for attempts < maxAttempts {
			// 添加序号和随机延迟确保唯一性
			uniqueName := fmt.Sprintf("%s-%d-%d", name, i, time.Now().UnixNano())
			code = g.GenerateActivationCode(uniqueName, points)

			// 检查是否重复
			if !usedCodes[code] {
				usedCodes[code] = true
				break
			}

			attempts++
			// 短暂延迟，增加时间差异
			time.Sleep(time.Nanosecond)
		}

		codes[i] = code
	}

	return codes
}

// ValidateActivationCodeFormat 验证激活码格式
func (g *ActivationCodeGenerator) ValidateActivationCodeFormat(code string) bool {
	// 检查格式：21位纯数字
	if len(code) < 19 {
		return false
	}

	// 检查每个字符是否为数字
	for _, char := range code {
		if char < '0' || char > '9' {
			return false
		}
	}

	// 验证校验位
	return g.validateChecksum(code)
}

// validateChecksum 验证校验位
func (g *ActivationCodeGenerator) validateChecksum(code string) bool {
	if len(code) < 19 {
		return false
	}

	calculatedChecksum := -1
	expectedChecksum := 0

	if len(code) == 19 {
		// 提取前18位数据
		dataCode := code[:18]
		// 提取校验位
		expectedChecksum = int(code[18] - '0')

		// 计算校验位
		calculatedChecksum = g.calculateChecksum(dataCode)
	} else if len(code) == 21 {
		// 提取前20位数据
		dataCode := code[:20]
		// 提取校验位
		expectedChecksum = int(code[20] - '0')

		// 计算校验位
		calculatedChecksum = g.calculateChecksum(dataCode)
	}

	return calculatedChecksum == expectedChecksum
}

// ExtractInfoFromCode 从激活码中提取信息（用于调试）
func (g *ActivationCodeGenerator) ExtractInfoFromCode(code string) map[string]interface{} {
	if !g.ValidateActivationCodeFormat(code) {
		return nil
	}

	// 解析21位数字激活码：时间戳(6位) + 名称哈希(4位) + 点数(2位) + 随机数(7位) + 保留位(1位) + 校验位(1位)
	timestampPart := code[0:6]
	nameHashPart := code[6:10]
	pointsPart := code[10:12]
	randomPart := code[12:19]
	reservedPart := code[19:20]
	checksumPart := code[20:21]

	// 转换为数值
	timestamp, _ := strconv.ParseUint(timestampPart, 10, 64)
	nameHash, _ := strconv.ParseUint(nameHashPart, 10, 32)
	points, _ := strconv.ParseUint(pointsPart, 10, 32)
	random, _ := strconv.ParseUint(randomPart, 10, 32)
	reserved, _ := strconv.ParseUint(reservedPart, 10, 32)
	checksum, _ := strconv.ParseUint(checksumPart, 10, 32)

	return map[string]any{
		"format":    "21位数字激活码（含校验位）",
		"timestamp": timestampPart,
		"nameHash":  nameHashPart,
		"points":    pointsPart,
		"random":    randomPart,
		"reserved":  reservedPart,
		"checksum":  checksumPart,
		"values": map[string]uint64{
			"timestamp": timestamp,
			"nameHash":  nameHash,
			"points":    points,
			"random":    random,
			"reserved":  reserved,
			"checksum":  checksum,
		},
	}
}
