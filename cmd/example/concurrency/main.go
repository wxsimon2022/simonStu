package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Task struct {
	ID   int
	Name string
}

type Result struct {
	TaskID  int
	Success bool
	Elapsed time.Duration
}

// worker 从 tasks 通道读取任务，处理后发送结果到 results 通道
func worker(ctx context.Context, id int, tasks <-chan Task, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("  worker %d 收到退出信号，停止\n", id)
			return
		case task, ok := <-tasks:
			if !ok {
				return // tasks 通道已关闭
			}

			start := time.Now()
			// 模拟处理耗时（100-600ms）
			workTime := time.Duration(rand.Intn(500)+100) * time.Millisecond
			time.Sleep(workTime)

			// 模拟 10% 概率失败
			success := rand.Float64() > 0.1

			results <- Result{
				TaskID:  task.ID,
				Success: success,
				Elapsed: time.Since(start),
			}

			status := "ok"
			if !success {
				status = "fail"
			}
			fmt.Printf("  worker %d 处理 task %d [%s] (%v)\n", id, task.ID, status, time.Since(start).Round(time.Millisecond))
		}
	}
}

func main() {
	const (
		numTasks  = 20
		numWorker = 3
	)

	rand.New(rand.NewSource(time.Now().UnixNano()))

	// --- 1. 生成任务 ---
	tasks := make(chan Task, numTasks)
	for i := range numTasks {
		tasks <- Task{ID: i + 1, Name: fmt.Sprintf("任务-%d", i+1)}
	}
	close(tasks) // 发送完毕关闭通道，worker 会自然退出

	// --- 2. 启动 worker ---
	results := make(chan Result, numTasks)
	var wg sync.WaitGroup

	// 带超时的 context，模拟优雅退出
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	for i := range numWorker {
		wg.Add(1)
		go worker(ctx, i+1, tasks, results, &wg)
	}

	// --- 3. 收集结果 ---
	// 开一个 goroutine 等待所有 worker 结束，然后关闭 results 通道
	go func() {
		wg.Wait()
		close(results)
	}()

	// --- 4. 主 goroutine 汇总 ---
	var total, success, fail int
	for r := range results {
		total++
		if r.Success {
			success++
		} else {
			fail++
		}
	}

	// --- 5. 输出汇总 ---
	fmt.Println("\n===== 执行汇总 =====")
	fmt.Printf("总任务: %d\n", total)
	fmt.Printf("成功:   %d\n", success)
	fmt.Printf("失败:   %d\n", fail)
	fmt.Printf("worker: %d\n", numWorker)
	fmt.Println("====================")
}
