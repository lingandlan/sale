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
	// 数据库连接
	dsn := "sale:sale123@tcp(localhost:3306)/sale_dev?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}

	// 查询用户
	ctx := context.Background()
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
		log.Fatal("查询用户失败:", err)
	}

	fmt.Printf("用户信息: ID=%d, Phone=%s, Name=%s, Role=%s, Status=%d\n",
		user.ID, user.Phone, user.Name, user.Role, user.Status)

	// 测试密码
	password := "admin123"
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Println("密码验证失败:", err)
	} else {
		fmt.Println("密码验证成功！")
	}
}
