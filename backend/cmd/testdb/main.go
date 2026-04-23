package main

import (
	"context"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func main() {
	// 使用后端的配置
	dsn := "sale:sale123@tcp(localhost:3306)/sale_dev?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	// 测试查询
	type User struct {
		ID     int64  `db:"id"`
		Phone  string `db:"phone"`
		Status int    `db:"status"`
	}

	user := &User{}
	query := `SELECT id, phone, status FROM users WHERE phone = ? AND deleted_at IS NULL`
	err = db.GetContext(ctx, user, query, "13800000000")

	fmt.Printf("查询结果:\n")
	fmt.Printf("错误: %v\n", err)
	fmt.Printf("用户: %+v\n", user)

	if err == nil && user.ID == 1 {
		fmt.Printf("\n✅ 数据库连接正常，可以查询到用户\n")
	} else {
		fmt.Printf("\n❌ 数据库查询有问题\n")
	}
}
