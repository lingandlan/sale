package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"marketplace/backend/internal/service"
	"marketplace/backend/pkg/response"
)

const testJWTSecret = "test-secret-key-for-testing"

// generateTestToken creates a valid JWT token for testing.
func generateTestToken(t *testing.T, secret string, userID int64, phone, role string, expiresAt time.Time) string {
	t.Helper()
	claims := service.Claims{
		UserID: userID,
		Phone:  phone,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	require.NoError(t, err)
	return tokenString
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	defer gin.SetMode(gin.ReleaseMode)

	middleware := NewAuthMiddleware(testJWTSecret)
	token := generateTestToken(t, testJWTSecret, 42, "13800138000", "admin", time.Now().Add(time.Hour))

	called := false
	var capturedUserID int64
	var capturedPhone, capturedRole string

	r := gin.New()
	r.Use(middleware.Auth())
	r.GET("/test", func(c *gin.Context) {
		called = true
		val, _ := c.Get("user_id")
		capturedUserID = val.(int64)
		val, _ = c.Get("phone")
		capturedPhone = val.(string)
		val, _ = c.Get("role")
		capturedRole = val.(string)
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.True(t, called, "handler should have been called")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, int64(42), capturedUserID)
	assert.Equal(t, "13800138000", capturedPhone)
	assert.Equal(t, "admin", capturedRole)
}

func TestAuthMiddleware_MissingAuthorizationHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	defer gin.SetMode(gin.ReleaseMode)

	middleware := NewAuthMiddleware(testJWTSecret)

	r := gin.New()
	r.Use(middleware.Auth())
	r.GET("/test", func(c *gin.Context) {
		t.Fatal("handler should not be called")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var body response.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, 401, body.Code)
	assert.Equal(t, "missing token", body.Message)
}

func TestAuthMiddleware_InvalidFormat_NoBearerPrefix(t *testing.T) {
	gin.SetMode(gin.TestMode)
	defer gin.SetMode(gin.ReleaseMode)

	middleware := NewAuthMiddleware(testJWTSecret)

	r := gin.New()
	r.Use(middleware.Auth())
	r.GET("/test", func(c *gin.Context) {
		t.Fatal("handler should not be called")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "SomeTokenWithoutBearer")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var body response.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, 401, body.Code)
	assert.Equal(t, "invalid token format", body.Message)
}

func TestAuthMiddleware_InvalidFormat_WrongScheme(t *testing.T) {
	gin.SetMode(gin.TestMode)
	defer gin.SetMode(gin.ReleaseMode)

	middleware := NewAuthMiddleware(testJWTSecret)

	r := gin.New()
	r.Use(middleware.Auth())
	r.GET("/test", func(c *gin.Context) {
		t.Fatal("handler should not be called")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Basic abcdef123456")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var body response.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, 401, body.Code)
	assert.Equal(t, "invalid token format", body.Message)
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	defer gin.SetMode(gin.ReleaseMode)

	middleware := NewAuthMiddleware(testJWTSecret)

	r := gin.New()
	r.Use(middleware.Auth())
	r.GET("/test", func(c *gin.Context) {
		t.Fatal("handler should not be called")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer this.is.not.a.valid.token")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var body response.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, 401, body.Code)
	assert.Equal(t, "invalid token", body.Message)
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	defer gin.SetMode(gin.ReleaseMode)

	middleware := NewAuthMiddleware(testJWTSecret)
	token := generateTestToken(t, testJWTSecret, 1, "13800138000", "user", time.Now().Add(-time.Hour))

	r := gin.New()
	r.Use(middleware.Auth())
	r.GET("/test", func(c *gin.Context) {
		t.Fatal("handler should not be called")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var body response.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, 401, body.Code)
	assert.Equal(t, "invalid token", body.Message)
}

func TestAuthMiddleware_WrongSecret(t *testing.T) {
	gin.SetMode(gin.TestMode)
	defer gin.SetMode(gin.ReleaseMode)

	middleware := NewAuthMiddleware(testJWTSecret)
	token := generateTestToken(t, "wrong-secret", 1, "13800138000", "user", time.Now().Add(time.Hour))

	r := gin.New()
	r.Use(middleware.Auth())
	r.GET("/test", func(c *gin.Context) {
		t.Fatal("handler should not be called")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var body response.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, 401, body.Code)
	assert.Equal(t, "invalid token", body.Message)
}

func TestRequireRole_AllowedRole(t *testing.T) {
	gin.SetMode(gin.TestMode)
	defer gin.SetMode(gin.ReleaseMode)

	called := false

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("role", "admin")
		c.Next()
	})
	r.Use(RequireRole("admin", "superadmin"))
	r.GET("/test", func(c *gin.Context) {
		called = true
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.True(t, called, "handler should have been called")
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRequireRole_DeniedRole(t *testing.T) {
	gin.SetMode(gin.TestMode)
	defer gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("role", "user")
		c.Next()
	})
	r.Use(RequireRole("admin", "superadmin"))
	r.GET("/test", func(c *gin.Context) {
		t.Fatal("handler should not be called")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var body response.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, 403, body.Code)
	assert.Equal(t, "insufficient permission", body.Message)
}

func TestRequireRole_NoRoleInContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	defer gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(RequireRole("admin"))
	r.GET("/test", func(c *gin.Context) {
		t.Fatal("handler should not be called")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var body response.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, 403, body.Code)
	assert.Equal(t, "forbidden", body.Message)
}

func TestGetUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	defer gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.GET("/test", func(c *gin.Context) {
		c.Set("user_id", int64(99))
		id := GetUserID(c)
		c.JSON(http.StatusOK, gin.H{"user_id": id})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "99")
}

func TestGetPhone(t *testing.T) {
	gin.SetMode(gin.TestMode)
	defer gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.GET("/test", func(c *gin.Context) {
		c.Set("phone", "13900139000")
		phone := GetPhone(c)
		c.JSON(http.StatusOK, gin.H{"phone": phone})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "13900139000")
}

func TestGetRole(t *testing.T) {
	gin.SetMode(gin.TestMode)
	defer gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.GET("/test", func(c *gin.Context) {
		c.Set("role", "admin")
		role := GetRole(c)
		c.JSON(http.StatusOK, gin.H{"role": role})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "admin")
}
