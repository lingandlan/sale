## Context

充值操作员信息存储在 `c_recharges.operator_id` 和 `c_recharges.operator_name`，但写入时使用硬编码 mock 值。需改为从当前请求的 JWT token 解析用户身份，查询数据库获取真实用户名。

## Goals / Non-Goals

**Goals:**

- `CreateCRecharge` 写入真实 `operatorId`（用户表主键）和 `operatorName`（用户表 name 字段）
- 复用已有的 `userRepo.GetByID()` 查询

**Non-Goals:**

- 不修改其他接口（CardTransaction、CardIssueRecord 等仍通过 handler 层传参）
- 不改变 JWT token 结构

## Decisions

**1. 复用 `userRepo` 而非新增 service**

`RechargeHandler` 已注入 `userRepo`，直接调用 `userRepo.GetByID()` 查询用户名，无需新增依赖。

**2. `getOperatorInfo` 替代 `getOperatorCenter`**

原函数仅返回 `userID + role + centerID`，增加第4个返回值 `name`。所有调用方适配新签名。

**3. 降级策略**

若 `userRepo.GetByID()` 查询失败（用户被删等），`operatorName` 降级为 `"未知用户"`，不影响业务写入。

## Risks / Trade-offs

- [最小风险] 纯读取逻辑，不涉及数据写入变更，失败可降级
