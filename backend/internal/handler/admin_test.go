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
	apperrors "marketplace/backend/pkg/errors"
)

func setupAdminRouter(h *AdminHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	admin := r.Group("/admin")
	{
		admin.GET("/users", h.ListUsers)
		admin.POST("/users", h.CreateUser)
		admin.PUT("/users/:id", h.UpdateUser)
		admin.POST("/users/:id/reset-password", h.ResetPassword)
		admin.PUT("/users/:id/status", h.UpdateUserStatus)
		admin.DELETE("/users/:id", h.DeleteUser)
	}
	return r
}

func assertResponseCode(t *testing.T, body []byte, expectedCode float64) {
	t.Helper()
	var resp map[string]interface{}
	json.Unmarshal(body, &resp)
	assert.Equal(t, expectedCode, resp["code"])
}

func TestAdminHandler_ListUsers(t *testing.T) {
	t.Run("success with filters", func(t *testing.T) {
		mockSvc := new(MockUserService)
		h := NewAdminHandler(mockSvc, nil)
		router := setupAdminRouter(h)
		ctx := context.Background()

		resp := &model.ListUsersResponse{
			Items:    []*model.User{{ID: 1, Phone: "13800138000", Name: "Admin", Role: model.RoleSuperAdmin}},
			Total:    1,
			Page:     1,
			PageSize: 20,
		}
		mockSvc.On("ListWithFilters", ctx, mock.AnythingOfType("*model.ListUsersRequest")).Return(resp, nil).Once()

		req, _ := http.NewRequest("GET", "/admin/users?keyword=test&role=super_admin&page=1&page_size=20", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 0)
		mockSvc.AssertExpectations(t)
	})

	t.Run("default pagination", func(t *testing.T) {
		mockSvc := new(MockUserService)
		h := NewAdminHandler(mockSvc, nil)
		router := setupAdminRouter(h)
		ctx := context.Background()

		resp := &model.ListUsersResponse{Items: []*model.User{}, Total: 0, Page: 1, PageSize: 20}
		mockSvc.On("ListWithFilters", ctx, mock.AnythingOfType("*model.ListUsersRequest")).Return(resp, nil).Once()

		req, _ := http.NewRequest("GET", "/admin/users", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestAdminHandler_CreateUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockUserService)
		h := NewAdminHandler(mockSvc, nil)
		router := setupAdminRouter(h)
		ctx := context.Background()

		createdUser := &model.User{ID: 1, Phone: "13800138000", Name: "New User", Role: model.RoleOperator}
		mockSvc.On("Create", ctx, mock.AnythingOfType("*model.CreateUserRequest")).Return(createdUser, nil).Once()

		body := map[string]interface{}{
			"phone":    "13800138000",
				"username": "testuser",
			"password": "123456",
			"name":     "New User",
			"role":     "operator",
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/admin/users", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 0)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid json returns error code", func(t *testing.T) {
		mockSvc := new(MockUserService)
		h := NewAdminHandler(mockSvc, nil)
		router := setupAdminRouter(h)

		req, _ := http.NewRequest("POST", "/admin/users", bytes.NewBufferString("{invalid}"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 400)
	})
}

func TestAdminHandler_UpdateUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockUserService)
		h := NewAdminHandler(mockSvc, nil)
		router := setupAdminRouter(h)
		ctx := context.Background()

		name := "Updated Name"
		updatedUser := &model.User{ID: 1, Phone: "13800138000", Name: name, Role: model.RoleOperator}
		mockSvc.On("Update", ctx, int64(1), mock.AnythingOfType("*model.UpdateUserRequest")).Return(updatedUser, nil).Once()

		body := map[string]interface{}{"name": name}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("PUT", "/admin/users/1", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 0)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid user id returns error code", func(t *testing.T) {
		mockSvc := new(MockUserService)
		h := NewAdminHandler(mockSvc, nil)
		router := setupAdminRouter(h)

		req, _ := http.NewRequest("PUT", "/admin/users/abc", bytes.NewBufferString(`{"name":"test"}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 400)
	})

	t.Run("user not found", func(t *testing.T) {
		mockSvc := new(MockUserService)
		h := NewAdminHandler(mockSvc, nil)
		router := setupAdminRouter(h)
		ctx := context.Background()

		mockSvc.On("Update", ctx, int64(999), mock.AnythingOfType("*model.UpdateUserRequest")).Return(nil, apperrors.ErrNotFound).Once()

		req, _ := http.NewRequest("PUT", "/admin/users/999", bytes.NewBufferString(`{"name":"test"}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 404)
		mockSvc.AssertExpectations(t)
	})
}

func TestAdminHandler_ResetPassword(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockUserService)
		h := NewAdminHandler(mockSvc, nil)
		router := setupAdminRouter(h)
		ctx := context.Background()

		mockSvc.On("ResetPassword", ctx, int64(1), mock.AnythingOfType("*model.ResetPasswordRequest")).Return(nil).Once()

		body := map[string]interface{}{"new_password": "newpass123"}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/admin/users/1/reset-password", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 0)
		mockSvc.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockSvc := new(MockUserService)
		h := NewAdminHandler(mockSvc, nil)
		router := setupAdminRouter(h)
		ctx := context.Background()

		mockSvc.On("ResetPassword", ctx, int64(999), mock.AnythingOfType("*model.ResetPasswordRequest")).Return(apperrors.ErrNotFound).Once()

		req, _ := http.NewRequest("POST", "/admin/users/999/reset-password", bytes.NewBufferString(`{"new_password":"newpass123"}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 404)
		mockSvc.AssertExpectations(t)
	})
}

func TestAdminHandler_UpdateUserStatus(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockUserService)
		h := NewAdminHandler(mockSvc, nil)
		router := setupAdminRouter(h)
		ctx := context.Background()

		mockSvc.On("UpdateStatus", ctx, int64(1), mock.AnythingOfType("*model.UpdateUserStatusRequest")).Return(nil).Once()
		disabledUser := &model.User{ID: 1, Phone: "13800138000", Name: "User", Role: model.RoleOperator, Status: 1}
		mockSvc.On("GetByID", ctx, int64(1)).Return(disabledUser, nil).Once()

		body := map[string]interface{}{"status": 1}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("PUT", "/admin/users/1/status", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 0)
		mockSvc.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockSvc := new(MockUserService)
		h := NewAdminHandler(mockSvc, nil)
		router := setupAdminRouter(h)
		ctx := context.Background()

		mockSvc.On("UpdateStatus", ctx, int64(999), mock.AnythingOfType("*model.UpdateUserStatusRequest")).Return(apperrors.ErrNotFound).Once()

		req, _ := http.NewRequest("PUT", "/admin/users/999/status", bytes.NewBufferString(`{"status": 1}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 404)
		mockSvc.AssertExpectations(t)
	})
}

func TestAdminHandler_DeleteUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockUserService)
		h := NewAdminHandler(mockSvc, nil)
		router := setupAdminRouter(h)
		ctx := context.Background()

		mockSvc.On("Delete", ctx, int64(1)).Return(nil).Once()

		req, _ := http.NewRequest("DELETE", "/admin/users/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 0)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid id returns error code", func(t *testing.T) {
		mockSvc := new(MockUserService)
		h := NewAdminHandler(mockSvc, nil)
		router := setupAdminRouter(h)

		req, _ := http.NewRequest("DELETE", "/admin/users/abc", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 400)
	})

	t.Run("user not found", func(t *testing.T) {
		mockSvc := new(MockUserService)
		h := NewAdminHandler(mockSvc, nil)
		router := setupAdminRouter(h)
		ctx := context.Background()

		mockSvc.On("Delete", ctx, int64(999)).Return(apperrors.ErrNotFound).Once()

		req, _ := http.NewRequest("DELETE", "/admin/users/999", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 404)
		mockSvc.AssertExpectations(t)
	})
}
