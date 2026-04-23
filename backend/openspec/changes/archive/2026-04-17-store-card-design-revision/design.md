## Context

门店卡发放（CardIssue 页面，路由 `/card/issue`）当前是直接填表提交模式：手动输入卡号和手机号，调用 `POST /card/bind` 绑定。存在三个核心问题：

1. **无商城用户验证** — 后端 `user_id` 硬编码为空字符串，前端不验证手机号是否属于 WSY 商城用户
2. **无权限隔离** — center_admin 角色能看到所有充值中心的卡片，后端也不校验操作员所属中心
3. **operatorID 硬编码** — 所有 handler 中 `operatorID := "op123"`

现有架构：
- 商城用户通过 WSY 外部 API 查询，已有 `GET /recharge/c-entry/search-member?phone=xxx` 端点
- JWT claims 包含 `user_id`、`phone`、`role`，但**不含** `center_id`（需从 DB 查 `users` 表）
- StoreCard 状态流转：InStock(1) → Issued(2) → Active(3)，冻结/作废等走 `transitionCardStatus()`
- BindCardToUser 当前非事务操作（状态更新 + issue 记录 + transaction 分开写入）

## Goals / Non-Goals

**Goals:**
- 发放前验证手机号为 WSY 商城用户，展示用户信息供确认
- 卡号按充值中心过滤，支持搜索，center_admin 只能看自己中心
- 后端校验操作员身份和中心权限，替换硬编码 operatorID
- BindCardToUser 使用数据库事务保证一致性

**Non-Goals:**
- 不修改批量入库（Excel 导入）和划拨（按数量）功能
- 不修改 JWT claims 结构（不在 token 中加 center_id，改为 handler 层查 DB）
- 不修改 WSY API 集成方式
- 不做批量发放（一次只发一张卡）

## Decisions

### 1. 操作员身份获取：handler 层查 DB 而非扩展 JWT

**选择：** handler 中通过 `userRepo.GetByID(ctx, userID)` 从 DB 获取完整 User 记录（含 CenterID），而非修改 JWT claims 加入 center_id。

**理由：**
- 修改 JWT 需要重新登录所有用户，影响线上
- center_id 使用频率低（仅 card 相关 handler），不值得污染全局 claims
- DB 查询开销可接受（单行主键查询，<1ms）

### 2. 充值中心权限校验：抽取 helper 函数

**选择：** 在 handler 包中创建 `getOperatorCenter(c *gin.Context) (userID int64, role string, centerID string, err error)` helper，统一从 context 获取 user_id/role，按需查 DB 获取 center_id。

**校验逻辑：**
- `super_admin` / `hq_admin`：不限制，可操作任何中心
- `center_admin` / `operator`：必须 `card.RechargeCenterID == operator.CenterID`，否则返回 403

### 3. 可发放卡号查询接口

**选择：** 新增 `GET /api/v1/card/available?center_id=xxx&keyword=TJ000` 接口，返回指定中心下 status=1（InStock）的卡号列表。

**理由：**
- 前端 el-select remote-search 需要按中心过滤 + 关键字搜索
- center_admin 调用时后端忽略请求的 center_id，强制使用操作员自己的 center_id

### 4. 前端交互流程：四步顺序式

**选择：** 页面改为有序的四步流程：

```
步骤1: 输入手机号 → 调用 search-member → 显示用户姓名/等级
步骤2: 选择充值中心（center_admin 自动填充且禁用）
步骤3: 搜索并选择卡号（按中心过滤的 remote-select）
步骤4: 填写发放原因 → 确认提交
```

每步依赖上一步完成，使用 `el-steps` 组件或简单的条件渲染。

### 5. BindCardToUser 事务包裹

**选择：** 将现有的三步操作（状态更新、issue 记录、transaction 记录）包在一个 `db.Transaction()` 中。

**理由：** 如果 transaction 记录写入失败，卡片状态已经变了，数据不一致。事务保证要么全成功要么全回滚。

## Risks / Trade-offs

**[WSY API 不可用]** → 前端 search-member 调用失败时，显示明确错误提示"商城服务暂时不可用，请稍后重试"，不允许跳过验证

**[DB 查 center_id 增加 QPS]** → 每次发卡多一次 DB 查询。但发卡频率低（人工操作），QPS 可忽略。如后续需要优化可加 JWT center_id

**[center_admin 中心变更]** → 如果管理员被调到其他中心，需要更新 users 表的 center_id。当前无此管理功能，但不在本次范围内
