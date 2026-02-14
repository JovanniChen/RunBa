package services

import (
	"errors"
	"fmt"
	"steam-backend/config"
	"steam-backend/logger"
	"steam-backend/models"
	"time"

	"gorm.io/gorm"
)

// AccountService Steam账户服务结构体
// 处理Steam账户相关的业务逻辑
type AccountService struct {
	db           *gorm.DB
	proxyService *ProxyService // 代理服务
}

// NewAccountService 创建新的账户服务实例
func NewAccountService(proxyService *ProxyService) *AccountService {
	logger.Info("正在初始化AccountService...")

	db := config.GetDB()
	if db == nil {
		logger.Error("数据库连接获取失败，AccountService可能无法正常工作")
	} else {
		logger.Debug("数据库连接获取成功")
	}

	service := &AccountService{
		db:           db,
		proxyService: proxyService,
	}

	logger.Info("AccountService初始化完成")
	return service
}

// CreateAccount 创建新的Steam账户（异步验证模式）
func (s *AccountService) CreateAccount(req *models.CreateAccountRequest) (*models.SteamAccount, error) {
	// 验证TokenContent是否有效
	if req.TokenContent == nil {
		return nil, errors.New("令牌内容不能为空")
	}

	// 检查用户名是否已存在
	var existingAccount models.SteamAccount
	err := s.db.Unscoped().Where("username = ?", req.Username).First(&existingAccount).Error
	if err == nil {
		if existingAccount.DeletedAt.Time.IsZero() {
			return nil, errors.New("用户名已存在")
		}
		// 如果是软删除的记录，可以考虑硬删除
		if err := s.DeleteAccount(existingAccount.ID); err != nil {
			return nil, errors.New("硬删除失败，请联系管理员")
		}
	}

	// 先创建账户记录（状态为验证中）
	account := models.SteamAccount{
		Username:    req.Username,
		Password:    req.Password,
		Status:      models.StatusIsLoginingOfAccount, // 登录中状态
		AccountNote: req.AccountNote,
	}

	// 设置令牌内容
	if err := account.SetTokenContent(req.TokenContent); err != nil {
		return nil, fmt.Errorf("设置令牌内容失败: %v", err)
	}

	// 保存到数据库
	if err := s.db.Create(&account).Error; err != nil {
		return nil, fmt.Errorf("创建账户记录失败: %v", err)
	}

	logger.Info("账户记录已创建，开始异步验证Steam登录: ID=%d, Username=%s", account.ID, account.Username)

	return &account, nil
}

// GetAccounts 获取Steam账户列表（带分页和搜索）
func (s *AccountService) GetAccounts(params *models.AccountListRequest) (*models.AccountListResponse, error) {
	var accounts []models.SteamAccount
	var total int64

	// 构建查询
	query := s.db.Model(&models.SteamAccount{})

	// 添加搜索条件
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("username LIKE ? OR nickname LIKE ?", keyword, keyword)
	}

	// 添加状态筛选
	if params.Status != 0 {
		// 使用模运算检查是否包含指定状态
		query = query.Where("status % ? = 0", uint64(params.Status))
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
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
	if err := query.Offset(offset).Limit(params.PageSize).Order("created_at DESC").Find(&accounts).Error; err != nil {
		return nil, errors.New("查询Steam账户列表失败")
	}

	// 计算总页数
	totalPage := int((total + int64(params.PageSize) - 1) / int64(params.PageSize))

	return &models.AccountListResponse{
		Accounts:  accounts,
		Page:      params.Page,
		PageSize:  params.PageSize,
		Total:     total,
		TotalPage: totalPage,
	}, nil
}

// GetAllAccounts 获取所有Steam账户列表
func (s *AccountService) GetAllAccounts() (*models.AccountListResponse, error) {
	var accounts []models.SteamAccount
	var total int64

	// 构建查询
	query := s.db.Model(&models.SteamAccount{})

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, errors.New("查询总数失败")
	}

	// 查询所有账户，按创建时间倒序
	if err := query.Order("created_at DESC").Find(&accounts).Error; err != nil {
		return nil, errors.New("查询所有Steam账户失败")
	}

	// 创建分页信息（显示所有数据）
	pageSize := int(total) // 设置为总数，表示显示所有
	totalPage := 1

	return &models.AccountListResponse{
		Accounts:  accounts,
		Page:      1,
		PageSize:  pageSize,
		Total:     total,
		TotalPage: totalPage,
	}, nil
}

// GetAccountByID 根据ID获取Steam账户信息
func (s *AccountService) GetAccountByID(accountID uint) (*models.SteamAccount, error) {
	var account models.SteamAccount

	if err := s.db.First(&account, accountID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("Steam账户不存在")
		}
		return nil, errors.New("数据库查询错误")
	}

	return &account, nil
}

// UpdateAccount 更新Steam账户信息
func (s *AccountService) UpdateAccount(accountID uint, req *models.UpdateAccountRequest) (*models.SteamAccount, error) {
	var account models.SteamAccount

	// 查询账户
	if err := s.db.First(&account, accountID).Error; err != nil {
		return nil, errors.New("钱包不存在")
	}

	// 更新账户信息
	updates := make(map[string]interface{})
	if req.Username != "" {
		// 检查新用户名是否已存在
		var existingAccount models.SteamAccount
		if err := s.db.Where("username = ? AND id != ?", req.Username, accountID).First(&existingAccount).Error; err == nil {
			return nil, errors.New("用户名已存在")
		}
		updates["username"] = req.Username
	}
	if req.Password != "" {
		updates["password"] = req.Password
	}
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.CountryCode != "" {
		updates["country_code"] = req.CountryCode
	}

	// 处理TokenContent更新
	if req.TokenContent != nil {
		if err := account.SetTokenContent(req.TokenContent); err != nil {
			return nil, errors.New("设置令牌内容失败: " + err.Error())
		}
		updates["token_content"] = account.TokenContent
	}

	if len(updates) > 0 {
		if err := s.db.Model(&account).Updates(updates).Error; err != nil {
			return nil, errors.New("Steam账户信息更新失败")
		}
	}

	// 重新查询更新后的账户信息
	if err := s.db.First(&account, accountID).Error; err != nil {
		return nil, errors.New("获取更新后的账户信息失败")
	}

	return &account, nil
}

// DeleteAccount 删除Steam账户（软删除）
func (s *AccountService) DeleteAccount(accountID uint) error {
	var account models.SteamAccount

	// 检查账户是否存在
	if err := s.db.Unscoped().First(&account, accountID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("Steam账户不存在")
		}
		return errors.New("数据库查询错误")
	}

	// 执行软删除
	if err := s.db.Unscoped().Delete(&account).Error; err != nil {
		return errors.New("删除Steam账户失败")
	}

	return nil
}

// UpdateAccountStatus 更新Steam账户状态（直接设置，会覆盖所有状态）
func (s *AccountService) UpdateAccountStatus(accountID uint, status uint64) (*models.SteamAccount, error) {
	var account models.SteamAccount

	// 查询账户
	if err := s.db.First(&account, accountID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("Steam账户不存在")
		}
		return nil, errors.New("数据库查询错误")
	}

	// 直接设置状态（覆盖模式，用于管理员完全重置状态）
	if err := s.db.Model(&account).Update("status", status).Error; err != nil {
		return nil, errors.New("Steam账户状态更新失败")
	}

	// 重新查询更新后的账户信息
	if err := s.db.First(&account, accountID).Error; err != nil {
		return nil, errors.New("获取更新后的账户信息失败")
	}

	return &account, nil
}

// AddAccountStatus 添加账户状态（支持复合状态）
func (s *AccountService) AddAccountStatus(accountID uint, status uint64) (*models.SteamAccount, error) {
	var account models.SteamAccount

	// 查询账户
	if err := s.db.First(&account, accountID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("Steam账户不存在")
		}
		return nil, errors.New("数据库查询错误")
	}

	// 添加状态（保持其他状态）
	account.AddStatus(status)

	// 更新数据库
	if err := s.db.Model(&account).Update("status", account.Status).Error; err != nil {
		return nil, errors.New("添加Steam账户状态失败")
	}

	return &account, nil
}

// RemoveAccountStatus 移除账户状态（支持复合状态）
func (s *AccountService) RemoveAccountStatus(accountID uint, status uint64) (*models.SteamAccount, error) {
	var account models.SteamAccount

	// 查询账户
	if err := s.db.First(&account, accountID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("Steam账户不存在")
		}
		return nil, errors.New("数据库查询错误")
	}

	// 移除状态（保持其他状态）
	account.RemoveStatus(status)

	// 更新数据库
	if err := s.db.Model(&account).Update("status", account.Status).Error; err != nil {
		return nil, errors.New("移除Steam账户状态失败")
	}

	return &account, nil
}

// SetAccountStatusSafe 安全设置账户状态（智能状态转换）
func (s *AccountService) SetAccountStatusSafe(accountID uint, newStatus uint64) (*models.SteamAccount, error) {
	var account models.SteamAccount

	// 查询账户
	if err := s.db.First(&account, accountID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("Steam账户不存在")
		}
		return nil, errors.New("数据库查询错误")
	}

	// 智能状态转换：移除冲突状态，添加新状态
	// 使用switch语句处理不同状态类型的冲突逻辑
	switch newStatus {
	case models.StatusNormalOfAccount:
		// 设置正常状态：移除禁用状态
		account.RemoveStatus(models.StatusBannedOfAccount)
		account.AddStatus(models.StatusNormalOfAccount)
	case models.StatusBannedOfAccount:
		// 设置禁用状态：移除正常状态
		account.RemoveStatus(models.StatusNormalOfAccount)
		account.AddStatus(models.StatusBannedOfAccount)
	}

	fmt.Println("anewStatus", newStatus)
	fmt.Println("account.Status", account.Status)

	// 更新数据库
	if err := s.db.Model(&account).Update("status", account.Status).Error; err != nil {
		return nil, errors.New("设置Steam账户状态失败")
	}

	return &account, nil
}

// UpdateLastLogin 更新最后登录时间
func (s *AccountService) UpdateLastLogin(accountID uint) error {
	now := time.Now()
	if err := s.db.Model(&models.SteamAccount{}).Where("id = ?", accountID).Update("last_login_at", &now).Error; err != nil {
		return errors.New("更新最后登录时间失败")
	}
	return nil
}

// AddPoints 增加积分
func (s *AccountService) AddPoints(accountID uint, points int32) error {
	if points <= 0 {
		return errors.New("积分必须大于0")
	}

	if err := s.db.Model(&models.SteamAccount{}).Where("id = ?", accountID).Update("points", gorm.Expr("points + ?", points)).Error; err != nil {
		return errors.New("增加积分失败")
	}
	return nil
}

// DeductPoints 扣除积分
func (s *AccountService) DeductPoints(accountID uint, points int32) error {
	if points <= 0 {
		return errors.New("积分必须大于0")
	}

	var account models.SteamAccount
	if err := s.db.First(&account, accountID).Error; err != nil {
		return errors.New("Steam账户不存在")
	}

	if account.Points < points {
		return errors.New("积分不足")
	}

	if err := s.db.Model(&account).Update("points", gorm.Expr("points - ?", points)).Error; err != nil {
		return errors.New("扣除积分失败")
	}
	return nil
}

// GetAccountStats 获取Steam账户统计信息
func (s *AccountService) GetAccountStats() (map[string]interface{}, error) {
	var stats = make(map[string]interface{})

	// 总账户数
	var totalAccounts int64
	if err := s.db.Model(&models.SteamAccount{}).Count(&totalAccounts).Error; err != nil {
		return nil, errors.New("查询总账户数失败")
	}
	stats["total_accounts"] = totalAccounts

	// 活跃账户数
	var activeAccounts int64
	if err := s.db.Model(&models.SteamAccount{}).Where("status % ? = 0", models.StatusNormalOfAccount).Count(&activeAccounts).Error; err != nil {
		return nil, errors.New("查询活跃账户数失败")
	}
	stats["active_accounts"] = activeAccounts

	// 今日新增账户数
	today := time.Now().Format("2006-01-02")
	var todayNewAccounts int64
	if err := s.db.Model(&models.SteamAccount{}).Where("DATE(created_at) = ?", today).Count(&todayNewAccounts).Error; err != nil {
		return nil, errors.New("查询今日新增账户数失败")
	}
	stats["today_new_accounts"] = todayNewAccounts

	// 停用账户数
	var inactiveAccounts int64
	if err := s.db.Model(&models.SteamAccount{}).Where("status % ? = 0", models.StatusBannedOfAccount).Count(&inactiveAccounts).Error; err != nil {
		return nil, errors.New("查询停用账户数失败")
	}
	stats["inactive_accounts"] = inactiveAccounts

	return stats, nil
}

// UpdateTokenContent 更新账户的令牌内容
func (s *AccountService) UpdateTokenContent(accountID uint, tokenContent *models.TokenContent) error {
	var account models.SteamAccount

	// 查询账户
	if err := s.db.First(&account, accountID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("Steam账户不存在")
		}
		return errors.New("数据库查询错误")
	}

	// 设置令牌内容
	if err := account.SetTokenContent(tokenContent); err != nil {
		return errors.New("设置令牌内容失败: " + err.Error())
	}

	// 更新数据库
	if err := s.db.Model(&account).Update("token_content", account.TokenContent).Error; err != nil {
		return errors.New("更新令牌内容失败")
	}

	return nil
}

// GetTokenContent 获取账户的令牌内容
func (s *AccountService) GetTokenContent(accountID uint) (*models.TokenContent, error) {
	var account models.SteamAccount

	// 查询账户
	if err := s.db.First(&account, accountID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("Steam账户不存在")
		}
		return nil, errors.New("数据库查询错误")
	}

	// 获取令牌内容
	tokenContent, err := account.GetTokenContent()
	if err != nil {
		return nil, errors.New("解析令牌内容失败: " + err.Error())
	}

	return tokenContent, nil
}

// ValidateToken 验证账户的令牌是否有效
func (s *AccountService) ValidateToken(accountID uint) (bool, error) {
	var account models.SteamAccount

	// 查询账户
	if err := s.db.First(&account, accountID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, errors.New("Steam账户不存在")
		}
		return false, errors.New("数据库查询错误")
	}

	// 检查令牌有效性
	return account.HasValidToken(), nil
}

// GetAccountWithTokenInfo 获取带有令牌信息的账户响应
func (s *AccountService) GetAccountWithTokenInfo(accountID uint, includeTokenContent bool) (*models.AccountTokenResponse, error) {
	var account models.SteamAccount

	// 查询账户
	if err := s.db.First(&account, accountID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("Steam账户不存在")
		}
		return nil, errors.New("数据库查询错误")
	}

	// 创建响应结构
	response := &models.AccountTokenResponse{
		ID:            account.ID,
		SteamID:       account.SteamID,
		Username:      account.Username,
		HasValidToken: account.HasValidToken(),
		LastLoginAt:   account.LastLoginAt,
		CreatedAt:     account.CreatedAt,
		UpdatedAt:     account.UpdatedAt,
	}

	// 根据需要包含令牌内容
	if includeTokenContent {
		tokenContent, err := account.GetTokenContent()
		if err == nil {
			response.TokenContent = tokenContent
		}
	}

	return response, nil
}
