## Why

`CreateCRecharge`（C端充值录入）接口中，操作员信息硬编码为 mock 值（`operatorId: "op123"`, `operatorName: "张出纳"`），未从 JWT 获取当前登录用户的真实身份，导致充值记录中操作员字段不可信。

## What Changes

- **新增**：`getOperatorInfo()` 替代 `getOperatorCenter()`，从 JWT 获取 `userID`，并通过 `userRepo.GetByID()` 查询数据库获取用户真实 `Name`
- **修改**：`CreateCRecharge` 从 `getOperatorInfo()` 获取真实 `operatorId` + `operatorName`，替换硬编码 mock 值
- **适配**：所有调用方适配 `getOperatorInfo()` 的新返回值（增加第4个返回值 `name`）
- **后端其他接口**：CardTransaction、CardIssueRecord 等已直接使用 `operatorID`，无需修改

## Capabilities

### Modified Capabilities

- `recharge-c-entry`: C端充值录入接口的操作员信息获取逻辑，从硬编码 mock 改为从 JWT + DB 查询真实用户

## Impact

- **后端**：`backend/internal/handler/recharge.go` — `getOperatorCenter` → `getOperatorInfo`，`CreateCRecharge` 逻辑调整
- **前端**：无需修改
- **数据库**：依赖 `userRepo.GetByID()` 可用（已存在）
