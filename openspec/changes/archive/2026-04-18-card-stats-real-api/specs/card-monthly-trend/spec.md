## ADDED Requirements

### Requirement: 月度发放/核销趋势统计

系统 SHALL 提供按月统计卡发放和核销数量的接口 `GET /card/monthly-trend`，返回最近 6 个月的发放与核销卡数量。

#### Scenario: 总部用户查询趋势
- **WHEN** super_admin/hq_admin 请求 `/card/monthly-trend`（不传 centerId）
- **THEN** 返回最近 6 个月全局的发放和核销统计，按月份升序排列

#### Scenario: 中心用户查询趋势
- **WHEN** center_admin/operator 请求 `/card/monthly-trend`
- **THEN** 自动使用该用户所属充值中心过滤，仅返回该中心的月度数据

#### Scenario: 总部用户按中心过滤
- **WHEN** super_admin/hq_admin 请求 `/card/monthly-trend?centerId=xxx`
- **THEN** 仅返回指定中心的月度数据

#### Scenario: 无数据月份
- **WHEN** 某月无发放或核销记录
- **THEN** 该月对应字段返回 0，月份仍展示

### Requirement: 充值中心维度卡统计

系统 SHALL 提供按充值中心分组的卡统计接口 `GET /card/center-stats`，返回每个中心的各状态卡数和总余额。

#### Scenario: 总部用户查询全部中心
- **WHEN** super_admin/hq_admin 请求 `/card/center-stats`
- **THEN** 返回所有充值中心的统计，每条包含 centerName、totalCards、各状态卡数、totalBalance

#### Scenario: 中心用户查询
- **WHEN** center_admin/operator 请求 `/card/center-stats`
- **THEN** 仅返回该用户所属充值中心的统计

### Requirement: GetCardStats 增加中心过滤

现有 `GET /card/stats` 接口 SHALL 支持 centerId 参数过滤统计范围。

#### Scenario: 总部用户无过滤
- **WHEN** super_admin/hq_admin 请求 `/card/stats`（不传 centerId）
- **THEN** 返回全局统计数据（与现有行为一致）

#### Scenario: 总部用户按中心过滤
- **WHEN** super_admin/hq_admin 请求 `/card/stats?centerId=xxx`
- **THEN** 仅统计指定中心的卡数据

#### Scenario: 中心用户自动过滤
- **WHEN** center_admin/operator 请求 `/card/stats`
- **THEN** 自动使用该用户所属中心过滤，仅返回该中心统计
