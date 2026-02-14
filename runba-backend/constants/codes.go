package constants

// HTTP状态码常量
// 提供统一的状态码定义，避免在代码中直接使用数字
const (
	// 成功状态码
	CodeSuccess = 200 // 请求成功

	CodeParamFailure = -301 // 请求参数错误

	// account.go
	CodeCreateAccountFailure       = -302 // 创建钱包失败
	CodeGetAccountsFailure         = -303 // 创建钱包失败
	CodeGetAccountFailure          = -304 // 获取钱包信息失败
	CodeUpdateAccountFailure       = -305 // 更新钱包信息失败
	CodeDeleteAccountFailure       = -306 // 更新钱包信息失败
	CodeUpdateAccountStatusFailure = -307 // 更新钱包状态失败

	// certificate.go
	CodeCreateCertificateFailure = -338 // 创建证书失败
	CodeGetCertificatesFailure   = -339 // 获取证书列表失败
	CodeGetCertificateFailure    = -340 // 获取证书信息失败
	CodeUpdateCertificateFailure = -341 // 更新证书信息失败
	CodeDeleteCertificateFailure = -342 // 删除证书失败

	// forge.go
	CodeCreateForgeFailure = -343 // 创建锻刀所失败
	CodeGetForgesFailure   = -344 // 获取锻刀所列表失败
	CodeGetForgeFailure    = -345 // 获取锻刀所信息失败
	CodeUpdateForgeFailure = -346 // 更新锻刀所信息失败
	CodeDeleteForgeFailure = -347 // 删除锻刀所失败

	// swordsmith.go
	CodeCreateSwordsmithFailure = -348 // 创建刀匠失败
	CodeGetSwordsmithsFailure   = -349 // 获取刀匠列表失败
	CodeGetSwordsmithFailure    = -350 // 获取刀匠信息失败
	CodeUpdateSwordsmithFailure = -351 // 更新刀匠信息失败
	CodeDeleteSwordsmithFailure = -352 // 删除刀匠失败

	// record.go
	CodeCreateExchangeRecordFailure = -308 // 创建兑换记录失败
	CodeGetExchangeRecordsFailure   = -309 // 获取兑换记录列表失败
	CodeGetExchangeRecordFailure    = -310 // 获取兑换记录失败

	// activation_code.go
	CodeCreateActivationCodeFailure       = -311 // 创建激活码失败
	CodeCreateBatchActivationCodesFailure = -312 // 批量创建激活码失败
	CodeGetActivationCodesFailure         = -313 // 批量创建激活码失败
	CodeActivationCodeNotExists           = -314 // 激活码不存在
	CodeGetActivationCodeFailure          = -315 // 获取激活码失败
	CodeGetActivationCodeExists           = -316 // 激活码已存在
	CodeUpdateActivationCodeFailure       = -317 // 更新激活码失败
	CodeDeleteActivationCodeFailure       = -318 // 删除激活码失败

	// proxy.go
	CodeCreateProxyFailure       = -319 // 创建代理失败
	CodeProxyNotExists           = -320 // 代理不存在
	CodeGetProxiesFailure        = -321 // 获取代理列表失败
	CodeUpdateProxyFailure       = -322 // 更新代理失败
	CodeDeleteProxyFailure       = -323 // 删除代理失败
	CodeUpdateProxyStatusFailure = -324 // 更新代理状态失败

	// conf.go
	CodeCreateOrUpdateConfFailure = -325 // 配置操作失败
	CodeConfNotExists             = -326 // 获取最新配置失败

	// user.go
	CodeUsernameOrPasswordFailure = -327 // 用户名或密码错误
	CodeLoginFailure              = -328 // 登录失败
	CodeUserNoExists              = -329 // 用户不存在
	CodeGetUserFailure            = -330 // 获取用户信息失败
	CodeUpdateUserFailure         = -331 // 更新用户信息失败
	CodeRepeatPasswordFailure     = -332 // 两次输入密码不一致
	CodePasswordFailure           = -333 // 原密码错误
	CodeModifyPasswordFailure     = -334 // 密码修改失败
	CodeGetUsersFailure           = -335 // 获取用户列表失败
	CodeDeleteUserFailure         = -336 // 删除用户失败
	CodeUpdateUserStatusFailure   = -337 // 删除用户状态失败

	CodeUnauthorized = -400 // 用户未认证
)
