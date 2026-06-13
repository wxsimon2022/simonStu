package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/wxsimon2022/simonStu/internal/database"
	"github.com/wxsimon2022/simonStu/internal/response"
)

// DashboardStats 返回仪表盘统计数据（管理员/角色/权限总量）。
func DashboardStats(c *gin.Context) {
	if database.DB == nil {
		response.Error(c, http.StatusInternalServerError, "数据库未连接")
		return
	}

	type statItem struct {
		AdminCount int64 `json:"admin_count"`
		UserCount  int64 `json:"user_count"`
		RoleCount  int64 `json:"role_count"`
		PermCount  int64 `json:"perm_count"`
		MenuCount  int64 `json:"menu_count"`
	}

	var s statItem
	database.DB.Table("c_admin").Where("is_delete = 0").Select("COUNT(*)").Scan(&s.AdminCount)
	database.DB.Table("c_users").Where("is_delete = 0").Select("COUNT(*)").Scan(&s.UserCount)
	database.DB.Table("c_roles").Where("is_delete = 0").Select("COUNT(*)").Scan(&s.RoleCount)
	database.DB.Table("c_permissions").Where("is_delete = 0").Select("COUNT(*)").Scan(&s.PermCount)
	database.DB.Table("c_permissions").Where("is_delete = 0 AND type IN ('dir', 'menu')").Select("COUNT(*)").Scan(&s.MenuCount)

	response.Success(c, s)
}
