package redis

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/wxsimon8888/simonStu/internal/database"
	"github.com/wxsimon8888/simonStu/internal/logger"
	"github.com/wxsimon8888/simonStu/internal/response"
)

// RedisMultiGet 并发批量查询 Redis 的多个 key。同时展示两种并发模式：
//   - Fan-Out / Fan-In：每个 key 开一个 goroutine
//   - Worker Pool：固定 5 个 worker 从通道取任务
//
// 请求示例：
//
//	POST /commonApi/redis/mget
//	{"keys": ["key1", "key2", "key3"]}
func RedisMultiGet(c *gin.Context) {
	var req RedisMultiGetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf(c, "RedisMultiGet 参数解析失败: %v", err)
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if database.RedisClient == nil {
		logger.Errorf(c, "RedisMultiGet Redis 未连接")
		response.Error(c, http.StatusServiceUnavailable, "Redis 未连接")
		return
	}

	total := len(req.Keys)
	logger.Infof(c, "RedisMultiGet 开始 keys=%d keys=%v", total, req.Keys)
	start := time.Now()

	// ———— 方式一：Fan-Out / Fan-In ————
	ch := make(chan redisGetResult, total)
	ctx := c.Request.Context()

	for _, key := range req.Keys {
		go func(k string) {
			t := time.Now()
			val, err := database.RedisClient.Get(ctx, k).Result()
			r := redisGetResult{Key: k, CostMs: time.Since(t).Milliseconds()}
			if err != nil {
				r.Err = err.Error()
			} else {
				r.Value = val
			}
			ch <- r
		}(key)
	}

	var fanInResults []redisGetResult
	for range req.Keys {
		fanInResults = append(fanInResults, <-ch)
	}
	close(ch)

	// ———— 方式二：Worker Pool ————
	const workerCount = 5
	taskCh := make(chan string, total)
	resultCh := make(chan redisGetResult, total)
	var wg sync.WaitGroup

	for _, key := range req.Keys {
		taskCh <- key
	}
	close(taskCh)

	ctx2 := c.Request.Context()
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for k := range taskCh {
				t := time.Now()
				val, err := database.RedisClient.Get(ctx2, k).Result()
				r := redisGetResult{Key: k, CostMs: time.Since(t).Milliseconds()}
				if err != nil {
					r.Err = err.Error()
				} else {
					r.Value = val
				}
				resultCh <- r
				logger.Warn.Printf("worker %d 查询 key=%s cost=%dms", workerID, k, r.CostMs)
			}
		}(i + 1)
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	var poolResults []redisGetResult
	for r := range resultCh {
		poolResults = append(poolResults, r)
	}

	elapsed := time.Since(start).Milliseconds()

	var success, fail int
	for _, r := range fanInResults {
		if r.Err == "" {
			success++
		} else {
			fail++
		}
	}

	logger.Infof(c, "RedisMultiGet 完成 total=%d success=%d fail=%d cost=%dms", total, success, fail, elapsed)

	response.Success(c, gin.H{
		"total":          total,
		"success":        success,
		"failed":         fail,
		"cost_ms":        elapsed,
		"fan_in_results": fanInResults,
		"pool_results":   poolResults,
	})
}
