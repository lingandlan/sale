package service

import (
	"context"

	"marketplace/backend/internal/model"
)

type UserServiceInterface interface {
	GetByID(ctx context.Context, id int64) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	Create(ctx context.Context, req *model.CreateUserRequest) (*model.User, error)
	Update(ctx context.Context, id int64, req *model.UpdateUserRequest) (*model.User, error)
	List(ctx context.Context, page, pageSize int) ([]*model.User, int64, error)
}

var _ UserServiceInterface = (*UserService)(nil)

type AuthServiceInterface interface {
	Login(ctx context.Context, username, password string) (*model.LoginResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
	Logout(ctx context.Context, userID int64) error
	GenerateToken(user *model.User) (string, error)
	GenerateRefreshToken(user *model.User) (string, error)
	ParseToken(tokenString string) (*Claims, error)
}

var _ AuthServiceInterface = (*AuthService)(nil)
