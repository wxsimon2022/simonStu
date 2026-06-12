package repository

import (
	"github.com/wxsimon8888/simonStu/internal/database"
	"github.com/wxsimon8888/simonStu/internal/model"
)

// 全局用户仓储实例
var UserRepo = NewBaseRepo[model.Users](database.DB)

// 还可以在这里加用户专属的查询方法
// func FindByUsername(db *gorm.DB, username string) (*model.Users, error) { ... }
