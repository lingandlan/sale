package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// config 从环境变量读取，与 docker-compose.prod.yml 中的 APP_DATABASE_* 保持一致
var (
	dbHost     = envOrDefault("DATABASE_HOST", "localhost")
	dbPort     = envOrDefault("DATABASE_PORT", "3306")
	dbUser     = envOrDefault("DATABASE_USER", "sale")
	dbPassword = envOrDefault("DATABASE_PASSWORD", "sale123")
	dbName     = envOrDefault("DATABASE_NAME", "sale_dev")
	sqlDir     = "migrations/sql"
)

func main() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local&multiStatements=true",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("连接数据库失败: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Printf("数据库连接测试失败: %v\n", err)
		os.Exit(1)
	}

	// 创建 schema_migrations 表
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		version VARCHAR(255) PRIMARY KEY,
		executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`); err != nil {
		fmt.Printf("创建 schema_migrations 表失败: %v\n", err)
		os.Exit(1)
	}

	// 读取 SQL 文件
	files, err := filepath.Glob(filepath.Join(sqlDir, "*.sql"))
	if err != nil {
		fmt.Printf("读取 migration 目录失败: %v\n", err)
		os.Exit(1)
	}
	sort.Strings(files)

	// 查询已执行的
	executed := make(map[string]bool)
	rows, err := db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		fmt.Printf("查询已执行记录失败: %v\n", err)
		os.Exit(1)
	}
	for rows.Next() {
		var v string
		rows.Scan(&v)
		executed[v] = true
	}
	rows.Close()

	// 执行未执行的
	executedCount := 0
	for _, f := range files {
		name := filepath.Base(f)
		if executed[name] {
			fmt.Printf("  SKIP  %s (已执行)\n", name)
			continue
		}

		content, err := os.ReadFile(f)
		if err != nil {
			fmt.Printf("  ERROR 读取 %s 失败: %v\n", name, err)
			continue
		}

		fmt.Printf("  EXEC  %s\n", name)
		if _, err := db.Exec(string(content)); err != nil {
			fmt.Printf("  FAIL  %s: %v\n", name, err)
			// 记录失败但不退出，继续下一个
			continue
		}

		db.Exec("INSERT INTO schema_migrations (version) VALUES (?)", name)
		executedCount++
	}

	if executedCount == 0 {
		fmt.Println("\n所有 migration 已是最新，无需执行。")
	} else {
		fmt.Printf("\n执行完成，共 %d 个 migration。\n", executedCount)
	}
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		// 支持 APP_ 前缀的环境变量
		return v
	}
	if v := os.Getenv("APP_" + strings.ToUpper(key)); v != "" {
		return v
	}
	return fallback
}
