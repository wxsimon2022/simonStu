package model

// Users 对应表 c_users
type Users struct {
	BaseModel
	Username     string `gorm:"column:username;type:varchar(64);not null;uniqueIndex:ix_c_users_username"`
	PasswordHash string `gorm:"column:password_hash;type:varchar(255);not null"`
	IsAdmin      bool   `gorm:"column:is_admin;type:tinyint(1);not null;default:0"`
}

// TableName 指定 GORM 的表名
func (Users) TableName() string {
	return "c_users"
}
