package repository

import (
	"github.com/wxsimon8888/simonStu/internal/database"
	"github.com/wxsimon8888/simonStu/internal/model"
)

// UserRepo 用户 CRUD 实例。
var UserRepo = NewBaseRepo[model.Users](database.DB)

// 扩展示例：
// func FindByUsername(db *gorm.DB, username string) (*model.Users, error) { ... }
