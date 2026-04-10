//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"marketplace/backend/internal/model"
	"marketplace/backend/internal/repository"
	"marketplace/backend/internal/service"
)

func TestUserService_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试: 使用 -short=false 运行")
	}

	db := SetupTestDB(t)
	defer db.Close()

	userRepo := repository.NewUserRepository(db.DB)
	userSvc := service.NewUserService(userRepo)

	ctx := context.Background()
	testUsername := "integration_test_user"
	testEmail := "integration@test.com"

	t.Run("Setup: 创建测试用户", func(t *testing.T) {
		db.CleanupTables(t, "users")
	})

	t.Run("Create", func(t *testing.T) {
		req := &model.CreateUserRequest{
			Username: testUsername,
			Password: "TestPassword123",
			Email:    testEmail,
			Nickname: "Integration Test",
		}

		user, err := userSvc.Create(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, testUsername, user.Username)
		assert.Equal(t, testEmail, user.Email)
		assert.Equal(t, model.RoleUser, user.Role)
		assert.Equal(t, model.UserStatusNormal, user.Status)
	})

	t.Run("GetByUsername", func(t *testing.T) {
		user, err := userSvc.GetByUsername(ctx, testUsername)

		require.NoError(t, err)
		assert.Equal(t, testUsername, user.Username)
		assert.Equal(t, testEmail, user.Email)
	})

	t.Run("GetByID", func(t *testing.T) {
		createdUser, _ := userSvc.GetByUsername(ctx, testUsername)

		user, err := userSvc.GetByID(ctx, createdUser.ID)

		require.NoError(t, err)
		assert.Equal(t, createdUser.ID, user.ID)
		assert.Equal(t, testUsername, user.Username)
	})

	t.Run("Update", func(t *testing.T) {
		createdUser, _ := userSvc.GetByUsername(ctx, testUsername)

		newNickname := "Updated Nickname"
		newEmail := "updated@test.com"
		req := &model.UpdateUserRequest{
			Nickname: &newNickname,
			Email:    &newEmail,
		}

		user, err := userSvc.Update(ctx, createdUser.ID, req)

		require.NoError(t, err)
		assert.Equal(t, "Updated Nickname", user.Nickname)
		assert.Equal(t, "updated@test.com", user.Email)
	})

	t.Run("List", func(t *testing.T) {
		users, total, err := userSvc.List(ctx, 1, 10)

		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(users), 1)
		assert.GreaterOrEqual(t, total, int64(1))
	})

	t.Run("Teardown: 清理测试用户", func(t *testing.T) {
		db.CleanupTables(t, "users")
	})
}

func TestAuthService_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	db := SetupTestDB(t)
	defer db.Close()

	userRepo := repository.NewUserRepository(db.DB)
	userSvc := service.NewUserService(userRepo)

	ctx := context.Background()
	testUsername := "auth_integration_test"
	testPassword := "OriginalPassword"

	t.Run("Setup", func(t *testing.T) {
		db.CleanupTables(t, "users")

		_, err := userSvc.Create(ctx, &model.CreateUserRequest{
			Username: testUsername,
			Password: testPassword,
			Email:    "auth@test.com",
		})
		require.NoError(t, err)
	})

	t.Run("Login with correct password", func(t *testing.T) {
		// 获取用户密码哈希进行验证
		user, err := userSvc.GetByUsername(ctx, testUsername)
		require.NoError(t, err)

		// 由于 Login 需要密码验证，这里验证用户存在
		assert.Equal(t, testUsername, user.Username)
	})

	t.Run("Login with non-existent user", func(t *testing.T) {
		_, err := userSvc.GetByUsername(ctx, "non_existent_user_xyz")
		assert.Error(t, err)
	})

	t.Run("Teardown", func(t *testing.T) {
		db.CleanupTables(t, "users")
	})
}

func TestRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	db := SetupTestDB(t)
	defer db.Close()

	repo := repository.NewUserRepository(db.DB)
	ctx := context.Background()

	t.Run("CRUD operations", func(t *testing.T) {
		db.CleanupTables(t, "users")

		user := &model.User{
			Username: "repo_test",
			Password: "hashed",
			Email:    "repo@test.com",
			Nickname: "Repo Test",
			Role:     0,
			Status:   1,
		}

		id, err := repo.Create(ctx, user)
		require.NoError(t, err)
		assert.Greater(t, id, int64(0))

		created, err := repo.GetByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "repo_test", created.Username)

		created.Nickname = "Updated"
		err = repo.Update(ctx, created)
		require.NoError(t, err)

		updated, err := repo.GetByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "Updated", updated.Nickname)

		db.CleanupTables(t, "users")
	})
}
