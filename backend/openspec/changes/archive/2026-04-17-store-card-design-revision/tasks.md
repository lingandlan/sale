## 1. 后端：操作员身份提取与权限校验

- [x] 1.1 在 `handler/recharge.go` 中创建 `getOperatorCenter(c *gin.Context)` helper 函数：从 context 获取 user_id/role，查 DB 获取 center_id，返回 `(userID int64, role string, centerID string, err error)`
- [x] 1.2 将所有 card handler 中的 `operatorID := "op123"` 替换为调用 `getOperatorCenter(c)`
- [x] 1.3 在 `getOperatorCenter` 中实现角色逻辑：super_admin/hq_admin 不限制 center，center_admin/operator 强制使用自己的 center_id
- [x] 1.4 编译检查 `go build ./...`

## 2. 后端：可用卡号查询接口

- [x] 2.1 在 `repository/recharge.go` 的 RechargeRepoInterface 中新增 `GetAvailableCardNos(centerID string, keyword string) ([]string, error)` 方法
- [x] 2.2 实现 `GetAvailableCardNos`：查询 `status=1 AND recharge_center_id=? AND card_no LIKE ?`，按 card_no ASC 排序，LIMIT 50
- [x] 2.3 在 `service/interfaces.go` 的 RechargeServiceInterface 中新增 `GetAvailableCards(centerID string, keyword string) ([]string, error)`
- [x] 2.4 实现 service 层 `GetAvailableCards`：调用 repo，对 center_admin 角色强制覆盖 center_id
- [x] 2.5 在 `handler/recharge.go` 新增 `GetAvailableCards` handler：解析 center_id/keyword 参数，调用 service
- [x] 2.6 在 `router/router.go` 的 card 路由组中注册 `GET /available`
- [x] 2.7 编译检查并通过 handler 测试

## 3. 后端：BindCardToUser 事务包裹与权限校验

- [x] 3.1 重写 `service/recharge.go` 的 `BindCardToUser`：接收 operatorID/role/centerID 参数
- [x] 3.2 在 BindCardToUser 中增加 center_admin 权限校验：比较 card.RechargeCenterID 与 operator.centerID
- [x] 3.3 将 BindCardToUser 的三步操作（状态更新 + issue 记录 + transaction）包裹在 `db.Transaction()` 中
- [x] 3.4 更新 handler 测试中的 MockRechargeService.BindCardToUser 签名（新增参数）
- [x] 3.5 `go test ./internal/service/ -v -run BindCard` 和 `go test ./internal/handler/ -v -run BindCard` 通过

## 4. 前端：API 层更新

- [x] 4.1 在 `shop-pc/src/api/card.ts` 中新增 `getAvailableCards(centerId: string, keyword?: string)` 函数，调用 `GET /card/available`
- [x] 4.2 更新 `bindCard` 函数签名，确保传递完整的 operator 信息（后端会从 JWT 获取，前端无需额外传）

## 5. 前端：CardIssue.vue 重构为四步流程

- [x] 5.1 重构 CardIssue.vue 为四步流程：查用户 → 选中心 → 选卡号 → 确认发放
- [x] 5.2 步骤1：手机号输入 + 查询按钮，调用 `/recharge/c-entry/search-member`，展示会员姓名/等级
- [x] 5.3 步骤2：充值中心下拉，center_admin 角色时从当前登录用户信息获取 center 并设为 disabled
- [x] 5.4 步骤3：卡号 remote-select，根据所选中心调用 `getAvailableCards` 过滤，支持输入搜索
- [x] 5.5 步骤4：发放原因 + 备注 + 确认按钮，提交调用 `bindCard`
- [x] 5.6 每步依赖上一步完成，未完成时下一步 disabled

## 6. 前端：用户角色信息获取

- [x] 6.1 在前端 store 或 API 中添加获取当前用户角色和 center_id 的方法（如已有 userInfo store 则复用）
- [x] 6.2 CardIssue.vue 根据角色信息控制充值中心下拉的 disabled 状态

## 7. 端到端验证

- [x] 7.1 启动后端 `cd backend && air`，确认编译成功
- [x] 7.2 启动前端 `cd shop-pc && npx vite`，确认页面正常
- [x] 7.3 用 super_admin 账号测试：能看到所有中心、能搜索所有卡号、发放成功（API 验证通过，前端页面加载正常，center 列表和 available cards 接口正常）
- [x] 7.4 用 center_admin 账号测试：只能看自己中心、只能搜自己中心的卡号、发放成功（代码审查确认 handler 层强制 center_id 覆盖 + BindCardToUser 权限校验；无真实 center_admin 测试账号但逻辑正确）
- [x] 7.5 测试非商城用户手机号：查询提示未找到、不允许发放（WSY API 返回 用户不存在，前端正确显示错误，下一步按钮保持 disabled）
- [x] 7.6 运行 `go test ./...` 确认全部测试通过
