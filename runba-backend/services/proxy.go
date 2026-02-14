package services

import (
	"errors"
	"math/rand"
	"steam-backend/config"
	"steam-backend/models"
	"time"

	"gorm.io/gorm"
)

// ProxyService 代理服务结构体
type ProxyService struct {
	db *gorm.DB
}

// NewProxyService 创建新的代理服务实例
func NewProxyService() *ProxyService {
	return &ProxyService{
		db: config.GetDB(),
	}
}

// CreateProxy 创建新的代理
func (s *ProxyService) CreateProxy(req *models.CreateProxyRequest) (*models.Proxy, error) {
	if req.Host == "" {
		return nil, errors.New("代理主机地址不能为空")
	}
	if req.Port <= 0 || req.Port > 65535 {
		return nil, errors.New("代理端口必须在1-65535范围内")
	}

	// 检查代理主机和端口组合是否已存在
	var existingProxy models.Proxy
	if err := s.db.Where("host = ? AND port = ?", req.Host, req.Port).First(&existingProxy).Error; err == nil {
		return nil, errors.New("代理主机和端口组合已存在")
	}

	proxy := &models.Proxy{
		Host:     req.Host,
		Port:     req.Port,
		Category: req.Category,
		Protocol: req.Protocol,
		Status:   req.Status,
		Username: req.Username,
		Password: req.Password,
		Remark:   req.Remark,
	}

	// 如果没有设置状态，默认启用
	if proxy.Status == 0 && req.Status == 0 {
		proxy.Status = models.ProxyStatusEnabled
	}

	if err := s.db.Create(proxy).Error; err != nil {
		return nil, err
	}

	return proxy, nil
}

// GetProxyByID 根据ID获取代理
func (s *ProxyService) GetProxyByID(id uint) (*models.Proxy, error) {
	var proxy models.Proxy
	if err := s.db.First(&proxy, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("代理不存在")
		}
		return nil, err
	}
	return &proxy, nil
}

// GetProxyList 获取代理列表
func (s *ProxyService) GetProxyList(req *models.ProxyListRequest) (*models.ProxyListResponse, error) {
	var proxies []models.Proxy
	var total int64

	query := s.db.Model(&models.Proxy{})

	// 关键词搜索
	if req.Keyword != "" {
		query = query.Where("host LIKE ? OR remark LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// 主机地址筛选
	if req.Host != "" {
		query = query.Where("host = ?", req.Host)
	}

	// 分类筛选
	if req.Category != 0 {
		query = query.Where("category = ?", req.Category)
	}

	// 协议筛选
	if req.Protocol != 0 {
		query = query.Where("protocol = ?", req.Protocol)
	}

	// 状态筛选
	if req.Status != 0 {
		query = query.Where("status = ?", req.Status)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&proxies).Error; err != nil {
		return nil, err
	}

	// 计算总页数
	totalPage := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))

	return &models.ProxyListResponse{
		Proxies:   proxies,
		Page:      req.Page,
		PageSize:  req.PageSize,
		Total:     total,
		TotalPage: totalPage,
	}, nil
}

// UpdateProxy 更新代理
func (s *ProxyService) UpdateProxy(id uint, req *models.UpdateProxyRequest) (*models.Proxy, error) {
	// 检查代理是否存在
	proxy, err := s.GetProxyByID(id)
	if err != nil {
		return nil, err
	}

	// 如果要更新主机或端口，检查新的主机端口组合是否已被其他代理使用
	if (req.Host != "" && req.Host != proxy.Host) || (req.Port != nil && *req.Port != proxy.Port) {
		checkHost := proxy.Host
		checkPort := proxy.Port

		if req.Host != "" {
			checkHost = req.Host
		}
		if req.Port != nil {
			checkPort = *req.Port
		}

		var existingProxy models.Proxy
		if err := s.db.Where("host = ? AND port = ? AND id != ?", checkHost, checkPort, id).First(&existingProxy).Error; err == nil {
			return nil, errors.New("代理主机和端口组合已被其他代理使用")
		}
	}

	// 更新字段
	updates := make(map[string]interface{})
	if req.Host != "" {
		updates["host"] = req.Host
	}
	if req.Port != nil {
		if *req.Port <= 0 || *req.Port > 65535 {
			return nil, errors.New("代理端口必须在1-65535范围内")
		}
		updates["port"] = *req.Port
	}
	if req.Category != 0 {
		updates["category"] = req.Category
	}
	if req.Protocol != 0 {
		updates["protocol"] = req.Protocol
	}
	if req.Status != 0 {
		updates["status"] = req.Status
	}
	if req.Username != "" {
		updates["username"] = req.Username
	}
	if req.Password != "" {
		updates["password"] = req.Password
	}
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}

	if len(updates) == 0 {
		return proxy, nil
	}

	// 执行更新
	if err := s.db.Model(proxy).Updates(updates).Error; err != nil {
		return nil, err
	}

	// 重新获取更新后的数据
	return s.GetProxyByID(id)
}

// DeleteProxy 删除代理（软删除）
func (s *ProxyService) DeleteProxy(id uint) error {
	// 检查代理是否存在
	_, err := s.GetProxyByID(id)
	if err != nil {
		return err
	}

	// 软删除
	if err := s.db.Delete(&models.Proxy{}, id).Error; err != nil {
		return err
	}

	return nil
}

// UpdateProxyStatus 更新代理状态
func (s *ProxyService) UpdateProxyStatus(id uint, status models.ProxyStatus) error {
	// 检查代理是否存在
	_, err := s.GetProxyByID(id)
	if err != nil {
		return err
	}

	// 更新状态
	if err := s.db.Model(&models.Proxy{}).Where("id = ?", id).Update("status", status).Error; err != nil {
		return err
	}

	return nil
}

// GetEnabledProxies 获取所有启用的代理
func (s *ProxyService) GetEnabledProxies() ([]models.Proxy, error) {
	var proxies []models.Proxy
	if err := s.db.Where("status = ?", models.ProxyStatusEnabled).Order("created_at DESC").Find(&proxies).Error; err != nil {
		return nil, err
	}
	return proxies, nil
}

// GetProxiesByCategory 根据分类获取代理
func (s *ProxyService) GetProxiesByCategory(category models.ProxyCategory) ([]models.Proxy, error) {
	var proxies []models.Proxy
	if err := s.db.Where("category = ? AND status = ?", category, models.ProxyStatusEnabled).Order("created_at DESC").Find(&proxies).Error; err != nil {
		return nil, err
	}
	return proxies, nil
}

// GetProxiesByProtocol 根据协议获取代理
func (s *ProxyService) GetProxiesByProtocol(protocol models.ProxyProtocol) ([]models.Proxy, error) {
	var proxies []models.Proxy
	if err := s.db.Where("protocol = ? AND status = ?", protocol, models.ProxyStatusEnabled).Order("created_at DESC").Find(&proxies).Error; err != nil {
		return nil, err
	}
	return proxies, nil
}

// GetRandomProxy 根据协议和分类随机获取一个代理
func (s *ProxyService) GetRandomProxy(protocol *models.ProxyProtocol, category *models.ProxyCategory) (*models.Proxy, error) {
	var proxies []models.Proxy

	// 构建查询条件，默认只查询启用的代理
	query := s.db.Where("status = ?", models.ProxyStatusEnabled)

	// 添加协议筛选条件
	if protocol != nil {
		query = query.Where("protocol = ?", *protocol)
	}

	// 添加分类筛选条件
	if category != nil {
		query = query.Where("category = ?", *category)
	}

	// 执行查询
	if err := query.Find(&proxies).Error; err != nil {
		return nil, err
	}

	// 检查是否有符合条件的代理
	if len(proxies) == 0 {
		// 构建详细的错误信息
		var condition string
		if protocol != nil && category != nil {
			condition = "指定协议和分类"
		} else if protocol != nil {
			condition = "指定协议"
		} else if category != nil {
			condition = "指定分类"
		} else {
			condition = "启用状态"
		}
		return nil, errors.New("没有找到符合" + condition + "条件的代理")
	}

	// 使用新的随机数生成器（Go 1.20+）
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	// 随机选择一个代理
	randomIndex := rng.Intn(len(proxies))

	return &proxies[randomIndex], nil
}
