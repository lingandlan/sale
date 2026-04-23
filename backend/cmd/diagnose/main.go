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
	fmt.Println("=== 登录问题诊断工具 ===")

	dsn := "sale:sale123@tcp(localhost:3306)/sale_dev?charset=utf8mb4&parseTime=True&loc=Local"

	// 1. 测试数据库连接
	fmt.Println("【步骤1】测试数据库连接")
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalf("❌ 数据库连接失败: %v\n", err)
	}
	fmt.Println("✅ 数据库连接成功")
	defer db.Close()

	ctx := context.Background()

	// 2. 查询用户（完全模仿后端代码）
	fmt.Println("【步骤2】查询用户（模仿Repository.GetByPhone）")
	type User struct {
		ID          int64   `db:"id"`
		Phone       string  `db:"phone"`
		Password    string  `db:"password"`
		Name        string  `db:"name"`
		Role        string  `db:"role"`
		CenterID    *int64  `db:"center_id"`
		CenterName  *string `db:"center_name"`
		Status      int     `db:"status"`
		LastLoginAt *string `db:"last_login_at"`
		LastLoginIP *string `db:"last_login_ip"`
		CreatedAt   string  `db:"created_at"`
		UpdatedAt   string  `db:"updated_at"`
	}

	user := &User{}
	query := `SELECT id, phone, password, name, role, center_id, center_name, status,
	             last_login_at, last_login_ip, created_at, updated_at
	         FROM users WHERE phone = ? AND deleted_at IS NULL`

	err = db.GetContext(ctx, user, query, "13800000000")
	if err != nil {
		log.Fatalf("❌ 查询用户失败: %v\n", err)
	}
	fmt.Printf("✅ 用户查询成功:\n")
	fmt.Printf("   ID: %d\n", user.ID)
	fmt.Printf("   Phone: %s\n", user.Phone)
	fmt.Printf("   Name: %s\n", user.Name)
	fmt.Printf("   Role: %s\n", user.Role)
	fmt.Printf("   Status: %d\n", user.Status)
	centerName := ""
	if user.CenterName != nil {
		centerName = *user.CenterName
	}
	fmt.Printf("   CenterName: %s\n", centerName)
	fmt.Printf("   Password: %s\n\n", user.Password)

	// 3. 检查用户状态
	fmt.Println("【步骤3】检查用户状态（模仿AuthService.Login）")
	const UserStatusNormal = 1
	if user.Status != UserStatusNormal {
		log.Fatalf("❌ 用户状态异常: %d (期望: %d)\n", user.Status, UserStatusNormal)
	}
	fmt.Println("✅ 用户状态正常")

	// 4. 验证密码
	fmt.Println("【步骤4】验证密码（模仿bcrypt.CompareHashAndPassword）")
	testPassword := "Test123456"
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(testPassword))
	if err != nil {
		log.Fatalf("❌ 密码验证失败: %v\n", err)
	}
	fmt.Printf("✅ 密码验证成功 ('%s')\n\n", testPassword)

	// 5. 测试错误的手机号
	fmt.Println("【步骤5】测试错误的手机号")
	wrongUser := &User{}
	err = db.GetContext(ctx, wrongUser, query, "99999999999")
	if err != nil {
		fmt.Printf("✅ 错误手机号正确返回错误: %v\n\n", err)
	} else {
		fmt.Printf("❌ 错误手机号应该返回错误，但查询到了用户\n\n")
	}

	// 6. 测试deleted_at
	fmt.Println("【步骤6】检查deleted_at条件")
	// 查询包括deleted_at的记录
	var deletedAtCount int
	countQuery := `SELECT COUNT(*) FROM users WHERE phone = ?`
	err = db.GetContext(ctx, &deletedAtCount, countQuery, "13800000000")
	if err != nil {
		log.Printf("查询总数失败: %v\n", err)
	} else {
		fmt.Printf("✅ 该手机号总记录数: %d\n", deletedAtCount)
	}

	activeQuery := `SELECT COUNT(*) FROM users WHERE phone = ? AND deleted_at IS NULL`
	var activeCount int
	err = db.GetContext(ctx, &activeCount, activeQuery, "13800000000")
	if err != nil {
		log.Printf("查询活跃记录失败: %v\n", err)
	} else {
		fmt.Printf("✅ 该手机号活跃记录数: %d\n\n", activeCount)
	}

	// 7. 模拟完整登录流程
	fmt.Println("【步骤7】模拟完整登录流程")
	fmt.Println("🔍 步骤1: Repository.GetByPhone(\"13800000000\")")
	loginUser := &User{}
	err = db.GetContext(ctx, loginUser, query, "13800000000")
	if err != nil {
		fmt.Printf("❌ 返回错误: %v (应该被处理为 ErrNotFound)\n", err)
	} else {
		fmt.Printf("✅ 返回用户: ID=%d, Phone=%s\n", loginUser.ID, loginUser.Phone)
	}

	fmt.Println("\n🔍 步骤2: 检查用户状态")
	if loginUser.Status != UserStatusNormal {
		fmt.Printf("❌ 用户状态: %d (应该返回 ErrUserDisabled)\n", loginUser.Status)
	} else {
		fmt.Printf("✅ 用户状态: %d (正常)\n", loginUser.Status)
	}

	fmt.Println("\n🔍 步骤3: 验证密码")
	err = bcrypt.CompareHashAndPassword([]byte(loginUser.Password), []byte(testPassword))
	if err != nil {
		fmt.Printf("❌ 密码错误: %v (应该返回 ErrPasswordIncorrect)\n", err)
	} else {
		fmt.Printf("✅ 密码正确 (应该继续生成Token)\n")
	}

	// 8. 总结
	fmt.Println("\n==================================================")
	fmt.Println("【诊断结论】")
	fmt.Println("==================================================")
	fmt.Println("✅ 数据库连接正常")
	fmt.Println("✅ 用户数据存在")
	fmt.Println("✅ 密码hash正确")
	fmt.Println("✅ 用户状态正常")
	fmt.Println("")
	fmt.Println("❌ 后端API仍返回401")
	fmt.Println("")
	fmt.Println("【可能原因】")
	fmt.Println("1. 后端连接的数据库配置不对")
	fmt.Println("2. Repository层的查询逻辑有问题")
	fmt.Println("3. 中间件或日志导致错误处理异常")
	fmt.Println("4. 后端代码的sqlx配置与测试工具不同")
	fmt.Println("")
	fmt.Println("【建议下一步】")
	fmt.Println("1. 检查后端日志: tail -f /tmp/server.log")
	fmt.Println("2. 在后端添加更多日志输出")
	fmt.Println("3. 使用pprof或dlv调试后端代码")
}
