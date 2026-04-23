## NEW Requirements

### Requirement: 充值趋势从 DB 按日期聚合真实数据

`GET /dashboard/recharge-trends?days=N` SHALL 从数据库按日期聚合 C端充值数据，不再返回硬编码值。

#### Scenario: 默认 7 天趋势

- **WHEN** 请求 `/dashboard/recharge-trends`（不传 days 参数）
- **THEN** 返回最近 7 天的 C端充值数据：
  - `dates`: 日期数组，格式 "MM-DD"（如 ["04-12", "04-13", ...]）
  - `values`: 每日充值金额总和数组，与 dates 一一对应

#### Scenario: 自定义天数趋势

- **WHEN** 请求 `/dashboard/recharge-trends?days=30`
- **THEN** 返回最近 30 天的聚合数据，格式同上

#### Scenario: center_admin/operator 限定本中心

- **WHEN** center_admin 或 operator 角色请求
- **THEN** 仅聚合本中心（`center_id` 匹配）的 C端充值数据

#### Scenario: 某天无充值记录

- **WHEN** 某天没有任何充值记录
- **THEN** 该天 `value` 为 0，日期仍需出现在 `dates` 数组中（保证连续性）

#### Scenario: 数据库查询失败时降级

- **WHEN** 查询过程中发生数据库错误
- **THEN** 返回 HTTP 500，错误信息 "趋势数据加载失败"
