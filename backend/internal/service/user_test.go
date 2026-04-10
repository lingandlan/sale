package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"marketplace/backend/internal/model"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) GetByID(ctx context.Context, id int64) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepo) GetByPhone(ctx context.Context, phone string) (*model.User, error) {
	args := m.Called(ctx, phone)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepo) Create(ctx context.Context, user *model.User) (int64, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepo) Update(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) UpdatePassword(ctx context.Context, id int64, password string) error {
	args := m.Called(ctx, id, password)
	return args.Error(0)
}

func (m *MockUserRepo) List(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]*model.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepo) ListWithFilters(ctx context.Context, page, pageSize int, keyword, role string, status *int8) ([]*model.User, int64, error) {
	args := m.Called(ctx, page, pageSize, keyword, role, status)
	return args.Get(0).([]*model.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepo) UpdateStatus(ctx context.Context, id int64, status int8) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockUserRepo) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestUserService_GetByID(t *testing.T) {
	mockRepo := new(MockUserRepo)
	svc := NewUserService(mockRepo)
	ctx := context.Background()
	expectedUser := &model.User{ID: 1, Phone: "13800138000", Name: "Test User"}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, int64(1)).Return(expectedUser, nil).Once()

		result, err := svc.GetByID(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("sql: no rows")).Once()

		result, err := svc.GetByID(ctx, 999)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_GetByPhone(t *testing.T) {
	mockRepo := new(MockUserRepo)
	svc := NewUserService(mockRepo)
	ctx := context.Background()
	expectedUser := &model.User{ID: 1, Phone: "13800138000", Name: "Test User"}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetByPhone", ctx, "13800138000").Return(expectedUser, nil).Once()

		result, err := svc.GetByPhone(ctx, "13800138000")

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("GetByPhone", ctx, "13900000000").Return(nil, errors.New("sql: no rows")).Once()

		result, err := svc.GetByPhone(ctx, "13900000000")

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_Create(t *testing.T) {
	mockRepo := new(MockUserRepo)
	svc := NewUserService(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo.On("Create", ctx, mock.AnythingOfType("*model.User")).Return(int64(1), nil).Once()

		req := &model.CreateUserRequest{
			Phone:    "13800138000",
			Password: "password123",
			Name:     "Test User",
			Role:     model.RoleOperator,
		}
		result, err := svc.Create(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(1), result.ID)
		assert.Equal(t, "13800138000", result.Phone)
		assert.Equal(t, model.RoleOperator, result.Role)
		assert.Equal(t, int8(model.UserStatusNormal), result.Status)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_Update(t *testing.T) {
	mockRepo := new(MockUserRepo)
	svc := NewUserService(mockRepo)
	ctx := context.Background()
	existingUser := &model.User{ID: 1, Phone: "13800138000", Name: "Old Name"}
	newName := "New Name"

	t.Run("success update name", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, int64(1)).Return(existingUser, nil).Once()
		mockRepo.On("Update", ctx, mock.AnythingOfType("*model.User")).Return(nil).Once()

		req := &model.UpdateUserRequest{Name: &newName}
		result, err := svc.Update(ctx, 1, req)

		assert.NoError(t, err)
		assert.Equal(t, "New Name", result.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("sql: no rows")).Once()

		req := &model.UpdateUserRequest{Name: &newName}
		result, err := svc.Update(ctx, 999, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_List(t *testing.T) {
	mockRepo := new(MockUserRepo)
	svc := NewUserService(mockRepo)
	ctx := context.Background()
	expectedUsers := []*model.User{
		{ID: 1, Phone: "13800138000"},
		{ID: 2, Phone: "13800138001"},
	}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("List", ctx, 1, 10).Return(expectedUsers, int64(2), nil).Once()

		users, total, err := svc.List(ctx, 1, 10)

		assert.NoError(t, err)
		assert.Equal(t, 2, len(users))
		assert.Equal(t, int64(2), total)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_Delete(t *testing.T) {
	mockRepo := new(MockUserRepo)
	svc := NewUserService(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, int64(1)).Return(&model.User{ID: 1}, nil).Once()
		mockRepo.On("Delete", ctx, int64(1)).Return(nil).Once()

		err := svc.Delete(ctx, 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("sql: no rows")).Once()

		err := svc.Delete(ctx, 999)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
