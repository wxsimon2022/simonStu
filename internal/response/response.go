// Package response 统一 HTTP 响应格式。所有 API 返回固定结构：{ code, data, message }。
// data 字段总是对象（即使是空对象 {}），避免前端处理 null。
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response API 统一响应结构。
type Response struct {
	Code      int         `json:"code"`
	Data      interface{} `json:"data"`
	Message   string      `json:"message"`
	ErrorCode int         `json:"error_code,omitempty"`
}

// Success 返回 200 成功响应。
func Success(c *gin.Context, data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Data:    data,
		Message: "success",
	})
}

// Error 返回错误响应。errorCode 不传时默认为 1。
func Error(c *gin.Context, httpStatus int, message string, errorCode ...int) {
	code := 1
	if len(errorCode) > 0 {
		code = errorCode[0]
	}
	c.JSON(httpStatus, Response{
		Code:      httpStatus,
		Data:      gin.H{},
		Message:   message,
		ErrorCode: code,
	})
}
