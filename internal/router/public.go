// Package router 公开路由：无需登录即可访问。
package router

import (
	"github.com/gin-gonic/gin"

	"github.com/wxsimon8888/simonStu/internal/handler"
)

// RegisterPublicRoutes 注册所有不需要 JWT 认证的公开接口。
//
// 路由一览：
//
//	POST /auth/login            普通用户登录，返回 JWT 令牌
//	POST /auth/admin/login      管理后台管理员登录，返回 JWT 令牌
//	GET  /hello                 服务健康检查
//	GET  /ping                  连通性测试
//	POST /redis/set             向 Redis 写入键值对（测试用）
//	GET  /redis/get             从 Redis 读取值（测试用）
//	POST /redis/stock/deduct    使用 Lua 脚本扣减 Redis 库存（测试用）
//	POST /concurrent/process    并发任务处理示例
//	GET  /concurrent/fetch      并发多源数据抓取示例
func RegisterPublicRoutes(api *gin.RouterGroup) {

	// GET /hello — 服务健康检查
	api.GET("/hello", handler.Hello)
	// GET /ping — 连通性测试
	api.GET("/ping", handler.Ping)
	// POST /redis/set — 向 Redis 写入键值对
	api.POST("/redis/set", handler.RedisSet)
	// GET /redis/get — 从 Redis 读取指定键的值
	api.GET("/redis/get", handler.RedisGet)
	// POST /redis/stock/deduct — 使用 Lua 脚本原子扣减 Redis 库存
	api.POST("/redis/stock/deduct", handler.StockDeduct)
	// POST /concurrent/process — 并发任务处理示例
	api.POST("/concurrent/process", handler.ConcurrentProcess)
	// GET /concurrent/fetch — 并发多源数据抓取示例
	api.GET("/concurrent/fetch", handler.ConcurrentMultiFetch)

	api.GET("/test", handler.Test)
}
