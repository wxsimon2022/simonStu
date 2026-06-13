package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wxsimon8888/simonStu/internal/logger"
	"github.com/wxsimon8888/simonStu/internal/response"
)

type TestRequest struct {
	ID int `json:"id" binding:"required"`
}

func Test(c *gin.Context) {

	var req TestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf(c, "TestRequest 参数解析失败: %v", err)
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}
	response.Success(c, gin.H{
		"id": req.ID,
	})
}
