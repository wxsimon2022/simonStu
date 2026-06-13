package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/wxsimon8888/simonStu/internal/database"
	"github.com/wxsimon8888/simonStu/internal/model"
	"github.com/wxsimon8888/simonStu/internal/response"
)

type TestRequest struct {
	ID int `json:"id" binding:"required"`
}

type userInfo struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

func Test(c *gin.Context) {

	var req TestRequest
	//if err := c.ShouldBindJSON(&req); err != nil {
	//	logger.Errorf(c, "TestRequest 参数解析失败: %v", err)
	//	response.Error(c, http.StatusBadRequest, "参数错误", 9)
	//	//return
	//}

	// 查询总数
	var total int64
	database.DB.Table("c_admin").Where("is_delete = 0").Count(&total)

	var rows []userInfo
	database.DB.Table("c_admin").Select("id,username").Where("is_delete = 0").Scan(&rows)

	// 通过 Model 查询示例（GORM 自动识别表名和字段）
	var modelRows []model.Admin
	database.DB.Where("is_delete = 0").Find(&modelRows)

	var modelList []userInfo
	for _, v := range modelRows {
		modelList = append(modelList, userInfo{
			ID:       v.ID,
			Username: v.Username,
		})
	}

	response.Success(c, gin.H{
		"id":         req.ID,
		"total":      total,
		"list":       rows,
		"model_list": modelList,
	})
}
