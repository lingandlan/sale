## Why

充值记录列表页（/recharge/records）和详情页（/recharge/records/:id）存在两个体验问题：1）列表页多展示了业务上不需要的"支付方式"列，干扰信息密度；2）列表页未展示"操作员"字段，无法识别操作来源，而详情页已正确展示操作员姓名。

## What Changes

- **移除**：`RecordList.vue` 表格中的"支付方式"列
- **移除**：`RecordDetail.vue` 详情卡中的"支付方式"行
- **新增**：`RecordList.vue` 表格中新增"操作员"列，展示操作员姓名（后端已返回 `operatorId` + `operatorName`，前端直接引用）
- **接口调整**：`api/recharge.ts` 中 `RechargeRecordItem` 接口移除 `paymentMethod` 字段，补充 `operatorId` + `operatorName`

## Capabilities

### New Capabilities

- `recharge-record-list`: 充值记录列表页，展示交易单号、会员姓名、手机号、充值中心、充值金额、操作员、充值时间，支持分页筛选
- `recharge-record-detail`: 充值记录详情页，展示充值全量信息（不含支付方式）

### Modified Capabilities

- `recharge-api`: `RechargeRecordItem` 接口字段调整（移除 paymentMethod，新增 operatorId/operatorName），后端接口不变

## Impact

- **前端**：`shop-pc/src/views/recharge-record/RecordList.vue`、`RecordDetail.vue`；`shop-pc/src/api/recharge.ts`
- **后端**：无需修改（`GetCRechargeList` / `GetCRechargeDetail` 已返回所需字段）
