package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"marketplace/backend/internal/config"
	"marketplace/backend/internal/model"
	"marketplace/backend/internal/repository"

	"gorm.io/gorm"
)

func main() {
	// 加载配置
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 连接数据库（GORM）
	if host := os.Getenv("APP_DATABASE_HOST"); host != "" {
		cfg.Database.Host = host
	}
	if port := os.Getenv("APP_DATABASE_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Database.Port = p
		}
	}
	if user := os.Getenv("APP_DATABASE_USER"); user != "" {
		cfg.Database.User = user
	}
	if name := os.Getenv("APP_DATABASE_NAME"); name != "" {
		cfg.Database.Name = name
	}

	fmt.Printf("[DEBUG] migrate 使用数据库: host=%s port=%d user=%s db=%s\n",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Name)

	// 连接数据库（GORM）
	db, err := repository.NewGormDB(&cfg.Database)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	fmt.Println("开始数据库迁移...")

	// 先删除旧 SQL migration 创建的、与 GORM model 冲突的表
	dropOldTables(db)

	// 迁移所有表
	err = db.AutoMigrate(
		&model.User{},
		&model.RechargeApplication{},
		&model.CRecharge{},
		&model.StoreCard{},
		&model.CardIssueRecord{},
		&model.CardTransaction{},
		&model.RechargeCenter{},
		&model.RechargeOperator{},
		&model.CenterMonthlyConsumption{},
	)
	if err != nil {
		log.Fatalf("迁移失败: %v", err)
	}

	fmt.Println("数据库迁移成功!")

	// 创建索引
	createIndexes(db)

	// 插入初始数据
	insertSeedData(db)

	fmt.Println("初始化完成!")
}

// dropOldTables 删除旧 SQL migration 创建的、与 GORM model 冲突的表
func dropOldTables(db *gorm.DB) {
	fmt.Println("清理旧表...")
	oldTables := []string{
		"member_card_transactions",
		"member_cards",
		"recharge_records",
		"stores",
	}
	for _, table := range oldTables {
		db.Exec("DROP TABLE IF EXISTS " + table)
	}
	fmt.Println("旧表清理完成!")
}

// createIndexes 创建数据库索引
func createIndexes(db *gorm.DB) {
	fmt.Println("索引已通过 GORM AutoMigrate 自动创建，跳过手动索引")
}

// insertSeedData 插入初始数据
func insertSeedData(db *gorm.DB) {
	fmt.Println("插入初始数据...")

	// 创建默认管理员用户
	admin := &model.User{
		Phone:   "13800000000",
		Password: "$2a$10$vXPjbfC511sMp3zdk1uFzOfxRWmtsZXnNIX7buP4C9Aq6In5YhV5S", // 123456
		Name:   "系统管理员",
		Role:   model.RoleHQAdmin,
		Status: model.UserStatusNormal,
	}
	result := db.Where("phone = ?", admin.Phone).FirstOrCreate(admin)
	if result.Error == nil {
		fmt.Println("创建管理员用户: 13800000000 / 123456")
	}

	// 创建测试充值中心
	centers := []model.RechargeCenter{
		{ID: "center-bj-cy", Code: "BJ_CY", Name: "北京朝阳中心", Address: "北京市朝阳区", Status: "active"},
		{ID: "center-bj-hd", Code: "BJ_HD", Name: "北京海淀中心", Address: "北京市海淀区", Status: "active"},
		{ID: "center-sh-pd", Code: "SH_PD", Name: "上海浦东中心", Address: "上海市浦东新区", Status: "active"},
	}
	for _, center := range centers {
		db.Where("code = ?", center.Code).FirstOrCreate(&center)
	}
	fmt.Println("创建充值中心数据")

	// 创建测试操作员
	operators := []model.RechargeOperator{
		{ID: "op-zhangcw", Name: "张财务", Phone: "13800138001", CenterID: "center-bj-cy", Role: "财务", Status: "active"},
		{ID: "op-licn", Name: "李出纳", Phone: "13800138002", CenterID: "center-bj-cy", Role: "出纳", Status: "active"},
	}
	for _, op := range operators {
		op.Password = "$2a$10$vXPjbfC511sMp3zdk1uFzOfxRWmtsZXnNIX7buP4C9Aq6In5YhV5S" // 123456
		db.Where("phone = ?", op.Phone).FirstOrCreate(&op)
	}
	fmt.Println("创建操作员数据")

	fmt.Println("初始数据插入完成!")
}
