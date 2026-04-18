## ADDED Requirements

### Requirement: 充值记录详情页展示

充值记录详情页 SHALL 展示单条充值记录的完整信息，包含交易单号、会员姓名、手机号、充值金额、获得积分、充值前余额、充值后余额、充值中心、操作员姓名、充值时间、备注。

#### Scenario: 详情页正确展示充值信息

- **WHEN** 用户访问 /recharge/records/{id}
- **THEN** 系统展示充值记录的完整信息（不含支付方式字段）

#### Scenario: 详情页展示操作员姓名

- **WHEN** 用户访问 /recharge/records/{id}
- **THEN** 系统展示操作员姓名（来自 operatorName 字段）

#### Scenario: 点击返回回到列表页

- **WHEN** 用户点击"返回"按钮
- **THEN** 系统导航回 /recharge/records 列表页
