// Package router 系统管理路由：管理员、角色、菜单 CRUD，需 JWT + 权限校验。
package router

import (
	"github.com/gin-gonic/gin"

	"github.com/wxsimon2022/simonStu/internal/handler"
	"github.com/wxsimon2022/simonStu/internal/middleware"
)

// RegisterSystemRoutes 注册系统管理相关路由，每个路由通过 PermissionRequired 绑定权限标识。
//
// 管理员管理（c_admin）：
//
//	GET    /admin          管理员列表（需 admin:list 权限）
//	POST   /admin          创建管理员（需 admin:create 权限）
//	PUT    /admin          修改管理员（需 admin:update 权限）
//	DELETE /admin/:id      软删除管理员（需 admin:delete 权限）
//
// 角色管理（c_roles）：
//
//	GET    /role           角色列表（需 role:list 权限）
//	POST   /role           创建角色（需 role:create 权限）
//	PUT    /role           修改角色（需 role:update 权限）
//	DELETE /role/:id       软删除角色（需 role:delete 权限）
//	PUT    /role/permissions 更新角色权限关联（需 role:update 权限）
//
// 菜单/权限管理（c_permissions）：
//
//	GET    /permission      权限树 + 平铺列表（需 perm:list 权限）
//	POST   /permission      创建菜单/权限（需 perm:create 权限）
//	PUT    /permission      修改菜单/权限（需 perm:update 权限）
//	DELETE /permission/:id  软删除菜单/权限（需 perm:delete 权限）
//	PUT    /permission/reorder  批量更新菜单排序（需 perm:update 权限）
//
// 菜单导航：
//
//	GET    /menus           当前用户的菜单树（基于角色权限过滤，无需额外权限校验）
func RegisterSystemRoutes(auth *gin.RouterGroup) {
	// ======================== 管理员 ========================
	// GET /admin — 管理员列表
	auth.GET("/admin", middleware.PermissionRequired("admin:list"), handler.SystemAdminList)
	// POST /admin — 创建管理员
	auth.POST("/admin", middleware.PermissionRequired("admin:create"), handler.SystemAdminCreate)
	// PUT /admin — 修改管理员
	auth.PUT("/admin", middleware.PermissionRequired("admin:update"), handler.SystemAdminUpdate)
	// DELETE /admin/:id — 软删除管理员
	auth.DELETE("/admin/:id", middleware.PermissionRequired("admin:delete"), handler.SystemAdminDelete)

	// ======================== 角色 ========================
	// GET /role — 角色列表
	auth.GET("/role", middleware.PermissionRequired("role:list"), handler.SystemRoleList)
	// POST /role — 创建角色
	auth.POST("/role", middleware.PermissionRequired("role:create"), handler.SystemRoleCreate)
	// PUT /role — 修改角色
	auth.PUT("/role", middleware.PermissionRequired("role:update"), handler.SystemRoleUpdate)
	// DELETE /role/:id — 软删除角色
	auth.DELETE("/role/:id", middleware.PermissionRequired("role:delete"), handler.SystemRoleDelete)
	// PUT /role/permissions — 更新角色的权限关联
	auth.PUT("/role/permissions", middleware.PermissionRequired("role:update"), handler.SystemRolePermissions)

	// ======================== 菜单/权限 ========================
	// GET /permission — 权限树 + 平铺列表
	auth.GET("/permission", middleware.PermissionRequired("perm:list"), handler.SystemPermissionList)
	// POST /permission — 创建菜单
	auth.POST("/permission", middleware.PermissionRequired("perm:create"), handler.SystemPermissionCreate)
	// PUT /permission — 修改菜单
	auth.PUT("/permission", middleware.PermissionRequired("perm:update"), handler.SystemPermissionUpdate)
	// DELETE /permission/:id — 软删除菜单
	auth.DELETE("/permission/:id", middleware.PermissionRequired("perm:delete"), handler.SystemPermissionDelete)
	// PUT /permission/reorder — 批量更新菜单排序
	auth.PUT("/permission/reorder", middleware.PermissionRequired("perm:update"), handler.SystemPermissionReorder)
	// GET /menus — 当前用户的菜单树（侧边栏导航，无额外权限校验）
	auth.GET("/menus", handler.SystemMenuList)
}
