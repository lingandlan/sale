package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"marketplace/backend/internal/model"
	"marketplace/backend/internal/repository"
)

// ========== Mock ==========

type MockRechargeRepo struct {
	mock.Mock
}

func (m *MockRechargeRepo) CreateRechargeApplication(app *model.RechargeApplication) error {
	args := m.Called(app)
	return args.Error(0)
}

func (m *MockRechargeRepo) GetRechargeApplications(status string, page, pageSize int) ([]model.RechargeApplication, int64, error) {
	args := m.Called(status, page, pageSize)
	return args.Get(0).([]model.RechargeApplication), args.Get(1).(int64), args.Error(2)
}

func (m *MockRechargeRepo) GetRechargeApplicationByID(id string) (*model.RechargeApplication, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.RechargeApplication), args.Error(1)
}

func (m *MockRechargeRepo) UpdateRechargeApplicationStatus(id, status, approvedBy, remark string) error {
	args := m.Called(id, status, approvedBy, remark)
	return args.Error(0)
}

func (m *MockRechargeRepo) CreateCRecharge(recharge *model.CRecharge) error {
	args := m.Called(recharge)
	return args.Error(0)
}

func (m *MockRechargeRepo) GetCRechargeList(memberPhone, centerID string, page, pageSize int) ([]model.CRecharge, int64, error) {
	args := m.Called(memberPhone, centerID, page, pageSize)
	return args.Get(0).([]model.CRecharge), args.Get(1).(int64), args.Error(2)
}

func (m *MockRechargeRepo) GetCRechargeByID(id string) (*model.CRecharge, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CRecharge), args.Error(1)
}

func (m *MockRechargeRepo) CreateCard(card *model.StoreCard) error {
	args := m.Called(card)
	return args.Error(0)
}

func (m *MockRechargeRepo) GetCardList(status, cardNo, holderPhone string, page, pageSize int) ([]model.StoreCard, int64, error) {
	args := m.Called(status, cardNo, holderPhone, page, pageSize)
	return args.Get(0).([]model.StoreCard), args.Get(1).(int64), args.Error(2)
}

func (m *MockRechargeRepo) GetCardByCardNo(cardNo string) (*model.StoreCard, error) {
	args := m.Called(cardNo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.StoreCard), args.Error(1)
}

func (m *MockRechargeRepo) UpdateCardBalance(cardNo string, balance float64) error {
	args := m.Called(cardNo, balance)
	return args.Error(0)
}

func (m *MockRechargeRepo) UpdateCardStatus(cardNo, status string) error {
	args := m.Called(cardNo, status)
	return args.Error(0)
}

func (m *MockRechargeRepo) CreateCardTransaction(transaction *model.CardTransaction) error {
	args := m.Called(transaction)
	return args.Error(0)
}

func (m *MockRechargeRepo) GetCardTransactions(cardNo string) ([]model.CardTransaction, error) {
	args := m.Called(cardNo)
	return args.Get(0).([]model.CardTransaction), args.Error(1)
}

func (m *MockRechargeRepo) GetCardStats() (map[string]int64, error) {
	args := m.Called()
	return args.Get(0).(map[string]int64), args.Error(1)
}

func (m *MockRechargeRepo) GetCenters() ([]model.RechargeCenter, error) {
	args := m.Called()
	return args.Get(0).([]model.RechargeCenter), args.Error(1)
}

func (m *MockRechargeRepo) CreateCenter(center *model.RechargeCenter) error {
	args := m.Called(center)
	return args.Error(0)
}

func (m *MockRechargeRepo) UpdateCenter(id string, updates map[string]interface{}) error {
	args := m.Called(center)
	return args.Error(0)
}

func (m *MockRechargeRepo) DeleteCenter(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRechargeRepo) GetOperators() ([]model.RechargeOperator, error) {
	args := m.Called()
	return args.Get(0).([]model.RechargeOperator), args.Error(1)
}

func (m *MockRechargeRepo) CreateOperator(operator *model.RechargeOperator) error {
	args := m.Called(operator)
	return args.Error(0)
}

func (m *MockRechargeRepo) UpdateOperator(operator *model.RechargeOperator) error {
	args := m.Called(operator)
	return args.Error(0)
}

func (m *MockRechargeRepo) DeleteOperator(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// Verify interface compliance
var _ repository.RechargeRepoInterface = (*MockRechargeRepo)(nil)

func newTestRechargeService(repo *MockRechargeRepo) *RechargeService {
	return NewRechargeService(repo)
}

// ========== Tests ==========

func TestRechargeService_CalculatePoints(t *testing.T) {
	svc := newTestRechargeService(new(MockRechargeRepo))

	tests := []struct {
		name                 string
		amount               float64
		lastMonthConsumption float64
		wantBase             int
		wantRebate           int
		wantTotal            int
	}{
		{"high consumption 2% rebate", 10000, 120000, 10000, 200, 10200},
		{"low consumption 1% rebate", 50000, 50000, 50000, 500, 50500},
		{"boundary exactly 100000", 100000, 100000, 100000, 2000, 102000},
		{"zero amount", 0, 0, 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base, rebate, total := svc.CalculatePoints(tt.amount, tt.lastMonthConsumption)
			assert.Equal(t, tt.wantBase, base)
			assert.Equal(t, tt.wantRebate, rebate)
			assert.Equal(t, tt.wantTotal, total)
		})
	}
}

func TestRechargeService_CreateBRechargeApplication(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(MockRechargeRepo)
		svc := newTestRechargeService(repo)

		repo.On("CreateRechargeApplication", mock.AnythingOfType("*model.RechargeApplication")).Return(nil)

		data := map[string]interface{}{
			"centerId":            "center-1",
			"centerName":          "北京中心",
			"amount":              float64(10000),
			"lastMonthConsumption": float64(120000),
			"applicantId":         "user-1",
			"applicantName":       "张三",
			"transactionNo":       "TX123",
			"screenshot":           "http://img.png",
			"remark":              "test",
		}

		app, err := svc.CreateBRechargeApplication(data)
		require.NoError(t, err)
		assert.NotEmpty(t, app.ID)
		assert.Equal(t, "pending", app.Status)
		assert.Equal(t, "center-1", app.CenterID)
		// Verify points calculated correctly (10000 amount, 120000 consumption => 2%)
		assert.Equal(t, 10000, app.BasePoints)
		assert.Equal(t, 200, app.RebatePoints)
		assert.Equal(t, 10200, app.Points)

		repo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		repo := new(MockRechargeRepo)
		svc := newTestRechargeService(repo)

		repo.On("CreateRechargeApplication", mock.AnythingOfType("*model.RechargeApplication")).Return(assert.AnError)

		data := map[string]interface{}{
			"centerId":            "center-1",
			"centerName":          "北京中心",
			"amount":              float64(10000),
			"lastMonthConsumption": float64(50000),
			"applicantId":         "user-1",
			"applicantName":       "张三",
			"transactionNo":       "TX123",
			"screenshot":           "",
			"remark":              "",
		}

		app, err := svc.CreateBRechargeApplication(data)
		assert.Nil(t, app)
		assert.Error(t, err)

		repo.AssertExpectations(t)
	})
}

func TestRechargeService_ApproveRechargeApplication(t *testing.T) {
	t.Run("approve success", func(t *testing.T) {
		repo := new(MockRechargeRepo)
		svc := newTestRechargeService(repo)

		pendingApp := &model.RechargeApplication{ID: "app-1", Status: "pending"}
		repo.On("GetRechargeApplicationByID", "app-1").Return(pendingApp, nil)
		repo.On("UpdateRechargeApplicationStatus", "app-1", "approved", "admin", "").Return(nil)

		err := svc.ApproveRechargeApplication("app-1", "approve", "admin", "")
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("reject success", func(t *testing.T) {
		repo := new(MockRechargeRepo)
		svc := newTestRechargeService(repo)

		pendingApp := &model.RechargeApplication{ID: "app-1", Status: "pending"}
		repo.On("GetRechargeApplicationByID", "app-1").Return(pendingApp, nil)
		repo.On("UpdateRechargeApplicationStatus", "app-1", "rejected", "admin", "金额有误").Return(nil)

		err := svc.ApproveRechargeApplication("app-1", "reject", "admin", "金额有误")
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("application not found", func(t *testing.T) {
		repo := new(MockRechargeRepo)
		svc := newTestRechargeService(repo)

		repo.On("GetRechargeApplicationByID", "nonexist").Return(nil, assert.AnError)

		err := svc.ApproveRechargeApplication("nonexist", "approve", "admin", "")
		assert.EqualError(t, err, "充值申请不存在")
		repo.AssertExpectations(t)
	})

	t.Run("invalid action", func(t *testing.T) {
		repo := new(MockRechargeRepo)
		svc := newTestRechargeService(repo)

		pendingApp := &model.RechargeApplication{ID: "app-1", Status: "pending"}
		repo.On("GetRechargeApplicationByID", "app-1").Return(pendingApp, nil)

		err := svc.ApproveRechargeApplication("app-1", "invalid", "admin", "")
		assert.EqualError(t, err, "invalid action")
		repo.AssertExpectations(t)
	})
}

func TestRechargeService_VerifyCard(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(MockRechargeRepo)
		svc := newTestRechargeService(repo)

		futureTime := time.Now().AddDate(1, 0, 0)
		card := &model.StoreCard{CardNo: "TJ2600001", Status: "active", ExpiryDate: futureTime}
		repo.On("GetCardByCardNo", "TJ2600001").Return(card, nil)

		result, err := svc.VerifyCard("TJ2600001")
		assert.NoError(t, err)
		assert.Equal(t, "TJ2600001", result.CardNo)
		repo.AssertExpectations(t)
	})

	t.Run("card not found", func(t *testing.T) {
		repo := new(MockRechargeRepo)
		svc := newTestRechargeService(repo)

		repo.On("GetCardByCardNo", "INVALID").Return(nil, assert.AnError)

		result, err := svc.VerifyCard("INVALID")
		assert.Nil(t, result)
		assert.EqualError(t, err, "卡号不存在")
		repo.AssertExpectations(t)
	})

	t.Run("card inactive", func(t *testing.T) {
		repo := new(MockRechargeRepo)
		svc := newTestRechargeService(repo)

		futureTime := time.Now().AddDate(1, 0, 0)
		card := &model.StoreCard{CardNo: "TJ2600002", Status: "inactive", ExpiryDate: futureTime}
		repo.On("GetCardByCardNo", "TJ2600002").Return(card, nil)

		result, err := svc.VerifyCard("TJ2600002")
		assert.Nil(t, result)
		assert.EqualError(t, err, "卡状态异常")
		repo.AssertExpectations(t)
	})

	t.Run("card expired", func(t *testing.T) {
		repo := new(MockRechargeRepo)
		svc := newTestRechargeService(repo)

		pastTime := time.Now().AddDate(-1, 0, 0)
		card := &model.StoreCard{CardNo: "TJ2600003", Status: "active", ExpiryDate: pastTime}
		repo.On("GetCardByCardNo", "TJ2600003").Return(card, nil)

		result, err := svc.VerifyCard("TJ2600003")
		assert.Nil(t, result)
		assert.EqualError(t, err, "卡已过期")
		repo.AssertExpectations(t)
	})
}

func TestRechargeService_ConsumeCard(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(MockRechargeRepo)
		svc := newTestRechargeService(repo)

		futureTime := time.Now().AddDate(1, 0, 0)
		card := &model.StoreCard{CardNo: "TJ2600001", Status: "active", Balance: 5000, ExpiryDate: futureTime}
		repo.On("GetCardByCardNo", "TJ2600001").Return(card, nil)
		repo.On("UpdateCardBalance", "TJ2600001", 3000.0).Return(nil)
		repo.On("CreateCardTransaction", mock.AnythingOfType("*model.CardTransaction")).Return(nil)

		err := svc.ConsumeCard("TJ2600001", 2000, "test consume", "op-1")
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("card not found", func(t *testing.T) {
		repo := new(MockRechargeRepo)
		svc := newTestRechargeService(repo)

		repo.On("GetCardByCardNo", "INVALID").Return(nil, assert.AnError)

		err := svc.ConsumeCard("INVALID", 100, "", "op-1")
		assert.EqualError(t, err, "卡号不存在")
		repo.AssertExpectations(t)
	})

	t.Run("card inactive", func(t *testing.T) {
		repo := new(MockRechargeRepo)
		svc := newTestRechargeService(repo)

		futureTime := time.Now().AddDate(1, 0, 0)
		card := &model.StoreCard{CardNo: "TJ2600001", Status: "inactive", Balance: 5000, ExpiryDate: futureTime}
		repo.On("GetCardByCardNo", "TJ2600001").Return(card, nil)

		err := svc.ConsumeCard("TJ2600001", 100, "", "op-1")
		assert.EqualError(t, err, "卡状态异常，无法核销")
		repo.AssertExpectations(t)
	})

	t.Run("insufficient balance", func(t *testing.T) {
		repo := new(MockRechargeRepo)
		svc := newTestRechargeService(repo)

		futureTime := time.Now().AddDate(1, 0, 0)
		card := &model.StoreCard{CardNo: "TJ2600001", Status: "active", Balance: 100, ExpiryDate: futureTime}
		repo.On("GetCardByCardNo", "TJ2600001").Return(card, nil)

		err := svc.ConsumeCard("TJ2600001", 200, "", "op-1")
		assert.EqualError(t, err, "余额不足")
		repo.AssertExpectations(t)
	})

	t.Run("repo error on update balance", func(t *testing.T) {
		repo := new(MockRechargeRepo)
		svc := newTestRechargeService(repo)

		futureTime := time.Now().AddDate(1, 0, 0)
		card := &model.StoreCard{CardNo: "TJ2600001", Status: "active", Balance: 5000, ExpiryDate: futureTime}
		repo.On("GetCardByCardNo", "TJ2600001").Return(card, nil)
		repo.On("UpdateCardBalance", "TJ2600001", 4000.0).Return(assert.AnError)

		err := svc.ConsumeCard("TJ2600001", 1000, "", "op-1")
		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}

func TestRechargeService_IssueCard(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(MockRechargeRepo)
		svc := newTestRechargeService(repo)

		repo.On("CreateCard", mock.AnythingOfType("*model.StoreCard")).Return(nil)
		repo.On("CreateCardTransaction", mock.AnythingOfType("*model.CardTransaction")).Return(nil)

		data := map[string]interface{}{
			"holderId":    "member-1",
			"holderName":  "李四",
			"holderPhone": "13800001111",
			"amount":      float64(5000),
			"centerId":    "center-1",
			"centerName":  "北京中心",
			"operatorId":  "op-1",
		}

		card, err := svc.IssueCard(data)
		require.NoError(t, err)
		assert.NotEmpty(t, card.ID)
		assert.Contains(t, card.CardNo, "TJ")
		assert.Equal(t, "active", card.Status)
		assert.Equal(t, 5000.0, card.Balance)
		assert.WithinDuration(t, time.Now().AddDate(1, 0, 0), card.ExpiryDate, 5*time.Second)

		repo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		repo := new(MockRechargeRepo)
		svc := newTestRechargeService(repo)

		repo.On("CreateCard", mock.AnythingOfType("*model.StoreCard")).Return(assert.AnError)

		data := map[string]interface{}{
			"holderId":    "member-1",
			"holderName":  "李四",
			"holderPhone": "13800001111",
			"amount":      float64(5000),
			"centerId":    "center-1",
			"centerName":  "北京中心",
			"operatorId":  "op-1",
		}

		card, err := svc.IssueCard(data)
		assert.Nil(t, card)
		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}

func TestRechargeService_UpdateCardStatus(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(MockRechargeRepo)
		svc := newTestRechargeService(repo)

		repo.On("UpdateCardStatus", "TJ2600001", "inactive").Return(nil)

		err := svc.UpdateCardStatus("TJ2600001", "inactive")
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})
}

func TestRechargeService_CreateCenter(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(MockRechargeRepo)
		svc := newTestRechargeService(repo)

		repo.On("CreateCenter", mock.AnythingOfType("*model.RechargeCenter")).Return(nil)

		data := map[string]interface{}{
			"name":    "北京中心",
			"code":    "BJ001",
			"address": "北京市朝阳区",
			"phone":   "010-12345678",
		}

		center, err := svc.CreateCenter(data)
		require.NoError(t, err)
		assert.NotEmpty(t, center.ID)
		assert.Equal(t, "active", center.Status)
		assert.Equal(t, "北京中心", center.Name)

		repo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		repo := new(MockRechargeRepo)
		svc := newTestRechargeService(repo)

		repo.On("CreateCenter", mock.AnythingOfType("*model.RechargeCenter")).Return(assert.AnError)

		data := map[string]interface{}{
			"name":    "北京中心",
			"code":    "BJ001",
			"address": "北京市朝阳区",
			"phone":   "010-12345678",
		}

		center, err := svc.CreateCenter(data)
		assert.Nil(t, center)
		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}

func TestRechargeService_CreateOperator(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(MockRechargeRepo)
		svc := newTestRechargeService(repo)

		repo.On("CreateOperator", mock.AnythingOfType("*model.RechargeOperator")).Return(nil)

		data := map[string]interface{}{
			"name":     "王五",
			"phone":    "13800002222",
			"password": "123456",
			"centerId": "center-1",
			"role":     "operator",
		}

		op, err := svc.CreateOperator(data)
		require.NoError(t, err)
		assert.NotEmpty(t, op.ID)
		assert.Equal(t, "active", op.Status)
		assert.Equal(t, "王五", op.Name)

		repo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		repo := new(MockRechargeRepo)
		svc := newTestRechargeService(repo)

		repo.On("CreateOperator", mock.AnythingOfType("*model.RechargeOperator")).Return(assert.AnError)

		data := map[string]interface{}{
			"name":     "王五",
			"phone":    "13800002222",
			"password": "123456",
			"centerId": "center-1",
			"role":     "operator",
		}

		op, err := svc.CreateOperator(data)
		assert.Nil(t, op)
		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}
