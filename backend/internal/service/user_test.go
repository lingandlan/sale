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

func (m *MockUserRepo) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	args := m.Called(ctx, username)
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

func (m *MockUserRepo) List(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]*model.User), args.Get(1).(int64), args.Error(2)
}

func TestUserService_GetByID(t *testing.T) {
	mockRepo := new(MockUserRepo)
	svc := NewUserService(mockRepo)
	ctx := context.Background()
	expectedUser := &model.User{ID: 1, Username: "testuser"}

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

func TestUserService_GetByUsername(t *testing.T) {
	mockRepo := new(MockUserRepo)
	svc := NewUserService(mockRepo)
	ctx := context.Background()
	expectedUser := &model.User{ID: 1, Username: "testuser"}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetByUsername", ctx, "testuser").Return(expectedUser, nil).Once()

		result, err := svc.GetByUsername(ctx, "testuser")

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("GetByUsername", ctx, "nonexistent").Return(nil, errors.New("sql: no rows")).Once()

		result, err := svc.GetByUsername(ctx, "nonexistent")

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
			Username: "newuser",
			Password: "password123",
			Email:    "new@test.com",
			Nickname: "New User",
		}
		result, err := svc.Create(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(1), result.ID)
		assert.Equal(t, "newuser", result.Username)
		assert.Equal(t, model.RoleUser, result.Role)
		assert.Equal(t, model.UserStatusNormal, result.Status)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_Update(t *testing.T) {
	mockRepo := new(MockUserRepo)
	svc := NewUserService(mockRepo)
	ctx := context.Background()
	existingUser := &model.User{ID: 1, Username: "testuser", Nickname: "Old Nick"}
	newNickname := "New Nickname"

	t.Run("success update nickname", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, int64(1)).Return(existingUser, nil).Once()
		mockRepo.On("Update", ctx, mock.AnythingOfType("*model.User")).Return(nil).Once()

		req := &model.UpdateUserRequest{Nickname: &newNickname}
		result, err := svc.Update(ctx, 1, req)

		assert.NoError(t, err)
		assert.Equal(t, "New Nickname", result.Nickname)
		mockRepo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("sql: no rows")).Once()

		req := &model.UpdateUserRequest{Nickname: &newNickname}
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
		{ID: 1, Username: "user1"},
		{ID: 2, Username: "user2"},
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
