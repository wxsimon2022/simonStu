// Package repository 缓存层。提供基于 Redis + JSON 的模型数据缓存，减少数据库查询。
package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/wxsimon8888/simonStu/internal/database"
)

// CacheTimeouts 缓存过期时长，按场景区分。
var CacheTimeouts = struct {
	Short  time.Duration // 短缓存，如列表（默认 2 分钟）
	Medium time.Duration // 中等缓存，如统计数据（默认 5 分钟）
	Long   time.Duration // 长缓存，如配置类数据（默认 30 分钟）
}{
	Short:  2 * time.Minute,
	Medium: 5 * time.Minute,
	Long:   30 * time.Minute,
}

// CacheSet 将任意可 JSON 序列化的数据存入 Redis，key 建议遵循 `cache:{type}:{suffix}` 格式。
func CacheSet(ctx context.Context, key string, data any, ttl time.Duration) error {
	if database.RedisClient == nil {
		return nil // Redis 未连接时不报错，业务继续
	}
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return database.RedisClient.Set(ctx, key, b, ttl).Err()
}

// CacheGet 从 Redis 读取并 JSON 反序列化到 dest。dest 必须是指针（如 &slice）。
// 返回 true 表示命中缓存，false 表示未命中或出错。
func CacheGet[T any](ctx context.Context, key string) (*T, bool) {
	if database.RedisClient == nil {
		return nil, false
	}
	val, err := database.RedisClient.Get(ctx, key).Bytes()
	if err != nil {
		return nil, false
	}
	var data T
	if err := json.Unmarshal(val, &data); err != nil {
		return nil, false
	}
	return &data, true
}

// CacheDelete 删除指定 key 的缓存。
func CacheDelete(ctx context.Context, key string) {
	if database.RedisClient == nil {
		return
	}
	database.RedisClient.Del(ctx, key)
}

// CacheClearByPrefix 按前缀批量删除缓存（如删除 cache:user: 开头的所有 key）。
func CacheClearByPrefix(ctx context.Context, prefix string) {
	if database.RedisClient == nil {
		return
	}
	keys, _ := database.RedisClient.Keys(ctx, prefix+"*").Result()
	if len(keys) > 0 {
		database.RedisClient.Del(ctx, keys...)
	}
}
