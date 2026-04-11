package service

import (
	"context"

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

// GetByPhone 根据手机号获取用户
func (s *UserService) GetByPhone(ctx context.Context, phone string) (*model.User, error) {
	user, err := s.userRepo.GetByPhone(ctx, phone)
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
		Phone:    req.Phone,
		Password: hashedPassword,
		Name:     req.Name,
		Role:     req.Role,
		CenterID: req.CenterID,
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

	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Phone != nil {
		user.Phone = *req.Phone
	}
	if req.Role != nil {
		user.Role = *req.Role
	}
	if req.CenterID != nil {
		user.CenterID = req.CenterID
	}
	if req.CenterName != nil {
		user.CenterName = req.CenterName
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

// ListWithFilters 带筛选条件的用户列表（管理员）
func (s *UserService) ListWithFilters(ctx context.Context, req *model.ListUsersRequest) (*model.ListUsersResponse, error) {
	// 设置默认值
	page := req.Page
	if page <= 0 {
		page = 1
	}

	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	users, total, err := s.userRepo.ListWithFilters(ctx, page, pageSize, req.Keyword, req.Role, req.Status)
	if err != nil {
		return nil, err
	}

	return &model.ListUsersResponse{
		Items:    users,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// ResetPassword 重置用户密码（管理员）
func (s *UserService) ResetPassword(ctx context.Context, id int64, req *model.ResetPasswordRequest) error {
	// 检查用户是否存在
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return apperrors.ErrNotFound
	}

	// 哈希新密码
	hashedPassword, err := HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	// 更新密码
	return s.userRepo.UpdatePassword(ctx, user.ID, hashedPassword)
}

// UpdateStatus 更新用户状态（管理员）
func (s *UserService) UpdateStatus(ctx context.Context, id int64, req *model.UpdateUserStatusRequest) error {
	// 检查用户是否存在
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return apperrors.ErrNotFound
	}

	// 更新状态
	return s.userRepo.UpdateStatus(ctx, user.ID, req.Status)
}

// Delete 删除用户（管理员，软删除）
func (s *UserService) Delete(ctx context.Context, id int64) error {
	// 检查用户是否存在
	_, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return apperrors.ErrNotFound
	}

	// 软删除
	return s.userRepo.Delete(ctx, id)
}

// UpdatePassword 更新密码
func (s *UserService) UpdatePassword(ctx context.Context, id int64, hashedPassword string) error {
	return s.userRepo.UpdatePassword(ctx, id, hashedPassword)
}
