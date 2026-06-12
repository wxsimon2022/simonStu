// 数据模型。定义数据库表对应的 Go 结构体。
package model

import "time"

// BaseModel 公共基础字段，嵌入到业务模型中统一管理主键和时间戳。
type BaseModel struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement" json:"-"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"-"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"-"`
	// DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"` // 软删除，按需开启
}
