// Package router 公开路由：无需登录即可访问。
package router

import (
	"github.com/gin-gonic/gin"

	"github.com/wxsimon8888/simonStu/internal/handler"
	"github.com/wxsimon8888/simonStu/internal/handler/httpcall"
	"github.com/wxsimon8888/simonStu/internal/handler/redis"
	"github.com/wxsimon8888/simonStu/internal/handler/upload"
)

// RegisterPublicRoutes 注册所有不需要 JWT 认证的公开接口。
func RegisterPublicRoutes(api *gin.RouterGroup) {

	// GET /hello — 服务健康检查
	api.GET("/hello", handler.Hello)
	// GET /ping — 连通性测试
	api.GET("/ping", handler.Ping)

	// ———— Redis 操作示例 ————
	// POST /redis/set — 向 Redis 写入键值对
	api.POST("/redis/set", handler.RedisSet)
	// GET /redis/get — 从 Redis 读取指定键的值
	api.GET("/redis/get", handler.RedisGet)
	// POST /redis/stock/deduct — 使用 Lua 脚本原子扣减 Redis 库存
	api.POST("/redis/stock/deduct", handler.StockDeduct)

	// ———— 并发编程示例 ————
	// POST /concurrent/process — 并发任务处理示例
	api.POST("/concurrent/process", handler.ConcurrentProcess)
	// GET /concurrent/fetch — 并发多源数据抓取示例
	api.GET("/concurrent/fetch", handler.ConcurrentMultiFetch)

	// ———— HTTP 调用外部 API 示例 ————
	api.Any("/http/call", httpcall.HttpCall)
	api.GET("/http/call/local", httpcall.HttpCallLocal)

	// ———— 文件上传示例 ————
	api.POST("/upload", upload.Upload)
	api.POST("/upload/multiple", upload.UploadMultiple)

	// ———— 多协程 Redis 操作示例 ————
	api.POST("/redis/mget", redis.RedisMultiGet)
	api.POST("/redis/mget/serial", redis.RedisMultiGetOneByOne)
	api.POST("/redis/mset", redis.RedisMultiSet)

	// ———— 综合测试 ————
	api.GET("/test", handler.Test)
}
