package handler

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/wxsimon8888/simonStu/internal/logger"
	"github.com/wxsimon8888/simonStu/internal/response"
)

type ProcessRequest struct {
	Items       []string `json:"items" binding:"required,min=1"`
	Concurrency int      `json:"concurrency" binding:"min=1"`
}

type taskResult struct {
	Item   string `json:"item"`
	Status string `json:"status"` // success / fail
	CostMs int64  `json:"cost_ms"`
}

// ConcurrentProcess 并发处理示例 — Worker Pool 模式
func ConcurrentProcess(c *gin.Context) {
	var req ProcessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf(c, "ConcurrentProcess 参数解析失败: %v", err)
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	// 默认并发数 = 3
	if req.Concurrency < 1 {
		req.Concurrency = 3
	}
	// 不超过任务总数
	if req.Concurrency > len(req.Items) {
		req.Concurrency = len(req.Items)
	}

	total := len(req.Items)
	logger.Infof(c, "ConcurrentProcess 开始 items=%d concurrency=%d", total, req.Concurrency)

	// --- 1. 任务通道（缓冲）---
	tasks := make(chan string, total)
	for _, item := range req.Items {
		tasks <- item
	}
	close(tasks)

	// --- 2. 结果通道 + WaitGroup ---
	results := make(chan taskResult, total)
	var wg sync.WaitGroup

	// 带超时的 context（给所有 worker 共享）
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	start := time.Now()

	// --- 3. 启动 N 个 worker goroutine ---
	for i := 0; i < req.Concurrency; i++ {
		wg.Add(1)
		go worker(ctx, i+1, tasks, results, &wg)
	}

	// --- 4. 等待所有 worker 结束，关闭结果通道 ---
	go func() {
		wg.Wait()
		close(results)
	}()

	// --- 5. 主 goroutine 收集结果 ---
	var list []taskResult
	var success, fail int

	for r := range results {
		list = append(list, r)
		if r.Status == "success" {
			success++
		} else {
			fail++
		}
	}

	elapsed := time.Since(start).Milliseconds()

	logger.Infof(c, "ConcurrentProcess 完成 total=%d success=%d fail=%d cost=%dms",
		total, success, fail, elapsed)

	response.Success(c, gin.H{
		"total":       total,
		"success":     success,
		"failed":      fail,
		"concurrency": req.Concurrency,
		"cost_ms":     elapsed,
		"results":     list,
	})
}

// worker 从 tasks 通道取任务，处理后发送结果到 results
func worker(ctx context.Context, id int, tasks <-chan string, results chan<- taskResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			// context 超时或被取消，优雅退出
			logger.Warn.Printf("worker %d 收到退出信号，提前停止", id)
			return
		case item, ok := <-tasks:
			if !ok {
				return // tasks 通道已关闭，正常结束
			}

			// 模拟处理耗时（50-300ms）
			workTime := time.Duration(50+rand.Intn(250)) * time.Millisecond
			time.Sleep(workTime)

			costMs := workTime.Milliseconds()
			status := "success"
			// 模拟 15% 概率失败
			if rand.Float64() < 0.15 {
				status = "fail"
			}

			results <- taskResult{
				Item:   item,
				Status: status,
				CostMs: costMs,
			}
			logger.Warn.Printf("worker %d 处理完成 item=%s status=%s cost=%dms", id, item, status, costMs)
		}
	}
}

// ConcurrentMultiFetch 并发请求示例 — Fan-Out / Fan-In 模式
func ConcurrentMultiFetch(c *gin.Context) {
	sources := []string{"数据源-A", "数据源-B", "数据源-C", "数据源-D", "数据源-E"}
	results := make(chan string, len(sources))

	// 并发发起多个请求
	for _, src := range sources {
		go func(name string) {
			workTime := time.Duration(100+rand.Intn(400)) * time.Millisecond
			time.Sleep(workTime)
			results <- fmt.Sprintf("%s 响应 (%v)", name, workTime.Round(time.Millisecond))
		}(src)
	}

	// 收集所有结果
	var list []string
	for range sources {
		list = append(list, <-results)
	}
	close(results)

	response.Success(c, gin.H{
		"sources": list,
	})
}
