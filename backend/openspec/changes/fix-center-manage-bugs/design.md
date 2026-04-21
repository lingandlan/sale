## Context

充值中心管理页面 `/center/manage` 存在 4 个功能 bug。当前代码状态：
- 后端 `GetCenters` 按 `centerID` 匹配 operator（错误），应按 `managerId` 匹配
- 后端 `UpdateCenterRequest` 只有 name/address/phone，无 status/province/city/district 字段
- 前端 `handleDetail` 只弹消息，不打开详情
- 前端城市区县用硬编码 `regionData.ts`

## Goals / Non-Goals

**Goals:**
- 修复 4 个 bug，使管理页面功能完整可用
- 冻结/解冻通过现有 UpdateCenter API 传递 status
- 详情弹窗展示中心完整信息（含操作员列表）
- 管理员列显示 managerId 对应的用户名和手机号
- 城市区县改用后端 API（如后端暂无此接口，先保留前端数据但做好抽象）

**Non-Goals:**
- 不重构页面 UI/样式
- 不新增独立的冻结 API 端点
- 不做地区数据后端接口（如不存在），保留前端 regionData

## Decisions

### D1: 管理员关联 — 修复后端匹配逻辑

**现状**: `GetCenters` 用 `opMap[c.ID]`（centerID 匹配 operator 表的 centerId）查找管理员。
**问题**: managerId 是操作员 ID，不是 centerId。需要用 `opMap[managerId]` 按 ID 查找操作员。
**方案**: 将 opMap 改为以 operator.ID 为 key，通过 managerId 查找对应操作员。

### D2: 冻结修复 — 扩展 UpdateCenterRequest

**现状**: `UpdateCenterRequest` 只有 name/address/phone，handler 只传递这 3 个字段。
**方案**: 在 `UpdateCenterRequest` 增加 `Status string` 字段（binding:"omitempty,oneof=active frozen"），handler 中传递 status 到 service 层。前端已有逻辑传 `{ status: "frozen"/"active" }`。

### D3: 详情弹窗 — 前端新增详情 Dialog

**现状**: `handleDetail` 只 `ElMessage.info`。
**方案**: 新增详情 Dialog，调用 `GET /center/:id` 获取完整信息，展示中心基本信息 + 关联操作员列表。后端 `GetCenterDetail` 已存在且可用。

### D4: 城市区县 — 保留前端数据

**现状**: 后端无地区列表 API。新增 API 工作量大且非核心 bug。
**方案**: 保留 `regionData.ts`，后续有需要再迁移到后端。

### D5: UpdateCenter 补全字段

**现状**: handler 只传递 name/address/phone，编辑时省市区等信息无法更新。
**方案**: `UpdateCenterRequest` 增加 province/city/district/managerId 字段，handler 传递所有非空字段。

## Risks / Trade-offs

- **[风险] UpdateCenterRequest 扩字段**: 现有调用方可能传多余字段 → 用 `omitempty` 标签确保向后兼容
- **[风险] opMap key 变更**: 改用 operator.ID 为 key → 需确认 managerId 存储的是操作员 ID，不是 centerId
- **[取舍] 保留前端地区数据**: 短期减少工作量，但数据更新需改前端代码 → 可接受
