// Package router 系统管理路由（管理员、角色、权限）。
package router

import (
	"github.com/gin-gonic/gin"

	"github.com/wxsimon8888/simonStu/internal/handler"
	"github.com/wxsimon8888/simonStu/internal/middleware"
)

// RegisterSystemRoutes 注册系统管理相关路由。
func RegisterSystemRoutes(auth *gin.RouterGroup) {
	auth.GET("/admin", middleware.PermissionRequired("admin:list"), handler.SystemAdminList)
	auth.POST("/admin", middleware.PermissionRequired("admin:create"), handler.SystemAdminCreate)
	auth.PUT("/admin", middleware.PermissionRequired("admin:update"), handler.SystemAdminUpdate)
	auth.DELETE("/admin/:id", middleware.PermissionRequired("admin:delete"), handler.SystemAdminDelete)

	auth.GET("/role", middleware.PermissionRequired("role:list"), handler.SystemRoleList)
	auth.POST("/role", middleware.PermissionRequired("role:create"), handler.SystemRoleCreate)
	auth.PUT("/role", middleware.PermissionRequired("role:update"), handler.SystemRoleUpdate)
	auth.DELETE("/role/:id", middleware.PermissionRequired("role:delete"), handler.SystemRoleDelete)
	auth.PUT("/role/permissions", middleware.PermissionRequired("role:update"), handler.SystemRolePermissions)

	auth.GET("/permission", middleware.PermissionRequired("perm:list"), handler.SystemPermissionList)
	auth.POST("/permission", middleware.PermissionRequired("perm:create"), handler.SystemPermissionCreate)
	auth.PUT("/permission", middleware.PermissionRequired("perm:update"), handler.SystemPermissionUpdate)
	auth.DELETE("/permission/:id", middleware.PermissionRequired("perm:delete"), handler.SystemPermissionDelete)
}
