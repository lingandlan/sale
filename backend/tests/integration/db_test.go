package integration

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type TestDatabase struct {
	DB  *sqlx.DB
	DSN string
}

func SetupTestDB(t *testing.T) *TestDatabase {
	t.Helper()

	dsn := os.Getenv("TEST_DSN")
	if dsn == "" {
		dsn = "root:root123@tcp(localhost:3306)/test_db?parseTime=true&charset=utf8mb4"
	}

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		t.Skipf("跳过集成测试: 无法连接数据库: %v", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	return &TestDatabase{
		DB:  db,
		DSN: dsn,
	}
}

func (td *TestDatabase) Close() {
	if td.DB != nil {
		td.DB.Close()
	}
}

func (td *TestDatabase) CleanupTables(t *testing.T, tables ...string) {
	t.Helper()

	for _, table := range tables {
		_, err := td.DB.Exec(fmt.Sprintf("DELETE FROM %s", table))
		if err != nil {
			t.Logf("清理表 %s 失败: %v", table, err)
		}
	}
}

func (td *TestDatabase) ExecSQLFile(t *testing.T, filePath string) {
	t.Helper()

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("读取 SQL 文件失败: %v", err)
	}

	statements := strings.Split(string(content), ";")
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" || strings.HasPrefix(stmt, "--") {
			continue
		}
		_, err := td.DB.Exec(stmt)
		if err != nil && !strings.Contains(err.Error(), "Duplicate") {
			t.Logf("执行 SQL 失败: %v\nSQL: %s", err, stmt)
		}
	}
}

func (td *TestDatabase) CreateTestUser(t *testing.T, username, email string) int64 {
	t.Helper()

	result, err := td.DB.Exec(`
		INSERT INTO users (username, password, email, nickname, role, status) 
		VALUES (?, ?, ?, ?, 0, 1)
	`, username, "$2a$10$hashedpassword", email, username)

	if err != nil {
		t.Fatalf("创建测试用户失败: %v", err)
	}

	id, _ := result.LastInsertId()
	return id
}

func (td *TestDatabase) GetUserByUsername(t *testing.T, username string) *TestUser {
	t.Helper()

	var user TestUser
	err := td.DB.GetContext(context.Background(), &user,
		"SELECT id, username, email FROM users WHERE username = ?", username)

	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		t.Fatalf("查询用户失败: %v", err)
	}

	return &user
}

type TestUser struct {
	ID       int64  `db:"id"`
	Username string `db:"username"`
	Email    string `db:"email"`
}

func GetMigrationsDir() string {
	execPath, _ := os.Executable()
	return filepath.Join(filepath.Dir(execPath), "..", "..", "migrations")
}
