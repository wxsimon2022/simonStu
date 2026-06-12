package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/wxsimon8888/simonStu/internal/database"
	"github.com/wxsimon8888/simonStu/internal/logger"
	"github.com/wxsimon8888/simonStu/internal/response"
)

type RedisSetRequest struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
}

var stockDeductScript = redis.NewScript(`
    local stock = redis.call('GET', KEYS[1])
    if not stock then
        return {0, "库存不存在"}
    end
    local num = tonumber(stock)
    local qty = tonumber(ARGV[1])
    if num < qty then
        return {0, "库存不足"}
    end
    redis.call('DECRBY', KEYS[1], qty)
    return {1, num - qty}
`)

// RedisSet 写入 key-value 到 Redis
func RedisSet(c *gin.Context) {
	var req RedisSetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf(c, "RedisSet 参数解析失败: %v", err)
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	ctx := context.Background()
	if err := database.RedisClient.Set(ctx, req.Key, req.Value, 0).Err(); err != nil {
		logger.Errorf(c, "RedisSet 写入失败 key=%s err=%v", req.Key, err)
		response.Error(c, http.StatusInternalServerError, "写入 Redis 失败: "+err.Error())
		return
	}

	logger.Infof(c, "RedisSet 写入成功 key=%s value=%s", req.Key, req.Value)
	response.Success(c, gin.H{
		"key":   req.Key,
		"value": req.Value,
	})
}

// RedisGet 从 Redis 读取 key 的值
func RedisGet(c *gin.Context) {
	key := c.DefaultQuery("key", "simon")
	if key == "" {
		logger.Errorf(c, "RedisGet key 为空")
		response.Error(c, http.StatusBadRequest, "参数 key 不能为空")
		return
	}

	ctx := context.Background()
	val, err := database.RedisClient.Get(ctx, key).Result()
	if err != nil {
		logger.Errorf(c, "RedisGet 读取失败 key=%s err=%v", key, err)
		response.Error(c, http.StatusNotFound, "key 不存在或读取失败: "+err.Error())
		return
	}

	logger.Infof(c, "RedisGet 读取成功 key=%s", key)
	response.Success(c, gin.H{
		"key":   key,
		"value": val,
	})
}

type StockDeductRequest struct {
	Key      string `json:"key" binding:"required"`
	Quantity int    `json:"quantity" binding:"required,min=1"`
}

// StockDeduct 原子扣减库存
func StockDeduct(c *gin.Context) {
	var req StockDeductRequest

	if c.Request.Method == http.MethodGet {
		req.Key = c.DefaultQuery("key", "")
		qtyStr := c.DefaultQuery("quantity", "1")
		qty, err := strconv.Atoi(qtyStr)
		if err != nil || qty < 1 {
			logger.Errorf(c, "StockDeduct quantity 参数无效: %s", qtyStr)
			response.Error(c, http.StatusBadRequest, "quantity 必须是正整数")
			return
		}
		req.Quantity = qty
	} else {
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.Errorf(c, "StockDeduct 参数解析失败: %v", err)
			response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
			return
		}
	}

	if req.Key == "" {
		logger.Errorf(c, "StockDeduct key 为空")
		response.Error(c, http.StatusBadRequest, "参数 key 不能为空")
		return
	}

	ctx := context.Background()
	result, err := stockDeductScript.Run(ctx, database.RedisClient,
		[]string{req.Key}, req.Quantity,
	).Slice()

	if err != nil {
		logger.Errorf(c, "StockDeduct 执行 Lua 脚本失败 key=%s err=%v", req.Key, err)
		response.Error(c, http.StatusInternalServerError, "扣减库存失败: "+err.Error())
		return
	}

	code := int(result[0].(int64))
	if code == 0 {
		msg := result[1].(string)
		logger.Errorf(c, "StockDeduct 扣减失败 key=%s qty=%d reason=%s", req.Key, req.Quantity, msg)
		response.Error(c, http.StatusBadRequest, msg)
		return
	}

	remaining := int(result[1].(int64))
	logger.Infof(c, "StockDeduct 扣减成功 key=%s qty=%d remaining=%d", req.Key, req.Quantity, remaining)
	response.Success(c, gin.H{
		"key":       req.Key,
		"quantity":  req.Quantity,
		"remaining": remaining,
	})
}
