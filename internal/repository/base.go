package repository

import (
	"errors"

	"gorm.io/gorm"
)

var ErrNotFound = errors.New("记录不存在")

// BaseRepo 通用 CRUD，T 为任意 model 类型
type BaseRepo[T any] struct {
	DB *gorm.DB
}

func NewBaseRepo[T any](db *gorm.DB) *BaseRepo[T] {
	return &BaseRepo[T]{DB: db}
}

// GetByID 根据主键查询单条
func (r *BaseRepo[T]) GetByID(id int) (*T, error) {
	if r.DB == nil {
		return nil, errors.New("数据库未连接")
	}
	var m T
	if err := r.DB.First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &m, nil
}

// List 分页列表，按 id DESC 排序
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
	r.DB.Model(new(T)).Count(&total)
	r.DB.Order("id DESC").Offset((page - 1) * size).Limit(size).Find(&list)
	return list, total, nil
}

// Create 创建记录
func (r *BaseRepo[T]) Create(data *T) error {
	if r.DB == nil {
		return errors.New("数据库未连接")
	}
	return r.DB.Create(data).Error
}

// Update 按 map 更新指定 ID 的记录
func (r *BaseRepo[T]) Update(id int, updates map[string]interface{}) error {
	if r.DB == nil {
		return errors.New("数据库未连接")
	}
	result := r.DB.Model(new(T)).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// Delete 软删除（依赖 BaseModel.DeletedAt）
func (r *BaseRepo[T]) Delete(id int) error {
	if r.DB == nil {
		return errors.New("数据库未连接")
	}
	result := r.DB.Delete(new(T), id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
