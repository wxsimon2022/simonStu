package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/wxsimon2022/simonStu/config"
	"github.com/wxsimon2022/simonStu/internal/database"
	"github.com/wxsimon2022/simonStu/internal/dubbo"
	"github.com/wxsimon2022/simonStu/internal/handler"
	"github.com/wxsimon2022/simonStu/internal/logger"
	"github.com/wxsimon2022/simonStu/internal/nacos"
	"github.com/wxsimon2022/simonStu/internal/router"
	"github.com/wxsimon2022/simonStu/internal/service"
)

func main() {
	start := time.Now()

	cfg := config.Load()

	loc, err := time.LoadLocation(cfg.TZ)
	if err != nil {
		log.Fatalf("时区设置失败: %v", err)
	}
	time.Local = loc

	logger.Init(cfg.LogDir)

	// 并行初始化网络依赖（Redis、MySQL、Nacos 无相互依赖）
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		database.InitRedis(cfg)
	}()
	go func() {
		defer wg.Done()
		database.InitMySQL(cfg)
	}()
	go func() {
		defer wg.Done()
		if err := nacos.InitNacos(cfg); err != nil {
			log.Printf("Nacos 初始化失败（非关键，继续启动）: %v", err)
		}
	}()

	// 等待核心网络依赖就绪（从 ~3 次 RTT 降为 ~1 次 RTT 时长）
	wg.Wait()
	log.Printf("核心网络依赖初始化完成 (%v)", time.Since(start))

	// Dubbo 后台异步初始化，不阻塞 HTTP 启动
	// HTTP 先可用，Dubbo 调用在就绪前返回 503
	go func() {
		if err := dubbo.DubboInit(cfg); err != nil {
			log.Printf("Dubbo 后台初始化失败（非关键，服务已启动）: %v", err)
		} else {
			log.Printf("Dubbo 后台初始化完成")
		}
	}()

	handler.Auth = service.NewAuthService(cfg.JWTSecret, database.RedisClient)

	r := router.Setup(cfg.Mode)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("HTTP 服务启动在 %s（启动耗时 %v）", addr, time.Since(start))
	if err := r.Run(addr); err != nil {
		log.Fatalf("HTTP 服务启动失败: %v", err)
	}
}
