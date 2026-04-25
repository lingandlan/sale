## Why

门店卡统计页面存在大量硬编码 mock 数据（饼图、柱状图、充值中心表格），无法展示真实业务数据。同时缺少充值中心维度的数据隔离——总部用户应查看全部数据，中心用户仅查看所属中心。

## What Changes

- 饼图（卡状态分布）：从 `getCardStats` 接口返回的状态计数动态生成，替换硬编码
- 柱状图（月度发放/核销趋势）：新增后端接口，按月 GROUP BY 统计发放和核销卡数量
- 充值中心统计表格：新增后端接口，按 `recharge_center_id` 分组统计各中心卡数和余额
- 所有接口增加充值中心过滤：总部角色看全部，中心角色只看自己中心
- 现有 `GetCardStats` 增加 `centerId` 参数支持按中心过滤

## Capabilities

### New Capabilities
- `card-monthly-trend`: 月度发放/核销趋势统计接口，按月份分组返回发放和核销卡数量
- `card-center-stats`: 充值中心维度的卡统计接口，按中心分组返回各状态的卡数和总余额

### Modified Capabilities
-（无现有 spec）

## Impact

- `backend/internal/repository/recharge.go` — 新增两个查询方法，`GetCardStats` 增加 centerId 过滤
- `backend/internal/service/recharge.go` — 新增两个 service 方法，`GetCardStats` 透传 centerId
- `backend/internal/handler/recharge.go` — 新增两个 handler，`GetCardStats` 从 context 获取用户中心
- `backend/internal/router/router.go` — 注册新路由
- `shop-pc/src/api/card.ts` — 新增两个 API 函数
- `shop-pc/src/views/card/CardStats.vue` — 替换所有 mock 数据为真实接口调用
