package redis

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/wxsimon8888/simonStu/internal/database"
	"github.com/wxsimon8888/simonStu/internal/logger"
	"github.com/wxsimon8888/simonStu/internal/response"
)

// RedisMultiSet 并发批量写入 Redis。配合 RedisMultiGet 使用，先写入再查询。
//
// 请求示例：
//
//	POST /commonApi/redis/mset
//	{"values": {"key1":"val1", "key2":"val2", "key3":"val3"}}
func RedisMultiSet(c *gin.Context) {
	var req RedisMultiSetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf(c, "RedisMultiSet 参数解析失败: %v", err)
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if database.RedisClient == nil {
		logger.Errorf(c, "RedisMultiSet Redis 未连接")
		response.Error(c, http.StatusServiceUnavailable, "Redis 未连接")
		return
	}

	ctx := c.Request.Context()
	start := time.Now()

	type setResult struct {
		Key string `json:"key"`
		Err string `json:"err,omitempty"`
	}
	ch := make(chan setResult, len(req.Values))

	for k, v := range req.Values {
		go func(key, val string) {
			err := database.RedisClient.Set(ctx, key, val, 0).Err()
			r := setResult{Key: key}
			if err != nil {
				r.Err = err.Error()
			}
			ch <- r
		}(k, v)
	}

	var results []setResult
	var success, fail int
	for range req.Values {
		r := <-ch
		results = append(results, r)
		if r.Err == "" {
			success++
		} else {
			fail++
		}
	}
	close(ch)

	elapsed := time.Since(start).Milliseconds()

	logger.Infof(c, "RedisMultiSet 完成 total=%d success=%d fail=%d cost=%dms",
		len(req.Values), success, fail, elapsed)

	response.Success(c, gin.H{
		"keys":    len(req.Values),
		"success": success,
		"failed":  fail,
		"cost_ms": elapsed,
		"results": results,
	})
}
