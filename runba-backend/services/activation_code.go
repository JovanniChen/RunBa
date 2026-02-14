package services

import (
	"errors"
	"fmt"
	"steam-backend/config"
	"steam-backend/logger"
	"steam-backend/models"
	"steam-backend/utils"
	"time"

	"gorm.io/gorm"
)

// ActivationCodeService 激活码服务结构体
type ActivationCodeService struct {
	db            *gorm.DB
	codeGenerator *utils.ActivationCodeGenerator
}

// NewActivationCodeService 创建新的激活码服务实例
func NewActivationCodeService() *ActivationCodeService {
	return &ActivationCodeService{
		db:            config.GetDB(),
		codeGenerator: utils.NewActivationCodeGenerator(),
	}
}

// CreateActivationCode 创建新的激活码
func (s *ActivationCodeService) CreateActivationCode(req *models.CreateActivationCodeRequest) (*models.ActivationCode, error) {
	// 调用批量创建方法，设置count为1
	batchReq := &models.CreateBatchActivationCodeRequest{
		Count:  1,
		Points: req.Points,
	}

	activationCodes, err := s.CreateBatchActivationCodes(batchReq)
	if err != nil {
		return nil, err
	}

	// 返回第一个（也是唯一一个）激活码
	return &activationCodes[0], nil
}

// generateUniqueActivationCode 生成唯一的激活码
func (s *ActivationCodeService) generateUniqueActivationCode(name string, points int32) string {
	maxAttempts := 10
	for attempt := 0; attempt < maxAttempts; attempt++ {
		// 生成激活码
		code := s.codeGenerator.GenerateActivationCode(name, points)

		// 检查是否已存在
		var existingCode models.ActivationCode
		if err := s.db.Where("code = ?", code).First(&existingCode).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// 激活码不存在，可以使用
				return code
			}
			// 数据库错误，记录日志但继续尝试
			logger.Error("检查激活码唯一性时发生错误: %v", err)
		}

		// 激活码已存在，继续尝试
		logger.Info("生成的激活码已存在，重新生成: %s", code)
	}

	// 如果多次尝试都失败，使用时间戳作为后缀
	timestamp := time.Now().UnixNano()
	fallbackCode := s.codeGenerator.GenerateActivationCode(fmt.Sprintf("%s-%d", name, timestamp), points)
	logger.Warn("多次生成激活码失败，使用备用方案: %s", fallbackCode)

	return fallbackCode
}

// CreateBatchActivationCodes 批量创建激活码
func (s *ActivationCodeService) CreateBatchActivationCodes(req *models.CreateBatchActivationCodeRequest) ([]models.ActivationCode, error) {
	name := fmt.Sprintf("激活码(%d点数)", req.Points)
	var activationCodes []models.ActivationCode

	// 使用事务确保数据一致性
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 预分配切片容量，避免频繁扩容
	activationCodes = make([]models.ActivationCode, 0, req.Count)

	// 批量生成激活码
	for i := 0; i < req.Count; i++ {
		// 生成唯一的激活码
		generatedCode := s.generateUniqueActivationCode(name, req.Points)

		// 创建激活码对象（不立即插入数据库）
		activationCode := models.ActivationCode{
			Name:   name,
			Code:   generatedCode,
			Points: req.Points,
		}
		activationCodes = append(activationCodes, activationCode)
	}

	// 使用批量插入，一次性插入所有记录
	if err := tx.CreateInBatches(activationCodes, 100).Error; err != nil {
		tx.Rollback()
		logger.Error("批量创建激活码失败: %v", err)
		return nil, errors.New("批量创建激活码失败")
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		logger.Error("提交批量创建激活码事务失败: %v", err)
		return nil, errors.New("提交事务失败")
	}

	logger.Info("成功批量创建激活码: 数量=%d, 点数=%d",
		len(activationCodes), req.Points)

	return activationCodes, nil
}

// GetActivationCodes 获取激活码列表（带分页和搜索）
func (s *ActivationCodeService) GetActivationCodes(params *models.ActivationCodeListRequest) (*models.ActivationCodeListResponse, error) {
	var activationCodes []models.ActivationCode
	var total int64

	// 构建查询
	query := s.db.Model(&models.ActivationCode{})

	// 添加状态条件
	if params.Status != 0 {
		query = query.Where("status = ?", params.Status)
	}

	// 添加搜索条件
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("name LIKE ? OR code LIKE ?", keyword, keyword)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		logger.Error("查询激活码总数失败: %v", err)
		return nil, errors.New("查询总数失败")
	}

	// 设置默认分页参数
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}

	// 分页查询
	offset := (params.Page - 1) * params.PageSize
	if err := query.Offset(offset).Limit(params.PageSize).Order("created_at DESC").Find(&activationCodes).Error; err != nil {
		logger.Error("查询激活码列表失败: %v", err)
		return nil, errors.New("查询激活码列表失败")
	}

	// 计算总页数
	totalPage := int((total + int64(params.PageSize) - 1) / int64(params.PageSize))

	logger.Info("查询激活码列表成功: 总数=%d, 当前页=%d, 每页大小=%d",
		total, params.Page, params.PageSize)

	return &models.ActivationCodeListResponse{
		ActivationCodes: activationCodes,
		Page:            params.Page,
		PageSize:        params.PageSize,
		Total:           total,
		TotalPage:       totalPage,
	}, nil
}

// GetActivationCodeByID 根据ID获取激活码信息
func (s *ActivationCodeService) GetActivationCodeByID(codeID uint) (*models.ActivationCode, error) {
	var activationCode models.ActivationCode

	if err := s.db.First(&activationCode, codeID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("激活码不存在")
		}
		logger.Error("查询激活码失败: %v", err)
		return nil, errors.New("数据库查询错误")
	}

	logger.Info("成功查询激活码: ID=%d, 名称=%s, 代码=%s",
		activationCode.ID, activationCode.Name, activationCode.Code)

	return &activationCode, nil
}

// GetActivationCodeByCode 根据激活码获取信息
func (s *ActivationCodeService) GetActivationCodeByCode(code string) (*models.ActivationCode, error) {
	var activationCode models.ActivationCode

	if err := s.db.Where("code = ?", code).First(&activationCode).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("激活码不存在")
		}
		logger.Error("查询激活码失败: %v", err)
		return nil, errors.New("数据库查询错误")
	}

	logger.Info("成功查询激活码: ID=%d, 名称=%s, 代码=%s",
		activationCode.ID, activationCode.Name, activationCode.Code)

	return &activationCode, nil
}

// UpdateActivationCode 更新激活码信息
func (s *ActivationCodeService) UpdateActivationCode(codeID uint, req *models.UpdateActivationCodeRequest) (*models.ActivationCode, error) {
	// 检查激活码是否存在
	var activationCode models.ActivationCode
	if err := s.db.First(&activationCode, codeID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("激活码不存在")
		}
		logger.Error("查询激活码失败: %v", err)
		return nil, errors.New("数据库查询错误")
	}

	if activationCode.Status == models.StatusUsedOfActivationCode || activationCode.Status == models.StatusDropOfActivationCode {
		return nil, errors.New("该激活码已使用或已废弃")
	}

	// 更新字段
	updates := make(map[string]interface{})
	if req.Status > 0 {
		updates["status"] = req.Status
	}

	// 执行更新
	if err := s.db.Model(&activationCode).Updates(updates).Error; err != nil {
		logger.Error("更新激活码失败: %v", err)
		return nil, errors.New("更新激活码失败")
	}

	// 重新查询更新后的数据
	if err := s.db.First(&activationCode, codeID).Error; err != nil {
		logger.Error("查询更新后的激活码失败: %v", err)
		return nil, errors.New("查询更新后的激活码失败")
	}

	logger.Info("成功更新激活码: ID=%d, 名称=%s, 代码=%s, 点数=%d",
		activationCode.ID, activationCode.Name, activationCode.Code, activationCode.Points)

	return &activationCode, nil
}

// DeleteActivationCode 删除激活码（软删除）
func (s *ActivationCodeService) DeleteActivationCode(codeID uint) error {
	// 检查激活码是否存在
	var activationCode models.ActivationCode
	if err := s.db.First(&activationCode, codeID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("激活码不存在")
		}
		logger.Error("查询激活码失败: %v", err)
		return errors.New("数据库查询错误")
	}

	// 软删除
	if err := s.db.Delete(&activationCode).Error; err != nil {
		logger.Error("删除激活码失败: %v", err)
		return errors.New("删除激活码失败")
	}

	logger.Info("成功删除激活码: ID=%d, 名称=%s, 代码=%s",
		activationCode.ID, activationCode.Name, activationCode.Code)

	return nil
}

// GetActivationCodeStats 获取激活码统计信息
func (s *ActivationCodeService) GetActivationCodeStats() (map[string]interface{}, error) {
	var totalCount int64
	var totalPoints int64
	var avgPoints float64

	// 统计总数
	if err := s.db.Model(&models.ActivationCode{}).Count(&totalCount).Error; err != nil {
		logger.Error("统计激活码总数失败: %v", err)
		return nil, errors.New("统计激活码总数失败")
	}

	// 统计总点数
	if err := s.db.Model(&models.ActivationCode{}).Select("SUM(points)").Scan(&totalPoints).Error; err != nil {
		logger.Error("统计激活码总点数失败: %v", err)
		return nil, errors.New("统计激活码总点数失败")
	}

	// 计算平均点数
	if totalCount > 0 {
		avgPoints = float64(totalPoints) / float64(totalCount)
	}

	stats := map[string]interface{}{
		"total_count":  totalCount,
		"total_points": totalPoints,
		"avg_points":   avgPoints,
	}

	logger.Info("激活码统计信息: 总数=%d, 总点数=%d, 平均点数=%.2f",
		totalCount, totalPoints, avgPoints)

	return stats, nil
}

// ValidateActivationCode 验证激活码是否有效
func (s *ActivationCodeService) ValidateActivationCode(code string) (bool, *models.ActivationCode, error) {
	activationCode, err := s.GetActivationCodeByCode(code)
	if err != nil {
		return false, nil, err
	}

	// 这里可以添加更多的验证逻辑，比如：
	// - 检查激活码是否过期
	// - 检查激活码是否已被使用
	// - 检查激活码的使用次数限制等

	return true, activationCode, nil
}
