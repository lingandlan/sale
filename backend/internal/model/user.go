package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型（太积堂系统）
type User struct {
	ID           int64          `gorm:"primaryKey;autoIncrement" json:"id" db:"id"`
	Username     string         `gorm:"size:50;uniqueIndex:uk_username;not null;comment:用户名（登录账号）" json:"username" db:"username"`
	Phone        string         `gorm:"size:20;uniqueIndex:uk_phone;not null;comment:手机号" json:"phone" db:"phone"`
	Password     string         `gorm:"size:255;not null;comment:密码（bcrypt加密）" json:"-" db:"password"`
	Name         string         `gorm:"size:50;not null;comment:姓名" json:"name" db:"name"`
	Role         string         `gorm:"size:20;default:'operator';comment:角色：super_admin=超管, admin=管理员, operator=操作员" json:"role" db:"role"`
	CenterID     *uint          `gorm:"comment:所属充值中心ID" json:"center_id" db:"center_id"`
	CenterName   *string        `gorm:"size:100;comment:所属充值中心名称（冗余字段）" json:"center_name" db:"center_name"`
	Status       int8           `gorm:"default:1;comment:状态：1=启用, 0=禁用" json:"status" db:"status"`
	LastLoginAt  *time.Time     `gorm:"comment:最后登录时间" json:"last_login_at" db:"last_login_at"`
	LastLoginIP  *string        `gorm:"size:50;comment:最后登录IP" json:"last_login_ip" db:"last_login_ip"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at" db:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at" db:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-" db:"deleted_at"`
}

// TableName 表名
func (User) TableName() string {
	return "users"
}

// UserRole 用户角色（太积堂系统）
const (
	RoleSuperAdmin  = "super_admin"  // 超级管理员
	RoleHQAdmin     = "hq_admin"     // 总部管理员
	RoleFinance     = "finance"      // 财务运营
	RoleCenterAdmin = "center_admin" // 充值中心管理员
	RoleOperator    = "operator"     // 操作员
)

// UserStatus 用户状态
const (
	UserStatusDisabled = 0 // 禁用
	UserStatusNormal   = 1 // 启用
)

// LoginRequest 登录请求
type LoginRequest struct {
	Phone    string `json:"phone" binding:"required,len=11"`
	Password string `json:"password" binding:"required,min=6,max=32"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

// RefreshTokenRequest 刷新 Token 请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username  string `json:"username" binding:"required,min=2,max=50"`
	Phone     string `json:"phone" binding:"required,len=11"`
	Password  string `json:"password" binding:"required,min=6,max=32"`
	Name      string `json:"name" binding:"required,min=2,max=50"`
	Role      string `json:"role" binding:"required,oneof=super_admin hq_admin finance center_admin operator"`
	CenterID  *uint  `json:"center_id" binding:"omitempty,min=1"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Username   *string `json:"username" binding:"omitempty,min=2,max=50"`
	Name       *string `json:"name" binding:"omitempty,min=2,max=50"`
	Phone      *string `json:"phone" binding:"omitempty,len=11"`
	Role       *string `json:"role" binding:"omitempty,oneof=super_admin hq_admin finance center_admin operator"`
	CenterID   *uint   `json:"center_id" binding:"omitempty,min=1"`
	CenterName *string `json:"center_name" binding:"omitempty,max=100"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=6"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=32"`
}

// ResetPasswordRequest 重置密码请求（管理员）
type ResetPasswordRequest struct {
	NewPassword string `json:"new_password" binding:"required,min=6,max=32"`
}

// UpdateUserStatusRequest 更新用户状态请求（管理员）
type UpdateUserStatusRequest struct {
	Status *int8 `json:"status" binding:"required,oneof=0 1"`
}

// ListUsersRequest 获取用户列表请求（管理员）
type ListUsersRequest struct {
	Page      int    `form:"page" binding:"omitempty,min=1"`
	PageSize  int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Keyword   string `form:"keyword" binding:"omitempty,max=50"`
	Role      string `form:"role" binding:"omitempty,oneof=super_admin hq_admin finance center_admin operator"`
	Status    *int8  `form:"status" binding:"omitempty,oneof=0 1"`
}

// ListUsersResponse 获取用户列表响应（管理员）
type ListUsersResponse struct {
	Items    []*User `json:"items"`
	Total    int64   `json:"total"`
	Page     int     `json:"page"`
	PageSize int     `json:"page_size"`
}
