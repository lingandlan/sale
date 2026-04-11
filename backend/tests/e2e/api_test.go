//go:build e2e
// +build e2e

package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"marketplace/backend/internal/config"
	"marketplace/backend/internal/handler"
	"marketplace/backend/internal/middleware"
	"marketplace/backend/internal/model"
	"marketplace/backend/internal/service"
)

// ========== Mock implementations ==========

// MockE2EUserRepo implements repository.UserRepoInterface
type MockE2EUserRepo struct {
	mock.Mock
}

func (m *MockE2EUserRepo) GetByID(ctx context.Context, id int64) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockE2EUserRepo) GetByPhone(ctx context.Context, phone string) (*model.User, error) {
	args := m.Called(ctx, phone)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockE2EUserRepo) Create(ctx context.Context, user *model.User) (int64, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockE2EUserRepo) Update(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockE2EUserRepo) UpdatePassword(ctx context.Context, id int64, password string) error {
	args := m.Called(ctx, id, password)
	return args.Error(0)
}

func (m *MockE2EUserRepo) List(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]*model.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockE2EUserRepo) ListWithFilters(ctx context.Context, page, pageSize int, keyword, role string, status *int8) ([]*model.User, int64, error) {
	args := m.Called(ctx, page, pageSize, keyword, role, status)
	return args.Get(0).([]*model.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockE2EUserRepo) UpdateStatus(ctx context.Context, id int64, status int8) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockE2EUserRepo) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockE2ERechargeRepo implements repository.RechargeRepoInterface
type MockE2ERechargeRepo struct {
	mock.Mock
}

func (m *MockE2ERechargeRepo) CreateRechargeApplication(app *model.RechargeApplication) error {
	args := m.Called(app)
	return args.Error(0)
}
func (m *MockE2ERechargeRepo) GetRechargeApplications(status string, page, pageSize int) ([]model.RechargeApplication, int64, error) {
	args := m.Called(status, page, pageSize)
	return args.Get(0).([]model.RechargeApplication), args.Get(1).(int64), args.Error(2)
}
func (m *MockE2ERechargeRepo) GetRechargeApplicationByID(id string) (*model.RechargeApplication, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.RechargeApplication), args.Error(1)
}
func (m *MockE2ERechargeRepo) UpdateRechargeApplicationStatus(id, status, approvedBy, remark string) error {
	args := m.Called(id, status, approvedBy, remark)
	return args.Error(0)
}
func (m *MockE2ERechargeRepo) CreateCRecharge(recharge *model.CRecharge) error {
	args := m.Called(recharge)
	return args.Error(0)
}
func (m *MockE2ERechargeRepo) GetCRechargeList(memberPhone, centerID string, page, pageSize int) ([]model.CRecharge, int64, error) {
	args := m.Called(memberPhone, centerID, page, pageSize)
	return args.Get(0).([]model.CRecharge), args.Get(1).(int64), args.Error(2)
}
func (m *MockE2ERechargeRepo) GetCRechargeByID(id string) (*model.CRecharge, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CRecharge), args.Error(1)
}
func (m *MockE2ERechargeRepo) CreateCard(card *model.StoreCard) error {
	args := m.Called(card)
	return args.Error(0)
}
func (m *MockE2ERechargeRepo) GetCardList(status, cardNo, holderPhone string, page, pageSize int) ([]model.StoreCard, int64, error) {
	args := m.Called(status, cardNo, holderPhone, page, pageSize)
	return args.Get(0).([]model.StoreCard), args.Get(1).(int64), args.Error(2)
}
func (m *MockE2ERechargeRepo) GetCardByCardNo(cardNo string) (*model.StoreCard, error) {
	args := m.Called(cardNo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.StoreCard), args.Error(1)
}
func (m *MockE2ERechargeRepo) UpdateCardBalance(cardNo string, balance float64) error {
	args := m.Called(cardNo, balance)
	return args.Error(0)
}
func (m *MockE2ERechargeRepo) UpdateCardStatus(cardNo, status string) error {
	args := m.Called(cardNo, status)
	return args.Error(0)
}
func (m *MockE2ERechargeRepo) CreateCardTransaction(transaction *model.CardTransaction) error {
	args := m.Called(transaction)
	return args.Error(0)
}
func (m *MockE2ERechargeRepo) GetCardTransactions(cardNo string) ([]model.CardTransaction, error) {
	args := m.Called(cardNo)
	return args.Get(0).([]model.CardTransaction), args.Error(1)
}
func (m *MockE2ERechargeRepo) GetCardStats() (map[string]int64, error) {
	args := m.Called()
	return args.Get(0).(map[string]int64), args.Error(1)
}
func (m *MockE2ERechargeRepo) GetCenters() ([]model.RechargeCenter, error) {
	args := m.Called()
	return args.Get(0).([]model.RechargeCenter), args.Error(1)
}
func (m *MockE2ERechargeRepo) CreateCenter(center *model.RechargeCenter) error {
	args := m.Called(center)
	return args.Error(0)
}
func (m *MockE2ERechargeRepo) UpdateCenter(center *model.RechargeCenter) error {
	args := m.Called(center)
	return args.Error(0)
}
func (m *MockE2ERechargeRepo) DeleteCenter(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *MockE2ERechargeRepo) GetOperators() ([]model.RechargeOperator, error) {
	args := m.Called()
	return args.Get(0).([]model.RechargeOperator), args.Error(1)
}
func (m *MockE2ERechargeRepo) CreateOperator(operator *model.RechargeOperator) error {
	args := m.Called(operator)
	return args.Error(0)
}
func (m *MockE2ERechargeRepo) UpdateOperator(operator *model.RechargeOperator) error {
	args := m.Called(operator)
	return args.Error(0)
}
func (m *MockE2ERechargeRepo) DeleteOperator(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// ========== Helpers ==========

const e2eJWTSecret = "e2e-test-secret-key"

func generateTestToken(t *testing.T, userID int64, role string) string {
	t.Helper()
	claims := service.Claims{
		UserID: userID,
		Phone:  "13800138000",
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(e2eJWTSecret))
	require.NoError(t, err)
	return signed
}

func makeRequest(t *testing.T, r *gin.Engine, method, path string, body interface{}, token string) *httptest.ResponseRecorder {
	t.Helper()

	var bodyReader *bytes.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		require.NoError(t, err)
		bodyReader = bytes.NewReader(jsonBody)
	} else {
		bodyReader = bytes.NewReader([]byte{})
	}

	req := httptest.NewRequest(method, path, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func parseResponse(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	return resp
}

func setupE2EEngine(t *testing.T) (*gin.Engine, *MockE2EUserRepo, *MockE2ERechargeRepo) {
	gin.SetMode(gin.TestMode)

	userRepo := new(MockE2EUserRepo)
	rechargeRepo := new(MockE2ERechargeRepo)

	userSvc := service.NewUserService(userRepo)
	authSvc := service.NewAuthService(
		&config.JWTConfig{Secret: e2eJWTSecret, ExpireHours: 24, RefreshExpireHours: 168},
		nil,
		userRepo,
	)
	rechargeSvc := service.NewRechargeService(rechargeRepo)

	authMW := middleware.NewAuthMiddleware(e2eJWTSecret)

	authHandler := handler.NewAuthHandler(authSvc, userSvc)
	adminHandler := handler.NewAdminHandler(userSvc)
	rechargeHandler := handler.NewRechargeHandler(rechargeSvc)

	r := gin.New()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	v1 := r.Group("/api/v1")

	// Public auth routes
	auth := v1.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
	}

	// Auth-protected routes (no RBAC for e2e simplicity)
	protected := v1.Group("")
	protected.Use(authMW.Auth())
	{
		// Dashboard
		dashboard := protected.Group("/dashboard")
		{
			dashboard.GET("/statistics", rechargeHandler.GetDashboardStatistics)
		}

		// B-recharge
		bApply := protected.Group("/recharge/b-apply")
		{
			bApply.POST("", rechargeHandler.CreateBRechargeApplication)
		}
		bApproval := protected.Group("/recharge/b-approval")
		{
			bApproval.GET("", rechargeHandler.GetRechargeApplicationList)
			bApproval.GET("/:id", rechargeHandler.GetRechargeApplicationDetail)
			bApproval.POST("/action", rechargeHandler.ApprovalRechargeApplication)
		}

		// C-recharge
		cEntry := protected.Group("/recharge/c-entry")
		{
			cEntry.POST("", rechargeHandler.CreateCRecharge)
			cEntry.GET("", rechargeHandler.GetCRechargeList)
		}

		// Records
		records := protected.Group("/recharge/records")
		{
			records.GET("", rechargeHandler.GetCRechargeList)
			records.GET("/:id", rechargeHandler.GetCRechargeDetail)
		}

		// Card
		card := protected.Group("/card")
		{
			card.GET("/verify/:cardNo", rechargeHandler.VerifyCard)
			card.POST("/consume", rechargeHandler.ConsumeCard)
			card.GET("/list", rechargeHandler.GetCardList)
			card.GET("/stats", rechargeHandler.GetCardStats)
			card.POST("/issue", rechargeHandler.IssueCard)
			card.POST("/:cardNo/status", rechargeHandler.UpdateCardStatus)
		}

		// Center
		center := protected.Group("/center")
		{
			center.GET("", rechargeHandler.GetCenters)
			center.POST("", rechargeHandler.CreateCenter)
			center.PUT("/:id", rechargeHandler.UpdateCenter)
			center.DELETE("/:id", rechargeHandler.DeleteCenter)
		}

		// Operator
		operator := protected.Group("/operator")
		{
			operator.GET("", rechargeHandler.GetOperators)
			operator.POST("", rechargeHandler.CreateOperator)
		}

		// Admin (skip RBAC in e2e)
		admin := protected.Group("/admin")
		{
			admin.GET("/users", adminHandler.ListUsers)
			admin.POST("/users", adminHandler.CreateUser)
			admin.PUT("/users/:id", adminHandler.UpdateUser)
			admin.DELETE("/users/:id", adminHandler.DeleteUser)
		}
	}

	return r, userRepo, rechargeRepo
}

// ========== E2E Tests ==========

func TestE2E_HealthCheck(t *testing.T) {
	r, _, _ := setupE2EEngine(t)

	w := makeRequest(t, r, "GET", "/health", nil, "")
	assert.Equal(t, http.StatusOK, w.Code)

	resp := parseResponse(t, w)
	assert.Equal(t, "ok", resp["status"])
}

func TestE2E_Login_Flow(t *testing.T) {
	r, userRepo, _ := setupE2EEngine(t)

	t.Run("Login with non-existent user returns 401", func(t *testing.T) {
		userRepo.On("GetByPhone", mock.Anything, "19900000000").Return(nil, fmt.Errorf("not found")).Once()

		w := makeRequest(t, r, "POST", "/api/v1/auth/login", map[string]string{
			"phone":    "19900000000",
			"password": "wrongpass",
		}, "")
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Login with missing fields returns error", func(t *testing.T) {
		w := makeRequest(t, r, "POST", "/api/v1/auth/login", map[string]string{
			"phone": "13800138000",
		}, "")
		assert.Equal(t, http.StatusOK, w.Code)

		resp := parseResponse(t, w)
		assert.NotEqual(t, float64(0), resp["code"])
	})

	// Note: Login success requires Redis for refresh token storage.
	// Full login flow is tested in handler/auth_test.go with mocked service.
}

func TestE2E_AdminUsers_CRUD(t *testing.T) {
	r, userRepo, _ := setupE2EEngine(t)
	token := generateTestToken(t, 1, "hq_admin")

	t.Run("List users with pagination", func(t *testing.T) {
		userRepo.On("ListWithFilters", mock.Anything, 1, 10, "", "", (*int8)(nil)).
			Return([]*model.User{
				{ID: 1, Phone: "13800138000", Name: "Admin", Role: "hq_admin", Status: 1},
			}, int64(1), nil).Once()

		w := makeRequest(t, r, "GET", "/api/v1/admin/users?page=1&page_size=10", nil, token)
		assert.Equal(t, http.StatusOK, w.Code)

		resp := parseResponse(t, w)
		assert.Equal(t, float64(0), resp["code"])
	})

	t.Run("Create user returns success", func(t *testing.T) {
		userRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.User")).Return(int64(2), nil).Once()

		w := makeRequest(t, r, "POST", "/api/v1/admin/users", map[string]interface{}{
			"phone":    "13800138001",
			"password": "Test123456",
			"name":     "新用户",
			"role":     "operator",
		}, token)
		assert.Equal(t, http.StatusOK, w.Code)

		resp := parseResponse(t, w)
		assert.Equal(t, float64(0), resp["code"])
	})

	t.Run("Create user with missing fields returns error", func(t *testing.T) {
		w := makeRequest(t, r, "POST", "/api/v1/admin/users", map[string]interface{}{
			"phone": "13800138002",
		}, token)
		assert.Equal(t, http.StatusOK, w.Code)

		resp := parseResponse(t, w)
		assert.NotEqual(t, float64(0), resp["code"])
	})
}

func TestE2E_BRecharge_Flow(t *testing.T) {
	r, _, rechargeRepo := setupE2EEngine(t)
	token := generateTestToken(t, 1, "hq_admin")

	t.Run("Create B-recharge application", func(t *testing.T) {
		rechargeRepo.On("CreateRechargeApplication", mock.AnythingOfType("*model.RechargeApplication")).Return(nil).Once()

		w := makeRequest(t, r, "POST", "/api/v1/recharge/b-apply", map[string]interface{}{
			"centerId":            "center-1",
			"centerName":          "北京中心",
			"amount":              10000,
			"lastMonthConsumption": 120000,
			"applicantId":         "user-1",
			"applicantName":       "张三",
			"transactionNo":       "TX-E2E-001",
			"screenshot":          "",
			"remark":              "E2E测试",
		}, token)
		assert.Equal(t, http.StatusOK, w.Code)

		resp := parseResponse(t, w)
		assert.Equal(t, float64(0), resp["code"])
	})

	t.Run("List B-recharge applications", func(t *testing.T) {
		rechargeRepo.On("GetRechargeApplications", "pending", 1, 10).
			Return([]model.RechargeApplication{
				{ID: "app-1", Status: "pending", Amount: 10000, BasePoints: 10000, RebatePoints: 200},
			}, int64(1), nil).Once()

		w := makeRequest(t, r, "GET", "/api/v1/recharge/b-approval?status=pending&page=1&page_size=10", nil, token)
		assert.Equal(t, http.StatusOK, w.Code)

		resp := parseResponse(t, w)
		assert.Equal(t, float64(0), resp["code"])
	})
}

func TestE2E_Card_Operations(t *testing.T) {
	r, _, rechargeRepo := setupE2EEngine(t)
	token := generateTestToken(t, 1, "operator")

	t.Run("Get card stats", func(t *testing.T) {
		rechargeRepo.On("GetCardStats").Return(map[string]int64{
			"totalCards":    100,
			"activeCards":   60,
			"totalBalance":  50000,
			"todayConsume":  0,
			"todayIssue":    0,
			"expireIn7Days": 0,
		}, nil).Once()

		w := makeRequest(t, r, "GET", "/api/v1/card/stats", nil, token)
		assert.Equal(t, http.StatusOK, w.Code)

		resp := parseResponse(t, w)
		assert.Equal(t, float64(0), resp["code"])
	})

	t.Run("Verify card", func(t *testing.T) {
		rechargeRepo.On("GetCardByCardNo", "TJ20260411001").Return(&model.StoreCard{
			ID:         "card-1",
			CardNo:     "TJ20260411001",
			Status:     "active",
			Balance:    5000.0,
			ExpiryDate: time.Now().AddDate(1, 0, 0),
		}, nil).Once()

		w := makeRequest(t, r, "GET", "/api/v1/card/verify/TJ20260411001", nil, token)
		assert.Equal(t, http.StatusOK, w.Code)

		resp := parseResponse(t, w)
		assert.Equal(t, float64(0), resp["code"])
	})

	t.Run("Issue card", func(t *testing.T) {
		rechargeRepo.On("CreateCard", mock.AnythingOfType("*model.StoreCard")).Return(nil).Once()
		rechargeRepo.On("CreateCardTransaction", mock.AnythingOfType("*model.CardTransaction")).Return(nil).Once()

		w := makeRequest(t, r, "POST", "/api/v1/card/issue", map[string]interface{}{
			"holderId":    "member-1",
			"holderName":  "李四",
			"holderPhone": "13800001111",
			"amount":      5000,
			"centerId":    "center-1",
			"centerName":  "北京中心",
			"operatorId":  "op-1",
		}, token)
		assert.Equal(t, http.StatusOK, w.Code)

		resp := parseResponse(t, w)
		assert.Equal(t, float64(0), resp["code"])
	})
}

func TestE2E_Center_CRUD(t *testing.T) {
	r, _, rechargeRepo := setupE2EEngine(t)
	token := generateTestToken(t, 1, "hq_admin")

	t.Run("Get centers list", func(t *testing.T) {
		rechargeRepo.On("GetCenters").Return([]model.RechargeCenter{
			{ID: "c1", Name: "北京中心", Code: "BJ001", Status: "active"},
		}, nil).Once()

		w := makeRequest(t, r, "GET", "/api/v1/center", nil, token)
		assert.Equal(t, http.StatusOK, w.Code)

		resp := parseResponse(t, w)
		assert.Equal(t, float64(0), resp["code"])
	})

	t.Run("Create center", func(t *testing.T) {
		rechargeRepo.On("CreateCenter", mock.AnythingOfType("*model.RechargeCenter")).Return(nil).Once()

		w := makeRequest(t, r, "POST", "/api/v1/center", map[string]interface{}{
			"name":    "上海中心",
			"code":    "SH001",
			"address": "上海市浦东新区",
			"phone":   "021-12345678",
		}, token)
		assert.Equal(t, http.StatusOK, w.Code)

		resp := parseResponse(t, w)
		assert.Equal(t, float64(0), resp["code"])
	})
}

func TestE2E_Unauthorized_Access(t *testing.T) {
	r, _, _ := setupE2EEngine(t)

	t.Run("No token returns 401", func(t *testing.T) {
		w := makeRequest(t, r, "GET", "/api/v1/admin/users", nil, "")
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Invalid token returns 401", func(t *testing.T) {
		w := makeRequest(t, r, "GET", "/api/v1/admin/users", nil, "invalid-token")
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Dashboard without token returns 401", func(t *testing.T) {
		w := makeRequest(t, r, "GET", "/api/v1/dashboard/statistics", nil, "")
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Card operations without token returns 401", func(t *testing.T) {
		w := makeRequest(t, r, "GET", "/api/v1/card/stats", nil, "")
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
