## NEW Requirements

### Requirement: 仪表盘统计数据从 DB 查询真实数据

`GET /dashboard/statistics` SHALL 从数据库聚合真实统计数据，不再返回硬编码值。

#### Scenario: super_admin/hq_admin 查看全局统计

- **WHEN** super_admin 或 hq_admin 角色请求 `/dashboard/statistics`
- **THEN** 返回全局聚合数据：
  - `todayRecharge`: 当日所有 C端充值金额之和
  - `todayConsumption`: 当日所有门店卡核销金额之和
  - `activeCenters`: 近 30 天有充值记录的充值中心数量
  - `memberCount`: 0（待 WSY 接口集成）
- **AND** 每个 `trend` 字段为对比前一天的增减百分比（如 "+5%"、"-3%"），无前一日数据时为 "—"

#### Scenario: center_admin/operator 查看本中心统计

- **WHEN** center_admin 或 operator 角色请求 `/dashboard/statistics`
- **THEN** 返回限定到本充值中心的聚合数据：
  - `todayRecharge`: 仅本中心当日充值金额
  - `todayConsumption`: 仅本中心当日核销金额
  - `activeCenters`: 1（本中心）
  - `memberCount`: 0
- **AND** trend 字段同全局逻辑，但限定本中心数据

#### Scenario: 数据库查询失败时降级

- **WHEN** 统计查询过程中发生数据库错误
- **THEN** 返回 HTTP 500，错误信息 "统计数据加载失败"
