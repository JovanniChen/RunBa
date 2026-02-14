package models

import "time"

// RegisterRequest 用户注册请求结构体
// 定义用户注册时需要提供的数据
type RegisterRequest struct {
	Username        string `json:"username" form:"username" binding:"required,min=3,max=20"`                     // 用户名，必填，长度3-20
	Password        string `json:"password" form:"password" binding:"required,min=6,max=100"`                    // 密码，必填，长度6-100
	ConfirmPassword string `json:"confirm_password" form:"confirm_password" binding:"required,eqfield=Password"` // 确认密码，必填，需与密码相同
	Nickname        string `json:"nickname" form:"nickname" binding:"max=50"`                                    // 昵称，可选，最大长度50
}

// LoginRequest 用户登录请求结构体
type LoginRequest struct {
	Username string `json:"account" form:"account" binding:"required"`   // 用户名，必填
	Password string `json:"password" form:"password" binding:"required"` // 密码，必填
}

// LoginResponse 登录响应结构体
// 返回登录成功后的用户信息和token
type LoginResponse struct {
	User  UserInfo `json:"user"`  // 用户基本信息
	Token string   `json:"token"` // JWT访问令牌
}

// UserInfo 用户基本信息结构体
// 用于API返回，不包含敏感信息如密码
type UserInfo struct {
	ID        uint       `json:"id"`         // 用户ID
	Username  string     `json:"username"`   // 用户名
	Nickname  string     `json:"nickname"`   // 昵称
	Status    UserStatus `json:"status"`     // 用户状态
	LastLogin *time.Time `json:"last_login"` // 最后登录时间
	CreatedAt time.Time  `json:"created_at"` // 创建时间
	UpdatedAt time.Time  `json:"updated_at"` // 更新时间
}

type GetUserRequest struct {
	ID uint `json:"id" form:"id" binding:"required,min=1"`
}

// UpdateUserRequest 更新用户信息请求结构体
// 用于用户更新个人信息
type UpdateUserRequest struct {
	ID       uint   `json:"id" form:"id" binding:"required,min=1"`
	Nickname string `json:"nickname" form:"nickname" binding:"omitempty,max=50"` // 昵称，可选，最大长度50
}

// ChangePasswordRequest 修改密码请求结构体
type ChangePasswordRequest struct {
	OldPassword        string `json:"old_password" form:"old_password" binding:"required"`                                     // 旧密码，必填
	NewPassword        string `json:"new_password" form:"new_password" binding:"required,min=6,max=100"`                       // 新密码，必填，长度6-100
	ConfirmNewPassword string `json:"confirm_new_password" form:"confirm_new_password" binding:"required,eqfield=NewPassword"` // 确认新密码，必填，需与新密码相同
}

type DeleteUserRequest struct {
	ID uint `json:"id" form:"id" binding:"required,min=1"`
}

// UpdateUserStatusRequest 更新用户状态请求结构体
type UpdateUserStatusRequest struct {
	ID     uint       `json:"id" form:"id" binding:"required,min=1"`
	Status UserStatus `json:"status" form:"status" binding:"required"` // 用户状态，必填
}

// UserListRequest 用户列表查询请求结构体
// 支持分页和搜索功能
type UserListRequest struct {
	Page     int    `form:"page,default=1" binding:"min=1"`               // 页码，默认为1，最小值1
	PageSize int    `form:"page_size,default=10" binding:"min=1,max=100"` // 每页大小，默认10，范围1-100
	Keyword  string `form:"keyword"`                                      // 搜索关键词，可搜索用户名、昵称
	Status   int    `form:"status"`                                       // 用户状态筛选，0-全部，1-激活，2-禁用
}

// UserListResponse 用户列表响应结构体
type UserListResponse struct {
	Users     []UserInfo `json:"rows"`      // 用户列表
	Page      int        `json:"page"`      // 当前页码
	PageSize  int        `json:"page_size"` // 每页大小
	Total     int64      `json:"total"`     // 总记录数
	TotalPage int        `json:"totalpage"` // 总页数
}

// CreateAccountRequest 创建Steam账户请求结构体
type CreateAccountRequest struct {
	Username     string        `json:"username" form:"username" binding:"required,min=3,max=100"` // Steam用户名，必填，长度3-100
	Password     string        `json:"password" form:"password" binding:"required,min=6,max=255"` // Steam密码，必填，长度6-255
	TokenContent *TokenContent `json:"token_content" form:"token_content" binding:"required"`     // 令牌内容，必填
	AccountNote  string        `json:"note" form:"note" binding:"omitempty,max=500"`              // 备注信息，可选，最大500字符
}

// UpdateAccountRequest 更新Steam账户信息请求结构体
type UpdateAccountRequest struct {
	ID           uint          `json:"id" form:"id" binding:"required,min=1"`
	Username     string        `json:"username" form:"username" binding:"omitempty,min=3,max=100"`  // Steam用户名，可选，长度3-100
	Password     string        `json:"password" form:"password" binding:"omitempty,min=6,max=255"`  // Steam密码，可选，长度6-255
	Nickname     string        `json:"nickname" form:"nickname" binding:"omitempty,max=100"`        // Steam昵称，可选，最大长度100
	CountryCode  string        `json:"country_code" form:"country_code" binding:"omitempty,max=10"` // 国家代码，可选，最大长度10
	TokenContent *TokenContent `json:"token_content,omitempty" form:"token_content,omitempty"`      // 令牌内容，可选
}

// AccountListRequest Steam账户列表查询请求结构体
// 支持分页和搜索功能
type AccountListRequest struct {
	Page     int    `form:"page,default=1" binding:"min=1"`               // 页码，默认为1，最小值1
	PageSize int    `form:"page_size,default=10" binding:"min=1,max=100"` // 每页大小，默认10，范围1-100
	Keyword  string `form:"keyword"`                                      // 搜索关键词，可搜索用户名、昵称
	Status   int    `form:"status"`                                       // 账户状态筛选，1-正常，2-停用，nil-全部
}

// AccountListResponse Steam账户列表响应结构体
type AccountListResponse struct {
	Accounts  []SteamAccount `json:"rows"`      // Steam账户列表
	Page      int            `json:"page"`      // 当前页码
	PageSize  int            `json:"page_size"` // 每页大小
	Total     int64          `json:"total"`     // 总记录数
	TotalPage int            `json:"totalpage"` // 总页数
}

type GetAccountRequest struct {
	ID uint `json:"id" form:"id" binding:"required,min=1"`
}

// UpdateAccountStatusRequest 更新Steam账户状态请求结构体
type UpdateAccountStatusRequest struct {
	ID     uint   `json:"id" form:"id" binding:"required,min=1"`
	Status uint64 `json:"status" form:"status" binding:"required"` // 账户状态
}

// AddAccountStatusRequest 添加账户状态请求结构体
type AddAccountStatusRequest struct {
	ID     uint   `json:"id" form:"id" binding:"required,min=1"`
	Status uint64 `json:"status" binding:"required"` // 要添加的状态
}

// RemoveAccountStatusRequest 移除账户状态请求结构体
type RemoveAccountStatusRequest struct {
	ID     uint   `json:"id" form:"id" binding:"required,min=1"`
	Status uint64 `json:"status" binding:"required"` // 要移除的状态
}

type DeleteAccountRequest struct {
	ID uint `json:"id" form:"id" binding:"required,min=1"`
}

// CreateCertificateRequest 创建证书请求结构体
type CreateCertificateRequest struct {
	Number           string `json:"number" form:"number" binding:"required"` // 证书编号
	OverallLength    string `json:"overall_length" form:"overall_length"`
	BladeLength      string `json:"blade_length" form:"blade_length"`
	HandleLength     string `json:"handle_length" form:"handle_length"`
	BladeMaterial    string `json:"blade_material" form:"blade_material"`
	BladeThickness   string `json:"blade_thickness" form:"blade_thickness"`
	KissanLength     string `json:"kissan_length" form:"kissan_length"`
	Hamon            string `json:"hamon" form:"hamon"`
	Mekugi           int8   `json:"mekugi" form:"mekugi"`
	Sakihaba         string `json:"sakihaba" form:"sakihaba"`
	Motohaba         string `json:"motohaba" form:"motohaba"`
	SayaMaterial     string `json:"saya_material" form:"saya_material"`
	TsubaMaterial    string `json:"tsuba_material" form:"tsuba_material"`
	HabakiMaterial   string `json:"habaki_material" form:"habaki_material"`
	ItoSageoMaterial string `json:"ito_sageo_material" form:"ito_sageo_material"`
	ForgeID          uint   `json:"forge_id" form:"forge_id" binding:"required,min=1"`
	SwordsmithID     uint   `json:"swordsmith_id" form:"swordsmith_id" binding:"required,min=1"`
	DateCompleted    string `json:"date_completed" form:"date_completed"`
	No               string `json:"no" form:"no"`
}

// UpdateCertificateRequest 更新证书请求结构体
type UpdateCertificateRequest struct {
	ID               uint   `json:"id" form:"id" binding:"required,min=1"`
	Number           string `json:"number" form:"number" binding:"omitempty"`
	OverallLength    string `json:"overall_length" form:"overall_length" binding:"omitempty"`
	BladeLength      string `json:"blade_length" form:"blade_length" binding:"omitempty"`
	HandleLength     string `json:"handle_length" form:"handle_length" binding:"omitempty"`
	BladeMaterial    string `json:"blade_material" form:"blade_material" binding:"omitempty"`
	BladeThickness   string `json:"blade_thickness" form:"blade_thickness" binding:"omitempty"`
	KissanLength     string `json:"kissan_length" form:"kissan_length" binding:"omitempty"`
	Hamon            string `json:"hamon" form:"hamon" binding:"omitempty"`
	Mekugi           *int8  `json:"mekugi" form:"mekugi"`
	Sakihaba         string `json:"sakihaba" form:"sakihaba" binding:"omitempty"`
	Motohaba         string `json:"motohaba" form:"motohaba" binding:"omitempty"`
	SayaMaterial     string `json:"saya_material" form:"saya_material" binding:"omitempty"`
	TsubaMaterial    string `json:"tsuba_material" form:"tsuba_material" binding:"omitempty"`
	HabakiMaterial   string `json:"habaki_material" form:"habaki_material" binding:"omitempty"`
	ItoSageoMaterial string `json:"ito_sageo_material" form:"ito_sageo_material" binding:"omitempty"`
	ForgeID          *uint  `json:"forge_id" form:"forge_id"`
	SwordsmithID     *uint  `json:"swordsmith_id" form:"swordsmith_id"`
	DateCompleted    string `json:"date_completed" form:"date_completed" binding:"omitempty"`
	No               string `json:"no" form:"no" binding:"omitempty"`
}

// CertificateListRequest 证书列表请求结构体
type CertificateListRequest struct {
	Page     int    `form:"page,default=1" binding:"min=1"`               // 页码
	PageSize int    `form:"page_size,default=10" binding:"min=1,max=100"` // 每页大小
	Keyword  string `form:"keyword"`                                      // 关键词
}

// CertificateListResponse 证书列表响应结构体
type CertificateListResponse struct {
	Certificates []Certificate `json:"rows"`
	Page         int           `json:"page"`
	PageSize     int           `json:"page_size"`
	Total        int64         `json:"total"`
	TotalPage    int           `json:"totalpage"`
}

type GetCertificateRequest struct {
	ID uint `json:"id" form:"id" binding:"required,min=1"`
}

type GetCertificateByNoRequest struct {
	No string `json:"no" form:"no" binding:"required"`
}

type DeleteCertificateRequest struct {
	ID uint `json:"id" form:"id" binding:"required,min=1"`
}

// CreateForgeRequest 创建锻刀所请求
type CreateForgeRequest struct {
	Name     string `json:"name" form:"name" binding:"required"`
	ImageURL string `json:"image_url" form:"image_url"`
}

// UpdateForgeRequest 更新锻刀所请求
type UpdateForgeRequest struct {
	ID       uint   `json:"id" form:"id" binding:"required,min=1"`
	Name     string `json:"name" form:"name" binding:"omitempty"`
	ImageURL string `json:"image_url" form:"image_url" binding:"omitempty"`
}

// ForgeListRequest 锻刀所列表请求
type ForgeListRequest struct {
	Page     int    `form:"page,default=1" binding:"min=1"`
	PageSize int    `form:"page_size,default=10" binding:"min=1,max=100"`
	Keyword  string `form:"keyword"`
}

// ForgeListResponse 锻刀所列表响应
type ForgeListResponse struct {
	Forges    []Forge `json:"rows"`
	Page      int     `json:"page"`
	PageSize  int     `json:"page_size"`
	Total     int64   `json:"total"`
	TotalPage int     `json:"totalpage"`
}

type GetForgeRequest struct {
	ID uint `json:"id" form:"id" binding:"required,min=1"`
}

type DeleteForgeRequest struct {
	ID uint `json:"id" form:"id" binding:"required,min=1"`
}

// CreateSwordsmithRequest 创建刀匠请求
type CreateSwordsmithRequest struct {
	Name     string `json:"name" form:"name" binding:"required"`
	ImageURL string `json:"image_url" form:"image_url"`
}

// UpdateSwordsmithRequest 更新刀匠请求
type UpdateSwordsmithRequest struct {
	ID       uint   `json:"id" form:"id" binding:"required,min=1"`
	Name     string `json:"name" form:"name" binding:"omitempty"`
	ImageURL string `json:"image_url" form:"image_url" binding:"omitempty"`
}

// SwordsmithListRequest 刀匠列表请求
type SwordsmithListRequest struct {
	Page     int    `form:"page,default=1" binding:"min=1"`
	PageSize int    `form:"page_size,default=10" binding:"min=1,max=100"`
	Keyword  string `form:"keyword"`
}

// SwordsmithListResponse 刀匠列表响应
type SwordsmithListResponse struct {
	Swordsmiths []Swordsmith `json:"rows"`
	Page        int          `json:"page"`
	PageSize    int          `json:"page_size"`
	Total       int64        `json:"total"`
	TotalPage   int          `json:"totalpage"`
}

type GetSwordsmithRequest struct {
	ID uint `json:"id" form:"id" binding:"required,min=1"`
}

type DeleteSwordsmithRequest struct {
	ID uint `json:"id" form:"id" binding:"required,min=1"`
}

// PointsOperationRequest 积分操作请求结构体
type PointsOperationRequest struct {
	Points int32 `json:"points" binding:"required,min=1"` // 积分数量，必填，最小值1
}

// UpdateTokenContentRequest 更新令牌内容请求结构体
type UpdateTokenContentRequest struct {
	TokenContent *TokenContent `json:"token_content" binding:"required"` // 令牌内容，必填
}

// AccountTokenResponse 账户令牌响应结构体
type AccountTokenResponse struct {
	ID            uint          `json:"id"`                      // 账户ID
	SteamID       uint64        `json:"steam_id"`                // Steam用户ID
	Username      string        `json:"username"`                // Steam用户名
	HasValidToken bool          `json:"has_valid_token"`         // 是否有有效令牌
	TokenContent  *TokenContent `json:"token_content,omitempty"` // 令牌内容（只在需要时返回）
	LastLoginAt   *time.Time    `json:"last_login_at"`           // 最后登录时间
	CreatedAt     time.Time     `json:"created_at"`              // 创建时间
	UpdatedAt     time.Time     `json:"updated_at"`              // 更新时间
}

// TokenValidationRequest 令牌验证请求结构体
type TokenValidationRequest struct {
	AccountID uint `json:"account_id" binding:"required"` // 账户ID，必填
}

// CreateExchangeRecordRequest 创建兑换记录请求结构体
type CreateExchangeRecordRequest struct {
	Code             string `json:"code" form:"code" binding:"required"`                             // 激活码，必填
	FriendCode       string `json:"friend_code" form:"friend_code" binding:"required"`               // 好友码
	RepeatFriendCode string `json:"repeat_friend_code" form:"repeat_friend_code" binding:"required"` // 重复好友码
}

// ExchangeRecordListRequest 兑换记录列表查询请求结构体
type ExchangeRecordListRequest struct {
	Page     int    `form:"page,default=1" binding:"min=1"`               // 页码，默认为1，最小值1
	PageSize int    `form:"page_size,default=10" binding:"min=1,max=100"` // 每页大小，默认10，范围1-100
	Code     string `form:"code"`                                         // 激活码
	SteamID  uint64 `form:"steam_id"`                                     // Steam用户ID筛选
	Status   int8   `form:"status"`                                       // 状态筛选
}

// ExchangeRecordListResponse 兑换记录列表响应结构体
type ExchangeRecordListResponse struct {
	Records   []ExchangeRecord `json:"rows"`      // 兑换记录列表
	Page      int              `json:"page"`      // 当前页码
	PageSize  int              `json:"page_size"` // 每页大小
	Total     int64            `json:"total"`     // 总记录数
	TotalPage int              `json:"totalpage"` // 总页数
}

type GetExchangeRecord struct {
	ID int64 `json:"id" form:"id" binding:"required,min=1"`
}

// CreateActivationCodeRequest 创建激活码请求结构体
type CreateActivationCodeRequest struct {
	Points int32 `json:"points" form:"points" binding:"required,min=1"` // 可激活点数，必填，最小值1
}

// CreateBatchActivationCodeRequest 批量创建激活码请求结构体
type CreateBatchActivationCodeRequest struct {
	Points int32 `json:"points" form:"points" binding:"required,min=100"`      // 可激活点数，必填，最小值1
	Count  int   `json:"count" form:"count" binding:"required,min=1,max=1000"` // 生成数量，必填，范围1-100
}

// UpdateActivationCodeRequest 更新激活码请求结构体
type UpdateActivationCodeRequest struct {
	ID     int64 `json:"id" form:"id" binding:"required,min=1"`
	Status int8  `json:"status" form:"status" binding:"required,oneof=1 2 3"` // 状态
}

type GetActivationCodeRequest struct {
	ID int64 `json:"id" form:"id" binding:"required,min=1"`
}

// ActivationCodeListRequest 激活码列表查询请求结构体
type ActivationCodeListRequest struct {
	Page     int    `form:"page,default=1" binding:"min=1"`               // 页码，默认为1，最小值1
	PageSize int    `form:"page_size,default=10" binding:"min=1,max=100"` // 每页大小，默认10，范围1-100
	Keyword  string `form:"keyword"`                                      // 关键词
	Status   int8   `form:"status"`                                       // 状态
}

// ActivationCodeListResponse 激活码列表响应结构体
type ActivationCodeListResponse struct {
	ActivationCodes []ActivationCode `json:"rows"`      // 激活码列表
	Page            int              `json:"page"`      // 当前页码
	PageSize        int              `json:"page_size"` // 每页大小
	Total           int64            `json:"total"`     // 总记录数
	TotalPage       int              `json:"totalpage"` // 总页数
}

type DeleteActivationCodeRequest struct {
	ID int64 `json:"id" form:"id" binding:"required,min=1"`
}

// CreateConfRequest 创建配置请求结构体
type CreateConfRequest struct {
	TutorialUrl string `json:"tutorial_url" form:"tutorial_url" binding:"required,url"` // 提货教程地址，必填，需要是有效URL
}

// UpdateConfRequest 更新配置请求结构体
type UpdateConfRequest struct {
	TutorialUrl string `json:"tutorial_url" form:"tutorial_url" binding:"omitempty,url"` // 提货教程地址，可选，需要是有效URL
}

// ConfListRequest 配置列表查询请求结构体
type ConfListRequest struct {
	Page     int `form:"page,default=1" binding:"min=1"`               // 页码，默认为1，最小值1
	PageSize int `form:"page_size,default=10" binding:"min=1,max=100"` // 每页大小，默认10，范围1-100
}

// ConfListResponse 配置列表响应结构体
type ConfListResponse struct {
	Confs     []Conf `json:"rows"`      // 配置列表
	Page      int    `json:"page"`      // 当前页码
	PageSize  int    `json:"page_size"` // 每页大小
	Total     int64  `json:"total"`     // 总记录数
	TotalPage int    `json:"totalpage"` // 总页数
}

// CreateProxyRequest 创建代理请求结构体
type CreateProxyRequest struct {
	Host     string        `json:"host" form:"host" binding:"required"`                     // 代理主机地址，必填
	Port     int32         `json:"port" form:"port" binding:"required,min=1,max=65535"`     // 代理端口，必填，范围1-65535
	Category ProxyCategory `json:"category" form:"category" binding:"required,oneof=1 2"`   // 代理分类，必填 (1=登录,2=支付)
	Protocol ProxyProtocol `json:"protocol" form:"protocol" binding:"required,oneof=1 2 3"` // 协议类型，必填 (1=HTTP,2=HTTPS,3=SOCKS5)
	Status   ProxyStatus   `json:"status" form:"status" binding:"omitempty,oneof=0 1"`      // 状态，可选
	Username string        `json:"username" form:"username" binding:"omitempty"`            // 用户名，可选
	Password string        `json:"password" form:"password" binding:"omitempty"`            // 密码，可选
	Remark   string        `json:"remark" form:"remark" binding:"omitempty,max=500"`        // 备注，可选，最大500字符
}

type GetProxyRequest struct {
	ID uint `json:"id" form:"id" binding:"required,min=1"`
}

// UpdateProxyRequest 更新代理请求结构体
type UpdateProxyRequest struct {
	ID       uint          `json:"id" form:"id" binding:"required,min=1"`
	Host     string        `json:"host" form:"host" binding:"omitempty"`                     // 代理主机地址，可选
	Port     *int32        `json:"port" form:"port" binding:"omitempty,min=1,max=65535"`     // 代理端口，可选，范围1-65535
	Category ProxyCategory `json:"category" form:"category" binding:"omitempty,oneof=1 2"`   // 代理分类，可选 (1=登录,2=支付)
	Protocol ProxyProtocol `json:"protocol" form:"protocol" binding:"omitempty,oneof=1 2 3"` // 协议类型，可选 (1=HTTP,2=HTTPS,3=SOCKS5)
	Status   ProxyStatus   `json:"status" form:"status" binding:"omitempty,oneof=1 2"`       // 状态，可选，(1=启用,2=禁用)
	Username string        `json:"username" form:"username" binding:"omitempty"`             // 用户名，可选
	Password string        `json:"password" form:"password" binding:"omitempty"`             // 密码，可选
	Remark   string        `json:"remark" form:"remark" binding:"omitempty,max=500"`         // 备注，可选，最大500字符
}

// ProxyListRequest 代理列表查询请求结构体
type ProxyListRequest struct {
	Page     int           `form:"page,default=1" binding:"min=1"`               // 页码，默认为1，最小值1
	PageSize int           `form:"page_size,default=10" binding:"min=1,max=100"` // 每页大小，默认10，范围1-100
	Keyword  string        `form:"keyword"`                                      // 搜索关键词，可搜索主机地址、备注
	Host     string        `form:"host"`                                         // 主机地址筛选
	Category ProxyCategory `form:"category"`                                     // 代理分类筛选
	Protocol ProxyProtocol `form:"protocol"`                                     // 协议类型筛选
	Status   ProxyStatus   `form:"status"`                                       // 状态筛选
}

// ProxyListResponse 代理列表响应结构体
type ProxyListResponse struct {
	Proxies   []Proxy `json:"rows"`      // 代理列表
	Page      int     `json:"page"`      // 当前页码
	PageSize  int     `json:"page_size"` // 每页大小
	Total     int64   `json:"total"`     // 总记录数
	TotalPage int     `json:"totalpage"` // 总页数
}

type DeleteProxyRequest struct {
	ID uint `json:"id" form:"id" binding:"required,min=1"`
}

// UpdateProxyStatusRequest 更新代理状态请求结构体
type UpdateProxyStatusRequest struct {
	ID     uint        `json:"id" form:"id" binding:"required,min=1"`
	Status ProxyStatus `json:"status" form:"status" binding:"required,oneof=0 1"` // 代理状态，必填
}
