## 1. 后端 - Repository 层

- [x] 1.1 `GetCardStats` 增加 `centerId` 参数，SQL 增加 `WHERE recharge_center_id = ?` 条件（空则不过滤）
- [x] 1.2 新增 `GetMonthlyTrend(centerId string)` — 按 `DATE_FORMAT(created_at, '%Y-%m')` 分组统计 card_transactions 中 issue/consume 数量，返回最近 6 个月
- [x] 1.3 新增 `GetCenterCardStats(centerId string)` — 按 `recharge_center_id` 分组统计 store_cards 各状态计数 + SUM(balance)，JOIN recharge_centers 获取中心名称

## 2. 后端 - Service 层

- [x] 2.1 `GetCardStats` 签名增加 `centerId string` 参数，透传到 repository
- [x] 2.2 新增 `GetMonthlyTrend(centerId string)` 方法
- [x] 2.3 新增 `GetCenterCardStats(centerID string)` 方法
- [x] 2.4 更新 `interfaces.go` 中对应接口定义

## 3. 后端 - Handler + Router

- [x] 3.1 `GetCardStats` handler 通过 `getOperatorInfo` 获取用户中心，总部支持 query param `centerId`
- [x] 3.2 新增 `GetMonthlyTrend` handler，同样处理中心过滤逻辑
- [x] 3.3 新增 `GetCenterCardStats` handler，同样处理中心过滤逻辑
- [x] 3.4 注册新路由 `GET /card/monthly-trend` 和 `GET /card/center-stats`

## 4. 前端 - API 层

- [x] 4.1 新增 `getMonthlyTrend(centerId?: string)` API 函数
- [x] 4.2 新增 `getCenterCardStats()` API 函数（中心过滤由后端自动处理）
- [x] 4.3 `getCardStats` 增加 `centerId` 可选参数

## 5. 前端 - CardStats.vue 页面

- [x] 5.1 饼图：从 `getCardStats` 返回的各状态计数动态生成 series data，替换硬编码
- [x] 5.2 柱状图：调用 `getMonthlyTrend` 接口获取真实月度数据，替换硬编码
- [x] 5.3 中心表格：调用 `getCenterCardStats` 获取真实数据，替换硬编码
- [x] 5.4 移除所有 mock 数据（pieOption 默认值、barOption 默认值、centerStats 默认值）

## 6. 验证

- [x] 6.1 用 agent-browser 打开 http://localhost:5175/card/stats 验证页面正常加载
- [x] 6.2 验证总部账号看到全部中心数据，中心账号只看自己中心
