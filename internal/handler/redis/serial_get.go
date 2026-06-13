package redis

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/wxsimon2022/simonStu/internal/database"
	"github.com/wxsimon2022/simonStu/internal/logger"
	"github.com/wxsimon2022/simonStu/internal/response"
)

// RedisMultiGetOneByOne 串行查询 Redis 的多个 key。
//
// 请求示例：
//
//	POST /commonApi/redis/mget/serial
//	{"keys": ["key1", "key2", "key3"]}
//
// 设计意图：
//
//	与 RedisMultiGet 配合使用，让调用方能直观对比"串行"和"并发"的性能差距。
//	在 key 数量较大的场景下，串行总耗时 ≈ 各 key 耗时之和，并发总耗时 ≈ 最慢的 key 耗时。
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

	// 串行遍历：一个一个查，前一个查询返回后才发起下一个。
	// 这是最直观但最慢的方式，性能瓶颈在于网络往返（RTT）。
	// 假设 Redis RTT 是 1ms，查 100 个 key 至少需要 100ms（理想情况），
	// 而并发模式下同样 100 个 key 只需要 ≈ 1ms（全部并行）。
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
