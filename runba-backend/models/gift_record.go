// Steam礼物记录模型
package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// GiftRecord 已赠送礼物记录模型
type GiftRecord struct {
	ID                uint           `json:"id" gorm:"primaryKey"`
	AccountID         uint           `json:"account_id" gorm:"not null;index;comment:Steam账户ID"`
	Code              string         `json:"code" gorm:"not null;index;comment:激活码"`
	Username          string         `json:"username" gorm:"not null;size:50"`
	SendSteamID       uint64         `json:"send_steam_id" gorm:"not null;index;comment:发送者SteamID"`
	TargetSteamID     uint64         `json:"target_steam_id" gorm:"not null;index;comment:接收者SteamID"`
	ReactionID        int32          `json:"reaction_id" gorm:"not null;comment:礼物反应ID"`
	GiftType          string         `json:"gift_type" gorm:"size:50;not null;comment:礼物类型标识"`
	GiftName          string         `json:"gift_name" gorm:"size:100;not null;comment:礼物名称"`
	PointsCost        int32          `json:"points_cost" gorm:"not null;comment:消耗的积分"`
	PointsTransferred int32          `json:"points_transferred" gorm:"not null;comment:转移的积分"`
	Status            string         `json:"status" gorm:"size:20;not null;default:'success';comment:赠送状态(success/failed)"`
	SentAt            time.Time      `json:"sent_at" gorm:"not null;comment:赠送时间"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Account SteamAccount `json:"account" gorm:"foreignKey:AccountID"`
}

// TableName 设置GiftRecord表名
func (GiftRecord) TableName() string {
	return "gift_records"
}

// BeforeCreate 创建前的钩子函数
func (g *GiftRecord) BeforeCreate(tx *gorm.DB) error {
	if g.SentAt.IsZero() {
		g.SentAt = time.Now()
	}
	return nil
}

// GetDisplayInfo 获取显示信息
func (g *GiftRecord) GetDisplayInfo() string {
	return fmt.Sprintf("账户ID:%d 向 SteamID:%d 赠送了 %s (消耗%d积分,转移%d积分)",
		g.AccountID, g.TargetSteamID, g.GiftName, g.PointsCost, g.PointsTransferred)
}
