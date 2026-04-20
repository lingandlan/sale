package handler

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"marketplace/backend/internal/model"
	"marketplace/backend/internal/rbac"
	"marketplace/backend/internal/service"
	apperrors "marketplace/backend/pkg/errors"
	"marketplace/backend/pkg/errmsg"
	"marketplace/backend/pkg/response"
)

// AdminHandler 管理员处理器
type AdminHandler struct {
	userSvc   service.UserServiceInterface
	casbinSvc *rbac.CasbinService
}

// NewAdminHandler 创建 AdminHandler
func NewAdminHandler(userSvc service.UserServiceInterface, casbinSvc *rbac.CasbinService) *AdminHandler {
	return &AdminHandler{
		userSvc:   userSvc,
		casbinSvc: casbinSvc,
	}
}

// ListUsers 获取用户列表（管理员）
func (h *AdminHandler) ListUsers(c *gin.Context) {
	var req model.ListUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ParamsError(c, err.Error())
		return
	}

	result, err := h.userSvc.ListWithFilters(c.Request.Context(), &req)
	if err != nil {
		response.InternalError(c, errmsg.Get("admin.user_list_failed"))
		return
	}

	response.Success(c, result)
}

// ResetPassword 重置用户密码（管理员）
func (h *AdminHandler) ResetPassword(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.ParamsError(c, errmsg.Get("user.id_format_error"))
		return
	}

	var req model.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamsError(c, err.Error())
		return
	}

	if err := h.userSvc.ResetPassword(c.Request.Context(), id, &req); err != nil {
		if err == apperrors.ErrNotFound {
			response.NotFound(c, errmsg.Get("admin.user_not_found"))
		} else {
			response.InternalError(c, errmsg.Get("admin.reset_password"))
		}
		return
	}

	response.Success(c, nil)
}

// UpdateUserStatus 启用/禁用用户（管理员）
func (h *AdminHandler) UpdateUserStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.ParamsError(c, errmsg.Get("user.id_format_error"))
		return
	}

	var req model.UpdateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamsError(c, err.Error())
		return
	}

	if err := h.userSvc.UpdateStatus(c.Request.Context(), id, &req); err != nil {
		if err == apperrors.ErrNotFound {
			response.NotFound(c, errmsg.Get("admin.user_not_found"))
		} else {
			response.InternalError(c, errmsg.Get("admin.update_status"))
		}
		return
	}

	// 获取更新后的用户信息
	user, err := h.userSvc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.InternalError(c, errmsg.Get("admin.get_user_failed"))
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
		response.InternalError(c, errmsg.Get("admin.create_user_failed"))
		return
	}

	// 同步 Casbin 角色
	if h.casbinSvc != nil && req.Role != "" {
		_ = h.casbinSvc.AddRoleForUser(fmt.Sprintf("%d", user.ID), req.Role)
	}

	response.Success(c, user)
}

// UpdateUser 更新用户（管理员）
func (h *AdminHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.ParamsError(c, errmsg.Get("user.id_format_error"))
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
			response.NotFound(c, errmsg.Get("admin.user_not_found"))
		} else {
			response.InternalError(c, errmsg.Get("admin.update_user_failed"))
		}
		return
	}

	// 角色变更时同步 Casbin
	if h.casbinSvc != nil && req.Role != nil {
		_ = h.casbinSvc.UpdateUserRole(fmt.Sprintf("%d", id), *req.Role)
	}

	response.Success(c, user)
}

// DeleteUser 删除用户（管理员）
func (h *AdminHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.ParamsError(c, errmsg.Get("user.id_format_error"))
		return
	}

	if err := h.userSvc.Delete(c.Request.Context(), id); err != nil {
		if err == apperrors.ErrNotFound {
			response.NotFound(c, errmsg.Get("admin.user_not_found"))
		} else {
			response.InternalError(c, errmsg.Get("admin.delete_user_failed"))
		}
		return
	}

	response.Success(c, nil)
}
