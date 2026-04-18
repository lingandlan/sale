## NEW Requirements

### Requirement: 待办事项按角色过滤返回

`GET /dashboard/todos` SHALL 根据用户角色返回不同范围的待办事项，不再返回硬编码值。

#### Scenario: super_admin/hq_admin/finance 查看全部待办

- **WHEN** super_admin、hq_admin 或 finance 角色请求 `/dashboard/todos`
- **THEN** 返回全局待办：
  - `pendingApprovals`: 状态为 `pending` 的 B端充值申请数量及描述（"N笔充值申请待审批"）
  - `expiringCards`: 7 天内将过期的门店卡数量及描述（"N张门店卡将在7天内到期"），门店卡过期条件为 `expired_at BETWEEN NOW() AND DATE_ADD(NOW(), INTERVAL 7 DAY)`

#### Scenario: center_admin/operator 查看本中心待办

- **WHEN** center_admin 或 operator 角色请求 `/dashboard/todos`
- **THEN** 返回限定到本中心的待办：
  - `pendingApprovals`: 仅本中心的待审批申请数量
  - `expiringCards`: 仅本中心（`recharge_center_id` 匹配）的即将过期门店卡数量

#### Scenario: 无待办事项

- **WHEN** 各项待办计数为 0
- **THEN** 对应项 `count` 返回 0，`description` 返回空字符串

#### Scenario: 数据库查询失败时降级

- **WHEN** 查询过程中发生数据库错误
- **THEN** 返回 HTTP 500，错误信息 "待办事项加载失败"
