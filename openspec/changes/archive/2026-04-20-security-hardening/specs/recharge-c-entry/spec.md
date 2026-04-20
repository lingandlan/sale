## MODIFIED Requirements

### Requirement: C端充值录入写入真实操作员信息

`CreateCRecharge` 接口 SHALL 从 JWT 获取当前登录用户 ID，并通过 `userRepo.GetByID()` 查询数据库获取用户真实姓名，写入 `operatorId` 和 `operatorName` 字段，不再使用硬编码 mock 值。

操作员密码 SHALL 使用 bcrypt 哈希存储，与用户密码保持一致的安全策略。

#### Scenario: 正常流程写入真实操作员

- **WHEN** 已登录用户调用 `POST /recharge/c-entry` 创建充值
- **THEN** `operatorId` 为当前用户的数据库主键（int64），`operatorName` 为用户在 `users.name` 字段的值

#### Scenario: 用户名为空时降级为 username

- **WHEN** `users.name` 字段为空
- **THEN** `operatorName` 使用 `users.username` 字段值

#### Scenario: 用户查询失败时降级

- **WHEN** `userRepo.GetByID()` 查询失败
- **THEN** `operatorName` 降级为 `"未知用户"`，`operatorId` 仍使用 JWT 中的用户 ID

#### Scenario: 创建操作员时密码哈希存储

- **WHEN** 管理员调用 `POST /operator` 创建操作员，传入明文密码
- **THEN** 系统使用 bcrypt 对密码哈希后存入数据库，数据库中不保留明文密码

#### Scenario: 更新操作员密码时哈希存储

- **WHEN** 管理员调用 `PUT /operator/:id` 更新操作员密码
- **THEN** 系统使用 bcrypt 对新密码哈希后更新，数据库中不保留明文密码

#### Scenario: 操作员登录验证哈希密码

- **WHEN** 操作员使用明文密码登录
- **THEN** 系统使用 bcrypt.CompareHashAndPassword 验证，不依赖明文比对

#### Scenario: 现有明文密码迁移

- **WHEN** 执行数据库迁移脚本
- **THEN** 所有操作员的明文密码被转换为 bcrypt 哈希，迁移脚本支持 dry-run 模式验证
