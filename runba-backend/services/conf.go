package services

import (
	"errors"
	"steam-backend/config"
	"steam-backend/models"
	"sync"

	"gorm.io/gorm"
)

// ConfService 配置服务结构体
type ConfService struct {
	db         *gorm.DB
	cache      *models.Conf // 缓存最新配置
	cacheMutex sync.RWMutex // 缓存读写锁
}

// NewConfService 创建新的配置服务实例
func NewConfService() *ConfService {
	return &ConfService{
		db:         config.GetDB(),
		cache:      nil,
		cacheMutex: sync.RWMutex{},
	}
}

// CreateConf 创建新的配置
func (s *ConfService) CreateConf(req *models.CreateConfRequest) (*models.Conf, error) {
	if req.TutorialUrl == "" {
		return nil, errors.New("提货教程地址不能为空")
	}

	conf := &models.Conf{
		TutorialUrl: req.TutorialUrl,
	}

	if err := s.db.Create(conf).Error; err != nil {
		return nil, err
	}

	// 清除缓存，因为有新的配置创建
	s.clearCache()

	return conf, nil
}

// GetConfByID 根据ID获取配置
func (s *ConfService) GetConfByID(id uint) (*models.Conf, error) {
	var conf models.Conf
	if err := s.db.First(&conf, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("配置不存在")
		}
		return nil, err
	}
	return &conf, nil
}

// GetConfList 获取配置列表
func (s *ConfService) GetConfList(req *models.ConfListRequest) (*models.ConfListResponse, error) {
	var confs []models.Conf
	var total int64

	query := s.db.Model(&models.Conf{})

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&confs).Error; err != nil {
		return nil, err
	}

	// 计算总页数
	totalPage := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))

	return &models.ConfListResponse{
		Confs:     confs,
		Page:      req.Page,
		PageSize:  req.PageSize,
		Total:     total,
		TotalPage: totalPage,
	}, nil
}

// UpdateConf 更新配置
func (s *ConfService) UpdateConf(id uint, req *models.UpdateConfRequest) (*models.Conf, error) {
	// 检查配置是否存在
	conf, err := s.GetConfByID(id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	updates := make(map[string]interface{})
	if req.TutorialUrl != "" {
		updates["tutorial_url"] = req.TutorialUrl
	}

	if len(updates) == 0 {
		return conf, nil
	}

	// 执行更新
	if err := s.db.Model(conf).Updates(updates).Error; err != nil {
		return nil, err
	}

	// 清除缓存，因为配置已更新
	s.clearCache()

	// 重新获取更新后的数据
	return s.GetConfByID(id)
}

// DeleteConf 删除配置（软删除）
func (s *ConfService) DeleteConf(id uint) error {
	// 检查配置是否存在
	_, err := s.GetConfByID(id)
	if err != nil {
		return err
	}

	// 软删除
	if err := s.db.Delete(&models.Conf{}, id).Error; err != nil {
		return err
	}

	// 清除缓存，因为配置已删除
	s.clearCache()

	return nil
}

// CreateOrUpdateConf 创建或更新配置
func (s *ConfService) CreateOrUpdateConf(req *models.CreateConfRequest) (*models.Conf, error) {
	if req.TutorialUrl == "" {
		return nil, errors.New("提货教程地址不能为空")
	}

	// 尝试获取最新配置
	existingConf, err := s.GetLatestConf()
	if err != nil && err.Error() != "暂无配置数据" {
		return nil, err
	}

	// 如果存在配置，则更新
	if existingConf != nil {
		updateReq := &models.UpdateConfRequest{
			TutorialUrl: req.TutorialUrl,
		}
		return s.UpdateConf(existingConf.ID, updateReq)
	}

	// 如果不存在配置，则创建新配置
	conf := &models.Conf{
		TutorialUrl: req.TutorialUrl,
	}

	if err := s.db.Create(conf).Error; err != nil {
		return nil, err
	}

	// 清除缓存，因为有新的配置创建
	s.clearCache()

	return conf, nil
}

// GetLatestConf 获取最新的配置（带缓存）
func (s *ConfService) GetLatestConf() (*models.Conf, error) {
	// 尝试从缓存中读取
	s.cacheMutex.RLock()
	if s.cache != nil {
		// 复制缓存数据以避免并发修改
		cachedConf := &models.Conf{
			ID:          s.cache.ID,
			TutorialUrl: s.cache.TutorialUrl,
			CreatedAt:   s.cache.CreatedAt,
			UpdatedAt:   s.cache.UpdatedAt,
			DeletedAt:   s.cache.DeletedAt,
		}
		s.cacheMutex.RUnlock()
		return cachedConf, nil
	}
	s.cacheMutex.RUnlock()

	// 缓存未命中，从数据库读取
	var conf models.Conf
	if err := s.db.Order("created_at DESC").First(&conf).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("暂无配置数据")
		}
		return nil, err
	}

	// 更新缓存
	s.setCache(&conf)

	return &conf, nil
}

// setCache 设置缓存（内部方法）
func (s *ConfService) setCache(conf *models.Conf) {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	// 复制配置数据到缓存
	s.cache = &models.Conf{
		ID:          conf.ID,
		TutorialUrl: conf.TutorialUrl,
		CreatedAt:   conf.CreatedAt,
		UpdatedAt:   conf.UpdatedAt,
		DeletedAt:   conf.DeletedAt,
	}
}

// clearCache 清除缓存（内部方法）
func (s *ConfService) clearCache() {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()
	s.cache = nil
}
