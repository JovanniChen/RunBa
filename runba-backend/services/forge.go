package services

import (
	"errors"
	"strings"

	"steam-backend/config"
	"steam-backend/models"

	"gorm.io/gorm"
)

// ForgeService 锻刀所服务
// 处理锻刀所相关的业务逻辑
// 说明：数据库表结构由外部维护，不依赖GORM自动迁移
type ForgeService struct {
	db *gorm.DB
}

// NewForgeService 创建新的锻刀所服务实例
func NewForgeService() *ForgeService {
	return &ForgeService{db: config.GetDB()}
}

// CreateForge 创建锻刀所
func (s *ForgeService) CreateForge(req *models.CreateForgeRequest) (*models.Forge, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, errors.New("锻刀所名称不能为空")
	}

	forge := &models.Forge{
		Name:     name,
		ImageURL: strings.TrimSpace(req.ImageURL),
	}

	if err := s.db.Create(forge).Error; err != nil {
		return nil, err
	}

	return forge, nil
}

// GetForges 获取锻刀所列表
func (s *ForgeService) GetForges(params *models.ForgeListRequest) (*models.ForgeListResponse, error) {
	var forges []models.Forge
	var total int64

	query := s.db.Model(&models.Forge{})
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("name LIKE ?", keyword)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, errors.New("查询总数失败")
	}

	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}

	offset := (params.Page - 1) * params.PageSize
	if err := query.Offset(offset).Limit(params.PageSize).Order("created_at DESC").Find(&forges).Error; err != nil {
		return nil, errors.New("查询锻刀所列表失败")
	}

	totalPage := int((total + int64(params.PageSize) - 1) / int64(params.PageSize))

	return &models.ForgeListResponse{
		Forges:    forges,
		Page:      params.Page,
		PageSize:  params.PageSize,
		Total:     total,
		TotalPage: totalPage,
	}, nil
}

// GetForgeByID 根据ID获取锻刀所
func (s *ForgeService) GetForgeByID(id uint) (*models.Forge, error) {
	var forge models.Forge
	if err := s.db.First(&forge, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("锻刀所不存在")
		}
		return nil, errors.New("数据库查询错误")
	}
	return &forge, nil
}

// UpdateForge 更新锻刀所
func (s *ForgeService) UpdateForge(id uint, req *models.UpdateForgeRequest) (*models.Forge, error) {
	forge, err := s.GetForgeByID(id)
	if err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})
	if strings.TrimSpace(req.Name) != "" {
		updates["name"] = strings.TrimSpace(req.Name)
	}
	if strings.TrimSpace(req.ImageURL) != "" {
		updates["image_url"] = strings.TrimSpace(req.ImageURL)
	}

	if len(updates) == 0 {
		return forge, nil
	}

	if err := s.db.Model(forge).Updates(updates).Error; err != nil {
		return nil, errors.New("更新锻刀所失败")
	}

	return s.GetForgeByID(id)
}

// DeleteForge 删除锻刀所
func (s *ForgeService) DeleteForge(id uint) error {
	forge, err := s.GetForgeByID(id)
	if err != nil {
		return err
	}

	if err := s.db.Delete(forge).Error; err != nil {
		return errors.New("删除锻刀所失败")
	}

	return nil
}
