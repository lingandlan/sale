## Why

当前门店卡状态机存在"作废"状态和对应的转换逻辑，但实际业务中不需要作废功能。现有白名单规则导致已激活(3)和已过期(5)的卡无法被冻结，且前端"作废"按钮对已激活卡可见但后端拒绝操作，前后端不一致。需要简化状态机：任何状态的卡都可以冻结，冻结后只能解冻，移除作废相关逻辑。

## What Changes

- **BREAKING**: 移除"作废"状态(status=6)的所有逻辑，包括前端按钮、后端转换白名单、API 路由
- 允许所有非冻结状态（已入库、已发放、已激活、已过期）的卡执行冻结操作
- 冻结后的卡只能执行解冻，不能执行其他任何操作（核销、发放等）
- 前端 CardManage.vue 操作列：所有非冻结卡显示"冻结"按钮，冻结卡显示"解冻"按钮，移除"作废"按钮
- 前端 CardInventory.vue 统计区移除"已作废"计数
- 后端状态转换白名单简化为：任何状态 → 可冻结(4)；冻结(4) → 可解冻(3)
- 卡状态枚举从 6 种减少为 5 种（已入库/已发放/已激活/已冻结/已过期）

## Capabilities

### New Capabilities

（无新能力）

### Modified Capabilities

- `card-status-machine`: 卡状态转换规则简化 — 移除作废，所有状态可冻结，冻结后仅可解冻

## Impact

- **后端**: `service/recharge.go` 状态转换白名单、`VoidCard` 方法移除、`VerifyCard` 状态检查、handler 路由和 RBAC 规则
- **前端**: `CardManage.vue` 操作按钮、`CardInventory.vue` 统计展示、`api/card.ts` 类型定义和状态枚举
- **数据库**: 已有 status=6 的历史数据需要考虑（保持原值，不再新增）
- **API**: 移除 `/card/:cardNo/void` 端点，**BREAKING** 变更
