package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"marketplace/backend/internal/model"
)

// MockRechargeService implements service.RechargeServiceInterface for testing
type MockRechargeService struct {
	mock.Mock
}

func (m *MockRechargeService) CalculatePoints(amount float64, lastMonthConsumption float64) (int, int, int) {
	args := m.Called(amount, lastMonthConsumption)
	return args.Int(0), args.Int(1), args.Int(2)
}

func (m *MockRechargeService) CreateBRechargeApplication(data map[string]interface{}) (*model.RechargeApplication, error) {
	args := m.Called(data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.RechargeApplication), args.Error(1)
}

func (m *MockRechargeService) GetRechargeApplicationList(status string, page, pageSize int) (map[string]interface{}, error) {
	args := m.Called(status, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockRechargeService) GetRechargeApplicationDetail(id string) (*model.RechargeApplication, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.RechargeApplication), args.Error(1)
}

func (m *MockRechargeService) ApproveRechargeApplication(id, action, approvedBy, remark string) error {
	args := m.Called(id, action, approvedBy, remark)
	return args.Error(0)
}

func (m *MockRechargeService) CreateCRecharge(data map[string]interface{}) (*model.CRecharge, error) {
	args := m.Called(data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CRecharge), args.Error(1)
}

func (m *MockRechargeService) GetCRechargeList(memberPhone, centerID string, page, pageSize int) (map[string]interface{}, error) {
	args := m.Called(memberPhone, centerID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockRechargeService) GetCRechargeDetail(id string) (*model.CRecharge, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CRecharge), args.Error(1)
}

func (m *MockRechargeService) IssueCard(data map[string]interface{}) (*model.StoreCard, error) {
	args := m.Called(data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.StoreCard), args.Error(1)
}

func (m *MockRechargeService) VerifyCard(cardNo string) (*model.StoreCard, error) {
	args := m.Called(cardNo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.StoreCard), args.Error(1)
}

func (m *MockRechargeService) ConsumeCard(cardNo string, amount float64, remark, operatorID string) error {
	args := m.Called(cardNo, amount, remark, operatorID)
	return args.Error(0)
}

func (m *MockRechargeService) UpdateCardStatus(cardNo, status string) error {
	args := m.Called(cardNo, status)
	return args.Error(0)
}

func (m *MockRechargeService) GetCardList(status, cardNo, holderPhone string, page, pageSize int) (map[string]interface{}, error) {
	args := m.Called(status, cardNo, holderPhone, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockRechargeService) GetCardDetail(cardNo string) (map[string]interface{}, error) {
	args := m.Called(cardNo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockRechargeService) GetCardStats() (map[string]interface{}, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockRechargeService) GetCenters() ([]model.RechargeCenter, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.RechargeCenter), args.Error(1)
}

func (m *MockRechargeService) CreateCenter(data map[string]interface{}) (*model.RechargeCenter, error) {
	args := m.Called(data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.RechargeCenter), args.Error(1)
}

func (m *MockRechargeService) UpdateCenter(id string, data map[string]interface{}) (*model.RechargeCenter, error) {
	args := m.Called(id, data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.RechargeCenter), args.Error(1)
}

func (m *MockRechargeService) DeleteCenter(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRechargeService) GetOperators() ([]model.RechargeOperator, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.RechargeOperator), args.Error(1)
}

func (m *MockRechargeService) CreateOperator(data map[string]interface{}) (*model.RechargeOperator, error) {
	args := m.Called(data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.RechargeOperator), args.Error(1)
}

func (m *MockRechargeService) UpdateOperator(id string, data map[string]interface{}) (*model.RechargeOperator, error) {
	args := m.Called(id, data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.RechargeOperator), args.Error(1)
}

func (m *MockRechargeService) DeleteOperator(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// setupRechargeRouter creates a test router for RechargeHandler
func setupRechargeRouter(h *RechargeHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	recharge := r.Group("/recharge")
	{
		recharge.POST("/b-application", h.CreateBRechargeApplication)
		recharge.POST("/b-application/approval", h.ApprovalRechargeApplication)
		recharge.GET("/cards/verify/:cardNo", h.VerifyCard)
		recharge.POST("/cards/consume", h.ConsumeCard)
		recharge.GET("/cards", h.GetCardList)
		recharge.POST("/centers", h.CreateCenter)
		recharge.GET("/dashboard/statistics", h.GetDashboardStatistics)
	}
	return r
}

// ========== Tests ==========

func TestRechargeHandler_CreateBRechargeApplication(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		app := &model.RechargeApplication{
			ID:     "app-001",
			Status: "pending",
		}
		mockSvc.On("CreateBRechargeApplication", mock.Anything).Return(app, nil).Once()

		body := map[string]interface{}{
			"centerId":      "center-001",
			"amount":        5000,
			"transactionNo": "TXN001",
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/recharge/b-application", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 0)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid json returns error code", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		req, _ := http.NewRequest("POST", "/recharge/b-application", bytes.NewBufferString("{invalid}"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 400)
	})
}

func TestRechargeHandler_ApprovalRechargeApplication(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		mockSvc.On("ApproveRechargeApplication", "app-001", "approve", "admin", "ok").Return(nil).Once()

		body := map[string]string{
			"id":     "app-001",
			"action": "approve",
			"remark": "ok",
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/recharge/b-application/approval", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 0)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid json returns error code", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		req, _ := http.NewRequest("POST", "/recharge/b-application/approval", bytes.NewBufferString("{invalid}"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 400)
	})
}

func TestRechargeHandler_VerifyCard(t *testing.T) {
	t.Run("success with active card", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		card := &model.StoreCard{
			ID:      "card-001",
			CardNo:  "CN001",
			Status:  "active",
			Balance: 1000,
		}
		mockSvc.On("VerifyCard", "CN001").Return(card, nil).Once()

		req, _ := http.NewRequest("GET", "/recharge/cards/verify/CN001", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 0)
		mockSvc.AssertExpectations(t)
	})

	t.Run("card not found", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		mockSvc.On("VerifyCard", "CN999").Return(nil, assert.AnError).Once()

		req, _ := http.NewRequest("GET", "/recharge/cards/verify/CN999", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 404)
		mockSvc.AssertExpectations(t)
	})
}

func TestRechargeHandler_ConsumeCard(t *testing.T) {
	t.Run("success with valid data", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		mockSvc.On("ConsumeCard", "CN001", 100.0, "test consume", "op123").Return(nil).Once()

		body := map[string]interface{}{
			"cardNo": "CN001",
			"amount": 100.0,
			"remark": "test consume",
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/recharge/cards/consume", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 0)
		mockSvc.AssertExpectations(t)
	})

	t.Run("missing cardNo returns params error", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		body := map[string]interface{}{
			"amount": 100.0,
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/recharge/cards/consume", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 400)
	})

	t.Run("missing amount returns params error", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		body := map[string]interface{}{
			"cardNo": "CN001",
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/recharge/cards/consume", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 400)
	})
}

func TestRechargeHandler_GetCardList(t *testing.T) {
	t.Run("success with filters", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		result := map[string]interface{}{
			"items": []interface{}{},
			"total": 0,
			"page":  1,
		}
		mockSvc.On("GetCardList", "active", "CN001", "13800138000", 1, 10).Return(result, nil).Once()

		req, _ := http.NewRequest("GET", "/recharge/cards?status=active&cardNo=CN001&holderPhone=13800138000&page=1&pageSize=10", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 0)
		mockSvc.AssertExpectations(t)
	})

	t.Run("default pagination", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		result := map[string]interface{}{
			"items": []interface{}{},
			"total": 0,
			"page":  1,
		}
		mockSvc.On("GetCardList", "", "", "", 1, 10).Return(result, nil).Once()

		req, _ := http.NewRequest("GET", "/recharge/cards", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 0)
		mockSvc.AssertExpectations(t)
	})
}

func TestRechargeHandler_CreateCenter(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		center := &model.RechargeCenter{
			ID:   "center-001",
			Name: "Test Center",
			Code: "TC001",
		}
		mockSvc.On("CreateCenter", mock.Anything).Return(center, nil).Once()

		body := map[string]interface{}{
			"name":    "Test Center",
			"code":    "TC001",
			"address": "Test Address",
			"phone":   "13800138000",
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/recharge/centers", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 0)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid json returns error code", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		req, _ := http.NewRequest("POST", "/recharge/centers", bytes.NewBufferString("{invalid}"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 400)
	})
}

func TestRechargeHandler_GetDashboardStatistics(t *testing.T) {
	t.Run("success verify response structure", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		req, _ := http.NewRequest("GET", "/recharge/dashboard/statistics", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, float64(0), resp["code"])

		// Verify the data structure
		data, ok := resp["data"].(map[string]interface{})
		assert.True(t, ok, "response should contain data object")
		assert.Contains(t, data, "todayRecharge")
		assert.Contains(t, data, "monthRecharge")
		assert.Contains(t, data, "totalCards")
		assert.Contains(t, data, "pendingApprovals")
	})
}
