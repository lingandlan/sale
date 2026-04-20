package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"marketplace/backend/pkg/response"
)

func TestRecovery_NormalRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	defer gin.SetMode(gin.ReleaseMode)

	called := false

	r := gin.New()
	r.Use(Recovery())
	r.GET("/test", func(c *gin.Context) {
		called = true
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.True(t, called, "handler should have been called")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "ok")
}

func TestRecovery_HandlerPanics(t *testing.T) {
	gin.SetMode(gin.TestMode)
	defer gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(Recovery())
	r.GET("/test", func(c *gin.Context) {
		panic("something went wrong")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var body response.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, 500, body.Code)
	assert.Equal(t, "internal server error", body.Message)
	assert.Empty(t, body.Detail)
}

func TestRecovery_HandlerPanicsWithRuntimeError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	defer gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(Recovery())
	r.GET("/test", func(c *gin.Context) {
		var ptr *string
		// This will panic with a nil pointer dereference
		_ = *ptr
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var body response.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, 500, body.Code)
	assert.Equal(t, "internal server error", body.Message)
	assert.Empty(t, body.Detail)
}

func TestRecovery_MultipleRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)
	defer gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(Recovery())
	r.GET("/panic", func(c *gin.Context) {
		panic("panic!")
	})
	r.GET("/ok", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// First request panics
	w1 := httptest.NewRecorder()
	req1 := httptest.NewRequest(http.MethodGet, "/panic", nil)
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusInternalServerError, w1.Code)

	// Second request should still work (server didn't crash)
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/ok", nil)
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}
