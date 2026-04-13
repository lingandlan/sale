//go:build integration
// +build integration

package integration

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"marketplace/backend/internal/model"
	"marketplace/backend/internal/repository"
)

func TestRechargeRepo_CRechargeCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	gormDB, td := setupRechargeTestDB(t)
	defer td.Close()

	repo := repository.NewRechargeRepository(gormDB)

	td.CleanupTables(t, "c_recharges")

	t.Run("Create and GetByID", func(t *testing.T) {
		recharge := &model.CRecharge{
			ID:            "test-c-recharge-1",
			MemberID:      "member-1",
			MemberName:    "张三",
			MemberPhone:   "13800001111",
			CenterID:      "center-1",
			CenterName:    "北京中心",
			Amount:        1000,
			Points:        1000,
			PaymentMethod: "wechat",
			OperatorID:    "op-1",
			OperatorName:  "李出纳",
			BalanceBefore: 0,
			BalanceAfter:  1000,
		}
		err := repo.CreateCRecharge(recharge)
		require.NoError(t, err)

		got, err := repo.GetCRechargeByID("test-c-recharge-1")
		require.NoError(t, err)
		assert.Equal(t, "test-c-recharge-1", got.ID)
		assert.Equal(t, "张三", got.MemberName)
		assert.Equal(t, float64(1000), got.Amount)
	})

	t.Run("GetCRechargeList", func(t *testing.T) {
		list, total, err := repo.GetCRechargeList("", "", 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, list, 1)
	})

	t.Run("GetCRechargeList filter by phone", func(t *testing.T) {
		_, total, err := repo.GetCRechargeList("13800001111", "", 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(1), total)

		// 不匹配的手机号
		_, total, err = repo.GetCRechargeList("19900000000", "", 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(0), total)
	})

	t.Run("GetCRechargeByID not found", func(t *testing.T) {
		_, err := repo.GetCRechargeByID("nonexistent")
		assert.Error(t, err)
	})
}

func TestRechargeRepo_StoreCardCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	gormDB, td := setupRechargeTestDB(t)
	defer td.Close()

	repo := repository.NewRechargeRepository(gormDB)

	td.CleanupTables(t, "card_transactions", "store_cards")

	t.Run("Create and GetByCardNo", func(t *testing.T) {
		card := &model.StoreCard{
			ID:              "test-card-1",
			CardNo:          "TJ26000001",
			HolderID:        "member-1",
			HolderName:      "李四",
			HolderPhone:     "13800002222",
			Balance:         5000,
			Status:          "active",
			IssueCenterID:   "center-1",
			IssueCenterName: "北京中心",
			IssueDate:       time.Now(),
			ExpiryDate:      time.Now().AddDate(1, 0, 0),
		}
		err := repo.CreateCard(card)
		require.NoError(t, err)

		got, err := repo.GetCardByCardNo("TJ26000001")
		require.NoError(t, err)
		assert.Equal(t, "test-card-1", got.ID)
		assert.Equal(t, 5000.0, got.Balance)
	})

	t.Run("UpdateCardBalance", func(t *testing.T) {
		err := repo.UpdateCardBalance("TJ26000001", 3000)
		require.NoError(t, err)

		got, _ := repo.GetCardByCardNo("TJ26000001")
		assert.Equal(t, 3000.0, got.Balance)
	})

	t.Run("CreateCardTransaction", func(t *testing.T) {
		tx := &model.CardTransaction{
			ID:           "test-tx-1",
			CardNo:       "TJ26000001",
			Type:         "consume",
			Amount:       2000,
			BalanceAfter: 3000,
			Remark:       "消费",
			OperatorID:   "op-1",
		}
		err := repo.CreateCardTransaction(tx)
		require.NoError(t, err)

		txs, err := repo.GetCardTransactions("TJ26000001")
		require.NoError(t, err)
		assert.Len(t, txs, 1)
		assert.Equal(t, "consume", txs[0].Type)
	})

	t.Run("UpdateCardStatus", func(t *testing.T) {
		err := repo.UpdateCardStatus("TJ26000001", "inactive")
		require.NoError(t, err)

		got, _ := repo.GetCardByCardNo("TJ26000001")
		assert.Equal(t, "inactive", got.Status)
	})

	t.Run("GetCardList", func(t *testing.T) {
		_, total, err := repo.GetCardList("", "", "", 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(1), total)

		// 按状态筛选
		_, total, err = repo.GetCardList("inactive", "", "", 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(1), total)

		_, total, err = repo.GetCardList("active", "", "", 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(0), total)
	})

	t.Run("GetCardStats", func(t *testing.T) {
		stats, err := repo.GetCardStats()
		require.NoError(t, err)
		assert.Equal(t, int64(1), stats["totalCards"])
		assert.Equal(t, int64(0), stats["activeCards"]) // 刚才改成了 inactive
	})
}

func TestRechargeRepo_RechargeApplicationCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	gormDB, td := setupRechargeTestDB(t)
	defer td.Close()

	repo := repository.NewRechargeRepository(gormDB)

	td.CleanupTables(t, "recharge_applications")

	t.Run("Create and GetByID", func(t *testing.T) {
		app := &model.RechargeApplication{
			ID:            "test-app-1",
			CenterID:      "center-1",
			CenterName:    "北京中心",
			Amount:        10000,
			BasePoints:    10000,
			RebatePoints:  200,
			Points:        10200,
			RebateRate:    2,
			ApplicantID:   "user-1",
			ApplicantName: "张三",
			TransactionNo: "TX-TEST-001",
			Status:        "pending",
		}
		err := repo.CreateRechargeApplication(app)
		require.NoError(t, err)

		got, err := repo.GetRechargeApplicationByID("test-app-1")
		require.NoError(t, err)
		assert.Equal(t, "pending", got.Status)
		assert.Equal(t, float64(10000), got.Amount)
	})

	t.Run("GetRechargeApplications filter", func(t *testing.T) {
		// 按状态筛选
		list, total, err := repo.GetRechargeApplications("pending", 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, list, 1)

		// 不存在的状态
		list, total, err = repo.GetRechargeApplications("approved", 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(0), total)
	})

	t.Run("UpdateRechargeApplicationStatus", func(t *testing.T) {
		err := repo.UpdateRechargeApplicationStatus("test-app-1", "approved", "admin", "通过")
		require.NoError(t, err)

		got, _ := repo.GetRechargeApplicationByID("test-app-1")
		assert.Equal(t, "approved", got.Status)
		assert.Equal(t, "admin", got.ApprovedBy)
	})
}

func TestRechargeRepo_CenterOperatorCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	gormDB, td := setupRechargeTestDB(t)
	defer td.Close()

	repo := repository.NewRechargeRepository(gormDB)

	td.CleanupTables(t, "recharge_operators", "recharge_centers")

	t.Run("CreateCenter and GetCenters", func(t *testing.T) {
		center := &model.RechargeCenter{
			ID:      "center-test-1",
			Name:    "测试中心",
			Code:    "TEST001",
			Address: "北京市",
			Phone:   "010-12345678",
			Status:  "active",
		}
		err := repo.CreateCenter(center)
		require.NoError(t, err)

		centers, err := repo.GetCenters()
		require.NoError(t, err)
		assert.Len(t, centers, 1)
	})

	t.Run("UpdateCenter", func(t *testing.T) {
		center := &model.RechargeCenter{
			ID:      "center-test-1",
			Name:    "更新中心",
			Code:    "TEST001",
			Address: "上海市",
			Phone:   "021-87654321",
			Status:  "active",
		}
		err := repo.UpdateCenter(center)
		require.NoError(t, err)

		centers, _ := repo.GetCenters()
		assert.Equal(t, "更新中心", centers[0].Name)
	})

	t.Run("DeleteCenter", func(t *testing.T) {
		err := repo.DeleteCenter("center-test-1")
		require.NoError(t, err)

		centers, _ := repo.GetCenters()
		assert.Len(t, centers, 0) // GetCenters only returns active
	})

	t.Run("CreateOperator and GetOperators", func(t *testing.T) {
		op := &model.RechargeOperator{
			ID:       "op-test-1",
			Name:     "测试操作员",
			Phone:    "13800009999",
			Password: "hashed",
			CenterID: "center-test-1",
			Role:     "operator",
			Status:   "active",
		}
		err := repo.CreateOperator(op)
		require.NoError(t, err)

		ops, err := repo.GetOperators()
		require.NoError(t, err)
		assert.Len(t, ops, 1)
		assert.Equal(t, "测试操作员", ops[0].Name)
	})

	t.Run("DeleteOperator", func(t *testing.T) {
		err := repo.DeleteOperator("op-test-1")
		require.NoError(t, err)

		ops, _ := repo.GetOperators()
		assert.Len(t, ops, 0)
	})
}
