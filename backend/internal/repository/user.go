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
	GetByPhone(ctx context.Context, phone string) (*model.User, error)
	Create(ctx context.Context, user *model.User) (int64, error)
	Update(ctx context.Context, user *model.User) error
	UpdatePassword(ctx context.Context, id int64, password string) error
	List(ctx context.Context, page, pageSize int) ([]*model.User, int64, error)
	ListWithFilters(ctx context.Context, page, pageSize int, keyword, role string, status *int8) ([]*model.User, int64, error)
	UpdateStatus(ctx context.Context, id int64, status int8) error
	UpdateLoginInfo(ctx context.Context, id int64, loginIP string) error
	Delete(ctx context.Context, id int64) error
}

// UserRepository 用户数据访问层
type UserRepository struct {
	db *sqlx.DB
}

var _ UserRepoInterface = (*UserRepository)(nil)

// UserRepositoryInterface is an alias for UserRepoInterface, used by handler layer
type UserRepositoryInterface = UserRepoInterface

// NewUserRepository 创建用户 Repository
func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetByID 根据 ID 获取用户
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	user := &model.User{}
	query := `SELECT id, username, phone, password, name, role, center_id, center_name, status,
	             last_login_at, last_login_ip, created_at, updated_at
	         FROM users WHERE id = ? AND deleted_at IS NULL`
	if err := r.db.GetContext(ctx, user, query, id); err != nil {
		return nil, err
	}
	return user, nil
}

// GetByPhone 根据手机号获取用户
func (r *UserRepository) GetByPhone(ctx context.Context, phone string) (*model.User, error) {
	user := &model.User{}
	query := `SELECT id, username, phone, password, name, role, center_id, center_name, status,
	             last_login_at, last_login_ip, created_at, updated_at
	         FROM users WHERE phone = ? AND deleted_at IS NULL`
	if err := r.db.GetContext(ctx, user, query, phone); err != nil {
		return nil, err
	}
	return user, nil
}

// Create 创建用户
func (r *UserRepository) Create(ctx context.Context, user *model.User) (int64, error) {
	query := `INSERT INTO users (username, phone, password, name, role, center_id, center_name, status, created_at, updated_at)
	         VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	now := time.Now()
	result, err := r.db.ExecContext(ctx, query,
		user.Username, user.Phone, user.Password, user.Name, user.Role,
		user.CenterID, user.CenterName, user.Status, now, now)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// Update 更新用户
func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	query := `UPDATE users SET username = ?, name = ?, phone = ?, role = ?, center_id = ?, center_name = ?, updated_at = ?
	         WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, user.Username, user.Name, user.Phone, user.Role, user.CenterID, user.CenterName, time.Now(), user.ID)
	return err
}

// UpdatePassword 更新密码
func (r *UserRepository) UpdatePassword(ctx context.Context, id int64, password string) error {
	query := "UPDATE users SET password = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, password, id)
	return err
}

// UpdateLoginInfo 更新登录信息
func (r *UserRepository) UpdateLoginInfo(ctx context.Context, id int64, loginIP string) error {
	query := "UPDATE users SET last_login_at = CURRENT_TIMESTAMP, last_login_ip = ? WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, loginIP, id)
	return err
}

// List 分页获取用户列表
func (r *UserRepository) List(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	countQuery := "SELECT COUNT(*) FROM users WHERE deleted_at IS NULL"
	if err := r.db.GetContext(ctx, &total, countQuery); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	query := `SELECT id, username, phone, name, role, center_id, center_name, status,
	             last_login_at, last_login_ip, created_at, updated_at
	         FROM users WHERE deleted_at IS NULL
	         ORDER BY created_at DESC
	         LIMIT ? OFFSET ?`
	if err := r.db.SelectContext(ctx, &users, query, pageSize, offset); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// ListWithFilters 带筛选条件的用户列表（管理员）
func (r *UserRepository) ListWithFilters(ctx context.Context, page, pageSize int, keyword, role string, status *int8) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	// 构建WHERE条件
	where := "deleted_at IS NULL"
	args := []interface{}{}

	if keyword != "" {
		where += " AND (username LIKE ? OR phone LIKE ? OR name LIKE ?)"
		args = append(args, "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	if role != "" {
		where += " AND role = ?"
		args = append(args, role)
	}

	if status != nil {
		where += " AND status = ?"
		args = append(args, *status)
	}

	// 获取总数
	countQuery := "SELECT COUNT(*) FROM users WHERE " + where
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	query := `SELECT id, username, phone, name, role, center_id, center_name, status,
	             last_login_at, last_login_ip, created_at, updated_at
	         FROM users WHERE ` + where + `
	         ORDER BY created_at DESC
	         LIMIT ? OFFSET ?`
	args = append(args, pageSize, offset)

	if err := r.db.SelectContext(ctx, &users, query, args...); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// Delete 软删除用户
func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	query := "UPDATE users SET deleted_at = CURRENT_TIMESTAMP WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// UpdateStatus 更新用户状态
func (r *UserRepository) UpdateStatus(ctx context.Context, id int64, status int8) error {
	query := "UPDATE users SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}
