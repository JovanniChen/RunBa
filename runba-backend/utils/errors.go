package utils

import (
	"steam-backend/models"
)

// AppError 类型别名，使用统一的错误接口
type AppError = models.AppError

// 创建各种错误的便捷函数
var (
	NewAppError          = models.NewAppError
	NewValidationError   = models.NewValidationError
	NewNotFoundError     = models.NewNotFoundError
	NewUnauthorizedError = models.NewUnauthorizedError
	NewForbiddenError    = models.NewForbiddenError
	NewConflictError     = models.NewConflictError
	NewInternalError     = models.NewInternalError
	WrapError            = models.WrapError
)
