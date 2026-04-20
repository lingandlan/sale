## Context

后端 Go 服务（Gin + GORM + Redis）当前存在三类稳定性隐患：
1. **Panic 风险**：service 层通过 `map[string]interface{}` 接收参数，使用裸类型断言 `data["key"].(type)` 访问，客户端缺参或类型错误直接触发 panic
2. **数据不一致**：C 充值流程"创建记录→扣中心余额→加会员积分"无事务保护；卡状态变更与审计记录分开写入
3. **可观测性缺失**：recovery 中间件捕获了 panic 但不记录堆栈日志，生产环境 crash 无从排查

涉及文件：`internal/service/recharge.go`、`internal/middleware/recovery.go`、`internal/middleware/rbac.go`、`internal/middleware/auth.go`、`internal/handler/admin.go`、`internal/repository/recharge.go`、`internal/repository/gorm.go`、`internal/repository/redis.go`、`cmd/server/main.go`、`pkg/logger/logger.go`

## Goals / Non-Goals

**Goals:**
- 消除所有裸类型断言引发的 panic，改为返回业务错误码
- Recovery 中间件记录 panic 堆栈到 zap logger
- C 充值核心流程（创建记录 + 扣中心余额 + 加会员积分）具备事务一致性
- 卡状态变更与审计记录在同一事务中完成
- Repository 层 GORM 错误不再被静默忽略
- 全局 logger 增加 nil 保护
- Server 启动/关闭流程健壮化

**Non-Goals:**
- 不重构 handler 层入参方式（仍使用 `map[string]interface{}`，仅在解析时增加安全检查）
- 不新增 request timeout 中间件（可作为后续优化）
- 不修改 API 契约，前端无需变更
- 不解决安全问题（CORS、明文密码、密码验证等属于安全域，另立 change）
- 不优化 N+1 查询和全表扫描（属于性能域，另立 change）

## Decisions

### D1: map 参数安全解析 — helper 函数

**决策**：在 `internal/service/recharge.go` 中添加 `getFloat64`、`getString`、`getInt` 三个 helper 函数，用于从 `map[string]interface{}` 安全取值。取值失败时返回 `errno` 业务错误而非 panic。

```go
func getFloat64(data map[string]interface{}, key string) (float64, error)
func getString(data map[string]interface{}, key string) (string, error)
func getInt(data map[string]interface{}, key string) (int, error)
```

**替代方案**：
- 定义 typed request struct → 更优雅但改动量大，需同步改 handler 层，当前不涉及
- 使用 `cast` 库（spf13/cast）→ 引入新依赖，且默认零值会掩盖缺失字段问题

**选择理由**：helper 函数改动最小，仅替换断言语法，不改变数据流。返回明确错误码，前端可直接展示。

### D2: Recovery 中间件恢复日志记录

**决策**：取消注释 `logger.Error(...)` 调用，记录 panic 错误信息和完整堆栈。移除 `_ = stack`。

**注意**：`pkg/logger` 包的 `log` 变量需先完成 nil 保护（D5），否则 recovery 中调用 logger 可能触发二次 panic。

### D3: C 充值事务保护 — 补偿方案

**决策**：C 充值流程（CreateCRecharge → DeductCenterBalance → AddIntegral）**不使用数据库事务包裹外部 API 调用**（WSY AddIntegral 是 HTTP 接口，无法参与 DB 事务）。改为以下补偿策略：

1. 创建 CRecharge 记录（status = "pending"）
2. 扣减中心余额（数据库操作）
3. 调用 WSY AddIntegral
4. 若成功 → 更新 CRecharge status 为 "success"，记录 balance_after
5. 若失败 → 记录错误日志，CRecharge status 保持 "pending"，返回业务错误（不回滚扣款，由运营手动处理或定时任务补偿）

**替代方案**：
- DB 事务包裹全部操作 → 不可行，外部 API 无法 rollback
- Saga 模式 → 过度设计，当前只有一个外部调用
- 不处理（现状）→ 已导致数据不一致

**选择理由**：补偿方案简单可靠，通过 status 字段标记中间状态，支持后续自动化补偿。不改变外部 API 调用链路。

### D4: 卡状态变更事务

**决策**：在 repository 层新增 `TransitionCardStatusTX` 方法，使用 GORM 事务在同一 TX 中完成：
1. UPDATE card SET status = ? WHERE card_no = ?
2. INSERT INTO card_transaction (...)

修改 `service/recharge.go` 的 `transitionCardStatus` 调用此新方法。

### D5: Logger nil 保护

**决策**：`GetLogger()` 和所有包级函数（Info/Error/Warn/Debug/Fatal）增加 nil 检查。若 `log == nil`，`GetLogger()` 返回 `zap.NewNop()`，包级函数静略。

```go
func GetLogger() *zap.Logger {
    if log == nil {
        return zap.NewNop()
    }
    return log
}
```

**替代方案**：`sync.Once` + 初始化检查 → 过度设计，当前 Init 在 main 最早期调用

### D6: Server 启动/关闭健壮化

**决策**：
- 将 `log.Fatal` 替换为 `log.Error` + 向 `quit` channel 发送信号，触发优雅关闭流程
- `gormDB.DB()` 返回值增加 nil 检查
- HTTP Server 增加 `ReadHeaderTimeout: 5 * time.Second`
- `repository/gorm.go` 增加 `SetConnMaxIdleLifetime(30 * time.Minute)`
- `repository/redis.go` 增加 `MinIdleConns: 5`、`DialTimeout`、`ReadTimeout`、`WriteTimeout`

### D7: Context helper 安全取值

**决策**：`middleware/auth.go` 中 `GetUserID`、`GetPhone`、`GetRole` 三个函数改为安全断言（comma-ok 模式），取值失败时返回零值而非 panic。调用方（handler）已有 auth 中间件前置校验，所以理论上不会走到零值分支，但防御性编程更安全。

同样修复 `middleware/rbac.go` 中 `role.(string)` 的裸断言。

### D8: Repository 层错误处理

**决策**：以下被静默忽略的错误必须修复：
- `recharge.go` 中 `.Count(&total)` 的返回值必须检查
- `.Scan(&result)` 的返回值必须检查
- `GetCenterTotalRecharge` / `GetCenterTotalConsumed` 等函数签名改为返回 error
- `gorm.go` 中 `SET NAMES utf8mb4` 的 Exec 错误需检查

**注意**：函数签名变更会影响 service 层调用方，需同步修改。

### D9: Admin handler Casbin 错误处理

**决策**：`handler/admin.go` 中 `_ = h.casbinSvc.AddRoleForUser(...)` 改为检查错误。若 Casbin 同步失败，回滚用户创建（返回错误），避免用户存在但无权限的不一致状态。

### D10: 调试日志清理

**决策**：`handler/auth.go` 中 `fmt.Printf("登录失败: phone=%s, ...")` 替换为 `log.Warn("login failed", zap.String("phone", req.Phone), zap.Error(err))`。

## Risks / Trade-offs

**[风险] C 充值补偿方案不自动回滚** → Mitigation：通过 CRecharge.status = "pending" 标记失败记录，后续可加定时任务自动重试 AddIntegral。当前由运营在后台人工处理。

**[风险] Repository 函数签名变更影响面大** → Mitigation：逐函数修改，每个函数的调用方通常只有 1-2 处，改动可控。编译器会捕获所有遗漏。

**[风险] Recovery 中间件调用 logger 可能循环 panic** → Mitigation：D5 的 logger nil 保护确保 GetLogger() 永远返回有效实例。Recovery 中间件用 recover + 独立 zap.NewNop() 兜底。

**[权衡] 不引入 typed struct 而用 helper 函数** → 牺牲了类型安全性，换取最小改动范围。未来如有大规模重构可统一改为 struct binding。

## Migration Plan

1. 所有改动均在一次合并中完成，无数据库 schema 变更
2. 无需修改前端代码
3. 部署后验证：登录、B端充值申请、C端充值、卡操作、Dashboard 等核心流程
4. 回滚策略：git revert 即可，无数据迁移需要回退
