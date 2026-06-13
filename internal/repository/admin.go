// Package repository 管理员仓储。
package repository

import (
	"github.com/wxsimon2022/simonStu/internal/database"
	"github.com/wxsimon2022/simonStu/internal/model"
)

// AdminRepo 管理员 CRUD 实例。
var AdminRepo = NewBaseRepo[model.Admin](database.DB)
