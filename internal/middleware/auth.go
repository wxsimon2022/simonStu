// JWT 认证中间件。从 Authorization 头提取并验证令牌，将用户信息注入请求上下文。
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/wxsimon8888/simonStu/internal/database"
	"github.com/wxsimon8888/simonStu/internal/handler"
	"github.com/wxsimon8888/simonStu/internal/response"
)

// AuthRequired 验证 JWT 令牌，通过后将 user_id / username / is_admin 写入 context。
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			response.Error(c, http.StatusUnauthorized, "未登录")
			c.Abort()
			return
		}

		claims, err := handler.Auth.ParseToken(strings.TrimPrefix(authHeader, "Bearer "))
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "令牌无效或已过期")
			c.Abort()
			return
		}

		userID := int(claims["user_id"].(float64))
		username := claims["username"].(string)
		isAdmin := claims["is_admin"].(bool)

		c.Set("user_id", userID)
		c.Set("username", username)
		c.Set("is_admin", isAdmin)
		c.Next()
	}
}

// PermissionRequired 校验当前用户是否拥有指定权限。依赖 AuthRequired 先执行。
func PermissionRequired(perm string) gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, _ := c.Get("is_admin")
		if isAdmin.(bool) {
			c.Next() // 管理员跳过校验
			return
		}

		userID, _ := c.Get("user_id")
		var count int64
		database.DB.Table("c_admin_roles").
			Joins("JOIN c_role_permissions ON c_admin_roles.role_id = c_role_permissions.role_id").
			Joins("JOIN c_permissions ON c_role_permissions.permission_id = c_permissions.id").
			Where("c_admin_roles.user_id = ? AND c_permissions.name = ?", userID, perm).
			Count(&count)

		if count == 0 {
			response.Error(c, http.StatusForbidden, "无权限")
			c.Abort()
			return
		}
		c.Next()
	}
}
