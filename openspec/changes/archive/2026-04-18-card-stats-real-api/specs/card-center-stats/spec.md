## ADDED Requirements

### Requirement: 充值中心统计表格数据

前端 CardStats 页面的充值中心统计表格 SHALL 从 `GET /card/center-stats` 接口获取真实数据，替换硬编码 mock 数据。

#### Scenario: 总部用户查看中心表格
- **WHEN** 总部用户打开门店卡统计页面
- **THEN** 表格展示所有充值中心的统计数据，包含 centerName、totalCards、issuedCards、activeCards、frozenCards、expiredCards、totalBalance

#### Scenario: 中心用户查看中心表格
- **WHEN** 中心用户打开门店卡统计页面
- **THEN** 表格仅展示该用户所属充值中心的数据

### Requirement: 饼图动态数据

前端卡状态分布饼图 SHALL 从 `getCardStats` 接口返回的各状态计数动态生成，替换硬编码数据。

#### Scenario: 有统计数据
- **WHEN** getCardStats 返回 inStockCards、issuedCards、activeCards、frozenCards、expiredCards 各值
- **THEN** 饼图展示各状态的实际数量和占比

#### Scenario: 全部为 0
- **WHEN** getCardStats 返回所有状态计数为 0
- **THEN** 饼图显示"暂无数据"或空状态

### Requirement: 月度趋势柱状图动态数据

前端月度趋势柱状图 SHALL 从 `GET /card/monthly-trend` 接口获取真实数据，替换硬编码数据。

#### Scenario: 有趋势数据
- **WHEN** monthly-trend 接口返回最近 6 个月数据
- **THEN** 柱状图展示对应月份的发放和核销数量

#### Scenario: 无趋势数据
- **WHEN** monthly-trend 接口返回空数组
- **THEN** 柱状图展示最近 6 个月，发放和核销均为 0
