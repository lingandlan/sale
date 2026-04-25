## Why

首页仪表盘目前对所有角色展示相同内容，但充值中心人员只需关注本中心数据。同时"总会员数"依赖商城接口无法获取，需要移除。

## What Changes

- 移除"总会员数"统计卡片（数据不可获取）
- 确认并修复"今日充值金额"和"今日核销金额"的数据准确性
- 仪表盘按角色区分展示：
  - **总部人员**（super_admin / hq_admin / finance）：展示全部统计数据、全部快捷操作、全部充值趋势
  - **充值中心人员**（center_admin / operator）：只展示本中心的充值金额和核销金额；快捷操作仅保留 C端充值录入、门店卡核销、绑定卡号；充值趋势只展示本中心数据
- 后端 Dashboard API 需根据用户角色/所属中心过滤数据

## Capabilities

### New Capabilities

- `dashboard-role-view`: 仪表盘按角色区分数据范围和可见操作项

### Modified Capabilities

无

## Impact

- **后端**: `service/recharge.go` GetDashboardStatistics / GetDashboardRechargeTrends 需接收用户角色和 centerId 参数，按角色过滤数据
- **前端**: `Dashboard.vue` 需根据用户角色条件渲染统计卡片、快捷操作和趋势图表
- **API**: Dashboard 相关接口可能需要调整请求参数或响应结构
