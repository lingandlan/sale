## Why

充值中心管理页面存在 4 个功能缺陷：管理员显示与用户数据未关联、地区选择使用硬编码数据、详情按钮无效、冻结操作不生效。需要修复这些 bug 以确保页面功能完整可用。

## What Changes

- **Bug1 - 管理员未关联真实用户**: 后端 GetCenters 返回的 managerId 未解析为用户姓名/手机号，前端仅显示字段值
- **Bug2 - 操作员/地区数据**: 城市区县选择使用前端硬编码 `regionData.ts`，需改为调用后端地区接口；操作员列表已使用真实 API（此项仅改地区数据源）
- **Bug3 - 详情按钮无效**: `handleDetail` 只弹出消息不跳转/弹窗，需改为调用 `GetCenterDetail` API 并展示详情弹窗
- **Bug4 - 冻结无效**: 后端 `UpdateCenterRequest` 缺少 `status` 字段，handler 也不处理 status 参数，导致前端传 `status: "frozen"` 被忽略

## Capabilities

### New Capabilities

- `center-freeze`: 充值中心冻结/解冻功能，后端 UpdateCenter 支持状态变更
- `center-detail`: 充值中心详情查看功能，前端详情弹窗 + 后端详情 API

### Modified Capabilities

无

## Impact

- **后端**: `model/recharge.go` UpdateCenterRequest 增加 status 字段；`handler/recharge.go` UpdateCenter handler 传递 status；`service/recharge.go` / repo 层支持 status 更新；新增地区 API 端点
- **前端**: `CenterManage.vue` 修改详情/冻结逻辑；新增地区 API 调用替代硬编码数据；管理员列显示关联用户信息
- **数据库**: 无 schema 变更（center 表已有 status 字段）
