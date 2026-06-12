// Package model 数据模型。定义数据库表对应的 Go 结构体。
package model

import "time"

// BaseModel 公共基础字段，所有业务模型嵌入此结构。
// 数据库字段统一使用 is_delete / create_time / update_time，无外键。
type BaseModel struct {
	ID         int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	IsDelete   int       `gorm:"column:is_delete;type:tinyint;not null;default:0" json:"-"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime" json:"-"`
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime" json:"-"`
}
