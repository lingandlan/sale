## Why

充值中心的"管理员"字段当前关联 `recharge_operators` 表（独立操作员表），但主设计中操作员应来自 `users` 表。管理员下拉选的是测试操作员数据而非真实用户，需要改为从用户表读取。

## What Changes

- 前端"管理员"下拉数据源从 `GET /operator` 改为 `GET /user`（过滤合适角色）
- 后端 `GetCenters` 关联管理员信息时查 `users` 表而非 `recharge_operators` 表
- `RechargeCenter.managerId` 存储 `User.ID`（int64 的字符串形式）
- 清除 `recharge_operators` 表中仅用于管理员显示的测试数据

## Capabilities

### New Capabilities

- `center-user-manager`: 充值中心管理员关联到 users 表，前端下拉从用户列表读取

### Modified Capabilities

无

## Impact

- **后端**: `service/recharge.go` GetCenters 方法改为查 users 表；可能需要 userRepo 依赖
- **前端**: `CenterManage.vue` 管理员下拉改为调用用户列表 API；`api/user.ts` 可能需要新增列表接口
- **数据库**: `recharge_centers.manager_id` 值需迁移为 user ID；`recharge_operators` 表保留但不再用于管理员关联
