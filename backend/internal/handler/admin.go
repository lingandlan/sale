package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"marketplace/backend/internal/model"
	"marketplace/backend/internal/service"
	apperrors "marketplace/backend/pkg/errors"
	"marketplace/backend/pkg/response"
)

// AdminHandler 管理员处理器
type AdminHandler struct {
	userSvc service.UserServiceInterface
}

// NewAdminHandler 创建 AdminHandler
func NewAdminHandler(userSvc service.UserServiceInterface) *AdminHandler {
	return &AdminHandler{
		userSvc: userSvc,
	}
}

// ListUsers 获取用户列表（管理员）
// @Summary 获取用户列表（管理员）
// @Tags 用户管理
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param keyword query string false "搜索关键词"
// @Param role query string false "角色筛选"
// @Param status query int false "状态筛选"
// @Success 200 {object} response.Response{data=model.ListUsersResponse}
// @Router /api/v1/admin/users [get]
func (h *AdminHandler) ListUsers(c *gin.Context) {
	var req model.ListUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ParamsError(c, err.Error())
		return
	}

	result, err := h.userSvc.ListWithFilters(c.Request.Context(), &req)
	if err != nil {
		response.InternalError(c, "获取用户列表失败")
		return
	}

	response.Success(c, result)
}

// ResetPassword 重置用户密码（管理员）
// @Summary 重置用户密码（管理员）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param request body model.ResetPasswordRequest true "请求体"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/users/{id}/reset-password [post]
func (h *AdminHandler) ResetPassword(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.ParamsError(c, "用户ID格式错误")
		return
	}

	var req model.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamsError(c, err.Error())
		return
	}

	if err := h.userSvc.ResetPassword(c.Request.Context(), id, &req); err != nil {
		if err == apperrors.ErrNotFound {
			response.NotFound(c, "用户不存在")
		} else {
			response.InternalError(c, "重置密码失败")
		}
		return
	}

	response.Success(c, nil)
}

// UpdateUserStatus 启用/禁用用户（管理员）
// @Summary 启用/禁用用户（管理员）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param request body model.UpdateUserStatusRequest true "请求体"
// @Success 200 {object} response.Response{data=model.User}
// @Router /api/v1/admin/users/{id}/status [put]
func (h *AdminHandler) UpdateUserStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.ParamsError(c, "用户ID格式错误")
		return
	}

	var req model.UpdateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamsError(c, err.Error())
		return
	}

	if err := h.userSvc.UpdateStatus(c.Request.Context(), id, &req); err != nil {
		if err == apperrors.ErrNotFound {
			response.NotFound(c, "用户不存在")
		} else {
			response.InternalError(c, "更新用户状态失败")
		}
		return
	}

	// 获取更新后的用户信息
	user, err := h.userSvc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, "获取用户信息失败")
		return
	}

	response.Success(c, user)
}

// CreateUser 创建用户（管理员）
func (h *AdminHandler) CreateUser(c *gin.Context) {
	var req model.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamsError(c, err.Error())
		return
	}

	user, err := h.userSvc.Create(c.Request.Context(), &req)
	if err != nil {
		response.InternalError(c, "创建用户失败")
		return
	}

	response.Success(c, user)
}

// UpdateUser 更新用户（管理员）
func (h *AdminHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.ParamsError(c, "用户ID格式错误")
		return
	}

	var req model.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamsError(c, err.Error())
		return
	}

	user, err := h.userSvc.Update(c.Request.Context(), id, &req)
	if err != nil {
		if err == apperrors.ErrNotFound {
			response.NotFound(c, "用户不存在")
		} else {
			response.InternalError(c, "更新用户失败")
		}
		return
	}

	response.Success(c, user)
}

// DeleteUser 删除用户（管理员）
// @Summary 删除用户（管理员）
// @Tags 用户管理
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} response.Response
// @Router /api/v1/admin/users/{id} [delete]
func (h *AdminHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.ParamsError(c, "用户ID格式错误")
		return
	}

	if err := h.userSvc.Delete(c.Request.Context(), id); err != nil {
		if err == apperrors.ErrNotFound {
			response.NotFound(c, "用户不存在")
		} else {
			response.InternalError(c, "删除用户失败")
		}
		return
	}

	response.Success(c, nil)
}
