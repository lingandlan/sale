## ADDED Requirements

### Requirement: 可发放卡号按充值中心过滤
系统 SHALL 提供 `GET /api/v1/card/available` 接口，按充值中心过滤 status=InStock(1) 的卡号列表，支持关键字搜索。

#### Scenario: 按中心查询可用卡号
- **WHEN** 请求 `GET /card/available?center_id=xxx`
- **THEN** 返回该中心下 status=1 的卡号列表（最多返回 50 条）

#### Scenario: 按中心 + 关键字搜索
- **WHEN** 请求 `GET /card/available?center_id=xxx&keyword=TJ0001`
- **THEN** 返回该中心下卡号包含 "TJ0001" 且 status=1 的卡号列表

#### Scenario: center_admin 查询强制过滤自己中心
- **WHEN** center_admin（center_id="A"）请求 `GET /card/available?center_id=B`
- **THEN** 后端忽略请求的 center_id=B，强制返回 center_id=A 的卡号列表

### Requirement: 前端充值中心下拉按角色过滤
系统 SHALL 根据操作员角色过滤充值中心下拉列表。center_admin/operator 只能看到自己所属的中心（下拉框只有一个选项且禁用）。

#### Scenario: super_admin 选择充值中心
- **WHEN** super_admin 打开发卡页面
- **THEN** 充值中心下拉显示所有中心，可自由选择

#### Scenario: center_admin 选择充值中心
- **WHEN** center_admin 打开发卡页面
- **THEN** 充值中心下拉只有一个选项（所属中心），下拉框 disabled

### Requirement: 后端校验充值中心数据访问权限
系统 SHALL 在所有返回中心级别数据的接口中校验操作员角色权限。center_admin/operator 角色 MUST 只能访问所属中心的数据。

#### Scenario: center_admin 查询其他中心数据
- **WHEN** center_admin 请求非所属中心的数据
- **THEN** 后端返回 403 权限错误

#### Scenario: super_admin 查询任意中心数据
- **WHEN** super_admin 请求任意中心的数据
- **THEN** 正常返回数据
