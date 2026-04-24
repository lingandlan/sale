## Why

WSY商城对接目前已实现 Token认证、手机号查询用户、查询积分余额（`GetAccessToken`/`PhoneToUserID`/`GetUserIntegral`）。但充值核心流程需要的积分操作（添加/扣除）以及零钱（余额）查询/变动接口尚未对接，导致充值后无法实际变更商城用户资产。

## What Changes

- **新增积分添加接口** — 调用 `10000_integral_add`，充值场景使用 `change_type=recharge`，支持幂等防重复（union_id）
- **新增积分扣除接口** — 调用 `10000_integral_reduce`，用于退款/冲正等场景
- **新增零钱查询接口** — 查询用户零钱余额（act 值待确认）
- **新增零钱增加接口** — 增加用户零钱（act 值待确认）
- **新增零钱扣除接口** — 扣除用户零钱（act 值待确认）
- **扩展 WSYClient 接口** — 在 `pkg/mall/client.go` 新增上述方法
- **扩展 MemberService** — 暴露积分/零钱操作给业务层

## Capabilities

### New Capabilities

- `wsy-integral-ops`: WSY商城积分操作 — 添加积分、扣除积分，包含幂等控制、batchcode生成规则
- `wsy-wallet-ops`: WSY商城零钱操作 — 查询余额、增加零钱、扣除零钱（act 值确认后实现）

### Modified Capabilities

（无现有 spec 需要修改）

## Impact

- **代码**: `backend/pkg/mall/client.go` 新增方法；`backend/internal/service/member.go` 新增业务方法
- **配置**: 无新增配置（复用现有 MallConfig）
- **API**: 无新增对外的 HTTP 接口（WSY 为内部调用）
- **依赖**: 依赖 WSY 商城零钱接口文档确认 act 值
