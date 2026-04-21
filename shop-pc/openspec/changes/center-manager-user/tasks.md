## 1. 后端 — GetCenters 关联 users 表

- [x] 1.1 `service/recharge.go` RechargeService 注入 userRepo 依赖（如未注入），在 GetCenters 中查 users 表替代 recharge_operators 表
- [x] 1.2 清除测试数据：将现有 recharge_centers.managerId 中的 operator ID 清空
- [x] 1.3 验证：`go build ./...` 和 `go test ./internal/...` 通过

## 2. 前端 — 管理员下拉改用用户列表

- [x] 2.1 `api/user.ts` 新增 `getUserList()` 函数（如不存在），调用 `GET /user`
- [x] 2.2 `CenterManage.vue` 管理员下拉数据源从 `getOperatorList` 改为 `getUserList`，显示 `name（phone）`，value 为 user ID
- [x] 2.3 前端类型检查通过

## 3. 集成验证

- [x] 3.1 联调：编辑弹窗管理员下拉显示真实用户
- [x] 3.2 联调：列表管理员列显示关联用户名
- [x] 3.3 联调：详情弹窗管理员显示正确
