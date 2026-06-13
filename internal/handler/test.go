package handler

import (
	"encoding/json"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/wxsimon8888/simonStu/internal/database"
	"github.com/wxsimon8888/simonStu/internal/model"
	"github.com/wxsimon8888/simonStu/internal/repository"
	"github.com/wxsimon8888/simonStu/internal/response"
	"github.com/wxsimon8888/simonStu/internal/util"
)

// Test 综合示例：查询 model rows 并演示 Redis 缓存的完整流程。
func Test(c *gin.Context) {
	// ==================== 1. 从 DB 查询 model rows ====================
	var modelRows []model.Admin
	database.DB.Where("is_delete = 0").Find(&modelRows)

	// ==================== 2. model 转成前端需要的结构体 ====================
	type userItem struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
	}
	var userList []userItem
	for _, v := range modelRows {
		userList = append(userList, userItem{
			ID:       v.ID,
			Username: v.Username,
		})
	}

	// ==================== 3. 转 JSON → 存 Redis ====================
	ctx := c.Request.Context()
	cacheKey := "test:admin:list"

	// 方法一：原始方式，json.Marshal → Redis Set
	b, _ := json.Marshal(userList)
	_ = database.RedisClient.Set(ctx, cacheKey, b, repository.CacheTimeouts.Short).Err()

	// ==================== 4. 从 Redis 读取 → JSON 反序列化 ====================
	val, err := database.RedisClient.Get(ctx, cacheKey).Bytes()
	if err == nil {
		var cachedList []userItem
		json.Unmarshal(val, &cachedList)
		_ = cachedList
	}

	// ==================== 5. 方法二：用封装好的 CacheSet / CacheGet（推荐）====================
	repository.CacheSet(ctx, "cache:admin:list_v2", userList, repository.CacheTimeouts.Short)

	cachedV2, ok := repository.CacheGet[[]userItem](ctx, "cache:admin:list_v2")
	var cacheHit bool
	if ok {
		cacheHit = true
		_ = cachedV2 // 拿到 []userItem，可直接用
	}

	// ==================== 6. model 本身也可以直接存 ====================
	// model.Admin 里 password_hash / is_delete 有 json:"-" tag，不会序列化出去
	repository.CacheSet(ctx, "cache:admin:raw", modelRows, repository.CacheTimeouts.Short)
	cachedRaw, rawOk := repository.CacheGet[[]model.Admin](ctx, "cache:admin:raw")

	// ==================== 7. 缓存淘汰 ====================
	// 更新数据后调用以下函数清除缓存，下次查询重新从 DB 拉取
	// repository.CacheDelete(ctx, "cache:admin:list_v2")
	// repository.CacheClearByPrefix(ctx, "cache:admin:")

	// ==================== 8. model rows → map[key]row ====================
	// 用工具函数一行搞定，不用手写 for 循环
	mapById := util.SliceToMap(modelRows, func(v model.Admin) int { return v.ID })
	// map[int]model.Admin{ 1: {ID:1, Username:"admin", ...}, ... }
	// 查 ID=5: mapById[5].RealName

	mapByUsername := util.SliceToMap(modelRows, func(v model.Admin) string { return v.Username })
	// map[string]model.Admin{ "admin": {ID:1, Username:"admin", ...}, ... }
	// 查 username="admin": mapByUsername["admin"].Email

	// 泛型版本：按 json tag 转为通用 map
	mapByIdGeneric := util.SliceToMapGeneric(modelRows, func(v map[string]any) int {
		return int(v["id"].(float64))
	})

	// ==================== 其他演示（原样保留）====================
	var idArray []int
	for _, v := range modelRows {
		idArray = append(idArray, v.ID)
	}
	queueKey := "test:queue"
	database.RedisClient.Del(ctx, queueKey)
	database.RedisClient.RPush(ctx, queueKey, util.ToAny(idArray)...)
	_, _ = database.RedisClient.LLen(ctx, queueKey).Result()
	queueItems, _ := database.RedisClient.LRange(ctx, queueKey, 0, -1).Result()
	var queueConsumed []string
	for {
		val, err := database.RedisClient.LPop(ctx, queueKey).Result()
		if errors.Is(err, redis.Nil) {
			break
		}
		if err != nil {
			break
		}
		queueConsumed = append(queueConsumed, val)
	}

	// ==================== 返回 ====================
	response.Success(c, gin.H{
		"model_rows_count": len(modelRows),
		"user_list":        userList,
		"cache_key":        cacheKey,
		"cache_hit":        cacheHit,
		"cached_via_repo":  cachedV2,
		"cached_raw_model": map[string]any{
			"hit":     rawOk,
			"records": cachedRaw,
		},
		"map_by_id":         mapById,        // map[int]model.Admin
		"map_by_username":   mapByUsername,  // map[string]model.Admin
		"map_by_id_generic": mapByIdGeneric, // map[int]map[string]any
		"id_array":          idArray,
		"queue_items":       queueItems,
		"queue_lrange":      queueConsumed,
	})
}
