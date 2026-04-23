## Context

当前门店卡状态机定义了 6 种状态（已入库1/已发放2/已激活3/已冻结4/已过期5/已作废6），通过白名单 `allowedTransitions` 控制转换。作废功能实际业务中不需要，且存在前后端不一致问题（已激活卡前端显示作废按钮但后端拒绝）。

涉及文件：
- 后端：`service/recharge.go`（白名单、VoidCard）、`handler/recharge.go`（VoidCard handler）、`repository/recharge.go`（TransitionCardStatusTX）、`router/router.go`（路由）、`rbac/casbin.go`（权限）
- 前端：`CardManage.vue`（操作按钮）、`CardInventory.vue`（统计）、`api/card.ts`（类型、枚举、API 函数）

## Goals / Non-Goals

**Goals:**
- 简化状态机：移除作废，所有非冻结卡可冻结，冻结卡只能解冻
- 前后端保持一致

**Non-Goals:**
- 不处理数据库中已有的 status=6 历史数据（保留原值）
- 不改变冻结/解冻的底层事务逻辑（TransitionCardStatusTX）
- 不改变卡核销、发放等其他业务流程

## Decisions

1. **白名单简化为统一规则**：不再用 map 逐状态枚举，改为两条规则判断：
   - 目标是冻结(4)：当前非冻结即可
   - 目标是解冻(3)：当前必须是冻结(4)
   - 其他转换：不允许（由业务逻辑自然触发，如过期由 VerifyCard 自动标记）

2. **VoidCard 方法保留但返回错误**：而非删除方法签名（避免 interface 变更），handler 返回 404 让 API 优雅退役

3. **前端完全移除作废 UI**：操作列只保留"冻结"和"解冻"按钮

## Risks / Trade-offs

- [已有 status=6 的卡] → 前端统计区不再展示作废计数，但数据库中历史数据保留不变，筛选"已作废"状态仍可查到
- [API 移除 /card/:cardNo/void] → 如果有外部调用方会报 404，属 BREAKING 变更
