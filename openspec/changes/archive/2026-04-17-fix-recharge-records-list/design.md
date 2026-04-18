## Context

充值记录列表页和详情页的字段展示需要调整，属于纯前端改动，不涉及后端接口变更。后端 `GetCRechargeList` / `GetCRechargeDetail` 已返回 `operatorId` + `operatorName`，无需修改。

## Goals / Non-Goals

**Goals:**

- 列表页移除"支付方式"列，新增"操作员"列
- 详情页移除"支付方式"行
- TypeScript 接口与实际返回字段对齐

**Non-Goals:**

- 不修改后端接口或数据库结构
- 不涉及权限逻辑变更

## Decisions

**1. 操作员姓名来源**

直接使用后端返回的 `operatorName`（已通过 `JOIN` 或 service 层注入），不在前端单独查询用户信息。理由：后端已处理联表，操作员姓名随列表接口一并返回，无需额外请求。

**2. 接口字段清理**

`RechargeRecordItem` 移除 `paymentMethod`，补充 `operatorId` + `operatorName` 与后端 `model.CRecharge` 保持一致。

## Risks / Trade-offs

- [最小风险] 纯展示层改动，不涉及数据写入或权限逻辑
