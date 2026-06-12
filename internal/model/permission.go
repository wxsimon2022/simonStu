package model

// Permission 对应表 c_permissions
type Permission struct {
	BaseModel
	Name        string `gorm:"column:name;type:varchar(128);not null;uniqueIndex" json:"name"`
	Description string `gorm:"column:description;type:varchar(255);default:''" json:"description"`
	ParentID    *int   `gorm:"column:parent_id" json:"parent_id"`
}

func (Permission) TableName() string { return "c_permissions" }
