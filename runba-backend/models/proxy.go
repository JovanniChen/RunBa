// 代理模型
package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// ProxyStatus 代理状态枚举
type ProxyStatus int8

const (
	ProxyStatusEnabled  ProxyStatus = 1 // 启用
	ProxyStatusDisabled ProxyStatus = 2 // 禁用
)

// ProxyProtocol 代理协议类型枚举
type ProxyProtocol int8

const (
	ProxyProtocolHTTP   ProxyProtocol = 1 // HTTP
	ProxyProtocolHTTPS  ProxyProtocol = 2 // HTTPS
	ProxyProtocolSOCKS5 ProxyProtocol = 3 // SOCKS5
)

// ProxyCategory 代理分类枚举
type ProxyCategory int8

const (
	ProxyCategoryLogin   ProxyCategory = 1 // 登录代理
	ProxyCategoryPayment ProxyCategory = 2 // 支付代理
)

// Proxy 代理记录模型
type Proxy struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Host      string         `json:"host" gorm:"unique;size:50;not null;comment:代理地址"`
	Port      int32          `json:"port" gorm:"not null;comment:端口"`
	Category  ProxyCategory  `json:"category" gorm:"not null;comment:代理分类(1=登录,2=支付)"`
	Protocol  ProxyProtocol  `json:"protocol" gorm:"not null;comment:协议类型(1=HTTP,2=HTTPS,3=SOCKS5)"`
	Status    ProxyStatus    `json:"status" gorm:"not null;default:1;comment:状态(1=启用,2=禁用,)"`
	Username  string         `json:"username" gorm:"comment:用户名"`
	Password  string         `json:"password" gorm:"comment:密码"`
	Remark    string         `json:"remark" gorm:"comment:备注信息"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 设置Proxy表名
func (Proxy) TableName() string {
	return "proxies"
}

// IsEnabled 检查代理是否启用
func (p *Proxy) IsEnabled() bool {
	return p.Status == ProxyStatusEnabled
}

// GetFullAddress 获取完整代理地址
func (p *Proxy) GetFullAddress() string {
	var protocol string
	switch p.Protocol {
	case ProxyProtocolHTTP:
		protocol = "http"
	case ProxyProtocolHTTPS:
		protocol = "https"
	case ProxyProtocolSOCKS5:
		protocol = "socks5"
	default:
		protocol = "unknown"
	}
	return fmt.Sprintf("%s://%s:%d", protocol, p.Host, p.Port)
}
