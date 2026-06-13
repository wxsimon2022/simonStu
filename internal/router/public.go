// Package router 公开路由：无需登录即可访问。
package router

import (
	"github.com/gin-gonic/gin"

	"github.com/wxsimon8888/simonStu/internal/handler"
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
	// GET  /http/call       调用外部 HTTP API（GET 方式，传 url query 参数）
	// POST /http/call       调用外部 HTTP API（POST 方式，传 JSON body）
	// GET  /http/call/local 调用本项目的其他接口示例
	api.Any("/http/call", handler.HttpCall)
	api.GET("/http/call/local", handler.HttpCallLocal)

	// ———— 文件上传示例 ————
	// POST /upload           单文件上传（字段名 file）
	// POST /upload/multiple  多文件上传（字段名 files）
	api.POST("/upload", handler.Upload)
	api.POST("/upload/multiple", handler.UploadMultiple)

	// ———— 综合测试 ————
	api.GET("/test", handler.Test)
}
