package repository

import (
	"fmt"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"marketplace/backend/internal/model"
)

func setupRechargeTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(
		&model.RechargeApplication{},
		&model.CRecharge{},
		&model.StoreCard{},
		&model.CardTransaction{},
		&model.RechargeCenter{},
		&model.RechargeOperator{},
	)
	require.NoError(t, err)

	return db
}

// ========== B端充值申请 ==========

func TestRechargeRepository_BRechargeApplication(t *testing.T) {
	db := setupRechargeTestDB(t)
	repo := NewRechargeRepository(db)

	t.Run("create and get by ID", func(t *testing.T) {
		app := &model.RechargeApplication{
			ID:             "app_001",
			CenterID:       "center_1",
			CenterName:     "Test Center",
			Amount:         1000.0,
			Points:         1000,
			BasePoints:     800,
			RebatePoints:   200,
			RebateRate:     25,
			ApplicantID:    "user_1",
			ApplicantName:  "张三",
			TransactionNo:  "TXN_20260411_001",
			Screenshot:     "https://example.com/screenshot.png",
			Remark:         "测试充值申请",
			Status:         "pending",
		}

		err := repo.CreateRechargeApplication(app)
		require.NoError(t, err)

		result, err := repo.GetRechargeApplicationByID("app_001")
		require.NoError(t, err)
		assert.Equal(t, "app_001", result.ID)
		assert.Equal(t, "center_1", result.CenterID)
		assert.Equal(t, "Test Center", result.CenterName)
		assert.Equal(t, 1000.0, result.Amount)
		assert.Equal(t, 1000, result.Points)
		assert.Equal(t, "pending", result.Status)
		assert.Equal(t, "张三", result.ApplicantName)
	})

	t.Run("get applications with status filter", func(t *testing.T) {
		pendingApp := &model.RechargeApplication{
			ID:            "app_pending",
			CenterID:      "center_1",
			CenterName:    "Center A",
			Amount:        500.0,
			Points:        500,
			ApplicantID:   "user_2",
			ApplicantName: "李四",
			Status:        "pending",
		}
		approvedApp := &model.RechargeApplication{
			ID:            "app_approved",
			CenterID:      "center_2",
			CenterName:    "Center B",
			Amount:        2000.0,
			Points:        2000,
			ApplicantID:   "user_3",
			ApplicantName: "王五",
			Status:        "approved",
		}

		err := repo.CreateRechargeApplication(pendingApp)
		require.NoError(t, err)
		err = repo.CreateRechargeApplication(approvedApp)
		require.NoError(t, err)

		// Filter by pending
		list, total, err := repo.GetRechargeApplications("pending", 1, 10)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, total, int64(1))
		for _, item := range list {
			assert.Equal(t, "pending", item.Status)
		}

		// Filter by approved
		list, total, err = repo.GetRechargeApplications("approved", 1, 10)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, total, int64(1))
		for _, item := range list {
			assert.Equal(t, "approved", item.Status)
		}

		// No filter - returns all
		_, total, err = repo.GetRechargeApplications("", 1, 10)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, total, int64(2))
	})

	t.Run("update application status", func(t *testing.T) {
		app := &model.RechargeApplication{
			ID:            "app_status",
			CenterID:      "center_1",
			CenterName:    "Center A",
			Amount:        3000.0,
			Points:        3000,
			ApplicantID:   "user_4",
			ApplicantName: "赵六",
			Status:        "pending",
		}
		err := repo.CreateRechargeApplication(app)
		require.NoError(t, err)

		err = repo.UpdateRechargeApplicationStatus("app_status", "approved", "admin_1", "审核通过")
		require.NoError(t, err)

		result, err := repo.GetRechargeApplicationByID("app_status")
		require.NoError(t, err)
		assert.Equal(t, "approved", result.Status)
		assert.Equal(t, "admin_1", result.ApprovedBy)
		assert.Equal(t, "审核通过", result.ApprovalRemark)
		assert.NotNil(t, result.ApprovedAt)
	})
}

// ========== C端充值 ==========

func TestRechargeRepository_CRecharge(t *testing.T) {
	db := setupRechargeTestDB(t)
	repo := NewRechargeRepository(db)

	t.Run("create and get by ID", func(t *testing.T) {
		recharge := &model.CRecharge{
			ID:            "crecharge_001",
			MemberID:      "member_1",
			MemberName:    "张三",
			MemberPhone:   "13800138000",
			CenterID:      "center_1",
			CenterName:    "Center A",
			Amount:        500.0,
			Points:        500,
			PaymentMethod: "wechat",
			OperatorID:    "op_1",
			OperatorName:  "操作员A",
			Remark:        "首次充值",
			BalanceBefore: 0,
			BalanceAfter:  500,
		}

		err := repo.CreateCRecharge(recharge)
		require.NoError(t, err)

		result, err := repo.GetCRechargeByID("crecharge_001")
		require.NoError(t, err)
		assert.Equal(t, "crecharge_001", result.ID)
		assert.Equal(t, "member_1", result.MemberID)
		assert.Equal(t, "13800138000", result.MemberPhone)
		assert.Equal(t, 500.0, result.Amount)
		assert.Equal(t, 500, result.Points)
		assert.Equal(t, "wechat", result.PaymentMethod)
	})

	t.Run("get list with phone filter", func(t *testing.T) {
		r1 := &model.CRecharge{
			ID:            "crecharge_phone1",
			MemberID:      "member_2",
			MemberName:    "李四",
			MemberPhone:   "13800138001",
			CenterID:      "center_1",
			CenterName:    "Center A",
			Amount:        1000.0,
			Points:        1000,
			PaymentMethod: "cash",
			BalanceBefore: 0,
			BalanceAfter:  1000,
		}
		r2 := &model.CRecharge{
			ID:            "crecharge_phone2",
			MemberID:      "member_3",
			MemberName:    "王五",
			MemberPhone:   "13800138002",
			CenterID:      "center_2",
			CenterName:    "Center B",
			Amount:        2000.0,
			Points:        2000,
			PaymentMethod: "alipay",
			BalanceBefore: 0,
			BalanceAfter:  2000,
		}

		err := repo.CreateCRecharge(r1)
		require.NoError(t, err)
		err = repo.CreateCRecharge(r2)
		require.NoError(t, err)

		// Filter by phone
		list, total, err := repo.GetCRechargeList("13800138001", "", 1, 10)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, total, int64(1))
		for _, item := range list {
			assert.Equal(t, "13800138001", item.MemberPhone)
		}

		// Filter by centerID
		list, total, err = repo.GetCRechargeList("", "center_2", 1, 10)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, total, int64(1))
		for _, item := range list {
			assert.Equal(t, "center_2", item.CenterID)
		}

		// No filter
		_, total, err = repo.GetCRechargeList("", "", 1, 10)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, total, int64(2))
	})
}

// ========== 门店卡 ==========

func TestRechargeRepository_StoreCard(t *testing.T) {
	db := setupRechargeTestDB(t)
	repo := NewRechargeRepository(db)
	now := time.Now()
	future := now.Add(365 * 24 * time.Hour)

	t.Run("create and get by card no", func(t *testing.T) {
		card := &model.StoreCard{
			ID:              "card_001",
			CardNo:          "SC_2026_0001",
			HolderID:        "holder_1",
			HolderName:      "张三",
			HolderPhone:     "13800138000",
			Balance:         5000.0,
			Status:          "active",
			IssueCenterID:   "center_1",
			IssueCenterName: "Center A",
			IssueDate:       now,
			ExpiryDate:      future,
		}

		err := repo.CreateCard(card)
		require.NoError(t, err)

		result, err := repo.GetCardByCardNo("SC_2026_0001")
		require.NoError(t, err)
		assert.Equal(t, "card_001", result.ID)
		assert.Equal(t, "SC_2026_0001", result.CardNo)
		assert.Equal(t, "张三", result.HolderName)
		assert.Equal(t, 5000.0, result.Balance)
		assert.Equal(t, "active", result.Status)
	})

	t.Run("get card list with status filter", func(t *testing.T) {
		activeCard := &model.StoreCard{
			ID:              "card_active",
			CardNo:          fmt.Sprintf("SC_ACTIVE_%d", now.UnixNano()),
			HolderID:        "holder_2",
			HolderName:      "李四",
			HolderPhone:     "13800138001",
			Balance:         3000.0,
			Status:          "active",
			IssueCenterID:   "center_1",
			IssueCenterName: "Center A",
			IssueDate:       now,
			ExpiryDate:      future,
		}
		inactiveCard := &model.StoreCard{
			ID:              "card_inactive",
			CardNo:          fmt.Sprintf("SC_INACTIVE_%d", now.UnixNano()),
			HolderID:        "holder_3",
			HolderName:      "王五",
			HolderPhone:     "13800138002",
			Balance:         0,
			Status:          "inactive",
			IssueCenterID:   "center_1",
			IssueCenterName: "Center A",
			IssueDate:       now,
			ExpiryDate:      future,
		}

		err := repo.CreateCard(activeCard)
		require.NoError(t, err)
		err = repo.CreateCard(inactiveCard)
		require.NoError(t, err)

		// Filter by active
		list, total, err := repo.GetCardList("active", "", "", 1, 10)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, total, int64(1))
		for _, item := range list {
			assert.Equal(t, "active", item.Status)
		}

		// Filter by inactive
		list, total, err = repo.GetCardList("inactive", "", "", 1, 10)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, total, int64(1))
		for _, item := range list {
			assert.Equal(t, "inactive", item.Status)
		}
	})

	t.Run("update card balance", func(t *testing.T) {
		card := &model.StoreCard{
			ID:              "card_balance",
			CardNo:          "SC_BALANCE_001",
			HolderID:        "holder_4",
			HolderName:      "赵六",
			HolderPhone:     "13800138003",
			Balance:         5000.0,
			Status:          "active",
			IssueCenterID:   "center_1",
			IssueCenterName: "Center A",
			IssueDate:       now,
			ExpiryDate:      future,
		}
		err := repo.CreateCard(card)
		require.NoError(t, err)

		err = repo.UpdateCardBalance("SC_BALANCE_001", 3000.0)
		require.NoError(t, err)

		result, err := repo.GetCardByCardNo("SC_BALANCE_001")
		require.NoError(t, err)
		assert.Equal(t, 3000.0, result.Balance)
	})

	t.Run("update card status", func(t *testing.T) {
		card := &model.StoreCard{
			ID:              "card_status",
			CardNo:          "SC_STATUS_001",
			HolderID:        "holder_5",
			HolderName:      "孙七",
			HolderPhone:     "13800138004",
			Balance:         1000.0,
			Status:          "active",
			IssueCenterID:   "center_1",
			IssueCenterName: "Center A",
			IssueDate:       now,
			ExpiryDate:      future,
		}
		err := repo.CreateCard(card)
		require.NoError(t, err)

		err = repo.UpdateCardStatus("SC_STATUS_001", "inactive")
		require.NoError(t, err)

		result, err := repo.GetCardByCardNo("SC_STATUS_001")
		require.NoError(t, err)
		assert.Equal(t, "inactive", result.Status)
	})

	t.Run("get card stats", func(t *testing.T) {
		// Create cards with different statuses
		cards := []*model.StoreCard{
			{
				ID:              "stat_card_1",
				CardNo:          "STAT_001",
				HolderID:        "h_1",
				HolderName:      "Stat1",
				HolderPhone:     "13900000001",
				Balance:         1000.0,
				Status:          "active",
				IssueCenterID:   "center_1",
				IssueCenterName: "Center A",
				IssueDate:       now,
				ExpiryDate:      future,
			},
			{
				ID:              "stat_card_2",
				CardNo:          "STAT_002",
				HolderID:        "h_2",
				HolderName:      "Stat2",
				HolderPhone:     "13900000002",
				Balance:         2000.0,
				Status:          "active",
				IssueCenterID:   "center_1",
				IssueCenterName: "Center A",
				IssueDate:       now,
				ExpiryDate:      future,
			},
			{
				ID:              "stat_card_3",
				CardNo:          "STAT_003",
				HolderID:        "h_3",
				HolderName:      "Stat3",
				HolderPhone:     "13900000003",
				Balance:         0,
				Status:          "inactive",
				IssueCenterID:   "center_1",
				IssueCenterName: "Center A",
				IssueDate:       now,
				ExpiryDate:      future,
			},
			{
				ID:              "stat_card_4",
				CardNo:          "STAT_004",
				HolderID:        "h_4",
				HolderName:      "Stat4",
				HolderPhone:     "13900000004",
				Balance:         500.0,
				Status:          "expired",
				IssueCenterID:   "center_1",
				IssueCenterName: "Center A",
				IssueDate:       now,
				ExpiryDate:      now.Add(-24 * time.Hour), // already expired
			},
		}
		for _, c := range cards {
			err := repo.CreateCard(c)
			require.NoError(t, err)
		}

		stats, err := repo.GetCardStats()
		require.NoError(t, err)
		assert.GreaterOrEqual(t, stats["totalCards"], int64(4))
		assert.GreaterOrEqual(t, stats["activeCards"], int64(2))
		assert.GreaterOrEqual(t, stats["frozenCards"], int64(1))
		assert.GreaterOrEqual(t, stats["expiredCards"], int64(1))
		// Total balance should include at least 1000+2000+500 = 3500
		assert.GreaterOrEqual(t, stats["totalBalance"], int64(3500))
	})
}

// ========== 门店卡交易 ==========

func TestRechargeRepository_CardTransaction(t *testing.T) {
	db := setupRechargeTestDB(t)
	repo := NewRechargeRepository(db)

	t.Run("create and get transactions ordered DESC", func(t *testing.T) {
		tx1 := &model.CardTransaction{
			ID:           "txn_001",
			CardNo:       "SC_TXN_001",
			Type:         "issue",
			Amount:       5000.0,
			BalanceAfter: 5000.0,
			Remark:       "开卡充值",
			OperatorID:   "op_1",
		}
		tx2 := &model.CardTransaction{
			ID:           "txn_002",
			CardNo:       "SC_TXN_001",
			Type:         "consume",
			Amount:       1000.0,
			BalanceAfter: 4000.0,
			Remark:       "消费扣款",
			OperatorID:   "op_1",
		}

		err := repo.CreateCardTransaction(tx1)
		require.NoError(t, err)
		err = repo.CreateCardTransaction(tx2)
		require.NoError(t, err)

		list, err := repo.GetCardTransactions("SC_TXN_001")
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(list), 2)

		// Verify ordered DESC by created_at (most recent first)
		for i := 1; i < len(list); i++ {
			assert.True(t, !list[i].CreatedAt.After(list[i-1].CreatedAt),
				"transactions should be ordered by created_at DESC")
		}
	})
}

// ========== 充值中心 ==========

func TestRechargeRepository_Center(t *testing.T) {
	db := setupRechargeTestDB(t)
	repo := NewRechargeRepository(db)

	t.Run("create and get centers (only active)", func(t *testing.T) {
		activeCenter := &model.RechargeCenter{
			ID:      "center_active",
			Name:    "Active Center",
			Code:    "AC_001",
			Address: "地址A",
			Phone:   "010_12345678",
			Status:  "active",
		}
		inactiveCenter := &model.RechargeCenter{
			ID:      "center_inactive",
			Name:    "Inactive Center",
			Code:    "IC_001",
			Address: "地址B",
			Phone:   "010_87654321",
			Status:  "inactive",
		}

		err := repo.CreateCenter(activeCenter)
		require.NoError(t, err)
		err = repo.CreateCenter(inactiveCenter)
		require.NoError(t, err)

		list, err := repo.GetCenters()
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(list), 1)
		for _, c := range list {
			assert.Equal(t, "active", c.Status)
		}
	})

	t.Run("update center", func(t *testing.T) {
		center := &model.RechargeCenter{
			ID:      "center_update",
			Name:    "Old Name",
			Code:    "UC_001",
			Address: "旧地址",
			Phone:   "010_11111111",
			Status:  "active",
		}
		err := repo.CreateCenter(center)
		require.NoError(t, err)

		err = repo.UpdateCenter(center.ID, map[string]interface{}{"name": "New Name", "address": "新地址"})
		require.NoError(t, err)

		// Verify update via raw query since GetCenters only returns active
		var updated model.RechargeCenter
		err = db.Where("id = ?", "center_update").First(&updated).Error
		require.NoError(t, err)
		assert.Equal(t, "New Name", updated.Name)
		assert.Equal(t, "新地址", updated.Address)
	})

	t.Run("delete center", func(t *testing.T) {
		center := &model.RechargeCenter{
			ID:     "1",
			Name:   "Delete Center",
			Code:   "DC_001",
			Status: "active",
		}
		err := repo.CreateCenter(center)
		require.NoError(t, err)

		err = repo.DeleteCenter("1")
		require.NoError(t, err)

		var found model.RechargeCenter
		err = db.Where("id = ?", "1").First(&found).Error
		assert.Error(t, err, "expected record not found after delete")
	})
}

// ========== 操作员 ==========

func TestRechargeRepository_Operator(t *testing.T) {
	db := setupRechargeTestDB(t)
	repo := NewRechargeRepository(db)

	t.Run("create and get operators", func(t *testing.T) {
		op := &model.RechargeOperator{
			ID:       "op_001",
			Name:     "操作员A",
			Phone:    "13800000001",
			Password: "hashed_password",
			CenterID: "center_1",
			Role:     "operator",
			Status:   "active",
		}

		err := repo.CreateOperator(op)
		require.NoError(t, err)

		list, err := repo.GetOperators()
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(list), 1)

		found := false
		for _, o := range list {
			if o.ID == "op_001" {
				found = true
				assert.Equal(t, "操作员A", o.Name)
				assert.Equal(t, "13800000001", o.Phone)
				assert.Equal(t, "center_1", o.CenterID)
				assert.Equal(t, "operator", o.Role)
				assert.Equal(t, "active", o.Status)
				break
			}
		}
		assert.True(t, found, "operator op_001 should be found in list")
	})

	t.Run("update operator", func(t *testing.T) {
		op := &model.RechargeOperator{
			ID:       "op_002",
			Name:     "旧名字",
			Phone:    "13800000002",
			Password: "hashed_password",
			CenterID: "center_1",
			Role:     "operator",
			Status:   "active",
		}
		err := repo.CreateOperator(op)
		require.NoError(t, err)

		op.Name = "新名字"
		op.Role = "admin"
		err = repo.UpdateOperator(op)
		require.NoError(t, err)

		var updated model.RechargeOperator
		err = db.Where("id = ?", "op_002").First(&updated).Error
		require.NoError(t, err)
		assert.Equal(t, "新名字", updated.Name)
		assert.Equal(t, "admin", updated.Role)
	})

	t.Run("delete operator", func(t *testing.T) {
		op := &model.RechargeOperator{
			ID:       "3",
			Name:     "待删除",
			Phone:    "13800000003",
			Password: "hashed_password",
			CenterID: "center_1",
			Role:     "operator",
			Status:   "active",
		}
		err := repo.CreateOperator(op)
		require.NoError(t, err)

		err = repo.DeleteOperator("3")
		require.NoError(t, err)

		var found model.RechargeOperator
		err = db.Where("id = ?", "3").First(&found).Error
		assert.Error(t, err, "expected record not found after delete")
	})
}
