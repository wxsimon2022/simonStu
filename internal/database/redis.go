package database

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"

	"github.com/wxsimon8888/simonStu/config"
)

var RedisClient *redis.Client

// InitRedis 初始化 Redis 客户端
func InitRedis(cfg *config.Config) {
	addr := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	// 测试连接
	ctx := context.Background()
	if err := RedisClient.Ping(ctx).Err(); err != nil {
		log.Printf("Redis 连接失败: %v", err)
	} else {
		log.Printf("Redis 连接成功 (%s)", addr)
	}
}
