package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"marketplace/backend/internal/model"
	"marketplace/backend/internal/service"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Login(ctx context.Context, phone, password string) (*model.LoginResponse, error) {
	args := m.Called(ctx, phone, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.LoginResponse), args.Error(1)
}

func (m *MockAuthService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	args := m.Called(ctx, refreshToken)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockAuthService) Logout(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockAuthService) GenerateToken(user *model.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) GenerateRefreshToken(user *model.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) ParseToken(tokenString string) (*service.Claims, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.Claims), args.Error(1)
}

func setupAuthRouter(h *AuthHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/auth/login", h.Login)
	r.POST("/auth/refresh", h.Refresh)
	r.POST("/auth/register", h.Register)
	return r
}

func TestAuthHandler_Login(t *testing.T) {
	mockAuthSvc := new(MockAuthService)
	mockUserSvc := new(MockUserService)
	h := NewAuthHandler(mockAuthSvc, mockUserSvc)
	router := setupAuthRouter(h)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		expectedResp := &model.LoginResponse{
			AccessToken:  "access-token",
			RefreshToken: "refresh-token",
			ExpiresIn:    3600,
		}
		mockAuthSvc.On("Login", ctx, "13800138000", "password123").Return(expectedResp, nil).Once()

		body := map[string]interface{}{"phone": "13800138000", "password": "password123"}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, float64(0), resp["code"])

		mockAuthSvc.AssertExpectations(t)
	})
}

func TestAuthHandler_Refresh(t *testing.T) {
	mockAuthSvc := new(MockAuthService)
	mockUserSvc := new(MockUserService)
	h := NewAuthHandler(mockAuthSvc, mockUserSvc)
	router := setupAuthRouter(h)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockAuthSvc.On("RefreshToken", ctx, "old-refresh-token").Return("new-access-token", "new-refresh-token", nil).Once()

		body := map[string]interface{}{"refresh_token": "old-refresh-token"}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/auth/refresh", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		mockAuthSvc.AssertExpectations(t)
	})
}

func TestAuthHandler_Register(t *testing.T) {
	mockAuthSvc := new(MockAuthService)
	mockUserSvc := new(MockUserService)
	h := NewAuthHandler(mockAuthSvc, mockUserSvc)
	router := setupAuthRouter(h)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		newUser := &model.User{ID: 1, Phone: "13800138000", Name: "New User", Role: model.RoleOperator}
		mockUserSvc.On("Create", ctx, mock.AnythingOfType("*model.CreateUserRequest")).Return(newUser, nil).Once()

		body := map[string]interface{}{
			"phone":    "13800138000",
			"password": "password123",
			"name":     "New User",
			"role":     "operator",
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		mockUserSvc.AssertExpectations(t)
	})
}
