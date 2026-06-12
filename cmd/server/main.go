// Package main 应用入口。初始化配置、日志、数据库后启动 HTTP 服务。
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/wxsimon8888/simonStu/config"
	"github.com/wxsimon8888/simonStu/internal/database"
	"github.com/wxsimon8888/simonStu/internal/logger"
	"github.com/wxsimon8888/simonStu/internal/router"
)

func main() {
	cfg := config.Load()

	// 设置时区（影响 time.Now 等）
	loc, err := time.LoadLocation(cfg.TZ)
	if err != nil {
		log.Fatalf("时区设置失败: %v", err)
	}
	time.Local = loc

	logger.Init(cfg.LogDir)

	database.InitRedis(cfg)
	database.InitMySQL(cfg)

	r := router.Setup(cfg.Mode)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("HTTP 服务启动在 %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("HTTP 服务启动失败: %v", err)
	}
}
