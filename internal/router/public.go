// Package router 公开路由（无需登录）。
package router

import (
	"github.com/gin-gonic/gin"

	"github.com/wxsimon8888/simonStu/internal/handler"
)

// RegisterPublicRoutes 注册公开路由。
func RegisterPublicRoutes(api *gin.RouterGroup) {
	api.POST("/auth/login", handler.Login)
	api.POST("/auth/admin/login", handler.AdminLogin)
	api.GET("/hello", handler.Hello)
	api.GET("/ping", handler.Ping)
	api.POST("/redis/set", handler.RedisSet)
	api.GET("/redis/get", handler.RedisGet)
	api.POST("/redis/stock/deduct", handler.StockDeduct)
	api.POST("/concurrent/process", handler.ConcurrentProcess)
	api.GET("/concurrent/fetch", handler.ConcurrentMultiFetch)
}
