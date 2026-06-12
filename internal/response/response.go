// Package response 统一 HTTP 响应格式。所有 API 返回固定结构：{ code, data, message }。
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response API 统一响应结构。
type Response struct {
	Code    int         `json:"code"`    // 业务码（与 HTTP 状态码保持一致）
	Data    interface{} `json:"data"`    // 业务数据
	Message string      `json:"message"` // 提示信息
}

// Success 返回 200 成功响应。
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Data:    data,
		Message: "success",
	})
}

// Error 返回指定 HTTP 状态码的错误响应，data 为 nil。
func Error(c *gin.Context, httpStatus int, message string) {
	c.JSON(httpStatus, Response{
		Code:    httpStatus,
		Data:    nil,
		Message: message,
	})
}
