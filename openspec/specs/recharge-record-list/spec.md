# recharge-record-list Specification

## Purpose
TBD - created by archiving change fix-recharge-records-list. Update Purpose after archive.
## Requirements
### Requirement: 充值记录列表展示

充值记录列表页 SHALL 支持展示充值记录列表，包含交易单号、会员姓名、手机号、充值中心、充值金额、操作员姓名、充值时间，支持分页和筛选。

#### Scenario: 列表页默认展示

- **WHEN** 用户访问 /recharge/records
- **THEN** 系统展示充值记录列表，每行包含：交易单号、会员姓名、手机号、充值中心、充值金额、操作员姓名、充值时间、详情操作按钮

#### Scenario: 列表页支持按手机号筛选

- **WHEN** 用户输入手机号并点击查询
- **THEN** 系统筛选并展示匹配的充值记录

#### Scenario: 列表页支持按充值中心筛选

- **WHEN** 用户选择充值中心并点击查询
- **THEN** 系统筛选并展示该中心的充值记录

#### Scenario: 列表页支持按日期范围筛选

- **WHEN** 用户选择日期范围并点击查询
- **THEN** 系统筛选并展示该日期范围内的充值记录

#### Scenario: 点击详情跳转到详情页

- **WHEN** 用户点击某行的"详情"按钮
- **THEN** 系统导航到 /recharge/records/{id} 详情页

