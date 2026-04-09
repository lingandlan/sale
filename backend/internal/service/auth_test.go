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

type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *MockRedisClient) Del(ctx context.Context, keys ...string) error {
	args := m.Called(ctx, keys)
	return args.Error(0)
}

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

func (m *MockAuthUserRepo) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	args := m.Called(ctx, username)
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

func TestAuthService_GenerateToken(t *testing.T) {
	cfg := &config.JWTConfig{
		Secret:             "test-secret",
		ExpireHours:        1,
		RefreshExpireHours: 24,
	}
	svc := NewAuthService(cfg, nil, nil)
	user := &model.User{ID: 1, Username: "testuser", Role: model.RoleUser}

	t.Run("success", func(t *testing.T) {
		token, err := svc.GenerateToken(user)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		claims, err := svc.ParseToken(token)
		require.NoError(t, err)
		assert.Equal(t, user.ID, claims.UserID)
		assert.Equal(t, user.Username, claims.Username)
	})
}

func TestAuthService_ParseToken(t *testing.T) {
	cfg := &config.JWTConfig{
		Secret:             "test-secret",
		ExpireHours:        1,
		RefreshExpireHours: 24,
	}
	svc := NewAuthService(cfg, nil, nil)
	user := &model.User{ID: 1, Username: "testuser", Role: model.RoleUser}

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
	user := &model.User{ID: 1, Username: "testuser", Role: model.RoleUser}

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

func TestClaims_ClaimsType(t *testing.T) {
	cfg := &config.JWTConfig{
		Secret:             "test-secret",
		ExpireHours:        1,
		RefreshExpireHours: 24,
	}
	svc := NewAuthService(cfg, nil, nil)
	user := &model.User{ID: 1, Username: "admin", Role: model.RoleAdmin}

	token, err := svc.GenerateToken(user)
	require.NoError(t, err)

	claims, err := svc.ParseToken(token)
	require.NoError(t, err)

	assert.Equal(t, int64(1), claims.UserID)
	assert.Equal(t, "admin", claims.Username)
	assert.Equal(t, model.RoleAdmin, claims.Role)

	_, ok := interface{}(claims).(jwt.Claims)
	assert.True(t, ok)
}

type MockRedisForAuth struct {
	mock.Mock
}

func (m *MockRedisForAuth) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockRedisForAuth) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *MockRedisForAuth) Del(ctx context.Context, keys ...string) error {
	args := m.Called(ctx, keys)
	return args.Error(0)
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
		mockRepo.On("GetByUsername", ctx, "nonexistent").Return(nil, errors.New("sql: no rows")).Once()

		resp, err := svc.Login(ctx, "nonexistent", "password")

		assert.Error(t, err)
		assert.Nil(t, resp)
		mockRepo.AssertExpectations(t)
	})

	t.Run("user disabled", func(t *testing.T) {
		disabledUser := &model.User{
			ID:       1,
			Username: "disabled",
			Password: "$2a$10$hashedpassword",
			Status:   model.UserStatusDisabled,
		}
		mockRepo.On("GetByUsername", ctx, "disabled").Return(disabledUser, nil).Once()

		resp, err := svc.Login(ctx, "disabled", "password")

		assert.Error(t, err)
		assert.Nil(t, resp)
		mockRepo.AssertExpectations(t)
	})
}
