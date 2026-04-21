## Why

后端代码中存在多处可导致服务崩溃（裸类型断言 panic）、数据不一致（充值扣款但积分未到账、状态变更无事务保护）、以及生产环境难以排查问题（recovery 中间件不记录堆栈）的隐患。系统即将进入正式运营阶段，需要在上线前系统性地修复这些稳定性问题。

## What Changes

- 修复所有 `map[string]interface{}` 裸类型断言，替换为带安全检查的解析方式，防止客户端参数缺失/类型错误导致 panic
- 修复 recovery 中间件，恢复 panic 堆栈日志记录，确保生产环境可追溯崩溃原因
- 修复 C 充值流程中 AddIntegral 失败被静默吞掉的问题，确保扣款与积分发放的数据一致性
- 为状态变更 + 审计记录创建操作添加事务保护（卡状态流转、批量创建卡事务记录）
- 修复 admin handler 中 Casbin 角色同步错误被忽略的问题
- 修复全局 logger 无 nil 保护的问题
- 替换 server 启动 goroutine 中的 `log.Fatal`，改为优雅关闭
- 补充 HTTP server `ReadHeaderTimeout`、数据库 `SetConnMaxIdleLifetime`、Redis 超时配置
- 修复 GetCenters / GetOperators 无分页的全表扫描，以及 GetCardList 的 N+1 查询问题
- 清理 repository 层大量被静默忽略的 GORM 错误（Count、Scan）
- 移除 auth handler 中的 `fmt.Printf` 调试日志，改用 zap logger

## Capabilities

### New Capabilities
- `panic-safety`: 裸类型断言修复 + recovery 中间件堆栈日志 + context helper 安全取值，防止运行时 panic 导致服务不可用
- `data-consistency`: 关键业务操作的事务保护——C充值积分发放、卡状态流转审计、批量卡创建事务包装
- `error-handling`: repository 层 GORM 错误处理规范化 + service 层被忽略的 Casbin/Redis 错误处理 + 全局 logger nil 保护

### Modified Capabilities
(无已有 spec 需要修改)

## Impact

- **后端代码**: internal/service/、internal/handler/、internal/repository/、internal/middleware/、cmd/server/、pkg/logger/
- **API 行为**: 部分接口错误响应更明确（之前返回 500 的 panic 改为返回业务错误码）
- **运维**: recovery 日志增强后可观测性提升；数据库连接池和 Redis 配置优化后资源管理更合理
- **无前端变更**: 本次改动仅涉及后端稳定性，不改变 API 契约
