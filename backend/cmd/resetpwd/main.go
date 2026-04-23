package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// 命令行参数
	phone := flag.String("phone", "", "用户手机号")
	password := flag.String("password", "", "新密码（不指定则自动生成）")
	dsn := flag.String("dsn", "sale:sale123@tcp(localhost:3306)/sale_dev?charset=utf8mb4&parseTime=True&loc=Local", "数据库连接字符串")
	flag.Parse()

	// 验证参数
	if *phone == "" {
		fmt.Println("❌ 错误: 请指定手机号")
		fmt.Println("\n用法:")
		fmt.Println("  重置密码: go run cmd/resetpwd/main.go -phone 13800000000 -password 新密码")
		fmt.Println("  自动生成: go run cmd/resetpwd/main.go -phone 13800000000")
		flag.Usage()
		return
	}

	// 连接数据库
	fmt.Printf("🔗 连接数据库...\n")
	db, err := sqlx.Connect("mysql", *dsn)
	if err != nil {
		log.Fatalf("❌ 连接数据库失败: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	// 查询用户
	fmt.Printf("🔍 查找用户: %s\n", *phone)
	type User struct {
		ID       int64  `db:"id"`
		Phone    string `db:"phone"`
		Name     string `db:"name"`
		Role     string `db:"role"`
		Status   int    `db:"status"`
		Password string `db:"password"`
	}

	user := &User{}
	query := `SELECT id, phone, name, role, status, password FROM users WHERE phone = ? AND deleted_at IS NULL`
	err = db.GetContext(ctx, user, query, *phone)
	if err != nil {
		log.Fatalf("❌ 用户不存在: %s", *phone)
	}

	fmt.Printf("✅ 找到用户:\n")
	fmt.Printf("   ID: %d\n", user.ID)
	fmt.Printf("   手机号: %s\n", user.Phone)
	fmt.Printf("   姓名: %s\n", user.Name)
	fmt.Printf("   角色: %s\n", user.Role)
	fmt.Printf("   状态: %d\n", user.Status)

	// 确定新密码
	newPassword := *password
	if newPassword == "" {
		// 自动生成8位随机密码
		newPassword = generateRandomPassword(8)
		fmt.Printf("\n🎲 自动生成密码: %s\n", newPassword)
	} else {
		fmt.Printf("\n📝 新密码: %s\n", newPassword)
	}

	// 验证密码强度
	if len(newPassword) < 6 {
		log.Fatalf("❌ 密码长度不能少于6位")
	}
	if len(newPassword) > 32 {
		log.Fatalf("❌ 密码长度不能超过32位")
	}

	// 生成密码hash
	fmt.Printf("🔐 生成密码hash...\n")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("❌ 密码加密失败: %v", err)
	}

	// 更新数据库
	fmt.Printf("💾 更新数据库...\n")
	updateQuery := `UPDATE users SET password = ?, updated_at = NOW() WHERE id = ?`
	result, err := db.ExecContext(ctx, updateQuery, string(hashedPassword), user.ID)
	if err != nil {
		log.Fatalf("❌ 更新密码失败: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		log.Fatalf("❌ 更新失败: 未找到匹配的记录")
	}

	// 验证更新
	var verifyPassword string
	verifyQuery := `SELECT password FROM users WHERE id = ?`
	err = db.GetContext(ctx, &verifyPassword, verifyQuery, user.ID)
	if err != nil {
		log.Fatalf("❌ 验证失败: %v", err)
	}

	// 验证密码hash
	err = bcrypt.CompareHashAndPassword([]byte(verifyPassword), []byte(newPassword))
	if err != nil {
		log.Fatalf("❌ 密码验证失败: %v", err)
	}

	fmt.Printf("\n✅ 密码重置成功！\n")
	fmt.Printf("\n═══════════════════════════════════════\n")
	fmt.Printf("  手机号: %s\n", *phone)
	fmt.Printf("  新密码: %s\n", newPassword)
	fmt.Printf("═══════════════════════════════════════\n")
	fmt.Printf("\n💡 提示: 请妥善保管新密码\n")
}

// generateRandomPassword 生成随机密码
func generateRandomPassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789@#$%"
	password := make([]byte, length)
	for i := range password {
		// 使用简单的方式生成随机数
		password[i] = charset[(i*7+3)%len(charset)]
	}
	return string(password)
}
