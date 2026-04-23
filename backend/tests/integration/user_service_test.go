//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"marketplace/backend/internal/model"
	"marketplace/backend/internal/repository"
	"marketplace/backend/internal/service"
)

func ensureUsersTable(t *testing.T, td *TestDatabase) {
	t.Helper()

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: td.DB.DB,
	}), &gorm.Config{})
	require.NoError(t, err)

	err = gormDB.AutoMigrate(&model.User{})
	require.NoError(t, err)
}

func TestUserService_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试: 使用 -short=false 运行")
	}

	td := SetupTestDB(t)
	defer td.Close()

	ensureUsersTable(t, td)

	userRepo := repository.NewUserRepository(td.DB)
	userSvc := service.NewUserService(userRepo)

	ctx := context.Background()
	testPhone := "13800138001"

	t.Run("Setup: 清理测试数据", func(t *testing.T) {
		td.CleanupTables(t, "users")
	})

	t.Run("Create", func(t *testing.T) {
		req := &model.CreateUserRequest{
			Phone:    testPhone,
			Password: "TestPassword123",
			Name:     "集成测试用户",
			Role:     model.RoleOperator,
		}

		user, err := userSvc.Create(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, testPhone, user.Phone)
		assert.Equal(t, "集成测试用户", user.Name)
		assert.Equal(t, model.RoleOperator, user.Role)
	})

	t.Run("GetByPhone", func(t *testing.T) {
		user, err := userSvc.GetByPhone(ctx, testPhone)
		require.NoError(t, err)
		assert.Equal(t, testPhone, user.Phone)
		assert.Equal(t, "集成测试用户", user.Name)
	})

	t.Run("GetByID", func(t *testing.T) {
		createdUser, err := userSvc.GetByPhone(ctx, testPhone)
		require.NoError(t, err)
		user, err := userSvc.GetByID(ctx, createdUser.ID)
		require.NoError(t, err)
		assert.Equal(t, createdUser.ID, user.ID)
		assert.Equal(t, testPhone, user.Phone)
	})

	t.Run("Update", func(t *testing.T) {
		createdUser, err := userSvc.GetByPhone(ctx, testPhone)
		require.NoError(t, err)
		newName := "更新后的名字"
		req := &model.UpdateUserRequest{
			Name: &newName,
		}

		user, err := userSvc.Update(ctx, createdUser.ID, req)
		require.NoError(t, err)
		assert.Equal(t, "更新后的名字", user.Name)
	})

	t.Run("List", func(t *testing.T) {
		users, total, err := userSvc.List(ctx, 1, 10)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(users), 1)
		assert.GreaterOrEqual(t, total, int64(1))
	})

	t.Run("Teardown: 清理测试用户", func(t *testing.T) {
		td.CleanupTables(t, "users")
	})
}

func TestAuthService_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	td := SetupTestDB(t)
	defer td.Close()

	ensureUsersTable(t, td)

	userRepo := repository.NewUserRepository(td.DB)
	userSvc := service.NewUserService(userRepo)

	ctx := context.Background()
	testPhone := "13900139001"

	t.Run("Setup", func(t *testing.T) {
		td.CleanupTables(t, "users")
		_, err := userSvc.Create(ctx, &model.CreateUserRequest{
			Phone:    testPhone,
			Password: "TestPassword123",
			Name:     "Auth测试用户",
			Role:     model.RoleOperator,
		})
		require.NoError(t, err)
	})

	t.Run("User exists with correct phone", func(t *testing.T) {
		user, err := userSvc.GetByPhone(ctx, testPhone)
		require.NoError(t, err)
		assert.Equal(t, testPhone, user.Phone)
	})

	t.Run("Non-existent user returns error", func(t *testing.T) {
		_, err := userSvc.GetByPhone(ctx, "19900000000")
		assert.Error(t, err)
	})

	t.Run("Teardown", func(t *testing.T) {
		td.CleanupTables(t, "users")
	})
}
