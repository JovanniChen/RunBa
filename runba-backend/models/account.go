// Steam用户模型
package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// Session Steam会话信息
type Session struct {
	SteamID      uint64  `json:"SteamID"`
	AccessToken  string  `json:"AccessToken"`
	RefreshToken string  `json:"RefreshToken"`
	SessionID    *string `json:"SessionID"`
}

// TokenContent Steam令牌内容结构体
type TokenContent struct {
	SharedSecret   string  `json:"shared_secret"`   // Steam Guard共享密钥
	SerialNumber   string  `json:"serial_number"`   // 序列号
	RevocationCode string  `json:"revocation_code"` // 撤销代码
	URI            string  `json:"uri"`             // TOTP URI
	ServerTime     int64   `json:"server_time"`     // 服务器时间
	AccountName    string  `json:"account_name"`    // 账户名
	TokenGID       string  `json:"token_gid"`       // 令牌GID
	IdentitySecret string  `json:"identity_secret"` // 身份密钥
	Secret1        string  `json:"secret_1"`        // 密钥1
	Status         int     `json:"status"`          // 状态
	DeviceID       string  `json:"device_id"`       // 设备ID
	FullyEnrolled  bool    `json:"fully_enrolled"`  // 是否完全注册
	Session        Session `json:"Session"`         // 会话信息
}

// SteamAccount Steam账户模型 - 独立实体，不关联User
type SteamAccount struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	SteamID      uint64         `json:"steam_id" gorm:"unique;not null;comment:Steam用户ID"`
	Username     string         `json:"username" gorm:"size:100;not null;comment:Steam用户名"`
	Password     string         `json:"-" gorm:"size:255;comment:Steam密码"`
	Nickname     string         `json:"nickname" gorm:"size:100;comment:Steam昵称"`
	CountryCode  string         `json:"country_code" gorm:"size:10;comment:国家代码"`
	Points       int32          `json:"points" gorm:"default:0;comment:Steam积分"`
	TokenContent string         `json:"token_content" gorm:"type:text;comment:Steam令牌内容JSON格式"`
	Status       uint64         `json:"status" gorm:"default:2;comment:账户状态"`
	LastLoginAt  *time.Time     `json:"last_login_at" gorm:"comment:最后登录时间"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
	AccountNote  string         `json:"note" gorm:"type:varchar(500);default:'';comment:备注信息"`
}

// SteamAccount 状态质数常量 - 使用质数支持复合状态
const (
	StatusNormalOfAccount       uint64 = 2  // 正常 (质数: 2)
	StatusBannedOfAccount       uint64 = 3  // 禁用 (质数: 3)
	StatusIsLoginingOfAccount   uint64 = 5  // 登录中 (质数: 5)
	StatusFailedLoginOfAccount  uint64 = 7  // 登录失败 (质数: 7)
	StatusSuccessLoginOfAccount uint64 = 11 // 登录成功 (质数: 11)
)

// StatusManager 状态管理器
type StatusManager struct{}

// AddStatus 添加状态
func (sm *StatusManager) AddStatus(currentStatus uint64, newStatus uint64) uint64 {
	if currentStatus%newStatus == 0 {
		return currentStatus // 已存在该状态
	}
	return currentStatus * newStatus
}

// RemoveStatus 移除状态
func (sm *StatusManager) RemoveStatus(currentStatus uint64, statusToRemove uint64) uint64 {
	if currentStatus%statusToRemove != 0 {
		return currentStatus // 不存在该状态
	}
	return currentStatus / statusToRemove
}

// HasStatus 检查是否包含状态
func (sm *StatusManager) HasStatus(currentStatus uint64, checkStatus uint64) bool {
	return currentStatus > 0 && currentStatus%checkStatus == 0
}

// GetAllStatuses 获取所有状态
func (sm *StatusManager) GetAllStatuses(currentStatus uint64) []uint64 {
	var statuses []uint64
	allPossibleStatuses := []uint64{
		StatusNormalOfAccount, StatusBannedOfAccount, StatusIsLoginingOfAccount,
		StatusFailedLoginOfAccount, StatusSuccessLoginOfAccount,
	}

	for _, status := range allPossibleStatuses {
		if currentStatus > 0 && currentStatus%status == 0 {
			statuses = append(statuses, status)
		}
	}
	return statuses
}

// 全局状态管理器实例
var AccountStatusManager = &StatusManager{}

// TableName 设置SteamAccount表名
func (SteamAccount) TableName() string {
	return "steam_accounts"
}

// UpdateLastLogin 更新最后登录时间
func (u *SteamAccount) UpdateLastLogin() {
	now := time.Now()
	u.LastLoginAt = &now
}

// GetDisplayName 获取显示名称
func (u *SteamAccount) GetDisplayName() string {
	if u.Nickname != "" {
		return u.Nickname
	}
	return u.Username
}

// IsActive 检查账户是否激活
func (u *SteamAccount) IsActive() bool {
	return AccountStatusManager.HasStatus(u.Status, StatusNormalOfAccount)
}

// AddPoints 增加积分
func (u *SteamAccount) AddPoints(points int32) {
	u.Points += points
}

// DeductPoints 扣除积分
func (u *SteamAccount) DeductPoints(points int32) bool {
	if u.Points >= points {
		u.Points -= points
		return true
	}
	return false
}

// SetTokenContent 设置令牌内容
func (u *SteamAccount) SetTokenContent(tokenContent *TokenContent) error {
	if tokenContent == nil {
		u.TokenContent = ""
		return nil
	}

	jsonData, err := json.Marshal(tokenContent)
	if err != nil {
		return err
	}

	u.TokenContent = string(jsonData)
	return nil
}

// GetTokenContent 获取令牌内容
func (u *SteamAccount) GetTokenContent() (*TokenContent, error) {
	if u.TokenContent == "" {
		return nil, nil
	}

	var tokenContent TokenContent
	err := json.Unmarshal([]byte(u.TokenContent), &tokenContent)
	if err != nil {
		return nil, err
	}

	return &tokenContent, nil
}

// HasValidToken 检查是否有有效的令牌
func (u *SteamAccount) HasValidToken() bool {
	if u.TokenContent == "" {
		return false
	}

	tokenContent, err := u.GetTokenContent()
	if err != nil {
		return false
	}

	// 检查必要字段是否存在
	return tokenContent != nil &&
		tokenContent.SharedSecret != "" &&
		tokenContent.IdentitySecret != "" &&
		tokenContent.AccountName != ""
}

// GetSharedSecret 获取共享密钥
func (u *SteamAccount) GetSharedSecret() string {
	tokenContent, err := u.GetTokenContent()
	if err != nil || tokenContent == nil {
		return ""
	}
	return tokenContent.SharedSecret
}

// GetIdentitySecret 获取身份密钥
func (u *SteamAccount) GetIdentitySecret() string {
	tokenContent, err := u.GetTokenContent()
	if err != nil || tokenContent == nil {
		return ""
	}
	return tokenContent.IdentitySecret
}

// GetAccessToken 获取访问令牌
func (u *SteamAccount) GetAccessToken() string {
	tokenContent, err := u.GetTokenContent()
	if err != nil || tokenContent == nil {
		return ""
	}
	return tokenContent.Session.AccessToken
}

// GetRefreshToken 获取刷新令牌
func (u *SteamAccount) GetRefreshToken() string {
	tokenContent, err := u.GetTokenContent()
	if err != nil || tokenContent == nil {
		return ""
	}
	return tokenContent.Session.RefreshToken
}

// UpdateTokenSession 更新令牌会话信息
func (u *SteamAccount) UpdateTokenSession(steamID uint64, accessToken, refreshToken string) error {
	tokenContent, err := u.GetTokenContent()
	if err != nil {
		return err
	}

	if tokenContent == nil {
		// 如果没有令牌内容，创建一个基本的
		tokenContent = &TokenContent{}
	}

	// 更新会话信息
	tokenContent.Session.SteamID = steamID
	tokenContent.Session.AccessToken = accessToken
	tokenContent.Session.RefreshToken = refreshToken

	return u.SetTokenContent(tokenContent)
}

// 状态管理便捷方法

// AddStatus 添加状态
func (u *SteamAccount) AddStatus(status uint64) {
	u.Status = AccountStatusManager.AddStatus(u.Status, status)
}

// RemoveStatus 移除状态
func (u *SteamAccount) RemoveStatus(status uint64) {
	u.Status = AccountStatusManager.RemoveStatus(u.Status, status)
}

// HasStatus 检查是否包含指定状态
func (u *SteamAccount) HasStatus(status uint64) bool {
	return AccountStatusManager.HasStatus(u.Status, status)
}

// GetAllStatuses 获取所有状态
func (u *SteamAccount) GetAllStatuses() []uint64 {
	return AccountStatusManager.GetAllStatuses(u.Status)
}

// IsBanned 检查是否被禁用
func (u *SteamAccount) IsBanned() bool {
	return u.HasStatus(StatusBannedOfAccount)
}

// IsVerifying 检查是否登录中
func (u *SteamAccount) IsLogining() bool {
	return u.HasStatus(StatusIsLoginingOfAccount)
}

// SetNormalStatus 设置为正常状态
func (u *SteamAccount) SetNormalStatus() {
	u.Status = StatusNormalOfAccount
}

// SetVerifyingStatus 设置为登录中状态
func (u *SteamAccount) SetVerifyingStatus() {
	u.SetNormalStatus()
	u.AddStatus(StatusIsLoginingOfAccount)
}

// SetFailedStatus 设置为验证失败状态
func (u *SteamAccount) SetFailedStatus() {
	u.Status = StatusFailedLoginOfAccount
}

// SetBannedStatus 设置为禁用状态
func (u *SteamAccount) SetBannedStatus() {
	u.Status = StatusBannedOfAccount
}
