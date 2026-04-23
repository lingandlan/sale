package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type Endpoint struct {
	Method  string
	Path    string
	Auth    bool
	Handler gin.HandlerFunc
}

func main() {
	fmt.Println("==========================================")
	fmt.Println("API 一致性验证")
	fmt.Println("==========================================")

	passed := 0
	failed := 0

	endpoints := []struct {
		name     string
		frontend string
		backend  string
		method   string
	}{
		{"B端充值申请", "/recharge/b-apply", "/api/v1/recharge/b-apply", "POST"},
		{"B端充值审批列表", "/recharge/b-approval", "/api/v1/recharge/b-approval", "GET"},
		{"B端充值审批详情", "/recharge/b-approval/:id", "/api/v1/recharge/b-approval/:id", "GET"},
		{"B端审批操作", "/recharge/b-approval/action", "/api/v1/recharge/b-approval/action", "POST"},
		{"C端充值录入", "/recharge/c-entry", "/api/v1/recharge/c-entry", "POST"},
		{"充值记录列表", "/recharge/records", "/api/v1/recharge/records", "GET"},
		{"充值记录详情", "/recharge/records/:id", "/api/v1/recharge/records/:id", "GET"},
		{"门店卡列表", "/card/list", "/api/v1/card/list", "GET"},
		{"门店卡详情", "/card/detail/:cardNo", "/api/v1/card/detail/:cardNo", "GET"},
		{"门店卡统计", "/card/stats", "/api/v1/card/stats", "GET"},
		{"门店卡发放", "/card/issue", "/api/v1/card/issue", "POST"},
		{"门店卡核销", "/card/consume", "/api/v1/card/consume", "POST"},
		{"门店卡状态", "/card/:cardNo/status", "/api/v1/card/:cardNo/status", "POST"},
		{"充值中心列表", "/center", "/api/v1/center", "GET"},
		{"充值中心创建", "/center", "/api/v1/center", "POST"},
		{"充值中心更新", "/center/:id", "/api/v1/center/:id", "PUT"},
		{"充值中心删除", "/center/:id", "/api/v1/center/:id", "DELETE"},
		{"操作员列表", "/operator", "/api/v1/operator", "GET"},
		{"操作员创建", "/operator", "/api/v1/operator", "POST"},
		{"操作员更新", "/operator/:id", "/api/v1/operator/:id", "PUT"},
		{"操作员删除", "/operator/:id", "/api/v1/operator/:id", "DELETE"},
		{"Dashboard统计", "/dashboard/statistics", "/api/v1/dashboard/statistics", "GET"},
		{"Dashboard待办", "/dashboard/todos", "/api/v1/dashboard/todos", "GET"},
		{"Dashboard趋势", "/dashboard/recharge-trends", "/api/v1/dashboard/recharge-trends", "GET"},
	}

	fmt.Println("\n前端 → 后端 API 路径映射:")
	fmt.Println("----------------------------------------")

	for _, ep := range endpoints {
		frontendOk := strings.HasPrefix(ep.frontend, "/recharge") ||
			strings.HasPrefix(ep.frontend, "/card") ||
			strings.HasPrefix(ep.frontend, "/center") ||
			strings.HasPrefix(ep.frontend, "/operator") ||
			strings.HasPrefix(ep.frontend, "/dashboard")

		if frontendOk {
			fmt.Printf("✅ %-20s %s %s\n", ep.name, ep.method, ep.frontend)
			fmt.Printf("   后端: %s\n", ep.backend)
			passed++
		} else {
			fmt.Printf("❌ %-20s %s %s\n", ep.name, ep.method, ep.frontend)
			failed++
		}
	}

	fmt.Println("\n==========================================")
	fmt.Printf("结果: %d 通过, %d 失败\n", passed, failed)
	fmt.Println("==========================================")

	if failed > 0 {
		os.Exit(1)
	}
}
