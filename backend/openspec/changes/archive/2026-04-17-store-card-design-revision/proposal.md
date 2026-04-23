## Why

当前门店卡发放（CardIssue / 绑定卡号）页面存在三个核心缺陷：
1. **未验证商城用户** — `user_id` 硬编码为空字符串，没有查询手机号是否为 WSY 商城用户
2. **无充值中心权限控制** — center_admin 角色不应看到其他充值中心的卡片，但当前无任何限制
3. **operatorID 硬编码** — 所有 handler 中 `operatorID := "op123"`，未从 JWT 获取真实操作员信息

## What Changes

- **发放前查询商城用户**：页面增加手机号搜索步骤，调用已有的 `search-member` 接口验证是否为商城用户，显示用户姓名/等级后才能继续发卡
- **卡号改为可搜索下拉**：根据所选充值中心过滤可发放的卡号列表（status=已入库 且属于该中心），支持输入搜索
- **充值中心角色隔离**：center_admin 只能看到自己的充值中心（select 只有一个选项且 disabled），后端也需校验操作员是否有权操作该中心的卡
- **从 JWT 提取操作员信息**：所有 card handler 从 `c.Get("user_id")` / `c.Get("role")` 获取真实身份，替换硬编码 `"op123"`
- **事务一致性**：BindCardToUser 的卡状态更新 + issue 记录创建 + transaction 记录创建统一在一个数据库事务内

## Capabilities

### New Capabilities
- `card-issue-with-member-check`: 发卡前通过手机号查询商城会员并确认身份的完整流程

### Modified Capabilities
- `card-bind`: 卡号绑定接口增加操作员身份提取、充值中心权限校验、事务包裹
- `card-center-access`: 充值中心数据访问增加角色过滤（center_admin 只能访问所属中心）

## Impact

- **后端 handler**: `recharge.go` — BindCardToUser、AllocateCards 等所有 card handler 提取 JWT 信息
- **后端 service**: `recharge.go` — BindCardToUser 增加中心权限校验、事务包裹
- **后端 middleware**: `auth.go` — JWT claims 增加 center_id（或从 DB 查）
- **前端 CardIssue.vue**: 重构为"查用户 → 选中心 → 选卡号 → 确认发放"四步流程
- **前端 API**: `card.ts` 增加按中心搜索可用卡号接口
- **后端路由**: 新增 `GET /card/available?center_id=xxx` 查询可发放卡号列表
