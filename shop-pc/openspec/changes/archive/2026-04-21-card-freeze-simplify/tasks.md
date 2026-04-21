## 1. 后端：状态转换白名单简化

- [x] 1.1 修改 `service/recharge.go` 中的 `allowedTransitions` 白名单，改为两条规则判断（非冻结可冻结、冻结可解冻），移除作废相关转换
- [x] 1.2 修改 `transitionCardStatus` 方法，用新的规则判断替换白名单 map 查找逻辑
- [x] 1.3 移除 `service/recharge.go` 中的 `VoidCard` 方法
- [x] 1.4 移除 `handler/recharge.go` 中的 `VoidCard` handler 方法
- [x] 1.5 移除 `service/interfaces.go` 中 `VoidCard` 的 interface 声明
- [x] 1.6 移除 `router/router.go` 中的 `/card/:cardNo/void` 路由
- [x] 1.7 移除 `rbac/casbin.go` 中 void 相关的 RBAC 权限规则

## 2. 前端：操作按钮和统计调整

- [x] 2.1 修改 `CardManage.vue` 操作列：非冻结显示"详情+冻结"，冻结显示"详情+解冻"，移除作废按钮
- [x] 2.2 修改 `CardManage.vue` 统计区：移除"已作废"统计卡片，调整为 6 项（grid 改为 6 列）
- [x] 2.3 修改 `CardInventory.vue` 统计区：移除"已作废"展示，保留 6 项（总卡数+5种状态）
- [x] 2.4 修改 `api/card.ts`：移除 `voidCard` 函数、`CardStatusMap` 中移除 `6: '已作废'`、`CardStatusTagType` 中移除 `6: 'info'`、`CardStatsResponse` 中移除 `voidedCards`

## 3. 验证

- [x] 3.1 编译后端确认无报错
- [x] 3.2 前端页面确认按钮和统计展示正确
- [x] 3.3 测试冻结/解冻流程：对已入库、已发放、已激活、已过期的卡分别测试冻结，冻结后测试解冻
