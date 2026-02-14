package models

import "time"

// Swordsmith 刀匠模型
type Swordsmith struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:100;not null;comment:刀匠名称"`
	ImageURL  string    `json:"image_url" gorm:"size:500;comment:图片地址"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 设置Swordsmith表名
func (Swordsmith) TableName() string {
	return "swordsmiths"
}
