package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型结构体
//
// 定义用户表的数据结构，包含用户的基本信息
// 支持软删除、自动时间戳、唯一索引等特性
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`                         // 用户ID，主键，自动递增
	Username  string         `json:"username" gorm:"uniqueIndex;not null;size:50"` // 用户名，唯一索引，不能为空
	Password  string         `json:"-" gorm:"not null;size:255"`                   // 密码，不在JSON中显示，不能为空
	Nickname  string         `json:"nickname" gorm:"size:50"`                      // 昵称，可选字段
	Status    UserStatus     `json:"status" gorm:"default:1"`                      // 用户状态，默认为激活状态
	LastLogin *time.Time     `json:"last_login"`                                   // 最后登录时间，可为空
	CreatedAt time.Time      `json:"created_at"`                                   // 创建时间，自动管理
	UpdatedAt time.Time      `json:"updated_at"`                                   // 更新时间，自动管理
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`                               // 软删除时间，支持软删除
}

// UserStatus 用户状态枚举
//
// 定义用户的各种状态，用于用户管理和权限控制
type UserStatus int

const (
	UserStatusInactive UserStatus = 0 // 未激活状态：用户注册但未激活
	UserStatusActive   UserStatus = 1 // 激活状态：正常使用
	UserStatusBanned   UserStatus = 2 // 封禁状态：被管理员封禁
)

// String 实现Stringer接口，返回用户状态的字符串表示
//
// 返回值：
// - string: 用户状态的字符串表示
//
// 状态映射：
// - 0 -> "inactive" (未激活)
// - 1 -> "active" (激活)
// - 2 -> "banned" (封禁)
// - 其他 -> "unknown" (未知状态)
func (s UserStatus) String() string {
	switch s {
	case UserStatusInactive:
		return "inactive" // 未激活
	case UserStatusActive:
		return "active" // 激活
	case UserStatusBanned:
		return "banned" // 封禁
	default:
		return "unknown" // 未知状态
	}
}

// TableName 设置User模型对应的表名
//
// 返回值：
// - string: 表名 "users"
//
// 说明：
// - GORM会自动将结构体名转换为复数形式作为表名
// - 这里显式指定表名以确保一致性
func (User) TableName() string {
	return "users"
}
