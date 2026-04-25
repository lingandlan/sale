## Why

返还积分比例前后端不一致：前端按"当前充值金额"判断（≥10万返2%），后端按"上月净消费"判断但前端传的 `lastMonthConsumption` 固定为 0，导致返还比例永远按 1% 计算。需要从后端查询充值中心的真实上月净消费数据，前后端统一规则。

## What Changes

- 后端新增接口：查询指定充值中心的上月净消费金额
- B端申请页面加载时，选中充值中心后自动查询该中心的上月净消费
- 前端根据真实上月净消费计算返还比例并展示
- 后端 `CreateBRechargeApplication` 使用前端传入的真实 `lastMonthConsumption` 计算积分
- 申请记录表新增 `last_month_consumption` 字段存档

## Capabilities

### New Capabilities
- `center-last-month-consumption`: 查询充值中心上月净消费金额，用于计算B端申请返还比例

### Modified Capabilities


## Impact

- 后端：新增 API 接口、新增 DB 字段、修改积分计算逻辑
- 前端：B端申请页调用新接口、修正返还比例展示
- 数据库：recharge_applications 表新增 last_month_consumption 列
