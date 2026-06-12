// Package repository 通用仓储层。基于 Go 泛型提供 CRUD 基础操作。
package repository

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// ErrNotFound 标记查询结果为空。
var ErrNotFound = errors.New("记录不存在")

// BaseRepo 通用 CRUD 操作。T 为任意 model 类型。
type BaseRepo[T any] struct {
	DB *gorm.DB
}

// NewBaseRepo 创建泛型仓储实例。
func NewBaseRepo[T any](db *gorm.DB) *BaseRepo[T] {
	return &BaseRepo[T]{DB: db}
}

// GetByID 根据主键查询单条记录（不含已删除）。未找到时返回 ErrNotFound。
func (r *BaseRepo[T]) GetByID(id int) (*T, error) {
	if r.DB == nil {
		return nil, errors.New("数据库未连接")
	}
	var m T
	if err := r.DB.Where("is_delete = 0").First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &m, nil
}

// List 分页列表，按 id DESC 排序，只返回未删除记录。
func (r *BaseRepo[T]) List(page, size int) ([]T, int64, error) {
	if r.DB == nil {
		return nil, 0, errors.New("数据库未连接")
	}
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	var list []T
	var total int64
	r.DB.Model(new(T)).Where("is_delete = 0").Count(&total)
	r.DB.Where("is_delete = 0").Order("id DESC").Offset((page - 1) * size).Limit(size).Find(&list)
	return list, total, nil
}

// Create 创建单条记录。
func (r *BaseRepo[T]) Create(data *T) error {
	if r.DB == nil {
		return errors.New("数据库未连接")
	}
	return r.DB.Create(data).Error
}

// Update 按 map 更新指定 ID 的记录（仅更新未删除的）。未找到时返回 ErrNotFound。
func (r *BaseRepo[T]) Update(id int, updates map[string]interface{}) error {
	if r.DB == nil {
		return errors.New("数据库未连接")
	}
	updates["update_time"] = time.Now()
	result := r.DB.Model(new(T)).Where("id = ? AND is_delete = 0", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// Delete 软删除（设置 is_delete = 1）。
func (r *BaseRepo[T]) Delete(id int) error {
	if r.DB == nil {
		return errors.New("数据库未连接")
	}
	result := r.DB.Model(new(T)).Where("id = ?", id).Updates(map[string]interface{}{
		"is_delete":   1,
		"update_time": time.Now(),
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
