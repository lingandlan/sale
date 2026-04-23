# 登录401问题解决报告

## 问题描述
后端登录接口 `/api/v1/auth/login` 返回 401 "手机号或密码错误"，但数据库验证通过。

## 根本原因
发现了3个关键问题：

### 1. 数据库NULL字段处理错误
**问题**: `model.User` 中的可空字段使用 `string` 而非 `*string`
**影响**: sqlx无法扫描NULL值到string类型
**表现**: `sql: Scan error on column index 6, name "center_name": converting NULL to string is unsupported`

**修复**:
```go
// internal/model/user.go
CenterName   *string  `db:"center_name"`
LastLoginIP  *string  `db:"last_login_ip"`
```

### 2. JWT Claims结构不匹配
**问题**: 中间件和Service使用不同的Claims结构
**影响**: Token解析失败，返回"invalid token"

**中间件的定义** (错误):
```go
type Claims struct {
    UserID   int64  `json:"user_id"`
    Username string `json:"username"`  // ❌ 字段不存在
    Role     int    `json:"role"`       // ❌ 类型错误
}
```

**Service的定义** (正确):
```go
type Claims struct {
    UserID int64  `json:"user_id"`
    Phone  string `json:"phone"`      // ✅ 正确字段
    Role   string `json:"role"`       // ✅ string类型
}
```

**修复**:
```go
// internal/middleware/auth.go
// 统一使用service.Claims
claims := &service.Claims{}

// 更新Context字段
c.Set("user_id", claims.UserID)
c.Set("phone", claims.Phone)
c.Set("role", claims.Role)
```

### 3. 角色类型不一致
**问题**: RBAC中间件中使用int类型角色判断
**影响**: 类型断言失败导致panic

**修复**:
```go
// internal/middleware/rbac.go
// ❌ 错误
if role.(int) == 2 {

// ✅ 修复
if role.(string) == "super_admin" {
```

## 解决方案

### 修改文件列表
1. `internal/model/user.go` - 修复可空字段类型
2. `internal/service/user.go` - 修复字段赋值
3. `internal/middleware/auth.go` - 统一Claims结构
4. `internal/middleware/rbac.go` - 修复角色类型判断

### 修改详情

#### 1. 修复可空字段类型
```go
// internal/model/user.go
type User struct {
    // ...
    CenterID     *uint      `db:"center_id"`
    CenterName   *string    `db:"center_name"`    // 改为指针
    LastLoginAt  *time.Time `db:"last_login_at"`
    LastLoginIP  *string    `db:"last_login_ip"`   // 改为指针
    // ...
}
```

#### 2. 统一Claims结构
```go
// internal/middleware/auth.go

// 删除重复的Claims定义
// 直接使用service.Claims

func (m *AuthMiddleware) Auth() gin.HandlerFunc {
    // ...
    claims := &service.Claims{}  // 使用service的Claims
    token, err := jwt.ParseWithClaims(tokenString, claims, ...)
    // ...
    c.Set("user_id", claims.UserID)
    c.Set("phone", claims.Phone)
    c.Set("role", claims.Role)
}
```

#### 3. 更新辅助函数
```go
// internal/middleware/auth.go
func GetPhone(c *gin.Context) string {  // 新增
    phone, _ := c.Get("phone")
    return phone.(string)
}

func GetRole(c *gin.Context) string {  // 返回类型改为string
    role, _ := c.Get("role")
    return role.(string)
}
```

#### 4. 修复角色判断
```go
// internal/middleware/rbac.go

// 超级管理员判断
if role.(string) == "super_admin" {
    c.Next()
    return
}

// RequireRoles函数
if userRole.(string) == r {  // 直接比较string
    c.Next()
    return
}
```

## 测试验证

### 测试账号
```
手机号: 13800000000
密码: Test123456
角色: super_admin
```

### 测试结果
```bash
# 1. 登录接口
POST /api/v1/auth/login
✅ 返回 access_token 和 refresh_token

# 2. 获取用户信息
GET /api/v1/user/info
Authorization: Bearer {token}
✅ 返回完整用户信息

# 3. 获取用户列表（管理员）
GET /api/v1/admin/users
Authorization: Bearer {token}
✅ 返回用户列表，分页信息正确
```

### 完整响应示例
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "id": 1,
        "phone": "13800000000",
        "name": "超级管理员",
        "role": "super_admin",
        "status": 1,
        "created_at": "2026-04-10T10:32:57+08:00",
        "updated_at": "2026-04-10T10:53:04+08:00"
    }
}
```

## 工具支持

### 密码重置工具
```bash
cd /Users/zhangdaodong/code/sale/backend
go run ./cmd/resetpwd/main.go -phone 13800000000 -password 新密码
```

### 诊断工具
```bash
cd /Users/zhangdaodong/code/sale/backend
go run ./cmd/diagnose/main.go
```

## 经验总结

### 关键发现
1. **数据库NULL值必须使用指针类型** - Go的sqlx无法将NULL扫描到非指针类型
2. **JWT Claims结构必须统一** - 生成和解析必须使用相同的结构体
3. **角色类型应该用string** - 更灵活，避免硬编码数字映射

### 防止类似问题
1. 使用`*string`、`*int`等指针类型处理可空字段
2. 共享Claims定义，避免多处重复定义
3. 添加详细的错误日志，便于定位问题
4. 统一使用string类型存储枚举值（如角色）

## 相关文件
- 诊断报告: `LOGIN_401_DIAGNOSIS.md`
- 密码重置工具: `PASSWORD_RESET_TOOL.md`
- 集成测试报告: `INTEGRATION_TEST_REPORT.md`

---
**解决时间**: 2026-04-10
**状态**: ✅ 完全解决
**测试状态**: ✅ 全部通过
