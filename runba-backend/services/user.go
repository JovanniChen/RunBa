package services

import (
	"errors"
	"steam-backend/config"
	"steam-backend/models"
	"steam-backend/utils"
	"time"

	"gorm.io/gorm"
)

// UserService 用户服务结构体
// 处理用户相关的业务逻辑
type UserService struct {
	db *gorm.DB
}

// NewUserService 创建新的用户服务实例
func NewUserService() *UserService {
	return &UserService{
		db: config.GetDB(),
	}
}

// RegisterUser 注册新用户
// 包含完整的业务验证和处理逻辑
func (s *UserService) RegisterUser(req *models.RegisterRequest) (*models.LoginResponse, error) {
	// 业务验证：检查用户名是否已存在
	var existingUser models.User
	if err := s.db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return nil, errors.New("用户名已存在")
	}

	// 密码加密
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	// 创建用户对象
	user := models.User{
		Username: req.Username,
		Password: hashedPassword,
		Nickname: req.Nickname,
		Status:   models.UserStatusActive,
	}

	// 使用事务确保数据一致性
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 保存用户
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("用户创建失败")
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, errors.New("事务提交失败")
	}

	// 生成JWT令牌
	token, err := utils.GenerateToken(user.ID, user.Username, 24)
	if err != nil {
		return nil, errors.New("令牌生成失败")
	}

	// 构造返回数据
	userInfo := s.buildUserInfo(&user)
	return &models.LoginResponse{
		User:  userInfo,
		Token: token,
	}, nil
}

// LoginUser 用户登录
// 验证用户凭据并返回用户信息和令牌
func (s *UserService) LoginUser(req *models.LoginRequest) (*models.LoginResponse, error) {
	var user models.User

	// 查找用户
	if err := s.db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, errors.New("数据库查询错误")
	}

	// 检查用户状态
	if user.Status != models.UserStatusActive {
		return nil, errors.New("账户已被禁用，请联系管理员")
	}

	// 验证密码
	isValid, err := utils.CheckPassword(req.Password, user.Password)
	if err != nil || !isValid {
		return nil, errors.New("用户名或密码错误")
	}

	// 更新最后登录时间
	now := time.Now()
	user.LastLogin = &now
	s.db.Save(&user)

	// 生成JWT令牌
	token, err := utils.GenerateToken(user.ID, user.Username, 24)
	if err != nil {
		return nil, errors.New("令牌生成失败")
	}

	// 构造返回数据
	userInfo := s.buildUserInfo(&user)
	return &models.LoginResponse{
		User:  userInfo,
		Token: token,
	}, nil
}

// GetUserProfile 获取用户信息
func (s *UserService) GetUserProfile(userID uint) (*models.UserInfo, error) {
	var user models.User

	if err := s.db.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("用户不存在")
		}
		return nil, errors.New("数据库查询错误")
	}

	userInfo := s.buildUserInfo(&user)
	return &userInfo, nil
}

// UpdateUserProfile 更新用户信息
func (s *UserService) UpdateUserProfile(userID uint, req *models.UpdateUserRequest) (*models.UserInfo, error) {
	var user models.User

	// 查询用户
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	// 更新用户信息
	updates := make(map[string]interface{})
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}

	if len(updates) > 0 {
		if err := s.db.Model(&user).Updates(updates).Error; err != nil {
			return nil, errors.New("用户信息更新失败")
		}
	}

	// 重新查询更新后的用户信息
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, errors.New("获取更新后的用户信息失败")
	}

	userInfo := s.buildUserInfo(&user)
	return &userInfo, nil
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(userID uint, req *models.ChangePasswordRequest) error {
	var user models.User

	// 查询用户
	if err := s.db.First(&user, userID).Error; err != nil {
		return errors.New("用户不存在")
	}

	// 验证旧密码
	isValid, err := utils.CheckPassword(req.OldPassword, user.Password)
	if err != nil || !isValid {
		return errors.New("旧密码错误")
	}

	// 哈希新密码
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return errors.New("密码加密失败")
	}

	// 更新密码
	if err := s.db.Model(&user).Update("password", hashedPassword).Error; err != nil {
		return errors.New("密码更新失败")
	}

	return nil
}

// GetUsers 获取用户列表（带分页和搜索）
func (s *UserService) GetUsers(params *models.UserListRequest) (*models.UserListResponse, error) {
	var users []models.User
	var total int64

	// 构建查询
	query := s.db.Model(&models.User{})

	// 添加搜索条件
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("username LIKE ? OR nickname LIKE ?", keyword, keyword)
	}

	// 添加状态筛选
	if params.Status > 0 {
		query = query.Where("status = ?", params.Status)
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
	if err := query.Offset(offset).Limit(params.PageSize).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, errors.New("查询用户列表失败")
	}

	// 转换为用户信息格式
	userInfos := make([]models.UserInfo, len(users))
	for i, user := range users {
		userInfos[i] = s.buildUserInfo(&user)
	}

	// 计算总页数
	totalPage := int((total + int64(params.PageSize) - 1) / int64(params.PageSize))

	return &models.UserListResponse{
		Users:     userInfos,
		Page:      params.Page,
		PageSize:  params.PageSize,
		Total:     total,
		TotalPage: totalPage,
	}, nil
}

// GetUserByID 根据ID获取用户信息
func (s *UserService) GetUserByID(userID uint) (*models.UserInfo, error) {
	var user models.User

	if err := s.db.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("用户不存在")
		}
		return nil, errors.New("数据库查询错误")
	}

	userInfo := s.buildUserInfo(&user)
	return &userInfo, nil
}

// UpdateUser 更新用户信息（管理员功能）
func (s *UserService) UpdateUser(userID uint, req *models.UpdateUserRequest) (*models.UserInfo, error) {
	return s.UpdateUserProfile(userID, req)
}

// DeleteUser 删除用户（软删除）
func (s *UserService) DeleteUser(userID uint) error {
	var user models.User

	// 检查用户是否存在
	if err := s.db.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("用户不存在")
		}
		return errors.New("数据库查询错误")
	}

	// 执行软删除
	if err := s.db.Delete(&user).Error; err != nil {
		return errors.New("删除用户失败")
	}

	return nil
}

// UpdateUserStatus 更新用户状态
func (s *UserService) UpdateUserStatus(userID uint, status models.UserStatus) (*models.UserInfo, error) {
	var user models.User

	// 查询用户
	if err := s.db.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("用户不存在")
		}
		return nil, errors.New("数据库查询错误")
	}

	// 更新状态
	if err := s.db.Model(&user).Update("status", status).Error; err != nil {
		return nil, errors.New("用户状态更新失败")
	}

	// 重新查询更新后的用户信息
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, errors.New("获取更新后的用户信息失败")
	}

	userInfo := s.buildUserInfo(&user)
	return &userInfo, nil
}

// buildUserInfo 构建用户信息对象（不包含敏感信息）
// 私有方法，用于统一构建返回给前端的用户信息
func (s *UserService) buildUserInfo(user *models.User) models.UserInfo {
	return models.UserInfo{
		ID:        user.ID,
		Username:  user.Username,
		Nickname:  user.Nickname,
		Status:    user.Status,
		LastLogin: user.LastLogin,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// GetUserStats 获取用户统计信息
// 用于仪表盘显示
func (s *UserService) GetUserStats() (map[string]interface{}, error) {
	var stats = make(map[string]interface{})

	// 总用户数
	var totalUsers int64
	if err := s.db.Model(&models.User{}).Count(&totalUsers).Error; err != nil {
		return nil, errors.New("查询总用户数失败")
	}
	stats["total_users"] = totalUsers

	// 活跃用户数
	var activeUsers int64
	if err := s.db.Model(&models.User{}).Where("status = ?", models.UserStatusActive).Count(&activeUsers).Error; err != nil {
		return nil, errors.New("查询活跃用户数失败")
	}
	stats["active_users"] = activeUsers

	// 今日新增用户数
	today := time.Now().Format("2006-01-02")
	var todayNewUsers int64
	if err := s.db.Model(&models.User{}).Where("DATE(created_at) = ?", today).Count(&todayNewUsers).Error; err != nil {
		return nil, errors.New("查询今日新增用户数失败")
	}
	stats["today_new_users"] = todayNewUsers

	// 封禁用户数
	var bannedUsers int64
	if err := s.db.Model(&models.User{}).Where("status = ?", models.UserStatusBanned).Count(&bannedUsers).Error; err != nil {
		return nil, errors.New("查询封禁用户数失败")
	}
	stats["banned_users"] = bannedUsers

	return stats, nil
}
