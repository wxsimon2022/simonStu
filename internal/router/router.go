package router

import (
	"github.com/gin-gonic/gin"

	"github.com/wxsimon8888/simonStu/internal/handler"
	"github.com/wxsimon8888/simonStu/internal/middleware"
)

// Setup 配置所有路由并返回 Engine
func Setup(mode string) *gin.Engine {
	gin.SetMode(mode)

	r := gin.Default()

	// 全局中间件
	r.Use(middleware.CORS())

	// 路由组
	api := r.Group("/api")
	{
		api.GET("/hello", handler.Hello)
		api.GET("/ping", handler.Ping)

		api.POST("/redis/set", handler.RedisSet)
		api.GET("/redis/get", handler.RedisGet)
		api.POST("/redis/stock/deduct", handler.StockDeduct)

		api.GET("/user", handler.UserList)
		api.PUT("/user", handler.UserUpdate)
	}

	return r
}
