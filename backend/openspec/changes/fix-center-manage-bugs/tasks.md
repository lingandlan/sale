## 1. 后端 — UpdateCenterRequest 扩展

- [x] 1.1 `model/recharge.go` UpdateCenterRequest 增加 Status、Province、City、District、ManagerID 字段，带 binding 验证标签
- [x] 1.2 `handler/recharge.go` UpdateCenter handler 传递所有非空字段（status/province/city/district/managerId）到 service data map
- [x] 1.3 验证：运行 `go test ./internal/...` 确保编译通过、现有测试不破

## 2. 后端 — GetCenters 管理员关联修复

- [x] 2.1 `service/recharge.go` GetCenters 方法：将 opMap 的 key 从 `operator.CenterID` 改为 `operator.ID`，通过 `managerId` 查找管理员信息
- [x] 2.2 验证：调用 `GET /center` 确认返回的 managerName/managerPhone 与 managerId 对应的操作员匹配

## 3. 前端 — 详情弹窗

- [x] 3.1 `api/center.ts` 新增 `getCenterDetail(id)` 函数
- [x] 3.2 `CenterManage.vue` 新增详情 Dialog，展示中心完整信息（名称、级别、省市区、地址、管理员、余额、累计充值、已消耗、状态）
- [x] 3.3 `handleDetail` 改为调用 `getCenterDetail(row.id)` 并打开详情弹窗

## 4. 前端 — 冻结/解冻修复

- [x] 4.1 `CenterManage.vue` handleToggleFreeze 已有正确逻辑（传 `{ status: "frozen"/"active" }`），确认后端 1.1/1.2 完成后该功能生效
- [x] 4.2 验证：点击冻结按钮，确认列表刷新后状态变为"冻结"；点击解冻按钮，确认恢复为"正常"

## 5. 集成验证

- [x] 5.1 前后端联调：确认编辑弹窗可更新省市区和管理员
- [x] 5.2 前后端联调：确认详情弹窗可正常打开并展示数据
- [x] 5.3 前后端联调：确认管理员列显示关联用户名和手机号
