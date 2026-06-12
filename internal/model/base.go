package model

import (
	"time"
)

// BaseModel 公共基础字段，其他 model 通过嵌入继承
type BaseModel struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement" json:"-"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"-"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"-"`
}
