package service

import (
	"context"
	"database/sql"

	"marketplace/backend/internal/model"
	"marketplace/backend/internal/repository"
	apperrors "marketplace/backend/pkg/errors"
)

// UserService 用户服务
type UserService struct {
	userRepo repository.UserRepoInterface
}

// NewUserService 创建 UserService
func NewUserService(userRepo repository.UserRepoInterface) *UserService {
	return &UserService{userRepo: userRepo}
}

// GetByID 根据 ID 获取用户
func (s *UserService) GetByID(ctx context.Context, id int64) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}
	return user, nil
}

// GetByUsername 根据用户名获取用户
func (s *UserService) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}
	return user, nil
}

// Create 创建用户
func (s *UserService) Create(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
	// 密码哈希
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
		Nickname: req.Nickname,
		Role:     model.RoleUser,
		Status:   model.UserStatusNormal,
	}

	id, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	user.ID = id
	return user, nil
}

// Update 更新用户
func (s *UserService) Update(ctx context.Context, id int64, req *model.UpdateUserRequest) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}

	if req.Nickname != nil {
		user.Nickname = *req.Nickname
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Avatar != nil {
		user.Avatar = sql.NullString{String: *req.Avatar, Valid: true}
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// List 分页获取用户列表
func (s *UserService) List(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	return s.userRepo.List(ctx, page, pageSize)
}
