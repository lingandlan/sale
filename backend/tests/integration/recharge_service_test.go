//go:build integration
// +build integration

package integration

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"marketplace/backend/internal/model"
	"marketplace/backend/internal/repository"
	"marketplace/backend/internal/service"
)

func setupRechargeTestDB(t *testing.T) (*gorm.DB, *TestDatabase) {
	t.Helper()

	td := SetupTestDB(t)

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: td.DB.DB,
	}), &gorm.Config{})
	require.NoError(t, err)

	err = gormDB.AutoMigrate(
		&model.RechargeApplication{},
		&model.CRecharge{},
		&model.StoreCard{},
		&model.CardTransaction{},
		&model.RechargeCenter{},
		&model.RechargeOperator{},
	)
	require.NoError(t, err)

	return gormDB, td
}

func TestRechargeService_BRechargeFlow_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	gormDB, td := setupRechargeTestDB(t)
	defer td.Close()

	repo := repository.NewRechargeRepository(gormDB)
	svc := service.NewRechargeService(repo)

	td.CleanupTables(t, "recharge_applications")

	t.Run("Create application", func(t *testing.T) {
		data := map[string]interface{}{
			"centerId":            "center-1",
			"centerName":          "北京中心",
			"amount":              float64(10000),
			"lastMonthConsumption": float64(120000),
			"applicantId":         "user-1",
			"applicantName":       "张三",
			"transactionNo":       "TX-INTEG-001",
			"screenshot":           "",
			"remark":              "集成测试",
		}

		app, err := svc.CreateBRechargeApplication(data)
		require.NoError(t, err)
		assert.NotEmpty(t, app.ID)
		assert.Equal(t, "pending", app.Status)
		assert.Equal(t, 10000, app.BasePoints)
		assert.Equal(t, 200, app.RebatePoints)
	})

	t.Run("List applications", func(t *testing.T) {
		result, err := svc.GetRechargeApplicationList("pending", 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(1), result["total"])
	})

	t.Run("Approve application", func(t *testing.T) {
		// Get the pending application
		result, err := svc.GetRechargeApplicationList("pending", 1, 10)
		require.NoError(t, err)
		list := result["list"].([]model.RechargeApplication)
		require.Len(t, list, 1)

		err = svc.ApproveRechargeApplication(list[0].ID, "approve", "admin-test", "通过")
		require.NoError(t, err)

		// Verify status changed
		app, err := svc.GetRechargeApplicationDetail(list[0].ID)
		require.NoError(t, err)
		assert.Equal(t, "approved", app.Status)
	})
}

func TestRechargeService_CardFlow_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	gormDB, td := setupRechargeTestDB(t)
	defer td.Close()

	repo := repository.NewRechargeRepository(gormDB)
	svc := service.NewRechargeService(repo)

	td.CleanupTables(t, "store_cards", "card_transactions")

	t.Run("Issue card", func(t *testing.T) {
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
		assert.Equal(t, 5000.0, card.Balance)
	})

	t.Run("Verify card", func(t *testing.T) {
		// Get the issued card
		result, err := svc.GetCardList("active", "", "", 1, 10)
		require.NoError(t, err)
		list := result["list"].([]model.StoreCard)
		require.Len(t, list, 1)

		card, err := svc.VerifyCard(list[0].CardNo)
		require.NoError(t, err)
		assert.Equal(t, "active", card.Status)
		assert.WithinDuration(t, time.Now().AddDate(1, 0, 0), card.ExpiryDate, 5*time.Second)
	})

	t.Run("Consume card", func(t *testing.T) {
		// Get the card
		result, err := svc.GetCardList("active", "", "", 1, 10)
		require.NoError(t, err)
		list := result["list"].([]model.StoreCard)
		require.Len(t, list, 1)
		cardNo := list[0].CardNo

		err = svc.ConsumeCard(cardNo, 2000, "消费", "op-1")
		require.NoError(t, err)

		// Verify balance updated
		card, err := svc.VerifyCard(cardNo)
		require.NoError(t, err)
		assert.Equal(t, 3000.0, card.Balance)
	})

	t.Run("Freeze card", func(t *testing.T) {
		result, err := svc.GetCardList("active", "", "", 1, 10)
		require.NoError(t, err)
		list := result["list"].([]model.StoreCard)
		require.Len(t, list, 1)

		err = svc.UpdateCardStatus(list[0].CardNo, "inactive")
		require.NoError(t, err)

		// Verify frozen
		_, err = svc.VerifyCard(list[0].CardNo)
		assert.Error(t, err) // inactive card should fail verify
	})
}

func TestRechargeService_CenterCRUD_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	gormDB, td := setupRechargeTestDB(t)
	defer td.Close()

	repo := repository.NewRechargeRepository(gormDB)
	svc := service.NewRechargeService(repo)

	td.CleanupTables(t, "recharge_centers")

	t.Run("Create center", func(t *testing.T) {
		data := map[string]interface{}{
			"name":    "集成测试中心",
			"code":    "TEST001",
			"address": "北京市海淀区",
			"phone":   "010-12345678",
		}

		center, err := svc.CreateCenter(data)
		require.NoError(t, err)
		assert.NotEmpty(t, center.ID)
		assert.Equal(t, "active", center.Status)
	})

	t.Run("List centers", func(t *testing.T) {
		centers, err := svc.GetCenters()
		require.NoError(t, err)
		assert.Len(t, centers, 1)
		assert.Equal(t, "集成测试中心", centers[0].Name)
	})

	t.Run("Update center", func(t *testing.T) {
		centers, _ := svc.GetCenters()
		data := map[string]interface{}{
			"name": "更新中心名称",
			"code":    "TEST001",
			"address": "北京市朝阳区",
			"phone":   "010-87654321",
			"status":  "active",
		}

		center, err := svc.UpdateCenter(centers[0].ID, data)
		require.NoError(t, err)
		assert.Equal(t, "更新中心名称", center.Name)
	})

	t.Run("Delete center", func(t *testing.T) {
		centers, _ := svc.GetCenters()
		err := svc.DeleteCenter(centers[0].ID)
		require.NoError(t, err)

		// Verify empty
		centers, err = svc.GetCenters()
		require.NoError(t, err)
		assert.Len(t, centers, 0)
	})
}
