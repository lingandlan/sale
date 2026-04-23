## ADDED Requirements

### Requirement: 前端管理员下拉使用用户列表
CenterManage 页面编辑弹窗的"管理员"下拉 SHALL 从 `GET /user` 接口获取数据，替代 `GET /operator`。下拉项显示用户姓名和手机号，值为用户 ID。

#### Scenario: 加载管理员下拉
- **WHEN** 页面挂载时加载用户列表
- **THEN** 管理员下拉显示 `users` 表中的用户（排除 super_admin），格式为 `姓名（手机号）`

#### Scenario: 选择管理员并保存
- **WHEN** 用户选择某用户作为管理员并保存
- **THEN** 前端将 user ID（int64 字符串）作为 managerId 传给后端

### Requirement: 后端 GetCenters 从 users 表关联管理员
`GetCenters` 返回的 `managerName` 和 `managerPhone` SHALL 通过 `users` 表查询，使用 `managerId` 匹配 `users.id`。

#### Scenario: 中心有管理员
- **WHEN** 充值中心的 `managerId` 能匹配到 `users` 表中的用户
- **THEN** 返回 `managerName`（用户姓名）和 `managerPhone`（用户手机号）

#### Scenario: 中心无管理员
- **WHEN** 充值中心的 `managerId` 为空或不匹配任何用户
- **THEN** 不返回 managerName 和 managerPhone 字段

### Requirement: 详情弹窗显示用户管理员
详情弹窗中管理员信息 SHALL 显示关联用户的姓名和手机号。

#### Scenario: 详情中有管理员
- **WHEN** 充值中心的 managerId 关联到用户
- **THEN** 详情弹窗管理员字段显示 `姓名（手机号）`
