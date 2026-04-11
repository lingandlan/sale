package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"marketplace/backend/internal/config"
	"marketplace/backend/internal/model"
)

type MockAuthUserRepo struct {
	mock.Mock
}

func (m *MockAuthUserRepo) GetByID(ctx context.Context, id int64) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockAuthUserRepo) GetByPhone(ctx context.Context, phone string) (*model.User, error) {
	args := m.Called(ctx, phone)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockAuthUserRepo) Create(ctx context.Context, user *model.User) (int64, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockAuthUserRepo) Update(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockAuthUserRepo) List(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]*model.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockAuthUserRepo) ListWithFilters(ctx context.Context, page, pageSize int, keyword, role string, status *int8) ([]*model.User, int64, error) {
	args := m.Called(ctx, page, pageSize, keyword, role, status)
	return args.Get(0).([]*model.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockAuthUserRepo) UpdateStatus(ctx context.Context, id int64, status int8) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockAuthUserRepo) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAuthUserRepo) UpdatePassword(ctx context.Context, id int64, password string) error {
	args := m.Called(ctx, id, password)
	return args.Error(0)
}

func TestAuthService_GenerateToken(t *testing.T) {
	cfg := &config.JWTConfig{
		Secret:             "test-secret",
		ExpireHours:        1,
		RefreshExpireHours: 24,
	}
	svc := NewAuthService(cfg, nil, nil)
	user := &model.User{ID: 1, Phone: "13800138000", Name: "Test User", Role: model.RoleOperator}

	t.Run("success", func(t *testing.T) {
		token, err := svc.GenerateToken(user)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		claims, err := svc.ParseToken(token)
		require.NoError(t, err)
		assert.Equal(t, user.ID, claims.UserID)
		assert.Equal(t, user.Phone, claims.Phone)
	})
}

func TestAuthService_ParseToken(t *testing.T) {
	cfg := &config.JWTConfig{
		Secret:             "test-secret",
		ExpireHours:        1,
		RefreshExpireHours: 24,
	}
	svc := NewAuthService(cfg, nil, nil)
	user := &model.User{ID: 1, Phone: "13800138000", Role: model.RoleOperator}

	t.Run("valid token", func(t *testing.T) {
		token, _ := svc.GenerateToken(user)
		claims, err := svc.ParseToken(token)

		assert.NoError(t, err)
		assert.Equal(t, user.ID, claims.UserID)
	})

	t.Run("invalid token", func(t *testing.T) {
		claims, err := svc.ParseToken("invalid-token")

		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	t.Run("token with wrong secret", func(t *testing.T) {
		otherCfg := &config.JWTConfig{
			Secret:             "other-secret",
			ExpireHours:        1,
			RefreshExpireHours: 24,
		}
		otherSvc := NewAuthService(otherCfg, nil, nil)
		token, _ := otherSvc.GenerateToken(user)

		claims, err := svc.ParseToken(token)

		assert.Error(t, err)
		assert.Nil(t, claims)
	})
}

func TestAuthService_GenerateRefreshToken(t *testing.T) {
	cfg := &config.JWTConfig{
		Secret:             "test-secret",
		ExpireHours:        1,
		RefreshExpireHours: 24,
	}
	svc := NewAuthService(cfg, nil, nil)
	user := &model.User{ID: 1, Phone: "13800138000", Role: model.RoleOperator}

	t.Run("success", func(t *testing.T) {
		token, err := svc.GenerateRefreshToken(user)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		claims, err := svc.ParseToken(token)
		require.NoError(t, err)
		assert.Equal(t, user.ID, claims.UserID)
	})
}

func TestHashPassword(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		hash, err := HashPassword("password123")

		assert.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.NotEqual(t, "password123", hash)
	})

	t.Run("different hashes for same password", func(t *testing.T) {
		hash1, _ := HashPassword("password123")
		hash2, _ := HashPassword("password123")

		assert.NotEqual(t, hash1, hash2)
	})
}

func TestClaims_Type(t *testing.T) {
	cfg := &config.JWTConfig{
		Secret:             "test-secret",
		ExpireHours:        1,
		RefreshExpireHours: 24,
	}
	svc := NewAuthService(cfg, nil, nil)
	user := &model.User{ID: 1, Phone: "13800138000", Name: "Admin", Role: model.RoleHQAdmin}

	token, err := svc.GenerateToken(user)
	require.NoError(t, err)

	claims, err := svc.ParseToken(token)
	require.NoError(t, err)

	assert.Equal(t, int64(1), claims.UserID)
	assert.Equal(t, "13800138000", claims.Phone)
	assert.Equal(t, model.RoleHQAdmin, claims.Role)

	_, ok := interface{}(claims).(jwt.Claims)
	assert.True(t, ok)
}

func TestAuthService_Login(t *testing.T) {
	cfg := &config.JWTConfig{
		Secret:             "test-secret",
		ExpireHours:        1,
		RefreshExpireHours: 24,
	}
	mockRepo := new(MockAuthUserRepo)
	svc := &AuthService{
		jwtSecret:          []byte(cfg.Secret),
		expireHours:        time.Duration(cfg.ExpireHours) * time.Hour,
		refreshExpireHours: time.Duration(cfg.RefreshExpireHours) * time.Hour,
		redis:              nil,
		userRepo:           mockRepo,
	}
	ctx := context.Background()

	t.Run("user not found", func(t *testing.T) {
		mockRepo.On("GetByPhone", ctx, "13900000000").Return(nil, errors.New("sql: no rows")).Once()

		resp, err := svc.Login(ctx, "13900000000", "password")

		assert.Error(t, err)
		assert.Nil(t, resp)
		mockRepo.AssertExpectations(t)
	})

	t.Run("user disabled", func(t *testing.T) {
		disabledUser := &model.User{
			ID:       1,
			Phone:    "13800138000",
			Password: "$2a$10$hashedpassword",
			Status:   model.UserStatusDisabled,
		}
		mockRepo.On("GetByPhone", ctx, "13800138000").Return(disabledUser, nil).Once()

		resp, err := svc.Login(ctx, "13800138000", "password")

		assert.Error(t, err)
		assert.Nil(t, resp)
		mockRepo.AssertExpectations(t)
	})
}
