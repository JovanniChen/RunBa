// Steam用户模型
package models

import (
	"time"

	"gorm.io/gorm"
)

// ExchangeRecord 状态常量
const (
	StatusPendingOfExchangeRecord int8 = 1 // 待兑换
	StatusSuccessOfExchangeRecord int8 = 2 // 已兑换
	StatusFailedOfExchangeRecord  int8 = 3 // 兑换失败
	StatusRetryOfExchangeRecord   int8 = 4 // 待重试
)

// ExchangeRecord 点数兑换记录模型
type ExchangeRecord struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Code      string         `json:"code" gorm:"unique;size:100;not null;comment:激活码"`
	SteamID   uint64         `json:"steam_id" gorm:"not null;index;comment:申请兑换点数者的SteamID"`
	Points    int32          `json:"points" gorm:"default:0;comment:要兑换的Steam积分"`
	Status    int8           `json:"status" gorm:"default:1;comment:记录状态[待兑换和已兑换]"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系 - 一个兑换记录对应多个礼物记录
	GiftRecords []GiftRecord `json:"gift_records,omitempty" gorm:"foreignKey:Code;references:Code"`
}

// TableName 设置ExchangeRecord表名
func (ExchangeRecord) TableName() string {
	return "exchange_records"
}

// GetStatusText 获取状态文本描述
func (r *ExchangeRecord) GetStatusText() string {
	switch r.Status {
	case StatusPendingOfExchangeRecord:
		return "待兑换"
	case StatusSuccessOfExchangeRecord:
		return "已兑换"
	case StatusFailedOfExchangeRecord:
		return "兑换失败"
	case StatusRetryOfExchangeRecord:
		return "待重试"
	default:
		return "未知状态"
	}
}

// IsProcessable 检查记录是否可以处理
func (r *ExchangeRecord) IsProcessable() bool {
	return r.Status == StatusPendingOfExchangeRecord
}
