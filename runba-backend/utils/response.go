package utils

import (
	"steam-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// ErrorResponse 发送错误响应
// 统一的错误响应格式，使用models.NewErrorResponse
func ErrorResponse(c *gin.Context, statusCode int, message string, errorDetail string) {
	var err error
	if errorDetail != "" {
		err = &simpleError{message: errorDetail}
	}
	response := models.NewErrorResponse(statusCode, message, err)
	c.JSON(statusCode, response)
}

// SuccessResponse 发送成功响应
// 统一的成功响应格式，使用models.NewSuccessResponse
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	response := models.NewSuccessResponse(message, data)
	response.Code = statusCode
	c.JSON(statusCode, response)
}

// ErrorResponseWithCode 发送带自定义错误码的错误响应
func ErrorResponseWithCode(c *gin.Context, statusCode int, message string, err error) {
	response := models.NewErrorResponse(statusCode, message, err)
	c.JSON(statusCode, response)
}

// simpleError 简单错误实现
type simpleError struct {
	message string
}

func (e *simpleError) Error() string {
	return e.message
}

// BindRequest 统一处理请求数据绑定，支持 JSON 和 urlencoded 格式
// 根据 Content-Type 自动选择合适的绑定方式
func BindRequest(c *gin.Context, obj interface{}) error {
	contentType := c.GetHeader("Content-Type")

	if contentType == "application/x-www-form-urlencoded" {
		return c.ShouldBindWith(obj, binding.Form)
	} else {
		return c.ShouldBind(obj)
	}
}
