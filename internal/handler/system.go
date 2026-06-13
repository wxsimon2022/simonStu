// Package handler 系统管理 — 管理员、角色、权限 CRUD。
package handler

import (
	"net/http"
	"strconv"
	"strings"

	"time"

	"github.com/gin-gonic/gin"

	"github.com/wxsimon8888/simonStu/internal/database"
	"github.com/wxsimon8888/simonStu/internal/logger"
	"github.com/wxsimon8888/simonStu/internal/model"
	"github.com/wxsimon8888/simonStu/internal/response"
	"github.com/wxsimon8888/simonStu/internal/service"
)

// ========================= 管理员 CRUD =========================

type adminItem struct {
	ID         int      `json:"id"`
	Username   string   `json:"username"`
	RealName   string   `json:"real_name"`
	Phone      string   `json:"phone"`
	Email      string   `json:"email"`
	Status     int      `json:"status"`
	CreateTime string   `json:"create_time"`
	Roles      []string `json:"roles"`
	RoleIDs    []int    `json:"role_ids"`
}

type AdminCreateRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	RealName string `json:"real_name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Status   *int   `json:"status"`
	RoleIDs  []int  `json:"role_ids"`
}

type AdminUpdateRequest struct {
	ID       int    `json:"id" binding:"required"`
	Password string `json:"password"`
	RealName string `json:"real_name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Status   *int   `json:"status"`
	RoleIDs  *[]int `json:"role_ids"`
}

// SystemAdminList 管理员列表（含角色）
func SystemAdminList(c *gin.Context) {
	if database.DB == nil {
		response.Error(c, http.StatusInternalServerError, "数据库未连接")
		return
	}

	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "15"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 15
	}
	offset := (page - 1) * pageSize

	// 查询总数
	var total int64
	database.DB.Table("c_admin").Where("is_delete = 0").Count(&total)

	// 分页扫描原始数据（time.Time 类型，后续用 Go 逻辑格式化）
	type adminDBItem struct {
		ID         int
		Username   string
		RealName   string
		Phone      string
		Email      string
		Status     int
		CreateTime time.Time
	}
	var dbItems []adminDBItem
	database.DB.Table("c_admin").
		Select("id, username, real_name, phone, email, status, create_time").
		Where("is_delete = 0").Order("id DESC").
		Limit(pageSize).Offset(offset).Find(&dbItems)

	// 转换为响应结构，Go 逻辑格式化时间（去掉时区）
	list := make([]adminItem, len(dbItems))
	for i, item := range dbItems {
		list[i] = adminItem{
			ID:         item.ID,
			Username:   item.Username,
			RealName:   item.RealName,
			Phone:      item.Phone,
			Email:      item.Email,
			Status:     item.Status,
			CreateTime: item.CreateTime.Format("2006-01-02 15:04:05"),
		}
	}

	// 一次性查所有管理员的角色
	type userRole struct {
		UserID   int
		RoleName string
		RoleID   int
	}
	var rows []userRole
	database.DB.Table("c_admin_roles").
		Select("c_admin_roles.user_id, c_roles.name AS role_name, c_roles.id AS role_id").
		Joins("JOIN c_roles ON c_admin_roles.role_id = c_roles.id").
		Where("c_admin_roles.is_delete = 0 AND c_roles.is_delete = 0").
		Scan(&rows)

	roleNames := map[int][]string{}
	roleIDs := map[int][]int{}
	for _, r := range rows {
		roleNames[r.UserID] = append(roleNames[r.UserID], r.RoleName)
		roleIDs[r.UserID] = append(roleIDs[r.UserID], r.RoleID)
	}
	for i := range list {
		list[i].Roles = roleNames[list[i].ID]
		list[i].RoleIDs = roleIDs[list[i].ID]
	}
	response.Success(c, gin.H{"list": list, "total": total})
}

// SystemAdminCreate 创建管理员（可选关联角色）
func SystemAdminCreate(c *gin.Context) {
	var req AdminCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	hash, err := service.HashPassword(req.Password)
	if err != nil {
		logger.Errorf(c, "SystemAdminCreate 密码哈希失败: %v", err)
		response.Error(c, http.StatusInternalServerError, "密码加密失败")
		return
	}
	status := 1
	if req.Status != nil {
		status = *req.Status
	}
	admin := model.Admin{
		Username:     req.Username,
		PasswordHash: hash,
		RealName:     req.RealName,
		Phone:        req.Phone,
		Email:        req.Email,
		Status:       status,
	}
	if err := database.DB.Create(&admin).Error; err != nil {
		logger.Errorf(c, "SystemAdminCreate 创建失败: %v", err)
		response.Error(c, http.StatusInternalServerError, "创建失败: "+err.Error())
		return
	}
	// 关联角色
	for _, roleID := range req.RoleIDs {
		database.DB.Create(&model.AdminRole{UserID: admin.ID, RoleID: roleID})
	}
	logger.Infof(c, "SystemAdminCreate 成功 username=%s", req.Username)
	response.Success(c, gin.H{"id": admin.ID})
}

// SystemAdminUpdate 修改管理员（含角色更新）
func SystemAdminUpdate(c *gin.Context) {
	var req AdminUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	updates := map[string]interface{}{}
	if req.Password != "" {
		hash, err := service.HashPassword(req.Password)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "密码加密失败")
			return
		}
		updates["password_hash"] = hash
	}
	if req.RealName != "" {
		updates["real_name"] = req.RealName
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if len(updates) > 0 {
		updates["update_time"] = time.Now()
		database.DB.Model(&model.Admin{}).Where("id = ? AND is_delete = 0", req.ID).Updates(updates)
	}
	// 更新角色关联
	if req.RoleIDs != nil {
		database.DB.Where("user_id = ?", req.ID).Delete(&model.AdminRole{})
		for _, roleID := range *req.RoleIDs {
			database.DB.Create(&model.AdminRole{UserID: req.ID, RoleID: roleID})
		}
	}
	logger.Infof(c, "SystemAdminUpdate 成功 id=%d", req.ID)
	response.Success(c, nil)
}

// SystemAdminDelete 软删除管理员
func SystemAdminDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效 ID")
		return
	}
	result := database.DB.Model(&model.Admin{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_delete": 1, "update_time": time.Now(),
	})
	if result.RowsAffected == 0 {
		response.Error(c, http.StatusNotFound, "管理员不存在")
		return
	}
	logger.Infof(c, "SystemAdminDelete 成功 id=%d", id)
	response.Success(c, nil)
}

// ========================= 角色 CRUD =========================

type roleItem struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Status      int      `json:"status"`
	Permissions []string `json:"permissions"`
}

type RoleCreateRequest struct {
	Name          string `json:"name" binding:"required"`
	Description   string `json:"description"`
	Status        *int   `json:"status"`
	PermissionIDs []int  `json:"permission_ids"`
}

type RoleUpdateRequest struct {
	ID            int    `json:"id" binding:"required"`
	Description   string `json:"description"`
	Status        *int   `json:"status"`
	PermissionIDs []int  `json:"permission_ids"`
}

// SystemRoleCreate 创建角色
func SystemRoleCreate(c *gin.Context) {
	var req RoleCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	role := model.Role{Name: req.Name, Description: req.Description}
	if req.Status != nil {
		role.Status = *req.Status
	}
	if err := database.DB.Create(&role).Error; err != nil {
		logger.Errorf(c, "SystemRoleCreate 失败: %v", err)
		response.Error(c, http.StatusInternalServerError, "创建失败: "+err.Error())
		return
	}
	for _, pid := range req.PermissionIDs {
		database.DB.Create(&model.RolePermission{RoleID: role.ID, PermissionID: pid})
	}
	logger.Infof(c, "SystemRoleCreate 成功 name=%s", req.Name)
	response.Success(c, gin.H{"id": role.ID})
}

// SystemRoleUpdate 更新角色（可同时更新权限）
func SystemRoleUpdate(c *gin.Context) {
	var req RoleUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	updates := map[string]interface{}{}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if len(updates) > 0 {
		updates["update_time"] = time.Now()
		database.DB.Model(&model.Role{}).Where("id = ? AND is_delete = 0", req.ID).Updates(updates)
	}
	if req.PermissionIDs != nil {
		database.DB.Where("role_id = ?", req.ID).Delete(&model.RolePermission{})
		for _, pid := range req.PermissionIDs {
			database.DB.Create(&model.RolePermission{RoleID: req.ID, PermissionID: pid})
		}
	}
	logger.Infof(c, "SystemRoleUpdate 成功 id=%d", req.ID)
	response.Success(c, nil)
}

// SystemRoleDelete 软删除角色
func SystemRoleDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效 ID")
		return
	}
	result := database.DB.Model(&model.Role{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_delete": 1, "update_time": time.Now(),
	})
	if result.RowsAffected == 0 {
		response.Error(c, http.StatusNotFound, "角色不存在")
		return
	}
	logger.Infof(c, "SystemRoleDelete 成功 id=%d", id)
	response.Success(c, nil)
}

// SystemRoleList 角色列表（含权限）
func SystemRoleList(c *gin.Context) {
	if database.DB == nil {
		response.Error(c, http.StatusInternalServerError, "数据库未连接")
		return
	}
	var roles []model.Role
	database.DB.Where("is_delete = 0").Order("id ASC").Find(&roles)

	type rpRow struct {
		RoleID         int
		PermissionName string
	}
	var rows []rpRow
	database.DB.Table("c_role_permissions").
		Select("c_role_permissions.role_id, c_permissions.name AS permission_name").
		Joins("JOIN c_permissions ON c_role_permissions.permission_id = c_permissions.id").
		Where("c_role_permissions.is_delete = 0 AND c_permissions.is_delete = 0").
		Find(&rows)

	permMap := map[int][]string{}
	for _, r := range rows {
		permMap[r.RoleID] = append(permMap[r.RoleID], r.PermissionName)
	}

	list := make([]roleItem, len(roles))
	for i, r := range roles {
		list[i] = roleItem{
			ID:          r.ID,
			Name:        r.Name,
			Description: r.Description,
			Status:      r.Status,
			Permissions: permMap[r.ID],
		}
	}
	response.Success(c, gin.H{"list": list})
}

// SystemRolePermissions 更新角色权限（旧接口，保留兼容）
func SystemRolePermissions(c *gin.Context) {
	var req struct {
		RoleID        int   `json:"role_id" binding:"required"`
		PermissionIDs []int `json:"permission_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	tx := database.DB.Begin()
	tx.Where("role_id = ?", req.RoleID).Delete(&model.RolePermission{})
	for _, pid := range req.PermissionIDs {
		tx.Create(&model.RolePermission{RoleID: req.RoleID, PermissionID: pid})
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		logger.Errorf(c, "SystemRolePermissions 更新失败 role_id=%d err=%v", req.RoleID, err)
		response.Error(c, http.StatusInternalServerError, "更新失败")
		return
	}
	logger.Infof(c, "SystemRolePermissions 更新成功 role_id=%d", req.RoleID)
	response.Success(c, nil)
}

// ========================= 权限 CRUD =========================

// ========================= 权限 CRUD =========================

// ========================= 权限（菜单）CRUD =========================

type PermissionTreeItem struct {
	ID          int                  `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Type        string               `json:"type"` // dir / menu / btn
	Children    []PermissionTreeItem `json:"children"`
}

type PermissionCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Type        string `json:"type"`
	ParentID    *int   `json:"parent_id"`
}

type PermissionUpdateRequest struct {
	ID          int    `json:"id" binding:"required"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	ParentID    *int   `json:"parent_id"`
}

// SystemPermissionCreate 创建权限
func SystemPermissionCreate(c *gin.Context) {
	if database.DB == nil {
		response.Error(c, http.StatusInternalServerError, "数据库未连接")
		return
	}
	var req PermissionCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	perm := model.Permission{
		Name: req.Name, Description: req.Description, Type: req.Type, ParentID: req.ParentID,
	}
	if err := database.DB.Create(&perm).Error; err != nil {
		logger.Errorf(c, "SystemPermissionCreate 失败: %v", err)
		response.Error(c, http.StatusInternalServerError, "创建失败: "+err.Error())
		return
	}
	logger.Infof(c, "SystemPermissionCreate 成功 name=%s", req.Name)
	response.Success(c, gin.H{"id": perm.ID})
}

// SystemPermissionUpdate 更新权限
func SystemPermissionUpdate(c *gin.Context) {
	if database.DB == nil {
		response.Error(c, http.StatusInternalServerError, "数据库未连接")
		return
	}
	var req PermissionUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Type != "" {
		updates["type"] = req.Type
	}
	if req.ParentID != nil {
		if *req.ParentID > 0 {
			updates["parent_id"] = *req.ParentID
		} else {
			updates["parent_id"] = nil
		}
	}
	if len(updates) == 0 {
		response.Error(c, http.StatusBadRequest, "没有可更新的字段")
		return
	}
	updates["update_time"] = time.Now()
	if err := database.DB.Model(&model.Permission{}).Where("id = ? AND is_delete = 0", req.ID).Updates(updates).Error; err != nil {
		logger.Errorf(c, "SystemPermissionUpdate 失败 id=%d err=%v", req.ID, err)
		response.Error(c, http.StatusInternalServerError, "更新失败: "+err.Error())
		return
	}
	logger.Infof(c, "SystemPermissionUpdate 成功 id=%d", req.ID)
	response.Success(c, nil)
}

// SystemPermissionDelete 软删除权限
func SystemPermissionDelete(c *gin.Context) {
	if database.DB == nil {
		response.Error(c, http.StatusInternalServerError, "数据库未连接")
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效 ID")
		return
	}
	result := database.DB.Model(&model.Permission{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_delete": 1, "update_time": time.Now(),
	})
	if result.RowsAffected == 0 {
		response.Error(c, http.StatusNotFound, "权限不存在")
		return
	}
	logger.Infof(c, "SystemPermissionDelete 成功 id=%d", id)
	response.Success(c, nil)
}

// SystemPermissionList 权限树（同时返回平铺列表供角色编辑使用）
func SystemPermissionList(c *gin.Context) {
	if database.DB == nil {
		response.Error(c, http.StatusInternalServerError, "数据库未连接")
		return
	}
	var flat []model.Permission
	database.DB.Where("is_delete = 0").Order("id ASC").Find(&flat)

	// 按 parent_id 分组
	group := map[int][]PermissionTreeItem{}
	for _, p := range flat {
		pid := 0
		if p.ParentID != nil {
			pid = *p.ParentID
		}
		group[pid] = append(group[pid], PermissionTreeItem{
			ID: p.ID, Name: p.Name, Description: p.Description, Type: p.Type,
		})
	}

	// 构建树
	var tree []PermissionTreeItem
	for _, p := range flat {
		if p.ParentID == nil {
			node := PermissionTreeItem{
				ID: p.ID, Name: p.Name, Description: p.Description, Type: p.Type,
			}
			node.Children = group[p.ID]
			tree = append(tree, node)
		}
	}
	if tree == nil {
		tree = []PermissionTreeItem{}
	}

	response.Success(c, gin.H{"tree": tree, "list": flat})
}

// MenuItem 菜单树节点
type MenuItem struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Type        string     `json:"type"`
	Route       string     `json:"route"`
	Children    []MenuItem `json:"children"`
}

// SystemMenuList 返回当前用户的菜单树（基于权限过滤）
func SystemMenuList(c *gin.Context) {
	if database.DB == nil {
		response.Error(c, http.StatusInternalServerError, "数据库未连接")
		return
	}
	isAdmin, _ := c.Get("is_admin")
	userID, _ := c.Get("user_id")

	var perms []model.Permission
	query := database.DB.Where("type IN ('dir', 'menu', 'btn') AND is_delete = 0").Order("id ASC")
	if !isAdmin.(bool) {
		query = query.Joins("JOIN c_role_permissions ON c_permissions.id = c_role_permissions.permission_id").
			Joins("JOIN c_admin_roles ON c_role_permissions.role_id = c_admin_roles.role_id").
			Where("c_admin_roles.user_id = ? AND c_admin_roles.is_delete = 0 AND c_role_permissions.is_delete = 0", userID)
	}
	query.Find(&perms)

	group := map[int][]MenuItem{}
	for _, p := range perms {
		route := ""
		if p.Type == "menu" {
			parts := strings.SplitN(p.Name, ":", 2)
			if len(parts) == 2 && parts[1] == "list" {
				route = "/" + parts[0]
			}
			if parts[0] == "perm" {
				route = "/permission"
			}
			if !strings.Contains(p.Name, ":") {
				switch p.Name {
				case "admin":
					route = "/admin"
				case "role":
					route = "/role"
				case "perm":
					route = "/permission"
				}
			}
		}
		pid := 0
		if p.ParentID != nil {
			pid = *p.ParentID
		}
		group[pid] = append(group[pid], MenuItem{
			ID: p.ID, Name: p.Name, Description: p.Description, Type: p.Type, Route: route,
		})
	}

	var tree []MenuItem
	for _, p := range perms {
		if p.ParentID == nil {
			node := MenuItem{
				ID: p.ID, Name: p.Name, Description: p.Description, Type: p.Type,
			}
			node.Children = group[p.ID]
			tree = append(tree, node)
		}
	}
	if tree == nil {
		tree = []MenuItem{}
	}

	response.Success(c, gin.H{"menus": tree})
}
