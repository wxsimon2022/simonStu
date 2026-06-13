package database

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"

	"github.com/wxsimon8888/simonStu/config"
)

// RedisClient 全局 Redis 客户端。包外部通过 database.RedisClient 直接使用。
var RedisClient *redis.Client

// InitRedis 根据配置初始化 Redis 客户端并测试连通性。连接失败时仅打日志，不阻塞启动。
func InitRedis(cfg *config.Config) {
	addr := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	ctx := context.Background()
	if err := RedisClient.Ping(ctx).Err(); err != nil {
		log.Printf("Redis 连接失败: %v", err)
		RedisClient = nil
	} else {
		log.Printf("Redis 连接成功 (%s)", addr)
	}
}
