## Context

- **页面**：`shop-pc/src/views/card/CardStats.vue`
- **现有接口**：`GET /card/stats` 返回全局统计（totalCards、各状态计数、totalBalance、todayConsume、expireIn7Days），无中心过滤
- **问题**：饼图/柱状图/中心表格全部硬编码 mock 数据；接口不区分用户所属充值中心
- **权限模型**：super_admin/hq_admin/finance 看全部，center_admin/operator 只看自己中心

## Goals / Non-Goals

**Goals:**
- 饼图从 `getCardStats` 返回的各状态计数动态生成
- 新增月度趋势接口（按月统计发放/核销数量）
- 新增充值中心维度统计接口（按中心分组）
- 所有统计接口支持中心过滤
- 前端完全去除 mock 数据

**Non-Goals:**
- 不修改卡业务流程
- 不新增前端页面或路由
- 不引入新的第三方图表库（已使用 ECharts）

## Decisions

1. **在现有 `GetCardStats` 上增加 centerId 参数**
   - handler 层通过 `getOperatorInfo` 获取用户中心，中心角色强制使用自己的 centerId
   - 总部角色支持 query param `centerId` 可选过滤
   - repository 层 SQL WHERE 条件增加 `recharge_center_id = ?`

2. **月度趋势用 `card_transactions` 表聚合**
   - 按 `DATE_FORMAT(created_at, '%Y-%m')` 分组
   - type IN ('issue', 'consume') 分别统计
   - 返回最近 6 个月数据
   - 同样支持 centerId 过滤（JOIN store_cards 表）

3. **充值中心统计用 `store_cards` 表聚合**
   - 按 `recharge_center_id` 分组，各状态计数 + SUM(balance)
   - JOIN recharge_centers 获取中心名称
   - 总部角色返回所有中心，中心角色只返回自己的

4. **前端 ECharts 用 computed reactive**
   - 饼图数据从 stats 响应的各状态字段映射
   - 柱状图数据从新接口响应直接赋值

## Risks / Trade-offs

- [Risk] 月度趋势依赖 `card_transactions` 表数据完整性 → 已有 issue/consume 记录，无需迁移
- [Risk] 大量卡数据时聚合查询性能 → 当前数据量小（千级），暂无需优化
