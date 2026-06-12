package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// Success 返回成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Data:    data,
		Message: "success",
	})
}

// Error 返回错误响应
func Error(c *gin.Context, httpStatus int, message string) {
	c.JSON(httpStatus, Response{
		Code:    httpStatus,
		Data:    nil,
		Message: message,
	})
}
