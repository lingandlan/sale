# Marketplace API 代码规范

## 1. 项目结构

```
marketplace/
├── cmd/
│   └── server/
│       └── main.go              # 应用入口
├── internal/
│   ├── config/                  # 配置
│   │   └── config.go
│   ├── handler/                 # HTTP 层（Controllers）
│   │   ├── auth.go
│   │   └── user.go
│   ├── service/                 # 业务逻辑层
│   │   ├── auth.go
│   │   └── user.go
│   ├── repository/              # 数据访问层
│   │   ├── user.go
│   │   └── product.go
│   ├── model/                   # 数据模型
│   │   ├── user.go
│   │   └── product.go
│   ├── middleware/              # 中间件
│   │   ├── auth.go
│   │   ├── logger.go
│   │   └── cors.go
│   └── pkg/                    # 公共工具
│       ├── response/            # 统一响应
│       │   └── response.go
│       └── errors/              # 统一错误
│           └── errors.go
├── migrations/                  # 数据库迁移
├── configs/                     # 配置文件
├── docs/                        # 文档
├── Makefile
└── docker-compose.yml
```

**原则：**
- `cmd/` 只能有一个入口，不放业务代码
- `internal/` 放置私有代码，对外不可见
- `pkg/` 放置可被外部引用的库代码
- 按功能分层，不按业务模块分（handler/service/repository 三层）

---

## 2. 命名规范

### 2.1 包名

```go
// ✅ 小写字母，无下划线，与目录名一致
package user        // user/user.go
package product     // product/product.go

// ❌ 错误示例
package UserService
package user_service
```

### 2.2 变量名

```go
// ✅ 驼峰命名，简洁表达含义
userID      int64
userName    string
totalCount  int
isActive    bool

// ❌ 错误示例
UserID           // 不要用 PascalCase
iUserId          // 不要用前缀
user_id          // 不要用下划线
strUserName      // 不要用类型前缀
```

### 2.3 函数名

```go
// ✅ 动词前缀 + 名词
func GetUserByID(id int64) (*User, error)
func CreateOrder(req *CreateOrderReq) (*Order, error)
func DeleteUser(id int64) error
func ListProducts(filter *ProductFilter) ([]*Product, int64, error)

// ❌ 错误示例
func User(id int64)                    // 不清晰
func GetAllUserData(userID int64)      // 冗余
```

### 2.4 接口名

```go
// ✅ 单数名词，或 "er" 后缀
type UserRepository interface {}
type CacheService interface {}
type Logger interface {}

// 如果接口只有一个方法，以方法名命名
type Finder interface {
    FindByID(id int64) (*Entity, error)
}
```

### 2.5 常量名

```go
// ✅ 全大写下划线分隔
const MaxRetryCount = 3
const DefaultPageSize = 20
const TokenExpireHours = 24

// ✅ 枚举类型用 iota + 类型名
type OrderStatus int
const (
    OrderStatusPending   OrderStatus = 1
    OrderStatusPaid     OrderStatus = 2
    OrderStatusCompleted OrderStatus = 3
)
```

---

## 3. 错误处理

### 3.1 错误返回

```go
// ✅ 返回 error，调用方决定如何处理
func GetUser(id int64) (*User, error) {
    user, err := repo.FindByID(id)
    if err != nil {
        return nil, err  // 直接返回
    }
    return user, nil
}

// ❌ 不要忽略错误
result, _ := repo.FindByID(id)  // 错误被忽略！
```

### 3.2 错误包装

```go
import "fmt"

func GetUser(id int64) (*User, error) {
    user, err := repo.FindByID(id)
    if err != nil {
        return nil, fmt.Errorf("GetUser: %w", err)
    }
    return user, nil
}
```

### 3.3 自定义错误类型

```go
// pkg/errors/errors.go
var (
    ErrNotFound     = errors.New("resource not found")
    ErrUnauthorized = errors.New("unauthorized")
    ErrForbidden    = errors.New("forbidden")
    ErrDuplicate    = errors.New("resource already exists")
)

// 在业务中使用
func GetUser(id int64) (*User, error) {
    user, err := repo.FindByID(id)
    if errors.Is(err, sql.ErrNoRows) {
        return nil, ErrNotFound
    }
    return user, nil
}
```

### 3.4 错误判断

```go
// ✅ 使用 errors.Is / errors.As
if errors.Is(err, ErrNotFound) {
    // 处理 404
}

// ❌ 不要这样判断
if err.Error() == "resource not found" {
}
```

---

## 4. API 设计规范

### 4.1 URL 路径

```
✅ RESTful 风格
GET    /api/v1/users           # 列表
GET    /api/v1/users/:id       # 详情
POST   /api/v1/users           # 创建
PUT    /api/v1/users/:id       # 更新
DELETE /api/v1/users/:id       # 删除

✅ 资源嵌套（适度）
GET    /api/v1/users/:id/orders      # 某用户的订单
GET    /api/v1/orders/:id/items      # 某订单的商品

❌ 避免
GET    /api/v1/getUser
POST   /api/v1/createUser
```

### 4.2 HTTP 方法

| 方法 | 用途 | 幂等 |  Body |
|------|------|------|-------|
| GET | 查询 | ✅ | 无 |
| POST | 创建 | ❌ | 有 |
| PUT | 全量更新 | ✅ | 有 |
| PATCH | 部分更新 | ❌ | 有 |
| DELETE | 删除 | ✅ | 无 |

### 4.3 请求与响应

```go
// ✅ 统一请求结构
type CreateUserReq struct {
    Username string `json:"username" binding:"required,min=3,max=32"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}

type UpdateUserReq struct {
    Nickname *string `json:"nickname"`           // 可选字段用指针
    Email    *string `json:"email" binding:"omitempty,email"`
}

// ✅ 统一响应格式
type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

// ✅ 分页响应
type ListResponse struct {
    Items      interface{} `json:"items"`
    Total      int64       `json:"total"`
    Page       int         `json:"page"`
    PageSize   int         `json:"page_size"`
}

// ✅ 错误响应
type ErrorResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Detail  string `json:"detail,omitempty"`
}
```

### 4.4 HTTP 状态码

| 状态码 | 用途 |
|--------|------|
| 200 | 成功 |
| 201 | 创建成功 |
| 204 | 删除成功（无内容）|
| 400 | 请求参数错误 |
| 401 | 未认证 |
| 403 | 无权限 |
| 404 | 资源不存在 |
| 409 | 资源冲突（如重复）|
| 500 | 服务器内部错误 |

---

## 5. 注释规范

### 5.1 包注释

```go
// Package user 提供用户相关的业务逻辑和数据访问。
//
// 在 internal/model/user.go 中定义用户模型。
// 在 internal/service/user.go 中实现业务逻辑。
package user
```

### 5.2 函数注释

```go
// GetUserByID 根据用户 ID 获取用户信息。
//
// 如果用户不存在，返回 ErrNotFound。
//
// 参数:
//   - id: 用户 ID
//
// 返回:
//   - *User: 用户信息
//   - error: 错误信息
func GetUserByID(id int64) (*User, error) {
    // ...
}
```

### 5.3 注释风格

```go
// ✅ 句子式注释（Go 官方风格）
// GetUserByID returns user by id.
func GetUserByID(id int64) (*User, error) {}

// ❌ 不要用 Markdown 列表或完整句子
// - id: 用户ID
// - returns: 用户信息 or error
```

### 5.4 公共 API 必须注释

```go
// ✅ 导出的函数/类型必须注释
func PublicFunction() {}

// ❌ 未导出的函数可以不注释
func privateFunction() {}
```

---

## 6. 数据库模型规范

### 6.1 模型定义

```go
type User struct {
    ID        int64          `gorm:"primaryKey;autoIncrement" json:"id"`
    Username  string         `gorm:"size:64;uniqueIndex;not null" json:"username"`
    Password  string         `gorm:"size:255;not null" json:"-"`          // json:"-" 不暴露
    Nickname  string         `gorm:"size:64" json:"nickname"`
    Email     string         `gorm:"size:128;uniqueIndex" json:"email"`
    Role      int            `gorm:"default:0" json:"role"`
    Status    int            `gorm:"default:1" json:"status"`
    CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`                       // 软删除
}

func (User) TableName() string {
    return "users"
}
```

### 6.2 索引规范

```go
// ✅ 有唯一性要求的字段用 uniqueIndex
Username string `gorm:"uniqueIndex"`

// ✅ 常用查询字段加普通索引
Email string `gorm:"index"`

// ✅ 复合索引
UserID int64 `gorm:"index:idx_user_status"`
Status int   `gorm:"index:idx_user_status"`
```

---

## 7. 日志规范

### 7.1 日志级别

```go
logger.Debug("debug message")   // 开发调试
logger.Info("info message")      // 一般信息
logger.Warn("warn message")     // 警告（可恢复的错误）
logger.Error("error message")   // 错误
```

### 7.2 日志内容

```go
// ✅ 包含上下文信息
logger.Info("user login",
    zap.Int64("user_id", userID),
    zap.String("ip", ip),
    zap.Duration("cost", cost),
)

// ❌ 避免
logger.Info("login")                    // 无上下文
logger.Error(err.Error())              // 没有字段
```

---

## 8. 测试规范

### 8.1 测试文件命名

```
user_test.go        // 单元测试
user_integration_test.go  // 集成测试
```

### 8.2 测试函数命名

```go
func TestGetUserByID(t *testing.T) {}           // 单元测试
func TestGetUserByID_WithInvalidID(t *testing.T) {}  // 特殊情况

func TestIntegration_UserFlow(t *testing.T) {}  // 集成测试
```

### 8.3 表格驱动测试

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name   string
        a      int
        b      int
        expect int
    }{
        {"positive + positive", 1, 2, 3},
        {"positive + negative", 1, -1, 0},
        {"zero + zero", 0, 0, 0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Add(tt.a, tt.b)
            if got != tt.expect {
                t.Errorf("Add(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.expect)
            }
        })
    }
}
```

---

## 9. 配置规范

### 9.1 环境变量命名

```yaml
# ✅ 全大写下划线分隔
DATABASE_HOST=localhost
DATABASE_PORT=3306
JWT_SECRET=secret
REDIS_HOST=127.0.0.1

# ❌ 避免
DatabaseHost
database-host
dbhost
```

### 9.2 敏感信息

```go
// ✅ 不在代码中硬编码敏感信息
// ❌ 错误
password := "123456"

// ✅ 从环境变量或配置中读取
password := viper.GetString("database.password")
```

---

## 10. Git 提交规范

### 10.1 格式

```
<type>(<scope>): <subject>

feat(auth): add login endpoint
fix(cart): resolve quantity calculation bug
docs(readme): update installation guide
style(handler): format code style
refactor(service): simplify user lookup logic
test(user): add unit tests for GetUserByID
chore(deps): upgrade gin to v1.9.0
perf(order): optimize query with index
ci: add github actions workflow
```

### 10.2 Type 说明

| Type | 说明 |
|------|------|
| feat | 新功能 |
| fix | Bug 修复 |
| docs | 文档更新 |
| style | 代码格式（不影响功能）|
| refactor | 重构 |
| test | 测试相关 |
| chore | 构建/工具/依赖 |
| perf | 性能优化 |
| ci | CI/CD |
| build | 构建系统 |

---

## 11. 目录组织原则

### 三层架构（推荐）

```
handler (HTTP 层)
    ↓ 调用
service (业务逻辑层)
    ↓ 调用
repository (数据访问层)
    ↓ 调用
database / cache
```

**职责划分：**
- `handler`: 参数校验、调用 service、返回响应
- `service`: 业务逻辑、事务控制
- `repository`: 数据库 CRUD、查询

**禁止：**
- ❌ handler 直接操作数据库
- ❌ service 处理 HTTP 请求/响应
- ❌ 跨层调用（handler → repository）
