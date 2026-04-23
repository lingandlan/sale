## Context

Dashboard 页面现有 3 个 API 端点（`/dashboard/statistics`、`/dashboard/todos`、`/dashboard/recharge-trends`），全部返回硬编码 mock 数据。前端有 fallback 硬编码值作为初始状态。后端 handler 直接返回 `gin.H{}`，未经过 service/repository 层。

相关数据表：
- `c_recharges` — C端充值记录（amount, created_at, center_id）
- `recharge_applications` — B端充值申请（status: pending/approved/rejected）
- `store_cards` — 门店卡（status, expired_at, recharge_center_id）
- `recharge_centers` — 充值中心
- `users` — 用户表（role, center_id）

## Goals / Non-Goals

**Goals:**

- 三个 dashboard API 从 DB 查询真实数据
- 待办事项按用户角色过滤（super_admin/hq_admin/finance 看全部，center_admin/operator 仅看本中心）
- 前端移除硬编码 fallback，API 失败时显示 0/空
- 快捷操作名称与侧边栏菜单一致

**Non-Goals:**

- 不新增前端页面或路由
- 不修改 GORM model 或新增数据表
- 不做前端图表组件重构
- 不做缓存优化（数据量小，直接查 DB）

## Decisions

**1. 复用现有 repository，不新建 DashboardRepository**

所需查询（统计金额、计数、按日期聚合）都是对已有表的简单聚合，直接在 `RechargeRepoInterface` 上新增 dashboard 查询方法即可。避免引入新的依赖注入。

**2. Service 层新增 3 个 dashboard 方法到 RechargeServiceInterface**

- `GetDashboardStatistics(userID int64, role string, centerID string) map[string]interface{}`
- `GetDashboardTodos(userID int64, role string, centerID string) map[string]interface{}`
- `GetDashboardRechargeTrends(days int, userID int64, role string, centerID string) map[string]interface{}`

统一传 `userID/role/centerID`，在 service 层做角色判断和数据过滤。

**3. 待办事项角色策略**

| 角色 | 可见待办 |
|------|---------|
| super_admin / hq_admin / finance | 全部待审批申请 + 全部即将过期卡 |
| center_admin | 本中心待审批 + 本中心即将过期卡 |
| operator | 本中心待审批（可见但不审批）+ 本中心即将过期卡 |

**4. 统计数据定义**

- **会员数（memberCount）**：暂不实现（需要 WSY 外部接口），返回 0
- **今日充值金额（todayRecharge）**：`SUM(c_recharges.amount) WHERE DATE(created_at) = CURDATE()`
- **今日核销金额（todayConsumption）**：`SUM(card_transactions.amount) WHERE type='consume' AND DATE(created_at) = CURDATE()`
- **活跃中心数（activeCenters）**：`COUNT(DISTINCT center_id) FROM c_recharges WHERE DATE(created_at) >= DATE_SUB(CURDATE(), INTERVAL 30 DAY)`
- **趋势值（trend）**：对比昨日/上周同日计算百分比，简化为固定格式 "+N%" 或 "-N%"

**5. 快捷操作名称修正**

| 原名称 | 修正为 | 对应菜单 |
|--------|--------|---------|
| C端充值 | C端充值录入 | C端充值录入 |
| 门店卡核销 | 门店卡核销 | 门店卡核销 |
| 门店卡发放 | 绑定卡号 | 绑定卡号 |
| 充值申请 | B端充值申请 | B端充值申请 |

## Risks / Trade-offs

- **[低风险] 统计查询性能**：当前数据量小，无需缓存。后续数据量增长时可加 Redis 缓存。
- **[低风险] 会员数字段**：暂返回 0，待 WSY 接口集成后补充。前端需展示合理 fallback（如 "—"）。
- **[低风险] 前端去掉 mock fallback**：API 失败时显示 0，用户体验略有下降但数据可信。
