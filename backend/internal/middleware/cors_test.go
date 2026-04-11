package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCORS_OptionsPreflight(t *testing.T) {
	gin.SetMode(gin.TestMode)
	defer gin.SetMode(gin.ReleaseMode)

	called := false

	r := gin.New()
	r.Use(CORS())
	r.OPTIONS("/test", func(c *gin.Context) {
		called = true
		c.JSON(http.StatusOK, nil)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodOptions, "/test", nil)
	r.ServeHTTP(w, req)

	// OPTIONS should be intercepted by CORS middleware, not reach handler
	assert.False(t, called, "handler should not be called for OPTIONS preflight")
	assert.Equal(t, http.StatusNoContent, w.Code)

	// Verify CORS headers
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "GET, POST, PUT, PATCH, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "Origin, Content-Type, Accept, Authorization, X-Requested-With", w.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(t, "Content-Length, Content-Type", w.Header().Get("Access-Control-Expose-Headers"))
	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
	assert.Equal(t, "86400", w.Header().Get("Access-Control-Max-Age"))
}

func TestCORS_NormalGetRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	defer gin.SetMode(gin.ReleaseMode)

	called := false

	r := gin.New()
	r.Use(CORS())
	r.GET("/test", func(c *gin.Context) {
		called = true
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.True(t, called, "handler should have been called for GET request")
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify CORS headers are still present on normal requests
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "GET, POST, PUT, PATCH, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "Origin, Content-Type, Accept, Authorization, X-Requested-With", w.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(t, "Content-Length, Content-Type", w.Header().Get("Access-Control-Expose-Headers"))
	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
	assert.Equal(t, "86400", w.Header().Get("Access-Control-Max-Age"))
}

func TestCORS_PostRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	defer gin.SetMode(gin.ReleaseMode)

	called := false

	r := gin.New()
	r.Use(CORS())
	r.POST("/test", func(c *gin.Context) {
		called = true
		c.JSON(http.StatusCreated, gin.H{"status": "created"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	r.ServeHTTP(w, req)

	assert.True(t, called, "handler should have been called for POST request")
	assert.Equal(t, http.StatusCreated, w.Code)

	// Verify CORS headers
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
}
