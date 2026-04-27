## Context

当前充值中心管理员（managerId）关联 `recharge_operators` 表，这是一张独立的操作员表，包含 mock 测试数据（如 op-licn、op-zhangcw）。主设计中所有人员都应在 `users` 表中管理，`users` 表已有 `name`、`phone`、`role`、`center_id` 字段。

## Goals / Non-Goals

**Goals:**
- 管理员下拉从 `users` 表读取，显示真实用户
- GetCenters 返回的 managerName/managerPhone 通过 users 表关联
- managerId 存储 user ID

**Non-Goals:**
- 不删除 `recharge_operators` 表（其他业务仍在使用）
- 不修改用户管理页面的逻辑

## Decisions

### D1: 前端管理员下拉改为用户列表

**现状**: 调用 `GET /operator`（`recharge_operators` 表）
**方案**: 改为调用 `GET /user`（已有接口，`UserService.List`），前端过滤 role 不是 super_admin 的用户

### D2: 后端 GetCenters 关联改为 users 表

**现状**: `GetCenters` 中查 `recharge_operators` 表，用 `opMap[managerId]` 匹配
**方案**: 注入 `userSvc UserServiceInterface`，通过 `userSvc.GetByID` 按 managerId 查用户。或在 repo 层批量查 users 表。

### D3: managerId 类型保持 string

User.ID 是 int64，但 managerId 在 RechargeCenter 中是 string 类型。保持 string，存 int64 的字符串形式（如 "1"、"2"）。前端不需要改类型。

## Risks / Trade-offs

- **[风险] 现有 managerId 数据失效**: 当前 managerId 存的是 operator ID（如 "op-licn"），改为 user ID 后需要清理旧数据 → 直接在数据库中清空旧 managerId，通过前端重新关联
- **[取舍] 前端加载全量用户列表**: 用户量小时可行，量大时需后端提供精简接口 → 当前阶段可接受
