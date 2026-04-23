## 1. 后端 Repository — 新增 dashboard 查询方法

- [x] 1.1 在 `RechargeRepoInterface` 新增方法：
  - `GetTodayRechargeTotal(centerID string) (float64, error)` — 当日充值金额合计
  - `GetTodayConsumptionTotal(centerID string) (float64, error)` — 当日核销金额合计
  - `GetActiveCenterCount(centerID string) (int64, error)` — 活跃中心数
  - `GetYesterdayRechargeTotal(centerID string) (float64, error)` — 昨日充值（算 trend）
  - `GetYesterdayConsumptionTotal(centerID string) (float64, error)` — 昨日核销（算 trend）
  - `CountPendingApprovals(centerID string) (int64, error)` — 待审批申请数
  - `CountExpiringCards(centerID string) (int64, error)` — 7天内过期卡数
  - `GetRechargeTrends(days int, centerID string) ([]string, []float64, error)` — 日期+金额趋势

- [x] 1.2 在 `RechargeRepository` 实现上述方法，centerID 为空时不加 center 过滤条件

## 2. 后端 Service — 新增 dashboard 业务方法

- [x] 2.1 在 `RechargeServiceInterface` 新增 3 个方法：
  - `GetDashboardStatistics(role, centerID string) (map[string]interface{}, error)`
  - `GetDashboardTodos(role, centerID string) (map[string]interface{}, error)`
  - `GetDashboardRechargeTrends(days int, role, centerID string) (map[string]interface{}, error)`

- [x] 2.2 在 `RechargeService` 实现上述方法：
  - super_admin/hq_admin/finance 传空 centerID（全局）
  - center_admin/operator 传真实 centerID（本中心）
  - Statistics: 调 repo 聚合，计算 trend 百分比
  - Todos: 调 repo 计数，生成 description
  - Trends: 调 repo 取日期+金额数组

## 3. 后端 Handler — 替换 mock 为 service 调用

- [x] 3.1 修改 `GetDashboardStatistics`：从 JWT 获取 role/centerID，调用 `rechargeService.GetDashboardStatistics`
- [x] 3.2 修改 `GetDashboardTodos`：同上，调用 `rechargeService.GetDashboardTodos`
- [x] 3.3 修改 `GetDashboardRechargeTrends`：同上，调用 `rechargeService.GetDashboardRechargeTrends`

## 4. 前端 — 修正快捷操作名称 + 移除 mock fallback

- [x] 4.1 修改 `Dashboard.vue` 快捷操作名称：
  - "C端充值" → "C端充值录入"
  - "门店卡发放" → "绑定卡号"
  - "充值申请" → "B端充值申请"

- [x] 4.2 移除 `statistics` 的硬编码初始值，改为全 0 / 空字符串
- [x] 4.3 移除 `chartData` 的硬编码初始值，改为空数组

## 5. 验证

- [x] 5.1 API 测试：不同角色调用 3 个 dashboard 接口，验证数据来源和角色过滤
