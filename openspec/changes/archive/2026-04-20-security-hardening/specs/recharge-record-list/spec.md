## ADDED Requirements

### Requirement: 充值记录路由组添加 RBAC 中间件

`/api/v1/recharge/records` 路由组 SHALL 在现有 JWT 认证中间件基础上，增加 RBAC 中间件进行角色权限校验。只有 `center_admin` 及以上角色的用户才能访问充值记录列表和详情。

#### Scenario: center_admin 访问充值记录

- **WHEN** 角色为 `center_admin` 的用户请求 `GET /api/v1/recharge/records`
- **THEN** RBAC 中间件校验通过，正常返回充值记录列表

#### Scenario: operator 访问充值记录

- **WHEN** 角色为 `operator` 的用户请求 `GET /api/v1/recharge/records`
- **THEN** 根据 Casbin 策略决定是否允许（如策略允许则通过，否则返回 403）

#### Scenario: 未认证用户访问充值记录

- **WHEN** 未携带 JWT token 的请求访问 `GET /api/v1/recharge/records`
- **THEN** 返回 401，不进入 RBAC 校验

### Requirement: 全部业务路由组添加 RBAC 中间件

以下路由组 SHALL 添加 `rbacMiddleware.Auth()` 中间件（在 `authMiddleware.Auth()` 之后）：

| 路由组 | 最低允许角色 |
|--------|-------------|
| `/dashboard` | `center_admin` |
| `/recharge/b-apply` | `operator` |
| `/recharge/b-approval` | `center_admin` |
| `/recharge/c-entry` | `operator` |
| `/recharge/records` | `operator` |
| `/card` | `operator` |
| `/center` | `center_admin`（创建/修改/删除需 `hq_admin`） |
| `/operator` | `center_admin` |

#### Scenario: 超级管理员访问所有路由

- **WHEN** 角色为 `super_admin` 的用户访问任意业务路由
- **THEN** RBAC 中间件直接放行（现有逻辑不变）

#### Scenario: 低权限角色访问高权限路由

- **WHEN** 角色为 `operator` 的用户访问 `POST /center`（创建充值中心）
- **THEN** 返回 403 Forbidden

#### Scenario: Casbin 策略未配置的路由

- **WHEN** 用户角色通过 JWT 认证但 Casbin 中无对应策略
- **THEN** RBAC 中间件返回 403，拒绝访问
