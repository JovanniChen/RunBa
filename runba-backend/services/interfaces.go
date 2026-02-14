package services

import "steam-backend/models"

// UserServiceInterface 用户服务接口
type UserServiceInterface interface {
	RegisterUser(req *models.RegisterRequest) (*models.LoginResponse, error)
	LoginUser(req *models.LoginRequest) (*models.LoginResponse, error)
	GetUserProfile(userID uint) (*models.UserInfo, error)
	UpdateUserProfile(userID uint, req *models.UpdateUserRequest) (*models.UserInfo, error)
	ChangePassword(userID uint, req *models.ChangePasswordRequest) error
	GetUsers(params *models.UserListRequest) (*models.UserListResponse, error)
	GetUserByID(userID uint) (*models.UserInfo, error)
	UpdateUser(userID uint, req *models.UpdateUserRequest) (*models.UserInfo, error)
	DeleteUser(userID uint) error
	UpdateUserStatus(userID uint, status models.UserStatus) (*models.UserInfo, error)
	GetUserStats() (map[string]interface{}, error)
}

// AccountServiceInterface Steam账户服务接口
type AccountServiceInterface interface {
	CreateAccount(req *models.CreateAccountRequest) (*models.SteamAccount, error)
	GetAccounts(params *models.AccountListRequest) (*models.AccountListResponse, error)
	GetAllAccounts() (*models.AccountListResponse, error)
	GetAccountByID(accountID uint) (*models.SteamAccount, error)
	UpdateAccount(accountID uint, req *models.UpdateAccountRequest) (*models.SteamAccount, error)
	DeleteAccount(accountID uint) error
	UpdateAccountStatus(accountID uint, status uint64) (*models.SteamAccount, error)
	AddAccountStatus(accountID uint, status uint64) (*models.SteamAccount, error)
	RemoveAccountStatus(accountID uint, status uint64) (*models.SteamAccount, error)
	SetAccountStatusSafe(accountID uint, status uint64) (*models.SteamAccount, error)
	UpdateLastLogin(accountID uint) error
	AddPoints(accountID uint, points int32) error
	DeductPoints(accountID uint, points int32) error
	GetAccountStats() (map[string]interface{}, error)
	// TokenContent相关方法
	UpdateTokenContent(accountID uint, tokenContent *models.TokenContent) error
	GetTokenContent(accountID uint) (*models.TokenContent, error)
	ValidateToken(accountID uint) (bool, error)
	GetAccountWithTokenInfo(accountID uint, includeTokenContent bool) (*models.AccountTokenResponse, error)
}

// RecordServiceInterface 兑换记录服务接口
type RecordServiceInterface interface {
	CreateExchangeRecord(req *models.CreateExchangeRecordRequest) (*models.ExchangeRecord, error)
	GetExchangeRecords(params *models.ExchangeRecordListRequest) (*models.ExchangeRecordListResponse, error)
	GetExchangeRecordByID(recordID uint) (*models.ExchangeRecord, error)
	GetExchangeRecordByIDWithOptions(recordID uint, includeGifts bool) (*models.ExchangeRecord, error)
}

// ActivationCodeServiceInterface 激活码服务接口
type ActivationCodeServiceInterface interface {
	CreateActivationCode(req *models.CreateActivationCodeRequest) (*models.ActivationCode, error)
	CreateBatchActivationCodes(req *models.CreateBatchActivationCodeRequest) ([]models.ActivationCode, error)
	GetActivationCodes(params *models.ActivationCodeListRequest) (*models.ActivationCodeListResponse, error)
	GetActivationCodeByID(codeID uint) (*models.ActivationCode, error)
	GetActivationCodeByCode(code string) (*models.ActivationCode, error)
	UpdateActivationCode(codeID uint, req *models.UpdateActivationCodeRequest) (*models.ActivationCode, error)
	DeleteActivationCode(codeID uint) error
	GetActivationCodeStats() (map[string]interface{}, error)
	ValidateActivationCode(code string) (bool, *models.ActivationCode, error)
}

// CertificateServiceInterface 证书服务接口
type CertificateServiceInterface interface {
	CreateCertificate(req *models.CreateCertificateRequest) (*models.Certificate, error)
	GetCertificates(params *models.CertificateListRequest) (*models.CertificateListResponse, error)
	GetCertificateByID(certificateID uint) (*models.Certificate, error)
	GetCertificateByNo(no string) (*models.Certificate, error)
	UpdateCertificate(certificateID uint, req *models.UpdateCertificateRequest) (*models.Certificate, error)
	DeleteCertificate(certificateID uint) error
}

// ForgeServiceInterface 锻刀所服务接口
type ForgeServiceInterface interface {
	CreateForge(req *models.CreateForgeRequest) (*models.Forge, error)
	GetForges(params *models.ForgeListRequest) (*models.ForgeListResponse, error)
	GetForgeByID(forgeID uint) (*models.Forge, error)
	UpdateForge(forgeID uint, req *models.UpdateForgeRequest) (*models.Forge, error)
	DeleteForge(forgeID uint) error
}

// SwordsmithServiceInterface 刀匠服务接口
type SwordsmithServiceInterface interface {
	CreateSwordsmith(req *models.CreateSwordsmithRequest) (*models.Swordsmith, error)
	GetSwordsmiths(params *models.SwordsmithListRequest) (*models.SwordsmithListResponse, error)
	GetSwordsmithByID(swordsmithID uint) (*models.Swordsmith, error)
	UpdateSwordsmith(swordsmithID uint, req *models.UpdateSwordsmithRequest) (*models.Swordsmith, error)
	DeleteSwordsmith(swordsmithID uint) error
}
