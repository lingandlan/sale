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
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetByID(ctx context.Context, id int64) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) GetByPhone(ctx context.Context, phone string) (*model.User, error) {
	args := m.Called(ctx, phone)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) Create(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) Update(ctx context.Context, id int64, req *model.UpdateUserRequest) (*model.User, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) UpdateStatus(ctx context.Context, id int64, req *model.UpdateUserStatusRequest) error {
	args := m.Called(ctx, id, req)
	return args.Error(0)
}

func (m *MockUserService) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserService) ResetPassword(ctx context.Context, id int64, req *model.ResetPasswordRequest) error {
	args := m.Called(ctx, id, req)
	return args.Error(0)
}

func (m *MockUserService) UpdatePassword(ctx context.Context, id int64, hashedPassword string) error {
	args := m.Called(ctx, id, hashedPassword)
	return args.Error(0)
}

func (m *MockUserService) List(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]*model.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserService) ListWithFilters(ctx context.Context, req *model.ListUsersRequest) (*model.ListUsersResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ListUsersResponse), args.Error(1)
}

func setupUserRouter(h *UserHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/user/info", func(c *gin.Context) {
		c.Set("user_id", int64(1))
		h.GetUserInfo(c)
	})
	r.PUT("/user/info", func(c *gin.Context) {
		c.Set("user_id", int64(1))
		h.UpdateUserInfo(c)
	})
	r.GET("/user/:id", h.GetUserByID)
	r.GET("/users", h.ListUsers)
	return r
}

func TestUserHandler_GetUserByID(t *testing.T) {
	mockSvc := new(MockUserService)
	h := NewUserHandler(mockSvc, nil)
	router := setupUserRouter(h)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		expectedUser := &model.User{ID: 1, Phone: "13800138000", Name: "Test User", Role: model.RoleOperator}
		mockSvc.On("GetByID", ctx, int64(1)).Return(expectedUser, nil).Once()

		req, _ := http.NewRequest("GET", "/user/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, float64(0), resp["code"])

		mockSvc.AssertExpectations(t)
	})
}

func TestUserHandler_UpdateUserInfo(t *testing.T) {
	mockSvc := new(MockUserService)
	h := NewUserHandler(mockSvc, nil)
	router := setupUserRouter(h)
	ctx := context.Background()
	newName := "New Name"

	t.Run("success", func(t *testing.T) {
		updatedUser := &model.User{ID: 1, Phone: "13800138000", Name: newName, Role: model.RoleOperator}
		mockSvc.On("Update", ctx, int64(1), mock.AnythingOfType("*model.UpdateUserRequest")).Return(updatedUser, nil).Once()

		body := map[string]interface{}{"name": newName}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("PUT", "/user/info", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		mockSvc.AssertExpectations(t)
	})
}

func TestUserHandler_ListUsers(t *testing.T) {
	mockSvc := new(MockUserService)
	h := NewUserHandler(mockSvc, nil)
	router := setupUserRouter(h)
	ctx := context.Background()

	t.Run("success with pagination", func(t *testing.T) {
		users := []*model.User{
			{ID: 1, Phone: "13800138000", Name: "User 1", Role: model.RoleOperator},
			{ID: 2, Phone: "13800138001", Name: "User 2", Role: model.RoleHQAdmin},
		}
		mockSvc.On("List", ctx, 1, 20).Return(users, int64(2), nil).Once()

		req, _ := http.NewRequest("GET", "/users?page=1&page_size=20", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, float64(0), resp["code"])

		mockSvc.AssertExpectations(t)
	})

	t.Run("default pagination", func(t *testing.T) {
		users := []*model.User{}
		mockSvc.On("List", ctx, 1, 20).Return(users, int64(0), nil).Once()

		req, _ := http.NewRequest("GET", "/users", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})
}
