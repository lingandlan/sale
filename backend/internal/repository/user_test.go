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
			username TEXT NOT NULL,
			password TEXT NOT NULL,
			nickname TEXT DEFAULT '',
			email TEXT DEFAULT '',
			avatar TEXT DEFAULT NULL,
			role INTEGER DEFAULT 0,
			status INTEGER DEFAULT 1,
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
			Username: "cruduser",
			Password: "password",
			Email:    "crud@test.com",
			Nickname: "CRUD User",
			Role:     0,
			Status:   1,
		}

		id, err := repo.Create(ctx, user)
		require.NoError(t, err)
		assert.Greater(t, id, int64(0))

		result, err := repo.GetByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "cruduser", result.Username)
		assert.Equal(t, "crud@test.com", result.Email)
	})

	t.Run("get by username", func(t *testing.T) {
		user := &model.User{
			Username: "getbyuser",
			Password: "password",
			Email:    "getbyuser@test.com",
			Role:     0,
			Status:   1,
		}
		repo.Create(ctx, user)

		result, err := repo.GetByUsername(ctx, "getbyuser")
		require.NoError(t, err)
		assert.Equal(t, "getbyuser@test.com", result.Email)
	})

	t.Run("get non-existent user", func(t *testing.T) {
		_, err := repo.GetByUsername(ctx, "nonexistent")
		assert.Error(t, err)
	})

	t.Run("update user", func(t *testing.T) {
		user := &model.User{
			Username: "updateuser",
			Password: "password",
			Email:    "old@test.com",
			Nickname: "Old Nick",
			Role:     0,
			Status:   1,
		}
		id, _ := repo.Create(ctx, user)

		updatedUser := &model.User{
			ID:       id,
			Username: "updateuser",
			Password: "password",
			Email:    "new@test.com",
			Nickname: "New Nick",
			Role:     0,
			Status:   1,
		}
		err := repo.Update(ctx, updatedUser)
		require.NoError(t, err)

		result, _ := repo.GetByID(ctx, id)
		assert.Equal(t, "New Nick", result.Nickname)
		assert.Equal(t, "new@test.com", result.Email)
	})

	t.Run("list users", func(t *testing.T) {
		users := []*model.User{
			{Username: "list1", Password: "p", Email: "l1@t.com", Role: 0, Status: 1},
			{Username: "list2", Password: "p", Email: "l2@t.com", Role: 0, Status: 1},
			{Username: "list3", Password: "p", Email: "l3@t.com", Role: 0, Status: 1},
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
			Password: "oldpass",
			Email:    "pwd@t.com",
			Role:     0,
			Status:   1,
		}
		id, _ := repo.Create(ctx, user)

		err := repo.UpdatePassword(ctx, id, "newpass")
		require.NoError(t, err)
	})
}
