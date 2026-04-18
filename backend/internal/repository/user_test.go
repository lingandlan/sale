package repository

import (
	"context"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"marketplace/backend/internal/model"
)

func setupTestDB(t *testing.T) (*sqlx.DB, func()) {
	gormDB, err := gorm.Open(sqlite.Open("file::memory:?cache=shared&_foreign_keys=ON"), &gorm.Config{})
	require.NoError(t, err)

	err = gormDB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL DEFAULT '',
			phone TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			name TEXT NOT NULL,
			role TEXT DEFAULT 'operator',
			center_id TEXT DEFAULT NULL,
			center_name TEXT DEFAULT NULL,
			status INTEGER DEFAULT 1,
			last_login_at DATETIME DEFAULT NULL,
			last_login_ip TEXT DEFAULT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			deleted_at DATETIME DEFAULT NULL
		)
	`).Error
	require.NoError(t, err)

	sqlDB, err := gormDB.DB()
	require.NoError(t, err)

	cleanup := func() {
		sqlDB.Close()
	}

	return sqlx.NewDb(sqlDB, "sqlite3"), cleanup
}

func TestUserRepository_CRUD(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewUserRepository(db)
	ctx := context.Background()

	t.Run("create and get by id", func(t *testing.T) {
		user := &model.User{
			Username: "testuser1",
			Phone:    "13800138000",
			Password: "password",
			Name:     "CRUD User",
			Role:     model.RoleOperator,
			Status:   model.UserStatusNormal,
		}

		id, err := repo.Create(ctx, user)
		require.NoError(t, err)
		assert.Greater(t, id, int64(0))

		result, err := repo.GetByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "testuser1", result.Username)
		assert.Equal(t, "13800138000", result.Phone)
		assert.Equal(t, "CRUD User", result.Name)
	})

	t.Run("get by phone", func(t *testing.T) {
		user := &model.User{
			Username: "testuser2",
			Phone:    "13800138001",
			Password: "password",
			Name:     "GetByPhone User",
			Role:     model.RoleOperator,
			Status:   model.UserStatusNormal,
		}
		repo.Create(ctx, user)

		result, err := repo.GetByPhone(ctx, "13800138001")
		require.NoError(t, err)
		assert.Equal(t, "GetByPhone User", result.Name)
	})

	t.Run("get non-existent user", func(t *testing.T) {
		_, err := repo.GetByPhone(ctx, "13900000000")
		assert.Error(t, err)
	})

	t.Run("update user", func(t *testing.T) {
		user := &model.User{
			Username: "testuser3",
			Phone:    "13800138002",
			Password: "password",
			Name:     "Old Name",
			Role:     model.RoleOperator,
			Status:   model.UserStatusNormal,
		}
		id, _ := repo.Create(ctx, user)

		updatedUser := &model.User{
			ID:       id,
			Username: "testuser3_new",
			Phone:    "13800138002",
			Password: "password",
			Name:     "New Name",
			Role:     model.RoleOperator,
			Status:   model.UserStatusNormal,
		}
		err := repo.Update(ctx, updatedUser)
		require.NoError(t, err)

		result, _ := repo.GetByID(ctx, id)
		assert.Equal(t, "New Name", result.Name)
		assert.Equal(t, "testuser3_new", result.Username)
	})

	t.Run("list users", func(t *testing.T) {
		users := []*model.User{
			{Username: "list1", Phone: "13800138003", Password: "p", Name: "List User 1", Role: model.RoleOperator, Status: model.UserStatusNormal},
			{Username: "list2", Phone: "13800138004", Password: "p", Name: "List User 2", Role: model.RoleHQAdmin, Status: model.UserStatusNormal},
			{Username: "list3", Phone: "13800138005", Password: "p", Name: "List User 3", Role: model.RoleOperator, Status: model.UserStatusNormal},
		}
		for _, u := range users {
			repo.Create(ctx, u)
		}

		result, total, err := repo.List(ctx, 1, 10)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(result), 3)
		assert.GreaterOrEqual(t, total, int64(3))
	})

	t.Run("update password", func(t *testing.T) {
		user := &model.User{
			Username: "pwduser",
			Phone:    "13800138006",
			Password: "oldpass",
			Name:     "Pwd User",
			Role:     model.RoleOperator,
			Status:   model.UserStatusNormal,
		}
		id, _ := repo.Create(ctx, user)

		err := repo.UpdatePassword(ctx, id, "newpass")
		require.NoError(t, err)
	})

	t.Run("delete user", func(t *testing.T) {
		user := &model.User{
			Username: "deluser",
			Phone:    "13800138007",
			Password: "password",
			Name:     "Delete User",
			Role:     model.RoleOperator,
			Status:   model.UserStatusNormal,
		}
		id, _ := repo.Create(ctx, user)

		err := repo.Delete(ctx, id)
		require.NoError(t, err)
	})

	t.Run("update status", func(t *testing.T) {
		user := &model.User{
			Username: "statususer",
			Phone:    "13800138008",
			Password: "password",
			Name:     "Status User",
			Role:     model.RoleOperator,
			Status:   model.UserStatusNormal,
		}
		id, _ := repo.Create(ctx, user)

		err := repo.UpdateStatus(ctx, id, model.UserStatusDisabled)
		require.NoError(t, err)

		result, _ := repo.GetByID(ctx, id)
		assert.Equal(t, int8(model.UserStatusDisabled), result.Status)
	})
}
