## 1. Logger nil 保护（前置依赖，其他修复依赖此模块）

- [x] 1.1 `pkg/logger/logger.go`: `GetLogger()` 增加 nil 检查，返回 `zap.NewNop()`
- [x] 1.2 `pkg/logger/logger.go`: 所有包级函数（Info/Error/Warn/Debug/Fatal）增加 nil 检查

## 2. Recovery 中间件恢复堆栈日志

- [x] 2.1 `internal/middleware/recovery.go`: 取消注释 `logger.Error(...)` 调用，记录 panic 错误 + 完整堆栈
- [x] 2.2 `internal/middleware/recovery.go`: 移除 `_ = stack` 行

## 3. Map 参数安全解析（panic-safety 核心）

- [x] 3.1 `internal/service/recharge.go`: 添加 `getFloat64`/`getString` helper 函数
- [x] 3.2 `internal/service/recharge.go`: 替换 `CreateBRechargeApplication` 中所有裸断言为 helper 调用
- [x] 3.3 `internal/service/recharge.go`: 替换 `CreateCRecharge` 中所有裸断言为 helper 调用
- [x] 3.4 `internal/service/recharge.go`: 替换 `CreateCenter` 中所有裸断言为 helper 调用
- [x] 3.5 `internal/service/recharge.go`: 替换 `CreateOperator` 中所有裸断言为 helper 调用

## 4. Context / RBAC 安全取值

- [x] 4.1 `internal/middleware/auth.go`: `GetUserID`/`GetPhone`/`GetRole` 改为 comma-ok 模式
- [x] 4.2 `internal/middleware/auth.go`: `RequireRole` 中 `roleVal.(string)` 改为 comma-ok 模式
- [x] 4.3 `internal/middleware/rbac.go`: `Auth()` 中 `role.(string)` 改为 comma-ok 模式
- [x] 4.4 `internal/middleware/rbac.go`: `RequireRoles` 中 `userRole.(string)` 改为 comma-ok 模式

## 5. C 充值补偿流程（data-consistency 核心）

- [x] 5.1 确认 `model.CRecharge` 是否已有 `status` 字段，若无则添加 + 写 migration SQL
- [x] 5.2 `internal/service/recharge.go` `CreateCRecharge`: 创建记录时设置 status = "pending"
- [x] 5.3 `internal/service/recharge.go` `CreateCRecharge`: AddIntegral 成功后更新 status = "success"
- [x] 5.4 `internal/service/recharge.go` `CreateCRecharge`: AddIntegral 失败时记录错误日志，保持 status = "pending"，返回业务错误

## 6. 卡状态变更事务

- [x] 6.1 `internal/repository/recharge.go`: 新增 `TransitionCardStatusTX(cardNo string, updates map[string]interface{}, txn *model.CardTransaction) error` 方法，使用 GORM 事务
- [x] 6.2 `internal/service/recharge.go`: `transitionCardStatus` 调用新的 `TransitionCardStatusTX`

## 7. 批量卡创建事务

- [x] 7.1 `internal/service/recharge.go`: 将 BatchCreateCards 之后的卡交易记录创建循环包裹在 GORM 事务中

## 8. Repository 层 GORM 错误处理

- [x] 8.1 `internal/repository/recharge.go`: `GetRechargeApplications` / `GetCRechargeList` / `GetCardList` / `GetCardTransactions` 中 `.Count(&total)` 错误检查
- [x] 8.2 `internal/repository/recharge.go`: `GetCardStats` / `GetMonthlyTrend` / `GetCenterCardStats` 中 `.Scan()` 错误检查
- [x] 8.3 `internal/repository/recharge.go`: `GetCenterTotalRecharge` 签名改为返回 `(int64, error)`，同步修改 service 调用方
- [x] 8.4 `internal/repository/recharge.go`: `GetCenterTotalConsumed` 签名改为返回 `(float64, error)`，同步修改 service 调用方
- [x] 8.5 `internal/repository/recharge.go`: `DeductCenterBalance` 中 `First` 错误检查
- [x] 8.6 `internal/repository/gorm.go`: `SET NAMES utf8mb4` Exec 错误检查并记录日志

## 9. Casbin 角色同步错误处理

- [x] 9.1 `internal/handler/admin.go`: `CreateUser` 中 `AddRoleForUser` 错误不再忽略，失败时返回错误响应
- [x] 9.2 `internal/handler/admin.go`: `UpdateUser` 中 `UpdateUserRole` 错误不再忽略，失败时返回错误响应

## 10. Server 启动/关闭健壮化

- [x] 10.1 `cmd/server/main.go`: 替换 ListenAndServe goroutine 中 `log.Fatal` 为 `log.Error` + 发送关闭信号
- [x] 10.2 `cmd/server/main.go`: 添加 `ReadHeaderTimeout: 5 * time.Second`
- [x] 10.3 `cmd/server/main.go`: `gormDB.DB()` 返回值 nil 检查后再调用 `.Close()`

## 11. 连接池配置优化

- [x] 11.1 `internal/repository/gorm.go`: 添加 `SetConnMaxIdleTime(30 * time.Minute)`
- [x] 11.2 `internal/repository/redis.go`: 配置 `MinIdleConns: 5`、`DialTimeout: 5s`、`ReadTimeout: 3s`、`WriteTimeout: 3s`，`PoolSize` 改为 50

## 12. 调试日志清理

- [x] 12.1 `internal/handler/auth.go`: 替换 `fmt.Printf` 为 `logger.Warn("login failed", ...)`
