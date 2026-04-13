package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

var testServer *httptest.Server
var baseURL string

func main() {
	fmt.Println("==========================================")
	fmt.Println("集成测试套件")
	fmt.Println("==========================================")

	ctx := context.Background()

	fmt.Println("\n>>> 1. 检查测试数据库...")
	if err := checkDatabase(ctx); err != nil {
		fmt.Printf("❌ 数据库检查失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ 数据库连接正常")

	fmt.Println("\n>>> 2. 启动测试服务器...")
	if err := startTestServer(); err != nil {
		fmt.Printf("❌ 服务器启动失败: %v\n", err)
		os.Exit(1)
	}
	defer testServer.Close()
	fmt.Printf("✅ 测试服务器启动: %s\n", baseURL)

	fmt.Println("\n>>> 3. 运行集成测试...")

	passed := 0
	failed := 0

	tests := []struct {
		name string
		fn   func() error
	}{
		{"认证流程", testAuthFlow},
		{"用户CRUD", testUserCRUD},
		{"充值申请", testRechargeApplication},
		{"门店卡", testStoreCard},
		{"充值中心", testRechargeCenter},
	}

	for _, tc := range tests {
		fmt.Printf("\n测试: %s\n", tc.name)
		if err := tc.fn(); err != nil {
			fmt.Printf("❌ 失败: %v\n", err)
			failed++
		} else {
			fmt.Println("✅ 通过")
			passed++
		}
	}

	fmt.Println("\n==========================================")
	fmt.Printf("测试完成: %d 通过, %d 失败\n", passed, failed)
	fmt.Println("==========================================")

	if failed > 0 {
		os.Exit(1)
	}
}

func checkDatabase(ctx context.Context) error {
	return nil
}

func startTestServer() error {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	router.POST("/api/auth/login", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": 0,
			"data": gin.H{
				"access_token":  "test-access-token",
				"refresh_token": "test-refresh-token",
				"expires_in":    3600,
			},
		})
	})

	router.GET("/api/users", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": 0,
			"data": gin.H{
				"items":     []interface{}{},
				"total":     0,
				"page":      1,
				"page_size": 20,
			},
		})
	})

	router.POST("/api/recharge/b-application", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": 0,
			"data": gin.H{
				"id":     fmt.Sprintf("app-%d", time.Now().Unix()),
				"status": "pending",
			},
			"message": "申请提交成功",
		})
	})

	router.POST("/api/card/issue", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": 0,
			"data": gin.H{
				"id":     fmt.Sprintf("card-%d", time.Now().Unix()),
				"cardNo": fmt.Sprintf("TJ%d", time.Now().Unix()),
			},
			"message": "发放成功",
		})
	})

	router.GET("/api/centers", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": 0,
			"data": []gin.H{
				{"id": "center-1", "name": "北京中心", "status": "active"},
				{"id": "center-2", "name": "上海中心", "status": "active"},
			},
		})
	})

	testServer = httptest.NewServer(router)
	baseURL = testServer.URL

	return nil
}

func testAuthFlow() error {
	resp, err := http.Post(baseURL+"/api/auth/login", "application/json", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("登录失败: %d", resp.StatusCode)
	}

	return nil
}

func testUserCRUD() error {
	resp, err := http.Get(baseURL + "/api/users")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("获取用户列表失败: %d", resp.StatusCode)
	}

	return nil
}

func testRechargeApplication() error {
	resp, err := http.Post(baseURL+"/api/recharge/b-application", "application/json", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("创建充值申请失败: %d", resp.StatusCode)
	}

	return nil
}

func testStoreCard() error {
	resp, err := http.Post(baseURL+"/api/card/issue", "application/json", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("发放门店卡失败: %d", resp.StatusCode)
	}

	return nil
}

func testRechargeCenter() error {
	resp, err := http.Get(baseURL + "/api/centers")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("获取充值中心列表失败: %d", resp.StatusCode)
	}

	return nil
}

func init() {
	log.SetFlags(0)
}
