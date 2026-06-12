package main

import (
	"fmt"
	"log"

	"github.com/wxsimon8888/simonStu/config"
	"github.com/wxsimon8888/simonStu/internal/database"
	"github.com/wxsimon8888/simonStu/internal/logger"
	"github.com/wxsimon8888/simonStu/internal/router"
)

func main() {
	cfg := config.Load()

	// 初始化日志（会按 LOG_DIR 创建目录和文件）
	logger.Init(cfg.LogDir)

	// 初始化 Redis
	database.InitRedis(cfg)
	// 初始化 MySQL
	database.InitMySQL(cfg)

	r := router.Setup(cfg.Mode)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("服务启动在 %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
