package constants

// API路径常量定义
// 统一管理所有路由路径，避免硬编码字符串

// API版本
const (
	APIVersion1 = "/api/v1"
)

// 认证相关路径
const (
	AuthGroup          = "/auth"
	AuthRegister       = "/register"
	AuthLogin          = "/login"
	AuthProfile        = "/profile"
	AuthChangePassword = "/change-password"
)

// 用户管理路径
const (
	UsersGroup       = "/users"
	GetUsers         = "/getUsers"
	GetUserByID      = "/getUserByID"
	UpdateUser       = "/updateUser"
	DeleteUserByID   = "/deleteUserByID"
	UpdateUserStatus = "/updateUserStatus"
)

// Steam账号管理路径
const (
	AccountsGroup       = "/accounts"
	CreateAccount       = "/createAccount"
	GetAccounts         = "/getAccounts"
	GetAccountByID      = "/getAccountByID"
	UpdateAccount       = "/updateAccount"
	DeleteAccountByID   = "/deleteAccountByID"
	UpdateAccountStatus = "/updateAccountStatus"
	GetAccountStats     = "/stats"
	AddPoints           = "/points/add"
	DeductPoints        = "/points/deduct"
	UpdateTokenContent  = "/token"
	GetTokenContent     = "/token"
	ValidateToken       = "/token/validate"
	GetAccountWithToken = "/tokeninfo"
)

// 证书管理路径
const (
	CertificatesGroup     = "/certificates"
	CreateCertificate     = "/createCertificate"
	GetCertificates       = "/getCertificates"
	GetCertificateByID    = "/getCertificateByID"
	UpdateCertificate     = "/updateCertificate"
	DeleteCertificateByID = "/deleteCertificateByID"
	GetCertificateByNo    = "/public"
)

// 锻刀所管理路径
const (
	ForgesGroup     = "/forges"
	CreateForge     = "/createForge"
	GetForges       = "/getForges"
	GetForgeByID    = "/getForgeByID"
	UpdateForge     = "/updateForge"
	DeleteForgeByID = "/deleteForgeByID"
)

// 刀匠管理路径
const (
	SwordsmithsGroup     = "/swordsmiths"
	CreateSwordsmith     = "/createSwordsmith"
	GetSwordsmiths       = "/getSwordsmiths"
	GetSwordsmithByID    = "/getSwordsmithByID"
	UpdateSwordsmith     = "/updateSwordsmith"
	DeleteSwordsmithByID = "/deleteSwordsmithByID"
)

// 兑换记录路径
const (
	RecordsGroup          = "/records"
	CreateExchangeRecord  = "/createExchangeRecord"
	GetExchangeRecords    = "/getExchangeRecords"
	GetExchangeRecordByID = "/getExchangeRecordByID"
)

// 激活码管理路径
const (
	ActivationCodesGroup       = "/activation-codes"
	CreateActivationCode       = "/createActivationCode"
	CreateBatchActivationCodes = "/createBatchActivationCodes"
	GetActivationCodes         = "/getActivationCodes"
	GetActivationCode          = "/getActivationCode"
	UpdateActivationCode       = "/updateActivationCode"
	DeleteActivationCode       = "/deleteActivationCode"
	GetActivationCodeStats     = "/stats"
	ValidateActivationCode     = "/validate"
)

// 代理管理路径
const (
	ProxiesGroup         = "/proxies"
	CreateProxy          = "/createProxy"
	GetProxies           = "/getProxies"
	GetProxyByID         = "/getProxyByID"
	UpdateProxy          = "/updateProxy"
	DeleteProxy          = "/deleteProxy"
	GetEnabledProxies    = "/enabled"
	GetProxiesByCategory = "/category"
	GetProxiesByProtocol = "/protocol"
	UpdateProxyStatus    = "/status"
)

// 系统配置路径
const (
	ConfsGroup         = "/confs"
	CreateConf         = "/createConf"
	CreateOrUpdateConf = "/createOrUpdateConf"
	GetConfs           = "/getConfs"
	GetConf            = "/getConf"
	GetConfByID        = "/getConfByID"
	UpdateConf         = "/updateConf"
	DeleteConf         = "/deleteConf"
	GetLatestConf      = "/latest"
)

// 系统路径
const (
	Health = "/health"
	Root   = "/"
)
