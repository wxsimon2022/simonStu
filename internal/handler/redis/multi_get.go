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
//
// 适用场景：
//   - 批量查用户信息、权限列表等，一次请求替代 N 次串行查询
//   - key 数量较大的时候（100+），并发收益明显
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
	logger.Infof(c, "RedisMultiGet 开始 keys=%d", total)
	start := time.Now()

	// =========================================================================
	// 方式一：Fan-Out / Fan-In
	// =========================================================================
	// 思路：给每个 key 启动一个 goroutine，各自查询 Redis 后把结果写入 channel，
	// 主 goroutine 从 channel 收齐所有结果。
	//
	// 优点：每个 key 独立并发，最快的 key 先返回，不会被慢的 key 阻塞。
	// 缺点：N 个 key 就开 N 个 goroutine，key 数量过大（数千+）时 goroutine 调度开销可观。
	// 适用：几十到几百个 key，且各 key 查询延迟差异较大的场景。
	//
	// 这里用缓冲 channel（容量 = total），防止 goroutine 发数据时阻塞。
	ch := make(chan redisGetResult, total)

	// 从请求 context 派生子 context，请求取消时所有 goroutine 能感知退出。
	// 关键细节：闭包捕获的是变量 k 的副本（通过函数参数传值），
	// 而不是循环变量 key，避免经典的"循环变量捕获陷阱"。
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
			ch <- r // 结果写入 channel，不关心谁在读
		}(key)
	}

	// 主 goroutine 收结果：从 channel 读 total 次，不管顺序。
	// Fan-In 模式的核心：多个发送者，一个接收者。
	var fanInResults []redisGetResult
	for range req.Keys {
		// 阻塞直到有结果可读。先完成的 goroutine 先被读到，
		// 所以 fanInResults 的顺序不一定与 req.Keys 一致。
		fanInResults = append(fanInResults, <-ch)
	}
	close(ch)

	// =========================================================================
	// 方式二：Worker Pool（固定 worker 数）
	// =========================================================================
	// 思路：把所有 key 放入任务通道，启动固定数量的 worker goroutine 从通道取任务，
	// 处理完一个再取下个，直到通道为空。
	//
	// 优点：goroutine 数量可控（本例 5 个），不会因 key 太多打满调度器。
	// 缺点：如果某个 key 查询特别慢，会阻塞当前 worker 处理后续 key，整体耗时 ≈ (总任务数/worker数) * 平均耗时。
	// 适用：几千上万个 key，或者需要控制 Redis 连接数防止打爆的场景。
	//
	// 关键设计：任务通道和结果通道都带缓冲，避免生产者和消费者相互阻塞。
	const workerCount = 5
	taskCh := make(chan string, total)
	resultCh := make(chan redisGetResult, total)
	var wg sync.WaitGroup

	// 把 key 全部放入任务通道后关闭，worker 遍历完自然退出。
	for _, key := range req.Keys {
		taskCh <- key
	}
	close(taskCh)

	// 启动 workerCount 个 worker，每个 worker 是一个独立的 goroutine。
	// worker 通过 for-range 从 taskCh 读取任务，通道关闭且读完时 for-range 自动退出。
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

	// 另起一个 goroutine 等待所有 worker 结束，然后关闭结果通道。
	// 如果直接在主 goroutine 里 wg.Wait()，会导致无法同时从 resultCh 收结果（死锁）。
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// 主 goroutine 收结果：worker 边处理边发，这边边收。
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
		"fan_in_results": fanInResults, // 方式一结果（goroutine 数量 = key 数量）
		"pool_results":   poolResults,  // 方式二结果（固定 5 个 goroutine）
	})
}
