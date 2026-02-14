// 配置模型
package models

import (
	"time"

	"gorm.io/gorm"
)

// Conf 配置记录模型
type Conf struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	TutorialUrl string         `json:"tutorial_url" gorm:"not null;comment:提货教程地址"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 设置Conf表名
func (Conf) TableName() string {
	return "confs"
}
