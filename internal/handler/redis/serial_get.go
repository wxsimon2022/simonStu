package redis

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/wxsimon8888/simonStu/internal/database"
	"github.com/wxsimon8888/simonStu/internal/logger"
	"github.com/wxsimon8888/simonStu/internal/response"
)

// RedisMultiGetOneByOne 串行查询 Redis 的多个 key（用来对比并发和串行的性能差异）。
//
// 请求示例：
//
//	POST /commonApi/redis/mget/serial
//	{"keys": ["key1", "key2", "key3"]}
func RedisMultiGetOneByOne(c *gin.Context) {
	var req RedisMultiGetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf(c, "RedisMultiGetSerial 参数解析失败: %v", err)
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if database.RedisClient == nil {
		logger.Errorf(c, "RedisMultiGetSerial Redis 未连接")
		response.Error(c, http.StatusServiceUnavailable, "Redis 未连接")
		return
	}

	start := time.Now()
	ctx := c.Request.Context()

	var results []redisGetResult
	for _, key := range req.Keys {
		t := time.Now()
		val, err := database.RedisClient.Get(ctx, key).Result()
		r := redisGetResult{Key: key, CostMs: time.Since(t).Milliseconds()}
		if err != nil {
			r.Err = err.Error()
		} else {
			r.Value = val
		}
		results = append(results, r)
	}

	elapsed := time.Since(start).Milliseconds()

	logger.Infof(c, "RedisMultiGetSerial 完成 keys=%d cost=%dms", len(req.Keys), elapsed)

	response.Success(c, gin.H{
		"total":   len(req.Keys),
		"cost_ms": elapsed,
		"results": results,
		"note":    fmt.Sprintf("串行查询 %d 个 key 耗时 %dms", len(req.Keys), elapsed),
	})
}
