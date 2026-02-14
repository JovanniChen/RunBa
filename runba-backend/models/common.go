// 通用数据模型定义
//
// 本包定义了应用程序中使用的通用数据结构，包括：
// - API响应格式：统一的响应结构，确保API的一致性
// - 错误处理：统一的错误接口和实现，支持丰富的错误信息
// - 分页结构：通用的分页信息，用于列表查询接口
// - 工具函数：创建响应和错误的便捷方法
//
// 设计原则：
// 1. 一致性：所有API使用相同的响应格式
// 2. 可扩展性：错误接口支持多种错误类型和判断方法
// 3. 易用性：提供便捷的创建函数，减少样板代码
// 4. 类型安全：使用接口和强类型，避免运行时错误
//
// Example:
//
//	// 创建成功响应
//	response := NewSuccessResponse("操作成功", userData)
//
//	// 创建错误响应
//	errorResp := NewErrorResponse(400, "请求参数错误", validationErr)
//
//	// 创建应用错误
//	appErr := NewValidationError("字段验证失败", details)
package models

import (
	"fmt"
	"net/http"
	"time"
)

// APIResponse 统一API响应结构体
//
// APIResponse定义了所有API接口的统一响应格式，确保前端能够
// 以一致的方式处理API响应。响应包含以下字段：
//
// 状态信息：
//   - Code: HTTP状态码，用于表示请求的处理结果
//   - Success: 布尔值，表示操作是否成功
//   - Message: 人类可读的消息，用于用户界面显示
//
// 数据和错误：
//   - Data: 成功响应的数据，可以是任何类型
//   - Error: 错误信息字符串，在出错时提供详细信息
//
// 使用场景：
// - 成功响应：包含数据和成功消息
// - 错误响应：包含错误信息和错误代码
// - 状态查询：返回操作状态和相关信息
//
// Example:
//
//	// 成功响应
//	response := &APIResponse{
//	    Code: 200,
//	    Success: true,
//	    Message: "用户创建成功",
//	    Data: user,
//	}
//
//	// 错误响应
//	response := &APIResponse{
//	    Code: 400,
//	    Success: false,
//	    Message: "请求参数错误",
//	    Error: "用户名不能为空",
//	}
type APIResponse struct {
	Code    int         `json:"code"`            // HTTP状态码，如200、400、500等
	Success bool        `json:"success"`         // 操作是否成功，true为成功，false为失败
	Message string      `json:"msg"`             // 响应消息，用于用户界面显示
	Data    interface{} `json:"data,omitempty"`  // 响应数据，成功时包含具体数据，失败时为空
	Error   string      `json:"error,omitempty"` // 错误信息，失败时提供详细错误描述
}

// ErrorResponse 统一错误响应结构体
// 用于返回详细的错误信息
type ErrorResponse struct {
	Code      int                    `json:"code"`                 // 错误状态码
	Success   bool                   `json:"success"`              // 固定为false
	Message   string                 `json:"msg"`                  // 错误消息
	Error     string                 `json:"error,omitempty"`      // 错误详细信息
	Details   map[string]interface{} `json:"details,omitempty"`    // 错误详情
	Timestamp time.Time              `json:"timestamp"`            // 错误发生时间
	RequestID string                 `json:"request_id,omitempty"` // 请求ID用于追踪
}

// PaginationInfo 分页信息结构体
// 通用的分页信息，可在各种列表响应中使用
type PaginationInfo struct {
	Page      int   `json:"page"`      // 当前页码
	PageSize  int   `json:"page_size"` // 每页大小
	Total     int64 `json:"total"`     // 总记录数
	TotalPage int   `json:"totalpage"` // 总页数
}

// CalculateTotalPage 计算总页数
func (p *PaginationInfo) CalculateTotalPage() {
	if p.PageSize > 0 {
		p.TotalPage = int((p.Total + int64(p.PageSize) - 1) / int64(p.PageSize))
	}
}

// NewSuccessResponse 创建成功响应
func NewSuccessResponse(message string, data interface{}) *APIResponse {
	return &APIResponse{
		Code:    200,
		Success: true,
		Message: message,
		Data:    data,
	}
}

// NewErrorResponse 创建错误响应
func NewErrorResponse(code int, message string, err error) *ErrorResponse {
	response := &ErrorResponse{
		Code:      code,
		Success:   false,
		Message:   message,
		Timestamp: time.Now(),
	}

	if err != nil {
		response.Error = err.Error()
	}

	return response
}

// ValidationError 验证错误结构体
type ValidationError struct {
	Field   string `json:"field"`           // 字段名
	Message string `json:"msg"`             // 错误消息
	Value   string `json:"value,omitempty"` // 字段值
}

// AppError 统一错误接口
type AppError interface {
	Error() string
	Code() int
	HTTPStatusCode() int
	IsRetryable() bool
	IsAuthError() bool
	IsClientError() bool
	IsServerError() bool
}

// BaseError 基础错误实现
type BaseError struct {
	ErrCode    int                    `json:"code"`
	ErrMessage string                 `json:"msg"`
	Details    map[string]interface{} `json:"details,omitempty"`
	StatusCode int                    `json:"-"`
	Detail     string                 `json:"detail,omitempty"`
	Retryable  bool                   `json:"-"`
}

// Error 实现error接口
func (e *BaseError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("错误 [%d]: %s - %s", e.ErrCode, e.ErrMessage, e.Detail)
	}
	return fmt.Sprintf("错误 [%d]: %s", e.ErrCode, e.ErrMessage)
}

// Code 返回错误代码
func (e *BaseError) Code() int {
	return e.ErrCode
}

// HTTPStatusCode 返回HTTP状态码
func (e *BaseError) HTTPStatusCode() int {
	if e.StatusCode == 0 {
		return e.ErrCode
	}
	return e.StatusCode
}

// IsRetryable 判断错误是否可重试
func (e *BaseError) IsRetryable() bool {
	if e.Retryable {
		return true
	}

	switch e.HTTPStatusCode() {
	case http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout,
		http.StatusRequestTimeout:
		return true
	default:
		return false
	}
}

// IsAuthError 判断是否为认证错误
func (e *BaseError) IsAuthError() bool {
	return e.HTTPStatusCode() == http.StatusUnauthorized
}

// IsClientError 判断是否为客户端错误
func (e *BaseError) IsClientError() bool {
	code := e.HTTPStatusCode()
	return code >= 400 && code < 500
}

// IsServerError 判断是否为服务器错误
func (e *BaseError) IsServerError() bool {
	return e.HTTPStatusCode() >= 500
}

// NewAppError 创建应用错误
func NewAppError(code int, message string, statusCode int) AppError {
	return &BaseError{
		ErrCode:    code,
		ErrMessage: message,
		StatusCode: statusCode,
	}
}

// NewValidationError 创建验证错误
func NewValidationError(message string, details map[string]interface{}) AppError {
	return &BaseError{
		ErrCode:    400,
		ErrMessage: message,
		Details:    details,
		StatusCode: http.StatusBadRequest,
	}
}

// NewNotFoundError 创建资源未找到错误
func NewNotFoundError(message string) AppError {
	return &BaseError{
		ErrCode:    404,
		ErrMessage: message,
		StatusCode: http.StatusNotFound,
	}
}

// NewUnauthorizedError 创建未授权错误
func NewUnauthorizedError(message string) AppError {
	return &BaseError{
		ErrCode:    401,
		ErrMessage: message,
		StatusCode: http.StatusUnauthorized,
	}
}

// NewForbiddenError 创建禁止访问错误
func NewForbiddenError(message string) AppError {
	return &BaseError{
		ErrCode:    403,
		ErrMessage: message,
		StatusCode: http.StatusForbidden,
	}
}

// NewConflictError 创建冲突错误
func NewConflictError(message string) AppError {
	return &BaseError{
		ErrCode:    409,
		ErrMessage: message,
		StatusCode: http.StatusConflict,
	}
}

// NewInternalError 创建内部服务器错误
func NewInternalError(message string) AppError {
	return &BaseError{
		ErrCode:    500,
		ErrMessage: message,
		StatusCode: http.StatusInternalServerError,
	}
}

// NewSteamError 创建Steam相关错误
func NewSteamError(code int, message, detail string) AppError {
	return &BaseError{
		ErrCode:    code,
		ErrMessage: message,
		Detail:     detail,
		StatusCode: code,
	}
}

// WrapError 包装普通错误
func WrapError(err error, message string) AppError {
	return &BaseError{
		ErrCode:    http.StatusInternalServerError,
		ErrMessage: message,
		Detail:     err.Error(),
		StatusCode: http.StatusInternalServerError,
	}
}
