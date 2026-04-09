package model

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        int64          `gorm:"primaryKey;autoIncrement" json:"id" db:"id"`
	Username  string         `gorm:"size:64;uniqueIndex;not null" json:"username" db:"username"`
	Password  string         `gorm:"size:255;not null" json:"-" db:"password"`
	Nickname  string         `gorm:"size:64" json:"nickname" db:"nickname"`
	Email     string         `gorm:"size:128;uniqueIndex" json:"email" db:"email"`
	Avatar    sql.NullString `gorm:"size:512" json:"avatar" db:"avatar"`
	Role      int            `gorm:"default:0;comment:'0:普通用户 1:商家 2:管理员'" json:"role" db:"role"`
	Status    int            `gorm:"default:1;comment:'0:禁用 1:正常'" json:"status" db:"status"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at" db:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at" db:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-" db:"deleted_at"`
}

// TableName 表名
func (User) TableName() string {
	return "users"
}

// UserRole 用户角色
const (
	RoleUser     = 0 // 普通用户
	RoleMerchant = 1 // 商家
	RoleAdmin    = 2 // 管理员
)

// UserStatus 用户状态
const (
	UserStatusDisabled = 0 // 禁用
	UserStatusNormal   = 1 // 正常
)

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         *User  `json:"user"`
}

// RefreshTokenRequest 刷新 Token 请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
	Nickname string `json:"nickname" binding:"max=64"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Nickname *string `json:"nickname" binding:"omitempty,max=64"`
	Email    *string `json:"email" binding:"omitempty,email"`
	Avatar   *string `json:"avatar" binding:"omitempty,max=512"`
}
