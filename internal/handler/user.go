package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/wxsimon8888/simonStu/internal/logger"
	"github.com/wxsimon8888/simonStu/internal/repository"
	"github.com/wxsimon8888/simonStu/internal/response"
)

type userItem struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	IsAdmin   bool   `json:"is_admin"`
	CreatedAt string `json:"create_time"`
}

type UserUpdateRequest struct {
	ID       int    `json:"id" binding:"required"`
	Username string `json:"username"`
	IsAdmin  *bool  `json:"is_admin"`
}

// UserUpdate 修改用户信息（仅允许修改 username、is_admin）
func UserUpdate(c *gin.Context) {
	var req UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf(c, "UserUpdate 参数解析失败: %v", err)
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	user, err := repository.UserRepo.GetByID(req.ID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			logger.Errorf(c, "UserUpdate 用户不存在 id=%d", req.ID)
			response.Error(c, http.StatusNotFound, "用户不存在")
		} else {
			logger.Errorf(c, "UserUpdate 查询用户失败 id=%d err=%v", req.ID, err)
			response.Error(c, http.StatusInternalServerError, "查询用户失败: "+err.Error())
		}
		return
	}

	updates := map[string]interface{}{}
	if req.Username != "" {
		updates["username"] = req.Username
	}
	if req.IsAdmin != nil {
		updates["is_admin"] = *req.IsAdmin
	}

	if len(updates) == 0 {
		logger.Errorf(c, "UserUpdate 没有可更新的字段 id=%d", req.ID)
		response.Error(c, http.StatusBadRequest, "没有可更新的字段")
		return
	}

	if err := repository.UserRepo.Update(req.ID, updates); err != nil {
		logger.Errorf(c, "UserUpdate 更新失败 id=%d updates=%v err=%v", req.ID, updates, err)
		response.Error(c, http.StatusInternalServerError, "更新失败: "+err.Error())
		return
	}

	response.Success(c, userItem{
		ID:        user.ID,
		Username:  user.Username,
		IsAdmin:   user.IsAdmin,
		CreatedAt: user.CreateTime.Format(time.DateTime),
	})
}
