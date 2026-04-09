package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"marketplace/backend/internal/middleware"
	"marketplace/backend/internal/model"
	"marketplace/backend/internal/service"
	"marketplace/backend/pkg/response"
)

// UserHandler 用户处理器
type UserHandler struct {
	userSvc service.UserServiceInterface
}

// NewUserHandler 创建 UserHandler
func NewUserHandler(userSvc service.UserServiceInterface) *UserHandler {
	return &UserHandler{userSvc: userSvc}
}

// GetUserInfo 获取当前用户信息
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	userID := middleware.GetUserID(c)

	user, err := h.userSvc.GetByID(c.Request.Context(), userID)
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	response.Success(c, user)
}

// UpdateUserInfo 更新当前用户信息
func (h *UserHandler) UpdateUserInfo(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req model.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamsError(c, err.Error())
		return
	}

	user, err := h.userSvc.Update(c.Request.Context(), userID, &req)
	if err != nil {
		response.InternalError(c, "更新失败")
		return
	}

	response.Success(c, user)
}

// GetUserByID 根据 ID 获取用户
func (h *UserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.ParamsError(c, "invalid user id")
		return
	}

	user, err := h.userSvc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	response.Success(c, user)
}

// ListUsers 获取用户列表（管理员）
func (h *UserHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	users, total, err := h.userSvc.List(c.Request.Context(), page, pageSize)
	if err != nil {
		response.InternalError(c, "获取列表失败")
		return
	}

	response.Success(c, response.ListResponse{
		Items:    users,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}
