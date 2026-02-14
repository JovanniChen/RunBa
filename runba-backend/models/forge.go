package models

import "time"

// Forge 锻刀所模型
type Forge struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:100;not null;comment:锻刀所名称"`
	ImageURL  string    `json:"image_url" gorm:"size:500;comment:图片地址"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 设置Forge表名
func (Forge) TableName() string {
	return "forges"
}
