// Steam礼物记录模型
package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// ActivationCode 状态常量
const (
	StatusNormalOfActivationCode int8 = 1 // 正常
	StatusUsedOfActivationCode   int8 = 2 // 已使用
	StatusDropOfActivationCode   int8 = 3 // 已废弃
)

// ActivationCode 激活码模型
type ActivationCode struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null;comment:激活码名称"`
	Code      string         `json:"code" gorm:"not null;comment:激活码"`
	Points    int32          `json:"points" gorm:"not null;comment:可激活点数"`
	Status    int8           `json:"status" gorm:"default:1;comment:记录状态[正常和已使用]"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 设置ActivationCode表名
func (ActivationCode) TableName() string {
	return "activation_codes"
}

// GetDisplayInfo 获取显示信息
func (g *ActivationCode) GetDisplayInfo() string {
	return fmt.Sprintf("激活码:%s 可激活点数:%d", g.Code, g.Points)
}

// GetStatusText 获取状态文本描述
func (r *ActivationCode) GetStatusText() string {
	switch r.Status {
	case StatusNormalOfActivationCode:
		return "正常"
	case StatusUsedOfActivationCode:
		return "已使用"
	case StatusDropOfActivationCode:
		return "已废弃"
	default:
		return "未知状态"
	}
}

// IsProcessable 检查记录是否可以处理
func (r *ActivationCode) IsProcessable() bool {
	return r.Status == StatusNormalOfActivationCode
}
