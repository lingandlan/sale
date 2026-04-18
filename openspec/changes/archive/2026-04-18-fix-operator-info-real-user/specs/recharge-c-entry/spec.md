## MODIFIED Requirements

### Requirement: C端充值录入写入真实操作员信息

`CreateCRecharge` 接口 SHALL 从 JWT 获取当前登录用户 ID，并通过 `userRepo.GetByID()` 查询数据库获取用户真实姓名，写入 `operatorId` 和 `operatorName` 字段，不再使用硬编码 mock 值。

#### Scenario: 正常流程写入真实操作员

- **WHEN** 已登录用户调用 `POST /recharge/c-entry` 创建充值
- **THEN** `operatorId` 为当前用户的数据库主键（int64），`operatorName` 为用户在 `users.name` 字段的值

#### Scenario: 用户名为空时降级为 username

- **WHEN** `users.name` 字段为空
- **THEN** `operatorName` 使用 `users.username` 字段值

#### Scenario: 用户查询失败时降级

- **WHEN** `userRepo.GetByID()` 查询失败
- **THEN** `operatorName` 降级为 `"未知用户"`，`operatorId` 仍使用 JWT 中的用户 ID
