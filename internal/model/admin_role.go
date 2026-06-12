package model

import "time"

// AdminRole 对应表 c_admin_roles
type AdminRole struct {
	ID         int       `gorm:"column:id;primaryKey;autoIncrement"`
	UserID     int       `gorm:"column:user_id;not null;uniqueIndex:uk_admin_role"`
	RoleID     int       `gorm:"column:role_id;not null;uniqueIndex:uk_admin_role"`
	IsDelete   int       `gorm:"column:is_delete;type:tinyint;not null;default:0"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime"`
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime"`
}

func (AdminRole) TableName() string { return "c_admin_roles" }

// RolePermission 对应表 c_role_permissions
type RolePermission struct {
	ID           int       `gorm:"column:id;primaryKey;autoIncrement"`
	RoleID       int       `gorm:"column:role_id;not null;uniqueIndex:uk_role_perm"`
	PermissionID int       `gorm:"column:permission_id;not null;uniqueIndex:uk_role_perm"`
	IsDelete     int       `gorm:"column:is_delete;type:tinyint;not null;default:0"`
	CreateTime   time.Time `gorm:"column:create_time;autoCreateTime"`
	UpdateTime   time.Time `gorm:"column:update_time;autoUpdateTime"`
}

func (RolePermission) TableName() string { return "c_role_permissions" }
