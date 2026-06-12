package model

// Role 对应表 c_roles
type Role struct {
	BaseModel
	Name        string `gorm:"column:name;type:varchar(64);not null;uniqueIndex"`
	Description string `gorm:"column:description;type:varchar(255);default:''"`
	Status      int    `gorm:"column:status;type:tinyint;not null;default:1"`
}

func (Role) TableName() string { return "c_roles" }
