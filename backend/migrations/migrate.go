package main

import (
	"fmt"
	"log"

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
	db, err := repository.NewGormDB(&cfg.Database)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	fmt.Println("开始数据库迁移...")

	// 迁移所有表
	err = db.AutoMigrate(
		&model.User{},
		&model.RechargeApplication{},
		&model.CRecharge{},
		&model.StoreCard{},
		&model.CardTransaction{},
		&model.RechargeCenter{},
		&model.RechargeOperator{},
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

// createIndexes 创建数据库索引
func createIndexes(db *gorm.DB) {
	fmt.Println("创建索引...")

	db.Exec("CREATE INDEX IF NOT EXISTS idx_recharge_applications_status ON recharge_applications(status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_recharge_applications_center ON recharge_applications(center_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_recharge_applications_applicant ON recharge_applications(applicant_id)")

	db.Exec("CREATE INDEX IF NOT EXISTS idx_c_recharges_member_phone ON c_recharges(member_phone)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_c_recharges_center ON c_recharges(center_id)")

	db.Exec("CREATE INDEX IF NOT EXISTS idx_store_cards_card_no ON store_cards(card_no)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_store_cards_holder_phone ON store_cards(holder_phone)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_store_cards_status ON store_cards(status)")

	db.Exec("CREATE INDEX IF NOT EXISTS idx_card_transactions_card_no ON card_transactions(card_no)")

	fmt.Println("索引创建成功!")
}

// insertSeedData 插入初始数据
func insertSeedData(db *gorm.DB) {
	fmt.Println("插入初始数据...")

	// 创建默认管理员用户
	admin := &model.User{
		Phone:   "13800000000",
		Password: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // admin123
		Name:   "系统管理员",
		Role:   model.RoleHQAdmin,
		Status: model.UserStatusNormal,
	}
	result := db.Where("phone = ?", admin.Phone).FirstOrCreate(admin)
	if result.Error == nil {
		fmt.Println("创建管理员用户: 13800000000 / admin123")
	}

	// 创建测试充值中心
	centers := []model.RechargeCenter{
		{Code: "BJ_CY", Name: "北京朝阳中心", Address: "北京市朝阳区", Status: "active"},
		{Code: "BJ_HD", Name: "北京海淀中心", Address: "北京市海淀区", Status: "active"},
		{Code: "SH_PD", Name: "上海浦东中心", Address: "上海市浦东新区", Status: "active"},
	}
	for _, center := range centers {
		db.Where("code = ?", center.Code).FirstOrCreate(&center)
	}
	fmt.Println("创建充值中心数据")

	// 创建测试操作员
	operators := []model.RechargeOperator{
		{Name: "张财务", Phone: "13800138001", CenterID: "1", Role: "财务", Status: "active"},
		{Name: "李出纳", Phone: "13800138002", CenterID: "1", Role: "出纳", Status: "active"},
	}
	for _, op := range operators {
		op.Password = "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy" // admin123
		db.Where("phone = ?", op.Phone).FirstOrCreate(&op)
	}
	fmt.Println("创建操作员数据")

	fmt.Println("初始数据插入完成!")
}
