package handler

import (
	"context"

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

	var idList []int64
	database.DB.Table("c_admin").Select("id").Where("is_delete = 0").Pluck("id", &idList)

	var idArray []int
	for _, v := range modelRows {
		idArray = append(idArray, v.ID)
	}

	// Redis 存数据示例
	ctx := context.Background()
	redisKey := "test:hello"
	redisVal := "Hello Redis!"
	_ = database.RedisClient.Set(ctx, redisKey, redisVal, 0).Err()

	cached, _ := database.RedisClient.Get(ctx, redisKey).Result()

	// Redis 存储结构体示例（序列化为 JSON）
	_ = database.RedisClient.Set(ctx, "test:user:1", `{"id":1,"username":"admin"}`, 0).Err()
	userCached, _ := database.RedisClient.Get(ctx, "test:user:1").Result()

	redisCountKey := "redis:count"
	redisCount, _ := database.RedisClient.Incr(ctx, redisCountKey).Result()

	response.Success(c, gin.H{
		"id":          req.ID,
		"total":       total,
		"list":        rows,
		"model_list":  modelList,
		"idList":      idList,
		"idArray":     idArray,
		"redis_key":   redisKey,
		"redis_value": cached,
		"redis_user":  userCached,
		"redisCount":  redisCount,
	})
}
