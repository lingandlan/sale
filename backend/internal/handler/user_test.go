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

func (m *MockUserService) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	args := m.Called(ctx, username)
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

func (m *MockUserService) List(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]*model.User), args.Get(1).(int64), args.Error(2)
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
	h := NewUserHandler(mockSvc)
	router := setupUserRouter(h)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		expectedUser := &model.User{ID: 1, Username: "testuser"}
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
	h := NewUserHandler(mockSvc)
	router := setupUserRouter(h)
	ctx := context.Background()
	newNickname := "New Nickname"

	t.Run("success", func(t *testing.T) {
		updatedUser := &model.User{ID: 1, Username: "testuser", Nickname: newNickname}
		mockSvc.On("Update", ctx, int64(1), mock.AnythingOfType("*model.UpdateUserRequest")).Return(updatedUser, nil).Once()

		body := map[string]interface{}{"nickname": newNickname}
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
	h := NewUserHandler(mockSvc)
	router := setupUserRouter(h)
	ctx := context.Background()

	t.Run("success with pagination", func(t *testing.T) {
		users := []*model.User{
			{ID: 1, Username: "user1"},
			{ID: 2, Username: "user2"},
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
