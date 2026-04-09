package handler

import (
	"github.com/gin-gonic/gin"

	"marketplace/backend/internal/middleware"
	"marketplace/backend/internal/model"
	"marketplace/backend/internal/service"
	apperrors "marketplace/backend/pkg/errors"
	"marketplace/backend/pkg/response"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authSvc service.AuthServiceInterface
	userSvc service.UserServiceInterface
}

// NewAuthHandler 创建 AuthHandler
func NewAuthHandler(authSvc service.AuthServiceInterface, userSvc service.UserServiceInterface) *AuthHandler {
	return &AuthHandler{
		authSvc: authSvc,
		userSvc: userSvc,
	}
}

// Login 登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamsError(c, err.Error())
		return
	}

	result, err := h.authSvc.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		switch err {
		case apperrors.ErrNotFound:
			response.Unauthorized(c, "用户名或密码错误")
		case apperrors.ErrUserDisabled:
			response.Forbidden(c, "用户已被禁用")
		default:
			response.InternalError(c, "登录失败")
		}
		return
	}

	response.Success(c, result)
}

// Refresh 刷新 Token
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req model.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamsError(c, err.Error())
		return
	}

	accessToken, refreshToken, err := h.authSvc.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		response.Unauthorized(c, "token 无效或已过期")
		return
	}

	response.Success(c, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// Logout 登出
func (h *AuthHandler) Logout(c *gin.Context) {
	userID := middleware.GetUserID(c)

	if err := h.authSvc.Logout(c.Request.Context(), userID); err != nil {
		response.InternalError(c, "登出失败")
		return
	}

	response.SuccessWithMessage(c, "logout success", nil)
}

// Register 注册
func (h *AuthHandler) Register(c *gin.Context) {
	var req model.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamsError(c, err.Error())
		return
	}

	user, err := h.userSvc.Create(c.Request.Context(), &req)
	if err != nil {
		response.InternalError(c, "注册失败")
		return
	}

	response.Created(c, user)
}
