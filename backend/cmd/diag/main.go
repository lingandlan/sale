package main

import (
	"context"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	dsn := "sale:sale123@tcp(localhost:3306)/sale_dev?charset=utf8mb4&parseTime=True&loc=Local"

	// 连接数据库
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	// 查询用户
	type User struct {
		ID       int64  `db:"id"`
		Phone    string `db:"phone"`
		Password string `db:"password"`
		Name     string `db:"name"`
		Role     string `db:"role"`
		Status   int    `db:"status"`
	}

	user := &User{}
	query := `SELECT id, phone, password, name, role, status FROM users WHERE phone = ? AND deleted_at IS NULL`
	err = db.GetContext(ctx, user, query, "13800000000")
	if err != nil {
		log.Fatalf("查询用户失败: %v", err)
	}

	fmt.Printf("=== 用户信息 ===\n")
	fmt.Printf("ID: %d\n", user.ID)
	fmt.Printf("手机号: %s\n", user.Phone)
	fmt.Printf("姓名: %s\n", user.Name)
	fmt.Printf("角色: %s\n", user.Role)
	fmt.Printf("状态: %d\n", user.Status)
	fmt.Printf("密码Hash: %s\n", user.Password)

	// 测试密码验证
	testPasswords := []string{"admin123", "Test123456", "123456"}

	fmt.Printf("\n=== 密码验证测试 ===\n")
	for _, pwd := range testPasswords {
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pwd))
		if err == nil {
			fmt.Printf("✅ 密码 '%s' 验证成功\n", pwd)
		} else {
			fmt.Printf("❌ 密码 '%s' 验证失败: %v\n", pwd, err)
		}
	}

	// 生成新的密码hash
	fmt.Printf("\n=== 生成新密码hash ===\n")
	newPassword := "NewPass123"
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("生成hash失败: %v", err)
	}
	fmt.Printf("新密码: %s\n", newPassword)
	fmt.Printf("新Hash: %s\n", string(hash))

	// 验证新hash
	err = bcrypt.CompareHashAndPassword(hash, []byte(newPassword))
	if err == nil {
		fmt.Printf("✅ 新hash验证成功\n")
	} else {
		fmt.Printf("❌ 新hash验证失败: %v\n", err)
	}

	// 检查用户状态
	fmt.Printf("\n=== 用户状态检查 ===\n")
	const UserStatusNormal = 1
	if user.Status == UserStatusNormal {
		fmt.Printf("✅ 用户状态正常 (Status=%d)\n", user.Status)
	} else {
		fmt.Printf("❌ 用户状态异常 (Status=%d)\n", user.Status)
	}
}
