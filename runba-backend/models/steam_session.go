// Steam会话持久化模型
package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// SteamSessionRecord Steam会话记录 - 用于数据库持久化
type SteamSessionRecord struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	AccountID     uint      `json:"account_id" gorm:"uniqueIndex;not null;comment:Steam账户ID"`
	SteamID       uint64    `json:"steam_id" gorm:"not null;comment:Steam用户ID"`
	Username      string    `json:"username" gorm:"size:100;not null;comment:Steam用户名"`
	Nickname      string    `json:"nickname" gorm:"size:100;comment:Steam昵称"`
	SessionData   string    `json:"session_data" gorm:"type:longtext;comment:会话数据JSON格式"`
	AccessToken   string    `json:"access_token" gorm:"type:text;comment:访问令牌"`
	LoginTime     time.Time `json:"login_time" gorm:"not null;comment:登录时间"`
	LastUsed      time.Time `json:"last_used" gorm:"not null;comment:最后使用时间"`
	ExpiresAt     time.Time `json:"expires_at" gorm:"not null;index;comment:过期时间"`
	Status        string    `json:"status" gorm:"size:20;default:'active';comment:会话状态"`
	ClientVersion string    `json:"client_version" gorm:"size:50;comment:客户端版本"`
	IPAddress     string    `json:"ip_address" gorm:"size:45;comment:登录IP地址"`
	UserAgent     string    `json:"user_agent" gorm:"size:500;comment:用户代理"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// TableName 指定表名
func (SteamSessionRecord) TableName() string {
	return "steam_session_records"
}

// SessionMetadata 会话元数据结构
type SessionMetadata struct {
	SteamID       uint64        `json:"steam_id"`
	Username      string        `json:"username"`
	Nickname      string        `json:"nickname"`
	LoginTime     time.Time     `json:"login_time"`
	LastUsed      time.Time     `json:"last_used"`
	SessionCookie string        `json:"session_cookie"`
	AccessToken   string        `json:"access_token"`
	TokenContent  *TokenContent `json:"token_content,omitempty"`
	Status        string        `json:"status"`
	ClientInfo    *ClientInfo   `json:"client_info,omitempty"`
}

// ClientInfo 客户端信息
type ClientInfo struct {
	Version   string `json:"version"`
	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
	Platform  string `json:"platform"`
}

// GetSessionMetadata 获取会话元数据
func (s *SteamSessionRecord) GetSessionMetadata() (*SessionMetadata, error) {
	if s.SessionData == "" {
		return nil, nil
	}

	var metadata SessionMetadata
	if err := json.Unmarshal([]byte(s.SessionData), &metadata); err != nil {
		return nil, err
	}

	// 填充基础字段
	metadata.SteamID = s.SteamID
	metadata.Username = s.Username
	metadata.Nickname = s.Nickname
	metadata.LoginTime = s.LoginTime
	metadata.LastUsed = s.LastUsed
	metadata.AccessToken = s.AccessToken
	metadata.Status = s.Status

	return &metadata, nil
}

// SetSessionMetadata 设置会话元数据
func (s *SteamSessionRecord) SetSessionMetadata(metadata *SessionMetadata) error {
	if metadata == nil {
		s.SessionData = ""
		return nil
	}

	// 更新基础字段
	s.SteamID = metadata.SteamID
	s.Username = metadata.Username
	s.Nickname = metadata.Nickname
	s.LoginTime = metadata.LoginTime
	s.LastUsed = metadata.LastUsed
	s.AccessToken = metadata.AccessToken
	s.Status = metadata.Status

	// 提取客户端信息
	if metadata.ClientInfo != nil {
		s.ClientVersion = metadata.ClientInfo.Version
		s.IPAddress = metadata.ClientInfo.IPAddress
		s.UserAgent = metadata.ClientInfo.UserAgent
	}

	// 序列化完整元数据
	sessionBytes, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	s.SessionData = string(sessionBytes)
	return nil
}

// IsExpired 检查会话是否已过期
func (s *SteamSessionRecord) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// IsActive 检查会话是否活跃
func (s *SteamSessionRecord) IsActive() bool {
	return s.Status == "active" && !s.IsExpired()
}

// UpdateLastUsed 更新最后使用时间
func (s *SteamSessionRecord) UpdateLastUsed() {
	s.LastUsed = time.Now()
}

// SetExpiry 设置过期时间
func (s *SteamSessionRecord) SetExpiry(duration time.Duration) {
	s.ExpiresAt = time.Now().Add(duration)
}

// Deactivate 停用会话
func (s *SteamSessionRecord) Deactivate() {
	s.Status = "inactive"
}

// BeforeCreate GORM钩子 - 创建前
func (s *SteamSessionRecord) BeforeCreate(tx *gorm.DB) error {
	if s.LoginTime.IsZero() {
		s.LoginTime = time.Now()
	}
	if s.LastUsed.IsZero() {
		s.LastUsed = time.Now()
	}
	if s.ExpiresAt.IsZero() {
		s.SetExpiry(2 * time.Hour) // 默认2小时过期
	}
	if s.Status == "" {
		s.Status = "active"
	}
	return nil
}

// BeforeUpdate GORM钩子 - 更新前
func (s *SteamSessionRecord) BeforeUpdate(tx *gorm.DB) error {
	// 如果更新了LastUsed，自动延长过期时间
	if s.LastUsed.After(s.ExpiresAt.Add(-30 * time.Minute)) {
		s.SetExpiry(2 * time.Hour)
	}
	return nil
}
