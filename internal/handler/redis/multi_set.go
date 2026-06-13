package redis

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/wxsimon8888/simonStu/internal/database"
	"github.com/wxsimon8888/simonStu/internal/logger"
	"github.com/wxsimon8888/simonStu/internal/response"
)

// RedisMultiSet 并发批量写入 Redis。
//
// 请求示例：
//
//	POST /commonApi/redis/mset
//	{"values": {"key1":"val1", "key2":"val2", "key3":"val3"}}
//
// 实现思路：
//
//	每个 key-value 对开一个 goroutine 并发写入，主 goroutine 通过 channel 收集结果。
//	注意：这是"最终一致"写入——部分写入成功、部分写入失败是可能的，
//	如果要求原子性写入，应使用 Redis 的 MSET 命令或 Lua 脚本。
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

	// setResult 是内部类型，只在这个函数里用，所以定义在函数内部。
	// 这样职责清晰：这个结构体只在并发 Set 的场景下有意义。
	type setResult struct {
		Key string `json:"key"`
		Err string `json:"err,omitempty"`
	}

	// 缓冲 channel 容量 = 写入总数，确保所有 goroutine 写入时不阻塞。
	ch := make(chan setResult, len(req.Values))

	// Fan-Out：每个 key-value 对启动一个 goroutine 执行 Redis SET。
	// 循环变量 k, v 通过函数参数传值给闭包，避免捕获循环变量的问题。
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

	// Fan-In：收齐所有写入结果，统计成功/失败数量。
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
