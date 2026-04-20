package main

import (
	"flag"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Operator 只映射需要的字段
type Operator struct {
	ID       string `gorm:"primaryKey;size:64"`
	Name     string
	Password string
}

func (Operator) TableName() string {
	return "recharge_operators"
}

func main() {
	dryRun := flag.Bool("dry-run", false, "只打印将要执行的变更，不实际修改")
	dsn := flag.String("dsn", "sale:sale123@tcp(localhost:3306)/sale_dev?charset=utf8mb4&parseTime=True&loc=Local", "数据库 DSN")
	flag.Parse()

	db, err := gorm.Open(mysql.Open(*dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	var operators []Operator
	if err := db.Find(&operators).Error; err != nil {
		log.Fatalf("查询操作员失败: %v", err)
	}

	updated := 0
	for _, op := range operators {
		// 跳过已经是 bcrypt 哈希的密码
		if len(op.Password) > 4 && op.Password[:4] == "$2a$" {
			fmt.Printf("[SKIP] %s (%s): 已经是 bcrypt 哈希\n", op.ID, op.Name)
			continue
		}

		// 明文密码哈希
		hash, err := bcrypt.GenerateFromPassword([]byte(op.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("[ERROR] %s: 哈希失败: %v", op.ID, err)
			continue
		}

		if *dryRun {
			fmt.Printf("[DRY-RUN] %s (%s): '%s' -> '%s'\n", op.ID, op.Name, op.Password, string(hash))
		} else {
			if err := db.Model(&Operator{}).Where("id = ?", op.ID).Update("password", string(hash)).Error; err != nil {
				log.Printf("[ERROR] %s: 更新失败: %v", op.ID, err)
				continue
			}
			fmt.Printf("[OK] %s (%s): 密码已哈希\n", op.ID, op.Name)
		}
		updated++
	}

	action := "将要更新"
	if !*dryRun {
		action = "已更新"
	}
	fmt.Printf("\n%s %d 个操作员密码\n", action, updated)
}
