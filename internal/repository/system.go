package repository

import (
	"github.com/wxsimon8888/simonStu/internal/database"
	"github.com/wxsimon8888/simonStu/internal/model"
)

var RoleRepo = NewBaseRepo[model.Role](database.DB)
var PermissionRepo = NewBaseRepo[model.Permission](database.DB)
