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

func (m *MockRechargeService) GetCRechargeList(memberPhone, centerID, startDate, endDate string, page, pageSize int) (map[string]interface{}, error) {
	args := m.Called(memberPhone, centerID, startDate, endDate, page, pageSize)
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

func (m *MockRechargeService) GetCenters() ([]map[string]interface{}, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]map[string]interface{}), args.Error(1)
}

func (m *MockRechargeService) GetCenterDetail(id string) (*model.RechargeCenter, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.RechargeCenter), args.Error(1)
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

// setupRechargeRouter creates a test router matching actual router.go routes
func setupRechargeRouter(h *RechargeHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	v1 := r.Group("/api/v1")

	// B端充值申请
	bApply := v1.Group("/recharge/b-apply")
	{
		bApply.POST("", h.CreateBRechargeApplication)
	}

	// B端充值审批
	bApproval := v1.Group("/recharge/b-approval")
	{
		bApproval.GET("", h.GetRechargeApplicationList)
		bApproval.GET("/:id", h.GetRechargeApplicationDetail)
		bApproval.POST("/action", h.ApprovalRechargeApplication)
	}

	// C端充值
	cEntry := v1.Group("/recharge/c-entry")
	{
		cEntry.POST("", h.CreateCRecharge)
		cEntry.GET("", h.GetCRechargeList)
		cEntry.GET("/:id", h.GetCRechargeDetail)
	}

	// 充值记录（复用C端handler）
	records := v1.Group("/recharge/records")
	{
		records.GET("", h.GetCRechargeList)
		records.GET("/:id", h.GetCRechargeDetail)
	}

	// 门店卡
	card := v1.Group("/card")
	{
		card.GET("/verify/:cardNo", h.VerifyCard)
		card.POST("/consume", h.ConsumeCard)
		card.GET("/list", h.GetCardList)
		card.GET("/detail/:cardNo", h.GetCardDetail)
		card.GET("/stats", h.GetCardStats)
		card.POST("/issue", h.IssueCard)
		card.POST("/:cardNo/status", h.UpdateCardStatus)
	}

	// 充值中心
	center := v1.Group("/center")
	{
		center.GET("", h.GetCenters)
		center.POST("", h.CreateCenter)
		center.PUT("/:id", h.UpdateCenter)
		center.DELETE("/:id", h.DeleteCenter)
	}

	// 操作员
	operator := v1.Group("/operator")
	{
		operator.GET("", h.GetOperators)
		operator.POST("", h.CreateOperator)
		operator.PUT("/:id", h.UpdateOperator)
		operator.DELETE("/:id", h.DeleteOperator)
	}

	// Dashboard
	dashboard := v1.Group("/dashboard")
	{
		dashboard.GET("/statistics", h.GetDashboardStatistics)
		dashboard.GET("/todos", h.GetDashboardTodos)
		dashboard.GET("/recharge-trends", h.GetDashboardRechargeTrends)
	}

	return r
}

// ========== B端充值申请 ==========

func TestRechargeHandler_CreateBRechargeApplication(t *testing.T) {
	t.Run("success with all required fields", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		app := &model.RechargeApplication{ID: "app-001", Status: "pending"}
		mockSvc.On("CreateBRechargeApplication", mock.Anything).Return(app, nil).Once()

		body := map[string]interface{}{
			"centerId":             "center-001",
			"centerName":           "北京朝阳中心",
			"amount":               50000,
			"lastMonthConsumption": 0,
			"transactionNo":        "BK20260413001",
			"screenshot":           "",
			"remark":               "",
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/api/v1/recharge/b-apply", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 0)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid json returns params error", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		req, _ := http.NewRequest("POST", "/api/v1/recharge/b-apply", bytes.NewBufferString("{invalid}"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 400)
	})

	t.Run("response contains id and status", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		app := &model.RechargeApplication{ID: "app-002", Status: "pending"}
		mockSvc.On("CreateBRechargeApplication", mock.Anything).Return(app, nil).Once()

		body := map[string]interface{}{
			"centerId": "c1", "centerName": "中心", "amount": 1000,
			"lastMonthConsumption": 0, "transactionNo": "T1", "screenshot": "", "remark": "",
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/api/v1/recharge/b-apply", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		data := resp["data"].(map[string]interface{})
		assert.Equal(t, "app-002", data["id"])
		assert.Equal(t, "pending", data["status"])
		mockSvc.AssertExpectations(t)
	})
}

// ========== B端充值审批 ==========

func TestRechargeHandler_GetRechargeApplicationList(t *testing.T) {
	t.Run("success with status filter", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		result := map[string]interface{}{"list": []interface{}{}, "total": 0}
		mockSvc.On("GetRechargeApplicationList", "pending", 1, 10).Return(result, nil).Once()

		req, _ := http.NewRequest("GET", "/api/v1/recharge/b-approval?status=pending&page=1&pageSize=10", nil)
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

		result := map[string]interface{}{"list": []interface{}{}, "total": 0}
		mockSvc.On("GetRechargeApplicationList", "", 1, 10).Return(result, nil).Once()

		req, _ := http.NewRequest("GET", "/api/v1/recharge/b-approval", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 0)
		mockSvc.AssertExpectations(t)
	})
}

func TestRechargeHandler_GetRechargeApplicationDetail(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		app := &model.RechargeApplication{ID: "app-001", CenterName: "北京朝阳中心", Amount: 50000}
		mockSvc.On("GetRechargeApplicationDetail", "app-001").Return(app, nil).Once()

		req, _ := http.NewRequest("GET", "/api/v1/recharge/b-approval/app-001", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 0)
		mockSvc.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		mockSvc.On("GetRechargeApplicationDetail", "not-exist").Return(nil, assert.AnError).Once()

		req, _ := http.NewRequest("GET", "/api/v1/recharge/b-approval/not-exist", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 404)
		mockSvc.AssertExpectations(t)
	})
}

func TestRechargeHandler_ApprovalRechargeApplication(t *testing.T) {
	t.Run("approve success", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		mockSvc.On("ApproveRechargeApplication", "app-001", "approve", "admin", "ok").Return(nil).Once()

		body := map[string]string{"id": "app-001", "action": "approve", "remark": "ok"}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/api/v1/recharge/b-approval/action", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 0)
		mockSvc.AssertExpectations(t)
	})

	t.Run("reject success", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		mockSvc.On("ApproveRechargeApplication", "app-002", "reject", "admin", "金额不对").Return(nil).Once()

		body := map[string]string{"id": "app-002", "action": "reject", "remark": "金额不对"}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/api/v1/recharge/b-approval/action", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 0)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid json returns params error", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		req, _ := http.NewRequest("POST", "/api/v1/recharge/b-approval/action", bytes.NewBufferString("{invalid}"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 400)
	})
}

// ========== C端充值 ==========

func TestRechargeHandler_CreateCRecharge(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		recharge := &model.CRecharge{ID: "rec-001"}
		mockSvc.On("CreateCRecharge", mock.Anything).Return(recharge, nil).Once()

		body := map[string]interface{}{
			"memberId": "m001", "memberName": "张三", "memberPhone": "13900139000",
			"centerId": "c1", "centerName": "北京", "amount": 1000,
			"paymentMethod": "wechat", "remark": "",
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/api/v1/recharge/c-entry", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 0)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid json", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		req, _ := http.NewRequest("POST", "/api/v1/recharge/c-entry", bytes.NewBufferString("{bad}"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assertResponseCode(t, w.Body.Bytes(), 400)
	})
}

func TestRechargeHandler_GetCRechargeList(t *testing.T) {
	t.Run("success with filters", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		result := map[string]interface{}{"list": []interface{}{}, "total": 0}
		mockSvc.On("GetCRechargeList", "13900139000", "c1", "", "", 1, 10).Return(result, nil).Once()

		req, _ := http.NewRequest("GET", "/api/v1/recharge/c-entry?memberPhone=13900139000&centerId=c1&page=1&pageSize=10", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 0)
		mockSvc.AssertExpectations(t)
	})
}

func TestRechargeHandler_GetCRechargeDetail(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		rec := &model.CRecharge{ID: "rec-001", MemberName: "张三"}
		mockSvc.On("GetCRechargeDetail", "rec-001").Return(rec, nil).Once()

		req, _ := http.NewRequest("GET", "/api/v1/recharge/c-entry/rec-001", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 0)
		mockSvc.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		mockSvc.On("GetCRechargeDetail", "not-exist").Return(nil, assert.AnError).Once()

		req, _ := http.NewRequest("GET", "/api/v1/recharge/c-entry/not-exist", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assertResponseCode(t, w.Body.Bytes(), 404)
		mockSvc.AssertExpectations(t)
	})
}

// ========== 充值记录路由（复用C端handler） ==========

func TestRechargeHandler_RecordsList(t *testing.T) {
	t.Run("records route reuses C端 handler", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		result := map[string]interface{}{"list": []interface{}{}, "total": 0}
		mockSvc.On("GetCRechargeList", "", "", "", "", 1, 10).Return(result, nil).Once()

		req, _ := http.NewRequest("GET", "/api/v1/recharge/records", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 0)
		mockSvc.AssertExpectations(t)
	})
}

// ========== 门店卡 ==========

func TestRechargeHandler_IssueCard(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		card := &model.StoreCard{ID: "card-001", CardNo: "TJ2612345"}
		mockSvc.On("IssueCard", mock.Anything).Return(card, nil).Once()

		body := map[string]interface{}{
			"holderId": "m001", "holderName": "张三", "holderPhone": "13900139000",
			"amount": 5000, "centerId": "c1", "centerName": "北京",
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/api/v1/card/issue", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 0)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid json", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		req, _ := http.NewRequest("POST", "/api/v1/card/issue", bytes.NewBufferString("{bad}"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assertResponseCode(t, w.Body.Bytes(), 400)
	})
}

func TestRechargeHandler_VerifyCard(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		card := &model.StoreCard{ID: "card-001", CardNo: "TJ2612345", Status: "active", Balance: 1000}
		mockSvc.On("VerifyCard", "TJ2612345").Return(card, nil).Once()

		req, _ := http.NewRequest("GET", "/api/v1/card/verify/TJ2612345", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 0)
		mockSvc.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		mockSvc.On("VerifyCard", "BAD").Return(nil, assert.AnError).Once()

		req, _ := http.NewRequest("GET", "/api/v1/card/verify/BAD", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assertResponseCode(t, w.Body.Bytes(), 404)
		mockSvc.AssertExpectations(t)
	})
}

func TestRechargeHandler_ConsumeCard(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		mockSvc.On("ConsumeCard", "TJ2612345", 100.0, "消费", "op123").Return(nil).Once()

		body := map[string]interface{}{"cardNo": "TJ2612345", "amount": 100.0, "remark": "消费"}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/api/v1/card/consume", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 0)
		mockSvc.AssertExpectations(t)
	})

	t.Run("missing cardNo", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		body := map[string]interface{}{"amount": 100.0}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/api/v1/card/consume", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assertResponseCode(t, w.Body.Bytes(), 400)
	})

	t.Run("missing amount", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		body := map[string]interface{}{"cardNo": "TJ2612345"}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/api/v1/card/consume", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assertResponseCode(t, w.Body.Bytes(), 400)
	})
}

func TestRechargeHandler_UpdateCardStatus(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		mockSvc.On("UpdateCardStatus", "TJ2612345", "inactive").Return(nil).Once()

		body := map[string]interface{}{"status": "inactive"}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/api/v1/card/TJ2612345/status", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 0)
		mockSvc.AssertExpectations(t)
	})

	t.Run("missing status", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		body := map[string]interface{}{}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/api/v1/card/TJ2612345/status", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assertResponseCode(t, w.Body.Bytes(), 400)
	})
}

func TestRechargeHandler_GetCardList(t *testing.T) {
	t.Run("with filters", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		result := map[string]interface{}{"list": []interface{}{}, "total": 0}
		mockSvc.On("GetCardList", "active", "TJ2612345", "13900139000", 1, 10).Return(result, nil).Once()

		req, _ := http.NewRequest("GET", "/api/v1/card/list?status=active&cardNo=TJ2612345&holderPhone=13900139000", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 0)
		mockSvc.AssertExpectations(t)
	})
}

func TestRechargeHandler_GetCardDetail(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		result := map[string]interface{}{"card": map[string]interface{}{}, "transactions": []interface{}{}}
		mockSvc.On("GetCardDetail", "TJ2612345").Return(result, nil).Once()

		req, _ := http.NewRequest("GET", "/api/v1/card/detail/TJ2612345", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 0)
		mockSvc.AssertExpectations(t)
	})
}

func TestRechargeHandler_GetCardStats(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		result := map[string]interface{}{"totalCards": 100, "activeCards": 80}
		mockSvc.On("GetCardStats").Return(result, nil).Once()

		req, _ := http.NewRequest("GET", "/api/v1/card/stats", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 0)
		mockSvc.AssertExpectations(t)
	})
}

// ========== 充值中心 ==========

func TestRechargeHandler_GetCenters(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		centers := []map[string]interface{}{{"id": "c1", "name": "北京"}}
		mockSvc.On("GetCenters").Return(centers, nil).Once()

		req, _ := http.NewRequest("GET", "/api/v1/center", nil)
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

		center := &model.RechargeCenter{ID: "c1", Name: "北京", Code: "BJ"}
		mockSvc.On("CreateCenter", mock.Anything).Return(center, nil).Once()

		body := map[string]interface{}{"name": "北京", "code": "BJ", "address": "xxx", "phone": "010-1234"}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/api/v1/center", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 0)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid json", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		req, _ := http.NewRequest("POST", "/api/v1/center", bytes.NewBufferString("{bad}"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assertResponseCode(t, w.Body.Bytes(), 400)
	})
}

func TestRechargeHandler_UpdateCenter(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		center := &model.RechargeCenter{ID: "c1", Name: "北京(更新)"}
		mockSvc.On("UpdateCenter", "c1", mock.Anything).Return(center, nil).Once()

		body := map[string]interface{}{"name": "北京(更新)", "code": "BJ", "address": "xxx", "phone": "010", "status": "active"}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("PUT", "/api/v1/center/c1", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 0)
		mockSvc.AssertExpectations(t)
	})
}

func TestRechargeHandler_DeleteCenter(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		mockSvc.On("DeleteCenter", "c1").Return(nil).Once()

		req, _ := http.NewRequest("DELETE", "/api/v1/center/c1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 0)
		mockSvc.AssertExpectations(t)
	})
}

// ========== 操作员 ==========

func TestRechargeHandler_GetOperators(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		ops := []model.RechargeOperator{{ID: "op1", Name: "小李"}}
		mockSvc.On("GetOperators").Return(ops, nil).Once()

		req, _ := http.NewRequest("GET", "/api/v1/operator", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 0)
		mockSvc.AssertExpectations(t)
	})
}

func TestRechargeHandler_CreateOperator(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		op := &model.RechargeOperator{ID: "op1", Name: "小李"}
		mockSvc.On("CreateOperator", mock.Anything).Return(op, nil).Once()

		body := map[string]interface{}{"name": "小李", "phone": "138", "password": "123", "centerId": "c1", "role": "operator"}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/api/v1/operator", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 0)
		mockSvc.AssertExpectations(t)
	})
}

func TestRechargeHandler_UpdateOperator(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		op := &model.RechargeOperator{ID: "op1", Name: "小李(更新)"}
		mockSvc.On("UpdateOperator", "op1", mock.Anything).Return(op, nil).Once()

		body := map[string]interface{}{"name": "小李(更新)", "phone": "138", "centerId": "c1", "role": "admin", "status": "active"}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("PUT", "/api/v1/operator/op1", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 0)
		mockSvc.AssertExpectations(t)
	})
}

func TestRechargeHandler_DeleteOperator(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		mockSvc.On("DeleteOperator", "op1").Return(nil).Once()

		req, _ := http.NewRequest("DELETE", "/api/v1/operator/op1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assertResponseCode(t, w.Body.Bytes(), 0)
		mockSvc.AssertExpectations(t)
	})
}

// ========== Dashboard ==========

func TestRechargeHandler_GetDashboardStatistics(t *testing.T) {
	t.Run("verify response structure", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		req, _ := http.NewRequest("GET", "/api/v1/dashboard/statistics", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, float64(0), resp["code"])

		data := resp["data"].(map[string]interface{})
		assert.Contains(t, data, "todayRecharge")
		assert.Contains(t, data, "monthRecharge")
		assert.Contains(t, data, "totalCards")
		assert.Contains(t, data, "pendingApprovals")
	})
}

func TestRechargeHandler_GetDashboardTodos(t *testing.T) {
	t.Run("verify response structure", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		req, _ := http.NewRequest("GET", "/api/v1/dashboard/todos", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		data := resp["data"].(map[string]interface{})
		assert.Contains(t, data, "pendingApprovals")
		assert.Contains(t, data, "expiringCards")
	})
}

func TestRechargeHandler_GetDashboardRechargeTrends(t *testing.T) {
	t.Run("verify response structure", func(t *testing.T) {
		mockSvc := new(MockRechargeService)
		h := NewRechargeHandler(mockSvc)
		router := setupRechargeRouter(h)

		req, _ := http.NewRequest("GET", "/api/v1/dashboard/recharge-trends?days=7", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		data := resp["data"].(map[string]interface{})
		assert.Contains(t, data, "dates")
		assert.Contains(t, data, "values")
	})
}
