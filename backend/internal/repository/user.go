package repository

import (
	"context"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"marketplace/backend/internal/config"
	"marketplace/backend/internal/model"
)

// NewDB 创建数据库连接
func NewDB(cfg *config.DatabaseConfig) (*sqlx.DB, error) {
	db, err := sqlx.Connect("mysql", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("connect database failed: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(time.Hour)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping database failed: %w", err)
	}

	return db, nil
}

// UserRepoInterface 用户仓库接口
type UserRepoInterface interface {
	GetByID(ctx context.Context, id int64) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	Create(ctx context.Context, user *model.User) (int64, error)
	Update(ctx context.Context, user *model.User) error
	List(ctx context.Context, page, pageSize int) ([]*model.User, int64, error)
}

// UserRepository 用户数据访问层
type UserRepository struct {
	db *sqlx.DB
}

var _ UserRepoInterface = (*UserRepository)(nil)

// NewUserRepository 创建用户 Repository
func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetByID 根据 ID 获取用户
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	user := &model.User{}
	query := "SELECT id, username, password, nickname, email, avatar, role, status, created_at, updated_at FROM users WHERE id = ?"
	if err := r.db.GetContext(ctx, user, query, id); err != nil {
		return nil, err
	}
	return user, nil
}

// GetByUsername 根据用户名获取用户
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	user := &model.User{}
	query := "SELECT id, username, password, nickname, email, avatar, role, status, created_at, updated_at FROM users WHERE username = ?"
	if err := r.db.GetContext(ctx, user, query, username); err != nil {
		return nil, err
	}
	return user, nil
}

// Create 创建用户
func (r *UserRepository) Create(ctx context.Context, user *model.User) (int64, error) {
	query := "INSERT INTO users (username, password, email, nickname, role, status) VALUES (?, ?, ?, ?, ?, ?)"
	result, err := r.db.ExecContext(ctx, query, user.Username, user.Password, user.Email, user.Nickname, user.Role, user.Status)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// Update 更新用户
func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	query := "UPDATE users SET nickname = ?, email = ?, avatar = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, user.Nickname, user.Email, user.Avatar, user.ID)
	return err
}

// UpdatePassword 更新密码
func (r *UserRepository) UpdatePassword(ctx context.Context, id int64, password string) error {
	query := "UPDATE users SET password = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, password, id)
	return err
}

// List 分页获取用户列表
func (r *UserRepository) List(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	countQuery := "SELECT COUNT(*) FROM users"
	if err := r.db.GetContext(ctx, &total, countQuery); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	query := "SELECT id, username, password, nickname, email, avatar, role, status, created_at, updated_at FROM users LIMIT ? OFFSET ?"
	if err := r.db.SelectContext(ctx, &users, query, pageSize, offset); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
