// 路由注册。组装所有中间件和 handler，暴露 Setup 函数给 main 调用。
package router

import (
	"github.com/gin-gonic/gin"

	"github.com/wxsimon8888/simonStu/internal/handler"
	"github.com/wxsimon8888/simonStu/internal/middleware"
)

// Setup 配置所有路由并返回 *gin.Engine。
func Setup(mode string) *gin.Engine {
	gin.SetMode(mode)

	r := gin.Default()
	r.Use(middleware.CORS())

	api := r.Group("/api")
	{
		api.GET("/hello", handler.Hello)
		api.GET("/ping", handler.Ping)

		api.POST("/redis/set", handler.RedisSet)
		api.GET("/redis/get", handler.RedisGet)
		api.POST("/redis/stock/deduct", handler.StockDeduct)

		api.GET("/user", handler.UserList)
		api.PUT("/user", handler.UserUpdate)

		api.POST("/concurrent/process", handler.ConcurrentProcess)
		api.GET("/concurrent/fetch", handler.ConcurrentMultiFetch)
	}

	return r
}
