// Package handler 接收请求、调用下层、返回标准响应。
package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/wxsimon2022/simonStu/internal/response"
)

// Hello 返回基础问候信息。
func Hello(c *gin.Context) {
	response.Success(c, gin.H{
		"say": "你好，世界！",
	})
}

// Ping 健康检查。
func Ping(c *gin.Context) {
	response.Success(c, gin.H{
		"say": "pong",
	})
}
