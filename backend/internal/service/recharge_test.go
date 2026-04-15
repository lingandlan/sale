package service

import (
	"errors"
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

// 门店卡 mock methods
func (m *MockRechargeRepo) CreateCard(card *model.StoreCard) error {
	args := m.Called(card)
	return args.Error(0)
}

func (m *MockRechargeRepo) BatchCreateCards(cards []*model.StoreCard) error {
	args := m.Called(cards)
	return args.Error(0)
}

func (m *MockRechargeRepo) GetCardList(status int, cardNo, centerID string, page, pageSize int) ([]model.StoreCard, int64, error) {
	args := m.Called(status, cardNo, centerID, page, pageSize)
	return args.Get(0).([]model.StoreCard), args.Get(1).(int64), args.Error(2)
}

func (m *MockRechargeRepo) GetCardByCardNo(cardNo string) (*model.StoreCard, error) {
	args := m.Called(cardNo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.StoreCard), args.Error(1)
}

func (m *MockRechargeRepo) GetMaxCardSequence() (int, error) {
	args := m.Called()
	return args.Get(0).(int), args.Error(1)
}

func (m *MockRechargeRepo) UpdateCardByMap(cardNo string, updates map[string]interface{}) error {
	args := m.Called(cardNo, updates)
	return args.Error(0)
}

func (m *MockRechargeRepo) AllocateCardsToCenter(centerID, startCardNo, endCardNo string) (int, error) {
	args := m.Called(centerID, startCardNo, endCardNo)
	return args.Get(0).(int), args.Error(1)
}

func (m *MockRechargeRepo) BindCardToUser(cardNo string, updates map[string]interface{}, record *model.CardIssueRecord) error {
	args := m.Called(cardNo, updates, record)
	return args.Error(0)
}

func (m *MockRechargeRepo) ConsumeCardInTx(cardNo string, amount int, operatorID, remark string) error {
	args := m.Called(cardNo, amount, operatorID, remark)
	return args.Error(0)
}

func (m *MockRechargeRepo) CreateCardTransaction(transaction *model.CardTransaction) error {
	args := m.Called(transaction)
	return args.Error(0)
}

func (m *MockRechargeRepo) GetCardTransactions(cardNo string, page, pageSize int) ([]model.CardTransaction, int64, error) {
	args := m.Called(cardNo, page, pageSize)
	return args.Get(0).([]model.CardTransaction), args.Get(1).(int64), args.Error(2)
}

func (m *MockRechargeRepo) GetCardStats() (map[string]int64, error) {
	args := m.Called()
	return args.Get(0).(map[string]int64), args.Error(1)
}

func (m *MockRechargeRepo) GetCardInventoryStats() (map[string]int64, error) {
	args := m.Called()
	return args.Get(0).(map[string]int64), args.Error(1)
}

// 充值中心 mock methods
func (m *MockRechargeRepo) GetCenterByID(id string) (*model.RechargeCenter, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.RechargeCenter), args.Error(1)
}

func (m *MockRechargeRepo) AddCenterBalance(id string, amount float64) error {
	args := m.Called(id, amount)
	return args.Error(0)
}

func (m *MockRechargeRepo) DeductCenterBalance(id string, amount float64) (float64, error) {
	args := m.Called(id, amount)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockRechargeRepo) GetCenterTotalRecharge(centerID string) int64 {
	args := m.Called(centerID)
	return args.Get(0).(int64)
}

func (m *MockRechargeRepo) GetCenterTotalConsumed(centerID string) float64 {
	args := m.Called(centerID)
	return args.Get(0).(float64)
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
	args := m.Called(id, updates)
	return args.Error(0)
}

func (m *MockRechargeRepo) DeleteCenter(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// 操作员 mock methods
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

		pendingApp := &model.RechargeApplication{ID: "app-1", Status: "pending", CenterID: "center-1", Points: 100}
		repo.On("GetRechargeApplicationByID", "app-1").Return(pendingApp, nil)
		repo.On("AddCenterBalance", "center-1", float64(100)).Return(nil)
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
	t.Run("success issued card", func(t *testing.T) {
		repo := new(MockRechargeRepo)
		svc := newTestRechargeService(repo)

		card := &model.StoreCard{CardNo: "TJ2600001", Status: model.CardStatusIssued, Balance: 1000}
		repo.On("GetCardByCardNo", "TJ2600001").Return(card, nil)

		result, err := svc.VerifyCard("TJ2600001")
		assert.NoError(t, err)
		assert.Equal(t, "TJ2600001", result.CardNo)
		repo.AssertExpectations(t)
	})

	t.Run("success active card", func(t *testing.T) {
		repo := new(MockRechargeRepo)
		svc := newTestRechargeService(repo)

		futureTime := time.Now().AddDate(1, 0, 0)
		card := &model.StoreCard{CardNo: "TJ2600002", Status: model.CardStatusActive, Balance: 1000, ExpiredAt: &futureTime}
		repo.On("GetCardByCardNo", "TJ2600002").Return(card, nil)

		result, err := svc.VerifyCard("TJ2600002")
		assert.NoError(t, err)
		assert.Equal(t, "TJ2600002", result.CardNo)
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

	t.Run("card frozen", func(t *testing.T) {
		repo := new(MockRechargeRepo)
		svc := newTestRechargeService(repo)

		card := &model.StoreCard{CardNo: "TJ2600003", Status: model.CardStatusFrozen}
		repo.On("GetCardByCardNo", "TJ2600003").Return(card, nil)

		result, err := svc.VerifyCard("TJ2600003")
		assert.Nil(t, result)
		assert.EqualError(t, err, "卡已冻结")
		repo.AssertExpectations(t)
	})

	t.Run("card in stock", func(t *testing.T) {
		repo := new(MockRechargeRepo)
		svc := newTestRechargeService(repo)

		card := &model.StoreCard{CardNo: "TJ2600004", Status: model.CardStatusInStock}
		repo.On("GetCardByCardNo", "TJ2600004").Return(card, nil)

		result, err := svc.VerifyCard("TJ2600004")
		assert.Nil(t, result)
		assert.EqualError(t, err, "卡未发放")
		repo.AssertExpectations(t)
	})

	t.Run("card expired", func(t *testing.T) {
		repo := new(MockRechargeRepo)
		svc := newTestRechargeService(repo)

		pastTime := time.Now().AddDate(-1, 0, 0)
		card := &model.StoreCard{CardNo: "TJ2600005", Status: model.CardStatusActive, ExpiredAt: &pastTime}
		repo.On("GetCardByCardNo", "TJ2600005").Return(card, nil)
		repo.On("UpdateCardByMap", "TJ2600005", mock.AnythingOfType("map[string]interface {}")).Return(nil)

		result, err := svc.VerifyCard("TJ2600005")
		assert.Nil(t, result)
		assert.EqualError(t, err, "卡已过期")
		repo.AssertExpectations(t)
	})
}

func TestRechargeService_ConsumeCard(t *testing.T) {
	t.Run("success - delegate to repo", func(t *testing.T) {
		repo := new(MockRechargeRepo)
		svc := newTestRechargeService(repo)

		repo.On("ConsumeCardInTx", "TJ2600001", 200, "op-1", "test consume").Return(nil)

		err := svc.ConsumeCard("TJ2600001", 200, "op-1", "test consume")
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := new(MockRechargeRepo)
		svc := newTestRechargeService(repo)

		repo.On("ConsumeCardInTx", "TJ2600001", 100, "op-1", "").Return(errors.New("余额不足"))

		err := svc.ConsumeCard("TJ2600001", 100, "op-1", "")
		assert.EqualError(t, err, "余额不足")
		repo.AssertExpectations(t)
	})
}

func TestRechargeService_BatchImportCards(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(MockRechargeRepo)
		svc := newTestRechargeService(repo)

		repo.On("GetMaxCardSequence").Return(0, nil)
		repo.On("BatchCreateCards", mock.AnythingOfType("[]*model.StoreCard")).Return(nil)
		repo.On("CreateCardTransaction", mock.AnythingOfType("*model.CardTransaction")).Return(nil)

		cardNos, err := svc.BatchImportCards(1, 5, model.CardTypePhysical, "op-1")
		require.NoError(t, err)
		assert.Len(t, cardNos, 5)
		assert.Equal(t, "TJ00000001", cardNos[0])
		assert.Equal(t, "TJ00000005", cardNos[4])

		repo.AssertExpectations(t)
	})

	t.Run("invalid sequence range", func(t *testing.T) {
		repo := new(MockRechargeRepo)
		svc := newTestRechargeService(repo)

		cardNos, err := svc.BatchImportCards(5, 1, model.CardTypePhysical, "op-1")
		assert.Nil(t, cardNos)
		assert.EqualError(t, err, "序号范围无效")
	})

	t.Run("sequence conflict", func(t *testing.T) {
		repo := new(MockRechargeRepo)
		svc := newTestRechargeService(repo)

		repo.On("GetMaxCardSequence").Return(10, nil)

		cardNos, err := svc.BatchImportCards(5, 10, model.CardTypePhysical, "op-1")
		assert.Nil(t, cardNos)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "冲突")

		repo.AssertExpectations(t)
	})

	t.Run("too many cards", func(t *testing.T) {
		repo := new(MockRechargeRepo)
		svc := newTestRechargeService(repo)

		cardNos, err := svc.BatchImportCards(1, 1002, model.CardTypePhysical, "op-1")
		assert.Nil(t, cardNos)
		assert.EqualError(t, err, "单次入库不能超过1000张")
	})
}

func TestRechargeService_FreezeCard(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(MockRechargeRepo)
		svc := newTestRechargeService(repo)

		card := &model.StoreCard{CardNo: "TJ2600001", Status: model.CardStatusIssued, Balance: 1000}
		repo.On("GetCardByCardNo", "TJ2600001").Return(card, nil)
		repo.On("UpdateCardByMap", "TJ2600001", mock.AnythingOfType("map[string]interface {}")).Return(nil)
		repo.On("CreateCardTransaction", mock.AnythingOfType("*model.CardTransaction")).Return(nil)

		err := svc.FreezeCard("TJ2600001", "op-1")
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("card not found", func(t *testing.T) {
		repo := new(MockRechargeRepo)
		svc := newTestRechargeService(repo)

		repo.On("GetCardByCardNo", "INVALID").Return(nil, assert.AnError)

		err := svc.FreezeCard("INVALID", "op-1")
		assert.EqualError(t, err, "卡号不存在")
		repo.AssertExpectations(t)
	})
}

func TestRechargeService_UnfreezeCard(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(MockRechargeRepo)
		svc := newTestRechargeService(repo)

		card := &model.StoreCard{CardNo: "TJ2600001", Status: model.CardStatusFrozen, Balance: 1000}
		repo.On("GetCardByCardNo", "TJ2600001").Return(card, nil)
		repo.On("UpdateCardByMap", "TJ2600001", mock.AnythingOfType("map[string]interface {}")).Return(nil)
		repo.On("CreateCardTransaction", mock.AnythingOfType("*model.CardTransaction")).Return(nil)

		err := svc.UnfreezeCard("TJ2600001", "op-1")
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})
}

func TestRechargeService_VoidCard(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(MockRechargeRepo)
		svc := newTestRechargeService(repo)

		card := &model.StoreCard{CardNo: "TJ2600001", Status: model.CardStatusInStock, Balance: 1000}
		repo.On("GetCardByCardNo", "TJ2600001").Return(card, nil)
		repo.On("UpdateCardByMap", "TJ2600001", mock.AnythingOfType("map[string]interface {}")).Return(nil)
		repo.On("CreateCardTransaction", mock.AnythingOfType("*model.CardTransaction")).Return(nil)

		err := svc.VoidCard("TJ2600001", "op-1")
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})
}

func TestRechargeService_GetCardStats(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(MockRechargeRepo)
		svc := newTestRechargeService(repo)

		stats := map[string]int64{
			"totalCards":  10,
			"activeCards": 5,
		}
		repo.On("GetCardStats").Return(stats, nil)

		result, err := svc.GetCardStats()
		require.NoError(t, err)
		assert.Equal(t, int64(10), result["totalCards"])
		assert.Equal(t, int64(5), result["activeCards"])
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
