// Package middleware JWT 认证 + 权限校验中间件。
package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/wxsimon8888/simonStu/internal/database"
	"github.com/wxsimon8888/simonStu/internal/handler"
	"github.com/wxsimon8888/simonStu/internal/response"
)

// AuthRequired 验证 JWT 令牌与 Redis Token 有效性。
// 通过后将 token / user_id / username / is_admin 写入 context。
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			response.Error(c, http.StatusUnauthorized, "未登录")
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := handler.Auth.ParseToken(tokenStr)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "令牌无效或已过期")
			c.Abort()
			return
		}

		valid, err := handler.Auth.ValidateToken(c.Request.Context(), tokenStr)
		if err != nil || !valid {
			response.Error(c, http.StatusUnauthorized, "令牌已失效，请重新登录")
			c.Abort()
			return
		}

		c.Set("token", tokenStr)
		c.Set("user_id", int(claims["user_id"].(float64)))
		c.Set("username", claims["username"].(string))
		c.Set("is_admin", claims["is_admin"].(bool))
		c.Next()
	}
}

// PermissionRequired 校验当前用户是否拥有指定权限（管理员跳过）。
// PermissionRequired 校验当前用户是否拥有指定权限（管理员跳过，优先走 Redis 缓存）。
func PermissionRequired(perm string) gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, _ := c.Get("is_admin")
		if isAdmin.(bool) {
			c.Next()
			return
		}

		userID, _ := c.Get("user_id")
		uid := userID.(int)

		// 1. 优先从 Redis 缓存校验
		cached, err := handler.Auth.HasPermissionCache(c.Request.Context(), uid)
		if err == nil && cached {
			has, _ := handler.Auth.HasCachedPermission(c.Request.Context(), uid, perm)
			if has {
				c.Next()
			} else {
				response.Error(c, http.StatusForbidden, "无权限")
				c.Abort()
			}
			return
		}

		// 2. 缓存未命中，回查数据库并缓存
		var perms []string
		database.DB.Table("c_permissions").
			Select("DISTINCT c_permissions.name").
			Joins("JOIN c_role_permissions ON c_permissions.id = c_role_permissions.permission_id").
			Joins("JOIN c_admin_roles ON c_role_permissions.role_id = c_admin_roles.role_id").
			Where("c_admin_roles.user_id = ? AND c_permissions.is_delete = 0 AND c_role_permissions.is_delete = 0 AND c_admin_roles.is_delete = 0", uid).
			Pluck("c_permissions.name", &perms)

		// 缓存到 Redis（即使无权限也缓存空集合，避免每次回查 DB）
		handler.Auth.CachePermissions(c.Request.Context(), uid, perms, 24*time.Hour)

		for _, p := range perms {
			if p == perm {
				c.Next()
				return
			}
		}
		response.Error(c, http.StatusForbidden, "无权限")
		c.Abort()
	}
}
