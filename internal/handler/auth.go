// Package handler 认证 handler。登录、注销、获取用户信息。
package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/wxsimon8888/simonStu/internal/database"
	"github.com/wxsimon8888/simonStu/internal/logger"
	"github.com/wxsimon8888/simonStu/internal/model"
	"github.com/wxsimon8888/simonStu/internal/response"
	"github.com/wxsimon8888/simonStu/internal/service"
)

// Auth 全局认证服务实例，在 main 中初始化。
var Auth *service.AuthService

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login 用户登录。
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请输入用户名和密码")
		return
	}
	if database.DB == nil {
		logger.Errorf(c, "Login 数据库未连接")
		response.Error(c, http.StatusInternalServerError, "数据库未连接")
		return
	}

	var user model.Users
	if err := database.DB.Where("username = ? AND is_delete = 0", req.Username).First(&user).Error; err != nil {
		logger.Errorf(c, "Login 用户不存在 username=%s", req.Username)
		response.Error(c, http.StatusUnauthorized, "用户名或密码错误")
		return
	}
	if !service.CheckPassword(req.Password, user.PasswordHash) {
		logger.Errorf(c, "Login 密码错误 username=%s", req.Username)
		response.Error(c, http.StatusUnauthorized, "用户名或密码错误")
		return
	}

	token, err := Auth.GenerateToken(user.ID, user.Username, user.IsAdmin)
	if err != nil {
		logger.Errorf(c, "Login 生成令牌失败: %v", err)
		response.Error(c, http.StatusInternalServerError, "生成令牌失败")
		return
	}

	// 令牌存入 Redis
	Auth.StoreToken(context.Background(), user.ID, token, 24*time.Hour)

	logger.Infof(c, "Login 成功 username=%s", req.Username)
	response.Success(c, gin.H{
		"token": token,
		"user":  gin.H{"id": user.ID, "username": user.Username, "is_admin": user.IsAdmin},
	})
}

// AdminLogin 管理员登录。
func AdminLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请输入用户名和密码")
		return
	}
	if database.DB == nil {
		logger.Errorf(c, "AdminLogin 数据库未连接")
		response.Error(c, http.StatusInternalServerError, "数据库未连接")
		return
	}

	var admin model.Admin
	if err := database.DB.Where("username = ? AND is_delete = 0", req.Username).First(&admin).Error; err != nil {
		logger.Errorf(c, "AdminLogin 管理员不存在 username=%s", req.Username)
		response.Error(c, http.StatusUnauthorized, "用户名或密码错误")
		return
	}
	if admin.Status == 0 {
		logger.Errorf(c, "AdminLogin 账号已禁用 username=%s", req.Username)
		response.Error(c, http.StatusForbidden, "账号已被禁用")
		return
	}
	if !service.CheckPassword(req.Password, admin.PasswordHash) {
		logger.Errorf(c, "AdminLogin 密码错误 username=%s", req.Username)
		response.Error(c, http.StatusUnauthorized, "用户名或密码错误")
		return
	}

	token, err := Auth.GenerateToken(admin.ID, admin.Username, true)
	if err != nil {
		logger.Errorf(c, "AdminLogin 生成令牌失败: %v", err)
		response.Error(c, http.StatusInternalServerError, "生成令牌失败")
		return
	}

	// 令牌存入 Redis
	Auth.StoreToken(context.Background(), admin.ID, token, 24*time.Hour)

	// 缓存当前用户权限到 Redis，后续权限校验优先走缓存
	var permList []string
	database.DB.Table("c_permissions").
		Select("DISTINCT c_permissions.name").
		Joins("JOIN c_role_permissions ON c_permissions.id = c_role_permissions.permission_id").
		Joins("JOIN c_admin_roles ON c_role_permissions.role_id = c_admin_roles.role_id").
		Where("c_admin_roles.user_id = ? AND c_permissions.is_delete = 0 AND c_role_permissions.is_delete = 0 AND c_admin_roles.is_delete = 0", admin.ID).
		Pluck("c_permissions.name", &permList)
	if len(permList) > 0 {
		Auth.CachePermissions(context.Background(), admin.ID, permList, 24*time.Hour)
	}

	logger.Infof(c, "AdminLogin 成功 username=%s", req.Username)
	response.Success(c, gin.H{
		"token": token,
		"user":  gin.H{"id": admin.ID, "username": admin.Username, "real_name": admin.RealName},
	})
}

// Logout 退出登录，从 Redis 撤销令牌。
func Logout(c *gin.Context) {
	tokenStr, _ := c.Get("token")
	if tokenStr != nil {
		Auth.RevokeToken(c.Request.Context(), tokenStr.(string))
	}
	logger.Infof(c, "Logout 成功")
	response.Success(c, nil)
}

// UserInfo 获取当前登录用户信息（含角色、权限）。
func UserInfo(c *gin.Context) {
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")

	var roles []string
	database.DB.Table("c_roles").
		Select("c_roles.name").
		Joins("JOIN c_admin_roles ON c_roles.id = c_admin_roles.role_id").
		Where("c_admin_roles.user_id = ? AND c_roles.is_delete = 0 AND c_admin_roles.is_delete = 0", userID).
		Pluck("c_roles.name", &roles)

	var perms []string
	database.DB.Table("c_permissions").
		Select("DISTINCT c_permissions.name").
		Joins("JOIN c_role_permissions ON c_permissions.id = c_role_permissions.permission_id").
		Joins("JOIN c_admin_roles ON c_role_permissions.role_id = c_admin_roles.role_id").
		Where("c_admin_roles.user_id = ? AND c_permissions.is_delete = 0 AND c_role_permissions.is_delete = 0 AND c_admin_roles.is_delete = 0", userID).
		Pluck("c_permissions.name", &perms)

	response.Success(c, gin.H{
		"user_id":     userID,
		"username":    username,
		"roles":       roles,
		"permissions": perms,
	})
}
