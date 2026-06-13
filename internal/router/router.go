// Package router 路由注册入口，按功能拆分为独立文件：
//
//	public.go — 无需登录的公开接口（登录、Redis、并发示例）
//	system.go — 系统管理接口（管理员、角色、权限），需 JWT + 权限校验
package router

import (
	"github.com/gin-gonic/gin"

	"github.com/wxsimon2022/simonStu/internal/handler"
	"github.com/wxsimon2022/simonStu/internal/middleware"
)

// Setup 创建 Gin 引擎实例，注册全局中间件和全部路由分组。
//
// 路由结构：
//
//	GET  /api/hello               健康检查
//	GET  /api/ping                连通性测试
//	POST /api/auth/login          用户登录
//	POST /api/auth/admin/login    管理员登录
//	...
//	[以下需 JWT 认证]
//	GET  /api/auth/userinfo       当前用户信息
//	POST /api/auth/logout         退出登录
//	PUT  /api/user                修改用户信息
//	[系统管理路由 → RegisterSystemRoutes]
func Setup(mode string) *gin.Engine {
	gin.SetMode(mode)
	r := gin.Default()
	r.SetTrustedProxies(nil)
	r.Use(middleware.CORS())
	r.Static("/uploads", "storage/uploads")

	// 公共接口（无需登录）
	commonApi := r.Group("/commonApi")
	RegisterPublicRoutes(commonApi)

	// 需登录接口（保持 /api 前缀）
	api := r.Group("/api")
	{
		// POST /auth/login — 普通用户登录
		api.POST("/auth/login", handler.Login)
		// POST /auth/admin/login — 管理后台管理员登录
		api.POST("/auth/admin/login", handler.AdminLogin)

		// 认证路由组：需携带有效 JWT 令牌
		auth := api.Group("")
		auth.Use(middleware.AuthRequired())
		{
			// GET /auth/userinfo — 获取当前登录用户信息（角色、权限列表）
			auth.GET("/auth/userinfo", handler.UserInfo)
			// POST /auth/logout — 退出登录（从 Redis 吊销令牌）
			auth.POST("/auth/logout", handler.Logout)
			// GET /auth/profile — 获取当前管理员个人信息
			auth.GET("/auth/profile", handler.SystemAdminProfile)
			// PUT /auth/profile — 修改当前管理员个人信息
			auth.PUT("/auth/profile", handler.SystemAdminProfileUpdate)
			// GET /dashboard/stats — 仪表盘统计数据
			auth.GET("/dashboard/stats", handler.DashboardStats)

			// PUT /user — 修改普通用户信息（需 user:update 权限）
			auth.PUT("/user", middleware.PermissionRequired("user:update"), handler.UserUpdate)

			// 系统管理路由组：管理员、角色、菜单 CRUD（每个路由配独立权限校验）
			RegisterSystemRoutes(auth)

		}
	}

	return r
}
