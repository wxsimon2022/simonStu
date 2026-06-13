package repository

import (
	"github.com/wxsimon2022/simonStu/internal/database"
	"github.com/wxsimon2022/simonStu/internal/model"
)

var RoleRepo = NewBaseRepo[model.Role](database.DB)
var PermissionRepo = NewBaseRepo[model.Permission](database.DB)
