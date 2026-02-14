package services

import (
	"errors"
	"strings"

	"steam-backend/config"
	"steam-backend/models"

	"gorm.io/gorm"
)

// SwordsmithService 刀匠服务
// 处理刀匠相关的业务逻辑
// 说明：数据库表结构由外部维护，不依赖GORM自动迁移
type SwordsmithService struct {
	db *gorm.DB
}

// NewSwordsmithService 创建新的刀匠服务实例
func NewSwordsmithService() *SwordsmithService {
	return &SwordsmithService{db: config.GetDB()}
}

// CreateSwordsmith 创建刀匠
func (s *SwordsmithService) CreateSwordsmith(req *models.CreateSwordsmithRequest) (*models.Swordsmith, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, errors.New("刀匠名称不能为空")
	}

	swordsmith := &models.Swordsmith{
		Name:     name,
		ImageURL: strings.TrimSpace(req.ImageURL),
	}

	if err := s.db.Create(swordsmith).Error; err != nil {
		return nil, err
	}

	return swordsmith, nil
}

// GetSwordsmiths 获取刀匠列表
func (s *SwordsmithService) GetSwordsmiths(params *models.SwordsmithListRequest) (*models.SwordsmithListResponse, error) {
	var swordsmiths []models.Swordsmith
	var total int64

	query := s.db.Model(&models.Swordsmith{})
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
	if err := query.Offset(offset).Limit(params.PageSize).Order("created_at DESC").Find(&swordsmiths).Error; err != nil {
		return nil, errors.New("查询刀匠列表失败")
	}

	totalPage := int((total + int64(params.PageSize) - 1) / int64(params.PageSize))

	return &models.SwordsmithListResponse{
		Swordsmiths: swordsmiths,
		Page:        params.Page,
		PageSize:    params.PageSize,
		Total:       total,
		TotalPage:   totalPage,
	}, nil
}

// GetSwordsmithByID 根据ID获取刀匠
func (s *SwordsmithService) GetSwordsmithByID(id uint) (*models.Swordsmith, error) {
	var swordsmith models.Swordsmith
	if err := s.db.First(&swordsmith, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("刀匠不存在")
		}
		return nil, errors.New("数据库查询错误")
	}
	return &swordsmith, nil
}

// UpdateSwordsmith 更新刀匠
func (s *SwordsmithService) UpdateSwordsmith(id uint, req *models.UpdateSwordsmithRequest) (*models.Swordsmith, error) {
	swordsmith, err := s.GetSwordsmithByID(id)
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
		return swordsmith, nil
	}

	if err := s.db.Model(swordsmith).Updates(updates).Error; err != nil {
		return nil, errors.New("更新刀匠失败")
	}

	return s.GetSwordsmithByID(id)
}

// DeleteSwordsmith 删除刀匠
func (s *SwordsmithService) DeleteSwordsmith(id uint) error {
	swordsmith, err := s.GetSwordsmithByID(id)
	if err != nil {
		return err
	}

	if err := s.db.Delete(swordsmith).Error; err != nil {
		return errors.New("删除刀匠失败")
	}

	return nil
}
