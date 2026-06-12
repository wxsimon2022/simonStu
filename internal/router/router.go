// Package router 路由注册。按功能拆分到独立文件：
//
//	public.go  — 无需登录的公开接口（登录、Redis、并发示例）
//	system.go — 系统管理接口（管理员、角色、权限），需 JWT + 权限校验
package router

import (
	"github.com/gin-gonic/gin"

	"github.com/wxsimon8888/simonStu/internal/handler"
	"github.com/wxsimon8888/simonStu/internal/middleware"
)

// Setup 创建 *gin.Engine，注册全部路由并返回。
func Setup(mode string) *gin.Engine {
	gin.SetMode(mode)
	r := gin.Default()
	r.Use(middleware.CORS())

	api := r.Group("/api")
	{
		// 公开接口（public.go）
		RegisterPublicRoutes(api)

		// 需 JWT 认证的接口
		auth := api.Group("")
		auth.Use(middleware.AuthRequired())
		{
			auth.GET("/auth/userinfo", handler.UserInfo)
			auth.POST("/auth/logout", handler.Logout)

			// 普通用户管理
			auth.PUT("/user", middleware.PermissionRequired("user:update"), handler.UserUpdate)

			// 系统管理（system.go）
			RegisterSystemRoutes(auth)
		}
	}
	return r
}
