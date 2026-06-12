package model

// Admin 对应表 c_admin，管理后台登录用户。
type Admin struct {
	BaseModel
	Username     string `gorm:"column:username;type:varchar(64);not null;uniqueIndex:uk_admin_username"`
	PasswordHash string `gorm:"column:password_hash;type:varchar(255);not null"`
	RealName     string `gorm:"column:real_name;type:varchar(64);default:''"`
	Phone        string `gorm:"column:phone;type:varchar(20);default:''"`
	Email        string `gorm:"column:email;type:varchar(128);default:''"`
	Status       int    `gorm:"column:status;type:tinyint;not null;default:1"`
}

func (Admin) TableName() string { return "c_admin" }
