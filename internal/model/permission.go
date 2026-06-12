package model

// Permission 对应表 c_permissions
type Permission struct {
	BaseModel
	Name        string `gorm:"column:name;type:varchar(128);not null;uniqueIndex"`
	Description string `gorm:"column:description;type:varchar(255);default:''"`
	ParentID    *int   `gorm:"column:parent_id"`
}

func (Permission) TableName() string { return "c_permissions" }
