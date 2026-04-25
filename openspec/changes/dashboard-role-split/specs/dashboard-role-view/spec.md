## ADDED Requirements

### Requirement: 移除总会员数统计卡片

前端 Dashboard 移除"总会员数"（memberCount）统计卡片。后端 `GetDashboardStatistics` 不再返回 `memberCount` 和 `memberTrend` 字段。

#### Scenario: 总部人员访问仪表盘
- **WHEN** super_admin / hq_admin / finance 用户访问 Dashboard
- **THEN** 显示 3 张统计卡片：今日充值金额、今日核销金额、活跃中心数；不显示总会员数

#### Scenario: 充值中心人员访问仪表盘
- **WHEN** center_admin / operator 用户访问 Dashboard
- **THEN** 显示 2 张统计卡片：今日充值金额、今日核销金额；不显示总会员数，也不显示活跃中心数

### Requirement: 充值中心人员数据范围限定

后端 Dashboard API 已有基于 JWT 的角色过滤。center_admin / operator 的数据（充值金额、核销金额、充值趋势）限定为用户所属 centerId 的数据。前端无需额外传参。

#### Scenario: 充值中心人员查看今日充值金额
- **WHEN** center_admin 用户请求 `GET /dashboard/statistics`
- **THEN** `todayRecharge` 仅返回该用户所属中心的今日充值总额

#### Scenario: 充值中心人员查看充值趋势
- **WHEN** center_admin 用户请求 `GET /dashboard/recharge-trends`
- **THEN** 返回的数据仅包含该用户所属中心的趋势数据

### Requirement: 快捷操作按角色区分

#### Scenario: 总部人员快捷操作
- **WHEN** super_admin / hq_admin / finance 用户访问 Dashboard
- **THEN** 显示全部 4 个快捷操作：C端充值录入、门店卡核销、绑定卡号、B端充值申请

#### Scenario: 充值中心人员快捷操作
- **WHEN** center_admin / operator 用户访问 Dashboard
- **THEN** 仅显示 3 个快捷操作：C端充值录入、门店卡核销、绑定卡号

### Requirement: 前端按角色条件渲染

前端通过 userStore 的 `isSuperAdmin` / `isHQAdmin` / `isFinance` / `canSelectAllCenters` 判断角色，条件渲染统计卡片、快捷操作和待办事项。

#### Scenario: 充值中心人员不显示待办事项
- **WHEN** center_admin / operator 用户访问 Dashboard
- **THEN** 不显示"待审批充值申请"待办项（该功能属于总部）
