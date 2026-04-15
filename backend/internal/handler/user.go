package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"marketplace/backend/internal/middleware"
	"marketplace/backend/internal/model"
	"marketplace/backend/internal/service"
	"marketplace/backend/pkg/errmsg"
	"marketplace/backend/pkg/response"
)

// UserHandler 用户处理器
type UserHandler struct {
	userSvc    service.UserServiceInterface
	memberSvc  service.MemberServiceInterface
}

// NewUserHandler 创建 UserHandler
func NewUserHandler(userSvc service.UserServiceInterface, memberSvc service.MemberServiceInterface) *UserHandler {
	return &UserHandler{userSvc: userSvc, memberSvc: memberSvc}
}

// GetUserInfo 获取当前用户信息
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	userID := middleware.GetUserID(c)

	user, err := h.userSvc.GetByID(c.Request.Context(), userID)
	if err != nil {
		response.NotFound(c, errmsg.Get("user.not_found"))
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
		response.InternalError(c, errmsg.Get("user.update_failed"))
		return
	}

	response.Success(c, user)
}

// GetUserByID 根据 ID 获取用户
func (h *UserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.ParamsError(c, errmsg.Get("user.id_format_error"))
		return
	}

	user, err := h.userSvc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, errmsg.Get("user.not_found"))
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
		response.InternalError(c, errmsg.Get("user.list_failed"))
		return
	}

	response.Success(c, response.ListResponse{
		Items:    users,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// ChangePassword 修改密码
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req model.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamsError(c, err.Error())
		return
	}

	// 验证旧密码（需要先获取用户）
	user, err := h.userSvc.GetByID(c.Request.Context(), userID)
	if err != nil {
		response.NotFound(c, errmsg.Get("user.not_found"))
		return
	}

	// TODO: 验证旧密码
	// 更新密码
	hashedPassword, err := service.HashPassword(req.NewPassword)
	if err != nil {
		response.InternalError(c, errmsg.Get("user.password_encrypt"))
		return
	}

	if err := h.userSvc.UpdatePassword(c.Request.Context(), user.ID, hashedPassword); err != nil {
		response.InternalError(c, errmsg.Get("user.password_change"))
		return
	}

	response.Success(c, nil)
}

// SearchMember 根据手机号查询商城会员信息
func (h *UserHandler) SearchMember(c *gin.Context) {
	phone := c.Query("phone")
	if phone == "" {
		response.ParamsError(c, "手机号不能为空")
		return
	}

	member, err := h.memberSvc.SearchByPhone(phone)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, member)
}
