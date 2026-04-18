## Why

Dashboard 页面的统计数据、待办事项、充值趋势三个 API 全部返回硬编码 mock 数据，无法反映真实业务状态。快捷操作名称与侧边栏菜单不一致，待办事项未按角色区分显示内容。

## What Changes

- **替换 mock API**：后端三个 dashboard 接口（statistics、todos、recharge-trends）改为从数据库查询真实数据
- **修正快捷操作名称**：与侧边栏菜单对齐（"C端充值" → "C端充值录入"，"充值申请" → "B端充值申请"，"门店卡发放" → "绑定卡号"）
- **待办事项角色区分**：按用户角色过滤待办项（super_admin/hq_admin/finance 可见待审批申请，center_admin/operator 仅可见本中心待办）

## Capabilities

### New Capabilities

- `dashboard-statistics`: 仪表盘统计数据接口，从 DB 聚合真实数据（会员数、今日充值、今日核销、活跃中心数）
- `dashboard-todos`: 待办事项接口，按角色返回不同待办列表
- `dashboard-recharge-trends`: 充值趋势图表接口，按日期范围聚合真实充值数据

### Modified Capabilities

## Impact

- **后端**：`handler/recharge.go` 中三个 dashboard handler 改为调用 service/repository 查询 DB
- **后端**：需在 `RechargeServiceInterface` 新增 dashboard 查询方法
- **前端**：`Dashboard.vue` 快捷操作名称修改，移除 statistics/chart 的硬编码 fallback 数据
- **数据库**：纯读取，无 schema 变更
