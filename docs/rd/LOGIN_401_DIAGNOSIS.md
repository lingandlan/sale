# 登录401问题诊断报告

## 问题现状
后端登录接口 `/api/v1/auth/login` 持续返回 401 "手机号或密码错误"

## 已验证的正常项 ✅

### 1. 数据库层面（完全正常）
```
✅ MySQL容器运行正常（端口3306）
✅ Redis容器运行正常（端口6379）
✅ 用户表数据存在
✅ 用户状态: 1 (正常)
✅ 密码hash正确 (bcrypt验证通过)
✅ deleted_at字段: NULL (未软删除)
```

### 2. 测试账号信息
```
手机号: 13800000000
密码: Test123456
角色: super_admin
```

### 3. 数据库查询测试
使用诊断工具（`cmd/diagnose/main.go`）验证：
- 数据库连接成功
- 用户查询成功
- 密码验证成功
- 用户状态检查通过

### 4. 后端服务状态
```
✅ 后端进程运行中（端口8080）
✅ 所有路由注册成功
✅ 中间件加载成功（Recovery、ZapLogger、CORS）
✅ HTTP请求可达（curl收到401响应）
```

## 问题定位 ❌

### 可能的原因

根据诊断结果，问题**不在数据库层面**，而在**后端代码的业务逻辑中**：

1. **后端连接的数据库配置不对**
   - 后端可能连接到了其他数据库实例
   - 配置文件路径或加载有问题

2. **Repository层的查询逻辑问题**
   - sqlx查询与测试工具的行为不一致
   - 可能存在字段映射问题

3. **错误处理逻辑异常**
   - 某个环节的错误被误判为密码错误
   - 错误类型判断有问题

4. **日志输出缺失**
   - 添加的调试日志没有输出
   - 日志被重定向或未正确初始化

### 已排除的原因

```
❌ 数据库服务未启动 → 已验证运行正常
❌ 用户不存在 → 已验证存在
❌ 密码hash错误 → bcrypt验证通过
❌ 用户状态异常 → status=1 (正常)
❌ 用户被软删除 → deleted_at=NULL
❌ 网络连接问题 → curl可达端口8080
```

## 当前状态

### 可以正常工作的部分
- ✅ 数据库连接和查询
- ✅ 密码重置工具（`cmd/resetpwd/main.go`）
- ✅ 诊断工具（`cmd/diagnose/main.go`）
- ✅ 后端服务启动和路由注册
- ✅ HTTP请求可达后端

### 卡住的部分
- ❌ 后端登录接口返回401
- ❌ 调试日志未输出
- ❌ 无法确定具体是哪个环节出错

## 建议的调试方向

### 方案1: 检查后端日志输出
```bash
# 查看后端实时日志
tail -f /tmp/server.log

# 启动后端并查看详细日志
go run ./cmd/server/main.go 2>&1 | tee /tmp/debug.log
```

### 方案2: 使用dlv调试器
```bash
# 安装dlv
go install github.com/go-delve/delve/cmd/dlv@latest

# 调试后端
dlv debug ./cmd/server/main.go
```

### 方案3: 添加HTTP日志中间件
```go
// 在main.go中添加
r.Use(gin.Logger())
```

### 方案4: 直接测试Repository层
```bash
# 运行Repository测试（如果有）
go test ./internal/repository/ -v
```

### 方案5: 检查配置加载
```bash
# 确认配置文件路径
ls -la configs/config.yaml

# 测试配置加载
go run -mod=mod ./cmd/testconfig/main.go
```

## 工具使用

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

### 数据库测试
```bash
cd /Users/zhangdaodong/code/sale/backend
go run ./cmd/testdb/main.go
```

## 已创建的工具文件

1. `cmd/resetpwd/main.go` - 密码重置工具
2. `cmd/diagnose/main.go` - 完整诊断工具
3. `cmd/testdb/main.go` - 数据库连接测试
4. `cmd/genpass/main.go` - 密码hash生成
5. `cmd/testlogin/main.go` - 登录流程测试
6. `PASSWORD_RESET_TOOL.md` - 工具使用文档

## 总结

**问题卡在**: 数据库验证通过，但后端业务逻辑返回401

**证据**: 诊断工具显示数据库层面一切正常

**下一步**: 需要深入调试后端代码，确定具体是哪个环节出错

**优先级建议**:
1. 启用Gin的Logger中间件，查看HTTP请求日志
2. 在Repository、Service、Handler层添加详细日志
3. 使用dlv调试器单步跟踪登录流程
4. 检查配置文件的加载路径和内容

---
**生成时间**: 2026-04-10
**诊断工具**: `/Users/zhangdaodong/code/sale/backend/cmd/diagnose/main.go`
