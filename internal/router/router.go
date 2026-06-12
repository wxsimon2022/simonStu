// Package router 路由注册。组装所有中间件和 handler。
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
		// 公开接口（无需登录）
		api.POST("/auth/login", handler.Login)
		api.POST("/auth/admin/login", handler.AdminLogin)

		// 需登录
		auth := api.Group("")
		auth.Use(middleware.AuthRequired())
		{
			auth.GET("/auth/userinfo", handler.UserInfo)

			auth.GET("/user", middleware.PermissionRequired("user:list"), handler.UserList)
			auth.PUT("/user", middleware.PermissionRequired("user:update"), handler.UserUpdate)
		}

		// 基础接口（无需权限）
		api.GET("/hello", handler.Hello)
		api.GET("/ping", handler.Ping)

		api.POST("/redis/set", handler.RedisSet)
		api.GET("/redis/get", handler.RedisGet)
		api.POST("/redis/stock/deduct", handler.StockDeduct)

		api.POST("/concurrent/process", handler.ConcurrentProcess)
		api.GET("/concurrent/fetch", handler.ConcurrentMultiFetch)
	}

	return r
}
