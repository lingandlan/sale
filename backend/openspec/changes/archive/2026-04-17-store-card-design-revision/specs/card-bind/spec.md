## MODIFIED Requirements

### Requirement: 操作员身份从 JWT 提取
系统 SHALL 从 Gin Context 中提取 `user_id` 和 `role` 替换硬编码的 `operatorID := "op123"`。所有 card 相关 handler MUST 使用真实操作员 ID。

#### Scenario: handler 获取操作员信息
- **WHEN** 请求经过 auth middleware 到达 card handler
- **THEN** handler 从 `c.Get("user_id")` 和 `c.Get("role")` 获取操作员身份，传递给 service 层

### Requirement: BindCardToUser 使用数据库事务
系统 SHALL 将 BindCardToUser 中的状态更新、issue 记录创建、transaction 记录创建包裹在一个数据库事务中。任一步骤失败时 SHALL 全部回滚。

#### Scenario: 事务成功
- **WHEN** 卡状态更新、issue 记录、transaction 记录全部写入成功
- **THEN** 事务提交，卡状态变为 Issued(2)

#### Scenario: 事务中某步失败
- **WHEN** transaction 记录写入失败
- **THEN** 事务回滚，卡状态保持不变（仍为 InStock），返回错误

### Requirement: BindCardToUser 校验充值中心权限
系统 SHALL 在 BindCardToUser 中校验操作员是否有权操作该充值中心的卡。center_admin/operator 角色 MUST 只能发放所属中心的卡。

#### Scenario: center_admin 发放自己中心的卡
- **WHEN** center_admin 操作员（center_id="A"）发放 recharge_center_id="A" 的卡
- **THEN** 校验通过，正常发放

#### Scenario: center_admin 发放其他中心的卡
- **WHEN** center_admin 操作员（center_id="A"）尝试发放 recharge_center_id="B" 的卡
- **THEN** 后端返回 403 权限错误

#### Scenario: super_admin 发放任意中心的卡
- **WHEN** super_admin 操作员发放任意中心的卡
- **THEN** 不受中心限制，正常发放
