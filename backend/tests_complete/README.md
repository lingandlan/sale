# 测试方案总结

## 📋 概述

本文档总结项目的完整测试方案，包括单元测试、集成测试和 E2E 测试。

---

## 🧪 测试分层

```
┌─────────────────────────────────────────────────────────────┐
│                    测试金字塔                                │
├─────────────────────────────────────────────────────────────┤
│  E2E 测试 (10%)                                            │
│  - Playwright 浏览器自动化                                  │
│  - 验证前后端集成                                           │
│  - 测试完整用户流程                                          │
├─────────────────────────────────────────────────────────────┤
│  集成测试 (20%)                                            │
│  - API 路径一致性                                          │
│  - 真实数据库操作                                          │
│  - Repository 层测试                                        │
├─────────────────────────────────────────────────────────────┤
│  单元测试 (70%)                                            │
│  - Handler 层 (Mock Service)                              │
│  - Service 层 (Mock Repository)                           │
│  - Repository 层 (SQLite 内存数据库)                       │
└─────────────────────────────────────────────────────────────┘
```

---

## 📁 测试文件结构

```
backend/
├── internal/
│   ├── handler/
│   │   ├── auth_test.go        # 认证 Handler 测试
│   │   ├── user_test.go        # 用户 Handler 测试
│   │   ├── recharge_test.go    # 充值 Handler 测试
│   │   └── admin_test.go       # 管理员 Handler 测试
│   │
│   ├── service/
│   │   ├── auth_test.go        # 认证 Service 测试
│   │   ├── user_test.go        # 用户 Service 测试
│   │   ├── recharge_test.go     # 充值 Service 测试
│   │   └── bench_test.go        # 性能基准测试
│   │
│   ├── repository/
│   │   ├── user_test.go         # 用户 Repository 测试
│   │   └── recharge_test.go     # 充值 Repository 测试
│   │
│   └── middleware/
│       └── auth_test.go          # 中间件测试
│
├── tests_complete/               # 完整测试套件
│   ├── e2e/
│   │   ├── api.spec.ts         # E2E API 测试
│   │   ├── playwright.config.ts # Playwright 配置
│   │   └── package.json
│   │
│   ├── integration/
│   │   └── main.go              # 集成测试入口
│   │
│   └── api_consistency_check.go  # API 一致性验证
│
└── tests/                       # 原有测试文件
    ├── integration/             # 集成测试
    └── e2e/                    # E2E 测试
```

---

## 🚀 运行测试

### 1. 单元测试

```bash
# 运行所有单元测试
go test ./internal/... -v

# 带覆盖率
go test ./internal/... -coverprofile=coverage.out -covermode=atomic
go tool cover -html=coverage.out -o coverage.html

# 运行特定包
go test ./internal/handler/... -v
go test ./internal/service/... -v
go test ./internal/repository/... -v
```

### 2. 集成测试

```bash
# 启动测试数据库
docker-compose up -d

# 运行集成测试
go run ./tests_complete/integration/main.go

# API 一致性检查
go run ./tests_complete/api_consistency_check.go
```

### 3. E2E 测试

```bash
cd tests_complete/e2e

# 安装依赖
npm install

# 安装 Playwright 浏览器
npx playwright install

# 运行测试
npx playwright test

# 查看报告
npx playwright show-report
```

### 4. 运行全部测试

```bash
# 使用集成测试脚本
./tests_complete/run_all_tests.sh

# 或分别运行
./tests_complete/run_unit_tests.sh
```

---

## 📊 当前覆盖率

| 模块 | 覆盖率 | 状态 |
|-----|-------|------|
| handler | 41.2% | ✅ |
| middleware | 45.3% | ✅ |
| repository | 62.6% | ✅ |
| service | 49.4% | ✅ |
| **总计** | **~50%** | 🟡 |

**目标覆盖率**: 70%

---

## 🧩 测试类型详解

### 1. 单元测试 (Unit Tests)

#### Handler 层测试
- **策略**: Mock Service 层
- **目的**: 测试 HTTP 请求/响应处理
- **示例**:
```go
func TestAuthHandler_Login_Success(t *testing.T) {
    mockAuthSvc := new(MockAuthService)
    h := NewAuthHandler(mockAuthSvc, mockUserSvc)
    
    expectedResp := &model.LoginResponse{
        AccessToken: "test-token",
    }
    mockAuthSvc.On("Login", ctx, phone, password).Return(expectedResp, nil)
    
    // 模拟 HTTP 请求
    req, _ := http.NewRequest("POST", "/auth/login", body)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
}
```

#### Service 层测试
- **策略**: Mock Repository 层
- **目的**: 测试业务逻辑
- **示例**:
```go
func TestRechargeService_CalculatePoints(t *testing.T) {
    svc := &RechargeService{}
    
    points, basePoints, rebatePoints := svc.CalculatePoints(10000, 120000)
    assert.Equal(t, 10200, points)  // 10000 + 200 (2% rebate)
}
```

#### Repository 层测试
- **策略**: 使用 SQLite 内存数据库
- **目的**: 测试数据库操作
- **示例**:
```go
func TestUserRepository_CRUD(t *testing.T) {
    db, cleanup := setupTestDB(t)
    defer cleanup()
    repo := NewUserRepository(db)
    
    user := &model.User{Phone: "13800000000", Name: "Test"}
    id, err := repo.Create(ctx, user)
    require.NoError(t, err)
    assert.Greater(t, id, int64(0))
}
```

### 2. 集成测试 (Integration Tests)

#### API 路径一致性测试
```go
// tests_complete/api_consistency_check.go
// 验证前端调用的 API 路径与后端路由一致
func main() {
    endpoints := []struct {
        name     string
        frontend string  // 前端定义
        backend  string  // 后端实际路径
    }{
        {"B端充值申请", "/recharge/b-apply", "/api/v1/recharge/b-apply"},
        {"充值记录列表", "/recharge/records", "/api/v1/recharge/records"},
        // ...
    }
}
```

#### Repository 集成测试
- 使用真实数据库连接
- 验证 SQL 查询正确性
- 测试事务处理

### 3. E2E 测试 (End-to-End)

#### Playwright API 测试
```typescript
// tests_complete/e2e/api.spec.ts
test.describe('充值记录', () => {
    let authToken: string;

    test.beforeAll(async ({ request }) => {
        const response = await request.post(`${API_BASE}/auth/login`, {
            data: { phone: '13800000000', password: '123456' }
        });
        const body = await response.json();
        authToken = body.data?.access_token;
    });

    test('获取充值记录列表', async ({ request }) => {
        const response = await request.get(`${API_BASE}/recharge/records`, {
            headers: { Authorization: `Bearer ${authToken}` }
        });
        expect(response.ok()).toBeTruthy();
    });
});
```

---

## 🔍 API 路径映射

### 前端 → 后端

| 功能模块 | 前端路径 | 后端路径 | 方法 |
|---------|---------|---------|------|
| 登录 | `/auth/login` | `/api/v1/auth/login` | POST |
| B端充值申请 | `/recharge/b-apply` | `/api/v1/recharge/b-apply` | POST |
| B端充值审批 | `/recharge/b-approval` | `/api/v1/recharge/b-approval` | GET |
| C端充值录入 | `/recharge/c-entry` | `/api/v1/recharge/c-entry` | POST |
| 充值记录 | `/recharge/records` | `/api/v1/recharge/records` | GET |
| 门店卡列表 | `/card/list` | `/api/v1/card/list` | GET |
| 充值中心 | `/center` | `/api/v1/center` | GET/POST |
| Dashboard | `/dashboard/statistics` | `/api/v1/dashboard/statistics` | GET |

---

## 🐛 测试发现的问题

### 已修复
1. ✅ `RoleAdmin` 常量不存在 → 使用 `RoleOperator`, `RoleFinance` 等
2. ✅ `Username/Email/Nickname` 字段 → 改为 `Phone/Name`
3. ✅ `GetByUsername` → `GetByPhone`

### 待修复
1. ❌ `/api/recharge/records` 返回 404 - 后端路由可能缺失
2. ❌ `AdminHandler` 缺少 `CreateUser/UpdateUser` 方法
3. ❌ Repository 层部分方法未实现

---

## ⚙️ CI/CD 集成

### Pre-push Hook
```bash
# .git/hooks/pre-push
- 运行单元测试
- 检查代码覆盖率
- 执行 golangci-lint
```

### GitHub Actions
```yaml
# .github/workflows/test.yml
- 单元测试
- 集成测试
- E2E 测试
- 覆盖率报告
```

---

## 📝 添加新测试

### 1. Handler 测试
```go
// internal/handler/{module}_test.go
func Test{Module}_{Handler}_{Scenario}(t *testing.T) {
    mockService := new(MockService)
    h := NewHandler(mockService)
    // 测试逻辑
}
```

### 2. Service 测试
```go
// internal/service/{module}_test.go
func Test{Module}_{Method}_{Scenario}(t *testing.T) {
    mockRepo := new(MockRepo)
    svc := NewService(mockRepo)
    // 测试逻辑
}
```

### 3. E2E 测试
```typescript
// tests_complete/e2e/api.spec.ts
test('{描述}', async ({ request }) => {
    const response = await request.{method}(`${API_BASE}{path}`);
    expect(response.ok()).toBeTruthy();
});
```

---

## 📚 参考资料

- [Go Testing](https://pkg.go.dev/testing)
- [Testify](https://github.com/stretchr/testify)
- [Playwright](https://playwright.dev/)
- [SQLite Testing](https://github.com/glebarez/sqlite)
